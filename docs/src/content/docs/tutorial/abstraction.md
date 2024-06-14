---
title: Abstraction
sidebar:
  order: 5
---

While we won't need to touch the resource definitions directly that frequently
anymore now that our deployments definitions are parametrized, the
`main.jsonnet` file is still very long and hard to read. Especially because of
all the brackets, it's even worse than yaml at the moment.

## Splitting it up

Let's start cleaning this up by separating logical pieces into distinct files:

- `main.jsonnet`: Still our main file, importing the other files
- `grafana.libsonnet`: `Deployment` and `Service` for the Grafana instance
- `prometheus.libsonnet`: `Deployment` and `Service` for the Prometheus server

:::note
The extension for Jsonnet libraries is `.libsonnet`. While you do
not have to use it, it distinguishes helper code from actual configuration.
:::

```jsonnet
// /environments/default/grafana.libsonnet
{
  new(name, port)::{
    deployment: {
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      metadata: {
        name: name,
      },
      spec: {
        selector: {
          matchLabels: {
            name: name,
          },
        },
        template: {
          metadata: {
            labels: {
              name: name,
            },
          },
          spec: {
            containers: [
              {
                image: 'grafana/grafana',
                name: name,
                ports: [{
                    containerPort: port,
                    name: 'ui',
                }],
              },
            ],
          },
        },
      },
    },
    service: {
      apiVersion: 'v1',
      kind: 'Service',
      metadata: {
        labels: {
          name: name,
        },
        name: name,
      },
      spec: {
        ports: [{
            name: '%s-ui' % name,
            port: port,
            targetPort: port,
        }],
        selector: {
          name: name,
        },
        type: 'NodePort',
      },
    },
  }
}
```

The file should contain an object with the same function that was defined under the `grafana` in `/environments/default/main.jsonnet`, but called `new` instead of `grafana`.
Do the same for `/environments/default/prometheus.libsonnet` as well.

```jsonnet
// /environments/default/main.jsonnet
local grafana = import "grafana.libsonnet";
local prometheus = import "prometheus.libsonnet";

{
  grafana: grafana.new("grafana", 3000),
  prometheus: prometheus.new("prometheus", 9090),
}
```

## Helper utilities

While `main.jsonnet` is now short and very readable, the other two files are not
really an improvement over regular yaml, mostly because they are still full of
boilerplate.

Let's use functions to create some useful helpers to reduce the amount of
repetition. For that, we create a new file called `kubernetes.libsonnet`, which
will hold our Kubernetes utilities.

### A Deployment constructor

Creating a `Deployment` requires some mandatory information and a lot of
boilerplate. A function that creates one could look like this:

```jsonnet
// /environments/default/kubernetes.libsonnet
{
  deployment: {
    new(name, containers):: {
      apiVersion: "apps/v1",
      kind: "Deployment",
      metadata: {
        name: name,
      },
      spec: {
        selector: { matchLabels: {
          name: name,
        }},
        template: {
          metadata: { labels: {
            name: name,
          }},
          spec: { containers: containers }
        }
      }
    }
  }
}
```

Invoking this function will substitute all the variables with the respective
passed function parameters and return the assembled object.

Let's simplify our `grafana.libsonnet` a bit:

```jsonnet
local k = import "kubernetes.libsonnet";

{
  new(name, port):: {
    deployment: k.deployment.new(name, [{
      image: 'grafana/grafana',
      name: name,
      ports: [{
          containerPort: port,
          name: 'ui',
      }],
    }]),
    service: {
      apiVersion: 'v1',
      kind: 'Service',
      metadata: {
        labels: {
          name: name,
        },
        name: name,
      },
      spec: {
        ports: [{
            name: '%s-ui' % name,
            port: port,
            targetPort: port,
        }],
        selector: {
          name: name,
        },
        type: 'NodePort',
      },
    },
  }
}
```

This drastically simplified the creation of the `Deployment`, because we do not
need to remember how exactly a `Deployment` is structured anymore. Just use
our helper and you are good to go.

:::tip[Task]
Now try adding a constructor for a `Service` to `kubernetes.libsonnet`
and use both helpers to recreate the other objects as well.
:::
