---
name: "main.jsonnet"
route: "/jsonnet/main"
menu: "Writing Jsonnet"
---

# `main.jsonnet`

The most important file is called `main.jsonnet`, because this is where Tanka
invokes the Jsonnet compiler on. Every single line of Jsonnet, including
imports, functions and whatnot is then evaluated until a single, very big JSON
object is left.  
This object is returned to Tanka and includes all of your Kubernetes manifests
somewhere in it, most probably deeply nested.

But as `kubectl` expects a yaml stream, and not a nested tree, Tanka needs to
extract your objects first. To do this, it traverses the tree until it finds
something that looks like a Kubernetes manifest. An object is considered valid
when it has both, `kind` and `apiVersion` set.

> This behaviour is going to change in the future, `metadata.name` will also
> become required.

To ensure Tanka can find your manifests, the output of your Jsonnet needs to
have one of the following structures:

## Deeply nested object (Recommended)

Most commonly used is a single big object that includes all manifests as
leaf-nodes.

How deeply encapsulated the actual object is does not matter, Tanka will
traverse down until it finds something that is valid.

```json
{
  "prometheus": {
    "service": {
      // Service nested one level
      "apiVersion": "v1",
      "kind": "Service",
      "metadata": {
        "name": "promSvc"
      }
    },
    "deployment": {
      "apiVersion": "apps/v1", // apiVersion ..
      "kind": "Deployment", // .. and kind are required to identify an object.
      "metadata": {
        "name": "prom"
      }
    }
  },
  "web": {
    "nginx": {
      "deployment": {
        // Deployment nested two levels
        "apiVersion": "apps/v1",
        "kind": "Deployment",
        "metadata": {
          "name": "nginx"
        }
      }
    }
  }
}
```

Using this technique has the big benefit that it is self-documentary, as the
nesting of keys can be used to logically group related manifests, for example by
application.

An encapsulation level of zero is also possible, which means nothing else than
regular object like it could be obtained from `kubectl show -o json`:

```json
{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "name": "foo"
  }
}
```

## Array

Using an array of objects is also fine:

```json
[
  {
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
      "name": "promSvc"
    }
  },
  {
    "apiVersion": "apps/v1",
    "kind": "Deployment",
    "metadata": {
      "name": "prom"
    }
  }
]
```

### `List` type

Users of `kubectl` might have had contact with a type called `List`. It is not
part of the official Kubernetes API but rather a pseudo-type introduced by
`kubectl` for dealing with multiple objects at once. Thus, Tanka does not
support it out of the box.

To take full advantage of Tankas features, you can manually flatten it:

```jsonnet
local list = import "list.libsonnet";

# expose the `items` array on the top level:
list.items
```
