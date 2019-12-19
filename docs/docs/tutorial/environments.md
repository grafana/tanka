---
name: Environments
route: /tutorial/environments
menu: Tutorial
---

# Environments

At this point, our is configuration is already flexible and concise, but not
really reusable. Let's take a look at Tanka's third buzzword as well: **Environments**.

In these days, the same piece of software is usually deployed many times inside
of one organziation. This can be `dev`, `testing` and `prod` environments, but
also regions (`europe`, `us`, `asia`) or individual customers (`foo-corp`,
`bar-gmbh`, `baz-inc`).

Most of the application however is exactly the same across those environments ..
usually only configuration, scaling or small details are different after all.
YAML (and thus `kubectl`) provides us only one solution here: Duplicating the
directory, changing the details, maintaining both. But if you have 32
environments? Correct! Then you have to maintain 32 directories of YAML. And we can all
imagine the nightmare of these files drifting apart from each other.

But again, **Jsonnet can be the solution**: By extracting the actual objects
into a library, you can import them in as many environments as you need!

## Creating a library
A library is nothing special, just a folder of `.libsonnet` files somewhere in the import paths:

| Path      | Description                                           |
|-----------|-------------------------------------------------------|
| `/lib`    | Custom, user-created libraries only for this project. |
| `/vendor` | External libraries installed using Jsonnet-bundler    |

So for our purpose `/lib` fits best, as we are only creating it for our current
project. Let's set one up:

```bash
/$ mkdir lib/prom-grafana # a folder for our prom-grafana library
/$ cd lib/prom-grafana

/lib/prom-grafana$ touch prom-grafana.libsonnet # library file that will be imported
/lib/prom-grafana$ touch config.libsonnet # _config and images
```

##### config.libsonnet
For documentation purposes it is handy to have a separate file for parameters and used images:

```jsonnet
{
  // +:: is important (we don't want to override the
  // _config object, just add to it)
  _config+:: {
    // define a namespace for this library
    promgrafana: {
      grafana: {
        port: 3000,
        name: "grafana",
      },
      prometheus: {
        port: 9090,
        name: "prometheus"
      }
    }
  },

  // again, make sure to use +::
  _images+:: {
    promgrafana: {
      grafana: "grafana/grafana",
      prometheus: "prom/prometheus",
    }
  }
}
```

##### prom-grafana.libsonnet
```jsonnet
(import "ksonnet-util/kausal.libsonnet") +
(import "./config.libsonnet") +
{
  local deployment = $.apps.v1.deployment,
  local container = $.core.v1.container,
  local port = $.core.v1.containerPort,
  local service = $.core.v1.service,

  // alias our params, too long to type every time
  local c = $._config.promgrafana,

  promgrafana: {
    prometheus: {
      deployment: deployment.new(
        name=c.prometheus.name, replicas=1,
        containers=[
          container.new(c.prometheus.name, $._images.promgrafana.prometheus)
          + container.withPorts([port.new("api", c.prometheus.port)]),
        ],
      ),
      service: $.util.serviceFor(self.deployment),
    },

    grafana: {
      deployment: deployment.new(
        name=c.grafana.name, replicas=1,
        containers=[
          container.new(c.grafana.name, $._images.promgrafana.grafana)
          + container.withPorts([port.new("ui", c.grafana.port)]),
        ],
      ),
      service: $.util.serviceFor(self.deployment) + service.mixin.spec.withType("NodePort"),
    },
  }
}
```

## Dev and Prod
So far we have only used the `environments/default` environment. Let's create some real ones:

```bash
/$ tk env add environments/prom-grafana/dev --namespace=prom-grafana-dev # one for dev ...
/$ tk env add environments/prom-grafana/prod --namespace=prom-grafana-prod # and one for prod
```

> **Note**: Remember to set up the cluster's IP in the respective `spec.json`!

All that's left now is importing the library and configuring it. For `dev`, the defaults defined in `/lib/prom-grafana/config.libsonnet` should be sufficient, so we do not override anything:

```jsonnet
// environments/prom-grafana/dev
(import "ksonnet-util/kausal.libsonnet") +
(import "prom-grafana/prom-grafana.libsonnet")
```

For `prod` however, it is a bad idea to rely on `latest` for the images .. let's
add some proper tags:

```jsonnet
// environments/prom-grafana/prod
(import "ksonnet-util/kausal.libsonnet") +
(import "prom-grafana/prom-grafana.libsonnet") + 
{
  // again, we only want to patch, not replace, thus +::
  _images+:: {
    // we update this one entirely, so we can replace this one (:)
    promgrafana: {
      prometheus: "prom/prometheus:v2.14",
      grafana: "grafana/grafana:6.5.2"
    }
  }
}
```

## Patching
The above works well for libraries we control ourselves, but what when another
team wrote the library, it was installed using `jb` from GitHub or you can't
change it easily?

Here comes the already familiar `+:` (or `+::`) syntax into play. It allows to
**partially** override values of an object. Let's say we wanted to add some labels to the Prometheus `Deployment`, but our `_config` params don't allow us to. We can still do this in our `main.jsonnet`:

```jsonnet
(import "ksonnet-util/kausal.libsonnet") +
(import "prom-grafana/prom-grafana.libsonnet") + 
{
  promgrafana+: {
    prometheus+: {
      deployment+: {
        metadata+: {
          labels+: {
            foo: "bar"
          }
        }
      }
    }
  }
}
```

By using the `+:` operator all the time and only `foo: "bar"` uses "`:`", we only
override the value of `"foo"`, while leaving the rest of the object like it was.

Let's check it worked:

```yaml
$ tk show environments/prom-grafana/patched -t deployment/prometheus
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    foo: bar # <- There it is!
  name: prometheus
  namespace: default
spec:
  minReadySeconds: 10
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      name: prometheus
  template:
    metadata:
      labels:
        name: prometheus
    spec:
      containers:
      - image: prom/prometheus
        imagePullPolicy: IfNotPresent
        name: prometheus
        ports:
        - containerPort: 9090
          name: api
```
