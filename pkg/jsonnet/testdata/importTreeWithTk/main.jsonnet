local tk = import "tk";
local helper = import "helper.libsonnet";

{
  apiVersion: 'v1',
  kind: 'ConfigMap',
  metadata: {
    name: 'test',
  },
  data: helper,
}
