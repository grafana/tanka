"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const fs = require("fs");
const url = require("url");
const im = require("immutable");
const editor = require("../compiler/editor");
const lexer = require("../compiler/lexical-analysis/lexer");
const lexical = require("../compiler/lexical-analysis/lexical");
const parser = require("../compiler/lexical-analysis/parser");
const _static = require("../compiler/static");
class VsDocumentManager {
    constructor(documents, libResolver) {
        this.documents = documents;
        this.libResolver = libResolver;
        this.get = (fileSpec) => {
            const parsedFileUri = this.libResolver.resolvePath(fileSpec);
            if (parsedFileUri == null) {
                throw new Error(`Could not open file`);
            }
            const fileUri = parsedFileUri.href;
            const filePath = parsedFileUri.path;
            if (fileUri == null || filePath == null) {
                throw new Error(`INTERNAL ERROR: ill-formed, null href or path`);
            }
            const version = fs.statSync(filePath).mtime.valueOf();
            const doc = this.documents.get(fileUri);
            if (doc == null) {
                const doc = this.fsCache.get(fileUri);
                if (doc != null && version == doc.version) {
                    // Return cached version if modified time is the same.
                    return {
                        text: doc.text,
                        version: doc.version,
                        resolvedPath: fileUri,
                    };
                }
                // Else, cache it.
                const text = fs.readFileSync(filePath).toString();
                const cached = {
                    text: text,
                    version: version,
                    resolvedPath: fileUri
                };
                this.fsCache = this.fsCache.set(fileUri, cached);
                return cached;
            }
            else {
                // Delete from `fsCache` just in case we were `import`'ing a
                // file and have since opened it.
                this.fsCache = this.fsCache.delete(fileUri);
                return {
                    text: doc.getText(),
                    version: doc.version,
                    resolvedPath: fileUri,
                };
            }
        };
        this.fsCache = im.Map();
    }
}
exports.VsDocumentManager = VsDocumentManager;
class VsCompilerService {
    constructor() {
        //
        // CompilerService implementation.
        //
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
exports.VsCompilerService = VsCompilerService;
class VsPathResolver extends editor.LibPathResolver {
    constructor() {
        super(...arguments);
        this.pathExists = (path) => {
            try {
                return fs.existsSync(path);
            }
            catch (err) {
                return false;
            }
        };
    }
}
exports.VsPathResolver = VsPathResolver;
//# sourceMappingURL=local.js.map