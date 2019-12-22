"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const tree = require("./tree");
class VisitorBase {
    constructor(rootNode, parent = null, env = tree.emptyEnvironment) {
        this.rootNode = rootNode;
        this.parent = parent;
        this.env = env;
        this.rootObject = null;
        this.visit = () => {
            this.visitHelper(this.rootNode, this.parent, this.env);
        };
        this.visitHelper = (node, parent, currEnv) => {
            if (node == null) {
                throw Error("INTERNAL ERROR: Can't visit a null node");
            }
            this.previsit(node, parent, currEnv);
            switch (node.type) {
                case "CommentNode": {
                    this.visitComment(node);
                    return;
                }
                case "CompSpecNode": {
                    const castedNode = node;
                    this.visitCompSpec(castedNode);
                    castedNode.varName && this.visitHelper(castedNode.varName, castedNode, currEnv);
                    this.visitHelper(castedNode.expr, castedNode, currEnv);
                    return;
                }
                case "ApplyNode": {
                    const castedNode = node;
                    this.visitApply(castedNode);
                    this.visitHelper(castedNode.target, castedNode, currEnv);
                    castedNode.args.forEach((arg) => {
                        this.visitHelper(arg, castedNode, currEnv);
                    });
                    return;
                }
                case "ApplyBraceNode": {
                    const castedNode = node;
                    this.visitApplyBrace(castedNode);
                    this.visitHelper(castedNode.left, castedNode, currEnv);
                    this.visitHelper(castedNode.right, castedNode, currEnv);
                    return;
                }
                case "ApplyParamAssignmentNode": {
                    const castedNode = node;
                    this.visitApplyParamAssignmentNode(castedNode);
                    this.visitHelper(castedNode.right, castedNode, currEnv);
                    return;
                }
                case "ArrayNode": {
                    const castedNode = node;
                    this.visitArray(castedNode);
                    castedNode.headingComment && this.visitHelper(castedNode.headingComment, castedNode, currEnv);
                    castedNode.elements.forEach((e) => {
                        this.visitHelper(e, castedNode, currEnv);
                    });
                    castedNode.trailingComment && this.visitHelper(castedNode.trailingComment, castedNode, currEnv);
                    return;
                }
                case "ArrayCompNode": {
                    const castedNode = node;
                    this.visitArrayComp(castedNode);
                    this.visitHelper(castedNode.body, castedNode, currEnv);
                    castedNode.specs.forEach((spec) => this.visitHelper(spec, castedNode, currEnv));
                    return;
                }
                case "AssertNode": {
                    const castedNode = node;
                    this.visitAssert(castedNode);
                    this.visitHelper(castedNode.cond, castedNode, currEnv);
                    castedNode.message && this.visitHelper(castedNode.message, castedNode, currEnv);
                    this.visitHelper(castedNode.rest, castedNode, currEnv);
                    return;
                }
                case "BinaryNode": {
                    const castedNode = node;
                    this.visitBinary(castedNode);
                    this.visitHelper(castedNode.left, castedNode, currEnv);
                    this.visitHelper(castedNode.right, castedNode, currEnv);
                    return;
                }
                case "BuiltinNode": {
                    const castedNode = node;
                    this.visitBuiltin(castedNode);
                    return;
                }
                case "ConditionalNode": {
                    const castedNode = node;
                    this.visitConditional(castedNode);
                    this.visitHelper(castedNode.cond, castedNode, currEnv);
                    this.visitHelper(castedNode.branchTrue, castedNode, currEnv);
                    castedNode.branchFalse && this.visitHelper(castedNode.branchFalse, castedNode, currEnv);
                    return;
                }
                case "DollarNode": {
                    const castedNode = node;
                    this.visitDollar(castedNode);
                    return;
                }
                case "ErrorNode": {
                    const castedNode = node;
                    this.visitError(castedNode);
                    this.visitHelper(castedNode.expr, castedNode, currEnv);
                    return;
                }
                case "FunctionNode": {
                    const castedNode = node;
                    this.visitFunction(castedNode);
                    if (castedNode.headingComment != null) {
                        this.visitHelper(castedNode.headingComment, castedNode, currEnv);
                    }
                    // Add params to environment before visiting body.
                    const envWithParams = currEnv.merge(tree.envFromParams(castedNode.parameters));
                    castedNode.parameters.forEach((param) => {
                        this.visitHelper(param, castedNode, envWithParams);
                    });
                    // Visit body.
                    this.visitHelper(castedNode.body, castedNode, envWithParams);
                    castedNode.trailingComment.forEach((comment) => {
                        // NOTE: Using `currEnv` instead of `envWithparams`.
                        this.visitHelper(comment, castedNode, currEnv);
                    });
                    return;
                }
                case "FunctionParamNode": {
                    const castedNode = node;
                    castedNode.defaultValue && this.visitHelper(castedNode.defaultValue, castedNode, currEnv);
                    return;
                }
                case "IdentifierNode": {
                    this.visitIdentifier(node);
                    return;
                }
                case "ImportNode": {
                    this.visitImport(node);
                    return;
                }
                case "ImportStrNode": {
                    this.visitImportStr(node);
                    return;
                }
                case "IndexNode": {
                    const castedNode = node;
                    this.visitIndex(castedNode);
                    castedNode.id != null && this.visitHelper(castedNode.id, castedNode, currEnv);
                    castedNode.target != null && this.visitHelper(castedNode.target, castedNode, currEnv);
                    castedNode.index != null && this.visitHelper(castedNode.index, castedNode, currEnv);
                    return;
                }
                case "LocalBindNode": {
                    const castedNode = node;
                    this.visitLocalBind(node);
                    // NOTE: If `functionSugar` is false, the params will be
                    // empty.
                    const envWithParams = currEnv.merge(tree.envFromParams(castedNode.params));
                    castedNode.params.forEach((param) => {
                        this.visitHelper(param, castedNode, envWithParams);
                    });
                    this.visitHelper(castedNode.body, castedNode, envWithParams);
                    return;
                }
                case "LocalNode": {
                    const castedNode = node;
                    this.visitLocal(castedNode);
                    // NOTE: The binds of a `local` are in scope for both the
                    // binds themselves, as well as the body of the `local`.
                    const envWithBinds = currEnv.merge(tree.envFromLocalBinds(castedNode));
                    castedNode.env = envWithBinds;
                    castedNode.binds.forEach((bind) => {
                        this.visitHelper(bind, castedNode, envWithBinds);
                    });
                    this.visitHelper(castedNode.body, castedNode, envWithBinds);
                    return;
                }
                case "LiteralBooleanNode": {
                    const castedNode = node;
                    this.visitLiteralBoolean(castedNode);
                    return;
                }
                case "LiteralNullNode": {
                    const castedNode = node;
                    this.visitLiteralNull(castedNode);
                    return;
                }
                case "LiteralNumberNode": {
                    return this.visitLiteralNumber(node);
                }
                case "LiteralStringNode": {
                    const castedNode = node;
                    this.visitLiteralString(castedNode);
                    return;
                }
                case "ObjectFieldNode": {
                    const castedNode = node;
                    this.visitObjectField(castedNode);
                    // NOTE: If `methodSugar` is false, the params will be empty.
                    let envWithParams = currEnv.merge(tree.envFromParams(castedNode.ids));
                    castedNode.id != null && this.visitHelper(castedNode.id, castedNode, envWithParams);
                    castedNode.expr1 != null && this.visitHelper(castedNode.expr1, castedNode, envWithParams);
                    castedNode.ids.forEach((param) => {
                        this.visitHelper(param, castedNode, envWithParams);
                    });
                    castedNode.expr2 != null && this.visitHelper(castedNode.expr2, castedNode, envWithParams);
                    castedNode.expr3 != null && this.visitHelper(castedNode.expr3, castedNode, envWithParams);
                    if (castedNode.headingComments != null) {
                        this.visitHelper(castedNode.headingComments, castedNode, currEnv);
                    }
                    return;
                }
                case "ObjectNode": {
                    const castedNode = node;
                    if (this.rootObject == null) {
                        this.rootObject = castedNode;
                        castedNode.rootObject = castedNode;
                    }
                    this.visitObject(castedNode);
                    // `local` object fields are scoped with order-independence,
                    // so something like this is legal:
                    //
                    // {
                    //    bar: {baz: foo},
                    //    local foo = 3,
                    // }
                    //
                    // Since this case requires `foo` to be in the environment of
                    // `bar`'s body, we here collect up the `local` fields first,
                    // create a new environment that includes them, and pass that
                    // on to each field we visit.
                    const envWithLocals = currEnv.merge(tree.envFromFields(castedNode.fields));
                    castedNode.fields.forEach((field) => {
                        // NOTE: If this is a `local` field, there is no need to
                        // remove current field from environment. It is perfectly
                        // legal to do something like `local foo = foo; foo` (though
                        // it will cause a stack overflow).
                        this.visitHelper(field, castedNode, envWithLocals);
                    });
                    return;
                }
                case "DesugaredObjectFieldNode": {
                    const castedNode = node;
                    this.visitDesugaredObjectField(castedNode);
                    this.visitHelper(castedNode.name, castedNode, currEnv);
                    this.visitHelper(castedNode.body, castedNode, currEnv);
                    return;
                }
                case "DesugaredObjectNode": {
                    const castedNode = node;
                    this.visitDesugaredObject(castedNode);
                    castedNode.asserts.forEach((a) => {
                        this.visitHelper(a, castedNode, currEnv);
                    });
                    castedNode.fields.forEach((field) => {
                        this.visitHelper(field, castedNode, currEnv);
                    });
                    return;
                }
                case "ObjectCompNode": {
                    const castedNode = node;
                    this.visitObjectComp(castedNode);
                    castedNode.specs.forEach((spec) => {
                        this.visitHelper(spec, castedNode, currEnv);
                    });
                    castedNode.fields.forEach((field) => {
                        this.visitHelper(field, castedNode, currEnv);
                    });
                    return;
                }
                case "ObjectComprehensionSimpleNode": {
                    const castedNode = node;
                    this.visitObjectComprehensionSimple(castedNode);
                    this.visitHelper(castedNode.id, castedNode, currEnv);
                    this.visitHelper(castedNode.field, castedNode, currEnv);
                    this.visitHelper(castedNode.value, castedNode, currEnv);
                    this.visitHelper(castedNode.array, castedNode, currEnv);
                    return;
                }
                case "SelfNode": {
                    const castedNode = node;
                    this.visitSelf(castedNode);
                    return;
                }
                case "SuperIndexNode": {
                    const castedNode = node;
                    this.visitSuperIndex(castedNode);
                    castedNode.index && this.visitHelper(castedNode.index, castedNode, currEnv);
                    castedNode.id && this.visitHelper(castedNode.id, castedNode, currEnv);
                    return;
                }
                case "UnaryNode": {
                    const castedNode = node;
                    this.visitUnary(castedNode);
                    this.visitHelper(castedNode.expr, castedNode, currEnv);
                    return;
                }
                case "VarNode": {
                    const castedNode = node;
                    this.visitVar(castedNode);
                    castedNode.id != null && this.visitHelper(castedNode.id, castedNode, currEnv);
                    return;
                }
                default: throw new Error(`Visitor could not traverse tree; unknown node type '${node.type}'`);
            }
        };
        this.previsit = (node, parent, currEnv) => { };
        this.visitComment = (node) => { };
        this.visitCompSpec = (node) => { };
        this.visitApply = (node) => { };
        this.visitApplyBrace = (node) => { };
        this.visitApplyParamAssignmentNode = (node) => { };
        this.visitArray = (node) => { };
        this.visitArrayComp = (node) => { };
        this.visitAssert = (node) => { };
        this.visitBinary = (node) => { };
        this.visitBuiltin = (node) => { };
        this.visitConditional = (node) => { };
        this.visitDollar = (node) => { };
        this.visitError = (node) => { };
        this.visitFunction = (node) => { };
        this.visitIdentifier = (node) => { };
        this.visitImport = (node) => { };
        this.visitImportStr = (node) => { };
        this.visitIndex = (node) => { };
        this.visitLocalBind = (node) => { };
        this.visitLocal = (node) => { };
        this.visitLiteralBoolean = (node) => { };
        this.visitLiteralNull = (node) => { };
        this.visitLiteralNumber = (node) => { };
        this.visitLiteralString = (node) => { };
        this.visitObjectField = (node) => { };
        this.visitObject = (node) => { };
        this.visitDesugaredObjectField = (node) => { };
        this.visitDesugaredObject = (node) => { };
        this.visitObjectComp = (node) => { };
        this.visitObjectComprehensionSimple = (node) => { };
        this.visitSelf = (node) => { };
        this.visitSuperIndex = (node) => { };
        this.visitUnary = (node) => { };
        this.visitVar = (node) => { };
    }
}
exports.VisitorBase = VisitorBase;
// ----------------------------------------------------------------------------
// Initializing visitor.
// ----------------------------------------------------------------------------
// InitializingVisitor initializes an AST by populating the `parent`
// and `env` values in every node.
class InitializingVisitor extends VisitorBase {
    constructor() {
        super(...arguments);
        this.previsit = (node, parent, currEnv) => {
            node.parent = parent;
            node.env = currEnv;
            node.rootObject = this.rootObject;
        };
    }
}
exports.InitializingVisitor = InitializingVisitor;
exports.isFindFailure = (thing) => {
    return thing instanceof UnanalyzableFindFailure ||
        thing instanceof AnalyzableFindFailure;
};
// AnalyzableFindFailure represents a failure to find a node whose
// range wraps a cursor location, but which is amenable to static
// analysis.
//
// In particular, this means that the cursor lies in the range of the
// document's AST, and it is therefore possible to inspect the AST
// surrounding the cursor.
class AnalyzableFindFailure {
    // IMPLEMENTATION NOTES: Currently we consider the kind
    // `"AfterDocEnd"` to be unanalyzable, but as our static analysis
    // features become more featureful, we can probably revisit this
    // corner case and get better results in the general case.
    constructor(kind, tightestEnclosingNode, terminalNodeOnCursorLine) {
        this.kind = kind;
        this.tightestEnclosingNode = tightestEnclosingNode;
        this.terminalNodeOnCursorLine = terminalNodeOnCursorLine;
    }
}
exports.AnalyzableFindFailure = AnalyzableFindFailure;
exports.isAnalyzableFindFailure = (thing) => {
    return thing instanceof AnalyzableFindFailure;
};
// UnanalyzableFindFailrue represents a failure to find a node whose
// range wraps a cursor location, and is not amenable to static
// analysis.
//
// In particular, this means that the cursor lies outside of the range
// of a document's AST, which means we cannot inspect the context of
// where the cursor lies in an AST.
class UnanalyzableFindFailure {
    constructor(kind) {
        this.kind = kind;
    }
}
exports.UnanalyzableFindFailure = UnanalyzableFindFailure;
exports.isUnanalyzableFindFailure = (thing) => {
    return thing instanceof UnanalyzableFindFailure;
};
// CursorVisitor finds a node whose range some cursor lies in, or the
// closest node to it.
class CursorVisitor extends VisitorBase {
    // IMPLEMENTATION NOTES: The goal of this class is to map the corner
    // cases into `ast.Node | FindFailure`. Broadly, this mapping falls
    // into a few cases:
    //
    // * Cursor in the range of an identifier.
    //   * Return the identifier.
    // * Cursor in the range of a node that is not an identifier (e.g.,
    //   number literal, multi-line object with no members, and so on).
    //   * Return a find failure with kind `"NotIdentifier"`.
    // * Cursor lies inside document range, the last node on the line
    //   of the cursor ends before the cursor's position.
    //   * Return find failure with kind `"AfterLineEnd"`.
    // * Cursor lies outside document range.
    //   * Return find failure with kind `"BeforeDocStart"` or
    //     `"AfterDocEnd"`.
    constructor(cursor, root) {
        super(root);
        this.cursor = cursor;
        // Identifier whose range encloses the cursor, if there is one. This
        // can be a multi-line node (e.g., perhaps an empty object), or a
        // single line node (e.g., a number literal).
        this.enclosingNode = null;
        // Last node in the line our cursor lies on, if there is one.
        this.terminalNodeOnCursorLine = null;
        this.previsit = (node, parent, currEnv) => {
            const nodeEnd = node.loc.end;
            if (this.cursor.inRange(node.loc)) {
                if (this.enclosingNode == null ||
                    node.loc.rangeIsTighter(this.enclosingNode.loc)) {
                    this.enclosingNode = node;
                }
            }
            if (nodeEnd.afterRangeOrEqual(this.terminalNode.loc)) {
                this.terminalNode = node;
            }
            if (nodeEnd.line === this.cursor.line) {
                if (this.terminalNodeOnCursorLine == null) {
                    this.terminalNodeOnCursorLine = node;
                }
                else if (nodeEnd.afterRangeOrEqual(this.terminalNodeOnCursorLine.loc)) {
                    this.terminalNodeOnCursorLine = node;
                }
            }
        };
        this.terminalNode = root;
    }
    get nodeAtPosition() {
        if (this.enclosingNode == null) {
            if (this.cursor.strictlyBeforeRange(this.rootNode.loc)) {
                return new UnanalyzableFindFailure("BeforeDocStart");
            }
            else if (this.cursor.strictlyAfterRange(this.terminalNode.loc)) {
                return new UnanalyzableFindFailure("AfterDocEnd");
            }
            throw new Error("INTERNAL ERROR: No wrapping identifier was found, but node didn't lie outside of document range");
        }
        else if (!tree.isIdentifier(this.enclosingNode)) {
            if (this.terminalNodeOnCursorLine != null &&
                this.cursor.strictlyAfterRange(this.terminalNodeOnCursorLine.loc)) {
                return new AnalyzableFindFailure("AfterLineEnd", this.enclosingNode, this.terminalNodeOnCursorLine);
            }
            return new AnalyzableFindFailure("NotIdentifier", this.enclosingNode, this.terminalNodeOnCursorLine);
        }
        return this.enclosingNode;
    }
}
exports.CursorVisitor = CursorVisitor;
// nodeRangeIsCloser checks whether `thisNode` is closer to `pos` than
// `thatNode`.
//
// NOTE: Function currently works for expressions that are on one
// line.
const nodeRangeIsCloser = (pos, thisNode, thatNode) => {
    const thisLoc = thisNode.loc;
    const thatLoc = thatNode.loc;
    if (thisLoc.begin.line == pos.line && thisLoc.end.line == pos.line) {
        if (thatLoc.begin.line == pos.line && thatLoc.end.line == pos.line) {
            // `thisNode` and `thatNode` lie on the same line, and
            // `thisNode` begins closer to the position.
            //
            // NOTE: We use <= here so that we always choose the last node
            // that begins at a point. For example, a `Var` and `Identifier`
            // might begin in the same place, but we'd like to choose the
            // `Identifier`, as it would be a child of the `Var`.
            return Math.abs(thisLoc.begin.column - pos.column) <=
                Math.abs(thatLoc.begin.column - pos.column);
        }
        else {
            return true;
        }
    }
    return false;
};
//# sourceMappingURL=visitor.js.map