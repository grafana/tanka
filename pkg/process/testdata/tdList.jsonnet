local k = (import "./k8s.libsonnet");

local deployment = (import './resources.jsonnet').deployment;
local service = (import './resources.jsonnet').service;

// NOTE: This testdata also needs Unwrap() in addition to Process()
{
  deep: {
    foo: k.list([deployment, service]),
  },
  flat: {
    "foo.items[0]": deployment,
    "foo.items[1]": service,
  },
}
