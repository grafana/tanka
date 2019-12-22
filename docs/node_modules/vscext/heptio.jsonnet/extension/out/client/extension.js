"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const child_process_1 = require("child_process");
const client = require("vscode-languageclient");
const os = require("os");
const path = require("path");
const fs = require("fs");
const vs = require("vscode");
const yaml = require("js-yaml");
const im = require("immutable");
const lexical = require("../compiler/lexical-analysis/lexical");
// activate registers the Jsonnet language server with vscode, and
// configures it based on the contents of the workspace JSON file.
exports.activate = (context) => {
    register.jsonnetClient(context);
    const diagProvider = register.diagnostics(context);
    register.previewCommands(context, diagProvider);
};
exports.deactivate = () => { };
var register;
(function (register) {
    // jsonnetClient registers the Jsonnet language client with vscode.
    register.jsonnetClient = (context) => {
        // The server is implemented in node
        let languageClient = jsonnet.languageClient(context.asAbsolutePath(path.join('out', 'server', 'server.js')));
        // Push the disposable to the context's subscriptions so that the
        // client can be deactivated on extension deactivation
        context.subscriptions.push(languageClient.start());
        // Configure the workspace.
        workspace.configure(vs.workspace.getConfiguration('jsonnet'));
    };
    // diagnostics registers a `jsonnet.DiagnosticProvider` with vscode.
    // This will cause vscode to render errors and warnings for users as
    // they save their code.
    register.diagnostics = (context) => {
        const diagnostics = vs.languages.createDiagnosticCollection("jsonnet");
        context.subscriptions.push(diagnostics);
        return new jsonnet.DiagnosticProvider(diagnostics);
    };
    // previewCommands will register the commands that allow people to
    // open a "preview" pane that renders their Jsonnet, similar to the
    // markdown preview pane.
    register.previewCommands = (context, diagProvider) => {
        // Create Jsonnet provider, register it to provide for documents
        // with `PREVIEW_SCHEME` URI scheme.
        const docProvider = new jsonnet.DocumentProvider();
        const registration = vs.workspace.registerTextDocumentContentProvider(jsonnet.PREVIEW_SCHEME, docProvider);
        // Subscribe to document updates. This allows us to detect (e.g.)
        // when a document was saved.
        context.subscriptions.push(registration);
        // Expand Jsonnet, register errors as diagnostics with vscode, and
        // generate preview if a preview tab is open.
        const preview = (doc) => {
            if (doc.languageId === "jsonnet") {
                const result = docProvider.cachePreview(doc);
                if (jsonnet.isRuntimeFailure(result)) {
                    diagProvider.report(doc.uri, result.error);
                }
                else {
                    diagProvider.clear(doc.uri);
                }
                docProvider.update(jsonnet.canonicalPreviewUri(doc.uri));
            }
        };
        // Register Jsonnet preview commands.
        context.subscriptions.push(vs.commands.registerCommand('jsonnet.previewToSide', () => display.previewJsonnet(true)));
        context.subscriptions.push(vs.commands.registerCommand('jsonnet.preview', () => display.previewJsonnet(false)));
        // Call `preview` any time we save or open a document.
        context.subscriptions.push(vs.workspace.onDidSaveTextDocument(preview));
        context.subscriptions.push(vs.workspace.onDidOpenTextDocument(preview));
        context.subscriptions.push(vs.workspace.onDidCloseTextDocument(doc => {
            docProvider.delete(doc);
        }));
        // Call `preview` when we open the editor.
        const active = vs.window.activeTextEditor;
        if (active != null) {
            preview(active.document);
        }
    };
})(register || (register = {}));
var workspace;
(function (workspace) {
    const extStrsProp = "extStrs";
    const execPathProp = "executablePath";
    workspace.extStrs = () => {
        const extStrsObj = vs.workspace.getConfiguration('jsonnet')[extStrsProp];
        return extStrsObj == null
            ? ""
            : Object.keys(extStrsObj)
                .map(key => `--ext-str ${key}="${extStrsObj[key]}"`)
                .join(" ");
    };
    workspace.libPaths = () => {
        const libPaths = vs.workspace.getConfiguration('jsonnet')["libPaths"];
        if (libPaths == null) {
            return "";
        }
        // Add executable to the beginning of the library paths, because
        // the Jsonnet CLI will look there first.
        //
        // TODO(hausdorff): Consider adding support for Jsonnet's
        // (undocumented) search paths `/usr/share/{jsonnet version}` and
        // `/usr/local/share/{jsonnet version}`. We don't support them
        // currently because (1) they're undocumented and therefore not
        // widely-used, and (2) it requires shelling out to the Jsonnet
        // command line, which complicates the extension.
        const jsonnetExecutable = vs.workspace.getConfiguration[execPathProp];
        if (jsonnetExecutable != null) {
            libPaths.unshift(jsonnetExecutable);
        }
        return libPaths
            .map(path => `-J ${path}`)
            .join(" ");
    };
    workspace.outputFormat = () => {
        return vs.workspace.getConfiguration('jsonnet')["outputFormat"];
    };
    workspace.configure = (config) => {
        if (os.type() === "Windows_NT") {
            return configureWindows(config);
        }
        else {
            return configureUnix(config);
        }
    };
    const configureUnix = (config) => {
        if (config[execPathProp] != null) {
            jsonnet.executable = config[execPathProp];
        }
        else {
            try {
                // If this doesn't throw, 'jsonnet' was found on
                // $PATH.
                //
                // TODO: Probably should find a good non-shell way of
                // doing this.
                child_process_1.execSync(`which jsonnet`);
            }
            catch (e) {
                alert.jsonnetCommandNotOnPath();
                return false;
            }
        }
        return true;
    };
    const configureWindows = (config) => {
        if (config[execPathProp] == null) {
            alert.jsonnetCommandIsNull();
            return false;
        }
        jsonnet.executable = config[execPathProp];
        return true;
    };
})(workspace || (workspace = {}));
var alert;
(function (alert_1) {
    const alert = vs.window.showErrorMessage;
    alert_1.noActiveWindow = () => {
        alert("Can't open Jsonnet preview because there is no active window");
    };
    alert_1.documentNotJsonnet = (languageId) => {
        alert(`Can't generate Jsonnet document preview for document with language id '${languageId}'`);
    };
    alert_1.couldNotRenderJsonnet = (reason) => {
        alert(`Error: Could not render Jsonnet; ${reason}`);
    };
    alert_1.jsonnetCommandNotOnPath = () => {
        alert(`Error: could not find 'jsonnet' command on path`);
    };
    alert_1.jsonnetCommandIsNull = () => {
        alert(`Error: 'jsonnet.executablePath' must be set in vscode settings`);
    };
})(alert || (alert = {}));
var html;
(function (html) {
    html.body = (body) => {
        return `<html><body>${body}</body></html>`;
    };
    html.codeLiteral = (code) => {
        return `<pre><code>${code}</code></pre>`;
    };
    html.errorMessage = (message) => {
        return `<i><pre>${message}</pre></i>`;
    };
    html.prettyPrintObject = (json, outputFormat) => {
        if (outputFormat == "yaml") {
            return html.codeLiteral(yaml.safeDump(JSON.parse(json)));
        }
        else {
            return html.codeLiteral(JSON.stringify(JSON.parse(json), null, 4));
        }
    };
})(html || (html = {}));
var jsonnet;
(function (jsonnet) {
    jsonnet.executable = "jsonnet";
    jsonnet.PREVIEW_SCHEME = "jsonnet-preview";
    jsonnet.DOCUMENT_FILTER = {
        language: 'jsonnet',
        scheme: 'file'
    };
    jsonnet.languageClient = (serverModule) => {
        // The debug options for the server
        let debugOptions = { execArgv: ["--nolazy", "--inspect=6009"] };
        // If the extension is launched in debug mode then the debug
        // server options are used. Otherwise the run options are used
        let serverOptions = {
            run: {
                module: serverModule,
                transport: client.TransportKind.ipc,
            },
            debug: {
                module: serverModule,
                transport: client.TransportKind.ipc,
                options: debugOptions
            }
        };
        // Options to control the language client
        let clientOptions = {
            // Register the server for plain text documents
            documentSelector: [jsonnet.DOCUMENT_FILTER.language],
            synchronize: {
                // Synchronize the workspace/user settings sections
                // prefixed with 'jsonnet' to the server.
                configurationSection: jsonnet.DOCUMENT_FILTER.language,
                // Notify the server about file changes to '.clientrc
                // files contain in the workspace.
                fileEvents: vs.workspace.createFileSystemWatcher('**/.clientrc')
            }
        };
        // Create the language client and start the client.
        return new client.LanguageClient("JsonnetLanguageServer", 'Jsonnet Language Server', serverOptions, clientOptions);
    };
    jsonnet.canonicalPreviewUri = (fileUri) => {
        return fileUri.with({
            scheme: jsonnet.PREVIEW_SCHEME,
            path: `${fileUri.path}.rendered`,
            query: fileUri.toString(),
        });
    };
    jsonnet.fileUriFromPreviewUri = (previewUri) => {
        const file = previewUri.fsPath.slice(0, -(".rendered".length));
        return vs.Uri.file(file);
    };
    // RuntimeError represents a runtime failure in a Jsonnet program.
    class RuntimeFailure {
        constructor(error) {
            this.error = error;
        }
    }
    jsonnet.RuntimeFailure = RuntimeFailure;
    jsonnet.isRuntimeFailure = (thing) => {
        return thing instanceof RuntimeFailure;
    };
    // DocumentProvider compiles Jsonnet code to JSON or YAML, and
    // provides that to vscode for rendering in the preview pane.
    //
    // DESIGN NOTES: This class optionally exposes `cachePreview` and
    // `delete` so that the caller can get the results of the document
    // compilation for purposes of (e.g.) reporting diagnostic issues.
    class DocumentProvider {
        constructor() {
            this.provideTextDocumentContent = (previewUri) => {
                const sourceUri = vs.Uri.parse(previewUri.query);
                return vs.workspace.openTextDocument(sourceUri)
                    .then(sourceDoc => {
                    const result = this.previewCache.has(sourceUri.toString())
                        ? this.previewCache.get(sourceUri.toString())
                        : this.cachePreview(sourceDoc);
                    if (jsonnet.isRuntimeFailure(result)) {
                        return html.body(html.errorMessage(result.error));
                    }
                    const outputFormat = workspace.outputFormat();
                    return html.body(html.prettyPrintObject(result, outputFormat));
                });
            };
            this.cachePreview = (sourceDoc) => {
                const sourceUri = sourceDoc.uri.toString();
                const sourceFile = sourceDoc.uri.fsPath;
                let codePaths = '';
                if (ksonnet.isInApp(sourceFile)) {
                    const dir = path.dirname(sourceFile);
                    const paramsPath = path.join(dir, "params.libsonnet");
                    const rootDir = ksonnet.rootPath(sourceFile);
                    const envParamsPath = path.join(rootDir, "environments", "default", "params.libsonnet");
                    let codeImports = {
                        '__ksonnet/params': path.join(dir, "params.libsonnet"),
                        '__ksonnet/environments': envParamsPath,
                    };
                    codePaths = Object.keys(codeImports)
                        .map(k => `--ext-code-file "${k}"=${codeImports[k]}`)
                        .join(' ');
                    console.log(codePaths);
                }
                try {
                    // Compile the preview Jsonnet file.
                    const extStrs = workspace.extStrs();
                    const libPaths = workspace.libPaths();
                    const jsonOutput = child_process_1.execSync(`${jsonnet.executable} ${libPaths} ${extStrs} ${codePaths} ${sourceFile}`).toString();
                    // Cache.
                    this.previewCache = this.previewCache.set(sourceUri, jsonOutput);
                    return jsonOutput;
                }
                catch (e) {
                    const failure = new RuntimeFailure(e.message);
                    this.previewCache = this.previewCache.set(sourceUri, failure);
                    return failure;
                }
            };
            this.delete = (document) => {
                const previewUri = document.uri.query.toString();
                this.previewCache = this.previewCache.delete(previewUri);
            };
            this.update = (uri) => {
                this._onDidChange.fire(uri);
            };
            //
            // Private members.
            //
            this._onDidChange = new vs.EventEmitter();
            this.previewCache = im.Map();
        }
        //
        // Document update API.
        //
        get onDidChange() {
            return this._onDidChange.event;
        }
    }
    jsonnet.DocumentProvider = DocumentProvider;
    // DiagnosticProvider will consume the output of the Jsonnet CLI and
    // either (1) report diagnostics issues (e.g., errors, warnings) to
    // the user, or (2) clear them if the compilation was successful.
    class DiagnosticProvider {
        constructor(diagnostics) {
            this.diagnostics = diagnostics;
            this.report = (fileUri, message) => {
                const messageLines = im.List(message.split(os.EOL)).rest();
                // Start over.
                this.diagnostics.clear();
                const errorMessage = messageLines.get(0);
                if (errorMessage.startsWith(lexical.staticErrorPrefix)) {
                    return this.reportStaticErrorDiagnostics(errorMessage);
                }
                else if (errorMessage.startsWith(lexical.runtimeErrorPrefix)) {
                    const stackTrace = messageLines.rest().toList();
                    return this.reportRuntimeErrorDiagnostics(fileUri, errorMessage, stackTrace);
                }
            };
            this.clear = (fileUri) => {
                this.diagnostics.delete(fileUri);
            };
            //
            // Private members.
            //
            this.reportStaticErrorDiagnostics = (message) => {
                const staticError = message.slice(lexical.staticErrorPrefix.length);
                const match = DiagnosticProvider.fileFromStackFrame(staticError);
                if (match == null) {
                    console.log(`Could not parse filename from Jsonnet error: '${message}'`);
                    return;
                }
                const locAndMessage = staticError.slice(match.fullMatch.length);
                const range = DiagnosticProvider.parseRange(locAndMessage);
                if (range == null) {
                    console.log(`Could not parse location range from Jsonnet error: '${message}'`);
                    return;
                }
                const diag = new vs.Diagnostic(range, locAndMessage, vs.DiagnosticSeverity.Error);
                this.diagnostics.set(vs.Uri.file(match.file), [diag]);
            };
            this.reportRuntimeErrorDiagnostics = (fileUri, message, messageLines) => {
                const diagnostics = messageLines
                    .reduce((acc, line) => {
                    // Filter error lines that we know aren't stack frames.
                    const trimmed = line.trim();
                    if (trimmed == "" || trimmed.startsWith("During manifestation")) {
                        return acc;
                    }
                    // Log when we think a line is a stack frame, but we can't
                    // parse it.
                    const match = DiagnosticProvider.fileFromStackFrame(line);
                    if (match == null) {
                        console.log(`Could not parse filename from Jsonnet error: '${line}'`);
                        return acc;
                    }
                    const loc = line.slice(match.fileWithLeadingWhitespace.length);
                    const range = DiagnosticProvider.parseRange(loc);
                    if (range == null) {
                        console.log(`Could not parse filename from Jsonnet error: '${line}'`);
                        return acc;
                    }
                    // Generate and emit diagnostics.
                    const diag = new vs.Diagnostic(range, `${message}`, vs.DiagnosticSeverity.Error);
                    const prev = acc.get(match.file, undefined);
                    return prev == null
                        ? acc.set(match.file, im.List([diag]))
                        : acc.set(match.file, prev.push(diag));
                }, im.Map());
                const fileDiags = diagnostics.get(fileUri.fsPath, undefined);
                fileDiags != null && this.diagnostics.set(fileUri, fileDiags.toArray());
            };
        }
    }
    DiagnosticProvider.parseRange = (range) => {
        const lr = lexical.LocationRange.fromString("Dummy name", range);
        if (lr == null) {
            return null;
        }
        const start = new vs.Position(lr.begin.line - 1, lr.begin.column - 1);
        // NOTE: Don't subtract 1 from `lr.end.column` because the range
        // is exclusive at the end.
        const end = new vs.Position(lr.end.line - 1, lr.end.column);
        return new vs.Range(start, end);
    };
    DiagnosticProvider.fileFromStackFrame = (frameMessage) => {
        const fileMatch = frameMessage.match(/(\s*)(.*?):/);
        return fileMatch == null
            ? null
            : {
                fullMatch: fileMatch[0],
                fileWithLeadingWhitespace: fileMatch[1] + fileMatch[2],
                file: fileMatch[2],
            };
    };
    jsonnet.DiagnosticProvider = DiagnosticProvider;
})(jsonnet || (jsonnet = {}));
var display;
(function (display) {
    display.previewJsonnet = (sideBySide) => {
        const editor = vs.window.activeTextEditor;
        if (editor == null) {
            alert.noActiveWindow();
            return;
        }
        const languageId = editor.document.languageId;
        if (!(editor.document.languageId === "jsonnet")) {
            alert.documentNotJsonnet(languageId);
            return;
        }
        const previewUri = jsonnet.canonicalPreviewUri(editor.document.uri);
        return vs.commands.executeCommand('vscode.previewHtml', previewUri, display.getViewColumn(sideBySide), `Jsonnet preview '${path.basename(editor.document.fileName)}'`).then((success) => { }, (reason) => {
            alert.couldNotRenderJsonnet(reason);
        });
    };
    display.getViewColumn = (sideBySide) => {
        const active = vs.window.activeTextEditor;
        if (!active) {
            return vs.ViewColumn.One;
        }
        if (!sideBySide) {
            return active.viewColumn;
        }
        switch (active.viewColumn) {
            case vs.ViewColumn.One:
                return vs.ViewColumn.Two;
            case vs.ViewColumn.Two:
                return vs.ViewColumn.Three;
        }
        return active.viewColumn;
    };
})(display || (display = {}));
var ksonnet;
(function (ksonnet) {
    // find the root of the components structure.
    function isInApp(filePath, fsRoot = '/') {
        const currentPath = path.join(fsRoot, filePath);
        return checkForKsonnet(currentPath);
    }
    ksonnet.isInApp = isInApp;
    function rootPath(filePath, fsRoot = '/') {
        const currentPath = path.join(fsRoot, filePath);
        return findRootPath(currentPath);
    }
    ksonnet.rootPath = rootPath;
    function checkForKsonnet(filePath) {
        if (filePath === "/") {
            return false;
        }
        const dir = path.dirname(filePath);
        const parts = dir.split(path.sep);
        if (parts[parts.length - 1] === "components") {
            const root = path.dirname(dir);
            const ksConfig = path.join(root, "app.yaml");
            try {
                const stats = fs.statSync(ksConfig);
                return true;
            }
            catch (err) {
                return false;
            }
        }
        return checkForKsonnet(dir);
    }
    function findRootPath(filePath) {
        if (filePath === "/") {
            return '';
        }
        const dir = path.dirname(filePath);
        const parts = dir.split(path.sep);
        if (parts[parts.length - 1] === "components") {
            const root = path.dirname(dir);
            const ksConfig = path.join(root, "app.yaml");
            try {
                const stats = fs.statSync(ksConfig);
                return root;
            }
            catch (err) {
                return '';
            }
        }
        return findRootPath(dir);
    }
})(ksonnet = exports.ksonnet || (exports.ksonnet = {}));
//# sourceMappingURL=extension.js.map