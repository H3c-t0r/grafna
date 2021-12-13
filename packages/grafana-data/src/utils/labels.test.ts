import { parseLabels, formatLabels, findCommonLabels, findUniqueLabels, matchAllLabels } from './labels';
import { Labels } from '../types/data';
import { renderLabelsTemplate } from '.';

describe('parseLabels()', () => {
  it('returns no labels on empty labels string', () => {
    expect(parseLabels('')).toEqual({});
    expect(parseLabels('{}')).toEqual({});
  });

  it('returns labels on labels string', () => {
    expect(parseLabels('{foo="bar", baz="42"}')).toEqual({ foo: 'bar', baz: '42' });
  });
});

describe('formatLabels()', () => {
  it('returns no labels on empty label set', () => {
    expect(formatLabels({})).toEqual('');
    expect(formatLabels({}, 'foo')).toEqual('foo');
  });

  it('returns label string on label set', () => {
    expect(formatLabels({ foo: 'bar', baz: '42' })).toEqual('{baz="42", foo="bar"}');
  });
});

describe('findCommonLabels()', () => {
  it('returns no common labels on empty sets', () => {
    expect(findCommonLabels([{}])).toEqual({});
    expect(findCommonLabels([{}, {}])).toEqual({});
  });

  it('returns no common labels on differing sets', () => {
    expect(findCommonLabels([{ foo: 'bar' }, {}])).toEqual({});
    expect(findCommonLabels([{}, { foo: 'bar' }])).toEqual({});
    expect(findCommonLabels([{ baz: '42' }, { foo: 'bar' }])).toEqual({});
    expect(findCommonLabels([{ foo: '42', baz: 'bar' }, { foo: 'bar' }])).toEqual({});
  });

  it('returns the single labels set as common labels', () => {
    expect(findCommonLabels([{ foo: 'bar' }])).toEqual({ foo: 'bar' });
  });
});

describe('findUniqueLabels()', () => {
  it('returns no uncommon labels on empty sets', () => {
    expect(findUniqueLabels({}, {})).toEqual({});
  });

  it('returns all labels given no common labels', () => {
    expect(findUniqueLabels({ foo: '"bar"' }, {})).toEqual({ foo: '"bar"' });
  });

  it('returns all labels except the common labels', () => {
    expect(findUniqueLabels({ foo: '"bar"', baz: '"42"' }, { foo: '"bar"' })).toEqual({ baz: '"42"' });
  });
});

describe('matchAllLabels()', () => {
  it('empty labels do math', () => {
    expect(matchAllLabels({}, {})).toBeTruthy();
  });

  it('missing labels', () => {
    expect(matchAllLabels({ foo: 'bar' }, {})).toBeFalsy();
  });

  it('extra labels should match', () => {
    expect(matchAllLabels({ foo: 'bar' }, { foo: 'bar', baz: '22' })).toBeTruthy();
  });

  it('be graceful with null values (match)', () => {
    expect(matchAllLabels({ foo: 'bar' })).toBeFalsy();
  });

  it('be graceful with null values (match)', () => {
    expect(matchAllLabels((undefined as unknown) as Labels, { foo: 'bar' })).toBeTruthy();
  });
});

describe('renderLabelsTemplate()', () => {
  it('works with empty labels', () => {
    expect(renderLabelsTemplate('hello', {})).toEqual('hello');
  });

  it('Simple replace', () => {
    expect(renderLabelsTemplate('value: {{a}}', { a: 'AAA' })).toEqual('value: AAA');

    // spaces ok
    expect(renderLabelsTemplate('value: {{ a }}', { a: 'AAA' })).toEqual('value: AAA');
  });

  it('Bad syntax', () => {
    expect(renderLabelsTemplate('value: {{a}', { a: 'AAA' })).toEqual('value: {{a}');
    expect(renderLabelsTemplate('value: {a}}}', { a: 'AAA' })).toEqual('value: {a}}}');

    // Current behavior -- not sure if expected or not
    expect(renderLabelsTemplate('value: {{{a}}}', { a: 'AAA' })).toEqual('value: {a}');
  });
});
