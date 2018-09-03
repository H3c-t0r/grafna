import React from 'react';
import _ from 'lodash';
import { hot } from 'react-hot-loader';
import { inject, observer } from 'mobx-react';
import config from 'app/core/config';
import PageHeader from 'app/core/components/PageHeader/PageHeader';
import { NavStore } from 'app/stores/NavStore/NavStore';
import { TeamsStore, Team } from 'app/stores/TeamsStore/TeamsStore';
import { ViewStore } from 'app/stores/ViewStore/ViewStore';
import { BackendSrv } from 'app/core/services/backend_srv';
import TeamMembers from './TeamMembers';
import TeamSettings from './TeamSettings';
import TeamGroupSync from './TeamGroupSync';

interface Props {
  nav: typeof NavStore.Type;
  teams: typeof TeamsStore.Type;
  view: typeof ViewStore.Type;
  backendSrv: BackendSrv;
}

@inject('nav', 'teams', 'view')
@observer
export class TeamPages extends React.Component<Props, any> {
  isSyncEnabled: boolean;
  currentPage: string;
  isGrafanaAdmin: boolean;

  constructor(props) {
    super(props);

    this.isSyncEnabled = config.buildInfo.isEnterprise;
    this.isGrafanaAdmin = config.bootData.user.isGrafanaAdmin;
    this.currentPage = this.getCurrentPage();

    this.loadTeam();
  }

  async loadTeam() {
    const { teams, nav, view } = this.props;

    await teams.loadById(view.routeParams.get('id'));

    nav.initTeamPage(this.getCurrentTeam(), this.currentPage, this.isSyncEnabled);
  }

  getCurrentTeam(): Team {
    const { teams, view } = this.props;
    return teams.map.get(view.routeParams.get('id'));
  }

  getCurrentPage() {
    const pages = ['members', 'settings', 'groupsync'];
    const currentPage = this.props.view.routeParams.get('page');
    return _.includes(pages, currentPage) ? currentPage : pages[0];
  }

  render() {
    const { nav, backendSrv } = this.props;
    const currentTeam = this.getCurrentTeam();

    if (!nav.main) {
      return null;
    }

    return (
      <div>
        <PageHeader model={nav as any} />
        {currentTeam && (
          <div className="page-container page-body">
            {this.currentPage === 'members' && <TeamMembers team={currentTeam} backendSrv={backendSrv} />}
            {this.currentPage === 'settings' && this.isGrafanaAdmin && <TeamSettings team={currentTeam} />}
            {this.currentPage === 'groupsync' && this.isSyncEnabled && <TeamGroupSync team={currentTeam} />}
          </div>
        )}
      </div>
    );
  }
}

export default hot(module)(TeamPages);
