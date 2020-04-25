import React from 'react';
import { GrafanaTheme, PanelPluginMeta } from '@grafana/data';
import { stylesFactory, useTheme } from '@grafana/ui';
import { css, cx } from 'emotion';
import { e2e } from '@grafana/e2e';
import { PanelPluginBadge } from '../../plugins/PluginSignatureBadge';

interface Props {
  isCurrent: boolean;
  plugin: PanelPluginMeta;
  onClick: () => void;
  disabled: boolean;
}

const VizTypePickerPlugin: React.FC<Props> = ({ isCurrent, plugin, onClick, disabled }) => {
  const theme = useTheme();
  const styles = getStyles(theme);
  const cssClass = cx({
    [styles.item]: true,
    [styles.current]: isCurrent,
    [styles.disabled]: disabled,
  });

  return (
    <div className={styles.wrapper}>
      <div
        className={cssClass}
        onClick={disabled ? () => {} : onClick}
        title={plugin.name}
        aria-label={e2e.components.PluginVisualization.selectors.item(plugin.name)}
      >
        <div className={styles.name} title={plugin.name}>
          {plugin.name}
        </div>
        <img className={styles.img} src={plugin.info.logos.small} />
      </div>
      <div className={styles.badge}>
        <PanelPluginBadge plugin={plugin} />
      </div>
    </div>
  );
};

const getStyles = stylesFactory((theme: GrafanaTheme) => {
  return {
    wrapper: css`
      position: relative;
    `,
    item: css`
      background: ${theme.colors.bg2};
      border: 1px solid ${theme.colors.border2};
      border-radius: 3px;
      height: 100px;
      width: 100%;
      max-width: 200px;
      flex-shrink: 0;
      flex-direction: column;
      text-align: center;
      cursor: pointer;
      display: flex;
      margin-right: 10px;
      align-items: center;
      justify-content: center;
      padding-bottom: 6px;

      &:hover {
        box-shadow: 0 0 4px ${theme.palette.blue95};
        border: 1px solid ${theme.palette.blue95};
      }
    `,
    current: css`
      label: currentVisualizationItem;
      box-shadow: 0 0 6px ${theme.palette.orange} !important;
      border: 1px solid ${theme.palette.orange} !important;
    `,
    disabled: css`
      opacity: 0.2;
      filter: grayscale(1);
      cursor: default;
      &:hover {
        box-shadow: none;
        border: 1px solid ${theme.colors.border2};
      }
    `,
    name: css`
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
      font-size: ${theme.typography.size.sm};
      text-align: center;
      height: 23px;
      font-weight: ${theme.typography.weight.semibold};
      padding: 0 10px;
      width: 100%;
    `,
    img: css`
      height: 55px;
    `,
    badge: css`
      position: absolute;
      bottom: ${theme.spacing.xs};
      right: ${theme.spacing.xs};
    `,
  };
});

export default VizTypePickerPlugin;
