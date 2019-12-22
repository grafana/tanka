"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const im = require("immutable");
const lexical = require("./lexical");
exports.makeCodePoints = (str) => {
    // NOTE: Splitting a string by Unicode code points is tricky in
    // JavaScript. We use the spread operator here, because it does
    // the split correctly (versus, say, `String.split).
    return im.List([...str]);
};
// TODO: Refactor this into a member of `CodePoints`.
exports.stringSlice = (points, begin, end) => {
    return points.slice(begin, end).join("");
};
exports.TokenKindStrings = im.Map({
    // Symbols
    TokenBraceL: "\"{\"",
    TokenBraceR: "\"}\"",
    TokenBracketL: "\"[\"",
    TokenBracketR: "\"]\"",
    TokenComma: "\",\"",
    TokenDollar: "\"$\"",
    TokenDot: "\".\"",
    TokenParenL: "\"(\"",
    TokenParenR: "\")\"",
    TokenSemicolon: "\";\"",
    // Arbitrary length lexemes
    TokenIdentifier: "IDENTIFIER",
    TokenNumber: "NUMBER",
    TokenOperator: "OPERATOR",
    TokenStringBlock: "STRING_BLOCK",
    TokenStringDouble: "STRING_DOUBLE",
    TokenStringSingle: "STRING_SINGLE",
    TokenCommentCpp: "CPP_COMMENT",
    TokenCommentC: "C_COMMENT",
    TokenCommentHash: "HASH_COMMENT",
    // Keywords
    TokenAssert: "assert",
    TokenElse: "else",
    TokenError: "error",
    TokenFalse: "false",
    TokenFor: "for",
    TokenFunction: "function",
    TokenIf: "if",
    TokenImport: "import",
    TokenImportStr: "importstr",
    TokenIn: "in",
    TokenLocal: "local",
    TokenNullLit: "null",
    TokenSelf: "self",
    TokenSuper: "super",
    TokenTailStrict: "tailstrict",
    TokenThen: "then",
    TokenTrue: "true",
    // A special token that holds line/column information about the end
    // of the file.
    TokenEndOfFile: "end of file",
});
class Token {
    constructor(kind, // The type of the token
    fodder, // Any fodder the occurs before this token
    data, // Content of the token if it is not a keyword
    // Extra info for when kind == tokenStringBlock
    stringBlockIndent, // The sequence of whitespace that indented the block.
    stringBlockTermIndent, // This is always fewer whitespace characters than in stringBlockIndent.
    loc) {
        this.kind = kind;
        this.fodder = fodder;
        this.data = data;
        this.stringBlockIndent = stringBlockIndent;
        this.stringBlockTermIndent = stringBlockTermIndent;
        this.loc = loc;
    }
    toString() {
        const tokenKind = exports.TokenKindStrings.get(this.kind);
        if (this.data == "") {
            return tokenKind;
        }
        else if (this.kind == "TokenOperator") {
            return `"${this.data}"`;
        }
        else {
            return `(${tokenKind}, "${this.data}")`;
        }
    }
}
exports.Token = Token;
exports.emptyTokens = im.List();
// ---------------------------------------------------------------------------
// Helpers
exports.isUpper = (r) => {
    return r.data >= 'A' && r.data <= 'Z';
};
exports.isLower = (r) => {
    return r.data >= 'a' && r.data <= 'z';
};
exports.isNumber = (r) => {
    return r.data >= '0' && r.data <= '9';
};
exports.isIdentifierFirst = (r) => {
    return exports.isUpper(r) || exports.isLower(r) || r.data === '_';
};
exports.isIdentifier = (r) => {
    return exports.isIdentifierFirst(r) || exports.isNumber(r);
};
exports.isSymbol = (r) => {
    switch (r.data) {
        case '!':
        case '$':
        case ':':
        case '~':
        case '+':
        case '-':
        case '&':
        case '|':
        case '^':
        case '=':
        case '<':
        case '>':
        case '*':
        case '/':
        case '%':
            return true;
    }
    return false;
};
// Check that b has at least the same whitespace prefix as a and
// returns the amount of this whitespace, otherwise returns 0.  If a
// has no whitespace prefix than return 0.
exports.checkWhitespace = (a, b) => {
    let i = 0;
    for (; i < a.length; i++) {
        if (a[i] != ' ' && a[i] != '\t') {
            // a has run out of whitespace and b matched up to this point.
            // Return result.
            return i;
        }
        if (i >= b.length) {
            // We ran off the edge of b while a still has whitespace.
            // Return 0 as failure.
            return 0;
        }
        if (a[i] != b[i]) {
            // a has whitespace but b does not.  Return 0 as failure.
            return 0;
        }
    }
    // We ran off the end of a and b kept up
    return i;
};
;
// NOTE: `pos` is the index of the code point, not the index of a byte
// in the string.
exports.runeFromString = (str, pos) => {
    return exports.runeFromCodePoints(exports.makeCodePoints(str), pos);
};
exports.runeFromCodePoints = (str, pos) => {
    const r = str.get(pos), codePoint = r.codePointAt(0);
    return {
        codePoint: codePoint,
        data: r,
    };
};
const LexEOF = {
    codePoint: -1,
    data: "\0",
};
// TODO: Replace this. This is because we need a special rune type in
// TS, but it should be phased out.
const LexEOFPos = -1;
class lexer {
    constructor(fn, input) {
        // next returns the next rune in the input.
        this.next = () => {
            if (this.pos >= this.input.count()) {
                this.prevPos = this.pos;
                return LexEOF;
            }
            const r = exports.runeFromCodePoints(this.input, this.pos);
            // NOTE: Because `CodePoints` is essentially an array of distinct
            // code points, rather than an array of bytes. So unlike the Go
            // implementation of this code, `pos` only ever needs to be
            // advanced by 1 (rather than the number of bytes a code point
            // takes up).
            this.prevPos = this.pos;
            this.pos += 1;
            if (r.data === '\n') {
                this.prevLineNumber = this.lineNumber;
                this.prevLineStart = this.lineStart;
                this.lineNumber++;
                this.lineStart = this.pos;
            }
            return r;
        };
        this.acceptN = (n) => {
            for (let i = 0; i < n; i++) {
                this.next();
            }
        };
        // peek returns but does not consume the next rune in the input.
        this.peek = () => {
            if (this.pos >= this.input.count()) {
                this.prevPos = this.pos;
                return LexEOF;
            }
            const r = exports.runeFromCodePoints(this.input, this.pos);
            return r;
        };
        // backup steps back one rune. Can only be called once per call of
        // next.
        this.backup = () => {
            if (this.prevPos === LexEOFPos) {
                throw new Error("INTERNAL ERROR: backup called with no valid previous rune");
            }
            if ((this.prevPos - this.lineStart) < 0) {
                this.lineNumber = this.prevLineNumber;
                this.lineStart = this.prevLineStart;
            }
            this.pos = this.prevPos;
            this.prevPos = LexEOFPos;
        };
        this.location = () => {
            return new lexical.Location(this.lineNumber, this.pos - this.lineStart + 1);
        };
        this.prevLocation = () => {
            if (this.prevPos == LexEOFPos) {
                throw new Error("INTERNAL ERROR: prevLocation called with no valid previous rune");
            }
            return new lexical.Location(this.prevLineNumber, this.prevPos - this.prevLineStart + 1);
        };
        // Reset the current working token start to the current cursor
        // position.  This may throw away some characters.  This does not
        // throw away any accumulated fodder.
        this.resetTokenStart = () => {
            this.tokenStart = this.pos;
            this.tokenStartLoc = this.location();
        };
        this.emitFullToken = (kind, data, stringBlockIndent, stringBlockTermIndent) => {
            this.tokens = this.tokens.push(new Token(kind, this.fodder, data, stringBlockIndent, stringBlockTermIndent, lexical.MakeLocationRange(this.fileName, this.tokenStartLoc, this.location())));
            this.fodder = [];
        };
        this.emitToken = (kind) => {
            this.emitFullToken(kind, exports.stringSlice(this.input, this.tokenStart, this.pos), "", "");
            this.resetTokenStart();
        };
        this.addWhitespaceFodder = () => {
            const fodderData = exports.stringSlice(this.input, this.tokenStart, this.pos);
            if (this.fodder.length == 0 || this.fodder[this.fodder.length - 1].kind != "FodderWhitespace") {
                this.fodder.push({
                    kind: "FodderWhitespace",
                    data: fodderData
                });
            }
            else {
                this.fodder[this.fodder.length - 1].data += fodderData;
            }
            this.resetTokenStart();
        };
        this.addCommentFodder = (kind) => {
            const fodderData = exports.stringSlice(this.input, this.tokenStart, this.pos);
            this.fodder.push({ kind: kind, data: fodderData });
            this.resetTokenStart();
        };
        this.addFodder = (kind, data) => {
            this.fodder.push({ kind: kind, data: data });
        };
        // lexNumber will consume a number and emit a token.  It is assumed
        // that the next rune to be served by the lexer will be a leading
        // digit.
        this.lexNumber = () => {
            // This function should be understood with reference to the linked
            // image: http://www.json.org/number.gif
            let state = "numBegin";
            outerLoop: while (true) {
                const r = this.next();
                switch (state) {
                    case "numBegin": {
                        if (r.data === '0') {
                            state = "numAfterZero";
                        }
                        else if (r.data >= '1' && r.data <= '9') {
                            state = "numAfterOneToNine";
                        }
                        else {
                            // The caller should ensure the first rune is a digit.
                            throw new Error("INTERNAL ERROR: Couldn't lex number");
                        }
                        break;
                    }
                    case "numAfterZero": {
                        if (r.data === '.') {
                            state = "numAfterDot";
                        }
                        else if (r.data === 'e' || r.data === 'E') {
                            state = "numAfterE";
                        }
                        else {
                            break outerLoop;
                        }
                        break;
                    }
                    case "numAfterOneToNine": {
                        if (r.data === '.') {
                            state = "numAfterDot";
                        }
                        else if (r.data === 'e' || r.data === 'E') {
                            state = "numAfterE";
                        }
                        else if (r.data >= '0' && r.data <= '9') {
                            state = "numAfterOneToNine";
                        }
                        else {
                            break outerLoop;
                        }
                        break;
                    }
                    case "numAfterDot": {
                        if (r.data >= '0' && r.data <= '9') {
                            state = "numAfterDigit";
                        }
                        else {
                            return lexical.MakeStaticErrorPoint(`Couldn't lex number, junk after decimal point: '${r.data}'`, this.fileName, this.prevLocation());
                        }
                        break;
                    }
                    case "numAfterDigit": {
                        if (r.data === 'e' || r.data === 'E') {
                            state = "numAfterE";
                        }
                        else if (r.data >= '0' && r.data <= '9') {
                            state = "numAfterDigit";
                        }
                        else {
                            break outerLoop;
                        }
                        break;
                    }
                    case "numAfterE": {
                        if (r.data === '+' || r.data === '-') {
                            state = "numAfterExpSign";
                        }
                        else if (r.data >= '0' && r.data <= '9') {
                            state = "numAfterExpDigit";
                        }
                        else {
                            return lexical.MakeStaticErrorPoint(`Couldn't lex number, junk after 'E': '${r.data}'`, this.fileName, this.prevLocation());
                        }
                        break;
                    }
                    case "numAfterExpSign": {
                        if (r.data >= '0' && r.data <= '9') {
                            state = "numAfterExpDigit";
                        }
                        else {
                            return lexical.MakeStaticErrorPoint(`Couldn't lex number, junk after exponent sign: '${r.data}'`, this.fileName, this.prevLocation());
                        }
                        break;
                    }
                    case "numAfterExpDigit": {
                        if (r.data >= '0' && r.data <= '9') {
                            state = "numAfterExpDigit";
                        }
                        else {
                            break outerLoop;
                        }
                        break;
                    }
                }
            }
            this.backup();
            this.emitToken("TokenNumber");
            return null;
        };
        // lexIdentifier will consume a identifer and emit a token. It is
        // assumed that the next rune to be served by the lexer will be a
        // leading digit. This may emit a keyword or an identifier.
        this.lexIdentifier = () => {
            let r = this.next();
            if (!exports.isIdentifierFirst(r)) {
                throw new Error("INTERNAL ERROR: Unexpected character in lexIdentifier");
            }
            for (; r.codePoint != LexEOF.codePoint; r = this.next()) {
                if (!exports.isIdentifier(r)) {
                    break;
                }
            }
            this.backup();
            switch (exports.stringSlice(this.input, this.tokenStart, this.pos)) {
                case "assert":
                    this.emitToken("TokenAssert");
                    break;
                case "else":
                    this.emitToken("TokenElse");
                    break;
                case "error":
                    this.emitToken("TokenError");
                    break;
                case "false":
                    this.emitToken("TokenFalse");
                    break;
                case "for":
                    this.emitToken("TokenFor");
                    break;
                case "function":
                    this.emitToken("TokenFunction");
                    break;
                case "if":
                    this.emitToken("TokenIf");
                    break;
                case "import":
                    this.emitToken("TokenImport");
                    break;
                case "importstr":
                    this.emitToken("TokenImportStr");
                    break;
                case "in":
                    this.emitToken("TokenIn");
                    break;
                case "local":
                    this.emitToken("TokenLocal");
                    break;
                case "null":
                    this.emitToken("TokenNullLit");
                    break;
                case "self":
                    this.emitToken("TokenSelf");
                    break;
                case "super":
                    this.emitToken("TokenSuper");
                    break;
                case "tailstrict":
                    this.emitToken("TokenTailStrict");
                    break;
                case "then":
                    this.emitToken("TokenThen");
                    break;
                case "true":
                    this.emitToken("TokenTrue");
                    break;
                default:
                    // Not a keyword, assume it is an identifier
                    this.emitToken("TokenIdentifier");
                    break;
            }
            ;
        };
        // locBeforeLastTokenRange checks whether a location specified by
        // `loc` exists before the range of coordinates of the last token
        // terminates.
        this.locBeforeLastTokenRange = (loc) => {
            const numTokens = this.tokens.count();
            if (loc.line == -1 && loc.column == -1) {
                return false;
            }
            else if (numTokens == 0) {
                return false;
            }
            const lastLocRange = this.tokens.get(numTokens - 1).loc;
            return loc.line < lastLocRange.begin.line ||
                (loc.line == lastLocRange.begin.line && loc.column < lastLocRange.begin.column);
        };
        // locInLastTokenRange checks whether a location specified by `loc`
        // exists within the range of coordinates of the last token.
        this.locInLastTokenRange = (loc) => {
            const numTokens = this.tokens.count();
            if (loc.line == -1 && loc.column == -1) {
                return false;
            }
            else if (numTokens == 0) {
                return false;
            }
            const lastToken = this.tokens.get(numTokens - 1);
            const lastLocRange = lastToken.loc;
            if ((lastLocRange.begin.line == loc.line) &&
                loc.line == lastLocRange.end.line &&
                lastLocRange.begin.column <= loc.column &&
                loc.column <= lastLocRange.end.column) {
                return true;
            }
            else if ((lastLocRange.begin.line < loc.line) &&
                loc.line == lastLocRange.end.line &&
                loc.column <= lastLocRange.end.column) {
                return true;
            }
            else if ((lastLocRange.begin.line == loc.line) &&
                loc.line < lastLocRange.end.line &&
                loc.column >= lastLocRange.begin.column) {
                return true;
            }
            else if ((lastLocRange.begin.line < loc.line) &&
                loc.line < lastLocRange.end.line) {
                return true;
            }
            else {
                return false;
            }
        };
        // checkTruncateTokenRange truncates the token stream if it exceeds
        // the token range. Note that a corner case is if the range max
        // happens to occur in whitespace; in this case, we will truncate at
        // the last token that occurs before the whitespace begins.
        this.checkTruncateTokenRange = (rangeMax) => {
            if (rangeMax.line == -1 && rangeMax.column == -1) {
                return false;
            }
            const numTokens = this.tokens.count();
            // Lex at least one token before returning.
            if (numTokens == 0) {
                return false;
            }
            while (true) {
                // If we have truncated all the tokens in the stream, return.
                if (numTokens == 0) {
                    return true;
                }
                // Return if location is in the range of the last token.
                if (this.locInLastTokenRange(rangeMax)) {
                    return true;
                }
                // If token range max occurs before the last token range starts,
                // truncate and return.
                if (this.locBeforeLastTokenRange(rangeMax)) {
                    this.tokens = this.tokens.pop();
                    // Stop truncating after the condition is no longer true.
                    if (!this.locBeforeLastTokenRange(rangeMax)) {
                        return true;
                    }
                    continue;
                }
                return false;
            }
        };
        this.fileName = fn;
        this.input = exports.makeCodePoints(input);
        this.lineNumber = 1;
        this.prevPos = LexEOFPos;
        this.prevLineNumber = 1;
        this.tokenStartLoc = new lexical.Location(1, 1);
        this.tokens = im.List();
        this.fodder = [];
        this.pos = 0;
        this.lineStart = 0;
        this.tokenStart = 0;
        this.prevLineStart = 0;
    }
    ;
    // lexSymbol will lex a token that starts with a symbol. This could
    // be a C or C++ comment, block quote or an operator. This function
    // assumes that the next rune to be served by the lexer will be the
    // first rune of the new token.
    lexSymbol() {
        let r = this.next();
        // Single line C++ style comment
        if (r.data === '/' && this.peek().data === '/') {
            this.next();
            this.resetTokenStart(); // Throw out the leading //
            for (r = this.next(); r.codePoint != LexEOF.codePoint && r.data !== '\n'; r = this.next()) {
            }
            // Leave the '\n' in the lexer to be fodder for the next round
            this.backup();
            this.emitToken("TokenCommentCpp");
            return null;
        }
        if (r.data === '/' && this.peek().data === '*') {
            const commentStartLoc = this.tokenStartLoc;
            this.next(); // consume the '*'
            this.resetTokenStart(); // Throw out the leading /*
            for (r = this.next();; r = this.next()) {
                if (r.codePoint == LexEOF.codePoint) {
                    return lexical.MakeStaticErrorPoint("Multi-line comment has no terminating */", this.fileName, commentStartLoc);
                }
                if (r.data === '*' && this.peek().data === '/') {
                    // Don't include trailing */
                    this.backup();
                    this.emitToken("TokenCommentC");
                    this.next(); // Skip past '*'
                    this.next(); // Skip past '/'
                    this.resetTokenStart(); // Start next token at this point
                    return null;
                }
            }
        }
        if (r.data === '|' && exports.stringSlice(this.input, this.pos).startsWith("||\n")) {
            const commentStartLoc = this.tokenStartLoc;
            this.acceptN(3); // Skip "||\n"
            var cb = im.List();
            // Skip leading blank lines
            for (r = this.next(); r.data === '\n'; r = this.next()) {
                cb = cb.push(r);
            }
            this.backup();
            let numWhiteSpace = exports.checkWhitespace(exports.stringSlice(this.input, this.pos), exports.stringSlice(this.input, this.pos));
            const stringBlockIndent = exports.stringSlice(this.input, this.pos, this.pos + numWhiteSpace);
            if (numWhiteSpace == 0) {
                return lexical.MakeStaticErrorPoint("Text block's first line must start with whitespace", this.fileName, commentStartLoc);
            }
            while (true) {
                if (numWhiteSpace <= 0) {
                    throw new Error("INTERNAL ERROR: Unexpected value for numWhiteSpace");
                }
                this.acceptN(numWhiteSpace);
                for (r = this.next(); r.data !== '\n'; r = this.next()) {
                    if (r.codePoint == LexEOF.codePoint) {
                        return lexical.MakeStaticErrorPoint("Unexpected EOF", this.fileName, commentStartLoc);
                    }
                    cb = cb.push(r);
                }
                cb = cb.push(exports.runeFromString("\n", 0));
                // Skip any blank lines
                for (r = this.next(); r.data === '\n'; r = this.next()) {
                    cb = cb.push(r);
                }
                this.backup();
                // Look at the next line
                numWhiteSpace = exports.checkWhitespace(stringBlockIndent, exports.stringSlice(this.input, this.pos));
                if (numWhiteSpace == 0) {
                    // End of the text block
                    let stringBlockTermIndent = "";
                    for (r = this.next(); r.data === ' ' || r.data === '\t'; r = this.next()) {
                        stringBlockTermIndent = stringBlockIndent.concat(r.data);
                    }
                    this.backup();
                    if (!exports.stringSlice(this.input, this.pos).startsWith("|||")) {
                        return lexical.MakeStaticErrorPoint("Text block not terminated with |||", this.fileName, commentStartLoc);
                    }
                    this.acceptN(3); // Skip '|||'
                    const tokenData = cb
                        .map((rune) => {
                        return rune.data;
                    })
                        .join("");
                    this.emitFullToken("TokenStringBlock", tokenData, stringBlockIndent, stringBlockTermIndent);
                    this.resetTokenStart();
                    return null;
                }
            }
        }
        // Assume any string of symbols is a single operator.
        for (r = this.next(); exports.isSymbol(r); r = this.next()) {
            // Not allowed // in operators
            if (r.data === '/' && exports.stringSlice(this.input, this.pos).startsWith("/")) {
                break;
            }
            // Not allowed /* in operators
            if (r.data === '/' && exports.stringSlice(this.input, this.pos).startsWith("*")) {
                break;
            }
            // Not allowed ||| in operators
            if (r.data === '|' && exports.stringSlice(this.input, this.pos).startsWith("||")) {
                break;
            }
        }
        this.backup();
        // Operators are not allowed to end with + - ~ ! unless they are
        // one rune long. So, wind it back if we need to, but stop at the
        // first rune. This relies on the hack that all operator symbols
        // are ASCII and thus there is no need to treat this substring as
        // general UTF-8.
        for (r = exports.runeFromCodePoints(this.input, this.pos - 1); this.pos > this.tokenStart + 1; this.pos--) {
            switch (r.data) {
                case '+':
                case '-':
                case '~':
                case '!':
                    continue;
            }
            break;
        }
        if (exports.stringSlice(this.input, this.tokenStart, this.pos) == "$") {
            this.emitToken("TokenDollar");
        }
        else {
            this.emitToken("TokenOperator");
        }
        return null;
    }
    ;
}
exports.lexer = lexer;
exports.Lex = (fn, input) => {
    const unlimitedRange = new lexical.Location(-1, -1);
    return exports.LexRange(fn, input, unlimitedRange);
};
exports.LexRange = (fn, input, tokenRange) => {
    const l = new lexer(fn, input);
    let err = null;
    for (let r = l.next(); r.codePoint != LexEOF.codePoint; r = l.next()) {
        // Terminate lexing if we're past the token range. If we've lexed
        // past the desired range, we will truncate the token stream.
        if (l.checkTruncateTokenRange(tokenRange)) {
            return l.tokens;
        }
        switch (r.data) {
            case ' ':
            case '\t':
            case '\r':
            case '\n':
                l.addWhitespaceFodder();
                continue;
            case '{':
                l.emitToken("TokenBraceL");
                break;
            case '}':
                l.emitToken("TokenBraceR");
                break;
            case '[':
                l.emitToken("TokenBracketL");
                break;
            case ']':
                l.emitToken("TokenBracketR");
                break;
            case ',':
                l.emitToken("TokenComma");
                break;
            case '.':
                l.emitToken("TokenDot");
                break;
            case '(':
                l.emitToken("TokenParenL");
                break;
            case ')':
                l.emitToken("TokenParenR");
                break;
            case ';':
                l.emitToken("TokenSemicolon");
                break;
            case '0':
            case '1':
            case '2':
            case '3':
            case '4':
            case '5':
            case '6':
            case '7':
            case '8':
            case '9': {
                l.backup();
                err = l.lexNumber();
                if (err != null) {
                    return err;
                }
                break;
            }
            // String literals
            case '"': {
                const stringStartLoc = l.prevLocation();
                // Don't include the quotes in the token data
                l.resetTokenStart();
                for (r = l.next();; r = l.next()) {
                    if (r.codePoint == LexEOF.codePoint) {
                        return lexical.MakeStaticErrorPoint("Unterminated String", l.fileName, stringStartLoc);
                    }
                    if (r.data === '"') {
                        l.backup();
                        l.emitToken("TokenStringDouble");
                        /*_ =*/ l.next();
                        l.resetTokenStart();
                        break;
                    }
                    if (r.data === '\\' && l.peek().codePoint != LexEOF.codePoint) {
                        r = l.next();
                    }
                }
                break;
            }
            case '\'': {
                const stringStartLoc = l.prevLocation();
                l.resetTokenStart(); // Don't include the quotes in the token data
                for (r = l.next();; r = l.next()) {
                    if (r.codePoint == LexEOF.codePoint) {
                        return lexical.MakeStaticErrorPoint("Unterminated String", l.fileName, stringStartLoc);
                    }
                    if (r.data === '\'') {
                        l.backup();
                        l.emitToken("TokenStringSingle");
                        r = l.next();
                        l.resetTokenStart();
                        break;
                    }
                    if (r.data === '\\' && l.peek().codePoint != LexEOF.codePoint) {
                        r = l.next();
                    }
                }
                break;
            }
            case '#': {
                l.resetTokenStart(); // Throw out the leading #
                for (r = l.next(); r.codePoint != LexEOF.codePoint && r.data !== '\n'; r = l.next()) {
                }
                // Leave the '\n' in the lexer to be fodder for the next round
                l.backup();
                l.emitToken("TokenCommentHash");
                break;
            }
            default: {
                if (exports.isIdentifierFirst(r)) {
                    l.backup();
                    l.lexIdentifier();
                }
                else if (exports.isSymbol(r)) {
                    l.backup();
                    err = l.lexSymbol();
                    if (err != null) {
                        return err;
                    }
                }
                else {
                    return lexical.MakeStaticErrorPoint(`Could not lex the character '${r.data}'`, l.fileName, l.prevLocation());
                }
                break;
            }
        }
    }
    // We are currently at the EOF.  Emit a special token to capture any
    // trailing fodder
    l.emitToken("TokenEndOfFile");
    return l.tokens;
};
//# sourceMappingURL=lexer.js.map