import _ from 'lodash';
import {
  AzureDataSourceJsonData,
  AzureSubscription as IAzureSubscription,
  AzureResourceGraphQuery as IAzureResourceGraphQuery,
} from '../types';
import { DataSourceInstanceSettings } from '@grafana/data';
import { BackendSrv } from 'app/core/services/backend_srv';
import { TemplateSrv } from 'app/features/templating/template_srv';

export class AzureSubscription {
  id: string;
  subscriptionId: string;
  tenantId: string;
  displayName: string;
  name: string;
  constructor(result: IAzureSubscription) {
    this.id = result.id || result.subscriptionId;
    this.subscriptionId = result.subscriptionId || result.id;
    this.tenantId = result.tenantId;
    this.displayName = result.displayName;
    this.name = result.name || result.displayName || '';
  }
}

export class AzureSubscriptions {
  subscriptions: IAzureSubscription[] = [];
  constructor(results: any) {
    if (results && results.data && results.data.value && results.data.value.length > 0) {
      this.subscriptions = results.data.value.map((sub: IAzureSubscription) => new AzureSubscription(sub));
    }
  }
  getSubscriptionIds(): string[] {
    return this.subscriptions.map((sub: IAzureSubscription) => sub.subscriptionId.toString());
  }
}

export class AzureResourceGraphResponseParser {
  output: any = {
    type: 'table',
    columns: [],
    rows: [],
  };
  constructor(response: any[]) {
    _.each(response, res => {
      if (res.result && res.result.data && res.result.data.data) {
        const tmpOutput = res.result.data.data || {};
        const output = this.parseTableResult(tmpOutput.columns, tmpOutput.rows);
        this.output.columns = output.columns;
        this.output.rows = _.concat(this.output.rows, output.rows);
      }
    });
  }
  parseTableResult(columns: any[], rows: any[]): any {
    columns = (columns || []).map((column: any, index: number) => {
      column.text = column.name || index;
      column.type = column.type === 'integer' ? 'number' : column.type || 'string';
      return column;
    });
    rows = (rows || []).map((row: any) => {
      return row.map((rowItem: any, index: number) => {
        if (columns[index] && columns[index].type === 'number') {
          return +rowItem;
        } else if (typeof rowItem === 'string') {
          return rowItem;
        } else {
          return JSON.stringify(rowItem);
        }
      });
    });
    const output: any = {
      type: 'table',
      columns,
      rows,
    };
    return output;
  }
  getResultsAsVariablesList() {
    const returnvalues: any[] = [];
    _.each(this.output.rows, row => {
      _.each(row, col => {
        returnvalues.push({
          value: col,
          text: col,
        });
      });
    });
    return returnvalues;
  }
}

export class AzureResourceGraphQuery implements IAzureResourceGraphQuery {
  query: string;
  top: number;
  skip: number;
  subscriptions: string[];
  constructor(query: string, top: number, skip: number, subscriptions: string[]) {
    this.query = query;
    this.top = top;
    this.skip = skip;
    this.subscriptions = subscriptions;
  }
}

export default class AzureResourceGraphDatasource {
  id: number;
  url: string;
  cloudName: string;
  baseUrl: string;
  allSubscriptions: IAzureSubscription[];

  /** @ngInject */
  constructor(
    instanceSettings: DataSourceInstanceSettings<AzureDataSourceJsonData>,
    private backendSrv: BackendSrv,
    private templateSrv: TemplateSrv
  ) {
    this.id = instanceSettings.id;
    this.url = instanceSettings.url;
    this.cloudName = instanceSettings.jsonData.cloudName || 'azuremonitor';
    this.baseUrl = `/resourcegraph`;
    this.allSubscriptions = [];
  }

  getSubscriptionIds() {
    const url = `/${this.cloudName}/subscriptions?api-version=2019-03-01`;
    return this.doSubscriptionsRequest(url).then((result: any) => {
      const subscriptions = new AzureSubscriptions(result);
      this.allSubscriptions = subscriptions.subscriptions;
      return subscriptions.getSubscriptionIds();
    });
  }

  doSubscriptionsRequest(url: string, maxRetries = 1): any {
    return this.backendSrv
      .datasourceRequest({
        url: this.url + url,
        method: 'GET',
      })
      .catch((error: any) => {
        if (maxRetries > 0) {
          return this.doSubscriptionsRequest(url, maxRetries - 1);
        }
        throw error;
      });
  }

  async doResourceGraphRequest(query: AzureResourceGraphQuery, maxRetries = 1): Promise<any> {
    let subscriptions = [];
    if (this.allSubscriptions.length === 0) {
      subscriptions = await this.getSubscriptionIds();
    } else {
      subscriptions = this.allSubscriptions.map((sub: any) => sub.subscriptionId.toString());
    }
    if (query.subscriptions.indexOf('all') < 0) {
      const filteredSubscriptions = query.subscriptions.filter(s => s !== 'all');
      if (filteredSubscriptions.length > 0) {
        subscriptions = filteredSubscriptions;
      }
    }
    return this.backendSrv
      .datasourceRequest({
        url: this.url + this.baseUrl + '?api-version=2019-04-01',
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        data: {
          query: query.query,
          subscriptions,
          options: {
            $top: query.top || 100,
            $skip: query.skip || 0,
            resultFormat: 'table',
          },
        },
      })
      .catch((error: any) => {
        if (maxRetries > 0) {
          return this.doResourceGraphRequest(query, maxRetries - 1);
        }
        throw error;
      });
  }

  doQueries(queries: AzureResourceGraphQuery[]) {
    return _.map(queries, query => {
      return this.doResourceGraphRequest(query)
        .then((result: any) => {
          return { result, query };
        })
        .catch((error: any) => {
          throw { error, query };
        });
    });
  }

  async query(options: any) {
    const queries: AzureResourceGraphQuery[] = _.filter(options.targets, item => {
      return item.hide !== true && item.azureResourceGraph && item.azureResourceGraph.query;
    }).map((target: any) => {
      const item: AzureResourceGraphQuery = target.azureResourceGraph;
      const queryOption = new AzureResourceGraphQuery(
        this.templateSrv.replace(item.query, options.scopedVars),
        item.top || 1000,
        item.skip || 0,
        item.subscriptions || ['all']
      );
      return queryOption;
    });
    if (!queries || queries.length === 0) {
      return;
    }
    const promises = this.doQueries(queries);
    return Promise.all(promises).then(results => {
      const responseParser = new AzureResourceGraphResponseParser(results);
      return responseParser.output;
    });
  }

  metricFindQuery(query: string) {
    if (query.startsWith(`ResourceGraph(`) && query.endsWith(`)`)) {
      const resourceGraphQuery = query.replace(`ResourceGraph(`, ``).slice(0, -1);
      const queryOption = new AzureResourceGraphQuery(this.templateSrv.replace(resourceGraphQuery), 1000, 0, ['all']);
      const promises = this.doQueries([queryOption]);
      return Promise.all(promises).then(results => {
        const responseParser = new AzureResourceGraphResponseParser(results);
        return responseParser.getResultsAsVariablesList();
      });
    }
    return undefined;
  }
}
