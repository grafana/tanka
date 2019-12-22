import * as fs from 'fs';
import * as path from 'path';
import * as proc from 'child_process';
import * as url from 'url';

import * as im from 'immutable';
import * as server from 'vscode-languageserver';

import * as ast from '../compiler/lexical-analysis/ast';
import * as editor from '../compiler/editor';
import * as lexer from '../compiler/lexical-analysis/lexer';
import * as lexical from '../compiler/lexical-analysis/lexical';
import * as parser from '../compiler/lexical-analysis/parser';
import * as _static from "../compiler/static";

export class VsDocumentManager implements editor.DocumentManager {
  constructor(
    private readonly documents: server.TextDocuments,
    private readonly libResolver: editor.LibPathResolver,
  ) { }

  get = (
    fileSpec: editor.FileUri | ast.Import | ast.ImportStr,
  ): {text: string, version?: number, resolvedPath: string} => {
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
    } else {
      // Delete from `fsCache` just in case we were `import`'ing a
      // file and have since opened it.
      this.fsCache = this.fsCache.delete(fileUri);
      return {
        text: doc.getText(),
        version: doc.version,
        resolvedPath: fileUri,
      }
    }
  }

  private fsCache = im.Map<string, {text: string, version: number}>();
}

export class VsCompilerService implements _static.LexicalAnalyzerService {
  //
  // CompilerService implementation.
  //

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

export class VsPathResolver extends editor.LibPathResolver {
  protected pathExists = (path: string): boolean => {
    try {
      return fs.existsSync(path);
    } catch (err) {
      return false;
    }
  }
}
