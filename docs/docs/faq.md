---
name: "Frequently asked questions"
route: "/faq/"
---

# Frequently asked questions

## What is Jsonnet?

Jsonnet is a data templating language, originally created by Google.

It is a superset of JSON, which adds common structures from full programming
languages to data modeling. Because it being a superset of JSON and ultimately
always compiling to JSON, it is guaranteed that the output will be valid JSON
(or YAML).

By allowing _functions_ and _imports_, rich abstraction is possible, even across
project boundaries.

For more, refer to the official documentation: https://jsonnet.org/

## How is this different from ksonnet?

Tanka aims to be a fully compatible, drop-in replacement for the main workflow
of `ksonnet` (`show`, `diff`, `apply`).

In general, both tools are very similar when it comes to how they handle Jsonnet
and apply to a Kubernetes cluster.

However, `ksonnet` included a rich code generator for establishing a CLI based
workflow for editing Kubernetes objects. It also used to manage dependencies
itself and had a lot of concepts for different levels of abstractions. When
designing Tanka, we felt these add more complexity for the user than they
provide additional value. To keep Tanka as minimal as possible, these are **not
available** and are not likely to be ever added.

## What about kubecfg ?

Tanka development has started at the time when kubecfg was a part of
already-deprecated `ksonnet` project. Although these projects are similar, Tanka
aims to provide continuity for `ksonnet` users, whereas `kubecfg` is (according
to the project's [README.md](https://github.com/bitnami/kubecfg/blob/master/README.md))
really just a thin Kubernetes-specific wrapper around jsonnet evaluation.

## Why not Helm?

Helm relies heavily on _string templating_ `.yaml` files. We feel this is the
wrong way to approach the absence of abstractions inside of `yaml`, because the
templating part of the application has no idea of the structure and syntax of
yaml.

This makes debugging very hard. Furthermore, `helm` is not able to provide an
adequate solution for edge cases. If I wanted to set some parameters that are
not already implemented by the Chart, I have no choice but to modify the Chart
first.

Jsonnet on the other hand got you covered by supporting mixing (patching,
deep-merging) objects on top of the libraries output if required.
