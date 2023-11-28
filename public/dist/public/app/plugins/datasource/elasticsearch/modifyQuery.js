import { isEqual } from 'lodash';
import lucene from 'lucene';
/**
 * Checks for the presence of a given label:"value" filter in the query.
 */
export function queryHasFilter(query, key, value, modifier = '') {
    return findFilterNode(query, key, value, modifier) !== null;
}
/**
 * Given a query, find the NodeTerm that matches the given field and value.
 */
export function findFilterNode(query, key, value, modifier = '') {
    const field = `${modifier}${lucene.term.escape(key)}`;
    value = lucene.phrase.escape(value);
    let ast = parseQuery(query);
    if (!ast) {
        return null;
    }
    return findNodeInTree(ast, field, value);
}
function findNodeInTree(ast, field, value) {
    // {}
    if (Object.keys(ast).length === 0) {
        return null;
    }
    // { left: {}, right: {} } or { left: {} }
    if (isAST(ast.left)) {
        return findNodeInTree(ast.left, field, value);
    }
    if (isNodeTerm(ast.left) && ast.left.field === field && ast.left.term === value) {
        return ast.left;
    }
    if (isLeftOnlyAST(ast)) {
        return null;
    }
    if (isNodeTerm(ast.right) && ast.right.field === field && ast.right.term === value) {
        return ast.right;
    }
    if (isBinaryAST(ast.right)) {
        return findNodeInTree(ast.right, field, value);
    }
    return null;
}
/**
 * Adds a label:"value" expression to the query.
 */
export function addFilterToQuery(query, key, value, modifier = '') {
    if (queryHasFilter(query, key, value, modifier)) {
        return query;
    }
    key = lucene.term.escape(key);
    value = lucene.phrase.escape(value);
    const filter = `${modifier}${key}:"${value}"`;
    return query === '' ? filter : `${query} AND ${filter}`;
}
/**
 * Removes a label:"value" expression from the query.
 */
export function removeFilterFromQuery(query, key, value, modifier = '') {
    const node = findFilterNode(query, key, value, modifier);
    const ast = parseQuery(query);
    if (!node || !ast) {
        return query;
    }
    return lucene.toString(removeNodeFromTree(ast, node));
}
function removeNodeFromTree(ast, node) {
    // {}
    if (Object.keys(ast).length === 0) {
        return ast;
    }
    // { left: {}, right: {} } or { left: {} }
    if (isAST(ast.left)) {
        ast.left = removeNodeFromTree(ast.left, node);
        return ast;
    }
    if (isNodeTerm(ast.left) && isEqual(ast.left, node)) {
        Object.assign(ast, {
            left: undefined,
            operator: undefined,
            right: undefined,
        }, 'right' in ast ? ast.right : {});
        return ast;
    }
    if (isLeftOnlyAST(ast)) {
        return ast;
    }
    if (isNodeTerm(ast.right) && isEqual(ast.right, node)) {
        Object.assign(ast, {
            right: undefined,
            operator: undefined,
        });
        return ast;
    }
    if (isBinaryAST(ast.right)) {
        ast.right = removeNodeFromTree(ast.right, node);
        return ast;
    }
    return ast;
}
/**
 * Filters can possibly reserved characters such as colons which are part of the Lucene syntax.
 * Use this function to escape filter keys.
 */
export function escapeFilter(value) {
    return lucene.term.escape(value);
}
/**
 * Values can possibly reserved special characters such as quotes.
 * Use this function to escape filter values.
 */
export function escapeFilterValue(value) {
    return lucene.phrase.escape(value);
}
/**
 * Normalizes the query by removing whitespace around colons, which breaks parsing.
 */
function normalizeQuery(query) {
    return query.replace(/(\w+)\s(:)/gi, '$1$2');
}
function isLeftOnlyAST(ast) {
    if (!ast) {
        return false;
    }
    if ('left' in ast && !('right' in ast)) {
        return true;
    }
    return false;
}
function isBinaryAST(ast) {
    if (!ast) {
        return false;
    }
    if ('left' in ast && 'right' in ast) {
        return true;
    }
    return false;
}
function isAST(ast) {
    return isLeftOnlyAST(ast) || isBinaryAST(ast);
}
function isNodeTerm(ast) {
    if (!ast) {
        return false;
    }
    if ('term' in ast) {
        return true;
    }
    return false;
}
function parseQuery(query) {
    try {
        return lucene.parse(normalizeQuery(query));
    }
    catch (e) {
        return null;
    }
}
//# sourceMappingURL=modifyQuery.js.map