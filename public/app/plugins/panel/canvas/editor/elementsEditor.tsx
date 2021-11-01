import { NestedPanelOptions } from '@grafana/data/src/utils/OptionsUIBuilders';
import { Scene } from 'app/features/canvas/runtime/scene';

export interface CanvasEditorGroupOptions {
  scene: Scene;
  category?: string[];
}

export const getElementsEditor = (opts: CanvasEditorGroupOptions): NestedPanelOptions<any> => {
  return {
    category: opts.category,
    path: '--',
    build: (builder, context) => {},
  };
};
