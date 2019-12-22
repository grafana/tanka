import * as lexer from '../lexical-analysis/lexer';
import * as lexical from '../lexical-analysis/lexical';

// UiEventListener  listens to events emitted by the UI that a user
// interacts with, and responds to those events with (e.g.) completion
// suggestions, or information about a symbol. For example, when a
// user hovers over a symbol, a `Hover` event is fired, and the
// UiEventListener will dispatch that event to the hook registered
// with `onHover`.
export interface UiEventListener {
  onHover: (fileUri: string, cursorLoc: lexical.Location) => Promise<HoverInfo>
  onComplete: (
    fileUri: string, cursorLoc: lexical.Location
  ) => Promise<CompletionInfo[]>
}

// HoverInfo represents data we want to emit when a `Hover` event is
// fired. For example, this might contain the syntax-highlighted
// definition of that symbol, and whatever comments accompany it.
export interface HoverInfo {
  contents: LanguageString | LanguageString[]
}

// LanguageString represents a string that is meant to be rendered
// (e.g., colored) using syntax highlighting for `language`. This is
// mainly used in response to a `Hover` event, when we want to show
// (e.g.) the definition of a symbol, along with some documentation
// about it.
export interface LanguageString {
  language: string
  value: string
}

// CompletionType represents all the possible autocomplete
// suggestions. For example, when a user `.`'s into an object, we
// might suggest a `Field` that completes it.
export type CompletionType = "Field" | "Variable" | "Method";

// CompletionInfo represents an auto-complete suggestion. Typically
// this consists of a `label` (i.e., the suggested completion text), a
// `kind` (i.e., the type of suggestion, like a file or a field), and
// documentation, if any.
export interface CompletionInfo {
  label: string
  kind: CompletionType
  documentation?: string
}
