import React from 'react';

import { Drawer } from '@grafana/ui';

import { RuleInspectorTabs } from '../rule-editor/RuleInspector';

import { ExportFormats, ExportProvider } from './providers';

interface GrafanaExportDrawerProps {
  activeTab: ExportFormats;
  onTabChange: (tab: ExportFormats) => void;
  children: React.ReactNode;
  onClose: () => void;
  formatProviders: Array<ExportProvider<ExportFormats>>;
}

export function GrafanaExportDrawer({
  activeTab,
  onTabChange,
  children,
  onClose,
  formatProviders,
}: GrafanaExportDrawerProps) {
  const grafanaRulesTabs = Object.values(formatProviders).map((provider) => ({
    label: provider.name,
    value: provider.exportFormat,
  }));

  return (
    <Drawer
      title="Export"
      subtitle="Select the format and download the file or copy the contents to clipboard"
      tabs={
        <RuleInspectorTabs<ExportFormats> tabs={grafanaRulesTabs} setActiveTab={onTabChange} activeTab={activeTab} />
      }
      onClose={onClose}
      size="md"
    >
      {children}
    </Drawer>
  );
}
