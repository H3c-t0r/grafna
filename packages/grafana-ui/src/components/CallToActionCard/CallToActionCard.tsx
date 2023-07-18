import { css, cx } from '@emotion/css';
import React from 'react';

import { GrafanaTheme2 } from '@grafana/data';

import { useStyles2 } from '../../themes/ThemeContext';

export interface CallToActionCardProps {
  message?: string | JSX.Element;
  callToActionElement: JSX.Element;
  footer?: string | JSX.Element;
  className?: string;
}

export const CallToActionCard = ({ message, callToActionElement, footer, className }: CallToActionCardProps) => {
  const css = useStyles2(getStyles);

  return (
    <div className={cx([css.wrapper, className])}>
      {message && <div className={css.message}>{message}</div>}
      {callToActionElement}
      {footer && <div className={css.footer}>{footer}</div>}
    </div>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  wrapper: css({
    label: 'call-to-action-card',
    padding: theme.spacing(3),
    background: theme.colors.background.secondary,
    borderRadius: theme.shape.radius.default,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    flexGrow: 1,
  }),
  message: css({
    marginBottom: theme.spacing(3),
    fontStyle: 'italic',
  }),
  footer: css({
    marginTop: theme.spacing(3),
  }),
});
