import React from 'react';

import {
  DataTransformerConfig,
  FieldConfigSource,
  PanelData,
  PanelProps,
  VisualizationSuggestion,
} from '@grafana/data';

/**
 * Describes the properties that can be passed to the PanelDataErrorView.
 *
 * @alpha
 */
export interface PanelDataErrorViewProps {
  panelProps?: PanelProps; // The raw panel props
  message?: string;
  panelId: number;
  data: PanelData;
  fieldConfig?: FieldConfigSource;
  needsTimeField?: boolean;
  needsNumberField?: boolean;
  needsStringField?: boolean;
  transformSuggestion?: DataTransformerConfig; // array?
  suggestions?: VisualizationSuggestion[];
}

/**
 * Simplified type with defaults that describes the PanelDataErrorView.
 *
 * @internal
 */
export type PanelDataErrorViewType = React.ComponentType<PanelDataErrorViewProps>;

/**
 * PanelDataErrorView allows panels to show a consistent error message when
 * the result structure does not meet expected criteria
 *
 * @alpha
 */
export let PanelDataErrorView: PanelDataErrorViewType = ({ message }) => {
  return <div>Unable to render data: {message}.</div>;
};

/**
 * Used to bootstrap the PanelDataErrorView during application start so the
 * PanelDataErrorView is exposed via runtime.
 *
 * @internal
 */
export function setPanelDataErrorView(renderer: PanelDataErrorViewType) {
  PanelDataErrorView = renderer;
}
