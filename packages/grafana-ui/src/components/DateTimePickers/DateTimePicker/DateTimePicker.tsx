import React, { FC, FormEvent, ReactNode, useCallback, useState } from 'react';
import { useMedia } from 'react-use';
import Calendar from 'react-calendar/dist/entry.nostyle';
import { css } from '@emotion/css';
import { dateTimeFormat, DateTime, dateTime, GrafanaTheme2 } from '@grafana/data';
import { Button, Field, Icon, Input, Portal } from '../..';
import { ClickOutsideWrapper } from '../../ClickOutsideWrapper/ClickOutsideWrapper';
import { isValid } from '../utils';
import { getBodyStyles, getStyles as getCalendarStyles } from '../TimeRangePicker/TimePickerCalendar';
import { useStyles2, useTheme2 } from '../../../themes';
import { TimeOfDayPicker } from '../TimeOfDayPicker';

export interface Props {
  label: ReactNode;
  date: DateTime;
  onChange: (date: DateTime) => void;
}

const stopPropagation = (event: React.MouseEvent<HTMLDivElement>) => event.stopPropagation();

export const DateTimePicker: FC<Props> = ({ date, label, onChange }) => {
  const [isOpen, setOpen] = useState(false);
  const [internalDate, setInternalDate] = useState<InputState>(() => {
    return { value: dateTimeFormat(date), invalid: false };
  });
  const theme = useTheme2();
  const isFullscreen = useMedia(`(min-width: ${theme.breakpoints.values.lg}px)`);
  const containerStyles = useStyles2(getCalendarStyles);

  const onApply = useCallback(
    (date: DateTime) => {
      setOpen(false);
      onChange(date);
    },
    [onChange]
  );

  const onOpen = useCallback(
    (event: FormEvent<HTMLElement>) => {
      event.preventDefault();
      setOpen(true);
    },
    [setOpen]
  );

  const onFocus = useCallback(
    (event: FormEvent<HTMLElement>) => {
      if (!isFullscreen) {
        return;
      }
      onOpen(event);
    },
    [isFullscreen, onOpen]
  );

  const onChangeDate = useCallback((event: FormEvent<HTMLInputElement>) => {
    const isInvalid = !isValid(event.currentTarget.value);
    setInternalDate({
      value: event.currentTarget.value,
      invalid: isInvalid,
    });
  }, []);

  const icon = <Button icon="calendar-alt" variant="secondary" onClick={onOpen} />;

  return (
    <div>
      <Field label={label} onClick={stopPropagation} invalid={internalDate.invalid} error="Incorrect date format">
        <Input
          onClick={stopPropagation}
          onChange={onChangeDate}
          addonAfter={icon}
          value={internalDate.value}
          onFocus={onFocus}
          onBlur={() => onChange(dateTime(internalDate.value))}
        />
      </Field>
      {isOpen ? (
        isFullscreen ? (
          <ClickOutsideWrapper onClick={() => setOpen(false)}>
            <DateTimeCalendar date={date} onChange={onApply} />
          </ClickOutsideWrapper>
        ) : (
          <Portal>
            <div className={containerStyles.modal} onClick={stopPropagation}>
              <DateTimeCalendar date={date} onChange={onApply} />
            </div>
            <div className={containerStyles.backdrop} onClick={stopPropagation} />
          </Portal>
        )
      ) : null}
    </div>
  );
};

interface DateTimeCalendarProps {
  date: DateTime;
  onChange: (date: DateTime) => void;
}

interface InputProps {
  label?: ReactNode;
  date: DateTime;
  onChange: (date: DateTime) => void;
  onFocus: (event: FormEvent<HTMLElement>) => void;
  onOpen: (event: FormEvent<HTMLElement>) => void;
}

type InputState = {
  value: string;
  invalid: boolean;
};

const DateTimeInput: FC<InputProps> = ({ date, label, onChange, onFocus, onOpen }) => {
  console.log('DateTimeInput', date);

  return (
    <Field label={label} onClick={stopPropagation} invalid={internalDate.invalid} error="Incorrect date format">
      <Input
        onClick={stopPropagation}
        onChange={onChangeDate}
        addonAfter={icon}
        value={internalDate.value}
        onFocus={onFocus}
        onBlur={() => onChange(dateTime(internalDate.value))}
      />
    </Field>
  );
};

const DateTimeCalendar: FC<DateTimeCalendarProps> = ({ date, onChange }) => {
  const calendarStyles = useStyles2(getBodyStyles);
  const styles = useStyles2(getStyles);
  const [internalDate, setInternalDate] = useState<Date>(date.toDate() || Date.now());

  const onChangeDate = useCallback((date: Date | Date[]) => {
    if (!Array.isArray(date)) {
      setInternalDate(date);
    }
  }, []);

  const onChangeTime = useCallback((date: DateTime) => {
    setInternalDate(date.toDate());
  }, []);

  return (
    <div className={styles.container} onClick={stopPropagation}>
      <Calendar
        next2Label={null}
        prev2Label={null}
        value={internalDate}
        nextLabel={<Icon name="angle-right" />}
        prevLabel={<Icon name="angle-left" />}
        onChange={onChangeDate}
        locale="en"
        className={calendarStyles.body}
        tileClassName={calendarStyles.title}
      />
      <div className={styles.time}>
        <TimeOfDayPicker showSeconds={true} onChange={onChangeTime} value={dateTime(internalDate)} />
      </div>
      <div>
        <Button type="button" onClick={() => onChange(dateTime(internalDate))}>
          Apply
        </Button>
      </div>
    </div>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  container: css`
    position: absolute;
    right: 510px;
    bottom: 435px;
    padding: ${theme.spacing(1)};
    border: 1px ${theme.colors.border.weak} solid;
    border-radius: ${theme.shape.borderRadius(1)};
  `,
  time: css`
    margin-bottom: ${theme.spacing(2)};
  `,
});
