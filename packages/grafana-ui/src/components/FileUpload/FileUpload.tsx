import React, { FC, FormEvent } from 'react';
import { GrafanaTheme } from '@grafana/data';
import { css, cx } from 'emotion';
import { getFormStyles, Icon } from '../index';
import { stylesFactory, useTheme } from '../../themes';

export interface Props {
  onFileUpload: (event: FormEvent<HTMLInputElement>) => void;
  /** Custom upload label text */
  label?: string;
  /** Accepted file extensions */
  accept?: string;
  className?: string;
}

export const FileUpload: FC<Props> = ({ onFileUpload, className, label = 'Upload file', accept = '*' }) => {
  const theme = useTheme();
  const style = getStyles(theme);

  return (
    <label className={cx(style.button, className)}>
      <Icon name="upload" className={style.icon} />
      {label}
      <input
        type="file"
        id="fileUpload"
        className={style.fileUpload}
        onChange={onFileUpload}
        multiple={false}
        accept={accept}
      />
    </label>
  );
};

const getStyles = stylesFactory((theme: GrafanaTheme) => {
  const buttonFormStyle = getFormStyles(theme, { variant: 'primary', invalid: false, size: 'md' }).button.button;
  return {
    fileUpload: css`
      display: none;
    `,
    button: css`
      ${buttonFormStyle}
    `,
    icon: css`
      margin-right: ${theme.spacing.xs};
    `,
  };
});
