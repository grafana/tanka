---
title: Configuration Reference
sidebar:
  order: 1
---

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

    // The Kubernetes context name(s) to use.
    // This field supports regular expressions and is mutually exclusive with apiServer field.
    "contextNames": ["<string>"],

    // Default namespace for objects that don't explicitely specify one
    "namespace": "<string>" | default = "default",

    // diffStrategy to use. Automatically chosen by default based on
    // the availability of "kubectl diff".
    // - native: uses "kubectl diff". Recommended
    // - validate: uses "kubectl diff --server-side". Safest, but slower than "native"
    // - subset: fallback for k8s versions below 1.13.0
    "diffStrategy": "[native, validate, subset]" | default = "auto",

    // Whether to add a "tanka.dev/environment" label to each created resource.
    // Required for garbage collection ("tk prune").
    "injectLabels": <boolean> | default = false
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

## Injecting labels from the command line

`tk apply`, `tk diff`, `tk show` and `tk export` accept a `--inject-label
key=value` flag that sets a label on **every** rendered resource. Unlike
`spec.resourceDefaults.labels` (which only fills in missing keys), an injected
label always overrides any pre-existing value. The flag can be repeated to set
multiple labels:

```console
$ tk apply environments/default --inject-label created-by=alice --inject-label team=infra
```

This is handy for ad-hoc identification of the resources a given run touched,
without changing the Environment's `spec.json`. Note this is unrelated to the
`spec.injectLabels` boolean, which only toggles the `tanka.dev/environment`
label used by garbage collection.
