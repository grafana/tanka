---
name: Syntax overview
route: /jsonnet/overview
menu: Writing Jsonnet
---

# Language overview

[Jsonnet](https://jsonnet.org) is the data templating language Tanka uses for
expressing what shall be deployed to your Kubernetes cluster. Understanding
Jsonnet is crucial to using Tanka effectively.

This page covers the Jsonnet language itself. For more information on how to use Jsonnet with Kubernetes, see [the tutorial](/tutorial/jsonnet).

## Syntax

Being a superset of JSON, the syntax is very similar:

```jsonnet
// Line comment
/* Block comment */

// a local variable (not exported)
local greeting = "hello world!";

// the exported/returned object
{
  foo: "bar", // string
  bar: 5, // int
  baz: false, // bool
  list: [1,2,3], // array
  // object
  dict: {
    nested: greeting, // using the local
  }
  hidden:: "incognito!" // an unexported field
}
```

## Abstraction

Jsonnet has rich abstraction features, which makes it interesting for
configuring Kubernetes, as it allows to keep configurations concise, yet
readable.

- [Imports](#imports)
- [Merging](#merging)
- [Functions](#functions)

### Imports

Just as other languages, Jsonnet allows code to be imported from other files:

```jsonnet
local secret = import "./secret.libsonnet";
```

The exported object (the only non-local one) of `secret.libsonnet` is now
available as a `local` variable called `secret`.

When using Tanka, it is also possible to directly import `.json` and `.yaml`
files, as if they were a `.libsonnet`.

Make sure to take also take a look on [Libraries](libraries.md) and
[Vendoring](vendoring.md) to learn how to use `import` to re-use code.

### Merging

Deep merging allows you to change parts of an object without touching all of it.
Consider the following example:

```jsonnet{5,1-2}
local secret = {
  kind: Secret,
  metadata: {
    name: "mySecret",
    namespace: "default", // need to change that
  },
  data: {
    foo: std.base64("foo")
  }
};
```

To change the namespace only, we can use the special merge key `+::` like so:

```jsonnet
// define the patch:
local patch = {
  metadata+:: {
    namespace: "myApp"
  }
}
```

The difference between `:` and `+::` is that the former replaces the original
data at that key, while the latter applies the new object as a patch on top,
meaning that values will be updated if possible but all other stay like they
are.  
To merge those two, just add (`+`) the patch to the original:

```jsonnet
secret + patch
```

The output of this is the following JSON object:

```json
{
  "kind": "Secret",
  "metadata": {
    "name": "mySecret",
    "namespace": "myApp"
  },
  "data": {
    "foo": "Zm9vCg=="
  }
}
```

### Functions

Jsonnet supports functions, similar to how Python does. They can be defined in
two different ways:

```jsonnet
local add(x,y) = x + y;
local mul = (function(x, y) x * y);
```

Objects can have methods:

```jsonnet
{
  greet(who): "hello " + who,
}
```

Default values, keyword-args and more examples can be found at
[jsonnet.org](https://jsonnet.org/learning/tutorial.html#functions).

## Standard library

The Jsonnet standard library includes many helper methods ranging from object
and array mutation, over string utils to computation helpers.

Documentation is available at
[jsonnet.org](https://jsonnet.org/ref/stdlib.html).

## Conditionals

Jsonnet supports a conditionals in a fashion similar to a ternary operator:

```jsonnet
local tag = if prod then "v1.0" else "latest";
```

More on [jsonnet.org](https://jsonnet.org/learning/tutorial.html#conditionals).

## References

Jsonnet has multiple options to refer to parts of an object:

```jsonnet
{ // this is $
  junk: "foo",
  nested: { // this is self
    app: "Tanka",
    msg: self.app + " rocks!" // "Tanka rocks!"
  },
  children: { // this is also self
    baz: "bar",
    junk: $.junk + self.baz, // "foobar"
  }
}
```

For more information take a look at
[jsonnet.org](https://jsonnet.org/learning/tutorial.html#references)
