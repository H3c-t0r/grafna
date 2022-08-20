import { cloneDeep } from 'lodash';
import { firstValueFrom } from 'rxjs';

import {
  dateTimeFormat,
  TimeRange,
  DataQuery,
  PanelData,
  DataTransformerConfig,
  getValueFormat,
  formattedValueToString,
  DataFrameJSON,
} from '@grafana/data';
import { config } from '@grafana/runtime';

import { PanelModel } from '../dashboard/state';

import { getPanelDataFrames } from './InspectJSONTab';
import { Randomize, randomizeData } from './randomizer';

export async function getTroubleshootingDashboard(panel: PanelModel, rand: Randomize, timeRange: TimeRange) {
  const saveModel = panel.getSaveModel();
  const dashboard = cloneDeep(embeddedDataTemplate);
  const info = {
    panelType: saveModel.type,
    datasource: '??',
  };

  // reproducable
  const data = await firstValueFrom(
    panel.getQueryRunner().getData({
      withFieldConfig: false,
      withTransforms: false,
    })
  );
  const frames = randomizeData(getPanelDataFrames(data), rand);
  const rawFrameContent = JSON.stringify(frames);
  const grafanaVersion = `${config.buildInfo.version} (${config.buildInfo.commit})`;
  const html = `<table width="100%">
    <tr>
      <th width="2%">Panel</th>
      <td >${info.panelType} @ ${saveModel.pluginVersion ?? grafanaVersion}</td>
    </tr>
    <tr>
      <th>Queries&nbsp;(${saveModel.targets?.length ?? 0})</th>
      <td>${saveModel.targets
        .map((t: DataQuery) => {
          return `${t.refId}[${t.datasource?.type}]`;
        })
        .join(', ')}</td>
    </tr>
    ${getTransformsRow(saveModel)}
    ${getDataRow(data, rawFrameContent)}
    ${getAnnotationsRow(data)}
    <tr>
      <th>Grafana</th>
      <td>${grafanaVersion} // ${config.buildInfo.edition}</td>
    </tr>
  </table>`;

  // Replace the panel with embedded data
  dashboard.panels[0] = {
    ...saveModel,
    ...dashboard.panels[0],
    targets: [
      {
        refId: 'A',
        datasource: {
          type: 'testdata',
          uid: 'nVPrVUQGk',
        },
        rawFrameContent,
        scenarioId: 'raw_frame',
      },
    ],
  };

  if (data.annotations?.length) {
    const anno: DataFrameJSON[] = [];
    for (const f of frames) {
      if (f.schema?.meta?.dataTopic) {
        delete f.schema.meta.dataTopic;
        anno.push(f);
      }
    }

    dashboard.panels.push({
      id: 7,
      gridPos: {
        h: 6,
        w: 24,
        x: 0,
        y: 20,
      },
      type: 'table',
      title: 'Annotations',
      datasource: {
        type: 'testdata',
        uid: 'nVPrVUQGk',
      },
      targets: [
        {
          refId: 'A',
          rawFrameContent: JSON.stringify(anno),
          scenarioId: 'raw_frame',
        },
      ],
    });
  }

  dashboard.panels[1].options.content = html;
  dashboard.panels[2].options.content = `<pre>${JSON.stringify(saveModel, null, 2)}</pre>`;
  dashboard.title = `Troubleshooting: ${saveModel.title} // ${dateTimeFormat(new Date())}`;
  dashboard.tags = ['troubleshoot', `troubleshoot-${info.panelType}`];
  dashboard.time = {
    from: timeRange.from.toISOString(),
    to: timeRange.to.toISOString(),
  };

  return dashboard;
}

function getTransformsRow(saveModel: any): string {
  if (!saveModel.transformations) {
    return '';
  }
  return `<tr>
      <th>Transforms (${saveModel.transformations.length})</th>
      <td>${saveModel.transformations.map((t: DataTransformerConfig) => t.id).join(', ')}</td>
  </tr>`;
}

function getDataRow(data: PanelData, raw: string): string {
  const size = getValueFormat('decbytes')(raw?.length);
  return `<tr>
  <th>Data</th>
  <td>${data.state} // Frames: ${data.series?.length} (${formattedValueToString(size)} JSON)</td>
</tr>`;
}

function getAnnotationsRow(data: PanelData): string {
  if (!data.annotations?.length) {
    return '';
  }

  return `<tr>
  <th>Annotations</th>
  <td>${data.annotations.length}</td>
</tr>`;
}

const embeddedDataTemplate: any = {
  // should be dashboard model when that is accurate enough
  panels: [
    {
      id: 2,
      title: 'Reproduced with embedded data',
      datasource: {
        type: 'testdata',
        uid: 'nVPrVUQGk',
      },
      gridPos: {
        h: 13,
        w: 15,
        x: 0,
        y: 0,
      },
    },
    {
      gridPos: {
        h: 7,
        w: 9,
        x: 15,
        y: 0,
      },
      id: 5,
      options: {
        content: 'enter HTLM here',
        mode: 'html',
      },
      title: 'Troubleshooting info',
      type: 'text',
    },
    {
      gridPos: {
        h: 13,
        w: 9,
        x: 15,
        y: 7,
      },
      id: 6,
      options: {
        content: 'enter HTLM here',
        mode: 'html',
      },
      title: 'Original Panel JSON',
      type: 'text',
    },
    {
      id: 3,
      title: 'Data from panel above',
      type: 'table',
      datasource: {
        type: 'datasource',
        uid: '-- Dashboard --',
      },
      gridPos: {
        h: 7,
        w: 15,
        x: 0,
        y: 13,
      },
      targets: [
        {
          datasource: {
            type: 'datasource',
            uid: '-- Dashboard --',
          },
          panelId: 2,
          refId: 'A',
        },
      ],
    },
  ],
  schemaVersion: 37,
};
