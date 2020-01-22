{
  deep: {
    deploy: (import "./resources.jsonnet").deployment,
    service: {
      "note": "invalid because apiVersion and kind are missing",
    }
  }
}
