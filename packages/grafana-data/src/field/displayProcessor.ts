// Libraries
import _ from 'lodash';

// Types
import { Field, FieldType } from '../types/dataFrame';
import { GrafanaTheme } from '../types/theme';
import { DecimalCount, DecimalInfo, DisplayProcessor, DisplayValue } from '../types/displayValue';
import { getValueFormat } from '../valueFormats/valueFormats';
import { getMappedValue } from '../utils/valueMappings';
import { dateTime } from '../datetime';
import { KeyValue, TimeZone } from '../types';
import { getScaleCalculator } from './scale';
import { getTestTheme } from '../utils/testdata/testTheme';

interface DisplayProcessorOptions {
  field: Partial<Field>;
  /**
   * Will pick browser timezone if not defined
   */
  timeZone?: TimeZone;
  /**
   * Will pick 'dark' if not defined
   */
  theme?: GrafanaTheme;
  /**
   * Used by auto decimals logic
   **/
  tickSize?: number;
}

// Reasonable units for time
const timeFormats: KeyValue<boolean> = {
  dateTimeAsIso: true,
  dateTimeAsIsoSmart: true,
  dateTimeAsUS: true,
  dateTimeAsUSSmart: true,
  dateTimeFromNow: true,
};

export function getDisplayProcessor(options?: DisplayProcessorOptions): DisplayProcessor {
  if (!options || _.isEmpty(options) || !options.field) {
    return toStringProcessor;
  }

  const { field } = options;
  const config = field.config ?? {};

  // Theme should be required or we need access to default theme instance from here
  const theme = options.theme ?? getTestTheme();

  let unit = config.unit;
  let hasDateUnit = unit && (timeFormats[unit] || unit.startsWith('time:'));

  if (field.type === FieldType.time && !hasDateUnit) {
    unit = `dateTimeAsSystem`;
    hasDateUnit = true;
  }

  const formatFunc = getValueFormat(unit || 'none');
  const scaleFunc = getScaleCalculator(field as Field, theme);

  return (value: any) => {
    const { mappings } = config;
    const isStringUnit = unit === 'string';

    if (hasDateUnit && typeof value === 'string') {
      value = dateTime(value).valueOf();
    }

    let text = _.toString(value);
    let numeric = isStringUnit ? NaN : toNumber(value);
    let prefix: string | undefined = undefined;
    let suffix: string | undefined = undefined;
    let shouldFormat = true;

    if (mappings && mappings.length > 0) {
      const mappedValue = getMappedValue(mappings, value);

      if (mappedValue) {
        text = mappedValue.text;
        const v = isStringUnit ? NaN : toNumber(text);

        if (!isNaN(v)) {
          numeric = v;
        }

        shouldFormat = false;
      }
    }

    if (!isNaN(numeric)) {
      if (shouldFormat && !_.isBoolean(value)) {
        const { decimals, scaledDecimals } = getDecimalsForValue(value, config.decimals);
        const v = formatFunc(numeric, decimals, scaledDecimals, options.timeZone);
        text = v.text;
        suffix = v.suffix;
        prefix = v.prefix;

        // Check if the formatted text mapped to a different value
        if (mappings && mappings.length > 0) {
          const mappedValue = getMappedValue(mappings, text);
          if (mappedValue) {
            text = mappedValue.text;
          }
        }
      }

      // Return the value along with scale info
      if (text) {
        return { text, numeric, prefix, suffix, ...scaleFunc(numeric) };
      }
    }

    if (!text) {
      if (config.noValue) {
        text = config.noValue;
      } else {
        text = ''; // No data?
      }
    }

    return { text, numeric, prefix, suffix, ...scaleFunc(-Infinity) };
  };
}

/** Will return any value as a number or NaN */
function toNumber(value: any): number {
  if (typeof value === 'number') {
    return value;
  }
  if (value === '' || value === null || value === undefined || Array.isArray(value)) {
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

function getSignificantDigitCount(n: number) {
  //remove decimal and make positive
  n = Math.abs(+String(n).replace('.', ''));
  if (n === 0) {
    return 0;
  }

  // kill the 0s at the end of n
  while (n !== 0 && n % 10 === 0) {
    n /= 10;
  }

  // get number of digits
  return Math.floor(Math.log(n) / Math.LN10) + 1;
}

export function getDecimalsForValue(value: number, decimalOverride?: DecimalCount): DecimalInfo {
  if (_.isNumber(decimalOverride)) {
    // It's important that scaledDecimals is null here
    return { decimals: decimalOverride, scaledDecimals: null };
  }

  if (value === 0) {
    return { decimals: 0, scaledDecimals: 0 };
  }

  const digits = getSignificantDigitCount(value);
  const log10 = Math.floor(Math.log(Math.abs(value)) / Math.LN10);
  let dec = -log10 + 1;
  const magn = Math.pow(10, -dec);
  const norm = value / magn; // norm is between 1.0 and 10.0

  // special case for 2.5, requires an extra decimal
  if (norm > 2.25) {
    ++dec;
  }

  if (value % 1 === 0) {
    dec = 0;
  }

  const decimals = Math.max(0, dec);
  const scaledDecimals = decimals - log10 + digits - 1;

  return { decimals, scaledDecimals };
}

export function getRawDisplayProcessor(): DisplayProcessor {
  return (value: any) => ({
    text: `${value}`,
    numeric: (null as unknown) as number,
  });
}
