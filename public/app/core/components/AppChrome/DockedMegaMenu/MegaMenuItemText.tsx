import { css, cx } from '@emotion/css';
import React from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { selectors } from '@grafana/e2e-selectors';
import { Icon, IconName, Link, useTheme2 } from '@grafana/ui';

export interface Props {
  children: React.ReactNode;
  icon?: IconName;
  isActive?: boolean;
  isChild?: boolean;
  onClick?: () => void;
  target?: HTMLAnchorElement['target'];
  url?: string;
  level: number;
}

export function MegaMenuItemText({ children, icon, isActive, isChild, onClick, target, url, level }: Props) {
  const theme = useTheme2();
  const styles = getStyles(theme, isActive, isChild, level);

  const linkContent = (
    <div className={styles.linkContent}>
      {icon && <Icon data-testid="dropdown-child-icon" name={icon} />}

      {children}

      {
        // As nav links are supposed to link to internal urls this option should be used with caution
        target === '_blank' && (
          <Icon data-testid="external-link-icon" name="external-link-alt" className={styles.externalLinkIcon} />
        )
      }
    </div>
  );

  let element = (
    <button
      data-testid={selectors.components.NavMenu.item}
      className={cx(styles.button, styles.element)}
      onClick={onClick}
    >
      {linkContent}
    </button>
  );

  if (url) {
    element =
      !target && url.startsWith('/') ? (
        <Link
          data-testid={selectors.components.NavMenu.item}
          className={styles.element}
          href={url}
          target={target}
          onClick={onClick}
        >
          {linkContent}
        </Link>
      ) : (
        <a
          data-testid={selectors.components.NavMenu.item}
          href={url}
          target={target}
          className={styles.element}
          onClick={onClick}
        >
          {linkContent}
        </a>
      );
  }

  return <div className={styles.wrapper}>{element}</div>;
}

MegaMenuItemText.displayName = 'MegaMenuItemText';

const getStyles = (
  theme: GrafanaTheme2,
  isActive: Props['isActive'],
  isChild: Props['isActive'],
  level: Props['level']
) => ({
  button: css({
    backgroundColor: 'unset',
    borderStyle: 'unset',
  }),
  linkContent: css({
    alignItems: 'center',
    display: 'flex',
    gap: '0.5rem',
    height: '100%',
    width: '100%',
  }),
  externalLinkIcon: css({
    color: theme.colors.text.secondary,
  }),
  element: css({
    alignItems: 'center',
    boxSizing: 'border-box',
    position: 'relative',
    color: isActive ? theme.colors.text.primary : theme.colors.text.secondary,
    padding: theme.spacing(1, 1, 1, isChild ? 5 : 0),
    ...(isChild && {
      borderRadius: theme.shape.radius.default,
    }),
    width: '100%',
    '&:hover, &:focus-visible': {
      ...(isChild && {
        background: theme.colors.emphasize(theme.colors.background.primary, 0.03),
      }),
      textDecoration: 'underline',
      color: theme.colors.text.primary,
    },
    '&:focus-visible': {
      boxShadow: 'none',
      outline: `2px solid ${theme.colors.primary.main}`,
      outlineOffset: '-2px',
      transition: 'none',
    },
  }),
  wrapper: css({
    boxSizing: 'border-box',
    position: 'relative',
    display: 'flex',
    width: '100%',
    ...(isChild && {
      padding: theme.spacing(0, 2),
    }),
  }),
});
