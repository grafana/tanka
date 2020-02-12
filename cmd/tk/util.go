package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func pageln(i ...interface{}) {
	fPageln(strings.NewReader(fmt.Sprint(i...)))
}

// fPageln invokes the systems pager with the supplied data
// falls back to fmt.Println() when paging fails or non-interactive
func fPageln(r io.Reader) {
	// get system pager, fallback to `less`
	pager := os.Getenv("PAGER")
	var args []string
	if pager == "" || pager == "less" {
		// --RAW-CONTROL-CHARS  Honors colors from diff. Must be in all caps, otherwise display issues occur.
		// --quit-if-one-screen Closer to the git experience.
		// --no-init            Don't clear the screen when exiting.
		pager = "less"
		args = []string{"--RAW-CONTROL-CHARS", "--quit-if-one-screen", "--no-init"}
	}

	// invoke pager
	cmd := exec.Command(pager, args...)
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// if this fails, just print it
	if err := cmd.Run(); err != nil {
		io.Copy(os.Stdout, r)
	}
}

func colordiff(d string) io.Reader {
	exps := map[string]func(s string) bool{
		"add":  regexp.MustCompile(`^\+.*`).MatchString,
		"del":  regexp.MustCompile(`^\-.*`).MatchString,
		"head": regexp.MustCompile(`^diff -u -N.*`).MatchString,
		"hid":  regexp.MustCompile(`^@.*`).MatchString,
	}

	buf := bytes.Buffer{}
	lines := strings.Split(d, "\n")

	for _, l := range lines {
		switch {
		case exps["add"](l):
			color.New(color.FgGreen).Fprintln(&buf, l)
		case exps["del"](l):
			color.New(color.FgRed).Fprintln(&buf, l)
		case exps["head"](l):
			color.New(color.FgBlue, color.Bold).Fprintln(&buf, l)
		case exps["hid"](l):
			color.New(color.FgMagenta, color.Bold).Fprintln(&buf, l)
		default:
			fmt.Fprintln(&buf, l)
		}
	}

	return &buf
}

// writeJSON writes the given object to the path as a JSON file
func writeJSON(i interface{}, path string) error {
	out, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling: %s", err)
	}

	if err := ioutil.WriteFile(path, append(out, '\n'), 0644); err != nil {
		return fmt.Errorf("writing %s: %s", path, err)
	}

	return nil
}
