define([
  'angular',
  'lodash',
  'config',
],
function (angular, _, config) {
  'use strict';

  var module = angular.module('grafana.services');

  module.service('panelSrv', function($rootScope, $timeout, datasourceSrv, $q, $location, healthSrv, contextSrv, alertSrv, integrateSrv) {

    this.init = function($scope) {

      if (!$scope.panel.span) { $scope.panel.span = 12; }

      $scope.menuItemShow = false;

      $scope.inspector = {};

      $scope.showRightMenu = function() {
        $scope.menuItemShow = !$scope.menuItemShow;
      };

      $scope.editPanel = function() {
        $scope.toggleFullscreen(true);
      };

      $scope.sharePanel = function() {
        $scope.appEvent('show-modal', {
          src: './app/features/dashboard/partials/shareModal.html',
          scope: $scope.$new()
        });
      };

      $scope.editPanelJson = function() {
        $scope.appEvent('show-json-editor', { object: $scope.panel, updateHandler: $scope.replacePanel });
      };

      $scope.duplicatePanel = function() {
        $scope.dashboard.duplicatePanel($scope.panel, $scope.row);
      };

      $scope.updateColumnSpan = function(span) {
        $scope.updatePanelSpan($scope.panel, span);

        $timeout(function() {
          $scope.$broadcast('render');
        });
      };

      $scope.addDataQuery = function(datasource) {
        $scope.dashboard.addDataQueryTo($scope.panel, datasource);
      };

      $scope.removeDataQuery = function (query) {
        $scope.dashboard.removeDataQuery($scope.panel, query);
        $scope.get_data();
      };

      $scope.duplicateDataQuery = function(query) {
        $scope.dashboard.duplicateDataQuery($scope.panel, query);
      };

      $scope.moveDataQuery = function(fromIndex, toIndex) {
        $scope.dashboard.moveDataQuery($scope.panel, fromIndex, toIndex);
      };

      $scope.setDatasource = function(datasource) {
        // switching to mixed
        if (datasource.meta.mixed) {
          _.each($scope.panel.targets, function(target) {
            target.datasource = $scope.panel.datasource;
            if (target.datasource === null) {
              target.datasource = config.defaultDatasource;
            }
          });
        }
        // switching from mixed
        else if ($scope.datasource && $scope.datasource.meta.mixed) {
          _.each($scope.panel.targets, function(target) {
            delete target.datasource;
          });
        }

        $scope.panel.datasource = datasource.value;
        $scope.datasource = null;
        $scope.get_data();
      };

      $scope.toggleEditorHelp = function(index) {
        if ($scope.editorHelpIndex === index) {
          $scope.editorHelpIndex = null;
          return;
        }
        $scope.editorHelpIndex = index;
      };

      $scope.isNewPanel = function() {
        return $scope.panel.title === config.new_panel_title;
      };

      $scope.toggleFullscreen = function(edit) {
        $scope.dashboardViewState.update({ fullscreen: true, edit: edit, panelId: $scope.panel.id });
      };

      $scope.otherPanelInFullscreenMode = function() {
        return $scope.dashboardViewState.fullscreen && !$scope.fullscreen;
      };

      $scope.getCurrentDatasource = function() {
        if ($scope.datasource) {
          return $q.when($scope.datasource);
        }

        return datasourceSrv.get($scope.panel.datasource);
      };

      $scope.panelRenderingComplete = function() {
        $rootScope.performance.panelsRendered++;
      };

      $scope.decompose = function() {
        $scope.dashboardViewState.update({ fullscreen: false, edit: false, panelId: null });
        window.decomposeTarget = $scope.panel.targets[0];
        $location.path("/decompose");
      };

      $scope.get_data = function() {
        if ($scope.otherPanelInFullscreenMode()) { return; }

        if ($scope.panel.snapshotData) {
          if ($scope.loadSnapshot) {
            $scope.loadSnapshot($scope.panel.snapshotData);
          }
          return;
        }

        delete $scope.panelMeta.error;
        $scope.panelMeta.loading = true;
        $scope.panelMeta.info = false;

        healthSrv.transformMetricType($scope.dashboard).then(function () {
          $scope.getCurrentDatasource().then(function (datasource) {
            $scope.datasource = datasource;
            return $scope.refreshData($scope.datasource) || $q.when({});
          }).then(function () {
            $scope.panelMeta.loading = false;
            $scope.titleInit($scope.panel);
          }, function (err) {
            console.log('Panel data error:', err);
            $scope.panelMeta.loading = false;
            $scope.panelMeta.error = err.message || "Timeseries data request error";
            $scope.inspector.error = err;
          });
        });
      };

      $scope.isShowInfo = function (event) {
        if (event.type === "click") {
          $scope.helpShow = true;
        } else {
          $scope.helpShow = false;
        }
      };

      $scope.titleInit = function (panel) {
        $scope.helpInfo = panel.helpInfo;
        $scope.showMenu = (panel.type == "graph" && $scope.panelMeta.loading == false ? true : false);
        var path = $location.path();
        $scope.associateMenu = /^\/anomaly/.test(path);
      };

      $scope.systemsMap = contextSrv.systemsMap[0];

      $scope.associateLink = function () {
        try{
          var host = $scope.panel.targets[0].tags.host;
          var metric = $scope.panel.targets[0].metric;
          if(host && metric) {
            var link = '/alerts/association/'+host+'/100/'+ $scope.systemsMap.Id+'.'+$scope.systemsMap.OrgId+'.'+metric;
            $location.path(link);
          }
        }catch(err){
          var reg = /\'(.*?)\'/g;
          var msg = "图表中缺少"+err.toString().match(reg)[0]+"配置";
          alertSrv.set("参数缺失", msg, "warning", 4000);
        }
      };

      if ($scope.refreshData) {
        $scope.$on("refresh", $scope.get_data);
      }

      // Post init phase
      $scope.fullscreen = false;
      $scope.editor = { index: 1 };

      $scope.dashboardViewState.registerPanel($scope);
      $scope.datasources = datasourceSrv.getMetricSources();

      if (!$scope.skipDataOnInit) {
        $timeout(function() {
          $scope.get_data();
        }, 30);
      }

      $scope.toIntegrate = function() {
        try{
          if (_.isNull($scope.panel.targets[0].metric) || _.isNull($scope.panel.targets[0].tags.host)) {
            return;
          }
          integrateSrv.format.targets = $scope.panel.targets;
          $location.path("/integrate");
        }catch(e){
          $scope.appEvent('alert-warning', ['日志分析跳转失败', '可能缺少指标名/主机名']);
        };
      }
    };
  });

});