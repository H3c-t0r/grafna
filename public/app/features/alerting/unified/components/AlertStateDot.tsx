import { css } from '@emotion/css';
import React from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { ComponentSize, Stack, useStyles2 } from '@grafana/ui';
import { RuleHealth } from 'app/types/unified-alerting';
import { PromAlertingRuleState } from 'app/types/unified-alerting-dto';

interface DotStylesProps {
  state: PromAlertingRuleState;
  health?: RuleHealth;
  includeState?: boolean;
  size?: ComponentSize; // TODO support this
}

const AlertStateDot = (props: DotStylesProps) => {
  const styles = useStyles2(getDotStyles, props);

  return (
    <Stack direction="row" gap={0.5}>
      <div className={styles.dot} />
    </Stack>
  );
};

const getDotStyles = (theme: GrafanaTheme2, props: DotStylesProps) => {
  const size = theme.spacing(1.25);
  const outlineSize = `calc(${size} / 2.5)`;

  const isError = props.health === 'error';
  const errorStyle = props.state === PromAlertingRuleState.Firing || isError;
  const successStyle = props.state === PromAlertingRuleState.Inactive;
  const pendingStyle = props.state === PromAlertingRuleState.Pending;

  return {
    dot: css`
      width: ${size};
      height: ${size};

      border-radius: 100%;

      background-color: ${theme.colors.secondary.main};
      outline: solid ${outlineSize} ${theme.colors.secondary.transparent};
      margin: ${outlineSize};

      ${successStyle &&
      css({
        backgroundColor: theme.colors.success.main,
        outlineColor: theme.colors.success.transparent,
      })}

      ${pendingStyle &&
      css({
        backgroundColor: theme.colors.warning.main,
        outlineColor: theme.colors.warning.transparent,
      })}

      ${errorStyle &&
      css({
        backgroundColor: theme.colors.error.main,
        outlineColor: theme.colors.error.transparent,
      })}
    `,
  };
};

export { AlertStateDot };
