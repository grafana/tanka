"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const fs_1 = require("./fs");
const splitVersion = (input) => {
    const [major, minor, patch] = input.split('.').map(i => parseInt(i, 10));
    return { major, minor, patch };
};
exports.default = () => {
    const out = {
        isUpdate: false,
        isFirstInstall: false
    };
    const defaults = fs_1.getDefaultValues();
    const packageJSON = fs_1.getPackageJSON();
    const isFirstInstall = defaults.changelog === undefined ||
        (defaults.changelog !== undefined && typeof defaults.changelog.lastversion !== 'string');
    if (isFirstInstall) {
        return Object.assign({}, out, { isFirstInstall });
    }
    const versionCurrent = splitVersion(packageJSON.version);
    const versionOld = isFirstInstall ? null : splitVersion(defaults.changelog.lastversion);
    const isUpdate = !versionOld ||
        versionCurrent.major > versionOld.major ||
        versionCurrent.minor > versionOld.minor ||
        versionCurrent.patch > versionOld.patch;
    return Object.assign({}, out, { isUpdate });
};
//# sourceMappingURL=check-installation.js.map