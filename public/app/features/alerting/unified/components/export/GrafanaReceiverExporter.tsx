import React, { useState } from 'react';

import { LoadingPlaceholder } from '@grafana/ui';

import { alertRuleApi } from '../../api/alertRuleApi';

import { FileExportPreview } from './FileExportPreview';
import { GrafanaExportDrawer } from './GrafanaExportDrawer';
import { ExportFormats, ExportFormatsWithoutHCL } from './providers';

interface GrafanaReceiverExportPreviewProps {
  exportFormat: ExportFormatsWithoutHCL | ExportFormats;
  onClose: () => void;
  receiverName: string;
  decrypt: string;
}

const GrafanaReceiverExportPreview = ({
  receiverName,
  decrypt,
  exportFormat,
  onClose,
}: GrafanaReceiverExportPreviewProps) => {
  const { currentData: receiverDefinition = '', isFetching } = alertRuleApi.useExportReceiverQuery({
    receiverName: receiverName,
    decrypt: decrypt,
    format: exportFormat,
  });

  const downloadFileName = `cp-${receiverName}-${new Date().getTime()}`;

  if (isFetching) {
    return <LoadingPlaceholder text="Loading...." />;
  }

  return (
    <FileExportPreview
      format={exportFormat}
      textDefinition={receiverDefinition}
      downloadFileName={downloadFileName}
      onClose={onClose}
    />
  );
};

interface GrafanaReceiverExporterProps {
  onClose: () => void;
  receiverName: string;
  decrypt: string;
}

export const GrafanaReceiverExporter = ({ onClose, receiverName, decrypt }: GrafanaReceiverExporterProps) => {
  const [activeTab, setActiveTab] = useState<ExportFormatsWithoutHCL | ExportFormats>('yaml');

  return (
    <GrafanaExportDrawer activeTab={activeTab} onTabChange={setActiveTab} onClose={onClose}>
      <GrafanaReceiverExportPreview
        receiverName={receiverName}
        decrypt={decrypt}
        exportFormat={activeTab}
        onClose={onClose}
      />
    </GrafanaExportDrawer>
  );
};
