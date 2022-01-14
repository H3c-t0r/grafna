import React, { useMemo } from 'react';

import { useStyles } from '@grafana/ui';

import { CheckPermissions } from '../shared/components/Elements/CheckPermissions/CheckPermissions';
import { TabbedContent, ContentTab } from '../shared/components/Elements/TabbedContent';
import PageWrapper from '../shared/components/PageWrapper/PageWrapper';

import { PAGE_MODEL } from './Inventory.constants';
import { Messages } from './Inventory.messages';
import { getStyles } from './Inventory.styles';
import { TabKeys } from './Inventory.types';
import { Agents, NodesTab, Services } from './Tabs';

export const DEFAULT_TAB = 'services';

export const InventoryPanel = () => {
  const { path: basePath } = PAGE_MODEL;
  const styles = useStyles(getStyles);

  const tabs: ContentTab[] = useMemo(
    (): ContentTab[] => [
      {
        label: Messages.tabs.services,
        key: TabKeys.services,
        component: <Services />,
      },
      {
        label: Messages.tabs.agents,
        key: TabKeys.agents,
        component: <Agents />,
      },
      {
        label: Messages.tabs.nodes,
        key: TabKeys.nodes,
        component: <NodesTab />,
      },
    ],
    []
  );

  return (
    <PageWrapper pageModel={PAGE_MODEL}>
      <div className={styles.inventoryWrapper}>
        <TabbedContent
          tabs={tabs}
          basePath={basePath}
          renderTab={({ Content }) => (
            <CheckPermissions>
              <Content />
            </CheckPermissions>
          )}
        />
      </div>
    </PageWrapper>
  );
};

export default InventoryPanel;
