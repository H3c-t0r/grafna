package state

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/grafana/grafana/pkg/infra/log"
	gModels "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/annotations"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/ngalert/eval"
	"github.com/grafana/grafana/pkg/services/ngalert/image"
	"github.com/grafana/grafana/pkg/services/ngalert/metrics"
	"github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/grafana/grafana/pkg/services/ngalert/store"
	"github.com/grafana/grafana/pkg/services/screenshot"
)

var ResendDelay = 30 * time.Second

// AlertInstanceManager defines the interface for querying the current alert instances.
type AlertInstanceManager interface {
	GetAll(orgID int64) []*State
	GetStatesForRuleUID(orgID int64, alertRuleUID string) []*State
}

type Manager struct {
	log     log.Logger
	metrics *metrics.State

	clock       clock.Clock
	cache       *cache
	quit        chan struct{}
	ResendDelay time.Duration

	ruleStore        store.RuleStore
	instanceStore    store.InstanceStore
	dashboardService dashboards.DashboardService
	imageService     image.ImageService
}

func NewManager(logger log.Logger, metrics *metrics.State, externalURL *url.URL,
	ruleStore store.RuleStore, instanceStore store.InstanceStore,
	dashboardService dashboards.DashboardService, imageService image.ImageService, clock clock.Clock) *Manager {
	manager := &Manager{
		cache:            newCache(logger, metrics, externalURL),
		quit:             make(chan struct{}),
		ResendDelay:      ResendDelay, // TODO: make this configurable
		log:              logger,
		metrics:          metrics,
		ruleStore:        ruleStore,
		instanceStore:    instanceStore,
		dashboardService: dashboardService,
		imageService:     imageService,
		clock:            clock,
	}
	go manager.recordMetrics()
	return manager
}

func (st *Manager) Close(ctx context.Context) {
	st.quit <- struct{}{}
	st.flushState(ctx)
}

func (st *Manager) Warm(ctx context.Context) {
	st.log.Info("warming cache for startup")
	st.ResetCache()

	orgIds, err := st.instanceStore.FetchOrgIds(ctx)
	if err != nil {
		st.log.Error("unable to fetch orgIds", "msg", err.Error())
	}

	var states []*State
	for _, orgId := range orgIds {
		// Get Rules
		ruleCmd := models.ListAlertRulesQuery{
			OrgID: orgId,
		}
		if err := st.ruleStore.ListAlertRules(ctx, &ruleCmd); err != nil {
			st.log.Error("unable to fetch previous state", "msg", err.Error())
		}

		ruleByUID := make(map[string]*models.AlertRule, len(ruleCmd.Result))
		for _, rule := range ruleCmd.Result {
			ruleByUID[rule.UID] = rule
		}

		// Get Instances
		cmd := models.ListAlertInstancesQuery{
			RuleOrgID: orgId,
		}
		if err := st.instanceStore.ListAlertInstances(ctx, &cmd); err != nil {
			st.log.Error("unable to fetch previous state", "msg", err.Error())
		}

		for _, entry := range cmd.Result {
			ruleForEntry, ok := ruleByUID[entry.RuleUID]
			if !ok {
				st.log.Error("rule not found for instance, ignoring", "rule", entry.RuleUID)
				continue
			}

			lbs := map[string]string(entry.Labels)
			cacheId, err := entry.Labels.StringKey()
			if err != nil {
				st.log.Error("error getting cacheId for entry", "msg", err.Error())
			}
			stateForEntry := &State{
				AlertRuleUID:         entry.RuleUID,
				OrgID:                entry.RuleOrgID,
				CacheId:              cacheId,
				Labels:               lbs,
				State:                translateInstanceState(entry.CurrentState),
				StateReason:          entry.CurrentReason,
				LastEvaluationString: "",
				StartsAt:             entry.CurrentStateSince,
				EndsAt:               entry.CurrentStateEnd,
				LastEvaluationTime:   entry.LastEvalTime,
				Annotations:          ruleForEntry.Annotations,
			}
			states = append(states, stateForEntry)
		}
	}

	for _, s := range states {
		st.set(s)
	}
}

func (st *Manager) getOrCreate(ctx context.Context, alertRule *models.AlertRule, result eval.Result, extraLabels data.Labels) *State {
	return st.cache.getOrCreate(ctx, alertRule, result, extraLabels)
}

func (st *Manager) set(entry *State) {
	st.cache.set(entry)
}

func (st *Manager) Get(orgID int64, alertRuleUID, stateId string) (*State, error) {
	return st.cache.get(orgID, alertRuleUID, stateId)
}

// ResetCache is used to ensure a clean cache on startup.
func (st *Manager) ResetCache() {
	st.cache.reset()
}

// RemoveByRuleUID deletes all entries in the state manager that match the given rule UID.
func (st *Manager) RemoveByRuleUID(orgID int64, ruleUID string) {
	st.cache.removeByRuleUID(orgID, ruleUID)
}

// ProcessEvalResults updates the current states that belong to a rule with the evaluation results.
// if extraLabels is not empty, those labels will be added to every state. The extraLabels take precedence over rule labels and result labels
func (st *Manager) ProcessEvalResults(ctx context.Context, evaluatedAt time.Time, alertRule *models.AlertRule, results eval.Results, extraLabels data.Labels) []*State {
	logger := st.log.New(alertRule.GetKey().LogContext())
	logger.Debug("state manager processing evaluation results", "resultCount", len(results))
	var states []*State
	processedResults := make(map[string]*State, len(results))
	for _, result := range results {
		s := st.setNextState(ctx, alertRule, result, extraLabels)
		states = append(states, s)
		processedResults[s.CacheId] = s
	}

	st.cleanupStaleResults(ctx, evaluatedAt, alertRule, processedResults)
	if len(states) > 0 {
		logger.Debug("saving new states to the database", "count", len(states))
		_, _ = st.saveAlertStates(ctx, states...)
	}

	return states
}

// Maybe take a screenshot. Do it if:
// 1. The alert state is transitioning into the "Alerting" state from something else.
// 2. The alert state has just transitioned to the resolved state.
// 3. The state is alerting and there is no screenshot annotation on the alert state.
func (st *Manager) maybeTakeScreenshot(
	ctx context.Context,
	alertRule *models.AlertRule,
	state *State,
	oldState eval.State,
) error {
	shouldScreenshot := state.Resolved ||
		state.State == eval.Alerting && oldState != eval.Alerting ||
		state.State == eval.Alerting && state.Image == nil
	if !shouldScreenshot {
		return nil
	}

	img, err := st.imageService.NewImage(ctx, alertRule)
	if err != nil &&
		errors.Is(err, screenshot.ErrScreenshotsUnavailable) ||
		errors.Is(err, image.ErrNoDashboard) ||
		errors.Is(err, image.ErrNoPanel) {
		// It's not an error if screenshots are disabled, or our rule isn't allowed to generate screenshots.
		return nil
	} else if err != nil {
		return err
	}
	state.Image = img
	return nil
}

// Set the current state based on evaluation results
func (st *Manager) setNextState(ctx context.Context, alertRule *models.AlertRule, result eval.Result, extraLabels data.Labels) *State {
	currentState := st.getOrCreate(ctx, alertRule, result, extraLabels)

	currentState.LastEvaluationTime = result.EvaluatedAt
	currentState.EvaluationDuration = result.EvaluationDuration
	currentState.Results = append(currentState.Results, Evaluation{
		EvaluationTime:  result.EvaluatedAt,
		EvaluationState: result.State,
		Values:          NewEvaluationValues(result.Values),
		Condition:       alertRule.Condition,
	})
	currentState.LastEvaluationString = result.EvaluationString
	currentState.TrimResults(alertRule)
	oldState := currentState.State
	oldReason := currentState.StateReason

	st.log.Debug("setting alert state", "uid", alertRule.UID)
	switch result.State {
	case eval.Normal:
		currentState.resultNormal(alertRule, result)
	case eval.Alerting:
		currentState.resultAlerting(alertRule, result)
	case eval.Error:
		currentState.resultError(alertRule, result)
	case eval.NoData:
		currentState.resultNoData(alertRule, result)
	case eval.Pending: // we do not emit results with this state
	}

	// Set reason iff: result is different than state, reason is not Alerting or Normal
	currentState.StateReason = ""

	if currentState.State != result.State &&
		result.State != eval.Normal &&
		result.State != eval.Alerting {
		currentState.StateReason = result.State.String()
	}

	// Set Resolved property so the scheduler knows to send a postable alert
	// to Alertmanager.
	currentState.Resolved = oldState == eval.Alerting && currentState.State == eval.Normal

	err := st.maybeTakeScreenshot(ctx, alertRule, currentState, oldState)
	if err != nil {
		st.log.Warn("failed to generate a screenshot for an alert instance",
			"alert_rule", alertRule.UID,
			"dashboard", alertRule.DashboardUID,
			"panel", alertRule.PanelID,
			"err", err)
	}

	st.set(currentState)

	shouldUpdateAnnotation := oldState != currentState.State || oldReason != currentState.StateReason
	if shouldUpdateAnnotation {
		go st.annotateState(ctx, alertRule, currentState.Labels, result.EvaluatedAt, InstanceStateAndReason{State: currentState.State, Reason: currentState.StateReason}, InstanceStateAndReason{State: oldState, Reason: oldReason})
	}
	return currentState
}

func (st *Manager) GetAll(orgID int64) []*State {
	return st.cache.getAll(orgID)
}

func (st *Manager) GetStatesForRuleUID(orgID int64, alertRuleUID string) []*State {
	return st.cache.getStatesForRuleUID(orgID, alertRuleUID)
}

func (st *Manager) recordMetrics() {
	// TODO: parameterize?
	// Setting to a reasonable default scrape interval for Prometheus.
	dur := time.Duration(15) * time.Second
	ticker := st.clock.Ticker(dur)
	for {
		select {
		case <-ticker.C:
			st.log.Debug("recording state cache metrics", "now", st.clock.Now())
			st.cache.recordMetrics()
		case <-st.quit:
			st.log.Debug("stopping state cache metrics recording", "now", st.clock.Now())
			ticker.Stop()
			return
		}
	}
}

func (st *Manager) Put(states []*State) {
	for _, s := range states {
		st.set(s)
	}
}

// flushState dumps the entire state to the database
func (st *Manager) flushState(ctx context.Context) {
	t := st.clock.Now()
	st.log.Info("flushing the state")
	st.cache.mtxStates.Lock()
	defer st.cache.mtxStates.Unlock()
	totalStates, errorsCnt := 0, 0
	var stateSlice []*State
	for _, orgStates := range st.cache.states {
		for _, ruleStates := range orgStates {
			stateSlice = stateSlice[:0]
			for _, state := range ruleStates {
				stateSlice = append(stateSlice, state)
			}
			savedCount, failedCount := st.saveAlertStates(ctx, stateSlice...)
			if failedCount != 0 {
				st.log.Error("failed to save alert states.", "count", failedCount)
			}
			errorsCnt += failedCount
			totalStates += savedCount
		}
	}

	st.log.Info("the state has been flushed", "total_instances", totalStates, "errors", errorsCnt, "took", st.clock.Since(t))
}

// TODO: Is the `State` type necessary? Should it embed the instance?
func (st *Manager) saveAlertStates(ctx context.Context, states ...*State) (saved, failed int) {
	st.log.Debug("saving alert states", "count", len(states))
	instances := make([]models.AlertInstance, 0, len(states))

	type debugInfo struct {
		OrgID  int64
		Uid    string
		State  string
		Labels string
	}
	debug := make([]debugInfo, 0)

	for _, s := range states {
		labels := models.InstanceLabels(s.Labels)
		_, hash, err := labels.StringAndHash()
		if err != nil {
			debug = append(debug, debugInfo{s.OrgID, s.AlertRuleUID, s.State.String(), s.Labels.String()})
			st.log.Error("failed to save alert instance with invalid labels", "orgID", s.OrgID, "ruleUID", s.AlertRuleUID, "err", err)
			continue
		}
		fields := models.AlertInstance{
			AlertInstanceKey: models.AlertInstanceKey{
				RuleOrgID:  s.OrgID,
				RuleUID:    s.AlertRuleUID,
				LabelsHash: hash,
			},
			Labels:            models.InstanceLabels(s.Labels),
			CurrentState:      models.InstanceStateType(s.State.String()),
			CurrentReason:     s.StateReason,
			LastEvalTime:      s.LastEvaluationTime,
			CurrentStateSince: s.StartsAt,
			CurrentStateEnd:   s.EndsAt,
		}
		instances = append(instances, fields)
	}
	err := st.instanceStore.SaveAlertInstances(ctx, instances...)
	if err != nil {
		for _, inst := range instances {
			debug = append(debug, debugInfo{inst.RuleOrgID, inst.RuleUID, string(inst.CurrentState), data.Labels(inst.Labels).String()})
		}
		st.log.Error("failed to save alert states", "states", debug, "err", err)
		return 0, len(debug)
	}

	return len(instances), len(debug)
}

// TODO: why wouldn't you allow other types like NoData or Error?
func translateInstanceState(state models.InstanceStateType) eval.State {
	switch {
	case state == models.InstanceStateFiring:
		return eval.Alerting
	case state == models.InstanceStateNormal:
		return eval.Normal
	default:
		return eval.Error
	}
}

// This struct provides grouping of state with reason, and string formatting.
type InstanceStateAndReason struct {
	State  eval.State
	Reason string
}

func (i InstanceStateAndReason) String() string {
	s := fmt.Sprintf("%v", i.State)
	if len(i.Reason) > 0 {
		s += fmt.Sprintf(" (%v)", i.Reason)
	}
	return s
}

func (st *Manager) annotateState(ctx context.Context, alertRule *models.AlertRule, labels data.Labels, evaluatedAt time.Time, currentData, previousData InstanceStateAndReason) {
	st.log.Debug("alert state changed creating annotation", "alertRuleUID", alertRule.UID, "newState", currentData.String(), "oldState", previousData.String())

	labels = removePrivateLabels(labels)
	annotationText := fmt.Sprintf("%s {%s} - %s", alertRule.Title, labels.String(), currentData.String())

	item := &annotations.Item{
		AlertId:   alertRule.ID,
		OrgId:     alertRule.OrgID,
		PrevState: previousData.String(),
		NewState:  currentData.String(),
		Text:      annotationText,
		Epoch:     evaluatedAt.UnixNano() / int64(time.Millisecond),
	}

	dashUid, ok := alertRule.Annotations[models.DashboardUIDAnnotation]
	if ok {
		panelUid := alertRule.Annotations[models.PanelIDAnnotation]

		panelId, err := strconv.ParseInt(panelUid, 10, 64)
		if err != nil {
			st.log.Error("error parsing panelUID for alert annotation", "panelUID", panelUid, "alertRuleUID", alertRule.UID, "err", err.Error())
			return
		}

		query := &gModels.GetDashboardQuery{
			Uid:   dashUid,
			OrgId: alertRule.OrgID,
		}

		err = st.dashboardService.GetDashboard(ctx, query)
		if err != nil {
			st.log.Error("error getting dashboard for alert annotation", "dashboardUID", dashUid, "alertRuleUID", alertRule.UID, "err", err.Error())
			return
		}

		item.PanelId = panelId
		item.DashboardId = query.Result.Id
	}

	annotationRepo := annotations.GetRepository()
	if err := annotationRepo.Save(item); err != nil {
		st.log.Error("error saving alert annotation", "alertRuleUID", alertRule.UID, "err", err.Error())
		return
	}
}

func (st *Manager) cleanupStaleResults(
	ctx context.Context,
	evaluatedAt time.Time,
	alertRule *models.AlertRule,
	newStates map[string]*State,
) {
	allStates := st.GetStatesForRuleUID(alertRule.OrgID, alertRule.UID)
	toDelete := make([]models.AlertInstanceKey, 0)

	for _, s := range allStates {
		// Is the cached state in our recently processed results? If not, is it stale?
		if _, ok := newStates[s.CacheId]; !ok && stateIsStale(evaluatedAt, s.LastEvaluationTime, alertRule.IntervalSeconds) {
			st.log.Debug("removing stale state entry", "orgID", s.OrgID, "alertRuleUID", s.AlertRuleUID, "cacheID", s.CacheId)
			st.cache.deleteEntry(s.OrgID, s.AlertRuleUID, s.CacheId)
			ilbs := models.InstanceLabels(s.Labels)
			_, labelsHash, err := ilbs.StringAndHash()
			if err != nil {
				st.log.Error("unable to get labelsHash", "err", err.Error(), "orgID", s.OrgID, "alertRuleUID", s.AlertRuleUID)
			}

			toDelete = append(toDelete, models.AlertInstanceKey{RuleOrgID: s.OrgID, RuleUID: s.AlertRuleUID, LabelsHash: labelsHash})

			if s.State == eval.Alerting {
				st.annotateState(ctx, alertRule, s.Labels, evaluatedAt,
					InstanceStateAndReason{State: eval.Normal, Reason: ""},
					InstanceStateAndReason{State: s.State, Reason: s.StateReason})
			}
		}
	}

	if err := st.instanceStore.DeleteAlertInstances(ctx, toDelete...); err != nil {
		st.log.Error("unable to delete stale instances from database", "err", err.Error(),
			"orgID", alertRule.OrgID, "alertRuleUID", alertRule.UID, "count", len(toDelete))
	}
}

func stateIsStale(evaluatedAt time.Time, lastEval time.Time, intervalSeconds int64) bool {
	return !lastEval.Add(2 * time.Duration(intervalSeconds) * time.Second).After(evaluatedAt)
}

func removePrivateLabels(labels data.Labels) data.Labels {
	result := make(data.Labels)
	for k, v := range labels {
		if !strings.HasPrefix(k, "__") && !strings.HasSuffix(k, "__") {
			result[k] = v
		}
	}
	return result
}
