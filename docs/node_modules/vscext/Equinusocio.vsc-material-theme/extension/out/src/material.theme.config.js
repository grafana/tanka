"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const vscode_1 = require("vscode");
const ThemeCommands = require("./commands");
const settings_1 = require("./helpers/settings");
const messages_1 = require("./helpers/messages");
const check_installation_1 = require("./helpers/check-installation");
const write_changelog_1 = require("./helpers/write-changelog");
const ReleaseNotes_1 = require("./webviews/ReleaseNotes");
function activate(context) {
    return __awaiter(this, void 0, void 0, function* () {
        const installationType = check_installation_1.default();
        const releaseNotesView = new ReleaseNotes_1.ReleaseNotesWebview(context);
        write_changelog_1.default();
        if (installationType.isFirstInstall) {
            yield messages_1.installationMessage();
        }
        const shouldShowChangelog = (installationType.isFirstInstall || installationType.isUpdate) && (yield messages_1.changelogMessage());
        if (shouldShowChangelog) {
            releaseNotesView.show();
        }
        // Registering commands
        vscode_1.commands.registerCommand('materialTheme.setAccent', () => __awaiter(this, void 0, void 0, function* () {
            const accentPicked = yield ThemeCommands.accentsQuickPick();
            yield ThemeCommands.accentsSetter(accentPicked);
            yield settings_1.updateAccent(accentPicked);
        }));
        vscode_1.commands.registerCommand('materialTheme.showReleaseNotes', () => releaseNotesView.show());
    });
}
exports.activate = activate;
//# sourceMappingURL=material.theme.config.js.map