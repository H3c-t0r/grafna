import { CombinedRuleGroup, CombinedRuleNamespace } from 'app/types/unified-alerting';
import React, { FC, useMemo, useState, Fragment } from 'react';
import { Icon, Tooltip, useStyles2 } from '@grafana/ui';
import { GrafanaTheme2 } from '@grafana/data';
import { css } from '@emotion/css';
import { isAlertingRule, isGrafanaRulerRule } from '../../utils/rules';
import { PromAlertingRuleState } from 'app/types/unified-alerting-dto';
import { StateColoredText } from '../StateColoredText';
import { CollapseToggle } from '../CollapseToggle';
import { RulesTable } from './RulesTable';
import { GRAFANA_RULES_SOURCE_NAME, isCloudRulesSource } from '../../utils/datasource';
import { ActionIcon } from './ActionIcon';
import pluralize from 'pluralize';
import { useHasRuler } from '../../hooks/useHasRuler';
import kbn from 'app/core/utils/kbn';
import { useFolder } from '../../hooks/useFolder';

interface Props {
  namespace: CombinedRuleNamespace;
  group: CombinedRuleGroup;
}

export const RulesGroup: FC<Props> = React.memo(({ group, namespace }) => {
  const { rulesSource } = namespace;
  const styles = useStyles2(getStyles);

  const [isCollapsed, setIsCollapsed] = useState(true);

  const hasRuler = useHasRuler();
  const rulerRule = group.rules[0]?.rulerRule;
  const folderUID = (rulerRule && isGrafanaRulerRule(rulerRule) && rulerRule.grafana_alert.namespace_uid) || undefined;
  const { folder } = useFolder(folderUID);

  const stats = useMemo(
    (): Record<PromAlertingRuleState, number> =>
      group.rules.reduce<Record<PromAlertingRuleState, number>>(
        (stats, rule) => {
          if (rule.promRule && isAlertingRule(rule.promRule)) {
            stats[rule.promRule.state] += 1;
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
      <StateColoredText key="pending" status={PromAlertingRuleState.Pending}>
        {stats[PromAlertingRuleState.Pending]} pending
      </StateColoredText>
    );
  }

  const actionIcons: React.ReactNode[] = [];

  // for grafana, link to folder views
  if (rulesSource === GRAFANA_RULES_SOURCE_NAME) {
    if (folderUID) {
      const baseUrl = `/dashboards/f/${folderUID}/${kbn.slugifyForUrl(namespace.name)}`;
      if (folder?.canSave) {
        actionIcons.push(
          <ActionIcon key="edit" icon="pen" tooltip="edit" to={baseUrl + '/settings'} target="__blank" />
        );
      }
      if (folder?.canAdmin) {
        actionIcons.push(
          <ActionIcon
            key="manage-perms"
            icon="lock"
            tooltip="manage permissions"
            to={baseUrl + '/permissions'}
            target="__blank"
          />
        );
      }
    } else if (hasRuler(rulesSource)) {
      actionIcons.push(<ActionIcon key="edit" icon="pen" tooltip="edit" />); // @TODO
    }
  }

  const groupName = isCloudRulesSource(rulesSource) ? `${namespace.name} > ${group.name}` : namespace.name;

  return (
    <div className={styles.wrapper} data-testid="rule-group">
      <div className={styles.header} data-testid="rule-group-header">
        <CollapseToggle
          className={styles.collapseToggle}
          isCollapsed={isCollapsed}
          onToggle={setIsCollapsed}
          data-testid="group-collapse-toggle"
        />
        <Icon name={isCollapsed ? 'folder-open' : 'folder'} />
        {isCloudRulesSource(rulesSource) && (
          <Tooltip content={rulesSource.name} placement="top">
            <img className={styles.dataSourceIcon} src={rulesSource.meta.info.logos.small} />
          </Tooltip>
        )}
        <h6 className={styles.heading}>{groupName}</h6>
        <div className={styles.spacer} />
        <div className={styles.headerStats}>
          {group.rules.length} {pluralize('rule', group.rules.length)}
          {!!statsComponents.length && (
            <>
              :{' '}
              {statsComponents.reduce<React.ReactNode[]>(
                (prev, curr, idx) => (prev.length ? [prev, <Fragment key={idx}>, </Fragment>, curr] : [curr]),
                []
              )}
            </>
          )}
        </div>
        {!!actionIcons.length && (
          <>
            <div className={styles.actionsSeparator}>|</div>
            <div className={styles.actionIcons}>{actionIcons}</div>
          </>
        )}
      </div>
      {!isCollapsed && (
        <RulesTable showSummaryColumn={true} className={styles.rulesTable} showGuidelines={true} rules={group.rules} />
      )}
    </div>
  );
});

RulesGroup.displayName = 'RulesGroup';

export const getStyles = (theme: GrafanaTheme2) => ({
  wrapper: css`
    & + & {
      margin-top: ${theme.spacing(2)};
    }
  `,
  header: css`
    display: flex;
    flex-direction: row;
    align-items: center;
    padding: ${theme.spacing(1)} ${theme.spacing(1)} ${theme.spacing(1)} 0;
    background-color: ${theme.colors.background.secondary};
  `,
  headerStats: css`
    span {
      vertical-align: middle;
    }
  `,
  heading: css`
    margin-left: ${theme.spacing(1)};
    margin-bottom: 0;
  `,
  spacer: css`
    flex: 1;
  `,
  collapseToggle: css`
    background: none;
    border: none;
    margin-top: -${theme.spacing(1)};
    margin-bottom: -${theme.spacing(1)};

    svg {
      margin-bottom: 0;
    }
  `,
  dataSourceIcon: css`
    width: ${theme.spacing(2)};
    height: ${theme.spacing(2)};
    margin-left: ${theme.spacing(2)};
  `,
  dataSourceOrigin: css`
    margin-right: 1em;
    color: ${theme.colors.text.disabled};
  `,
  actionsSeparator: css`
    margin: 0 ${theme.spacing(2)};
  `,
  actionIcons: css`
    & > * + * {
      margin-left: ${theme.spacing(1)};
    }
  `,
  rulesTable: css`
    margin-top: ${theme.spacing(3)};
  `,
});
