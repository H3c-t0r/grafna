import '@testing-library/jest-dom';
import i18next from 'i18next';
import failOnConsole from 'jest-fail-on-console';
import { initReactI18next } from 'react-i18next';

import { matchers } from './matchers';

if (process.env.CI) {
  failOnConsole({
    shouldFailOnLog: true,
    shouldFailOnDebug: true,
    shouldFailOnInfo: true,
  });
}

expect.extend(matchers);

i18next.use(initReactI18next).init({
  resources: {},
  returnEmptyString: false,
  lng: 'en-US', // this should be the locale of the phrases in our source JSX
});

// Mock out the worker that detects changes in the dashboard
// The mock is needed because JSDOM does not support workers and
// the factory uses import.meta.url so we can't use it in CommonJS modules.
jest.mock('app/features/dashboard-scene/saving/createDetectChangesWorker.ts');
