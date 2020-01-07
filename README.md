# Grafana Tanka

![Tanka Banner](docs/img/banner.png)

[![Build Status](https://cloud.drone.io/api/badges/grafana/tanka/status.svg)](https://cloud.drone.io/grafana/tanka)
![Golang](https://img.shields.io/badge/language-Go-blue)
![GitHub contributors](https://img.shields.io/github/contributors/grafana/tanka)
![GitHub release](https://img.shields.io/github/release/grafana/tanka)
![License](https://img.shields.io/github/license/grafana/tanka)

Tanka is a composable configuration utility for
[Kubernetes](https://kubernetes.io/). It leverages the
[Jsonnet](https://jsonnet.org) language to realize flexible, reusable and
concise configuration.

## Highlights

- **:wrench: Flexible**: The
  [Jsonnet data templating language](https://jsonnet.org) gives us much smarter
  ways to express our Kubernetes configuration than YAML does.
- **:books: Reusable**: Code can be refactored into libraries, they can be
  imported wherever you like and even shared on GitHub!
- **:pushpin: Concise**: Using the Kubernetes library and abstraction, you will
  never see boilerplate again!
- **:dart: Work with confidence**: `tk diff` allows to check all changes before
  they will be applied and `tk apply` makes sure you always select the correct
  cluster. Stop guessing and make sure it's all good.
- **:rocket: Used in production**: While still a very young project, Tanka is
  used internally at Grafana Labs for all of their Kubernetes configuration needs.
- **:heart: Fully open source**: This is an open-source project. It is free as
  in beer and as in speech and this will never change.

## Getting started

To get started, [install Tanka](https://tanka.dev/install) first, and then
[follow the tutorial](https://tanka.dev/tutorial/overview2). This should get you
on track quickly.

## Additional resources

- https://jsonnet.org/, the official Jsonnet documentation provides lots of
  examples on how to use the language.
- https://github.com/grafana/jsonnet-libs: Grafana Labs' Jsonnet libraries are a
  rich set of configuration examples compatible with Tanka.

## License

Licensed Apache 2.0, see [LICENSE](LICENSE).
