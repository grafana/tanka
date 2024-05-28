// A generated module for Tanka functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
)

type Tanka struct{}

func (m *Tanka) Build(ctx context.Context, rootDir *Directory) *Container {
	return dag.Container().
		Build(rootDir)
}

func (m *Tanka) AcceptanceTests(ctx context.Context, rootDir *Directory) error {
	_, err := dag.Container().
		From("golang:1.22-alpine").
		WithMountedFile("/usr/bin/tk", m.Build(ctx, rootDir).File("/usr/local/bin/tk")).
		WithMountedDirectory("/tests", rootDir.Directory("acceptance-tests")).
		WithWorkdir("/tests").
		WithExec([]string{"go", "test", "./...", "-v"}).
		Sync(ctx)
	return err

}
