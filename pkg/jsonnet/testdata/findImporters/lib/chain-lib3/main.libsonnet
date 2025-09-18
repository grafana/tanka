{
  // lib3 imports env2
  importedFromEnv2: (import '../../environments/chain-env2/main.jsonnet') + {
    processed_by: "lib3"
  }
}
