---
name: Using Jsonnet
menu: Tutorial
route: /tutorial/jsonnet
---

# Using Jsonnet

The most powerful piece of Tanka is the [Jsonnet data templating
language](https://jsonnet.org). Jsonnet is a superset of JSON, adding variables,
functions, patching (deep merging), arithmetic, conditionals and many more to
it.

It has a lot in common with more _real_ programming languages such as JavaScript
than with markup languages, still it is tailored specifically to representing
data and configuration. Opposing to JSON (and YAML) it is a language meant for
humans, not for computers.

## Creating a new project

To get started with Tanka and Jsonnet, let's initiate a new project:

```bash
$ mkdir prom-grafana && cd prom-grafana # create a new folder for the project and change to it
$ tk init # initiate a new project
```

This gives us the following directory structure:

```sh
├── environments
│   └── default # default environment
│       ├── main.jsonnet # main file (important!)
│       └── spec.json # environment's config
├── jsonnetfile.json
├── lib # libraries
└── vendor # external libraries
```

For the moment, we only really care about the `environments/default` folder. The
purpose of the other directories will be explained later in this guide (mostly
related to libraries).

## Environments

When using Tanka, you apply **configuration** for an **Environment** to a
Kubernetes **cluster**. An Environment is some logical group of pieces that form
an application stack.

Grafana for example runs [Loki](https://grafana.com/loki),
[Cortex](https://cortexmetrics.io) and of course
[Grafana](https://grafana.com/grafana) for our [Grafana
Cloud](https://grafana.com/cloud) hosted offering. For each of these, we have a
separate environment. Furthermore, we like to see changes to our code in
separate `dev` setups to make sure they are all good for production usage – so
we have `dev` and `prod` environments for each app as well, as `prod`
environments usually require other configuration (secrets, scale, etc) than
`dev`. This roughly leaves us with the following:

|        | Loki                                                          | Cortex                                                            | Grafana                                                             |
|--------|---------------------------------------------------------------|-------------------------------------------------------------------|---------------------------------------------------------------------|
| `prod` | Name: `/environments/loki/prod` <br /> Namespace: `loki-prod` | Name: `/environments/cortex/prod` <br /> Namespace: `cortex-prod` | Name: `/environments/grafana/prod` <br /> Namespace: `grafana-prod` |
| `dev` | Name: `/environments/loki/dev` <br /> Namespace: `loki-dev` | Name: `/environments/cortex/dev` <br /> Namespace: `cortex-dev` | Name: `/environments/grafana/dev` <br /> Namespace: `grafana-dev` |

There is no limit in Environment complexity, create as many as you need to model
your own requirements. Grafana for example also has all of these multiplied per
high-availability region.
