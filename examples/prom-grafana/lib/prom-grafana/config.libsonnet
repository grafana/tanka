{
  // +:: is important (we don't want to override the
  // _config object, just add to it)
  _config+:: {
    // define a namespace for this library
    promgrafana: {
      grafana: {
        port: 3000,
        name: "grafana",
      },
      prometheus: {
        port: 9090,
        name: "prometheus"
      }
    }
  },

  // again, make sure to use +::
  _images+:: {
    promgrafana: {
      grafana: "grafana/grafana",
      prometheus: "prom/prometheus",
    }
  }
}
