import { AnyAction } from 'redux';

import { llms } from '@grafana/experimental';
import { PrometheusDatasource } from 'app/plugins/datasource/prometheus/datasource';
import { getMetadataHelp } from 'app/plugins/datasource/prometheus/language_provider';

import { promQueryModeller } from '../../../PromQueryModeller';
import { buildVisualQueryFromString } from '../../../parsing';
import { PromVisualQuery } from '../../../types';
import {
  ExplainSystemPrompt,
  ExplainUserPromptParams,
  GetExplainUserPrompt,
  GetSuggestSystemPrompt,
  SuggestSystemPromptParams,
} from '../prompts';
import { Interaction, QuerySuggestion, SuggestionType } from '../types';

import { createInteraction, stateSlice } from './state';

const OPENAI_MODEL_NAME = 'gpt-3.5-turbo';
const promQLTemplatesCollection = 'promql:templates';
// actions to update the state
const { updateInteraction } = stateSlice.actions;

export const querySuggestions: QuerySuggestion[] = [
  {
    query: 'min(mlapi_http_requests_total{instance="localhost:3005"})',
    explanation: '',
  },
  {
    query: 'sum(mlapi_http_requests_total{instance="localhost:3005"})',
    explanation: '',
  },
  {
    query: 'avg(mlapi_http_requests_total{job="go-app"})',
    explanation: '',
  },
  {
    query: 'mlapi_http_requests_total{job="go-app"}',
    explanation: '',
  },
  {
    query: 'mlapi_http_requests_total{instance="localhost:3005"}',
    explanation: '',
  },
];

interface TemplateSearchResult {
  description: string | null;
  metric_type: string | null;
  promql: string | null;
}

function getExplainMessage({
  documentation,
  metricType,
  description,
  query,
}: ExplainUserPromptParams): llms.openai.Message[] {
  return [
    { role: 'system', content: ExplainSystemPrompt },
    { role: 'user', content: GetExplainUserPrompt({ documentation, metricType, description, query }) },
  ];
}

function getSuggestMessage({ promql, question, labels, templates }: SuggestSystemPromptParams): llms.openai.Message[] {
  return [{ role: 'system', content: GetSuggestSystemPrompt({ promql, question, labels, templates }) }];
}

/**
 * Calls the API and adds suggestions to the interaction
 *
 * @param dispatch
 * @param idx
 * @param interaction
 * @returns
 */
export async function promQailExplain(
  dispatch: React.Dispatch<AnyAction>,
  idx: number,
  query: PromVisualQuery,
  interaction: Interaction,
  suggIdx: number,
  datasource: PrometheusDatasource
) {
  // const enabled = await llms.openai.enabled();

  // if (!enabled) {
  //   return false;
  // }

  const suggestedQuery = interaction.suggestions[suggIdx].query;

  let metricMetadata: string | undefined;

  if (datasource.languageProvider.metricsMetadata) {
    if (interaction.suggestionType === SuggestionType.Historical) {
      // parse the suggested query
      // get the metric
      // then check the metadata
      const pvq = buildVisualQueryFromString(suggestedQuery);
      metricMetadata = getMetadataHelp(pvq.query.metric, datasource.languageProvider.metricsMetadata!);
    } else {
      // for the AI we already have a metric selected
      metricMetadata = getMetadataHelp(query.metric, datasource.languageProvider.metricsMetadata!);
    }
  }

  const promptMessages = getExplainMessage({
    documentation: '',
    metricType: query.metric,
    description: metricMetadata ?? '',
    query: suggestedQuery,
  });

  const interactionToUpdate = interaction;

  return llms.openai
    .streamChatCompletions({
      model: OPENAI_MODEL_NAME,
      messages: promptMessages,
    })
    .pipe(llms.openai.accumulateContent())
    .subscribe((response) => {
      const updatedSuggestions = interactionToUpdate.suggestions.map((sg: QuerySuggestion, sidx: number) => {
        if (suggIdx === sidx) {
          return {
            query: interactionToUpdate.suggestions[suggIdx].query,
            explanation: response,
          };
        }

        return sg;
      });

      const payload = {
        idx,
        interaction: {
          ...interactionToUpdate,
          suggestions: updatedSuggestions,
          explanationIsLoading: false,
        },
      };
      dispatch(updateInteraction(payload));
    });
}

/**
 * Calls the API and adds suggestions to the interaction
 *
 * @param dispatch
 * @param idx
 * @param interaction
 * @returns
 */
export async function promQailSuggest(
  dispatch: React.Dispatch<AnyAction>,
  idx: number,
  query: PromVisualQuery,
  interaction?: Interaction
) {
  // when you're not running promqail
  const check = (await llms.openai.enabled()) && (await llms.vector.enabled());

  if (!check) {
    return new Promise<void>((resolve) => {
      return setTimeout(() => {
        const interactionToUpdate = interaction ? interaction : createInteraction(SuggestionType.Historical);

        const suggestions =
          interactionToUpdate.suggestionType === SuggestionType.Historical ? querySuggestions : [querySuggestions[0]];

        const payload = {
          idx,
          interaction: { ...interactionToUpdate, suggestions: suggestions, isLoading: false },
        };
        dispatch(updateInteraction(payload));
        resolve();
      }, 1000);
    });
  } else {
    type SuggestionBody = {
      metric: string;
      labels: string;
      prompt?: string;
    };

    let feedTheAI: SuggestionBody = {
      metric: query.metric,
      labels: promQueryModeller.renderLabels(query.labels),
    };

    let results: Array<llms.vector.SearchResult<TemplateSearchResult>> = [];
    if (interaction?.suggestionType === SuggestionType.AI) {
      feedTheAI = { ...feedTheAI, prompt: interaction.prompt };

      results = await llms.vector.search<TemplateSearchResult>({
        query: interaction.prompt,
        collection: promQLTemplatesCollection,
        topK: 5,
      });
      // TODO: handle errors from vector search
    }

    const resultsString = JSON.stringify(
      results.map((r) => {
        return {
          metric_type: r.payload.metric_type,
          promql: r.payload.promql,
          description: r.payload.description,
        };
      })
    );

    const interactionToUpdate = interaction ? interaction : createInteraction(SuggestionType.Historical);
    const promptMessages = getSuggestMessage({
      promql: query.metric,
      question: interaction ? interaction.prompt : '',
      labels: promQueryModeller.renderLabels(query.labels),
      templates: resultsString,
    });

    return llms.openai
      .streamChatCompletions({
        model: OPENAI_MODEL_NAME,
        messages: promptMessages,
      })
      .pipe(llms.openai.accumulateContent())
      .subscribe((response) => {
        const payload = {
          idx,
          interaction: {
            ...interactionToUpdate,
            suggestions: [
              {
                query: response,
                explanation: '',
              },
            ],
            isLoading: false,
          },
        };
        dispatch(updateInteraction(payload));
      });
  }
}
