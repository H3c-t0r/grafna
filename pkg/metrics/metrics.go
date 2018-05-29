package metrics

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/prometheus/client_golang/prometheus"
)

const exporterName = "grafana"

var (
	M_Instance_Start       prometheus.Counter
	M_Page_Status          *prometheus.CounterVec
	M_Api_Status           *prometheus.CounterVec
	M_Proxy_Status         *prometheus.CounterVec
	M_Http_Request_Total   *prometheus.CounterVec
	M_Http_Request_Summary *prometheus.SummaryVec

	M_Api_User_SignUpStarted   prometheus.Counter
	M_Api_User_SignUpCompleted prometheus.Counter
	M_Api_User_SignUpInvite    prometheus.Counter
	M_Api_Dashboard_Save       prometheus.Summary
	M_Api_Dashboard_Get        prometheus.Summary
	M_Api_Dashboard_Search     prometheus.Summary
	M_Api_Admin_User_Create    prometheus.Counter
	M_Api_Login_Post           prometheus.Counter
	M_Api_Login_OAuth          prometheus.Counter
	M_Api_Org_Create           prometheus.Counter

	M_Api_Dashboard_Snapshot_Create      prometheus.Counter
	M_Api_Dashboard_Snapshot_External    prometheus.Counter
	M_Api_Dashboard_Snapshot_Get         prometheus.Counter
	M_Api_Dashboard_Insert               prometheus.Counter
	M_Alerting_Result_State              *prometheus.CounterVec
	M_Alerting_Notification_Sent         *prometheus.CounterVec
	M_Aws_CloudWatch_GetMetricStatistics prometheus.Counter
	M_Aws_CloudWatch_ListMetrics         prometheus.Counter
	M_DB_DataSource_QueryById            prometheus.Counter

	// Timers
	M_DataSource_ProxyReq_Timer prometheus.Summary
	M_Alerting_Execution_Time   prometheus.Summary

	// Webhook
	M_Webhook_Execution_Time *prometheus.SummaryVec

	// StatTotals
	M_Alerting_Active_Alerts prometheus.Gauge
	M_StatTotal_Dashboards   prometheus.Gauge
	M_StatTotal_Users        prometheus.Gauge
	M_StatActive_Users       prometheus.Gauge
	M_StatTotal_Orgs         prometheus.Gauge
	M_StatTotal_Playlists    prometheus.Gauge
	M_Grafana_Version        *prometheus.GaugeVec
)

func init() {
	M_Instance_Start = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "instance_start_total",
		Help:      "counter for started instances",
		Namespace: exporterName,
	})

	M_Page_Status = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "page_response_status_total",
			Help:      "page http response status",
			Namespace: exporterName,
		},
		[]string{"code"},
	)

	M_Api_Status = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "api_response_status_total",
			Help:      "api http response status",
			Namespace: exporterName,
		},
		[]string{"code"},
	)

	M_Proxy_Status = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "proxy_response_status_total",
			Help:      "proxy http response status",
			Namespace: exporterName,
		},
		[]string{"code"},
	)

	M_Http_Request_Total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "http request counter",
		},
		[]string{"handler", "statuscode", "method"},
	)

	M_Http_Request_Summary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_duration_milliseconds",
			Help: "http request summary",
		},
		[]string{"handler", "statuscode", "method"},
	)

	M_Api_User_SignUpStarted = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "api_user_signup_started_total",
		Help:      "amount of users who started the signup flow",
		Namespace: exporterName,
	})

	M_Api_User_SignUpCompleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "api_user_signup_completed_total",
		Help:      "amount of users who completed the signup flow",
		Namespace: exporterName,
	})

	M_Api_User_SignUpInvite = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "api_user_signup_invite_total",
		Help:      "amount of users who have been invited",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Save = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "api_dashboard_save_milliseconds",
		Help:      "summary for dashboard save duration",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Get = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "api_dashboard_get_milliseconds",
		Help:      "summary for dashboard get duration",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Search = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "api_dashboard_search_milliseconds",
		Help:      "summary for dashboard search duration",
		Namespace: exporterName,
	})

	M_Api_Admin_User_Create = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "api_admin_user_created_total",
		Help:      "api admin user created counter",
		Namespace: exporterName,
	})

	M_Api_Login_Post = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "api_login_post_total",
		Help:      "api login post counter",
		Namespace: exporterName,
	})

	M_Api_Login_OAuth = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "api_login_oauth_total",
		Help:      "api login oauth counter",
		Namespace: exporterName,
	})

	M_Api_Org_Create = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "api_org_create_total",
		Help:      "api org created counter",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Snapshot_Create = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "api_dashboard_snapshot_create_total",
		Help:      "dashboard snapshots created",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Snapshot_External = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "api_dashboard_snapshot_external_total",
		Help:      "external dashboard snapshots created",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Snapshot_Get = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "api_dashboard_snapshot_get_total",
		Help:      "loaded dashboards",
		Namespace: exporterName,
	})

	M_Api_Dashboard_Insert = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "api_models_dashboard_insert_total",
		Help:      "dashboards inserted ",
		Namespace: exporterName,
	})

	M_Alerting_Result_State = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "alerting_result_total",
		Help:      "alert execution result counter",
		Namespace: exporterName,
	}, []string{"state"})

	M_Alerting_Notification_Sent = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "alerting_notification_sent_total",
		Help:      "counter for how many alert notifications been sent",
		Namespace: exporterName,
	}, []string{"type"})

	M_Aws_CloudWatch_GetMetricStatistics = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "aws_cloudwatch_get_metric_statistics_total",
		Help:      "counter for getting metric statistics from aws",
		Namespace: exporterName,
	})

	M_Aws_CloudWatch_ListMetrics = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "aws_cloudwatch_list_metrics_total",
		Help:      "counter for getting list of metrics from aws",
		Namespace: exporterName,
	})

	M_DB_DataSource_QueryById = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "db_datasource_query_by_id_total",
		Help:      "counter for getting datasource by id",
		Namespace: exporterName,
	})

	M_DataSource_ProxyReq_Timer = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "api_dataproxy_request_all_milliseconds",
		Help:      "summary for dataproxy request duration",
		Namespace: exporterName,
	})

	M_Alerting_Execution_Time = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      "alerting_execution_time_milliseconds",
		Help:      "summary of alert exeuction duration",
		Namespace: exporterName,
	})

	M_Webhook_Execution_Time = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:      "webhook_execution_time_milliseconds",
			Help:      "summary of webhook request duration",
			Namespace: exporterName,
		},
		[]string{"type", "statuscode"},
	)

	M_Alerting_Active_Alerts = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "alerting_active_alerts",
		Help:      "amount of active alerts",
		Namespace: exporterName,
	})

	M_StatTotal_Dashboards = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_totals_dashboard",
		Help:      "total amount of dashboards",
		Namespace: exporterName,
	})

	M_StatTotal_Users = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_total_users",
		Help:      "total amount of users",
		Namespace: exporterName,
	})

	M_StatActive_Users = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_active_users",
		Help:      "number of active users",
		Namespace: exporterName,
	})

	M_StatTotal_Orgs = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_total_orgs",
		Help:      "total amount of orgs",
		Namespace: exporterName,
	})

	M_StatTotal_Playlists = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "stat_total_playlists",
		Help:      "total amount of playlists",
		Namespace: exporterName,
	})

	M_Grafana_Version = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "info",
		Help:      "Information about the Grafana",
		Namespace: exporterName,
	}, []string{"version"})
}

func initMetricVars() {
	prometheus.MustRegister(
		M_Instance_Start,
		M_Page_Status,
		M_Api_Status,
		M_Proxy_Status,
		M_Http_Request_Total,
		M_Http_Request_Summary,
		M_Api_User_SignUpStarted,
		M_Api_User_SignUpCompleted,
		M_Api_User_SignUpInvite,
		M_Api_Dashboard_Save,
		M_Api_Dashboard_Get,
		M_Api_Dashboard_Search,
		M_DataSource_ProxyReq_Timer,
		M_Alerting_Execution_Time,
		M_Webhook_Execution_Time,
		M_Api_Admin_User_Create,
		M_Api_Login_Post,
		M_Api_Login_OAuth,
		M_Api_Org_Create,
		M_Api_Dashboard_Snapshot_Create,
		M_Api_Dashboard_Snapshot_External,
		M_Api_Dashboard_Snapshot_Get,
		M_Api_Dashboard_Insert,
		M_Alerting_Result_State,
		M_Alerting_Notification_Sent,
		M_Aws_CloudWatch_GetMetricStatistics,
		M_Aws_CloudWatch_ListMetrics,
		M_DB_DataSource_QueryById,
		M_Alerting_Active_Alerts,
		M_StatTotal_Dashboards,
		M_StatTotal_Users,
		M_StatActive_Users,
		M_StatTotal_Orgs,
		M_StatTotal_Playlists,
		M_Grafana_Version)

}

func updateTotalStats() {
	statsQuery := models.GetSystemStatsQuery{}
	if err := bus.Dispatch(&statsQuery); err != nil {
		metricsLogger.Error("Failed to get system stats", "error", err)
		return
	}

	M_StatTotal_Dashboards.Set(float64(statsQuery.Result.Dashboards))
	M_StatTotal_Users.Set(float64(statsQuery.Result.Users))
	M_StatActive_Users.Set(float64(statsQuery.Result.ActiveUsers))
	M_StatTotal_Playlists.Set(float64(statsQuery.Result.Playlists))
	M_StatTotal_Orgs.Set(float64(statsQuery.Result.Orgs))
}

var usageStatsURL = "https://stats.grafana.org/grafana-usage-report"

func sendUsageStats() {
	if !setting.ReportingEnabled {
		return
	}

	metricsLogger.Debug("Sending anonymous usage stats to stats.grafana.org")

	version := strings.Replace(setting.BuildVersion, ".", "_", -1)

	metrics := map[string]interface{}{}
	report := map[string]interface{}{
		"version": version,
		"metrics": metrics,
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
	}

	statsQuery := models.GetSystemStatsQuery{}
	if err := bus.Dispatch(&statsQuery); err != nil {
		metricsLogger.Error("Failed to get system stats", "error", err)
		return
	}

	metrics["stats.dashboards.count"] = statsQuery.Result.Dashboards
	metrics["stats.users.count"] = statsQuery.Result.Users
	metrics["stats.orgs.count"] = statsQuery.Result.Orgs
	metrics["stats.playlist.count"] = statsQuery.Result.Playlists
	metrics["stats.plugins.apps.count"] = len(plugins.Apps)
	metrics["stats.plugins.panels.count"] = len(plugins.Panels)
	metrics["stats.plugins.datasources.count"] = len(plugins.DataSources)
	metrics["stats.alerts.count"] = statsQuery.Result.Alerts
	metrics["stats.active_users.count"] = statsQuery.Result.ActiveUsers
	metrics["stats.datasources.count"] = statsQuery.Result.Datasources
	metrics["stats.stars.count"] = statsQuery.Result.Stars
	metrics["stats.folders.count"] = statsQuery.Result.Folders
	metrics["stats.dashboard_permissions.count"] = statsQuery.Result.DashboardPermissions
	metrics["stats.folder_permissions.count"] = statsQuery.Result.FolderPermissions
	metrics["stats.provisioned_dashboards.count"] = statsQuery.Result.ProvisionedDashboards
	metrics["stats.snapshots.count"] = statsQuery.Result.Snapshots
	metrics["stats.teams.count"] = statsQuery.Result.Teams

	dsStats := models.GetDataSourceStatsQuery{}
	if err := bus.Dispatch(&dsStats); err != nil {
		metricsLogger.Error("Failed to get datasource stats", "error", err)
		return
	}

	// send counters for each data source
	// but ignore any custom data sources
	// as sending that name could be sensitive information
	dsOtherCount := 0
	for _, dsStat := range dsStats.Result {
		if models.IsKnownDataSourcePlugin(dsStat.Type) {
			metrics["stats.ds."+dsStat.Type+".count"] = dsStat.Count
		} else {
			dsOtherCount += dsStat.Count
		}
	}
	metrics["stats.ds.other.count"] = dsOtherCount

	dsAccessStats := models.GetDataSourceAccessStatsQuery{}
	if err := bus.Dispatch(&dsAccessStats); err != nil {
		metricsLogger.Error("Failed to get datasource access stats", "error", err)
		return
	}

	// send access counters for each data source
	// but ignore any custom data sources
	// as sending that name could be sensitive information
	dsAccessOtherCount := make(map[string]int64)
	for _, dsAccessStat := range dsAccessStats.Result {
		if dsAccessStat.Access == "" {
			continue
		}

		access := strings.ToLower(dsAccessStat.Access)

		if models.IsKnownDataSourcePlugin(dsAccessStat.Type) {
			metrics["stats.ds_access."+dsAccessStat.Type+"."+access+".count"] = dsAccessStat.Count
		} else {
			old := dsAccessOtherCount[access]
			dsAccessOtherCount[access] = old + dsAccessStat.Count
		}
	}

	for access, count := range dsAccessOtherCount {
		metrics["stats.ds_access.other."+access+".count"] = count
	}

	out, _ := json.MarshalIndent(report, "", " ")
	data := bytes.NewBuffer(out)

	client := http.Client{Timeout: 5 * time.Second}
	go client.Post(usageStatsURL, "application/json", data)
}

func NewSummaryTimer(summaryVec *prometheus.SummaryVec, status *int, labels ...string) *prometheus.Timer {
	return prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		if *status == 0 {
			return
		}
		switch {
		case *status >= 500:
			labels = append(labels, "5xx")
		case *status >= 400:
			labels = append(labels, "4xx")
		case *status >= 300:
			labels = append(labels, "3xx")
		case *status >= 200:
			labels = append(labels, "2xx")
		default:
			labels = append(labels, "1xx")
		}
		summaryVec.WithLabelValues(labels...).Observe(v)
	}))
}
