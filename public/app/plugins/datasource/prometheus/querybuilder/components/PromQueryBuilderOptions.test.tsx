import { render, screen } from '@testing-library/react';
import React from 'react';

import { CoreApp } from '@grafana/data';
import { selectOptionInTest } from '@grafana/ui';

import { PromQuery } from '../../types';
import { getQueryWithDefaults } from '../state';

import { PromQueryBuilderOptions } from './PromQueryBuilderOptions';

describe('PromQueryBuilderOptions', () => {
  it('Can change query type', async () => {
    const { props } = setup();

    screen.getByTitle('Click to edit options').click();
    expect(screen.getByLabelText('Range')).toBeChecked();

    screen.getByLabelText('Instant').click();

    expect(props.onChange).toHaveBeenCalledWith({
      ...props.query,
      instant: true,
      range: false,
      exemplar: false,
    });
  });

  it('Legend format default to Auto', async () => {
    setup();
    expect(screen.getByText('Legend: Auto')).toBeInTheDocument();
  });

  it('Can change legend format to verbose', async () => {
    const { props } = setup();

    screen.getByTitle('Click to edit options').click();

    let legendModeSelect = screen.getByText('Auto').parentElement!;
    legendModeSelect.click();

    await selectOptionInTest(legendModeSelect as HTMLElement, 'Verbose');

    expect(props.onChange).toHaveBeenCalledWith({
      ...props.query,
      legendFormat: '',
    });
  });

  it('Can change legend format to custom', async () => {
    const { props } = setup();

    screen.getByTitle('Click to edit options').click();

    let legendModeSelect = screen.getByText('Auto').parentElement!;
    legendModeSelect.click();

    await selectOptionInTest(legendModeSelect as HTMLElement, 'Custom');

    expect(props.onChange).toHaveBeenCalledWith({
      ...props.query,
      legendFormat: '{{label_name}}',
    });
  });

  it('Handle defaults with undefined range', async () => {
    setup(getQueryWithDefaults({ refId: 'A', expr: '', range: undefined, instant: true }, CoreApp.Dashboard, '9.0.0'));

    expect(screen.getByText('Type: Instant')).toBeInTheDocument();
  });
});

function setup(queryOverrides: Partial<PromQuery> = {}) {
  const props = {
    query: {
      ...getQueryWithDefaults({ refId: 'A' } as PromQuery, CoreApp.PanelEditor, '9.0.0'),
      ...queryOverrides,
    },
    onRunQuery: jest.fn(),
    onChange: jest.fn(),
  };

  const { container } = render(<PromQueryBuilderOptions {...props} />);
  return { container, props };
}
