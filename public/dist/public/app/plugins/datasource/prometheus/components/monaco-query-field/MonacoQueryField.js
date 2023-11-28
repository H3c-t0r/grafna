import { css } from '@emotion/css';
import { debounce } from 'lodash';
import { promLanguageDefinition } from 'monaco-promql';
import React, { useRef, useEffect } from 'react';
import { useLatest } from 'react-use';
import { v4 as uuidv4 } from 'uuid';
import { selectors } from '@grafana/e2e-selectors';
import { useTheme2, ReactMonacoEditor } from '@grafana/ui';
import { getOverrideServices } from './getOverrideServices';
import { getCompletionProvider, getSuggestOptions } from './monaco-completion-provider';
const options = {
    codeLens: false,
    contextmenu: false,
    // we need `fixedOverflowWidgets` because otherwise in grafana-dashboards
    // the popup is clipped by the panel-visualizations.
    fixedOverflowWidgets: true,
    folding: false,
    fontSize: 14,
    lineDecorationsWidth: 8,
    lineNumbers: 'off',
    minimap: { enabled: false },
    overviewRulerBorder: false,
    overviewRulerLanes: 0,
    padding: {
        // these numbers were picked so that visually this matches the previous version
        // of the query-editor the best
        top: 4,
        bottom: 5,
    },
    renderLineHighlight: 'none',
    scrollbar: {
        vertical: 'hidden',
        verticalScrollbarSize: 8,
        horizontal: 'hidden',
        horizontalScrollbarSize: 0,
    },
    scrollBeyondLastLine: false,
    suggest: getSuggestOptions(),
    suggestFontSize: 12,
    wordWrap: 'on',
};
// this number was chosen by testing various values. it might be necessary
// because of the width of the border, not sure.
//it needs to do 2 things:
// 1. when the editor is single-line, it should make the editor height be visually correct
// 2. when the editor is multi-line, the editor should not be "scrollable" (meaning,
//    you do a scroll-movement in the editor, and it will scroll the content by a couple pixels
//    up & down. this we want to avoid)
const EDITOR_HEIGHT_OFFSET = 2;
const PROMQL_LANG_ID = promLanguageDefinition.id;
// we must only run the promql-setup code once
let PROMQL_SETUP_STARTED = false;
function ensurePromQL(monaco) {
    if (PROMQL_SETUP_STARTED === false) {
        PROMQL_SETUP_STARTED = true;
        const { aliases, extensions, mimetypes, loader } = promLanguageDefinition;
        monaco.languages.register({ id: PROMQL_LANG_ID, aliases, extensions, mimetypes });
        loader().then((mod) => {
            monaco.languages.setMonarchTokensProvider(PROMQL_LANG_ID, mod.language);
            monaco.languages.setLanguageConfiguration(PROMQL_LANG_ID, mod.languageConfiguration);
        });
    }
}
const getStyles = (theme, placeholder) => {
    return {
        container: css `
      border-radius: ${theme.shape.radius.default};
      border: 1px solid ${theme.components.input.borderColor};
    `,
        placeholder: css `
      ::after {
        content: '${placeholder}';
        font-family: ${theme.typography.fontFamilyMonospace};
        opacity: 0.3;
      }
    `,
    };
};
const MonacoQueryField = (props) => {
    const id = uuidv4();
    // we need only one instance of `overrideServices` during the lifetime of the react component
    const overrideServicesRef = useRef(getOverrideServices());
    const containerRef = useRef(null);
    const { languageProvider, history, onBlur, onRunQuery, initialValue, placeholder, onChange } = props;
    const lpRef = useLatest(languageProvider);
    const historyRef = useLatest(history);
    const onRunQueryRef = useLatest(onRunQuery);
    const onBlurRef = useLatest(onBlur);
    const onChangeRef = useLatest(onChange);
    const autocompleteDisposeFun = useRef(null);
    const theme = useTheme2();
    const styles = getStyles(theme, placeholder);
    useEffect(() => {
        // when we unmount, we unregister the autocomplete-function, if it was registered
        return () => {
            var _a;
            (_a = autocompleteDisposeFun.current) === null || _a === void 0 ? void 0 : _a.call(autocompleteDisposeFun);
        };
    }, []);
    return (React.createElement("div", { "aria-label": selectors.components.QueryField.container, className: styles.container, 
        // NOTE: we will be setting inline-style-width/height on this element
        ref: containerRef },
        React.createElement(ReactMonacoEditor, { overrideServices: overrideServicesRef.current, options: options, language: "promql", value: initialValue, beforeMount: (monaco) => {
                ensurePromQL(monaco);
            }, onMount: (editor, monaco) => {
                var _a;
                const isEditorFocused = editor.createContextKey('isEditorFocused' + id, false);
                // we setup on-blur
                editor.onDidBlurEditorWidget(() => {
                    isEditorFocused.set(false);
                    onBlurRef.current(editor.getValue());
                });
                editor.onDidFocusEditorText(() => {
                    isEditorFocused.set(true);
                });
                // we construct a DataProvider object
                const getHistory = () => Promise.resolve(historyRef.current.map((h) => h.query.expr).filter((expr) => expr !== undefined));
                const getAllMetricNames = () => {
                    const { metrics, metricsMetadata } = lpRef.current;
                    const result = metrics.map((m) => {
                        var _a, _b;
                        const metaItem = metricsMetadata === null || metricsMetadata === void 0 ? void 0 : metricsMetadata[m];
                        return {
                            name: m,
                            help: (_a = metaItem === null || metaItem === void 0 ? void 0 : metaItem.help) !== null && _a !== void 0 ? _a : '',
                            type: (_b = metaItem === null || metaItem === void 0 ? void 0 : metaItem.type) !== null && _b !== void 0 ? _b : '',
                        };
                    });
                    return Promise.resolve(result);
                };
                const getAllLabelNames = () => Promise.resolve(lpRef.current.getLabelKeys());
                const getLabelValues = (labelName) => lpRef.current.getLabelValues(labelName);
                const getSeriesValues = lpRef.current.getSeriesValues;
                const getSeriesLabels = lpRef.current.getSeriesLabels;
                const dataProvider = {
                    getHistory,
                    getAllMetricNames,
                    getAllLabelNames,
                    getLabelValues,
                    getSeriesValues,
                    getSeriesLabels,
                };
                const completionProvider = getCompletionProvider(monaco, dataProvider);
                // completion-providers in monaco are not registered directly to editor-instances,
                // they are registered to languages. this makes it hard for us to have
                // separate completion-providers for every query-field-instance
                // (but we need that, because they might connect to different datasources).
                // the trick we do is, we wrap the callback in a "proxy",
                // and in the proxy, the first thing is, we check if we are called from
                // "our editor instance", and if not, we just return nothing. if yes,
                // we call the completion-provider.
                const filteringCompletionProvider = Object.assign(Object.assign({}, completionProvider), { provideCompletionItems: (model, position, context, token) => {
                        var _a;
                        // if the model-id does not match, then this call is from a different editor-instance,
                        // not "our instance", so return nothing
                        if (((_a = editor.getModel()) === null || _a === void 0 ? void 0 : _a.id) !== model.id) {
                            return { suggestions: [] };
                        }
                        return completionProvider.provideCompletionItems(model, position, context, token);
                    } });
                const { dispose } = monaco.languages.registerCompletionItemProvider(PROMQL_LANG_ID, filteringCompletionProvider);
                autocompleteDisposeFun.current = dispose;
                // this code makes the editor resize itself so that the content fits
                // (it will grow taller when necessary)
                // FIXME: maybe move this functionality into CodeEditor, like:
                // <CodeEditor resizingMode="single-line"/>
                const updateElementHeight = () => {
                    const containerDiv = containerRef.current;
                    if (containerDiv !== null) {
                        const pixelHeight = editor.getContentHeight();
                        containerDiv.style.height = `${pixelHeight + EDITOR_HEIGHT_OFFSET}px`;
                        containerDiv.style.width = '100%';
                        const pixelWidth = containerDiv.clientWidth;
                        editor.layout({ width: pixelWidth, height: pixelHeight });
                    }
                };
                editor.onDidContentSizeChange(updateElementHeight);
                updateElementHeight();
                // Whenever the editor changes, lets save the last value so the next query for this editor will be up-to-date.
                // This change is being introduced to fix a bug where you can submit a query via shift+enter:
                // If you clicked into another field and haven't un-blurred the active field,
                // then the query that is run will be stale, as the reference is only updated
                // with the value of the last blurred input.
                // This can run quite slowly, so we're debouncing this which should accomplish two things
                // 1. Should prevent this function from blocking the current call stack by pushing into the web API callback queue
                // 2. Should prevent a bunch of duplicates of this function being called as the user is typing
                const updateCurrentEditorValue = debounce(() => {
                    const editorValue = editor.getValue();
                    onChangeRef.current(editorValue);
                }, lpRef.current.datasource.getDebounceTimeInMilliseconds());
                (_a = editor.getModel()) === null || _a === void 0 ? void 0 : _a.onDidChangeContent(() => {
                    updateCurrentEditorValue();
                });
                // handle: shift + enter
                // FIXME: maybe move this functionality into CodeEditor?
                editor.addCommand(monaco.KeyMod.Shift | monaco.KeyCode.Enter, () => {
                    onRunQueryRef.current(editor.getValue());
                }, 'isEditorFocused' + id);
                /* Something in this configuration of monaco doesn't bubble up [mod]+K, which the
                command palette uses. Pass the event out of monaco manually
                */
                editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyK, function () {
                    global.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', metaKey: true }));
                });
                if (placeholder) {
                    const placeholderDecorators = [
                        {
                            range: new monaco.Range(1, 1, 1, 1),
                            options: {
                                className: styles.placeholder,
                                isWholeLine: true,
                            },
                        },
                    ];
                    let decorators = [];
                    const checkDecorators = () => {
                        const model = editor.getModel();
                        if (!model) {
                            return;
                        }
                        const newDecorators = model.getValueLength() === 0 ? placeholderDecorators : [];
                        decorators = model.deltaDecorations(decorators, newDecorators);
                    };
                    checkDecorators();
                    editor.onDidChangeModelContent(checkDecorators);
                }
            } })));
};
// we will lazy-load this module using React.lazy,
// and that only supports default-exports,
// so we have to default-export this, even if
// it is against the style-guidelines.
export default MonacoQueryField;
//# sourceMappingURL=MonacoQueryField.js.map