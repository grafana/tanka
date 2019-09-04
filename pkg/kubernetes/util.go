package kubernetes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// diff computes the differences between the strings `is` and `should` using the
// UNIX `diff(1)` utility.
func diff(name, is, should string) (string, error) {
	dir, err := ioutil.TempDir("", "diff")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	if err := ioutil.WriteFile(filepath.Join(dir, "LIVE-"+name), []byte(is), os.ModePerm); err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "MERGED-"+name), []byte(should), os.ModePerm); err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	merged := filepath.Join(dir, "MERGED-"+name)
	live := filepath.Join(dir, "LIVE-"+name)
	cmd := exec.Command("diff", "-u", "-N", live, merged)
	cmd.Stdout = &buf
	err = cmd.Run()

	// the diff utility exits with `1` if there are differences. We need to not fail there.
	if exitError, ok := err.(*exec.ExitError); ok && err != nil {
		if exitError.ExitCode() != 1 {
			return "", err
		}
	}

	out := buf.String()
	if out != "" {
		out = fmt.Sprintf("diff -u -N %s %s\n%s", live, merged, out)
	}

	return out, nil
}

// diffstat uses `diffstat(1)` utility to summarize a `diff(1)` output
func diffstat(d string) (*string, error) {
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
