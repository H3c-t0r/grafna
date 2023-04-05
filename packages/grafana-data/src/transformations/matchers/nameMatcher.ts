import { getFieldDisplayName } from '../../field/fieldState';
import { stringToJsRegex } from '../../text/string';
import { Field, DataFrame, FieldType, TIME_SERIES_VALUE_FIELD_NAME } from '../../types/dataFrame';
import { FieldMatcherInfo, FrameMatcherInfo, FieldMatcher } from '../../types/transformations';

import { FieldMatcherID, FrameMatcherID } from './ids';
export interface RegexpOrNamesMatcherOptions {
  pattern?: string;
  names?: string[];
}

/**
 * Mode to be able to toggle if the names matcher should match fields in provided
 * list or all except provided names.
 * @public
 */
export enum ByNamesMatcherMode {
  exclude = 'exclude',
  include = 'include',
}

/**
 * Options to instruct the by names matcher to either match all fields in given list
 * or all except the fields in the list.
 * @public
 */
export interface ByNamesMatcherOptions {
  mode?: ByNamesMatcherMode;
  names?: string[];
  readOnly?: boolean;
  prefix?: string;
}

// General Field matcher
const fieldNameMatcher: FieldMatcherInfo<string> = {
  id: FieldMatcherID.byName,
  name: 'Field Name',
  description: 'match the field name',
  defaultOptions: '',

  get: (name: string): FieldMatcher => {
    // how to get the feature toggle in grafana-data because we cannot import config from grafana-runtime
    const fieldNameMatcherFallback = window.grafanaBootData.settings.featureToggles.fieldNameMatcherFallback;
    // wrapped in a feature toggle
    // so we may turn this off, with prometheusDataplane toggle, if we have unexpected behavior/errors
    if (fieldNameMatcherFallback) {
      let fallback: FieldMatcher | undefined = undefined;

      // In an effor to support migrating to a consistent data contract, the
      // naming conventions need to get normalized.  However many existing setups
      // exist that would no longer match names if that changes.  This injects
      // fallback logic when when the data frame has not type version specified
      if (name === TIME_SERIES_VALUE_FIELD_NAME) {
        fallback = (field: Field, frame: DataFrame) => {
          return (
            Boolean(field.labels) && // Value was reasonable when the name was set in labels or on the frame
            field.labels?.__name__ === field.name
          );
        };
      } else if (name === 'Time' || name === 'time') {
        fallback = (field: Field, frame: DataFrame) => {
          return frame.meta?.typeVersion == null && field.type === FieldType.time;
        };
      }

      return (field: Field, frame: DataFrame, allFrames: DataFrame[]) => {
        return (
          name === field.name ||
          name === getFieldDisplayName(field, frame, allFrames) ||
          Boolean(fallback && fallback(field, frame, allFrames))
        );
      };
    } else {
      return (field: Field, frame: DataFrame, allFrames: DataFrame[]) => {
        return name === field.name || getFieldDisplayName(field, frame, allFrames) === name;
      };
    }
  },

  getOptionsDisplayText: (name: string) => {
    return `Field name: ${name}`;
  },
};

const multipleFieldNamesMatcher: FieldMatcherInfo<ByNamesMatcherOptions> = {
  id: FieldMatcherID.byNames,
  name: 'Field Names',
  description: 'match any of the given the field names',
  defaultOptions: {
    mode: ByNamesMatcherMode.include,
    names: [],
  },

  get: (options: ByNamesMatcherOptions): FieldMatcher => {
    const { names, mode = ByNamesMatcherMode.include } = options;
    const uniqueNames = new Set<string>(names ?? []);

    const matcher = (field: Field, frame: DataFrame, frames: DataFrame[]) => {
      return uniqueNames.has(field.name) || uniqueNames.has(getFieldDisplayName(field, frame, frames));
    };

    if (mode === ByNamesMatcherMode.exclude) {
      return (field: Field, frame: DataFrame, frames: DataFrame[]) => {
        return !matcher(field, frame, frames);
      };
    }
    return matcher;
  },

  getOptionsDisplayText: (options: ByNamesMatcherOptions): string => {
    const { names, mode } = options;
    const displayText = (names ?? []).join(', ');
    if (mode === ByNamesMatcherMode.exclude) {
      return `All except: ${displayText}`;
    }
    return `All of: ${displayText}`;
  },
};

const regexpFieldNameMatcher: FieldMatcherInfo<string> = {
  id: FieldMatcherID.byRegexp,
  name: 'Field Name by Regexp',
  description: 'match the field name by a given regexp pattern',
  defaultOptions: '/.*/',

  get: (pattern: string): FieldMatcher => {
    const regexp = patternToRegex(pattern);

    return (field: Field, frame: DataFrame, allFrames: DataFrame[]) => {
      const displayName = getFieldDisplayName(field, frame, allFrames);
      return !!regexp && regexp.test(displayName);
    };
  },

  getOptionsDisplayText: (pattern: string): string => {
    return `Field name by pattern: ${pattern}`;
  },
};

/**
 * Field matcher that will match all fields that exists in a
 * data frame with configured refId.
 * @public
 */
const fieldsInFrameMatcher: FieldMatcherInfo<string> = {
  id: FieldMatcherID.byFrameRefID,
  name: 'Fields by frame refId',
  description: 'match all fields returned in data frame with refId.',
  defaultOptions: '',

  get: (refId: string): FieldMatcher => {
    return (field: Field, frame: DataFrame, allFrames: DataFrame[]) => {
      return frame.refId === refId;
    };
  },

  getOptionsDisplayText: (refId: string): string => {
    return `Math all fields returned by query with reference ID: ${refId}`;
  },
};

const regexpOrMultipleNamesMatcher: FieldMatcherInfo<RegexpOrNamesMatcherOptions> = {
  id: FieldMatcherID.byRegexpOrNames,
  name: 'Field Name by Regexp or Names',
  description: 'match the field name by a given regexp pattern or given names',
  defaultOptions: {
    pattern: '/.*/',
    names: [],
  },

  get: (options: RegexpOrNamesMatcherOptions): FieldMatcher => {
    const regexpMatcher = regexpFieldNameMatcher.get(options?.pattern || '');
    const namesMatcher = multipleFieldNamesMatcher.get({
      mode: ByNamesMatcherMode.include,
      names: options?.names ?? [],
    });

    return (field: Field, frame: DataFrame, allFrames: DataFrame[]) => {
      return namesMatcher(field, frame, allFrames) || regexpMatcher(field, frame, allFrames);
    };
  },

  getOptionsDisplayText: (options: RegexpOrNamesMatcherOptions): string => {
    const pattern = options?.pattern ?? '';
    const names = options?.names?.join(',') ?? '';
    return `Field name by pattern: ${pattern} or names: ${names}`;
  },
};

const patternToRegex = (pattern?: string): RegExp | undefined => {
  if (!pattern) {
    return undefined;
  }

  try {
    return stringToJsRegex(pattern);
  } catch (error) {
    console.error(error);
    return undefined;
  }
};

// General Frame matcher
const frameNameMatcher: FrameMatcherInfo<string> = {
  id: FrameMatcherID.byName,
  name: 'Frame Name',
  description: 'match the frame name',
  defaultOptions: '/.*/',

  get: (pattern: string) => {
    const regex = stringToJsRegex(pattern);
    return (frame: DataFrame) => {
      return regex.test(frame.name || '');
    };
  },

  getOptionsDisplayText: (pattern: string) => {
    return `Frame name: ${pattern}`;
  },
};

/**
 * Registry Initialization
 */
export function getFieldNameMatchers(): FieldMatcherInfo[] {
  return [
    fieldNameMatcher,
    regexpFieldNameMatcher,
    multipleFieldNamesMatcher,
    regexpOrMultipleNamesMatcher,
    fieldsInFrameMatcher,
  ];
}

export function getFrameNameMatchers(): FrameMatcherInfo[] {
  return [frameNameMatcher];
}
