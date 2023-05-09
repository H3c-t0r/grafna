import { DataSourceInstanceSettings, PluginMetaInfo, PluginType } from '@grafana/data';
import { monacoTypes } from '@grafana/ui';

import { TempoDatasource } from '../datasource';
import TempoLanguageProvider from '../language_provider';
import { TempoJsonData, Tags } from '../types';

import { CompletionProvider } from './autocomplete';
import { intrinsics, scopes } from './traceql';

jest.mock('@grafana/runtime', () => ({
  ...jest.requireActual('@grafana/runtime'),
}));

describe('CompletionProvider', () => {
  it('suggests tags, intrinsics and scopes (API v1)', async () => {
    const { provider, model } = setup('{}', 1, defaultTags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual([
      ...scopes.map((s) => expect.objectContaining({ label: s, insertText: s })),
      ...intrinsics.map((s) => expect.objectContaining({ label: s, insertText: s })),
      expect.objectContaining({ label: 'bar', insertText: '.bar' }),
      expect.objectContaining({ label: 'foo', insertText: '.foo' }),
    ]);
  });

  it('suggests tags, intrinsics and scopes (API v2)', async () => {
    const { provider, model } = setup('{}', 1, v2Tags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual([
      ...scopes.map((s) => expect.objectContaining({ label: s, insertText: s })),
      ...intrinsics.map((s) => expect.objectContaining({ label: s, insertText: s })),
      expect.objectContaining({ label: 'cluster', insertText: '.cluster' }),
      expect.objectContaining({ label: 'container', insertText: '.container' }),
      expect.objectContaining({ label: 'db', insertText: '.db' }),
    ]);
  });

  it('does not wrap the tag value in quotes if the type in the response is something other than "string"', async () => {
    const { provider, model } = setup('{foo=}', 5, defaultTags);

    jest.spyOn(provider.languageProvider, 'getOptionsV2').mockImplementation(
      () =>
        new Promise((resolve) => {
          resolve([
            {
              type: 'int',
              value: 'foobar',
              label: 'foobar',
            },
          ]);
        })
    );

    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual([
      expect.objectContaining({ label: 'foobar', insertText: 'foobar' }),
    ]);
  });

  it('wraps the tag value in quotes if the type in the response is set to "string"', async () => {
    const { provider, model } = setup('{foo=}', 5, defaultTags);

    jest.spyOn(provider.languageProvider, 'getOptionsV2').mockImplementation(
      () =>
        new Promise((resolve) => {
          resolve([
            {
              type: 'string',
              value: 'foobar',
              label: 'foobar',
            },
          ]);
        })
    );

    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual([
      expect.objectContaining({ label: 'foobar', insertText: '"foobar"' }),
    ]);
  });

  it('inserts the tag value without quotes if the user has entered quotes', async () => {
    const { provider, model } = setup('{foo="}', 6, defaultTags);

    jest.spyOn(provider.languageProvider, 'getOptionsV2').mockImplementation(
      () =>
        new Promise((resolve) => {
          resolve([
            {
              value: 'foobar',
              label: 'foobar',
            },
          ]);
        })
    );

    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual([
      expect.objectContaining({ label: 'foobar', insertText: 'foobar' }),
    ]);
  });

  it('suggests nothing without tags', async () => {
    const { provider, model } = setup('{foo="}', 7, emptyTags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual([]);
  });

  it('suggests tags on empty input (API v1)', async () => {
    const { provider, model } = setup('', 0, defaultTags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual([
      ...scopes.map((s) => expect.objectContaining({ label: s, insertText: `{ ${s}` })),
      ...intrinsics.map((s) => expect.objectContaining({ label: s, insertText: `{ ${s}` })),
      expect.objectContaining({ label: 'bar', insertText: '{ .bar' }),
      expect.objectContaining({ label: 'foo', insertText: '{ .foo' }),
    ]);
  });

  it('suggests tags on empty input (API v2)', async () => {
    const { provider, model } = setup('', 0, v2Tags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual([
      ...scopes.map((s) => expect.objectContaining({ label: s, insertText: `{ ${s}` })),
      ...intrinsics.map((s) => expect.objectContaining({ label: s, insertText: `{ ${s}` })),
      expect.objectContaining({ label: 'cluster', insertText: '{ .cluster' }),
      expect.objectContaining({ label: 'container', insertText: '{ .container' }),
      expect.objectContaining({ label: 'db', insertText: '{ .db' }),
    ]);
  });

  it('only suggests tags after typing the global attribute scope (API v1)', async () => {
    const { provider, model } = setup('{.}', 2, defaultTags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual(
      defaultTags.v1.map((s) => expect.objectContaining({ label: s, insertText: s }))
    );
  });

  it('only suggests tags after typing the global attribute scope (API v2)', async () => {
    const { provider, model } = setup('{.}', 2, v2Tags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual(
      ['cluster', 'container', 'db'].map((s) => expect.objectContaining({ label: s, insertText: s }))
    );
  });

  it('suggests operators after a space after the tag name', async () => {
    const { provider, model } = setup('{ foo }', 6, defaultTags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual(
      CompletionProvider.operators.map((s) => expect.objectContaining({ label: s, insertText: s }))
    );
  });

  it('suggests tags after a scope (API v1)', async () => {
    const { provider, model } = setup('{ resource. }', 11, defaultTags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual(
      defaultTags.v1.map((s) => expect.objectContaining({ label: s, insertText: s }))
    );
  });

  it('suggests correct tags after the resource scope (API v2)', async () => {
    const { provider, model } = setup('{ resource. }', 11, v2Tags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual(
      ['cluster', 'container'].map((s) => expect.objectContaining({ label: s, insertText: s }))
    );
  });

  it('suggests correct tags after the span scope (API v2)', async () => {
    const { provider, model } = setup('{ span. }', 7, v2Tags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual(
      ['db'].map((s) => expect.objectContaining({ label: s, insertText: s }))
    );
  });

  it('suggests logical operators and close bracket after the value', async () => {
    const { provider, model } = setup('{foo=300 }', 9, defaultTags);
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual([
      ...CompletionProvider.logicalOps.map((s) => expect.objectContaining({ label: s, insertText: s })),
      expect.objectContaining({ label: '}', insertText: '}' }),
    ]);
  });

  it('suggests tag values after a space inside a string', async () => {
    const { provider, model } = setup('{foo="bar test " }', 15, defaultTags);

    jest.spyOn(provider.languageProvider, 'getOptionsV2').mockImplementation(
      () =>
        new Promise((resolve) => {
          resolve([
            {
              value: 'foobar',
              label: 'foobar',
            },
          ]);
        })
    );
    const result = await provider.provideCompletionItems(
      model as unknown as monacoTypes.editor.ITextModel,
      {} as monacoTypes.Position
    );
    expect((result! as monacoTypes.languages.CompletionList).suggestions).toEqual([
      expect.objectContaining({ label: 'foobar', insertText: 'foobar' }),
    ]);
  });
});

export const emptyTags = {
  v1: [],
  v2: [],
};

export const defaultTags = {
  v1: ['bar', 'foo'],
  v2: undefined,
};

export const v2Tags = {
  v1: undefined,
  v2: [
    {
      name: 'resource',
      tags: ['cluster', 'container'],
    },
    {
      name: 'span',
      tags: ['db'],
    },
    {
      name: 'intrinsic',
      tags: ['duration', 'kind', 'name', 'status'],
    },
  ],
};

function setup(value: string, offset: number, tags?: Tags) {
  const ds = new TempoDatasource(defaultSettings);
  const provider = new CompletionProvider({ languageProvider: new TempoLanguageProvider(ds) });
  if (tags) {
    provider.setTags(tags);
  }
  const model = makeModel(value, offset);
  provider.monaco = {
    Range: {
      fromPositions() {
        return null;
      },
    },
    languages: {
      CompletionItemKind: {
        Enum: 1,
        EnumMember: 2,
      },
    },
  } as any;
  provider.editor = {
    getModel() {
      return model;
    },
  } as any;

  return { provider, model };
}

function makeModel(value: string, offset: number) {
  return {
    id: 'test_monaco',
    getWordAtPosition() {
      return null;
    },
    getOffsetAt() {
      return offset;
    },
    getValue() {
      return value;
    },
  };
}

const defaultSettings: DataSourceInstanceSettings<TempoJsonData> = {
  id: 0,
  uid: 'gdev-tempo',
  type: 'tracing',
  name: 'tempo',
  access: 'proxy',
  meta: {
    id: 'tempo',
    name: 'tempo',
    type: PluginType.datasource,
    info: {} as PluginMetaInfo,
    module: '',
    baseUrl: '',
  },
  readOnly: false,
  jsonData: {
    nodeGraph: {
      enabled: true,
    },
  },
};
