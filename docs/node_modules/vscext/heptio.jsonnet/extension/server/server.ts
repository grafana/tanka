import * as fs from 'fs';
import * as os from 'os';
import * as path from 'path';
import * as server from 'vscode-languageserver';
import * as url from 'url';

import * as im from 'immutable';

import * as ast from '../compiler/lexical-analysis/ast';
import * as diagnostic from './diagnostic';
import * as editor from '../compiler/editor';
import * as lexer from '../compiler/lexical-analysis/lexer';
import * as lexical from '../compiler/lexical-analysis/lexical';
import * as local from './local';
import * as _static from '../compiler/static';

// Create a connection for the server. The connection uses Node's IPC
// as a transport
const connection: server.IConnection = server.createConnection(
  new server.IPCMessageReader(process),
  new server.IPCMessageWriter(process));

// Create a simple text document manager. The text document manager
// supports full document sync only
const docs = new server.TextDocuments();

const compiler = new local.VsCompilerService();

const libResolver = new local.VsPathResolver();
const documentManager = new local.VsDocumentManager(docs, libResolver);

const analyzer: _static.EventedAnalyzer = new _static.Analyzer(
  documentManager, compiler);

const reportDiagnostics = (doc: server.TextDocument) => {
  const text = doc.getText();
  const results = compiler.cache(doc.uri, text, doc.version);

  if (_static.isParsedDocument(results)) {
    connection.sendDiagnostics({
      uri: doc.uri,
      diagnostics: diagnostic.fromAst(results.parse, libResolver),
    });
  } else {
    connection.sendDiagnostics({
      uri: doc.uri,
      diagnostics: [diagnostic.fromFailure(results.parse)],
    });
  }
};

//
// TODO: We should find a way to move these hooks to a "init doc
// manager" method, or something.
//
// TODO: We should abstract over these hooks with
// `workspace.DocumentManager`.
//

docs.onDidOpen(openEvent => {
  const doc = openEvent.document;
  if (doc.languageId === "jsonnet") {
    reportDiagnostics(doc);
    return analyzer.onDocumentOpen(doc.uri, doc.getText(), doc.version);
  }
});
docs.onDidSave(saveEvent => {
  // TODO: relay once we can get the "last good" parse, we can check
  // the env of the position, and then evaluate that identifier, or
  // splice it into the tree. We can perhaps split changes into those
  // that are single-line, and those that are multi-line.
  const doc = saveEvent.document;
  if (doc.languageId === "jsonnet") {
    reportDiagnostics(doc);
    return analyzer.onDocumentOpen(doc.uri, doc.getText(), doc.version);
  }
});
docs.onDidClose(closeEvent => {
  // TODO: This is a bit simplistic. We'll need to have a graph of
  // files, eventually, so that we can reload any previews whose
  // dependencies we save.
  if (closeEvent.document.languageId === "jsonnet") {
    return analyzer.onDocumentClose(closeEvent.document.uri);
  }
});
docs.onDidChangeContent(changeEvent => {
  if (changeEvent.document.languageId === "jsonnet") {
    reportDiagnostics(changeEvent.document);
  }
});

// Make the text document manager listen on the connection
// for open, change and close text document events
docs.listen(connection);

connection.onInitialize((params) => initializer(docs, params));
connection.onDidChangeConfiguration(params => configUpdateProvider(params));
connection.onHover(position => {
  const fileUri = position.textDocument.uri;
  return analyzer.onHover(fileUri, positionToLocation(position));
});

connection.onCompletion(position => {
  return analyzer
    .onComplete(position.textDocument.uri, positionToLocation(position))
    .then<server.CompletionItem[]>(
      completions => completions.map(completionInfoToCompletionItem));
});
// Prevent the language server from complaining that
// `onCompletionResolve` handle is not implemented.
connection.onCompletionResolve(item => item);

// Listen on the connection
connection.listen();


export const initializer = (
  documents: server.TextDocuments,
  params: server.InitializeParams,
): server.InitializeResult => {
  return {
    capabilities: {
      // Tell the client that the server works in FULL text
      // document sync mode
      textDocumentSync: documents.syncKind,
      // Tell the client that the server support code complete
      completionProvider: {
        resolveProvider: true,
        triggerCharacters: ["."],
      },
      hoverProvider: true,
    }
  }
}

export const configUpdateProvider = (
  change: server.DidChangeConfigurationParams,
): void => {
  if (
    change.settings == null || change.settings.jsonnet == null ||
    change.settings.jsonnet.libPaths == null
  ) {
    return;
  }

  const jsonnet = change.settings.jsonnet;

  if (jsonnet.executablePath != null) {
    jsonnet.libPaths.unshift(jsonnet.executablePath);
  }

  libResolver.libPaths = im.List<string>(jsonnet.libPaths);
}

const positionToLocation = (
  posParams: server.TextDocumentPositionParams
): lexical.Location => {
  return new lexical.Location(
    posParams.position.line + 1,
    posParams.position.character + 1);
}

const completionInfoToCompletionItem = (
  completionInfo: editor.CompletionInfo
): server.CompletionItem => {
    let kindMapping: server.CompletionItemKind;
    switch (completionInfo.kind) {
      case "Field": {
        kindMapping = server.CompletionItemKind.Field;
        break;
      }
      case "Variable": {
        kindMapping = server.CompletionItemKind.Variable;
        break;
      }
      case "Method": {
        kindMapping = server.CompletionItemKind.Method;
        break;
      }
      default: throw new Error(
        `Unrecognized completion type '${completionInfo.kind}'`);
    }

    // Black magic type coercion. This allows us to avoid doing a
    // deep copy over to a new `CompletionItem` object, and
    // instead only re-assign the `kindMapping`.
    const completionItem = (<server.CompletionItem>(<object>completionInfo));
    completionItem.kind = kindMapping;
    return completionItem;
}
