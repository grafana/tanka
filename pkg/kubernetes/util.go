package kubernetes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

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
