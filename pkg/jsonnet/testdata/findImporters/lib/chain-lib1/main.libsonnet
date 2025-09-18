{
  // lib1 imports env1
  processedConfig: (import '../../environments/chain-env1/config.jsonnet') + {
    processed_by: "lib1"
  }
}
