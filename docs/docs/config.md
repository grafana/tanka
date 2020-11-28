---
name: "Configuration Reference"
route: "/config"
menu: "References"
---

# Configuration Reference

Tanka's behavior can be customized per Environment using a file called `spec.json`

## File format

```json
{
  // Config format revision. Currently only "v1alpha1"
  "apiVersion": "v1alpha1",
  // Always "Environment". Reserved for future use
  "kind": "Environment",

  // Descriptive fields
  "metadata": {
    // Name of the Environment. Automatically set to the relative
    // path from the project root
    "name": "<string>",

    // Arbitrary key:value string pairs. Not parsed by Tanka
    "labels": { "<string>": "<string>" }
  },

  // Properties influencing Tanka's behavior
  "spec": {
    // The Kubernetes cluster to use.
    // Must be the full URL, e.g. https://cluster.fqdn:6443
    "apiServer": "<url>",

    // Default namespace for objects that don't explicitely specify one
    "namespace": "<string>" | default = "default",

    // diffStrategy to use. Automatically chosen by default based on
    // the availability of "kubectl diff".
    // - native: uses "kubectl diff". Recommended
    // - subset: fallback for k8s versions below 1.13.0
    "diffStrategy": "[native, subset]" | default = "auto",

    // Whether to add a "tanka.dev/environment" label to each created resource.
    "injectLabels": <boolean> | default = false,

    // Whether to add a "tanka.dev/prune-mark" label to each created resource.
    // Required for garbage collection ("tk prune").
    "pruneMark": <boolean> | default = false
  }
}
```

## Jsonnet access

It is possible to access above data from Jsonnet:

```jsonnet
local tk = import "tk";

{
  // The cluster IP
  cluster: tk.env.spec.apiServer,
  // The labels of your Environment
  labels: tk.env.metadata.labels,
}
```
