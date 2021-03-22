import React, { FC, MouseEvent } from 'react';
import { css, cx } from 'emotion';
import { Icon, IconButton, useStyles } from '@grafana/ui';
import { GrafanaTheme } from '@grafana/data';

import { TagBadge } from '../../core/components/TagFilter/TagBadge';
import { PlaylistItem } from './types';

interface PlaylistTableRowProps {
  first: boolean;
  last: boolean;
  item: PlaylistItem;
  onMoveUp: (item: PlaylistItem) => void;
  onMoveDown: (item: PlaylistItem) => void;
  onDelete: (item: PlaylistItem) => void;
}

export const PlaylistTableRow: FC<PlaylistTableRowProps> = ({ item, onDelete, onMoveDown, onMoveUp, first, last }) => {
  const styles = useStyles(getStyles);
  const onDeleteClick = (event: MouseEvent) => {
    event.preventDefault();
    onDelete(item);
  };
  const onMoveDownClick = (event: MouseEvent) => {
    event.preventDefault();
    onMoveDown(item);
  };
  const onMoveUpClick = (event: MouseEvent) => {
    event.preventDefault();
    onMoveUp(item);
  };

  return (
    <tr>
      {item.type === 'dashboard_by_id' ? (
        <td className={cx(styles.td, styles.item)}>
          <Icon name="apps" />
          <span>{item.title}</span>
        </td>
      ) : null}
      {item.type === 'dashboard_by_tag' ? (
        <td className={cx(styles.td, styles.item)}>
          <Icon name="tag-alt" />
          <TagBadge key={item.id} label={item.title} removeIcon={false} count={0} />
        </td>
      ) : null}
      <td className={cx(styles.td, styles.settings)}>
        {!first ? <IconButton name="arrow-up" size="md" onClick={onMoveUpClick} /> : null}
        {!last ? <IconButton name="arrow-down" size="md" onClick={onMoveDownClick} /> : null}
        <IconButton name="times" size="md" onClick={onDeleteClick} />
      </td>
    </tr>
  );
};

function getStyles(theme: GrafanaTheme) {
  return {
    td: css`
      label: td;
      line-height: 28px;
      max-width: 335px;
      white-space: nowrap;
      text-overflow: ellipsis;
      overflow: hidden;
    `,
    item: css`
      label: item;
      span {
        margin-left: ${theme.spacing.xs};
      }
    `,
    settings: css`
      label: settings;
      text-align: right;
    `,
  };
}
