---
name: Parameterizing
menu: Tutorial
route: /tutorial/parameters
---

# Parameterizing

Deploying using Tanka worked well, but it did not really improve the situation
in terms of maintainability and readability.

To do so, the following sections will explore some ways Jsonnet provides us with.

## Config object

The most straightforward thing to do is creating a hidden object that holds all
actual values in a single place to be consumed by the actual resources.

Luckily, Jsonnet has the `key:: "value"` stanza for private fields. Such are
only available during compiling and will be removed from the actual output.

Such an object could look like this:

```jsonnet
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
  }
}
```

We can then replace hardcoded values with a reference to this object:

```diff
{ // <- This is $
  _config:: { /* .. */ },
  grafana: {
    service: {
      apiVersion: 'v1',
      kind: 'Service',
      metadata: {
        labels: {
-         name: 'grafana',
+         name: $._config.grafana.name, // $ refers to the outermost object
        },
-       name: 'grafana',
+       name: $._config.grafana.name,
      },
      spec: {
        ports: [{
-           name: 'grafana-ui',
+           name: '%s-ui' % $._config.grafana.port, // printf-style formatting
-           port: 3000,
+           port: $._config.grafana.port,
-           targetPort: 3000,
+           targetPort: $._config.grafana.port,
        }],
        selector: {
-          name: 'grafana',
+          name: $._config.grafana.name,
        },
        type: 'NodePort',
      },
    },
  },
}
```

Here we see that we can easily refer to other parts of the configuration using
the outer-most object `$` (the root level). Every value is just a regular
variable that you can refer to using the same familiar syntax from other C-like
languages.

Now we do not only have a single place to change tunables, but also won't suffer
from mismatching labels and selectors anymore, as they are defined in a single
place and all changed at once.
