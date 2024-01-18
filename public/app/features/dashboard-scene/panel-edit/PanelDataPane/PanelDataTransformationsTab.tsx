import { css } from '@emotion/css';
import React, { useState } from 'react';
import { DragDropContext, Droppable } from 'react-beautiful-dnd';

import { DataTransformerConfig, GrafanaTheme2, IconName, PanelData } from '@grafana/data';
import { selectors } from '@grafana/e2e-selectors';
import { SceneObjectBase, SceneComponentProps, SceneDataTransformer, SceneQueryRunner } from '@grafana/scenes';
import { Button, ButtonGroup, ConfirmModal, useStyles2 } from '@grafana/ui';
import { TransformationOperationRows } from 'app/features/dashboard/components/TransformationsEditor/TransformationOperationRows';

import { VizPanelManager } from '../VizPanelManager';

import { EmptyTransformationsMessage } from './EmptyTransformationsMessage';
import { TransformationsDrawer } from './TransformationsDrawer';
import { PanelDataPaneTabState, PanelDataPaneTab } from './types';

interface PanelDataTransformationsTabState extends PanelDataPaneTabState {}

export class PanelDataTransformationsTab
  extends SceneObjectBase<PanelDataTransformationsTabState>
  implements PanelDataPaneTab
{
  static Component = PanelDataTransformationsTabRendered;
  tabId = 'transformations';
  icon: IconName = 'process';
  private _panelManager: VizPanelManager;

  getTabLabel() {
    return 'Transformations';
  }

  getItemsCount() {
    return this.getDataTransformer().state.transformations.length;
  }

  constructor(panelManager: VizPanelManager) {
    super({});

    this._panelManager = panelManager;
  }

  public getDataTransformer(): SceneDataTransformer {
    const provider = this._panelManager.state.panel.state.$data;
    if (!provider || !(provider instanceof SceneDataTransformer)) {
      throw new Error('Could not find SceneDataTransformer for panel');
    }
    return provider;
  }

  public getQueryRunner(): SceneQueryRunner {
    return this._panelManager.queryRunner;
  }

  public changeTransformations(transformations: DataTransformerConfig[]) {
    const dataProvider = this.getDataTransformer();
    dataProvider.setState({ transformations });
    dataProvider.reprocessTransformations();
  }
}

export function PanelDataTransformationsTabRendered({ model }: SceneComponentProps<PanelDataTransformationsTab>) {
  const styles = useStyles2(getStyles);
  const sourceData = model.getQueryRunner().useState();
  const { data, transformations: transformsWrongType } = model.getDataTransformer().useState();
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
  const transformations: DataTransformerConfig[] = transformsWrongType as unknown as DataTransformerConfig[];

  const [drawerOpen, setDrawerOpen] = useState<boolean>(false);
  const [confirmModalOpen, setConfirmModalOpen] = useState<boolean>(false);

  const openDrawer = () => setDrawerOpen(true);
  const closeDrawer = () => setDrawerOpen(false);

  if (!data || !sourceData.data) {
    return;
  }

  const transformationsDrawer = (
    <TransformationsDrawer
      onClose={() => setDrawerOpen(false)}
      onTransformationAdd={(selected) => {
        if (selected.value === undefined) {
          return;
        }
        model.changeTransformations([...transformations, { id: selected.value, options: {} }]);
        closeDrawer();
      }}
      isOpen={drawerOpen}
      series={data.series}
    ></TransformationsDrawer>
  );

  if (transformations.length < 1) {
    return (
      <>
        <EmptyTransformationsMessage onShowPicker={openDrawer}></EmptyTransformationsMessage>
        {transformationsDrawer}
      </>
    );
  }

  return (
    <>
      <TransformationsEditor data={sourceData.data} transformations={transformations} model={model} />
      <ButtonGroup>
        <Button
          icon="plus"
          variant="secondary"
          onClick={openDrawer}
          data-testid={selectors.components.Transforms.addTransformationButton}
        >
          Add another transformation
        </Button>
        <Button
          data-testid={selectors.components.Transforms.removeAllTransformationsButton}
          className={styles.removeAll}
          icon="times"
          variant="secondary"
          onClick={() => setConfirmModalOpen(true)}
        >
          Delete all transformations
        </Button>
      </ButtonGroup>
      <ConfirmModal
        isOpen={confirmModalOpen}
        title="Delete all transformations?"
        body="By deleting all transformations, you will go back to the main selection screen."
        confirmText="Delete all"
        onConfirm={() => {
          model.changeTransformations([]);
          setConfirmModalOpen(false);
        }}
        onDismiss={() => setConfirmModalOpen(false)}
      />
      {transformationsDrawer}
    </>
  );
}

interface TransformationEditorProps {
  transformations: DataTransformerConfig[];
  model: PanelDataTransformationsTab;
  data: PanelData;
}

function TransformationsEditor({ transformations, model, data }: TransformationEditorProps) {
  const transformationEditorRows = transformations.map((t, i) => ({ id: `${i} - ${t.id}`, transformation: t }));

  return (
    <DragDropContext onDragEnd={() => {}}>
      <Droppable droppableId="transformations-list" direction="vertical">
        {(provided) => {
          return (
            <div ref={provided.innerRef} {...provided.droppableProps}>
              <TransformationOperationRows
                onChange={(index, transformation) => {
                  const newTransformations = transformations.slice();
                  newTransformations[index] = transformation;
                  model.changeTransformations(newTransformations);
                }}
                onRemove={(index) => {
                  const newTransformations = transformations.slice();
                  newTransformations.splice(index);
                  model.changeTransformations(newTransformations);
                }}
                configs={transformationEditorRows}
                data={data}
              ></TransformationOperationRows>
              {provided.placeholder}
            </div>
          );
        }}
      </Droppable>
    </DragDropContext>
  );
}

const getStyles = (theme: GrafanaTheme2) => ({
  removeAll: css({
    marginLeft: theme.spacing(2),
  }),
});
