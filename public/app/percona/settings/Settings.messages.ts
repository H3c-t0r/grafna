export const Messages = {
  advanced: {
    action: 'Apply changes',
    retentionLabel: 'Data retention',
    retentionTooltip: 'This is the value for how long data will be stored.',
    retentionUnits: 'days',
    retentionLink: `https://per.co.na/data_retention`,
    telemetryLabel: 'Telemetry',
    telemetryLink: `https://per.co.na/telemetry`,
    telemetryTooltip: 'Option to send usage data back to Percona to let us make our product better.',
    telemetryDisclaimer:
      'Disabling Telemetry while Advisors/Percona Alerting is on will prevent PMM from downloading checks and alert templates.',
    telemetrySummaryTitle: 'We gather and send the following information to Percona:',
    updatesLabel: 'Check for updates',
    updatesLink: `https://per.co.na/updates`,
    updatesTooltip: 'Option to check new versions and ability to update PMM from UI.',
    advisorsLabel: 'Advisors',
    sttRareIntervalLabel: 'Rare interval',
    sttStandardIntervalLabel: 'Standard interval',
    sttFrequentIntervalLabel: 'Frequent interval',
    sttCheckIntervalsLabel: 'Execution Intervals',
    sttCheckIntervalTooltip: 'Interval between check runs',
    sttCheckIntervalUnit: 'hours',
    advisorsLink: `https://per.co.na/advisors`,
    advisorsTooltip: 'Enable Advisors and get updated checks from Percona.',
    dbaasLabel: 'Database as a Service (DBaaS)',
    azureDiscoverLabel: 'Microsoft Azure monitoring',
    azureDiscoverTooltip: 'Option to enable/disable Microsoft Azure DB instanced  discovery and monitoring',
    azureDiscoverLink: `https://per.co.na/azure_monitoring`,
    dbaasTooltip:
      'Option to enable/disable DBaaS features. Disabling DBaaS does not suspend or remove running clusters.',
    dbaasLink: `https://per.co.na/dbaas`,
    accessControl: 'Access control',
    accessControlTooltip: 'Option to enable/disable Access control.',
    accessControlLink: 'https://per.co.na/roles_permissions',
    publicAddressLabel: 'Public Address',
    publicAddressTooltip: 'Public Address to this PMM server.',
    publicAddressButton: 'Get from browser',
    alertingLabel: 'Percona Alerting',
    alertingTooltip: 'Option to enable/disable Percona Alerting features.',
    alertingLink: `https://per.co.na/alerting`,
    backupLabel: 'Backup Management',
    backupTooltip: 'Option to enable/disable Backup Management features.',
    backupLink: `https://per.co.na/backup_management`,
    technicalPreviewLegend: 'Technical preview features',
    technicalPreviewDescription:
      'These are technical preview features, not recommended to be used in production environments. Read more\n' +
      '                  about feature status',
    technicalPreviewLinkText: 'here',
    deprecatedFeatures: 'Deprecated features',
  },
  alertmanager: {
    warningPre: "Note: integration with Alertmanager is needed only in cases when you can't use",
    warningLinkContent: 'Percona Alerting',
    warningPost: 'but you need to manage and configure alerts.',
    action: 'Apply Alertmanager settings',
    rulesLabel: 'Prometheus Alerting rules',
    rulesLink: `https://per.co.na/alertmanager`,
    rulesTooltip: 'Alerting rules in the YAML configuration format.',
    urlLabel: 'Alertmanager URL',
    urlLink: `https://per.co.na/alertmanager`,
    urlTooltip: 'The URL of the external Alertmanager to use.',
  },
  diagnostics: {
    action: 'Download server diagnostics',
    label: 'Diagnostics',
    tooltip:
      'You can download server logs to make the problem detection simpler. Please include this file if you are submitting a bug report.',
  },
  metrics: {
    action: 'Apply changes',
    label: 'Metrics resolution, sec',
    link: `https://per.co.na/metrics_resolution`,
    options: {
      rare: 'Rare',
      standard: 'Standard',
      frequent: 'Frequent',
      custom: 'Custom',
    },
    intervals: {
      low: 'Low',
      medium: 'Medium',
      high: 'High',
    },
    tooltip: 'This setting defines how frequently the data will be collected.',
  },
  ssh: {
    action: 'Apply SSH key',
    label: 'SSH key',
    link: `https://per.co.na/ssh_key`,
    tooltip: 'Public SSH key to let you login into the server using SSH.',
  },
  service: {
    success: 'Settings updated',
  },
  tabs: {
    metrics: 'Metrics Resolution',
    advanced: 'Advanced Settings',
    ssh: 'SSH Key',
    alertManager: 'Alertmanager Integration',
    perconaPlatform: 'Percona Platform',
    communication: 'Communication',
  },
  tooltipLinkText: 'Read more',
  unauthorized: 'Insufficient access permissions.',
};
