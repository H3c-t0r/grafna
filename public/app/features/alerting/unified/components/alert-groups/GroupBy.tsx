import { AlertmanagerGroup } from 'app/plugins/datasource/alertmanager/types';
import React from 'react';
import { uniq } from 'lodash';
import { Icon, Label, MultiSelect } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

interface Props {
  className?: string;
  groups: AlertmanagerGroup[];
  groupBy: string[];
  onGroupingChange: (keys: string[]) => void;
}

export const GroupBy = ({ className, groups, groupBy, onGroupingChange }: Props) => {
  const labelKeyOptions = uniq(
    groups.reduce((keys, group) => {
      group.alerts.forEach(({ labels }) => {
        Object.keys(labels).forEach((label) => keys.push(label));
      });

      return keys;
    }, [] as string[])
  ).map(
    (key) =>
      ({
        label: key,
        value: key,
      } as SelectableValue)
  );

  return (
    <div className={className}>
      <Label>Custom group by</Label>
      <MultiSelect
        value={groupBy}
        placeholder="Group by"
        prefix={<Icon name={'tag-alt'} />}
        onChange={(items) => {
          onGroupingChange(items.map(({ value }) => value as string));
        }}
        options={labelKeyOptions}
      />
    </div>
  );
};
