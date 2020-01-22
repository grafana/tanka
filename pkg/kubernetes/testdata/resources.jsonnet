local k = (import "./k8s.libsonnet");
{
  deployment: k.deployment(),
  service: k.service(),
  namespace: k.namespace(),
}
