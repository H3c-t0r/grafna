import { useLayoutEffect } from 'react';
import { Branding } from 'app/core/components/Branding/Branding';
import { useGrafana } from 'app/core/context/GrafanaContext';
/**
 * We're definitely doing something odd here, and it all boils down to
 * 1. we have a page layout that is different from what <Page /> forces us to do with pageNav
 * 2. because of 1. we don't get to update the pageNav that way and circumvents
 *     the `usePageTitle` hook in the <Page /> component
 *
 * Therefore we are manually setting the breadcrumb and the page title.
 */
export function useRuleViewerPageTitle(rule) {
    const { chrome } = useGrafana();
    useLayoutEffect(() => {
        if (rule === null || rule === void 0 ? void 0 : rule.name) {
            chrome.update({ pageNav: { text: rule.name } });
        }
    }, [chrome, rule]);
    if (!rule) {
        return;
    }
    document.title = `${rule.name} - Alerting - ${Branding.AppTitle}`;
}
//# sourceMappingURL=useRuleViewerPageTitle.js.map