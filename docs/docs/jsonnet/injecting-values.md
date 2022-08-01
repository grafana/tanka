---
name: "Injecting Values"
route: "/jsonnet/injecting-values/"
menu: "Writing Jsonnet"
---

# Injecting Values

Sometimes it might be required to pass externally acquired data into Jsonnet.

There are three ways of doing so:

- [Injecting Values](#injecting-values)
  - [JSON files](#json-files)
  - [External variables](#external-variables)
  - [Top Level Arguments](#top-level-arguments)

Also check out the [official Jsonnet docs on this
topic](https://jsonnet.org/ref/language.html#passing-data-to-jsonnet).

## JSON files

Jsonnet is a superset of JSON, it treats any JSON as valid Jsonnet. Because many
systems can be told to output their data in JSON format, this provides a pretty
good interface between those.

For example, your build tooling like `make` could acquire secrets from systems such as
[Vault](https://www.vaultproject.io/), etc. and write that into `secrets.json`.

```jsonnet
local secrets = import "secrets.json";

{
  foo: secrets.myPassword,
}
```

> **Note**: Using `import` with JSON treats it as Jsonnet, so make sure to not
> use it with untrusted code.  
> A safer, but more verbose, alternative is `std.parseJson(importstr 'path_to_json.json')`

## External variables

Another way of passing values from the outside are external variables, which are specified like so:

```bash
# strings
$ tk show . --ext-str hello=world

# any Jsonnet snippet
$ tk show . --ext-code foo=4 --ext-code bar='[ 1, 3 ]'
```

They can be accessed using `std.extVar` and the name given to them on the command line:

```jsonnet
{
  foo: std.extVar('foo'), // 4, integer
  bar: std.extVar('bar'), // [ 1, 3 ], array
}
```

> **Warning**: External variables are directly accessible in all parts of the
> configuration, which can make it difficult to track where they are used and
> what effect they have on the final result.
> Try to use [Top Level Arguments](#top-level-arguments) instead.

## Top Level Arguments

Usually with Tanka, your `main.jsonnet` holds an object at the top level (most
outer type in the generated JSON):

```jsonnet
// main.jsonnet
{
  /* your resources */
}
```

Another type of Jsonnet that naturally accepts parameters is the `function`.
When the Jsonnet compiler finds a function at the top level, it invokes it and
allows passing parameter values from the command line:

```jsonnet
// Actual output (object) returned by function, which is taking parameters and default values
function(who, msg="Hello %s!") {
  hello: msg % who
}
```

Here, `who` needs a value while `msg` has a default. This can be invoked like so:

```bash
$ tk show . --tla-str who=John
```
