package state

import (
	"fmt"
	"sync"
	"time"

	prometheusModel "github.com/prometheus/common/model"

	"github.com/grafana/grafana/pkg/infra/log"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana/pkg/services/ngalert/eval"
	ngModels "github.com/grafana/grafana/pkg/services/ngalert/models"
)

type AlertState struct {
	AlertRuleUID       string
	OrgID              int64
	CacheId            string
	Labels             data.Labels
	State              eval.State
	Results            []StateEvaluation
	StartsAt           time.Time
	EndsAt             time.Time
	LastEvaluationTime time.Time
	EvaluationDuration time.Duration
	Annotations        map[string]string
}

type StateEvaluation struct {
	EvaluationTime  time.Time
	EvaluationState eval.State
}

type cache struct {
	states    map[string]AlertState
	mtxStates sync.Mutex
}

type StateTracker struct {
	cache cache
	quit  chan struct{}
	Log   log.Logger
}

func NewStateTracker(logger log.Logger) *StateTracker {
	tracker := &StateTracker{
		cache: cache{
			states: make(map[string]AlertState),
		},
		quit: make(chan struct{}),
		Log:  logger,
	}
	go tracker.cleanUp()
	return tracker
}

func (st *StateTracker) getOrCreate(alertRule *ngModels.AlertRule, result eval.Result, evaluationDuration time.Duration) AlertState {
	st.cache.mtxStates.Lock()
	defer st.cache.mtxStates.Unlock()

	// if duplicate labels exist, alertRule label will take precedence
	lbs := mergeLabels(alertRule.Labels, result.Instance)
	lbs["__alert_rule_uid__"] = alertRule.UID
	lbs["__alert_rule_namespace_uid__"] = alertRule.NamespaceUID
	lbs[prometheusModel.AlertNameLabel] = alertRule.Title

	id := fmt.Sprintf("%s", map[string]string(lbs))
	if state, ok := st.cache.states[id]; ok {
		return state
	}

	annotations := map[string]string{}
	if len(alertRule.Annotations) > 0 {
		annotations = alertRule.Annotations
	}

	// If the first result we get is alerting, set StartsAt to EvaluatedAt because we
	// do not have data for determining StartsAt otherwise
	st.Log.Debug("adding new alert state cache entry", "cacheId", id, "state", result.State.String(), "evaluatedAt", result.EvaluatedAt.String())
	newState := AlertState{
		AlertRuleUID:       alertRule.UID,
		OrgID:              alertRule.OrgID,
		CacheId:            id,
		Labels:             lbs,
		State:              result.State,
		Annotations:        annotations,
		EvaluationDuration: evaluationDuration,
	}
	if result.State == eval.Alerting {
		newState.StartsAt = result.EvaluatedAt
	}
	st.cache.states[id] = newState
	return newState
}

func (st *StateTracker) set(entry AlertState) {
	st.cache.mtxStates.Lock()
	defer st.cache.mtxStates.Unlock()
	st.cache.states[entry.CacheId] = entry
}

func (st *StateTracker) Get(id string) AlertState {
	st.cache.mtxStates.Lock()
	defer st.cache.mtxStates.Unlock()
	return st.cache.states[id]
}

//Used to ensure a clean cache on startup
func (st *StateTracker) ResetCache() {
	st.cache.mtxStates.Lock()
	defer st.cache.mtxStates.Unlock()
	st.cache.states = make(map[string]AlertState)
}

func (st *StateTracker) ProcessEvalResults(alertRule *ngModels.AlertRule, results eval.Results, evaluationDuration time.Duration) []AlertState {
	st.Log.Info("state tracker processing evaluation results", "uid", alertRule.UID, "resultCount", len(results))
	var states []AlertState
	for _, result := range results {
		s := st.setNextState(alertRule, result, evaluationDuration)
		states = append(states, s)
	}
	st.Log.Debug("returning changed states to scheduler", "count", len(states))
	return states
}

//TODO: When calculating if an alert should not be firing anymore, we should take three things into account:
// 1. The re-send the delay if any, we don't want to send every firing alert every time, we should have a fixed delay across all alerts to avoid saturating the notification system
//Set the current state based on evaluation results
func (st *StateTracker) setNextState(alertRule *ngModels.AlertRule, result eval.Result, evaluationDuration time.Duration) AlertState {
	currentState := st.getOrCreate(alertRule, result, evaluationDuration)

	currentState.LastEvaluationTime = result.EvaluatedAt
	currentState.EvaluationDuration = evaluationDuration
	currentState.Results = append(currentState.Results, StateEvaluation{
		EvaluationTime:  result.EvaluatedAt,
		EvaluationState: result.State,
	})

	st.Log.Debug("setting alert state", "uid", alertRule.UID)
	switch result.State {
	case eval.Normal:
		currentState = resultNormal(currentState, result)
	case eval.Alerting:
		currentState = currentState.resultAlerting(alertRule, result)
	case eval.Error:
		currentState = currentState.resultError(alertRule, result)
	case eval.NoData:
		currentState = currentState.resultNoData(alertRule, result)
	case eval.Pending: // we do not emit results with this state
	}

	st.set(currentState)
	return currentState
}

func resultNormal(alertState AlertState, result eval.Result) AlertState {
	newState := alertState
	if alertState.State != eval.Normal {
		newState.EndsAt = result.EvaluatedAt
	}
	newState.State = eval.Normal
	return newState
}

func (a AlertState) resultAlerting(alertRule *ngModels.AlertRule, result eval.Result) AlertState {
	switch a.State {
	case eval.Alerting:
		if !(alertRule.For > 0) {
			// If there is not For set, we will set EndsAt to be twice the evaluation interval
			// to avoid flapping with every evaluation
			a.EndsAt = result.EvaluatedAt.Add(time.Duration(alertRule.IntervalSeconds*2) * time.Second)
			return a
		}
		a.EndsAt = result.EvaluatedAt.Add(alertRule.For)
	case eval.Pending:
		if result.EvaluatedAt.Sub(a.StartsAt) > alertRule.For {
			a.State = eval.Alerting
			a.StartsAt = result.EvaluatedAt
			a.EndsAt = result.EvaluatedAt.Add(alertRule.For)
			a.Annotations["alerting_at"] = result.EvaluatedAt.String()
		}
	default:
		a.StartsAt = result.EvaluatedAt
		if !(alertRule.For > 0) {
			a.EndsAt = result.EvaluatedAt.Add(time.Duration(alertRule.IntervalSeconds*2) * time.Second)
			a.State = eval.Alerting
			a.Annotations["alerting_at"] = result.EvaluatedAt.String()
		} else {
			a.EndsAt = result.EvaluatedAt.Add(alertRule.For)
			if result.EvaluatedAt.Sub(a.StartsAt) > alertRule.For {
				a.State = eval.Alerting
				a.Annotations["alerting_at"] = result.EvaluatedAt.String()
			} else {
				a.State = eval.Pending
			}
		}
	}
	return a
}

func (a AlertState) resultError(alertRule *ngModels.AlertRule, result eval.Result) AlertState {
	if a.StartsAt.IsZero() {
		a.StartsAt = result.EvaluatedAt
	}
	if !(alertRule.For > 0) {
		a.EndsAt = result.EvaluatedAt.Add(time.Duration(alertRule.IntervalSeconds*2) * time.Second)
	} else {
		a.EndsAt = result.EvaluatedAt.Add(alertRule.For)
	}
	if a.State != eval.Error {
		a.Annotations["last_error"] = result.EvaluatedAt.String()
	}

	switch alertRule.ExecErrState {
	case ngModels.AlertingErrState:
		a.State = eval.Alerting
	case ngModels.KeepLastStateErrState:
	}
	return a
}

func (a AlertState) resultNoData(alertRule *ngModels.AlertRule, result eval.Result) AlertState {
	if a.StartsAt.IsZero() {
		a.StartsAt = result.EvaluatedAt
	}
	if !(alertRule.For > 0) {
		a.EndsAt = result.EvaluatedAt.Add(time.Duration(alertRule.IntervalSeconds*2) * time.Second)
	} else {
		a.EndsAt = result.EvaluatedAt.Add(alertRule.For)
	}
	if a.State != eval.NoData {
		a.Annotations["no_data"] = result.EvaluatedAt.String()
	}

	switch alertRule.NoDataState {
	case ngModels.Alerting:
		a.State = eval.Alerting
	case ngModels.NoData:
		a.State = eval.NoData
	case ngModels.KeepLastState:
	case ngModels.OK:
		a.State = eval.Normal
	}
	return a
}

func (st *StateTracker) GetAll() []AlertState {
	var states []AlertState
	st.cache.mtxStates.Lock()
	defer st.cache.mtxStates.Unlock()
	for _, v := range st.cache.states {
		states = append(states, v)
	}
	return states
}

func (st *StateTracker) GetStatesByRuleUID() map[string][]AlertState {
	ruleMap := make(map[string][]AlertState)
	st.cache.mtxStates.Lock()
	defer st.cache.mtxStates.Unlock()
	for _, state := range st.cache.states {
		if ruleStates, ok := ruleMap[state.AlertRuleUID]; ok {
			ruleStates = append(ruleStates, state)
			ruleMap[state.AlertRuleUID] = ruleStates
		} else {
			ruleStates := []AlertState{state}
			ruleMap[state.AlertRuleUID] = ruleStates
		}
	}
	return ruleMap
}

func (st *StateTracker) cleanUp() {
	ticker := time.NewTicker(time.Duration(60) * time.Minute)
	st.Log.Debug("starting cleanup process", "intervalMinutes", 60)
	for {
		select {
		case <-ticker.C:
			st.trim()
		case <-st.quit:
			st.Log.Debug("stopping cleanup process", "now", time.Now())
			ticker.Stop()
			return
		}
	}
}

func (st *StateTracker) trim() {
	st.Log.Info("trimming alert state cache", "now", time.Now())
	st.cache.mtxStates.Lock()
	defer st.cache.mtxStates.Unlock()
	for _, v := range st.cache.states {
		if len(v.Results) > 100 {
			st.Log.Debug("trimming result set", "cacheId", v.CacheId, "count", len(v.Results)-100)
			newResults := make([]StateEvaluation, 100)
			copy(newResults, v.Results[100:])
			v.Results = newResults
			st.set(v)
		}
	}
}

func (a AlertState) Equals(b AlertState) bool {
	return a.AlertRuleUID == b.AlertRuleUID &&
		a.OrgID == b.OrgID &&
		a.CacheId == b.CacheId &&
		a.Labels.String() == b.Labels.String() &&
		a.State.String() == b.State.String() &&
		a.StartsAt == b.StartsAt &&
		a.EndsAt == b.EndsAt &&
		a.LastEvaluationTime == b.LastEvaluationTime
}

func (st *StateTracker) Put(states []AlertState) {
	for _, s := range states {
		st.set(s)
	}
}

// if duplicate labels exist, keep the value from the first set
func mergeLabels(a, b data.Labels) data.Labels {
	newLbs := data.Labels{}
	for k, v := range a {
		newLbs[k] = v
	}
	for k, v := range b {
		if _, ok := newLbs[k]; !ok {
			newLbs[k] = v
		}
	}
	return newLbs
}
