// To change feature flags, edit:
//  pkg/services/featuremgmt/registry.go
// Then run tests in:
//  pkg/services/featuremgmt/toggles_gen_test.go
// twice to generate and validate the feature flag files

package featuremgmt

var (
	falsePtr = boolPtr(false)
	truePtr  = boolPtr(true)
	// Register each toggle here
	standardFeatureFlags = []FeatureFlag{
		{
			Name:              "disableEnvelopeEncryption",
			Description:       "Disable envelope encryption (emergency only)",
			Stage:             FeatureStageGeneralAvailability,
			Owner:             grafanaAsCodeSquad,
			HideFromAdminPage: true,
			AllowSelfServe:    falsePtr,
		},
		{
			Name:         "live-service-web-worker",
			Description:  "This will use a webworker thread to processes events rather than the main thread",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaAppPlatformSquad,
		},
		{
			Name:         "queryOverLive",
			Description:  "Use Grafana Live WebSocket to execute backend queries",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaAppPlatformSquad,
		},
		{
			Name:              "panelTitleSearch",
			Description:       "Search for dashboards using panel title",
			Stage:             FeatureStagePublicPreview,
			Owner:             grafanaAppPlatformSquad,
			HideFromAdminPage: true,
		},
		{
			Name:           "publicDashboards",
			Description:    "Enables public access to dashboards",
			Stage:          FeatureStageGeneralAvailability,
			Owner:          grafanaSharingSquad,
			Expression:     "true", // enabled by default
			AllowSelfServe: truePtr,
		},
		{
			Name:              "publicDashboardsEmailSharing",
			Description:       "Enables public dashboard sharing to be restricted to only allowed emails",
			Stage:             FeatureStagePublicPreview,
			RequiresLicense:   true,
			Owner:             grafanaSharingSquad,
			HideFromDocs:      true,
			HideFromAdminPage: true,
		},
		{
			Name:        "lokiExperimentalStreaming",
			Description: "Support new streaming approach for loki (prototype, needs special loki build)",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaObservabilityLogsSquad,
		},
		{
			Name:           "featureHighlights",
			Description:    "Highlight Grafana Enterprise features",
			Stage:          FeatureStageGeneralAvailability,
			Owner:          grafanaAsCodeSquad,
			AllowSelfServe: truePtr,
		},
		{
			Name:        "migrationLocking",
			Description: "Lock database during migrations",
			Stage:       FeatureStagePublicPreview,
			Owner:       grafanaBackendPlatformSquad,
		},
		{
			Name:        "storage",
			Description: "Configurable storage for dashboards, datasources, and resources",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaAppPlatformSquad,
		},
		{
			Name:        "correlations",
			Description: "Correlations page",
			Stage:       FeatureStagePublicPreview,
			Owner:       grafanaExploreSquad,
		},
		{
			Name:           "exploreContentOutline",
			Description:    "Content outline sidebar",
			Stage:          FeatureStageGeneralAvailability,
			Owner:          grafanaExploreSquad,
			Expression:     "true", // enabled by default
			FrontendOnly:   true,
			AllowSelfServe: truePtr,
		},
		{
			Name:        "datasourceQueryMultiStatus",
			Description: "Introduce HTTP 207 Multi Status for api/ds/query",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaPluginsPlatformSquad,
		},
		{
			Name:         "traceToMetrics",
			Description:  "Enable trace to metrics links",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityTracesAndProfilingSquad,
		},
		{
			Name:         "autoMigrateOldPanels",
			Description:  "Migrate old angular panels to supported versions (graph, table-old, worldmap, etc)",
			Stage:        FeatureStagePublicPreview,
			FrontendOnly: true,
			Owner:        grafanaDatavizSquad,
		},
		{
			Name:              "disableAngular",
			Description:       "Dynamic flag to disable angular at runtime. The preferred method is to set `angular_support_enabled` to `false` in the [security] settings, which allows you to change the state at runtime.",
			Stage:             FeatureStagePublicPreview,
			FrontendOnly:      true,
			Owner:             grafanaDatavizSquad,
			HideFromAdminPage: true,
		},
		{
			Name:              "canvasPanelNesting",
			Description:       "Allow elements nesting",
			Stage:             FeatureStageExperimental,
			FrontendOnly:      true,
			Owner:             grafanaDatavizSquad,
			HideFromAdminPage: true,
		},
		{
			Name:           "newVizTooltips",
			Description:    "New visualizations tooltips UX",
			Stage:          FeatureStageGeneralAvailability,
			FrontendOnly:   true,
			Owner:          grafanaDatavizSquad,
			AllowSelfServe: truePtr,
		},
		{
			Name:         "scenes",
			Description:  "Experimental framework to build interactive dashboards",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaDashboardsSquad,
		},
		{
			Name:            "disableSecretsCompatibility",
			Description:     "Disable duplicated secret storage in legacy tables",
			Stage:           FeatureStageExperimental,
			RequiresRestart: true,
			Owner:           hostedGrafanaTeam,
		},
		{
			Name:        "logRequestsInstrumentedAsUnknown",
			Description: "Logs the path for requests that are instrumented as unknown",
			Stage:       FeatureStageExperimental,
			Owner:       hostedGrafanaTeam,
		},
		{
			Name:           "dataConnectionsConsole",
			Description:    "Enables a new top-level page called Connections. This page is an experiment that provides a better experience when you install and configure data sources and other plugins.",
			Stage:          FeatureStageGeneralAvailability,
			Expression:     "true", // turned on by default
			Owner:          grafanaPluginsPlatformSquad,
			AllowSelfServe: truePtr,
		},
		{
			// Some plugins rely on topnav feature flag being enabled, so we cannot remove this until we
			// can afford the breaking change, or we've detemined no one else is relying on it
			Name:        "topnav",
			Description: "Enables topnav support in external plugins. The new Grafana navigation cannot be disabled.",
			Stage:       FeatureStageDeprecated,
			Expression:  "true", // enabled by default
			Owner:       grafanaFrontendPlatformSquad,
		},
		{
			Name:         "dockedMegaMenu",
			Description:  "Enable support for a persistent (docked) navigation menu",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaFrontendPlatformSquad,
		},
		{
			Name:              "grpcServer",
			Description:       "Run the GRPC server",
			Stage:             FeatureStagePublicPreview,
			Owner:             grafanaAppPlatformSquad,
			HideFromAdminPage: true,
		},
		{
			Name:            "entityStore",
			Description:     "SQL-based entity store (requires storage flag also)",
			Stage:           FeatureStageExperimental,
			RequiresDevMode: true,
			Owner:           grafanaAppPlatformSquad,
		},
		{
			Name:           "cloudWatchCrossAccountQuerying",
			Description:    "Enables cross-account querying in CloudWatch datasources",
			Stage:          FeatureStageGeneralAvailability,
			Expression:     "true", // enabled by default
			Owner:          awsDatasourcesSquad,
			AllowSelfServe: truePtr,
		},
		{
			Name:           "redshiftAsyncQueryDataSupport",
			Description:    "Enable async query data support for Redshift",
			Stage:          FeatureStageGeneralAvailability,
			Expression:     "true", // enabled by default
			Owner:          awsDatasourcesSquad,
			AllowSelfServe: falsePtr,
		},
		{
			Name:           "athenaAsyncQueryDataSupport",
			Description:    "Enable async query data support for Athena",
			Stage:          FeatureStageGeneralAvailability,
			Expression:     "true", // enabled by default
			FrontendOnly:   true,
			Owner:          awsDatasourcesSquad,
			AllowSelfServe: falsePtr,
		},
		{
			Name:        "cloudwatchNewRegionsHandler",
			Description: "Refactor of /regions endpoint, no user-facing changes",
			Stage:       FeatureStageExperimental,
			Owner:       awsDatasourcesSquad,
		},
		{
			Name:        "showDashboardValidationWarnings",
			Description: "Show warnings when dashboards do not validate against the schema",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaDashboardsSquad,
		},
		{
			Name:        "mysqlAnsiQuotes",
			Description: "Use double quotes to escape keyword in a MySQL query",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaBackendPlatformSquad,
		},
		{
			Name:              "accessControlOnCall",
			Description:       "Access control primitives for OnCall",
			Stage:             FeatureStagePublicPreview,
			Owner:             identityAccessTeam,
			HideFromAdminPage: true,
		},
		{
			Name:        "nestedFolders",
			Description: "Enable folder nesting",
			Stage:       FeatureStagePublicPreview,
			Owner:       grafanaBackendPlatformSquad,
		},
		{
			Name:           "nestedFolderPicker",
			Description:    "Enables the new folder picker to work with nested folders. Requires the nestedFolders feature toggle",
			Stage:          FeatureStageGeneralAvailability,
			Owner:          grafanaFrontendPlatformSquad,
			FrontendOnly:   true,
			Expression:     "true", // enabled by default
			AllowSelfServe: truePtr,
		},
		{
			Name:           "accessTokenExpirationCheck",
			Description:    "Enable OAuth access_token expiration check and token refresh using the refresh_token",
			Stage:          FeatureStageGeneralAvailability,
			Owner:          identityAccessTeam,
			AllowSelfServe: falsePtr,
		},
		{
			Name:              "emptyDashboardPage",
			Description:       "Enable the redesigned user interface of a dashboard page that includes no panels",
			Stage:             FeatureStageGeneralAvailability,
			FrontendOnly:      true,
			Expression:        "true", // enabled by default
			Owner:             grafanaDashboardsSquad,
			AllowSelfServe:    falsePtr,
			HideFromAdminPage: true,
		},
		{
			Name:           "disablePrometheusExemplarSampling",
			Description:    "Disable Prometheus exemplar sampling",
			Stage:          FeatureStageGeneralAvailability,
			Owner:          grafanaObservabilityMetricsSquad,
			AllowSelfServe: truePtr,
		},
		{
			Name:        "alertingBacktesting",
			Description: "Rule backtesting API for alerting",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaAlertingSquad,
		},
		{
			Name:         "editPanelCSVDragAndDrop",
			Description:  "Enables drag and drop for CSV and Excel files",
			FrontendOnly: true,
			Stage:        FeatureStageExperimental,
			Owner:        grafanaBiSquad,
		},
		{
			Name:              "alertingNoNormalState",
			Description:       "Stop maintaining state of alerts that are not firing",
			Stage:             FeatureStagePublicPreview,
			RequiresRestart:   false,
			Owner:             grafanaAlertingSquad,
			HideFromAdminPage: true,
		},
		{
			Name:           "logsContextDatasourceUi",
			Description:    "Allow datasource to provide custom UI for context view",
			Stage:          FeatureStageGeneralAvailability,
			FrontendOnly:   true,
			Owner:          grafanaObservabilityLogsSquad,
			Expression:     "true", // turned on by default
			AllowSelfServe: truePtr,
		},
		{
			Name:           "lokiQuerySplitting",
			Description:    "Split large interval queries into subqueries with smaller time intervals",
			Stage:          FeatureStageGeneralAvailability,
			FrontendOnly:   true,
			Owner:          grafanaObservabilityLogsSquad,
			Expression:     "true", // turned on by default
			AllowSelfServe: truePtr,
		},
		{
			Name:         "lokiQuerySplittingConfig",
			Description:  "Give users the option to configure split durations for Loki queries",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:        "individualCookiePreferences",
			Description: "Support overriding cookie preferences per user",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaBackendPlatformSquad,
		},
		{
			Name:           "prometheusMetricEncyclopedia",
			Description:    "Adds the metrics explorer component to the Prometheus query builder as an option in metric select",
			Expression:     "true",
			Stage:          FeatureStageGeneralAvailability,
			FrontendOnly:   true,
			Owner:          grafanaObservabilityMetricsSquad,
			AllowSelfServe: truePtr,
		},
		{
			Name:           "influxdbBackendMigration",
			Description:    "Query InfluxDB InfluxQL without the proxy",
			Stage:          FeatureStageGeneralAvailability,
			FrontendOnly:   true,
			Owner:          grafanaObservabilityMetricsSquad,
			Expression:     "true", // enabled by default
			AllowSelfServe: falsePtr,
		},
		{
			Name:           "clientTokenRotation",
			Description:    "Replaces the current in-request token rotation so that the client initiates the rotation",
			Stage:          FeatureStageGeneralAvailability,
			Expression:     "true",
			Owner:          identityAccessTeam,
			AllowSelfServe: falsePtr,
		},
		{
			Name:           "prometheusDataplane",
			Description:    "Changes responses to from Prometheus to be compliant with the dataplane specification. In particular, when this feature toggle is active, the numeric `Field.Name` is set from 'Value' to the value of the `__name__` label.",
			Expression:     "true",
			Stage:          FeatureStageGeneralAvailability,
			Owner:          grafanaObservabilityMetricsSquad,
			AllowSelfServe: truePtr,
		},
		{
			Name:           "lokiMetricDataplane",
			Description:    "Changes metric responses from Loki to be compliant with the dataplane specification.",
			Stage:          FeatureStageGeneralAvailability,
			Expression:     "true",
			Owner:          grafanaObservabilityLogsSquad,
			AllowSelfServe: truePtr,
		},
		{
			Name:        "lokiLogsDataplane",
			Description: "Changes logs responses from Loki to be compliant with the dataplane specification.",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaObservabilityLogsSquad,
		},
		{
			Name:           "dataplaneFrontendFallback",
			Description:    "Support dataplane contract field name change for transformations and field name matchers where the name is different",
			Stage:          FeatureStageGeneralAvailability,
			FrontendOnly:   true,
			Expression:     "true",
			Owner:          grafanaObservabilityMetricsSquad,
			AllowSelfServe: truePtr,
		},
		{
			Name:        "disableSSEDataplane",
			Description: "Disables dataplane specific processing in server side expressions.",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaObservabilityMetricsSquad,
		},
		{
			Name:        "alertStateHistoryLokiSecondary",
			Description: "Enable Grafana to write alert state history to an external Loki instance in addition to Grafana annotations.",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaAlertingSquad,
		},
		{
			Name:        "alertStateHistoryLokiPrimary",
			Description: "Enable a remote Loki instance as the primary source for state history reads.",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaAlertingSquad,
		},
		{
			Name:        "alertStateHistoryLokiOnly",
			Description: "Disable Grafana alerts from emitting annotations when a remote Loki instance is available.",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaAlertingSquad,
		},
		{
			Name:        "unifiedRequestLog",
			Description: "Writes error logs to the request logger",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaBackendPlatformSquad,
		},
		{
			Name:              "renderAuthJWT",
			Description:       "Uses JWT-based auth for rendering instead of relying on remote cache",
			Stage:             FeatureStagePublicPreview,
			Owner:             grafanaAsCodeSquad,
			HideFromAdminPage: true,
		},
		{
			Name:            "externalServiceAuth",
			Description:     "Starts an OAuth2 authentication provider for external services",
			Stage:           FeatureStageExperimental,
			RequiresDevMode: true,
			Owner:           identityAccessTeam,
		},
		{
			Name:              "refactorVariablesTimeRange",
			Description:       "Refactor time range variables flow to reduce number of API calls made when query variables are chained",
			Stage:             FeatureStagePublicPreview,
			Owner:             grafanaDashboardsSquad,
			HideFromAdminPage: true, // Non-feature, used to test out a bug fix that impacts the performance of template variables.
		},
		{
			Name:            "useCachingService",
			Description:     "When active, the new query and resource caching implementation using a wire service inject replaces the previous middleware implementation.",
			Stage:           FeatureStageGeneralAvailability,
			Owner:           grafanaOperatorExperienceSquad,
			RequiresRestart: true,
			Expression:      "true", // enabled by default
			AllowSelfServe:  falsePtr,
		},
		{
			Name:           "enableElasticsearchBackendQuerying",
			Description:    "Enable the processing of queries and responses in the Elasticsearch data source through backend",
			Stage:          FeatureStageGeneralAvailability,
			Owner:          grafanaObservabilityLogsSquad,
			Expression:     "true", // enabled by default
			AllowSelfServe: truePtr,
		},
		{
			Name:              "advancedDataSourcePicker",
			Description:       "Enable a new data source picker with contextual information, recently used order and advanced mode",
			Stage:             FeatureStageGeneralAvailability,
			FrontendOnly:      true,
			Expression:        "true", // enabled by default
			Owner:             grafanaDashboardsSquad,
			AllowSelfServe:    falsePtr,
			HideFromAdminPage: true,
		},
		{
			Name:         "faroDatasourceSelector",
			Description:  "Enable the data source selector within the Frontend Apps section of the Frontend Observability",
			Stage:        FeatureStagePublicPreview,
			FrontendOnly: true,
			Owner:        appO11ySquad,
		},
		{
			Name:         "enableDatagridEditing",
			Description:  "Enables the edit functionality in the datagrid panel",
			FrontendOnly: true,
			Stage:        FeatureStagePublicPreview,
			Owner:        grafanaBiSquad,
		},
		{
			Name:         "extraThemes",
			Description:  "Enables extra themes",
			FrontendOnly: true,
			Stage:        FeatureStageExperimental,
			Owner:        grafanaFrontendPlatformSquad,
		},
		{
			Name:         "lokiPredefinedOperations",
			Description:  "Adds predefined query operations to Loki query editor",
			FrontendOnly: true,
			Stage:        FeatureStageExperimental,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:         "pluginsFrontendSandbox",
			Description:  "Enables the plugins frontend sandbox",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaPluginsPlatformSquad,
		},
		{
			Name:         "dashboardEmbed",
			Description:  "Allow embedding dashboard for external use in Code editors",
			FrontendOnly: true,
			Stage:        FeatureStageExperimental,
			Owner:        grafanaAsCodeSquad,
		},
		{
			Name:         "frontendSandboxMonitorOnly",
			Description:  "Enables monitor only in the plugin frontend sandbox (if enabled)",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaPluginsPlatformSquad,
		},
		{
			Name:              "sqlDatasourceDatabaseSelection",
			Description:       "Enables previous SQL data source dataset dropdown behavior",
			FrontendOnly:      true,
			Stage:             FeatureStagePublicPreview,
			Owner:             grafanaBiSquad,
			HideFromAdminPage: true,
		},
		{
			Name:         "lokiFormatQuery",
			Description:  "Enables the ability to format Loki queries",
			FrontendOnly: true,
			Stage:        FeatureStageExperimental,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:           "cloudWatchLogsMonacoEditor",
			Description:    "Enables the Monaco editor for CloudWatch Logs queries",
			Stage:          FeatureStageGeneralAvailability,
			FrontendOnly:   true,
			Expression:     "true", // enabled by default
			Owner:          awsDatasourcesSquad,
			AllowSelfServe: truePtr,
		},
		{
			Name:         "exploreScrollableLogsContainer",
			Description:  "Improves the scrolling behavior of logs in Explore",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:           "recordedQueriesMulti",
			Description:    "Enables writing multiple items from a single query within Recorded Queries",
			Stage:          FeatureStageGeneralAvailability,
			Expression:     "true",
			Owner:          grafanaObservabilityMetricsSquad,
			AllowSelfServe: falsePtr,
		},
		{
			Name:         "pluginsDynamicAngularDetectionPatterns",
			Description:  "Enables fetching Angular detection patterns for plugins from GCOM and fallback to hardcoded ones",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        grafanaPluginsPlatformSquad,
		},
		{
			Name:         "vizAndWidgetSplit",
			Description:  "Split panels between visualizations and widgets",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaDashboardsSquad,
		},
		{
			Name:         "prometheusIncrementalQueryInstrumentation",
			Description:  "Adds RudderStack events to incremental queries",
			FrontendOnly: true,
			Stage:        FeatureStageExperimental,
			Owner:        grafanaObservabilityMetricsSquad,
		},
		{
			Name:         "logsExploreTableVisualisation",
			Description:  "A table visualisation for logs in Explore",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:        "awsDatasourcesTempCredentials",
			Description: "Support temporary security credentials in AWS plugins for Grafana Cloud customers",
			Stage:       FeatureStageExperimental,
			Owner:       awsDatasourcesSquad,
		},
		{
			Name:           "transformationsRedesign",
			Description:    "Enables the transformations redesign",
			Stage:          FeatureStageGeneralAvailability,
			FrontendOnly:   true,
			Expression:     "true", // enabled by default
			Owner:          grafanaObservabilityMetricsSquad,
			AllowSelfServe: truePtr,
		},
		{
			Name:         "mlExpressions",
			Description:  "Enable support for Machine Learning in server-side expressions",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        grafanaAlertingSquad,
		},
		{
			Name:         "traceQLStreaming",
			Description:  "Enables response streaming of TraceQL queries of the Tempo data source",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityTracesAndProfilingSquad,
		},
		{
			Name:         "metricsSummary",
			Description:  "Enables metrics summary queries in the Tempo data source",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityTracesAndProfilingSquad,
		},
		{
			Name:         "grafanaAPIServer",
			Description:  "Enable Kubernetes API Server for Grafana resources",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        grafanaAppPlatformSquad,
		},
		{
			Name:         "grafanaAPIServerWithExperimentalAPIs",
			Description:  "Register experimental APIs with the k8s API server",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        grafanaAppPlatformSquad,
		},
		{
			Name:            "featureToggleAdminPage",
			Description:     "Enable admin page for managing feature toggles from the Grafana front-end",
			Stage:           FeatureStageExperimental,
			FrontendOnly:    false,
			Owner:           grafanaOperatorExperienceSquad,
			RequiresRestart: true,
		},
		{
			Name:        "awsAsyncQueryCaching",
			Description: "Enable caching for async queries for Redshift and Athena. Requires that the `useCachingService` feature toggle is enabled and the datasource has caching and async query support enabled",
			Stage:       FeatureStagePublicPreview,
			Owner:       awsDatasourcesSquad,
		},
		{
			Name:              "splitScopes",
			Description:       "Support faster dashboard and folder search by splitting permission scopes into parts",
			Stage:             FeatureStagePublicPreview,
			FrontendOnly:      false,
			Owner:             identityAccessTeam,
			RequiresRestart:   true,
			HideFromAdminPage: true, // This is internal work to speed up dashboard search, and is not ready for wider use
		},
		{
			Name:         "traceToProfiles",
			Description:  "Enables linking between traces and profiles",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityTracesAndProfilingSquad,
		},
		{
			Name:         "tracesEmbeddedFlameGraph",
			Description:  "Enables embedding a flame graph in traces",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityTracesAndProfilingSquad,
		},
		{
			Name:        "permissionsFilterRemoveSubquery",
			Description: "Alternative permission filter implementation that does not use subqueries for fetching the dashboard folder",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaBackendPlatformSquad,
		},
		{
			Name:           "prometheusConfigOverhaulAuth",
			Description:    "Update the Prometheus configuration page with the new auth component",
			Owner:          grafanaObservabilityMetricsSquad,
			Stage:          FeatureStageGeneralAvailability,
			Expression:     "true", // on by default
			AllowSelfServe: falsePtr,
		},
		{
			Name:            "configurableSchedulerTick",
			Description:     "Enable changing the scheduler base interval via configuration option unified_alerting.scheduler_tick_interval",
			Stage:           FeatureStageExperimental,
			FrontendOnly:    false,
			Owner:           grafanaAlertingSquad,
			RequiresRestart: true,
			HideFromDocs:    true,
		},
		{
			Name:            "influxdbSqlSupport",
			Description:     "Enable InfluxDB SQL query language support with new querying UI",
			Stage:           FeatureStageExperimental,
			FrontendOnly:    false,
			Owner:           grafanaObservabilityMetricsSquad,
			RequiresRestart: false,
		},
		{
			Name:            "alertingNoDataErrorExecution",
			Description:     "Changes how Alerting state manager handles execution of NoData/Error",
			Stage:           FeatureStagePrivatePreview,
			FrontendOnly:    false,
			Owner:           grafanaAlertingSquad,
			RequiresRestart: true,
			Enabled:         true,
		},
		{
			Name:         "angularDeprecationUI",
			Description:  "Display new Angular deprecation-related UI features",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaPluginsPlatformSquad,
		},
		{
			Name:         "dashgpt",
			Description:  "Enable AI powered features in dashboards",
			Stage:        FeatureStagePublicPreview,
			FrontendOnly: true,
			Owner:        grafanaDashboardsSquad,
		},
		{
			Name:            "reportingRetries",
			Description:     "Enables rendering retries for the reporting feature",
			Stage:           FeatureStagePublicPreview,
			FrontendOnly:    false,
			Owner:           grafanaSharingSquad,
			RequiresRestart: true,
		},
		{
			Name:        "sseGroupByDatasource",
			Description: "Send query to the same datasource in a single request when using server side expressions",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaObservabilityMetricsSquad,
		},
		{
			Name:         "requestInstrumentationStatusSource",
			Description:  "Include a status source label for request metrics and logs",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        grafanaPluginsPlatformSquad,
		},
		{
			Name:            "libraryPanelRBAC",
			Description:     "Enables RBAC support for library panels",
			Stage:           FeatureStageExperimental,
			FrontendOnly:    false,
			Owner:           grafanaDashboardsSquad,
			RequiresRestart: true,
		},
		{
			Name:         "lokiRunQueriesInParallel",
			Description:  "Enables running Loki queries in parallel",
			Stage:        FeatureStagePrivatePreview,
			FrontendOnly: false,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:         "wargamesTesting",
			Description:  "Placeholder feature flag for internal testing",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        hostedGrafanaTeam,
		},
		{
			Name:              "alertingInsights",
			Description:       "Show the new alerting insights landing page",
			FrontendOnly:      true,
			Stage:             FeatureStageGeneralAvailability,
			Owner:             grafanaAlertingSquad,
			Expression:        "true", // enabled by default
			AllowSelfServe:    falsePtr,
			HideFromAdminPage: true, // This is moving away from being a feature toggle.
		},
		{
			Name:        "externalCorePlugins",
			Description: "Allow core plugins to be loaded as external",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaPluginsPlatformSquad,
		},
		{
			Name:         "pluginsAPIMetrics",
			Description:  "Sends metrics of public grafana packages usage by plugins",
			FrontendOnly: true,
			Stage:        FeatureStageExperimental,
			Owner:        grafanaPluginsPlatformSquad,
		},
		{
			Name:            "httpSLOLevels",
			Description:     "Adds SLO level to http request metrics",
			Stage:           FeatureStageExperimental,
			FrontendOnly:    false,
			Owner:           hostedGrafanaTeam,
			RequiresRestart: true,
		},
		{
			Name:            "idForwarding",
			Description:     "Generate signed id token for identity that can be forwarded to plugins and external services",
			Stage:           FeatureStageExperimental,
			Owner:           identityAccessTeam,
			RequiresDevMode: true,
		},
		{
			Name:           "cloudWatchWildCardDimensionValues",
			Description:    "Fetches dimension values from CloudWatch to correctly label wildcard dimensions",
			Stage:          FeatureStageGeneralAvailability,
			Expression:     "true", // enabled by default
			Owner:          awsDatasourcesSquad,
			AllowSelfServe: truePtr,
		},
		{
			Name:            "externalServiceAccounts",
			Description:     "Automatic service account and token setup for plugins",
			Stage:           FeatureStageExperimental,
			RequiresDevMode: true,
			Owner:           identityAccessTeam,
		},
		{
			Name:         "panelMonitoring",
			Description:  "Enables panel monitoring through logs and measurements",
			Stage:        FeatureStageExperimental,
			Owner:        grafanaDatavizSquad,
			FrontendOnly: true,
		},
		{
			Name:         "enableNativeHTTPHistogram",
			Description:  "Enables native HTTP Histograms",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        hostedGrafanaTeam,
		},
		{
			Name:         "formatString",
			Description:  "Enable format string transformer",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaBiSquad,
		},
		{
			Name:         "transformationsVariableSupport",
			Description:  "Allows using variables in transformations",
			FrontendOnly: true,
			Stage:        FeatureStagePublicPreview,
			Owner:        grafanaBiSquad,
		},
		{
			Name:            "kubernetesPlaylists",
			Description:     "Use the kubernetes API in the frontend for playlists, and route /api/playlist requests to k8s",
			Stage:           FeatureStageExperimental,
			Owner:           grafanaAppPlatformSquad,
			RequiresRestart: true, // changes the API routing
		},
		{
			Name:        "cloudWatchBatchQueries",
			Description: "Runs CloudWatch metrics queries as separate batches",
			Stage:       FeatureStagePublicPreview,
			Owner:       awsDatasourcesSquad,
		},
		{
			Name:            "recoveryThreshold",
			Description:     "Enables feature recovery threshold (aka hysteresis) for threshold server-side expression",
			Stage:           FeatureStageExperimental,
			FrontendOnly:    false,
			Owner:           grafanaAlertingSquad,
			RequiresRestart: true,
		},
		{
			Name:         "lokiStructuredMetadata",
			Description:  "Enables the loki data source to request structured metadata from the Loki server",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:         "teamHttpHeaders",
			Description:  "Enables datasources to apply team headers to the client requests",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        identityAccessTeam,
		},
		{
			Name:         "awsDatasourcesNewFormStyling",
			Description:  "Applies new form styling for configuration and query editors in AWS plugins",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        awsDatasourcesSquad,
		},
		{
			Name:         "cachingOptimizeSerializationMemoryUsage",
			Description:  "If enabled, the caching backend gradually serializes query responses for the cache, comparing against the configured `[caching]max_value_mb` value as it goes. This can can help prevent Grafana from running out of memory while attempting to cache very large query responses.",
			Stage:        FeatureStageExperimental,
			Owner:        grafanaOperatorExperienceSquad,
			FrontendOnly: false,
		},
		{
			Name:            "panelTitleSearchInV1",
			Description:     "Enable searching for dashboards using panel title in search v1",
			RequiresDevMode: true,
			Stage:           FeatureStageExperimental,
			Owner:           grafanaBackendPlatformSquad,
		},
		{
			Name:         "pluginsInstrumentationStatusSource",
			Description:  "Include a status source label for plugin request metrics and logs",
			FrontendOnly: false,
			Stage:        FeatureStageExperimental,
			Owner:        grafanaPluginsPlatformSquad,
		},
		{
			Name:         "costManagementUi",
			Description:  "Toggles the display of the cost management ui plugin",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        grafanaDatabasesFrontend,
		},
		{
			Name:            "managedPluginsInstall",
			Description:     "Install managed plugins directly from plugins catalog",
			Stage:           FeatureStageExperimental,
			RequiresDevMode: false,
			Owner:           grafanaPluginsPlatformSquad,
		},
		{
			Name:         "prometheusPromQAIL",
			Description:  "Prometheus and AI/ML to assist users in creating a query",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityMetricsSquad,
		},
		{
			Name:         "addFieldFromCalculationStatFunctions",
			Description:  "Add cumulative and window functions to the add field from calculation transformation",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaBiSquad,
		},
		{
			Name:        "alertmanagerRemoteSecondary",
			Description: "Enable Grafana to sync configuration and state with a remote Alertmanager.",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaAlertingSquad,
		},
		{
			Name:        "alertmanagerRemotePrimary",
			Description: "Enable Grafana to have a remote Alertmanager instance as the primary Alertmanager.",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaAlertingSquad,
		},
		{
			Name:        "alertmanagerRemoteOnly",
			Description: "Disable the internal Alertmanager and only use the external one defined.",
			Stage:       FeatureStageExperimental,
			Owner:       grafanaAlertingSquad,
		},
		{
			Name:            "annotationPermissionUpdate",
			Description:     "Separate annotation permissions from dashboard permissions to allow for more granular control.",
			Stage:           FeatureStageExperimental,
			RequiresDevMode: false,
			Owner:           identityAccessTeam,
		},
		{
			Name:         "extractFieldsNameDeduplication",
			Description:  "Make sure extracted field names are unique in the dataframe",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaBiSquad,
		},
		{
			Name:         "dashboardSceneForViewers",
			Description:  "Enables dashboard rendering using Scenes for viewer roles",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaDashboardsSquad,
		},
		{
			Name:         "dashboardScene",
			Description:  "Enables dashboard rendering using scenes for all roles",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaDashboardsSquad,
		},
		{
			Name:         "panelFilterVariable",
			Description:  "Enables use of the `systemPanelFilterVar` variable to filter panels in a dashboard",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaDashboardsSquad,
			HideFromDocs: true,
		},
		{
			Name:         "pdfTables",
			Description:  "Enables generating table data as PDF in reporting",
			Stage:        FeatureStagePrivatePreview,
			FrontendOnly: false,
			Owner:        grafanaSharingSquad,
		},
		{
			Name:            "ssoSettingsApi",
			Description:     "Enables the SSO settings API",
			RequiresDevMode: true,
			Stage:           FeatureStageExperimental,
			FrontendOnly:    false,
			Owner:           identityAccessTeam,
		},
		{
			Name:         "logsInfiniteScrolling",
			Description:  "Enables infinite scrolling for the Logs panel in Explore and Dashboards",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:         "flameGraphItemCollapsing",
			Description:  "Allow collapsing of flame graph items",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityTracesAndProfilingSquad,
		},
		{
			Name:         "alertingDetailsViewV2",
			Description:  "Enables the preview of the new alert details view",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaAlertingSquad,
			HideFromDocs: true,
		},
		{
			Name:         "datatrails",
			Description:  "Enables the new core app datatrails",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaDashboardsSquad,
			HideFromDocs: true,
		},
		{
			Name:         "alertingSimplifiedRouting",
			Description:  "Enables the simplified routing for alerting",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        grafanaAlertingSquad,
			HideFromDocs: true,
		},
		{
			Name:         "logRowsPopoverMenu",
			Description:  "Enable filtering menu displayed when text of a log line is selected",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaObservabilityLogsSquad,
		},
		{
			Name:         "pluginsSkipHostEnvVars",
			Description:  "Disables passing host environment variable to plugin processes",
			Stage:        FeatureStageExperimental,
			FrontendOnly: false,
			Owner:        grafanaPluginsPlatformSquad,
		},
		{
			Name:         "regressionTransformation",
			Description:  "Enables regression analysis transformation",
			Stage:        FeatureStageExperimental,
			FrontendOnly: true,
			Owner:        grafanaBiSquad,
		},
	}
)

func boolPtr(b bool) *bool {
	return &b
}
