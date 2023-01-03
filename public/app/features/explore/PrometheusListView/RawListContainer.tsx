import { css } from '@emotion/css';
import { cloneDeep } from 'lodash';
import React, { useEffect, useRef, useState } from 'react';
import { useWindowSize } from 'react-use';
import { VariableSizeList as List } from 'react-window';

import { DataFrame, Field } from '@grafana/data/src';
import { Button, stylesFactory } from '@grafana/ui/src';

import {
  getRawPrometheusListItemsFromDataFrame,
  RawPrometheusListItemEmptyValue,
} from '../utils/getRawPrometheusListItemsFromDataFrame';

import { ItemLabels } from './ItemLabels';
import RawListItem from './RawListItem';

export type instantQueryRawVirtualizedListData = { Value: string; __name__: string; [index: string]: string };

export interface RawListContainerProps {
  tableResult: DataFrame;
}

const getRawListContainerStyles = stylesFactory(() => {
  return {
    wrapper: css`
      height: 100%;
      overflow: scroll;
    `,
    mobileWrapper: css`
      height: 100%;
      overflow: scroll;
    `,
    header: css`
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 10px 0;
    `,
  };
});

const mobileWidthThreshold = 480;
const numberOfColumnsBeforeExpandedViewIsDefault = 2;

/**
 * The container that provides the virtualized list to the child components
 * @param props
 * @constructor
 */
const RawListContainer = (props: RawListContainerProps) => {
  const { tableResult } = props;
  const dataFrame = cloneDeep(tableResult);
  const styles = getRawListContainerStyles();
  const listRef = useRef<List | null>(null);

  const valueLabels = dataFrame.fields.filter((field) => field.name.includes('Value'));
  const items = getRawPrometheusListItemsFromDataFrame(dataFrame);
  const { width } = useWindowSize();
  const [isExpandedView, setIsExpandedView] = useState(
    width <= mobileWidthThreshold || valueLabels.length > numberOfColumnsBeforeExpandedViewIsDefault
  );

  const onContentClick = () => {
    setIsExpandedView(!isExpandedView);
  };

  useEffect(() => {
    // After the expanded view has updated, tell the list to re-render
    listRef.current?.resetAfterIndex(0, true);
  }, [isExpandedView]);

  const calculateInitialHeight = (length: number): number => {
    const maxListHeight = 600;
    const shortListLength = 10;

    if (length < shortListLength) {
      let sum = 0;
      for (let i = 0; i < length; i++) {
        sum += getListItemHeight(i, true);
      }

      return Math.min(maxListHeight, sum);
    }

    return maxListHeight;
  };

  const getListItemHeight = (itemIndex: number, isExpandedView: boolean) => {
    const singleLineHeight = 32;
    const additionalLineHeight = 22;
    if (!isExpandedView) {
      return singleLineHeight;
    }
    const item = items[itemIndex];

    // Height of 1.5 lines, plus the number of non-value attributes times the height of additional lines
    return 1.5 * singleLineHeight + (Object.keys(item).length - valueLabels.length) * additionalLineHeight;
  };

  return (
    <section>
      <header className={styles.header}>
        <Button
          aria-hidden="true"
          variant="secondary"
          role={'button'}
          title={`${isExpandedView ? 'Contract results' : 'Expand results'}`}
          icon={isExpandedView ? 'minus' : 'plus'}
          type="button"
          onClick={onContentClick}
        >
          {isExpandedView ? 'Contract' : 'Expand'}
        </Button>
        <div>Result series: {items.length}</div>
      </header>

      <div role={'table'}>
        {
          <>
            {/* Show the value headings above all the values, but only if we're in the contracted view */}
            {valueLabels.length > 1 && !isExpandedView && (
              <ItemLabels valueLabels={valueLabels} expanded={isExpandedView} />
            )}
            <List
              ref={listRef}
              itemCount={items.length}
              className={isExpandedView ? styles.mobileWrapper : styles.wrapper}
              itemSize={(index) => getListItemHeight(index, isExpandedView)}
              height={calculateInitialHeight(items.length)}
              width="100%"
            >
              {({ index, style }) => {
                let filteredValueLabels: Field[] | undefined;
                if (isExpandedView) {
                  filteredValueLabels = valueLabels.filter((valueLabel) => {
                    const itemWithValue = items[index][valueLabel.name];
                    return itemWithValue && itemWithValue !== RawPrometheusListItemEmptyValue;
                  });
                }

                return (
                  <div role="row" style={{ ...style, overflow: 'hidden' }}>
                    <RawListItem
                      isExpandedView={isExpandedView}
                      valueLabels={filteredValueLabels}
                      totalNumberOfValues={valueLabels.length}
                      listKey={items[index].__name__}
                      listItemData={items[index]}
                    />
                  </div>
                );
              }}
            </List>
          </>
        }
      </div>
    </section>
  );
};

export default RawListContainer;
