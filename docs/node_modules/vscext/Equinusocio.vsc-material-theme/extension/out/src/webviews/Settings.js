"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const Webview_1 = require("./Webview");
const vscode_1 = require("vscode");
const settings_1 = require("../helpers/settings");
const fs_1 = require("../helpers/fs");
class SettingsWebview extends Webview_1.WebviewController {
    constructor(context) {
        super(context);
    }
    get filename() {
        return 'settings.html';
    }
    get id() {
        return 'materialTheme.settings';
    }
    get title() {
        return 'Material Theme Settings';
    }
    getAvailableScopes() {
        const scopes = [['user', 'User']];
        return scopes
            .concat(vscode_1.workspace.workspaceFolders !== undefined && vscode_1.workspace.workspaceFolders.length ?
            ['workspace', 'Workspace'] :
            []);
    }
    /**
     * This will be called by the WebviewController when init the view
     * passing as `window.bootstrap` to the view.
     */
    getBootstrap() {
        return {
            config: settings_1.getCustomSettings(),
            defaults: fs_1.getDefaultValues(),
            scope: 'user',
            scopes: this.getAvailableScopes()
        };
    }
}
exports.SettingsWebview = SettingsWebview;
//# sourceMappingURL=Settings.js.map