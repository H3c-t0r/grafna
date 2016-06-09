define([
  'angular',
  'lodash'
],
function (angular, _) {
  'use strict';

  var module = angular.module('grafana.controllers');

  module.controller('AlertEditCtrl', function($scope, $routeParams, $location, alertMgrSrv, alertSrv, datasourceSrv) {

    $scope.init = function() {
      $scope.datasource = null;
      _.each(datasourceSrv.getAll(), function(ds) {
        if (ds.type === 'opentsdb') {
          datasourceSrv.get(ds.name).then(function(datasource) {
            $scope.datasource = datasource;
          });
        }
      });

      $scope.alertDef = alertMgrSrv.get($routeParams.id);
      $scope.isNew = !$scope.alertDef;
      if ($scope.isNew) {
        $scope.alertDef = {};
        $scope.alertDef.service = "com.test";
        $scope.alertDef.alertDetails = {};
        $scope.alertDef.alertDetails.cluster = "cluster1";
        $scope.alertDef.alertDetails.hosts = null;
        $scope.alertDef.alertDetails.membership = "*";
        $scope.alertDef.alertDetails.monitoringScope = "HOST";
        $scope.alertDef.alertDetails.hostQuery = {};
        $scope.alertDef.alertDetails.hostQuery.metricQueries = [];
      }
    };

    $scope.saveChanges = function() {
      var milliseconds = (new Date).getTime();
      if ($scope.isNew) {
        //if it is new, we need to fill in some hard-coded value for now.
        $scope.alertDef.creationTime = milliseconds;
        $scope.alertDef.modificationTime = milliseconds;
      } else {
        $scope.alertDef.modificationTime = milliseconds;
      }

      alertMgrSrv.save($scope.alertDef).then(function onSuccess() {
        $location.path("alerts");
      }, function onFailed(response) {
        alertSrv.set("error", response.status + " " + (response.data || "Request failed"), response.severity, 10000);
      });
    };

    $scope.getTextValues = function(metricFindResult) {
      return _.map(metricFindResult, function(value) { return value.text; });
    };

    $scope.suggestMetrics = function(query, callback) {
      $scope.datasource.metricFindQuery('metrics(' + query + ')')
        .then($scope.getTextValues)
        .then(callback);
    };

    $scope.init();
  });
});
