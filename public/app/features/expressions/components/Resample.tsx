import React, { ChangeEvent, FC } from 'react';

import { SelectableValue } from '@grafana/data';
import { InlineField, InlineFieldRow, Input, Select } from '@grafana/ui';

import { downsamplingTypes, ExpressionQuery, upsamplingTypes } from '../types';

interface Props {
  refIds: Array<SelectableValue<string>>;
  query: ExpressionQuery;
  labelWidth: number;
  onChange: (query: ExpressionQuery) => void;
}

export const Resample: FC<Props> = ({ labelWidth, onChange, refIds, query }) => {
  const downsampler = downsamplingTypes.find((o) => o.value === query.downsampler);
  const upsampler = upsamplingTypes.find((o) => o.value === query.upsampler);

  const onWindowChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, window: event.target.value });
  };

  const onRefIdChange = (value: SelectableValue<string>) => {
    onChange({ ...query, expression: value.value });
  };

  const onSelectDownsampler = (value: SelectableValue<string>) => {
    onChange({ ...query, downsampler: value.value });
  };

  const onSelectUpsampler = (value: SelectableValue<string>) => {
    onChange({ ...query, upsampler: value.value });
  };

  return (
    <>
      <InlineFieldRow>
        <InlineField label="Input" labelWidth={labelWidth}>
          <Select menuShouldPortal onChange={onRefIdChange} options={refIds} value={query.expression} width={20} />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Resample to" labelWidth={labelWidth} tooltip="10s, 1m, 30m, 1h">
          <Input onChange={onWindowChange} value={query.window} width={15} />
        </InlineField>
        <InlineField label="Downsample">
          <Select
            menuShouldPortal
            options={downsamplingTypes}
            value={downsampler}
            onChange={onSelectDownsampler}
            width={25}
          />
        </InlineField>
        <InlineField label="Upsample">
          <Select
            menuShouldPortal
            options={upsamplingTypes}
            value={upsampler}
            onChange={onSelectUpsampler}
            width={25}
          />
        </InlineField>
      </InlineFieldRow>
    </>
  );
};
