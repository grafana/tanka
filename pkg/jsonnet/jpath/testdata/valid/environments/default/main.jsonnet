{
  apiVersion: 'v1',
  kind: 'ConfigMap',
  metadata: { name: 'myConfig' },
  data: (import "nestedDir/file.jsonnet"),
}
