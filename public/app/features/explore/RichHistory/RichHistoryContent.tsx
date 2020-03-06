import React, { useEffect, useState } from 'react';
import { css } from 'emotion';

// Types
import { RichHistoryQuery } from 'app/types/explore';

// Utils
import { stylesFactory, useTheme } from '@grafana/ui';
import { GrafanaTheme, SelectableValue } from '@grafana/data';
import { getExploreDatasources } from '../state/selectors';

import { SortOrder } from '../../../core/utils/explore';
import {
  sortQueries,
  mapNumbertoTimeInSlider,
  createRetentionPeriodBoundary,
  createQueryHeading,
} from '../../../core/utils/richHistory';

// Components
import { RichHistoryCard } from './RichHistoryCard';
import { Forms, Slider } from '@grafana/ui';

const sortOrderOptions = [
  { label: 'Time ascending', value: SortOrder.Ascending },
  { label: 'Time descending', value: SortOrder.Descending },
  { label: 'Datasource A-Z', value: SortOrder.DatasourceAZ },
  { label: 'Datasource Z-A', value: SortOrder.DatasourceZA },
];

interface RichHistoryContentProps {
  queries: RichHistoryQuery[];
  sortOrder: SortOrder;
  activeDatasourceOnly: boolean;
  activeDatasourceInstance: string;
  datasourceFilters: SelectableValue[] | null;
  onlyStarred: boolean;
  retentionPeriod: number;
  onChangeSortOrder: (sortOrder: SortOrder) => void;
  onChangeRichHistoryProperty: (ts: number, property: string, updatedProperty?: string) => void;
  onSelectDatasourceFilters: (value: SelectableValue[] | null) => void;
}

const getStyles = stylesFactory((theme: GrafanaTheme, onlyStarred: boolean) => {
  const bgColor = theme.isLight ? theme.colors.gray5 : theme.colors.dark4;

  /* 134px is based on the width of the Query history tabs bar, so the content is aligned to right side of the tab */
  const cardWidth = onlyStarred ? '100%' : '100% - 134px';
  return {
    container: css`
      display: flex;
      .label-slider {
        font-size: ${theme.typography.size.sm};
        &:last-of-type {
          margin-top: ${theme.spacing.lg};
        }
        &:first-of-type {
          margin-top: ${theme.spacing.sm};
          font-weight: ${theme.typography.weight.semibold};
          margin-bottom: ${theme.spacing.xs};
        }
      }
    `,
    containerContent: css`
      width: calc(${cardWidth});
    `,
    containerSlider: css`
      width: 125px;
      margin-right: ${theme.spacing.sm};
      .slider {
        bottom: 10px;
        height: 200px;
        width: 125px;
        padding: ${theme.spacing.xs} 0;
      }
    `,
    slider: css`
      height: 300px;
      position: absolute;
    `,
    selectors: css`
      display: flex;
      justify-content: space-between;
    `,
    multiselect: css`
      width: 60%;
      .gf-form-select-box__multi-value {
        background-color: ${bgColor};
        padding: ${theme.spacing.xxs} ${theme.spacing.xs} ${theme.spacing.xxs} ${theme.spacing.sm};
        border-radius: ${theme.border.radius.sm};
      }
    `,
    sort: css`
      width: 170px;
    `,
    sessionName: css`
      display: flex;
      align-items: flex-start;
      justify-content: flex-start;
      margin-top: ${theme.spacing.lg};
      h4 {
        margin: 0 10px 0 0;
      }
    `,
    heading: css`
      font-size: ${theme.typography.heading.h4};
      margin: ${theme.spacing.md} ${theme.spacing.xxs} ${theme.spacing.sm} ${theme.spacing.xxs};
    `,
  };
});

export function RichHistoryContent(props: RichHistoryContentProps) {
  const {
    datasourceFilters,
    onSelectDatasourceFilters,
    queries,
    onlyStarred,
    onChangeSortOrder,
    sortOrder,
    onChangeRichHistoryProperty,
    activeDatasourceOnly,
    activeDatasourceInstance,
    retentionPeriod,
  } = props;

  const [sliderRetentionFilter, setSliderRetentionFilter] = useState<number[]>([0, retentionPeriod]);

  const theme = useTheme();
  const styles = getStyles(theme, onlyStarred);
  const exploreDatasources = getExploreDatasources()?.map(d => {
    return { value: d.value!, label: d.value!, imgUrl: d.meta.info.logos.small };
  });
  const listOfDatasourceFilters = datasourceFilters?.map(d => d.value);

  /* If user selects activeDatasourceOnly === true, set datasource filter to currently active datasource.
   *  Filtering based on datasource won't be available. Otherwise set to null, as filtering will be
   * available for user.
   */
  useEffect(() => {
    activeDatasourceOnly && activeDatasourceInstance
      ? onSelectDatasourceFilters([{ label: activeDatasourceInstance, value: activeDatasourceInstance }])
      : onSelectDatasourceFilters(null);
  }, [activeDatasourceInstance, activeDatasourceOnly]);

  /* selectorsAndSorter id used several times, therefore extracted to variable */
  const selectorsAndSorter = () => (
    <div className={styles.selectors}>
      {!activeDatasourceOnly && (
        <div className={styles.multiselect}>
          <Forms.Select
            isMulti={true}
            options={exploreDatasources}
            value={datasourceFilters}
            placeholder="Filter queries for specific datasources(s)"
            onChange={onSelectDatasourceFilters}
          />
        </div>
      )}
      <div className={styles.sort}>
        <Forms.Select
          options={sortOrderOptions}
          placeholder="Sort queries by"
          onChange={e => onChangeSortOrder(e.value as SortOrder)}
        />
      </div>
    </div>
  );

  /* If in Starred tab, render following content */
  if (onlyStarred) {
    const starredQueries = queries.filter(q => q.starred === true);
    const starredQueriesFilteredByDatasource = datasourceFilters
      ? starredQueries?.filter(q => listOfDatasourceFilters?.includes(q.datasourceName))
      : starredQueries;
    const sortedStarredQueries = sortQueries(starredQueriesFilteredByDatasource, sortOrder);

    return (
      <div className={styles.container}>
        <div className={styles.containerContent}>
          {selectorsAndSorter()}
          {sortedStarredQueries.map(q => {
            return <RichHistoryCard query={q} key={q.ts} onChangeRichHistoryProperty={onChangeRichHistoryProperty} />;
          })}
        </div>
      </div>
    );

    /* If in History tab, render following content */
  } else {
    const filteredQueriesByDatasource = datasourceFilters
      ? queries?.filter(q => listOfDatasourceFilters?.includes(q.datasourceName))
      : queries;
    const sortedQueries = sortQueries(filteredQueriesByDatasource, sortOrder);
    const queriesWithinSelectedTimeline = sortedQueries?.filter(
      q =>
        q.ts < createRetentionPeriodBoundary(sliderRetentionFilter[0], true) &&
        q.ts > createRetentionPeriodBoundary(sliderRetentionFilter[1], false)
    );
    return (
      <div className={styles.container}>
        <div className={styles.containerSlider}>
          <div className={styles.slider}>
            <div className="label-slider">
              Filter history <br />
              between
            </div>
            <div className="label-slider">{mapNumbertoTimeInSlider(sliderRetentionFilter[0])}</div>
            <div className="slider">
              <Slider
                tooltipAlwaysVisible={false}
                min={0}
                max={retentionPeriod}
                value={sliderRetentionFilter}
                orientation="vertical"
                formatTooltipResult={mapNumbertoTimeInSlider}
                reverse={true}
                onAfterChange={setSliderRetentionFilter}
              />
            </div>
            <div className="label-slider">{mapNumbertoTimeInSlider(sliderRetentionFilter[1])}</div>
          </div>
        </div>

        <div className={styles.containerContent}>
          {selectorsAndSorter()}
          {queriesWithinSelectedTimeline.map((q, i) => {
            let currentHeading = createQueryHeading(q, sortOrder);
            let previousHeading = i > 0 ? createQueryHeading(queriesWithinSelectedTimeline[i - 1], sortOrder) : '';
            if (currentHeading !== previousHeading) {
              return (
                <div key={q.ts}>
                  <div className={styles.heading}>{currentHeading}</div>
                  <RichHistoryCard query={q} key={q.ts} onChangeRichHistoryProperty={onChangeRichHistoryProperty} />
                </div>
              );
            } else {
              return <RichHistoryCard query={q} key={q.ts} onChangeRichHistoryProperty={onChangeRichHistoryProperty} />;
            }
          })}
        </div>
      </div>
    );
  }
}
