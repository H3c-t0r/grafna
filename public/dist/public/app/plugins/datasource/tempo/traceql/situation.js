import { Aggregate, And, AttributeField, ComparisonOp, FieldExpression, FieldOp, GroupOperation, IntrinsicField, Or, parser, Pipe, ScalarFilter, SelectArgs, SpansetFilter, SpansetPipeline, SpansetPipelineExpression, Static, TraceQL, } from '@grafana/lezer-traceql';
function getErrorNode(tree, cursorPos) {
    const cur = tree.cursorAt(cursorPos);
    do {
        if (cur.from === cursorPos || cur.to === cursorPos) {
            const { node } = cur;
            if (node.type.isError) {
                return node;
            }
        }
    } while (cur.next());
    return null;
}
function move(node, direction) {
    return node[direction];
}
function walk(node, path) {
    let current = node;
    for (const [direction, expectedNodeIDs] of path) {
        current = move(current, direction);
        if (current === null) {
            // we could not move in the direction, we stop
            return null;
        }
        // note that the found value can be 0, which is acceptable
        if (expectedNodeIDs.find((id) => id === (current === null || current === void 0 ? void 0 : current.type.id)) === undefined) {
            // the reached node has wrong type, we stop
            return null;
        }
    }
    return current;
}
function getNodeText(node, text) {
    // if the from and to are them same (e.g. for an error node) we can subtract 1 from the start/from index
    return text.slice(node.from === node.to ? node.from - 1 : node.from, node.to);
}
function isPathMatch(resolverPath, cursorPath) {
    return resolverPath.every((item, index) => item === cursorPath[index]);
}
/**
 * Figure out where is the cursor and what kind of suggestions are appropriate.
 * @param text the user input
 * @param offset the position of the cursor (starting from 0) in the user input
 */
export function getSituation(text, offset) {
    // there is a special case when we are at the start of writing text,
    // so we handle that case first
    if (text === '') {
        return {
            query: text,
            type: 'EMPTY',
        };
    }
    const tree = parser.parse(text);
    // Whitespaces (especially when multiple) on the left of the text cursor can trick the Lezer parser,
    // causing a wrong tree cursor to be picked.
    // Example: `{ span.foo =    ↓ }`, with `↓` being the cursor, tricks the parser.
    // Quick and dirty hack: Shift the cursor to the left until we find a non-whitespace character on its left.
    let shiftedOffset = offset;
    while (shiftedOffset - 1 >= 0 && text[shiftedOffset - 1] === ' ') {
        shiftedOffset -= 1;
    }
    // if the tree contains error, it is very probable that
    // our node is one of those error nodes.
    // also, if there are errors, the node lezer finds us,
    // might not be the best node.
    // so first we check if there is an error node at the cursor position
    let maybeErrorNode = getErrorNode(tree, shiftedOffset);
    if (!maybeErrorNode) {
        // try again with the previous character
        maybeErrorNode = getErrorNode(tree, shiftedOffset - 1);
    }
    const cur = maybeErrorNode != null ? maybeErrorNode.cursor() : tree.cursorAt(shiftedOffset);
    const currentNode = cur.node;
    const ids = [cur.type.id];
    while (cur.parent()) {
        ids.push(cur.type.id);
    }
    let situationType = null;
    for (let resolver of RESOLVERS) {
        if (isPathMatch(resolver.path, ids)) {
            situationType = resolver.fun(currentNode, text, shiftedOffset, offset);
        }
    }
    return Object.assign({ query: text }, (situationType !== null && situationType !== void 0 ? situationType : { type: 'UNKNOWN' }));
}
const ERROR_NODE_ID = 0;
const RESOLVERS = [
    // Curson on error node cases
    {
        path: [ERROR_NODE_ID, AttributeField],
        fun: resolveAttribute,
    },
    {
        path: [ERROR_NODE_ID, FieldExpression],
        fun: resolveExpression,
    },
    {
        path: [ERROR_NODE_ID, SpansetFilter],
        fun: () => ({
            type: 'SPANSET_EXPRESSION_OPERATORS_WITH_MISSING_CLOSED_BRACE',
        }),
    },
    {
        path: [ERROR_NODE_ID, Aggregate],
        fun: resolveAttributeForFunction,
    },
    {
        path: [ERROR_NODE_ID, IntrinsicField],
        fun: resolveAttributeForFunction,
    },
    {
        path: [ERROR_NODE_ID, SpansetPipelineExpression],
        fun: resolveSpansetPipeline,
    },
    {
        path: [ERROR_NODE_ID, ScalarFilter, SpansetPipeline],
        fun: resolveArithmeticOperator,
    },
    {
        path: [ERROR_NODE_ID, TraceQL],
        fun: () => {
            return {
                type: 'UNKNOWN',
            };
        },
    },
    // Curson on valid node cases (the whole query could contain errors nevertheless)
    {
        path: [FieldExpression],
        fun: resolveSpanset,
    },
    {
        path: [SpansetFilter],
        fun: resolveSpanset,
    },
    {
        path: [SpansetPipelineExpression],
        fun: resolveNewSpansetExpression,
    },
    {
        path: [TraceQL],
        fun: resolveNewSpansetExpression,
    },
];
const resolveAttributeCompletion = (node, text, pos) => {
    // The user is completing an expression. We can take advantage of the fact that the Monaco editor is smart
    // enough to automatically detect that there are some characters before the cursor and to take them into
    // account when providing suggestions.
    const endOfPathNode = walk(node, [['firstChild', [FieldExpression]]]);
    if (endOfPathNode && text[pos - 1] !== ' ') {
        const attributeFieldParent = walk(endOfPathNode, [['firstChild', [AttributeField]]]);
        const attributeFieldParentText = attributeFieldParent ? getNodeText(attributeFieldParent, text) : '';
        const indexOfDot = attributeFieldParentText.indexOf('.');
        const attributeFieldUpToDot = attributeFieldParentText.slice(0, indexOfDot);
        return {
            type: 'SPANSET_IN_NAME_SCOPE',
            scope: attributeFieldUpToDot,
        };
    }
};
function resolveSpanset(node, text, _, originalPos) {
    const situation = resolveAttributeCompletion(node, text, originalPos);
    if (situation) {
        return situation;
    }
    let endOfPathNode = walk(node, [
        ['firstChild', [FieldExpression]],
        ['firstChild', [AttributeField]],
    ]);
    if (endOfPathNode) {
        return {
            type: 'SPANSET_EXPRESSION_OPERATORS',
        };
    }
    endOfPathNode = walk(node, [
        ['lastChild', [FieldExpression]],
        ['lastChild', [FieldExpression]],
        ['lastChild', [Static]],
    ]);
    if (endOfPathNode) {
        return {
            type: 'SPANFIELD_COMBINING_OPERATORS',
        };
    }
    endOfPathNode = walk(node, [['lastChild', [FieldExpression]]]);
    if (endOfPathNode) {
        return {
            type: 'SPANSET_EXPRESSION_OPERATORS',
        };
    }
    return {
        type: 'SPANSET_EMPTY',
    };
}
function resolveAttribute(node, text) {
    const attributeFieldParent = walk(node, [['parent', [AttributeField]]]);
    const attributeFieldParentText = attributeFieldParent ? getNodeText(attributeFieldParent, text) : '';
    if (attributeFieldParentText === '.') {
        return {
            type: 'SPANSET_ONLY_DOT',
        };
    }
    const indexOfDot = attributeFieldParentText.indexOf('.');
    const attributeFieldUpToDot = attributeFieldParentText.slice(0, indexOfDot);
    if (['span', 'resource', 'parent'].find((item) => item === attributeFieldUpToDot)) {
        return {
            type: 'SPANSET_IN_NAME_SCOPE',
            scope: attributeFieldUpToDot,
        };
    }
    return {
        type: 'SPANSET_IN_NAME',
    };
}
function resolveExpression(node, text, _, originalPos) {
    var _a, _b, _c;
    const situation = resolveAttributeCompletion(node, text, originalPos);
    if (situation) {
        return situation;
    }
    if (((_a = node.prevSibling) === null || _a === void 0 ? void 0 : _a.type.id) === FieldOp) {
        let attributeField = node.prevSibling.prevSibling;
        if (attributeField) {
            return {
                type: 'SPANSET_IN_VALUE',
                tagName: getNodeText(attributeField, text),
                betweenQuotes: false,
            };
        }
    }
    if (((_b = node.prevSibling) === null || _b === void 0 ? void 0 : _b.type.name) === 'And' || ((_c = node.prevSibling) === null || _c === void 0 ? void 0 : _c.type.name) === 'Or') {
        return {
            type: 'SPANSET_EMPTY',
        };
    }
    return {
        type: 'SPANSET_IN_THE_MIDDLE',
    };
}
function resolveArithmeticOperator(node, _0, _1) {
    var _a;
    if (((_a = node.prevSibling) === null || _a === void 0 ? void 0 : _a.type.id) === ComparisonOp) {
        return {
            type: 'UNKNOWN',
        };
    }
    return {
        type: 'SPANSET_COMPARISON_OPERATORS',
    };
}
function resolveNewSpansetExpression(node, text, offset) {
    // Select the node immediately before the one pointed by the cursor
    let previousNode = node.firstChild;
    try {
        previousNode = node.firstChild;
        while (previousNode.to < offset) {
            previousNode = previousNode.nextSibling;
        }
    }
    catch (error) {
        console.error('Unexpected error while searching for previous node', error);
    }
    if ((previousNode === null || previousNode === void 0 ? void 0 : previousNode.type.id) === And || (previousNode === null || previousNode === void 0 ? void 0 : previousNode.type.id) === Or) {
        return {
            type: 'NEW_SPANSET',
        };
    }
    return {
        type: 'SPANSET_COMBINING_OPERATORS',
    };
}
function resolveAttributeForFunction(node, _0, _1) {
    const parent = node === null || node === void 0 ? void 0 : node.parent;
    if (!!parent && [IntrinsicField, Aggregate, GroupOperation, SelectArgs].includes(parent.type.id)) {
        return {
            type: 'ATTRIBUTE_FOR_FUNCTION',
        };
    }
    return {
        type: 'UNKNOWN',
    };
}
function resolveSpansetPipeline(node, _1, _2) {
    var _a;
    if (((_a = node.prevSibling) === null || _a === void 0 ? void 0 : _a.type.id) === Pipe) {
        return {
            type: 'SPANSET_PIPELINE_AFTER_OPERATOR',
        };
    }
    return {
        type: 'NEW_SPANSET',
    };
}
//# sourceMappingURL=situation.js.map