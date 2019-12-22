"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const gulp = require("gulp");
const path = require("path");
const paths_1 = require("./../../src/consts/paths");
/*
 * > Watcher
 * Watches files and build the themes
 */
exports.default = gulp.task('watch', () => {
    // Commented due
    // gulp.watch(path.join(Paths.SRC, `./themes/**/*.json`), ['build:themes']);
    gulp.watch(path.join(paths_1.default.SRC, './themes/**/*.json'));
});
//# sourceMappingURL=watcher.js.map