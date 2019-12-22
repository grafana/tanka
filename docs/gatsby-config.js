const path = require("path");

module.exports = {
  siteMetadata: {
    title: `Grafana Tanka`,
    description: `Flexible, reusable and concise configuration for Kubernetes`,
    author: `@sh0rez`,
  },
  plugins: [
    `gatsby-plugin-netlify-cache`,
    {
      resolve: "gatsby-theme-docz",
      options: {
        gatsbyRemarkPlugins: [
          {
            resolve: "gatsby-remark-vscode",
            options: {
              colorTheme: "Material Theme Darker",
              injectStyles: false,
              extensionDataDirectory: path.resolve("node_modules/vscext"),
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
