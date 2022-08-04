declare let __webpack_public_path__: string;
declare let __webpack_nonce__: string;

if (window.__POWERED_BY_QIANKUN__) {
  __webpack_public_path__ = window.__INJECTED_PUBLIC_PATH_BY_QIANKUN__;
}

import React from 'react';
import ReactDOM from 'react-dom';

import { config } from '@grafana/runtime';

import { FNDashboard } from './FNDashboard';

// set featureToggles
config.featureToggles = {
  ...config.featureToggles,
  isPublicDashboardView: true,
};
config.isPublicDashboardView = true;

/**
 * The bootstrap will only be called once when the child application is initialized.
 * The next time the child application re-enters, the mount hook will be called directly, and bootstrap will not be triggered repeatedly.
 * Usually we can do some initialization of global variables here,
 * such as application-level caches that will not be destroyed during the unmount phase.
 */
export async function bootstrap() {
  console.log('react app bootstraped');
}

/**
 * The mount method is called every time the application enters,
 * usually we trigger the application's rendering method here.
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function mount(props: any) {
  console.log('react app mount: ', props.container);
  ReactDOM.render(
    React.createElement(FNDashboard),
    props.container ? props.container.querySelector('#reactRoot') : document.getElementById('reactRoot')
  );
}

/**
 * Methods that are called each time the application is switched/unloaded,
 * usually in this case we uninstall the application instance of the subapplication.
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function unmount(props: any) {
  const container = props.container
    ? props.container.querySelector('#reactRoot')
    : document.getElementById('reactRoot');
  if (container) {
    ReactDOM.unmountComponentAtNode(container);
  }
}

/**
 * Optional lifecycle，just available with loadMicroApp way
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function update(props: any) {
  console.log('update props', props);
}
