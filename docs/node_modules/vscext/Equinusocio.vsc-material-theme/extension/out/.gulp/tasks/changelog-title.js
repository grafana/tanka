"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const fs = require("fs");
const gulp = require("gulp");
const files_1 = require("../../src/consts/files");
exports.default = gulp.task('changelog-title', () => {
    fs.writeFileSync('./CHANGELOG.md', fs.readFileSync('CHANGELOG.md', files_1.CHARSET).replace('# Change Log', '# Material Theme Changelog'), { encoding: files_1.CHARSET });
});
//# sourceMappingURL=changelog-title.js.map