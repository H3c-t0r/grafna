import {
  ColumnDefinition,
  CompletionItemKind,
  CompletionItemPriority,
  LanguageCompletionProvider,
  LinkedToken,
  StatementPlacementProvider,
  StatementPosition,
  SuggestionKindProvider,
  TableDefinition,
  TokenType,
} from '@grafana/experimental';
import { PositionContext } from '@grafana/experimental/dist/sql-editor/types';
import { AGGREGATE_FNS, OPERATORS } from 'app/features/plugins/sql/constants';
import { DB, MetaDefinition, SQLQuery } from 'app/features/plugins/sql/types';

interface CompletionProviderGetterArgs {
  getColumns: React.MutableRefObject<(t: SQLQuery) => Promise<ColumnDefinition[]>>;
  getTables: React.MutableRefObject<(d?: string) => Promise<TableDefinition[]>>;
  fetchMeta: React.MutableRefObject<(d?: string) => Promise<MetaDefinition[]>>;
}

export const getSqlCompletionProvider: (args: CompletionProviderGetterArgs) => LanguageCompletionProvider =
  ({ getColumns, getTables, fetchMeta }) =>
  () => ({
    triggerCharacters: ['.', ' ', '$', ',', '(', "'"],
    supportedFunctions: () => AGGREGATE_FNS,
    supportedOperators: () => OPERATORS,
    customSuggestionKinds: customSuggestionKinds(getTables, getColumns, fetchMeta),
    customStatementPlacement,
  });

export enum CustomStatementPlacement {
  AfterDataset = 'afterDataset',
  AfterFrom = 'afterFrom',
}

export enum CustomSuggestionKind {
  TablesWithinDataset = 'tablesWithinDataset',
}

const TRIGGER_SUGGEST = 'editor.action.triggerSuggest';

enum Keyword {
  Where = 'WHERE',
  From = 'FROM',
}

export const customStatementPlacement: StatementPlacementProvider = () => [
  {
    id: CustomStatementPlacement.AfterDataset,
    resolve: (currentToken, previousKeyword) => {
      return Boolean(
        currentToken?.is(TokenType.Delimiter, '.') ||
          (currentToken?.is(TokenType.Whitespace) && currentToken?.previous?.is(TokenType.Delimiter, '.'))
      );
    },
  },
  {
    id: CustomStatementPlacement.AfterFrom,
    resolve: (currentToken, previousKeyword) => {
      return Boolean(isAfterFrom(currentToken));
    },
  },
];

export const customSuggestionKinds: (
  getTables: CompletionProviderGetterArgs['getTables'],
  getFields: CompletionProviderGetterArgs['getColumns'],
  fetchMeta: CompletionProviderGetterArgs['fetchMeta']
) => SuggestionKindProvider = (getTables, _, fetchMeta) => () =>
  [
    {
      id: CustomSuggestionKind.TablesWithinDataset,
      applyTo: [CustomStatementPlacement.AfterDataset],
      suggestionsResolver: async (ctx) => {
        const tablePath = ctx.currentToken ? getTablePath(ctx.currentToken) : '';
        const t = await getTables.current(tablePath);
        return t.map((table) => suggestion(table.name, table.completion ?? table.name, CompletionItemKind.Field, ctx));
      },
    },
    {
      id: `MYSQL${StatementPosition.WhereKeyword}`,
      applyTo: [StatementPosition.WhereKeyword],
      suggestionsResolver: async (ctx) => {
        const path = ctx.currentToken?.value || '';
        const t = await fetchMeta.current(path);
        return t.map((meta) => {
          const completion = meta.kind === CompletionItemKind.Class ? `${meta.completion}.` : meta.completion;
          return suggestion(meta.name, completion!, meta.kind, ctx);
        });
      },
    },
    {
      id: StatementPosition.WhereComparisonOperator,
      applyTo: [StatementPosition.WhereComparisonOperator],
      suggestionsResolver: async (ctx) => {
        if (!isAfterWhere(ctx.currentToken)) {
          return [];
        }
        const path = ctx.currentToken?.value || '';
        const t = await fetchMeta.current(path);
        const sugg = t.map((meta) => {
          const completion = meta.kind === CompletionItemKind.Class ? `${meta.completion}.` : meta.completion;
          return suggestion(meta.name, completion!, meta.kind, ctx);
        });
        return sugg;
      },
    },
    {
      id: 'metaAfterSelect',
      applyTo: [StatementPosition.AfterSelectKeyword],
      suggestionsResolver: async (ctx) => {
        let path = ctx.currentToken?.value || '';
        const fromValue = keywordValue(ctx.currentToken, Keyword.From);
        if (fromValue) {
          path = fromValue;
        }
        const t = await fetchMeta.current(path);
        return t.map((meta) => {
          const completion = meta.kind === CompletionItemKind.Class ? `${meta.completion}.` : meta.completion;
          return suggestion(meta.name, completion!, meta.kind, ctx);
        });
      },
    },
    {
      id: 'metaAfterFrom',
      applyTo: [CustomStatementPlacement.AfterFrom],
      suggestionsResolver: async (ctx) => {
        // TODO: why is this triggering when isAfterFrom is false
        if (!isAfterFrom(ctx.currentToken)) {
          return [];
        }
        const path = ctx.currentToken?.value || '';
        const t = await fetchMeta.current(path);
        return t.map((meta) => suggestion(meta.name, meta.completion!, meta.kind, ctx));
      },
    },
  ];

export function getTablePath(token: LinkedToken) {
  let processedToken = token;
  let tablePath = '';
  while (processedToken?.previous && !processedToken.previous.isWhiteSpace()) {
    processedToken = processedToken.previous;
    tablePath = processedToken.value + tablePath;
  }

  tablePath = tablePath.trim();
  return tablePath;
}

function suggestion(label: string, completion: string, kind: CompletionItemKind, ctx: PositionContext) {
  return {
    label,
    insertText: completion,
    command: { id: TRIGGER_SUGGEST, title: '' },
    kind,
    sortText: CompletionItemPriority.High,
    range: {
      ...ctx.range,
      startColumn: ctx.range.endColumn,
      endColumn: ctx.range.endColumn,
    },
  };
}

function isAfterFrom(token: LinkedToken | null) {
  return isAfter(token, Keyword.From);
}

function isAfterWhere(token: LinkedToken | null) {
  return isAfter(token, Keyword.Where);
}

function isAfter(token: LinkedToken | null, keyword: string) {
  return token?.is(TokenType.Whitespace) && token?.previous?.is(TokenType.Keyword, keyword);
}

function keywordValue(token: LinkedToken | null, keyword: Keyword) {
  let next = token;
  while (next) {
    if (next.is(TokenType.Keyword, keyword)) {
      return tokenValue(next);
    }
    next = next.next;
  }
  return false;
}

function tokenValue(token: LinkedToken | null): string | undefined {
  const ws = token?.next;
  if (ws?.isWhiteSpace()) {
    const v = ws.next;
    const delim = v?.next;
    if (!delim?.is(TokenType.Delimiter)) {
      return v?.value;
    }
    return `${v?.value}${delim?.value}${delim.next?.value}`;
  }
  return undefined;
}

export async function fetchColumns(db: DB, q: SQLQuery) {
  const cols = await db.fields(q);
  if (cols.length > 0) {
    return cols.map((c) => {
      return { name: c.value, type: c.value, description: c.value };
    });
  } else {
    return [];
  }
}

export async function fetchTables(db: DB, q: Partial<SQLQuery>) {
  const tables = await db.lookup(q.dataset);
  return tables;
}
