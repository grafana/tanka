package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/grafana/tanka/dagger/internal/dagger"
)

type Tanka struct{}

func (m *Tanka) Build(ctx context.Context, rootDir *dagger.Directory) *dagger.Container {
	buildArgs := make([]dagger.BuildArg, 0, 2)
	return dag.Container().
		Build(rootDir, dagger.ContainerBuildOpts{
			BuildArgs: buildArgs,
		})
}

func (m *Tanka) GetGoVersion(ctx context.Context, file *dagger.File) (string, error) {
	versionPattern := regexp.MustCompile(`^go ((\d+)\.(\d+)(\.(\d+))?)$`)
	content, err := file.Contents(ctx)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(content, "\n") {
		matches := versionPattern.FindStringSubmatch(line)
		if len(matches) < 2 {
			continue
		}
		return matches[1], nil
	}
	return "", fmt.Errorf("no Go version found")
}

func (m *Tanka) AcceptanceTests(ctx context.Context, rootDir *dagger.Directory, acceptanceTestsDir *dagger.Directory) (string, error) {
	goVersion, err := m.GetGoVersion(ctx, rootDir.File("go.mod"))
	if err != nil {
		return "", err
	}
	buildContainer := m.Build(ctx, rootDir)

	k3s := dag.K3S("k3sdemo")
	k3sSrv, err := k3s.Server().Start(ctx)
	if err != nil {
		return "", err
	}
	defer k3sSrv.Stop(ctx)

	goCache := dag.CacheVolume("acceptance-tests-gomodules")

	output, err := dag.Container().
		From(fmt.Sprintf("golang:%s-alpine", goVersion)).
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
		WithMountedCache("/go/pkg", goCache).
		WithExec([]string{"go", "test", "./...", "-v"}).
		Stdout(ctx)
	return output, err
}
