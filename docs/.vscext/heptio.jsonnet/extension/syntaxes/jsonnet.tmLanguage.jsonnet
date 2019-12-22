local identifier = "[a-zA-Z_][a-z0-9A-Z_]*";

local Include(id) = { include: "#%s" % id };

local string = {
  local escapeCharsPattern = "\\\\([%s\\\\/bfnrt]|(u[0-9a-fA-F]{4}))",
  local illegalCharsPattern = "\\\\[^%s\\\\/bfnrtu]",

  escape:: {
    single: escapeCharsPattern % "'",
    double: escapeCharsPattern % "\"",
  },

  illegal:: {
    single: illegalCharsPattern % "'",
    double: illegalCharsPattern % "\"",
  },
};

local match = {
  Simple(name, match):: {
    name: name,
    match: match,
  },

  Span(name, begin, end):: {
    name: name,
    begin: begin,
    end: end,
  },
};

{
  "$schema": "https://raw.githubusercontent.com/martinring/tmlanguage/master/tmlanguage.json",
  name: "Jsonnet",
  patterns: [
    Include("expression"),
    Include("keywords"),
  ],
  repository: {
    expression: {
      patterns: [
        Include("literals"),
        Include("comment"),
        Include("single-quoted-strings"),
        Include("double-quoted-strings"),
        Include("triple-quoted-strings"),
        Include("builtin-functions"),
        Include("functions"),
      ]
    },
    keywords: {
      patterns: [
        match.Simple("keyword.operator.jsonnet", "[!:~\\+\\-&\\|\\^=<>\\*\\/%]"),
        match.Simple("keyword.other.jsonnet", "\\$"),
        match.Simple("keyword.other.jsonnet", "\\b(self|super|import|importstr|local|tailstrict)\\b"),
        match.Simple("keyword.control.jsonnet", "\\b(if|then|else|for|in|error|assert)\\b"),
        match.Simple("storage.type.jsonnet", "\\b(function)\\b"),
        match.Simple("variable.parameter.jsonnet", "%s\\s*(:::|\\+:::)" % identifier),
        match.Simple("entity.name.type", "%s\\s*(::|\\+::)" % identifier,),
        match.Simple("variable.parameter.jsonnet", "%s\\s*(:|\\+:)" % identifier),

      ]
    },
    literals: {
      patterns: [
         match.Simple("constant.language.jsonnet", "\\b(true|false|null)\\b"),
         match.Simple("constant.numeric.jsonnet", "\\b(\\d+([Ee][+-]?\\d+)?)\\b"),
         match.Simple("constant.numeric.jsonnet", "\\b\\d+[.]\\d*([Ee][+-]?\\d+)?\\b"),
         match.Simple("constant.numeric.jsonnet", "\\b[.]\\d+([Ee][+-]?\\d+)?\\b"),
      ]
    },
    "builtin-functions": {
      patterns: [
        match.Simple("support.function.jsonnet", "\\bstd[.](acos|asin|atan|ceil|char|codepoint|cos|exp|exponent)\\b"),
        match.Simple("support.function.jsonnet", "\\bstd[.](filter|floor|force|length|log|makeArray|mantissa)\\b"),
        match.Simple("support.function.jsonnet", "\\bstd[.](objectFields|objectHas|pow|sin|sqrt|tan|type|thisFile)\\b"),
        match.Simple("support.function.jsonnet", "\\bstd[.](acos|asin|atan|ceil|char|codepoint|cos|exp|exponent)\\b"),
        match.Simple("support.function.jsonnet", "\\bstd[.](abs|assertEqual|escapeString(Bash|Dollars|Json|Python))\\b"),
        match.Simple("support.function.jsonnet", "\\bstd[.](filterMap|flattenArrays|foldl|foldr|format|join)\\b"),
        match.Simple("support.function.jsonnet", "\\bstd[.](lines|manifest(Ini|Python(Vars)?)|map|max|min|mod)\\b"),
        match.Simple("support.function.jsonnet", "\\bstd[.](set|set(Diff|Inter|Member|Union)|sort)\\b"),
        match.Simple("support.function.jsonnet", "\\bstd[.](range|split|stringChars|substr|toString|uniq)\\b"),
      ]
    },
    "single-quoted-strings":
      match.Span("string.quoted.double.jsonnet", "'", "'") {
        patterns: [
          match.Simple("constant.character.escape.jsonnet", string.escape.single),
          match.Simple("invalid.illegal.jsonnet", string.illegal.single),
        ]
      },
    "double-quoted-strings":
      match.Span("string.quoted.double.jsonnet", "\"", "\"") {
        patterns: [
          match.Simple("constant.character.escape.jsonnet", string.escape.double),
          match.Simple("invalid.illegal.jsonnet", string.illegal.double),
        ]
      },
    "triple-quoted-strings": {
      patterns: [
        match.Span("string.quoted.triple.jsonnet", "\\|\\|\\|", "\\|\\|\\|"),
      ]
    },
    functions: {
      patterns: [
        match.Span("meta.function", "\\b([a-zA-Z_][a-z0-9A-Z_]*)\\s*\\(", "\\)") {
          beginCaptures: {
            "1": { name: "entity.name.function.jsonnet" }
          },
          patterns: [
            Include("expression"),
          ],
        },
      ]
    },
    comment: {
      patterns: [
        match.Span("comment.block.jsonnet", "/\\*", "\\*/"),
        match.Simple("comment.line.jsonnet", "//.*$"),
        match.Simple("comment.block.jsonnet", "#.*$"),
      ]
    }
  },
  scopeName: "source.jsonnet"
}