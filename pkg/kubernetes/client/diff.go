package client

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// DiffServerSide takes the desired state and computes the differences, returning them in `diff(1)` format
// It also validates that manifests are valid server-side
func (k Kubectl) DiffServerSide(data manifest.List) (*string, error) {
	return k.diff(data, true)
}

// DiffClientSide takes the desired state and computes the differences, returning them in `diff(1)` format
func (k Kubectl) DiffClientSide(data manifest.List) (*string, error) {
	return k.diff(data, false)
}

func (k Kubectl) diff(data manifest.List, validate bool) (*string, error) {
	fw := FilterWriter{filters: []*regexp.Regexp{regexp.MustCompile(`exit status \d`)}}
	diffCmd := func(serverSide bool) (string, error) {
		args := []string{"-f", "-"}
		if serverSide {
			args = append(args, "--server-side")
		}
		cmd := k.ctl("diff", args...)

		raw := bytes.Buffer{}
		cmd.Stdout = &raw
		cmd.Stderr = &fw
		cmd.Stdin = strings.NewReader(data.String())
		err := cmd.Run()
		return raw.String(), err
	}

	if validate {
		// Running the diff server-side, this checks that the resource definitions are valid
		// However, it also diffs with server-side kubernetes elements, so it adds in a lot of elements that we shouldn't consider
		_, err := diffCmd(true)
		if diffErr := parseDiffErr(err, fw.buf, k.Info().ClientVersion); diffErr != nil {
			return nil, diffErr
		}
	}

	// Running the actual diff without considering server-side elements
	s, err := diffCmd(false)
	if diffErr := parseDiffErr(err, fw.buf, k.Info().ClientVersion); diffErr != nil {
		return nil, diffErr
	}

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
