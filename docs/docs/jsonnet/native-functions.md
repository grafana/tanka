---
name: "Native Functions"
route: "/jsonnet/native"
menu: "Writing Jsonnet"
---

# Native Functions

Tanka extends Jsonnet using _native functions_, offering additional functionality not yet available in the standard library.

To use them in your code, you need to access them using `std.native` from the standard library:

```jsonnet
{
  deployment:  std.native('<name>')('<arguments>'),
}
```

`std.native` takes the native function's name as a `string` argument and returns a `function`, which is called using the second set of parentheses.

## parseJson

### Signature

```ts
parseJSON(string json) Object
```

`parseJSON` wraps `json.Unmarshal` to convert a json string into a dict.

### Examples

```jsonnet
{
  parseJsonArray: std.native('parseJson')('[0, 1, 2]'),
  parseJsonObject: std.native('parseJson')('{ "foo": "bar" }'),
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
parseYAML(string yaml) []Object
```

`parseYAML` wraps `yaml.Unmarshal` to convert a string of yaml document(s) into
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
manifestJSONFromJSON(string json, int indent) string
```

`manifestJSONFromJSON` reserializes JSON and allows to change the indentation.

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
manifestYAMLFromJSON(string json) string
```

`manifestYamlFromJSON` serializes a JSON string as a YAML document.

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
[RE2](https://golang.org/s/re2syntax) regular expression.

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
