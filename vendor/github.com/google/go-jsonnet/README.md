# go-jsonnet

[![GoDoc Widget]][GoDoc] [![Travis Widget]][Travis] [![Coverage Status Widget]][Coverage Status]

[GoDoc]: https://godoc.org/github.com/google/go-jsonnet
[GoDoc Widget]: https://godoc.org/github.com/google/go-jsonnet?status.png
[Travis]: https://travis-ci.org/google/go-jsonnet
[Travis Widget]: https://travis-ci.org/google/go-jsonnet.svg?branch=master
[Coverage Status Widget]: https://coveralls.io/repos/github/google/go-jsonnet/badge.svg?branch=master
[Coverage Status]: https://coveralls.io/github/google/go-jsonnet?branch=master

This an implementation of [Jsonnet](http://jsonnet.org/) in pure Go.  It is
feature complete but is not as heavily exercised as the [Jsonnet C++
implementation](https://github.com/google/jsonnet).  Please try it out and give
feedback.

This code is known to work on Go 1.8 and above. We recommend always using the newest stable release of Go.

## Install instructions

```
go get github.com/google/go-jsonnet/cmd/jsonnet
```

## Build instructions (go 1.11+)

```bash
git clone github.com/google/go-jsonnet
cd go-jsonnet
go build ./cmd/jsonnet
```

## Build instructions (go 1.8 - 1.10)

```bash
go get -u github.com/google/go-jsonnet
cd $GOPATH/src/github.com/google/go-jsonnet
go get -u .
go build ./cmd/jsonnet
```

## Running tests

```bash
./tests.sh  # Also runs `go test ./...`
```

## Implementation Notes

We are generating some helper classes on types by using
http://clipperhouse.github.io/gen/.  Do the following to regenerate these if
necessary:

```bash
go get github.com/clipperhouse/gen
go get github.com/clipperhouse/set
export PATH=$PATH:$GOPATH/bin  # If you haven't already
go generate
```

## Updating and modifying the standard library

Standard library source code is kept in `cpp-jsonnet` submodule,
because it is shared with [Jsonnet C++
implementation](https://github.com/google/jsonnet).

For perfomance reasons we perform preprocessing on the standard library,
so for the changes to be visible, regeneration is necessary:

```bash
./reset_stdast_go.sh && go run cmd/dumpstdlibast/dumpstdlibast.go
```

The above command recreates `ast/stdast.go` which puts the desugared standard library into the right data structures, which lets us avoid the parsing overhead during execution.
