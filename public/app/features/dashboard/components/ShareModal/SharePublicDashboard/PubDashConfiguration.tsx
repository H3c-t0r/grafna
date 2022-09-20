import { css } from '@emotion/css';
import React from 'react';

import { DateTime, GrafanaTheme2 } from '@grafana/data/src';
import { selectors as e2eSelectors } from '@grafana/e2e-selectors/src';
import { reportInteraction } from '@grafana/runtime/src';
import { FieldSet, Label, Switch, TimeRangeInput, useStyles2, VerticalGroup } from '@grafana/ui/src';
import { Layout } from '@grafana/ui/src/components/Layout/Layout';
import { TimeModel } from 'app/features/dashboard/state/TimeModel';
import { useIsDesktop } from 'app/features/dashboard/utils/screen';
import { getTimeRange } from 'app/features/dashboard/utils/timeRange';

export const PubDashConfiguration = ({
  disabled,
  isPubDashEnabled,
  hasTemplateVariables,
  time,
  onToggleEnabled,
}: {
  disabled: boolean;
  isPubDashEnabled?: boolean;
  onToggleEnabled: () => void;
  hasTemplateVariables: boolean;
  time: { from: DateTime | string; to: DateTime | string; timeZone: TimeModel };
}) => {
  const selectors = e2eSelectors.pages.ShareDashboardModal.PublicDashboard;
  const styles = useStyles2(getStyles);
  const isDesktop = useIsDesktop();

  const timeRange = getTimeRange({ from: time.from, to: time.to }, time.timeZone);

  return (
    <>
      <h4 className="share-modal-info-text">Public dashboard configuration</h4>
      <FieldSet disabled={disabled} className={styles.dashboardConfig}>
        <VerticalGroup spacing="md">
          <Layout orientation={isDesktop ? 0 : 1} spacing="xs" justify="space-between">
            <Label description="The public dashboard uses the default time settings of the dashboard">Time Range</Label>
            <TimeRangeInput value={timeRange} disabled onChange={() => {}} />
          </Layout>
          <Layout orientation={isDesktop ? 0 : 1} spacing="xs" justify="space-between">
            <Label description="Configures whether current dashboard can be available publicly">Enabled</Label>
            <Switch
              disabled={hasTemplateVariables}
              data-testid={selectors.EnableSwitch}
              value={isPubDashEnabled}
              onChange={() => {
                reportInteraction('grafana_dashboards_public_enable_clicked', {
                  action: isPubDashEnabled ? 'disable' : 'enable',
                });

                onToggleEnabled();
              }}
            />
          </Layout>
        </VerticalGroup>
      </FieldSet>
    </>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  dashboardConfig: css`
    margin: ${theme.spacing(0, 0, 3, 0)};
  `,
});
