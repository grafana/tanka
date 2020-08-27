---
name: Namespaces
route: /namespaces
---

# Namespaces

When using Tanka, namespaces are handled slightly different compared to
`kubectl`, because environments offer more granular control than contexts used
by `kubectl`.

## Default namespaces

In the [`spec.json`](/config/#file-format) of each environment, you can set the
`spec.namespace` field, which is the default namespace. The default namespace is
set for every resource that **does not** have a namespace **set from Jsonnet**.

|     | Scenario                                                                           | Action                                                                          |
| --- | ---------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| 1.  | Your resource **lacks namespace** information (`metadata.namespace`) unset or `""` | Tanka sets `metadata.namespace` to the value of `spec.namespace` in `spec.json` |
| 2.  | Your resource **already has** namespace information                                | Tanka does nothing, accepting the explicit namespace                            |

While we recommend keeping environments limited to a single namespace, there are
legit cases where it's handy to have them span multiple namespaces, for example:

- Some other piece of software (Operators, etc) require resources to be in a specific namespace
- A rarely changing "base" environment holding resources deployed for many clusters in the same way
- etc.

## Cluster-wide resources

Some resources in Kubernetes are cluster-wide, meaning they don't belong to a single namespace at all.

Tanka will make an attempt to not add namespaces to *known* cluster-wide types. 
It does this with a short list of types in [the source code](https://github.com/grafana/tanka/blob/master/pkg/process/namespace.go).

Tanka cannot feasibly maintain this list for all known custom resource types. In those cases, resources will have namespaces added to their manifests,
and kubectl should happily apply them as non-namespaced resources.

If this presents a problem for your workflow, you can **override this** behavior
per-resource, by setting the `tanka.dev/namespaced` annotation to `"false"`
(must be of `string` type):

```jsonnet
thing: clusterRole.new("myClusterRole")
       + clusterRole.mixin.metadata.withAnnotationsMixin({ "tanka.dev/namespaced": "false" })
```
