"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const vscode = require("vscode");
const fs_1 = require("./fs");
/**
 * Gets saved accent
 */
function getAccent() {
    return getCustomSettings().accent;
}
exports.getAccent = getAccent;
/**
 * Gets custom settings
 */
function getCustomSettings() {
    return vscode.workspace.getConfiguration().get('materialTheme', {});
}
exports.getCustomSettings = getCustomSettings;
/**
 * Get showReloadNotification
 */
function isReloadNotificationEnable() {
    return vscode.workspace.getConfiguration().get('materialTheme.showReloadNotification');
}
exports.isReloadNotificationEnable = isReloadNotificationEnable;
/**
 * Checks if a given string could be an accent
 */
function isAccent(accentName, defaults) {
    return Boolean(Object.keys(defaults.accents).find(name => name === accentName));
}
exports.isAccent = isAccent;
/**
 * Determines if the passing theme id is a material theme
 */
function isMaterialTheme(themeName) {
    const packageJSON = fs_1.getPackageJSON();
    return Boolean(packageJSON.contributes.themes.find(contrib => contrib.label === themeName));
}
exports.isMaterialTheme = isMaterialTheme;
/**
 * Sets a custom property in custom settings
 */
function setCustomSetting(settingName, value) {
    return vscode.workspace.getConfiguration().update(`materialTheme.${settingName}`, value, true).then(() => settingName);
}
exports.setCustomSetting = setCustomSetting;
/**
 * Updates accent name
 */
function updateAccent(accentName) {
    return setCustomSetting('accent', accentName);
}
exports.updateAccent = updateAccent;
//# sourceMappingURL=settings.js.map