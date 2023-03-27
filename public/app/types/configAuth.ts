import { Settings } from 'app/types';

export interface AuthConfigState {
  settings: Settings;
  providerStatuses: { [key: string]: AuthProviderStatus };
  updateError?: SettingsError;
  warning?: SettingsError;
}

export interface AuthProviderStatus {
  enabled: boolean;
  configured: boolean;
}

export interface SettingsError {
  message: string;
  errors: string[];
}
