import { cx, css } from '@emotion/css';
import React, { ReactElement, useEffect, useRef, useState } from 'react';

import { GrafanaTheme2 } from '@grafana/data';

import { useStyles2 } from '../../themes';
import { ComponentSize } from '../../types/size';
import { Button, ButtonVariant } from '../Button';
import { Stack } from '../Layout/Stack/Stack';

export interface Props {
  /** Confirm action callback */
  onConfirm(): void;
  children: string | ReactElement;
  /** Custom button styles */
  className?: string;
  /** Button size */
  size?: ComponentSize;
  /** Text for the Confirm button */
  confirmText?: string;
  /** Disable button click action */
  disabled?: boolean;
  /** Variant of the Confirm button */
  confirmVariant?: ButtonVariant;
  /** Hide confirm actions when after of them is clicked */
  closeOnConfirm?: boolean;
  /** Optional on click handler for the original button */
  onClick?(): void;
  /** Callback for the cancel action */
  onCancel?(): void;
}

export const ConfirmButton = ({
  children,
  className,
  closeOnConfirm,
  confirmText = 'Save',
  confirmVariant = 'primary',
  disabled = false,
  onCancel,
  onClick,
  onConfirm,
  size = 'md',
}: Props) => {
  const mainButtonRef = useRef<HTMLButtonElement>(null);
  const confirmButtonRef = useRef<HTMLButtonElement>(null);
  const [showConfirm, setShowConfirm] = useState(false);
  const [shouldRestoreFocus, setShouldRestoreFocus] = useState(false);
  const styles = useStyles2(getStyles);

  useEffect(() => {
    if (showConfirm) {
      confirmButtonRef.current?.focus();
      setShouldRestoreFocus(true);
    } else {
      if (shouldRestoreFocus) {
        mainButtonRef.current?.focus();
        setShouldRestoreFocus(false);
      }
    }
  }, [shouldRestoreFocus, showConfirm]);

  const onClickButton = (event: React.MouseEvent<HTMLButtonElement>) => {
    if (event) {
      event.preventDefault();
    }

    setShowConfirm(true);
    onClick?.();
  };

  const onClickCancel = (event: React.MouseEvent<HTMLButtonElement>) => {
    if (event) {
      event.preventDefault();
    }
    setShowConfirm(false);
    mainButtonRef.current?.focus();
    onCancel?.();
  };

  const onClickConfirm = (event: React.MouseEvent<HTMLButtonElement>) => {
    if (event) {
      event.preventDefault();
    }
    onConfirm?.();
    if (closeOnConfirm) {
      setShowConfirm(false);
    }
  };

  const buttonClass = cx(className, styles.mainButton, {
    [styles.mainButtonHide]: showConfirm,
  });
  const confirmButtonClass = cx(styles.confirmButton, {
    [styles.confirmButtonHide]: !showConfirm,
  });

  return (
    <Stack gap={0} alignItems="center" justifyContent="flex-end">
      <span className={buttonClass}>
        {typeof children === 'string' ? (
          <Button disabled={disabled} size={size} fill="text" onClick={onClickButton} ref={mainButtonRef}>
            {children}
          </Button>
        ) : (
          React.cloneElement(children, { disabled, onClick: onClickButton, ref: mainButtonRef })
        )}
      </span>
      <span className={confirmButtonClass}>
        <Button size={size} variant={confirmVariant} onClick={onClickConfirm} ref={confirmButtonRef}>
          {confirmText}
        </Button>
        <Button size={size} fill="text" onClick={onClickCancel}>
          Cancel
        </Button>
      </span>
    </Stack>
  );
};
ConfirmButton.displayName = 'ConfirmButton';

const getStyles = (theme: GrafanaTheme2) => {
  return {
    mainButton: css({
      opacity: 1,
      transition: 'opacity 0.1s ease',
      zIndex: 2,
    }),
    mainButtonHide: css({
      opacity: 0,
      transition: 'opacity 0.1s ease, visibility 0 0.1s',
      visibility: 'hidden',
      zIndex: 0,
    }),
    confirmButton: css({
      alignItems: 'flex-start',
      background: theme.colors.background.primary,
      display: 'flex',
      opacity: 1,
      pointerEvents: 'all',
      position: 'absolute',
      transform: 'translateX(0)',
      transition: 'opacity 0.08s ease-out, transform 0.1s ease-out',
      zIndex: 1,
    }),
    confirmButtonHide: css({
      opacity: 0,
      pointerEvents: 'none',
      transform: 'translateX(100%)',
      transition: 'opacity 0.12s ease-in, transform 0.14s ease-in, visibility 0s 0.12s',
      visibility: 'hidden',
    }),
  };
};
