package main

import (
	"bytes"
	"fmt"
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
		fmt.Println(i...)
		return
	}

	// get system pager, fallback to `more`
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "more"
	}

	// invoke pager
	cmd := exec.Command(pager)
	cmd.Stdin = strings.NewReader(fmt.Sprintln(i...))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// if this fails, just print it
	if err := cmd.Run(); err != nil {
		fmt.Println(i...)
	}
}

func highlight(lang, s string) string {
	var buf bytes.Buffer
	if err := quick.Highlight(&buf, s, lang, "terminal", "vim"); err != nil {
		log.Fatalln("Highlighting:", err)
	}
	return buf.String()
}
