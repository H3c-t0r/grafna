define([
    'angular',
    'lodash'
  ],
  function (angular, _) {
    'use strict';

    var module = angular.module('grafana.controllers');

    module.controller('AnomalyCtrl', function ($scope, healthSrv) {

      $scope.init = function () {
        healthSrv.load().then(function (data) {
          $scope.anomalyList = data;
        });
      };
      $scope.init();
    });
  });
