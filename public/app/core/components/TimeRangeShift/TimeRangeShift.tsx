import React, { useState } from 'react';
import { useDebounce } from 'react-use';

import { InlineField, InlineFieldRow, Input } from '@grafana/ui';

interface Props {
  label: string;
  tooltip: string;
  value: string;
  onChange: (val: string) => void;
  isInvalid: boolean;
  disabled?: boolean;
}

export function TimeRangeShift(props: Props) {
  const [localValue, setLocalValue] = useState(props.value);

  useDebounce(
    () => {
      props.onChange(localValue);
    },
    500,
    [localValue]
  );

  return (
    <InlineFieldRow>
      <InlineField
        label={props.label}
        labelWidth={26}
        disabled={props.disabled ?? false}
        grow
        tooltip={props.tooltip}
        invalid={props.isInvalid}
        error={'Invalid time shift. See tooltip for examples.'}
      >
        <Input
          type="text"
          placeholder="0"
          width={40}
          onChange={(e) => setLocalValue(e.currentTarget.value)}
          value={localValue}
        />
      </InlineField>
    </InlineFieldRow>
  );
}
