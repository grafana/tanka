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
const path = require("path");
const vscode_1 = require("vscode");
const settings_1 = require("../helpers/settings");
class WebviewController extends vscode_1.Disposable {
    constructor(context) {
        // Applying dispose callback for our disposable function
        super(() => this.dispose());
        this.context = context;
    }
    dispose() {
        if (this.disposablePanel) {
            this.disposablePanel.dispose();
        }
    }
    getHtml() {
        return __awaiter(this, void 0, void 0, function* () {
            const doc = yield vscode_1.workspace
                .openTextDocument(this.context.asAbsolutePath(path.join('out/ui', this.filename)));
            return doc.getText();
        });
    }
    postMessage(message, invalidates = 'all') {
        if (this.panel === undefined) {
            return false;
        }
        const result = this.panel.webview.postMessage(message);
        // If post was ok, update invalidateOnVisible if different than default
        if (!result && this.invalidateOnVisible !== 'all') {
            this.invalidateOnVisible = invalidates;
        }
        return result;
    }
    postUpdatedConfiguration() {
        // Post full raw configuration
        return this.postMessage({
            type: 'settingsChanged',
            config: settings_1.getCustomSettings()
        }, 'config');
    }
    onPanelDisposed() {
        if (this.disposablePanel) {
            this.disposablePanel.dispose();
        }
        this.panel = undefined;
    }
    onViewStateChanged(event) {
        console.log('WebviewEditor.onViewStateChanged', event.webviewPanel.visible);
        if (!this.invalidateOnVisible || !event.webviewPanel.visible) {
            return;
        }
        // Update the view since it can be outdated
        const invalidContext = this.invalidateOnVisible;
        this.invalidateOnVisible = undefined;
        switch (invalidContext) {
            case 'config':
                // Post the new configuration to the view
                return this.postUpdatedConfiguration();
            default:
                return this.show();
        }
    }
    onMessageReceived(event) {
        return __awaiter(this, void 0, void 0, function* () {
            if (event === null) {
                return;
            }
            console.log(`WebviewEditor.onMessageReceived: type=${event.type}, data=${JSON.stringify(event)}`);
            switch (event.type) {
                case 'saveSettings':
                    // TODO: update settings
                    return;
                default:
                    return;
            }
        });
    }
    show() {
        return __awaiter(this, void 0, void 0, function* () {
            const html = yield this.getHtml();
            const rootPath = vscode_1.Uri
                .file(this.context.asAbsolutePath('./out'))
                .with({ scheme: 'vscode-resource' }).toString();
            // Replace placeholders in html content for assets and adding configurations as `window.bootstrap`
            const fullHtml = html
                .replace(/{{root}}/g, rootPath)
                .replace('\'{{bootstrap}}\'', JSON.stringify(this.getBootstrap()));
            // If panel already opened just reveal
            if (this.panel !== undefined) {
                this.panel.webview.html = fullHtml;
                return this.panel.reveal(vscode_1.ViewColumn.Active);
            }
            this.panel = vscode_1.window.createWebviewPanel(this.id, this.title, vscode_1.ViewColumn.Active, {
                retainContextWhenHidden: true,
                enableFindWidget: true,
                enableCommandUris: true,
                enableScripts: true
            });
            // Applying listeners
            this.disposablePanel = vscode_1.Disposable.from(this.panel, this.panel.onDidDispose(this.onPanelDisposed, this), this.panel.onDidChangeViewState(this.onViewStateChanged, this), this.panel.webview.onDidReceiveMessage(this.onMessageReceived, this));
            this.panel.webview.html = fullHtml;
        });
    }
}
exports.WebviewController = WebviewController;
//# sourceMappingURL=Webview.js.map