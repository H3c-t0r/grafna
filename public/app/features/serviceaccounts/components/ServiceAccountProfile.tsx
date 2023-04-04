import { css } from '@emotion/css';
import React from 'react';

import { dateTimeFormat, GrafanaTheme2, OrgRole, TimeZone } from '@grafana/data';
import { useStyles2 } from '@grafana/ui';
import { fetchRoleOptions } from 'app/core/components/RolePicker/api';
import { contextSrv } from 'app/core/core';
import { dispatch } from 'app/store/store';
import { AccessControlAction, Role, ServiceAccountDTO } from 'app/types';

import { acOptionsLoaded } from '../state/reducers';

import { ServiceAccountProfileRow } from './ServiceAccountProfileRow';
import { ServiceAccountRoleRow } from './ServiceAccountRoleRow';

interface Props {
  serviceAccount: ServiceAccountDTO;
  timeZone: TimeZone;
  roleOptions: Role[];
  onChange: (serviceAccount: ServiceAccountDTO) => void;
}

export function ServiceAccountProfile({ serviceAccount, timeZone, roleOptions, onChange }: Props): JSX.Element {
  const styles = useStyles2(getStyles);
  const ableToWrite = contextSrv.hasPermission(AccessControlAction.ServiceAccountsWrite);
  const [roles, setRoles] = React.useState<Role[]>(roleOptions);

  const onRoleChange = (role: OrgRole) => {
    onChange({ ...serviceAccount, role: role });
  };

  const onNameChange = (newValue: string) => {
    onChange({ ...serviceAccount, name: newValue });
  };
  React.useEffect(() => {
    if (roleOptions.length > 0) {
      return;
    }
    if (contextSrv.licensedAccessControlEnabled()) {
      if (contextSrv.hasPermission(AccessControlAction.ActionRolesList)) {
        fetchRoleOptions(serviceAccount.orgId)
          .then((roles) => {
            setRoles(roles);
            dispatch(acOptionsLoaded(roles));
          })
          .catch((err) => {
            console.log('fetchRoleOptions error: ', err);
          });
      }
    }
  }, [roleOptions, serviceAccount.orgId]);

  return (
    <div className={styles.section}>
      <h3>Information</h3>
      <table className="filter-table">
        <tbody>
          <ServiceAccountProfileRow
            label="Name"
            value={serviceAccount.name}
            onChange={onNameChange}
            disabled={!ableToWrite || serviceAccount.isDisabled}
          />
          <ServiceAccountProfileRow label="ID" value={serviceAccount.login} disabled={serviceAccount.isDisabled} />
          <ServiceAccountRoleRow
            label="Roles"
            serviceAccount={serviceAccount}
            onRoleChange={onRoleChange}
            roleOptions={roles}
          />
          <ServiceAccountProfileRow
            label="Creation date"
            value={dateTimeFormat(serviceAccount.createdAt, { timeZone })}
            disabled={serviceAccount.isDisabled}
          />
        </tbody>
      </table>
    </div>
  );
}

export const getStyles = (theme: GrafanaTheme2) => ({
  section: css`
    margin-bottom: ${theme.spacing(4)};
  `,
});
