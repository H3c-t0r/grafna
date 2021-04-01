import React, { FC, ReactElement } from 'react';
import { DisplayValue, Field, formattedValueToString } from '@grafana/data';

import { TableCellDisplayMode, TableCellProps } from './types';
import tinycolor from 'tinycolor2';
import { TableStyles } from './styles';
import { FilterActions } from './FilterActions';
import { getTextColorForBackground } from '../../utils';

export const DefaultCell: FC<TableCellProps> = (props) => {
  const { field, cell, tableStyles, cellProps } = props;

  const displayValue = field.display!(cell.value);

  let value: string | ReactElement;
  if (React.isValidElement(cell.value)) {
    value = cell.value;
  } else {
    value = formattedValueToString(displayValue);
  }

  const cellStyle = getCellStyle(tableStyles, field, displayValue);
  const showFilters = field.config.filterable;

  return (
    <div {...cellProps} className={cellStyle}>
      <div className={tableStyles.cellText}>{value}</div>
      {showFilters && cell.value !== undefined && <FilterActions {...props} />}
    </div>
  );
};

function getCellStyle(tableStyles: TableStyles, field: Field, displayValue: DisplayValue) {
  if (field.config.custom?.displayMode === TableCellDisplayMode.ColorText) {
    return tableStyles.buildCellContainerStyle(displayValue.color);
  }

  if (field.config.custom?.displayMode === TableCellDisplayMode.ColorBackground) {
    const themeFactor = tableStyles.theme.isDark ? 1 : -0.7;
    const bgColor2 = tinycolor(displayValue.color)
      .darken(10 * themeFactor)
      .spin(5)
      .toRgbString();

    const textColor = getTextColorForBackground(displayValue.color!);

    return tableStyles.buildCellContainerStyle(
      textColor,
      `linear-gradient(120deg, ${bgColor2}, ${displayValue.color})`
    );
  }

  return tableStyles.cellContainer;
}
