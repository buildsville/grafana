import angular from 'angular';
import config from 'app/core/config';
import { appEvents } from 'app/core/core';

export class ShareSlackCtrl {
  /** @ngInject */
  constructor($scope, $location, backendSrv, timeSrv) {
    $scope.theme = "current";

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
      const payload = {
        rawURL: $location.absUrl(),
        uid: uid,
        slug: slug,
        panelId: $scope.panel.id,
        panelName: $scope.panel.title,
        channel: channel,
        from: range.from.valueOf(),
        to: range.to.valueOf(),
        theme: theme,
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
