"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const location = require("./location");
exports.staticErrorPrefix = "STATIC ERROR: ";
exports.runtimeErrorPrefix = "RUNTIME ERROR: ";
//////////////////////////////////////////////////////////////////////////////
// StaticError
// StaticError represents an error during parsing/lexing some jsonnet.
class StaticError {
    constructor(
    // rest allows the parser to return a partial parse result. For
    // example, if the user types a `.`, it is likely the document
    // will not parse, and it is useful to the autocomplete mechanisms
    // to return the AST that preceeds the `.` character.
    rest, loc, msg) {
        this.rest = rest;
        this.loc = loc;
        this.msg = msg;
        this.Error = () => {
            const loc = this.loc.IsSet()
                ? this.loc.String()
                : "";
            return `${loc} ${this.msg}`;
        };
    }
}
exports.StaticError = StaticError;
exports.isStaticError = (x) => {
    return x instanceof StaticError;
};
exports.MakeStaticErrorMsg = (msg) => {
    return new StaticError(null, location.MakeLocationRangeMessage(""), msg);
};
exports.MakeStaticErrorPoint = (msg, fn, l) => {
    return new StaticError(null, location.MakeLocationRange(fn, l, l), msg);
};
exports.MakeStaticError = (msg, lr) => {
    return new StaticError(null, lr, msg);
};
exports.MakeStaticErrorRest = (rest, msg, lr) => {
    return new StaticError(rest, lr, msg);
};
//# sourceMappingURL=static_error.js.map