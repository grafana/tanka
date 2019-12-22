{

  Default():: {
    engines: {},
    categories: [],
    activationEvents: [],
    contributes: $.contributes.Default(),
  },

  Name(name):: {name: name},
  DisplayName(displayName):: {displayName: displayName},
  Description(description):: {description: description},
  Version(version):: {version: version},
  Publisher(publisher):: {publisher: publisher},
  License(license):: {license: license},
  Homepage(homepage):: {homepage: homepage},
  Category(category):: {categories+: [category]},
  ActivationEvent(event):: {activationEvents+: [event]},
  Main(main):: {main: main},

  repository:: {
    Default(type, url):: {
      repository: {
        type: type,
        url: url,
      },
    },
  },

  engines:: {
    VsCode(vscodeVersion):: {
      engines+: {
        vscode: vscodeVersion,
      },
    },
  },

  event:: {
    OnLanguage(languageId):: "onLanguage:%s" % languageId,
    OnCommand(id):: "onCommand:%s" % id,
  },

  languageSpec:: {
    Default(name, displayName, extensions):: {
      name: name,
      displayName: displayName,
      extensions: extensions,
    }
  },

  contributes:: {
    Default():: {
      languages: [],
      grammars: [],
      commands: [],
      keybindings: [],
    },

    Language(language):: {contributes+: {languages+: [language]}},
    Grammar(grammar):: {contributes+: {grammars+: [grammar]}},
    Command(command):: {contributes+: {commands+: [command]}},
    Keybinding(keybinding):: {contributes+: {keybindings+: [keybinding]}},

    DefaultConfiguration(title, properties):: {
      contributes+: {
          configuration: {
          type: "object",
          title: title,
          properties: properties,
        },
      },
    },

    configuration:: {
      DefaultStringProperty(property, description, default=null):: {
        [property]: {
          type: "string",
          default: default,
          description: description,
        },
      },

      DefaultObjectProperty(property, description, default=null):: {
        [property]: {
          type: "object",
          default: default,
          description: description,
        },
      },

      DefaultArrayProperty(property, description, default=[]):: {
        [property]: {
          type: "array",
          default: default,
          description: description,
        },
      },

      DefaultEnumProperty(property, description, enum=[], default=null):: {
       [property]: {
         default: default,
         enum: enum,
         description: description,
       },
      },

    },

    command:: {
      Default(command, title):: {
        command: command,
        title: title,
      },
    },

    keybinding:: {
      FromCommand(command, when, key, mac=null):: {
        command: command.command,
        key: key,
        [if !(mac == null) then "mac"]: mac,
        when: when,
      },
    },

    language:: {
      FromLanguageSpec(language, configurationFile):: {
        id: language.name,
        aliases: [language.displayName, language.name],
        extensions: language.extensions,
        configuration: configurationFile,
      },
    },

    grammar:: {
      FromLanguageSpec(language, scopeName, path):: {
        language: language.name,
        scopeName: scopeName,
        path: path,
      },
    },
  },
}
