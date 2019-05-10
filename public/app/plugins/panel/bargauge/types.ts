import { VizOrientation, SelectOptionItem, SingleStatBaseOptions } from '@grafana/ui';
import { standardGaugeFieldOptions } from '../gauge/types';

export interface BarGaugeOptions extends SingleStatBaseOptions {
  displayMode: 'basic' | 'lcd' | 'gradient';
}

export const displayModes: Array<SelectOptionItem<string>> = [
  { value: 'gradient', label: 'Gradient' },
  { value: 'lcd', label: 'Retro LCD' },
  { value: 'basic', label: 'Basic' },
];

export const orientationOptions: Array<SelectOptionItem<VizOrientation>> = [
  { value: VizOrientation.Horizontal, label: 'Horizontal' },
  { value: VizOrientation.Vertical, label: 'Vertical' },
];

export const defaults: BarGaugeOptions = {
  displayMode: 'lcd',
  orientation: VizOrientation.Horizontal,
  fieldOptions: standardGaugeFieldOptions,
};
