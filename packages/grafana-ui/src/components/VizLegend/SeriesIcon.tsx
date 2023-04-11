import React, { CSSProperties } from 'react';

import { fieldColorModeRegistry } from '@grafana/data';

import { useTheme2 } from '../../themes';

export interface Props extends React.HTMLAttributes<HTMLDivElement> {
  color?: string;
  gradient?: string;
}

export const SeriesIcon = React.memo(
  React.forwardRef<HTMLDivElement, Props>(({ color, className, gradient, ...restProps }, ref) => {
    const theme = useTheme2();
    let cssColor: string;

    if (gradient) {
      const colors = fieldColorModeRegistry.get(gradient).getColors?.(theme);
      if (colors?.length) {
        cssColor = `linear-gradient(90deg, ${colors.join(', ')})`;
      } else {
        // Not sure what to default to, this will return gray, this should not happen though.
        cssColor = theme.visualization.getColorByName('');
      }
    } else {
      cssColor = color!;
    }

    const styles: CSSProperties = {
      background: cssColor,
      width: '14px',
      height: '4px',
      borderRadius: theme.shape.radius.pill,
      display: 'inline-block',
      marginRight: '8px',
    };

    return <div data-testid="series-icon" ref={ref} className={className} style={styles} {...restProps} />;
  })
);

SeriesIcon.displayName = 'SeriesIcon';
