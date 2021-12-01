import React from 'react';
import { Container } from '@grafana/ui';
import { StandardEditorProps } from '@grafana/data';
import { DropResult } from 'react-beautiful-dnd';

import { GeomapPanelOptions, MapLayerState } from '../types';
import { GeomapInstanceState } from '../GeomapPanel';
import { AddLayerButton } from './LayerDragDropList/AddLayerButton';
import { LayerDragDropList } from './LayerDragDropList/LayerDragDropList';
import { geomapLayerRegistry } from '../layers/registry';
import { dataLayerFilter } from './layerEditor';

type LayersEditorProps = StandardEditorProps<any, any, GeomapPanelOptions, GeomapInstanceState>;

export const LayersEditor = (props: LayersEditorProps) => {
  const { layers, selected, actions } = props.context.instanceState ?? {};
  if (!layers || !actions) {
    return <div>No layers?</div>;
  }

  const onDragEnd = (result: DropResult) => {
    if (!result.destination) {
      return;
    }

    const { layers, actions } = props.context.instanceState ?? {};
    if (!layers || !actions) {
      return;
    }

    // account for the reverse order and offset (0 is baselayer)
    const count = layers.length - 1;
    const src = (result.source.index - count) * -1;
    const dst = (result.destination.index - count) * -1;

    actions.reorder(src, dst);
  };

  const onSelect = (element: MapLayerState<any>) => {
    actions.selectLayer(element.options.name);
  };

  const onDelete = (element: MapLayerState<any>) => {
    actions.deleteLayer(element.options.name);
  };

  const getLayerType = (element: MapLayerState<any>) => {
    return element.options.type;
  };

  const selection = selected ? [selected] : [];

  return (
    <>
      <Container>
        <AddLayerButton
          onChange={(v) => actions.addlayer(v.value!)}
          options={geomapLayerRegistry.selectOptions(undefined, dataLayerFilter).options}
          label={'Add layer'}
        />
      </Container>
      <br />

      <LayerDragDropList
        layers={layers}
        getLayerType={getLayerType}
        onDragEnd={onDragEnd}
        onSelect={onSelect}
        onDelete={onDelete}
        selection={selection}
        excludeBaseLayer
        selectByIndex
        verifyLayerNameUniqueness={actions.canRename}
      />
    </>
  );
};
