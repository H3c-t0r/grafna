import { Registry } from '@grafana/data';
import { FieldExtractorID } from './types';
const extJSON = {
    id: FieldExtractorID.JSON,
    name: 'JSON',
    description: 'Parse JSON string',
    parse: (v) => {
        return JSON.parse(v);
    },
};
function parseKeyValuePairs(raw) {
    const buff = []; // array of characters
    let esc = '';
    let key = '';
    const obj = {};
    for (let i = 0; i < raw.length; i++) {
        let c = raw[i];
        if (c === esc) {
            esc = '';
            c = raw[++i];
        }
        const isEscaped = c === '\\';
        if (isEscaped) {
            c = raw[++i];
        }
        // When escaped just append
        if (isEscaped || esc.length) {
            buff.push(c);
            continue;
        }
        if (c === `"` || c === `'`) {
            esc = c;
        }
        switch (c) {
            case ':':
            case '=':
                if (buff.length) {
                    if (key) {
                        obj[key] = '';
                    }
                    key = buff.join('');
                    buff.length = 0; // clear values
                }
                break;
            // escape chars
            case `"`:
            case `'`:
            // whitespace
            case ` `:
            case `\n`:
            case `\t`:
            case `\r`:
            case `\n`:
                if (buff.length && key === '') {
                    obj[buff.join('')] = '';
                    buff.length = 0;
                }
            // seperators
            case ',':
            case ';':
            case '&':
            case '{':
            case '}':
                if (buff.length) {
                    const val = buff.join('');
                    if (key.length) {
                        obj[key] = val;
                        key = '';
                    }
                    else {
                        key = val;
                    }
                    buff.length = 0; // clear values
                }
                break;
            // append our buffer
            default:
                buff.push(c);
                if (i === raw.length - 1) {
                    if (key === '' && buff.length) {
                        obj[buff.join('')] = '';
                        buff.length = 0;
                    }
                }
        }
    }
    if (key.length) {
        obj[key] = buff.join('');
    }
    return obj;
}
const extLabels = {
    id: FieldExtractorID.KeyValues,
    name: 'Key+value pairs',
    description: 'Look for a=b, c: d values in the line',
    parse: parseKeyValuePairs,
};
const fmts = [extJSON, extLabels];
const extAuto = {
    id: FieldExtractorID.Auto,
    name: 'Auto',
    description: 'parse new fields automatically',
    parse: (v) => {
        for (const f of fmts) {
            try {
                const r = f.parse(v);
                if (r != null) {
                    return r;
                }
            }
            catch (_a) { } // ignore errors
        }
        return undefined;
    },
};
export const fieldExtractors = new Registry(() => [...fmts, extAuto]);
//# sourceMappingURL=fieldExtractors.js.map