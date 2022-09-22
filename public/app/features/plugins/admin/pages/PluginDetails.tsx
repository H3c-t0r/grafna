import { css } from '@emotion/css';
import React, { useEffect } from 'react';
import { usePrevious } from 'react-use';

import { GrafanaTheme2, NavModelItem } from '@grafana/data';
import { locationService } from '@grafana/runtime';
import { useStyles2, TabContent, Alert } from '@grafana/ui';
import { Layout } from '@grafana/ui/src/components/Layout/Layout';
import { Page } from 'app/core/components/Page/Page';
import { GrafanaRouteComponentProps } from 'app/core/navigation/types';
import { AppNotificationSeverity } from 'app/types';

import { Loader } from '../components/Loader';
import { PluginDetailsBody } from '../components/PluginDetailsBody';
import { PluginDetailsDisabledError } from '../components/PluginDetailsDisabledError';
import { PluginDetailsHeader } from '../components/PluginDetailsHeader';
import { PluginDetailsSignature } from '../components/PluginDetailsSignature';
import { usePluginDetailsTabs } from '../hooks/usePluginDetailsTabs';
import { useGetSingle, useFetchStatus, useFetchDetailsStatus } from '../state/hooks';
import { PluginTabLabels, PluginTabIds, PluginDetailsTab } from '../types';

type Props = GrafanaRouteComponentProps<{ pluginId?: string }>;

export default function PluginDetails({ match, queryParams }: Props): JSX.Element | null {
  const {
    params: { pluginId = '' },
    url,
  } = match;
  const parentUrl = url.substring(0, url.lastIndexOf('/'));

  const defaultTabs: PluginDetailsTab[] = [
    {
      label: PluginTabLabels.OVERVIEW,
      icon: 'file-alt',
      id: PluginTabIds.OVERVIEW,
      href: `${url}?page=${PluginTabIds.OVERVIEW}`,
    },
  ];
  const plugin = useGetSingle(pluginId); // fetches the localplugin settings
  const { tabs, defaultTab } = usePluginDetailsTabs(plugin, defaultTabs);
  const { isLoading: isFetchLoading } = useFetchStatus();
  const { isLoading: isFetchDetailsLoading } = useFetchDetailsStatus();
  const styles = useStyles2(getStyles);
  const prevTabs = usePrevious(tabs);
  const pageId = (queryParams.page as PluginTabIds) || defaultTab;

  // If an app plugin is uninstalled we need to reset the active tab when the config / dashboards tabs are removed.
  useEffect(() => {
    const hasUninstalledWithConfigPages = prevTabs && prevTabs.length > tabs.length;
    const isViewingAConfigPage = pageId !== PluginTabIds.OVERVIEW && pageId !== PluginTabIds.VERSIONS;

    if (hasUninstalledWithConfigPages && isViewingAConfigPage) {
      locationService.replace(`${url}?page=${PluginTabIds.OVERVIEW}`);
    }
  }, [pageId, url, tabs, prevTabs]);

  if (isFetchLoading || isFetchDetailsLoading) {
    return (
      <Page navId="plugins">
        <Loader />
      </Page>
    );
  }

  if (!plugin) {
    return (
      <Layout justify="center" align="center">
        <Alert severity={AppNotificationSeverity.Warning} title="Plugin not found">
          That plugin cannot be found. Please check the url is correct or <br />
          go to the <a href={parentUrl}>plugin catalog</a>.
        </Alert>
      </Layout>
    );
  }

  const pageNav: NavModelItem = {
    text: plugin.name,
    img: plugin.info.logos.small,
    breadcrumbs: [{ title: 'Plugins', url: parentUrl }],
    children: tabs.map((tab) => ({
      text: tab.label,
      icon: tab.icon,
      url: tab.href,
      active: tab.id === pageId,
    })),
    headerExtra: () => <PluginDetailsHeader plugin={plugin} />,
  };

  return (
    <Page navId="plugins" pageNav={pageNav}>
      <Page.Contents>
        <TabContent className={styles.tabContent}>
          <PluginDetailsSignature plugin={plugin} className={styles.alert} />
          <PluginDetailsDisabledError plugin={plugin} className={styles.alert} />
          <PluginDetailsBody queryParams={queryParams} plugin={plugin} pageId={pageId} />
        </TabContent>
      </Page.Contents>
    </Page>
  );
}

export const getStyles = (theme: GrafanaTheme2) => {
  return {
    alert: css`
      margin-bottom: ${theme.spacing(2)};
    `,
    // Needed due to block formatting context
    tabContent: css`
      overflow: auto;
      height: 100%;
    `,
  };
};
