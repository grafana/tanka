"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const fs = require("fs");
const path = require("path");
const gulp = require("gulp");
const paths_1 = require("../../src/consts/paths");
const fs_1 = require("../../src/helpers/fs");
exports.default = gulp.task('build:copy-ui', callback => {
    try {
        fs_1.ensureDir(path.resolve(paths_1.PATHS.UI));
        fs.copyFileSync(path.join(paths_1.PATHS.SRC, 'webviews', 'ui', 'release-notes', 'release-notes.html'), path.join(paths_1.PATHS.UI, 'release-notes.html'));
        fs.copyFileSync(path.join(paths_1.PATHS.SRC, 'webviews', 'ui', 'release-notes', 'style.css'), path.join(paths_1.PATHS.UI, 'release-notes.css'));
    }
    catch (error) {
        return callback(error);
    }
    callback();
});
//# sourceMappingURL=copy-ui.js.map