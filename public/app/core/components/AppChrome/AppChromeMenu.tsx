import { css } from '@emotion/css';
import { useDialog } from '@react-aria/dialog';
import { FocusScope } from '@react-aria/focus';
import { OverlayContainer, useOverlay } from '@react-aria/overlays';
import React, { useEffect, useRef, useState } from 'react';
import CSSTransition from 'react-transition-group/CSSTransition';

import { GrafanaTheme2 } from '@grafana/data';
import { useStyles2, useTheme2 } from '@grafana/ui';
import { useGrafana } from 'app/core/context/GrafanaContext';
import { KioskMode } from 'app/types';

import { DockedMegaMenu, MENU_WIDTH } from './DockedMegaMenu/DockedMegaMenu';
import { TOP_BAR_LEVEL_HEIGHT } from './types';

interface Props {}

export function AppChromeMenu({}: Props) {
  const theme = useTheme2();
  const { chrome } = useGrafana();
  const state = chrome.useState();
  const searchBarHidden = state.searchBarHidden || state.kioskMode === KioskMode.TV;

  const ref = useRef(null);
  const backdropRef = useRef(null);
  const animationSpeed = theme.transitions.duration.shortest;
  const animationStyles = useStyles2(getAnimStyles, animationSpeed);

  // need this janky state/effect logic to prevent the mega menu from reopening
  // when the user clicks the hamburger icon whilst it is already open
  const [isOpen, setIsOpen] = useState(false);
  const onClose = () => setIsOpen(false);
  useEffect(() => {
    if (state.megaMenuOpen) {
      setIsOpen(true);
    }
  }, [state.megaMenuOpen]);

  const { overlayProps, underlayProps } = useOverlay(
    {
      isDismissable: true,
      isOpen: true,
      onClose,
    },
    ref
  );
  const { dialogProps } = useDialog({}, ref);
  const styles = useStyles2(getStyles, searchBarHidden);

  return (
    <div className={styles.menuWrapper}>
      <OverlayContainer>
        <CSSTransition
          nodeRef={ref}
          in={isOpen}
          unmountOnExit={true}
          classNames={animationStyles.overlay}
          timeout={{ enter: animationSpeed, exit: 0 }}
          onExited={() => chrome.setMegaMenu(false)}
        >
          <FocusScope contain autoFocus>
            <DockedMegaMenu
              className={styles.menuContainer}
              onClose={onClose}
              ref={ref}
              {...overlayProps}
              {...dialogProps}
            />
          </FocusScope>
        </CSSTransition>
        <CSSTransition
          nodeRef={backdropRef}
          in={isOpen}
          unmountOnExit={true}
          classNames={animationStyles.backdrop}
          timeout={{ enter: animationSpeed, exit: 0 }}
        >
          <div ref={backdropRef} className={styles.menuBackdrop} {...underlayProps} />
        </CSSTransition>
      </OverlayContainer>
    </div>
  );
}

const getStyles = (theme: GrafanaTheme2, searchBarHidden?: boolean) => {
  const topPosition = (searchBarHidden ? TOP_BAR_LEVEL_HEIGHT : TOP_BAR_LEVEL_HEIGHT * 2) + 1;

  return {
    menuBackdrop: css({
      backdropFilter: 'blur(1px)',
      backgroundColor: theme.components.overlay.background,
      bottom: 0,
      left: 0,
      position: 'fixed',
      right: 0,
      top: searchBarHidden ? 0 : TOP_BAR_LEVEL_HEIGHT,
      zIndex: theme.zIndex.modalBackdrop,

      [theme.breakpoints.up('md')]: {
        top: topPosition,
      },
    }),
    menuContainer: css({
      display: 'flex',
      bottom: 0,
      flexDirection: 'column',
      left: 0,
      right: 0,
      // Needs to below navbar should we change the navbarFixed? add add a new level?
      zIndex: theme.zIndex.modal,
      position: 'fixed',
      top: searchBarHidden ? 0 : TOP_BAR_LEVEL_HEIGHT,
      backgroundColor: theme.colors.background.primary,
      boxSizing: 'content-box',
      flex: '1 1 0',

      [theme.breakpoints.up('md')]: {
        right: 'unset',
        borderRight: `1px solid ${theme.colors.border.weak}`,
        top: topPosition,
      },
    }),
    menuWrapper: css({
      position: 'fixed',
      display: 'grid',
      gridAutoFlow: 'column',
      height: '100%',
      zIndex: theme.zIndex.sidemenu,
    }),
  };
};

const getAnimStyles = (theme: GrafanaTheme2, animationDuration: number) => {
  const commonTransition = {
    transitionDuration: `${animationDuration}ms`,
    transitionTimingFunction: theme.transitions.easing.easeInOut,
    [theme.breakpoints.down('md')]: {
      overflow: 'hidden',
    },
  };

  const overlayTransition = {
    ...commonTransition,
    transitionProperty: 'box-shadow, width',
    // this is needed to prevent a horizontal scrollbar during the animation on firefox
    '.scrollbar-view': {
      overflow: 'hidden !important',
    },
  };

  const backdropTransition = {
    ...commonTransition,
    transitionProperty: 'opacity',
  };

  const overlayOpen = {
    width: '100%',
    [theme.breakpoints.up('md')]: {
      boxShadow: theme.shadows.z3,
      width: MENU_WIDTH,
    },
  };

  const overlayClosed = {
    boxShadow: 'none',
    width: 0,
  };

  const backdropOpen = {
    opacity: 1,
  };

  const backdropClosed = {
    opacity: 0,
  };

  return {
    backdrop: {
      enter: css(backdropClosed),
      enterActive: css(backdropTransition, backdropOpen),
      enterDone: css(backdropOpen),
    },
    overlay: {
      enter: css(overlayClosed),
      enterActive: css(overlayTransition, overlayOpen),
      enterDone: css(overlayOpen),
    },
  };
};
