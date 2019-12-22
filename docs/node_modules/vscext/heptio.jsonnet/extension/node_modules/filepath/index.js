var UTIL = require('util');
var FS = require('fs');
var PATH = require('path');
var Promise = require('bluebird');
var slice = Array.prototype.slice;

//
// Error Types
//

// Base Error type for all operational errors.
function FilePathError(message) {
	Error.call(this);
	Error.captureStackTrace(this, this.constructor);
	this.name = this.constructor.name;
	this.message = message;
}
UTIL.inherits(FilePathError, Error);

// Rejection when an expected path does not exist.
function NotFoundError(message) {
	FilePathError.call(this);
	Error.captureStackTrace(this, this.constructor);
	this.name = this.constructor.name;
	this.message = message;
}
UTIL.inherits(NotFoundError, FilePathError);

// Rejection when a directory was expected for a method call.
function ExpectDirectoryError(message) {
	FilePathError.call(this);
	Error.captureStackTrace(this, this.constructor);
	this.name = this.constructor.name;
	this.message = message;
}
UTIL.inherits(ExpectDirectoryError, FilePathError);

// Rejection when a file was expected for a method call.
function ExpectFileError(message) {
	FilePathError.call(this);
	Error.captureStackTrace(this, this.constructor);
	this.name = this.constructor.name;
	this.message = message;
}
UTIL.inherits(ExpectFileError, FilePathError);

// Base FilePath constructor.
function FilePath(path) {
	this.path = path;
}

FilePath.prototype = {

	resolve: function resolve(to) {
		var p;
		if (typeof to === 'string') {
			p = PATH.resolve(this.path, to);
		} else {
			p = PATH.resolve(this.path);
		}
		return FilePath.create(p);
	},

	relative: function relative(to) {
		return PATH.relative(this.path, typeof to === 'string' ? to : process.cwd());
	},

	append: function append() {
		// Join an arbitrary number of arguments.
		var args = [this.path].concat(slice.call(arguments));
		return FilePath.create.apply(null, args);
	},

	split: function split() {
		return this.path
			.replace(/\\/g, '/')
			.split('/')
			.filter(FilePath.partsFilter);
	},

	basename: function basename(ext) {
		return PATH.basename(this.path, ext);
	},

	extname: function extname() {
		return PATH.extname(this.path);
	},

	dir: function dir() {
		var p = PATH.dirname(this.path);
		return FilePath.create(p);
	},

	exists: function exists() {
		return Boolean(FS.existsSync(this.path));
	},

	isFile: function isFile() {
		var stats;
		try {
			stats = FS.statSync(this.path);
		} catch (err) {
			if (err.code === 'ENOENT') {
				return false;
			}
			throw err;
		}
		return Boolean(stats.isFile());
	},

	isDirectory: function isDirectory() {
		var stats;
		try {
			stats = FS.statSync(this.path);
		} catch (err) {
			if (err.code === 'ENOENT') {
				return false;
			}
			throw err;
		}
		return Boolean(stats.isDirectory());
	},

	newReadStream: function newReadStream(opts) {
		return FS.createReadStream(this.path, opts);
	},

	newWriteStream: function newWriteStream(opts) {
		opts = opts || (opts || {});
		if (typeof opts.encoding === 'undefined') {
			opts.encoding = 'utf8';
		}
		return FS.createWriteStream(this.path, opts);
	},

	read: function read(opts) {
		opts = (opts || Object.create(null));

		if (typeof opts.encoding === 'undefined') {
			opts.encoding = 'utf8';
		}

		var self = this;
		var promise;

		function handleError(err) {
			if (err.code === 'ENOENT') {
				return null;
			} else if (err.code === 'EISDIR') {
				err = new ExpectFileError('Cannot read "' + self.path + '"; it is a directory.');
				err.code = 'PATH_IS_DIRECTORY';
				throw err;
			}
			throw err;
		}

		if (opts.sync || opts.synchronous) {
			try {
				return FS.readFileSync(this.path, opts);
			} catch (err) {
				return handleError(err);
			}
		}

		promise = new Promise(function (resolve, reject) {
			FS.readFile(self.path, opts, function (err, data) {
				if (err) {
					try {
						return resolve(handleError(err));
					} catch (e) {
						return reject(e);
					}
				}

				return resolve(data);
			});
		});

		return promise;
	},

	write: function write(data, opts) {
		opts = (opts || Object.create(null));

		var self = this;
		var promise;
		var dir = this.dirname();

		if (!dir.exists()) {
			dir.mkdir();
		}

		function handleError(err) {
			if (err.code === 'ENOENT') {
				return null;
			} else if (err.code === 'EISDIR') {
				err = new ExpectFileError('Cannot write to "' + self.path + '"; it is a directory.');
				err.code = 'PATH_IS_DIRECTORY';
				throw err;
			}
			throw err;
		}

		if (opts.sync || opts.synchronous) {
			try {
				FS.writeFileSync(self.path, data, opts);
				return self;
			} catch (err) {
				return handleError(err);
			}
		}

		promise = new Promise(function (resolve, reject) {
			FS.writeFile(self.path, data, opts, function (err) {
				if (err) {
					try {
						return resolve(handleError(err));
					} catch (e) {
						return reject(e);
					}
				}

				return resolve(self);
			});
		});

		return promise;
	},

	copy: function copy(opts) {
		var target;
		var args = slice.call(arguments);
		var lastArg = args[args.length - 1];

		if (!args.length || lastArg instanceof FilePath || typeof lastArg === 'string') {
			opts = Object.create(null);
		} else {
			opts = args.pop();
		}

		target = FilePath.create.apply(null, args);

		// Use a buffer.
		opts.encoding = null;

		if (opts.sync || opts.synchronous) {
			var contents = this.read(opts);
			return target.write(contents, opts);
		}

		function copyContents(contents) {
			return target.write(contents, opts);
		}

		return this.read(opts).then(copyContents);
	},

	remove: function remove() {
		try {
			FS.unlinkSync(this.path);
		} catch (e) {}
		return this;
	},

	require: function filepathRequire(contextualRequire) {
		var opError;
		if (typeof contextualRequire !== 'function') {
			var err = new Error('Must pass a require function to #require().');
			err.code = 'NO_REQUIRE_CONTEXT';
			throw err;
		}
		try {
			return contextualRequire(this.path);
		} catch (err) {
			if (err.code === 'MODULE_NOT_FOUND') {
				opError = new NotFoundError(err.message);
				opError.code = err.code;
				throw opError;
			}
			throw err;
		}
	},

	list: function list() {
		var list;
		try {
			list = FS.readdirSync(this.path);
		} catch (err) {
			var e;
			if (err.code === 'ENOTDIR') {
				e = new ExpectDirectoryError('Cannot list "' + this.path + '"; it is a file.');
				e.code = 'PATH_IS_FILE';
			} else if (err.code === 'ENOENT') {
				e = new NotFoundError('Cannot list "' + this.path + '"; it does not exist.');
				e.code = 'PATH_NO_EXIST';
			}

			if (e) {
				throw e;
			}
			throw err;
		}

		return list.map(function (item) {
			return FilePath.create(this.path, item);
		}, this);
	},

	mkdir: function mkdir() {
		var self = this;
		var parts = this.resolve().toString().split(PATH.sep);
		var fullpath;

		// Shift off the empty string.
		parts.shift();

		fullpath = parts.reduce(function (fullpath, part) {
			fullpath = fullpath.append(part);
			if (fullpath.exists()) {
				if (fullpath.isDirectory()) {
					return fullpath;
				}
				var e = new ExpectDirectoryError('Cannot create directory "' + self.path + '"; it is a file.');
				e.code = 'PATH_IS_FILE';
				throw e;
			}

			FS.mkdirSync(fullpath.toString());
			return fullpath;
		}, FilePath.root());

		return FilePath.create(fullpath);
	},

	recurse: function recurse(callback) {
		var listing;
		var p = this.resolve();

		if (!p.isDirectory()) {
			return callback(p);
		}

		try {
			listing = p.list();
		} catch (err) {
			if (err.code === 'PATH_IS_FILE') {
				return p;
			}

			throw err;
		}

		listing.sort(FilePath.alphaSort).forEach(function (li) {
			callback(li);
			if (li.isDirectory()) {
				li.recurse(callback);
			}
		});

		return this;
	},

	toString: function toString() {
		return this.path;
	},

	valueOf: function valueOf() {
		return this.path;
	}
};

// For backwards compatibility:
FilePath.prototype.dirname = FilePath.prototype.dir;

//
// Class methods
//

// Create a new FilePath instance.
FilePath.create = function create() {
	var path;
	var args;

	if (arguments.length === 1 && arguments[0]) {
		path = arguments[0];
	} else if (arguments.length < 1) {
		path = process.cwd();
	} else {
		args = slice.call(arguments).map(function (item) {
			if (typeof item === 'undefined' || item === null) {
				return '';
			}
			return String(item);
		}).filter(FilePath.partsFilter);

		if (args.length < 1) {
			path = process.cwd();
		} else {
			path = PATH.join.apply(PATH, args);
		}
	}

	return new FilePath(PATH.resolve(path.toString()));
};

// Create a new FilePath instance representing the root directory.
FilePath.root = function root() {
	return FilePath.create(process.platform === 'win32' ? '\\' : '/');
};

// Create a new FilePath instance representing the home directory.
FilePath.home = function home() {
	var path = process.platform === 'win32' ? process.env.USERPROFILE : process.env.HOME;
	return FilePath.create(path);
};

FilePath.alphaSort = function alphaSort(a, b) {
	a = a.toString();
	b = b.toString();
	if (a < b) {
		return -1;
	}
	if (a > b) {
		return 1;
	}
	return 0;
};

FilePath.partsFilter = function partsFilter(part) {
	return Boolean(part);
};

//
// Public API
//

exports.FilePath = FilePath;
exports.create = exports.newPath = FilePath.create;
exports.root = FilePath.root;
exports.home = FilePath.home;
exports.FilePathError = FilePathError;
exports.NotFoundError = NotFoundError;
exports.ExpectDirectoryError = ExpectDirectoryError;
exports.ExpectFileError = ExpectFileError;
