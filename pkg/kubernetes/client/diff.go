package client

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// DiffServerSide takes the desired state and computes the differences on the
// server, returning them in `diff(1)` format
func (k Kubectl) DiffServerSide(data manifest.List) (*string, error) {
	cmd := k.ctl("diff", "-f", "-")

	raw := bytes.Buffer{}
	cmd.Stdout = &raw

	fw := FilterWriter{filters: []*regexp.Regexp{regexp.MustCompile(`exit status \d`)}}
	cmd.Stderr = &fw

	cmd.Stdin = strings.NewReader(data.String())

	err := cmd.Run()
	if diffErr := parseDiffErr(err, fw.buf, k.Info().ClientVersion); diffErr != nil {
		return nil, diffErr
	}

	s := raw.String()
	if s == "" {
		return nil, nil
	}

	return &s, nil
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
