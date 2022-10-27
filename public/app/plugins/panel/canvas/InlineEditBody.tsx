import { css } from '@emotion/css';
import { get as lodashGet } from 'lodash';
import React, { useMemo, useState } from 'react';
import { useObservable } from 'react-use';

import {
  DataFrame,
  GrafanaTheme2,
  PanelOptionsEditorBuilder,
  SelectableValue,
  StandardEditorContext,
} from '@grafana/data';
import { PanelOptionsSupplier } from '@grafana/data/src/panel/PanelPlugin';
import { NestedValueAccess } from '@grafana/data/src/utils/OptionsUIBuilders';
import { useStyles2 } from '@grafana/ui/src';
import { AddLayerButton } from 'app/core/components/Layers/AddLayerButton';
import { notFoundItem } from 'app/features/canvas/elements/notFound';
import { ElementState } from 'app/features/canvas/runtime/element';
import { FrameState } from 'app/features/canvas/runtime/frame';
import { OptionsPaneCategory } from 'app/features/dashboard/components/PanelEditor/OptionsPaneCategory';
import { OptionsPaneCategoryDescriptor } from 'app/features/dashboard/components/PanelEditor/OptionsPaneCategoryDescriptor';
import { fillOptionsPaneItems } from 'app/features/dashboard/components/PanelEditor/getVisualizationOptions';
import { setOptionImmutably } from 'app/features/dashboard/components/PanelEditor/utils';

import { CanvasElementOptions, canvasElementRegistry } from '../../../features/canvas';

import { activePanelSubject, InstanceState } from './CanvasPanel';
import { TabsEditor } from './editor/TabsEditor';
import { getElementEditor } from './editor/elementEditor';
import { getLayerEditor } from './editor/layerEditor';
import { addStandardCanvasEditorOptions } from './module';
import { InlineEditTabs } from './types';
import { getElementTypes } from './utils';

export function InlineEditBody() {
  const activePanel = useObservable(activePanelSubject);
  const instanceState = activePanel?.panel.context?.instanceState;
  const styles = useStyles2(getStyles);

  const [activeTab, setActiveTab] = useState<string>(InlineEditTabs.SelectedElement);

  const pane = useMemo(() => {
    const p = activePanel?.panel;
    const state: InstanceState = instanceState;
    if (!state || !p) {
      return new OptionsPaneCategoryDescriptor({ id: 'root', title: 'root' });
    }

    const supplier = (builder: PanelOptionsEditorBuilder<any>, context: StandardEditorContext<any>) => {
      if (activeTab === InlineEditTabs.ElementManagement) {
        builder.addNestedOptions(getLayerEditor(instanceState));
      }

      const selection = state.selected;
      if (selection?.length === 1 && activeTab === InlineEditTabs.SelectedElement) {
        const element = selection[0];
        if (element && !(element instanceof FrameState)) {
          builder.addNestedOptions(
            getElementEditor({
              category: [`Selected element (${element.options.name})`],
              element,
              scene: state.scene,
            })
          );
        }
      }

      addStandardCanvasEditorOptions(builder);
    };

    return getOptionsPaneCategoryDescriptor(
      {
        options: p.props.options,
        onChange: p.props.onOptionsChange,
        data: p.props.data?.series,
      },
      supplier
    );
  }, [instanceState, activePanel, activeTab]);

  const topLevelItemsContainerStyle = {
    marginLeft: 15,
    marginTop: 10,
  };

  const onTabChange = (tab: string) => {
    setActiveTab(tab);
  };

  const typeOptions = getElementTypes(instanceState?.scene.shouldShowAdvancedTypes).options;
  const rootLayer: FrameState | undefined = instanceState?.layer;

  const onAddItem = (sel: SelectableValue<string>) => {
    const newItem = canvasElementRegistry.getIfExists(sel.value) ?? notFoundItem;
    const newElementOptions = newItem.getNewOptions() as CanvasElementOptions;
    newElementOptions.type = newItem.id;
    if (newItem.defaultSize) {
      newElementOptions.placement = { ...newElementOptions.placement, ...newItem.defaultSize };
    }
    if (rootLayer) {
      const newElement = new ElementState(newItem, newElementOptions, rootLayer);
      newElement.updateData(rootLayer.scene.context);
      rootLayer.elements.push(newElement);
      rootLayer.scene.save();

      rootLayer.reinitializeMoveable();
    }
  };

  return (
    <>
      <div style={topLevelItemsContainerStyle}>{pane.items.map((item) => item.render())}</div>
      <div style={topLevelItemsContainerStyle}>
        <AddLayerButton onChange={onAddItem} options={typeOptions} label={'Add item'} />
      </div>
      <div style={topLevelItemsContainerStyle}>
        <TabsEditor onTabChange={onTabChange} />
        {pane.categories.map((p) => renderOptionsPaneCategoryDescriptor(p, activeTab))}
        {instanceState && activeTab === InlineEditTabs.SelectedElement && instanceState.selected.length === 0 && (
          <div className={styles.selectElement}>Please select an element</div>
        )}
      </div>
    </>
  );
}

// Recursively render options
function renderOptionsPaneCategoryDescriptor(
  pane: OptionsPaneCategoryDescriptor,
  activeTab?: string,
  instanceState?: InstanceState
) {
  return (
    <OptionsPaneCategory {...pane.props} key={pane.props.id}>
      <div>{pane.items.map((v) => v.render())}</div>
      {pane.categories.map((c) => renderOptionsPaneCategoryDescriptor(c))}
    </OptionsPaneCategory>
  );
}

interface EditorProps<T> {
  onChange: (v: T) => void;
  options: T;
  data?: DataFrame[];
}

function getOptionsPaneCategoryDescriptor<T = any>(
  props: EditorProps<T>,
  supplier: PanelOptionsSupplier<T>
): OptionsPaneCategoryDescriptor {
  const context: StandardEditorContext<unknown, unknown> = {
    data: props.data ?? [],
    options: props.options,
  };

  const root = new OptionsPaneCategoryDescriptor({ id: 'root', title: 'root' });
  const getOptionsPaneCategory = (categoryNames?: string[]): OptionsPaneCategoryDescriptor => {
    if (categoryNames?.length) {
      const key = categoryNames[0];
      let sub = root.categories.find((v) => v.props.id === key);
      if (!sub) {
        sub = new OptionsPaneCategoryDescriptor({ id: key, title: key });
        root.categories.push(sub);
      }
      return sub;
    }
    return root;
  };

  const access: NestedValueAccess = {
    getValue: (path: string) => lodashGet(props.options, path),
    onChange: (path: string, value: any) => {
      props.onChange(setOptionImmutably(props.options as any, path, value));
    },
  };

  // Use the panel options loader
  fillOptionsPaneItems(supplier, access, getOptionsPaneCategory, context);
  return root;
}

const getStyles = (theme: GrafanaTheme2) => ({
  selectElement: css`
    color: ${theme.colors.text.secondary};
    padding: ${theme.spacing(2)};
  `,
});
