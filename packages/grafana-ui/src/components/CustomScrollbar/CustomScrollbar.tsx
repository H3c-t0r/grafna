import React, { FC, useCallback, useEffect, useRef } from 'react';
import isNil from 'lodash/isNil';
import classNames from 'classnames';
import { css } from '@emotion/css';
import Scrollbars from 'react-custom-scrollbars';
import { useStyles2 } from '../../themes';
import { GrafanaThemeV2 } from '@grafana/data';

interface Props {
  className?: string;
  autoHide?: boolean;
  autoHideTimeout?: number;
  autoHeightMax?: string;
  hideTracksWhenNotNeeded?: boolean;
  hideHorizontalTrack?: boolean;
  hideVerticalTrack?: boolean;
  scrollTop?: number;
  setScrollTop?: (event: any) => void;
  autoHeightMin?: number | string;
  updateAfterMountMs?: number;
}

/**
 * Wraps component into <Scrollbars> component from `react-custom-scrollbars`
 */
export const CustomScrollbar: FC<Props> = ({
  autoHide = false,
  autoHideTimeout = 200,
  setScrollTop,
  className,
  autoHeightMin = '0',
  autoHeightMax = '100%',
  hideTracksWhenNotNeeded = false,
  hideHorizontalTrack,
  hideVerticalTrack,
  updateAfterMountMs,
  scrollTop,
  children,
}) => {
  const ref = useRef<Scrollbars>(null);
  const styles = useStyles2(getStyles);

  const updateScroll = () => {
    if (ref.current && !isNil(scrollTop)) {
      ref.current.scrollTop(scrollTop);
    }
  };

  useEffect(() => {
    updateScroll();
  });

  /**
   * Special logic for doing a update a few milliseconds after mount to check for
   * updated height due to dynamic content
   */

  useEffect(() => {
    if (!updateAfterMountMs) {
      return;
    }
    setTimeout(() => {
      const scrollbar = ref.current as any;
      if (scrollbar?.update) {
        scrollbar.update();
      }
    }, updateAfterMountMs);
  }, [updateAfterMountMs]);

  function renderTrack(className: string, hideTrack: boolean | undefined, passedProps: any) {
    if (passedProps.style && hideTrack) {
      passedProps.style.display = 'none';
    }

    return <div {...passedProps} className={className} />;
  }

  const renderTrackHorizontal = useCallback(
    (passedProps: any) => {
      return renderTrack('track-horizontal', hideHorizontalTrack, passedProps);
    },
    [hideHorizontalTrack]
  );

  const renderTrackVertical = useCallback(
    (passedProps: any) => {
      return renderTrack('track-vertical', hideVerticalTrack, passedProps);
    },
    [hideVerticalTrack]
  );

  const renderThumbHorizontal = useCallback((passedProps: any) => {
    return <div {...passedProps} className="thumb-horizontal" />;
  }, []);

  const renderThumbVertical = useCallback((passedProps: any) => {
    return <div {...passedProps} className="thumb-vertical" />;
  }, []);

  const renderView = useCallback((passedProps: any) => {
    return <div {...passedProps} className="scrollbar-view" />;
  }, []);

  return (
    <Scrollbars
      ref={ref}
      className={classNames(styles.customScrollbar, className)}
      onScroll={setScrollTop}
      autoHeight={true}
      autoHide={autoHide}
      autoHideTimeout={autoHideTimeout}
      hideTracksWhenNotNeeded={hideTracksWhenNotNeeded}
      // These autoHeightMin & autoHeightMax options affect firefox and chrome differently.
      // Before these where set to inherit but that caused problems with cut of legends in firefox
      autoHeightMax={autoHeightMax}
      autoHeightMin={autoHeightMin}
      renderTrackHorizontal={renderTrackHorizontal}
      renderTrackVertical={renderTrackVertical}
      renderThumbHorizontal={renderThumbHorizontal}
      renderThumbVertical={renderThumbVertical}
      renderView={renderView}
    >
      {children}
    </Scrollbars>
  );
};

export default CustomScrollbar;

const getStyles = (theme: GrafanaThemeV2) => {
  return {
    customScrollbar: css`
      // Fix for Firefox. For some reason sometimes .view container gets a height of its content, but in order to
      // make scroll working it should fit outer container size (scroll appears only when inner container size is
      // greater than outer one).
      display: flex;
      flex-grow: 1;
      .scrollbar-view {
        display: flex;
        flex-grow: 1;
        flex-direction: column;
      }
      .track-vertical {
        border-radius: ${theme.shape.borderRadius(2)};
        width: ${theme.spacing(1)} !important;
        right: 0px;
        bottom: ${theme.spacing(0.25)};
        top: ${theme.spacing(0.25)};
      }
      .track-horizontal {
        border-radius: ${theme.shape.borderRadius(2)};
        height: ${theme.spacing(1)} !important;
        right: ${theme.spacing(0.25)};
        bottom: ${theme.spacing(0.25)};
        left: ${theme.spacing(0.25)};
      }
      .thumb-vertical {
        background: ${theme.palette.action.focus};
        border-radius: ${theme.shape.borderRadius(2)};
        opacity: 0;
      }
      .thumb-horizontal {
        background: ${theme.palette.action.focus};
        border-radius: ${theme.shape.borderRadius(2)};
        opacity: 0;
      }
      &:hover {
        .thumb-vertical,
        .thumb-horizontal {
          opacity: 1;
          transition: opacity 0.3s ease-in-out;
        }
      }
    `,
  };
};
