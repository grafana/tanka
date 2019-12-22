"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// ParsedDocument represents a successfully-parsed document.
class ParsedDocument {
    constructor(text, lex, parse, version) {
        this.text = text;
        this.lex = lex;
        this.parse = parse;
        this.version = version;
    }
}
exports.ParsedDocument = ParsedDocument;
exports.isParsedDocument = (testMe) => {
    return testMe instanceof ParsedDocument;
};
// FailedParsedDocument represents a document that failed to parse.
class FailedParsedDocument {
    constructor(text, parse, version) {
        this.text = text;
        this.parse = parse;
        this.version = version;
    }
}
exports.FailedParsedDocument = FailedParsedDocument;
exports.isFailedParsedDocument = (testMe) => {
    return testMe instanceof FailedParsedDocument;
};
// LexFailure represents a failure to lex a document.
class LexFailure {
    constructor(lex, lexError) {
        this.lex = lex;
        this.lexError = lexError;
    }
}
exports.LexFailure = LexFailure;
exports.isLexFailure = (testMe) => {
    return testMe instanceof LexFailure;
};
// ParseFailure represents a failure to parse a document.
class ParseFailure {
    constructor(lex, 
    // TODO: Enable this.
    // readonly parse: ast.Node,
    parseError) {
        this.lex = lex;
        this.parseError = parseError;
    }
}
exports.ParseFailure = ParseFailure;
exports.isParseFailure = (testMe) => {
    return testMe instanceof ParseFailure;
};
//# sourceMappingURL=service.js.map