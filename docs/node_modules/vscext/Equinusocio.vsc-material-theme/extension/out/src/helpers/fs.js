"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const fs = require("fs");
const path = require("path");
const files_1 = require("./../consts/files");
const paths_1 = require("../consts/paths");
function ensureDir(dirname) {
    if (!fs.existsSync(dirname)) {
        fs.mkdirSync(dirname);
    }
}
exports.ensureDir = ensureDir;
function getDefaultValues() {
    const defaults = require(path.join(paths_1.PATHS.VSIX_DIR, 'defaults.json'));
    if (defaults === undefined || defaults === null) {
        throw new Error('Cannot find defaults params');
    }
    return defaults;
}
exports.getDefaultValues = getDefaultValues;
function getAbsolutePath(input) {
    return path.join(paths_1.PATHS.VSIX_DIR, input);
}
exports.getAbsolutePath = getAbsolutePath;
function getAccentsProperties() {
    return getDefaultValues().accentsProperties;
}
exports.getAccentsProperties = getAccentsProperties;
/**
 * Gets package JSON
 */
function getPackageJSON() {
    return require(path.join(paths_1.PATHS.VSIX_DIR, './package.json'));
}
exports.getPackageJSON = getPackageJSON;
/**
 * Writes a file inside the vsix directory
 */
function writeFile(filename, filecontent) {
    const filePath = path.join(paths_1.PATHS.VSIX_DIR, filename);
    fs.writeFileSync(filePath, filecontent, { encoding: files_1.CHARSET });
}
exports.writeFile = writeFile;
//# sourceMappingURL=fs.js.map