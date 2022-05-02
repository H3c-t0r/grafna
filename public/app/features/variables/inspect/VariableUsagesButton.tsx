import React, { FC, useMemo } from 'react';

import { reportInteraction } from '@grafana/runtime';
import { IconButton } from '@grafana/ui';

import { NetworkGraphModal } from './NetworkGraphModal';
import { UsagesToNetwork } from './utils';

interface Props {
  id: string;
  usages: UsagesToNetwork[];
  isAdhoc: boolean;
}

export const VariableUsagesButton: FC<Props> = ({ id, usages, isAdhoc }) => {
  const network = useMemo(() => usages.find((n) => n.variable.id === id), [usages, id]);
  if (usages.length === 0 || isAdhoc || !network) {
    return null;
  }

  const nodes = network.nodes.map((n) => {
    if (n.label.includes(`$${id}`)) {
      return { ...n, color: '#FB7E81' };
    }
    return n;
  });

  return (
    <NetworkGraphModal show={false} title={`Showing usages for: $${id}`} nodes={nodes} edges={network.edges}>
      {({ showModal }) => {
        return (
          <IconButton
            onClick={() => {
              reportInteraction('Show variable usages');
              showModal();
            }}
            name="code-branch"
            title="Show usages"
          />
        );
      }}
    </NetworkGraphModal>
  );
};
