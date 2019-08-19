import React, { PureComponent } from 'react';
import classNames from 'classnames';

import { QueriesTab } from './QueriesTab';
import { ReferencedPanelEditor } from './ReferencedPanelEditor';
import VisualizationTab from './VisualizationTab';
import { GeneralTab } from './GeneralTab';
import { AlertTab } from '../../alerting/AlertTab';

import config from 'app/core/config';
import { store } from 'app/store/store';
import { updateLocation } from 'app/core/actions';
import { AngularComponent } from '@grafana/runtime';

import { PanelModel } from '../state/PanelModel';
import { DashboardModel } from '../state/DashboardModel';
import { Tooltip, PanelPlugin, PanelPluginMeta } from '@grafana/ui';

interface PanelEditorProps {
  panel: PanelModel;
  dashboard: DashboardModel;
  plugin: PanelPlugin;
  angularPanel?: AngularComponent;
  onPluginTypeChange: (newType: PanelPluginMeta) => void;
}

interface PanelEditorTab {
  id: string;
  text: string;
}

enum PanelEditorTabIds {
  Queries = 'queries',
  Visualization = 'visualization',
  Advanced = 'advanced',
  Alert = 'alert',
}

interface PanelEditorTab {
  id: string;
  text: string;
}

const panelEditorTabTexts = {
  [PanelEditorTabIds.Queries]: 'Queries',
  [PanelEditorTabIds.Visualization]: 'Visualization',
  [PanelEditorTabIds.Advanced]: 'General',
  [PanelEditorTabIds.Alert]: 'Alert',
};

const getPanelEditorTab = (tabId: PanelEditorTabIds): PanelEditorTab => {
  return {
    id: tabId,
    text: panelEditorTabTexts[tabId],
  };
};

export class PanelEditor extends PureComponent<PanelEditorProps> {
  constructor(props: PanelEditorProps) {
    super(props);
  }

  onReferenceChanged = () => {
    const { onPluginTypeChange, plugin, panel } = this.props;
    onPluginTypeChange(plugin.meta); // Will reload the reference
    panel.render();
    this.forceUpdate();
  };

  onChangeTab = (tab: PanelEditorTab) => {
    store.dispatch(
      updateLocation({
        query: { tab: tab.id, openVizPicker: null },
        partial: true,
      })
    );
    this.forceUpdate();
  };

  renderCurrentTab(activeTab: string) {
    const { panel, dashboard, onPluginTypeChange, plugin, angularPanel } = this.props;

    switch (activeTab) {
      case 'advanced':
        return <GeneralTab panel={panel} onReferenceChanged={this.onReferenceChanged} />;
      case 'queries':
        return <QueriesTab panel={panel} dashboard={dashboard} />;
      case 'alert':
        return <AlertTab angularPanel={angularPanel} dashboard={dashboard} panel={panel} />;
      case 'visualization':
        return (
          <VisualizationTab
            panel={panel}
            dashboard={dashboard}
            plugin={plugin}
            onPluginTypeChange={onPluginTypeChange}
            angularPanel={angularPanel}
          />
        );
      default:
        return null;
    }
  }

  render() {
    if (this.props.panel.reference) {
      return <ReferencedPanelEditor panel={this.props.panel} onReferenceChanged={this.onReferenceChanged} />;
    }

    const { plugin } = this.props;
    let activeTab: PanelEditorTabIds = store.getState().location.query.tab || PanelEditorTabIds.Queries;

    const tabs: PanelEditorTab[] = [
      getPanelEditorTab(PanelEditorTabIds.Queries),
      getPanelEditorTab(PanelEditorTabIds.Visualization),
      getPanelEditorTab(PanelEditorTabIds.Advanced),
    ];

    // handle panels that do not have queries tab
    if (plugin.meta.skipDataQuery) {
      // remove queries tab
      tabs.shift();
      // switch tab
      if (activeTab === PanelEditorTabIds.Queries) {
        activeTab = PanelEditorTabIds.Visualization;
      }
    }

    if (config.alertingEnabled && plugin.meta.id === 'graph') {
      tabs.push(getPanelEditorTab(PanelEditorTabIds.Alert));
    }

    return (
      <div className="panel-editor-container__editor">
        <div className="panel-editor-tabs">
          {tabs.map(tab => {
            return <TabItem tab={tab} activeTab={activeTab} onClick={this.onChangeTab} key={tab.id} />;
          })}
        </div>
        <div className="panel-editor__right">{this.renderCurrentTab(activeTab)}</div>
      </div>
    );
  }
}

interface TabItemParams {
  tab: PanelEditorTab;
  activeTab: string;
  onClick: (tab: PanelEditorTab) => void;
}

function TabItem({ tab, activeTab, onClick }: TabItemParams) {
  const tabClasses = classNames({
    'panel-editor-tabs__link': true,
    active: activeTab === tab.id,
  });

  return (
    <div className="panel-editor-tabs__item" onClick={() => onClick(tab)}>
      <a className={tabClasses} aria-label={`${tab.text} tab button`}>
        <Tooltip content={`${tab.text}`} placement="auto">
          <i className={`gicon gicon-${tab.id}${activeTab === tab.id ? '-active' : ''}`} />
        </Tooltip>
      </a>
    </div>
  );
}
