import * as os from 'os';
import * as path from 'path';

import * as im from 'immutable';

import * as ast from '../lexical-analysis/ast';
import * as lexer from '../lexical-analysis/lexer';
import * as lexical from '../lexical-analysis/lexical';
import * as service from './service';
import * as editor from '../editor';

//
// Analyzer.
//

export interface EventedAnalyzer
  extends editor.DocumentEventListener, editor.UiEventListener { }

// TODO: Rename this to `EventedAnalyzer`.
export class Analyzer implements EventedAnalyzer {
  constructor(
    private documents: editor.DocumentManager,
    private compilerService: service.LexicalAnalyzerService,
  ) { }

  //
  // WorkspaceEventListener implementation.
  //

  public onDocumentOpen = this.compilerService.cache;

  public onDocumentSave = this.compilerService.cache;

  public onDocumentClose = this.compilerService.delete;

  //
  // AnalysisEventListener implementation.
  //

  public onHover = (
    fileUri: string, cursorLoc: lexical.Location
  ): Promise<editor.HoverInfo> => {
    const emptyOnHover = Promise.resolve().then(
      () => <editor.HoverInfo>{
        contents: [],
      });

    const onHoverPromise = (
      node: ast.Node | ast.IndexedObjectFields,
    ): Promise<editor.HoverInfo> => {
      if (node == null) {
        return emptyOnHover;
      }

      try {
        const msg = this.renderOnhoverMessage(fileUri, node);
        return Promise.resolve().then(
          () => <editor.HoverInfo> {
            contents: msg,
          });
      } catch(err) {
        console.log(err);
        return emptyOnHover;
      }
    }

    try {
      const {text: docText, version: version, resolvedPath: resolvedUri} =
        this.documents.get(fileUri);
      const cached = this.compilerService.cache(fileUri, docText, version);
      if (service.isFailedParsedDocument(cached)) {
        return emptyOnHover;
      }

      // Get symbol we're hovering over.
      const nodeAtPos = getNodeAtPositionFromAst(cached.parse, cursorLoc);
      if (ast.isFindFailure(nodeAtPos)) {
        return emptyOnHover;
      }

      if (nodeAtPos.parent != null && ast.isFunctionParam(nodeAtPos.parent)) {
        // A function parameter is a free variable, so we can't resolve
        // it. Simply return.
        return onHoverPromise(nodeAtPos.parent);
      }

      if (!ast.isResolvable(nodeAtPos)) {
        return emptyOnHover;
      }

      const ctx = new ast.ResolutionContext(
        this.compilerService, this.documents, resolvedUri);
      const resolved = ast.tryResolveIndirections(nodeAtPos, ctx);

      // Handle the special cases. If we hover over a symbol that points
      // at a function of some sort (i.e., a `function` literal, a
      // `local` that has a bind that is a function, or an object field
      // that is a function), then we want to render the name and
      // parameters that function takes, rather than the definition of
      // the function itself.
      if (ast.isResolveFailure(resolved)) {
        if (ast.isUnresolved(resolved) || ast.isUnresolvedIndex(resolved)) {
          return emptyOnHover;
        } else if (ast.isResolvedFreeVar(resolved)) {
          return onHoverPromise(resolved.variable);
        } else if (ast.isResolvedFunction(resolved)) {
          return onHoverPromise(resolved.functionNode);
        } else {
          return onHoverPromise(resolved);
        }
      } else {
        return onHoverPromise(resolved.value);
      }
    } catch (err) {
      console.log(err);
      return emptyOnHover;
    }
  }

  public onComplete = (
    fileUri: editor.FileUri, cursorLoc: lexical.Location
  ): Promise<editor.CompletionInfo[]> => {
    const doc = this.documents.get(fileUri);

    return Promise.resolve().then(
      (): editor.CompletionInfo[] => {
        //
        // Generate suggestions. This process follows three steps:
        //
        // 1. Try to parse the document text.
        // 2. If we succeed, go to cursor, select that node, and if
        //    it's an identifier that can be completed, then return
        //    the environment.
        // 3. If we fail, go try to go to the "hole" where the
        //    identifier exists.
        //

        try {
          const compiled = this.compilerService.cache(
            fileUri, doc.text, doc.version);
          const lines = doc.text.split("\n");

          // Lets us know whether the user has typed something like
          // `foo` or `foo.` (i.e., whether they are "dotting into"
          // `foo`). In the case of the latter, we will want to emit
          // suggestions from the members of `foo`.
          const lastCharIsDot =
            lines[cursorLoc.line-1][cursorLoc.column-2] === ".";

          let node: ast.Node | null = null;
          if (service.isParsedDocument(compiled)) {
            // Success case. The document parses, and we can offer
            // suggestions from a well-formed document.

            return this.completionsFromParse(
              fileUri, compiled, cursorLoc, lastCharIsDot);
          } else {
            const lastParse = this.compilerService.getLastSuccess(fileUri);
            if (lastParse == null) {
              return [];
            }

            return this.completionsFromFailedParse(
              fileUri, compiled, lastParse, cursorLoc, lastCharIsDot);
          }
        } catch (err) {
          console.log(err);
          return [];
        }
      });
  }

  // --------------------------------------------------------------------------
  // Completion methods.
  // --------------------------------------------------------------------------

  // completionsFromParse takes a `ParsedDocument` (i.e., a
  // successfully-parsed document), a cursor location, and an
  // indication of whether the user is "dotting in" to a property, and
  // produces a list of autocomplete suggestions.
  public completionsFromParse = (
    fileUri: editor.FileUri, compiled: service.ParsedDocument,
    cursorLoc: lexical.Location,
    lastCharIsDot: boolean,
  ): editor.CompletionInfo[] => {
    // IMPLEMENTATION NOTES: We have kept this method relatively free
    // of calls to `this` so that we don't have to mock out more of
    // the analyzer to test it.

    let foundNode = getNodeAtPositionFromAst(
      compiled.parse, cursorLoc);
    if (ast.isAnalyzableFindFailure(foundNode)) {
      if (foundNode.kind === "NotIdentifier") {
        return [];
      }
      if (foundNode.terminalNodeOnCursorLine != null) {
        foundNode = foundNode.terminalNodeOnCursorLine;
      } else {
        foundNode = foundNode.tightestEnclosingNode;
      }
    } else if (ast.isUnanalyzableFindFailure(foundNode)) {
      return [];
    }

    return this.completionsFromNode(
      fileUri, foundNode, cursorLoc, lastCharIsDot);
  }

  // completionsFromFailedParse takes a `FailedParsedDocument` (i.e.,
  // a document that does not parse), a `ParsedDocument` (i.e., a
  // last-known good parse for the document), a cursor location, and
  // an indication of whether the user is "dotting in" to a property,
  // and produces a list of autocomplete suggestions.
  public completionsFromFailedParse = (
    fileUri: editor.FileUri, compiled: service.FailedParsedDocument,
    lastParse: service.ParsedDocument,
    cursorLoc: lexical.Location, lastCharIsDot: boolean,
  ): editor.CompletionInfo[] => {
    // IMPLEMENTATION NOTES: We have kept this method relatively free
    // of calls to `this` so that we don't have to mock out more of
    // the analyzer to test it.
    //
    // Failure case. The document does not parse, so we need
    // to:
    //
    // 1. Obtain a partial parse from the parser.
    // 2. Get our "best guess" for where in the AST the user's
    //    cursor would be, if the document did parse.
    // 3. Use the partial parse and the environment "best
    //    guess" to create suggestions based on the context
    //    of where the user is typing.

    if (
      service.isLexFailure(compiled.parse) ||
      compiled.parse.parseError.rest == null
    ) {
      return [];
    }

    // Step 1, get the "rest" of the parse, i.e., the partial
    // parse emitted by the parser.
    const rest = compiled.parse.parseError.rest;
    const restEnd = rest.loc.end;

    if (rest == null) {
      throw new Error(`INTERNAL ERROR: rest should never be null`);
    } else if (
      !cursorLoc.inRange(rest.loc) &&
      !(restEnd.line === cursorLoc.line && cursorLoc.column === restEnd.column + 1)
    ) {
      // NOTE: the `+ 1` correctly captures the case of the
      // user typing `.`.

      // Explicitly handle the case that the user has pressed a
      // newline and `.` character. For example, in the third line
      // below:
      //
      //   metadata.withAnnotations({foo: "bar"})
      //
      //     .;
      //
      // Return no suggestions if the parse is not broken at the
      // cursor.
      const lines = compiled.text.split(/\r\n|\r|\n/g);
      const gapLines = lines.slice(restEnd.line, cursorLoc.line);
      if (gapLines.length == 0) {
        return [];
      } else if (gapLines.length === 1) {
        const gap = gapLines[0].slice(cursorLoc.column, restEnd.column);
        if (gap.trim().length != 0) {
          return [];
        }
      } else {
        const firstGap = gapLines[0].slice(restEnd.column);
        const lastGap = gapLines[gapLines.length - 1]
          .slice(
            0,
            cursorLoc.column - (lastCharIsDot ? 2 : 1));
        const middleGapLengths = gapLines
          .slice(1, gapLines.length - 2)
          .reduce((gapLenAcc: number, line: string) => gapLenAcc + line.trim().length, 0);

        if (firstGap.trim().length !== 0 || middleGapLengths !== 0 || lastGap.trim().length !== 0) {
          return [];
        }
      }

      cursorLoc = restEnd;
    }

    // Step 2, try to find the "best guess".
    let foundNode = getNodeAtPositionFromAst(
      lastParse.parse, cursorLoc);
    if (ast.isAnalyzableFindFailure(foundNode)) {
      if (foundNode.terminalNodeOnCursorLine != null) {
        foundNode = foundNode.terminalNodeOnCursorLine;
      } else {
        foundNode = foundNode.tightestEnclosingNode;
      }
    } else if (ast.isUnanalyzableFindFailure(foundNode)) {
      return [];
    }

    // Step 3, combine the partial parse and the environment
    // of the "best guess" to attempt to create meaningful
    // suggestions for the user.
    if (foundNode.env == null) {
      throw new Error("INTERNAL ERROR: Node environment can't be null");
    }
    new ast
      .InitializingVisitor(rest, foundNode, foundNode.env)
      .visit();

    // Create suggestions.
    return this.completionsFromNode(fileUri, rest, cursorLoc, lastCharIsDot);
  }

  // completionsFromNode takes a `Node`, a cursor location, and an
  // indication of whether the user is "dotting in" to a property, and
  // produces a list of autocomplete suggestions.
  private completionsFromNode = (
    fileUri: editor.FileUri, node: ast.Node, cursorLoc: lexical.Location,
    lastCharIsDot: boolean,
  ): editor.CompletionInfo[] => {
    // Attempt to resolve the node.
    const ctx = new ast.ResolutionContext(
      this.compilerService, this.documents, fileUri);
    const resolved = ast.tryResolveIndirections(node, ctx);

    if (ast.isUnresolved(resolved)) {
      // If we could not even partially resolve a node (as we do,
      // e.g., when an index target resolves, but the ID doesn't),
      // then create suggestions from the environment.
      return node.env != null
        ? envToSuggestions(node.env)
        : [];
    } else if (ast.isUnresolvedIndexTarget(resolved)) {
      // One of the targets in some index expression failed to
      // resolve, so we have no suggestions. For example, in
      // `foo.bar.baz.bat`, if any of `foo`, `bar`, or `baz` fail,
      // then we have nothing to suggest as the user is typing `bat`.
      return [];
    } else if (ast.isUnresolvedIndexId(resolved)) {
      // We have successfully resolved index target, but not the index
      // ID, so generate suggestions from the resolved target. For
      // example, if the user types `foo.b`, then we would generate
      // suggestions from the members of `foo`.
      return this.completionsFromFields(resolved.resolvedTarget);
    } else if (
      ast.isResolvedFunction(resolved) ||
      ast.isResolvedFreeVar(resolved) ||
      (!lastCharIsDot && ast.isIndexedObjectFields(resolved.value) ||
      ast.isNode(resolved.value))
    ) {
      // Our most complex case. One of two things is true:
      //
      // 1. Resolved the ID to a function or a free param, in which
      //    case we do not want to emit any suggestions, or
      // 2. The user has NOT typed a dot, AND the resolve node is not
      //    fields addressable, OR it's a node. In other words, the
      //    user has typed something like `foo` (and specifically not
      //    `foo.`, which is covered in another case), and `foo`
      //    completely resolves, either to a value (e.g., a number
      //    like 3) or a set of fields (i.e., `foo` is an object). In
      //    both cases the user has type variable, and we don't want
      //    to suggest anything; if they wanted to see the members of
      //    `foo`, they should type `foo.`.
      return [];
    } else if (lastCharIsDot && ast.isIndexedObjectFields(resolved.value)) {
      // User has typed a dot, and the resolved symbol is
      // fields-resolvable, so we can return the fields of the
      // expression. For example, if the user types `foo.`, then we
      // can suggest the members of `foo`.
      return this.completionsFromFields(resolved.value);
    }

    // Catch-all case. Suggest nothing.
    return [];
  }

  private completionsFromFields = (
    fieldSet: ast.IndexedObjectFields
  ): editor.CompletionInfo[] => {
    // Attempt to get all the possible fields we could suggest. If the
    // resolved item is an `ObjectNode`, just use its fields; if it's
    // a mixin of two objects, merge them and use the merged fields
    // instead.

    return im.List(fieldSet.values())
      .filter((field: ast.ObjectField) =>
        field != null && field.id != null && field.expr2 != null && field.kind !== "ObjectLocal")
      .map((field: ast.ObjectField) => {
        if (field == null || field.id == null || field.expr2 == null) {
          throw new Error(
            `INTERNAL ERROR: Filtered out null fields, but found field null`);
        }

        let kind: editor.CompletionType = "Field";
        if (field.methodSugar) {
          kind = "Method";
        }

        const comments = this.getComments(field);
        return {
          label: field.id.name,
          kind: kind,
          documentation: comments || undefined,
        };
      })
      .toArray();
  }

  // --------------------------------------------------------------------------
  // Completion methods.
  // --------------------------------------------------------------------------

  private renderOnhoverMessage = (
    fileUri: editor.FileUri, node: ast.Node | ast.IndexedObjectFields,
  ): editor.LanguageString[] => {
    if (ast.isIndexedObjectFields(node)) {
      if (node.count() === 0) {
        return [];
      }

      const first = node.first();
      if (first.parent == null) {
        return [];
      }
      node = first.parent;
    }

    const commentText: string | null = this.resolveComments(node);

    const doc = this.documents.get(fileUri);
    let line: string = doc.text.split(os.EOL)
      .slice(node.loc.begin.line - 1, node.loc.end.line)
      .join("\n");

    if (ast.isFunctionParam(node)) {
      // A function parameter is either a free variable, or a free
      // variable with a default value. Either way, there's not more
      // we can know statically, so emit that.
      line = node.prettyPrint();
    }

    line = node.prettyPrint();

    return <editor.LanguageString[]>[
      {language: 'jsonnet', value: line},
      commentText,
    ];
  }

  // --------------------------------------------------------------------------
  // Comment resolution.
  // --------------------------------------------------------------------------

  // resolveComments takes a node as argument, and attempts to find the
  // comments that correspond to that node. For example, if the node
  // passed in exists inside an object field, we will explore the parent
  // nodes until we find the object field, and return the comments
  // associated with that (if any).
  public resolveComments = (node: ast.Node | null): string | null => {
    while (true) {
      if (node == null) { return null; }

      switch (node.type) {
        case "ObjectFieldNode": {
          // Only retrieve comments for.
          const field = <ast.ObjectField>node;
          if (field.kind != "ObjectFieldID" && field.kind == "ObjectFieldStr") {
            return null;
          }

          // Convert to field object, pull comments out.
          return this.getComments(field);
        }
        default: {
          node = node.parent;
          continue;
        }
      }
    }
  }

  private getComments = (field: ast.ObjectField): string | null => {
    // Convert to field object, pull comments out.
    const comments = field.headingComments;
    if (comments == null) {
      return null;
    }

    return comments.text.join(os.EOL);
  }
}

//
// Utilities.
//

export const getNodeAtPositionFromAst = (
  rootNode: ast.Node, pos: lexical.Location
): ast.Node | ast.FindFailure => {
  // Special case. Make sure that if the cursor is beyond the range
  // of text of the last good parse, we just return the last node.
  // For example, if the user types a `.` character at the end of
  // the document, the document now fails to parse, and the cursor
  // is beyond the range of text of the last good parse.
  const endLoc = rootNode.loc.end;
  if (endLoc.line < pos.line || (endLoc.line == pos.line && endLoc.column < pos.column)) {
    pos = endLoc;
  }

  const visitor = new ast.CursorVisitor(pos, rootNode);
  visitor.visit();
  const tightestNode = visitor.nodeAtPosition;
  return tightestNode;
}

const envToSuggestions = (env: ast.Environment): editor.CompletionInfo[] => {
    return env.map((value: ast.LocalBind | ast.FunctionParam, key: string) => {
      // TODO: Fill in documentation later. This might involve trying
      // to parse function comment to get comments about different
      // parameters.
      return <editor.CompletionInfo>{
        label: key,
        kind: "Variable",
      };
    })
    .toArray();
}
