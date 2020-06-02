{
  deep: (import './resources.jsonnet').deployment,
  flat: {
    '.': $.deep,
  },
}
