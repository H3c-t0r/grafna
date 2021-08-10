package notifier

import (
	"fmt"
	"regexp"
	"sort"
	"time"

	apimodels "github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	v2 "github.com/prometheus/alertmanager/api/v2"
	"github.com/prometheus/alertmanager/dispatch"
	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/prometheus/alertmanager/types"
	prometheus_model "github.com/prometheus/common/model"
)

type ErrGetAlertsInternal struct {
	OrgID int64
}

func (e ErrGetAlertsInternal) Error() string {
	return fmt.Sprintf("unable to retrieve alerts(s) for org %d due to an internal error", e.OrgID)
}

type ErrGetAlertsBadPayload struct {
	OrgID int64
}

func (e ErrGetAlertsBadPayload) Error() string {
	return fmt.Sprintf("unable to retrieve alerts(s) for org %d as alertmanager is not initialised yet", e.OrgID)
}

type ErrGetAlertsUnavailable struct {
	OrgID int64
}

func (e ErrGetAlertsUnavailable) Error() string {
	return fmt.Sprintf("unable to retrieve alerts for org %d", e.OrgID)
}

type ErrGetAlertGroupsBadPayload struct {
	OrgID int64
}

func (e ErrGetAlertGroupsBadPayload) Error() string {
	return fmt.Sprintf("unable to retrieve alerts groups for org %d", e.OrgID)
}

type ErrGetAlertsNotFound struct {
	OrgID int64
}

func (e ErrGetAlertsNotFound) Error() string {
	return fmt.Sprintf("unable to retrieve alertmanager for org %d", e.OrgID)
}

func (am *Alertmanager) GetAlerts(orgID int64, active, silenced, inhibited bool, filter []string, receivers string) (apimodels.GettableAlerts, error) {
	var (
		// Initialize result slice to prevent api returning `null` when there
		// are no alerts present
		res = apimodels.GettableAlerts{}
	)

	if !am.Ready(orgID) {
		return res, ErrGetAlertsUnavailable{OrgID: orgID}
	}

	matchers, err := parseFilter(filter)
	if err != nil {
		am.logger.Error("failed to parse matchers", "err", err)
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrGetAlertsBadPayload{OrgID: orgID})
	}

	receiverFilter, err := parseReceivers(receivers)
	if err != nil {
		am.logger.Error("failed to parse receiver regex", "err", err)
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrGetAlertsBadPayload{OrgID: orgID})
	}

	amInstance, err := am.GetAM(orgID)
	if err != nil {
		am.logger.Error("failed to get alertmanager instance", "error", err.Error(), "orgID", orgID)
		return nil, fmt.Errorf(" %s: %w", err.Error(), ErrGetAlertsNotFound{OrgID: orgID})
	}
	alerts := amInstance.alerts.GetPending()
	defer alerts.Close()

	alertFilter := amInstance.alertFilter(matchers, silenced, inhibited, active)
	now := time.Now()

	amInstance.reloadConfigMtx.RLock()
	for a := range alerts.Next() {
		if err = alerts.Err(); err != nil {
			break
		}

		routes := amInstance.route.Match(a.Labels)
		receivers := make([]string, 0, len(routes))
		for _, r := range routes {
			receivers = append(receivers, r.RouteOpts.Receiver)
		}

		if receiverFilter != nil && !receiversMatchFilter(receivers, receiverFilter) {
			continue
		}

		if !alertFilter(a, now) {
			continue
		}

		alert := v2.AlertToOpenAPIAlert(a, amInstance.marker.Status(a.Fingerprint()), receivers)

		res = append(res, alert)
	}
	amInstance.reloadConfigMtx.RUnlock()

	if err != nil {
		am.logger.Error("failed to iterate through the alerts", "err", err)
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrGetAlertsInternal{OrgID: orgID})
	}
	sort.Slice(res, func(i, j int) bool {
		return *res[i].Fingerprint < *res[j].Fingerprint
	})

	return res, nil
}

func (am *alertmanager) GetAlertGroups(active, silenced, inhibited bool, filter []string, receivers string) (apimodels.AlertGroups, error) {
	matchers, err := parseFilter(filter)
	if err != nil {
		am.logger.Error("msg", "failed to parse matchers", "err", err, "org", am.orgID)
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrGetAlertGroupsBadPayload{OrgID: am.orgID})
	}

	receiverFilter, err := parseReceivers(receivers)
	if err != nil {
		am.logger.Error("msg", "failed to compile receiver regex", "err", err)
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrGetAlertGroupsBadPayload{OrgID: am.orgID})
	}

	rf := func(receiverFilter *regexp.Regexp) func(r *dispatch.Route) bool {
		return func(r *dispatch.Route) bool {
			receiver := r.RouteOpts.Receiver
			if receiverFilter != nil && !receiverFilter.MatchString(receiver) {
				return false
			}
			return true
		}
	}(receiverFilter)

	af := am.alertFilter(matchers, silenced, inhibited, active)
	alertGroups, allReceivers := am.dispatcher.Groups(rf, af)

	res := make(apimodels.AlertGroups, 0, len(alertGroups))

	for _, alertGroup := range alertGroups {
		ag := &apimodels.AlertGroup{
			Receiver: &apimodels.Receiver{Name: &alertGroup.Receiver},
			Labels:   v2.ModelLabelSetToAPILabelSet(alertGroup.Labels),
			Alerts:   make([]*apimodels.GettableAlert, 0, len(alertGroup.Alerts)),
		}

		for _, alert := range alertGroup.Alerts {
			fp := alert.Fingerprint()
			receivers := allReceivers[fp]
			status := am.marker.Status(fp)
			apiAlert := v2.AlertToOpenAPIAlert(alert, status, receivers)
			ag.Alerts = append(ag.Alerts, apiAlert)
		}
		res = append(res, ag)
	}

	return res, nil
}

func (am *alertmanager) alertFilter(matchers []*labels.Matcher, silenced, inhibited, active bool) func(a *types.Alert, now time.Time) bool {
	return func(a *types.Alert, now time.Time) bool {
		if !a.EndsAt.IsZero() && a.EndsAt.Before(now) {
			return false
		}

		// Set alert's current status based on its label set.
		am.silencer.Mutes(a.Labels)
		am.inhibitor.Mutes(a.Labels)

		// Get alert's current status after seeing if it is suppressed.
		status := am.marker.Status(a.Fingerprint())

		if !active && status.State == types.AlertStateActive {
			return false
		}

		if !silenced && len(status.SilencedBy) != 0 {
			return false
		}

		if !inhibited && len(status.InhibitedBy) != 0 {
			return false
		}

		return alertMatchesFilterLabels(&a.Alert, matchers)
	}
}

func alertMatchesFilterLabels(a *prometheus_model.Alert, matchers []*labels.Matcher) bool {
	sms := make(map[string]string)
	for name, value := range a.Labels {
		sms[string(name)] = string(value)
	}
	return matchFilterLabels(matchers, sms)
}

func matchFilterLabels(matchers []*labels.Matcher, sms map[string]string) bool {
	for _, m := range matchers {
		v, prs := sms[m.Name]
		switch m.Type {
		case labels.MatchNotRegexp, labels.MatchNotEqual:
			if m.Value == "" && prs {
				continue
			}
			if !m.Matches(v) {
				return false
			}
		default:
			if m.Value == "" && !prs {
				continue
			}
			if !m.Matches(v) {
				return false
			}
		}
	}

	return true
}

func parseReceivers(receivers string) (*regexp.Regexp, error) {
	if receivers == "" {
		return nil, nil
	}

	return regexp.Compile("^(?:" + receivers + ")$")
}

func parseFilter(filter []string) ([]*labels.Matcher, error) {
	matchers := make([]*labels.Matcher, 0, len(filter))
	for _, matcherString := range filter {
		matcher, err := labels.ParseMatcher(matcherString)
		if err != nil {
			return nil, err
		}

		matchers = append(matchers, matcher)
	}
	return matchers, nil
}

func receiversMatchFilter(receivers []string, filter *regexp.Regexp) bool {
	for _, r := range receivers {
		if filter.MatchString(r) {
			return true
		}
	}

	return false
}
