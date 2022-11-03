import { css, cx } from '@emotion/css';
import React, { ReactNode, useState } from 'react';

import { DataQuery, DataSourceInstanceSettings, GrafanaTheme2 } from '@grafana/data';
import { selectors } from '@grafana/e2e-selectors';
import { DataSourcePicker } from '@grafana/runtime';
import { Icon, Input, FieldValidationMessage, useStyles2 } from '@grafana/ui';

export interface Props<TQuery extends DataQuery = DataQuery> {
  query: TQuery;
  queries: TQuery[];
  disabled?: boolean;
  dataSource: DataSourceInstanceSettings;
  renderExtras?: () => ReactNode;
  onChangeDataSource?: (settings: DataSourceInstanceSettings) => void;
  onChange: (query: TQuery) => void;
  onClick: (e: React.MouseEvent) => void;
  collapsedText: string | null;
  alerting?: boolean;
}

export const QueryEditorRowHeader = <TQuery extends DataQuery>(props: Props<TQuery>) => {
  const { query, queries, onClick, onChange, collapsedText, renderExtras, disabled } = props;

  const styles = useStyles2(getStyles);
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [validationError, setValidationError] = useState<string | null>(null);

  const onEditQuery = (event: React.SyntheticEvent) => {
    setIsEditing(true);
  };

  const onEndEditName = (newName: string) => {
    setIsEditing(false);

    // Ignore change if invalid
    if (validationError) {
      setValidationError(null);
      return;
    }

    if (query.refId !== newName) {
      onChange({
        ...query,
        refId: newName,
      });
    }
  };

  const onInputChange = (event: React.SyntheticEvent<HTMLInputElement>) => {
    const newName = event.currentTarget.value.trim();

    if (newName.length === 0) {
      setValidationError('An empty query name is not allowed');
      return;
    }

    for (const otherQuery of queries) {
      if (otherQuery !== query && newName === otherQuery.refId) {
        setValidationError('Query name already exists');
        return;
      }
    }

    if (validationError) {
      setValidationError(null);
    }
  };

  const onEditQueryBlur = (event: React.SyntheticEvent<HTMLInputElement>) => {
    onEndEditName(event.currentTarget.value.trim());
  };

  const onKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter') {
      onEndEditName((event.target as any).value);
    }
  };

  const onFocus = (event: React.FocusEvent<HTMLInputElement>) => {
    event.target.select();
  };

  return (
    <>
      <div className={styles.wrapper}>
        {!isEditing && (
          <button
            className={styles.queryNameWrapper}
            aria-label={selectors.components.QueryEditorRow.title(query.refId)}
            title="Edit query name"
            onClick={onEditQuery}
            data-testid="query-name-div"
            type="button"
          >
            <span className={styles.queryName}>{query.refId}</span>
            <Icon name="pen" className={styles.queryEditIcon} size="sm" />
          </button>
        )}

        {isEditing && (
          <>
            <Input
              type="text"
              defaultValue={query.refId}
              onBlur={onEditQueryBlur}
              autoFocus
              onKeyDown={onKeyDown}
              onFocus={onFocus}
              invalid={validationError !== null}
              onChange={onInputChange}
              className={styles.queryNameInput}
              data-testid="query-name-input"
            />
            {validationError && <FieldValidationMessage horizontal>{validationError}</FieldValidationMessage>}
          </>
        )}
        {renderDataSource(props, styles)}
        {renderExtras && <div className={styles.itemWrapper}>{renderExtras()}</div>}
        {disabled && <em className={styles.contextInfo}>Disabled</em>}
      </div>

      {collapsedText && (
        <div className={styles.collapsedText} onClick={onClick}>
          {collapsedText}
        </div>
      )}
    </>
  );
};

const renderDataSource = <TQuery extends DataQuery>(
  props: Props<TQuery>,
  styles: ReturnType<typeof getStyles>
): ReactNode => {
  const { alerting, dataSource, onChangeDataSource } = props;

  if (!onChangeDataSource) {
    return <em className={styles.contextInfo}>({dataSource.name})</em>;
  }

  return (
    <div className={styles.itemWrapper}>
      <DataSourcePicker variables={true} alerting={alerting} current={dataSource.name} onChange={onChangeDataSource} />
    </div>
  );
};

const getStyles = (theme: GrafanaTheme2) => {
  return {
    wrapper: css`
      label: Wrapper;
      display: flex;
      align-items: center;
      margin-left: ${theme.spacing(0.5)};
    `,
    queryNameWrapper: css`
      display: flex;
      cursor: pointer;
      border: 1px solid transparent;
      border-radius: ${theme.shape.borderRadius(2)};
      align-items: center;
      padding: 0 0 0 ${theme.spacing(0.5)};
      margin: 0;
      background: transparent;

      &:hover {
        background: ${theme.colors.action.hover};
        border: 1px dashed ${theme.colors.border.strong};
      }

      &:focus {
        border: 2px solid ${theme.colors.primary.border};
      }

      &:hover,
      &:focus {
        .query-name-edit-icon {
          visibility: visible;
        }
      }
    `,
    queryName: css`
      font-weight: ${theme.typography.fontWeightMedium};
      color: ${theme.colors.primary.text};
      cursor: pointer;
      overflow: hidden;
      margin-left: ${theme.spacing(0.5)};
    `,
    queryEditIcon: cx(
      css`
        margin-left: ${theme.spacing(2)};
        visibility: hidden;
      `,
      'query-name-edit-icon'
    ),
    queryNameInput: css`
      max-width: 300px;
      margin: -4px 0;
    `,
    collapsedText: css`
      font-weight: ${theme.typography.fontWeightRegular};
      font-size: ${theme.typography.bodySmall.fontSize};
      color: ${theme.colors.text.secondary};
      padding-left: ${theme.spacing(1)};
      align-items: center;
      overflow: hidden;
      font-style: italic;
      white-space: nowrap;
      text-overflow: ellipsis;
    `,
    contextInfo: css`
      font-size: ${theme.typography.bodySmall.fontSize};
      font-style: italic;
      color: ${theme.colors.text.secondary};
      padding-left: 10px;
    `,
    itemWrapper: css`
      display: flex;
      margin-left: 4px;
    `,
  };
};
