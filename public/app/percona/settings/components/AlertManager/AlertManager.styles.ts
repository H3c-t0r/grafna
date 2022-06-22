import { css } from '@emotion/css';

import { GrafanaTheme } from '@grafana/data';
import { stylesFactory } from '@grafana/ui';

export const getStyles = stylesFactory((theme: GrafanaTheme) => ({
  alertManagerWrapper: css`
    display: flex;
    flex-direction: column;
    width: 600px;
  `,
  textarea: css`
    margin: ${theme.spacing.md} 0;
    min-height: 150px;
  `,
  input: css`
    margin: ${theme.spacing.md} 0;
  `,
  rulesLabel: css`
    margin-top: ${theme.spacing.sm};
  `,
  warning: css`
    margin-bottom: ${theme.spacing.md};
  `,
  warningLink: css`
    text-decoration: underline;
  `,
}));
