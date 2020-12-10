import { reducerTester } from 'test/core/redux/reducerTester';
import { reducer } from './reducer';
import {
  addMetric,
  changeMetricAttribute,
  changeMetricField,
  changeMetricMeta,
  changeMetricSetting,
  changeMetricType,
  removeMetric,
  toggleMetricVisibility,
} from './actions';
import { Derivative, ExtendedStats, MetricAggregation } from '../aggregations';
import { defaultMetricAgg } from '../../../../query_def';
import { metricAggregationConfig } from '../utils';
import { initQuery } from '../../state';

describe('Metric Aggregations Reducer', () => {
  it('should correctly add new aggregations', () => {
    const firstAggregation: MetricAggregation = {
      id: '1',
      type: 'count',
    };
    const secondAggregation: MetricAggregation = {
      id: '2',
      type: 'count',
    };

    reducerTester()
      .givenReducer(reducer, [])
      .whenActionIsDispatched(addMetric(firstAggregation.id))
      .thenStateShouldEqual([firstAggregation])
      .whenActionIsDispatched(addMetric(secondAggregation.id))
      .thenStateShouldEqual([firstAggregation, secondAggregation]);
  });

  describe('When removing aggregations', () => {
    it('Should correctly remove aggregations', () => {
      const firstAggregation: MetricAggregation = {
        id: '1',
        type: 'count',
      };
      const secondAggregation: MetricAggregation = {
        id: '2',
        type: 'count',
      };

      reducerTester()
        .givenReducer(reducer, [firstAggregation, secondAggregation])
        .whenActionIsDispatched(removeMetric(firstAggregation.id))
        .thenStateShouldEqual([secondAggregation]);
    });

    it('Should insert a default aggregation when the last one is removed', () => {
      const initialState: MetricAggregation[] = [{ id: '2', type: 'avg' }];

      reducerTester()
        .givenReducer(reducer, initialState)
        .whenActionIsDispatched(removeMetric(initialState[0].id))
        .thenStateShouldEqual([defaultMetricAgg()]);
    });
  });

  describe("When changing existing aggregation's type", () => {
    it('Should correctly change type to selected aggregation', () => {
      const firstAggregation: MetricAggregation = {
        id: '1',
        type: 'count',
      };
      const secondAggregation: MetricAggregation = {
        id: '2',
        type: 'count',
      };

      const expectedSecondAggregation: MetricAggregation = { ...secondAggregation, type: 'avg' };

      reducerTester()
        .givenReducer(reducer, [firstAggregation, secondAggregation])
        .whenActionIsDispatched(changeMetricType(secondAggregation.id, expectedSecondAggregation.type))
        .thenStateShouldEqual([firstAggregation, { ...secondAggregation, type: expectedSecondAggregation.type }]);
    });

    it('Should remove all other aggregations when the newly selected one is `isSingleMetric`', () => {
      const firstAggregation: MetricAggregation = {
        id: '1',
        type: 'count',
      };
      const secondAggregation: MetricAggregation = {
        id: '2',
        type: 'count',
      };

      const expectedAggregation: MetricAggregation = {
        ...secondAggregation,
        type: 'raw_data',
        ...metricAggregationConfig['raw_data'].defaults,
      };

      reducerTester()
        .givenReducer(reducer, [firstAggregation, secondAggregation])
        .whenActionIsDispatched(changeMetricType(secondAggregation.id, expectedAggregation.type))
        .thenStateShouldEqual([expectedAggregation]);
    });
  });

  it("Should correctly change aggregation's field", () => {
    const firstAggregation: MetricAggregation = {
      id: '1',
      type: 'count',
    };
    const secondAggregation: MetricAggregation = {
      id: '2',
      type: 'moving_fn',
    };

    const expectedSecondAggregation = {
      ...secondAggregation,
      field: 'new field',
    };

    reducerTester()
      .givenReducer(reducer, [firstAggregation, secondAggregation])
      .whenActionIsDispatched(changeMetricField(secondAggregation.id, expectedSecondAggregation.field))
      .thenStateShouldEqual([firstAggregation, expectedSecondAggregation]);
  });

  it('Should correctly toggle `hide` field', () => {
    const firstAggregation: MetricAggregation = {
      id: '1',
      type: 'count',
    };

    const secondAggregation: MetricAggregation = {
      id: '2',
      type: 'count',
    };

    reducerTester()
      .givenReducer(reducer, [firstAggregation, secondAggregation])
      .whenActionIsDispatched(toggleMetricVisibility(firstAggregation.id))
      .thenStateShouldEqual([{ ...firstAggregation, hide: true }, secondAggregation])
      .whenActionIsDispatched(toggleMetricVisibility(firstAggregation.id))
      .thenStateShouldEqual([{ ...firstAggregation, hide: false }, secondAggregation]);
  });

  it("Should correctly change aggregation's settings", () => {
    const firstAggregation: Derivative = {
      id: '1',
      type: 'derivative',
      settings: {
        unit: 'Some unit',
      },
    };
    const secondAggregation: MetricAggregation = {
      id: '2',
      type: 'count',
    };

    const expectedSettings: typeof firstAggregation['settings'] = {
      unit: 'Changed unit',
    };

    reducerTester()
      .givenReducer(reducer, [firstAggregation, secondAggregation])
      .whenActionIsDispatched(changeMetricSetting(firstAggregation, 'unit', expectedSettings.unit!))
      .thenStateShouldEqual([{ ...firstAggregation, settings: expectedSettings }, secondAggregation]);
  });

  it("Should correctly change aggregation's meta", () => {
    const firstAggregation: ExtendedStats = {
      id: '1',
      type: 'extended_stats',
      meta: {
        avg: true,
      },
    };
    const secondAggregation: MetricAggregation = {
      id: '2',
      type: 'count',
    };

    const expectedMeta: typeof firstAggregation['meta'] = {
      avg: false,
    };

    reducerTester()
      .givenReducer(reducer, [firstAggregation, secondAggregation])
      .whenActionIsDispatched(changeMetricMeta(firstAggregation, 'avg', expectedMeta.avg!))
      .thenStateShouldEqual([{ ...firstAggregation, meta: expectedMeta }, secondAggregation]);
  });

  it("Should correctly change aggregation's attribute", () => {
    const firstAggregation: ExtendedStats = {
      id: '1',
      type: 'extended_stats',
    };
    const secondAggregation: MetricAggregation = {
      id: '2',
      type: 'count',
    };

    const expectedHide: typeof firstAggregation['hide'] = false;

    reducerTester()
      .givenReducer(reducer, [firstAggregation, secondAggregation])
      .whenActionIsDispatched(changeMetricAttribute(firstAggregation, 'hide', expectedHide))
      .thenStateShouldEqual([{ ...firstAggregation, hide: expectedHide }, secondAggregation]);
  });

  it('Should not change state with other action types', () => {
    const initialState: MetricAggregation[] = [
      {
        id: '1',
        type: 'count',
      },
    ];

    reducerTester()
      .givenReducer(reducer, initialState)
      .whenActionIsDispatched({ type: 'THIS ACTION SHOULD NOT HAVE ANY EFFECT IN THIS REDUCER' })
      .thenStateShouldEqual(initialState);
  });

  it('Should correctly initialize first Metric Aggregation', () => {
    reducerTester()
      .givenReducer(reducer, [])
      .whenActionIsDispatched(initQuery())
      .thenStateShouldEqual([defaultMetricAgg('1')]);
  });
});
