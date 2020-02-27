package client

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/kubernetes/util"
)

// DiffServerSide takes the desired state and computes the differences on the
// server, returning them in `diff(1)` format
func (k Kubectl) DiffServerSide(data manifest.List) (*string, error) {
	existingNamespaces, err := k.Namespaces()
	if err != nil {
		return nil, err
	}

	ready, missing := separateMissingNamespace(data, k.namespace, existingNamespaces)
	cmd := k.ctl("diff", "-f", "-")

	raw := bytes.Buffer{}
	cmd.Stdout = &raw

	fw := FilterWriter{filters: []*regexp.Regexp{regexp.MustCompile(`exit status \d`)}}
	cmd.Stderr = &fw

	cmd.Stdin = strings.NewReader(ready.String())

	err = cmd.Run()
	if diffErr := parseDiffErr(err, fw.buf, k.Info().ClientVersion); diffErr != nil {
		return nil, diffErr
	}

	s := raw.String()
	for _, r := range missing {
		d, err := util.DiffStr(util.DiffName(r), "", r.String())
		if err != nil {
			return nil, err
		}
		s += d
	}

	if s != "" {
		return &s, nil
	}

	// no diff -> nil
	return nil, nil
}

// parseDiffErr handles the exit status code of `kubectl diff`. It returns err
// when an error happened, nil otherwise.
// "Differences found (exit status 1)" is not an error.
//
// kubectl >= 1.18:
// 0: no error, no differences
// 1: differences found
// >1: error
//
// kubectl < 1.18:
// 0: no error, no differences
// 1: error OR differences found
func parseDiffErr(err error, stderr string, version *semver.Version) error {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		// this error is not kubectl related
		return err
	}

	// internal kubectl error
	if exitErr.ExitCode() != 1 {
		return err
	}

	// before 1.18 "exit status 1" meant error as well ... so we need to check stderr
	if version.LessThan(semver.MustParse("1.18.0")) && stderr != "" {
		return err
	}

	// differences found is not an error
	return nil
}

func separateMissingNamespace(in manifest.List, current string, exists map[string]bool) (ready, missingNamespace manifest.List) {
	for _, r := range in {
		// namespace does not exist, also ignore implicit default ("")
		if ns := r.Metadata().Namespace(); ns != "" && !exists[ns] ||
			r.Kind() == "Namespace" && !exists[r.Metadata().Name()] ||
			!exists[current] {
			missingNamespace = append(missingNamespace, r)
			continue
		}
		ready = append(ready, r)
	}
	return
}
