import { css } from '@emotion/css';
import { GrafanaTheme } from '@grafana/data';

export const getStyles = ({ palette, colors, isLight }: GrafanaTheme) => {
  const borderColor = isLight ? palette.gray85 : palette.dark7;

  return {
    panel: css`
      display: flex;
      flex-direction: column;
      height: 100%;
    `,
    spinner: css`
      display: flex;
      height: 10em;
      align-items: center;
      justify-content: center;
    `,
    empty: css`
      display: flex;
      width: 100%;
      height: 160px;
      justify-content: center;
      align-items: center;
      border: 1px solid ${borderColor};
      white-space: pre-wrap;
    `,
    link: css`
      color: ${colors.linkExternal};
      &:hover {
        color: ${colors.textBlue};
      }
    `,
  };
};
