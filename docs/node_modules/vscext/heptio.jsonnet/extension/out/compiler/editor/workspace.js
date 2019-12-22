"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const path = require("path");
const url = require("url");
const im = require("immutable");
const ast = require("../lexical-analysis/ast");
;
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
class LibPathResolver {
    constructor() {
        this.wd = path.resolve(".");
        this._libPaths = im.List([this.wd]);
        this.resolvePath = (fileSpec) => {
            // IMPLEMENTATION NOTE: We're using the RFC 1630/1738 file URI
            // specification for specifying files, largely because that's what
            // vscode expects. Specifically, we use the `file:///${filename}`
            // pattern rather than the `file://localhost/${filename}` pattern.
            let importPath = null;
            if (ast.isImport(fileSpec) || ast.isImportStr(fileSpec)) {
                if (path.isAbsolute(fileSpec.file)) {
                    // If path is absolute and exists, it's resolved.
                    importPath = this.pathExists(fileSpec.file)
                        ? `file://${fileSpec.file}`
                        : null;
                }
                else {
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
                    const pathToImportedFile = path.dirname(path.resolve(fileSpec.loc.fileName));
                    const paths = im.List([pathToImportedFile])
                        .concat(this.libPaths)
                        .toList();
                    importPath = this.searchPaths(fileSpec.file, paths);
                    // NOTE: Failing to set `importPath` at this point will cause
                    // us to return `null` below.
                }
            }
            else {
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
        };
        //---------------------------------------------------------------------------
        // Private members.
        //---------------------------------------------------------------------------
        this.searchPaths = (importPath, paths) => {
            for (let libPath of paths.toArray()) {
                try {
                    const resolvedPath = path.join(libPath, importPath);
                    if (this.pathExists(resolvedPath)) {
                        return path.isAbsolute(resolvedPath)
                            ? `file://${resolvedPath}`
                            : `file:///${resolvedPath}`;
                    }
                }
                catch (err) {
                    // Ignore.
                }
            }
            return null;
        };
    }
    set libPaths(libPaths) {
        this._libPaths = im.List([this.wd])
            .concat(libPaths)
            .toList();
    }
    get libPaths() {
        return this._libPaths;
    }
}
exports.LibPathResolver = LibPathResolver;
//# sourceMappingURL=workspace.js.map