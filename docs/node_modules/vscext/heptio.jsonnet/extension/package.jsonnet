local package = import "./package.libsonnet";

local contributes = package.contributes;
local event = package.event;
local grammar = package.contributes.grammar;
local keybinding = package.contributes.keybinding;
local language = package.contributes.language;
local languageSpec = package.languageSpec;

local jsonnetLanguage = languageSpec.Default(
  "jsonnet", "Jsonnet", [".jsonnet",".libsonnet"]);

local preview = contributes.command.Default(
  "jsonnet.preview",
  "Jsonnet: Open Preview");

local previewToSide = contributes.command.Default(
  "jsonnet.previewToSide",
  "Jsonnet: Open Preview to the Side");

local previewKeybinding = keybinding.FromCommand(
  previewToSide, "editorFocus", "shift+ctrl+i", mac="shift+cmd+i");

package.Default() +
package.Name(jsonnetLanguage.name) +
package.DisplayName(jsonnetLanguage.displayName) +
package.Description("Language support for Jsonnet") +
package.Version("0.0.15") +
package.Publisher("heptio") +
package.License("SEE LICENSE IN 'LICENSE' file") +
package.Homepage("https://github.com/heptio/vscode-jsonnet/blob/master/README.md") +
package.Category("Languages") +
package.ActivationEvent(event.OnLanguage(jsonnetLanguage.name)) +
package.ActivationEvent(event.OnCommand(previewToSide.command)) +
package.ActivationEvent(event.OnCommand(preview.command)) +
package.Main("./out/client/extension") +

// Repository.
package.repository.Default(
  "git", "https://github.com/heptio/vscode-jsonnet.git") +

// Engines.
package.engines.VsCode("^1.10.0") +

// Contribution points.
package.contributes.Language(language.FromLanguageSpec(
  jsonnetLanguage, "./language-configuration.json")) +
package.contributes.Grammar(grammar.FromLanguageSpec(
  jsonnetLanguage, "source.jsonnet", "./syntaxes/jsonnet.tmLanguage.json")) +
package.contributes.Command(previewToSide) +
package.contributes.Command(preview) +
package.contributes.Keybinding(previewKeybinding) +
package.contributes.DefaultConfiguration(
  "Jsonnet configuration",
  contributes.configuration.DefaultStringProperty(
    "jsonnet.executablePath", "Location of the `jsonnet` executable.") +
  contributes.configuration.DefaultArrayProperty(
    "jsonnet.libPaths",
    "Additional paths to search for libraries when compiling Jsonnet code.") +
  contributes.configuration.DefaultObjectProperty(
    "jsonnet.extStrs", "External strings to pass to `jsonnet` executable.") +
  contributes.configuration.DefaultEnumProperty(
    "jsonnet.outputFormat",
    "Preview output format (yaml / json)",
    ["json", "yaml"],
    "yaml")) +
// Everything else.
{
  scripts: {
    "vscode:prepublish": "tsc -p ./",
    compile: "tsc -watch -p ./",
    "compile-once": "tsc -p ./",
    "compile-site": "browserify ./out/site/main.js > ksonnet.js",
    postinstall: "node ./node_modules/vscode/bin/install",
    test: "node ./node_modules/vscode/bin/test"
  },
  dependencies: {
    "js-yaml": "^3.0.0",
    "immutable": "^3.8.1",
    "vscode-languageclient": "^3.1.0",
    "vscode-languageserver": "^3.1.0",
  },
  devDependencies: {
    browserify: "^14.3.0",
    typescript: "^2.3.2",
    vscode: "^1.0.0",
    mocha: "^2.3.3",
    chai: "^3.5.0",
    "@types/chai": "^3.5.0",
    "@types/node": "^6.0.40",
    "@types/mocha": "^2.2.32",
  }
}
