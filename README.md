# Overview

![Tanka Banner](docs/img/banner.png)

[![Build Status](https://cloud.drone.io/api/badges/grafana/tanka/status.svg)](https://cloud.drone.io/grafana/tanka)
![Golang](https://img.shields.io/badge/language-Go-blue)
![GitHub contributors](https://img.shields.io/github/contributors/grafana/tanka)
![GitHub release](https://img.shields.io/github/release/grafana/tanka)
![License](https://img.shields.io/github/license/grafana/tanka)

Tanka is a composable configuration utility for [Kubernetes](https://kubernetes.io/). It
leverages the [Jsonnet](https://jsonnet.org) language to realize flexible, reusable and
concise configuration.

- **:repeat: `ksonnet` drop-in replacement**: Tanka aims to provide the same
  workflow as `ksonnet`: `show`, `diff` and `apply` are just where you expect
  them.
- **:nut_and_bolt: integrates with the ecosystem**: Tanka doesn't re-invent the
  wheel. It rather makes heavy use of what is already there:
  [`jsonnet-bundler`](https://github.com/jsonnet-bundler/jsonnet-bundler) for
  package management and
  [`kubectl`](https://kubernetes.io/docs/reference/kubectl/overview/) for
  communicating with Kubernetes clusters.
- **:hammer: powerful:** Being a `jsonnet`-compatibility layer for Kubernetes,
  it removes the limitations of static (or template-based) configuration languages.
- **:rocket: used in production**: We use Tanka internally at
  Grafana Labs for all of our Kubernetes configuration needs.
- **:heart: fully open-source**: This is an open-source project. It is free as
  in beer and as in speech and this will never change.

## Getting started
Head over to the [Releases](https://github.com/grafana/tanka/releases) section
and download the most latest release of Tanka for your OS and arch.

Then check everything is working correctly with
```bash
$ tk --version
tk version v0.3.0
```

It is also recommended to install Jsonnet bundler:
```bash
$ go get -u github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb
```

### Creating a new project
To start from scratch with the recommended directory structure, do the following:

```bash
# create a directory and enter it
$ mkdir poetry && cd poetry

# initialize the Tanka application
$ tk init
```

### Deploying an application
As an example, [Promtail](https://github.com/grafana/loki/blob/master/docs/promtail/README.md) is being deployed using Tanka now.

After you initialized the directory structure, install the required libraries
using `jb`:
```bash
# Ksonnet kubernetes libraries
$ jb install github.com/ksonnet/ksonnet-lib/ksonnet.beta.4/k.libsonnet
$ jb install github.com/ksonnet/ksonnet-lib/ksonnet.beta.4/k8s.libsonnet

# Promtail library
$ jb install github.com/grafana/loki/production/ksonnet/promtail
```

Then, replace the contents of `environments/default/main.jsonnet` with the
following: 

```js
local promtail = import 'promtail/promtail.libsonnet';

promtail + {
  _config+:: {
    namespace: 'loki',

    promtail_config+: {
      clients: [
        {
          scheme:: 'https',
          hostname:: 'logs-us-west1.grafana.net',
          username:: 'user-id',
          password:: 'password',
          external_labels: {},
        }
      ],
      container_root_path: '/var/lib/docker',
    },
  },
}

```

As a last step, fill add the correct `spec.apiServer` and `spec.namespace` to
`environments/default/spec.json`:

```json
{
  "apiVersion": "tanka.dev/v1alpha1",
  "kind": "Environment",
  "spec": {
    "apiServer": "https://localhost:6443",
    "namespace": "default"
  }
}
```

Now use `tk show environments/default` to see the `yaml`, and
`tk apply environments/default` to apply it to the cluster.

Congratulations! You have successfully set up your first application using Tanka :tada:

## Additional resources

- https://jsonnet.org/, the official Jsonnet documentation provides lots of
  examples on how to use the language.
- https://github.com/grafana/jsonnet-libs: Grafana Labs' Jsonnet libraries are a
  rich set of configuration examples compatible with Tanka.

## License
Licensed Apache 2.0, see [LICENSE](LICENSE).
