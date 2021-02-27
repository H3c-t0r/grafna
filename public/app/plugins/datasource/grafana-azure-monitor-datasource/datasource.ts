import _ from 'lodash';
import AzureMonitorDatasource from './azure_monitor/azure_monitor_datasource';
import AppInsightsDatasource from './app_insights/app_insights_datasource';
import AzureLogAnalyticsDatasource from './azure_log_analytics/azure_log_analytics_datasource';
import {
  AzureDataSourceJsonData,
  AzureLogsVariable,
  AzureMonitorQuery,
  AzureQueryType,
  InsightsAnalyticsQuery,
} from './types';
import {
  DataFrame,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceApi,
  DataSourceInstanceSettings,
  LoadingState,
  ScopedVars,
} from '@grafana/data';
import { forkJoin, Observable, of } from 'rxjs';
import { DataSourceWithBackend, getBackendSrv } from '@grafana/runtime';
import InsightsAnalyticsDatasource from './insights_analytics/insights_analytics_datasource';
import { migrateMetricsDimensionFilters } from './query_ctrl';
import { map } from 'rxjs/operators';

export interface LogConfigs {
  subscriptionId?: string;
  queryType: string;
  postfix: string;
  azureMonitorPostfix: string;
  sameAsAzureMonitor?: boolean;
  getWorkspacesOrResources: (subscription: string, url: string) => Promise<AzureLogsVariable[]>;
}

export default class Datasource extends DataSourceApi<AzureMonitorQuery, AzureDataSourceJsonData> {
  azureMonitorDatasource: AzureMonitorDatasource;
  appInsightsDatasource: AppInsightsDatasource;
  azureLogAnalyticsDatasource: AzureLogAnalyticsDatasource;
  azureResourceLogAnalyticsDatasource: AzureLogAnalyticsDatasource;
  insightsAnalyticsDatasource: InsightsAnalyticsDatasource;

  pseudoDatasource: Record<AzureQueryType, DataSourceWithBackend>;
  optionsKey: Record<AzureQueryType, string>;

  constructor(instanceSettings: DataSourceInstanceSettings<AzureDataSourceJsonData>) {
    super(instanceSettings);
    const workspaceLogConfigs: LogConfigs = {
      subscriptionId: instanceSettings.jsonData.logAnalyticsSubscriptionId,
      queryType: AzureQueryType.LogAnalytics,
      postfix: 'loganalyticsazure',
      azureMonitorPostfix: 'workspacesloganalytics',
      sameAsAzureMonitor: instanceSettings.jsonData.azureLogAnalyticsSameAs,
      getWorkspacesOrResources: async (subscription: string, url: string) => {
        const response = await getBackendSrv().datasourceRequest({
          url:
            url +
            `/azuremonitor/subscriptions/${subscription}/providers/Microsoft.OperationalInsights/workspaces?api-version=2017-04-26-preview`,
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
        });
        return (
          _.map(response.data.value, (val: any) => {
            return { text: val.name, value: val.properties.customerId };
          }) || []
        );
      },
    };
    const resourceLogConfigs: LogConfigs = {
      subscriptionId: instanceSettings.jsonData.resourceLogAnalyticsSubscriptionId,
      queryType: AzureQueryType.ResourceLogAnalytics,
      postfix: 'resourceloganalyticsazure',
      azureMonitorPostfix: 'resourcesloganalytics',
      sameAsAzureMonitor: instanceSettings.jsonData.azureResourceLogAnalyticsSameAs,
      getWorkspacesOrResources: async (subscription: string, url: string) => {
        const result = await getBackendSrv().datasourceRequest({
          url: url + `/azuremonitor/providers/Microsoft.ResourceGraph/resources?api-version=2019-04-01`,
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          data: { subscriptions: [subscription], query: 'project id, name' },
        });
        return (
          _.map(result.data.data.rows, (val: any) => {
            return { text: val[1], value: val[0] };
          }) || []
        );
      },
    };
    this.azureMonitorDatasource = new AzureMonitorDatasource(instanceSettings);
    this.appInsightsDatasource = new AppInsightsDatasource(instanceSettings);
    this.azureLogAnalyticsDatasource = new AzureLogAnalyticsDatasource(instanceSettings, workspaceLogConfigs);
    this.azureResourceLogAnalyticsDatasource = new AzureLogAnalyticsDatasource(instanceSettings, resourceLogConfigs);
    this.insightsAnalyticsDatasource = new InsightsAnalyticsDatasource(instanceSettings);

    const pseudoDatasource: any = {};
    pseudoDatasource[AzureQueryType.ApplicationInsights] = this.appInsightsDatasource;
    pseudoDatasource[AzureQueryType.AzureMonitor] = this.azureMonitorDatasource;
    pseudoDatasource[AzureQueryType.InsightsAnalytics] = this.insightsAnalyticsDatasource;
    pseudoDatasource[AzureQueryType.LogAnalytics] = this.azureLogAnalyticsDatasource;
    pseudoDatasource[AzureQueryType.ResourceLogAnalytics] = this.azureResourceLogAnalyticsDatasource;
    this.pseudoDatasource = pseudoDatasource;

    const optionsKey: any = {};
    optionsKey[AzureQueryType.ApplicationInsights] = 'appInsights';
    optionsKey[AzureQueryType.AzureMonitor] = 'azureMonitor';
    optionsKey[AzureQueryType.InsightsAnalytics] = 'insightsAnalytics';
    optionsKey[AzureQueryType.LogAnalytics] = 'azureLogAnalytics';
    optionsKey[AzureQueryType.ResourceLogAnalytics] = 'azureResourceLogAnalytics';
    this.optionsKey = optionsKey;
  }

  query(options: DataQueryRequest<AzureMonitorQuery>): Observable<DataQueryResponse> {
    const byType: Record<AzureQueryType, DataQueryRequest<AzureMonitorQuery>> = ({} as unknown) as Record<
      AzureQueryType,
      DataQueryRequest<AzureMonitorQuery>
    >;

    for (const target of options.targets) {
      // Migrate old query structure
      if (target.queryType === AzureQueryType.ApplicationInsights) {
        if ((target.appInsights as any).rawQuery) {
          target.queryType = AzureQueryType.InsightsAnalytics;
          target.insightsAnalytics = (target.appInsights as unknown) as InsightsAnalyticsQuery;
          delete target.appInsights;
        }
      }
      if (!target.queryType) {
        target.queryType = AzureQueryType.AzureMonitor;
      }

      if (target.queryType === AzureQueryType.AzureMonitor) {
        migrateMetricsDimensionFilters(target.azureMonitor);
      }

      // Check that we have options
      const opts = (target as any)[this.optionsKey[target.queryType]];

      // Skip hidden queries or ones without properties
      if (target.hide || !opts) {
        continue;
      }

      // Initialize the list of queries
      let q = byType[target.queryType];
      if (!q) {
        q = _.cloneDeep(options);
        q.requestId = `${q.requestId}-${target.refId}`;
        q.targets = [];
        byType[target.queryType] = q;
      }
      q.targets.push(target);
    }

    // Distinct types are managed by distinct requests
    const obs = Object.keys(byType).map((type: AzureQueryType) => {
      const req = byType[type];
      return this.pseudoDatasource[type].query(req);
    });
    // Single query can skip merge
    if (obs.length === 1) {
      return obs[0];
    }

    if (obs.length > 1) {
      return forkJoin(obs).pipe(
        map((results: DataQueryResponse[]) => {
          const data: DataFrame[] = [];
          for (const result of results) {
            for (const frame of result.data) {
              data.push(frame);
            }
          }

          return { state: LoadingState.Done, data };
        })
      );
    }

    return of({ state: LoadingState.Done, data: [] });
  }

  async annotationQuery(options: any) {
    return this.azureLogAnalyticsDatasource.annotationQuery(options);
  }

  async metricFindQuery(query: string) {
    if (!query) {
      return Promise.resolve([]);
    }

    const aiResult = this.appInsightsDatasource.metricFindQueryInternal(query);
    if (aiResult) {
      return aiResult;
    }

    const amResult = this.azureMonitorDatasource.metricFindQueryInternal(query);
    if (amResult) {
      return amResult;
    }

    const arlaResult = query.includes('resource')
      ? this.azureResourceLogAnalyticsDatasource.metricFindQueryInternal(query)
      : null;
    if (arlaResult) {
      return arlaResult;
    }

    const alaResult = this.azureLogAnalyticsDatasource.metricFindQueryInternal(query);
    if (alaResult) {
      return alaResult;
    }

    return Promise.resolve([]);
  }

  async testDatasource() {
    const promises: any[] = [];

    if (this.azureMonitorDatasource.isConfigured()) {
      promises.push(this.azureMonitorDatasource.testDatasource());
    }

    if (this.appInsightsDatasource.isConfigured()) {
      promises.push(this.appInsightsDatasource.testDatasource());
    }

    if (this.azureLogAnalyticsDatasource.isConfigured()) {
      promises.push(this.azureLogAnalyticsDatasource.testDatasource());
    }

    if (this.azureResourceLogAnalyticsDatasource.isConfigured()) {
      promises.push(this.azureResourceLogAnalyticsDatasource.testDatasource());
    }

    if (promises.length === 0) {
      return {
        status: 'error',
        message: `Nothing configured. At least one of the API's must be configured.`,
        title: 'Error',
      };
    }

    return Promise.all(promises).then((results) => {
      let status = 'success';
      let message = '';

      for (let i = 0; i < results.length; i++) {
        if (results[i].status !== 'success') {
          status = results[i].status;
        }
        message += `${i + 1}. ${results[i].message} `;
      }

      return {
        status: status,
        message: message,
        title: _.upperFirst(status),
      };
    });
  }

  /* Azure Monitor REST API methods */
  getResourceGroups(subscriptionId: string) {
    return this.azureMonitorDatasource.getResourceGroups(subscriptionId);
  }

  getMetricDefinitions(subscriptionId: string, resourceGroup: string) {
    return this.azureMonitorDatasource.getMetricDefinitions(subscriptionId, resourceGroup);
  }

  getResourceNames(subscriptionId: string, resourceGroup: string, metricDefinition: string) {
    return this.azureMonitorDatasource.getResourceNames(subscriptionId, resourceGroup, metricDefinition);
  }

  getMetricNames(
    subscriptionId: string,
    resourceGroup: string,
    metricDefinition: string,
    resourceName: string,
    metricNamespace: string
  ) {
    return this.azureMonitorDatasource.getMetricNames(
      subscriptionId,
      resourceGroup,
      metricDefinition,
      resourceName,
      metricNamespace
    );
  }

  getMetricNamespaces(subscriptionId: string, resourceGroup: string, metricDefinition: string, resourceName: string) {
    return this.azureMonitorDatasource.getMetricNamespaces(
      subscriptionId,
      resourceGroup,
      metricDefinition,
      resourceName
    );
  }

  getMetricMetadata(
    subscriptionId: string,
    resourceGroup: string,
    metricDefinition: string,
    resourceName: string,
    metricNamespace: string,
    metricName: string
  ) {
    return this.azureMonitorDatasource.getMetricMetadata(
      subscriptionId,
      resourceGroup,
      metricDefinition,
      resourceName,
      metricNamespace,
      metricName
    );
  }

  /* Application Insights API method */
  getAppInsightsMetricNames() {
    return this.appInsightsDatasource.getMetricNames();
  }

  getAppInsightsMetricMetadata(metricName: string) {
    return this.appInsightsDatasource.getMetricMetadata(metricName);
  }

  getAppInsightsColumns(refId: string | number) {
    return this.appInsightsDatasource.logAnalyticsColumns[refId];
  }

  /*Azure Log Analytics */
  getAzureLogAnalyticsWorkspaces(subscriptionId: string) {
    return this.azureLogAnalyticsDatasource.getWorkspacesOrResources(subscriptionId);
  }

  getSubscriptions() {
    return this.azureMonitorDatasource.getSubscriptions();
  }

  interpolateVariablesInQueries(queries: AzureMonitorQuery[], scopedVars: ScopedVars): AzureMonitorQuery[] {
    return queries.map(
      (query) => this.pseudoDatasource[query.queryType].applyTemplateVariables(query, scopedVars) as AzureMonitorQuery
    );
  }
}
