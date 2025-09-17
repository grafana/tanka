local lib = import 'lib-imports-environment/main.libsonnet';

{
  // This environment uses a lib that imports from another environment
  result: lib.environment_config
}
