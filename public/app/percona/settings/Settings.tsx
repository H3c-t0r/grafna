import { logger } from '@percona/platform-core';
import React, { FC, useEffect, useMemo, useState, useCallback } from 'react';

import { Spinner, useTheme } from '@grafana/ui';

import { ContentTab, TabbedContent, TabOrientation } from '../shared/components/Elements/TabbedContent';
import PageWrapper from '../shared/components/PageWrapper/PageWrapper';
import { useCancelToken } from '../shared/components/hooks/cancelToken.hook';

import { GET_SETTINGS_CANCEL_TOKEN, SET_SETTINGS_CANCEL_TOKEN, PAGE_MODEL } from './Settings.constants';
import { Messages } from './Settings.messages';
import { LoadingCallback, SettingsService } from './Settings.service';
import { getSettingsStyles } from './Settings.styles';
import { Settings, SettingsAPIChangePayload } from './Settings.types';
import { Advanced, AlertManager, Diagnostics, MetricsResolution, PlatformLogin, SSHKey } from './components';
import { Communication } from './components/Communication/Communication';

export const SettingsPanel: FC = () => {
  const { path: basePath } = PAGE_MODEL;
  const [generateToken] = useCancelToken();

  const theme = useTheme();
  const [loading, setLoading] = useState(true);
  const styles = getSettingsStyles(theme);
  const { metrics, advanced, ssh, alertManager, perconaPlatform, communication } = Messages.tabs;
  const [settings, setSettings] = useState<Settings>();

  const updateSettings = useCallback(
    async (body: SettingsAPIChangePayload, callback: LoadingCallback, refresh?: boolean) => {
      const response = await SettingsService.setSettings(body, callback, generateToken(SET_SETTINGS_CANCEL_TOKEN));
      const { email_alerting_settings: { password = '' } = {} } = body;

      if (refresh && response) {
        window.location.reload();

        return;
      }

      if (response) {
        // password is not being returned by the API, hence this construction
        const newSettings: Settings = {
          ...response,
          alertingSettings: { ...response.alertingSettings, email: { ...response.alertingSettings.email, password } },
        };
        setSettings(newSettings);
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
      logger.error(e);
    } finally {
      setLoading(false);
    }
  }, [generateToken]);

  const getSettings = async () => {
    try {
      setLoading(true);
      const settings = await SettingsService.getSettings(generateToken(GET_SETTINGS_CANCEL_TOKEN));
      setSettings(settings);
    } catch (e) {
      logger.error(e);
    } finally {
      setLoading(false);
    }
  };

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
              component: <PlatformLogin userEmail={settings.platformEmail} getSettings={getSettings} />,
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
      <div className={styles.settingsWrapper}>
        {loading ? (
          <Spinner />
        ) : (
          <TabbedContent
            className={styles.tabsWrapper}
            tabs={tabs}
            basePath={basePath}
            orientation={TabOrientation.Vertical}
            tabsDataQa="settings-tabs"
            contentDataQa="settings-tab-content"
            renderTab={({ Content }) => <Content className={styles.tabContentWrapper} />}
          ></TabbedContent>
        )}
        <Diagnostics />
      </div>
    </PageWrapper>
  );
};
export default SettingsPanel;
