import React, { FunctionComponent } from 'react';

import { ColorPickerPopover, ColorPickerProps } from './ColorPickerPopover';
import { PopoverContentProps } from '../Tooltip/Tooltip';
import { Switch } from '../Forms/Legacy/Switch/Switch';
import { css } from '@emotion/css';
import { withTheme, useStyles } from '../../themes';

export interface SeriesColorPickerPopoverProps extends ColorPickerProps, PopoverContentProps {
  yaxis?: number;
  onToggleAxis?: () => void;
}

export const SeriesColorPickerPopover: FunctionComponent<SeriesColorPickerPopoverProps> = (props) => {
  const styles = useStyles(getStyles);
  const { yaxis, onToggleAxis, color, ...colorPickerProps } = props;

  const customPickers = onToggleAxis
    ? {
        yaxis: {
          name: 'Y-Axis',
          tabComponent() {
            return (
              <Switch
                key="yaxisSwitch"
                label="Use right y-axis"
                className={styles.colorPickerAxisSwitch}
                labelClass={styles.colorPickerAxisSwitchLabel}
                checked={yaxis === 2}
                onChange={() => {
                  if (onToggleAxis) {
                    onToggleAxis();
                  }
                }}
              />
            );
          },
        },
      }
    : undefined;
  return <ColorPickerPopover {...colorPickerProps} color={color || '#000000'} customPickers={customPickers} />;
};

interface AxisSelectorProps {
  yaxis: number;
  onToggleAxis?: () => void;
}

interface AxisSelectorState {
  yaxis: number;
}

export class AxisSelector extends React.PureComponent<AxisSelectorProps, AxisSelectorState> {
  constructor(props: AxisSelectorProps) {
    super(props);
    this.state = {
      yaxis: this.props.yaxis,
    };
    this.onToggleAxis = this.onToggleAxis.bind(this);
  }

  onToggleAxis() {
    this.setState({
      yaxis: this.state.yaxis === 2 ? 1 : 2,
    });

    if (this.props.onToggleAxis) {
      this.props.onToggleAxis();
    }
  }

  render() {
    const leftButtonClass = this.state.yaxis === 1 ? 'btn-primary' : 'btn-inverse';
    const rightButtonClass = this.state.yaxis === 2 ? 'btn-primary' : 'btn-inverse';

    return (
      <div className="p-b-1">
        <label className="small p-r-1">Y Axis:</label>
        <button onClick={this.onToggleAxis} className={'btn btn-small ' + leftButtonClass}>
          Left
        </button>
        <button onClick={this.onToggleAxis} className={'btn btn-small ' + rightButtonClass}>
          Right
        </button>
      </div>
    );
  }
}

// This component is to enable SeriesColorPickerPopover usage via series-color-picker-popover directive
export const SeriesColorPickerPopoverWithTheme = withTheme(SeriesColorPickerPopover);

const getStyles = () => {
  return {
    colorPickerAxisSwitch: css`
      width: 100%;
    `,
    colorPickerAxisSwitchLabel: css`
      display: flex;
      flex-grow: 1;
    `,
  };
};
