import { AzureCredentials, AzureDataSourceSettings } from './types';

function getSecret(options: AzureDataSourceSettings): undefined | string | object {
  if (options.secureJsonFields.clientSecret) {
    // The secret is concealed on server
    return {};
  } else {
    const secret = options.secureJsonData?.clientSecret;
    return typeof secret === 'string' && secret.length > 0 ? secret : undefined;
  }
}

function getLogAnalyticsSecret(options: AzureDataSourceSettings): undefined | string | object {
  if (options.secureJsonFields.logAnalyticsClientSecret) {
    // The secret is concealed on server
    return {};
  } else {
    const secret = options.secureJsonData?.logAnalyticsClientSecret;
    return typeof secret === 'string' && secret.length > 0 ? secret : undefined;
  }
}

export function isLogAnalyticsSameAs(options: AzureDataSourceSettings): boolean {
  return typeof options.jsonData.azureLogAnalyticsSameAs !== 'boolean' || options.jsonData.azureLogAnalyticsSameAs;
}

export function isCredentialsComplete(credentials: AzureCredentials) {
  return !!(credentials.tenantId && credentials.clientId && credentials.clientSecret);
}

export function getCredentials(options: AzureDataSourceSettings): AzureCredentials {
  return {
    azureCloud: options.jsonData.cloudName || 'azuremonitor',
    tenantId: options.jsonData.tenantId,
    clientId: options.jsonData.clientId,
    clientSecret: getSecret(options),
  };
}

export function getLogAnalyticsCredentials(options: AzureDataSourceSettings): AzureCredentials | undefined {
  if (isLogAnalyticsSameAs(options)) {
    return undefined;
  }

  return {
    tenantId: options.jsonData.logAnalyticsTenantId,
    clientId: options.jsonData.logAnalyticsClientId,
    clientSecret: getLogAnalyticsSecret(options),
  };
}

export function updateCredentials(
  options: AzureDataSourceSettings,
  credentials: AzureCredentials
): AzureDataSourceSettings {
  options = {
    ...options,
    jsonData: {
      ...options.jsonData,
      cloudName: credentials.azureCloud || 'azuremonitor',
      tenantId: credentials.tenantId,
      clientId: credentials.clientId,
    },
    secureJsonData: {
      ...options.secureJsonData,
      clientSecret:
        typeof credentials.clientSecret === 'string' && credentials.clientSecret.length > 0
          ? credentials.clientSecret
          : undefined,
    },
    secureJsonFields: {
      ...options.secureJsonFields,
      clientSecret: typeof credentials.clientSecret === 'object',
    },
  };

  if (isLogAnalyticsSameAs(options)) {
    options = updateLogAnalyticsCredentials(options, credentials);
  }

  return options;
}

export function updateLogAnalyticsCredentials(
  options: AzureDataSourceSettings,
  credentials: AzureCredentials
): AzureDataSourceSettings {
  options = {
    ...options,
    jsonData: {
      ...options.jsonData,
      logAnalyticsTenantId: credentials.tenantId,
      logAnalyticsClientId: credentials.clientId,
    },
    secureJsonData: {
      ...options.secureJsonData,
      logAnalyticsClientSecret:
        typeof credentials.clientSecret === 'string' && credentials.clientSecret.length > 0
          ? credentials.clientSecret
          : undefined,
    },
    secureJsonFields: {
      ...options.secureJsonFields,
      logAnalyticsClientSecret: typeof credentials.clientSecret === 'object',
    },
  };

  return options;
}

export function updateLogAnalyticsSameAs(options: AzureDataSourceSettings, sameAs: boolean): AzureDataSourceSettings {
  if (sameAs !== isLogAnalyticsSameAs(options)) {
    // Update the 'Same As' switch
    options = {
      ...options,
      jsonData: {
        ...options.jsonData,
        azureLogAnalyticsSameAs: sameAs,
      },
    };

    if (sameAs) {
      // Get the primary credentials
      let credentials = getCredentials(options);

      // Check whether the client secret is concealed
      if (typeof credentials.clientSecret === 'object') {
        // Log Analytics credentials need to be synchronized but the client secret is concealed,
        // so we have to reset the primary client secret to ensure that user enters a new secret
        credentials.clientSecret = undefined;
        options = updateCredentials(options, credentials);
      }

      // Synchronize the Log Analytics credentials with primary credentials
      options = updateLogAnalyticsCredentials(options, credentials);
    }
  }

  return options;
}
