import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data/src';
import { Stack } from '@grafana/experimental';
import { useStyles2 } from '@grafana/ui';
import { isEmpty } from 'lodash';
import React, { FC } from 'react';
import { useRulesSourcesWithRuler } from '../../../hooks/useRuleSourcesWithRuler';
import { RuleFormType } from '../../../types/rule-form';
import { GrafanaManagedRuleType } from './GrafanaManagedAlert';
import { MimirFlavoredType } from './MimirOrLokiAlert';
import { RecordingRuleType } from './MimirOrLokiRecordingRule';

interface RuleTypePickerProps {
  onChange: (value: RuleFormType) => void;
  selected: RuleFormType;
  enableGrafana: boolean;
  enableCloud: boolean;
}

const RuleTypePicker: FC<RuleTypePickerProps> = ({ selected, onChange, enableCloud, enableGrafana }) => {
  const rulesSourcesWithRuler = useRulesSourcesWithRuler();
  const hasLotexDatasources = !isEmpty(rulesSourcesWithRuler);

  const styles = useStyles2(getStyles);

  return (
    <>
      <Stack direction="row" gap={2}>
        {enableGrafana && <GrafanaManagedRuleType selected={selected === RuleFormType.grafana} onClick={onChange} />}
        {enableCloud && (
          <>
            <MimirFlavoredType
              selected={selected === RuleFormType.cloudAlerting}
              onClick={onChange}
              disabled={!hasLotexDatasources}
            />
            <RecordingRuleType
              selected={selected === RuleFormType.cloudRecording}
              onClick={onChange}
              disabled={!hasLotexDatasources}
            />
          </>
        )}
      </Stack>
      {enableGrafana && enableCloud && (
        <small className={styles.meta}>
          Select &ldquo;Grafana managed&rdquo; unless you have a Mimir, Loki or Cortex data source with the Ruler API
          enabled.
        </small>
      )}
    </>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  meta: css`
    color: ${theme.colors.text.disabled};
  `,
});

export { RuleTypePicker };
