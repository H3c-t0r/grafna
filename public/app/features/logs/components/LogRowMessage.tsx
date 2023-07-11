import React, { useCallback, useMemo, useState } from 'react';
import Highlighter from 'react-highlight-words';

import { CoreApp, findHighlightChunksInText, LogRowModel } from '@grafana/data';

import { LogMessageAnsi } from './LogMessageAnsi';
import { LogRowMenuCell } from './LogRowMenuCell';
import { LogRowStyles } from './getLogRowStyles';

export const MAX_CHARACTERS = 100000;

interface Props {
  row: LogRowModel;
  wrapLogMessage: boolean;
  prettifyLogMessage: boolean;
  app?: CoreApp;
  showContextToggle?: (row?: LogRowModel) => boolean;
  onOpenContext: (row: LogRowModel) => void;
  onPermalinkClick?: (row: LogRowModel) => Promise<void>;
  onPinLine?: (row: LogRowModel) => void;
  onUnpinLine?: (row: LogRowModel) => void;
  pinned?: boolean;
  styles: LogRowStyles;
}

interface LogMessageProps {
  hasAnsi: boolean;
  entry: string;
  highlights: string[] | undefined;
  styles: LogRowStyles;
}

const LogMessage = ({ hasAnsi, entry, highlights, styles }: LogMessageProps) => {
  const needsHighlighter =
    highlights && highlights.length > 0 && highlights[0] && highlights[0].length > 0 && entry.length < MAX_CHARACTERS;
  const searchWords = highlights ?? [];
  if (hasAnsi) {
    const highlight = needsHighlighter ? { searchWords, highlightClassName: styles.logsRowMatchHighLight } : undefined;
    return <LogMessageAnsi value={entry} highlight={highlight} />;
  } else if (needsHighlighter) {
    return (
      <Highlighter
        textToHighlight={entry}
        searchWords={searchWords}
        findChunks={findHighlightChunksInText}
        highlightClassName={styles.logsRowMatchHighLight}
      />
    );
  }
  return <>{entry}</>;
};

const restructureLog = (line: string, prettifyLogMessage: boolean): string => {
  if (prettifyLogMessage) {
    try {
      return JSON.stringify(JSON.parse(line), undefined, 2);
    } catch (error) {
      return line;
    }
  }
  return line;
};

export const LogRowMessage = React.memo((props: Props) => {
  const {
    row,
    wrapLogMessage,
    prettifyLogMessage,
    showContextToggle,
    styles,
    onOpenContext,
    onPermalinkClick,
    onUnpinLine,
    onPinLine,
    pinned,
  } = props;
  const { hasAnsi, raw } = row;
  const restructuredEntry = useMemo(() => restructureLog(raw, prettifyLogMessage), [raw, prettifyLogMessage]);
  const [hover, setHover] = useState(false);
  const showMenu = useCallback(() => {
    setHover(true);
  }, []);
  const hideMenu = useCallback(() => {
    setHover(false);
  }, []);
  const shouldShowMenu = useMemo(() => hover || pinned, [hover, pinned]);

  return (
    <>
      {
        // When context is open, the position has to be NOT relative. // Setting the postion as inline-style to
        // overwrite the more sepecific style definition from `styles.logsRowMessage`.
      }
      <td className={styles.logsRowMessage} onMouseEnter={showMenu} onMouseLeave={hideMenu}>
        <div className={wrapLogMessage ? styles.positionRelative : styles.horizontalScroll}>
          <button className={`${styles.logLine} ${styles.positionRelative}`}>
            <LogMessage hasAnsi={hasAnsi} entry={restructuredEntry} highlights={row.searchWords} styles={styles} />
          </button>
        </div>
        {shouldShowMenu && (
          <LogRowMenuCell
            logText={restructuredEntry}
            row={row}
            showContextToggle={showContextToggle}
            onOpenContext={onOpenContext}
            onPermalinkClick={onPermalinkClick}
            onPinLine={onPinLine}
            onUnpinLine={onUnpinLine}
            pinned={pinned}
            styles={styles}
            mouseIsOver={hover}
          />
        )}
      </td>
    </>
  );
});

LogRowMessage.displayName = 'LogRowMessage';
