import { css } from '@emotion/css';
import React, { useState, MouseEvent } from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { useStyles2 } from '@grafana/ui';

import { useTimeout } from '../../hooks/useTimeout';

type HotEdgePosition = 'left' | 'right';

interface HotEdgeProps {
  position: HotEdgePosition;
  isDragging: boolean;
  onActivate: () => void;
}

export const HotEdge = ({ position = 'right', isDragging, onActivate }: HotEdgeProps) => {
  const [isHovered, setIsHovered] = useState(false);
  const [isEnabled, setIsEnabled] = useState(false);
  const styles = useStyles2(getHotEdgeStyles, { position, isEnabled });

  const reset = () => {
    setIsHovered(false);
    setIsEnabled(false);
  };

  const handleMoveEnter = (e: MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.preventDefault();
    e.stopPropagation();

    setIsHovered(true);
    console.log('handleMoveEnter');
  };

  const handleMoveOver = (e: MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.preventDefault();
    e.stopPropagation();

    console.log('handleMoveOver');
  };

  const handleMoveLeave = (e: MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.preventDefault();
    e.stopPropagation();

    reset();
    console.log('handleMoveLeave');
  };

  const handleActivated = () => {
    onActivate();
    reset();
  };

  useTimeout(() => setIsEnabled(true), isHovered ? 300 : null);
  useTimeout(handleActivated, isEnabled && isDragging ? 300 : null);

  return (
    <div
      className={styles.area}
      onDragEnter={handleMoveEnter}
      onDragLeave={handleMoveLeave}
      onDragOver={handleMoveOver}
      onMouseEnter={handleMoveEnter}
      onMouseLeave={handleMoveLeave}
    >
      <div className={styles.tab} onClick={handleActivated} />
    </div>
  );
};

interface GetHotEdgeStylesProps {
  position: HotEdgePosition;
  isEnabled: boolean;
}

const getPosition = ({ isEnabled, position }: GetHotEdgeStylesProps) => {
  if (!isEnabled) {
    if (position === 'left') {
      return '-100%';
    }

    return '100%';
  }

  return '0%';
};

const getHotEdgeStyles = (theme: GrafanaTheme2, props: GetHotEdgeStylesProps) => ({
  area: css({
    position: 'fixed',
    top: 0,
    bottom: 0,
    width: '15px',
    zIndex: 9999,
    [props.position]: 0,
  }),
  tab: css({
    position: 'absolute',
    width: '20px',
    height: '100%',
    background: theme.colors.border.weak,
    transform: `translateX(${getPosition(props)})`,
    transition: 'transform 0.2s ease-in-out',
    top: 0,
    [props.position]: 0,
  }),
});
