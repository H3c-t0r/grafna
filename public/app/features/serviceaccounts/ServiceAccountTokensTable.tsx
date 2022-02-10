import React, { FC } from 'react';
import { DeleteButton, Icon, Tooltip, useTheme2 } from '@grafana/ui';
import { dateTimeFormat, GrafanaTheme2, TimeZone } from '@grafana/data';

import { ApiKey } from '../../types';
import { css } from '@emotion/css';

interface Props {
  tokens: ApiKey[];
  timeZone: TimeZone;
  onDelete: (token: ApiKey) => void;
}

export const ServiceAccountTokensTable: FC<Props> = ({ tokens, timeZone, onDelete }) => {
  const theme = useTheme2();
  const styles = getStyles(theme);

  return (
    <>
      <table className="filter-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Expires</th>
            <th>Created</th>
            <th style={{ width: '34px' }} />
          </tr>
        </thead>
        <tbody>
          {tokens.map((key) => {
            return (
              <tr key={key.id} className={styles.tableRow(key.hasExpired)}>
                <td>{key.name}</td>
                <td>
                  <TokenExpiration timeZone={timeZone} token={key} />
                </td>
                <td>{formatDate(timeZone, key.created)}</td>
                <td>
                  <DeleteButton aria-label="Delete API key" size="sm" onConfirm={() => onDelete(key)} />
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </>
  );
};

function formatDate(timeZone: TimeZone, expiration?: string): string {
  if (!expiration) {
    return 'No expiration date';
  }
  return dateTimeFormat(expiration, { timeZone });
}

function formatSecondsLeftUntilExpiration(secondsUntilExpiration: number): string {
  const days = Math.floor(secondsUntilExpiration / (3600 * 24));
  const daysFormat = days > 1 ? `${days} days` : `${days} day`;
  return `Expires in ${daysFormat}`;
}

interface TokenExpirationProps {
  timeZone: TimeZone;
  token: ApiKey;
}

const TokenExpiration = ({ timeZone, token }: TokenExpirationProps) => {
  const theme = useTheme2();
  const styles = getStyles(theme);
  if (!token.expiration) {
    return <span className={styles.neverExpire}>Never</span>;
  }
  if (token.secondsUntilExpiration) {
    return (
      <span className={styles.secondsUntilExpiration}>
        {formatSecondsLeftUntilExpiration(token.secondsUntilExpiration)}
      </span>
    );
  }
  if (token.hasExpired) {
    return (
      <span className={styles.hasExpired}>
        Expired
        <span className={styles.tooltipContainer}>
          <Tooltip content="This API key has expired.">
            <Icon name="exclamation-triangle" style={{ color: `${theme.colors.error.text}` }} />
          </Tooltip>
        </span>
      </span>
    );
  }
  return <span>{formatDate(timeZone, token.expiration)}</span>;
};

const getStyles = (theme: GrafanaTheme2) => ({
  tableRow: (hasExpired: boolean) => css`
    color: ${hasExpired ? theme.colors.text.secondary : theme.colors.text.primary};
  `,
  tooltipContainer: css`
    margin-left: ${theme.spacing(1)};
  `,
  secondsUntilExpiration: css`
    color: ${theme.colors.warning.text};
  `,
  hasExpired: css`
    color: ${theme.colors.error.text};
  `,
  neverExpire: css`
    color: ${theme.colors.text.secondary};
  `,
});
