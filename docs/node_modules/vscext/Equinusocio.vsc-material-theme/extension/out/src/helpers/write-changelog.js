"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const path = require("path");
const fs_1 = require("./fs");
const writeDefaults = (defaults) => fs_1.writeFile(path.join('./defaults.json'), JSON.stringify(defaults, null, 2));
exports.default = () => {
    const defaults = fs_1.getDefaultValues();
    const packageJSON = fs_1.getPackageJSON();
    const newChangelog = Object.assign({}, defaults.changelog, { lastversion: packageJSON.version });
    const newDefaults = Object.assign({}, defaults, { changelog: newChangelog });
    writeDefaults(newDefaults);
};
//# sourceMappingURL=write-changelog.js.map