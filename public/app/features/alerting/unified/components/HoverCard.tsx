import { css } from '@emotion/css';
import { Placement } from '@popperjs/core';
import classnames from 'classnames';
import React, { ReactElement, ReactNode, useRef } from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { Popover as GrafanaPopover, PopoverController, useStyles2, Stack } from '@grafana/ui';

export interface HoverCardProps {
  children: ReactElement;
  header?: ReactNode;
  content: ReactElement;
  footer?: ReactNode;
  wrapperClassName?: string;
  placement?: Placement;
  disabled?: boolean;
  showAfter?: number;
  arrow?: boolean;
}

export const HoverCard = ({
  children,
  header,
  content,
  footer,
  arrow,
  showAfter = 300,
  wrapperClassName,
  disabled = false,
  ...rest
}: HoverCardProps) => {
  const popoverRef = useRef<HTMLElement>(null);
  const styles = useStyles2(getStyles);

  if (disabled) {
    return children;
  }

  const body = (
    <Stack direction="column" gap={0}>
      {header && <div className={styles.card.header}>{header}</div>}
      <div className={styles.card.body}>{content}</div>
      {footer && <div className={styles.card.footer}>{footer}</div>}
    </Stack>
  );

  return (
    <PopoverController content={body} hideAfter={100}>
      {(showPopper, hidePopper, popperProps) => {
        return (
          <>
            {popoverRef.current && (
              <GrafanaPopover
                {...popperProps}
                {...rest}
                wrapperClassName={classnames(styles.popover(arrow ? 1.25 : 0), wrapperClassName)}
                onMouseLeave={hidePopper}
                onMouseEnter={showPopper}
                onFocus={showPopper}
                onBlur={hidePopper}
                referenceElement={popoverRef.current}
                renderArrow={arrow}
              />
            )}

            {React.cloneElement(children, {
              ref: popoverRef,
              onMouseEnter: showPopper,
              onMouseLeave: hidePopper,
              onFocus: showPopper,
              onBlur: hidePopper,
            })}
          </>
        );
      }}
    </PopoverController>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  popover: (offset: number) => css`
    border-radius: ${theme.shape.radius.default};
    box-shadow: ${theme.shadows.z3};
    background: ${theme.colors.background.primary};
    border: 1px solid ${theme.colors.border.medium};

    margin-bottom: ${theme.spacing(offset)};
  `,
  card: {
    body: css`
      padding: ${theme.spacing(1)};
    `,
    header: css`
      padding: ${theme.spacing(1)};
      background: ${theme.colors.background.secondary};
      border-bottom: solid 1px ${theme.colors.border.medium};
    `,
    footer: css`
      padding: ${theme.spacing(0.5)} ${theme.spacing(1)};
      background: ${theme.colors.background.secondary};
      border-top: solid 1px ${theme.colors.border.medium};
    `,
  },
});
