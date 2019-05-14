import angular from 'angular';
import config from 'app/core/config';
import moment from 'moment';
/** @ngInject */
export function ShareModalCtrl($scope, $rootScope, $location, $timeout, timeSrv, templateSrv, linkSrv) {
    $scope.options = {
        forCurrent: true,
        includeTemplateVars: true,
        theme: 'current',
    };
    $scope.editor = { index: $scope.tabIndex || 0 };
    $scope.init = function () {
        $scope.panel = $scope.model && $scope.model.panel ? $scope.model.panel : $scope.panel; // React pass panel and dashboard in the "model" property
        $scope.dashboard = $scope.model && $scope.model.dashboard ? $scope.model.dashboard : $scope.dashboard; // ^
        $scope.modeSharePanel = $scope.panel ? true : false;
        $scope.tabs = [{ title: 'Link', src: 'shareLink.html' }];
        if ($scope.modeSharePanel) {
            $scope.modalTitle = 'Share Panel';
            $scope.tabs.push({ title: 'Embed', src: 'shareEmbed.html' });
        }
        else {
            $scope.modalTitle = 'Share';
        }
        if (!$scope.dashboard.meta.isSnapshot) {
            $scope.tabs.push({ title: 'Snapshot', src: 'shareSnapshot.html' });
        }
        if (!$scope.dashboard.meta.isSnapshot && !$scope.modeSharePanel) {
            $scope.tabs.push({ title: 'Export', src: 'shareExport.html' });
        }
        $scope.buildUrl();
    };
    $scope.buildUrl = function () {
        var baseUrl = $location.absUrl();
        var queryStart = baseUrl.indexOf('?');
        if (queryStart !== -1) {
            baseUrl = baseUrl.substring(0, queryStart);
        }
        var params = angular.copy($location.search());
        var range = timeSrv.timeRange();
        params.from = range.from.valueOf();
        params.to = range.to.valueOf();
        params.orgId = config.bootData.user.orgId;
        if ($scope.options.includeTemplateVars) {
            templateSrv.fillVariableValuesForUrl(params);
        }
        if (!$scope.options.forCurrent) {
            delete params.from;
            delete params.to;
        }
        if ($scope.options.theme !== 'current') {
            params.theme = $scope.options.theme;
        }
        if ($scope.modeSharePanel) {
            params.panelId = $scope.panel.id;
            params.fullscreen = true;
        }
        else {
            delete params.panelId;
            delete params.fullscreen;
        }
        $scope.shareUrl = linkSrv.addParamsToUrl(baseUrl, params);
        var soloUrl = baseUrl.replace(config.appSubUrl + '/dashboard/', config.appSubUrl + '/dashboard-solo/');
        soloUrl = soloUrl.replace(config.appSubUrl + '/d/', config.appSubUrl + '/d-solo/');
        delete params.fullscreen;
        delete params.edit;
        soloUrl = linkSrv.addParamsToUrl(soloUrl, params);
        $scope.iframeHtml = '<iframe src="' + soloUrl + '" width="450" height="200" frameborder="0"></iframe>';
        $scope.imageUrl = soloUrl.replace(config.appSubUrl + '/dashboard-solo/', config.appSubUrl + '/render/dashboard-solo/');
        $scope.imageUrl = $scope.imageUrl.replace(config.appSubUrl + '/d-solo/', config.appSubUrl + '/render/d-solo/');
        $scope.imageUrl += '&width=1000&height=500' + $scope.getLocalTimeZone();
    };
    // This function will try to return the proper full name of the local timezone
    // Chrome does not handle the timezone offset (but phantomjs does)
    $scope.getLocalTimeZone = function () {
        var utcOffset = '&tz=UTC' + encodeURIComponent(moment().format('Z'));
        // Older browser does not the internationalization API
        if (!window.Intl) {
            return utcOffset;
        }
        var dateFormat = window.Intl.DateTimeFormat();
        if (!dateFormat.resolvedOptions) {
            return utcOffset;
        }
        var options = dateFormat.resolvedOptions();
        if (!options.timeZone) {
            return utcOffset;
        }
        return '&tz=' + encodeURIComponent(options.timeZone);
    };
    $scope.getShareUrl = function () {
        return $scope.shareUrl;
    };
}
angular.module('grafana.controllers').controller('ShareModalCtrl', ShareModalCtrl);
//# sourceMappingURL=ShareModalCtrl.js.map