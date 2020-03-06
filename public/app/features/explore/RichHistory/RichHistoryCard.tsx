import React, { FunctionComponent, useState } from 'react';
import { css, cx } from 'emotion';
import { stylesFactory, useTheme, Forms } from '@grafana/ui';
import { GrafanaTheme, AppEvents } from '@grafana/data';
import { RichHistoryQuery } from 'app/types/explore';
import { copyStringToClipboard, createUrlFromRichHistory } from '../../../core/utils/richHistory';
import appEvents from 'app/core/app_events';

interface Props {
  query: RichHistoryQuery;
  onChangeRichHistoryProperty: (ts: number, property: string, comment?: string) => void;
}

const getStyles = stylesFactory((theme: GrafanaTheme) => {
  const bgColor = theme.isLight ? theme.colors.gray5 : theme.colors.dark4;
  const cardColor = theme.isLight ? theme.colors.white : theme.colors.dark7;
  return {
    queryCard: css`
      display: flex;
      border: 1px solid ${bgColor};
      padding: ${theme.spacing.sm};
      margin: ${theme.spacing.sm} 0;
      box-shadow: 0px 2px 2px ${bgColor};
      background-color: ${cardColor};
      border-radius: ${theme.border.radius};
      .starred {
        color: ${theme.colors.blue77};
      }
    `,
    queryCardLeft: css`
      padding-right: 10px;
      width: calc(100% - 150px);
    `,
    queryCardRight: css`
      width: 150px;
      display: flex;
      justify-content: flex-end;
      i {
        font-size: ${theme.typography.size.lg};
        font-weight: ${theme.typography.weight.bold};
        margin: 3px;
        cursor: pointer;
      }
    `,
    queryRow: css`
      border-top: 2px solid ${bgColor};
      font-weight: ${theme.typography.weight.bold};
      word-break: break-all;
      padding: 4px 2px;
      :first-child {
        border-top: none;
        padding: 0 0 4px 0;
      }
    `,
    input: css`
      width: 100%;
      &:focus,
      &:active {
        outline: none;
      }
      &:placeholder {
        font-style: italic;
      }
    `,
    buttonRow: css`
      > * {
        margin-right: ${theme.spacing.xs};
      }
    `,
  };
});

export const RichHistoryCard: FunctionComponent<Props> = ({ query, onChangeRichHistoryProperty }) => {
  const [starred, setStared] = useState(query.starred);
  const [activeUpdateComment, setActiveUpdateComment] = useState(false);
  const [comment, setComment] = useState<string | undefined>(query.comment);

  const toggleActiveUpdateComment = () => setActiveUpdateComment(!activeUpdateComment);
  const theme = useTheme();
  const styles = getStyles(theme);

  return (
    <div className={styles.queryCard}>
      <div className={styles.queryCardLeft}>
        {query.queries.map((q, i) => {
          return (
            <div key={`${q}-${i}`} className={styles.queryRow}>
              {q}
            </div>
          );
        })}
        {!activeUpdateComment && query.comment && (
          <div
            className={css`
              overflow-wrap: break-word;
            `}
          >
            {query.comment}
          </div>
        )}
        {activeUpdateComment && (
          <div>
            <Forms.Input
              className={styles.input}
              value={comment}
              placeholder={comment ? undefined : 'add comment'}
              onChange={e => setComment(e.currentTarget.value)}
            />
            <div className={styles.buttonRow}>
              <Forms.Button
                onClick={e => {
                  e.preventDefault();
                  onChangeRichHistoryProperty(query.ts, 'comment', comment);
                  toggleActiveUpdateComment();
                }}
              >
                Save
              </Forms.Button>
              <Forms.Button
                variant="secondary"
                className={css`
                  margin-left: 8px;
                `}
                onClick={() => {
                  toggleActiveUpdateComment();
                  setComment(query.comment);
                }}
              >
                Cancel
              </Forms.Button>
            </div>
          </div>
        )}
      </div>
      <div className={styles.queryCardRight}>
        <i
          className="fa fa-fw fa-comment-o"
          onClick={() => {
            toggleActiveUpdateComment();
          }}
        ></i>
        <i
          className="fa fa-fw fa-copy"
          onClick={() => {
            const queries = query.queries.join('\n\n');
            copyStringToClipboard(queries);
            appEvents.emit(AppEvents.alertSuccess, ['Query copied to clipboard']);
          }}
        ></i>
        <i
          className="fa fa-fw fa-link"
          onClick={() => {
            const url = createUrlFromRichHistory(query);
            copyStringToClipboard(url);
            appEvents.emit(AppEvents.alertSuccess, ['Link copied to clipboard']);
          }}
          style={{ fontWeight: 'normal' }}
        ></i>
        <i
          className={cx('fa fa-fw', starred ? 'fa-star starred' : 'fa-star-o')}
          onClick={() => {
            onChangeRichHistoryProperty(query.ts, 'starred');
            setStared(!starred);
          }}
        ></i>
      </div>
    </div>
  );
};
