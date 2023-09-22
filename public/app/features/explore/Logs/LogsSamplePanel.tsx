import { css } from '@emotion/css';
import React from 'react';

import {
  DataQueryResponse,
  DataSourceApi,
  GrafanaTheme2,
  hasSupplementaryQuerySupport,
  LoadingState,
  LogsDedupStrategy,
  SplitOpen,
  SupplementaryQueryType,
} from '@grafana/data';
import { config, reportInteraction } from '@grafana/runtime';
import { DataQuery, TimeZone } from '@grafana/schema';
import { Button, Collapse, Icon, Tooltip, useStyles2 } from '@grafana/ui';
import { t, Trans } from 'app/core/internationalization';
import store from 'app/core/store';

import { LogRows } from '../../logs/components/LogRows';
import { dataFrameToLogsModel } from '../../logs/logsModel';
import { SupplementaryResultError } from '../SupplementaryResultError';

import { SETTINGS_KEYS } from './utils/logs';

type Props = {
  queryResponse: DataQueryResponse | undefined;
  enabled: boolean;
  timeZone: TimeZone;
  queries: DataQuery[];
  datasourceInstance: DataSourceApi | null | undefined;
  splitOpen: SplitOpen;
  setLogsSampleEnabled: (enabled: boolean) => void;
};

export function LogsSamplePanel(props: Props) {
  const { queryResponse, timeZone, enabled, setLogsSampleEnabled, datasourceInstance, queries, splitOpen } = props;

  const styles = useStyles2(getStyles);
  const onToggleLogsSampleCollapse = (isOpen: boolean) => {
    setLogsSampleEnabled(isOpen);
    reportInteraction('grafana_explore_logs_sample_toggle_clicked', {
      datasourceType: datasourceInstance?.type ?? 'unknown',
      type: isOpen ? 'open' : 'close',
    });
  };

  const OpenInSplitViewButton = () => {
    if (!hasSupplementaryQuerySupport(datasourceInstance, SupplementaryQueryType.LogsSample)) {
      return null;
    }

    const logSampleQueries = queries
      .map((query) => datasourceInstance.getSupplementaryQuery({ type: SupplementaryQueryType.LogsSample }, query))
      .filter((query): query is DataQuery => !!query);

    if (!logSampleQueries.length) {
      return null;
    }

    const onSplitOpen = () => {
      splitOpen({ queries: logSampleQueries, datasourceUid: datasourceInstance.uid });
      reportInteraction('grafana_explore_logs_sample_split_button_clicked', {
        datasourceType: datasourceInstance?.type ?? 'unknown',
        queriesCount: logSampleQueries.length,
      });
    };

    return (
      <Button size="sm" className={styles.logSamplesButton} onClick={onSplitOpen}>
        <Trans i18nKey="logs-sample-panel.open-in-split-view">Open logs in split view</Trans>
      </Button>
    );
  };

  let LogsSamplePanelContent: JSX.Element | null;

  if (queryResponse === undefined) {
    LogsSamplePanelContent = null;
  } else if (queryResponse.error !== undefined) {
    LogsSamplePanelContent = (
      <SupplementaryResultError
        error={queryResponse.error}
        title={t('logs-sample.failed-to-load-title', 'Failed to load logs sample for this query')}
      />
    );
  } else if (queryResponse.state === LoadingState.Loading) {
    LogsSamplePanelContent = (
      <span>
        <Trans i18nKey="logs-sample.loading">Logs sample is loading...</Trans>
      </span>
    );
  } else if (queryResponse.data.length === 0 || queryResponse.data.every((frame) => frame.length === 0)) {
    LogsSamplePanelContent = (
      <span>
        <Trans i18nKey="logs-sample.no-data">No logs sample data.</Trans>
      </span>
    );
  } else {
    const logs = dataFrameToLogsModel(queryResponse.data);
    LogsSamplePanelContent = (
      <>
        <OpenInSplitViewButton />
        <div className={styles.logContainer}>
          <LogRows
            logRows={logs.rows}
            dedupStrategy={LogsDedupStrategy.none}
            showLabels={store.getBool(SETTINGS_KEYS.showLabels, false)}
            showTime={store.getBool(SETTINGS_KEYS.showTime, true)}
            wrapLogMessage={store.getBool(SETTINGS_KEYS.wrapLogMessage, true)}
            prettifyLogMessage={store.getBool(SETTINGS_KEYS.prettifyLogMessage, false)}
            timeZone={timeZone}
            enableLogDetails={true}
          />
        </div>
      </>
    );
  }

  return queryResponse?.state !== LoadingState.NotStarted ? (
    <Collapse
      className={styles.logsSamplePanel}
      label={
        <div>
          <Trans i18nKey="logs-sample.logs-sample">Logs sample</Trans>
          <Tooltip content={t('logs-sample.tooltip', 'Show log lines that contributed to visualized metrics')}>
            <Icon name="info-circle" className={styles.infoTooltip} />
          </Tooltip>
        </div>
      }
      isOpen={enabled}
      collapsible={true}
      onToggle={onToggleLogsSampleCollapse}
    >
      {LogsSamplePanelContent}
    </Collapse>
  ) : null;
}

const getStyles = (theme: GrafanaTheme2) => {
  const scrollableLogsContainer = config.featureToggles.exploreScrollableLogsContainer;

  return {
    logsSamplePanel: css`
      ${scrollableLogsContainer && 'max-height: calc(100vh - 115px);'}
    `,
    logSamplesButton: css`
      position: absolute;
      top: ${theme.spacing(1)};
      right: ${theme.spacing(1)};
    `,
    logContainer: css`
      ${scrollableLogsContainer && 'position: relative;'}
      ${scrollableLogsContainer && 'height: 100%;'}
      overflow: scroll;
    `,
    infoTooltip: css`
      margin-left: ${theme.spacing(1)};
    `,
  };
};
