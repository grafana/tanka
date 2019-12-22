"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const http = require("http");
const url = require("url");
const im = require("immutable");
const ast = require("../compiler/lexical-analysis/ast");
const lexer = require("../compiler/lexical-analysis/lexer");
const lexical = require("../compiler/lexical-analysis/lexical");
const parser = require("../compiler/lexical-analysis/parser");
const _static = require("../compiler/static");
// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------
const getUrl = (url, cb) => {
    let text = '';
    http.get(url, function (res) {
        const { statusCode } = res;
        let error;
        if (statusCode !== 200) {
            res.resume();
            throw new Error(`Request Failed.\n Status Code: ${statusCode}`);
        }
        res.setEncoding('utf8');
        res.on('data', (chunk) => { text += chunk; });
        res.on('end', () => {
            cb(text);
        });
    });
};
// ----------------------------------------------------------------------------
// Browser-specific implementations of core analyzer constructs.
// ----------------------------------------------------------------------------
class BrowserDocumentManager {
    constructor() {
        // URI utilities.
        this.backsplicePrefix = `file:///`;
        this.windowDocUri = `${this.backsplicePrefix}window`;
        this.k8sUri = `${this.backsplicePrefix}k8s.libsonnet`;
        // The ksonnet files (e.g., `k.libsonnet`) never change. Because
        // static analysis of their files is expensive, we assign them
        // version 0 so that it's _always_ cached.
        this.staticsVersion = 0;
        this.get = (file) => {
            const fileUri = ast.isImport(file) || ast.isImportStr(file)
                ? `${this.backsplicePrefix}${file.file}`
                : file;
            if (fileUri === `${this.backsplicePrefix}ksonnet.beta.2/k.libsonnet`) {
                return {
                    text: this.k,
                    version: this.staticsVersion,
                    resolvedPath: fileUri,
                };
            }
            else if (fileUri === this.k8sUri) {
                return {
                    text: this.k8s,
                    version: this.staticsVersion,
                    resolvedPath: fileUri,
                };
            }
            else if (fileUri === this.windowDocUri) {
                return {
                    text: this.windowText,
                    version: this.version,
                    resolvedPath: fileUri,
                };
            }
            throw new Error(`Unrecognized file ${fileUri}`);
        };
        this.setWindowText = (text, version) => {
            this.windowText = text;
            this.version = version;
        };
        this.windowText = "";
        this.version = undefined;
    }
}
exports.BrowserDocumentManager = BrowserDocumentManager;
class BrowserCompilerService {
    constructor() {
        this.cache = (fileUri, text, version) => {
            //
            // There are 3 possible outcomes:
            //
            // 1. We successfully parse the document. Cache.
            // 2. We successfully lex but fail to parse. Return
            //    `PartialParsedDocument`.
            // 3. We fail to lex. Return `PartialParsedDocument`.
            //
            // Attempt to retrieve cached parse if document versions are the
            // same. If version is undefined, it comes from a source that
            // doesn't track document version, and we always re-parse.
            const tryGet = this.docCache.get(fileUri);
            if (tryGet !== undefined && tryGet.version !== undefined &&
                tryGet.version === version) {
                return tryGet;
            }
            // TODO: Replace this with a URL provider abstraction.
            const parsedUrl = url.parse(fileUri);
            if (!parsedUrl || !parsedUrl.path) {
                throw new Error(`INTERNAL ERROR: Failed to parse URI '${fileUri}'`);
            }
            const lex = lexer.Lex(parsedUrl.path, text);
            if (lexical.isStaticError(lex)) {
                // TODO: emptyTokens is not right. Fill it in.
                const fail = new _static.LexFailure(lexer.emptyTokens, lex);
                return new _static.FailedParsedDocument(text, fail, version);
            }
            const parse = parser.Parse(lex);
            if (lexical.isStaticError(parse)) {
                const fail = new _static.ParseFailure(lex, parse);
                return new _static.FailedParsedDocument(text, fail, version);
            }
            const parsedDoc = new _static.ParsedDocument(text, lex, parse, version);
            this.docCache = this.docCache.set(fileUri, parsedDoc);
            return parsedDoc;
        };
        this.getLastSuccess = (fileUri) => {
            return this.docCache.has(fileUri) && this.docCache.get(fileUri) || null;
        };
        this.delete = (fileUri) => {
            this.docCache = this.docCache.delete(fileUri);
        };
        //
        // Private members.
        //
        this.docCache = im.Map();
    }
}
exports.BrowserCompilerService = BrowserCompilerService;
// ----------------------------------------------------------------------------
// Set up analyzer in browser.
// ----------------------------------------------------------------------------
const docs = new BrowserDocumentManager();
const cs = new BrowserCompilerService();
const analyzer = new _static.Analyzer(docs, cs);
// ----------------------------------------------------------------------------
// Get ksonnet files.
// ----------------------------------------------------------------------------
getUrl('https://raw.githubusercontent.com/ksonnet/ksonnet-lib/bd6b2d618d6963ea6a81fcc5623900d8ba110a32/ksonnet.beta.2/k.libsonnet', text => { docs.k = text; });
getUrl("https://raw.githubusercontent.com/ksonnet/ksonnet-lib/bd6b2d618d6963ea6a81fcc5623900d8ba110a32/ksonnet.beta.2/k8s.libsonnet", text => {
    docs.k8s = text;
    // Static analysis on `k8s.libsonnet` takes multiple seconds to
    // complete, so do this immediately.
    cs.cache(docs.k8sUri, text, docs.staticsVersion);
});
;
// ----------------------------------------------------------------------------
// Public functions for the Monaco editor to call.
// ----------------------------------------------------------------------------
global.docOnChange = (text, version) => {
    docs.setWindowText(text, version);
    cs.cache(docs.windowDocUri, text, version);
};
global.onComplete = (text, position) => {
    return analyzer
        .onComplete(docs.windowDocUri, new lexical.Location(position.lineNumber, position.column));
};
//# sourceMappingURL=main.js.map