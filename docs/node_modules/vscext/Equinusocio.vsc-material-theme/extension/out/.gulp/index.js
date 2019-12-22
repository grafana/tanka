"use strict";
function __export(m) {
    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
}
Object.defineProperty(exports, "__esModule", { value: true });
// export the tasks
__export(require("./tasks/themes"));
__export(require("./tasks/watcher"));
__export(require("./tasks/changelog-title"));
__export(require("./tasks/copy-ui"));
// export default script
exports.default = ['build:themes'];
//# sourceMappingURL=index.js.map