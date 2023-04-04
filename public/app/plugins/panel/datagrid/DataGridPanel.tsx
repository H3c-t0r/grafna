import { css } from '@emotion/css';
import DataEditor, {
  GridCell,
  Item,
  GridColumn,
  GridCellKind,
  EditableGridCell,
  GridColumnIcon,
} from '@glideapps/glide-data-grid';
import React, { useCallback, useEffect, useState } from 'react';

import {
  ArrayVector,
  DataFrame,
  Field,
  MutableDataFrame,
  PanelProps,
  GrafanaTheme2,
  getFieldDisplayName,
  FieldType,
} from '@grafana/data';
import { PanelDataErrorView } from '@grafana/runtime';
// eslint-disable-next-line import/order
import { useTheme2 } from '@grafana/ui';

import '@glideapps/glide-data-grid/dist/index.css';

import { getDashboardSrv } from 'app/features/dashboard/services/DashboardSrv';

import { AddColumn } from './components/AddColumn';
import { PanelOptions } from './panelcfg.gen';
import { EMPTY_DF, GRAFANA_DS, publishSnapshot } from './utils';

const ICON_WIDTH = 30;

interface Props extends PanelProps<PanelOptions> {}

export const DataGridPanel: React.FC<Props> = ({ options, data, id, fieldConfig }) => {
  const [gridData, setGridData] = useState<DataFrame>(data.series[options.selectedSeries ?? 0]);
  const [columns, setColumns] = useState<GridColumn[]>([]);
  const [isSnapshotted, setIsSnapshotted] = useState<boolean>(false);

  const theme = useTheme2();
  const gridTheme = {
    accentColor: theme.colors.primary.main,
    accentFg: theme.colors.secondary.main,
    textDark: theme.colors.text.primary,
    textMedium: theme.colors.text.primary,
    textLight: theme.colors.text.primary,
    textBubble: theme.colors.text.primary,
    textHeader: theme.colors.text.primary,
    bgCell: theme.colors.background.primary,
    bgCellMedium: theme.colors.background.primary,
    bgHeader: theme.colors.background.secondary,
    bgHeaderHasFocus: theme.colors.background.secondary,
    bgHeaderHovered: theme.colors.background.secondary,
  };

  const setGridColumns = useCallback(() => {
    const typeToIconMap: Map<string, GridColumnIcon> = new Map([
      [FieldType.number, GridColumnIcon.HeaderNumber],
      [FieldType.string, GridColumnIcon.HeaderTextTemplate],
      [FieldType.boolean, GridColumnIcon.HeaderBoolean],
    ]);

    if (!gridData) {
      return;
    }

    setColumns(
      gridData.fields.map((f) => {
        const displayName = getFieldDisplayName(f, gridData);
        const width = displayName.length * theme.typography.fontSize + ICON_WIDTH;
        return { title: displayName, width: width, icon: typeToIconMap.get(f.type) };
      })
    );
  }, [gridData, theme.typography.fontSize]);

  useEffect(() => {
    setGridColumns();
  }, [gridData, setGridColumns]);

  useEffect(() => {
    const panelModel = getDashboardSrv().getCurrent()?.getPanelById(id);

    if (panelModel?.datasource?.type !== GRAFANA_DS.type) {
      setIsSnapshotted(false);
    }
  }, [id, options, data]);

  useEffect(() => {
    if (!isSnapshotted) {
      setGridData(data.series[options.selectedSeries ?? 0]);
    }
  }, [data, isSnapshotted, options.selectedSeries]);

  useEffect(() => {
    if (isSnapshotted) {
      publishSnapshot(gridData, id);
    }
  }, [gridData, id, isSnapshotted]);

  const getCellContent = ([col, row]: Item): GridCell => {
    const field: Field = gridData.fields[col];

    if (!field) {
      throw new Error('OH NO');
    }

    const value = field.values.get(row);

    if (value === undefined || value === null) {
      throw new Error('OH NO 2');
    }

    switch (field.type) {
      case FieldType.boolean:
        return {
          kind: GridCellKind.Boolean,
          data: value,
          allowOverlay: false,
        };
      case FieldType.number:
        return {
          kind: GridCellKind.Number,
          data: value,
          allowOverlay: true,
          readonly: false,
          displayData: value.toString(),
        };
      default:
        return {
          kind: GridCellKind.Text,
          data: value,
          allowOverlay: true,
          readonly: false,
          displayData: value.toString(),
        };
    }
  };

  const onCellEdited = (cell: Item, newValue: EditableGridCell) => {
    const [col, row] = cell;
    const field: Field = gridData.fields[col];

    if (!field || !newValue.data) {
      throw new Error('OH NO 3');
    }

    const values = field.values.toArray();

    //TODO REMOVE THIS AND LET THE USER CHOOSE THE TYPE BY APPLYING TRANSFORMS!
    // //todo maybe come back to this later and look for a better way
    // //Convert field type and value between string and number if needed
    // //If field type is number we check if the new value is numeric. If it isn't we change the field type
    // let val = newValue.data;
    // if (field.type === FieldType.number) {
    //   if (!isNumeric(val)) {
    //     field.type = FieldType.string;
    //   } else {
    //     val = Number(val);
    //   }
    //   //If field type is string we check if the new value is numeric. If it is numeric and all other fields are also numeric we change the field type
    //   //If we change the field type we also convert all other values to numbers
    // } else if (field.type === FieldType.string) {
    //   if (isNumeric(val) && values.filter((_, index) => index !== row).findIndex((v) => !isNumeric(v)) === -1) {
    //     field.type = FieldType.number;
    //     val = Number(val);

    //     if (values.findIndex((v) => typeof v === 'string') !== -1) {
    //       values.forEach((v, index) => {
    //         if (typeof v === 'string') {
    //           values[index] = Number(v);
    //         }
    //       });
    //     }
    //   }
    // }

    values[row] = newValue.data;
    field.values = new ArrayVector(values);

    setIsSnapshotted(true);
    setGridData(new MutableDataFrame(gridData));
  };

  const onColumnInputBlur = (columnName: string) => {
    const len = gridData.length ?? 50; //todo ?????? 50????

    const newFrame = new MutableDataFrame(gridData);

    const field: Field = {
      name: columnName,
      type: FieldType.string,
      config: {},
      values: new ArrayVector(new Array(len).fill('')),
    };

    newFrame.addField(field);

    setIsSnapshotted(true);
    setGridData(newFrame);
  };

  const addNewRow = () => {
    const newFrame = new MutableDataFrame(gridData);

    for (const field of newFrame.fields) {
      field.values.add('');
    }

    setIsSnapshotted(true);
    setGridData(newFrame);
  };

  const onColumnResize = (column: GridColumn, newSize: number, colIndex: number, newSizeWithGrow: number) => {
    setColumns((prevColumns) => {
      const newColumns = [...prevColumns];
      newColumns[colIndex] = { title: column.title, icon: column.icon, width: newSize };
      return newColumns;
    });
  };

  if (!document.getElementById('portal')) {
    const portal = document.createElement('div');
    portal.id = 'portal';
    document.body.appendChild(portal);
  }

  if (!gridData) {
    return <PanelDataErrorView panelId={id} fieldConfig={fieldConfig} data={data} />;
  }

  //TODO multiple series support
  const numRows = gridData.length;
  const styles = getStyles(theme);

  return (
    <>
      <DataEditor
        getCellContent={getCellContent}
        columns={columns}
        rows={numRows}
        width={'100%'}
        height={'100%'}
        theme={gridTheme}
        onCellEdited={onCellEdited}
        onHeaderClicked={() => {
          console.log('header clicked');
        }}
        onRowAppended={addNewRow}
        rowMarkers={'clickable-number'}
        onColumnResize={onColumnResize}
        onCellContextMenu={(cell) => {
          console.log(cell);
        }}
        trailingRowOptions={{
          sticky: false,
          tint: true,
          targetColumn: 0,
        }}
        rightElement={<AddColumn onColumnInputBlur={onColumnInputBlur} divStyle={styles.addColumnDiv} />}
        rightElementProps={{
          fill: true,
          sticky: false,
        }}
      />
    </>
  );
};

const getStyles = (theme: GrafanaTheme2) => {
  const height = '37px';
  const width = '120px';

  return {
    addColumnDiv: css`
      width: ${width};
      display: flex;
      flex-direction: column;
      background-color: ${theme.colors.background.primary};
      button {
        border: none;
        outline: none;
        height: ${height};
        font-size: 20px;
        background-color: ${theme.colors.background.secondary};
        color: ${theme.colors.text.primary};
        border-right: 1px solid ${theme.components.panel.borderColor};
        border-bottom: 1px solid ${theme.components.panel.borderColor};
        transition: background-color 200ms;
        cursor: pointer;
        :hover {
          background-color: ${theme.colors.secondary.shade};
        }
      }
      input {
        height: ${height};
        border: 1px solid ${theme.colors.primary.main};
        :focus {
          outline: none;
        }
      }
    `,
  };
};
