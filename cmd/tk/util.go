package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func pageln(i ...interface{}) {
	fPageln(strings.NewReader(fmt.Sprint(i...)))
}

// fPageln invokes the systems pager with the supplied data.
// If the PAGER environment variable is empty, no pager is used.
// If the PAGER environment variable is unset, use less with posix flags.
func fPageln(r io.Reader) {
	pagerEnv, ok := os.LookupEnv("PAGER")
	if !ok {
		// --RAW-CONTROL-CHARS  Honors colors from diff. Must be in all caps, otherwise display issues occur.
		// --quit-if-one-screen Closer to the git experience.
		// --no-init            Don't clear the screen when exiting.
		pagerEnv = "less --RAW-CONTROL-CHARS --quit-if-one-screen--no-init"
	}

	pager := strings.Split(pagerEnv, " ")

	if interactive && len(pager) > 0 {
		pagerCmd := pager[0]
		pagerArgs := strings.Split(os.Getenv("PAGER"), " ")[1:len(pager)]

		cmd := exec.Command(pagerCmd, pagerArgs...)
		cmd.Stdin = r
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			if _, err = io.Copy(os.Stdout, r); err != nil {
				log.Fatalln("Writing to Stdout:", err)
			}
		}
	} else {
		fmt.Println(r)
	}
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
