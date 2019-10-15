package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/alecthomas/chroma/quick"
)

// pageln invokes the systems pager with the supplied data
// falls back to fmt.Println() when paging fails or non-interactive
func pageln(i ...interface{}) {
	// no paging in non-interactive mode
	if !interactive {
		fmt.Print(i...)
		return
	}

	// get system pager, fallback to `less`
	pager := os.Getenv("PAGER")
	var args []string
	if pager == "" || pager == "less" {
		// --raw-control-chars  Honors colors from diff.
		// --quit-if-one-screen Closer to the git experience.
		// --no-init            Don't clear the screen when exiting.
		pager = "less"
		args = []string{"--raw-control-chars", "--quit-if-one-screen", "--no-init"}
	}

	// invoke pager
	cmd := exec.Command(pager, args...)
	cmd.Stdin = strings.NewReader(fmt.Sprint(i...))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// if this fails, just print it
	if err := cmd.Run(); err != nil {
		fmt.Print(i...)
	}
}

func highlight(lang, s string) string {
	var buf bytes.Buffer
	if err := quick.Highlight(&buf, s, lang, "terminal", "vim"); err != nil {
		log.Fatalln("Highlighting:", err)
	}
	return buf.String()
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
