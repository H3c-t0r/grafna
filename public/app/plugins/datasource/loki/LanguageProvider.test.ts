import { AbstractLabelOperator, DataFrame } from '@grafana/data';

import LanguageProvider from './LanguageProvider';
import { LokiDatasource } from './datasource';
import { createLokiDatasource, createMetadataRequest } from './mocks';
import {
  extractLogParserFromDataFrame,
  extractLabelKeysFromDataFrame,
  extractUnwrapLabelKeysFromDataFrame,
} from './responseUtils';
import { LokiQueryType } from './types';

jest.mock('./responseUtils');

jest.mock('app/store/store', () => ({
  store: {
    getState: jest.fn().mockReturnValue({
      explore: {
        left: {
          mode: 'Logs',
        },
      },
    }),
  },
}));

describe('Language completion provider', () => {
  describe('fetchSeries', () => {
    it('should use match[] parameter', () => {
      const datasource = setup({}, { '{foo="bar"}': [{ label1: 'label_val1' }] });
      const languageProvider = new LanguageProvider(datasource);
      const fetchSeries = languageProvider.fetchSeries;
      const requestSpy = jest.spyOn(languageProvider, 'request');
      fetchSeries('{job="grafana"}');
      expect(requestSpy).toHaveBeenCalledWith('series', {
        end: 1560163909000,
        'match[]': '{job="grafana"}',
        start: 1560153109000,
      });
    });
  });

  describe('fetchSeriesLabels', () => {
    it('should interpolate variable in series', () => {
      const datasource = setup({});
      jest.spyOn(datasource, 'getTimeRangeParams').mockReturnValue({ start: 0, end: 1 });
      jest
        .spyOn(datasource, 'interpolateString')
        .mockImplementation((string: string) => string.replace(/\$/, 'interpolated-'));

      const languageProvider = new LanguageProvider(datasource);
      const fetchSeriesLabels = languageProvider.fetchSeriesLabels;
      const requestSpy = jest.spyOn(languageProvider, 'request').mockResolvedValue([]);
      fetchSeriesLabels('$stream');
      expect(requestSpy).toHaveBeenCalled();
      expect(requestSpy).toHaveBeenCalledWith('series', {
        end: 1,
        'match[]': 'interpolated-stream',
        start: 0,
      });
    });
  });

  describe('label values', () => {
    it('should fetch label values if not cached', async () => {
      const datasource = setup({ testkey: ['label1_val1', 'label1_val2'], label2: [] });
      const provider = await getLanguageProvider(datasource);
      const requestSpy = jest.spyOn(provider, 'request');
      const labelValues = await provider.fetchLabelValues('testkey');
      expect(requestSpy).toHaveBeenCalled();
      expect(labelValues).toEqual(['label1_val1', 'label1_val2']);
    });

    it('should return cached values', async () => {
      const datasource = setup({ testkey: ['label1_val1', 'label1_val2'], label2: [] });
      const provider = await getLanguageProvider(datasource);
      const requestSpy = jest.spyOn(provider, 'request');
      const labelValues = await provider.fetchLabelValues('testkey');
      expect(requestSpy).toHaveBeenCalledTimes(1);
      expect(labelValues).toEqual(['label1_val1', 'label1_val2']);

      const nextLabelValues = await provider.fetchLabelValues('testkey');
      expect(requestSpy).toHaveBeenCalledTimes(1);
      expect(nextLabelValues).toEqual(['label1_val1', 'label1_val2']);
    });

    it('should encode special characters', async () => {
      const datasource = setup({ '`\\"testkey': ['label1_val1', 'label1_val2'], label2: [] });
      const provider = await getLanguageProvider(datasource);
      const requestSpy = jest.spyOn(provider, 'request');
      await provider.fetchLabelValues('`\\"testkey');

      expect(requestSpy).toHaveBeenCalledWith('label/%60%5C%22testkey/values', expect.any(Object));
    });
  });
});

describe('Request URL', () => {
  it('should contain range params', async () => {
    const datasourceWithLabels = setup({ other: [] });
    const rangeParams = datasourceWithLabels.getTimeRangeParams();
    const datasourceSpy = jest.spyOn(datasourceWithLabels, 'metadataRequest');

    const instance = new LanguageProvider(datasourceWithLabels);
    instance.fetchLabels();
    const expectedUrl = 'labels';
    expect(datasourceSpy).toHaveBeenCalledWith(expectedUrl, rangeParams);
  });
});

describe('fetchLabels', () => {
  it('should return labels', async () => {
    const datasourceWithLabels = setup({ other: [] });

    const instance = new LanguageProvider(datasourceWithLabels);
    const labels = await instance.fetchLabels();
    expect(labels).toEqual(['other']);
  });

  it('should set labels', async () => {
    const datasourceWithLabels = setup({ other: [] });

    const instance = new LanguageProvider(datasourceWithLabels);
    await instance.fetchLabels();
    expect(instance.labelKeys).toEqual(['other']);
  });

  it('should return empty array', async () => {
    const datasourceWithLabels = setup({});

    const instance = new LanguageProvider(datasourceWithLabels);
    const labels = await instance.fetchLabels();
    expect(labels).toEqual([]);
  });

  it('should set empty array', async () => {
    const datasourceWithLabels = setup({});

    const instance = new LanguageProvider(datasourceWithLabels);
    await instance.fetchLabels();
    expect(instance.labelKeys).toEqual([]);
  });
});

describe('Query imports', () => {
  const datasource = setup({});

  it('returns empty queries', async () => {
    const instance = new LanguageProvider(datasource);
    const result = await instance.importFromAbstractQuery({ refId: 'bar', labelMatchers: [] });
    expect(result).toEqual({ refId: 'bar', expr: '', queryType: LokiQueryType.Range });
  });

  describe('exporting to abstract query', () => {
    it('exports labels', async () => {
      const instance = new LanguageProvider(datasource);
      const abstractQuery = instance.exportToAbstractQuery({
        refId: 'bar',
        expr: '{label1="value1", label2!="value2", label3=~"value3", label4!~"value4"}',
        instant: true,
        range: false,
      });
      expect(abstractQuery).toMatchObject({
        refId: 'bar',
        labelMatchers: [
          { name: 'label1', operator: AbstractLabelOperator.Equal, value: 'value1' },
          { name: 'label2', operator: AbstractLabelOperator.NotEqual, value: 'value2' },
          { name: 'label3', operator: AbstractLabelOperator.EqualRegEx, value: 'value3' },
          { name: 'label4', operator: AbstractLabelOperator.NotEqualRegEx, value: 'value4' },
        ],
      });
    });
  });

  describe('getParserAndLabelKeys()', () => {
    let datasource: LokiDatasource, languageProvider: LanguageProvider;
    const extractLogParserFromDataFrameMock = jest.mocked(extractLogParserFromDataFrame);
    const extractedLabelKeys = ['extracted', 'label'];
    const unwrapLabelKeys = ['unwrap', 'labels'];

    beforeEach(() => {
      datasource = createLokiDatasource();
      languageProvider = new LanguageProvider(datasource);
      jest.mocked(extractLabelKeysFromDataFrame).mockReturnValue(extractedLabelKeys);
      jest.mocked(extractUnwrapLabelKeysFromDataFrame).mockReturnValue(unwrapLabelKeys);
    });

    it('identifies selectors with JSON parser data', async () => {
      jest.spyOn(datasource, 'getDataSamples').mockResolvedValue([{}] as DataFrame[]);
      extractLogParserFromDataFrameMock.mockReturnValueOnce({ hasLogfmt: false, hasJSON: true, hasPack: false });

      expect(await languageProvider.getParserAndLabelKeys('{place="luna"}')).toEqual({
        extractedLabelKeys,
        unwrapLabelKeys,
        hasJSON: true,
        hasLogfmt: false,
        hasPack: false,
      });
    });

    it('identifies selectors with Logfmt parser data', async () => {
      jest.spyOn(datasource, 'getDataSamples').mockResolvedValue([{}] as DataFrame[]);
      extractLogParserFromDataFrameMock.mockReturnValueOnce({ hasLogfmt: true, hasJSON: false, hasPack: false });

      expect(await languageProvider.getParserAndLabelKeys('{place="luna"}')).toEqual({
        extractedLabelKeys,
        unwrapLabelKeys,
        hasJSON: false,
        hasLogfmt: true,
        hasPack: false,
      });
    });

    it('correctly processes empty data', async () => {
      jest.spyOn(datasource, 'getDataSamples').mockResolvedValue([]);
      extractLogParserFromDataFrameMock.mockClear();

      expect(await languageProvider.getParserAndLabelKeys('{place="luna"}')).toEqual({
        extractedLabelKeys: [],
        unwrapLabelKeys: [],
        hasJSON: false,
        hasLogfmt: false,
        hasPack: false,
      });
      expect(extractLogParserFromDataFrameMock).not.toHaveBeenCalled();
    });
  });
});

async function getLanguageProvider(datasource: LokiDatasource) {
  const instance = new LanguageProvider(datasource);
  await instance.start();
  return instance;
}

function setup(
  labelsAndValues: Record<string, string[]>,
  series?: Record<string, Array<Record<string, string>>>
): LokiDatasource {
  const datasource = createLokiDatasource();

  const rangeMock = {
    start: 1560153109000,
    end: 1560163909000,
  };

  jest.spyOn(datasource, 'getTimeRangeParams').mockReturnValue(rangeMock);
  jest.spyOn(datasource, 'metadataRequest').mockImplementation(createMetadataRequest(labelsAndValues, series));
  jest.spyOn(datasource, 'interpolateString').mockImplementation((string: string) => string);

  return datasource;
}
