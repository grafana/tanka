import * as server from 'vscode-languageserver';

import * as im from 'immutable';

import * as lexical from '../lexical';
import * as tree from './tree';

export interface Visitor {
  visit(): void
}

export abstract class VisitorBase implements Visitor {
  protected rootObject: tree.Node | null = null;

  constructor(
    protected rootNode: tree.Node,
    private parent: tree.Node | null = null,
    private env: tree.Environment = tree.emptyEnvironment,
  ) {}

  public visit = () => {
    this.visitHelper(this.rootNode, this.parent, this.env);
  }

  protected visitHelper = (
    node: tree.Node, parent: tree.Node | null, currEnv: tree.Environment
  ): void => {
    if (node == null) {
      throw Error("INTERNAL ERROR: Can't visit a null node");
    }

    this.previsit(node, parent, currEnv);

    switch(node.type) {
      case "CommentNode": {
        this.visitComment(<tree.Comment>node);
        return;
      }
      case "CompSpecNode": {
        const castedNode = <tree.CompSpec>node;
        this.visitCompSpec(castedNode);
        castedNode.varName && this.visitHelper(castedNode.varName, castedNode, currEnv);
        this.visitHelper(castedNode.expr, castedNode, currEnv);
        return;
      }
      case "ApplyNode": {
        const castedNode = <tree.Apply>node;
        this.visitApply(castedNode);
        this.visitHelper(castedNode.target, castedNode, currEnv);
        castedNode.args.forEach((arg: tree.Node) => {
          this.visitHelper(arg, castedNode, currEnv);
        });
        return;
      }
      case "ApplyBraceNode": {
        const castedNode = <tree.ApplyBrace>node;
        this.visitApplyBrace(castedNode);
        this.visitHelper(castedNode.left, castedNode, currEnv);
        this.visitHelper(castedNode.right, castedNode, currEnv);
        return;
      }
      case "ApplyParamAssignmentNode": {
        const castedNode = <tree.ApplyParamAssignment>node;
        this.visitApplyParamAssignmentNode(castedNode);
        this.visitHelper(castedNode.right, castedNode, currEnv);
        return;
      }
      case "ArrayNode": {
        const castedNode = <tree.Array>node;
        this.visitArray(castedNode);
        castedNode.headingComment && this.visitHelper(
          castedNode.headingComment, castedNode, currEnv);
        castedNode.elements.forEach((e: tree.Node) => {
          this.visitHelper(e, castedNode, currEnv);
        });
        castedNode.trailingComment && this.visitHelper(
          castedNode.trailingComment, castedNode, currEnv);
        return;
      }
      case "ArrayCompNode": {
        const castedNode = <tree.ArrayComp>node;
        this.visitArrayComp(castedNode);
        this.visitHelper(castedNode.body, castedNode, currEnv);
        castedNode.specs.forEach((spec: tree.CompSpec) =>
          this.visitHelper(spec, castedNode, currEnv));
        return;
      }
      case "AssertNode": {
        const castedNode = <tree.Assert>node;
        this.visitAssert(castedNode);
        this.visitHelper(castedNode.cond, castedNode, currEnv);
        castedNode.message && this.visitHelper(
          castedNode.message, castedNode, currEnv);
        this.visitHelper(castedNode.rest, castedNode, currEnv);
        return;
      }
      case "BinaryNode": {
        const castedNode = <tree.Binary>node;
        this.visitBinary(castedNode);
        this.visitHelper(castedNode.left, castedNode, currEnv);
        this.visitHelper(castedNode.right, castedNode, currEnv);
        return;
      }
      case "BuiltinNode": {
        const castedNode = <tree.Builtin>node;
        this.visitBuiltin(castedNode);
        return;
      }
      case "ConditionalNode": {
        const castedNode = <tree.Conditional>node;
        this.visitConditional(castedNode);
        this.visitHelper(castedNode.cond, castedNode, currEnv);
        this.visitHelper(castedNode.branchTrue, castedNode, currEnv);
        castedNode.branchFalse && this.visitHelper(
          castedNode.branchFalse, castedNode, currEnv);
        return;
      }
      case "DollarNode": {
        const castedNode = <tree.Dollar>node;
        this.visitDollar(castedNode);
        return;
      }
      case "ErrorNode": {
        const castedNode = <tree.ErrorNode>node;
        this.visitError(castedNode);
        this.visitHelper(castedNode.expr, castedNode, currEnv);
        return;
      }
      case "FunctionNode": {
        const castedNode = <tree.Function>node;
        this.visitFunction(castedNode);

        if (castedNode.headingComment != null) {
          this.visitHelper(castedNode.headingComment, castedNode, currEnv);
        }

        // Add params to environment before visiting body.
        const envWithParams = currEnv.merge(
          tree.envFromParams(castedNode.parameters));

        castedNode.parameters.forEach((param: tree.FunctionParam) => {
          this.visitHelper(param, castedNode, envWithParams);
        });

        // Visit body.
        this.visitHelper(castedNode.body, castedNode, envWithParams);
        castedNode.trailingComment.forEach((comment: tree.Comment) => {
          // NOTE: Using `currEnv` instead of `envWithparams`.
          this.visitHelper(comment, castedNode, currEnv);
        });
        return;
      }
      case "FunctionParamNode": {
        const castedNode = <tree.FunctionParam>node;
        castedNode.defaultValue && this.visitHelper(
          castedNode.defaultValue, castedNode, currEnv);
        return;
      }
      case "IdentifierNode": {
        this.visitIdentifier(<tree.Identifier>node);
        return;
      }
      case "ImportNode": {
        this.visitImport(<tree.Import>node);
        return;
      }
      case "ImportStrNode": {
        this.visitImportStr(<tree.ImportStr>node);
        return;
      }
      case "IndexNode": {
        const castedNode = <tree.Index>node;
        this.visitIndex(castedNode);
        castedNode.id != null && this.visitHelper(castedNode.id, castedNode, currEnv);
        castedNode.target != null && this.visitHelper(
          castedNode.target, castedNode, currEnv);
        castedNode.index != null && this.visitHelper(
          castedNode.index, castedNode, currEnv);
        return;
      }
      case "LocalBindNode": {
        const castedNode = <tree.LocalBind>node;
        this.visitLocalBind(<tree.LocalBind>node);

        // NOTE: If `functionSugar` is false, the params will be
        // empty.
        const envWithParams = currEnv.merge(
          tree.envFromParams(castedNode.params));

        castedNode.params.forEach((param: tree.FunctionParam) => {
          this.visitHelper(param, castedNode, envWithParams)
        });

        this.visitHelper(castedNode.body, castedNode, envWithParams);
        return;
      }
      case "LocalNode": {
        const castedNode = <tree.Local>node;
        this.visitLocal(castedNode);

        // NOTE: The binds of a `local` are in scope for both the
        // binds themselves, as well as the body of the `local`.
        const envWithBinds = currEnv.merge(tree.envFromLocalBinds(castedNode));
        castedNode.env = envWithBinds;

        castedNode.binds.forEach((bind: tree.LocalBind) => {
          this.visitHelper(bind, castedNode, envWithBinds);
        });

        this.visitHelper(castedNode.body, castedNode, envWithBinds);
        return;
      }
      case "LiteralBooleanNode": {
        const castedNode = <tree.LiteralBoolean>node;
        this.visitLiteralBoolean(castedNode);
        return;
      }
      case "LiteralNullNode": {
        const castedNode = <tree.LiteralNull>node;
        this.visitLiteralNull(castedNode);
        return;
      }
      case "LiteralNumberNode": { return this.visitLiteralNumber(<tree.LiteralNumber>node); }
      case "LiteralStringNode": {
        const castedNode = <tree.LiteralString>node;
        this.visitLiteralString(castedNode);
        return;
      }
      case "ObjectFieldNode": {
        const castedNode = <tree.ObjectField>node;
        this.visitObjectField(castedNode);

        // NOTE: If `methodSugar` is false, the params will be empty.
        let envWithParams = currEnv.merge(tree.envFromParams(castedNode.ids));

        castedNode.id != null && this.visitHelper(
          castedNode.id, castedNode, envWithParams);
        castedNode.expr1 != null && this.visitHelper(
          castedNode.expr1, castedNode, envWithParams);

        castedNode.ids.forEach((param: tree.FunctionParam) => {
          this.visitHelper(param, castedNode, envWithParams);
        });

        castedNode.expr2 != null && this.visitHelper(
          castedNode.expr2, castedNode, envWithParams);
        castedNode.expr3 != null && this.visitHelper(
          castedNode.expr3, castedNode, envWithParams);
        if (castedNode.headingComments != null) {
          this.visitHelper(castedNode.headingComments, castedNode, currEnv);
        }
        return;
      }
      case "ObjectNode": {
        const castedNode = <tree.ObjectNode>node;
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
        const envWithLocals = currEnv.merge(
          tree.envFromFields(castedNode.fields));

        castedNode.fields.forEach((field: tree.ObjectField) => {
          // NOTE: If this is a `local` field, there is no need to
          // remove current field from environment. It is perfectly
          // legal to do something like `local foo = foo; foo` (though
          // it will cause a stack overflow).
          this.visitHelper(field, castedNode, envWithLocals);
        });
        return;
      }
      case "DesugaredObjectFieldNode": {
        const castedNode = <tree.DesugaredObjectField>node;
        this.visitDesugaredObjectField(castedNode);
        this.visitHelper(castedNode.name, castedNode, currEnv);
        this.visitHelper(castedNode.body, castedNode, currEnv);
        return;
      }
      case "DesugaredObjectNode": {
        const castedNode = <tree.DesugaredObject>node;
        this.visitDesugaredObject(castedNode);
        castedNode.asserts.forEach((a: tree.Assert) => {
          this.visitHelper(a, castedNode, currEnv);
        });
        castedNode.fields.forEach((field: tree.DesugaredObjectField) => {
          this.visitHelper(field, castedNode, currEnv);
        });
        return;
      }
      case "ObjectCompNode": {
        const castedNode = <tree.ObjectComp>node;
        this.visitObjectComp(castedNode);
        castedNode.specs.forEach((spec: tree.CompSpec) => {
          this.visitHelper(spec, castedNode, currEnv);
        });
        castedNode.fields.forEach((field: tree.ObjectField) => {
          this.visitHelper(field, castedNode, currEnv);
        });
        return;
      }
      case "ObjectComprehensionSimpleNode": {
        const castedNode = <tree.ObjectComprehensionSimple>node;
        this.visitObjectComprehensionSimple(castedNode);
        this.visitHelper(castedNode.id, castedNode, currEnv);
        this.visitHelper(castedNode.field, castedNode, currEnv);
        this.visitHelper(castedNode.value, castedNode, currEnv);
        this.visitHelper(castedNode.array, castedNode, currEnv);
        return;
      }
      case "SelfNode": {
        const castedNode = <tree.Self>node;
        this.visitSelf(castedNode);
        return;
      }
      case "SuperIndexNode": {
        const castedNode = <tree.SuperIndex>node;
        this.visitSuperIndex(castedNode);
        castedNode.index && this.visitHelper(castedNode.index, castedNode, currEnv);
        castedNode.id && this.visitHelper(castedNode.id, castedNode, currEnv);
        return;
      }
      case "UnaryNode": {
        const castedNode = <tree.Unary>node;
        this.visitUnary(castedNode);
        this.visitHelper(castedNode.expr, castedNode, currEnv);
        return;
      }
      case "VarNode": {
        const castedNode = <tree.Var>node;
        this.visitVar(castedNode);
        castedNode.id != null && this.visitHelper(castedNode.id, castedNode, currEnv);
        return
      }
      default: throw new Error(
        `Visitor could not traverse tree; unknown node type '${node.type}'`);
    }
  }

  protected previsit = (
    node: tree.Node, parent: tree.Node | null, currEnv: tree.Environment
  ): void => {}

  protected visitComment = (node: tree.Comment): void => {}
  protected visitCompSpec = (node: tree.CompSpec): void => {}
  protected visitApply = (node: tree.Apply): void => {}
  protected visitApplyBrace = (node: tree.ApplyBrace): void => {}
  protected visitApplyParamAssignmentNode = (node: tree.ApplyParamAssignment): void => {}
  protected visitArray = (node: tree.Array): void => {}
  protected visitArrayComp = (node: tree.ArrayComp): void => {}
  protected visitAssert = (node: tree.Assert): void => {}
  protected visitBinary = (node: tree.Binary): void => {}
  protected visitBuiltin = (node: tree.Builtin): void => {}
  protected visitConditional = (node: tree.Conditional): void => {}
  protected visitDollar = (node: tree.Dollar): void => {}
  protected visitError = (node: tree.ErrorNode): void => {}
  protected visitFunction = (node: tree.Function): void => {}

  protected visitIdentifier = (node: tree.Identifier): void => {}
  protected visitImport = (node: tree.Import): void => {}
  protected visitImportStr = (node: tree.ImportStr): void => {}
  protected visitIndex = (node: tree.Index): void => {}
  protected visitLocalBind = (node: tree.LocalBind): void => {}
  protected visitLocal = (node: tree.Local): void => {}

  protected visitLiteralBoolean = (node: tree.LiteralBoolean): void => {}
  protected visitLiteralNull = (node: tree.LiteralNull): void => {}

  protected visitLiteralNumber = (node: tree.LiteralNumber): void => {}
  protected visitLiteralString = (node: tree.LiteralString): void => {}
  protected visitObjectField = (node: tree.ObjectField): void => {}
  protected visitObject = (node: tree.ObjectNode): void => {}
  protected visitDesugaredObjectField = (node: tree.DesugaredObjectField): void => {}
  protected visitDesugaredObject = (node: tree.DesugaredObject): void => {}
  protected visitObjectComp = (node: tree.ObjectComp): void => {}
  protected visitObjectComprehensionSimple = (node: tree.ObjectComprehensionSimple): void => {}
  protected visitSelf = (node: tree.Self): void => {}
  protected visitSuperIndex = (node: tree.SuperIndex): void => {}
  protected visitUnary = (node: tree.Unary): void => {}
  protected visitVar = (node: tree.Var): void => {}
}

// ----------------------------------------------------------------------------
// Initializing visitor.
// ----------------------------------------------------------------------------

// InitializingVisitor initializes an AST by populating the `parent`
// and `env` values in every node.
export class InitializingVisitor extends VisitorBase {
  protected previsit = (
    node: tree.Node, parent: tree.Node | null, currEnv: tree.Environment
  ): void => {
    node.parent = parent;
    node.env = currEnv;
    node.rootObject = this.rootObject;
  }
}

// ----------------------------------------------------------------------------
// Cursor visitor.
// ----------------------------------------------------------------------------

// FindFailure represents a failure find a node whose range wraps a
// cursor location.
export type FindFailure = AnalyzableFindFailure | UnanalyzableFindFailure;

export const isFindFailure = (thing): thing is FindFailure => {
  return thing instanceof UnanalyzableFindFailure ||
    thing instanceof AnalyzableFindFailure;
}

export type FindFailureKind =
  "BeforeDocStart" | "AfterDocEnd" | "AfterLineEnd" | "NotIdentifier";

// AnalyzableFindFailure represents a failure to find a node whose
// range wraps a cursor location, but which is amenable to static
// analysis.
//
// In particular, this means that the cursor lies in the range of the
// document's AST, and it is therefore possible to inspect the AST
// surrounding the cursor.
export class AnalyzableFindFailure {
  // IMPLEMENTATION NOTES: Currently we consider the kind
  // `"AfterDocEnd"` to be unanalyzable, but as our static analysis
  // features become more featureful, we can probably revisit this
  // corner case and get better results in the general case.

  constructor(
    public readonly kind: "AfterLineEnd" | "NotIdentifier",
    public readonly tightestEnclosingNode: tree.Node,
    public readonly terminalNodeOnCursorLine: tree.Node | null,
  ) {}
}

export const isAnalyzableFindFailure = (
  thing
): thing is AnalyzableFindFailure => {
  return thing instanceof AnalyzableFindFailure;
}

// UnanalyzableFindFailrue represents a failure to find a node whose
// range wraps a cursor location, and is not amenable to static
// analysis.
//
// In particular, this means that the cursor lies outside of the range
// of a document's AST, which means we cannot inspect the context of
// where the cursor lies in an AST.
export class UnanalyzableFindFailure {
  constructor(public readonly kind: "BeforeDocStart" | "AfterDocEnd") {}
}

export const isUnanalyzableFindFailure = (
  thing
): thing is UnanalyzableFindFailure => {
  return thing instanceof UnanalyzableFindFailure;
}

// CursorVisitor finds a node whose range some cursor lies in, or the
// closest node to it.
export class CursorVisitor extends VisitorBase {
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

  constructor(
    private cursor: lexical.Location,
    root: tree.Node,
  ) {
    super(root);
    this.terminalNode = root;
  }

  // Identifier whose range encloses the cursor, if there is one. This
  // can be a multi-line node (e.g., perhaps an empty object), or a
  // single line node (e.g., a number literal).
  private enclosingNode: tree.Node | null = null;

  // Last node in the tree.
  private terminalNode: tree.Node;

  // Last node in the line our cursor lies on, if there is one.
  private terminalNodeOnCursorLine: tree.Node | null = null;

  get nodeAtPosition(): tree.Identifier | FindFailure {
    if (this.enclosingNode == null) {
      if (this.cursor.strictlyBeforeRange(this.rootNode.loc)) {
        return new UnanalyzableFindFailure("BeforeDocStart");
      } else if (this.cursor.strictlyAfterRange(this.terminalNode.loc)) {
        return new UnanalyzableFindFailure("AfterDocEnd");
      }
      throw new Error(
        "INTERNAL ERROR: No wrapping identifier was found, but node didn't lie outside of document range");
    } else if (!tree.isIdentifier(this.enclosingNode)) {
      if (
        this.terminalNodeOnCursorLine != null &&
        this.cursor.strictlyAfterRange(this.terminalNodeOnCursorLine.loc)
      ) {
        return new AnalyzableFindFailure(
          "AfterLineEnd", this.enclosingNode, this.terminalNodeOnCursorLine);
      }
      return new AnalyzableFindFailure(
        "NotIdentifier", this.enclosingNode, this.terminalNodeOnCursorLine);
    }
    return this.enclosingNode;
  }

  protected previsit = (
    node: tree.Node, parent: tree.Node | null, currEnv: tree.Environment,
  ): void => {
    const nodeEnd = node.loc.end;

    if (this.cursor.inRange(node.loc)) {
      if (
        this.enclosingNode == null ||
        node.loc.rangeIsTighter(this.enclosingNode.loc)
      ) {
        this.enclosingNode = node;
      }
    }

    if (nodeEnd.afterRangeOrEqual(this.terminalNode.loc)) {
      this.terminalNode = node;
    }

    if (nodeEnd.line === this.cursor.line) {
      if (this.terminalNodeOnCursorLine == null) {
        this.terminalNodeOnCursorLine = node;
      } else if (nodeEnd.afterRangeOrEqual(this.terminalNodeOnCursorLine.loc)) {
        this.terminalNodeOnCursorLine = node;
      }
    }
  }
}

// nodeRangeIsCloser checks whether `thisNode` is closer to `pos` than
// `thatNode`.
//
// NOTE: Function currently works for expressions that are on one
// line.
const nodeRangeIsCloser = (
  pos: lexical.Location, thisNode: tree.Node, thatNode: tree.Node
): boolean => {
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
        Math.abs(thatLoc.begin.column - pos.column)
    } else {
      return true;
    }
  }

  return false;
}
