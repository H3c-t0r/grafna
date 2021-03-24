import { RuleGroup, RulesSource } from 'app/types/unified-alerting/internal';
import React, { FC, useMemo, useState, Fragment } from 'react';
import { Icon, Tooltip, useStyles } from '@grafana/ui';
import { GrafanaTheme } from '@grafana/data';
import { css } from 'emotion';
import { isAlertingRule } from '../../utils/rules';
import { PromAlertingRuleState } from 'app/types/unified-alerting/dto';
import { StateColoredText } from '../StateColoredText';
import { CollapseToggle } from '../CollapseToggle';
import { RulesTable } from './RulesTable';
import { isCloudRulesSource } from '../../utils/datasource';
import { ActionIcon } from './ActionIcon';

interface Props {
  namespace: string;
  rulesSource: RulesSource;
  group: RuleGroup;
}

export const RulesGroup: FC<Props> = ({ group, namespace, rulesSource }) => {
  const styles = useStyles(getStyles);

  const [isCollapsed, setIsCollapsed] = useState(true);

  const stats = useMemo(
    (): Record<PromAlertingRuleState, number> =>
      group.rules.reduce<Record<PromAlertingRuleState, number>>(
        (stats, rule) => {
          if (isAlertingRule(rule)) {
            stats[rule.state] += 1;
          }
          return stats;
        },
        {
          [PromAlertingRuleState.Firing]: 0,
          [PromAlertingRuleState.Pending]: 0,
          [PromAlertingRuleState.Inactive]: 0,
        }
      ),
    [group]
  );

  const statsComponents: React.ReactNode[] = [];
  if (stats[PromAlertingRuleState.Firing]) {
    statsComponents.push(
      <StateColoredText key="firing" status={PromAlertingRuleState.Firing}>
        {stats[PromAlertingRuleState.Firing]} firing
      </StateColoredText>
    );
  }
  if (stats[PromAlertingRuleState.Pending]) {
    statsComponents.push(
      <StateColoredText key="firing" status={PromAlertingRuleState.Pending}>
        {stats[PromAlertingRuleState.Pending]} pending
      </StateColoredText>
    );
  }

  return (
    <div className={styles.wrapper} data-testid="rule-group">
      <div className={styles.header} data-testid="rule-group-header">
        <CollapseToggle className={styles.collapseToggle} isCollapsed={isCollapsed} onToggle={setIsCollapsed} />
        <Icon name={isCollapsed ? 'folder-open' : 'folder'} />
        {isCloudRulesSource(rulesSource) && (
          <Tooltip content={rulesSource.name} placement="top">
            <img className={styles.datasourceIcon} src={rulesSource.meta.info.logos.small} />
          </Tooltip>
        )}
        <h6 className={styles.heading}>
          {namespace} &gt; {group.name}
        </h6>
        <div className={styles.spacer} />
        <div className={styles.headerStats}>
          {group.rules.length} rules
          {!!statsComponents.length && (
            <>
              :{' '}
              {statsComponents.reduce<React.ReactNode[]>(
                (prev, curr, idx) => (prev.length ? [<Fragment key={idx}>, </Fragment>, curr] : [curr]),
                []
              )}
            </>
          )}
        </div>
        <div className={styles.actionsSeparator}>|</div>
        <div className={styles.actionIcons}>
          <ActionIcon icon="pen" tooltip="edit" />
          <ActionIcon icon="lock" tooltip="manage permissions" />
        </div>
      </div>
      {!isCollapsed && <RulesTable rulesSource={rulesSource} namespace={namespace} group={group} />}
    </div>
  );
};

export const getStyles = (theme: GrafanaTheme) => ({
  wrapper: css`
    & + & {
      margin-top: ${theme.spacing.md};
    }
  `,
  header: css`
    display: flex;
    flex-direction: row;
    align-items: center;
    padding: ${theme.spacing.sm} ${theme.spacing.sm} ${theme.spacing.sm} 0;
    background-color: ${theme.colors.bg2};
  `,
  headerStats: css`
    span {
      vertical-align: middle;
    }
  `,
  heading: css`
    margin-left: ${theme.spacing.sm};
    margin-bottom: 0;
  `,
  spacer: css`
    flex: 1;
  `,
  collapseToggle: css`
    background: none;
    border: none;
    margin-top: -${theme.spacing.sm};
    margin-bottom: -${theme.spacing.sm};

    svg {
      margin-bottom: 0;
    }
  `,
  datasourceIcon: css`
    width: ${theme.spacing.md};
    height: ${theme.spacing.md};
    margin-left: ${theme.spacing.md};
  `,
  datasourceOrigin: css`
    margin-right: 1em;
    color: ${theme.colors.textFaint};
  `,
  actionsSeparator: css`
    margin: 0 ${theme.spacing.sm};
  `,
  actionIcons: css`
    & > * + * {
      margin-left: ${theme.spacing.sm};
    }
  `,
});
