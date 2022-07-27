import { css } from '@emotion/css';
import pluralize from 'pluralize';
import React, { FC, useMemo, useState } from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { Button, ConfirmModal, Modal, useStyles2, Badge } from '@grafana/ui';
import { contextSrv } from 'app/core/services/context_srv';
import { AlertManagerCortexConfig } from 'app/plugins/datasource/alertmanager/types';
import { useDispatch, ContactPointsState } from 'app/types';

import { Authorize } from '../../components/Authorize';
import { useUnifiedAlertingSelector } from '../../hooks/useUnifiedAlertingSelector';
import { deleteReceiverAction } from '../../state/actions';
import { getAlertTableStyles } from '../../styles/table';
import { getNotificationsPermissions } from '../../utils/access-control';
import { isReceiverUsed } from '../../utils/alertmanager';
import { isVanillaPrometheusAlertManagerDataSource } from '../../utils/datasource';
import { makeAMLink } from '../../utils/misc';
import { extractNotifierTypeCounts } from '../../utils/receivers';
import { initialAsyncRequestState } from '../../utils/redux';
import { ProvisioningBadge } from '../Provisioning';
import { ActionIcon } from '../rules/ActionIcon';

import { ReceiversSection } from './ReceiversSection';

interface ReceiverErrorProps {
  errorCount: number;
}
function ReceiverError({ errorCount }: ReceiverErrorProps) {
  return (
    <Badge
      color="orange"
      icon="exclamation-triangle"
      text={`${errorCount} ${pluralize('error', errorCount)}`}
      tooltip={`${errorCount} ${pluralize('error', errorCount)} detected in this contact point`}
    />
  );
}
interface ReceiverHealthProps {
  errorsByReceiver: number;
}

function ReceiverHealth({ errorsByReceiver }: ReceiverHealthProps) {
  return errorsByReceiver > 0 ? (
    <ReceiverError errorCount={errorsByReceiver} />
  ) : (
    <Badge color="green" text="OK" tooltip="No errors detected" />
  );
}

interface Props {
  config: AlertManagerCortexConfig;
  alertManagerName: string;
}

const useContactPointsState = (alertManagerName: string) => {
  const contactPointsStateRequest = useUnifiedAlertingSelector((state) => state.contactPointsState);
  const { result: contactPointsState } = (alertManagerName && contactPointsStateRequest) || initialAsyncRequestState;
  const receivers = contactPointsState?.receivers ?? {};
  const errorStateAvailable = Object.keys(receivers).length > 0; // this logic can change depending on how we implement this in the BE
  return { contactPointsState, errorStateAvailable };
};

export const ReceiversTable: FC<Props> = ({ config, alertManagerName }) => {
  const dispatch = useDispatch();
  const tableStyles = useStyles2(getAlertTableStyles);
  const styles = useStyles2(getStyles);
  const isVanillaAM = isVanillaPrometheusAlertManagerDataSource(alertManagerName);
  const permissions = getNotificationsPermissions(alertManagerName);
  const grafanaNotifiers = useUnifiedAlertingSelector((state) => state.grafanaNotifiers);

  const { contactPointsState, errorStateAvailable } = useContactPointsState(alertManagerName);

  // receiver name slated for deletion. If this is set, a confirmation modal is shown. If user approves, this receiver is deleted
  const [receiverToDelete, setReceiverToDelete] = useState<string>();
  const [showCannotDeleteReceiverModal, setShowCannotDeleteReceiverModal] = useState(false);

  const onClickDeleteReceiver = (receiverName: string): void => {
    if (isReceiverUsed(receiverName, config)) {
      setShowCannotDeleteReceiverModal(true);
    } else {
      setReceiverToDelete(receiverName);
    }
  };

  const deleteReceiver = () => {
    if (receiverToDelete) {
      dispatch(deleteReceiverAction(receiverToDelete, alertManagerName));
    }
    setReceiverToDelete(undefined);
  };

  const rows = useMemo(
    () =>
      config.alertmanager_config.receivers?.map((receiver) => ({
        name: receiver.name,
        types: Object.entries(extractNotifierTypeCounts(receiver, grafanaNotifiers.result ?? [])).map(
          ([type, count]) => {
            if (count > 1) {
              return `${type} (${count})`;
            }
            return type;
          }
        ),
        provisioned: receiver.grafana_managed_receiver_configs?.some((receiver) => receiver.provenance),
      })) ?? [],
    [config, grafanaNotifiers.result]
  );

  const errorsByReceiver = (contactPointsState: ContactPointsState, receiverName: string) =>
    contactPointsState.receivers[receiverName]?.errorCount ?? 0;

  return (
    <ReceiversSection
      className={styles.section}
      title="Contact points"
      description="Define where the notifications will be sent to, for example email or Slack."
      showButton={!isVanillaAM && contextSrv.hasPermission(permissions.create)}
      addButtonLabel="New contact point"
      addButtonTo={makeAMLink('/alerting/notifications/receivers/new', alertManagerName)}
    >
      <table className={tableStyles.table} data-testid="receivers-table">
        <colgroup>
          <col />
          <col />
          <Authorize actions={[permissions.update, permissions.delete]}>
            <col />
          </Authorize>
        </colgroup>
        <thead>
          <tr>
            <th>Contact point name</th>
            <th>Type</th>
            {errorStateAvailable && <th>Health</th>}
            <Authorize actions={[permissions.update, permissions.delete]}>
              <th>Actions</th>
            </Authorize>
          </tr>
        </thead>
        <tbody>
          {!rows.length && (
            <tr className={tableStyles.evenRow}>
              <td colSpan={3}>No receivers defined.</td>
            </tr>
          )}
          {rows.map((receiver, idx) => (
            <tr key={receiver.name} className={idx % 2 === 0 ? tableStyles.evenRow : undefined}>
              <td>
                {receiver.name} {receiver.provisioned && <ProvisioningBadge />}
              </td>
              <td>{receiver.types.join(', ')}</td>
              {errorStateAvailable && (
                <td>
                  <ReceiverHealth
                    errorsByReceiver={contactPointsState ? errorsByReceiver(contactPointsState, receiver.name) : 0}
                  />
                </td>
              )}
              <Authorize actions={[permissions.update, permissions.delete]}>
                <td className={tableStyles.actionsCell}>
                  {!isVanillaAM && !receiver.provisioned && (
                    <>
                      <Authorize actions={[permissions.update]}>
                        <ActionIcon
                          aria-label="Edit"
                          data-testid="edit"
                          to={makeAMLink(
                            `/alerting/notifications/receivers/${encodeURIComponent(receiver.name)}/edit`,
                            alertManagerName
                          )}
                          tooltip="Edit contact point"
                          icon="pen"
                        />
                      </Authorize>
                      <Authorize actions={[permissions.delete]}>
                        <ActionIcon
                          onClick={() => onClickDeleteReceiver(receiver.name)}
                          tooltip="Delete contact point"
                          icon="trash-alt"
                        />
                      </Authorize>
                    </>
                  )}
                  {(isVanillaAM || receiver.provisioned) && (
                    <Authorize actions={[permissions.update]}>
                      <ActionIcon
                        data-testid="view"
                        to={makeAMLink(
                          `/alerting/notifications/receivers/${encodeURIComponent(receiver.name)}/edit`,
                          alertManagerName
                        )}
                        tooltip="View contact point"
                        icon="file-alt"
                      />
                    </Authorize>
                  )}
                </td>
              </Authorize>
            </tr>
          ))}
        </tbody>
      </table>
      {!!showCannotDeleteReceiverModal && (
        <Modal
          isOpen={true}
          title="Cannot delete contact point"
          onDismiss={() => setShowCannotDeleteReceiverModal(false)}
        >
          <p>
            Contact point cannot be deleted because it is used in more policies. Please update or delete these policies
            first.
          </p>
          <Modal.ButtonRow>
            <Button variant="secondary" onClick={() => setShowCannotDeleteReceiverModal(false)} fill="outline">
              Close
            </Button>
          </Modal.ButtonRow>
        </Modal>
      )}
      {!!receiverToDelete && (
        <ConfirmModal
          isOpen={true}
          title="Delete contact point"
          body={`Are you sure you want to delete contact point "${receiverToDelete}"?`}
          confirmText="Yes, delete"
          onConfirm={deleteReceiver}
          onDismiss={() => setReceiverToDelete(undefined)}
        />
      )}
    </ReceiversSection>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  section: css`
    margin-top: ${theme.spacing(4)};
  `,
  warning: css`
    color: ${theme.colors.warning.text};
  `,
  countMessage: css``,
});
