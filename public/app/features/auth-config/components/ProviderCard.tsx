import { css, cx } from '@emotion/css';
import React from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { IconButton, useStyles2 } from '@grafana/ui';

import { AuthBadge, OAuthBadge } from './AuthBadges';

// const configPath = 'admin/authentication';
export const LOGO_SIZE = '48px';

type Props = {
  providerId: string;
  displayName: string;
  enabled: boolean;
  configPath?: string;
  badges?: JSX.Element[];
};

export function ProviderCard({ providerId, displayName, enabled, configPath, badges }: Props) {
  const styles = useStyles2(getStyles);
  const basePath = configPath || 'admin/authentication';

  return (
    <a href={`${basePath}/${providerId}`} className={styles.container}>
      <div className={styles.header}>
        <div className={styles.badgesContainer}>{getProviderBadges(providerId)}</div>
        <IconButton
          name={enabled ? 'pen' : 'plus-circle'}
          size="lg"
          variant="secondary"
          className={styles.actionButton}
        />
      </div>
      <h2 className={cx(styles.name, { [styles.disabled]: !enabled })}>{displayName}</h2>
      <div className={styles.content}></div>
      <div className={styles.status}>{enabled ? 'Configured' : 'Not configured'}</div>
    </a>
  );
}

function getProviderBadges(providerId: string): JSX.Element[] {
  switch (providerId) {
    case 'generic_oauth':
      return [<OAuthBadge key="oauth" />];
    case 'azuread':
      return [<OAuthBadge key="oauth" />, <AuthBadge key="name" icon="microsoft" text="azuread" />];
    case 'github':
      return [<OAuthBadge key="oauth" />, <AuthBadge key="name" icon="github" text="github" />];
    case 'gitlab':
      return [<OAuthBadge key="oauth" />, <AuthBadge key="name" icon="gitlab" text="gitlab" />];
    case 'google':
      return [<OAuthBadge key="oauth" />, <AuthBadge key="name" icon="google" text="google" />];
    case 'okta':
      return [<OAuthBadge key="oauth" />, <AuthBadge key="name" icon="okta" text="okta" />];
    case 'grafana_com':
      return [<OAuthBadge key="oauth" />, <AuthBadge key="name" text="grafana.com" />];
    case 'saml':
      return [<AuthBadge key="name" icon="sitemap" text="SAML" />];
    default:
      return [];
  }
  return [];
}

export const getStyles = (theme: GrafanaTheme2) => {
  return {
    container: css`
      // display: grid;
      // grid-template-columns: ${LOGO_SIZE} 1fr ${theme.spacing(3)};
      // grid-template-rows: auto;
      // gap: ${theme.spacing(2)};
      // grid-auto-flow: row;
      background: ${theme.colors.background.secondary};
      border-radius: ${theme.shape.borderRadius()};
      padding: ${theme.spacing(2)};
      transition: ${theme.transitions.create(['background-color', 'box-shadow', 'border-color', 'color'], {
        duration: theme.transitions.duration.short,
      })};

      &:hover {
        background: ${theme.colors.emphasize(theme.colors.background.secondary, 0.03)};
      }
    `,
    header: css`
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: ${theme.spacing(2)};
    `,
    content: css`
      // grid-area: 3 / 1 / 4 / 3;
      color: ${theme.colors.text.secondary};
      margin-bottom: ${theme.spacing(2)};
    `,
    status: css`
      grid-area: 4 / 1 / 4 / 4;
      color: ${theme.colors.text.secondary};
    `,
    disabled: css`
      color: ${theme.colors.text.secondary};
    `,
    name: css`
      // grid-area: 2 / 1 / 3 / 3;
      align-self: center;
      font-size: ${theme.typography.h4.fontSize};
      color: ${theme.colors.text.primary};
      margin: 0;
    `,
    badgesContainer: css`
      display: flex;
      > div {
        margin-right: ${theme.spacing(1)};
      }
    `,
    actionButton: css`
      color: ${theme.colors.text.primary};
    `,
  };
};
