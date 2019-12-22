"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const server = require("vscode-languageserver");
const im = require("immutable");
const ast = require("../compiler/lexical-analysis/ast");
const _static = require("../compiler/static");
// fromFailure creates a diagnostic from a `LexFailure |
// ParseFailure`.
exports.fromFailure = (error) => {
    let begin = null;
    let end = null;
    let message = null;
    if (_static.isLexFailure(error)) {
        begin = error.lexError.loc.begin;
        end = error.lexError.loc.end;
        message = error.lexError.msg;
    }
    else {
        begin = error.parseError.loc.begin;
        end = error.parseError.loc.end;
        message = error.parseError.msg;
    }
    return {
        severity: server.DiagnosticSeverity.Error,
        range: {
            start: { line: begin.line - 1, character: begin.column - 1 },
            end: { line: end.line - 1, character: end.column - 1 },
        },
        message: `${message}`,
        source: `Jsonnet`,
    };
};
// fromAst takes a Jsonnet AST and returns an array of `Diagnostic`
// issues it finds.
exports.fromAst = (root, libResolver) => {
    const diags = new Visitor(root, libResolver);
    diags.visit();
    return diags.diagnostics;
};
// ----------------------------------------------------------------------------
// Private utilities.
// ----------------------------------------------------------------------------
// Visitor traverses the Jsonnet AST and accumulates `Diagnostic`
// errors for reporting.
class Visitor extends ast.VisitorBase {
    constructor(root, libResolver) {
        super(root);
        this.libResolver = libResolver;
        this.diags = im.List();
        this.visitImport = (node) => this.importDiagnostics(node);
        this.visitImportStr = (node) => this.importDiagnostics(node);
        this.importDiagnostics = (node) => {
            if (!this.libResolver.resolvePath(node)) {
                const begin = node.loc.begin;
                const end = node.loc.end;
                const diagnostic = {
                    severity: server.DiagnosticSeverity.Warning,
                    range: {
                        start: { line: begin.line - 1, character: begin.column - 1 },
                        end: { line: end.line - 1, character: end.column - 1 },
                    },
                    message: `Can't find path '${node.file}'. If the file is not in the ` +
                        `current directory, it may be necessary to add it to the ` +
                        `'jsonnet.libPaths'. If you are in vscode, you can press ` +
                        `'cmd/ctrl-,' and add the path this library is located at to the ` +
                        `'jsonnet.libPaths' array`,
                    source: `Jsonnet`,
                };
                this.diags = this.diags.push(diagnostic);
            }
        };
    }
    get diagnostics() {
        return this.diags.toArray();
    }
}
//# sourceMappingURL=diagnostic.js.map