---
title: Native Functions
sidebar:
  order: 3
---

Tanka extends Jsonnet using _native functions_, offering additional functionality not yet available in the standard library.

To use them in your code, you need to access them using `std.native` from the standard library:

```jsonnet
{
  someField:  std.native('<name>')(<arguments>),
}
```

`std.native` takes the native function's name as a `string` argument and returns a `function`, which is called using the second set of parentheses.

## sha256

### Signature

```ts
sha256(string str) string
```

`sha256` computes the SHA256 sum of the given string.

### Examples

```jsonnet
{
  sum: std.native('sha256')('Hello, World!'),
}
```

Evaluating with Tanka results in the JSON:

```json
{
  "sum": "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
}
```

## parseJson

### Signature

```ts
parseJson(string json) Object
```

`parseJson` parses a json string and returns the respective Jsonnet type (`Object`, `Array`, etc).

### Examples

```jsonnet
{
  array: std.native('parseJson')('[0, 1, 2]'),
  object: std.native('parseJson')('{ "foo": "bar" }'),
}
```

Evaluating with Tanka results in the JSON:

```json
{
  "array": [0, 1, 2],
  "object": {
    "foo": "bar"
  }
}
```

## parseYaml

### Signature

```ts
parseYaml(string yaml) []Object
```

`parseYaml` wraps `yaml.Unmarshal` to convert a string of yaml document(s) into
a set of dicts. If `yaml` only contains a single document, a single value array
will be returned.

### Examples

```jsonnet
{
  yaml: std.native('parseYaml')(|||
    ---
    foo: bar
    ---
    bar: baz
  |||),
}
```

Evaluating with Tanka results in the JSON:

```json
{
  "yaml": [
    {
      "foo": "bar"
    },
    {
      "bar": "baz"
    }
  ]
}
```

## manifestJsonFromJson

### Signature

```ts
manifestJsonFromJson(string json, int indent) string
```

`manifestJsonFromJson` reserializes JSON and allows to change the indentation.

### Examples

```jsonnet
{
  indentWithEightSpaces: std.native('manifestJsonFromJson')('{ "foo": { "bar": "baz" } }', 8),
}
```

Evaluating with Tanka results in the JSON:

```json
{
  "indentWithEightSpaces": "{\n        \"foo\": {\n                \"bar\": \"baz\"\n        }\n}\n"
}
```

## manifestYamlFromJson

### Signature

```ts
manifestYamlFromJson(string json) string
```

`manifestYamlFromJson` serializes a JSON string as a YAML document.

### Examples

```jsonnet
{
  yaml: std.native('manifestYamlFromJson')('{ "foo": { "bar": "baz" } }'),
}
```

Evaluating with Tanka results in the JSON:

```json
{
  "yaml": "foo:\n    bar: baz\n"
}
```

## escapeStringRegex

### Signature

```ts
escapeStringRegex(string s) string
```

`escapeStringRegex` escapes all regular expression metacharacters and returns a
regular expression that matches the literal text.

### Examples

```jsonnet
{
  escaped: std.native('escapeStringRegex')('"([0-9]+"'),
}
```

Evaluating with Tanka results in the JSON:

```json
{
  "escaped": "\"\\(\\[0-9\\]\\+\""
}
```

## regexMatch

### Signature

```ts
regexMatch(string regex, string s) boolean
```

`regexMatch` returns whether the given string is matched by the given
[RE2](https://github.com/google/re2/wiki/Syntax) regular expression.

### Examples

```jsonnet
{
  matched: std.native('regexMatch')('.', 'a'),
}
```

Evaluating with Tanka results in the JSON:

```json
{
  "matched": true
}
```

## regexSubst

### Signature

```ts
regexSubst(string regex, string src, string repl) string
```

`regexSubst` replaces all matches of the re2 regular expression with the
replacement string.

### Examples

```jsonnet
{
  substituted: std.native('regexSubst')('p[^m]*', 'pm', 'poe'),
}
```

Evaluating with Tanka results in the JSON:

```json
{
  "substituted": "poem"
}
```
