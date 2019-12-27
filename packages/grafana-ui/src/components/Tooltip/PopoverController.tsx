import React from 'react';
import * as PopperJS from 'popper.js';
import { PopoverContent } from './Tooltip';

// This API allows popovers to update Popper's position when e.g. popover content changes
// updatePopperPosition is delivered to content by react-popper

export interface UsingPopperProps {
  show?: boolean;
  placement?: PopperJS.Placement;
  content: PopoverContent;
  children: JSX.Element;
}

type PopperControllerRenderProp = (
  showPopper: () => void,
  hidePopper: () => void,
  popperProps: {
    show: boolean;
    placement: PopperJS.Placement;
    content: PopoverContent;
  }
) => JSX.Element;

interface Props {
  placement?: PopperJS.Placement;
  content: PopoverContent;
  className?: string;
  children: PopperControllerRenderProp;
  hideAfter?: number;
}

interface State {
  show: boolean;
}

class PopoverController extends React.Component<Props, State> {
  private hideTimeout: any;

  constructor(props: Props) {
    super(props);

    this.state = {
      show: false,
    };
  }

  showPopper = () => {
    if (this.hideTimeout) {
      clearTimeout(this.hideTimeout);
    }

    this.setState(prevState => ({
      ...prevState,
      show: true,
    }));
  };

  hidePopper = () => {
    if (this.props.hideAfter !== 0) {
      this.hideTimeout = setTimeout(() => {
        this.setState(prevState => ({
          ...prevState,
          show: false,
        }));
      }, this.props.hideAfter);
      return;
    }
    this.setState(prevState => ({
      ...prevState,
      show: false,
    }));
  };

  render() {
    const { children, content, placement = 'auto' } = this.props;
    const { show } = this.state;

    return children(this.showPopper, this.hidePopper, {
      show,
      placement,
      content,
    });
  }
}

export { PopoverController };
