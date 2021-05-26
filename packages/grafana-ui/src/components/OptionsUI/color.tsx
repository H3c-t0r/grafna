import React from 'react';
import { getColorForTheme, GrafanaTheme } from '@grafana/data';
import { ColorPicker } from '../ColorPicker/ColorPicker';
import { stylesFactory, useTheme } from '../../themes';
import { css } from '@emotion/css';
import { ColorSwatch } from '../ColorPicker/ColorSwatch';

/**
 * @alpha
 * */
export interface ColorValueEditorProps {
  value?: string;
  onChange: (value?: string) => void;
}

/**
 * @alpha
 * */
export const ColorValueEditor: React.FC<ColorValueEditorProps> = ({ value, onChange }) => {
  const theme = useTheme();
  const styles = getStyles(theme);

  return (
    <ColorPicker color={value ?? ''} onChange={onChange} enableNamedColors={true}>
      {({ ref, showColorPicker, hideColorPicker }) => {
        return (
          <div className={styles.spot} onBlur={hideColorPicker}>
            <div className={styles.colorPicker}>
              <ColorSwatch
                ref={ref}
                onClick={showColorPicker}
                onMouseLeave={hideColorPicker}
                color={value ? getColorForTheme(value, theme) : theme.colors.formInputBorder}
              />
            </div>
          </div>
        );
      }}
    </ColorPicker>
  );
};

const getStyles = stylesFactory((theme: GrafanaTheme) => {
  return {
    spot: css`
      color: ${theme.colors.text};
      background: ${theme.colors.formInputBg};
      padding: 3px;
      height: ${theme.spacing.formInputHeight}px;
      border: 1px solid ${theme.colors.formInputBorder};
      display: flex;
      flex-direction: row;
      align-items: center;
      &:hover {
        border: 1px solid ${theme.colors.formInputBorderHover};
      }
    `,
    colorPicker: css`
      padding: 0 ${theme.spacing.sm};
    `,
    colorText: css`
      cursor: pointer;
      flex-grow: 1;
    `,
    trashIcon: css`
      cursor: pointer;
      color: ${theme.colors.textWeak};
      &:hover {
        color: ${theme.colors.text};
      }
    `,
  };
});
