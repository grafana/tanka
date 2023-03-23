local deployment = (import './resources.jsonnet').deployment;

{
  deep: {
    deployment: deployment {
      kind: 3000,
    },
  },
  flat: {
    '.deployment': deployment,
  },
}
