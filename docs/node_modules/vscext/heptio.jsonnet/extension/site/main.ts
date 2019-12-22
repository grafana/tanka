import * as http from 'http';
import * as url from 'url';

import * as im from 'immutable';

import * as ast from '../compiler/lexical-analysis/ast';
import * as editor from '../compiler/editor';
import * as lexer from '../compiler/lexical-analysis/lexer';
import * as lexical from '../compiler/lexical-analysis/lexical';
import * as parser from '../compiler/lexical-analysis/parser';
import * as _static from '../compiler/static';

declare var global: any;

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

const getUrl = (url: string, cb) => {
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

export class BrowserDocumentManager implements editor.DocumentManager {
  public k: string;
  public k8s: string;

  // URI utilities.
  public readonly backsplicePrefix = `file:///`;
  public readonly windowDocUri = `${this.backsplicePrefix}window`;
  public readonly k8sUri = `${this.backsplicePrefix}k8s.libsonnet`;

  // The ksonnet files (e.g., `k.libsonnet`) never change. Because
  // static analysis of their files is expensive, we assign them
  // version 0 so that it's _always_ cached.
  public readonly staticsVersion = 0;

  get = (
    file: editor.FileUri | ast.Import | ast.ImportStr,
  ): {text: string, version?: number, resolvedPath: string} => {
    const fileUri = ast.isImport(file) || ast.isImportStr(file)
      ? `${this.backsplicePrefix}${file.file}`
      : file;

    if (fileUri === `${this.backsplicePrefix}ksonnet.beta.2/k.libsonnet`) {
      return {
        text: this.k,
        version: this.staticsVersion,
        resolvedPath: fileUri,
      };
    } else if (fileUri === this.k8sUri) {
      return {
        text: this.k8s,
        version: this.staticsVersion,
        resolvedPath: fileUri,
      };
    } else if (fileUri === this.windowDocUri) {
      return {
        text: this.windowText,
        version: this.version,
        resolvedPath: fileUri,
      };
    }

    throw new Error(`Unrecognized file ${fileUri}`);
  }

  public setWindowText = (text: string, version?: number) => {
    this.windowText = text;
    this.version = version;
  }

  private windowText: string = "";
  private version?: number = undefined;
}

export class BrowserCompilerService implements _static.LexicalAnalyzerService {
  public cache = (
    fileUri: string, text: string, version?: number
  ): _static.ParsedDocument | _static.FailedParsedDocument => {
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
      tryGet.version === version
    ) {
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
  }

  public getLastSuccess = (
    fileUri: string
  ): _static.ParsedDocument | null => {
    return this.docCache.has(fileUri) && this.docCache.get(fileUri) || null;
  }

  public delete = (fileUri: string): void => {
    this.docCache = this.docCache.delete(fileUri);
  }

  //
  // Private members.
  //

  private docCache = im.Map<string, _static.ParsedDocument>();
}

// ----------------------------------------------------------------------------
// Set up analyzer in browser.
// ----------------------------------------------------------------------------

const docs = new BrowserDocumentManager();
const cs = new BrowserCompilerService();
const analyzer = new _static.Analyzer(docs, cs);

// ----------------------------------------------------------------------------
// Get ksonnet files.
// ----------------------------------------------------------------------------

getUrl(
  'https://raw.githubusercontent.com/ksonnet/ksonnet-lib/bd6b2d618d6963ea6a81fcc5623900d8ba110a32/ksonnet.beta.2/k.libsonnet',
  text => {docs.k = text;});
getUrl(
  "https://raw.githubusercontent.com/ksonnet/ksonnet-lib/bd6b2d618d6963ea6a81fcc5623900d8ba110a32/ksonnet.beta.2/k8s.libsonnet",
  text => {
    docs.k8s = text;
    // Static analysis on `k8s.libsonnet` takes multiple seconds to
    // complete, so do this immediately.
    cs.cache(docs.k8sUri, text, docs.staticsVersion);
  });

interface MonacoPosition {
  lineNumber: number,
  column: number,
};

// ----------------------------------------------------------------------------
// Public functions for the Monaco editor to call.
// ----------------------------------------------------------------------------

global.docOnChange = (text: string, version?: number) => {
  docs.setWindowText(text, version);
  cs.cache(docs.windowDocUri, text, version);
}

global.onComplete = (
  text: string, position: MonacoPosition
): Promise<editor.CompletionInfo[]> => {
  return analyzer
    .onComplete(
      docs.windowDocUri, new lexical.Location(position.lineNumber, position.column));
}
