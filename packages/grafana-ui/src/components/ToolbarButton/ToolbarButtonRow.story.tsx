import { ComponentMeta, ComponentStory } from '@storybook/react';
import React from 'react';

import { withCenteredStory } from '../../utils/storybook/withCenteredStory';

import { ToolbarButton } from './ToolbarButton';
import mdx from './ToolbarButton.mdx';
import { ToolbarButtonRow } from './ToolbarButtonRow';

const meta: ComponentMeta<typeof ToolbarButtonRow> = {
  title: 'Buttons/ToolbarButton/ToolbarButtonRow',
  component: ToolbarButtonRow,
  decorators: [withCenteredStory],
  parameters: {
    docs: {
      page: mdx,
    },
    controls: {
      exclude: ['className'],
    },
  },
};

export const Basic: ComponentStory<typeof ToolbarButtonRow> = (args) => {
  return (
    <ToolbarButtonRow {...args}>
      <ToolbarButton>Just text</ToolbarButton>
      <ToolbarButton icon="sync" tooltip="Sync" />
      <ToolbarButton imgSrc="./grafana_icon.svg">With imgSrc</ToolbarButton>
      <ToolbarButton>Just text</ToolbarButton>
      <ToolbarButton icon="sync" tooltip="Sync" />
      <ToolbarButton imgSrc="./grafana_icon.svg">With imgSrc</ToolbarButton>
    </ToolbarButtonRow>
  );
};

export default meta;
