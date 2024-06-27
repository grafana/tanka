# Acceptance tests

These tests aim to cover some e2e use-cases like creating a new Tanka
environment and pushing it up to an ephemeral Kubernetes cluster.

To run these, you need to have the Dagger CLI >= 0.11 installed. Then you can
execute the tests like this from the *root directory* of the project:

```
make acceptance-tests
```
