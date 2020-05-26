local deployment = (import './resources.jsonnet').deployment;
local service = (import './resources.jsonnet').service;

{
  deep: {
    deployment: deployment,
    service: service,
  },
  flat: {
    '.deployment': deployment,
    '.service': service,
  },
}
