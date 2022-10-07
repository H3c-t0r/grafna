import { thunkTester } from 'test/core/thunk/thunkTester';

import { setBackendSrv } from '@grafana/runtime';

import { validateDashboardJson } from '../utils/validation';

import { importDashboard } from './actions';
import { DataSourceInput, ImportDashboardDTO, initialImportDashboardState, InputType } from './reducers';

describe('importDashboard', () => {
  it('Should send data source uid', async () => {
    const form: ImportDashboardDTO = {
      title: 'Asda',
      uid: '12',
      gnetId: 'asd',
      constants: [],
      dataSources: [
        {
          id: 1,
          uid: 'ds-uid',
          name: 'ds-name',
          type: 'prometheus',
        } as any,
      ],
      elements: [],
      folder: {
        id: 1,
        title: 'title',
      },
    };

    let postArgs: any;

    setBackendSrv({
      post: (url: string, args: any) => {
        postArgs = args;
        return Promise.resolve({
          importedUrl: '/my/dashboard',
        });
      },
    } as any);

    await thunkTester({
      importDashboard: {
        ...initialImportDashboardState,
        inputs: {
          dataSources: [
            {
              name: 'ds-name',
              pluginId: 'prometheus',
              type: InputType.DataSource,
            },
          ] as DataSourceInput[],
          constants: [],
          libraryPanels: [],
        },
      },
    })
      .givenThunk(importDashboard)
      .whenThunkIsDispatched(form);

    expect(postArgs).toEqual({
      dashboard: {
        title: 'Asda',
        uid: '12',
      },
      folderId: 1,
      inputs: [
        {
          name: 'ds-name',
          pluginId: 'prometheus',
          type: 'datasource',
          value: 'ds-uid',
        },
      ],
      overwrite: true,
    });
  });
});

describe('validateDashboardJson', () => {
  it('Should return true if correct json', async () => {
    const jsonImportCorrectFormat = '{"title": "Correct Format", "tags": ["tag1", "tag2"], "schemaVersion": 36}';
    const validateDashboardJsonCorrectFormat = await validateDashboardJson(jsonImportCorrectFormat);
    expect(validateDashboardJsonCorrectFormat).toBe(true);
  });
  it('Should not return true if nested tags', async () => {
    const jsonImportNestedTags =
      '{"title": "Nested tags","tags": ["tag1", "tag2", ["nestedTag1", "nestedTag2"]],"schemaVersion": 36}';
    const validateDashboardJsonNestedTags = await validateDashboardJson(jsonImportNestedTags);
    expect(validateDashboardJsonNestedTags).toBe('error: tags expected Array of Strings');
  });
  it('Should not return true if not an array', async () => {
    const jsonImportNotArray = '{"title": "Not Array","tags": "tag1","schemaVersion":36}';
    const validateDashboardJsonNotArray = await validateDashboardJson(jsonImportNotArray);
    expect(validateDashboardJsonNotArray).toBe('error: tags expected Array');
  });
});
