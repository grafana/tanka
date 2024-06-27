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
	"strings"
	"time"
)

type Tanka struct{}

func (m *Tanka) Build(ctx context.Context, rootDir *Directory) *Container {
	return dag.Container().
		Build(rootDir)
}

func (m *Tanka) AcceptanceTests(ctx context.Context, rootDir *Directory, acceptanceTestsDir *Directory) (string, error) {
	// Determine Go version through Dockerfile
	goVersion, err := dag.Container().
		From("busybox").
		WithMountedFile("/tmp/Dockerfile", rootDir.File("Dockerfile")).
		WithExec([]string{"/bin/sh", "-c", "cat /tmp/Dockerfile | grep 'FROM golang' | grep 'as build' | awk -e '{print $2}'"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	goVersion = strings.TrimSpace(goVersion)
	if !strings.HasSuffix(goVersion, "-alpine") {
		goVersion += "-alpine"
	}

	buildContainer := m.Build(ctx, rootDir)

	k3s := dag.K3S("k3sdemo")
	k3sSrv, err := k3s.Server().Start(ctx)
	if err != nil {
		return "", err
	}
	defer k3sSrv.Stop(ctx)

	output, err := dag.Container().
		From(goVersion).
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedFile("/usr/bin/tk", buildContainer.File("/usr/local/bin/tk")).
		WithMountedFile("/usr/bin/jb", buildContainer.File("/usr/local/bin/jb")).
		WithMountedFile("/usr/bin/helm", buildContainer.File("/usr/local/bin/helm")).
		WithMountedFile("/usr/bin/kustomize", buildContainer.File("/usr/local/bin/kustomize")).
		WithMountedFile("/usr/bin/kubectl", buildContainer.File("/usr/local/bin/kubectl")).
		WithMountedDirectory("/tests", acceptanceTestsDir).
		WithEnvVariable("CACHE", time.Now().String()).
		WithServiceBinding("kubernetes", k3sSrv).
		WithFile("/root/.kube/config", k3s.Config(false)).
		WithWorkdir("/tests").
		WithExec([]string{"sed", "-i", `s/https:.*:6443/https:\/\/kubernetes:6443/g`, "/root/.kube/config"}).
		WithExec([]string{"go", "test", "./...", "-v"}).
		Stdout(ctx)
	return output, err
}
