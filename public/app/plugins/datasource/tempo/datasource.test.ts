import { lastValueFrom, Observable, of } from 'rxjs';
import { createFetchResponse } from 'test/helpers/createFetchResponse';

import {
  ArrayVector,
  DataFrame,
  dataFrameToJSON,
  DataSourceInstanceSettings,
  FieldType,
  getDefaultTimeRange,
  LoadingState,
  MutableDataFrame,
  PluginType,
} from '@grafana/data';
import { BackendDataSourceResponse, FetchResponse, setBackendSrv, setDataSourceSrv } from '@grafana/runtime';
import config from 'app/core/config';

import {
  DEFAULT_LIMIT,
  TempoJsonData,
  TempoDatasource,
  TempoQuery,
  buildExpr,
  buildLinkExpr,
  getRateAlignedValues,
  makeApmRequest,
  makeTempoLink,
} from './datasource';
import mockJson from './mockJsonResponse.json';
import mockServiceGraph from './mockServiceGraph.json';

jest.mock('@grafana/runtime', () => {
  return {
    ...jest.requireActual('@grafana/runtime'),
    reportInteraction: jest.fn(),
  };
});

describe('Tempo data source', () => {
  // Mock the console error so that running the test suite doesnt throw the error
  const origError = console.error;
  const consoleErrorMock = jest.fn();
  afterEach(() => (console.error = origError));
  beforeEach(() => (console.error = consoleErrorMock));

  it('returns empty response when traceId is empty', async () => {
    const ds = new TempoDatasource(defaultSettings);
    const response = await lastValueFrom(
      ds.query({ targets: [{ refId: 'refid1', queryType: 'traceId', query: '' } as Partial<TempoQuery>] } as any),
      { defaultValue: 'empty' }
    );
    expect(response).toBe('empty');
  });

  describe('Variables should be interpolated correctly', () => {
    function getQuery(): TempoQuery {
      return {
        refId: 'x',
        queryType: 'traceId',
        linkedQuery: {
          refId: 'linked',
          expr: '{instance="$interpolationVar"}',
        },
        query: '$interpolationVar',
        search: '$interpolationVar',
        minDuration: '$interpolationVar',
        maxDuration: '$interpolationVar',
      };
    }

    it('when traceId query for dashboard->explore', async () => {
      const templateSrv: any = { replace: jest.fn() };
      const ds = new TempoDatasource(defaultSettings, templateSrv);
      const text = 'interpolationText';
      templateSrv.replace.mockReturnValue(text);

      const queries = ds.interpolateVariablesInQueries([getQuery()], {
        interpolationVar: { text: text, value: text },
      });
      expect(templateSrv.replace).toBeCalledTimes(5);
      expect(queries[0].linkedQuery?.expr).toBe(text);
      expect(queries[0].query).toBe(text);
      expect(queries[0].search).toBe(text);
      expect(queries[0].minDuration).toBe(text);
      expect(queries[0].maxDuration).toBe(text);
    });

    it('when traceId query for template variable', async () => {
      const templateSrv: any = { replace: jest.fn() };
      const ds = new TempoDatasource(defaultSettings, templateSrv);
      const text = 'interpolationText';
      templateSrv.replace.mockReturnValue(text);

      const resp = ds.applyTemplateVariables(getQuery(), {
        interpolationVar: { text: text, value: text },
      });
      expect(templateSrv.replace).toBeCalledTimes(5);
      expect(resp.linkedQuery?.expr).toBe(text);
      expect(resp.query).toBe(text);
      expect(resp.search).toBe(text);
      expect(resp.minDuration).toBe(text);
      expect(resp.maxDuration).toBe(text);
    });
  });

  it('parses json fields from backend', async () => {
    setupBackendSrv(
      new MutableDataFrame({
        fields: [
          { name: 'traceID', values: ['04450900759028499335'] },
          { name: 'spanID', values: ['4322526419282105830'] },
          { name: 'parentSpanID', values: [''] },
          { name: 'operationName', values: ['store.validateQueryTimeRange'] },
          { name: 'startTime', values: [1619712655875.4539] },
          { name: 'duration', values: [14.984] },
          { name: 'serviceTags', values: ['{"key":"servicetag1","value":"service"}'] },
          { name: 'logs', values: ['{"timestamp":12345,"fields":[{"key":"count","value":1}]}'] },
          { name: 'tags', values: ['{"key":"tag1","value":"val1"}'] },
          { name: 'serviceName', values: ['service'] },
        ],
      })
    );
    const templateSrv: any = { replace: jest.fn() };
    const ds = new TempoDatasource(defaultSettings, templateSrv);
    const response = await lastValueFrom(ds.query({ targets: [{ refId: 'refid1', query: '12345' }] } as any));

    expect(
      (response.data[0] as DataFrame).fields.map((f) => ({
        name: f.name,
        values: f.values.toArray(),
      }))
    ).toMatchObject([
      { name: 'traceID', values: ['04450900759028499335'] },
      { name: 'spanID', values: ['4322526419282105830'] },
      { name: 'parentSpanID', values: [''] },
      { name: 'operationName', values: ['store.validateQueryTimeRange'] },
      { name: 'startTime', values: [1619712655875.4539] },
      { name: 'duration', values: [14.984] },
      { name: 'serviceTags', values: [{ key: 'servicetag1', value: 'service' }] },
      { name: 'logs', values: [{ timestamp: 12345, fields: [{ key: 'count', value: 1 }] }] },
      { name: 'tags', values: [{ key: 'tag1', value: 'val1' }] },
      { name: 'serviceName', values: ['service'] },
    ]);

    expect(
      (response.data[1] as DataFrame).fields.map((f) => ({
        name: f.name,
        values: f.values.toArray(),
      }))
    ).toMatchObject([
      { name: 'id', values: ['4322526419282105830'] },
      { name: 'title', values: ['service'] },
      { name: 'subtitle', values: ['store.validateQueryTimeRange'] },
      { name: 'mainstat', values: ['14.98ms (100%)'] },
      { name: 'secondarystat', values: ['14.98ms (100%)'] },
      { name: 'color', values: [1.000007560204647] },
    ]);

    expect(
      (response.data[2] as DataFrame).fields.map((f) => ({
        name: f.name,
        values: f.values.toArray(),
      }))
    ).toMatchObject([
      { name: 'id', values: [] },
      { name: 'target', values: [] },
      { name: 'source', values: [] },
    ]);
  });

  it('should handle json file upload', async () => {
    const ds = new TempoDatasource(defaultSettings);
    ds.uploadedJson = JSON.stringify(mockJson);
    const response = await lastValueFrom(
      ds.query({
        targets: [{ queryType: 'upload', refId: 'A' }],
      } as any)
    );
    const field = response.data[0].fields[0];
    expect(field.name).toBe('traceID');
    expect(field.type).toBe(FieldType.string);
    expect(field.values.get(0)).toBe('60ba2abb44f13eae');
    expect(field.values.length).toBe(6);
  });

  it('should fail on invalid json file upload', async () => {
    const ds = new TempoDatasource(defaultSettings);
    ds.uploadedJson = JSON.stringify(mockInvalidJson);
    const response = await lastValueFrom(
      ds.query({
        targets: [{ queryType: 'upload', refId: 'A' }],
      } as any)
    );
    expect(response.error?.message).toBeDefined();
    expect(response.data.length).toBe(0);
  });

  it('should handle service graph upload', async () => {
    const ds = new TempoDatasource(defaultSettings);
    ds.uploadedJson = JSON.stringify(mockServiceGraph);
    const response = await lastValueFrom(
      ds.query({
        targets: [{ queryType: 'upload', refId: 'A' }],
      } as any)
    );
    expect(response.data).toHaveLength(2);
    const nodesFrame = response.data[0];
    expect(nodesFrame.name).toBe('Nodes');
    expect(nodesFrame.meta.preferredVisualisationType).toBe('nodeGraph');

    const edgesFrame = response.data[1];
    expect(edgesFrame.name).toBe('Edges');
    expect(edgesFrame.meta.preferredVisualisationType).toBe('nodeGraph');
  });

  it('should build search query correctly', () => {
    const templateSrv: any = { replace: jest.fn() };
    const ds = new TempoDatasource(defaultSettings, templateSrv);
    const duration = '10ms';
    templateSrv.replace.mockReturnValue(duration);
    const tempoQuery: TempoQuery = {
      queryType: 'search',
      refId: 'A',
      query: '',
      serviceName: 'frontend',
      spanName: '/config',
      search: 'root.http.status_code=500',
      minDuration: '$interpolationVar',
      maxDuration: '$interpolationVar',
      limit: 10,
    };
    const builtQuery = ds.buildSearchQuery(tempoQuery);
    expect(builtQuery).toStrictEqual({
      tags: 'root.http.status_code=500 service.name="frontend" name="/config"',
      minDuration: duration,
      maxDuration: duration,
      limit: 10,
    });
  });

  it('should include a default limit', () => {
    const ds = new TempoDatasource(defaultSettings);
    const tempoQuery: TempoQuery = {
      queryType: 'search',
      refId: 'A',
      query: '',
      search: '',
    };
    const builtQuery = ds.buildSearchQuery(tempoQuery);
    expect(builtQuery).toStrictEqual({
      tags: '',
      limit: DEFAULT_LIMIT,
    });
  });

  it('should include time range if provided', () => {
    const ds = new TempoDatasource(defaultSettings);
    const tempoQuery: TempoQuery = {
      queryType: 'search',
      refId: 'A',
      query: '',
      search: '',
    };
    const timeRange = { startTime: 0, endTime: 1000 };
    const builtQuery = ds.buildSearchQuery(tempoQuery, timeRange);
    expect(builtQuery).toStrictEqual({
      tags: '',
      limit: DEFAULT_LIMIT,
      start: timeRange.startTime,
      end: timeRange.endTime,
    });
  });

  it('formats native search query history correctly', () => {
    const ds = new TempoDatasource(defaultSettings);
    const tempoQuery: TempoQuery = {
      queryType: 'nativeSearch',
      refId: 'A',
      query: '',
      serviceName: 'frontend',
      spanName: '/config',
      search: 'root.http.status_code=500',
      minDuration: '1ms',
      maxDuration: '100s',
      limit: 10,
    };
    const result = ds.getQueryDisplayText(tempoQuery);
    expect(result).toBe(
      'Service Name: frontend, Span Name: /config, Search: root.http.status_code=500, Min Duration: 1ms, Max Duration: 100s, Limit: 10'
    );
  });

  it('should get loki search datasource', () => {
    // 1. Get lokiSearch.datasource if present
    const ds1 = new TempoDatasource({
      ...defaultSettings,
      jsonData: {
        lokiSearch: {
          datasourceUid: 'loki-1',
        },
      },
    });
    const lokiDS1 = ds1.getLokiSearchDS();
    expect(lokiDS1).toBe('loki-1');

    // 2. Get traceToLogs.datasource
    const ds2 = new TempoDatasource({
      ...defaultSettings,
      jsonData: {
        tracesToLogs: {
          lokiSearch: true,
          datasourceUid: 'loki-2',
        },
      },
    });
    const lokiDS2 = ds2.getLokiSearchDS();
    expect(lokiDS2).toBe('loki-2');

    // 3. Return undefined if neither is available
    const ds3 = new TempoDatasource(defaultSettings);
    const lokiDS3 = ds3.getLokiSearchDS();
    expect(lokiDS3).toBe(undefined);

    // 4. Return undefined if lokiSearch is undefined, even if traceToLogs is present
    // since this indicates the user cleared the fallback setting
    const ds4 = new TempoDatasource({
      ...defaultSettings,
      jsonData: {
        tracesToLogs: {
          lokiSearch: true,
          datasourceUid: 'loki-2',
        },
        lokiSearch: {
          datasourceUid: undefined,
        },
      },
    });
    const lokiDS4 = ds4.getLokiSearchDS();
    expect(lokiDS4).toBe(undefined);
  });
});

describe('Tempo apm table', () => {
  it('runs service graph queries', async () => {
    const ds = new TempoDatasource({
      ...defaultSettings,
      jsonData: {
        serviceMap: {
          datasourceUid: 'prom',
        },
      },
    });
    config.featureToggles.tempoApmTable = true;
    setDataSourceSrv(backendSrvWithPrometheus as any);
    const response = await lastValueFrom(
      ds.query({ targets: [{ queryType: 'serviceMap' }], range: getDefaultTimeRange() } as any)
    );

    expect(response.data).toHaveLength(3);
    expect(response.state).toBe(LoadingState.Done);

    // APM table
    expect(response.data[0].fields[0].name).toBe('Name');
    expect(response.data[0].fields[0].values.toArray().length).toBe(2);
    expect(response.data[0].fields[0].values.toArray()[0]).toBe('HTTP Client');
    expect(response.data[0].fields[0].values.toArray()[1]).toBe('HTTP GET - root');

    expect(response.data[0].fields[1].name).toBe('Rate');
    expect(response.data[0].fields[1].values.toArray().length).toBe(2);
    expect(response.data[0].fields[1].values.toArray()[0]).toBe(12.75164671814457);
    expect(response.data[0].fields[1].values.toArray()[1]).toBe(12.121331111401608);
    expect(response.data[0].fields[1].config.decimals).toBe(2);
    expect(response.data[0].fields[1].config.links[0].title).toBe('Rate');
    expect(response.data[0].fields[1].config.links[0].internal.query.expr).toBe(
      'topk(5, sum(rate(traces_spanmetrics_calls_total{span_name="${__data.fields[0]}"}[$__rate_interval])) by (span_name))'
    );
    expect(response.data[0].fields[1].config.links[0].internal.query.range).toBe(true);
    expect(response.data[0].fields[1].config.links[0].internal.query.exemplar).toBe(true);
    expect(response.data[0].fields[1].config.links[0].internal.query.instant).toBe(false);

    expect(response.data[0].fields[2].values.toArray().length).toBe(2);
    expect(response.data[0].fields[2].values.toArray()[0]).toBe(12.75164671814457);
    expect(response.data[0].fields[2].values.toArray()[1]).toBe(12.121331111401608);
    expect(response.data[0].fields[2].config.color.mode).toBe('continuous-BlPu');
    expect(response.data[0].fields[2].config.custom.displayMode).toBe('lcd-gauge');
    expect(response.data[0].fields[2].config.decimals).toBe(3);

    expect(response.data[0].fields[3].name).toBe('Error Rate');
    expect(response.data[0].fields[3].values.length).toBe(2);
    expect(response.data[0].fields[3].values[0]).toBe(3.75164671814457);
    expect(response.data[0].fields[3].values[1]).toBe(3.121331111401608);
    expect(response.data[0].fields[3].config.decimals).toBe(2);
    expect(response.data[0].fields[3].config.links[0].title).toBe('Error Rate');
    expect(response.data[0].fields[3].config.links[0].internal.query.expr).toBe(
      'topk(5, sum(rate(traces_spanmetrics_calls_total{span_status="STATUS_CODE_ERROR",span_name="${__data.fields[0]}"}[$__rate_interval])) by (span_name))'
    );
    expect(response.data[0].fields[3].config.links[0].internal.query.range).toBe(true);
    expect(response.data[0].fields[3].config.links[0].internal.query.exemplar).toBe(true);
    expect(response.data[0].fields[3].config.links[0].internal.query.instant).toBe(false);

    expect(response.data[0].fields[4].values.length).toBe(2);
    expect(response.data[0].fields[4].values[0]).toBe(3.75164671814457);
    expect(response.data[0].fields[4].values[1]).toBe(3.121331111401608);
    expect(response.data[0].fields[4].config.color.mode).toBe('continuous-RdYlGr');
    expect(response.data[0].fields[4].config.custom.displayMode).toBe('lcd-gauge');
    expect(response.data[0].fields[4].config.decimals).toBe(3);

    expect(response.data[0].fields[5].name).toBe('Duration (p90)');
    expect(response.data[0].fields[5].values.length).toBe(2);
    expect(response.data[0].fields[5].values[0]).toBe('0');
    expect(response.data[0].fields[5].values[1]).toBe(0.12003505696757232);
    expect(response.data[0].fields[5].config.unit).toBe('s');
    expect(response.data[0].fields[5].config.links[0].title).toBe('Duration');
    expect(response.data[0].fields[5].config.links[0].internal.query.expr).toBe(
      'histogram_quantile(.9, sum(rate(traces_spanmetrics_duration_seconds_bucket{span_status="STATUS_CODE_ERROR",span_name="${__data.fields[0]}"}[$__rate_interval])) by (le))'
    );
    expect(response.data[0].fields[5].config.links[0].internal.query.range).toBe(true);
    expect(response.data[0].fields[5].config.links[0].internal.query.exemplar).toBe(true);
    expect(response.data[0].fields[5].config.links[0].internal.query.instant).toBe(false);

    expect(response.data[0].fields[6].config.links[0].url).toBe('');
    expect(response.data[0].fields[6].config.links[0].title).toBe('Tempo');
    expect(response.data[0].fields[6].config.links[0].internal.query.queryType).toBe('nativeSearch');
    expect(response.data[0].fields[6].config.links[0].internal.query.spanName).toBe('${__data.fields[0]}');

    // Service graph
    expect(response.data[1].name).toBe('Nodes');
    expect(response.data[1].fields[0].values.length).toBe(3);
    expect(response.data[1].fields[0].config.links.length).toBeGreaterThan(0);
    expect(response.data[1].fields[0].config.links).toEqual(serviceGraphLinks);

    expect(response.data[2].name).toBe('Edges');
    expect(response.data[2].fields[0].values.length).toBe(2);
  });

  it('should build expr correctly', () => {
    let targets = { targets: [{ queryType: 'serviceMap' }] } as any;
    let builtQuery = buildExpr(
      { expr: 'topk(5, sum(rate(traces_spanmetrics_calls_total{}[$__range])) by (span_name))', params: [] },
      '',
      targets
    );
    expect(builtQuery).toBe('topk(5, sum(rate(traces_spanmetrics_calls_total{}[$__range])) by (span_name))');

    builtQuery = buildExpr(
      {
        expr: 'topk(5, sum(rate(traces_spanmetrics_calls_total{}[$__range])) by (span_name))',
        params: ['span_status="STATUS_CODE_ERROR"'],
      },
      'span_name=~"HTTP Client|HTTP GET|HTTP GET - root|HTTP POST|HTTP POST - post"',
      targets
    );
    expect(builtQuery).toBe(
      'topk(5, sum(rate(traces_spanmetrics_calls_total{span_status="STATUS_CODE_ERROR",span_name=~"HTTP Client|HTTP GET|HTTP GET - root|HTTP POST|HTTP POST - post"}[$__range])) by (span_name))'
    );

    builtQuery = buildExpr(
      {
        expr: 'histogram_quantile(.9, sum(rate(traces_spanmetrics_duration_seconds_bucket{}[$__range])) by (le))',
        params: ['span_status="STATUS_CODE_ERROR"'],
      },
      'span_name=~"HTTP Client"',
      targets
    );
    expect(builtQuery).toBe(
      'histogram_quantile(.9, sum(rate(traces_spanmetrics_duration_seconds_bucket{span_status="STATUS_CODE_ERROR",span_name=~"HTTP Client"}[$__range])) by (le))'
    );

    targets = { targets: [{ queryType: 'serviceMap', serviceMapQuery: '{client="app",service="app"}' }] } as any;
    builtQuery = buildExpr(
      { expr: 'topk(5, sum(rate(traces_spanmetrics_calls_total{}[$__range])) by (span_name))', params: [] },
      '',
      targets
    );
    expect(builtQuery).toBe(
      'topk(5, sum(rate(traces_spanmetrics_calls_total{service="app",service="app"}[$__range])) by (span_name))'
    );
  });

  it('should build link expr correctly', () => {
    let builtQuery = buildLinkExpr('topk(5, sum(rate(traces_spanmetrics_calls_total{}[$__range])) by (span_name))');
    expect(builtQuery).toBe('topk(5, sum(rate(traces_spanmetrics_calls_total{}[$__rate_interval])) by (span_name))');
  });

  it('should get rate aligned values correctly', () => {
    const resp = [
      {
        refId:
          'topk(5, sum(rate(traces_spanmetrics_calls_total{service="app",service="app"}[$__range])) by (span_name))',
        fields: [
          {
            name: 'Time',
            type: 'time',
            config: {},
            values: [1653828275000, 1653828275000, 1653828275000, 1653828275000, 1653828275000],
          },
          {
            name: 'span_name',
            config: {
              filterable: true,
            },
            type: 'string',
            values: new ArrayVector(['HTTP Client', 'HTTP GET', 'HTTP GET - root', 'HTTP POST', 'HTTP POST - post']),
          },
        ],
      },
    ];
    const objToAlign = {
      'HTTP GET - root': {
        name: 'HTTP GET - root',
        value: 0.2724936652307618,
      },
      'HTTP GET': {
        name: 'HTTP GET',
        value: 0.2724936652307618,
      },
      'HTTP POST - post': {
        name: 'HTTP POST - post',
        value: 0.03697421858453128,
      },
    };

    let value = getRateAlignedValues(resp, objToAlign as any);
    expect(value.toString()).toBe('0,0.2724936652307618,0.2724936652307618,0,0.03697421858453128');
  });

  it('should make apm request correctly', () => {
    const apmRequest = makeApmRequest([
      'topk(5, sum(rate(traces_spanmetrics_calls_total{service="app"}[$__range])) by (span_name))"',
      'histogram_quantile(.9, sum(rate(traces_spanmetrics_duration_seconds_bucket{span_status="STATUS_CODE_ERROR",service="app",service="app",span_name=~"HTTP Client"}[$__range])) by (le))',
    ]);
    expect(apmRequest).toEqual([
      {
        refId: 'topk(5, sum(rate(traces_spanmetrics_calls_total{service="app"}[$__range])) by (span_name))"',
        expr: 'topk(5, sum(rate(traces_spanmetrics_calls_total{service="app"}[$__range])) by (span_name))"',
        instant: true,
      },
      {
        refId:
          'histogram_quantile(.9, sum(rate(traces_spanmetrics_duration_seconds_bucket{span_status="STATUS_CODE_ERROR",service="app",service="app",span_name=~"HTTP Client"}[$__range])) by (le))',
        expr: 'histogram_quantile(.9, sum(rate(traces_spanmetrics_duration_seconds_bucket{span_status="STATUS_CODE_ERROR",service="app",service="app",span_name=~"HTTP Client"}[$__range])) by (le))',
        instant: true,
      },
    ]);
  });

  it('should make tempo link correctly', () => {
    const tempoLink = makeTempoLink('Tempo', '', '"${__data.fields[0]}"', 'gdev-tempo');
    expect(tempoLink).toEqual({
      url: '',
      title: 'Tempo',
      internal: {
        query: {
          queryType: 'nativeSearch',
          spanName: '"${__data.fields[0]}"',
        },
        datasourceUid: 'gdev-tempo',
        datasourceName: 'Tempo',
      },
    });
  });
});

const backendSrvWithPrometheus = {
  async get(uid: string) {
    if (uid === 'prom') {
      return {
        query() {
          return of({
            data: [rateMetric, errorRateMetric, durationMetric, totalsPromMetric, secondsPromMetric, failedPromMetric],
          });
        },
      };
    }
    throw new Error('unexpected uid');
  },
};

function setupBackendSrv(frame: DataFrame) {
  setBackendSrv({
    fetch(): Observable<FetchResponse<BackendDataSourceResponse>> {
      return of(
        createFetchResponse({
          results: {
            refid1: {
              frames: [dataFrameToJSON(frame)],
            },
          },
        })
      );
    },
  } as any);
}

const defaultSettings: DataSourceInstanceSettings<TempoJsonData> = {
  id: 0,
  uid: '0',
  type: 'tracing',
  name: 'tempo',
  access: 'proxy',
  meta: {
    id: 'tempo',
    name: 'tempo',
    type: PluginType.datasource,
    info: {} as any,
    module: '',
    baseUrl: '',
  },
  jsonData: {
    nodeGraph: {
      enabled: true,
    },
  },
};

const rateMetric = new MutableDataFrame({
  refId: 'topk(5, sum(rate(traces_spanmetrics_calls_total{}[$__range])) by (span_name))',
  fields: [
    { name: 'Time', values: [1653725618609, 1653725618609] },
    { name: 'span_name', values: ['HTTP Client', 'HTTP GET - root'] },
    {
      name: 'Value #topk(5, sum(rate(traces_spanmetrics_calls_total{}[$__range])) by (span_name))',
      values: [12.75164671814457, 12.121331111401608],
    },
  ],
});

const errorRateMetric = new MutableDataFrame({
  refId:
    'topk(5, sum(rate(traces_spanmetrics_calls_total{span_status="STATUS_CODE_ERROR",span_name=~"HTTP Client|HTTP GET - root"}[$__range])) by (span_name))',
  fields: [
    { name: 'Time', values: [1653725618609, 1653725618609] },
    { name: 'span_name', values: ['HTTP Client', 'HTTP GET - root'] },
    {
      name: 'Value #topk(5, sum(rate(traces_spanmetrics_calls_total{span_status="STATUS_CODE_ERROR"}[$__range])) by (span_name))',
      values: [3.75164671814457, 3.121331111401608],
    },
  ],
});

const durationMetric = new MutableDataFrame({
  refId:
    'histogram_quantile(.9, sum(rate(traces_spanmetrics_duration_seconds_bucket{span_status="STATUS_CODE_ERROR",span_name=~"HTTP GET - root"}[$__range])) by (le))',
  fields: [
    { name: 'Time', values: [1653725618609] },
    {
      name: 'Value #histogram_quantile(.9, sum(rate(traces_spanmetrics_duration_seconds_bucket{span_status="STATUS_CODE_ERROR",span_name=~"HTTP GET - root"}[$__range])) by (le))',
      values: [0.12003505696757232],
    },
  ],
});

const totalsPromMetric = new MutableDataFrame({
  refId: 'traces_service_graph_request_total',
  fields: [
    { name: 'Time', values: [1628169788000, 1628169788000] },
    { name: 'client', values: ['app', 'lb'] },
    { name: 'instance', values: ['127.0.0.1:12345', '127.0.0.1:12345'] },
    { name: 'job', values: ['local_scrape', 'local_scrape'] },
    { name: 'server', values: ['db', 'app'] },
    { name: 'tempo_config', values: ['default', 'default'] },
    { name: 'Value #traces_service_graph_request_total', values: [10, 20] },
  ],
});

const secondsPromMetric = new MutableDataFrame({
  refId: 'traces_service_graph_request_server_seconds_sum',
  fields: [
    { name: 'Time', values: [1628169788000, 1628169788000] },
    { name: 'client', values: ['app', 'lb'] },
    { name: 'instance', values: ['127.0.0.1:12345', '127.0.0.1:12345'] },
    { name: 'job', values: ['local_scrape', 'local_scrape'] },
    { name: 'server', values: ['db', 'app'] },
    { name: 'tempo_config', values: ['default', 'default'] },
    { name: 'Value #traces_service_graph_request_server_seconds_sum', values: [10, 40] },
  ],
});

const failedPromMetric = new MutableDataFrame({
  refId: 'traces_service_graph_request_failed_total',
  fields: [
    { name: 'Time', values: [1628169788000, 1628169788000] },
    { name: 'client', values: ['app', 'lb'] },
    { name: 'instance', values: ['127.0.0.1:12345', '127.0.0.1:12345'] },
    { name: 'job', values: ['local_scrape', 'local_scrape'] },
    { name: 'server', values: ['db', 'app'] },
    { name: 'tempo_config', values: ['default', 'default'] },
    { name: 'Value #traces_service_graph_request_failed_total', values: [2, 15] },
  ],
});

const mockInvalidJson = {
  batches: [
    {
      resource: {
        attributes: [],
      },
      instrumentation_library_spans: [
        {
          instrumentation_library: {},
          spans: [
            {
              trace_id: 'AAAAAAAAAABguiq7RPE+rg==',
              span_id: 'cmteMBAvwNA=',
              parentSpanId: 'OY8PIaPbma4=',
              name: 'HTTP GET - root',
              kind: 'SPAN_KIND_SERVER',
              startTimeUnixNano: '1627471657255809000',
              endTimeUnixNano: '1627471657256268000',
              attributes: [
                { key: 'http.status_code', value: { intValue: '200' } },
                { key: 'http.method', value: { stringValue: 'GET' } },
                { key: 'http.url', value: { stringValue: '/' } },
                { key: 'component', value: { stringValue: 'net/http' } },
              ],
              status: {},
            },
          ],
        },
      ],
    },
  ],
};

const serviceGraphLinks = [
  {
    url: '',
    title: 'Request rate',
    internal: {
      query: {
        expr: 'sum by (client, server)(rate(traces_service_graph_request_total{server="${__data.fields.id}"}[$__rate_interval]))',
        instant: false,
        range: true,
        exemplar: true,
      },
      datasourceUid: 'prom',
      datasourceName: 'Prometheus',
    },
  },
  {
    url: '',
    title: 'Request histogram',
    internal: {
      query: {
        expr: 'histogram_quantile(0.9, sum(rate(traces_service_graph_request_server_seconds_bucket{server="${__data.fields.id}"}[$__rate_interval])) by (le, client, server))',
        instant: false,
        range: true,
        exemplar: true,
      },
      datasourceUid: 'prom',
      datasourceName: 'Prometheus',
    },
  },
  {
    url: '',
    title: 'Failed request rate',
    internal: {
      query: {
        expr: 'sum by (client, server)(rate(traces_service_graph_request_failed_total{server="${__data.fields.id}"}[$__rate_interval]))',
        instant: false,
        range: true,
        exemplar: true,
      },
      datasourceUid: 'prom',
      datasourceName: 'Prometheus',
    },
  },
  {
    url: '',
    title: 'View traces',
    internal: {
      query: {
        queryType: 'nativeSearch',
        serviceName: '${__data.fields[0]}',
      } as TempoQuery,
      datasourceUid: 'tempo',
      datasourceName: 'Tempo',
    },
  },
];
