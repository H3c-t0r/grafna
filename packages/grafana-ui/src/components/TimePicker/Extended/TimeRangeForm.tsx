import React, { PureComponent, FormEvent, ReactNode } from 'react';
import Forms from '../../Forms';
import { TIME_FORMAT, TimeZone, isDateTime, TimeRange, DateTime } from '@grafana/data';
import { stringToDateTimeType, isValidTimeString } from '../time';
import TimePickerCalendar from './Calendar';

type CalendarTrigger = 'onFocus' | 'onButton';

interface Props {
  className?: string;
  calendarTrigger?: CalendarTrigger;
  value: TimeRange;
  onApply: (range: TimeRange) => void;
}
interface State {
  displayCalendar: boolean;
  from: { value: string; invalid: boolean };
  to: { value: string; invalid: boolean };
}

export default class TimeRangeForm extends PureComponent<Props, State> {
  constructor(props: Props, context?: any) {
    super(props, context);

    this.state = {
      from: {
        value: valueAsString(props.value.raw.from),
        invalid: false,
      },
      to: {
        value: valueAsString(props.value.raw.to),
        invalid: false,
      },
      displayCalendar: false,
    };
  }

  onInputFocus = (event: FormEvent<HTMLInputElement>) => {
    this.setState({ displayCalendar: true });
  };

  onInputChange = (event: FormEvent<HTMLInputElement>, name: keyof State) => {
    const { value } = event.currentTarget;
    const invalid = !isValid(value);

    this.setState({
      ...this.state,
      [name]: { value, invalid },
    });
  };

  onCalendarChange = (timeRange: TimeRange) => {
    this.setState({
      from: {
        value: valueAsString(timeRange.raw.from),
        invalid: false,
      },
      to: {
        value: valueAsString(timeRange.raw.to),
        invalid: false,
      },
    });
  };

  onApply = () => {
    const timeRange = toTimeRange(this.state);
    this.props.onApply(timeRange);
  };

  render() {
    const { from, to, displayCalendar } = this.state;
    const { calendarTrigger = 'onFocus' } = this.props;

    return (
      <>
        <TimePickerCalendar
          visible={displayCalendar}
          selected={toTimeRange(this.state)}
          onChange={this.onCalendarChange}
        />
        <Forms.Field label="From">
          <Forms.Input
            onFocus={this.onInputFocus}
            onChange={event => this.onInputChange(event, 'from')}
            addonAfter={calendarIcon(calendarTrigger)}
            {...from}
          />
        </Forms.Field>
        <Forms.Field label="To">
          <Forms.Input
            onFocus={this.onInputFocus}
            onChange={event => this.onInputChange(event, 'to')}
            addonAfter={calendarIcon(calendarTrigger)}
            {...to}
          />
        </Forms.Field>
        <Forms.Button onClick={this.onApply}>Apply time interval</Forms.Button>
      </>
    );
  }
}

function calendarIcon(calendarTrigger: CalendarTrigger): ReactNode {
  if (calendarTrigger === 'onButton') {
    return <Forms.Button icon="fa fa-calendar" variant="secondary" />;
  }
  return null;
}

function toTimeRange(state: State): TimeRange {
  const { from, to } = state;

  return {
    from: stringToDateTimeType(from.value),
    to: stringToDateTimeType(to.value),
    raw: {
      from: from.value,
      to: to.value,
    },
  };
}

function valueAsString(value: DateTime | string): string {
  if (isDateTime(value)) {
    return value.format(TIME_FORMAT);
  }
  return value;
}

function isValid(value: string, roundup?: boolean, timeZone?: TimeZone): boolean {
  if (value.indexOf('now') !== -1) {
    return isValidTimeString(value);
  }

  const parsed = stringToDateTimeType(value, roundup, timeZone);
  return parsed.isValid();
}
