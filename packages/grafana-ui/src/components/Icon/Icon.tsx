import React from 'react';
import { css, cx } from 'emotion';
import { camelCase } from 'lodash';
import { GrafanaTheme } from '@grafana/data';
import { stylesFactory } from '../../themes';
import { useTheme } from '../../themes/ThemeContext';
import { IconName, IconType } from '../../types';
import { ComponentSize } from '../../types/size';
//@ts-ignore
import * as DefaultIcon from '@iconscout/react-unicons';
import * as MonoIcon from './assets';

interface IconProps extends React.HTMLAttributes<HTMLDivElement> {
  name: IconName;
  size?: ComponentSize | 'xl' | 'xxl';
  type?: IconType;
}
export interface SvgProps extends React.HTMLAttributes<SVGElement> {
  size: number;
  secondaryColor?: string;
  className?: string;
}

const getIconStyles = stylesFactory((theme: GrafanaTheme, type: IconType) => {
  const defaultIconColor = type === 'default' ? 'currentColor' : theme.colors.orange;
  return {
    container: css`
      display: inline-block;
    `,
    icon: css`
      vertical-align: middle;
      display: inline-block;
      margin-bottom: ${theme.spacing.xxs};
      cursor: pointer;
      * {
        fill: ${defaultIconColor};
      }
    `,
  };
});

export const Icon = React.forwardRef<HTMLDivElement, IconProps>(
  ({ size = 'md', type = 'default', title, name, className, style, ...divElementProps }, ref) => {
    const theme = useTheme();
    const styles = getIconStyles(theme, type);
    const svgSize = getSvgSize(size, theme);
    const iconName = type === 'default' ? `Uil${pascalCase(name)}` : pascalCase(name);

    /* Unicons don't have type definitions */
    //@ts-ignore
    const Component = type === 'default' ? DefaultIcon[iconName] : MonoIcon[iconName];

    if (!Component) {
      return <div />;
    }

    return (
      <div className={styles.container} {...divElementProps} ref={ref}>
        {type === 'default' && <Component size={svgSize} className={cx(styles.icon, className)} style={style} />}
        {type === 'mono' && <Component size={svgSize} className={cx(styles.icon, className)} style={style} />}
      </div>
    );
  }
);

Icon.displayName = 'Icon';

const pascalCase = (string: string) => {
  const str = camelCase(string);
  return str.charAt(0).toUpperCase() + str.substring(1);
};

/* Transform string with px to number and add 2 pxs as path in svg is 2px smaller */
const getSvgSize = (size: ComponentSize | 'xl' | 'xxl', theme: GrafanaTheme) => {
  let svgSize;
  if (size === 'xl') {
    svgSize = Number(theme.typography.heading.h1.slice(0, -2));
  } else if (size === 'xxl') {
    svgSize = Number(theme.height.lg.slice(0, -2));
  } else {
    svgSize = Number(theme.typography.size[size].slice(0, -2)) + 2;
  }
  return svgSize;
};
