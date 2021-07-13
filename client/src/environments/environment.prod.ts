import packageInfo from '../../package.json';

export const environment = {
  production: true,
  serverEndpoint: '/api/v2',
  apiRequestTimeoutInSecs: 60,
  timezone: 'Europe/Berlin',
  appVersion: packageInfo.version
};
