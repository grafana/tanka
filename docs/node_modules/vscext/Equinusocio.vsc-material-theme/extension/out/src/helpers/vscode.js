"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const vscode = require("vscode");
/**
 * Gets your current theme ID
 */
function getCurrentThemeID() {
    return vscode.workspace.getConfiguration().get('workbench.colorTheme');
}
exports.getCurrentThemeID = getCurrentThemeID;
/**
 * Reloads current vscode window.
 */
function reloadWindow() {
    vscode.commands.executeCommand('workbench.action.reloadWindow');
}
exports.reloadWindow = reloadWindow;
//# sourceMappingURL=vscode.js.map