define([
  'angular',
  'lodash'
],
function (angular, _) {
  'use strict';

  var module = angular.module('grafana.controllers');

  module.controller('EndpointSummaryCtrl', function($scope, $http, backendSrv, $location, $routeParams) {
    $scope.init = function() {
      $scope.endpoints = [];
      $scope.monitors = {};
      $scope.monitor_health = {};
      $scope.monitor_types = {};
      $scope.monitor_types_by_name = {};
      $scope.endpoint = null;
      $scope.refreshTime = new Date();
      $scope.getMonitorTypes();
      var promise = $scope.getEndpoints();
      promise.then(function() {
        $scope.getEndpoint($routeParams.id);
      });

    };

    $scope.getEndpoints = function() {
      var promise = backendSrv.get('/api/endpoints')
      promise.then(function(endpoints) {
        $scope.endpoints = endpoints;
      });
      return promise;
    };
    $scope.tagsUpdated = function(newVal) {
      backendSrv.post("/api/endpoints", $scope.endpoint);
    }

    $scope.getMonitorTypes = function() {
      backendSrv.get('/api/monitor_types').then(function(types) {
        _.forEach(types, function(type) {
          $scope.monitor_types[type.id] = type;
          $scope.monitor_types_by_name[type.name] = type;
        });
        console.log("monitor_types ready");
      });
    };

    $scope.getEndpoint = function(id) {
      _.forEach($scope.endpoints, function(endpoint) {
        if (endpoint.id == id) {
          $scope.endpoint = endpoint;
          //get monitors for this endpoint.
          backendSrv.get('/api/monitors?endpoint_id='+id).then(function(monitors) {
            _.forEach(monitors, function(monitor) {
              $scope.monitors[monitor.monitor_type_id] = monitor;
            });
          });
          backendSrv.get('/api/endpoints/'+id+'/health').then(function(health) {
            var healthByCheck = {};
            _.forEach(health, function(checkState) {
              if (!(checkState.monitor_id in healthByCheck)) {
                healthByCheck[checkState.monitor_id] = [];
              }
              healthByCheck[checkState.monitor_id].push(checkState);
            });
            $scope.monitor_health = healthByCheck;
          });
        }
      });
    };

    $scope.getMonitorByTypeName = function(name) {
      if (name in $scope.monitor_types_by_name) {
        var type = $scope.monitor_types_by_name[name];
        return $scope.monitors[type.id];
      }
      return undefined;
    }

    $scope.monitorStateTxt = function(mon) {
      if (typeof(mon) != "object") {
        return "nodata";
      }
      if (mon.state < 0 || mon.state > 2) {
        return 'nodata';
      }
      var states = ["online", "warn", "critical"];
      return states[mon.state];
    }

    $scope.stateChangeStr = function(mon) {
      if (typeof(mon) != "object") {
        return "";
      }
      var duration = new Date().getTime() - new Date(mon.state_change).getTime();
      if (duration < 10000) {
        return "a few seconds ago";
      }
      if (duration < 60000) {
        var secs = Math.floor(duration/1000);
        return secs + " seconds ago";
      }
      if (duration < 3600000) {
        var mins = Math.floor(duration/1000/60);
        return mins + " minutes ago";
      }
      if (duration < 86400000) {
        var hours = Math.floor(duration/1000/60/60);
        return hours + " hours ago";
      }
      var days = Math.floor(duration/1000/60/60/24);
      return days + " days ago";
    }

    $scope.setEndpoint = function(id) {
      $location.path('/endpoints/summary/'+id);
    }

    $scope.slug = function(name) {
      var label = name.toLowerCase();
      var re = new RegExp("[^\\w-]+");
      var re2 = new RegExp("\\s");
      var slug = label.replace(re, "_").replace(re2, "-");
      return slug;
    }

    $scope.gotoDashboard = function(endpoint, type) {
      if (!type) {
        type = 'summary';
      }
      var search = {
        "var-collector": "All",
        "var-endpoint": $scope.slug($scope.endpoint.name)
      };
      switch(type) {
        case "summary":
          $location.path("/dashboard/raintank/statusboard").search(search);
          break;
        case "ping":
          $location.path("/dashboard/raintank/rt-endpoint-icmp").search(search);
          break;
        case "dns":
          $location.path("/dashboard/raintank/rt-endpoint-dns").search(search);
          break;
        case "http":
          search['var-protocol'] = "http";
          $location.path("/dashboard/raintank/rt-endpoint-web").search(search);
          break;
        case "https":
          search['var-protocol'] = "https";
          $location.path("/dashboard/raintank/rt-endpoint-web").search(search);
          break;
        default:
          $location.path("/dashboard/raintank/statusboard").search(search);
          break;
      }
    };

    $scope.refresh = function() {
      $scope.getEndpoint($scope.endpoint.id);
      $scope.refreshTime = new Date();
    }
    $scope.init();
  });
});
