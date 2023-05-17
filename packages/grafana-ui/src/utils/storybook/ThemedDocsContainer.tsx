// Wrap the DocsContainer for storybook-dark-mode theme switching support.
import { DocsContainer, DocsContextProps } from '@storybook/addon-docs';
import React from 'react';
import { useDarkMode } from 'storybook-dark-mode';

import { GrafanaLight, GrafanaDark } from '../../../.storybook/storybookTheme';

type Props = {
  context: DocsContextProps;
  children?: React.ReactNode;
};

export const ThemedDocsContainer = ({ children, context }: Props) => {
  const dark = useDarkMode();

  return (
    <DocsContainer theme={dark ? GrafanaDark : GrafanaLight} context={context}>
      {children}
    </DocsContainer>
  );
};
