package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// DiffName computes the filename for use with `DiffStr`
func DiffName(m manifest.Manifest) string {
	return strings.ReplaceAll(fmt.Sprintf("%s.%s.%s.%s",
		m.APIVersion(),
		m.Kind(),
		m.Metadata().Namespace(),
		m.Metadata().Name(),
	), "/", "-")
}

// DiffStr computes the differences between the strings `is` and `should` using diff
// command specified in `KUBECTL_EXTERNAL_DIFF` (if set) or the UNIX `diff(1)` utility.
func DiffStr(name, is, should string) (string, error) {
	dir, err := os.MkdirTemp("", "diff")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	if err := os.WriteFile(filepath.Join(dir, "LIVE-"+name), []byte(is), os.ModePerm); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(dir, "MERGED-"+name), []byte(should), os.ModePerm); err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	errBuf := bytes.Buffer{}
	merged := filepath.Join(dir, "MERGED-"+name)
	live := filepath.Join(dir, "LIVE-"+name)
	command, args := diffCommand(live, merged)
	cmd := exec.Command(command, args...)
	cmd.Stdout = &buf
	cmd.Stderr = &errBuf
	err = cmd.Run()

	// the diff utility exits with `1` if there are differences. We need to not fail there.
	if exitError, ok := err.(*exec.ExitError); ok && err != nil {
		if exitError.ExitCode() != 1 {
			return "", err
		}
	}
	out := buf.String()
	if out != "" {
		out = fmt.Sprintf("%s %s\n%s", command, strings.Join(args, " "), out)
	}
	errOut := errBuf.String()
	if errOut != "" {
		out += fmt.Sprintf("%s %s\n%s", command, strings.Join(args, " "), errOut)
	}

	return out, nil
}

// diffCommand returns command and arguments to run to compute differences between
// `live` and `merged` files.
// If set, env variable `KUBECTL_EXTERNAL_DIFF` is used. By default, "diff -u -N" is used.
// For consistency, we want to process KUBECTL_EXTERNAL_DIFF just like kubectl does.
// diffComand was adapted from this kubectl function:
// https://github.com/kubernetes/kubectl/blob/ac49920c0ccb0dd0899d5300fc43713ee2dfcdc9/pkg/cmd/diff/diff.go#L173
func diffCommand(live, merged string) (string, []string) {
	diff := ""
	args := []string{live, merged}
	if envDiff := os.Getenv("KUBECTL_EXTERNAL_DIFF"); envDiff != "" {
		diff = envDiff
		diffCommand := strings.Split(envDiff, " ")
		diff = diffCommand[0]

		if len(diffCommand) > 1 {
			// Regex accepts: Alphanumeric (case-insensitive) and dash
			isValidChar := regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString
			for i := 1; i < len(diffCommand); i++ {
				if isValidChar(diffCommand[i]) {
					args = append(args, diffCommand[i])
				}
			}
		}
	} else {
		diff = "diff"
		args = append([]string{"-u", "-N"}, args...)
	}
	return diff, args
}

// Diffstat uses `diffstat(1)` utility to summarize a `diff(1)` output
func Diffstat(d string) (*string, error) {
	cmd := exec.Command("diffstat", "-C")
	buf := bytes.Buffer{}
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(d)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("invoking diffstat(1): %s", err.Error())
	}

	out := buf.String()
	return &out, nil
}

// FilteredErr is a filtered Stderr. If one of the regular expressions match, the current input is discarded.
type FilteredErr []*regexp.Regexp

func (r FilteredErr) Write(p []byte) (n int, err error) {
	for _, re := range r {
		if re.Match(p) {
			// silently discard
			return len(p), nil
		}
	}
	return os.Stderr.Write(p)
}
