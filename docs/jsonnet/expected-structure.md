# Expected structure

Tanka evaluates the `main.jsonnet` file of your [Environment](/environments) and
filters the output (either Object or Array) for valid Kubernetes objects.  
An object is considered valid if it has both, a `kind` and a `apiVersion` set.

!!! warning
    This behaviour is going to change in the future, `metadata.name` will
    also become required.

## Deeply nested object (Recommended)
The most commonly used structure is a single big object that includes all of
your configs to be applied to the cluster nested under keys.  
How deeply encapsulated the actual object is does not matter, Tanka will
traverse down until it finds something that has both, a `kind` and an
`apiVersion`.  

??? Example
    ```json
    {
      "prometheus": {
        "service": { // Service nested one level
          "apiVersion": "v1",
          "kind": "Service",
          "metadata": {
            "name": "promSvc"
          }
        },
        "deployment": {
          "apiVersion": "apps/v1",
          "kind": "Deployment", // kind ..
          "metadata": {
            "name": "prom" // .. and metadata.name are required
                          // to indentify a valid object.
          }
        }
      },
      "web": {
        "nginx": {
          "deployment": { // Deployment nested two levels
            "apiVersion": "apps/v1",
            "kind": "Deployment",
            "metadata": {
              "name": "nginx"
            }
          }
        }
      }
    }
    ```
    
Using this technique has the big benefit that it is self-documentary, as the
nesting of keys can be used to logically group related manifests, for example by
application.

!!! info
    It is also valid to use an encapsulation level of zero, which means
    just a regular object like it could be obtained from `kubectl show -o json`:
    ```json
    {
      "apiVersion": "v1",
      "kind": "Service",
      "metadata": {
        "name": "foo"
      }
    }
    ```


## Array
Using an array of objects is also fine:
```json
[
  {
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
      "name": "promSvc"
    }
  },
  {
    "apiVersion": "apps/v1",
    "kind": "Deployment",
    "metadata": {
      "name": "prom"
    }
  }
]
```

### `List` type
Users of `kubectl` might have had contact with a type called `List`. It is not
part of the official Kubernetes API but rather a pseudo-type introduced by
`kubectl` for dealing with multiple objects at once. Thus, Tanka does not
support it out of the box.

To take full advantage of Tankas features, you can manually flatten it:

```jsonnet
local list = import "list.libsonnet";

# expose the `items` array on the top level:
list.items 
```
