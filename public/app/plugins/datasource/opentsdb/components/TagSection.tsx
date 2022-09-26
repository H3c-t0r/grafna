import { has, size } from 'lodash';
import React, { useCallback, useState } from 'react';

import { SelectableValue, toOption } from '@grafana/data';
import { Select, InlineFormLabel, Icon } from '@grafana/ui';

import { OpenTsdbQuery } from '../types';

export interface TagSectionProps {
  query: OpenTsdbQuery;
  onChange: (query: OpenTsdbQuery) => void;
  onRunQuery: () => void;
  suggestTagKeys: (query: OpenTsdbQuery) => Promise<string[]>;
  suggestTagValues: () => Promise<SelectableValue[]>;
  tsdbVersion: number;
}

export function TagSection({
  query,
  onChange,
  onRunQuery,
  suggestTagKeys,
  suggestTagValues,
  tsdbVersion,
}: TagSectionProps) {
  const [tagKeys, updTagKeys] = useState<Array<SelectableValue<string>>>();
  const [keyIsLoading, updKeyIsLoading] = useState<boolean>();

  const [tagValues, updTagValues] = useState<Array<SelectableValue<string>>>();
  const [valueIsLoading, updValueIsLoading] = useState<boolean>();

  const [addTagMode, updAddTagMode] = useState<boolean>(false);

  const [curTagKey, updCurTagKey] = useState<string | number>('');
  const [curTagValue, updCurTagValue] = useState<string>('');

  const [errors, setErrors] = useState<string>('');

  function changeAddTagMode() {
    updAddTagMode(!addTagMode);
  }

  function addTag() {
    if (query.filters && size(query.filters) > 0) {
      const err = 'Please remove filters to use tags, tags and filters are mutually exclusive.';
      setErrors(err);
      return;
    }

    // ?? not sure why this is here yet
    if (!addTagMode) {
      updAddTagMode(true);
      return;
    }

    // check for duplicate tags
    // used to be in validateTarget()
    if (query.tags && has(query.tags, curTagKey)) {
      const err = "Duplicate tag key '" + curTagKey + "'.";
      setErrors(err);
      return;
    }

    // tags may be undefined
    if (!query.tags) {
      query.tags = {};
    }

    // add tag to query
    query.tags[curTagKey] = curTagValue;

    // reset the inputs
    updCurTagKey('');
    updCurTagValue('');

    // fire the query
    onChange(query);
    onRunQuery();

    // close the tag ditor
    changeAddTagMode();
  }

  function removeTag(key: string | number) {
    delete query.tags[key];

    // fire off the query
    onChange(query);
    onRunQuery();
  }

  function editTag(key: string | number, value: string) {
    removeTag(key);
    updCurTagKey(key);
    updCurTagValue(value);
    addTag();
  }

  // We are matching words split with space
  const splitSeparator = ' ';
  const customTagOption = useCallback((option: SelectableValue<string>, searchQuery: string) => {
    const label = option.value ?? '';

    const searchWords = searchQuery.split(splitSeparator);
    return searchWords.reduce((acc, cur) => acc && label.toLowerCase().includes(cur.toLowerCase()), true);
  }, []);

  return (
    <div className="gf-form-inline" data-testid={testIds.section}>
      <div className="gf-form">
        <InlineFormLabel
          className="query-keyword"
          width={8}
          tooltip={tsdbVersion >= 2 ? <div>Please use filters, tags are deprecated in opentsdb 2.2</div> : undefined}
        >
          Tags
        </InlineFormLabel>
        {query.tags &&
          Object.keys(query.tags).map((tagKey: string | number, idx: number) => {
            const tagValue = query.tags[tagKey];
            return (
              <InlineFormLabel key={idx} width="auto" data-testid={testIds.list + idx}>
                {tagKey}={tagValue}
                <a onClick={() => editTag(tagKey, tagValue)}>
                  <Icon name={'pen'} />
                </a>
                <a onClick={() => removeTag(tagKey)} data-testid={testIds.remove}>
                  <Icon name={'times'} />
                </a>
              </InlineFormLabel>
            );
          })}
        {!addTagMode && (
          <label className="gf-form-label query-keyword">
            <a onClick={changeAddTagMode} data-testid={testIds.open}>
              <Icon name={'plus'} />
            </a>
          </label>
        )}
      </div>
      {addTagMode && (
        <div className="gf-form-inline">
          <div className="gf-form">
            <Select
              inputId="opentsdb-suggested-tagk-select"
              className="gf-form-input"
              value={curTagKey ? toOption('' + curTagKey) : undefined}
              placeholder="key"
              onOpenMenu={async () => {
                updKeyIsLoading(true);
                const tKs = await suggestTagKeys(query);
                const tKsOptions = tKs.map((value: string) => toOption(value));
                updTagKeys(tKsOptions);
                updKeyIsLoading(false);
              }}
              isLoading={keyIsLoading}
              options={tagKeys}
              onChange={({ value }) => {
                if (value) {
                  updCurTagKey(value);
                }
              }}
            />
          </div>

          <div className="gf-form">
            <Select
              inputId="opentsdb-suggested-tagv-select"
              className="gf-form-input"
              value={curTagValue ? toOption(curTagValue) : undefined}
              placeholder="value"
              allowCustomValue
              filterOption={customTagOption}
              onOpenMenu={async () => {
                if (!tagValues) {
                  updValueIsLoading(true);
                  const tVs = await suggestTagValues();
                  updTagValues(tVs);
                  updValueIsLoading(false);
                }
              }}
              isLoading={valueIsLoading}
              options={tagValues}
              onChange={({ value }) => {
                if (value) {
                  updCurTagValue(value);
                }
              }}
            />
          </div>

          <div className="gf-form">
            {errors && (
              <label className="gf-form-label" title={errors} data-testid={testIds.error}>
                <Icon name={'exclamation-triangle'} color={'rgb(229, 189, 28)'} />
              </label>
            )}

            <label className="gf-form-label">
              <a onClick={addTag}>add tag</a>
              <a onClick={changeAddTagMode}>
                <Icon name={'times'} />
              </a>
            </label>
          </div>
        </div>
      )}
    </div>
  );
}

export const testIds = {
  section: 'opentsdb-tag',
  open: 'opentsdb-tag-editor',
  list: 'opentsdb-tag-list',
  error: 'opentsdb-tag-error',
  remove: 'opentsdb-tag-remove',
};
