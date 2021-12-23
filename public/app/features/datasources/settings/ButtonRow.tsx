import React, { FC } from 'react';
import { selectors } from '@grafana/e2e-selectors';

import { Button, LinkButton } from '@grafana/ui';

import { AccessControlAction } from 'app/types/';
import { contextSrv } from 'app/core/core';

export interface Props {
  exploreUrl: string;
  canSave: boolean;
  canDelete: boolean;
  onDelete: () => void;
  onSubmit: (event: any) => void;
  onTest: (event: any) => void;
}

const ButtonRow: FC<Props> = ({ canSave, canDelete, onDelete, onSubmit, onTest, exploreUrl }) => {
  const canExploreDataSources = contextSrv.hasPermission(AccessControlAction.DataSourcesExplore);

  return (
    <div className="gf-form-button-row">
      <Button variant="secondary" fill="solid" type="button" onClick={() => history.back()}>
        Back
      </Button>
      <LinkButton variant="secondary" fill="solid" href={exploreUrl} disabled={!canExploreDataSources}>
        Explore
      </LinkButton>
      <Button
        type="button"
        variant="destructive"
        disabled={!canDelete}
        onClick={onDelete}
        aria-label={selectors.pages.DataSource.delete}
      >
        Delete
      </Button>
      {canSave && (
        <Button
          type="submit"
          variant="primary"
          disabled={!canSave}
          onClick={(event) => onSubmit(event)}
          aria-label={selectors.pages.DataSource.saveAndTest}
        >
          Save &amp; test
        </Button>
      )}
      {!canSave && (
        <Button type="submit" variant="primary" onClick={onTest}>
          Test
        </Button>
      )}
    </div>
  );
};

export default ButtonRow;
