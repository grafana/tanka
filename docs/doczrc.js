export default {
  title: "Tanka",
  description: "Flexible, reusable and concise configuration for Kubernetes",

  public: "/static",
  ignore: ["design/**", ".vscodeext/**"],

  themeConfig: {
    showDarkModeSwitch: false,
  },

  menu: [
    "Introduction",
    "Installation",
    {
      name: "Tutorial",
      menu: [
        "Overview",
        "Refresher on deploying",
        "Using Jsonnet",
        "Parameterizing",
        "Abstraction",
        "Kubernetes library",
        "Environments",
      ],
    },
    {
      name: "Writing Jsonnet",
      menu: [
        "Syntax overview",
        "main.jsonnet",
        // "The global object",
        "Native Functions",
      ],
    },
    {
      name: "Libraries",
      menu: [
        "Import paths",
        // "Using libraries",
        // "Creating and structure",
        "Installing and publishing",
        "Overriding",
      ],
    },
    {
      name: "Advanced features",
      menu: [
        "Garbage collection",
        "Helm support",
        "Output filtering",
        "Exporting as YAML",
      ],
    },
    {
      name: "References",
      menu: [
        "Configuration Reference",
        "Directory structure",
        "Environment variables",
        "Command-line completion",
        "Diff strategies",
        "Namespaces",
      ],
    },

    "Frequently asked questions",
    "Known issues",
  ],
}
