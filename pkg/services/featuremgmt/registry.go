// To change feature flags, edit:
//  pkg/services/featuremgmt/registry.go
// Then run tests in:
//  pkg/services/featuremgmt/toggles_gen_test.go
// twice to generate and validate the feature flag files

package featuremgmt

var (
	// Register each toggle here
	standardFeatureFlags = []FeatureFlag{
		{
			Name:        "alertingBigTransactions",
			Description: "Use big transactions for alerting database writes",
			State:       FeatureStateAlpha,
			Owner:       grafanaAlertingSquad,
		},
		{
			Name:        "trimDefaults",
			Description: "Use cue schema to remove values that will be applied automatically",
			State:       FeatureStateBeta,
			Owner:       grafanaAsCodeSquad,
		},
		{
			Name:        "disableEnvelopeEncryption",
			Description: "Disable envelope encryption (emergency only)",
			State:       FeatureStateStable,
			Owner:       grafanaAsCodeSquad,
		},
		{
			Name:        "database_metrics",
			Description: "Add Prometheus metrics for database tables",
			State:       FeatureStateStable,
			Owner:       hostedGrafanaTeam,
		},
		{
			Name:        "dashboardPreviews",
			Description: "Create and show thumbnails for dashboard search results",
			State:       FeatureStateAlpha,
			Owner:       grafanaAppPlatformSquad,
		},
		{
			Name:         "live-service-web-worker",
			Description:  "This will use a webworker thread to processes events rather than the main thread",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaAppPlatformSquad,
		},
		{
			Name:         "queryOverLive",
			Description:  "Use Grafana Live WebSocket to execute backend queries",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaAppPlatformSquad,
		},
		{
			Name:        "panelTitleSearch",
			Description: "Search for dashboards using panel title",
			State:       FeatureStateBeta,
			Owner:       grafanaAppPlatformSquad,
		},
		{
			Name:        "prometheusAzureOverrideAudience",
			Description: "Experimental. Allow override default AAD audience for Azure Prometheus endpoint",
			State:       FeatureStateBeta,
			Owner:       grafanaObservabilityMetricsSquad,
		},
		{
			Name:        "publicDashboards",
			Description: "Enables public access to dashboards",
			State:       FeatureStateAlpha,
			Owner:       grafanaDashboardsSquad,
		},
		{
			Name:            "publicDashboardsEmailSharing",
			Description:     "Enables public dashboard sharing to be restricted to only allowed emails",
			State:           FeatureStateAlpha,
			RequiresLicense: true,
			Owner:           grafanaDashboardsSquad,
		},
		{
			Name:        "lokiLive",
			Description: "Support WebSocket streaming for loki (early prototype)",
			State:       FeatureStateAlpha,
			Owner:       grafanaObservabilityLogsSquad,
		},
		{
			Name:        "lokiDataframeApi",
			Description: "Use experimental loki api for WebSocket streaming (early prototype)",
			State:       FeatureStateAlpha,
			Owner:       grafanaObservabilityLogsSquad,
		},
		{
			Name:        "featureHighlights",
			Description: "Highlight Grafana Enterprise features",
			State:       FeatureStateStable,
			Owner:       grafanaAsCodeSquad,
		},
		{
			Name:        "migrationLocking",
			Description: "Lock database during migrations",
			State:       FeatureStateBeta,
			Owner:       grafanaBackendPlatformSquad,
		},
		{
			Name:        "storage",
			Description: "Configurable storage for dashboards, datasources, and resources",
			State:       FeatureStateAlpha,
			Owner:       grafanaAppPlatformSquad,
		},
		{
			Name:            "k8s",
			Description:     "Explore native k8s integrations",
			State:           FeatureStateAlpha,
			RequiresDevMode: true,
			Owner:           grafanaAppPlatformSquad,
		},
		{
			Name:         "exploreMixedDatasource",
			Description:  "Enable mixed datasource in Explore",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaExploreSquad,
		},
		{
			Name:         "newTraceView",
			Description:  "Shows the new trace view design",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaObservabilityTracesAndProfilingSquad,
		},
		{
			Name:        "correlations",
			Description: "Correlations page",
			State:       FeatureStateAlpha,
			Owner:       grafanaExploreSquad,
		},
		{
			Name:        "cloudWatchDynamicLabels",
			Description: "Use dynamic labels instead of alias patterns in CloudWatch datasource",
			State:       FeatureStateStable,
			Expression:  "true", // enabled by default
			Owner:       awsPluginsSquad,
		},
		{
			Name:        "datasourceQueryMultiStatus",
			Description: "Introduce HTTP 207 Multi Status for api/ds/query",
			State:       FeatureStateAlpha,
			Owner:       grafanaPluginsPlatformSquad,
		},
		{
			Name:         "traceToMetrics",
			Description:  "Enable trace to metrics links",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaObservabilityTracesAndProfilingSquad,
		},
		{
			Name:        "newDBLibrary",
			Description: "Use jmoiron/sqlx rather than xorm for a few backend services",
			State:       FeatureStateBeta,
			Owner:       grafanaBackendPlatformSquad,
		},
		{
			Name:            "validateDashboardsOnSave",
			Description:     "Validate dashboard JSON POSTed to api/dashboards/db",
			State:           FeatureStateBeta,
			RequiresRestart: true,
			Owner:           grafanaAsCodeSquad,
		},
		{
			Name:         "autoMigrateOldPanels",
			Description:  "Migrate old angular panels to supported versions (graph, table-old, worldmap, etc)",
			State:        FeatureStateBeta,
			FrontendOnly: true,
			Owner:        grafanaDatavizSquad,
		},
		{
			Name:         "disableAngular",
			Description:  "Dynamic flag to disable angular at runtime. The preferred method is to set `angular_support_enabled` to `false` in the [security] settings, which allows you to change the state at runtime.",
			State:        FeatureStateBeta,
			FrontendOnly: true,
			Owner:        grafanaDatavizSquad,
		},
		{
			Name:        "prometheusWideSeries",
			Description: "Enable wide series responses in the Prometheus datasource",
			State:       FeatureStateAlpha,
			Owner:       grafanaObservabilityMetricsSquad,
		},
		{
			Name:         "canvasPanelNesting",
			Description:  "Allow elements nesting",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaDatavizSquad,
		},
		{
			Name:         "scenes",
			Description:  "Experimental framework to build interactive dashboards",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaDashboardsSquad,
		},
		{
			Name:            "disableSecretsCompatibility",
			Description:     "Disable duplicated secret storage in legacy tables",
			State:           FeatureStateAlpha,
			RequiresRestart: true,
			Owner:           hostedGrafanaTeam,
		},
		{
			Name:        "logRequestsInstrumentedAsUnknown",
			Description: "Logs the path for requests that are instrumented as unknown",
			State:       FeatureStateAlpha,
			Owner:       hostedGrafanaTeam,
		},
		{
			Name:        "dataConnectionsConsole",
			Description: "Enables a new top-level page called Connections. This page is an experiment that provides a better experience when you install and configure data sources and other plugins.",
			State:       FeatureStateStable,
			Expression:  "true", // turned on by default
			Owner:       grafanaPluginsPlatformSquad,
		},
		{
			Name:        "internationalization",
			Description: "Enables internationalization",
			State:       FeatureStateStable,
			Expression:  "true", // enabled by default
			Owner:       grafanaUserEssentialsSquad,
		},
		{
			Name:        "topnav",
			Description: "Displays new top nav and page layouts",
			State:       FeatureStateBeta,
			Owner:       grafanaUserEssentialsSquad,
		},
		{
			Name:            "grpcServer",
			Description:     "Run GRPC server",
			State:           FeatureStateAlpha,
			RequiresDevMode: true,
			Owner:           grafanaAppPlatformSquad,
		},
		{
			Name:            "entityStore",
			Description:     "SQL-based entity store (requires storage flag also)",
			State:           FeatureStateAlpha,
			RequiresDevMode: true,
			Owner:           grafanaAppPlatformSquad,
		},
		{
			Name:        "cloudWatchCrossAccountQuerying",
			Description: "Enables cross-account querying in CloudWatch datasources",
			State:       FeatureStateStable,
			Expression:  "true", // enabled by default
			Owner:       awsPluginsSquad,
		},
		{
			Name:         "redshiftAsyncQueryDataSupport",
			Description:  "Enable async query data support for Redshift",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        awsPluginsSquad,
		},
		{
			Name:         "athenaAsyncQueryDataSupport",
			Description:  "Enable async query data support for Athena",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        awsPluginsSquad,
		},
		{
			Name:         "newPanelChromeUI",
			Description:  "Show updated look and feel of grafana-ui PanelChrome: panel header, icons, and menu",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaDashboardsSquad,
		},
		{
			Name:        "showDashboardValidationWarnings",
			Description: "Show warnings when dashboards do not validate against the schema",
			State:       FeatureStateAlpha,
			Owner:       grafanaDashboardsSquad,
		},
		{
			Name:        "mysqlAnsiQuotes",
			Description: "Use double quotes to escape keyword in a MySQL query",
			State:       FeatureStateAlpha,
			Owner:       grafanaBackendPlatformSquad,
		},
		{
			Name:        "accessControlOnCall",
			Description: "Access control primitives for OnCall",
			State:       FeatureStateBeta,
			Owner:       grafanaAuthnzSquad,
		},
		{
			Name:            "nestedFolders",
			Description:     "Enable folder nesting",
			State:           FeatureStateAlpha,
			RequiresDevMode: true,
			Owner:           grafanaBackendPlatformSquad,
		},
		{
			Name:        "accessTokenExpirationCheck",
			Description: "Enable OAuth access_token expiration check and token refresh using the refresh_token",
			State:       FeatureStateStable,
			Owner:       grafanaAuthnzSquad,
		},
		{
			Name:        "elasticsearchBackendMigration",
			Description: "Use Elasticsearch as backend data source",
			State:       FeatureStateAlpha,
			Owner:       grafanaObservabilityLogsSquad,
		},
		{
			Name:        "datasourceOnboarding",
			Description: "Enable data source onboarding page",
			State:       FeatureStateAlpha,
			Owner:       grafanaDashboardsSquad,
		},
		{
			Name:         "emptyDashboardPage",
			Description:  "Enable the redesigned user interface of a dashboard page that includes no panels",
			State:        FeatureStateStable,
			FrontendOnly: true,
			Expression:   "true", // enabled by default
			Owner:        grafanaDashboardsSquad,
		},
		{
			Name:        "secureSocksDatasourceProxy",
			Description: "Enable secure socks tunneling for supported core datasources",
			State:       FeatureStateAlpha,
			Owner:       hostedGrafanaTeam,
		},
		{
			Name:        "authnService",
			Description: "Use new auth service to perform authentication",
			State:       FeatureStateAlpha,
			Owner:       grafanaAuthnzSquad,
		},
		{
			Name:        "disablePrometheusExemplarSampling",
			Description: "Disable Prometheus examplar sampling",
			State:       FeatureStateStable,
			Owner:       grafanaObservabilityMetricsSquad,
		},
		{
			Name:        "alertingBacktesting",
			Description: "Rule backtesting API for alerting",
			State:       FeatureStateAlpha,
			Owner:       grafanaAlertingSquad,
		},
		{
			Name:         "editPanelCSVDragAndDrop",
			Description:  "Enables drag and drop for CSV and Excel files",
			FrontendOnly: true,
			State:        FeatureStateAlpha,
			Owner:        grafanaBiSquad,
		},
		{
			Name:            "alertingNoNormalState",
			Description:     "Stop maintaining state of alerts that are not firing",
			State:           FeatureStateBeta,
			RequiresRestart: false,
			Owner:           grafanaAlertingSquad,
		},
		{

			Name:         "logsSampleInExplore",
			Description:  "Enables access to the logs sample feature in Explore",
			State:        FeatureStateStable,
			Expression:   "true", // turned on by default
			FrontendOnly: true,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:         "logsContextDatasourceUi",
			Description:  "Allow datasource to provide custom UI for context view",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:         "lokiQuerySplitting",
			Description:  "Split large interval queries into subqueries with smaller time intervals",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:         "lokiQuerySplittingConfig",
			Description:  "Give users the option to configure split durations for Loki queries",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:        "individualCookiePreferences",
			Description: "Support overriding cookie preferences per user",
			State:       FeatureStateAlpha,
			Owner:       grafanaBackendPlatformSquad,
		},
		{
			Name:        "onlyExternalOrgRoleSync",
			Description: "Prohibits a user from changing organization roles synced with external auth providers",
			State:       FeatureStateAlpha,
			Owner:       grafanaAuthnzSquad,
		},
		{
			Name:         "drawerDataSourcePicker",
			Description:  "Changes the user experience for data source selection to a drawer.",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaBiSquad,
		},
		{
			Name:         "traceqlSearch",
			Description:  "Enables the 'TraceQL Search' tab for the Tempo datasource which provides a UI to generate TraceQL queries",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaObservabilityTracesAndProfilingSquad,
		},
		{
			Name:         "prometheusMetricEncyclopedia",
			Description:  "Replaces the Prometheus query builder metric select option with a paginated and filterable component",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaObservabilityMetricsSquad,
		},
		{
			Name:         "timeSeriesTable",
			Description:  "Enable time series table transformer & sparkline cell type",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        appO11ySquad,
		},
		{
			Name:         "influxdbBackendMigration",
			Description:  "Query InfluxDB InfluxQL without the proxy",
			State:        FeatureStateAlpha,
			FrontendOnly: true,
			Owner:        grafanaObservabilityMetricsSquad,
		},
		{
			Name:        "clientTokenRotation",
			Description: "Replaces the current in-request token rotation so that the client initiates the rotation",
			State:       FeatureStateAlpha,
			Owner:       grafanaAuthnzSquad,
		},
	}
)
