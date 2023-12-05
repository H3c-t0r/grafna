import { isEqual, map } from 'lodash';
import React, { SyntheticEvent, useCallback, useEffect, useState } from 'react';

import { CoreApp, LoadingState, SelectableValue } from '@grafana/data';
import { selectors } from '@grafana/e2e-selectors';
import { EditorHeader, EditorRows, FlexItem, Space } from '@grafana/experimental';
import { config, reportInteraction } from '@grafana/runtime';
import { Button, ConfirmModal, Drawer } from '@grafana/ui';

import { PromQueryEditorProps } from '../../components/types';
import { PromQueryFormat } from '../../dataquery.gen';
import { PromQuery } from '../../types';
import { QueryPatternsModal } from '../QueryPatternsModal';
import { buildVisualQueryFromString } from '../parsing';
import { QueryEditorModeToggle } from '../shared/QueryEditorModeToggle';
import { QueryHeaderSwitch } from '../shared/QueryHeaderSwitch';
import { promQueryEditorExplainKey, useFlag } from '../shared/hooks/useFlag';
import { QueryEditorMode } from '../shared/types';
import { changeEditorMode, getQueryWithDefaults } from '../state';

import { PromQueryBuilderContainer } from './PromQueryBuilderContainer';
import { PromQueryBuilderOptions } from './PromQueryBuilderOptions';
import { PromQueryCodeEditor } from './PromQueryCodeEditor';
import { WizarDS } from './wizarDS/WizarDS';

export const FORMAT_OPTIONS: Array<SelectableValue<PromQueryFormat>> = [
  { label: 'Time series', value: 'time_series' },
  { label: 'Table', value: 'table' },
  { label: 'Heatmap', value: 'heatmap' },
];

export const INTERVAL_FACTOR_OPTIONS: Array<SelectableValue<number>> = map([1, 2, 3, 4, 5, 10], (value: number) => ({
  value,
  label: '1/' + value,
}));

type Props = PromQueryEditorProps;

export const PromQueryEditorSelector = React.memo<Props>((props) => {
  const {
    onChange,
    onRunQuery,
    data,
    app,
    onAddQuery,
    datasource: { defaultEditor },
    queries,
  } = props;

  const [parseModalOpen, setParseModalOpen] = useState(false);
  const [queryPatternsModalOpen, setQueryPatternsModalOpen] = useState(false);
  const [dataIsStale, setDataIsStale] = useState(false);
  const [wizarDSDrawerOpen, setWizarDSDrawerOpen] = useState(false);
  const { flag: explain, setFlag: setExplain } = useFlag(promQueryEditorExplainKey);

  const query = getQueryWithDefaults(props.query, app, defaultEditor);
  // This should be filled in from the defaults by now.
  const editorMode = query.editorMode!;

  const onEditorModeChange = useCallback(
    (newMetricEditorMode: QueryEditorMode) => {
      reportInteraction('user_grafana_prometheus_editor_mode_clicked', {
        newEditor: newMetricEditorMode,
        previousEditor: query.editorMode ?? '',
        newQuery: !query.expr,
        app: app ?? '',
      });

      if (newMetricEditorMode === QueryEditorMode.Builder) {
        const result = buildVisualQueryFromString(query.expr || '');
        // If there are errors, give user a chance to decide if they want to go to builder as that can lose some data.
        if (result.errors.length) {
          setParseModalOpen(true);
          return;
        }
      }
      changeEditorMode(query, newMetricEditorMode, onChange);
    },
    [onChange, query, app]
  );

  useEffect(() => {
    setDataIsStale(false);
  }, [data]);

  const onChangeInternal = (query: PromQuery) => {
    if (!isEqual(query, props.query)) {
      setDataIsStale(true);
    }
    onChange(query);
  };

  const onShowExplainChange = (e: SyntheticEvent<HTMLInputElement>) => {
    setExplain(e.currentTarget.checked);
  };

  const wizarDSToggle = config.featureToggles.wizarDSToggle;

  return (
    <>
      <ConfirmModal
        isOpen={parseModalOpen}
        title="Parsing error: Switch to the builder mode?"
        body="There is a syntax error, or the query structure cannot be visualized when switching to the builder mode. Parts of the query may be lost. "
        confirmText="Continue"
        onConfirm={() => {
          changeEditorMode(query, QueryEditorMode.Builder, onChange);
          setParseModalOpen(false);
        }}
        onDismiss={() => setParseModalOpen(false)}
      />
      <QueryPatternsModal
        isOpen={queryPatternsModalOpen}
        onClose={() => setQueryPatternsModalOpen(false)}
        query={query}
        queries={queries}
        app={app}
        onChange={onChange}
        onAddQuery={onAddQuery}
      />
      {wizarDSToggle && wizarDSDrawerOpen && (
        <Drawer closeOnMaskClick={false} onClose={() => setWizarDSDrawerOpen(false)}>
          <WizarDS
            // remove all references to the query
            query={{ metric: '', labels: [], operations: [] }}
            closeDrawer={() => setWizarDSDrawerOpen(false)}
            datasource={props.datasource}
          />
        </Drawer>
      )}
      <EditorHeader>
        <Button
          aria-label={selectors.components.QueryBuilder.queryPatterns}
          variant="secondary"
          size="sm"
          onClick={() => setQueryPatternsModalOpen((prevValue) => !prevValue)}
        >
          Kick start your query
        </Button>
        <Button variant="secondary" size="sm" onClick={() => setWizarDSDrawerOpen((prevValue) => !prevValue)}>
          Take me on a magical journey
        </Button>
        <QueryHeaderSwitch label="Explain" value={explain} onChange={onShowExplainChange} />
        <FlexItem grow={1} />
        {app !== CoreApp.Explore && app !== CoreApp.Correlations && (
          <Button
            variant={dataIsStale ? 'primary' : 'secondary'}
            size="sm"
            onClick={onRunQuery}
            icon={data?.state === LoadingState.Loading ? 'spinner' : undefined}
            disabled={data?.state === LoadingState.Loading}
          >
            Run queries
          </Button>
        )}
        <QueryEditorModeToggle mode={editorMode} onChange={onEditorModeChange} />
      </EditorHeader>
      <Space v={0.5} />
      <EditorRows>
        {editorMode === QueryEditorMode.Code && (
          <PromQueryCodeEditor {...props} query={query} showExplain={explain} onChange={onChangeInternal} />
        )}
        {editorMode === QueryEditorMode.Builder && (
          <PromQueryBuilderContainer
            query={query}
            datasource={props.datasource}
            onChange={onChangeInternal}
            onRunQuery={props.onRunQuery}
            data={data}
            showExplain={explain}
          />
        )}
        <PromQueryBuilderOptions query={query} app={props.app} onChange={onChange} onRunQuery={onRunQuery} />
      </EditorRows>
    </>
  );
});

PromQueryEditorSelector.displayName = 'PromQueryEditorSelector';
