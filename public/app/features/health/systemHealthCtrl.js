define([
    'angular',
    'lodash'
  ],
  function (angular, _) {
    'use strict';
    var module = angular.module('grafana.controllers');
    var panelMeta = {
      "collapse": false,
      "editable": false,
      "height": "300px",
      "panels": [
        {
          "aliasColors": {
            "test_health": "#EAB839"
          },
          "bars": false,
          "datasource": null,
          "editable": true,
          "error": false,
          "fill": 0,
          "id": 2,
          "lines": true,
          "linewidth": 2,
          "nullPointMode": "connected",
          "percentage": false,
          "pointradius": 4,
          "points": true,
          "renderer": "flot",
          "seriesOverrides": [],
          "span": 12,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "aggregator": "sum",
              "downsampleAggregator": "avg",
              "downsampleInterval": "1m",
              "errors": {},
              "metric": "test_health",
            }
          ],
          "title": "历史健康指数趋势",
          "tooltip": {
            "shared": true,
            "value_type": "cumulative"
          },
          "type": "graph",
          "x-axis": true,
          "y-axis": true,
          "y_formats": [
            "short",
            "short"
          ],
          "transparent": true
        }
      ],
      "showTitle": false,
      "title": "New row"
    };

    function getLeveal(score) {
      if (!_.isNumber(score) && _.isNaN(score) && _.isEmpty(score)) {
        return "无";
      }
      if (score > 75) {
        return "优";
      } else if (score > 50) {
        return "良";
      } else if (score > 25) {
        return "中";
      } else {
        return "差";
      }
    }

    module.controller('SystemHealthCtrl', function ($scope, contextSrv, healthSrv, backendSrv, $rootScope) {
      $scope.init = function () {
        $scope.system = backendSrv.getSystemById(contextSrv.system);
        $scope.initDashboard({
          meta: {canStar: false, canShare: false, canEdit: false},
          dashboard: {
            system: contextSrv.system,
            title: "健康状态",
            id: "name",
            rows: [panelMeta],
            time: {from: "now-2h", to: "now"}
          }
        }, $scope);

        $scope.reload();
      };

      $scope.reload = function() {
        healthSrv.load().then(function (data) {
          $scope.applicationHealth = Math.floor(data.health);
          $scope.leveal = getLeveal($scope.applicationHealth);

          $scope.includeMetricsData = healthSrv.floor(data.includedMetricHealths);
          $scope.excludeMetricsData = healthSrv.floor(data.excludedMetricHealths);
          $scope.excludeMetricLength = _.size($scope.excludeMetricsData);
        });
      };

      $scope.changeExcludeMetrics = function () {
        $scope.appEvent('show-modal', {
          src: './app/partials/exclude_metrics.html',
          modalClass: 'modal-no-header confirm-modal',
          scope: $scope.$new(),
        });
      };

      $scope.init();
    });
  });
