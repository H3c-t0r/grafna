import { Feature } from 'ol';
import { ComparisonOperation } from '../types';
import { checkFeatureMatchesStyleRule } from './checkFeatureMatchesStyleRule';

describe('check if feature matches style rule', () => {
  it('can compare with numbers', () => {
    const feature = new Feature({
      number: 3,
    });

    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.EQ,
          property: 'number',
          value: 3,
        },
        feature
      )
    ).toEqual(true);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.LT,
          property: 'number',
          value: 2,
        },
        feature
      )
    ).toEqual(false);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.LTE,
          property: 'number',
          value: 3,
        },
        feature
      )
    ).toEqual(true);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.GT,
          property: 'number',
          value: 3,
        },
        feature
      )
    ).toEqual(false);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.GTE,
          property: 'number',
          value: 3,
        },
        feature
      )
    ).toEqual(true);
  });
  it('can compare with strings', () => {
    const feature = new Feature({
      string: 'b',
    });

    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.EQ,
          property: 'string',
          value: 'B',
        },
        feature
      )
    ).toEqual(false);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.LT,
          property: 'string',
          value: 'c',
        },
        feature
      )
    ).toEqual(true);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.LTE,
          property: 'string',
          value: 'bc',
        },
        feature
      )
    ).toEqual(true);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.GT,
          property: 'string',
          value: 'ab',
        },
        feature
      )
    ).toEqual(true);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.GTE,
          property: 'string',
          value: 'abc',
        },
        feature
      )
    ).toEqual(true);
  });
  it('can compare with booleans', () => {
    const feature = new Feature({
      name: 'test polygon',
    });

    feature.setProperties({ boolean: false });

    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.EQ,
          property: 'boolean',
          value: false,
        },
        feature
      )
    ).toEqual(true);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.LT,
          property: 'boolean',
          value: true,
        },
        feature
      )
    ).toEqual(true);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.LTE,
          property: 'boolean',
          value: true,
        },
        feature
      )
    ).toEqual(true);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.GT,
          property: 'boolean',
          value: false,
        },
        feature
      )
    ).toEqual(false);
    expect(
      checkFeatureMatchesStyleRule(
        {
          fillColor: '#1F60C4',
          operation: ComparisonOperation.GTE,
          property: 'boolean',
          value: false,
        },
        feature
      )
    ).toEqual(true);
  });
});
