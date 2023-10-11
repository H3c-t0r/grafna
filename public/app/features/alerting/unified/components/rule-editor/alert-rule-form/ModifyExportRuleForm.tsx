import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import { useAsync } from 'react-use';

import { Stack } from '@grafana/experimental';
import { Button, CustomScrollbar, LinkButton, LoadingPlaceholder } from '@grafana/ui';
import { useAppNotification } from 'app/core/copy/appNotification';
import { useQueryParams } from 'app/core/hooks/useQueryParams';

import { AppChromeUpdate } from '../../../../../../core/components/AppChrome/AppChromeUpdate';
import { RulerRuleDTO } from '../../../../../../types/unified-alerting-dto';
import { alertRuleApi, ModifyExportPayload } from '../../../api/alertRuleApi';
import { fetchRulerRulesGroup } from '../../../api/ruler';
import { useDataSourceFeatures } from '../../../hooks/useCombinedRule';
import { RuleFormValues } from '../../../types/rule-form';
import { GRAFANA_RULES_SOURCE_NAME } from '../../../utils/datasource';
import { formValuesToRulerGrafanaRuleDTO, MINUTE } from '../../../utils/rule-form';
import { isGrafanaRulerRule } from '../../../utils/rules';
import { FileExportPreview } from '../../export/FileExportPreview';
import { GrafanaExportDrawer } from '../../export/GrafanaExportDrawer';
import { allGrafanaExportProviders, ExportFormats } from '../../export/providers';
import { AlertRuleNameInput } from '../AlertRuleNameInput';
import AnnotationsStep from '../AnnotationsStep';
import { GrafanaEvaluationBehavior } from '../GrafanaEvaluationBehavior';
import { NotificationsStep } from '../NotificationsStep';
import { QueryAndExpressionsStep } from '../query-and-alert-condition/QueryAndExpressionsStep';

interface ModifyExportRuleFormProps {
  alertUid?: string;
  ruleForm?: RuleFormValues;
}

type ModifyExportMode = 'rule' | 'group';

export function ModifyExportRuleForm({ ruleForm, alertUid }: ModifyExportRuleFormProps) {
  const formAPI = useForm<RuleFormValues>({
    mode: 'onSubmit',
    defaultValues: ruleForm,
    shouldFocusError: true,
  });
  const [queryParams] = useQueryParams();

  const existing = Boolean(ruleForm); // always should be true
  const notifyApp = useAppNotification();
  const returnTo = !queryParams['returnTo'] ? '/alerting/list' : String(queryParams['returnTo']);

  const [exportData, setExportData] = useState<
    { exportForm: RuleFormValues; exportMode: ModifyExportMode } | undefined
  >(undefined);

  const [conditionErrorMsg, setConditionErrorMsg] = useState('');
  const [evaluateEvery, setEvaluateEvery] = useState(ruleForm?.evaluateEvery ?? MINUTE);

  const checkAlertCondition = (msg = '') => {
    setConditionErrorMsg(msg);
  };

  const submit = (exportData: { exportForm: RuleFormValues; exportMode: ModifyExportMode } | undefined) => {
    if (conditionErrorMsg !== '') {
      notifyApp.error(conditionErrorMsg);
      return;
    }
    setExportData(exportData);
  };

  const onClose = useCallback(() => {
    setExportData(undefined);
  }, [setExportData]);

  const actionButtons = [
    <LinkButton href={returnTo} key="cancel" size="sm" variant="secondary" onClick={() => submit(undefined)}>
      Cancel
    </LinkButton>,
    <Button
      key="export-rule"
      size="sm"
      onClick={formAPI.handleSubmit((formValues) => submit({ exportForm: formValues, exportMode: 'rule' }))}
    >
      Export Rule
    </Button>,
    <Button
      key="export-group"
      size="sm"
      onClick={formAPI.handleSubmit((formValues) => submit({ exportForm: formValues, exportMode: 'group' }))}
    >
      Export Group
    </Button>,
  ];

  return (
    <>
      <FormProvider {...formAPI}>
        <AppChromeUpdate actions={actionButtons} />
        <form onSubmit={(e) => e.preventDefault()}>
          <div>
            <CustomScrollbar autoHeightMin="100%" hideHorizontalTrack={true}>
              <Stack direction="column" gap={3}>
                {/* Step 1 */}
                <AlertRuleNameInput />
                {/* Step 2 */}
                <QueryAndExpressionsStep editingExistingRule={existing} onDataChange={checkAlertCondition} />
                {/* Step 3-4-5 */}

                <GrafanaEvaluationBehavior
                  evaluateEvery={evaluateEvery}
                  setEvaluateEvery={setEvaluateEvery}
                  existing={Boolean(existing)}
                  includeProvisioned={true}
                />

                {/* Step 4 & 5 */}
                {/* Annotations only for cloud and Grafana */}
                <AnnotationsStep />
                {/* Notifications step*/}
                <NotificationsStep alertUid={alertUid} />
              </Stack>
            </CustomScrollbar>
          </div>
        </form>
        {exportData && (
          <GrafanaRuleDesignExporter
            exportValues={exportData.exportForm}
            exportMode={exportData.exportMode}
            onClose={onClose}
          />
        )}
      </FormProvider>
    </>
  );
}

const useGetGroup = (nameSpace: string, group: string) => {
  const { dsFeatures } = useDataSourceFeatures(GRAFANA_RULES_SOURCE_NAME);

  const rulerConfig = dsFeatures?.rulerConfig;

  const targetGroup = useAsync(async () => {
    return rulerConfig ? await fetchRulerRulesGroup(rulerConfig, nameSpace, group) : undefined;
  }, [rulerConfig, nameSpace, group]);

  return targetGroup;
};

interface GrafanaRuleDesignExportPreviewProps {
  exportFormat: ExportFormats;
  onClose: () => void;
  exportMode: ModifyExportMode;
  exportValues: RuleFormValues;
}

const useGetPayloadToExport = (values: RuleFormValues, exportMode: ModifyExportMode) => {
  const rulerGroupDto = useGetGroup(values.folder?.title ?? '', values.group);
  const grafanaRuleDto = useMemo(() => formValuesToRulerGrafanaRuleDTO(values), [values]);
  const includeRulesInGroup = exportMode === 'group';
  const routeParams = useParams<{ type: string; id: string }>();
  const uidFromParams = routeParams.id;

  const payload: ModifyExportPayload = useMemo(() => {
    const updatedRule = { ...grafanaRuleDto, grafana_alert: { ...grafanaRuleDto.grafana_alert, uid: uidFromParams } };
    if (rulerGroupDto?.value?.rules) {
      const rulesWithoutCurrent = rulerGroupDto.value.rules.filter(
        (rule: RulerRuleDTO) => isGrafanaRulerRule(rule) && rule.grafana_alert.uid !== uidFromParams
      );
      return {
        ...rulerGroupDto?.value,
        rules: [...(includeRulesInGroup ? rulesWithoutCurrent : []), updatedRule],
      };
    } else {
      return {
        name: values.group,
        rules: [updatedRule],
      };
    }
  }, [rulerGroupDto.value, grafanaRuleDto, values.group, includeRulesInGroup, uidFromParams]);
  return { payload, loadingGroup: rulerGroupDto.loading };
};

const GrafanaRuleDesignExportPreview = ({
  exportFormat,
  exportValues,
  onClose,
  exportMode,
}: GrafanaRuleDesignExportPreviewProps) => {
  const [getExport, exportData] = alertRuleApi.endpoints.exportModifiedRuleGroup.useMutation();
  const { loadingGroup, payload } = useGetPayloadToExport(exportValues, exportMode);

  const nameSpace = exportValues.folder?.title ?? '';

  useEffect(() => {
    !loadingGroup && getExport({ payload, format: exportFormat, nameSpace: nameSpace });
  }, [nameSpace, exportFormat, payload, getExport, loadingGroup]);

  const downloadFileName = `modify-export-${new Date().getTime()}`;

  if (exportData.isLoading) {
    return <LoadingPlaceholder text="Loading...." />;
  }

  return (
    <FileExportPreview
      format={exportFormat}
      textDefinition={exportData.data ?? ''}
      downloadFileName={downloadFileName}
      onClose={onClose}
    />
  );
};

interface GrafanaRuleDesignExporterProps {
  onClose: () => void;
  exportValues: RuleFormValues;
  exportMode: ModifyExportMode;
}

export const GrafanaRuleDesignExporter = React.memo(
  ({ onClose, exportMode, exportValues }: GrafanaRuleDesignExporterProps) => {
    const [activeTab, setActiveTab] = useState<ExportFormats>('yaml');
    const title = exportMode === 'rule' ? 'Export Rule' : 'Export Group';

    return (
      <GrafanaExportDrawer
        title={title}
        activeTab={activeTab}
        onTabChange={setActiveTab}
        onClose={onClose}
        formatProviders={Object.values(allGrafanaExportProviders)}
      >
        <GrafanaRuleDesignExportPreview
          exportFormat={activeTab}
          onClose={onClose}
          exportMode={exportMode}
          exportValues={exportValues}
        />
      </GrafanaExportDrawer>
    );
  }
);

GrafanaRuleDesignExporter.displayName = 'GrafanaRuleDesignExporter';
