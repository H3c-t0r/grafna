import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';
import FakeSchemaData from '../../azure_log_analytics/__mocks__/schema';

import createMockDatasource from '../../__mocks__/datasource';
import createMockQuery from '../../__mocks__/query';
import { TimeManagement } from './TimeManagement';

const variableOptionGroup = {
  label: 'Template variables',
  options: [],
};

describe('LogsQueryEditor.TimeManagement', () => {
  it('should render the column picker if Dashboard is chosen', async () => {
    const mockDatasource = createMockDatasource();
    const query = createMockQuery();
    const onChange = jest.fn();

    const { rerender } = render(
      <TimeManagement
        query={query}
        datasource={mockDatasource}
        variableOptionGroup={variableOptionGroup}
        onQueryChange={onChange}
        setError={() => {}}
        schema={FakeSchemaData.getLogAnalyticsFakeEngineSchema()}
      />
    );

    const dashboardTimeOption = await screen.findByLabelText('Dashboard');
    await userEvent.click(dashboardTimeOption);

    expect(onChange).toBeCalledWith(
      expect.objectContaining({
        azureLogAnalytics: expect.objectContaining({
          dashboardTime: true,
        }),
      })
    );

    rerender(
      <TimeManagement
        query={{ ...query, azureLogAnalytics: { ...query.azureLogAnalytics, dashboardTime: true } }}
        datasource={mockDatasource}
        variableOptionGroup={variableOptionGroup}
        onQueryChange={onChange}
        setError={() => {}}
        schema={FakeSchemaData.getLogAnalyticsFakeEngineSchema()}
      />
    );

    expect(onChange).toBeCalledWith(
      expect.objectContaining({
        azureLogAnalytics: expect.objectContaining({
          timeColumn: 'TimeGenerated',
        }),
      })
    );
  });

  it('should render the default value if no time columns exist', async () => {
    const mockDatasource = createMockDatasource();
    const query = createMockQuery();
    const onChange = jest.fn();

    render(
      <TimeManagement
        query={{ ...query, azureLogAnalytics: { ...query.azureLogAnalytics, dashboardTime: true } }}
        datasource={mockDatasource}
        variableOptionGroup={variableOptionGroup}
        onQueryChange={onChange}
        setError={() => {}}
        schema={FakeSchemaData.getLogAnalyticsFakeEngineSchema([
          {
            id: 't/Alert',
            name: 'Alert',
            timespanColumn: 'TimeGenerated',
            columns: [],
            related: {
              solutions: [],
            },
          },
        ])}
      />
    );

    expect(onChange).toBeCalledWith(
      expect.objectContaining({
        azureLogAnalytics: expect.objectContaining({
          timeColumn: 'TimeGenerated',
        }),
      })
    );
  });

  it('should render the first time column if no default exists', async () => {
    const mockDatasource = createMockDatasource();
    const query = createMockQuery();
    const onChange = jest.fn();

    render(
      <TimeManagement
        query={{ ...query, azureLogAnalytics: { ...query.azureLogAnalytics, dashboardTime: true } }}
        datasource={mockDatasource}
        variableOptionGroup={variableOptionGroup}
        onQueryChange={onChange}
        setError={() => {}}
        schema={FakeSchemaData.getLogAnalyticsFakeEngineSchema([
          {
            id: 't/Alert',
            name: 'Alert',
            timespanColumn: '',
            columns: [{ name: 'Timespan', type: 'datetime' }],
            related: {
              solutions: [],
            },
          },
        ])}
      />
    );

    expect(onChange).toBeCalledWith(
      expect.objectContaining({
        azureLogAnalytics: expect.objectContaining({
          timeColumn: 'Timespan',
        }),
      })
    );
  });
});
