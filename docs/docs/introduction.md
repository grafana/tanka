---
route: /
title: Introduction
---
# Grafana Tanka
![Tanka Banner](docs/img/banner.png)

Tanka is a composable configuration utility for [Kubernetes](https://kubernetes.io/). It
leverages the [Jsonnet](https://jsonnet.org) language to realize flexible, reusable and
concise configuration.

## Highlights
* **Flexible**: The [Jsonnet data templating language](https://jsonnet.org)
  gives us much smarter ways to express our Kubernetes configuration than YAML
  does.
* **Reusable**: Code can be refactored into libraries, they can be imported
  wherever you like and even shared on GitHub!
* **Concise**: Using the Kubernetes library and abstraction, you will never see boilerplate again!
* **Work with confidence**: `tk diff` allows to check all changes before they
  will be applied. Stop guessing and make sure it's all good.
* **Used in production**: While still a very young project, Tanka is used
  internally for all our Kubernetes configuration needs.
* **Fully open source**: This is an open-source project. It is free as in beer and as in speech and this will never change.


## Getting started
To get started, [install Tanka](/install) first, and then
[follow the tutorial](/tutorial/overview). This should get you
on track quickly.
