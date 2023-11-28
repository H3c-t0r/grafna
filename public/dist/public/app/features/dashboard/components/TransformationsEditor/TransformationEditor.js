import { css } from '@emotion/css';
import React, { useEffect, useMemo, useState } from 'react';
import { mergeMap } from 'rxjs/operators';
import { transformDataFrame, getFrameMatchers, } from '@grafana/data';
import { selectors } from '@grafana/e2e-selectors';
import { getTemplateSrv } from '@grafana/runtime';
import { Icon, JSONFormatter, useStyles2 } from '@grafana/ui';
export const TransformationEditor = ({ debugMode, index, data, uiConfig, configs, onChange, }) => {
    const styles = useStyles2(getStyles);
    const [input, setInput] = useState([]);
    const [output, setOutput] = useState([]);
    const config = useMemo(() => configs[index], [configs, index]);
    useEffect(() => {
        var _a;
        const config = configs[index].transformation;
        const matcher = ((_a = config.filter) === null || _a === void 0 ? void 0 : _a.options) ? getFrameMatchers(config.filter) : undefined;
        const inputTransforms = configs.slice(0, index).map((t) => t.transformation);
        const outputTransforms = configs.slice(index, index + 1).map((t) => t.transformation);
        const ctx = {
            interpolate: (v) => getTemplateSrv().replace(v),
        };
        const inputSubscription = transformDataFrame(inputTransforms, data, ctx).subscribe((v) => {
            if (matcher) {
                v = data.filter((v) => matcher(v));
            }
            setInput(v);
        });
        const outputSubscription = transformDataFrame(inputTransforms, data, ctx)
            .pipe(mergeMap((before) => transformDataFrame(outputTransforms, before, ctx)))
            .subscribe(setOutput);
        return function unsubscribe() {
            inputSubscription.unsubscribe();
            outputSubscription.unsubscribe();
        };
    }, [index, data, configs]);
    const editor = useMemo(() => React.createElement(uiConfig.editor, {
        options: Object.assign(Object.assign({}, uiConfig.transformation.defaultOptions), config.transformation.options),
        input,
        onChange: (opts) => {
            onChange(index, Object.assign(Object.assign({}, config.transformation), { options: opts }));
        },
    }), [uiConfig.editor, uiConfig.transformation.defaultOptions, config.transformation, input, onChange, index]);
    return (React.createElement("div", { className: styles.editor, "data-testid": selectors.components.TransformTab.transformationEditor(uiConfig.name) },
        editor,
        debugMode && (React.createElement("div", { className: styles.debugWrapper, "data-testid": selectors.components.TransformTab.transformationEditorDebugger(uiConfig.name) },
            React.createElement("div", { className: styles.debug },
                React.createElement("div", { className: styles.debugTitle }, "Transformation input data"),
                React.createElement("div", { className: styles.debugJson },
                    React.createElement(JSONFormatter, { json: input }))),
            React.createElement("div", { className: styles.debugSeparator },
                React.createElement(Icon, { name: "arrow-right" })),
            React.createElement("div", { className: styles.debug },
                React.createElement("div", { className: styles.debugTitle }, "Transformation output data"),
                React.createElement("div", { className: styles.debugJson }, output && React.createElement(JSONFormatter, { json: output })))))));
};
const getStyles = (theme) => {
    const debugBorder = theme.isLight ? theme.v1.palette.gray85 : theme.v1.palette.gray15;
    return {
        title: css `
      display: flex;
      padding: 4px 8px 4px 8px;
      position: relative;
      height: 35px;
      border-radius: 4px 4px 0 0;
      flex-wrap: nowrap;
      justify-content: space-between;
      align-items: center;
    `,
        name: css `
      font-weight: ${theme.typography.fontWeightMedium};
      color: ${theme.colors.primary.text};
    `,
        iconRow: css `
      display: flex;
    `,
        icon: css `
      background: transparent;
      border: none;
      box-shadow: none;
      cursor: pointer;
      color: ${theme.colors.text.secondary};
      margin-left: ${theme.spacing(1)};
      &:hover {
        color: ${theme.colors.text};
      }
    `,
        editor: css ``,
        debugWrapper: css `
      display: flex;
      flex-direction: row;
    `,
        debugSeparator: css `
      width: 48px;
      min-height: 300px;
      display: flex;
      align-items: center;
      align-self: stretch;
      justify-content: center;
      margin: 0 ${theme.spacing(0.5)};
      color: ${theme.colors.primary.text};
    `,
        debugTitle: css `
      padding: ${theme.spacing(1)} ${theme.spacing(0.25)};
      font-family: ${theme.typography.fontFamilyMonospace};
      font-size: ${theme.typography.bodySmall.fontSize};
      color: ${theme.colors.text};
      border-bottom: 1px solid ${debugBorder};
      flex-grow: 0;
      flex-shrink: 1;
    `,
        debug: css `
      margin-top: ${theme.spacing(1)};
      padding: 0 ${theme.spacing(1, 1, 1)};
      border: 1px solid ${debugBorder};
      background: ${theme.isLight ? theme.v1.palette.white : theme.v1.palette.gray05};
      border-radius: ${theme.shape.radius.default};
      width: 100%;
      min-height: 300px;
      display: flex;
      flex-direction: column;
      align-self: stretch;
    `,
        debugJson: css `
      flex-grow: 1;
      height: 100%;
      overflow: hidden;
      padding: ${theme.spacing(0.5)};
    `,
    };
};
//# sourceMappingURL=TransformationEditor.js.map