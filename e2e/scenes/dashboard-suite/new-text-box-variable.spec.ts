import { e2e } from '../utils';

const PAGE_UNDER_TEST = 'kVi2Gex7z/test-variable-output';
const DASHBOARD_NAME = 'Test variable output';

describe('Variables - Text box', () => {
  beforeEach(() => {
    e2e.flows.login(Cypress.env('USERNAME'), Cypress.env('PASSWORD'));
  });

  it('can add a new text box variable', () => {
    e2e.flows.openDashboard({ uid: `${PAGE_UNDER_TEST}?orgId=1&editview=variables` });
    cy.contains(DASHBOARD_NAME).should('be.visible');

    // Create a new "text box" variable
    e2e.components.CallToActionCard.buttonV2('Add variable').click();

    e2e.pages.Dashboard.Settings.Variables.Edit.General.generalTypeSelectV2().within(() => {
      cy.get('input').type('Textbox{enter}');
    });
    e2e.pages.Dashboard.Settings.Variables.Edit.General.generalNameInputV2().clear().type('VariableUnderTest').blur();
    e2e.pages.Dashboard.Settings.Variables.Edit.General.generalLabelInputV2().type('Variable under test').blur();
    e2e.pages.Dashboard.Settings.Variables.Edit.TextBoxVariable.textBoxOptionsQueryInputV2().type('cat-dog').blur();

    // Navigate back to the homepage and change the selected variable value
    e2e.pages.Dashboard.Settings.Variables.Edit.General.applyButton().click();
    e2e.components.NavToolbar.editDashboard.backToDashboardButton().click();
    e2e.pages.Dashboard.SubMenu.submenuItem().within(() => {
      cy.get('input').clear().type('dog-cat').blur();
    });
    // Assert it was rendered
    cy.get('.markdown-html').should('include.text', 'VariableUnderTest: dog-cat');
  });
});
