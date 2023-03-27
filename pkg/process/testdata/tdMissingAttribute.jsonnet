{
  deep: {
    deploy: (import './resources.jsonnet').deployment,
    service: {
      // Missing kind
      apiVersion: 'v1',
      spec: {
        selector: {
          app: 'deep',
        },
        ports: [
          {
            protocol: 'TCP',
            port: 80,
            targetPort: 8080,
          },
        ],
      },
    },
  },
}
