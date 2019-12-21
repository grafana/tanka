---
name: Kubernetes library
route: /tutorial/k-lib
menu: Tutorial
---

# Kubernetes library

The last section has shown that using a library for creating Kubernetes objects
can drastically simplify the code you need to write. However, there is a huge
amount of different kinds of objects and the Kubernetes API is evolving (and
thus changing) quite rapidly.

Writing and maintaining such a library could be a full-time job on it's own.
Luckily, it is possible to generate such a library from the Kubernetes OpenAPI
specification! Even better, it has already been done for you.

## k.libsonnet

The library is called `k.libsonnet` (sometimes also `ksonnet-lib`), currently
available at https://github.com/ksonnet/ksonnet-lib.

> **Note**: Being part of the discontinued `ksonnet` project, the library is not
> really maintained at the moment. However, Grafana Labs will soon pick this up and
> take care of it :D  
> Nevertheless, it has already proven to be stable enough for our own production
> setup to rely on it.

However while using it internally we have discovered that the exposed API has
several annoyances. To address them, we developed another library that builds on
top of the generated one but improves the developer experience:
https://github.com/grafana/jsonnet-libs/ksonnet-util

If you do not have any strong reasons against it, just adopt the wrapper as
well, it will ease your work. Ultimately, we hope to integrate our enhancements
in the original library as well.

## Installation

Like every other external library, `ksonnet-lib` can be installed using `jsonnet-bundler`. However, we need to pick a version first:

| Version          | OpenAPI version | Notes                                                                                                      |
| ---------------- | --------------- | ---------------------------------------------------------------------------------------------------------- |
| `ksonnet.beta.3` | `v1.8.0`        |                                                                                                            |
| `ksonnet.beta.4` | `v1.14.0`       | Required for 1.16+: includes `apps/v1`, which must be used for `Deployment`, etc. from this version and up |

For the time being, you will most probably want to go with `ksonnet.beta.4`, as
it should cover all current Kubernetes versions around (mostly).

> **Note**: Once our own edition of this library are available, there will be a
> pregenerated one for each Kubernetes version.

Let's install it then:

```bash
$ jb install github.com/ksonnet/ksonnet-lib/ksonnet.beta.4
$ jb install github.com/grafana/jsonnet-libs/ksonnet-util
```

This creates the following files in `/vendor`:

```bash
vendor
├── ksonnet.beta.4
│   ├── k.libsonnet # human friendly wrapper (this is what we use in our code)
│   └── k8s.libsonnet # literally the entire API as a library. Very huge file
└── ksonnet-util
    └── kausal.libsonnet # Grafana's wrapper
```

> **Info**: The `vendor/` is the location for external libraries, while `lib/`
> can be used for your own ones. Check [import paths](/libraries/import-paths) for more information.

## Aliasing

While you could already use the library by importing `ksonnet.beta.4/k.libsonnet`, this has a drawback: Because the Kubernetes API version is indirectly included in the import name, it makes it impossible to create version agnostic downstream libraries.

As a workaround, most libraries expect the correct version of `k.libsonnet` to be importable as a literal `k.libsonnet` (without any package name prefixes). While Jsonnet-bundler won't let you do that, you can alias it by hand:

First, create a file `/lib/k.libsonnet` and add the following line to it:

```jsonnet
import "ksonnet.beta.4/k.libsonnet"
```

> **More information**:
>
> - This works, because `import` behaves like copy-pasting. So
>   the contents of `ksonnet.beta.4` are "copied" into our new file, making them
>   behave exactly the same.
> - Make sure to use the `lib/` instead of the `vendor/` folder, because `jb`
>   cleans everything from `vendor/` it did not create itself on each run.

## Using it

First we need to import it in `main.jsonnet`:

```diff
- (import "kubernetes.libsonnet") +
+ (import "ksonnet-util/kausal.libsonnet") +
  (import "grafana.jsonnet") +
  (import "prometheus.jsonnet") +
  { /* ... */ }
```

> **Note**: `kausal.libsonnet` imports literal `k.libsonnet`, so
> [aliasing](#aliasing) is a must here. This works, because `/lib` and `/vendor`
> are automatically searched for libraries, and `k.libsonnet` can be found in
> `/lib` due to aforementioned aliasing.

Now that we have installed the correct version, let's use it in
`/environments/default/grafana.jsonnet` instead of our own helper:

```jsonnet
{
  // use locals to extract the parts we need
  local deploy = $.apps.v1.deployment,
  local container = $.core.v1.container,
  local port = $.core.v1.containerPort,
  local service = $.core.v1.service,
  // defining the objects:
  grafana: {
    // deployment constructor: name, replicas, containers
    deployment: deploy.new(name=$._config.grafana.name, replicas=1, containers=[
      // container constructor
      container.new($._config.grafana.name, "grafana/grafana")
      + container.withPorts( // add ports to the container
          [port.new("ui", $._config.grafana.port)] // port constructor
        ),
    ]),

    // instead of using a service constructor, our wrapper provides
    // a handy helper to automatically generate a service for a Deployment
    service: $.util.serviceFor(self.deployment)
             + service.mixin.spec.withType("NodePort"),
  }
}
```

## Full example

Now that creating the individual objects does not take more than 5 lines, we can
merge it all back into a single file (`main.jsonnet`) and take a look at the
whole picture:

```jsonnet
(import "ksonnet-util/kausal.libsonnet") +
{
  _config:: {
    grafana: {
      port: 3000,
      name: "grafana",
    },
    prometheus: {
      port: 9090,
      name: "prometheus"
    }
  },

  local deployment = $.apps.v1.deployment,
  local container = $.core.v1.container,
  local port = $.core.v1.containerPort,
  local service = $.core.v1.service,

  prometheus: {
    deployment: deployment.new(
      name=$._config.prometheus.name, replicas=1,
      containers=[
        container.new($._config.prometheus.name, "prom/prometheus")
        + container.withPorts([port.new("api", $._config.prometheus.port)]),
      ],
    ),
    service: $.util.serviceFor(self.deployment),
  },
  grafana: {
    deployment: deployment.new(
      name=$._config.grafana.name, replicas=1,
      containers=[
        container.new($._config.grafana.name, "grafana/grafana")
        + container.withPorts([port.new("ui", $._config.grafana.port)]),
      ],
    ),
    service: $.util.serviceFor(self.deployment) + service.mixin.spec.withType("NodePort"),
  },
}
```

That's a pretty big improvement, considering how verbose and error-prone it was
before!
