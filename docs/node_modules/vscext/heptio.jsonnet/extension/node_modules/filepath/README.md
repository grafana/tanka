Filepath
========

A cross platform interface for working with the file system in Node.js programs. Yes, it works with both posix and win32. So there.

[![NPM](https://nodei.co/npm/filepath.png?downloads=true)](https://nodei.co/npm/filepath/)

[![npm version](https://badge.fury.io/js/filepath.svg)](https://badge.fury.io/js/filepath)

__Built by [@kixxauth](https://twitter.com/kixxauth)__

Installation
------------
The most common use of Filepath is to use it as a library. In that case, just include it in your Node.js project by adding a line for "filepath" in your `pacakge.json` dependencies. For more information about your `package.json` file, you should check out the npm documentation by running `npm help json`.

Alternatively, you can quickly install Filepath for use in a project by running

	npm install filepath

which will install filepath in your `node_modules/` folder.

Quick Start
-----------
### Load the module
```JS
var filepath = require('filepath');
```

### Create a new FilePath instance
FilePath.create just takes a string to create a new path object:
```JS
var path = filepath.create(__filname);
```

It's important to remember that a FilePath instance is *not* a String. The 'path' property of a FilePath instance is the string representation of the FilePath instance, which is the same thing as calling .toString().
```JS
console.log(path);
// "{ [String: '/Users/kris/filepath/README.md'] path: '/Users/kris/filepath/README.md' }"

path.path; // "/Users/kris/filepath/README.md"
path.valueOf(); // "/Users/kris/filepath/README.md"
path.toString(); // "/Users/kris/filepath/README.md"
path + ''; // "/Users/kris/filepath/README.md"

assert(path.path === path.toString())
```

API Reference
-------------

#### Class Methods
* [.create()](#create)
* [.root()](#root)
* [.home()](#home)

#### Instance Methods

##### Manipulation
* [#append](#append)
* [#resolve](#resolve)
* [#dir](#dir)
* [#copy](#copy)
* [#remove](#remove)
* [#relative](#relative) (returns a String)

##### To String
* [#toString](#tostring)
* [#valueOf](#valueof)
* [#basename](#basename)
* [#extname](#extname)
* [#split](#split)
* [#relative](#relative)

##### Tests
* [#exists](#exists)
* [#isFile](#isfile)
* [#isDirectory](#isdirectory)

##### Reading and Writing
* [#read](#read)
* [#write](#write)
* [#require](#require)
* [#copy](#copy)

##### Streams
* [#newReadStream](#newreadstream)
* [#newWriteStream](#newwritestream)

##### Directories
* [#mkdir](#mkdir)
* [#list](#list)
* [#recurse](#recurse)

#### [About Promises](#promises)
#### [About Error Handling](#error-handling)

### Class Methods

#### .create()
Returns a new FilePath instance. Defaults to the current working directory if you don't pass any arguments.
```JS
var path = filepath.create();
assert(path.toString() === process.cwd());
path.toString();
// "/Users/kris/projects/filepath"
```
Joins multiple arguments into a single path object.
```JS
var path = filepath.create(__dirname, 'foo', 'bar');
assert(path.toString() === __dirname + '/foo/bar');
path.path; // "/Users/kris/projects/filepath/foo/bar"
```

#### .root()
Returns a FilePath instance representing the root system path.
```JS
// On a posix system:
assert(filepath.root().toString() === '/');
```

#### .home()
Returns a FilePath instance representing the users's home directory. This is achieved using environment variables `process.env.HOME` on posix and `process.env.USERPROFILE` on win32.
```JS
assert(filepath.home().toString() === '/home/kris');
```

### Instance Methods

#### #append()
Joins an arbitrary number of arguments and appends them onto the path. Returns a new FilePath instance, leaving the original intact.
```JS
var path1 = filepath.create(__dirname);
var path2 = path1.append('foo', 'bar');
var path3 = path1.append('baz');

assert(path1.toString() === __dirname);
assert(path2.toString() === __dirname + '/foo/bar');
assert(path3.toString() === __dirname + '/baz');
```

#### #resolve()
Resolves a relative path with this one. Returns a new FilePath instance, leaving the original intact.
```JS
var path = filepath
  .create('/home/kris/filepath/lib')
  .resolve('../README.md');

assert(path.toString() === '/home/kris/filepath/README.md');
```

#### #dir()
Pops off the file or directory basename. The same as doing `../` on a posix system. Returns a new FilePath instance.
```JS
var path = filepath.create('/home/kris/filepath').dir();

assert(path.toString() === '/home/kris');
```

#### #copy()
Copies the current path to the given path. Resolves with a new FilePath instance representing the new location. Also can be invoked synchronously.

See also: [Promises](#promises) and [Error Handling](#error-handling)
```JS
filepath
  .create(__filename)
  .copy('/tmp/README.md')
  .then(function (target) {
    // The callback value `target` is a new FilePath instance.
    assert(target.toString() === '/tmp/README.md');
  })
  .catch(console.error);
```

Pass in mixed parts as the target.
```JS
var targetDir = filepath.root().append('tmp');
filepath
  .create(__filename)
  .copy(targetDir, 'README.md');
```

Or you can copy a file *synchronously*:
```JS
var target = filepath
  .create(__filename)
  .copy('/tmp/README.md', {sync: true});
assert(target.toString() === '/tmp/README.md');
```

#### #remove()
Removes a FilePath. This is done by calling native Node.js `fs.unlinkSync`. There is no asynchronous pattern for #remove(). Returns the FilePath instance.
```JS
var path = filepath.create(__filename);
assert(path.exists());
path.remove();
assert(!path.exists());
```

#### #toString()
Returns the stringified version of a FilePath. This is the same thing as the `.path` attribute, and doing `path + ''`.
```JS
path.toString(); // "/Users/kris/filepath/README.md"
path.path; // "/Users/kris/filepath/README.md"
path.valueOf(); // "/Users/kris/filepath/README.md"
path + ''; // "/Users/kris/filepath/README.md"
```

#### #valueOf()
Same as #toString().

#### #basename()
Returns the last part of the path only. If you pass in the extension string, it will not be included in the returned part. Note that #basename() returns a *String* and *not* a FilePath instance.
```JS
var path = filepath.create('/home/kris/projects/filepath/README.md');
assert(path.basename() === 'README.md');
assert(path.basename('.md') === 'README');
```

#### #extname()
Returns the extension of the last part of the path. Note that #extname() returns a *String* and *not* a FilePath instance.
```JS
var ext = filepath.create('/home/kris/projects/filepath/README.md').extname();
assert(ext === '.md');
```

#### #split()
Splits a FilePath into an Array of parts. Each element in the Array is a String.
```JS
var parts = filepath.create('/home/kris/projects/filepath/README.md').split();
assert(Array.isArray(parts));
assert(parts[0] === 'home');
assert(parts.pop() === 'README.md');
```

#### #relative()
Returns the relative String required to reach the passed in path. Note that #relative() returns a *String* and *not* a FilePath instance.
```JS
var rel = filepath
  .create('/home/kris/projects/filepath/lib')
  .relative('/home/kris/projects/filepath/test');

assert(rel === '../test');
```

#### #exists()
Check to see if a FilePath is present on the filesystem. Returns a Boolean.
```JS
var path = filepath.create(__dirname);
assert(path.exists());
assert(!path.append('foo').exists());
```

#### #isFile()
Check to see if a FilePath is a file type on the filesystem. This is accomplished using `stats.isFile()`. If the FilePath does not exist, #isFile() will return false rather than throwing an Error. Returns a Boolean.
```JS
var path = filepath.create(__filename);
assert(path.isFile());
```

#### #isDirectory()
Check to see if a FilePath is a directory type on the filesystem. This is accomplished using `stats.isDirectory()`. If the FilePath does not exist, #isFile() will return false rather than throwing an Error. Returns a Boolean.
```JS
var path = filepath.create(__dirname)
assert(path.isDirectory())
```

#### #read()
Reads the contents of a file. This can be called asynchronously (by default) or synchronously. If the path is not a file type, it will throw an `ExpectFileError`. The default encoding is 'utf8', which means you will get a String back. Set the encoding to `null` to get a Buffer instead. Any options passed in (including 'encoding') will be passed to Node.js native `fs.readFile()`.

read() returns a Promise for the file contents unless you run it synchronously. If the file does not exist, it will return null. See also: [Promises](#promises) and [Error Handling](#error-handling)
```JS
var path = filepath.create(__filename);

// #read() returns a promise object with #then() and #catch() methods.
path.read()
  .then(function assertResult(contents) {
    assert(typeof contents === 'string');
  }).catch(console.error);

// Or you can read a file *synchronously*:
var readmeContents = path.read({sync: true, encoding: null});
assert(readmeContents instanceof Buffer);
```

#### #write()
Writes a file with the given content. This can be called asynchronously (by default) or synchronously. If the path is not a file type, it will throw an `ExpectFileError`. If the file does not exist, it will be created. The default encoding is 'utf8' which means you are passing #write() a String. Set the encoding to `null` to pass a Buffer instead. Any options passed in (including 'encoding') will be passed to Node.js native `fs.writeFile()`.

write() returns a Promise for the FilePath instance unless you run it synchronously. See also: [Promises](#promises) and [Error Handling](#error-handling)
```JS
var path = filepath.create('/tmp/new_file.txt')

path.write('Hello world!\n')
  .then(function assertResult(returnedPath) {
    assert(returnedPath === path);
    assert(path.read({sync: true}) === 'Hello world!\n');
  })
  .catch(console.error);

// Or you can write a file *synchronously*:
var syncPath = filepath.create('/tmp/new_file_sync.txt');
syncPath.write('Overwrite with this text', {sync: true});
assert(syncPath.read({sync: true}) === 'Hello world!\n');
```

#### #require()
Require a module using Node.js native `require()`. You'll need to pass in the current `require` function to have the correct context for requiring a module.
```JS
var path = filepath.create(__dirname, 'index.js');
var filepath = path.require(require);
assert(filepath === filepath);
```

#### #newReadStream()
Creates a new Read Stream from the FilePath. This can be used in streaming APIs.
```JS
var FS = require('fs');
var stream = filepath.create(__filename).newReadStream();
assert(stream instanceof FS.ReadStream);
```

#### #newWriteStream()
Creates a new Write Stream from the FilePath. This can be used in streaming APIs.
```JS
var FS = require('fs');
var stream = filepath.create('/tmp/new_file.txt').newWriteStream();
assert(stream instanceof FS.WriteStream);
```

#### #mkdir()
Create a directory, unless it already exists. Will create any parent directories which do not already exist. Works kinda like 'mkdir -P'. Returns the FilePath instance.
```JS
var path = filepath.create('/tmp/some/new/deep/dir').mkdir();
assert(path instanceof filepath.FilePath);
assert(path.isDirectory());
```

#### #list()
List paths in a directory. Listing a directory returns an Array of fully resolved FilePath instances.
```JS
var li = filepath.create(__dirname).list();
assert(li[4] instanceof filepath.FilePath);
assert(li[4].toString() === '/home/kris/projects/filepath/README.md');
```

#### #recurse()
Recursively walk a directory tree. The given callbak with be called with fully resolved FilePath instances. Returns the FilePath instance.
```JS
filepath.create(__dirname).recurse(function (path) {
  assert(path instanceof filepath.FilePath);
  assert(path.toString().indexOf(__dirname) === 0);
});
```

### Promises
FilePath uses [Bluebird](https://github.com/petkaantonov/bluebird/) promises from end to end. This provides a full featured and consistent API for some of the most common asynchronous flows in your program.

Bluebird also has a unique and very useful error handling mechanism (see below).

### Error Handling
FilePath defines 4 Error types:
* FilePathError
* NotFoundError
* ExpectDirectoryError
* ExpectFileError

`NotFoundError`, `ExpectDirectoryError`, and `ExpectFileError` are subtypes of `FilePathError`. This allows you to use [Bluebirds error handling](http://bluebirdjs.com/docs/api/catch.html) in .catch() handlers:

```JS
// Catch only an ExpectFileError
filepath.create(__dirname)
  .read()
  .then(function () { /* This will not be called */ })
  .catch(filepath.ExpectFileError, function (err) {
    // handle the ExpectFileError
  })
  .catch(function (err) {
    // handle the unexpected Error
  });
```

## Testing
To run the tests, just do

  npm test

You should see the test results output.


Copyright and License
---------------------
Copyright (c) 2013-2015 by Kris Walker <kris@kixx.name> (http://www.kixx.name).

Unless otherwise indicated, all source code is licensed under the MIT license.
See LICENSE for details.
