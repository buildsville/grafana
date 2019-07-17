import angular from 'angular';
import config from 'app/core/config';
import { appEvents } from 'app/core/core';

export class ShareSlackCtrl {
  /** @ngInject */
  constructor($scope, $location, backendSrv, timeSrv, templateSrv, linkSrv) {
    $scope.theme = "current";
    $scope.timeRange = true;
    $scope.tempVar = true;

    const options = config.slackShare.channels.split(',');

    $scope.channelOptions = [];
    $scope.channelOptions.push({ text: 'Default channel', value: '__default__' });
    for (let index = 0; index < options.length; index++) {
      const o = options[index];
      $scope.channelOptions.push({ text: o, value: o });
    }
    $scope.channel = '__default__';

    $scope.init = () => {
      $scope.panel = $scope.panel;
    };

    $scope.sharePanel = () => {
      const range = timeSrv.timeRange();
      const channel = $scope.channel;
      const theme = $scope.theme;

      const path = $location.path();
      const pathArr = path.split('/');
      const uid = pathArr[pathArr.length - 2];
      const slug = pathArr[pathArr.length - 1];

      const params = angular.copy($location.search());
      params.from = range.from.valueOf();
      params.to = range.to.valueOf();
      params.orgId = config.bootData.user.orgId;
      params.panelId = $scope.panel.id;
      params.fullscreen = true;
      if ($scope.tempVar) {
        templateSrv.fillVariableValuesForUrl(params);
      }
      if (!$scope.timeRange) {
        delete params.from;
        delete params.to;
      }
      if ($scope.theme !== "current") {
        params.theme = theme;
      }

      const paramStr = linkSrv.addParamsToUrl("", params);

      const payload = {
        rawURL: $location.absUrl(),
        uid: uid,
        slug: slug,
        panelName: $scope.panel.title,
        channel: channel,
        param: paramStr,
      };

      $scope.loading = true;
      backendSrv.post(`/share/slack`, payload).then(() => {
        const msg = 'Panel shared to ' + channel;
        appEvents.emit('alert-success', [msg, '']);
        $scope.loading = false;
      });
    };
  }
}

angular.module('grafana.controllers').controller('ShareSlackCtrl', ShareSlackCtrl);
