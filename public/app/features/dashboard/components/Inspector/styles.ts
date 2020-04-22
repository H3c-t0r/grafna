import { css } from 'emotion';
import { config } from 'app/core/config';
import { stylesFactory } from '@grafana/ui';

export const getPanelInspectorStyles = stylesFactory(() => {
  return {
    wrap: css`
      display: flex;
      flex-direction: column;
      height: 100%;
      width: 100%;
      flex: 1 1 0;
    `,
    toolbar: css`
      display: flex;
      width: 100%;
      flex-grow: 0;
      align-items: center;
      justify-content: flex-end;
    `,
    toolbarItem: css`
      margin-left: ${config.theme.spacing.md};
    `,
    content: css`
      flex-grow: 1;
      padding-bottom: 16px;
    `,
    contentQueryInspector: css`
      flex-grow: 1;
      padding: ${config.theme.spacing.md} 0;
    `,
    editor: css`
      font-family: monospace;
      height: 100%;
      flex-grow: 1;
    `,
    viewer: css`
      overflow: scroll;
    `,
    dataFrameSelect: css`
      flex-grow: 2;
    `,
    downloadCsv: css`
      margin-left: 16px;
    `,
    tabContent: css`
      height: 100%;
      display: flex;
      flex-direction: column;
    `,
    dataTabContent: css`
      display: flex;
      flex-direction: column;
      height: 100%;
      width: 100%;
    `,
  };
});
