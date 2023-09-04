import React, { useCallback, useState } from 'react';

import { Modal, ToolbarButton } from '@grafana/ui';
import { t, Trans } from 'app/core/internationalization';
import { useSelector } from 'app/types';

import { getExploreItemSelector } from '../../state/selectors';

import { AddToDashboardForm } from './AddToDashboardForm';
import { getAddToDashboardTitle } from './getAddToDashboardTitle';

interface Props {
  exploreId: string;
}

export const AddToDashboard = ({ exploreId }: Props) => {
  const [isOpen, setIsOpen] = useState(false);
  const selectExploreItem = getExploreItemSelector(exploreId);
  const explorePaneHasQueries = !!useSelector(selectExploreItem)?.queries?.length;
  const onClose = useCallback(() => setIsOpen(false), []);

  return (
    <>
      <ToolbarButton
        icon="apps"
        variant="canvas"
        onClick={() => setIsOpen(true)}
        aria-label={t('explore.add-to-dashboard', 'Add to dashboard')}
        disabled={!explorePaneHasQueries}
      >
        <Trans i18nKey="explore.add-to-dashboard">Add to dashboard</Trans>
      </ToolbarButton>

      {isOpen && (
        <Modal title={getAddToDashboardTitle()} onDismiss={onClose} isOpen>
          <AddToDashboardForm onClose={onClose} exploreId={exploreId} />
        </Modal>
      )}
    </>
  );
};
