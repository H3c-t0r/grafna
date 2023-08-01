import { css, cx } from '@emotion/css';
import React, { HTMLAttributes } from 'react';

import { GrafanaTheme2 } from '@grafana/data';

import { stylesFactory, useTheme2 } from '../../themes';
import { getChildId } from '../../utils/reactUtils';
import { Icon } from '../Icon/Icon';
import { Tooltip } from '../Tooltip/Tooltip';

import { FieldValidationMessage } from './FieldValidationMessage';
import { Label } from './Label';

export interface FieldProps extends HTMLAttributes<HTMLDivElement> {
  /** Form input element, i.e Input or Switch */
  children: React.ReactElement;
  /** Label for the field */
  label?: React.ReactNode;
  /** Description of the field */
  description?: React.ReactNode;
  /** Indicates if field is in invalid state */
  invalid?: boolean;
  /** Indicates if field is in loading state */
  loading?: boolean;
  /** Indicates if field is disabled */
  disabled?: boolean;
  /** Indicates if field is required */
  required?: boolean;
  /** Error message to display */
  error?: React.ReactNode;
  /** Indicates horizontal layout of the field */
  horizontal?: boolean;
  /** make validation message overflow horizontally. Prevents pushing out adjacent inline components */
  validationMessageHorizontalOverflow?: boolean;
  /** Adds tooltip next to the label */
  tooltipMessage?: string;

  className?: string;
  /**
   *  A unique id that associates the label of the Field component with the control with the unique id.
   *  If the `htmlFor` property is missing the `htmlFor` will be inferred from the `id` or `inputId` property of the first child.
   *  https://developer.mozilla.org/en-US/docs/Web/HTML/Element/label#attr-for
   */
  htmlFor?: string;
}

export const getFieldStyles = stylesFactory((theme: GrafanaTheme2) => {
  return {
    field: css({
      display: 'flex',
      flexDirection: 'column',
      marginBottom: theme.spacing(2),
    }),
    fieldHorizontal: css({
      flexDirection: 'row',
      justifyContent: 'space-between',
      flexWrap: 'wrap',
    }),
    fieldValidationWrapper: css({
      marginTop: theme.spacing(0.5),
    }),
    fieldValidationWrapperHorizontal: css({
      flex: '1 1 100%',
    }),
    validationMessageHorizontalOverflow: css({
      width: 0,
      overflowX: 'visible',

      '& > *': {
        whiteSpace: 'nowrap',
      },
    }),
    tooltipIcon: css({
      marginLeft: theme.spacing(0.5),
    }),
  };
});

export const Field = React.forwardRef<HTMLDivElement, FieldProps>(
  (
    {
      label,
      description,
      horizontal,
      invalid,
      loading,
      disabled,
      required,
      error,
      children,
      className,
      validationMessageHorizontalOverflow,
      htmlFor,
      tooltipMessage,
      ...otherProps
    }: FieldProps,
    ref
  ) => {
    const theme = useTheme2();
    const styles = getFieldStyles(theme);
    const inputId = htmlFor ?? getChildId(children);

    const labelElement =
      typeof label === 'string' ? (
        <Label htmlFor={inputId} description={description}>
          {`${label}${required ? ' *' : ''}`}
          {tooltipMessage && (
            <Tooltip placement="right-end" content={tooltipMessage}>
              <Icon name="question-circle" size="xs" className={styles.tooltipIcon} />
            </Tooltip>
          )}
        </Label>
      ) : (
        label
      );

    const childProps = deleteUndefinedProps({ invalid, disabled, loading });
    return (
      <div className={cx(styles.field, horizontal && styles.fieldHorizontal, className)} {...otherProps}>
        {labelElement}
        <div>
          <div ref={ref}>{React.cloneElement(children, childProps)}</div>
          {invalid && error && !horizontal && (
            <div
              className={cx(styles.fieldValidationWrapper, {
                [styles.validationMessageHorizontalOverflow]: !!validationMessageHorizontalOverflow,
              })}
            >
              <FieldValidationMessage>{error}</FieldValidationMessage>
            </div>
          )}
        </div>

        {invalid && error && horizontal && (
          <div
            className={cx(styles.fieldValidationWrapper, styles.fieldValidationWrapperHorizontal, {
              [styles.validationMessageHorizontalOverflow]: !!validationMessageHorizontalOverflow,
            })}
          >
            <FieldValidationMessage>{error}</FieldValidationMessage>
          </div>
        )}
      </div>
    );
  }
);

Field.displayName = 'Field';

function deleteUndefinedProps<T extends Object>(obj: T): Partial<T> {
  for (const key in obj) {
    if (obj[key] === undefined) {
      delete obj[key];
    }
  }

  return obj;
}
