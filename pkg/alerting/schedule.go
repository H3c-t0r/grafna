package alerting

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/grafana/grafana/pkg/api"
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

// Job is a job for an alert execution
// note that LastPointTs is a time denoting the timestamp of the last point to run against
// this way the check runs always on the right data, irrespective of execution delays
// that said, for convenience, we track the generatedAt timestamp
type Job struct {
	OrgId           int64
	MonitorId       int64
	EndpointId      int64
	EndpointName    string
	EndpointSlug    string
	Settings        map[string]string
	MonitorTypeName string
	Notifications   m.MonitorNotificationSetting
	Freq            int64
	Offset          int64 // offset on top of "even" minute/10s/.. intervals
	Definition      CheckDef
	GeneratedAt     time.Time
	LastPointTs     time.Time
	StoreMetricFunc func(m *m.MetricDefinition) `json:"-"`
	AssertMinSeries int                         // to verify during execution at least this many series are returned (would be nice at some point to include actual number of collectors)
	AssertStart     time.Time                   // to verify timestamps in response
	AssertStep      int                         // to verify step duration
	AssertSteps     int                         // to verify during execution this many points are included
}

func (job Job) String() string {
	return fmt.Sprintf("<Job> monitorId=%d generatedAt=%s lastPointTs=%s definition: %s", job.MonitorId, job.GeneratedAt, job.LastPointTs, job.Definition)
}

func (job Job) StoreResult(res m.CheckEvalResult) {
	if job.StoreMetricFunc == nil {
		return
	}
	if !setting.WriteIndividualAlertResults {
		return
	}
	metrics := make([]*m.MetricDefinition, 3)
	metricNames := [3]string{"ok_state", "warn_state", "error_state"}
	for pos, state := range metricNames {
		metrics[pos] = &m.MetricDefinition{
			OrgId:      job.OrgId,
			Name:       fmt.Sprintf("health.%s.%s.%s", job.EndpointSlug, strings.ToLower(job.MonitorTypeName), state),
			Metric:     fmt.Sprintf("health.%s.%s", strings.ToLower(job.MonitorTypeName), state),
			Interval:   job.Freq,
			Value:      0.0,
			Unit:       "state",
			Time:       job.LastPointTs.Unix(),
			TargetType: "gauge",
			Extra: map[string]interface{}{
				"endpoint_id": job.EndpointId,
				"monitor_id":  job.MonitorId,
			},
		}
	}
	if int(res) >= 0 {
		metrics[int(res)].Value = 1.0
	}
	for _, m := range metrics {
		job.StoreMetricFunc(m)
	}
}

func getJobs(lastPointAt int64) ([]*Job, error) {

	query := m.GetMonitorsForAlertsQuery{
		Timestamp: lastPointAt,
	}

	if err := bus.Dispatch(&query); err != nil {
		return nil, err
	}

	jobs := make([]*Job, 0)
	for _, monitor := range query.Result {
		job := buildJobForMonitor(monitor, lastPointAt)
		if job != nil {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func buildJobForMonitor(monitor *m.MonitorForAlertDTO, lastPointAt int64) *Job {
	//state could in theory be ok, warn, error, but we only use ok vs error for now

	if monitor.HealthSettings == nil {
		return nil
	}

	if monitor.Frequency == 0 || monitor.HealthSettings.Steps == 0 || monitor.HealthSettings.NumCollectors == 0 {
		//fmt.Printf("bad monitor definition given: %#v", monitor)
		return nil
	}

	// let's say it takes at least warmupPeriod, after job creation, to start getting data.
	period := int(monitor.Frequency) // what we call frequency is actually the period, like 60s
	warmupPeriod := time.Duration((monitor.HealthSettings.Steps+1)*period) * time.Second
	if monitor.Created.After(time.Unix(lastPointAt, 0).Add(-warmupPeriod)) {
		return nil
	}

	type Settings struct {
		EndpointSlug    string
		MonitorTypeName string
		Duration        string
		NumCollectors   int
		Steps           int
	}

	// graphite behaves like so:
	// from is exclusive (from=foo returns data at ts=foo+1 and higher)
	// until is inclusive (until=bar returns data at ts=bar and lower)
	// so if lastPointAt is 1000, and Steps = 3 and Frequency is 10
	// we want points with timestamps 980, 990, 1000
	// we can just query from 970

	settings := Settings{
		EndpointSlug:    monitor.EndpointSlug,
		MonitorTypeName: monitor.MonitorTypeName,
		Duration:        fmt.Sprintf("%d", int64(monitor.HealthSettings.Steps)*monitor.Frequency),
		NumCollectors:   monitor.HealthSettings.NumCollectors,
		Steps:           monitor.HealthSettings.Steps,
	}

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	// note: in graphite, using the series-wise sum(), sum(1+null) = 1, and sum(null+null) gets dropped from the result!
	// in bosun, series-wise sum() doesn't exist, you can only sum over time. (though functions like t() help)
	// when bosun pulls in graphite results, null values are removed from the series.
	// we get from graphite the raw (incl nulls) series, so that we can inspect and log/instrument nulls
	// bosun does all the logic as follows: see how many collectors are errorring .Steps in a row, using streak
	// transpose that, to get a 1/0 for each sufficiently-erroring collector, sum them together and compare to the threshold.

	// note: it may look like the end of the queried interval is ambiguous here, and if offset > frequency, may include "too recent" values by accident.
	// fear not, as when we execute the alert in the executor, we set the lastPointTs as end time

	target := `litmus.{{.EndpointSlug}}.*.{{.MonitorTypeName | ToLower }}.error_state`
	tpl := `sum(t(streak(graphite("` + target + `", "{{.Duration}}s", "", "")) == {{.Steps}} , "")) >= {{.NumCollectors}}`

	var t = template.Must(template.New("query").Funcs(funcMap).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, settings)
	if err != nil {
		panic(fmt.Sprintf("Could not execute alert query template: %q", err))
	}
	j := &Job{
		MonitorId:       monitor.Id,
		EndpointId:      monitor.EndpointId,
		EndpointName:    monitor.EndpointName,
		EndpointSlug:    monitor.EndpointSlug,
		Settings:        monitor.SettingsMap(),
		MonitorTypeName: monitor.MonitorTypeName,
		Notifications:   monitor.HealthSettings.Notifications,
		OrgId:           monitor.OrgId,
		Freq:            monitor.Frequency,
		Offset:          monitor.Offset,
		Definition: CheckDef{
			CritExpr: b.String(),
			WarnExpr: "0", // for now we have only good or bad. so only crit is needed
		},
		StoreMetricFunc: api.StoreMetric,
		AssertMinSeries: monitor.HealthSettings.NumCollectors,
		AssertStep:      int(monitor.Frequency),
		AssertSteps:     monitor.HealthSettings.Steps,
	}
	return j
}
