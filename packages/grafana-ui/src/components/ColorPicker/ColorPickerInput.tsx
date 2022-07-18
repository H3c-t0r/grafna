import { css } from '@emotion/css';
import React, { useState, forwardRef } from 'react';
import { RgbaStringColorPicker } from 'react-colorful';
import { useThrottleFn } from 'react-use';

import { colorManipulator, GrafanaTheme2 } from '@grafana/data';

import { useStyles2, useTheme2 } from '../../themes';
import { ClickOutsideWrapper } from '../ClickOutsideWrapper/ClickOutsideWrapper';
import { Props as InputProps } from '../Input/Input';

import ColorInput from './ColorInput';

export interface ColorPickerInputProps extends Omit<InputProps, 'value' | 'onChange'> {
  value?: string;
  onChange: (color: string) => void;
  /** Format for returning the color in onChange callback, defaults to 'rgb' */
  returnColorAs?: 'rgb' | 'hex';
}

export const ColorPickerInput = forwardRef<HTMLInputElement, ColorPickerInputProps>(
  ({ value = '', onChange, returnColorAs = 'rgb', ...inputProps }, ref) => {
    const [currentColor, setColor] = useState(value);
    const [isOpen, setIsOpen] = useState(false);
    const theme = useTheme2();
    const styles = useStyles2(getStyles);

    useThrottleFn(
      (c) => {
        const color = theme.visualization.getColorByName(c);
        if (returnColorAs === 'rgb') {
          onChange(colorManipulator.asRgbString(color));
        } else {
          onChange(colorManipulator.asHexString(color));
        }
      },
      500,
      [currentColor]
    );

    return (
      <ClickOutsideWrapper onClick={() => setIsOpen(false)}>
        <div className={styles.wrapper}>
          {isOpen && <RgbaStringColorPicker color={currentColor} onChange={setColor} className={styles.picker} />}
          <div onClick={() => setIsOpen(true)}>
            <ColorInput
              {...inputProps}
              theme={theme}
              color={currentColor}
              onChange={setColor}
              className={styles.input}
              ref={ref}
            />
          </div>
        </div>
      </ClickOutsideWrapper>
    );
  }
);

ColorPickerInput.displayName = 'ColorPickerInput';

const getStyles = (theme: GrafanaTheme2) => {
  return {
    wrapper: css`
      position: relative;
    `,
    picker: css`
      position: absolute;
      bottom: ${theme.spacing(0.5)};
      padding: ${theme.spacing(0.25)};
      border-radius: ${theme.shape.borderRadius(4)};
      background: ${theme.colors.border.medium};
    `,
    input: css``,
  };
};
