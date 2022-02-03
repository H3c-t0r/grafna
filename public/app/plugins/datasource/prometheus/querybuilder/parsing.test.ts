import { buildVisualQueryFromString } from './parsing';

describe('buildVisualQueryFromString', () => {
  it('parses simple query', () => {
    expect(buildVisualQueryFromString('counters_logins{app="frontend"}').query).toEqual({
      metric: 'counters_logins',
      labels: [
        {
          op: '=',
          value: 'frontend',
          label: 'app',
        },
      ],
      operations: [],
    });
  });

  it('parses query with rate and interval', () => {
    expect(buildVisualQueryFromString('rate(counters_logins{app="frontend"}[5m])').query).toEqual({
      metric: 'counters_logins',
      labels: [
        {
          op: '=',
          value: 'frontend',
          label: 'app',
        },
      ],
      operations: [
        {
          id: 'rate',
          params: ['5m'],
        },
      ],
    });
  });

  it('parses query with nested query and interval variable', () => {
    expect(
      buildVisualQueryFromString(
        'avg(rate(access_evaluation_duration_count{instance="host.docker.internal:3000"}[$__rate_interval]))'
      ).query
    ).toEqual({
      metric: 'access_evaluation_duration_count',
      labels: [
        {
          op: '=',
          value: 'host.docker.internal:3000',
          label: 'instance',
        },
      ],
      operations: [
        {
          id: 'rate',
          params: ['$__rate_interval'],
        },
        {
          id: 'avg',
          params: [],
        },
      ],
    });
  });

  it('parses query with aggregation by labels', () => {
    const visQuery = {
      metric: 'metric_name',
      labels: [
        {
          label: 'instance',
          op: '=',
          value: 'internal:3000',
        },
      ],
      operations: [
        {
          id: '__sum_by',
          params: ['app', 'version'],
        },
      ],
    };
    expect(buildVisualQueryFromString('sum(metric_name{instance="internal:3000"}) by (app, version)').query).toEqual(
      visQuery
    );
    expect(buildVisualQueryFromString('sum by (app, version)(metric_name{instance="internal:3000"})').query).toEqual(
      visQuery
    );
  });

  it('parses aggregation with params', () => {
    expect(buildVisualQueryFromString('topk(5, http_requests_total)').query).toEqual({
      metric: 'http_requests_total',
      labels: [],
      operations: [
        {
          id: 'topk',
          params: [5],
        },
      ],
    });
  });

  it('parses aggregation with params and labels', () => {
    expect(buildVisualQueryFromString('topk by(instance, job) (5, http_requests_total)').query).toEqual({
      metric: 'http_requests_total',
      labels: [],
      operations: [
        {
          id: '__topk_by',
          params: [5, 'instance', 'job'],
        },
      ],
    });
  });

  it('parses function with multiple arguments', () => {
    expect(
      buildVisualQueryFromString(
        'label_replace(avg_over_time(http_requests_total{instance="foo"}[$__interval]), "instance", "$1", "", "(.*)")'
      ).query
    ).toEqual({
      metric: 'http_requests_total',
      labels: [{ label: 'instance', op: '=', value: 'foo' }],
      operations: [
        {
          id: 'avg_over_time',
          params: ['$__interval'],
        },
        {
          id: 'label_replace',
          params: ['instance', '$1', '', '(.*)'],
        },
      ],
    });
  });

  it('parses binary operation with scalar', () => {
    expect(
      buildVisualQueryFromString('avg_over_time(http_requests_total{instance="foo"}[$__interval]) / 2').query
    ).toEqual({
      metric: 'http_requests_total',
      labels: [{ label: 'instance', op: '=', value: 'foo' }],
      operations: [
        {
          id: 'avg_over_time',
          params: ['$__interval'],
        },
        {
          id: '__divide_by',
          params: [2],
        },
      ],
    });
  });

  it('parses binary operation with 2 queries', () => {
    expect(
      buildVisualQueryFromString('avg_over_time(http_requests_total{instance="foo"}[$__interval]) / sum(logins_count)')
        .query
    ).toEqual({
      metric: 'http_requests_total',
      labels: [{ label: 'instance', op: '=', value: 'foo' }],
      operations: [{ id: 'avg_over_time', params: ['$__interval'] }],
      binaryQueries: [
        {
          operator: '/',
          query: {
            metric: 'logins_count',
            labels: [],
            operations: [{ id: 'sum', params: [] }],
          },
        },
      ],
    });
  });

  it('parses template variables in strings', () => {
    expect(buildVisualQueryFromString('http_requests_total{instance="$label_variable"}').query).toEqual({
      metric: 'http_requests_total',
      labels: [{ label: 'instance', op: '=', value: '$label_variable' }],
      operations: [],
    });
  });

  it('parses template variables for metric', () => {
    expect(buildVisualQueryFromString('$metric_variable{instance="foo"}').query).toEqual({
      metric: '$metric_variable',
      labels: [{ label: 'instance', op: '=', value: 'foo' }],
      operations: [],
    });

    expect(buildVisualQueryFromString('${metric_variable:fmt}{instance="foo"}').query).toEqual({
      metric: '${metric_variable:fmt}',
      labels: [{ label: 'instance', op: '=', value: 'foo' }],
      operations: [],
    });

    expect(buildVisualQueryFromString('[[metric_variable:fmt]]{instance="foo"}').query).toEqual({
      metric: '[[metric_variable:fmt]]',
      labels: [{ label: 'instance', op: '=', value: 'foo' }],
      operations: [],
    });
  });

  it('parses template variables in label name', () => {
    expect(buildVisualQueryFromString('metric{${variable_label}="foo"}').query).toEqual({
      metric: 'metric',
      labels: [{ label: '${variable_label}', op: '=', value: 'foo' }],
      operations: [],
    });
  });
});
