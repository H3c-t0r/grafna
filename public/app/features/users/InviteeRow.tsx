import React, { createRef, PureComponent } from 'react';
import { connect } from 'react-redux';
import { Invitee } from 'app/types';
import { revokeInvite } from './state/actions';
import { Button, ClipboardButton, Icon } from '@grafana/ui';

export interface Props {
  invitee: Invitee;
  revokeInvite: typeof revokeInvite;
}

class InviteeRow extends PureComponent<Props> {
  render() {
    const { invitee, revokeInvite } = this.props;
    return (
      <tr>
        <td>{invitee.email}</td>
        <td>{invitee.name}</td>
        <td className="text-right">
          <ClipboardButton variant="secondary" size="sm" getText={() => invitee.url}>
            Copy Invite
          </ClipboardButton>
          &nbsp;
        </td>
        <td>
          <Button variant="destructive" size="sm" icon="fa fa-remove" onClick={() => revokeInvite(invitee.code)}>
            <Icon name="times" style={{ marginBottom: 0 }} />
          </Button>
        </td>
      </tr>
    );
  }
}

const mapDispatchToProps = {
  revokeInvite,
};

export default connect(() => {
  return {};
}, mapDispatchToProps)(InviteeRow);
