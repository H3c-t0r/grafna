import { clone, each, map } from 'lodash';
export class QueryPartDef {
    constructor(options) {
        this.type = options.type;
        this.params = options.params;
        this.defaultParams = options.defaultParams;
        this.renderer = options.renderer;
        this.category = options.category;
        this.addStrategy = options.addStrategy;
    }
}
export class QueryPart {
    constructor(part, def) {
        this.part = part;
        this.def = def;
        if (!this.def) {
            throw { message: 'Could not find query part ' + part.type };
        }
        part.params = part.params || clone(this.def.defaultParams);
        this.params = part.params;
        this.text = '';
        this.updateText();
    }
    render(innerExpr) {
        return this.def.renderer(this, innerExpr);
    }
    hasMultipleParamsInString(strValue, index) {
        if (strValue.indexOf(',') === -1) {
            return false;
        }
        return this.def.params[index + 1] && this.def.params[index + 1].optional;
    }
    updateParam(strValue, index) {
        // handle optional parameters
        // if string contains ',' and next param is optional, split and update both
        if (this.hasMultipleParamsInString(strValue, index)) {
            each(strValue.split(','), (partVal, idx) => {
                this.updateParam(partVal.trim(), idx);
            });
            return;
        }
        if (strValue === '' && this.def.params[index].optional) {
            this.params.splice(index, 1);
        }
        else {
            this.params[index] = strValue;
        }
        this.part.params = this.params;
        this.updateText();
    }
    updateText() {
        if (this.params.length === 0) {
            this.text = this.def.type + '()';
            return;
        }
        let text = this.def.type + '(';
        text += this.params.join(', ');
        text += ')';
        this.text = text;
    }
}
export function functionRenderer(part, innerExpr) {
    const str = part.def.type + '(';
    const parameters = map(part.params, (value, index) => {
        const paramType = part.def.params[index];
        if (paramType.type === 'time') {
            if (value === 'auto') {
                value = '$__interval';
            }
        }
        if (paramType.quote === 'single') {
            return "'" + value + "'";
        }
        else if (paramType.quote === 'double') {
            return '"' + value + '"';
        }
        return value;
    });
    if (innerExpr) {
        parameters.unshift(innerExpr);
    }
    return str + parameters.join(', ') + ')';
}
export function suffixRenderer(part, innerExpr) {
    return innerExpr + ' ' + part.params[0];
}
export function identityRenderer(part, innerExpr) {
    return part.params[0];
}
export function quotedIdentityRenderer(part, innerExpr) {
    return '"' + part.params[0] + '"';
}
//# sourceMappingURL=query_part.js.map