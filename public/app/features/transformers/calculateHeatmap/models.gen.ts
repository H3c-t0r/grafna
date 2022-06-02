import { ScaleDistributionConfig } from '@grafana/schema';

export enum HeatmapCalculationMode {
  Size = 'size', // When exponential, this is "splitFactor"
  Count = 'count',
}

export const enum HeatmapBucketLayout {
  le = 'le',
  ge = 'ge',
  unknown = 'unknown', // unknown
  auto = 'auto', // becomes unknown
}

export interface HeatmapCalculationAxisConfig {
  mode?: HeatmapCalculationMode;
  value?: string; // number or interval string ie 10s
  scale?: ScaleDistributionConfig;
}

export interface HeatmapCalculationOptions {
  xAxis?: HeatmapCalculationAxisConfig;
  yAxis?: HeatmapCalculationAxisConfig;
  xAxisField?: string; // name of the x field
}
