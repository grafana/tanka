path: tree/master
source: pkg/native/funcs.go

# Native Functions

Tanka extends the Jsonnet standard library with some common helper utilities.

### `parseJSON`
```ts
parseJSON(string json) Object
```
`parseJSON` wraps `json.Unmarshal` to convert a json string into a dict

### `parseYAML`
```ts
parseYAML(string yaml) []Object
```
`parseYAML` wraps `yaml.Unmarshal` to convert a string of yaml document(s) into a set of dicts.
If `yaml` only contains a single document, a single value array will be returned.

### `manifestJSONFromJSON`
```ts
manifestJSONFromJSON(string json, int indent) string
```
`manifestJSONFromJSON` reserializes JSON and allows to change the indentation.

### `manifestYAMLFromJSON`
```ts
manifestYAMLFromJSON(string json) string
```
`manifestYamlFromJSON` serializes a JSON string as a YAML document

### `escapeStringRegex`
```ts
escapeStringRegex(string s) string
```
`escapeStringRegex` escapes all regular expression metacharacters and returns a
regular expression that matches the literal text.

### `regexMatch`
```ts
regexMatch(string regex, string s) boolean
```
`regexMatch` returns whether the given string is matched by the given [RE2](https://golang.org/s/re2syntax) regular expression.

### `regexSubst`
```ts
regexSubst(string regex, string src, string repl) string
```
`regexSubst` replaces all matches of the re2 regular expression with the
replacement string.
