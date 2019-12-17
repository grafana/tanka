---
name: "Output filtering"
route: "/other/output-filtering"
---

# Output Filtering

When a project becomes bigger over time and includes a lot of Kubernetes
objects, it may become required to operate on only a subset of them (e.g. apply
only a part of an application).

Tanka helps you with this, by allowing you to limit the used objects on the
command line using the `--target` flag. Say you are deploying an `nginx`
instance with a special `nginx.conf` and want to apply the `ConfigMap` first:

```bash
# show the ConfigMap
$ tk show -t configmap/nginx .

# all good? apply!
$ tk apply -t configmap/nginx .

# and apply everything else:
$ tk apply .
```

The syntax of the `--target` / `-t` flag is `--target=<kind>/<name>`. If
multiple objects match this pattern, all of them are used.

The `--target` / `-t` flag can be specified multiple times, to work with
multiple objects.

## Regular Expressions

The argument passed to the `--target` flag is interpreted as a
[RE2](https://golang.org/s/re2syntax) regular expression.

This allows you to use all sorts of wildcards and other advanced matching
functionality to select Kubernetes objects:

```bash
# show all deployments
$ tk show . -t 'deployment/.*'

# show all objects named "loki"
$ tk show . -t '.*/loki'
```

### Gotchas

When using regular expressions, there are some things to watch out for:

#### Line Anchors

Tanka automatically surrounds your regular expression with line anchors:

```text
^<your expression>$
```

For example, `--target 'deployment/.*'` becomes `^deployment/.*$`.

#### Quoting

Regular expressions may consist of characters that have special meanings in
shell. Always make sure to properly quote your regular expression using **single
quotes**.

```zsh
# shell attempts to match the wildcard itself:
zsh-5.4.2$ tk show . -t deployment/.*
zsh: no matches found: deployment/.*

# properly quoted:
zsh-5.4.2$ tk show . -t 'deployment/.*'
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
# ...
```
