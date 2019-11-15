import { ClickablePageObjectType, TestPage } from '@grafana/toolkit/src/e2e';

export interface VariablesPage {
  callToActionButton: ClickablePageObjectType;
}

export const variablesPage = new TestPage<VariablesPage>({
  pageObjects: {
    callToActionButton: 'Call to action button Add variable',
  },
});
