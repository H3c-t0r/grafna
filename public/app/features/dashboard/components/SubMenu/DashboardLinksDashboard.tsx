import { css, cx } from '@emotion/css';
import React from 'react';
import { useAsync } from 'react-use';

import { GrafanaTheme2 } from '@grafana/data';
import { sanitize, sanitizeUrl } from '@grafana/data/src/text/sanitize';
import { selectors } from '@grafana/e2e-selectors';
import { DashboardLink } from '@grafana/schema';
import { CustomScrollbar, Dropdown, Icon, LinkButton, Menu, useStyles2 } from '@grafana/ui';
import { ButtonLinkProps } from '@grafana/ui/src/components/Button';
import { getBackendSrv } from 'app/core/services/backend_srv';
import { DashboardSearchItem } from 'app/features/search/types';

import { getLinkSrv } from '../../../panel/panellinks/link_srv';

interface Props {
  link: DashboardLink;
  linkInfo: { title: string; href: string };
  dashboardUID: string;
}

interface DashboardLinksMenuProps {
  link: DashboardLink;
  dashboardUID: string;
}

function DashboardLinksMenu({ dashboardUID, link }: DashboardLinksMenuProps) {
  const styles = useStyles2(getStyles);
  const resolvedLinks = useResolvedLinks({ dashboardUID, link });

  if (!resolvedLinks || resolveLinks.length === 0) {
    return null;
  }

  return (
    <Menu>
      <div className={styles.dropdown}>
        <CustomScrollbar>
          {resolvedLinks.map((resolvedLink, index) => {
            return (
              <Menu.Item
                url={resolvedLink.url}
                target={link.targetBlank ? '_blank' : undefined}
                key={`dashlinks-dropdown-item-${resolvedLink.uid}-${index}`}
                label={resolvedLink.title}
                testId={selectors.components.DashboardLinks.link}
                aria-label={`${resolvedLink.title} dashboard`}
              />
            );
          })}
        </CustomScrollbar>
      </div>
    </Menu>
  );
}

export const DashboardLinksDashboard = (props: Props) => {
  const { link, linkInfo, dashboardUID } = props;
  const resolvedLinks = useResolvedLinks(props);
  const styles = useStyles2(getStyles);

  if (link.asDropdown) {
    return (
      <Dropdown overlay={<DashboardLinksMenu link={link} dashboardUID={dashboardUID} />}>
        <DashboardLinkButton
          data-placement="bottom"
          data-toggle="dropdown"
          aria-controls="dropdown-list"
          aria-haspopup="menu"
          fill="outline"
          variant="secondary"
          data-testid={selectors.components.DashboardLinks.dropDown}
        >
          <Icon aria-hidden name="bars" className={styles.iconMargin} />
          <span>{linkInfo.title}</span>
        </DashboardLinkButton>
      </Dropdown>
    );
  }

  return (
    <>
      {resolvedLinks.length > 0 &&
        resolvedLinks.map((resolvedLink, index) => {
          return (
            <DashboardLinkButton
              key={`dashlinks-list-item-${resolvedLink.uid}-${index}`}
              icon="apps"
              variant="secondary"
              fill="outline"
              href={resolvedLink.url}
              target={link.targetBlank ? '_blank' : undefined}
              rel="noreferrer"
              data-testid={selectors.components.DashboardLinks.link}
            >
              {resolvedLink.title}
            </DashboardLinkButton>
          );
        })}
    </>
  );
};

const useResolvedLinks = ({ link, dashboardUID }: Pick<Props, 'link' | 'dashboardUID'>): ResolvedLinkDTO[] => {
  const { tags } = link;
  const result = useAsync(() => searchForTags(tags), [tags]);
  if (!result.value) {
    return [];
  }
  return resolveLinks(dashboardUID, link, result.value);
};

interface ResolvedLinkDTO {
  uid: string;
  url: string;
  title: string;
}

export async function searchForTags(
  tags: string[],
  dependencies: { getBackendSrv: typeof getBackendSrv } = { getBackendSrv }
): Promise<DashboardSearchItem[]> {
  const limit = 100;
  const searchHits: DashboardSearchItem[] = await dependencies.getBackendSrv().search({ tag: tags, limit });

  return searchHits;
}

export function resolveLinks(
  dashboardUID: string,
  link: DashboardLink,
  searchHits: DashboardSearchItem[],
  dependencies: { getLinkSrv: typeof getLinkSrv; sanitize: typeof sanitize; sanitizeUrl: typeof sanitizeUrl } = {
    getLinkSrv,
    sanitize,
    sanitizeUrl,
  }
): ResolvedLinkDTO[] {
  return searchHits
    .filter((searchHit) => searchHit.uid !== dashboardUID)
    .map((searchHit) => {
      const uid = searchHit.uid;
      const title = dependencies.sanitize(searchHit.title);
      const resolvedLink = dependencies.getLinkSrv().getLinkUrl({ ...link, url: searchHit.url });
      const url = dependencies.sanitizeUrl(resolvedLink);

      return { uid, title, url };
    });
}

function getStyles(theme: GrafanaTheme2) {
  return {
    iconMargin: css({
      marginRight: theme.spacing(0.5),
    }),
    dropdown: css({
      maxWidth: 'max(30vw, 300px)',
      maxHeight: '70vh',
      overflowY: 'auto',
    }),
    button: css({
      color: theme.colors.text.primary,
    }),
    dashButton: css({
      fontSize: theme.typography.bodySmall.fontSize,
      paddingLeft: theme.spacing(1),
      paddingRight: theme.spacing(1),
    }),
  };
}

export const DashboardLinkButton = React.forwardRef<HTMLAnchorElement, ButtonLinkProps>(
  ({ className, ...otherProps }, ref) => {
    const styles = useStyles2(getStyles);
    return (
      <LinkButton
        {...otherProps}
        variant="secondary"
        fill="outline"
        className={cx(className, styles.dashButton)}
        ref={ref}
      />
    );
  }
);

DashboardLinkButton.displayName = 'DashboardLinkButton';
