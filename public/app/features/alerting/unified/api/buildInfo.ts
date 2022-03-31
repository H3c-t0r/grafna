import { getBackendSrv } from '@grafana/runtime';
import { PromApplication, PromBuildInfo, PromBuildInfoResponse } from 'app/types/unified-alerting-dto';
import { lastValueFrom } from 'rxjs';
import { isFetchError } from '../utils/alertmanager';
import { RULER_NOT_SUPPORTED_MSG } from '../utils/constants';
import { getDataSourceByName } from '../utils/datasource';
import { fetchRules } from './prometheus';
import { fetchRulerRulesGroup } from './ruler';

export async function fetchDataSourceBuildInfo(dsSettings: { url: string; name: string }): Promise<PromBuildInfo> {
  const { url, name } = dsSettings;

  const response = await lastValueFrom(
    getBackendSrv().fetch<PromBuildInfoResponse>({
      url: `${url}/api/v1/status/buildinfo`,
      showErrorAlert: false,
      showSuccessAlert: false,
    })
  ).catch((e) => {
    if ('status' in e && e.status === 404) {
      return null; // Cortex does not support buildinfo endpoint
    }

    throw e;
  });

  if (!response?.data.data) {
    const promRulesSupported = await hasPromRulesSupport(name);
    const rulerSupported = await hasRulerSupport(name);

    if (!promRulesSupported) {
      throw new Error(`Unable to fetch alert rules. Is the ${name} data source properly configured?`);
    }

    return {
      application: PromApplication.Cortex,
      features: {
        rulerApiEnabled: rulerSupported,
      },
    };
  }

  const { application, features } = response.data.data;

  return {
    application: PromApplication.Prometheus,
    features: {
      rulerApiEnabled: features?.ruler_config_api === 'true',
    },
  };
}

export async function fetchBuildInfo(dataSourceName: string): Promise<PromBuildInfo> {
  const dsConfig = getDataSourceByName(dataSourceName);
  if (!dsConfig) {
    throw new Error(`Cannot find data source configuration for ${dataSourceName}`);
  }
  const { url, name } = dsConfig;
  if (!url) {
    throw new Error(`The data souce url cannot be empty.`);
  }

  return fetchDataSourceBuildInfo({ name, url });
}

async function hasPromRulesSupport(dataSourceName: string) {
  try {
    await fetchRules(dataSourceName);
    return true;
  } catch (e) {
    return false;
  }
}

async function hasRulerSupport(dataSourceName: string) {
  try {
    await fetchRulerRulesGroup(dataSourceName, 'test', 'test');
    return true;
  } catch (e) {
    if (errorIndicatesMissingRulerSupport(e)) {
      return false;
    }
    throw e;
  }
}

function errorIndicatesMissingRulerSupport(error: any) {
  return (
    (isFetchError(error) &&
      (error.data.message?.includes('GetRuleGroup unsupported in rule local store') || // "local" rule storage
        error.data.message?.includes('page not found'))) || // ruler api disabled
    error.message?.includes('404 from rules config endpoint') || // ruler api disabled
    error.data.message?.includes(RULER_NOT_SUPPORTED_MSG) // ruler api not supported
  );
}
