import { now } from 'lodash';
export function ensureDefine(templateName, templateContent) {
    // notification template content must be wrapped in {{ define "name" }} tag,
    // but this is not obvious because user also has to provide name separately in the form.
    // so if user does not manually add {{ define }} tag, we do it automatically
    let content = templateContent.trim();
    if (!content.match(/\{\{\s*define/)) {
        const indentedContent = content
            .split('\n')
            .map((line) => '  ' + line)
            .join('\n');
        content = `{{ define "${templateName}" }}\n${indentedContent}\n{{ end }}`;
    }
    return content;
}
export function updateDefinesWithUniqueValue(templateContent) {
    const getNewValue = (match_, originalDefineName) => {
        return `{{ define "${originalDefineName}_NEW_${now()}" }}`;
    };
    return templateContent.replace(/\{\{\s*define\s*\"(?<defineName>.*)\"\s*\}\}/g, getNewValue);
}
//# sourceMappingURL=templates.js.map