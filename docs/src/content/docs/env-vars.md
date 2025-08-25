---
title: Environment variables
sidebar:
  order: 3
---

## TANKA_JB_PATH

**Description**: Path to the `jb` tool executable  
**Default**: `$PATH/jb`

## TANKA_KUBECTL_PATH

**Description**: Path to the `kubectl` tool executable  
**Default**: `$PATH/kubectl`

## TANKA_KUBECTL_TRACE

**Description**: Print all calls to `kubectl`  
**Default**: `false`

## TANKA_HELM_PATH

**Description**: Path to the `helm` executable  
**Default**: `$PATH/helm`

## TANKA_KUSTOMIZE_PATH

**Description**: Path to the `kustomize` executable  
**Default**: `$PATH/kustomize`

## TANKA_PAGER

**Description**: Pager to use when displaying output. Set to an empty string to disable paging.
**Default**: `$PAGER`

## PAGER

**Description**: Pager to use when displaying output. Only used if TANKA_PAGER is not set. Set to an empty string to disable paging.
**Default**: `less --RAW-CONTROL-CHARS --quit-if-one-screen --no-init`

## TANKA_DANGEROUS_ALLOW_REDIRECT

**Description**: Allow redirection of the output of `tk show` to other commands when set to `true`, same as the `--dangerous-allow-redirect ` flag. Redirection of the output of `tk show` is discouraged and disabled by default. If you want to export `.yaml` files for use with other tools, try `tk export`.
**Default**: `false`
