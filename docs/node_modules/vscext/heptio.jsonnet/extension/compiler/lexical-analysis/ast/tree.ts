import * as os from 'os';
import * as path from 'path';

import * as im from 'immutable';

import * as editor from '../../editor';
import * as lexical from '../lexical';
import * as _static from '../../static';

// ---------------------------------------------------------------------------

export type Environment = im.Map<string, LocalBind | FunctionParam>;

export const emptyEnvironment = im.Map<string, LocalBind | FunctionParam>();

export const envFromLocalBinds = (
  local: Local | ObjectField | FunctionParam
): Environment => {
  if (isLocal(local)) {
    const defaultLocal: {[key: string]: LocalBind} = {};
    const binds = local.binds
      .reduce(
        (acc: {[key: string]: LocalBind}, bind: LocalBind) => {
          acc[bind.variable.name] = bind;
          return acc;
        },
        defaultLocal);
    return im.Map(binds);
  } else if (isObjectField(local)) {
    if (local.expr2 == null || local.id == null) {
      throw new Error(`INTERNAL ERROR: Object local fields can't have a null expr2 or id field`);
    }

    const bind: LocalBind = new LocalBind(
      local.id,
      local.expr2,
      local.methodSugar,
      local.ids,
      local.trailingComma,
      local.loc,
    );
    return im.Map<string, LocalBind>().set(local.id.name, bind);
  }

  // Else, it's a `FunctionParam`, i.e., a free parameter (or a free
  // parameter with a default value). Either way, emit that.
  return im.Map<string, LocalBind | FunctionParam>().set(local.id, local);
}

export const envFromParams = (
  params: FunctionParams
): Environment => {
  return params
    .reduce(
      (acc: Environment, field: FunctionParam) => {
        return acc.merge(envFromLocalBinds(field));
      },
      emptyEnvironment
    );
}

export const envFromFields = (
  fields: ObjectFields,
): Environment => {
  return fields
    .filter((field: ObjectField) => {
      const localKind: ObjectFieldKind = "ObjectLocal";
      return field.kind === localKind;
    })
    .reduce(
      (acc: Environment, field: ObjectField) => {
        return acc.merge(envFromLocalBinds(field));
      },
      emptyEnvironment
    );
}

export const renderAsJson = (node: Node): string => {
  return "```\n" + JSON.stringify(
  node,
  (k, v) => {
    if (k === "parent") {
      return v == null
        ? "null"
        : (<Node>v).type;
    } else if (k === "env") {
      return v == null
        ? "null"
        : `${Object.keys(v).join(", ")}`;
    } else if (k === "rootObject") {
      return v == null
        ? "null"
        : (<Node>v).type;
    } else {
      return v;
    }
  },
  "  ") + "\n```";
}

// ---------------------------------------------------------------------------

// NodeKind captures the type of the node. Implementing this as a
// union of specific strings allows us to `switch` on node type.
// Additionally, specific nodes can specialize and restrict the `type`
// field to be something like `type: "ObjectNode" = "ObjectNode"`,
// which will cause a type error if something tries to instantiate on
// `ObjectNode` with a `type` that is not this specific string.
export type NodeKind =
  "CommentNode" |
  "CompSpecNode" |
  "ApplyNode" |
  "ApplyBraceNode" |
  "ApplyParamAssignmentNode" |
  "ArrayNode" |
  "ArrayCompNode" |
  "AssertNode" |
  "BinaryNode" |
  "BuiltinNode" |
  "ConditionalNode" |
  "DollarNode" |
  "ErrorNode" |
  "FunctionNode" |
  "FunctionParamNode" |
  "IdentifierNode" |
  "ImportNode" |
  "ImportStrNode" |
  "IndexNode" |
  "LocalBindNode" |
  "LocalNode" |
  "LiteralBooleanNode" |
  "LiteralNullNode" |
  "LiteralNumberNode" |
  "LiteralStringNode" |
  "ObjectFieldNode" |
  "ObjectNode" |
  "DesugaredObjectFieldNode" |
  "DesugaredObjectNode" |
  "ObjectCompNode" |
  "ObjectComprehensionSimpleNode" |
  "SelfNode" |
  "SuperIndexNode" |
  "UnaryNode" |
  "VarNode";

// isValueType returns true if the node is a computed value literal
// (e.g., a string literal, an object literal, and so on).
//
// Notably, this explicitly omits structures whose value must be
// computed at runtime: particularly object comprehension, array
// comprehension, self, super, and function types (whose value depends
// on parameter binds).
export const isValueType = (node: Node): boolean => {
  // TODO(hausdorff): Consider adding object comprehension here, too.
  return isLiteralBoolean(node) || isLiteralNull(node) ||
    isLiteralNumber(node) || isLiteralString(node) || isObjectNode(node);
}

// ---------------------------------------------------------------------------

export interface Node {
  readonly type:     NodeKind
  readonly loc:      lexical.LocationRange

  prettyPrint(): string

  rootObject: Node | null;
  parent: Node | null;     // Filled in by the visitor.
  env: Environment | null; // Filled in by the visitor.
}
export type Nodes = im.List<Node>

// NodeBase is a simple abstract base class that makes sure we're
// initializing the parent and env members to null. It is not exposed
// to the public because it is meant to be a transparent base blass
// for all `Node` implementations.
abstract class NodeBase implements Node {
  readonly type:     NodeKind
  readonly loc:      lexical.LocationRange

  constructor() {
    this.rootObject = null;
    this.parent = null;
    this.env = null;
  }

  abstract prettyPrint: () => string;

  rootObject: Node | null;
  parent: Node | null;     // Filled in by the visitor.
  env: Environment | null; // Filled in by the visitor.
}

export const isNode = (thing): thing is Node => {
  // TODO: Probably want to check the types of the properties instead.
  return thing instanceof NodeBase;
}

// ---------------------------------------------------------------------------

// Resolve represents a resolved node, including a fully-qualified RFC
// 1630/1738-compliant URI representing the absolute location of the
// Jsonnet file the symbol occurs in.
export class Resolve {
  constructor(
    public readonly fileUri:  editor.FileUri,
    public readonly value: Node | IndexedObjectFields,
  ) {}
}

export const isResolve = (thing): thing is Resolve => {
  return thing instanceof Resolve;
}

// ResolutionContext represents the context we carry along as we
// attempt to resolve symbols. For example, an `import` node will have
// a filename, and to locate it, we will need to (1) search for the
// import path relative to the current path, or (2) look in the
// `libPaths` for it if necessary. This "context" is carried along in
// this object.
export class ResolutionContext {
  constructor (
    public readonly compiler: _static.LexicalAnalyzerService,
    public readonly documents: editor.DocumentManager,
    public readonly currFile: editor.FileUri,
  ) {}

  public withUri = (currFile: editor.FileUri): ResolutionContext => {
    return new ResolutionContext(this.compiler, this.documents, currFile);
  }
}

export interface Resolvable extends NodeBase {
  resolve(context: ResolutionContext): Resolve | ResolveFailure
}

export const isResolvable = (node: NodeBase): node is Resolvable => {
  return node instanceof NodeBase && typeof node["resolve"] === "function";
}

export interface FieldsResolvable extends NodeBase {
  resolveFields(context: ResolutionContext): Resolve | ResolveFailure
}

export const isFieldsResolvable = (
  node: NodeBase
): node is FieldsResolvable => {
  return node instanceof NodeBase &&
    typeof node["resolveFields"] === "function";
}

export interface TypeGuessResolvable extends NodeBase {
  resolveTypeGuess(context: ResolutionContext): Resolve | ResolveFailure
}

export const isTypeGuessResolvable = (
  node: NodeBase
): node is TypeGuessResolvable => {
  return node instanceof NodeBase &&
    typeof node["resolveTypeGuess"] === "function";
}

// ---------------------------------------------------------------------------

// IdentifierName represents a variable / parameter / field name.
//+gen set
export type IdentifierName = string
export type IdentifierNames = im.List<IdentifierName>
export type IdentifierSet = im.Set<IdentifierName>;

export class Identifier extends NodeBase  {
  readonly type: "IdentifierNode" = "IdentifierNode";

  constructor(
    readonly name: IdentifierName,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return this.name;
  }

  public resolve = (context: ResolutionContext): Resolve | ResolveFailure => {
    if (this.parent == null) {
      // An identifier with no parent is not a valid Jsonnet file.
      return Unresolved.Instance;
    }

    return tryResolve(this.parent, context)
  }
}

export const isIdentifier = (node): node is Identifier => {
  return node instanceof Identifier;
}

// TODO(jbeda) implement interning of IdentifierNames if necessary.  The C++
// version does so.

// ---------------------------------------------------------------------------

export type CommentKind =
  "CppStyle" |
  "CStyle" |
  "HashStyle";

export interface Comment extends Node {
  // TODO: This Kind/Type part seems wrong, as it does in
  // `ObjectField`.
  readonly type: "CommentNode";
  readonly kind: CommentKind
  readonly text: im.List<string>
};
export type Comments = im.List<Comment>;

export type BindingComment = CppComment | CComment | HashComment | null;

export const isBindingComment = (node): node is BindingComment => {
  return isCppComment(node) || isCComment(node);
}

export const isComment = (node: Node): node is Comment => {
  const nodeType: NodeKind = "CommentNode";
  return node.type === nodeType;
}

export class CppComment extends NodeBase implements Comment {
  readonly type: "CommentNode" = "CommentNode";
  readonly kind: "CppStyle"    = "CppStyle";

  constructor(
    readonly text: im.List<string>,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return this.text.join(os.EOL);
  }
}

export const isCppComment = (node): node is CppComment => {
  return node instanceof CppComment;
}

export class CComment extends NodeBase implements Comment {
  readonly type: "CommentNode" = "CommentNode";
  readonly kind: "CStyle"    = "CStyle";

  constructor(
    readonly text: im.List<string>,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return this.text.join(os.EOL);
  }
}

export const isCComment = (node): node is CComment => {
  return node instanceof CComment;
}

export class HashComment extends NodeBase implements Comment {
  readonly type: "CommentNode" = "CommentNode";
  readonly kind: "HashStyle"    = "HashStyle";

  constructor(
    readonly text: im.List<string>,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return this.text.join(os.EOL);
  }
}

export const isHashComment = (node): node is HashComment => {
  return node instanceof HashComment;
}

// ---------------------------------------------------------------------------

export type CompKind =
  "CompFor" |
  "CompIf";

export interface CompSpec extends Node {
  readonly type:    "CompSpecNode"
  readonly kind:    CompKind
  readonly varName: Identifier | null // null when kind != compSpecFor
  readonly expr:    Node
};
export type CompSpecs = im.List<CompSpec>;

export const isCompSpec = (node: Node): node is CompSpec => {
  const nodeType: NodeKind = "CompSpecNode";
  return node.type === nodeType;
}

export class CompSpecIf extends NodeBase implements CompSpec {
  readonly type:    "CompSpecNode" = "CompSpecNode";
  readonly kind:    "CompIf"       = "CompIf";
  readonly varName: Identifier | null = null // null when kind != compSpecFor

  constructor(
    readonly expr: Node,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `if ${this.expr.prettyPrint()}`;
  }
}

export const isCompSpecIf = (node): node is CompSpec => {
  return node instanceof CompSpecIf;
}

export class CompSpecFor extends NodeBase implements CompSpec {
  readonly type:    "CompSpecNode" = "CompSpecNode";
  readonly kind:    "CompFor"      = "CompFor";

  constructor(
    readonly varName: Identifier, // null for `CompSpecIf`
    readonly expr:    Node,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `for ${this.varName.prettyPrint()} in ${this.expr.prettyPrint()}`;
  }
}

export const isCompSpecFor = (node): node is CompSpec => {
  return node instanceof CompSpecFor;
}

// ---------------------------------------------------------------------------

// Apply represents a function call
export class Apply extends NodeBase implements TypeGuessResolvable {
  readonly type: "ApplyNode" = "ApplyNode";

  constructor(
    readonly target:        Node,
    readonly args:          Nodes,
    readonly trailingComma: boolean,
    readonly tailStrict:    boolean,
    readonly loc:           lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    const argsString = this.args
      .map((arg: Node) => arg.prettyPrint())
      .join(", ");

    // NOTE: Space between `tailstrict` is important.
    const tailStrictString = this.tailStrict
      ? " tailstrict"
      : "";

    return `${this.target.prettyPrint()}(${argsString}${tailStrictString})`;
  }

  public resolveTypeGuess = (
    context: ResolutionContext
  ): Resolve | ResolveFailure => {
    if (!isResolvable(this.target)) {
      return Unresolved.Instance;
    }

    const fn = this.target.resolve(context);
    if (!isResolvedFunction(fn) || !isObjectField(fn.functionNode)) {
      return Unresolved.Instance;
    }

    const body = fn.functionNode.expr2;
    if (isBinary(body) && body.op == "BopPlus" && isSelf(body.left)) {
      return body.left.resolve(context);
    }
    return Unresolved.Instance;
  }
}

export const isApply = (node): node is Apply => {
  return node instanceof Apply;
}

export class ApplyParamAssignment extends NodeBase {
  readonly type: "ApplyParamAssignmentNode" = "ApplyParamAssignmentNode";

  constructor(
    readonly id:    IdentifierName,
    readonly right: Node,
    readonly loc:   lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `${this.id}=${this.right.prettyPrint()}`;
  }
}
export type ApplyParamAssignments = im.List<ApplyParamAssignment>

export const isApplyParamAssignment = (node): node is ApplyParamAssignment => {
  return node instanceof ApplyParamAssignment;
};

// ---------------------------------------------------------------------------

// ApplyBrace represents e { }.  Desugared to e + { }.
export class ApplyBrace extends NodeBase {
  readonly type: "ApplyBraceNode" = "ApplyBraceNode";

  constructor(
    readonly left:  Node,
    readonly right: Node,
    readonly loc:   lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `${this.left.prettyPrint()} ${this.right.prettyPrint()}`;
  }
}

export const isApplyBrace = (node): node is ApplyBrace => {
  return node instanceof ApplyBrace;
}

// ---------------------------------------------------------------------------

// Array represents array constructors [1, 2, 3].
export class Array extends NodeBase {
  readonly type: "ArrayNode" = "ArrayNode";

  constructor(
    readonly elements:        Nodes,
    readonly trailingComma:   boolean,
    readonly headingComment:  Comment | null,
    readonly trailingComment: Comment | null,
    readonly loc:             lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    const elementsString = this.elements
      .map((element: Node) => element.prettyPrint())
      .join(", ");
    return `[${elementsString}]`;
  }
}

export const isArray = (node): node is Array => {
  return node instanceof Array;
}

// ---------------------------------------------------------------------------

// ArrayComp represents array comprehensions (which are like Python list
// comprehensions)
export class ArrayComp extends NodeBase {
  readonly type: "ArrayCompNode" = "ArrayCompNode";

  constructor(
    readonly body:          Node,
    readonly trailingComma: boolean,
    readonly specs:         CompSpecs,
    readonly loc:           lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    const specsString = this.specs
      .map((spec: CompSpec) => spec.prettyPrint())
      .join(", ");
    return `[${specsString} ${this.body.prettyPrint()}]`;
  }
}

export const isArrayComp = (node): node is ArrayComp => {
  return node instanceof ArrayComp;
}

// ---------------------------------------------------------------------------

// Assert represents an assert expression (not an object-level assert).
//
// After parsing, message can be nil indicating that no message was
// specified. This AST is elimiated by desugaring.
export class Assert extends NodeBase {
  readonly type: "AssertNode" = "AssertNode";

  constructor(
    readonly cond:    Node,
    readonly message: Node | null,
    readonly rest:    Node,
    readonly loc:     lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `assert ${this.cond.prettyPrint()}`;
  }
}

export const isAssert = (node): node is Assert => {
  return node instanceof Assert;
}

// ---------------------------------------------------------------------------

export type BinaryOp =
  "BopMult" |
  "BopDiv" |
  "BopPercent" |

  "BopPlus" |
  "BopMinus" |

  "BopShiftL" |
  "BopShiftR" |

  "BopGreater" |
  "BopGreaterEq" |
  "BopLess" |
  "BopLessEq" |

  "BopManifestEqual" |
  "BopManifestUnequal" |

  "BopBitwiseAnd" |
  "BopBitwiseXor" |
  "BopBitwiseOr" |

  "BopAnd" |
  "BopOr";

const BopStrings = {
  BopMult:    "*",
  BopDiv:     "/",
  BopPercent: "%",

  BopPlus:  "+",
  BopMinus: "-",

  BopShiftL: "<<",
  BopShiftR: ">>",

  BopGreater:   ">",
  BopGreaterEq: ">=",
  BopLess:      "<",
  BopLessEq:    "<=",

  BopManifestEqual:   "==",
  BopManifestUnequal: "!=",

  BopBitwiseAnd: "&",
  BopBitwiseXor: "^",
  BopBitwiseOr:  "|",

  BopAnd: "&&",
  BopOr:  "||",
};

export const BopMap = im.Map<string, BinaryOp>({
  "*": "BopMult",
  "/": "BopDiv",
  "%": "BopPercent",

  "+": "BopPlus",
  "-": "BopMinus",

  "<<": "BopShiftL",
  ">>": "BopShiftR",

  ">":  "BopGreater",
  ">=": "BopGreaterEq",
  "<":  "BopLess",
  "<=": "BopLessEq",

  "==": "BopManifestEqual",
  "!=": "BopManifestUnequal",

  "&": "BopBitwiseAnd",
  "^": "BopBitwiseXor",
  "|": "BopBitwiseOr",

  "&&": "BopAnd",
  "||": "BopOr",
});

// Binary represents binary operators.
export class Binary extends NodeBase implements FieldsResolvable {
  readonly type: "BinaryNode" = "BinaryNode";

  constructor(
    readonly left:  Node,
    readonly op:    BinaryOp,
    readonly right: Node,
    readonly loc:     lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    const leftString = this.left.prettyPrint();
    const opString = BopStrings[this.op];
    const rightString = this.right.prettyPrint();
    return `${leftString} ${opString} ${rightString}`;
  }

  public resolveFields = (
    context: ResolutionContext,
  ): Resolve | ResolveFailure => {
    // Recursively merge fields if it's another mixin; if it's an
    // object, return fields; else, no fields to return.
    if (this.op !== "BopPlus") {
      return Unresolved.Instance;
    }

    const left = tryResolveIndirections(this.left, context);
    if (isResolveFailure(left) || !isIndexedObjectFields(left.value)) {
      return Unresolved.Instance;
    }

    const right = tryResolveIndirections(this.right, context);
    if (isResolveFailure(right) || !isIndexedObjectFields(right.value)) {
      return Unresolved.Instance;
    }

    let merged = left.value;
    right.value.forEach(
      (v: ObjectField, k: string) => {
        // TODO(hausdorff): Account for syntax sugar here. For
        // example:
        //
        //   `{foo: "bar"} + {foo+: "bar"}`
        //
        // should produce `{foo: "barbar"} but because we are merging
        // naively, we will report the value as simply `"bar"`. The
        // reason we have punted for now is mainly that we have to
        // implement an ad hoc version of Jsonnet's type coercion.
        merged = merged.set(k, v);
      });

    return new Resolve(context.currFile, merged);
  }
}

export const isBinary = (node): node is Binary => {
  return node instanceof Binary;
}

// ---------------------------------------------------------------------------

// Builtin represents built-in functions.
//
// There is no parse rule to build this AST.  Instead, it is used to build the
// std object in the interpreter.
export class Builtin extends NodeBase {
  readonly type: "BuiltinNode" = "BuiltinNode";

  constructor(
    readonly id:     number,
    readonly params: IdentifierNames,
    readonly loc:    lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    const paramsString = this.params.join(", ");
    return `std.${this.id}(${paramsString})`;
  }
}

export const isBuiltin = (node): node is Builtin => {
  return node instanceof Builtin;
}

// ---------------------------------------------------------------------------

// Conditional represents if/then/else.
//
// After parsing, branchFalse can be nil indicating that no else branch
// was specified.  The desugarer fills this in with a LiteralNull
export class Conditional extends NodeBase {
  readonly type: "ConditionalNode" = "ConditionalNode";

  constructor(
    readonly cond:        Node,
    readonly branchTrue:  Node,
    readonly branchFalse: Node | null,
    readonly loc:    lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    const trueClause = `then ${this.branchTrue.prettyPrint()}`;
    const falseClause = this.branchFalse == null
      ? ""
      : `else ${this.branchFalse.prettyPrint()}`;
    return `if ${this.cond.prettyPrint()} ${trueClause} ${falseClause}`;
  }
}

export const isConditional = (node): node is Conditional => {
  return node instanceof Conditional;
}

// ---------------------------------------------------------------------------

// Dollar represents the $ keyword
export class Dollar extends NodeBase implements Resolvable {
  readonly type: "DollarNode" = "DollarNode";

  constructor(
    readonly loc:    lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `$`;
  }

  public resolve = (context: ResolutionContext): Resolve | ResolveFailure => {
    if (this.rootObject == null) {
      return Unresolved.Instance;
    }
    return new Resolve(context.currFile, this.rootObject);
  }
};

export const isDollar = (node): node is Dollar => {
  return node instanceof Dollar;
}

// ---------------------------------------------------------------------------

// Error represents the error e.
export class ErrorNode extends NodeBase {
  readonly type: "ErrorNode" = "ErrorNode";

  constructor(
    readonly expr: Node,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `error ${this.expr.prettyPrint()}`;
  }
}

export const isError = (node): node is ErrorNode => {
  return node instanceof ErrorNode;
}

// ---------------------------------------------------------------------------

// Function represents a function call. (jbeda: or is it function defn?)
export class Function extends NodeBase {
  readonly type: "FunctionNode" = "FunctionNode";

  constructor(
    readonly parameters:      FunctionParams,
    readonly trailingComma:   boolean,
    readonly body:            Node,
    readonly headingComment:  BindingComment,
    readonly trailingComment: Comments,
    readonly loc:             lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    const params = this.parameters
      .map((param: FunctionParam) => param.prettyPrint())
      .join(", ");
    return `function (${params}) ${this.body.prettyPrint()}`;
  }
}

export const isFunction = (node): node is Function => {
  return node instanceof Function;
}

export class FunctionParam extends NodeBase {
  readonly type: "FunctionParamNode" = "FunctionParamNode";

  constructor(
    readonly id:           IdentifierName,
    readonly defaultValue: Node | null,
    readonly loc:             lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    const defaultValueString = this.defaultValue == null
      ? ""
      : `=${this.defaultValue.prettyPrint()}`;
    return `(parameter) ${this.id}${defaultValueString}`;
  }
}
export type FunctionParams = im.List<FunctionParam>

export const isFunctionParam = (node): node is FunctionParam => {
  return node instanceof FunctionParam;
}

// ---------------------------------------------------------------------------

// Import represents import "file".
export class Import extends NodeBase implements Resolvable {
  readonly type: "ImportNode" = "ImportNode";

  constructor(
    readonly file: string,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `import "${this.file}"`;
  }

  public resolve = (context: ResolutionContext): Resolve | ResolveFailure => {
    const {text: docText, version: version, resolvedPath: fileUri} =
      context.documents.get(this);
    const cached =
      context.compiler.cache(fileUri, docText, version);
    if (_static.isFailedParsedDocument(cached)) {
      return Unresolved.Instance;
    }

    let resolved = cached.parse;
    // If the var was pointing at an import, then resolution probably
    // has `local` definitions at the top of the file. Get rid of
    // them, since they are not useful for resolving the index
    // identifier.
    while (isLocal(resolved)) {
      resolved = resolved.body;
    }

    return new Resolve(fileUri, resolved);
  }
}

export const isImport = (node): node is Import => {
  return node instanceof Import;
}

// ---------------------------------------------------------------------------

// ImportStr represents importstr "file".
export class ImportStr extends NodeBase {
  readonly type: "ImportStrNode" = "ImportStrNode";

  constructor(
    readonly file: string,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `importstr "${this.file}"`;
  }
}

export const isImportStr = (node): node is ImportStr => {
  return node instanceof ImportStr;
}

// ---------------------------------------------------------------------------

// Index represents both e[e] and the syntax sugar e.f.
//
// One of index and id will be nil before desugaring.  After desugaring id
// will be nil.
export interface Index extends Node, Resolvable {
  readonly type:   "IndexNode"
  readonly target: Node
  readonly index:  Node | null
  readonly id:     Identifier | null
}

const resolveIndex = (
  index: Index, context: ResolutionContext,
): Resolve | ResolveFailure => {
  if (
    index.target == null ||
    (!isResolvable(index.target) && !isFieldsResolvable(index.target) && !isTypeGuessResolvable(index.target))
  ) {
    throw new Error(
      `INTERNAL ERROR: Index node must have a resolvable target:\n${renderAsJson(index)}`);
  } else if (index.id == null) {
    return Unresolved.Instance;
  }

  // Find root target, look up in environment.
  let resolvedTarget = tryResolveIndirections(index.target, context);
  if (isResolveFailure(resolvedTarget)) {
    return new UnresolvedIndexTarget(index);
  } else if (!isIndexedObjectFields(resolvedTarget.value)) {
    return new UnresolvedIndexTarget(index);
  }

  const filtered = resolvedTarget.value.filter((field: ObjectField) => {
    return field.id != null && index.id != null &&
      field.id.name == index.id.name;
  });

  if (filtered.count() == 0) {
    return new UnresolvedIndexId(index, resolvedTarget.value);
  } else if (filtered.count() != 1) {
    throw new Error(
      `INTERNAL ERROR: Object contained multiple fields with name '${index.id.name}'}`);
  }

  const field = filtered.first();
  if (field.methodSugar) {
    return new ResolvedFunction(field);
  } else if (field.expr2 == null) {
    throw new Error(
      `INTERNAL ERROR: Object field can't have null property expr2:\n${renderAsJson(field)}'}`);
  }
  return new Resolve(context.currFile, field.expr2);
}

export const isIndex = (node: Node): node is Index => {
  const nodeType: NodeKind = "IndexNode";
  return node.type === nodeType;
}

export class IndexSubscript extends NodeBase implements Index {
  readonly type: "IndexNode"        = "IndexNode";
  readonly id:    Identifier | null = null;

  constructor(
    readonly target: Node,
    readonly index:  Node,
    readonly loc:    lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `${this.target.prettyPrint()}[${this.index.prettyPrint()}]`;
  }

  public resolve = (context: ResolutionContext): Resolve | ResolveFailure =>
    resolveIndex(this, context);
}

export const isIndexSubscript = (node): node is Index => {
  return node instanceof IndexSubscript;
}

export class IndexDot extends NodeBase implements Index {
  readonly type:  "IndexNode" = "IndexNode";
  readonly index: Node | null = null;

  constructor(
    readonly target: Node,
    readonly id:     Identifier,
    readonly loc:    lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `${this.target.prettyPrint()}.${this.id.prettyPrint()}`;
  }

  public resolve = (context: ResolutionContext): Resolve | ResolveFailure =>
    resolveIndex(this, context);
}

export const isIndexDot = (node): node is Index => {
  return node instanceof IndexDot;
}

// ---------------------------------------------------------------------------

// LocalBind is a helper struct for Local
export class LocalBind extends NodeBase {
  readonly type: "LocalBindNode" = "LocalBindNode";

  constructor(
    readonly variable:      Identifier,
    readonly body:          Node,
    readonly functionSugar: boolean,
    readonly params:        FunctionParams, // if functionSugar is true
    readonly trailingComma: boolean,
    readonly loc:           lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    const idString = this.variable.prettyPrint();
    if (this.functionSugar) {
      const paramsString = this.params
        .map((param: FunctionParam) => param.id)
        .join(", ");
      return `${idString}(${paramsString})`;
    }
    return `${idString} = ${this.body.prettyPrint()}`;
  }
}
export type LocalBinds = im.List<LocalBind>

export const isLocalBind = (node): node is LocalBind => {
  return node instanceof LocalBind;
}

// Local represents local x = e; e.  After desugaring, functionSugar is false.
export class Local extends NodeBase {
  readonly type: "LocalNode" = "LocalNode";

  constructor(
    readonly binds: LocalBinds,
    readonly body:  Node,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    const bindsString = this.binds
      .map((bind: LocalBind) => bind.prettyPrint())
      .join(",\n  ");

    return `local ${bindsString}`;
  }
}

export const isLocal = (node): node is Local => {
  return node instanceof Local;
}

// ---------------------------------------------------------------------------

// LiteralBoolean represents true and false
export class LiteralBoolean extends NodeBase {
  readonly type: "LiteralBooleanNode" = "LiteralBooleanNode";

  constructor(
    readonly value: boolean,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `${this.value}`;
  }
}

export const isLiteralBoolean = (node): node is LiteralBoolean => {
  return node instanceof LiteralBoolean;
}

// ---------------------------------------------------------------------------

// LiteralNull represents the null keyword
export class LiteralNull extends NodeBase {
  readonly type: "LiteralNullNode" = "LiteralNullNode";

  constructor(
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `null`;
  }
}

export const isLiteralNull = (node): node is LiteralNull => {
  return node instanceof LiteralNull;
}

// ---------------------------------------------------------------------------

// LiteralNumber represents a JSON number
export class LiteralNumber extends NodeBase {
  readonly type: "LiteralNumberNode" = "LiteralNumberNode";

  constructor(
    readonly value:          number,
    readonly originalString: string,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `${this.originalString}`;
  }
}

export const isLiteralNumber = (node): node is LiteralNumber => {
  return node instanceof LiteralNumber;
}

// ---------------------------------------------------------------------------

export type LiteralStringKind =
  "StringSingle" |
  "StringDouble" |
  "StringBlock";

// LiteralString represents a JSON string
export interface LiteralString extends Node {
  readonly type:        "LiteralStringNode"
  readonly value:       string
  readonly kind:        LiteralStringKind
  readonly blockIndent: string
}

export const isLiteralString = (node: Node): node is LiteralString => {
  const nodeType: NodeKind = "LiteralStringNode";
  return node.type === nodeType;
}

export class LiteralStringSingle extends NodeBase implements LiteralString {
  readonly type:        "LiteralStringNode" = "LiteralStringNode";
  readonly kind:        "StringSingle"      = "StringSingle";
  readonly blockIndent: ""                  = "";

  constructor(
    readonly value:       string,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `'${this.value}'`;
  }
}

export const isLiteralStringSingle = (node): node is LiteralStringSingle => {
  return node instanceof LiteralStringSingle;
}

export class LiteralStringDouble extends NodeBase implements LiteralString {
  readonly type:        "LiteralStringNode" = "LiteralStringNode";
  readonly kind:        "StringDouble"      = "StringDouble";
  readonly blockIndent: ""                  = "";

  constructor(
    readonly value:       string,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `"${this.value}"`;
  }
}

export const isLiteralStringDouble = (node): node is LiteralString => {
  return node instanceof LiteralStringDouble;
}

export class LiteralStringBlock extends NodeBase implements LiteralString {
  readonly type: "LiteralStringNode" = "LiteralStringNode";
  readonly kind: "StringBlock"       = "StringBlock";

  constructor(
    readonly value:       string,
    readonly blockIndent: string,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `|||${this.value}|||`;
  }
}

export const isLiteralStringBlock = (node): node is LiteralStringBlock => {
  return node instanceof LiteralStringBlock;
}

// ---------------------------------------------------------------------------

export type ObjectFieldKind =
  "ObjectAssert" |    // assert expr2 [: expr3]  where expr3 can be nil
  "ObjectFieldID" |   // id:[:[:]] expr2
  "ObjectFieldExpr" | // '['expr1']':[:[:]] expr2
  "ObjectFieldStr" |  // expr1:[:[:]] expr2
  "ObjectLocal";      // local id = expr2

export type ObjectFieldHide =
  "ObjectFieldHidden" |  // f:: e
  "ObjectFieldInherit" | // f: e
  "ObjectFieldVisible";  // f::: e

const objectFieldHideStrings = im.Map<ObjectFieldHide, string>({
  "ObjectFieldHidden": "::",
  "ObjectFieldInherit": ":",
  "ObjectFieldVisible": ":::",
});

// export interface ObjectField extends NodeBase {
//   readonly type:            "ObjectFieldNode"
//   readonly kind:            ObjectFieldKind
//   readonly hide:            ObjectFieldHide // (ignore if kind != astObjectField*)
//   readonly superSugar:      boolean         // +:  (ignore if kind != astObjectField*)
//   readonly methodSugar:     boolean         // f(x, y, z): ...  (ignore if kind  == astObjectAssert)
//   readonly expr1:           Node | null     // Not in scope of the object
//   readonly id:              Identifier | null
//   readonly ids:             FunctionParams  // If methodSugar == true then holds the params.
//   readonly trailingComma:   boolean         // If methodSugar == true then remembers the trailing comma
//   readonly expr2:           Node | null     // In scope of the object (can see self).
//   readonly expr3:           Node | null     // In scope of the object (can see self).
//   readonly headingComments: Comments
// }

export class ObjectField extends NodeBase {
  readonly type: "ObjectFieldNode" = "ObjectFieldNode";

  constructor(
    readonly kind:            ObjectFieldKind,
    readonly hide:            ObjectFieldHide, // (ignore if kind != astObjectField*)
    readonly superSugar:      boolean,         // +:  (ignore if kind != astObjectField*)
    readonly methodSugar:     boolean,         // f(x, y, z): ...  (ignore if kind  == astObjectAssert)
    readonly expr1:           Node | null,     // Not in scope of the object
    readonly id:              Identifier | null,
    readonly ids:             FunctionParams,  // If methodSugar == true then holds the params.
    readonly trailingComma:   boolean,         // If methodSugar == true then remembers the trailing comma
    readonly expr2:           Node | null,     // In scope of the object (can see self).
    readonly expr3:           Node | null,     // In scope of the object (can see self).
    readonly headingComments: BindingComment,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    switch (this.kind) {
      case "ObjectAssert": return prettyPrintObjectAssert(this);
      case "ObjectFieldID": return prettyPrintObjectFieldId(this);
      case "ObjectLocal": return prettyPrintObjectLocal(this);
      case "ObjectFieldExpr":
      case "ObjectFieldStr":
      default: throw new Error(`INTERNAL ERROR: Unrecognized object field kind '${this.kind}':\n${renderAsJson(this)}`);
    }
  }
}

export const isObjectField = (node): node is ObjectField => {
  return node instanceof ObjectField;
}

const prettyPrintObjectAssert = (field: ObjectField): string => {
    if (field.expr2 == null) {
      throw new Error(`INTERNAL ERROR: object 'assert' must have expression to assert:\n${renderAsJson(field)}`);
    }
    return field.expr3 == null
      ? `assert ${field.expr2.prettyPrint()}`
      : `assert ${field.expr2.prettyPrint()} : ${field.expr3.prettyPrint()}`;
}

const prettyPrintObjectFieldId = (field: ObjectField): string => {
  if (field.id == null) {
    throw new Error(`INTERNAL ERROR: object field must have id:\n${renderAsJson(field)}`);
  }
  const idString = field.id.prettyPrint();
  const hide = objectFieldHideStrings.get(field.hide);

  if (field.methodSugar) {
    const argsList = field.ids
      .map((param: FunctionParam) => param.id)
      .join(", ");
    return `(method) ${idString}(${argsList})${hide}`;
  }
  return `(field) ${idString}${hide}`;
}

const prettyPrintObjectLocal = (field: ObjectField): string => {
  if (field.id == null) {
    throw new Error(`INTERNAL ERROR: object field must have id:\n${renderAsJson(field)}`);
  }
  const idString = field.id.prettyPrint();

  if (field.methodSugar) {
    const argsList = field.ids
      .map((param: FunctionParam) => param.id)
      .join(", ");
    return `(method) local ${idString}(${argsList})`;
  }
  return `(field) local ${idString}`;
}

// TODO(jbeda): Add the remaining constructor helpers here

export type ObjectFields = im.List<ObjectField>;
export type IndexedObjectFields = im.Map<string, ObjectField>;

// NOTE: Type parameters are erased at runtime, so we can't check them
// here.
export const isIndexedObjectFields = (thing): thing is IndexedObjectFields => {
  return im.Map.isMap(thing);
}

export const indexFields = (fields: ObjectFields): IndexedObjectFields => {
  return fields
    .reduce((
      acc: im.Map<string, ObjectField>, field: ObjectField
    ) => {
      return field.id != null && acc.set(field.id.name, field) || acc;
    },
    im.Map<string, ObjectField>()
  );
}

// ---------------------------------------------------------------------------

// Object represents object constructors { f: e ... }.
//
// The trailing comma is only allowed if len(fields) > 0.  Converted to
// DesugaredObject during desugaring.
export class ObjectNode extends NodeBase implements FieldsResolvable {
  readonly type: "ObjectNode" = "ObjectNode";

  constructor(
    readonly fields:          ObjectFields,
    readonly trailingComma:   boolean,
    readonly headingComments: BindingComment,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    const fields = this.fields
      .filter((field: ObjectField) => field.kind === "ObjectFieldID")
      .map((field: ObjectField) => `  ${field.prettyPrint()}`)
      .join(",\n");

    return `(module) {\n${fields}\n}`;
  }

  public resolveFields = (
    context: ResolutionContext,
  ): Resolve | ResolveFailure => {
    return new Resolve(context.currFile, indexFields(this.fields));
  }
}

export const isObjectNode = (node): node is ObjectNode => {
  return node instanceof ObjectNode;
}

// ---------------------------------------------------------------------------

export interface DesugaredObjectField extends NodeBase {
  readonly type: "DesugaredObjectFieldNode"
  readonly hide: ObjectFieldHide
  readonly name: Node
  readonly body: Node
}
export type DesugaredObjectFields = im.List<DesugaredObjectField>;

// DesugaredObject represents object constructors { f: e ... } after
// desugaring.
//
// The assertions either return true or raise an error.
export interface DesugaredObject extends NodeBase {
  readonly type:    "DesugaredObjectNode"
  readonly asserts: Nodes
  readonly fields:  DesugaredObjectFields
}

export const isDesugaredObject = (node: Node): node is DesugaredObject => {
  const nodeType: NodeKind = "DesugaredObjectNode";
  return node.type === nodeType;
}

// ---------------------------------------------------------------------------

// ObjectComp represents object comprehension
//   { [e]: e for x in e for.. if... }.
// export interface ObjectComp extends NodeBase {
//   readonly type: "ObjectCompNode"
//   readonly fields:        ObjectFields
//   readonly trailingComma: boolean
//   readonly specs:         CompSpecs
// }

export class ObjectComp extends NodeBase {
  readonly type: "ObjectCompNode" = "ObjectCompNode";

  constructor(
    readonly fields:        ObjectFields,
    readonly trailingComma: boolean,
    readonly specs:         CompSpecs,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `[OBJECT COMP]`
  }
}

export const isObjectComp = (node): node is ObjectComp => {
  return node instanceof ObjectComp;
}

// ---------------------------------------------------------------------------

// ObjectComprehensionSimple represents post-desugaring object
// comprehension { [e]: e for x in e }.
//
// TODO: Rename this to `ObjectCompSimple`
export interface ObjectComprehensionSimple extends NodeBase {
  readonly type: "ObjectComprehensionSimpleNode"
  readonly field: Node
  readonly value: Node
  readonly id:    Identifier
  readonly array: Node
}

export const isObjectComprehensionSimple = (
  node: Node
): node is ObjectComprehensionSimple => {
  const nodeType: NodeKind = "ObjectComprehensionSimpleNode";
  return node.type === nodeType;
}

// ---------------------------------------------------------------------------

// Self represents the self keyword.
export class Self extends NodeBase implements Resolvable {
  readonly type: "SelfNode" = "SelfNode";

  constructor(
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `self`;
  }

  public resolve = (context: ResolutionContext): Resolve | ResolveFailure => {
    let curr: Node | null = this;
    while (true) {
      if (curr == null || curr.parent == null) {
        return Unresolved.Instance
      }

      if (isObjectNode(curr)) {
        return curr.resolveFields(context);
      }
      curr = curr.parent;
    }
  }
};

export const isSelf = (node): node is Self => {
  return node instanceof Self;
}

// ---------------------------------------------------------------------------

// SuperIndex represents the super[e] and super.f constructs.
//
// Either index or identifier will be set before desugaring.  After desugaring, id will be
// nil.
export class SuperIndex extends NodeBase {
  readonly type: "SuperIndexNode" = "SuperIndexNode";

  constructor(
    readonly index: Node | null,
    readonly id:    Identifier | null,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    if (this.id != null) {
      return `super.${this.id.prettyPrint()}`;
    } else if (this.index != null) {
      return `super[${this.index.prettyPrint()}]`
    }
    throw new Error(`INTERNAL ERROR: Can't pretty-print super index if both 'id' and 'index' fields are null`);
  }
}

export const isSuperIndex = (node): node is SuperIndex => {
  return node instanceof SuperIndex;
}

// ---------------------------------------------------------------------------

export type UnaryOp =
  "UopNot" |
  "UopBitwiseNot" |
  "UopPlus" |
  "UopMinus";

export const UopStrings = {
  UopNot:        "!",
  UopBitwiseNot: "~",
  UopPlus:       "+",
  UopMinus:      "-",
};

export const UopMap = im.Map<string, UnaryOp>({
  "!": "UopNot",
  "~": "UopBitwiseNot",
  "+": "UopPlus",
  "-": "UopMinus",
});

// Unary represents unary operators.
export class Unary extends NodeBase {
  readonly type: "UnaryNode" = "UnaryNode";

  constructor(
    readonly op:   UnaryOp,
    readonly expr: Node,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return `${UopStrings[this.op]}${this.expr.prettyPrint()}`;
  }
}

export const isUnary = (node): node is Unary => {
  return node instanceof Unary;
}

// ---------------------------------------------------------------------------

// Var represents variables.
export class Var extends NodeBase implements Resolvable {
  readonly type: "VarNode" = "VarNode";

  constructor(
    readonly id: Identifier,
    readonly loc:  lexical.LocationRange,
  ) { super(); }

  public prettyPrint = (): string => {
    return this.id.prettyPrint();
  }

  public resolve = (context: ResolutionContext): Resolve | ResolveFailure => {
    // Look up in the environment, get docs for that definition.
    if (this.env == null) {
      throw new Error(
        `INTERNAL ERROR: AST improperly set up, property 'env' can't be null:\n${renderAsJson(this)}`);
    } else if (!this.env.has(this.id.name)) {
      return Unresolved.Instance;
    }

    return resolveFromEnv(this.id.name, this.env, context);
  }
}

export const isVar = (node): node is Var => {
  return node instanceof Var;
}

// ---------------------------------------------------------------------------

const resolveFromEnv = (
  idName: string, env: Environment, context: ResolutionContext,
): Resolve | ResolveFailure => {
  const bind = env.get(idName);
  if (bind == null) {
    return Unresolved.Instance;
  }

  if (isFunctionParam(bind)) {
    // A function param is either a free variable, or it has a default
    // value. We consider both of these to be free variables, since we
    // would not know the value until the function was applied.
    return new ResolvedFreeVar(bind);
  }

  if (bind.body == null) {
    throw new Error(`INTERNAL ERROR: Bind can't have null body:\n${bind}`);
  }

  return tryResolve(bind.body, context);
}

const tryResolve = (
  node: Node, context: ResolutionContext,
): Resolve | ResolveFailure => {
  if (isFunction(node)) {
    return new ResolvedFunction(node);
  } else if (isResolvable(node)) {
    return node.resolve(context);
  } else if (isFieldsResolvable(node)) {
    // Found an object or perhaps an object mixin. Break.
    return node.resolveFields(context);
  } else if (isValueType(node)) {
    return new Resolve(context.currFile, node);
  } else {
    return Unresolved.Instance;
  }
}


export const tryResolveIndirections = (
  node: Node, context: ResolutionContext,
): Resolve | ResolveFailure => {
  // This loop will try to "strip out the indirections" of an
  // argument to a mixin. For example, consider the expression
  // `foo1.bar + foo2.bar` in the following example:
  //
  //   local bar1 = {a: 1, b: 2},
  //   local bar2 = {b: 3, c: 4},
  //   local foo1 = {bar: bar1},
  //   local foo2 = {bar: bar2},
  //   useMerged: foo1.bar + foo2.bar,
  //
  // In this case, if we try to resolve `foo1.bar + foo2.bar`, we
  // will first need to resolve `foo1.bar`, and then the value of
  // that resolve, `bar1`, which resolves to an object, and so on.
  //
  // This loop follows these indirections: first, it resolves
  // `foo1.bar`, and then `bar1`, before encountering an object
  // and stopping.

  let resolved: Resolve | ResolveFailure = new Resolve(context.currFile, node);
  while (true) {
    if (isResolveFailure(resolved)) {
      return resolved;
    } else if (isIndexedObjectFields(resolved.value)) {
      // We've resolved to a set of fields. Return.
      return resolved;
    } else if (isResolvable(resolved.value)) {
      resolved = resolved.value.resolve(context.withUri(resolved.fileUri));
    } else if (isFieldsResolvable(resolved.value)) {
      resolved = resolved.value.resolveFields(context.withUri(resolved.fileUri));
    } else if (isTypeGuessResolvable(resolved.value)) {
      resolved = resolved.value.resolveTypeGuess(context.withUri(resolved.fileUri));
    } else if (isValueType(resolved.value)) {
      // We've resolved to a value. Return.
      return resolved;
    } else {
      return Unresolved.Instance;
    }
  }
}

// ---------------------------------------------------------------------------
// Failures.
// ---------------------------------------------------------------------------

// ResolveFailure represents a failure to resolve a symbol to a "value
// type", for our particular definition of that term, which is
// captured by `isValueType`.
//
// For example, a symbol might refer to a function, which we would not
// consider a "value type", and hence we would return a
// `ResolveFailure`.
export type ResolveFailure =
  ResolvedFunction | ResolvedFreeVar |        // Resolved to uncompletable nodes.
  UnresolvedIndexId | UnresolvedIndexTarget | // Failed to resolve `Index` node.
  Unresolved;                                 // Misc.

export const isResolveFailure = (thing): thing is ResolveFailure => {
  return thing instanceof ResolvedFunction ||
    thing instanceof ResolvedFreeVar ||
    thing instanceof UnresolvedIndexId ||
    thing instanceof UnresolvedIndexTarget ||
    thing instanceof Unresolved;
}

// ResolvedFunction represents the event that we have tried to resolve
// a symbol to a "value type" (as defined by `isValueType`), but
// failed since that value depends on the resolution of a function,
// which cannot be resolved to a value type without binding the
// parameters.
export class ResolvedFunction {
  constructor(
    public readonly functionNode: Function | ObjectField | LocalBind
  ) {}
};

export const isResolvedFunction = (thing): thing is ResolvedFunction => {
  return thing instanceof ResolvedFunction;
}

// ResolvedFreeVar represents the event that we have tried to resolve
// a value to a "value type" (as defined by `isValueType`), but failed
// since that value is a free parameter, and must be bound at runtime
// to be computed.
//
// A good example of such a situation is `self`, `super`, and
// function parameters.
export class ResolvedFreeVar {
  constructor(public readonly variable: Var | FunctionParam) {}
};

export const isResolvedFreeVar = (thing): thing is ResolvedFreeVar => {
  return thing instanceof ResolvedFreeVar;
}

export type UnresolvedIndex = UnresolvedIndexTarget | UnresolvedIndexId;

export const isUnresolvedIndex = (thing): thing is UnresolvedIndex => {
  return thing instanceof UnresolvedIndexTarget ||
    thing instanceof UnresolvedIndexId;
}

// UnresolvedIndexTarget represents a failure to resolve an `Index`
// node because the target has failed to resolve.
//
// For example, in `foo.bar.baz`, failure to resolve either `foo` or
// `bar`, would result in an `UnresolvedIndexTarget`.
//
// NOTE: If `bar` fails to resolve, then we will still report an
// `UnresolvedIndexTarget`, since `bar` is the target of `bar.baz`.
export class UnresolvedIndexTarget {
  constructor(
    public readonly index: Index,
  ) {}
}

export const isUnresolvedIndexTarget = (thing): thing is UnresolvedIndexTarget => {
  return thing instanceof UnresolvedIndexTarget;
}

// UnresolvedIndexId represents a failure to resolve the ID of an
// `Index` node.
//
// For example, in `foo.bar.baz`, `baz` is the ID, hence failing to
// resolve `baz` will result in this error.
//
// NOTE: Only `baz` can cause an `UnresolvedIndexId` failure in this
// example. The reason failing to resolve `bar` doesn't cause an
// `UnresolvedIndexId` is because `bar` is the target in `bar.baz`.
export class UnresolvedIndexId {
  constructor(
    public readonly index: Index,
    public readonly resolvedTarget: IndexedObjectFields,
  ) {}
}

export const isUnresolvedIndexId = (thing): thing is UnresolvedIndexId => {
  return thing instanceof UnresolvedIndexId;
}

// Unresolved represents a miscelleneous failure to resolve a symbol.
// Typically this occurs the structure of the AST is not amenable to
// static analysis, and we simply punt.
//
// TODO: Expand this to more cases as `onComplete` features require it.
export class Unresolved {
  public static readonly Instance = new Unresolved();

  // NOTE: This is a work around for a bug in the TypeScript type
  // checker. We have not had time to report this bug, but when this
  // line is commented out, then use of `isResolveFailure` will cause
  // the type we're checking to resolve to `never` (TypeScript's
  // bottom type), which causes compile to fail.
  private readonly foo = "foo";
}

export const isUnresolved = (thing): thing is Unresolved => {
  return thing instanceof Unresolved;
}
