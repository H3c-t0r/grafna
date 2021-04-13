// Prometheus API DTOs, possibly to be autogenerated from openapi spec in the near future

export type Labels = Record<string, string>;
export type Annotations = Record<string, string>;

export enum PromAlertingRuleState {
  Firing = 'firing',
  Inactive = 'inactive',
  Pending = 'pending',
}

export enum PromRuleType {
  Alerting = 'alerting',
  Recording = 'recording',
}

interface PromRuleDTOBase {
  health: string;
  name: string;
  query: string; // expr
  evaluationTime?: number;
  lastEvaluation?: string;
  lastError?: string;
}

export interface PromAlertingRuleDTO extends PromRuleDTOBase {
  alerts: Array<{
    labels: Labels;
    annotations: Annotations;
    state: Exclude<PromAlertingRuleState, PromAlertingRuleState.Inactive>;
    activeAt: string;
    value: string;
  }>;
  labels: Labels;
  annotations: Annotations;
  duration?: number; // for
  state: PromAlertingRuleState;
  type: PromRuleType.Alerting;
}

export interface PromRecordingRuleDTO extends PromRuleDTOBase {
  health: string;
  name: string;
  query: string; // expr
  type: PromRuleType.Recording;
  labels?: Labels;
}

export type PromRuleDTO = PromAlertingRuleDTO | PromRecordingRuleDTO;

export interface PromRuleGroupDTO {
  name: string;
  file: string;
  rules: PromRuleDTO[];
  interval: number;

  evaluationTime?: number; // these 2 are not in older prometheus payloads
  lastEvaluation?: string;
}

export interface PromResponse<T> {
  status: 'success' | 'error' | ''; // mocks return empty string
  data: T;
  errorType?: string;
  error?: string;
  warnings?: string[];
}

export type PromRulesResponse = PromResponse<{ groups: PromRuleGroupDTO[] }>;

// Ruler rule DTOs
interface RulerRuleBaseDTO {
  expr: string;
  labels?: Labels;
}

export interface RulerRecordingRuleDTO extends RulerRuleBaseDTO {
  record: string;
}

export interface RulerAlertingRuleDTO extends RulerRuleBaseDTO {
  alert: string;
  for?: string;
  annotations?: Annotations;
}

export enum GrafanaAlertState {
  Alerting = 'Alerting',
  NoData = 'NoData',
  KeepLastState = 'KeepLastState',
  OK = 'OK',
}

export interface GrafanaQueryModel {
  datasource: string;
  datasourceUid: string;

  refId: string;
  [key: string]: any;
}

export interface GrafanaQuery {
  refId: string;
  queryType: string;
  relativeTimeRange: {
    from: number;
    to: number;
  };
  model: GrafanaQueryModel;
}
export interface GrafanaRuleDefinition {
  uid?: string;
  title: string;
  condition: string;
  for: number; //@TODO Sofia will update to accept string
  no_data_state: GrafanaAlertState;
  exec_err_state: GrafanaAlertState;
  data: GrafanaQuery[];
  annotations: Annotations;
  labels: Labels;
}

export interface RulerGrafanaRuleDTO {
  grafana_alert: GrafanaRuleDefinition;
  // labels?: Labels; @TODO to be discussed
  // annotations?: Annotations;
}

export type RulerRuleDTO = RulerAlertingRuleDTO | RulerRecordingRuleDTO | RulerGrafanaRuleDTO;

export type RulerRuleGroupDTO = {
  name: string;
  interval?: string;
  rules: RulerRuleDTO[];
};

export type RulerRulesConfigDTO = { [namespace: string]: RulerRuleGroupDTO[] };
