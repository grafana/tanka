# Environments

An Environment describes a configuration tailored for a context of any kind
whatsoever.  
Such a context might be `dev` and `prod` environments, blue/green deployments or the
same application available in multiple regional zones / datacenters.

An environment does not need to be created, it rather just exists as soon as a
main.jsonnet is found somewhere in the tree of a `rootDir`. This also means that
an environment by definition is equivalent to a [`baseDir`](directory-structure.md#base-directory-basedir).

## Configuration
To correctly deal with an environment, Tanka needs some additional information
about it.

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
|----------------------|--------------------------------------------------|
| `apiVersion`         | currently only `tanka.dev/v1alpha1` is available |
| `kind`               | always `Environment`                             |
| `metadata.name`      | *automatically set to the directory name*        |
| `metadata.labels`    | descriptive `key:value` pairs                    |
| **`spec.apiServer`** | The Kubernetes endpoint to use                   |
| **`spec.namespace`** | All objects will be created in this namespace    |

Everything written in **bold** is required, the other fields may be omitted.
