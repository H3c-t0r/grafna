export type UpdateConfig = {
  [K in keyof UpdateCounters]: boolean;
};

export type UpdateCounters = {
  render: number;
  dataChanged: number;
  schemaChanged: number;
};

export interface DebugPanelOptions {
  counters?: UpdateConfig;
}
