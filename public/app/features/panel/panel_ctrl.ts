import _ from 'lodash';
import { escapeHtml, sanitize } from 'app/core/utils/text';

import config from 'app/core/config';
import { Emitter, profiler } from 'app/core/core';
import {
  copyPanel as copyPanelUtil,
  duplicatePanel,
  editPanelJson as editPanelJsonUtil,
  removePanel,
  sharePanel as sharePanelUtil,
} from 'app/features/dashboard/utils/panel';
import { auto } from 'angular';
import { TemplateSrv } from '../templating/template_srv';
import { getPanelLinksSupplier } from './panellinks/linkSuppliers';
import { AppEvent, PanelEvents, PanelPluginMeta, renderMarkdown, AngularPanelMenuItem } from '@grafana/data';
import { getLocationSrv } from '@grafana/runtime';
import { DashboardModel } from '../dashboard/state';

export class PanelCtrl {
  panel: any;
  error: any;
  dashboard: DashboardModel;
  pluginName: string;
  pluginId: string;
  editorTabs: any;
  $scope: any;
  $injector: auto.IInjectorService;
  $location: any;
  $timeout: any;
  editModeInitiated: boolean;
  height: any;
  containerHeight: any;
  events: Emitter;
  loading: boolean;
  timing: any;

  constructor($scope: any, $injector: auto.IInjectorService) {
    this.$injector = $injector;
    this.$location = $injector.get('$location');
    this.$scope = $scope;
    this.$timeout = $injector.get('$timeout');
    this.editorTabs = [];
    this.events = this.panel.events;
    this.timing = {}; // not used but here to not break plugins

    const plugin = config.panels[this.panel.type];
    if (plugin) {
      this.pluginId = plugin.id;
      this.pluginName = plugin.name;
    }

    $scope.$on(PanelEvents.componentDidMount.name, () => this.panelDidMount());
  }

  panelDidMount() {
    this.events.emit(PanelEvents.componentDidMount);
    this.dashboard.panelInitialized(this.panel);
  }

  renderingCompleted() {
    profiler.renderingCompleted();
  }

  refresh() {
    this.panel.refresh();
  }

  publishAppEvent<T>(event: AppEvent<T>, payload?: T) {
    this.$scope.$root.appEvent(event, payload);
  }

  changeView(fullscreen: boolean, edit: boolean) {
    this.publishAppEvent(PanelEvents.panelChangeView, {
      fullscreen,
      edit,
      panelId: this.panel.id,
    });
  }

  viewPanel() {
    this.changeView(true, false);
  }

  editPanel() {
    this.changeView(true, true);
  }

  exitFullscreen() {
    this.changeView(false, false);
  }

  initEditMode() {
    if (!this.editModeInitiated) {
      this.editModeInitiated = true;
      this.events.emit(PanelEvents.editModeInitialized);
    }
  }

  addEditorTab(title: string, directiveFn: any, index?: number, icon?: any) {
    const editorTab = { title, directiveFn, icon };

    if (_.isString(directiveFn)) {
      editorTab.directiveFn = () => {
        return { templateUrl: directiveFn };
      };
    }

    if (index) {
      this.editorTabs.splice(index, 0, editorTab);
    } else {
      this.editorTabs.push(editorTab);
    }
  }

  getExtendedMenu() {
    const menu: AngularPanelMenuItem[] = [];
    this.events.emit(PanelEvents.initPanelActions, menu);
    return menu;
  }

  // Override in sub-class to add items before extended menu
  async getAdditionalMenuItems(): Promise<any[]> {
    return [];
  }

  otherPanelInFullscreenMode() {
    return this.dashboard.meta.fullscreen && !this.panel.fullscreen;
  }

  render(payload?: any) {
    this.events.emit(PanelEvents.render, payload);
  }

  duplicate() {
    duplicatePanel(this.dashboard, this.panel);
  }

  removePanel() {
    removePanel(this.dashboard, this.panel, true);
  }

  editPanelJson() {
    editPanelJsonUtil(this.dashboard, this.panel);
  }

  copyPanel() {
    copyPanelUtil(this.panel);
  }

  sharePanel() {
    sharePanelUtil(this.dashboard, this.panel);
  }

  inspectPanel() {
    getLocationSrv().update({
      query: {
        inspect: this.panel.id,
      },
      partial: true,
    });
  }

  getInfoMode() {
    if (this.error) {
      return 'error';
    }
    if (!!this.panel.description) {
      return 'info';
    }
    if (this.panel.links && this.panel.links.length) {
      return 'links';
    }
    return '';
  }

  getInfoContent(options: { mode: string }) {
    const { panel } = this;
    let markdown = panel.description || '';

    if (options.mode === 'tooltip') {
      markdown = this.error || panel.description || '';
    }

    const templateSrv: TemplateSrv = this.$injector.get('templateSrv');
    const interpolatedMarkdown = templateSrv.replace(markdown, panel.scopedVars);
    let html = '<div class="markdown-html panel-info-content">';

    const md = renderMarkdown(interpolatedMarkdown);
    html += md;

    if (panel.links && panel.links.length > 0) {
      const interpolatedLinks = getPanelLinksSupplier(panel).getLinks();
      html += '<ul class="panel-info-corner-links">';
      for (const link of interpolatedLinks) {
        html +=
          '<li><a class="panel-menu-link" href="' +
          escapeHtml(link.href) +
          '" target="' +
          escapeHtml(link.target) +
          '">' +
          escapeHtml(link.title) +
          '</a></li>';
      }
      html += '</ul>';
    }

    html += '</div>';

    return config.disableSanitizeHtml ? html : sanitize(html);
  }

  // overriden from react
  onPluginTypeChange = (plugin: PanelPluginMeta) => {};
}
