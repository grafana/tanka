[
  {
    apiVersion: 'tanka.dev/v1alpha1',
    kind: 'Environment',
    metadata: {
      name: env.namespace,
      labels: {
        type: 'inline',
      },
    },
    spec: {
      apiServer: 'https://localhost',
      namespace: env.namespace,
    },
    data:
      {
        deployment: {
          apiVersion: 'apps/v1',
          kind: 'Deployment',
          metadata: {
            name: 'my-deployment',
          },
        },
        service: {
          apiVersion: 'v1',
          kind: 'Service',
          metadata: {
            name: 'my-service',
          },
        },
      } +
      (if env.hasConfigMap then {
         configMap: {
           apiVersion: 'v1',
           kind: 'ConfigMap',
           metadata: {
             name: 'my-configmap',
           },
         },
       } else {}),
  }

  for env in [
    {
      namespace: 'inline-namespace1',
      hasConfigMap: true,
    },
    {
      namespace: 'inline-namespace2',
      hasConfigMap: false,
    },
  ]
]
