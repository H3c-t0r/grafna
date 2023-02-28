import { DropzoneOptions } from 'react-dropzone';

import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceJsonData, DataSourceRef } from '@grafana/schema';

export interface DataSourceDrawerProps {
  datasources: Array<DataSourceInstanceSettings<DataSourceJsonData>>;
  onFileDrop?: () => void;
  onChange: (ds: string) => void;
  current: DataSourceInstanceSettings<DataSourceJsonData> | string | DataSourceRef | null | undefined;
  fileUploadOptions?: DropzoneOptions;
  recentlyUsed?: string[];
}

export interface DataSourceCardProps {
  onChange: (uid: string) => void;
  selected?: boolean;
  ds: DataSourceInstanceSettings<DataSourceJsonData>;
}

export interface PickerContentProps extends DataSourceDrawerProps {
  onDismiss: () => void;
}

export interface DataSourcePickerProps {
  drawer?: boolean;
  onChange: (ds: DataSourceInstanceSettings) => void;
  current: DataSourceRef | string | null; // uid
  tracing?: boolean;
  recentlyUsed?: string[];
  mixed?: boolean;
  dashboard?: boolean;
  metrics?: boolean;
  type?: string | string[];
  annotations?: boolean;
  variables?: boolean;
  alerting?: boolean;
  pluginId?: string;
  /** If true,we show only DSs with logs; and if true, pluginId shouldnt be passed in */
  logs?: boolean;
  // Does not set the default data source if there is no value.
  noDefault?: boolean;
  inputId?: string;
  filter?: (dataSource: DataSourceInstanceSettings) => boolean;
  onClear?: () => void;
  disabled?: boolean;
  fileUploadOptions?: DropzoneOptions;
}

export interface DataSourcePickerWithHistoryProps extends Omit<DataSourcePickerProps, 'recentlyUsed'> {
  key?: string;
}

export interface DataSourcePickerHistoryItem {
  lastUse: string;
  uid: string;
}
