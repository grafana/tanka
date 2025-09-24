local lib3 = import 'chain-lib3/main.libsonnet';

{
  // env3 imports lib3
  name: "env3",
  final_result: lib3.importedFromEnv2 + {
    finalized_by: "env3"
  }
}
