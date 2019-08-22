// Libraries
import _ from 'lodash';
import {
  Threshold,
  getMappedValue,
  FieldConfig,
  DisplayProcessor,
  DecimalInfo,
  DisplayValue,
  DecimalCount,
} from '@grafana/data';

// Utils
import { getValueFormat } from './valueFormats/valueFormats';
import { getColorFromHexRgbOrName } from './namedColorsPalette';

// Types
import { GrafanaTheme, GrafanaThemeType } from '../types/index';

interface DisplayProcessorOptions {
  field?: FieldConfig;

  // Context
  isUtc?: boolean;
  theme?: GrafanaTheme; // Will pick 'dark' if not defined
}

export function getDisplayProcessor(options?: DisplayProcessorOptions): DisplayProcessor {
  if (options && !_.isEmpty(options)) {
    const field = options.field ? options.field : {};
    const formatFunc = getValueFormat(field.unit || 'none');

    return (value: any) => {
      const { theme } = options;
      const { mappings, thresholds } = field;
      let color;

      let text = _.toString(value);
      let numeric = toNumber(value);

      let shouldFormat = true;
      if (mappings && mappings.length > 0) {
        const mappedValue = getMappedValue(mappings, value);

        if (mappedValue) {
          text = mappedValue.text;
          const v = toNumber(text);

          if (!isNaN(v)) {
            numeric = v;
          }

          shouldFormat = false;
        }
      }

      if (!isNaN(numeric)) {
        if (shouldFormat && !_.isBoolean(value)) {
          const { decimals, scaledDecimals } = getDecimalsForValue(value, field.decimals);
          text = formatFunc(numeric, decimals, scaledDecimals, options.isUtc);
        }
        if (thresholds && thresholds.length) {
          color = getColorFromThreshold(numeric, thresholds, theme);
        }
      }

      if (!text) {
        if (field && field.noValue) {
          text = field.noValue;
        } else {
          text = ''; // No data?
        }
      }
      return { text, numeric, color };
    };
  }

  return toStringProcessor;
}

/** Will return any value as a number or NaN */
function toNumber(value: any): number {
  if (typeof value === 'number') {
    return value;
  }
  if (value === null || value === undefined || Array.isArray(value)) {
    return NaN; // lodash calls them 0
  }
  if (typeof value === 'boolean') {
    return value ? 1 : 0;
  }
  return _.toNumber(value);
}

function toStringProcessor(value: any): DisplayValue {
  return { text: _.toString(value), numeric: toNumber(value) };
}

export function getColorFromThreshold(value: number, thresholds: Threshold[], theme?: GrafanaTheme): string {
  const themeType = theme ? theme.type : GrafanaThemeType.Dark;

  if (thresholds.length === 1) {
    return getColorFromHexRgbOrName(thresholds[0].color, themeType);
  }

  const atThreshold = thresholds.filter(threshold => value === threshold.value)[0];
  if (atThreshold) {
    return getColorFromHexRgbOrName(atThreshold.color, themeType);
  }

  const belowThreshold = thresholds.filter(threshold => value > threshold.value);

  if (belowThreshold.length > 0) {
    const nearestThreshold = belowThreshold.sort((t1, t2) => t2.value - t1.value)[0];
    return getColorFromHexRgbOrName(nearestThreshold.color, themeType);
  }

  // Use the first threshold as the default color
  return getColorFromHexRgbOrName(thresholds[0].color, themeType);
}

// function getSignificantDigitCount(n: number): number {
//   // remove decimal and make positive
//   n = Math.abs(parseInt(String(n).replace('.', ''), 10));
//   if (n === 0) {
//     return 0;
//   }
//
//   // kill the 0s at the end of n
//   while (n !== 0 && n % 10 === 0) {
//     n /= 10;
//   }
//
//   // get number of digits
//   return Math.floor(Math.log(n) / Math.LN10) + 1;
// }

export function getDecimalsForValue(value: number, decimalOverride?: DecimalCount): DecimalInfo {
  if (_.isNumber(decimalOverride)) {
    // It's important that scaledDecimals is null here
    return { decimals: decimalOverride, scaledDecimals: null };
  }

  let dec = -Math.floor(Math.log(value) / Math.LN10) + 1;
  const magn = Math.pow(10, -dec);
  const norm = value / magn; // norm is between 1.0 and 10.0
  let size;

  if (norm < 1.5) {
    size = 1;
  } else if (norm < 3) {
    size = 2;
    // special case for 2.5, requires an extra decimal
    if (norm > 2.25) {
      size = 2.5;
      ++dec;
    }
  } else if (norm < 7.5) {
    size = 5;
  } else {
    size = 10;
  }

  size *= magn;

  // reduce starting decimals if not needed
  if (value % 1 === 0) {
    dec = 0;
  }

  const decimals = Math.max(0, dec);
  const scaledDecimals = decimals - Math.floor(Math.log(size) / Math.LN10) + 2;

  return { decimals, scaledDecimals };
}
