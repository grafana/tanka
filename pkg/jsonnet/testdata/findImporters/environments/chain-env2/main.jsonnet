local lib1 = import 'chain-lib1/main.libsonnet';

{
  // env2 imports lib1
  name: "env2",
  inherited_config: lib1.processedConfig + {
    enhanced_by: "env2"
  }
}
