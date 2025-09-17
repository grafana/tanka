{
  // This lib file imports from an environment directory. So when we look for
  // the importers of `config.jsonnet` we should find environments that import
  // this file - those are the ones that need to be re-exported.
  environment_config: (import '../../environments/lib-imports-environment/config.jsonnet') + {
    lib_enhancement: 'added by lib',
  },
}
