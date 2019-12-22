import * as ast from "../ast";
import * as location from "./location";

export const staticErrorPrefix = "STATIC ERROR: ";
export const runtimeErrorPrefix = "RUNTIME ERROR: ";

//////////////////////////////////////////////////////////////////////////////
// StaticError

// StaticError represents an error during parsing/lexing some jsonnet.
export class StaticError {
  constructor (
    // rest allows the parser to return a partial parse result. For
    // example, if the user types a `.`, it is likely the document
    // will not parse, and it is useful to the autocomplete mechanisms
    // to return the AST that preceeds the `.` character.
    readonly rest: ast.Node | null,
    readonly loc: location.LocationRange,
    readonly msg: string,
  ) {}

  public Error = (): string => {
    const loc = this.loc.IsSet()
      ? this.loc.String()
      : "";
    return `${loc} ${this.msg}`;
  }
}

export const isStaticError = (x: any): x is StaticError => {
    return x instanceof StaticError;
}

export const MakeStaticErrorMsg = (msg: string): StaticError => {
  return new StaticError(null, location.MakeLocationRangeMessage(""), msg);
}

export const MakeStaticErrorPoint = (
  msg: string, fn: string, l: location.Location
): StaticError => {
  return new StaticError(null, location.MakeLocationRange(fn, l, l), msg);
}

export const MakeStaticError = (
  msg: string, lr: location.LocationRange
): StaticError => {
  return new StaticError(null, lr, msg);
}

export const MakeStaticErrorRest = (
  rest: ast.Node, msg: string, lr: location.LocationRange
): StaticError => {
  return new StaticError(rest, lr, msg);
}
