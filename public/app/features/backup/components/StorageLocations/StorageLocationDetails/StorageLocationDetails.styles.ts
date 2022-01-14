import { css } from 'emotion';

import { GrafanaTheme } from '@grafana/data';

export const getStyles = ({ spacing }: GrafanaTheme) => ({
  wrapper: css`
    display: flex;
    align-items: center;

    & > * {
      flex: 1 1 calc(100% / 3);

      &:not(:last-child) {
        padding-right: ${spacing.md};
      }

      &:not(:first-child) {
        padding-left: ${spacing.md};
      }
    }
  `,
});
