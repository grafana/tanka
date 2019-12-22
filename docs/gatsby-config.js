const path = require("path")

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
              logLevel: "debug",
              colorTheme: "Material Theme Darker",
              injectStyles: false,
              extensionDataDirectory: path.resolve(".vscext"),
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
    `gatsby-plugin-netlify-cache`,
    {
      resolve: `gatsby-plugin-manifest`,
      options: {
        name: "Grafana Tanka",
        short_name: "Tanka",
        start_url: "/",
        display: `standalone`,
        icon: `img/tk_black.png`,
        background_color: "#ffffff",
        theme_color: "#000000",
      },
    },
    {
      resolve: `gatsby-plugin-offline`,
      options: {
        precachePages: [`/`, `/install`, `/tutorial/overview`],
      },
    },
  ],
}
