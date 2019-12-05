import React, { useRef, useState } from 'react';
import { cx, css } from 'emotion';
import useClickAway from 'react-use/lib/useClickAway';
import { useExpandableLabel, SegmentProps } from '.';

export interface SegmentInputProps<T> extends SegmentProps<T> {
  value: string | number;
  onChange: (text: string | number) => void;
}

const textWidth = (text: string) => {
  const element = document.createElement('canvas');
  const context = element.getContext('2d')!;
  context.font = '14px Roboto';
  return context.measureText(text).width;
};

export function SegmentInput<T>({
  value,
  onChange,
  Component,
  className,
}: React.PropsWithChildren<SegmentInputProps<T>>) {
  const ref = useRef(null);
  const [inputWidth, setInputWidth] = useState<number>(textWidth(value.toString()));
  const [Label, , expanded, setExpanded] = useExpandableLabel(false);
  useClickAway(ref, () => setExpanded(false));

  if (!expanded) {
    return <Label Component={Component || <a className={cx('gf-form-label', 'query-part', className)}>{value}</a>} />;
  }

  const inputWidthStyle = css`
    width: ${Math.max(inputWidth + 20, 32)}px;
  `;

  return (
    <input
      ref={ref}
      autoFocus
      className={cx(`gf-form gf-form-input`, inputWidthStyle)}
      value={value}
      onChange={item => {
        setInputWidth(textWidth(item.target.value));
        onChange(item.target.value);
      }}
      onBlur={() => setExpanded(false)}
      onKeyDown={e => e.keyCode === 13 && setExpanded(false)}
    />
  );
}
