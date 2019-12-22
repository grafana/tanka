import * as path from 'path';
import * as url from 'url';

import * as im from 'immutable';

import * as ast from '../lexical-analysis/ast';

// DocumentEventListener listens to events emitted by a
// `DocumentManager` in response to changes to documents that it
// manages. For example, if a document is saved, an `Save` event
// would be fired by the `DocumentManager`, and subsequently processed
// by a hook registered with the `DocumentEventListener`.
export interface DocumentEventListener {
  onDocumentOpen: (uri: string, text: string, version?: number) => void
  onDocumentSave: (uri: string, text: string, version?: number) => void
  onDocumentClose: (uri: string) => void
};

export type FileUri = string;

// DocumentManager typically provides 2 important pieces of
// functionality:
//
// 1. It is the system of record for documents managed in a
//   "workspace"; if a document exists in the workspace, it should be
//   possible to `get` it by providing a `fileUri`. For example, in
//   the context of vscode, this should wrap an instance of
//   `TextDocuments`, which manages changes for all documents in a
//   vscode workspace.
// 2. When a document that is managed by the `DocumentManager` is
//   changed, we should be firing off an event, so that the
//   `DocumentEventListener` can call the appropriate hook. This is
//   important, as it allows users to (e.g.) update parse caches,
//   which allows the client to provide efficient support for
//   features like autocomplete.
//
//   IMPORTANT NOTE: Right now, this behavior is completely implicit.
//   This interface does not currently contain functions that express
//   hook registration (e.g., as `TextDocuments#onDidSave` does in
//   the case of vscode). This means that it is incumbent on the user
//   to actually implement this functionality and hook it up
//   correctly to the `DocumentEventListener`.
export interface DocumentManager {
  get: (
    fileSpec: FileUri | ast.Import | ast.ImportStr,
  ) => {text: string, version?: number, resolvedPath: string}
}

// LibPathResolver searches a set of library paths for a Jsonnet
// library file specified either by a `FileUri`, or AST nodes of type
// `Import` or `ImportStr`, returning an absolute, fully-qualified
// file path if (and only if) the path exists. These "resolved" files
// are represented with RFC 1630/1738-compliant file URIs, which can
// be consumed by (e.g.) vscode.
//
// LibPathResolver is an abstract representation of the problem of
// searching a set of paths for a Jsonnet, in the sense that only the
// search logic is specified, independent of any specific filesystem.
// All filesystem-specific logic (either real or an FS mock) must be
// specified by implementing `LibPathResolver#pathExists`.
export abstract class LibPathResolver {
  private readonly wd = path.resolve(".");
  private _libPaths = im.List<string>([this.wd]);

  set libPaths(libPaths: im.List<string>) {
    this._libPaths = im.List<string>([this.wd])
      .concat(libPaths)
      .toList();
  }

  get libPaths(): im.List<string> {
    return this._libPaths;
  }

  public resolvePath = (
    fileSpec: FileUri | ast.Import | ast.ImportStr,
  ): url.Url | null => {
    // IMPLEMENTATION NOTE: We're using the RFC 1630/1738 file URI
    // specification for specifying files, largely because that's what
    // vscode expects. Specifically, we use the `file:///${filename}`
    // pattern rather than the `file://localhost/${filename}` pattern.

    let importPath: string | null = null;
    if (ast.isImport(fileSpec) || ast.isImportStr(fileSpec)) {
      if (path.isAbsolute(fileSpec.file)) {
        // If path is absolute and exists, it's resolved.
        importPath = this.pathExists(fileSpec.file)
          ? `file://${fileSpec.file}`
          : null;
      } else {
        // Else the `import` path is either:
        //
        // 1. relative to the file that's doing the importing, or
        // 2. relative to one of the `libPaths`.
        //
        // If neither is true, fail.
        //
        // TODO(hausdorff): I think this might be a bug. The filename
        // might not be relative to workspace root, in which case
        // `resolve` seems like it should fail.
        const pathToImportedFile =
          path.dirname(path.resolve(fileSpec.loc.fileName));

        const paths = im.List<string>([pathToImportedFile])
          .concat(this.libPaths)
          .toList();

        importPath = this.searchPaths(fileSpec.file, paths);

        // NOTE: Failing to set `importPath` at this point will cause
        // us to return `null` below.
      }
    } else {
      // NOTE: No need to convert to URI, it was passed in as
      // `FileUri`.
      importPath = fileSpec;
    }

    if (importPath == null) {
      return null;
    }

    const parsed = url.parse(importPath);
    if (!parsed || !parsed.path || parsed.protocol !== "file:") {
      throw new Error(`INTERNAL ERROR: Failed to parse URI '${fileSpec}'`);
    }
    return parsed;
  }

  //---------------------------------------------------------------------------
  // Protected members.
  //---------------------------------------------------------------------------

  protected pathExists: (path: string) => boolean;

  //---------------------------------------------------------------------------
  // Private members.
  //---------------------------------------------------------------------------

  private searchPaths = (
    importPath: string, paths: im.List<string>,
  ): FileUri | null => {
    for (let libPath of paths.toArray()) {
      try {
        const resolvedPath = path.join(libPath, importPath);
        if (this.pathExists(resolvedPath)) {
          return path.isAbsolute(resolvedPath)
            ? `file://${resolvedPath}`
            : `file:///${resolvedPath}`;
        }
      } catch (err) {
        // Ignore.
      }
    }

    return null;
  }
}
