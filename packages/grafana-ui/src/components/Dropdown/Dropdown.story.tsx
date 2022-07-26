import { ComponentMeta } from '@storybook/react';
import React from 'react';

import { StoryExample } from '../../utils/storybook/StoryExample';
import { withCenteredStory } from '../../utils/storybook/withCenteredStory';
import { Button } from '../Button';
import { VerticalGroup } from '../Layout/Layout';
import { Menu } from '../Menu/Menu';
import { MenuItem } from '../Menu/MenuItem';

import { Dropdown } from './Dropdown';
import mdx from './Dropdown.mdx';

const meta: ComponentMeta<typeof Dropdown> = {
  title: 'Overlays/Dropdown',
  component: Dropdown,
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

export function Examples() {
  const menu = (
    <Menu>
      <MenuItem label="View settings" tabIndex={0} />
      <MenuItem label="Edit actions" tabIndex={1} />
      <MenuItem label="Share" tabIndex={2} />
      <MenuItem label="Delete" tabIndex={3} />
    </Menu>
  );

  return (
    <VerticalGroup>
      <StoryExample name="Button + defaults">
        <Dropdown overlay={menu}>
          <Button variant="secondary">Button</Button>
        </Dropdown>
      </StoryExample>
      <StoryExample name="Icon button, placement=bottom-start">
        <Dropdown overlay={menu} placement="bottom-start">
          <Button variant="secondary" icon="bars" />
        </Dropdown>
      </StoryExample>
    </VerticalGroup>
  );
}

Examples.parameters = {
  controls: {
    hideNoControlsWarning: true,
    include: [],
  },
};

export default meta;
