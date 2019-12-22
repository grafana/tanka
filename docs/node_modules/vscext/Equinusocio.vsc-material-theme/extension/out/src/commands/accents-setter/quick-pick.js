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
const fs_1 = require("../../helpers/fs");
const consts_1 = require("./consts");
exports.default = () => __awaiter(this, void 0, void 0, function* () {
    const themeConfigCommon = fs_1.getDefaultValues();
    const options = Object.keys(themeConfigCommon.accents).concat(consts_1.default.PURGE_KEY);
    return vscode.window.showQuickPick(options);
});
//# sourceMappingURL=quick-pick.js.map