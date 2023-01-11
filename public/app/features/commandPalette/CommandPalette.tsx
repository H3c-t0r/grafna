import { css } from '@emotion/css';
import { useDialog } from '@react-aria/dialog';
import { FocusScope } from '@react-aria/focus';
import { useOverlay } from '@react-aria/overlays';
import debounce from 'debounce-promise';
import {
  KBarAnimator,
  KBarPortal,
  KBarPositioner,
  KBarResults,
  KBarSearch,
  useMatches,
  VisualState,
  useRegisterActions,
  useKBar,
  ActionImpl,
} from 'kbar';
import React, { useEffect, useRef, useState } from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { useStyles2 } from '@grafana/ui';
import { t } from 'app/core/internationalization';

import { ResultItem } from './ResultItem';
import { getDashboardSearchResultActions } from './actions/dashboardActions';
import useActions from './actions/useActions';

const debouncedDashboardSearch = debounce(getDashboardSearchResultActions, 200);

export const CommandPalette = () => {
  const styles = useStyles2(getSearchStyles);

  const { query, showing, searchQuery } = useKBar((state) => ({
    showing: state.visualState === VisualState.showing,
    searchQuery: state.searchQuery,
  }));
  const [dashboardResults, setDashboardResults] = useState<ActionImpl[]>([]);

  // Hit dashboards API
  useEffect(() => {
    if (showing && searchQuery.length > 0) {
      debouncedDashboardSearch(searchQuery).then((resultActions) => {
        setDashboardResults(resultActions);
      });
    }
  }, [showing, searchQuery]);

  const actions = useActions(searchQuery, showing);
  useRegisterActions(actions, [actions]);

  const ref = useRef<HTMLDivElement>(null);
  const { overlayProps } = useOverlay(
    { isOpen: showing, onClose: () => query.setVisualState(VisualState.animatingOut) },
    ref
  );
  const { dialogProps } = useDialog({}, ref);

  // Report interaction when opened
  useEffect(() => {
    showing && reportInteraction('command_palette_opened');
  }, [showing]);

  return actions.length > 0 ? (
    <KBarPortal>
      <KBarPositioner className={styles.positioner}>
        <KBarAnimator className={styles.animator}>
          <FocusScope contain autoFocus restoreFocus>
            <div {...overlayProps} {...dialogProps}>
              <KBarSearch className={styles.search} />
              <RenderResults dashboardResults={dashboardResults} />
            </div>
          </FocusScope>
        </KBarAnimator>
      </KBarPositioner>
    </KBarPortal>
  ) : null;
};

interface RenderResultsProps {
  dashboardResults: ActionImpl[];
}

const RenderResults = ({ dashboardResults }: RenderResultsProps) => {
  const { results, rootActionId } = useMatches();
  const styles = useStyles2(getSearchStyles);
  const dashboardsSectionTitle = t('command-palette.section.dashboard-search-results', 'Dashboards');
  const dashboards = dashboardResults.length > 0 ? [dashboardsSectionTitle, ...dashboardResults] : [];

  return (
    <div className={styles.resultsContainer}>
      <KBarResults
        items={[...results, ...dashboards]}
        onRender={({ item, active }) =>
          typeof item === 'string' ? (
            <div className={styles.sectionHeader}>{item}</div>
          ) : (
            <ResultItem action={item} active={active} currentRootActionId={rootActionId!} />
          )
        }
      />
    </div>
  );
};

const getSearchStyles = (theme: GrafanaTheme2) => ({
  positioner: css({
    zIndex: theme.zIndex.portal,
    marginTop: '0px',
    '&::before': {
      content: '""',
      position: 'fixed',
      top: 0,
      right: 0,
      bottom: 0,
      left: 0,
      background: theme.components.overlay.background,
      backdropFilter: 'blur(1px)',
    },
  }),
  animator: css({
    maxWidth: theme.breakpoints.values.sm, // supposed to be 600...
    width: '100%',
    background: theme.colors.background.canvas,
    color: theme.colors.text.primary,
    borderRadius: theme.shape.borderRadius(4),
    overflow: 'hidden',
    boxShadow: theme.shadows.z3,
  }),
  search: css({
    padding: theme.spacing(2, 3),
    fontSize: theme.typography.fontSize,
    width: '100%',
    boxSizing: 'border-box',
    outline: 'none',
    border: 'none',
    background: theme.colors.background.canvas,
    color: theme.colors.text.primary,
    borderBottom: `1px solid ${theme.colors.border.weak}`,
  }),
  sectionHeader: css({
    padding: theme.spacing(1, 2),
    fontSize: theme.typography.h6.fontSize,
    fontWeight: theme.typography.body.fontWeight,
    color: theme.colors.text.secondary,
  }),
  resultsContainer: css({
    padding: theme.spacing(2, 0),
  }),
});
