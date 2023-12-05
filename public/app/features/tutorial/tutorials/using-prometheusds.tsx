import type { Step, Tutorial } from '../types';

const info = {
  id: 'using-prometheusds',
  name: 'Using the PrometheusDS',
  description: `This is a tutorial to help you get started with the PrometheusDS`,
};

const tutorialSteps: Step[] = [
  {
    route: `/`,
    target: `[data-testid="data-testid Toggle menu"]`,
    title: `Open the menu`,
    content: `Click the menu button to open the menu!`,
    placement: `right`,
    requiredActions: [
      {
        target: `[data-testid="data-testid Toggle menu"]`,
        action: 'click',
      },
    ],
  },
  {
    route: `/`,
    target: `[href="/explore"]`,
    title: `Let's go find the Prometheus datasource!`,
    content: `The Prometheus datasource is a great place to start!`,
    placement: `right`,
    requiredActions: [
      {
        target: `[href="/explore"]`,
        action: 'click',
      },
    ],
  },
  {
    route: `/explore`,
    target: `[data-testid*="Select a data source"]`,
    title: `Let's get started`,
    content: 'Pick the prometheus datasource!',
    requiredActions: [
      {
        target: `[data-testid*="Select a data source"]`,
        action: 'change',
        attribute: {
          name: 'placeholder',
          value: 'gdev-prometheus',
        },
      },
    ],
  },
  {
    route: `/explore`,
    target: `[data-testid*="Select a data source"]`,
    title: `The prometheus datasource is selected`,
    content: `Awesome, let's take a look at what you can do next!`,
  },
  {
    route: `/explore`,
    target: `[aria-label="Query patterns"]`,
    title: `Turbo charge`,
    content: `This is the 'Kick start your query' button. It will help you get started with your first query!`,
  },
  {
    route: `/explore`,
    target: `[aria-label="Toggle switch"]`,
    title: `ELI5`,
    content: `Give it a go!`,
    requiredActions: [
      {
        target: `[aria-label="Toggle switch"]`,
        action: 'click',
      },
    ],
  },
];

export const tutorial: Tutorial = {
  ...info,
  steps: tutorialSteps,
};
