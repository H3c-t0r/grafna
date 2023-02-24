import {
  DataLink,
  DataQuery,
  ExplorePanelsState,
  Field,
  InternalDataLink,
  InterpolateFunction,
  LinkModel,
  ScopedVars,
  SplitOpen,
  TimeRange,
} from '../types';

import { locationUtil } from './location';
import { serializeStateToUrlParam } from './url';

export const DataLinkBuiltInVars = {
  keepTime: '__url_time_range',
  timeRangeFrom: '__from',
  timeRangeTo: '__to',
  includeVars: '__all_variables',
  seriesName: '__series.name',
  fieldName: '__field.name',
  valueTime: '__value.time',
  valueNumeric: '__value.numeric',
  valueText: '__value.text',
  valueRaw: '__value.raw',
  // name of the calculation represented by the value
  valueCalc: '__value.calc',
};

// We inject these because we cannot import them directly as they reside inside grafana main package.
export type LinkToExploreOptions = {
  link: DataLink;
  scopedVars: ScopedVars;
  range: TimeRange;
  field: Field;
  internalLink: InternalDataLink;
  onClickFn?: SplitOpen;
  replaceVariables: InterpolateFunction;
};

export function mapInternalLinkToExplore(options: LinkToExploreOptions): LinkModel<Field> {
  const { onClickFn, replaceVariables, link, scopedVars, range, field, internalLink } = options;

  const interpolatedQuery = interpolateObject(link.internal?.query, scopedVars, replaceVariables);
  const interpolatedPanelsState = interpolateObject(link.internal?.panelsState, scopedVars, replaceVariables);
  const title = link.title ? link.title : internalLink.datasourceName;

  return {
    title: replaceVariables(title, scopedVars),
    // In this case this is meant to be internal link (opens split view by default) the href will also points
    // to explore but this way you can open it in new tab.
    href: generateInternalHref(internalLink.datasourceUid, interpolatedQuery, range, interpolatedPanelsState),
    onClick: onClickFn
      ? () => {
          onClickFn({
            datasourceUid: internalLink.datasourceUid,
            query: interpolatedQuery,
            panelsState: interpolatedPanelsState,
            range,
          });
        }
      : undefined,
    target: link?.targetBlank ? '_blank' : '_self',
    origin: field,
  };
}

/**
 * Generates href for internal derived field link.
 */
function generateInternalHref<T extends DataQuery = any>(
  datasourceUid: string,
  query: T,
  range: TimeRange,
  panelsState?: ExplorePanelsState
): string {
  return locationUtil.assureBaseUrl(
    `/explore?left=${encodeURIComponent(
      serializeStateToUrlParam({
        range: range.raw,
        datasource: datasourceUid,
        queries: [query],
        panelsState: panelsState,
      })
    )}`
  );
}

function interpolateObject<T>(
  obj: T | undefined,
  scopedVars: ScopedVars,
  replaceVariables: InterpolateFunction
): T | undefined {
  if (!obj) {
    return obj;
  }
  if (typeof obj === 'string') {
    // @ts-ignore this is complaining we are returning string, but we are checking if obj is a string so should be fine.
    return replaceVariables(obj, scopedVars);
  }
  const copy = JSON.parse(JSON.stringify(obj));
  return interpolateObjectRecursive(copy, scopedVars, replaceVariables);
}

function interpolateObjectRecursive<T extends Object>(
  obj: T,
  scopedVars: ScopedVars,
  replaceVariables: InterpolateFunction
): T {
  for (const k of Object.keys(obj)) {
    // Honestly not sure how to type this to make TS happy.
    // @ts-ignore
    if (typeof obj[k] === 'string') {
      // @ts-ignore
      obj[k] = replaceVariables(obj[k], scopedVars);
      // @ts-ignore
    } else if (typeof obj[k] === 'object' && obj[k] !== null) {
      // @ts-ignore
      obj[k] = interpolateObjectRecursive(obj[k], scopedVars, replaceVariables);
    }
  }
  return obj;
}

/**
 * This function takes some code from  template service replace() function to figure out if all variables are
 * interpolated. This is so we don't show links that do not work. This cuts a lots of corners though and that is why
 * it's a local function. We sort of don't care about the dashboard template variables for example. Also we only link
 * to loki/splunk/elastic, so it should be less probable that user needs part of a query that looks like a variable but
 * is actually part of the query language.
 * @param query
 * @param scopedVars
 * @param getVarMap
 */
export function dataLinkHasAllVariablesDefined<T extends DataLink>(
  query: T,
  scopedVars: ScopedVars,
  getVarMap: Function
): boolean {
  const vars = getVarMap(getStringsFromObject(query), scopedVars);
  console.log('varmap', vars, scopedVars);
  return Object.values(vars).every((val) => val !== undefined && val !== null);
}

function getStringsFromObject<T extends Object>(obj: T): string {
  let acc = '';
  for (const k of Object.keys(obj)) {
    // Honestly not sure how to type this to make TS happy.
    // @ts-ignore
    if (typeof obj[k] === 'string') {
      // @ts-ignore
      acc += ' ' + obj[k];
      // @ts-ignore
    } else if (typeof obj[k] === 'object' && obj[k] !== null) {
      // @ts-ignore
      acc += ' ' + getStringsFromObject(obj[k]);
    }
  }
  return acc;
}
