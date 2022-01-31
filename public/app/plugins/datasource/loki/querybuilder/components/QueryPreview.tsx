import React from 'react';
import { LokiVisualQuery } from '../types';
import { useTheme2 } from '@grafana/ui';
import { GrafanaTheme2 } from '@grafana/data';
import { css, cx } from '@emotion/css';
import { EditorField, EditorFieldGroup } from '@grafana/experimental';
import Prism from 'prismjs';
import { lokiGrammar } from '../../syntax';
import { lokiQueryModeller } from '../LokiQueryModeller';

export interface Props {
  query: LokiVisualQuery;
}

export function QueryPreview({ query }: Props) {
  const theme = useTheme2();
  const styles = getStyles(theme);
  const hightlighted = Prism.highlight(lokiQueryModeller.renderQuery(query), lokiGrammar, 'lokiql');

  return (
    <EditorFieldGroup>
      <EditorField label="Query text">
        <div
          className={cx(styles.editorField, 'prism-syntax-highlight')}
          aria-label="selector"
          dangerouslySetInnerHTML={{ __html: hightlighted }}
        />
      </EditorField>
    </EditorFieldGroup>
  );
}

const getStyles = (theme: GrafanaTheme2) => {
  return {
    editorField: css({
      padding: theme.spacing(0.25, 1),
      fontFamily: theme.typography.fontFamilyMonospace,
      fontSize: theme.typography.bodySmall.fontSize,
    }),
  };
};
