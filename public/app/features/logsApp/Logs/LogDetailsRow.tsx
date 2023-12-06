import { css, cx } from '@emotion/css';
import memoizeOne from 'memoize-one';
import React, { PureComponent, useState } from 'react';

import {
  CoreApp,
  DataFrame,
  Field,
  GrafanaTheme2,
  IconName,
  LinkModel,
  LogRowModel,
} from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { ClipboardButton, DataLinkButton, IconButton, Themeable2, withTheme2 } from '@grafana/ui';
import { getLogRowStyles } from 'app/features/logs/components/getLogRowStyles';
import { logRowToSingleRowDataFrame } from 'app/features/logs/logsModel';

export interface Props extends Themeable2 {
  parsedValues: string[];
  parsedKeys: string[];
  disableActions: boolean;
  wrapLogMessage?: boolean;
  isLabel?: boolean;
  onClickFilterLabel?: (key: string, value: string, frame?: DataFrame) => void;
  onClickFilterOutLabel?: (key: string, value: string, frame?: DataFrame) => void;
  links?: Array<LinkModel<Field>>;
  displayedFields?: string[];
  onClickShowField?: (key: string) => void;
  onClickHideField?: (key: string) => void;
  row: LogRowModel;
  app?: CoreApp;
  isFilterLabelActive?: (key: string, value: string, refId?: string) => Promise<boolean>;
}

const getStyles = memoizeOne((theme: GrafanaTheme2) => {
  return {
    wordBreakAll: css`
      label: wordBreakAll;
      word-break: break-all;
    `,
    copyButton: css`
      & > button {
        color: ${theme.colors.text.secondary};
        padding: 0;
        justify-content: center;
        border-radius: ${theme.shape.radius.circle};
        height: ${theme.spacing(theme.components.height.sm)};
        width: ${theme.spacing(theme.components.height.sm)};
        svg {
          margin: 0;
        }

        span > div {
          top: -5px;
          & button {
            color: ${theme.colors.success.main};
          }
        }
      }
    `,
    adjoiningLinkButton: css`
      padding-top: ${theme.spacing(1)};
    `,
    wrapLine: css`
      label: wrapLine;
      white-space: pre-wrap;
    `,
    logDetailsStats: css`
      padding: 0 ${theme.spacing(1)};
    `,
    logDetailsValue: css`
      display: flex;
      align-items: center;
      line-height: 22px;

      .log-details-value-copy {
        visibility: hidden;
      }
      &:hover {
        .log-details-value-copy {
          visibility: visible;
        }
      }
    `,
    buttonRow: css`
      display: flex;
      flex-direction: row;
      gap: ${theme.spacing(0.5)};
      margin-left: ${theme.spacing(0.5)};
    `,
  };
});

class UnThemedLogDetailsRow extends PureComponent<Props> {
  showField = () => {
    const { onClickShowField: onClickShowDetectedField, parsedKeys, row } = this.props;
    if (onClickShowDetectedField) {
      onClickShowDetectedField(parsedKeys[0]);
    }

    reportInteraction('grafana_explore_logs_log_details_replace_line_clicked', {
      datasourceType: row.datasourceType,
      logRowUid: row.uid,
      type: 'enable',
    });
  };

  hideField = () => {
    const { onClickHideField: onClickHideDetectedField, parsedKeys, row } = this.props;
    if (onClickHideDetectedField) {
      onClickHideDetectedField(parsedKeys[0]);
    }

    reportInteraction('grafana_explore_logs_log_details_replace_line_clicked', {
      datasourceType: row.datasourceType,
      logRowUid: row.uid,
      type: 'disable',
    });
  };

  isFilterLabelActive = async () => {
    const { isFilterLabelActive, parsedKeys, parsedValues, row } = this.props;
    if (isFilterLabelActive) {
      return await isFilterLabelActive(parsedKeys[0], parsedValues[0], row.dataFrame?.refId);
    }
    return false;
  };

  filterLabel = () => {
    const { onClickFilterLabel, parsedKeys, parsedValues, row } = this.props;
    if (onClickFilterLabel) {
      onClickFilterLabel(parsedKeys[0], parsedValues[0], logRowToSingleRowDataFrame(row) || undefined);
    }

    reportInteraction('grafana_explore_logs_log_details_filter_clicked', {
      datasourceType: row.datasourceType,
      filterType: 'include',
      logRowUid: row.uid,
    });
  };

  filterOutLabel = () => {
    const { onClickFilterOutLabel, parsedKeys, parsedValues, row } = this.props;
    if (onClickFilterOutLabel) {
      onClickFilterOutLabel(parsedKeys[0], parsedValues[0], logRowToSingleRowDataFrame(row) || undefined);
    }

    reportInteraction('grafana_explore_logs_log_details_filter_clicked', {
      datasourceType: row.datasourceType,
      filterType: 'exclude',
      logRowUid: row.uid,
    });
  };

  generateClipboardButton(val: string) {
    const { theme } = this.props;
    const styles = getStyles(theme);

    return (
      <div className={`log-details-value-copy ${styles.copyButton}`}>
        <ClipboardButton
          getText={() => val}
          title="Copy value to clipboard"
          fill="text"
          variant="secondary"
          icon="copy"
          size="md"
        />
      </div>
    );
  }

  generateMultiVal(value: string[], showCopy?: boolean) {
    return (
      <table>
        <tbody>
          {value?.map((val, i) => {
            return (
              <tr key={`${val}-${i}`}>
                <td>
                  {val}
                  {showCopy && val !== '' && this.generateClipboardButton(val)}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    );
  }

  render() {
    const {
      theme,
      parsedKeys,
      parsedValues,
      links,
      displayedFields,
      wrapLogMessage,
      onClickFilterLabel,
      onClickFilterOutLabel,
      disableActions,
      row,
    } = this.props;
    const styles = getStyles(theme);
    const rowStyles = getLogRowStyles(theme);
    const singleKey = parsedKeys == null ? false : parsedKeys.length === 1;
    const singleVal = parsedValues == null ? false : parsedValues.length === 1;
    const hasFilteringFunctionality = !disableActions && onClickFilterLabel && onClickFilterOutLabel;
    const refIdTooltip = row.dataFrame?.refId ? ` in query ${row.dataFrame?.refId}` : '';

    const toggleFieldButton =
      displayedFields && parsedKeys != null && displayedFields.includes(parsedKeys[0]) ? (
        <IconButton variant="primary" tooltip="Hide this field" name="eye" onClick={this.hideField} />
      ) : (
        <IconButton tooltip="Show this field instead of the message" name="eye" onClick={this.showField} />
      );

    return (
      <>
        <tr className={rowStyles.logDetailsValue}>
          <td className={rowStyles.logsDetailsIcon}>
            <div className={styles.buttonRow}>
              {hasFilteringFunctionality && (
                <>
                  <AsyncIconButton
                    name="search-plus"
                    onClick={this.filterLabel}
                    isActive={this.isFilterLabelActive}
                    tooltipSuffix={refIdTooltip}
                  />
                  <IconButton
                    name="search-minus"
                    tooltip={`Filter out value${refIdTooltip}`}
                    onClick={this.filterOutLabel}
                  />
                </>
              )}
              {!disableActions && displayedFields && toggleFieldButton}
            </div>
          </td>

          {/* Key - value columns */}
          <td className={rowStyles.logDetailsLabel}>{singleKey ? parsedKeys[0] : this.generateMultiVal(parsedKeys)}</td>
          <td className={cx(styles.wordBreakAll, wrapLogMessage && styles.wrapLine)}>
            <div className={styles.logDetailsValue}>
              {singleVal ? parsedValues[0] : this.generateMultiVal(parsedValues, true)}
              {singleVal && this.generateClipboardButton(parsedValues[0])}
            </div>
          </td>
        </tr>
        <tr>
          <td></td>
          <td></td>
          <td>
            <div className={styles.adjoiningLinkButton}>
              {links?.map((link, i) => (
                <span key={`${link.title}-${i}`}>
                  <DataLinkButton link={link} />
                </span>
              ))}
            </div>
          </td>
        </tr>
      </>
    );
  }
}

interface AsyncIconButtonProps extends Pick<React.ButtonHTMLAttributes<HTMLButtonElement>, 'onClick'> {
  name: IconName;
  isActive(): Promise<boolean>;
  tooltipSuffix: string;
}

const AsyncIconButton = ({ isActive, tooltipSuffix, ...rest }: AsyncIconButtonProps) => {
  const [active, setActive] = useState(false);
  const tooltip = active ? 'Remove filter' : 'Filter for value';

  /**t
   * We purposely want to run this on every render to allow the active state to be updated
   * when log details remains open between updates.
   */
  isActive().then(setActive);

  return <IconButton {...rest} variant={active ? 'primary' : undefined} tooltip={tooltip + tooltipSuffix} />;
};

export const LogDetailsRow = withTheme2(UnThemedLogDetailsRow);
LogDetailsRow.displayName = 'LogDetailsRow';
