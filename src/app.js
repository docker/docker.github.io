require.main.paths.splice(0, 0, process.env.NODE_PATH);
import remote from 'remote';
var Menu = remote.require('menu');
import React from 'react';
import SetupStore from './stores/SetupStore';
import ipc from 'ipc';
import machine from './utils/DockerMachineUtil';
import metrics from './utils/MetricsUtil';
import template from './menutemplate';
import webUtil from './utils/WebUtil';
import hubUtil from './utils/HubUtil';
import setupUtil from './utils/SetupUtil';
import request from 'request';
import docker from './utils/DockerUtil';
import hub from './utils/HubUtil';
import Router from 'react-router';
import createHashHistory from 'history/lib/createHashHistory'
import routes from './routes';
import routerContainer from './router';
import repositoryActions from './actions/RepositoryActions';
var app = remote.require('app');

hubUtil.init();

if (hubUtil.loggedin()) {
  repositoryActions.repos();
}

repositoryActions.recommended();

webUtil.addWindowSizeSaving();
webUtil.addLiveReload();
webUtil.addBugReporting();
webUtil.disableGlobalBackspace();

Menu.setApplicationMenu(Menu.buildFromTemplate(template()));

metrics.track('Started App');
metrics.track('app heartbeat');
setInterval(function () {
  metrics.track('app heartbeat');
}, 14400000);

let history = createHashHistory()
React.render(<Router history={history}>{routes}</Router>, document.body)

setupUtil.setup().then(() => {
  Menu.setApplicationMenu(Menu.buildFromTemplate(template()));
  docker.init();
  if (!hub.prompted() && !hub.loggedin()) {
    history.replaceState(null, '/account/login');
  } else {
    history.replaceState(null, '/containers/new');
  }
}).catch(err => {
  metrics.track('Setup Failed', {
    step: 'catch',
    message: err.message
  });
  throw err;
});

ipc.on('application:quitting', () => {
  if (localStorage.getItem('settings.closeVMOnQuit') === 'true') {
    machine.stop();
  }
});
