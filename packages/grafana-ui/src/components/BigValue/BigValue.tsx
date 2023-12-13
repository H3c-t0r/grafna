import { cx } from '@emotion/css';
import React, { PureComponent } from 'react';

import { DisplayValue, DisplayValueAlignmentFactors, FieldSparkline } from '@grafana/data';
import { VizTextDisplayOptions } from '@grafana/schema';

import { Themeable2 } from '../../types';
import { clearButtonStyles } from '../Button';
import { FormattedValueDisplay } from '../FormattedValueDisplay/FormattedValueDisplay';
import { Icon } from '../Icon/Icon';
import { HorizontalGroup } from '../Layout/Layout';

import { buildLayout } from './BigValueLayout';

export enum BigValueColorMode {
  Background = 'background',
  BackgroundSolid = 'background_solid',
  None = 'none',
  Value = 'value',
}

export enum BigValueGraphMode {
  None = 'none',
  Line = 'line',
  Area = 'area',
}

export enum BigValueJustifyMode {
  Auto = 'auto',
  Center = 'center',
}

/**
 * Options for how the value & title are to be displayed
 */
export enum BigValueTextMode {
  Auto = 'auto',
  Value = 'value',
  ValueAndName = 'value_and_name',
  Name = 'name',
  None = 'none',
}

export interface Props extends Themeable2 {
  /** Height of the component */
  height: number;
  /** Width of the component */
  width: number;
  /** Value displayed as Big Value */
  value: DisplayValue;
  /** Sparkline values for showing a graph under/behind the value  */
  sparkline?: FieldSparkline;
  /** onClick handler for the value */
  onClick?: React.MouseEventHandler<HTMLElement>;
  /** Custom styling */
  className?: string;
  /** Color mode for coloring the value or the background */
  colorMode: BigValueColorMode;
  /** Show a graph behind/under the value */
  graphMode: BigValueGraphMode;
  /** Auto justify value and text or center it */
  justifyMode?: BigValueJustifyMode;
  /** Factors that should influence the positioning of the text  */
  alignmentFactors?: DisplayValueAlignmentFactors;
  /** Explicit font size control */
  text?: VizTextDisplayOptions;
  /** Specify which text should be visible in the BigValue */
  textMode?: BigValueTextMode;
  /** If true disables the tooltip */
  hasLinks?: boolean;

  /**
   * If part of a series of stat panes, this is the total number.
   * Used by BigValueTextMode.Auto text mode.
   */
  count?: number;

  /**
   * Disable the wide layout for the BigValue
   */
  disableWideLayout?: boolean;
}

export class BigValue extends PureComponent<Props> {
  static defaultProps: Partial<Props> = {
    justifyMode: BigValueJustifyMode.Auto,
  };

  render() {
    const { onClick, className, hasLinks, theme, textMode } = this.props;
    const layout = buildLayout(this.props);
    const panelStyles = layout.getPanelStyles();
    const valueAndTitleContainerStyles = layout.getValueAndTitleContainerStyles();
    const valueStyles = layout.getValueStyles();
    const titleStyles = layout.getTitleStyles();
    const textValues = layout.textValues;
    const percentChange = this.props.value.percentChange;
    const percentChangeNaN = Number.isNaN(percentChange) || !percentChange;
    const percentChangeString =
      percentChange?.toLocaleString(undefined, { style: 'percent', maximumSignificantDigits: 3 }) ?? '';
    const percentChangeIcon =
      percentChange && (percentChange > 0 ? 'arrow-up' : percentChange < 0 ? 'arrow-down' : undefined);
    const percentChangeStyles = layout.getValueStyles(true);

    // When there is an outer data link this tooltip will override the outer native tooltip
    const tooltip = hasLinks ? undefined : textValues.tooltip;
    const metricHeight = layout.valueFontSize;
    const metricAlignment = panelStyles.flexDirection === 'row' ? 'center' : 'flex-start';
    const iconDim = metricHeight / 2;

    if (!onClick) {
      return (
        <div className={className} style={panelStyles} title={tooltip}>
          <div style={valueAndTitleContainerStyles}>
            {textValues.title && <div style={titleStyles}>{textValues.title}</div>}
            <FormattedValueDisplay value={textValues} style={valueStyles} />
            {!percentChangeNaN && textMode !== BigValueTextMode.None && (
              <div style={percentChangeStyles}>
                <HorizontalGroup height={metricHeight} align={metricAlignment}>
                  {percentChangeIcon && (
                    <Icon name={percentChangeIcon} height={iconDim} width={iconDim} viewBox="6 6 12 12" />
                  )}
                  {percentChangeString}
                </HorizontalGroup>
              </div>
            )}
          </div>
          {layout.renderChart()}
        </div>
      );
    }

    return (
      <button
        type="button"
        className={cx(clearButtonStyles(theme), className)}
        style={panelStyles}
        onClick={onClick}
        title={tooltip}
      >
        <div style={valueAndTitleContainerStyles}>
          {textValues.title && <div style={titleStyles}>{textValues.title}</div>}
          <FormattedValueDisplay value={textValues} style={valueStyles} />
        </div>
        {layout.renderChart()}
      </button>
    );
  }
}
