"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const os = require("os");
const im = require("immutable");
const _static = require("../../static");
exports.emptyEnvironment = im.Map();
exports.envFromLocalBinds = (local) => {
    if (exports.isLocal(local)) {
        const defaultLocal = {};
        const binds = local.binds
            .reduce((acc, bind) => {
            acc[bind.variable.name] = bind;
            return acc;
        }, defaultLocal);
        return im.Map(binds);
    }
    else if (exports.isObjectField(local)) {
        if (local.expr2 == null || local.id == null) {
            throw new Error(`INTERNAL ERROR: Object local fields can't have a null expr2 or id field`);
        }
        const bind = new LocalBind(local.id, local.expr2, local.methodSugar, local.ids, local.trailingComma, local.loc);
        return im.Map().set(local.id.name, bind);
    }
    // Else, it's a `FunctionParam`, i.e., a free parameter (or a free
    // parameter with a default value). Either way, emit that.
    return im.Map().set(local.id, local);
};
exports.envFromParams = (params) => {
    return params
        .reduce((acc, field) => {
        return acc.merge(exports.envFromLocalBinds(field));
    }, exports.emptyEnvironment);
};
exports.envFromFields = (fields) => {
    return fields
        .filter((field) => {
        const localKind = "ObjectLocal";
        return field.kind === localKind;
    })
        .reduce((acc, field) => {
        return acc.merge(exports.envFromLocalBinds(field));
    }, exports.emptyEnvironment);
};
exports.renderAsJson = (node) => {
    return "```\n" + JSON.stringify(node, (k, v) => {
        if (k === "parent") {
            return v == null
                ? "null"
                : v.type;
        }
        else if (k === "env") {
            return v == null
                ? "null"
                : `${Object.keys(v).join(", ")}`;
        }
        else if (k === "rootObject") {
            return v == null
                ? "null"
                : v.type;
        }
        else {
            return v;
        }
    }, "  ") + "\n```";
};
// isValueType returns true if the node is a computed value literal
// (e.g., a string literal, an object literal, and so on).
//
// Notably, this explicitly omits structures whose value must be
// computed at runtime: particularly object comprehension, array
// comprehension, self, super, and function types (whose value depends
// on parameter binds).
exports.isValueType = (node) => {
    // TODO(hausdorff): Consider adding object comprehension here, too.
    return exports.isLiteralBoolean(node) || exports.isLiteralNull(node) ||
        exports.isLiteralNumber(node) || exports.isLiteralString(node) || exports.isObjectNode(node);
};
// NodeBase is a simple abstract base class that makes sure we're
// initializing the parent and env members to null. It is not exposed
// to the public because it is meant to be a transparent base blass
// for all `Node` implementations.
class NodeBase {
    constructor() {
        this.rootObject = null;
        this.parent = null;
        this.env = null;
    }
}
exports.isNode = (thing) => {
    // TODO: Probably want to check the types of the properties instead.
    return thing instanceof NodeBase;
};
// ---------------------------------------------------------------------------
// Resolve represents a resolved node, including a fully-qualified RFC
// 1630/1738-compliant URI representing the absolute location of the
// Jsonnet file the symbol occurs in.
class Resolve {
    constructor(fileUri, value) {
        this.fileUri = fileUri;
        this.value = value;
    }
}
exports.Resolve = Resolve;
exports.isResolve = (thing) => {
    return thing instanceof Resolve;
};
// ResolutionContext represents the context we carry along as we
// attempt to resolve symbols. For example, an `import` node will have
// a filename, and to locate it, we will need to (1) search for the
// import path relative to the current path, or (2) look in the
// `libPaths` for it if necessary. This "context" is carried along in
// this object.
class ResolutionContext {
    constructor(compiler, documents, currFile) {
        this.compiler = compiler;
        this.documents = documents;
        this.currFile = currFile;
        this.withUri = (currFile) => {
            return new ResolutionContext(this.compiler, this.documents, currFile);
        };
    }
}
exports.ResolutionContext = ResolutionContext;
exports.isResolvable = (node) => {
    return node instanceof NodeBase && typeof node["resolve"] === "function";
};
exports.isFieldsResolvable = (node) => {
    return node instanceof NodeBase &&
        typeof node["resolveFields"] === "function";
};
exports.isTypeGuessResolvable = (node) => {
    return node instanceof NodeBase &&
        typeof node["resolveTypeGuess"] === "function";
};
class Identifier extends NodeBase {
    constructor(name, loc) {
        super();
        this.name = name;
        this.loc = loc;
        this.type = "IdentifierNode";
        this.prettyPrint = () => {
            return this.name;
        };
        this.resolve = (context) => {
            if (this.parent == null) {
                // An identifier with no parent is not a valid Jsonnet file.
                return Unresolved.Instance;
            }
            return tryResolve(this.parent, context);
        };
    }
}
exports.Identifier = Identifier;
exports.isIdentifier = (node) => {
    return node instanceof Identifier;
};
;
exports.isBindingComment = (node) => {
    return exports.isCppComment(node) || exports.isCComment(node);
};
exports.isComment = (node) => {
    const nodeType = "CommentNode";
    return node.type === nodeType;
};
class CppComment extends NodeBase {
    constructor(text, loc) {
        super();
        this.text = text;
        this.loc = loc;
        this.type = "CommentNode";
        this.kind = "CppStyle";
        this.prettyPrint = () => {
            return this.text.join(os.EOL);
        };
    }
}
exports.CppComment = CppComment;
exports.isCppComment = (node) => {
    return node instanceof CppComment;
};
class CComment extends NodeBase {
    constructor(text, loc) {
        super();
        this.text = text;
        this.loc = loc;
        this.type = "CommentNode";
        this.kind = "CStyle";
        this.prettyPrint = () => {
            return this.text.join(os.EOL);
        };
    }
}
exports.CComment = CComment;
exports.isCComment = (node) => {
    return node instanceof CComment;
};
class HashComment extends NodeBase {
    constructor(text, loc) {
        super();
        this.text = text;
        this.loc = loc;
        this.type = "CommentNode";
        this.kind = "HashStyle";
        this.prettyPrint = () => {
            return this.text.join(os.EOL);
        };
    }
}
exports.HashComment = HashComment;
exports.isHashComment = (node) => {
    return node instanceof HashComment;
};
;
exports.isCompSpec = (node) => {
    const nodeType = "CompSpecNode";
    return node.type === nodeType;
};
class CompSpecIf extends NodeBase {
    constructor(expr, loc) {
        super();
        this.expr = expr;
        this.loc = loc;
        this.type = "CompSpecNode";
        this.kind = "CompIf";
        this.varName = null; // null when kind != compSpecFor
        this.prettyPrint = () => {
            return `if ${this.expr.prettyPrint()}`;
        };
    }
}
exports.CompSpecIf = CompSpecIf;
exports.isCompSpecIf = (node) => {
    return node instanceof CompSpecIf;
};
class CompSpecFor extends NodeBase {
    constructor(varName, // null for `CompSpecIf`
    expr, loc) {
        super();
        this.varName = varName;
        this.expr = expr;
        this.loc = loc;
        this.type = "CompSpecNode";
        this.kind = "CompFor";
        this.prettyPrint = () => {
            return `for ${this.varName.prettyPrint()} in ${this.expr.prettyPrint()}`;
        };
    }
}
exports.CompSpecFor = CompSpecFor;
exports.isCompSpecFor = (node) => {
    return node instanceof CompSpecFor;
};
// ---------------------------------------------------------------------------
// Apply represents a function call
class Apply extends NodeBase {
    constructor(target, args, trailingComma, tailStrict, loc) {
        super();
        this.target = target;
        this.args = args;
        this.trailingComma = trailingComma;
        this.tailStrict = tailStrict;
        this.loc = loc;
        this.type = "ApplyNode";
        this.prettyPrint = () => {
            const argsString = this.args
                .map((arg) => arg.prettyPrint())
                .join(", ");
            // NOTE: Space between `tailstrict` is important.
            const tailStrictString = this.tailStrict
                ? " tailstrict"
                : "";
            return `${this.target.prettyPrint()}(${argsString}${tailStrictString})`;
        };
        this.resolveTypeGuess = (context) => {
            if (!exports.isResolvable(this.target)) {
                return Unresolved.Instance;
            }
            const fn = this.target.resolve(context);
            if (!exports.isResolvedFunction(fn) || !exports.isObjectField(fn.functionNode)) {
                return Unresolved.Instance;
            }
            const body = fn.functionNode.expr2;
            if (exports.isBinary(body) && body.op == "BopPlus" && exports.isSelf(body.left)) {
                return body.left.resolve(context);
            }
            return Unresolved.Instance;
        };
    }
}
exports.Apply = Apply;
exports.isApply = (node) => {
    return node instanceof Apply;
};
class ApplyParamAssignment extends NodeBase {
    constructor(id, right, loc) {
        super();
        this.id = id;
        this.right = right;
        this.loc = loc;
        this.type = "ApplyParamAssignmentNode";
        this.prettyPrint = () => {
            return `${this.id}=${this.right.prettyPrint()}`;
        };
    }
}
exports.ApplyParamAssignment = ApplyParamAssignment;
exports.isApplyParamAssignment = (node) => {
    return node instanceof ApplyParamAssignment;
};
// ---------------------------------------------------------------------------
// ApplyBrace represents e { }.  Desugared to e + { }.
class ApplyBrace extends NodeBase {
    constructor(left, right, loc) {
        super();
        this.left = left;
        this.right = right;
        this.loc = loc;
        this.type = "ApplyBraceNode";
        this.prettyPrint = () => {
            return `${this.left.prettyPrint()} ${this.right.prettyPrint()}`;
        };
    }
}
exports.ApplyBrace = ApplyBrace;
exports.isApplyBrace = (node) => {
    return node instanceof ApplyBrace;
};
// ---------------------------------------------------------------------------
// Array represents array constructors [1, 2, 3].
class Array extends NodeBase {
    constructor(elements, trailingComma, headingComment, trailingComment, loc) {
        super();
        this.elements = elements;
        this.trailingComma = trailingComma;
        this.headingComment = headingComment;
        this.trailingComment = trailingComment;
        this.loc = loc;
        this.type = "ArrayNode";
        this.prettyPrint = () => {
            const elementsString = this.elements
                .map((element) => element.prettyPrint())
                .join(", ");
            return `[${elementsString}]`;
        };
    }
}
exports.Array = Array;
exports.isArray = (node) => {
    return node instanceof Array;
};
// ---------------------------------------------------------------------------
// ArrayComp represents array comprehensions (which are like Python list
// comprehensions)
class ArrayComp extends NodeBase {
    constructor(body, trailingComma, specs, loc) {
        super();
        this.body = body;
        this.trailingComma = trailingComma;
        this.specs = specs;
        this.loc = loc;
        this.type = "ArrayCompNode";
        this.prettyPrint = () => {
            const specsString = this.specs
                .map((spec) => spec.prettyPrint())
                .join(", ");
            return `[${specsString} ${this.body.prettyPrint()}]`;
        };
    }
}
exports.ArrayComp = ArrayComp;
exports.isArrayComp = (node) => {
    return node instanceof ArrayComp;
};
// ---------------------------------------------------------------------------
// Assert represents an assert expression (not an object-level assert).
//
// After parsing, message can be nil indicating that no message was
// specified. This AST is elimiated by desugaring.
class Assert extends NodeBase {
    constructor(cond, message, rest, loc) {
        super();
        this.cond = cond;
        this.message = message;
        this.rest = rest;
        this.loc = loc;
        this.type = "AssertNode";
        this.prettyPrint = () => {
            return `assert ${this.cond.prettyPrint()}`;
        };
    }
}
exports.Assert = Assert;
exports.isAssert = (node) => {
    return node instanceof Assert;
};
const BopStrings = {
    BopMult: "*",
    BopDiv: "/",
    BopPercent: "%",
    BopPlus: "+",
    BopMinus: "-",
    BopShiftL: "<<",
    BopShiftR: ">>",
    BopGreater: ">",
    BopGreaterEq: ">=",
    BopLess: "<",
    BopLessEq: "<=",
    BopManifestEqual: "==",
    BopManifestUnequal: "!=",
    BopBitwiseAnd: "&",
    BopBitwiseXor: "^",
    BopBitwiseOr: "|",
    BopAnd: "&&",
    BopOr: "||",
};
exports.BopMap = im.Map({
    "*": "BopMult",
    "/": "BopDiv",
    "%": "BopPercent",
    "+": "BopPlus",
    "-": "BopMinus",
    "<<": "BopShiftL",
    ">>": "BopShiftR",
    ">": "BopGreater",
    ">=": "BopGreaterEq",
    "<": "BopLess",
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
class Binary extends NodeBase {
    constructor(left, op, right, loc) {
        super();
        this.left = left;
        this.op = op;
        this.right = right;
        this.loc = loc;
        this.type = "BinaryNode";
        this.prettyPrint = () => {
            const leftString = this.left.prettyPrint();
            const opString = BopStrings[this.op];
            const rightString = this.right.prettyPrint();
            return `${leftString} ${opString} ${rightString}`;
        };
        this.resolveFields = (context) => {
            // Recursively merge fields if it's another mixin; if it's an
            // object, return fields; else, no fields to return.
            if (this.op !== "BopPlus") {
                return Unresolved.Instance;
            }
            const left = exports.tryResolveIndirections(this.left, context);
            if (exports.isResolveFailure(left) || !exports.isIndexedObjectFields(left.value)) {
                return Unresolved.Instance;
            }
            const right = exports.tryResolveIndirections(this.right, context);
            if (exports.isResolveFailure(right) || !exports.isIndexedObjectFields(right.value)) {
                return Unresolved.Instance;
            }
            let merged = left.value;
            right.value.forEach((v, k) => {
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
        };
    }
}
exports.Binary = Binary;
exports.isBinary = (node) => {
    return node instanceof Binary;
};
// ---------------------------------------------------------------------------
// Builtin represents built-in functions.
//
// There is no parse rule to build this AST.  Instead, it is used to build the
// std object in the interpreter.
class Builtin extends NodeBase {
    constructor(id, params, loc) {
        super();
        this.id = id;
        this.params = params;
        this.loc = loc;
        this.type = "BuiltinNode";
        this.prettyPrint = () => {
            const paramsString = this.params.join(", ");
            return `std.${this.id}(${paramsString})`;
        };
    }
}
exports.Builtin = Builtin;
exports.isBuiltin = (node) => {
    return node instanceof Builtin;
};
// ---------------------------------------------------------------------------
// Conditional represents if/then/else.
//
// After parsing, branchFalse can be nil indicating that no else branch
// was specified.  The desugarer fills this in with a LiteralNull
class Conditional extends NodeBase {
    constructor(cond, branchTrue, branchFalse, loc) {
        super();
        this.cond = cond;
        this.branchTrue = branchTrue;
        this.branchFalse = branchFalse;
        this.loc = loc;
        this.type = "ConditionalNode";
        this.prettyPrint = () => {
            const trueClause = `then ${this.branchTrue.prettyPrint()}`;
            const falseClause = this.branchFalse == null
                ? ""
                : `else ${this.branchFalse.prettyPrint()}`;
            return `if ${this.cond.prettyPrint()} ${trueClause} ${falseClause}`;
        };
    }
}
exports.Conditional = Conditional;
exports.isConditional = (node) => {
    return node instanceof Conditional;
};
// ---------------------------------------------------------------------------
// Dollar represents the $ keyword
class Dollar extends NodeBase {
    constructor(loc) {
        super();
        this.loc = loc;
        this.type = "DollarNode";
        this.prettyPrint = () => {
            return `$`;
        };
        this.resolve = (context) => {
            if (this.rootObject == null) {
                return Unresolved.Instance;
            }
            return new Resolve(context.currFile, this.rootObject);
        };
    }
}
exports.Dollar = Dollar;
;
exports.isDollar = (node) => {
    return node instanceof Dollar;
};
// ---------------------------------------------------------------------------
// Error represents the error e.
class ErrorNode extends NodeBase {
    constructor(expr, loc) {
        super();
        this.expr = expr;
        this.loc = loc;
        this.type = "ErrorNode";
        this.prettyPrint = () => {
            return `error ${this.expr.prettyPrint()}`;
        };
    }
}
exports.ErrorNode = ErrorNode;
exports.isError = (node) => {
    return node instanceof ErrorNode;
};
// ---------------------------------------------------------------------------
// Function represents a function call. (jbeda: or is it function defn?)
class Function extends NodeBase {
    constructor(parameters, trailingComma, body, headingComment, trailingComment, loc) {
        super();
        this.parameters = parameters;
        this.trailingComma = trailingComma;
        this.body = body;
        this.headingComment = headingComment;
        this.trailingComment = trailingComment;
        this.loc = loc;
        this.type = "FunctionNode";
        this.prettyPrint = () => {
            const params = this.parameters
                .map((param) => param.prettyPrint())
                .join(", ");
            return `function (${params}) ${this.body.prettyPrint()}`;
        };
    }
}
exports.Function = Function;
exports.isFunction = (node) => {
    return node instanceof Function;
};
class FunctionParam extends NodeBase {
    constructor(id, defaultValue, loc) {
        super();
        this.id = id;
        this.defaultValue = defaultValue;
        this.loc = loc;
        this.type = "FunctionParamNode";
        this.prettyPrint = () => {
            const defaultValueString = this.defaultValue == null
                ? ""
                : `=${this.defaultValue.prettyPrint()}`;
            return `(parameter) ${this.id}${defaultValueString}`;
        };
    }
}
exports.FunctionParam = FunctionParam;
exports.isFunctionParam = (node) => {
    return node instanceof FunctionParam;
};
// ---------------------------------------------------------------------------
// Import represents import "file".
class Import extends NodeBase {
    constructor(file, loc) {
        super();
        this.file = file;
        this.loc = loc;
        this.type = "ImportNode";
        this.prettyPrint = () => {
            return `import "${this.file}"`;
        };
        this.resolve = (context) => {
            const { text: docText, version: version, resolvedPath: fileUri } = context.documents.get(this);
            const cached = context.compiler.cache(fileUri, docText, version);
            if (_static.isFailedParsedDocument(cached)) {
                return Unresolved.Instance;
            }
            let resolved = cached.parse;
            // If the var was pointing at an import, then resolution probably
            // has `local` definitions at the top of the file. Get rid of
            // them, since they are not useful for resolving the index
            // identifier.
            while (exports.isLocal(resolved)) {
                resolved = resolved.body;
            }
            return new Resolve(fileUri, resolved);
        };
    }
}
exports.Import = Import;
exports.isImport = (node) => {
    return node instanceof Import;
};
// ---------------------------------------------------------------------------
// ImportStr represents importstr "file".
class ImportStr extends NodeBase {
    constructor(file, loc) {
        super();
        this.file = file;
        this.loc = loc;
        this.type = "ImportStrNode";
        this.prettyPrint = () => {
            return `importstr "${this.file}"`;
        };
    }
}
exports.ImportStr = ImportStr;
exports.isImportStr = (node) => {
    return node instanceof ImportStr;
};
const resolveIndex = (index, context) => {
    if (index.target == null ||
        (!exports.isResolvable(index.target) && !exports.isFieldsResolvable(index.target) && !exports.isTypeGuessResolvable(index.target))) {
        throw new Error(`INTERNAL ERROR: Index node must have a resolvable target:\n${exports.renderAsJson(index)}`);
    }
    else if (index.id == null) {
        return Unresolved.Instance;
    }
    // Find root target, look up in environment.
    let resolvedTarget = exports.tryResolveIndirections(index.target, context);
    if (exports.isResolveFailure(resolvedTarget)) {
        return new UnresolvedIndexTarget(index);
    }
    else if (!exports.isIndexedObjectFields(resolvedTarget.value)) {
        return new UnresolvedIndexTarget(index);
    }
    const filtered = resolvedTarget.value.filter((field) => {
        return field.id != null && index.id != null &&
            field.id.name == index.id.name;
    });
    if (filtered.count() == 0) {
        return new UnresolvedIndexId(index, resolvedTarget.value);
    }
    else if (filtered.count() != 1) {
        throw new Error(`INTERNAL ERROR: Object contained multiple fields with name '${index.id.name}'}`);
    }
    const field = filtered.first();
    if (field.methodSugar) {
        return new ResolvedFunction(field);
    }
    else if (field.expr2 == null) {
        throw new Error(`INTERNAL ERROR: Object field can't have null property expr2:\n${exports.renderAsJson(field)}'}`);
    }
    return new Resolve(context.currFile, field.expr2);
};
exports.isIndex = (node) => {
    const nodeType = "IndexNode";
    return node.type === nodeType;
};
class IndexSubscript extends NodeBase {
    constructor(target, index, loc) {
        super();
        this.target = target;
        this.index = index;
        this.loc = loc;
        this.type = "IndexNode";
        this.id = null;
        this.prettyPrint = () => {
            return `${this.target.prettyPrint()}[${this.index.prettyPrint()}]`;
        };
        this.resolve = (context) => resolveIndex(this, context);
    }
}
exports.IndexSubscript = IndexSubscript;
exports.isIndexSubscript = (node) => {
    return node instanceof IndexSubscript;
};
class IndexDot extends NodeBase {
    constructor(target, id, loc) {
        super();
        this.target = target;
        this.id = id;
        this.loc = loc;
        this.type = "IndexNode";
        this.index = null;
        this.prettyPrint = () => {
            return `${this.target.prettyPrint()}.${this.id.prettyPrint()}`;
        };
        this.resolve = (context) => resolveIndex(this, context);
    }
}
exports.IndexDot = IndexDot;
exports.isIndexDot = (node) => {
    return node instanceof IndexDot;
};
// ---------------------------------------------------------------------------
// LocalBind is a helper struct for Local
class LocalBind extends NodeBase {
    constructor(variable, body, functionSugar, params, // if functionSugar is true
    trailingComma, loc) {
        super();
        this.variable = variable;
        this.body = body;
        this.functionSugar = functionSugar;
        this.params = params;
        this.trailingComma = trailingComma;
        this.loc = loc;
        this.type = "LocalBindNode";
        this.prettyPrint = () => {
            const idString = this.variable.prettyPrint();
            if (this.functionSugar) {
                const paramsString = this.params
                    .map((param) => param.id)
                    .join(", ");
                return `${idString}(${paramsString})`;
            }
            return `${idString} = ${this.body.prettyPrint()}`;
        };
    }
}
exports.LocalBind = LocalBind;
exports.isLocalBind = (node) => {
    return node instanceof LocalBind;
};
// Local represents local x = e; e.  After desugaring, functionSugar is false.
class Local extends NodeBase {
    constructor(binds, body, loc) {
        super();
        this.binds = binds;
        this.body = body;
        this.loc = loc;
        this.type = "LocalNode";
        this.prettyPrint = () => {
            const bindsString = this.binds
                .map((bind) => bind.prettyPrint())
                .join(",\n  ");
            return `local ${bindsString}`;
        };
    }
}
exports.Local = Local;
exports.isLocal = (node) => {
    return node instanceof Local;
};
// ---------------------------------------------------------------------------
// LiteralBoolean represents true and false
class LiteralBoolean extends NodeBase {
    constructor(value, loc) {
        super();
        this.value = value;
        this.loc = loc;
        this.type = "LiteralBooleanNode";
        this.prettyPrint = () => {
            return `${this.value}`;
        };
    }
}
exports.LiteralBoolean = LiteralBoolean;
exports.isLiteralBoolean = (node) => {
    return node instanceof LiteralBoolean;
};
// ---------------------------------------------------------------------------
// LiteralNull represents the null keyword
class LiteralNull extends NodeBase {
    constructor(loc) {
        super();
        this.loc = loc;
        this.type = "LiteralNullNode";
        this.prettyPrint = () => {
            return `null`;
        };
    }
}
exports.LiteralNull = LiteralNull;
exports.isLiteralNull = (node) => {
    return node instanceof LiteralNull;
};
// ---------------------------------------------------------------------------
// LiteralNumber represents a JSON number
class LiteralNumber extends NodeBase {
    constructor(value, originalString, loc) {
        super();
        this.value = value;
        this.originalString = originalString;
        this.loc = loc;
        this.type = "LiteralNumberNode";
        this.prettyPrint = () => {
            return `${this.originalString}`;
        };
    }
}
exports.LiteralNumber = LiteralNumber;
exports.isLiteralNumber = (node) => {
    return node instanceof LiteralNumber;
};
exports.isLiteralString = (node) => {
    const nodeType = "LiteralStringNode";
    return node.type === nodeType;
};
class LiteralStringSingle extends NodeBase {
    constructor(value, loc) {
        super();
        this.value = value;
        this.loc = loc;
        this.type = "LiteralStringNode";
        this.kind = "StringSingle";
        this.blockIndent = "";
        this.prettyPrint = () => {
            return `'${this.value}'`;
        };
    }
}
exports.LiteralStringSingle = LiteralStringSingle;
exports.isLiteralStringSingle = (node) => {
    return node instanceof LiteralStringSingle;
};
class LiteralStringDouble extends NodeBase {
    constructor(value, loc) {
        super();
        this.value = value;
        this.loc = loc;
        this.type = "LiteralStringNode";
        this.kind = "StringDouble";
        this.blockIndent = "";
        this.prettyPrint = () => {
            return `"${this.value}"`;
        };
    }
}
exports.LiteralStringDouble = LiteralStringDouble;
exports.isLiteralStringDouble = (node) => {
    return node instanceof LiteralStringDouble;
};
class LiteralStringBlock extends NodeBase {
    constructor(value, blockIndent, loc) {
        super();
        this.value = value;
        this.blockIndent = blockIndent;
        this.loc = loc;
        this.type = "LiteralStringNode";
        this.kind = "StringBlock";
        this.prettyPrint = () => {
            return `|||${this.value}|||`;
        };
    }
}
exports.LiteralStringBlock = LiteralStringBlock;
exports.isLiteralStringBlock = (node) => {
    return node instanceof LiteralStringBlock;
};
const objectFieldHideStrings = im.Map({
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
class ObjectField extends NodeBase {
    constructor(kind, hide, // (ignore if kind != astObjectField*)
    superSugar, // +:  (ignore if kind != astObjectField*)
    methodSugar, // f(x, y, z): ...  (ignore if kind  == astObjectAssert)
    expr1, // Not in scope of the object
    id, ids, // If methodSugar == true then holds the params.
    trailingComma, // If methodSugar == true then remembers the trailing comma
    expr2, // In scope of the object (can see self).
    expr3, // In scope of the object (can see self).
    headingComments, loc) {
        super();
        this.kind = kind;
        this.hide = hide;
        this.superSugar = superSugar;
        this.methodSugar = methodSugar;
        this.expr1 = expr1;
        this.id = id;
        this.ids = ids;
        this.trailingComma = trailingComma;
        this.expr2 = expr2;
        this.expr3 = expr3;
        this.headingComments = headingComments;
        this.loc = loc;
        this.type = "ObjectFieldNode";
        this.prettyPrint = () => {
            switch (this.kind) {
                case "ObjectAssert": return prettyPrintObjectAssert(this);
                case "ObjectFieldID": return prettyPrintObjectFieldId(this);
                case "ObjectLocal": return prettyPrintObjectLocal(this);
                case "ObjectFieldExpr":
                case "ObjectFieldStr":
                default: throw new Error(`INTERNAL ERROR: Unrecognized object field kind '${this.kind}':\n${exports.renderAsJson(this)}`);
            }
        };
    }
}
exports.ObjectField = ObjectField;
exports.isObjectField = (node) => {
    return node instanceof ObjectField;
};
const prettyPrintObjectAssert = (field) => {
    if (field.expr2 == null) {
        throw new Error(`INTERNAL ERROR: object 'assert' must have expression to assert:\n${exports.renderAsJson(field)}`);
    }
    return field.expr3 == null
        ? `assert ${field.expr2.prettyPrint()}`
        : `assert ${field.expr2.prettyPrint()} : ${field.expr3.prettyPrint()}`;
};
const prettyPrintObjectFieldId = (field) => {
    if (field.id == null) {
        throw new Error(`INTERNAL ERROR: object field must have id:\n${exports.renderAsJson(field)}`);
    }
    const idString = field.id.prettyPrint();
    const hide = objectFieldHideStrings.get(field.hide);
    if (field.methodSugar) {
        const argsList = field.ids
            .map((param) => param.id)
            .join(", ");
        return `(method) ${idString}(${argsList})${hide}`;
    }
    return `(field) ${idString}${hide}`;
};
const prettyPrintObjectLocal = (field) => {
    if (field.id == null) {
        throw new Error(`INTERNAL ERROR: object field must have id:\n${exports.renderAsJson(field)}`);
    }
    const idString = field.id.prettyPrint();
    if (field.methodSugar) {
        const argsList = field.ids
            .map((param) => param.id)
            .join(", ");
        return `(method) local ${idString}(${argsList})`;
    }
    return `(field) local ${idString}`;
};
// NOTE: Type parameters are erased at runtime, so we can't check them
// here.
exports.isIndexedObjectFields = (thing) => {
    return im.Map.isMap(thing);
};
exports.indexFields = (fields) => {
    return fields
        .reduce((acc, field) => {
        return field.id != null && acc.set(field.id.name, field) || acc;
    }, im.Map());
};
// ---------------------------------------------------------------------------
// Object represents object constructors { f: e ... }.
//
// The trailing comma is only allowed if len(fields) > 0.  Converted to
// DesugaredObject during desugaring.
class ObjectNode extends NodeBase {
    constructor(fields, trailingComma, headingComments, loc) {
        super();
        this.fields = fields;
        this.trailingComma = trailingComma;
        this.headingComments = headingComments;
        this.loc = loc;
        this.type = "ObjectNode";
        this.prettyPrint = () => {
            const fields = this.fields
                .filter((field) => field.kind === "ObjectFieldID")
                .map((field) => `  ${field.prettyPrint()}`)
                .join(",\n");
            return `(module) {\n${fields}\n}`;
        };
        this.resolveFields = (context) => {
            return new Resolve(context.currFile, exports.indexFields(this.fields));
        };
    }
}
exports.ObjectNode = ObjectNode;
exports.isObjectNode = (node) => {
    return node instanceof ObjectNode;
};
exports.isDesugaredObject = (node) => {
    const nodeType = "DesugaredObjectNode";
    return node.type === nodeType;
};
// ---------------------------------------------------------------------------
// ObjectComp represents object comprehension
//   { [e]: e for x in e for.. if... }.
// export interface ObjectComp extends NodeBase {
//   readonly type: "ObjectCompNode"
//   readonly fields:        ObjectFields
//   readonly trailingComma: boolean
//   readonly specs:         CompSpecs
// }
class ObjectComp extends NodeBase {
    constructor(fields, trailingComma, specs, loc) {
        super();
        this.fields = fields;
        this.trailingComma = trailingComma;
        this.specs = specs;
        this.loc = loc;
        this.type = "ObjectCompNode";
        this.prettyPrint = () => {
            return `[OBJECT COMP]`;
        };
    }
}
exports.ObjectComp = ObjectComp;
exports.isObjectComp = (node) => {
    return node instanceof ObjectComp;
};
exports.isObjectComprehensionSimple = (node) => {
    const nodeType = "ObjectComprehensionSimpleNode";
    return node.type === nodeType;
};
// ---------------------------------------------------------------------------
// Self represents the self keyword.
class Self extends NodeBase {
    constructor(loc) {
        super();
        this.loc = loc;
        this.type = "SelfNode";
        this.prettyPrint = () => {
            return `self`;
        };
        this.resolve = (context) => {
            let curr = this;
            while (true) {
                if (curr == null || curr.parent == null) {
                    return Unresolved.Instance;
                }
                if (exports.isObjectNode(curr)) {
                    return curr.resolveFields(context);
                }
                curr = curr.parent;
            }
        };
    }
}
exports.Self = Self;
;
exports.isSelf = (node) => {
    return node instanceof Self;
};
// ---------------------------------------------------------------------------
// SuperIndex represents the super[e] and super.f constructs.
//
// Either index or identifier will be set before desugaring.  After desugaring, id will be
// nil.
class SuperIndex extends NodeBase {
    constructor(index, id, loc) {
        super();
        this.index = index;
        this.id = id;
        this.loc = loc;
        this.type = "SuperIndexNode";
        this.prettyPrint = () => {
            if (this.id != null) {
                return `super.${this.id.prettyPrint()}`;
            }
            else if (this.index != null) {
                return `super[${this.index.prettyPrint()}]`;
            }
            throw new Error(`INTERNAL ERROR: Can't pretty-print super index if both 'id' and 'index' fields are null`);
        };
    }
}
exports.SuperIndex = SuperIndex;
exports.isSuperIndex = (node) => {
    return node instanceof SuperIndex;
};
exports.UopStrings = {
    UopNot: "!",
    UopBitwiseNot: "~",
    UopPlus: "+",
    UopMinus: "-",
};
exports.UopMap = im.Map({
    "!": "UopNot",
    "~": "UopBitwiseNot",
    "+": "UopPlus",
    "-": "UopMinus",
});
// Unary represents unary operators.
class Unary extends NodeBase {
    constructor(op, expr, loc) {
        super();
        this.op = op;
        this.expr = expr;
        this.loc = loc;
        this.type = "UnaryNode";
        this.prettyPrint = () => {
            return `${exports.UopStrings[this.op]}${this.expr.prettyPrint()}`;
        };
    }
}
exports.Unary = Unary;
exports.isUnary = (node) => {
    return node instanceof Unary;
};
// ---------------------------------------------------------------------------
// Var represents variables.
class Var extends NodeBase {
    constructor(id, loc) {
        super();
        this.id = id;
        this.loc = loc;
        this.type = "VarNode";
        this.prettyPrint = () => {
            return this.id.prettyPrint();
        };
        this.resolve = (context) => {
            // Look up in the environment, get docs for that definition.
            if (this.env == null) {
                throw new Error(`INTERNAL ERROR: AST improperly set up, property 'env' can't be null:\n${exports.renderAsJson(this)}`);
            }
            else if (!this.env.has(this.id.name)) {
                return Unresolved.Instance;
            }
            return resolveFromEnv(this.id.name, this.env, context);
        };
    }
}
exports.Var = Var;
exports.isVar = (node) => {
    return node instanceof Var;
};
// ---------------------------------------------------------------------------
const resolveFromEnv = (idName, env, context) => {
    const bind = env.get(idName);
    if (bind == null) {
        return Unresolved.Instance;
    }
    if (exports.isFunctionParam(bind)) {
        // A function param is either a free variable, or it has a default
        // value. We consider both of these to be free variables, since we
        // would not know the value until the function was applied.
        return new ResolvedFreeVar(bind);
    }
    if (bind.body == null) {
        throw new Error(`INTERNAL ERROR: Bind can't have null body:\n${bind}`);
    }
    return tryResolve(bind.body, context);
};
const tryResolve = (node, context) => {
    if (exports.isFunction(node)) {
        return new ResolvedFunction(node);
    }
    else if (exports.isResolvable(node)) {
        return node.resolve(context);
    }
    else if (exports.isFieldsResolvable(node)) {
        // Found an object or perhaps an object mixin. Break.
        return node.resolveFields(context);
    }
    else if (exports.isValueType(node)) {
        return new Resolve(context.currFile, node);
    }
    else {
        return Unresolved.Instance;
    }
};
exports.tryResolveIndirections = (node, context) => {
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
    let resolved = new Resolve(context.currFile, node);
    while (true) {
        if (exports.isResolveFailure(resolved)) {
            return resolved;
        }
        else if (exports.isIndexedObjectFields(resolved.value)) {
            // We've resolved to a set of fields. Return.
            return resolved;
        }
        else if (exports.isResolvable(resolved.value)) {
            resolved = resolved.value.resolve(context.withUri(resolved.fileUri));
        }
        else if (exports.isFieldsResolvable(resolved.value)) {
            resolved = resolved.value.resolveFields(context.withUri(resolved.fileUri));
        }
        else if (exports.isTypeGuessResolvable(resolved.value)) {
            resolved = resolved.value.resolveTypeGuess(context.withUri(resolved.fileUri));
        }
        else if (exports.isValueType(resolved.value)) {
            // We've resolved to a value. Return.
            return resolved;
        }
        else {
            return Unresolved.Instance;
        }
    }
};
exports.isResolveFailure = (thing) => {
    return thing instanceof ResolvedFunction ||
        thing instanceof ResolvedFreeVar ||
        thing instanceof UnresolvedIndexId ||
        thing instanceof UnresolvedIndexTarget ||
        thing instanceof Unresolved;
};
// ResolvedFunction represents the event that we have tried to resolve
// a symbol to a "value type" (as defined by `isValueType`), but
// failed since that value depends on the resolution of a function,
// which cannot be resolved to a value type without binding the
// parameters.
class ResolvedFunction {
    constructor(functionNode) {
        this.functionNode = functionNode;
    }
}
exports.ResolvedFunction = ResolvedFunction;
;
exports.isResolvedFunction = (thing) => {
    return thing instanceof ResolvedFunction;
};
// ResolvedFreeVar represents the event that we have tried to resolve
// a value to a "value type" (as defined by `isValueType`), but failed
// since that value is a free parameter, and must be bound at runtime
// to be computed.
//
// A good example of such a situation is `self`, `super`, and
// function parameters.
class ResolvedFreeVar {
    constructor(variable) {
        this.variable = variable;
    }
}
exports.ResolvedFreeVar = ResolvedFreeVar;
;
exports.isResolvedFreeVar = (thing) => {
    return thing instanceof ResolvedFreeVar;
};
exports.isUnresolvedIndex = (thing) => {
    return thing instanceof UnresolvedIndexTarget ||
        thing instanceof UnresolvedIndexId;
};
// UnresolvedIndexTarget represents a failure to resolve an `Index`
// node because the target has failed to resolve.
//
// For example, in `foo.bar.baz`, failure to resolve either `foo` or
// `bar`, would result in an `UnresolvedIndexTarget`.
//
// NOTE: If `bar` fails to resolve, then we will still report an
// `UnresolvedIndexTarget`, since `bar` is the target of `bar.baz`.
class UnresolvedIndexTarget {
    constructor(index) {
        this.index = index;
    }
}
exports.UnresolvedIndexTarget = UnresolvedIndexTarget;
exports.isUnresolvedIndexTarget = (thing) => {
    return thing instanceof UnresolvedIndexTarget;
};
// UnresolvedIndexId represents a failure to resolve the ID of an
// `Index` node.
//
// For example, in `foo.bar.baz`, `baz` is the ID, hence failing to
// resolve `baz` will result in this error.
//
// NOTE: Only `baz` can cause an `UnresolvedIndexId` failure in this
// example. The reason failing to resolve `bar` doesn't cause an
// `UnresolvedIndexId` is because `bar` is the target in `bar.baz`.
class UnresolvedIndexId {
    constructor(index, resolvedTarget) {
        this.index = index;
        this.resolvedTarget = resolvedTarget;
    }
}
exports.UnresolvedIndexId = UnresolvedIndexId;
exports.isUnresolvedIndexId = (thing) => {
    return thing instanceof UnresolvedIndexId;
};
// Unresolved represents a miscelleneous failure to resolve a symbol.
// Typically this occurs the structure of the AST is not amenable to
// static analysis, and we simply punt.
//
// TODO: Expand this to more cases as `onComplete` features require it.
class Unresolved {
    constructor() {
        // NOTE: This is a work around for a bug in the TypeScript type
        // checker. We have not had time to report this bug, but when this
        // line is commented out, then use of `isResolveFailure` will cause
        // the type we're checking to resolve to `never` (TypeScript's
        // bottom type), which causes compile to fail.
        this.foo = "foo";
    }
}
Unresolved.Instance = new Unresolved();
exports.Unresolved = Unresolved;
exports.isUnresolved = (thing) => {
    return thing instanceof Unresolved;
};
//# sourceMappingURL=tree.js.map