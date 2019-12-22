"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const accents_selector_1 = require("./lib/accents-selector");
const run = () => {
    bind();
    const { config, defaults } = window.bootstrap;
    accents_selector_1.default('[data-setting="accentSelector"]', defaults.accents, config.accent);
    console.log(defaults);
    console.log(config);
};
const bind = () => {
    document.querySelector('#fixIconsCTA').addEventListener('click', () => {
        console.log('Test click');
    });
};
run();
//# sourceMappingURL=index.js.map