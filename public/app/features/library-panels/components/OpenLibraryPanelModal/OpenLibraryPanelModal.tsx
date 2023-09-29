import { css } from '@emotion/css';
import debounce from 'debounce-promise';
import React, { MouseEvent, useCallback, useEffect, useMemo, useState } from 'react';

import { GrafanaTheme2, SelectableValue, urlUtil } from '@grafana/data';
import { locationService } from '@grafana/runtime';
import { AsyncSelect, Button, Modal, useStyles2 } from '@grafana/ui';
import { t, Trans } from 'app/core/internationalization';

import { DashboardSearchItem } from '../../../search/types';
import { getConnectedDashboards, getLibraryPanelConnectedDashboards } from '../../state/api';
import { LibraryElementDTO } from '../../types';

export interface OpenLibraryPanelModalProps {
  onDismiss: () => void;
  libraryPanel: LibraryElementDTO;
}

export function OpenLibraryPanelModal({ libraryPanel, onDismiss }: OpenLibraryPanelModalProps): JSX.Element {
  const styles = useStyles2(getStyles);
  const [loading, setLoading] = useState(false);
  const [connected, setConnected] = useState(0);
  const [option, setOption] = useState<SelectableValue<DashboardSearchItem> | undefined>(undefined);
  useEffect(() => {
    const getConnected = async () => {
      const connectedDashboards = await getLibraryPanelConnectedDashboards(libraryPanel.uid);
      setConnected(connectedDashboards.length);
    };
    getConnected();
  }, [libraryPanel.uid]);
  const loadOptions = useCallback(
    (searchString: string) => loadOptionsAsync(libraryPanel.uid, searchString, setLoading),
    [libraryPanel.uid]
  );
  const debouncedLoadOptions = useMemo(() => debounce(loadOptions, 300, { leading: true }), [loadOptions]);
  const onViewPanel = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    locationService.push(urlUtil.renderUrl(`/d/${option?.value?.uid}`, {}));
  };

  return (
    <Modal
      title={t('library-panels.modal.title', 'View panel in dashboard')}
      onDismiss={onDismiss}
      onClickBackdrop={onDismiss}
      isOpen
    >
      <div className={styles.container}>
        {connected === 0 ? (
          <span>
            <Trans i18nKey={'library-panels.modal.panel-not-linked'}>Add the panel to a dashboard and retry.</Trans>
          </span>
        ) : null}
        {connected > 0 ? (
          <>
            <p>
              <Trans i18nKey={'library-panels.modal.contianer-prefix'}>This panel is being used in </Trans>
              <strong>
                {connected}{' '}
                {connected > 1 ? (
                  <Trans i18nKey={'library-panels.modal.contianer-dashboards'}>dashboards</Trans>
                ) : (
                  <Trans i18nKey={'library-panels.modal.contianer-dashboard'}>dashboard</Trans>
                )}
              </strong>
              {t('library-panels.modal.contianer-suffix', '. Please choose which dashboard to view the panel in:')}
            </p>
            <AsyncSelect
              isClearable
              isLoading={loading}
              defaultOptions={true}
              loadOptions={debouncedLoadOptions}
              onChange={setOption}
              placeholder={t('library-panels.select.placeholder', 'Start typing to search for dashboard')}
              noOptionsMessage={t('library-panels.select.no-dashboard-message', 'No dashboards found')}
            />
          </>
        ) : null}
      </div>
      <Modal.ButtonRow>
        <Button variant="secondary" onClick={onDismiss} fill="outline">
          <Trans i18nKey={'library-panels.modal.button-cancel'}>Cancel</Trans>
        </Button>
        <Button onClick={onViewPanel} disabled={!Boolean(option)}>
          {option
            ? t('library-panels.modal.button-view-panel1', 'View panel in {{label}}...', { label: option?.label })
            : t('library-panels.modal.button-view-panel2', 'View panel in dashboard...')}
        </Button>
      </Modal.ButtonRow>
    </Modal>
  );
}

async function loadOptionsAsync(uid: string, searchString: string, setLoading: (loading: boolean) => void) {
  setLoading(true);
  const searchHits = await getConnectedDashboards(uid);
  const options = searchHits
    .filter((d) => d.title.toLowerCase().includes(searchString.toLowerCase()))
    .map((d) => ({ label: d.title, value: d }));
  setLoading(false);

  return options;
}

function getStyles(theme: GrafanaTheme2) {
  return {
    container: css``,
  };
}
