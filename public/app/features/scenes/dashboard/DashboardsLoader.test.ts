import { defaultDashboard, LoadingState, VariableType } from '@grafana/schema';
import { DashboardLoaderSrv, setDashboardLoaderSrv } from 'app/features/dashboard/services/DashboardLoaderSrv';
import { DashboardModel, PanelModel } from 'app/features/dashboard/state';

import { SceneGridLayout } from '../components';
import { SceneQueryRunner } from '../querying/SceneQueryRunner';
import { CustomVariable } from '../variables/variants/CustomVariable';
import { DataSourceVariable } from '../variables/variants/DataSourceVariable';
import { QueryVariable } from '../variables/variants/query/QueryVariable';

import { DashboardScene } from './DashboardScene';
import {
  createDashboardSceneFromDashboardModel,
  createSceneVariableFromVariableModel,
  DashboardLoader,
  createVizPanelFromPanelModel,
} from './DashboardsLoader';

describe('DashboardLoader', () => {
  describe('when fetching/loading a dashboard', () => {
    beforeEach(() => {
      new DashboardLoader({});
    });

    it('should load the dashboard from the cache if it exists', () => {
      const loader = new DashboardLoader({});
      const dashboard = new DashboardScene({
        title: 'cached',
        uid: 'fake-uid',
        body: new SceneGridLayout({ children: [] }),
      });
      // @ts-expect-error
      loader.cache['fake-uid'] = dashboard;
      loader.load('fake-uid');
      expect(loader.state.dashboard).toBe(dashboard);
    });

    it('should call dashboard loader server if the dashboard is not cached', async () => {
      const loadDashboardMock = jest.fn().mockResolvedValue({ dashboard: { uid: 'fake-dash' }, meta: {} });
      setDashboardLoaderSrv({
        loadDashboard: loadDashboardMock,
      } as unknown as DashboardLoaderSrv);

      const loader = new DashboardLoader({});
      await loader.load('fake-dash');

      expect(loadDashboardMock).toHaveBeenCalledWith('db', '', 'fake-dash');
    });

    it("should error when the dashboard doesn't exist", async () => {
      const loadDashboardMock = jest.fn().mockResolvedValue({ dashboard: undefined, meta: undefined });
      setDashboardLoaderSrv({
        loadDashboard: loadDashboardMock,
      } as unknown as DashboardLoaderSrv);

      const loader = new DashboardLoader({});
      await loader.load('fake-dash');

      expect(loader.state.dashboard).toBeUndefined();
      // @ts-expect-error
      expect(loader.cache['fake-dash']).toBeUndefined();
      expect(loader.state.loadError).toBe('Error: Dashboard not found');
    });

    it('should initialize the dashboard scene with the loaded dashboard', async () => {
      const loadDashboardMock = jest.fn().mockResolvedValue({ dashboard: { uid: 'fake-dash' }, meta: {} });
      setDashboardLoaderSrv({
        loadDashboard: loadDashboardMock,
      } as unknown as DashboardLoaderSrv);

      const loader = new DashboardLoader({});
      await loader.load('fake-dash');

      expect(loader.state.dashboard?.state.uid).toBe('fake-dash');
      // @ts-expect-error
      expect(loader.cache['fake-dash']).toBeDefined();
      expect(loader.state.loadError).toBe(undefined);
    });

    it('should use DashboardScene creator to create the scene', async () => {
      const loadDashboardMock = jest.fn().mockResolvedValue({ dashboard: { uid: 'fake-dash' }, meta: {} });
      setDashboardLoaderSrv({
        loadDashboard: loadDashboardMock,
      } as unknown as DashboardLoaderSrv);

      const loader = new DashboardLoader({});
      await loader.load('fake-dash');
      // TODO: Mock the createDashboardSceneFromDashboardModel
      //expect(createDashboardSceneFromDashboardModel).toBeCalledWith('fake-dash');
    });
  });

  describe('when creating dashboard scene', () => {
    it('should initialize the DashboardScene with the model state', () => {
      const dash = {
        ...defaultDashboard,
        title: 'test',
        uid: 'test-uid',
        time: { from: 'now-10h', to: 'now' },
        templating: {
          list: [
            {
              hide: 2,
              name: 'constant',
              skipUrlSync: false,
              type: 'constant' as VariableType,
              rootStateKey: 'N4XLmH5Vz',
              query: 'test',
              id: 'constant',
              global: false,
              index: 3,
              state: LoadingState.Done,
              error: null,
              description: '',
              datasource: null,
            },
          ],
        },
      };
      const oldModel = new DashboardModel(dash);

      const scene = createDashboardSceneFromDashboardModel(oldModel);

      expect(scene.state.title).toBe('test');
      expect(scene.state.uid).toBe('test-uid');
      expect(scene.state?.$timeRange?.state.value.raw).toEqual(dash.time);
      expect(scene.state?.$variables?.state.variables).toHaveLength(1);
      expect(scene.state.subMenu).toBeDefined();
    });

    it.todo('should call variable migrator for each variable');
  });

  describe('when organizing panels as scene children', () => {
    it.todo('should not create panels within collapsed rows');
    it.todo('should nest panels within their row');
  });

  describe('when creating viz panel objects', () => {
    it('should initalize properly the VizPanel scene object', () => {
      const panel = {
        title: 'test',
        type: 'test-plugin',
        gridPos: { x: 0, y: 0, w: 12, h: 8 },
        options: {
          fieldOptions: {
            defaults: {
              unit: 'none',
              decimals: 2,
            },
            overrides: [],
          },
        },
        fieldConfig: {
          defaults: {
            unit: 'none',
          },
        },
        pluginVersion: '1.0.0',
        transformations: [
          {
            id: 'reduce',
            options: {
              reducers: [
                {
                  id: 'mean',
                },
              ],
            },
          },
        ],
        targets: [
          {
            refId: 'A',
            queryType: 'randomWalk',
          },
        ],
      };
      const vizPanelSceneObject = createVizPanelFromPanelModel(new PanelModel(panel));

      expect(vizPanelSceneObject.state.title).toBe('test');
      expect(vizPanelSceneObject.state.pluginId).toBe('test-plugin');
      expect(vizPanelSceneObject.state.placement).toEqual({ x: 0, y: 0, width: 12, height: 8 });
      expect(vizPanelSceneObject.state.options).toEqual(panel.options);
      expect(vizPanelSceneObject.state.fieldConfig).toEqual(panel.fieldConfig);
      expect(vizPanelSceneObject.state.pluginVersion).toBe('1.0.0');
      expect((vizPanelSceneObject.state.$data as SceneQueryRunner)?.state.queries).toEqual(panel.targets);
      expect((vizPanelSceneObject.state.$data as SceneQueryRunner)?.state.transformations).toEqual(
        panel.transformations
      );
    });
  });

  describe('when creating variables objects', () => {
    it('should migrate custom variable', () => {
      const variable = {
        current: {
          selected: false,
          text: 'a',
          value: 'a',
        },
        hide: 0,
        includeAll: false,
        multi: false,
        name: 'query0',
        options: [
          {
            selected: true,
            text: 'a',
            value: 'a',
          },
          {
            selected: false,
            text: 'b',
            value: 'b',
          },
          {
            selected: false,
            text: 'c',
            value: 'c',
          },
          {
            selected: false,
            text: 'd',
            value: 'd',
          },
        ],
        query: 'a,b,c,d',
        skipUrlSync: false,
        type: 'custom' as VariableType,
        rootStateKey: 'N4XLmH5Vz',
        id: 'query0',
        global: false,
        index: 0,
        state: 'Done',
        error: null,
        description: null,
        allValue: null,
      };

      const migrated = createSceneVariableFromVariableModel(variable);
      const { key, ...rest } = migrated.state;

      expect(migrated).toBeInstanceOf(CustomVariable);
      expect(rest).toEqual({
        allValue: undefined,
        defaultToAll: false,
        description: null,
        includeAll: false,
        isMulti: false,
        label: undefined,
        name: 'query0',
        options: [],
        query: 'a,b,c,d',
        skipUrlSync: false,
        text: 'a',
        type: 'custom',
        value: 'a',
        hide: 0,
      });
    });
    it('should migrate query variable', () => {
      const variable = {
        allValue: null,
        current: {
          text: 'America',
          value: 'America',
          selected: false,
        },
        datasource: {
          uid: 'P15396BDD62B2BE29',
          type: 'influxdb',
        },
        definition: '',
        hide: 0,
        includeAll: false,
        label: 'Datacenter',
        multi: false,
        name: 'datacenter',
        options: [
          {
            text: 'America',
            value: 'America',
            selected: true,
          },
          {
            text: 'Africa',
            value: 'Africa',
            selected: false,
          },
          {
            text: 'Asia',
            value: 'Asia',
            selected: false,
          },
          {
            text: 'Europe',
            value: 'Europe',
            selected: false,
          },
        ],
        query: 'SHOW TAG VALUES  WITH KEY = "datacenter" ',
        refresh: 1,
        regex: '',
        skipUrlSync: false,
        sort: 0,
        tagValuesQuery: null,
        tagsQuery: null,
        type: 'query' as VariableType,
        useTags: false,
        rootStateKey: '000000002',
        id: 'datacenter',
        global: false,
        index: 0,
        state: 'Done',
        error: null,
        description: null,
      };

      const migrated = createSceneVariableFromVariableModel(variable);
      const { key, ...rest } = migrated.state;

      expect(migrated).toBeInstanceOf(QueryVariable);
      expect(rest).toEqual({
        allValue: undefined,
        datasource: {
          type: 'influxdb',
          uid: 'P15396BDD62B2BE29',
        },
        defaultToAll: false,
        description: null,
        includeAll: false,
        isMulti: false,
        label: 'Datacenter',
        name: 'datacenter',
        options: [],
        query: 'SHOW TAG VALUES  WITH KEY = "datacenter" ',
        refresh: 1,
        regex: '',
        skipUrlSync: false,
        sort: 0,
        text: 'America',
        type: 'query',
        value: 'America',
        hide: 0,
      });
    });

    it('should migrate datasource variable', () => {
      const variable = {
        id: 'query1',
        rootStateKey: 'N4XLmH5Vz',
        name: 'query1',
        type: 'datasource' as VariableType,
        global: false,
        index: 1,
        hide: 0,
        skipUrlSync: false,
        state: 'Done',
        error: null,
        description: null,
        current: {
          value: ['gdev-prometheus', 'gdev-slow-prometheus'],
          text: ['gdev-prometheus', 'gdev-slow-prometheus'],
          selected: true,
        },
        regex: '/^gdev/',
        options: [
          {
            text: 'All',
            value: '$__all',
            selected: false,
          },
          {
            text: 'gdev-prometheus',
            value: 'gdev-prometheus',
            selected: true,
          },
          {
            text: 'gdev-slow-prometheus',
            value: 'gdev-slow-prometheus',
            selected: false,
          },
        ],
        query: 'prometheus',
        multi: true,
        includeAll: true,
        refresh: 1,
        allValue: 'Custom all',
      };

      const migrated = createSceneVariableFromVariableModel(variable);
      const { key, ...rest } = migrated.state;

      expect(migrated).toBeInstanceOf(DataSourceVariable);
      expect(rest).toEqual({
        allValue: 'Custom all',
        defaultToAll: true,
        includeAll: true,
        label: undefined,
        name: 'query1',
        options: [],
        query: 'prometheus',
        regex: '/^gdev/',
        skipUrlSync: false,
        text: ['gdev-prometheus', 'gdev-slow-prometheus'],
        type: 'datasource',
        value: ['gdev-prometheus', 'gdev-slow-prometheus'],
        isMulti: true,
        description: null,
        hide: 0,
      });
    });

    it('should migrate constant variable', () => {
      const variable = {
        hide: 2,
        label: 'constant',
        name: 'constant',
        skipUrlSync: false,
        type: 'constant' as VariableType,
        rootStateKey: 'N4XLmH5Vz',
        current: {
          selected: true,
          text: 'test',
          value: 'test',
        },
        options: [
          {
            selected: true,
            text: 'test',
            value: 'test',
          },
        ],
        query: 'test',
        id: 'constant',
        global: false,
        index: 3,
        state: 'Done',
        error: null,
        description: null,
      };

      const migrated = createSceneVariableFromVariableModel(variable);
      const { key, ...rest } = migrated.state;

      expect(rest).toEqual({
        description: null,
        hide: 2,
        label: 'constant',
        name: 'constant',
        skipUrlSync: false,
        type: 'constant',
        value: 'test',
      });
    });

    it.each(['adhoc', 'interval', 'textbox', 'system'])('should throw for unsupported (yet) variables', (type) => {
      const variable = {
        name: 'query0',
        type: type as VariableType,
      };

      expect(() => createSceneVariableFromVariableModel(variable)).toThrow();
    });
  });
});
