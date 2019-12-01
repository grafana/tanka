---
name: "Configuration"
route: "/environments/configuration"
menu: "Environments"
---

# Configuration

An Environment describes a configuration tailored for a context of any kind
whatsoever.  
Such a context might be `dev` and `prod` environments, blue/green deployments or
the same application available in multiple regional zones / datacenters.

An environment does not need to be created, it rather just exists as soon as a
main.jsonnet is found somewhere in the tree of a `rootDir`. This also means that
an environment by definition is equivalent to a
[`baseDir`](directory-structure.md#base-directory-basedir).

## Configuration

To correctly deal with an environment, Tanka needs some additional information
about it. These are specified in a file called `spec.json` which is placed next
to `main.jsonnet`.

```json
{
  "apiVersion": "tanka.dev/v1alpha1",
  "kind": "Environment",
  "metadata": {
    "name": "auto",
    "labels": {}
  },
  "spec": {
    "apiServer": "https://localhost:6443",
    "namespace": "default"
  }
}
```

| Field                | Description                                      |
| -------------------- | ------------------------------------------------ |
| `apiVersion`         | currently only `tanka.dev/v1alpha1` is available |
| `kind`               | always `Environment`                             |
| `metadata.name`      | _automatically set to the directory name_        |
| `metadata.labels`    | descriptive `key:value` pairs                    |
| **`spec.apiServer`** | The Kubernetes endpoint to use                   |
| **`spec.namespace`** | Default namespace used if not set in jsonnet     |

Everything written in **bold** is required, the other fields may be omitted.

## Context discovery

To make sure you **never** apply to the wrong cluster, Tanka parses the output
of `kubectl config view` to select a context that matches the API Server
endpoint specified in the [Configuration](#configuration).

It first searches for a `cluster` matching the IP or hostname and then for a
context that uses this cluster. If one of those is missing, the apply fails.

So please make sure `$KUBECONFIG` and `kubectl` are set up correctly if you run
into any problems.
