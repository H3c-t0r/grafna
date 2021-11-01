import { e2e } from '@grafana/e2e';
import { addDays, addHours, differenceInCalendarDays, differenceInMinutes, format, isBefore, parse } from 'date-fns';

e2e.scenario({
  describeName: 'Dashboard time zone support',
  itName: 'Tests dashboard time zone scenarios',
  addScenarioDataSource: false,
  addScenarioDashBoard: false,
  skipScenario: false,
  scenario: () => {
    e2e.flows.openDashboard({ uid: '5SdHCasdf' });

    const fromTimeZone = 'Coordinated Universal Time';
    const toTimeZone = 'America/Chicago';
    const offset = -5;

    const panelsToCheck = [
      'Random walk series',
      'Millisecond res x-axis and tooltip',
      '2 yaxis and axis labels',
      'Stacking value ontop of nulls',
      'Null between points',
      'Legend Table No Scroll Visible',
    ];

    const timesInUtc: Record<string, string> = {};

    for (const title of panelsToCheck) {
      e2e.components.Panels.Panel.containerByTitle(title)
        .should('be.visible')
        .within(() =>
          e2e.components.Panels.Visualization.Graph.xAxis
            .labels()
            .should('be.visible')
            .last()
            .should((element) => {
              timesInUtc[title] = element.text();
            })
        );
    }

    e2e.components.PageToolbar.item('Dashboard settings').click();

    e2e.components.TimeZonePicker.container()
      .should('be.visible')
      .within(() => {
        e2e.components.Select.singleValue().should('be.visible').should('have.text', fromTimeZone);

        e2e.components.Select.input().should('be.visible').click();
      });

    e2e.components.Select.option().should('be.visible').contains(toTimeZone).click();

    // click to go back to the dashboard.
    e2e.components.BackButton.backArrow().click({ force: true }).wait(2000);

    for (const title of panelsToCheck) {
      e2e.components.Panels.Panel.containerByTitle(title)
        .should('be.visible')
        .within(() =>
          e2e.components.Panels.Visualization.Graph.xAxis
            .labels()
            .should('be.visible')
            .last()
            .should((element) => {
              const inUtc = timesInUtc[title];
              const inTz = element.text();
              const isCorrect = isTimeCorrect(inUtc, inTz, offset);
              assert.isTrue(isCorrect, `Panel with title: "${title}"`);
            })
        );
    }
  },
});

const isTimeCorrect = (inUtc: string, inTz: string, offset: number): boolean => {
  if (inUtc === inTz) {
    // we need to catch issues when timezone isn't changed for some reason like https://github.com/grafana/grafana/issues/35504
    return false;
  }

  const reference = format(new Date(), 'yyyy-LL-dd');

  const utcDate = parse(`${reference} ${inUtc}`);
  const utcDateWithOffset = addHours(parse(`${reference} ${inUtc}`), offset);
  const dayDifference = differenceInCalendarDays(utcDate, utcDateWithOffset); // if the utcDate +/- offset is the day before/after then we need to adjust reference
  const dayOffset = isBefore(utcDateWithOffset, utcDate) ? dayDifference * -1 : dayDifference;
  const tzDate = addDays(parse(`${reference} ${inTz}`), dayOffset); // adjust tzDate with any dayOffset
  const diff = Math.abs(differenceInMinutes(utcDate, tzDate)); // use Math.abs if tzDate is in future

  return diff <= Math.abs(offset * 60);
};
