import { action } from '@storybook/addon-actions';
import React from 'react';

import { RefreshPicker } from '@grafana/ui';

import { DashboardStoryCanvas } from '../../utils/storybook/DashboardStoryCanvas';
import { StoryExample } from '../../utils/storybook/StoryExample';
import { UseState } from '../../utils/storybook/UseState';
import { withCenteredStory } from '../../utils/storybook/withCenteredStory';
import { HorizontalGroup } from '../Layout/Layout';

import mdx from './RefreshPicker.mdx';

export default {
  title: 'Pickers and Editors/RefreshPicker',
  component: RefreshPicker,
  decorators: [withCenteredStory],
  parameters: {
    docs: {
      page: mdx,
    },
  },
};

export const Examples = () => {
  const intervals = ['5s', '10s', '30s', '1m', '5m', '15m', '30m', '1h', '2h', '1d'];
  const onIntervalChanged = (interval: string) => {
    action('onIntervalChanged fired')(interval);
  };

  const onRefresh = () => {
    action('onRefresh fired')();
  };

  return (
    <DashboardStoryCanvas>
      <UseState initialState={'1h'}>
        {(value, updateValue) => {
          return (
            <HorizontalGroup>
              <StoryExample name="Simple">
                <RefreshPicker
                  tooltip="Hello world"
                  value={value}
                  intervals={intervals}
                  onIntervalChanged={onIntervalChanged}
                  onRefresh={onRefresh}
                />
              </StoryExample>
              <StoryExample name="With text">
                <RefreshPicker
                  tooltip="Hello world"
                  value={value}
                  text="Run query"
                  intervals={intervals}
                  onIntervalChanged={onIntervalChanged}
                  onRefresh={onRefresh}
                />
              </StoryExample>
              <StoryExample name="With text and loading">
                <RefreshPicker
                  tooltip="Hello world"
                  value={value}
                  text="Run query"
                  isLoading={true}
                  intervals={intervals}
                  onIntervalChanged={onIntervalChanged}
                  onRefresh={onRefresh}
                />
              </StoryExample>
            </HorizontalGroup>
          );
        }}
      </UseState>
    </DashboardStoryCanvas>
  );
};
