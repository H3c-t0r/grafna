import { css } from '@emotion/css';
import React from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { Components } from '@grafana/e2e-selectors';
import { ToolbarButton, useTheme2 } from '@grafana/ui';
import { t, Trans } from 'app/core/internationalization';

type Props = {
  addQueryRowButtonDisabled?: boolean;
  addQueryRowButtonHidden?: boolean;
  richHistoryRowButtonHidden?: boolean;
  richHistoryButtonActive?: boolean;
  queryInspectorButtonActive?: boolean;

  onClickAddQueryRowButton: () => void;
  onClickRichHistoryButton: () => void;
  onClickQueryInspectorButton: () => void;
};

const getStyles = (theme: GrafanaTheme2) => {
  return {
    containerMargin: css`
      display: flex;
      flex-wrap: wrap;
      gap: ${theme.spacing(1)};
      margin-top: ${theme.spacing(2)};
    `,
  };
};

export function SecondaryActions(props: Props) {
  const theme = useTheme2();
  const styles = getStyles(theme);
  return (
    <div className={styles.containerMargin}>
      {!props.addQueryRowButtonHidden && (
        <ToolbarButton
          variant="canvas"
          aria-label={t('explore.secondary-actions.query-add-button-aria-label', 'Add query')}
          onClick={props.onClickAddQueryRowButton}
          disabled={props.addQueryRowButtonDisabled}
          icon="plus"
        >
          <Trans i18nKey="explore.secondary-actions.query-add-button">Add query</Trans>
        </ToolbarButton>
      )}
      {!props.richHistoryRowButtonHidden && (
        <ToolbarButton
          variant={props.richHistoryButtonActive ? 'active' : 'canvas'}
          aria-label={t('explore.secondary-actions.query-history-button-aria-label', 'Query history')}
          onClick={props.onClickRichHistoryButton}
          data-testid={Components.QueryTab.queryHistoryButton}
          icon="history"
        >
          <Trans i18nKey="explore.secondary-actions.query-history-button">Query history</Trans>
        </ToolbarButton>
      )}
      <ToolbarButton
        variant={props.queryInspectorButtonActive ? 'active' : 'canvas'}
        aria-label={t('explore.secondary-actions.query-inspector-button-aria-label', 'Query inspector')}
        onClick={props.onClickQueryInspectorButton}
        icon="info-circle"
      >
        <Trans i18nKey="explore.secondary-actions.query-inspector-button">Query inspector</Trans>
      </ToolbarButton>
    </div>
  );
}
