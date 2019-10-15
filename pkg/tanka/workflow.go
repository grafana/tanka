package tanka

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/kubernetes"
)

func Apply(baseDir string, mods ...Modifier) error {
	opts := parseModifiers(mods)

	rec, kube, err := parse(baseDir, opts)
	if err != nil {
		return err
	}

	diff, err := kube.Diff(rec, kubernetes.DiffOpts{})
	if err != nil {
		return errors.Wrap(err, "diffing")
	}
	if diff == nil {
		tmp := "Warning: There are no differences. Your apply may not do anything at all."
		diff = &tmp
	}

	if opts.wWarn == nil {
		opts.wWarn = os.Stderr
	}
	fmt.Fprintln(opts.wWarn, *diff)

	return kube.Apply(rec, opts.apply)
}

func Diff(baseDir string, mods ...Modifier) (*string, error) {
	opts := parseModifiers(mods)

	rec, kube, err := parse(baseDir, opts)
	if err != nil {
		return nil, err
	}

	return kube.Diff(rec, opts.diff)
}

func Show(baseDir string, mods ...Modifier) (string, error) {
	opts := parseModifiers(mods)

	rec, kube, err := parse(baseDir, opts)
	if err != nil {
		return "", err
	}

	s, err := kube.Fmt(rec)
	if err != nil {
		return "", errors.Wrap(err, "pretty printing state")
	}
	return s, nil
}
