module.exports = {
  siteMetadata: {
    title: `Grafana Tanka`,
    description: `Flexible, reusable and concise configuration for Kubernetes`,
    author: `@sh0rez`,
  },
  plugins: [
    {
      resolve: "gatsby-theme-docz",
      options: {
        gatsbyRemarkPlugins: [
          {
            resolve: "gatsby-remark-vscode",
            options: {
              colorTheme: "Material Theme Darker",
              injectStyles: false,
              extensions: [
                {
                  identifier: "heptio.jsonnet",
                  version: "0.1.0",
                },
                {
                  identifier: "Equinusocio.vsc-material-theme",
                  version: "30.0.0",
                },
              ],
            },
          },
        ],
      },
    },
  ],
}
