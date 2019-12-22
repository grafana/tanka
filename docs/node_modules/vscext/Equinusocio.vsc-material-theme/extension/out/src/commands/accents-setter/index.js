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
const vscode = require("vscode");
const fs_1 = require("./../../helpers/fs");
const consts_1 = require("./consts");
const REGEXP_HEX = /^#([0-9A-F]{6}|[0-9A-F]{8})$/i;
/**
 * Assigns colours
 */
const assignColorCustomizations = (colour) => {
    const accentsProperties = fs_1.getAccentsProperties();
    const newColour = isValidColour(colour) ? colour : undefined;
    return Object.keys(accentsProperties).reduce((acc, propName) => {
        const accent = accentsProperties[propName];
        let colorProp = newColour;
        if (colour && accent.alpha < 100) {
            colorProp = `${colour}${accent.alpha > 10 ? accent.alpha : `0${accent.alpha}`}`;
        }
        acc[propName] = colorProp;
        return acc;
    }, {});
};
/**
 * Determines if a string is a valid colour
 */
const isValidColour = (colour) => typeof colour === 'string' && REGEXP_HEX.test(colour);
/**
 * Sets workbench options
 */
const setWorkbenchOptions = (config) => vscode.workspace.getConfiguration().update('workbench.colorCustomizations', config, true)
    .then(() => true, reason => vscode.window.showErrorMessage(reason));
/**
 * VSCode command
 */
exports.default = (accent) => __awaiter(this, void 0, void 0, function* () {
    const themeConfigCommon = fs_1.getDefaultValues();
    const config = vscode.workspace.getConfiguration().get('workbench.colorCustomizations');
    switch (accent) {
        case consts_1.default.PURGE_KEY: {
            const newConfig = Object.assign({}, config, assignColorCustomizations(undefined));
            return setWorkbenchOptions(newConfig)
                .then(() => Promise.resolve(true));
        }
        default: {
            const newConfig = Object.assign({}, config, assignColorCustomizations(themeConfigCommon.accents[accent]));
            return setWorkbenchOptions(newConfig)
                .then(() => Boolean(accent));
        }
    }
});
//# sourceMappingURL=index.js.map