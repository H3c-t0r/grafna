import { css } from '@emotion/css';
import React from 'react';
import { useToggle } from 'react-use';

import { GrafanaTheme2 } from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { useStyles2, PanelContainer, CustomScrollbar } from '@grafana/ui';

import { useContentOutlineContext } from './ContentOutlineContext';
import { ContentOutlineItemButton } from './ContentOutlineItemButton';

const getStyles = (theme: GrafanaTheme2) => {
  return {
    wrapper: css({
      label: 'wrapper',
      position: 'relative',
      display: 'flex',
      justifyContent: 'center',
      marginRight: theme.spacing(1),
      height: '100%',
      backgroundColor: theme.colors.background.primary,
    }),
    content: css({
      label: 'content',
      top: 0,
    }),
    buttonStyles: css({
      textAlign: 'right',
      width: '100%',
      padding: theme.spacing(0, 1.5),
      '&:hover': {
        color: theme.colors.text.primary,
        textDecoration: 'underline',
      },
    }),
    iconButton: css({
      justifyContent: 'center',
      width: theme.spacing(4),
      '&:hover': {
        color: theme.colors.text.primary,
        background: theme.colors.background.secondary,
        textDecoration: 'underline',
      },
    }),
    sectionWrapper: css({
      display: 'flex',
    }),
    indent: css({
      width: '100%',
      paddingLeft: '68px',
      '&:hover': {
        color: theme.colors.text.primary,
        textDecoration: 'underline',
      },
    }),
  };
};

export function ContentOutline({ scroller, panelId }: { scroller: HTMLElement | undefined; panelId: string }) {
  const [contentOutlineExpanded, toggleContentOutlineExpanded] = useToggle(false);
  const [sectionExpanded, toggleSectionExpanded] = useToggle(false);
  const styles = useStyles2((theme) => getStyles(theme));
  const { outlineItems } = useContentOutlineContext();

  const scrollIntoView = (ref: HTMLElement | null, itemPanelId: string) => {
    let scrollValue = 0;
    let el: HTMLElement | null | undefined = ref;

    // This is to handle ContentOutlineItem wrapping each QueryEditorRow
    const customOffsetTop = itemPanelId === 'Queries' ? -10 : 0;

    do {
      scrollValue += el?.offsetTop || customOffsetTop;
      el = el?.offsetParent as HTMLElement;
    } while (el && el !== scroller);

    scroller?.scroll({
      top: scrollValue,
      behavior: 'smooth',
    });
    reportInteraction('explore_toolbar_contentoutline_clicked', {
      item: 'select_section',
      type: itemPanelId,
    });
  };

  const toggle = () => {
    toggleContentOutlineExpanded();
    reportInteraction('explore_toolbar_contentoutline_clicked', {
      item: 'outline',
      type: contentOutlineExpanded ? 'minimize' : 'expand',
    });
  };

  const toggleSection = () => {
    toggleSectionExpanded();
    reportInteraction('explore_toolbar_contentoutline_clicked', {
      item: 'section',
      type: sectionExpanded ? 'minimize' : 'expand',
    });
  };

  console.log('outlineItems', outlineItems);

  const outlineItemsHaveChildren = outlineItems.some((item) => item.children);

  console.log('outlineItemsHaveChildren', outlineItemsHaveChildren);

  // TODO: fix indenting for buttons that are not in a section
  // figure out why ellipsis is not working

  return (
    <PanelContainer className={styles.wrapper} id={panelId}>
      <CustomScrollbar>
        <div className={styles.content}>
          <ContentOutlineItemButton
            title={contentOutlineExpanded ? 'Collapse outline' : undefined}
            icon={contentOutlineExpanded ? 'angle-left' : 'angle-right'}
            onClick={toggle}
            tooltip={!contentOutlineExpanded ? 'Expand content outline' : undefined}
            className={styles.buttonStyles}
            aria-expanded={contentOutlineExpanded}
          />

          {outlineItems.map((item) => (
            <React.Fragment key={item.id}>
              <div className={styles.sectionWrapper}>
                {item.children && item.children.length > 0 && contentOutlineExpanded && (
                  <ContentOutlineItemButton
                    icon={sectionExpanded ? 'angle-down' : 'angle-right'}
                    className={styles.iconButton}
                    onClick={toggleSection}
                  />
                )}
                <ContentOutlineItemButton
                  key={item.id}
                  title={contentOutlineExpanded ? item.title : undefined}
                  className={styles.buttonStyles}
                  icon={item.icon}
                  onClick={() => scrollIntoView(item.ref, item.panelId)}
                  tooltip={!contentOutlineExpanded ? item.title : undefined}
                  topLeftIcon={
                    item.children && item.children.length > 0 && !contentOutlineExpanded
                      ? sectionExpanded
                        ? 'angle-down'
                        : 'angle-right'
                      : undefined
                  }
                  onTopLeftIconClick={() => toggleSection()}
                />
              </div>
              {item.children &&
                sectionExpanded &&
                item.children.map((child) => (
                  <ContentOutlineItemButton
                    key={child.id}
                    title={contentOutlineExpanded ? child.title : undefined}
                    icon={contentOutlineExpanded ? undefined : item.icon}
                    className={contentOutlineExpanded ? styles.indent : styles.buttonStyles}
                    onClick={() => scrollIntoView(child.ref, child.panelId)}
                    tooltip={!contentOutlineExpanded ? child.title : undefined}
                  />
                ))}
            </React.Fragment>
          ))}
        </div>
      </CustomScrollbar>
    </PanelContainer>
  );
}
