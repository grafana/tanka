export default {
  title: "Tanka",
  description: "Flexible, reusable and concise configuration for Kubernetes",

  public: "./public",
  ignore: ["design/**"],

  mdPlugins: [],
  hastPlugins: [],
  gatsbyRemarkPlugins: [
    {
      resolve: "gatsby-remark-vscode",
      options: {
        colorTheme: "Gruvbox Dark Medium",
        injectStyles: false,
        extensions: [
          {
            identifier: "heptio.jsonnet",
            version: "0.1.0"
          },
          {
            identifier: "jdinhlife.gruvbox",
            version: "1.4.0"
          }
        ]
      }
    }
  ],

  menu: [
    {
      name: "General",
      menu: ["Introduction", "Installation", "Getting started", "FAQ"]
    },
    {
      name: "Environments",
      menu: ["Overview", "Directory structure", "Configuration"]
    },
    {
      name: "Writing Jsonnet",
      menu: [
        "Language overview",
        "main.jsonnet",
        "Libraries",
        "Vendoring",
        "k.libsonnet",
        "Native Functions"
      ]
    },
    {
      name: "Other",
      menu: ["Command-line completion", "Output filtering"]
    }
  ]
};
