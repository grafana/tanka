---
name: "Injecting Values"
route: "/jsonnet/injecting-values"
menu: "Writing Jsonnet"
---

# Injecting Values

Sometimes it might be required to pass externally acquired data into Jsonnet.

There are three ways of doing so:

1. [JSON files](#json-files)
2. [External variables](#external-variables)
3. [Top level arguments](#top-level-arguments)

## JSON files

Jsonnet is a superset of JSON, it treats any JSON as valid Jsonnet. Because many
systems can be told to output their data in JSON format, this provides a pretty
good interface between those.

For example, your deployment pipeline could acquire secrets from systems such as
[Vault](https://www.vaultproject.io/), etc. and write that into `secrets.json`.

```jsonnet
local secrets = import "secrets.json";

{
  foo: secrets.myPassword,
}
```

## External variables

Another way of passing values from the outside are external variables, which are specified like so:

```bash
# strings
$ tk show . --ext-str hello=world

# any Jsonnet snippet
$ tk show . --ext-str foo=4 --ext-str bar='[ 1, 3 ]'
```

They can be accessed using `std.extVar` and the name given to them on the command line:

```jsonnet
{
  foo: std.extVar('foo'), // 4, integer
  bar: std.extVar('bar'), // [ 1, 3 ], array
}
```

> **Warning**: It's not possible to detect what extVars are available, and it's
> not possible to set default values.  
> Try to use [**Top Level Arguments**](#top-level-arguments) instead, to avoid surprising behaviour

## Top Level Arguments

To avoid the magic and possibly surprising nature of external variables, Top
Level Arguments can be used.

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
// function returns the final object, supporting parameters and default values
function(who, msg="Hello %s!") {
  hello: msg % who
}
```

Here, `who` needs a value while `msg` has a default. This can be invoked like so:

```bash
$ tk show . --tla-str who=John
```
