import * as os from 'os';

import * as im from 'immutable';

import * as ast from './ast';
import * as lexer from './lexer';
import * as lexical from './lexical';

// ---------------------------------------------------------------------------

type precedence = number;

const applyPrecedence: precedence = 2,  // Function calls and indexing.
      unaryPrecedence: precedence = 4,  // Logical and bitwise negation, unary + -
      maxPrecedence:   precedence = 16; // Local, If, Import, Function, Error

var bopPrecedence = im.Map<ast.BinaryOp, precedence>({
  "BopMult":            5,
  "BopDiv":             5,
  "BopPercent":         5,
  "BopPlus":            6,
  "BopMinus":           6,
  "BopShiftL":          7,
  "BopShiftR":          7,
  "BopGreater":         8,
  "BopGreaterEq":       8,
  "BopLess":            8,
  "BopLessEq":          8,
  "BopManifestEqual":   9,
  "BopManifestUnequal": 9,
  "BopBitwiseAnd":      10,
  "BopBitwiseXor":      11,
  "BopBitwiseOr":       12,
  "BopAnd":             13,
  "BopOr":              14,
});

// ---------------------------------------------------------------------------

const makeUnexpectedError = (
  t: lexer.Token, during: string
): lexical.StaticError => {
  return lexical.MakeStaticError(
    `Unexpected: ${t} while ${during}`,
    t.loc);
}

const locFromTokens = (
  begin: lexer.Token, end: lexer.Token
): lexical.LocationRange => {
  return lexical.MakeLocationRange(begin.loc.fileName, begin.loc.begin, end.loc.end)
}

const locFromTokenAST = (
  begin: lexer.Token, end: ast.Node
): lexical.LocationRange => {
  return lexical.MakeLocationRange(
    begin.loc.fileName, begin.loc.begin, end.loc.end)
}

// ---------------------------------------------------------------------------

class parser {
  private currT: number = 0;

  constructor(
    readonly t: lexer.Tokens
  ) {}

  public pop = (): lexer.Token => {
    const t = this.t.get(this.currT);
    this.currT++
    return t
  };

  public popExpect = (
    tk: lexer.TokenKind
  ): lexer.Token | lexical.StaticError => {
    const t = this.pop();
    if (t.kind !== tk) {
      return lexical.MakeStaticError(
        `Expected token ${lexer.TokenKindStrings.get(tk)} but got ${t}`,
        t.loc);
    }
    return t;
  };

  public popExpectOp = (
    op: string
  ): lexer.Token | lexical.StaticError => {
    const t = this.pop();
    if (t.kind !== "TokenOperator" || t.data != op) {
      return lexical.MakeStaticError(
        `Expected operator ${op} but got ${t}`, t.loc);
    }
    return t
  };

  public peek = (): lexer.Token => {
    return this.t.get(this.currT);
  };

  // parseOptionalComments parses a block of comments if they exist at
  // the current position in the token stream (as measured by
  // `this.peek()`), and has no effect if they don't.
  public parseOptionalComments = ()
  : ast.CppComment | ast.CComment | ast.HashComment | null => {
    const next = this.peek();
    switch (next.kind) {
      case "TokenCommentCpp": {
        return this.parseCppCommentBlock();
      }
      case "TokenCommentC": {
        return this.parseCComment();
      }
      case "TokenCommentHash": {
        return this.parseHashCommentBlock();
      }
      default: {
        return null;
      }
    }
  };

  public parseCppCommentBlock = ()
  : ast.CppComment | ast.CComment | ast.HashComment | null => {
    let lines = im.List<string>();
    const first = this.peek();
    let curr = this.peek();

    while (true) {
      curr = this.peek();
      switch (curr.kind) {
        case "TokenCommentCpp": {
          if (curr.fodder != null) {
            const anyNewlines = curr.fodder.filter(fodder =>
              fodder.data.match(/^\n\s*\n/) != null);
            if (anyNewlines.length) {
              curr = this.pop();
              lines = im.List<string>([curr.data]);
              break;
            }
          }

          curr = this.pop();
          lines = lines.push(curr.data);
          break;
        }
        case "TokenCommentC": {
          return this.parseCComment();
        }
        case "TokenCommentHash": {
          return this.parseHashCommentBlock();
        }
        default: {
          return lines.count() == 0
            ? null
            : new ast.CppComment(lines, locFromTokens(first, curr));
        }
      }
    }
  }

  public parseCComment = ()
  : ast.CppComment | ast.CComment | ast.HashComment | null => {
    let lines = im.List<string>();
    let next = this.peek();

    while (true) {
      next = this.peek();
      switch (next.kind) {
        case "TokenCommentCpp": {
          return this.parseCppCommentBlock();
        }
        case "TokenCommentC": {
          // NOTE: This does not trim the whitespace in the last line.
          // For example, if a multi-line comment ends: `   */`, then
          // the line will end with 3 space characters.
          const processedLines = next.data
            .split(os.EOL)
            .map(line => {
              const m = line.match(/^\s*\*/);
              if (m == null) {
                return line;
              }
              return line.slice(m[0].length);
            });

          next = this.pop();
          lines = im.List<string>(processedLines);
          break;
        }
        case "TokenCommentHash": {
          return this.parseHashCommentBlock();
        }
        default: {
          return lines.count() == 0
            ? null
            : new ast.CComment(lines, next.loc);
        }
      }
    }
  }

  public parseHashCommentBlock = ()
  : ast.CppComment | ast.CComment | ast.HashComment | null => {
    let lines = im.List<string>();
    const first = this.peek();
    let curr = this.peek();

    while (true) {
      curr = this.peek();
      switch (curr.kind) {
        case "TokenCommentCpp": {
          return this.parseCppCommentBlock();
        }
        case "TokenCommentC": {
          return this.parseCComment();
        }
        case "TokenCommentHash": {
          if (curr.fodder != null) {
            const anyNewlines = curr.fodder.filter(fodder =>
              fodder.data.match(/^\n\s*\n/) != null);
            if (anyNewlines.length) {
              curr = this.pop();
              lines = im.List<string>([curr.data]);
              break;
            }
          }

          curr = this.pop();
          lines = lines.push(curr.data);
          break;
        }
        default: {
          return lines.count() == 0
            ? null
            : new ast.HashComment(lines, locFromTokens(first, curr));
        }
      }
    }
  }

  public parseCommaList = <T extends ast.Node>(
    end: lexer.TokenKind, elementKind: string,
    elementCallback: (e: ast.Node) => T | lexical.StaticError = (e) => <T>e,
  ): {next: lexer.Token, exprs: im.List<T>, gotComma: boolean} | lexical.StaticError => {
    let exprs = im.List<T>();
    let gotComma = false;
    let first = true;
    while (true) {
      let next = this.peek();
      if (!first && !gotComma) {
        if (next.kind === "TokenComma") {
          this.pop();
          next = this.peek();
          gotComma = true;
        }
      }
      if (next.kind === end) {
        // gotComma can be true or false here.
        return {next: this.pop(), exprs: exprs, gotComma: gotComma};
      }

      if (!first && !gotComma) {
        return lexical.MakeStaticError(
          `Expected a comma before next ${elementKind}.`, next.loc);
      }

      const expr = this.parse(maxPrecedence, null);
      if (lexical.isStaticError(expr)) {
        return expr;
      }

      const mappedExpr = elementCallback(expr);
      if (lexical.isStaticError(mappedExpr)) {
        return mappedExpr;
      }
      exprs = exprs.push(mappedExpr);

      gotComma = false;
      first = false;
    }
  }

  public parseArgsList = (
    elementKind: string
  ): {next: lexer.Token, params: ast.Nodes, gotComma: boolean} | lexical.StaticError => {
    const result = this.parseCommaList<ast.Node>(
      "TokenParenR",
      elementKind,
      (expr): ast.Node | lexical.StaticError => {
        const next = this.peek();
        let rhs: ast.Node | null = null;
        if (ast.isVar(expr) && next.kind === "TokenOperator" &&
            next.data === "="
        ) {
          this.pop();
          const assignment = this.parse(maxPrecedence, null);
          if (lexical.isStaticError(assignment)) {
            return assignment;
          }
          return new ast.ApplyParamAssignment(
            expr.id.name, assignment, expr.loc);
        }

        return expr;
      });
    if (lexical.isStaticError(result)) {
      return result;
    }

    return {next: result.next, params: result.exprs, gotComma: result.gotComma};
  }

  public parseParamsList = (
    elementKind: string
  ): {next: lexer.Token, params: ast.FunctionParams, gotComma: boolean} | lexical.StaticError => {
    const result = this.parseCommaList<ast.FunctionParam>(
      "TokenParenR",
      elementKind,
      (expr): ast.FunctionParam | lexical.StaticError => {
        if (!ast.isVar(expr)) {
          return lexical.MakeStaticError(
            `Expected simple identifier but got a complex expression.`,
            expr.loc);
        }

        const next = this.peek();
        let rhs: ast.Node | null = null;
        if (next.kind === "TokenOperator" && next.data === "=") {
          this.pop();
          const assignment = this.parse(maxPrecedence, null);
          if (lexical.isStaticError(assignment)) {
            return assignment;
          }
          rhs = assignment;
        }
        return new ast.FunctionParam(expr.id.name, rhs, expr.loc);
      });
    if (lexical.isStaticError(result)) {
      return result;
    }

    return {next: result.next, params: result.exprs, gotComma: result.gotComma};
  }

  public parseBind = (
    localToken: lexer.Token, binds: ast.LocalBinds
  ): ast.LocalBinds | lexical.StaticError => {
    const varID = this.popExpect("TokenIdentifier");
    if (lexical.isStaticError(varID)) {
      return varID;
    }

    for (let b of binds.toArray()) {
      if (b.variable.name === varID.data) {
        return lexical.MakeStaticError(
          `Duplicate local var: ${varID.data}`, varID.loc);
      }
    }

    if (this.peek().kind === "TokenParenL") {
      this.pop();
      const result = this.parseParamsList("function parameter");
      if (lexical.isStaticError(result)) {
        return result;
      }

      const pop = this.popExpectOp("=")
      if (lexical.isStaticError(pop)) {
        return pop;
      }

      const body = this.parse(maxPrecedence, null);
      if (lexical.isStaticError(body)) {
        return body;
      }
      const id = new ast.Identifier(varID.data, varID.loc);
      const {params: params, gotComma: gotComma} = result;
      const location = locFromTokenAST(localToken, body);
      const bind = new ast.LocalBind(
        id, body, true, params, gotComma, location);
      binds = binds.push(bind);
    } else {
      const pop = this.popExpectOp("=");
      if (lexical.isStaticError(pop)) {
        return pop;
      }
      const body = this.parse(maxPrecedence, null);
      if (lexical.isStaticError(body)) {
        return body;
      }
      const id = new ast.Identifier(varID.data, varID.loc);
      const location = locFromTokenAST(localToken, body);
      const bind = new ast.LocalBind(
        id, body, false, im.List<ast.FunctionParam>(), false, location);
      binds = binds.push(bind);
    }

    return binds;
  };

  public parseObjectAssignmentOp = (
  ): {plusSugar: boolean, hide: ast.ObjectFieldHide} | lexical.StaticError => {
    let plusSugar = false;
    let hide: ast.ObjectFieldHide | null = null;

    const op = this.popExpect("TokenOperator");
    if (lexical.isStaticError(op)) {
      return op;
    }

    let opStr = op.data;
    if (opStr[0] === '+') {
      plusSugar = true;
      opStr = opStr.slice(1);
    }

    let numColons = 0
    while (opStr.length > 0) {
      if (opStr[0] !== ':') {
        return lexical.MakeStaticError(
          `Expected one of :, ::, :::, +:, +::, +:::, got: ${op.data}`,
          op.loc);
      }
      opStr = opStr.slice(1);
      numColons++
    }

    switch (numColons) {
      case 1:
        hide = "ObjectFieldInherit";
        break;
      case 2:
        hide = "ObjectFieldHidden"
        break;
      case 3:
        hide = "ObjectFieldVisible"
        break;
      default:
        return lexical.MakeStaticError(
          `Expected one of :, ::, :::, +:, +::, +:::, got: ${op.data}`,
          op.loc);
      }

    return {plusSugar: plusSugar, hide: hide};
  };

  // parseObjectCompRemainder parses the remainder of an object as an
  // object comprehension. This function is meant to act in conjunction
  // with `parseObjectCompRemainder`, and is meant to be called
  // immediately after the `for` token is encountered, since this is
  // typically the first indication that we are in an object
  // comprehension. Partially to enforce this condition, this function
  // takes as an argument the token representing the `for` keyword.
  public parseObjectCompRemainder = (
    first: lexer.Token, forTok: lexer.Token, gotComma: boolean,
    fields: ast.ObjectFields,
  ): {comp: ast.Node, last: lexer.Token} | lexical.StaticError => {
    let numFields = 0;
    let numAsserts = 0;
    let field = fields.first();
    for (field of fields.toArray()) {
      if (field.kind === "ObjectLocal") {
        continue;
      }
      if (field.kind === "ObjectAssert") {
        numAsserts++;
        continue;
      }
      numFields++;
    }

    if (numAsserts > 0) {
      return lexical.MakeStaticError(
        "Object comprehension cannot have asserts.", forTok.loc);
    }
    if (numFields != 1) {
      return lexical.MakeStaticError(
        "Object comprehension can only have one field.", forTok.loc);
    }
    if (field.hide != "ObjectFieldInherit") {
      return lexical.MakeStaticError(
        "Object comprehensions cannot have hidden fields.", forTok.loc);
    }
    if (field.kind !== "ObjectFieldExpr") {
      return lexical.MakeStaticError(
        "Object comprehensions can only have [e] fields.", forTok.loc);
    }
    const result = this.parseCompSpecs("TokenBraceR");
    if (lexical.isStaticError(result)) {
      return result;
    }

    const comp = new ast.ObjectComp(
      fields,
      gotComma,
      result.compSpecs,
      locFromTokens(first, result.maybeIf),
    );
    return {comp: comp, last: result.maybeIf};
  };

  // parseObjectField will parse a single field in an object.
  public parseObjectField = (
    headingComments: ast.BindingComment, next: lexer.Token,
    literalFields: im.Set<literalField>,
  ): {field: ast.ObjectField, literals: im.Set<literalField>} | lexical.StaticError => {
    let kind: ast.ObjectFieldKind;
    let expr1: ast.Node | null = null;
    let id: ast.Identifier | null = null;

    switch (next.kind) {
      case "TokenIdentifier": {
        kind = "ObjectFieldID";
        id = new ast.Identifier(next.data, next.loc);
        break;
      }
      case "TokenStringDouble": {
        kind = "ObjectFieldStr";
        expr1 = new ast.LiteralStringDouble(next.data, next.loc);
        break;
      }
      case "TokenStringSingle": {
        kind = "ObjectFieldStr";
        expr1 = new ast.LiteralStringSingle(next.data, next.loc);
        break;
      }
      case "TokenStringBlock": {
        kind = "ObjectFieldStr"
        expr1 = new ast.LiteralStringBlock(
          next.data, next.stringBlockIndent, next.loc);
        break;
      }
      default: {
        kind = "ObjectFieldExpr"
        const expr1 = this.parse(maxPrecedence, null);
        if (lexical.isStaticError(expr1)) {
          return expr1;
        }
        const pop = this.popExpect("TokenBracketR");
        if (lexical.isStaticError(pop)) {
          return pop;
        }
        break;
      }
    }

    let isMethod = false;
    let methComma = false;
    let params = im.List<ast.FunctionParam>();
    if (this.peek().kind === "TokenParenL") {
      this.pop();
      const result = this.parseParamsList("method parameter");
      if (lexical.isStaticError(result)) {
        return result;
      }
      params = result.params;
      isMethod = true
    }

    const result = this.parseObjectAssignmentOp();
    if (lexical.isStaticError(result)) {
      return result;
    }

    if (result.plusSugar && isMethod) {
      return lexical.MakeStaticError(
        `Cannot use +: syntax sugar in a method: ${next.data}`, next.loc);
    }

    if (kind !== "ObjectFieldExpr") {
      if (literalFields.contains(next.data)) {
        return lexical.MakeStaticError(
          `Duplicate field: ${next.data}`, next.loc);
      }
      literalFields = literalFields.add(next.data);
    }

    const body = this.parse(maxPrecedence, null);
    if (lexical.isStaticError(body)) {
      return body;
    }

    // TODO: The location range here is probably not quite correct.
    // For example, in cases where `body` is a string literal, the
    // location range will only reflect the string contents, not the
    // ending quote.
    return {
      field: new ast.ObjectField(
        kind,
        result.hide,
        result.plusSugar,
        isMethod,
        expr1,
        id,
        params,
        methComma,
        body,
        null,
        headingComments,
        locFromTokenAST(next, body),
      ),
      literals: literalFields
    };
  }

  // parseObjectLocal parses a `local` definition that appears in an
  // object, as an object field. `assertToken` is required to allow the
  // object to create an appropriate location range for the field.
  public parseObjectLocal = (
    localToken: lexer.Token, binds: ast.IdentifierSet,
  ): {field: ast.ObjectField, binds: ast.IdentifierSet} | lexical.StaticError => {
    const varID = this.popExpect("TokenIdentifier");
    if (lexical.isStaticError(varID)) {
      return varID;
    }
    const id = new ast.Identifier(varID.data, varID.loc);
    if (binds.contains(id.name)) {
      return lexical.MakeStaticError(
        `Duplicate local var: ${id.name}`, varID.loc);
    }

    let isMethod = false;
    let funcComma = false;
    let params = im.List<ast.FunctionParam>();
    if (this.peek().kind === "TokenParenL") {
      this.pop();
      const result = this.parseParamsList("function parameter");
      if (lexical.isStaticError(result)) {
        return result;
      }
      isMethod = true;
      params = result.params;
    }
    const pop = this.popExpectOp("=");
    if (lexical.isStaticError(pop)) {
      return pop;
    }

    const body = this.parse(maxPrecedence, null);
    if (lexical.isStaticError(body)) {
      return body;
    }

    binds = binds.add(id.name);

    return {
      field: new ast.ObjectField(
        "ObjectLocal",
        "ObjectFieldVisible",
        false,
        isMethod,
        null,
        id,
        params,
        funcComma,
        body,
        null,
        null,
        locFromTokenAST(localToken, body),
      ),
      binds: binds,
    };
  };

  // parseObjectAssert parses an `assert` that appears as an object
  // field. `localToken` is required to allow the object to create an
  // appropriate location range for the field.
  public parseObjectAssert = (
    localToken: lexer.Token,
  ): ast.ObjectField | lexical.StaticError => {
    const cond = this.parse(maxPrecedence, null)
    if (lexical.isStaticError(cond)) {
      return cond;
    }
    let msg: ast.Node | null = null;
    if (this.peek().kind === "TokenOperator" && this.peek().data == ":") {
      this.pop();
      const result = this.parse(maxPrecedence, null);
      if (lexical.isStaticError(result)) {
        return result;
      }
      msg = result;
    }

    // Message is optional, so location range changes based on whether
    // it's present.
    const loc: lexical.LocationRange = msg == null
      ?  locFromTokenAST(localToken, cond)
      : locFromTokenAST(localToken, msg);

    return new ast.ObjectField(
      "ObjectAssert",
      "ObjectFieldVisible",
      false,
      false,
      null,
      null,
      im.List<ast.FunctionParam>(),
      false,
      cond,
      msg,
      null,
      loc,
    );
  };

  // parseObjectRemainder parses "the rest" of an object, typically
  // immediately after we encounter the '{' character.
  public parseObjectRemainder = (
    tok: lexer.Token, heading: ast.BindingComment,
  ): {objRemainder: ast.Node, next: lexer.Token} | lexical.StaticError => {
    let fields = im.List<ast.ObjectField>();
    let literalFields = im.Set<literalField>();
    let binds = im.Set<ast.IdentifierName>()

    let gotComma = false
    let first = true

    while (true) {
      // Comments for an object field are allowed to be of either of
      // these forms:
      //
      //     // Explains `foo`.
      //     foo: "bar",
      //
      // or (note the leading comma before the field):
      //
      //     // Explains `foo`.
      //     , foo: "bar"
      //
      // To accomodate both styles, we attempt to parse comments
      // before and after the comma. If there are comments after, that
      // is becomes the heading comment for that field; if not, then
      // we use any comments that happen after the line that contains
      // the last field, but before the comma. So, for example, we
      // ignore the following comment:
      //
      //     , foo: "value1" // This comment is not a heading comment.
      //     // But this one is.
      //     , bar: "value2"
      let headingComments = this.parseOptionalComments();

      let next = this.peek();
      if (!gotComma && !first) {
        if (next.kind === "TokenComma") {
          this.pop();
          next = this.peek();
          gotComma = true
        }
      }

      const thisKind = this.peek().kind;
      if (
        thisKind === "TokenCommentCpp" || thisKind === "TokenCommentC" ||
        thisKind === "TokenCommentHash"
      ) {
        headingComments = this.parseOptionalComments();
      }
      next = this.pop();

      // Done parsing the object. Return.
      if (next.kind === "TokenBraceR") {
        return {
          objRemainder: new ast.ObjectNode(
            fields, gotComma, heading, locFromTokens(tok, next)),
          next: next
        };
      }

      // Object comprehension.
      if (next.kind === "TokenFor") {
        const result = this.parseObjectCompRemainder(
          tok, next, gotComma, fields)
        if (lexical.isStaticError(result)) {
          return result;
        }
        return {objRemainder: result.comp, next: result.last};
      }

      if (!gotComma && !first) {
        return lexical.MakeStaticError(
          "Expected a comma before next field.", next.loc);
      }
      first = false;

      // Start to parse an object field. There are basically 3 valid
      // cases:
      // 1. An object field. The key is a string, an identifier, or a
      //    computed field.
      // 2. A `local` definition.
      // 3. An `assert`.
      switch (next.kind) {
        case "TokenBracketL":
        case "TokenIdentifier":
        case "TokenStringDouble":
        case "TokenStringSingle":
        case "TokenStringBlock": {
          const result = this.parseObjectField(
            headingComments, next, literalFields);
          if (lexical.isStaticError(result)) {
            return result;
          }
          literalFields = result.literals;
          fields = fields.push(result.field);
          break;
        }

        case "TokenLocal": {
          const result = this.parseObjectLocal(next, binds);
          if (lexical.isStaticError(result)) {
            return result;
          }
          binds = result.binds;
          fields = fields.push(result.field);
          break;
        }

        case "TokenAssert": {
          const field = this.parseObjectAssert(next);
          if (lexical.isStaticError(field)) {
            return field;
          }
          fields = fields.push(field);
          break;
        }

        default: {
          return makeUnexpectedError(next, "parsing field definition");
        }
      }
      gotComma = false;
    }
  };


  // parseCompSpecs parses expressions of the form (e.g.) `for x in expr
  // for y in expr if expr for z in expr ...`
  public parseCompSpecs = (
    end: lexer.TokenKind
  ): {compSpecs: ast.CompSpecs, maybeIf: lexer.Token} | lexical.StaticError => {
    let specs = im.List<ast.CompSpec>();
    while (true) {
      const varID = this.popExpect("TokenIdentifier");
      if (lexical.isStaticError(varID)) {
        return varID;
      }

      const id: ast.Identifier = new ast.Identifier(varID.data, varID.loc);
      const pop = this.popExpect("TokenIn");
      if (lexical.isStaticError(pop)) {
        return pop;
      }
      const arr = this.parse(maxPrecedence, null);
      if (lexical.isStaticError(arr)) {
        return arr;
      }
      specs = specs.push(new ast.CompSpecFor(
        id, arr, locFromTokenAST(varID, arr)));

      let maybeIf = this.pop();
      for (; maybeIf.kind === "TokenIf"; maybeIf = this.pop()) {
        const cond = this.parse(maxPrecedence, null);
        if (lexical.isStaticError(cond)) {
          return cond;
        }
        specs = specs.push(new ast.CompSpecIf(
          cond, locFromTokenAST(maybeIf, cond)));
      }
      if (maybeIf.kind === end) {
        return {compSpecs: specs, maybeIf: maybeIf};
      }

      if (maybeIf.kind !== "TokenFor") {
        const tokenKind = lexer.TokenKindStrings.get(end);
        return lexical.MakeStaticError(
          `Expected for, if or ${tokenKind} after for clause, got: ${maybeIf}`, maybeIf.loc);
      }

    }
  };

  // parseArrayRemainder parses "the rest" of an array literal,
  // typically immediately after we encounter the '[' character.
  public parseArrayRemainder = (
    tok: lexer.Token
  ): ast.Node | lexical.StaticError => {
    let next = this.peek();
    if (next.kind === "TokenBracketR") {
      this.pop();
      return new ast.Array(
        im.List<ast.Node>(), false, null, null, locFromTokens(tok, next));
    }

    const first = this.parse(maxPrecedence, null);
    if (lexical.isStaticError(first)) {
      return first;
    }
    let gotComma = false;
    next = this.peek();
    if (next.kind === "TokenComma") {
      this.pop();
      next = this.peek();
      gotComma = true;
    }

    if (next.kind === "TokenFor") {
      // It's a comprehension
      this.pop();
      const result = this.parseCompSpecs("TokenBracketR");
      if (lexical.isStaticError(result)) {
        return result;
      }

      return new ast.ArrayComp(
        first, gotComma, result.compSpecs, locFromTokens(tok, result.maybeIf));
    }
    // Not a comprehension: It can have more elements.
    let elements = im.List<ast.Node>([first]);

    while (true) {
      if (next.kind === "TokenBracketR") {
        // TODO(dcunnin): SYNTAX SUGAR HERE (preserve comma)
        this.pop();
        break;
      }
      if (!gotComma) {
        return lexical.MakeStaticError(
          "Expected a comma before next array element.", next.loc);
      }
      const nextElem = this.parse(maxPrecedence, null);
      if (lexical.isStaticError(nextElem)) {
        return nextElem;
      }
      elements = elements.push(nextElem);

      // Throw away comments before the comma.
      this.parseOptionalComments();

      next = this.peek();
      if (next.kind === "TokenComma") {
        this.pop();
        next = this.peek();
        gotComma = true;
      } else {
        gotComma = false;
      }

      // Throw away comments after the comma.
      this.parseOptionalComments();
    }

    // TODO: Remove trailing whitespace here after we emit newlines
    // from the lexer. If we don't do that, we might accidentally kill
    // comments that correspond to, e.g., the next field of an object.

    return new ast.Array(
      elements, gotComma, null, null, locFromTokens(tok, next));
  };

  public parseTerminal = (
    heading: ast.BindingComment,
  ): ast.Node | lexical.StaticError => {
    let tok = this.pop();
    switch (tok.kind) {
      case "TokenAssert":
      case "TokenBraceR":
      case "TokenBracketR":
      case "TokenComma":
      case "TokenDot":
      case "TokenElse":
      case "TokenError":
      case "TokenFor":
      case "TokenFunction":
      case "TokenIf":
      case "TokenIn":
      case "TokenImport":
      case "TokenImportStr":
      case "TokenLocal":
      case "TokenOperator":
      case "TokenParenR":
      case "TokenSemicolon":
      case "TokenTailStrict":
      case "TokenThen":
        return makeUnexpectedError(tok, "parsing terminal");

      case "TokenEndOfFile":
        return lexical.MakeStaticError("Unexpected end of file.", tok.loc);

      case "TokenBraceL": {
        const result = this.parseObjectRemainder(tok, heading);
        if (lexical.isStaticError(result)) {
          return result;
        }
        return result.objRemainder;
      }

      case "TokenBracketL":
        return this.parseArrayRemainder(tok);

      case "TokenParenL": {
        const inner = this.parse(maxPrecedence, null);
        if (lexical.isStaticError(inner)) {
          return inner;
        }
        const pop = this.popExpect("TokenParenR");
        if (lexical.isStaticError(pop)) {
          return pop;
        }
        return inner;
      }

      // Literals
      case "TokenNumber": {
        // This shouldn't fail as the lexer should make sure we have
        // good input but we handle the error regardless.
        const num = Number(tok.data);
        // TODO: Figure out whether this is correct.
        if (isNaN(num) && tok.data !== "NaN") {
          return lexical.MakeStaticError(
            "Could not parse floating point number.", tok.loc);
        }
        return new ast.LiteralNumber(num, tok.data, tok.loc);
      }
      case "TokenStringSingle":
        return new ast.LiteralStringSingle(tok.data, tok.loc);
      case "TokenStringDouble":
        return new ast.LiteralStringDouble(tok.data, tok.loc);
      case "TokenStringBlock":
        return new ast.LiteralStringBlock(
          tok.data, tok.stringBlockIndent, tok.loc);
      case "TokenFalse":
        return new ast.LiteralBoolean(false, tok.loc);
      case "TokenTrue":
        return new ast.LiteralBoolean(true, tok.loc);
      case "TokenNullLit":
        return new ast.LiteralNull(tok.loc);

      // Variables
      case "TokenDollar":
        return new ast.Dollar(tok.loc);
      case "TokenIdentifier": {
        const id = new ast.Identifier(tok.data, tok.loc);
        return new ast.Var(id, tok.loc);
      }
      case "TokenSelf":
        return new ast.Self(tok.loc);
      case "TokenSuper": {
        const next = this.pop();
        let index: ast.Node | null = null;
        let id: ast.Identifier | null = null;
        switch (next.kind) {
          case "TokenDot": {
            const fieldID = this.popExpect("TokenIdentifier");
            if (lexical.isStaticError(fieldID)) {
              return fieldID;
            }
            id = new ast.Identifier(fieldID.data, fieldID.loc);
            break;
          }
          case "TokenBracketL": {
            let parseErr: lexical.StaticError | null;
            const result = this.parse(maxPrecedence, null);
            if (lexical.isStaticError(result)) {
              return result;
            }
            index = result;
            const pop = this.popExpect("TokenBracketR");
            if (lexical.isStaticError(pop)) {
              return pop;
            }
            break;
          }
        default:
          return lexical.MakeStaticError(
            "Expected . or [ after super.", tok.loc);
        }
        return new ast.SuperIndex(index, id, tok.loc);
      }
    }

    return lexical.MakeStaticError(
      `INTERNAL ERROR: Unknown tok kind: ${tok.kind}`, tok.loc);
  }

  // parse is the main parsing routine.
  public parse = (
    prec: precedence, heading: ast.BindingComment
  ): ast.Node | lexical.StaticError => {
    let begin = this.peek();

    // Consume heading comments if they exist.
    heading = this.parseOptionalComments();

    switch (begin.kind) {
      // These cases have effectively maxPrecedence as the first call
      // to parse will parse them.
      case "TokenAssert": {
        this.pop();
        const cond = this.parse(maxPrecedence, null);
        if (lexical.isStaticError(cond)) {
          return cond;
        }
        let msg: ast.Node | null = null;
        if (this.peek().kind === "TokenOperator" && this.peek().data === ":") {
          this.pop();
          const result = this.parse(maxPrecedence, null);
          if (lexical.isStaticError(result)) {
            return result;
          }
          msg = result;
        }
        const pop = this.popExpect("TokenSemicolon");
        if (lexical.isStaticError(pop)) {
          return pop;
        }
        const rest = this.parse(maxPrecedence, null);
        if (lexical.isStaticError(rest)) {
          return rest;
        }
        return new ast.Assert(cond, msg, rest, locFromTokenAST(begin, rest));
      }

      case "TokenError": {
        this.pop();
        const expr = this.parse(maxPrecedence, null);
        if (lexical.isStaticError(expr)) {
          return expr;
        }
        return new ast.ErrorNode(expr, locFromTokenAST(begin, expr));
      }

      case "TokenIf": {
        this.pop();
        const cond = this.parse(maxPrecedence, null);
        if (lexical.isStaticError(cond)) {
          return cond;
        }
        const pop = this.popExpect("TokenThen");
        if (lexical.isStaticError(pop)) {
          return pop;
        }
        const branchTrue = this.parse(maxPrecedence, null);
        if (lexical.isStaticError(branchTrue)) {
          return branchTrue;
        }
        let branchFalse: ast.Node | null = null;
        let lr = locFromTokenAST(begin, branchTrue);
        if (this.peek().kind === "TokenElse") {
          this.pop();
          const branchFalse = this.parse(maxPrecedence, null);
          if (lexical.isStaticError(branchFalse)) {
            return branchFalse;
          }
          lr = locFromTokenAST(begin, branchFalse)
        }
        return new ast.Conditional(cond, branchTrue, branchFalse, lr);
      }

      case "TokenFunction": {
        this.pop();
        const next = this.pop();
        if (next.kind === "TokenParenL") {
          const result = this.parseParamsList("function parameter");
          if (lexical.isStaticError(result)) {
            return result;
          }

          const body = this.parse(maxPrecedence, null);
          if (lexical.isStaticError(body)) {
            return body;
          }
          const fn = new ast.Function(
            result.params,
            result.gotComma,
            body,
            null,
            im.List<ast.Comment>(),
            locFromTokenAST(begin, body),
          );
          return fn;
        }
        return lexical.MakeStaticError(`Expected ( but got ${next}`, next.loc);
      }

      case "TokenImport": {
        this.pop();
        const body = this.parse(maxPrecedence, null);
        if (lexical.isStaticError(body)) {
          return body;
        }
        if (ast.isLiteralString(body)) {
          return new ast.Import(body.value, locFromTokenAST(begin, body));
        }
        return lexical.MakeStaticError(
          "Computed imports are not allowed", body.loc);
      }

      case "TokenImportStr": {
        this.pop();
        const body = this.parse(maxPrecedence, null);
        if (lexical.isStaticError(body)) {
          return body;
        }
        if (ast.isLiteralString(body)) {
          return new ast.ImportStr(body.value, locFromTokenAST(begin, body));
        }
        return lexical.MakeStaticError(
          "Computed imports are not allowed", body.loc);
      }

      case "TokenLocal": {
        this.pop();
        let binds = im.List<ast.LocalBind>();
        while (true) {
          const newBinds = this.parseBind(begin, binds);
          if (lexical.isStaticError(newBinds)) {
            return newBinds;
          }
          binds = newBinds;
          const delim = this.pop();
          if (delim.kind !== "TokenSemicolon" && delim.kind !== "TokenComma") {
            const msg = `Expected , or ; but got ${delim}`;
            const rest = restFromBinds(newBinds);
            if (rest == null) {
              return lexical.MakeStaticError(msg, delim.loc);
            }
            return lexical.MakeStaticErrorRest(
              rest, msg, delim.loc);
          }
          if (delim.kind === "TokenSemicolon") {
            break;
          }
        }
        const body = this.parse(maxPrecedence, null);
        if (lexical.isStaticError(body)) {
          return body;
        }
        return new ast.Local(binds, body, locFromTokenAST(begin, body));
      }

      default: {
        // Unary operator
        if (begin.kind === "TokenOperator") {
          const uop = ast.UopMap.get(begin.data);
          if (uop == undefined) {
            return lexical.MakeStaticError(
              `Not a unary operator: ${begin.data}`, begin.loc);
          }
          if (prec == unaryPrecedence) {
            const op = this.pop();
            const expr = this.parse(prec, null);
            if (lexical.isStaticError(expr)) {
              return expr;
            }
            return new ast.Unary(uop, expr, locFromTokenAST(op, expr));
          }
        }

        // Base case
        if (prec == 0) {
          return this.parseTerminal(heading);
        }

        let lhs = this.parse(prec-1, heading);
        if (lexical.isStaticError(lhs)) {
          return lhs;
        }

        while (true) {
          // Then next token must be a binary operator.

          let bop: ast.BinaryOp | null = null;

          // Check precedence is correct for this level.  If we're
          // parsing operators with higher precedence, then return lhs
          // and let lower levels deal with the operator.
          switch (this.peek().kind) {
            case "TokenOperator": {
              // _ = "breakpoint"
              if (this.peek().data === ":" || this.peek().data === "=") {
                // Special case for the colons in assert. Since COLON
                // is no-longer a special token, we have to make sure
                // it does not trip the op_is_binary test below.  It
                // should terminate parsing of the expression here,
                // returning control to the parsing of the actual
                // assert AST.
                return lhs;
              }
              bop = ast.BopMap.get(this.peek().data);
              if (bop == undefined) {
                return lexical.MakeStaticError(
                  `Not a binary operator: ${this.peek().data}`, this.peek().loc);
              }

              if (bopPrecedence.get(bop) != prec) {
                return lhs;
              }
              break;
            }

            case "TokenDot":
            case "TokenBracketL":
            case "TokenParenL":
            case "TokenBraceL": {
              if (applyPrecedence != prec) {
                return lhs;
              }
              break;
            }
            default:
              return lhs;
          }

          const op = this.pop();
          switch (op.kind) {
            case "TokenBracketL": {
              const index = this.parse(maxPrecedence, null);
              if (lexical.isStaticError(index)) {
                return index;
              }
              const end = this.popExpect("TokenBracketR");
              if (lexical.isStaticError(end)) {
                return end;
              }

              lhs = new ast.IndexSubscript(
                lhs, index, locFromTokens(begin, end));
              break;
            }
            case "TokenDot": {
              const fieldID = this.popExpect("TokenIdentifier");
              if (lexical.isStaticError(fieldID)) {
                // After the user types a `.`, the document very
                // likely doesn't parse. For autocomplete facilities,
                // it's useful to return the AST that precedes the `.`
                // character (typically a `Var` or `Index`
                // expression), so that it is easier to discern what
                // to complete.
                return lexical.MakeStaticErrorRest(lhs, fieldID.msg, fieldID.loc);
              }
              const id = new ast.Identifier(fieldID.data, fieldID.loc);
              lhs = new ast.IndexDot(lhs, id, locFromTokens(begin, fieldID));
              break;
            }
            case "TokenParenL": {
              const result = this.parseArgsList("function argument");
              if (lexical.isStaticError(result)) {
                return result;
              }

              const {next: end, params: args, gotComma: gotComma} = result;
              let tailStrict = false
              if (this.peek().kind === "TokenTailStrict") {
                this.pop();
                tailStrict = true;
              }
              lhs = new ast.Apply(
                lhs, args, gotComma, tailStrict, locFromTokens(begin, end));
              break;
            }
            case "TokenBraceL": {
              const result = this.parseObjectRemainder(op, heading);
              if (lexical.isStaticError(result)) {
                return result;
              }
              lhs = new ast.ApplyBrace(
                lhs, result.objRemainder, locFromTokens(begin, result.next));
              break;
            }
            default: {
              const rhs = this.parse(prec-1, null);
              if (lexical.isStaticError(rhs)) {
                return rhs;
              }
              if (bop == null) {
                throw new Error(
                  "INTERNAL ERROR: `parse` can't return a null node unless an `error` is populated");
              }
              lhs = new ast.Binary(lhs, bop, rhs, locFromTokenAST(begin, rhs));
              break;
            }
          }
        }
      }
    }
  }
}

type literalField = string;

// ---------------------------------------------------------------------------

const restFromBinds = (newBinds: ast.LocalBinds): ast.Var | null => {
  if (newBinds.count() == 0) {
    return null
  }

  const lastBody = newBinds.last().body;
  if (ast.isBinary(lastBody) && lastBody.op === "BopPlus" &&
    ast.isVar(lastBody.right)
  ) {
    return lastBody.right;
  } else if (ast.isVar(lastBody)) {
    return lastBody;
  }

  return null;
}

// ---------------------------------------------------------------------------

export const Parse = (
  t: lexer.Tokens
): ast.Node | lexical.StaticError => {
  const p = new parser(t);
  const expr = p.parse(maxPrecedence, null);
  if (lexical.isStaticError(expr)) {
    return expr;
  }

  // Get rid of any trailing comments.
  p.parseOptionalComments();

  if (p.peek().kind !== "TokenEndOfFile") {
    return lexical.MakeStaticError(`Did not expect: ${p.peek()}`, p.peek().loc);
  }
  new ast.InitializingVisitor(expr).visit();

  return expr;
}
