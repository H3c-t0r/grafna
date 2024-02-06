import { isEmpty } from 'lodash';
import React from 'react';
import { useAsync, useToggle } from 'react-use';

import { Dashboard, Panel, RowPanel } from '@grafana/schema';
import { Alert, Collapse, Column, InteractiveTable, TextLink } from '@grafana/ui';
import { getDashboardScenePageStateManager } from 'app/features/dashboard-scene/pages/DashboardScenePageStateManager';
import { DashboardRoutes } from 'app/types';

import { makePanelLink } from '../utils/misc';

interface DeprecationNoticeProps {
  dashboardUid?: string;
  dashboard?: Dashboard;
}

export default function LegacyAlertsDeprecationNotice({ dashboardUid, dashboard }: DeprecationNoticeProps) {
  const dashboardStateManager = getDashboardScenePageStateManager();

  const {
    loading,
    value: dashboardData,
    error,
  } = useAsync(() => {
    if (dashboard) {
      return Promise.resolve(dashboard);
    } else if (dashboardUid) {
      return dashboardStateManager
        .fetchDashboard({
          uid: dashboardUid,
          route: DashboardRoutes.Normal,
        })
        .then((data) => (data ? data.dashboard : undefined));
    } else {
      throw new Error('LegacyAlertsDeprecationNotice missing any of "dashboardUid" or "dashboardModel"');
    }
  }, [dashboardStateManager]);

  if (loading) {
    return null;
  }

  // we probably don't want to show the user an error if this fails because there's nothing they can do about it
  if (error) {
    console.error(error);
  }

  const panelsWithLegacyAlerts = getLegacyAlertPanelsFromDashboard(dashboardData);

  // don't show anything when the user has no legacy alerts defined
  const hasLegacyAlerts = !isEmpty(panelsWithLegacyAlerts);
  if (!hasLegacyAlerts) {
    return null;
  }

  return <LegacyAlertsWarning dashboardUid={dashboardData.uid} panels={panelsWithLegacyAlerts} />;
}

/**
 * This function uses two different ways to detect legacy alerts based on what dashboard system is being used.
 *
 * 1. if using the older (non-scenes) dashboard system we can simply check for "alert" in the panel definition.
 * 2. for dashboard scenes the alerts are no longer added to the model but we can check for "alertThreshold" in the panel options object
 */
function getLegacyAlertPanelsFromDashboard(dashboard: Dashboard): Panel[] {
  const panelsWithLegacyAlerts = dashboard.panels?.filter((panel) => {
    const hasAlertDefinition = 'alert' in panel;
    const hasAlertThreshold = 'options' in panel && panel.options ? 'alertThreshold' in panel.options : false;
    return hasAlertDefinition || hasAlertThreshold;
  });

  return panelsWithLegacyAlerts ?? [];
}

interface Props {
  dashboardUid: string;
  panels: Panel[];
}

function LegacyAlertsWarning({ dashboardUid, panels }: Props) {
  const [isOpen, toggleCollapsible] = useToggle(false);

  const columns: Array<Column<Panel | RowPanel>> = [
    { id: 'id', header: 'ID' },
    {
      id: 'title',
      header: 'Title',
      cell: (cell) => (
        <TextLink
          external
          href={makePanelLink(dashboardUid, String(cell.row.id), { editPanel: cell.row.id, tab: 'alert' })}
        >
          {cell.value}
        </TextLink>
      ),
    },
  ];

  return (
    <Alert severity="warning" title="Legacy alert rules are deprecated">
      <p>
        You have legacy alert rules in this dashboard that were deprecated in Grafana 11 and are no longer supported.
      </p>
      <p>
        Refer to{' '}
        <TextLink href="https://grafana.com/docs/grafana/latest/alerting/set-up/migrating-alerts/" external>
          our documentation
        </TextLink>{' '}
        on how to migrate legacy alert rules and how to import and export using Grafana Alerting.
      </p>

      <Collapse label={'List of panels using legacy alerts'} collapsible isOpen={isOpen} onToggle={toggleCollapsible}>
        <InteractiveTable columns={columns} data={panels} getRowId={(panel) => String(panel.id)} pageSize={5} />
      </Collapse>
    </Alert>
  );
}
