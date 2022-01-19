import { logger } from '@percona/platform-core';
import React, { FC, useEffect, useMemo, useState, useRef, useCallback } from 'react';
import { createPortal } from 'react-dom';

import { Spinner, useTheme } from '@grafana/ui';

import { EmptyBlock } from '../shared/components/Elements/EmptyBlock';
import { ContentTab, TabbedContent, TabOrientation } from '../shared/components/Elements/TabbedContent';
import { TechnicalPreview } from '../shared/components/Elements/TechnicalPreview/TechnicalPreview';
import PageWrapper from '../shared/components/PageWrapper/PageWrapper';
import { useCancelToken } from '../shared/components/hooks/cancelToken.hook';

import { GET_SETTINGS_CANCEL_TOKEN, PAGE_MODEL } from './Settings.constants';
import { Messages } from './Settings.messages';
import { LoadingCallback, SettingsService } from './Settings.service';
import { getSettingsStyles } from './Settings.styles';
import { Settings, SettingsAPIChangePayload, TabKeys } from './Settings.types';
import { Advanced, AlertManager, Diagnostics, MetricsResolution, SSHKey } from './components';
import { Communication } from './components/Communication/Communication';
import { Platform } from './components/Platform/Platform';

export const SettingsPanel: FC = () => {
  const { path: basePath } = PAGE_MODEL;
  const [generateToken] = useCancelToken();

  const theme = useTheme();
  const [loading, setLoading] = useState(true);
  const [hasNoAccess, setHasNoAccess] = useState(false);
  const styles = getSettingsStyles(theme);
  const { metrics, advanced, ssh, alertManager, perconaPlatform, communication } = Messages.tabs;
  const [settings, setSettings] = useState<Settings>();
  const techPreviewRef = useRef<HTMLDivElement | null>(null);

  const updateSettings = useCallback(
    async (body: SettingsAPIChangePayload, callback: LoadingCallback, refresh?: boolean, onError = () => {}) => {
      // we save the test email here so that we can sent it all the way down to the form again after re-render
      // the field is deleted from the payload so as not to be sent to the API
      let password = '';
      let testEmail = '';

      if ('email_alerting_settings' in body) {
        password = body.email_alerting_settings.password || '';
        testEmail = body.email_alerting_settings.test_email || '';

        if (testEmail) {
          body.email_alerting_settings.test_email = undefined;
        }
      }
      const response = await SettingsService.setSettings(body, callback, generateToken(SET_SETTINGS_CANCEL_TOKEN));

      if (refresh && response) {
        window.location.reload();

        return;
      }

      if (response) {
        // password is not being returned by the API, hence this construction
        const newSettings: Settings = {
          ...response,
          alertingSettings: {
            ...response.alertingSettings,
            email: { ...response.alertingSettings.email, password, test_email: testEmail },
          },
        };
        setSettings(newSettings);
      } else {
        onError();
      }
    },
    [generateToken]
  );

  const getSettings = useCallback(async () => {
    try {
      setLoading(true);
      const settings = await SettingsService.getSettings(generateToken(GET_SETTINGS_CANCEL_TOKEN));
      setSettings(settings);
    } catch (e) {
      if (e.response?.status === 401) {
        setHasNoAccess(true);
      }
      logger.error(e);
    } finally {
      setLoading(false);
    }
  }, [generateToken]);

  const tabs: ContentTab[] = useMemo(
    (): ContentTab[] =>
      settings
        ? [
            {
              label: metrics,
              key: TabKeys.metrics,
              component: (
                <MetricsResolution metricsResolutions={settings.metricsResolutions} updateSettings={updateSettings} />
              ),
            },
            {
              label: advanced,
              key: TabKeys.advanced,
              component: (
                <Advanced
                  dataRetention={settings.dataRetention}
                  telemetryEnabled={!!settings.telemetryEnabled}
                  updatesDisabled={!!settings.updatesDisabled}
                  sttEnabled={!!settings.sttEnabled}
                  dbaasEnabled={!!settings.dbaasEnabled}
                  alertingEnabled={!!settings.alertingEnabled}
                  backupEnabled={!!settings.backupEnabled}
                  azureDiscoverEnabled={!!settings.azureDiscoverEnabled}
                  publicAddress={settings.publicAddress}
                  updateSettings={updateSettings}
                  sttCheckIntervals={settings.sttCheckIntervals}
                />
              ),
            },
            {
              label: ssh,
              key: TabKeys.ssh,
              component: <SSHKey sshKey={settings.sshKey || ''} updateSettings={updateSettings} />,
            },
            {
              label: alertManager,
              key: TabKeys.alertManager,
              component: (
                <AlertManager
                  alertManagerUrl={settings.alertManagerUrl || ''}
                  alertManagerRules={settings.alertManagerRules || ''}
                  updateSettings={updateSettings}
                />
              ),
            },
            {
              label: perconaPlatform,
              key: TabKeys.perconaPlatform,
              component: (
                <>
                  {techPreviewRef.current && createPortal(<TechnicalPreview />, techPreviewRef.current)}
                  <Platform isConnected={settings.isConnectedToPortal} getSettings={getSettings} />
                </>
              ),
            },
            {
              label: communication,
              key: TabKeys.communication,
              hidden: !settings?.alertingEnabled,
              component: (
                <Communication
                  alertingSettings={settings.alertingSettings}
                  alertingEnabled={!!settings.alertingEnabled}
                  updateSettings={updateSettings}
                />
              ),
            },
          ]
        : [],
    [settings, advanced, alertManager, communication, metrics, perconaPlatform, ssh, getSettings, updateSettings]
  );

  useEffect(() => {
    getSettings();
  }, [getSettings]);

  return (
    <PageWrapper pageModel={PAGE_MODEL}>
      <div ref={(e) => (techPreviewRef.current = e)} />
      <div className={styles.settingsWrapper}>
        {(loading || hasNoAccess) && (
          <div className={styles.emptyBlock}>
            <EmptyBlock dataTestId="empty-block">
              {loading ? <Spinner /> : hasNoAccess && <div data-testid="unauthorized">{Messages.unauthorized}</div>}
            </EmptyBlock>
          </div>
        )}
        {!loading && !hasNoAccess && (
          <>
            <TabbedContent
              className={styles.tabsWrapper}
              tabs={tabs}
              basePath={basePath}
              orientation={TabOrientation.Vertical}
              tabsdataTestId="settings-tabs"
              contentdataTestId="settings-tab-content"
              renderTab={({ Content }) => <Content className={styles.tabContentWrapper} />}
            />
            <Diagnostics />
          </>
        )}
      </div>
    </PageWrapper>
  );
};
export default SettingsPanel;
