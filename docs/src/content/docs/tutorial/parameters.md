---
title: Parameterizing
sidebar:
  order: 4
---

Deploying using Tanka worked well, but it did not really improve the situation
in terms of maintainability and readability.

To do so, the following sections will explore some ways Jsonnet provides us with.

## Functions parameters

Defining our deployment in a single block is not the best solution.
Luckily with Jsonnet we can split our configuration into smaller, self-contained chunks.

Let's start by creating a new function in `main.jsonnet` responsible of creating a Grafana deployment:

```diff lang="jsonnet"
// envirnoments/default/main.jsonnet
local grafana() = {
  deployment: {
    apiVersion: 'apps/v1',
    kind: 'Deployment',
    metadata: {
      name: 'grafana',
    },
    spec: {
      selector: {
        matchLabels: {
          name: 'grafana',
        },
      },
      template: {
        metadata: {
          labels: {
            name: 'grafana',
          },
        },
        spec: {
          containers: [
            {
              image: 'grafana/grafana',
              name: 'grafana',
              ports: [{
                  containerPort: 3000,
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
        name: 'grafana',
      },
      name: 'grafana',
    },
    spec: {
      ports: [{
          name: 'grafana-ui',
          port: 3000,
          targetPort: 3000,
      }],
      selector: {
        name: 'grafana',
      },
      type: 'NodePort',
    },
  },
};
```

and let's use it in our main configuration:

```diff lang="jsonnet"
// environments/default/main.jsonnet
local grafana() = {
  #  ...
};

{
-  grafana: {
-    # ...
-  },
+  grafana: grafana(),
  prometheus: #...
};
```

We can then replace hardcoded values by adding parameters to our function:

```diff lang="jsonnet"
// environments/default/main.jsonnet
-local grafana() = {
+local grafana(name, port) = {
  deployment: {
    apiVersion: 'apps/v1',
    kind: 'Deployment',
    metadata: {
-      name: 'grafana',
+      name: name,
    },
    spec: {
      selector: {
        matchLabels: {
-          name: 'grafana',
+          name: name,
        },
      },
      template: {
        metadata: {
          labels: {
-            name: 'grafana',
+            name: name,
          },
        },
        spec: {
          containers: [
            {
              image: 'grafana/grafana',
-              name: 'grafana',
+              name: name,
              ports: [{
-                  containerPort: 3000,
+                  containerPort: port,
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
-        name: 'grafana',
+        name: name,
      },
-      name: 'grafana',
+      name: name,
    },
    spec: {
      ports: [{
-        name: 'grafana-ui',
-        port: 3000,
-        targetPort: 3000,
+        name: '%s-ui' % name, // printf-style formatting
+        port: port,
+        targetPort: port,
      }],
      selector: {
-        name: 'grafana',
+        name: name,
      },
      type: 'NodePort',
    },
  },
};
```

and update the usage accordingly:

```diff lang="jsonnet"
// environments/default/main.jsonnet
local grafana(name, port) = {
  # ...
};

{
-  grafana: grafana(),
+  grafana: grafana('grafana', 3000),
  prometheus: #...
};
```

:::tip
You can also set default values for function parameters:

```jsonnet
local grafana(name='grafana', port=3000) = {
  # ...
};
```

:::

Now we do not only have a single place to change tunables, but also won't suffer
from mismatching labels and selectors anymore, as they are defined in a single
place and all changed at once.

:::tip[Task]
Now do the same for the Prometheus deployment by creating a function `prometheus` that takes a `name` and a `port` as parameters.
:::
