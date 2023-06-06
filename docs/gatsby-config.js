const path = require("path")

module.exports = {
  pathPrefix: process.env.PATH_PREFIX || "",
  siteMetadata: {
    title: `Grafana Tanka`,
    description: `Flexible, reusable and concise configuration for Kubernetes`,
    author: `@sh0rez`,
  },
  plugins: [
    `gatsby-plugin-sharp`,
    `gatsby-plugin-catch-links`,
    {
      resolve: "gatsby-theme-docz",
      options: {
        gatsbyRemarkPlugins: [
          {
            resolve: `gatsby-remark-images`,
            options: {
              sizeByPixelDensity: true,
              withWebp: true,
            },
          },
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
    {
      resolve: `gatsby-plugin-manifest`,
      options: {
        name: "Grafana Tanka",
        short_name: "Tanka",
        start_url: "/",
        display: `standalone`,
        icon: `img/logo_black.svg`,
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

    {
      resolve: `gatsby-plugin-algolia-docsearch`,
      options: require("./algolia.json"),
    },
  ],
}
