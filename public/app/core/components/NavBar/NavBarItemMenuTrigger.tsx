import React, { ReactElement, useEffect, useState } from 'react';
import { css, cx } from '@emotion/css';
import { Icon, IconName, Link, useTheme2 } from '@grafana/ui';
import { GrafanaTheme2, NavModelItem } from '@grafana/data';
import { MenuTriggerProps } from '@react-types/menu';
import { useMenuTriggerState } from '@react-stately/menu';
import { useMenuTrigger } from '@react-aria/menu';
import { useFocusWithin, useHover, useKeyboard } from '@react-aria/interactions';
import { useButton } from '@react-aria/button';
import { useDialog } from '@react-aria/dialog';
import { DismissButton, OverlayContainer, useOverlay, useOverlayPosition } from '@react-aria/overlays';
import { FocusScope } from '@react-aria/focus';

import { NavBarItemMenuContext, useNavBarContext } from './context';
import { NavFeatureHighlight } from './NavFeatureHighlight';
import { reportExperimentView } from '@grafana/runtime';

export interface NavBarItemMenuTriggerProps extends MenuTriggerProps {
  children: ReactElement;
  item: NavModelItem;
  isActive?: boolean;
  label: string;
  reverseMenuDirection: boolean;
}

export function NavBarItemMenuTrigger(props: NavBarItemMenuTriggerProps): ReactElement {
  const { item, isActive, label, children: menu, reverseMenuDirection, ...rest } = props;
  const [menuHasFocus, setMenuHasFocus] = useState(false);
  const { menuIdOpen, setMenuIdOpen } = useNavBarContext();
  const theme = useTheme2();
  const styles = getStyles(theme, isActive);

  // Create state based on the incoming props
  const state = useMenuTriggerState({ ...rest });

  // Get props for the menu trigger and menu elements
  const ref = React.useRef<HTMLElement>(null);
  const { menuTriggerProps, menuProps } = useMenuTrigger({}, state, ref);

  useEffect(() => {
    if (item.highlightId) {
      reportExperimentView(`feature-highlights-${item.highlightId}-nav`, 'test', '');
    }
  }, [item.highlightId]);

  const { hoverProps } = useHover({
    onHoverChange: (isHovering) => {
      if (isHovering) {
        state.open();
        setMenuIdOpen(ref.current?.id || null);
      } else {
        state.close();
      }
    },
  });

  useEffect(() => {
    if (menuIdOpen !== ref.current?.id) {
      state.close();
      setMenuHasFocus(false);
    } else {
      state.open();
    }
  }, [menuIdOpen]);

  const { keyboardProps } = useKeyboard({
    onKeyDown: (e) => {
      switch (e.key) {
        case 'ArrowRight':
          if (!state.isOpen) {
            state.open();
          }
          setMenuHasFocus(true);
          break;
        case 'Tab':
          setMenuIdOpen(null);
          break;
        default:
          break;
      }
    },
  });

  // Get props for the button based on the trigger props from useMenuTrigger
  const { buttonProps } = useButton(menuTriggerProps, ref);
  const Wrapper = item.highlightText ? NavFeatureHighlight : React.Fragment;
  const itemContent = (
    <Wrapper>
      <span className={styles.icon}>
        {item?.icon && <Icon name={item.icon as IconName} size="xl" />}
        {item?.img && <img src={item.img} alt={`${item.text} logo`} />}
      </span>
    </Wrapper>
  );
  let element = (
    <button
      className={styles.element}
      {...buttonProps}
      {...keyboardProps}
      ref={ref as React.RefObject<HTMLButtonElement>}
      onClick={item?.onClick}
      aria-label={label}
    >
      {itemContent}
    </button>
  );

  if (item?.url) {
    element =
      !item.target && item.url.startsWith('/') ? (
        <Link
          {...buttonProps}
          {...keyboardProps}
          ref={ref as React.RefObject<HTMLAnchorElement>}
          href={item.url}
          target={item.target}
          onClick={item?.onClick}
          className={styles.element}
          aria-label={label}
        >
          {itemContent}
        </Link>
      ) : (
        <a
          href={item.url}
          target={item.target}
          onClick={item?.onClick}
          {...buttonProps}
          {...keyboardProps}
          ref={ref as React.RefObject<HTMLAnchorElement>}
          className={styles.element}
          aria-label={label}
        >
          {itemContent}
        </a>
      );
  }

  const overlayRef = React.useRef<HTMLDivElement>(null);
  const { dialogProps } = useDialog({}, overlayRef);
  const { overlayProps } = useOverlay(
    {
      onClose: () => state.close(),
      isOpen: state.isOpen,
      isDismissable: true,
    },
    overlayRef
  );
  let { overlayProps: overlayPositionProps } = useOverlayPosition({
    targetRef: ref,
    overlayRef,
    placement: reverseMenuDirection ? 'right bottom' : 'right top',
    isOpen: state.isOpen,
  });

  const { focusWithinProps } = useFocusWithin({
    onFocusWithin: (e) => {
      if (e.target.id === ref.current?.id) {
        // If focussing on the trigger itself, set the menu id that is open
        setMenuIdOpen(ref.current?.id);
        state.open();
      }
    },
    onBlurWithin: (e) => {
      if (e.target?.getAttribute('role') === 'menuitem' && !overlayRef.current?.contains(e.relatedTarget)) {
        // If it is blurring from a menuitem to an element outside the current overlay
        // close the menu that is open
        setMenuIdOpen(null);
      }
    },
  });

  return (
    <div className={cx(styles.element, 'dropdown')} {...focusWithinProps} {...hoverProps}>
      {element}
      {state.isOpen && (
        <OverlayContainer>
          <NavBarItemMenuContext.Provider
            value={{
              menuProps,
              menuHasFocus,
              onClose: () => state.close(),
              onLeft: () => {
                setMenuHasFocus(false);
                ref.current?.focus();
              },
            }}
          >
            <FocusScope restoreFocus>
              <div {...overlayProps} {...overlayPositionProps} {...dialogProps} {...hoverProps} ref={overlayRef}>
                <DismissButton onDismiss={() => state.close()} />
                {menu}
                <DismissButton onDismiss={() => state.close()} />
              </div>
            </FocusScope>
          </NavBarItemMenuContext.Provider>
        </OverlayContainer>
      )}
    </div>
  );
}

const getStyles = (theme: GrafanaTheme2, isActive?: boolean) => ({
  element: css`
    background-color: transparent;
    border: none;
    color: inherit;
    display: block;
    line-height: ${theme.components.sidemenu.width}px;
    padding: 0;
    text-align: center;
    width: ${theme.components.sidemenu.width}px;

    &::before {
      display: ${isActive ? 'block' : 'none'};
      content: ' ';
      position: absolute;
      left: 0;
      top: 0;
      bottom: 0;
      width: 4px;
      border-radius: 2px;
      background-image: ${theme.colors.gradients.brandVertical};
    }

    &:focus-visible {
      background-color: ${theme.colors.action.hover};
      box-shadow: none;
      color: ${theme.colors.text.primary};
      outline: 2px solid ${theme.colors.primary.main};
      outline-offset: -2px;
      transition: none;
    }
  `,
  icon: css`
    height: 100%;
    width: 100%;

    img {
      border-radius: 50%;
      height: ${theme.spacing(3)};
      width: ${theme.spacing(3)};
    }
  `,
});
