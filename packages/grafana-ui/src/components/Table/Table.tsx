import React, { FC, memo, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import {
  Cell,
  Column,
  TableState,
  useAbsoluteLayout,
  useFilters,
  usePagination,
  useResizeColumns,
  useSortBy,
  useTable,
} from 'react-table';
import { FixedSizeList } from 'react-window';

import { DataFrame, getFieldDisplayName } from '@grafana/data';

import { useStyles2 } from '../../themes';
import { Button } from '../Button';
import { CustomScrollbar } from '../CustomScrollbar/CustomScrollbar';
import { Pagination } from '../Pagination/Pagination';

import { FooterRow } from './FooterRow';
import { HeaderRow } from './HeaderRow';
import { TableCell } from './TableCell';
import { getTableStyles } from './styles';
import {
  TableColumnResizeActionCallback,
  TableFilterActionCallback,
  FooterItem,
  TableSortByActionCallback,
  TableSortByFieldState,
} from './types';
import { getColumns, sortCaseInsensitive, sortNumber } from './utils';

const COLUMN_MIN_WIDTH = 150;

export interface Props {
  ariaLabel?: string;
  data: DataFrame;
  width: number;
  height: number;
  /** Minimal column width specified in pixels */
  columnMinWidth?: number;
  noHeader?: boolean;
  showTypeIcons?: boolean;
  resizable?: boolean;
  initialSortBy?: TableSortByFieldState[];
  onColumnResize?: TableColumnResizeActionCallback;
  onSortByChange?: TableSortByActionCallback;
  onCellFilterAdded?: TableFilterActionCallback;
  footerValues?: FooterItem[];
  enablePagination?: boolean;
}

function useTableStateReducer({ onColumnResize, onSortByChange, data }: Props) {
  return useCallback(
    (newState: TableState, action: any) => {
      switch (action.type) {
        case 'columnDoneResizing':
          if (onColumnResize) {
            const info = (newState.columnResizing.headerIdWidths as any)[0];
            const columnIdString = info[0];
            const fieldIndex = parseInt(columnIdString, 10);
            const width = Math.round(newState.columnResizing.columnWidths[columnIdString] as number);

            const field = data.fields[fieldIndex];
            if (!field) {
              return newState;
            }

            const fieldDisplayName = getFieldDisplayName(field, data);
            onColumnResize(fieldDisplayName, width);
          }
        case 'toggleSortBy':
          if (onSortByChange) {
            const sortByFields: TableSortByFieldState[] = [];

            for (const sortItem of newState.sortBy) {
              const field = data.fields[parseInt(sortItem.id, 10)];
              if (!field) {
                continue;
              }

              sortByFields.push({
                displayName: getFieldDisplayName(field, data),
                desc: sortItem.desc,
              });
            }

            onSortByChange(sortByFields);
          }
          break;
      }

      return newState;
    },
    [data, onColumnResize, onSortByChange]
  );
}

function getInitialState(initialSortBy: Props['initialSortBy'], columns: Column[]): Partial<TableState> {
  const state: Partial<TableState> = {};

  if (initialSortBy) {
    state.sortBy = [];

    for (const sortBy of initialSortBy) {
      for (const col of columns) {
        if (col.Header === sortBy.displayName) {
          state.sortBy.push({ id: col.id as string, desc: sortBy.desc });
        }
      }
    }
  }

  return state;
}

export const Table: FC<Props> = memo((props: Props) => {
  const {
    ariaLabel,
    data,
    height,
    onCellFilterAdded,
    width,
    columnMinWidth = COLUMN_MIN_WIDTH,
    noHeader,
    resizable = true,
    initialSortBy,
    footerValues,
    showTypeIcons,
    enablePagination,
  } = props;

  const [showSpans, setShowSpans] = useState(false);

  const listRef = useRef<FixedSizeList>(null);
  const tableDivRef = useRef<HTMLDivElement>(null);
  const fixedSizeListScrollbarRef = useRef<HTMLDivElement>(null);
  const tableStyles = useStyles2(getTableStyles);
  const headerHeight = noHeader ? 0 : tableStyles.cellHeight;

  const footerHeight = useMemo(() => {
    const EXTENDED_ROW_HEIGHT = 33;
    let length = 0;

    if (!footerValues) {
      return 0;
    }

    for (const fv of footerValues) {
      if (Array.isArray(fv) && fv.length > length) {
        length = fv.length;
      }
    }

    if (length > 1) {
      return EXTENDED_ROW_HEIGHT * length;
    }

    return EXTENDED_ROW_HEIGHT;
  }, [footerValues]);

  // React table data array. This data acts just like a dummy array to let react-table know how many rows exist
  // The cells use the field to look up values
  const memoizedData = useMemo(() => {
    if (!data.fields.length) {
      return [];
    }
    // as we only use this to fake the length of our data set for react-table we need to make sure we always return an array
    // filled with values at each index otherwise we'll end up trying to call accessRow for null|undefined value in
    // https://github.com/tannerlinsley/react-table/blob/7be2fc9d8b5e223fc998af88865ae86a88792fdb/src/hooks/useTable.js#L585
    return Array(data.length).fill(0);
  }, [data]);

  // React-table column definitions
  const memoizedColumns = useMemo(
    () => getColumns(data, width, columnMinWidth, showSpans, footerValues),
    [data, width, columnMinWidth, showSpans, footerValues]
  );

  // Internal react table state reducer
  const stateReducer = useTableStateReducer(props);

  const options: any = useMemo(
    () => ({
      columns: memoizedColumns,
      data: memoizedData,
      disableResizing: !resizable,
      stateReducer: stateReducer,
      initialState: getInitialState(initialSortBy, memoizedColumns),
      autoResetFilters: false,
      sortTypes: {
        number: sortNumber, // the builtin number type on react-table does not handle NaN values
        'alphanumeric-insensitive': sortCaseInsensitive, // should be replace with the builtin string when react-table is upgraded, see https://github.com/tannerlinsley/react-table/pull/3235
      },
    }),
    [initialSortBy, memoizedColumns, memoizedData, resizable, stateReducer]
  );

  const {
    getTableProps,
    headerGroups,
    rows,
    prepareRow,
    totalColumnsWidth,
    footerGroups,
    page,
    state,
    gotoPage,
    setPageSize,
    pageOptions,
  } = useTable(options, useFilters, useSortBy, usePagination, useAbsoluteLayout, useResizeColumns);

  let listHeight = height - (headerHeight + footerHeight);

  if (enablePagination) {
    listHeight -= tableStyles.cellHeight;
  }
  const pageSize = Math.round(listHeight / tableStyles.cellHeight) - 1;
  useEffect(() => {
    // Don't update the page size if it is less than 1
    if (pageSize <= 0) {
      return;
    }
    setPageSize(pageSize);
  }, [pageSize, setPageSize]);

  useEffect(() => {
    // To have the custom vertical scrollbar always visible (https://github.com/grafana/grafana/issues/52136),
    // we need to bring the element from the FixedSizeList scope to the outer Table container scope,
    // because the FixedSizeList scope has overflow. By moving scrollbar to container scope we will have
    // it always visible since the entire width is in view.

    // Select the scrollbar element from the FixedSizeList scope
    const listVerticalScrollbarHTML = (fixedSizeListScrollbarRef.current as HTMLDivElement)?.querySelector(
      '.track-vertical'
    );

    // Select Table custom scrollbars
    const tableScrollbarView = (tableDivRef.current as HTMLDivElement)?.firstChild;

    //If they exists, move the scrollbar element to the Table container scope
    if (tableScrollbarView && listVerticalScrollbarHTML) {
      listVerticalScrollbarHTML?.remove();
      (tableScrollbarView as HTMLDivElement).querySelector(':scope > .track-vertical')?.remove();

      (tableScrollbarView as HTMLDivElement).append(listVerticalScrollbarHTML as Node);
    }
  });

  const RenderRow = React.useCallback(
    ({ index: rowIndex, style }) => {
      let row = rows[rowIndex];
      if (enablePagination) {
        row = page[rowIndex];
      }
      prepareRow(row);
      if (!showSpans && row.values[0]) {
        return (
          <div {...row.getRowProps({ style })} className={tableStyles.row}>
            {row.cells.map((cell: Cell, index: number) => (
              <TableCell
                key={index}
                tableStyles={tableStyles}
                cell={cell}
                onCellFilterAdded={onCellFilterAdded}
                columnIndex={index}
                columnCount={row.cells.length}
              />
            ))}
          </div>
        );
      } else if (showSpans) {
        return (
          <div {...row.getRowProps({ style })} className={tableStyles.row}>
            {row.cells.map((cell: Cell, index: number) => (
              <TableCell
                key={index}
                tableStyles={tableStyles}
                cell={cell}
                onCellFilterAdded={onCellFilterAdded}
                columnIndex={index}
                columnCount={row.cells.length}
              />
            ))}
          </div>
        );
      }
      return null;
    },
    [onCellFilterAdded, page, enablePagination, prepareRow, rows, showSpans, tableStyles]
  );

  const onNavigate = useCallback(
    (toPage: number) => {
      gotoPage(toPage - 1);
    },
    [gotoPage]
  );

  const itemCount = enablePagination ? page.length : rows.length;
  let paginationEl = null;
  if (enablePagination) {
    const itemsRangeStart = state.pageIndex * state.pageSize + 1;
    let itemsRangeEnd = itemsRangeStart + state.pageSize - 1;
    const isSmall = width < 550;
    if (itemsRangeEnd > data.length) {
      itemsRangeEnd = data.length;
    }
    paginationEl = (
      <div className={tableStyles.paginationWrapper}>
        {isSmall ? null : <div className={tableStyles.paginationItem} />}
        <div className={tableStyles.paginationCenterItem}>
          <Pagination
            currentPage={state.pageIndex + 1}
            numberOfPages={pageOptions.length}
            showSmallVersion={isSmall}
            onNavigate={onNavigate}
          />
        </div>
        {isSmall ? null : (
          <div className={tableStyles.paginationSummary}>
            {itemsRangeStart} - {itemsRangeEnd} of {data.length} rows
          </div>
        )}
      </div>
    );
  }

  const handleScroll: React.UIEventHandler = (event) => {
    const { scrollTop } = event.target as HTMLDivElement;

    if (listRef.current !== null) {
      listRef.current.scrollTo(scrollTop);
    }
  };

  return (
    <div {...getTableProps()} className={tableStyles.table} aria-label={ariaLabel} role="table" ref={tableDivRef}>
      {data.fields[1].name === 'spanID' ? ( // Only show the button if we're in the TraceQL tab, not when using Search.
        <div>
          <Button size="sm" onClick={() => setShowSpans((val) => !val)}>
            {showSpans ? 'Hide spans' : 'Show spans'}
          </Button>
          <br /> <br />
        </div>
      ) : null}

      <CustomScrollbar hideVerticalTrack={true}>
        <div className={tableStyles.tableContentWrapper(totalColumnsWidth)}>
          {!noHeader && <HeaderRow headerGroups={headerGroups} showTypeIcons={showTypeIcons} />}
          {itemCount > 0 ? (
            <div ref={fixedSizeListScrollbarRef}>
              <CustomScrollbar onScroll={handleScroll} hideHorizontalTrack={true}>
                <FixedSizeList
                  height={listHeight}
                  itemCount={itemCount}
                  itemSize={tableStyles.rowHeight}
                  width={'100%'}
                  ref={listRef}
                  style={{ overflow: undefined }}
                >
                  {RenderRow}
                </FixedSizeList>
              </CustomScrollbar>
            </div>
          ) : (
            <div style={{ height: height - headerHeight }} className={tableStyles.noData}>
              No data
            </div>
          )}
          {footerValues && (
            <FooterRow
              height={footerHeight}
              isPaginationVisible={Boolean(enablePagination)}
              footerValues={footerValues}
              footerGroups={footerGroups}
              totalColumnsWidth={totalColumnsWidth}
            />
          )}
        </div>
      </CustomScrollbar>
      {paginationEl}
    </div>
  );
});

Table.displayName = 'Table';
