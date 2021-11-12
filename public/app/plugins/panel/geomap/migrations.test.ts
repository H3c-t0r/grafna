import { PanelModel, FieldConfigSource } from '@grafana/data';
import { mapMigrationHandler, mapPanelChangedHandler } from './migrations';
describe('Worldmap Migrations', () => {
  let prevFieldConfig: FieldConfigSource;

  beforeEach(() => {
    prevFieldConfig = {
      defaults: {},
      overrides: [],
    };
  });

  it('simple worldmap', () => {
    const old: any = {
      angular: simpleWorldmapConfig,
    };
    const panel = {} as PanelModel;
    panel.options = mapPanelChangedHandler(panel, 'grafana-worldmap-panel', old, prevFieldConfig);
    expect(panel).toMatchInlineSnapshot(`
      Object {
        "fieldConfig": Object {
          "defaults": Object {
            "decimals": 3,
            "thresholds": Object {
              "mode": "absolute",
              "steps": Array [
                Object {
                  "color": "#37872D",
                  "value": -Infinity,
                },
                Object {
                  "color": "#E0B400",
                  "value": 0,
                },
                Object {
                  "color": "#C4162A",
                  "value": 50,
                },
                Object {
                  "color": "#8F3BB8",
                  "value": 100,
                },
              ],
            },
          },
          "overrides": Array [],
        },
        "options": Object {
          "basemap": Object {
            "type": "default",
          },
          "controls": Object {
            "mouseWheelZoom": true,
            "showZoom": true,
          },
          "layers": Array [],
          "view": Object {
            "id": "europe",
            "lat": 46,
            "lon": 14,
            "zoom": 6,
          },
        },
      }
    `);
  });
});

const simpleWorldmapConfig = {
  id: 23763571993,
  gridPos: {
    h: 8,
    w: 12,
    x: 0,
    y: 0,
  },
  type: 'grafana-worldmap-panel',
  title: 'Panel Title',
  thresholds: '0,50,100',
  maxDataPoints: 1,
  circleMaxSize: 30,
  circleMinSize: 2,
  colors: ['#37872D', '#E0B400', '#C4162A', '#8F3BB8'],
  decimals: 3,
  esMetric: 'Count',
  hideEmpty: false,
  hideZero: false,
  initialZoom: '6',
  locationData: 'countries',
  mapCenter: 'Europe',
  mapCenterLatitude: 46,
  mapCenterLongitude: 14,
  mouseWheelZoom: true,
  showLegend: true,
  stickyLabels: false,
  tableQueryOptions: {
    geohashField: 'geohash',
    latitudeField: 'latitude',
    longitudeField: 'longitude',
    metricField: 'metric',
    queryType: 'geohash',
  },
  unitPlural: '',
  unitSingle: '',
  valueName: 'total',
  datasource: null,
};

describe('geomap migrations', () => {
  it('updates marker', () => {
    const panel = ({
      type: 'geomap',
      options: {
        layers: [
          {
            config: {
              shape: 'circle', // should become resource path
            },
            type: 'markers',
          },
        ],
      },
      pluginVersion: '8.2.0',
    } as any) as PanelModel;
    panel.options = mapMigrationHandler(panel);

    expect(panel).toMatchInlineSnapshot(`
      Object {
        "options": Object {
          "layers": Array [
            Object {
              "config": Object {
                "markerSymbol": Object {
                  "fixed": "img/icons/marker/circle.svg",
                  "mode": "fixed",
                },
              },
              "type": "markers",
            },
          ],
        },
        "pluginVersion": "8.2.0",
        "type": "geomap",
      }
    `);
  });
});
