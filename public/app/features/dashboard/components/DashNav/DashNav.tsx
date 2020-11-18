// Libaries
import React, { PureComponent, FC, ReactNode } from 'react';
import { connect, MapDispatchToProps } from 'react-redux';
// Utils & Services
import { appEvents } from 'app/core/app_events';
import { PlaylistSrv } from 'app/features/playlist/playlist_srv';
// Components
import { DashNavButton } from './DashNavButton';
import { DashNavTimeControls } from './DashNavTimeControls';
import { ModalsController } from '@grafana/ui';
import { textUtil } from '@grafana/data';
import { BackButton } from 'app/core/components/BackButton/BackButton';
// State
import { updateLocation } from 'app/core/actions';
import { updateTimeZoneForSession } from 'app/features/profile/state/reducers';
// Types
import { DashboardModel } from '../../state';
import { CoreEvents, StoreState } from 'app/types';
import { SaveDashboardModalProxy } from 'app/features/dashboard/components/SaveDashboard/SaveDashboardModalProxy';

export interface OwnProps {
  dashboard: DashboardModel;
  isFullscreen: boolean;
  $injector: any;
  onAddPanel: () => void;
}

interface DispatchProps {
  updateTimeZoneForSession: typeof updateTimeZoneForSession;
  updateLocation: typeof updateLocation;
}

interface DashNavButtonModel {
  show: (props: Props) => boolean;
  component: FC<Partial<Props>>;
  index?: number | 'end';
}

const customLeftActions: DashNavButtonModel[] = [];
const customRightActions: DashNavButtonModel[] = [];

export function addCustomLeftAction(content: DashNavButtonModel) {
  customLeftActions.push(content);
}

export function addCustomRightAction(content: DashNavButtonModel) {
  customRightActions.push(content);
}

export interface StateProps {
  location: any;
}

type Props = StateProps & OwnProps & DispatchProps;

class DashNav extends PureComponent<Props> {
  playlistSrv: PlaylistSrv;

  constructor(props: Props) {
    super(props);
    this.playlistSrv = this.props.$injector.get('playlistSrv');
  }

  onFolderNameClick = () => {
    this.props.updateLocation({
      query: { search: 'open', folder: 'current' },
      partial: true,
    });
  };

  onClose = () => {
    this.props.updateLocation({
      query: { viewPanel: null },
      partial: true,
    });
  };

  onToggleTVMode = () => {
    appEvents.emit(CoreEvents.toggleKioskMode);
  };

  onOpenSettings = () => {
    this.props.updateLocation({
      query: { editview: 'settings' },
      partial: true,
    });
  };

  onStarDashboard = () => {
    const { dashboard, $injector } = this.props;
    const dashboardSrv = $injector.get('dashboardSrv');

    dashboardSrv.starDashboard(dashboard.id, dashboard.meta.isStarred).then((newState: any) => {
      dashboard.meta.isStarred = newState;
      this.forceUpdate();
    });
  };

  onPlaylistPrev = () => {
    this.playlistSrv.prev();
  };

  onPlaylistNext = () => {
    this.playlistSrv.next();
  };

  onPlaylistStop = () => {
    this.playlistSrv.stop();
    this.forceUpdate();
  };

  addCustomContent(actions: DashNavButtonModel[], buttons: ReactNode[]) {
    actions.map((action, index) => {
      const Component = action.component;
      const element = <Component {...this.props} key={`button-custom-${index}`} />;
      typeof action.index === 'number' ? buttons.splice(action.index, 0, element) : buttons.push(element);
    });
  }

  renderDashboardTitleSearchButton() {
    return (
      <>
        <div className="navbar__spacer" />
      </>
    );
  }

  renderBackButton() {
    return (
      <div className="navbar-edit">
        <BackButton surface="dashboard" onClick={this.onClose} />
      </div>
    );
  }

  renderRightActionsButton() {
    const { dashboard, onAddPanel } = this.props;
    const { canEdit, showSettings } = dashboard.meta;
    const { snapshot } = dashboard;
    const snapshotUrl = snapshot && snapshot.originalUrl;

    const buttons: ReactNode[] = [];
    if (canEdit) {
      buttons.push(
        <DashNavButton
          classSuffix="save"
          tooltip="Add panel"
          icon="panel-add"
          onClick={onAddPanel}
          iconType="mono"
          iconSize="xl"
          key="button-panel-add"
        />
      );
      buttons.push(
        <ModalsController key="button-save">
          {({ showModal, hideModal }) => (
            <DashNavButton
              tooltip="Save dashboard"
              classSuffix="save"
              icon="save"
              onClick={() => {
                showModal(SaveDashboardModalProxy, {
                  dashboard,
                  onDismiss: hideModal,
                });
              }}
            />
          )}
        </ModalsController>
      );
    }

    if (snapshotUrl) {
      buttons.push(
        <DashNavButton
          tooltip="Open original dashboard"
          classSuffix="snapshot-origin"
          href={textUtil.sanitizeUrl(snapshotUrl)}
          icon="link"
          key="button-snapshot"
        />
      );
    }

    if (showSettings) {
      buttons.push(
        <DashNavButton
          tooltip="Dashboard settings"
          classSuffix="settings"
          icon="cog"
          onClick={this.onOpenSettings}
          key="button-settings"
        />
      );
    }

    this.addCustomContent(customRightActions, buttons);
    return buttons;
  }

  render() {
    const { dashboard, location, isFullscreen, updateTimeZoneForSession } = this.props;

    return (
      <div className="navbar" style={location.query.hideTimeControls ? { display: 'none' } : {}}>
        {isFullscreen && this.renderBackButton()}
        {this.renderDashboardTitleSearchButton()}

        {this.playlistSrv.isPlaying && (
          <div className="navbar-buttons navbar-buttons--playlist">
            <DashNavButton
              tooltip="Go to previous dashboard"
              classSuffix="tight"
              icon="step-backward"
              onClick={this.onPlaylistPrev}
            />
            <DashNavButton
              tooltip="Stop playlist"
              classSuffix="tight"
              icon="square-shape"
              onClick={this.onPlaylistStop}
            />
            <DashNavButton
              tooltip="Go to next dashboard"
              classSuffix="tight"
              icon="forward"
              onClick={this.onPlaylistNext}
            />
          </div>
        )}

        <div className="navbar-buttons navbar-buttons--actions">{this.renderRightActionsButton()}</div>

        {!dashboard.timepicker.hidden && (
          <div className="navbar-buttons">
            <DashNavTimeControls
              dashboard={dashboard}
              location={location}
              updateLocation={updateLocation}
              onChangeTimeZone={updateTimeZoneForSession}
            />
          </div>
        )}
      </div>
    );
  }
}

const mapStateToProps = (state: StoreState) => ({
  location: state.location,
});

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  updateLocation,
  updateTimeZoneForSession,
};

export default connect(mapStateToProps, mapDispatchToProps)(DashNav);
