package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/go-clix/cli"
)

func prefixCmds(prefix string) []*cli.Command {
	ext_subcommands := map[string]string{}
	if path, ok := os.LookupEnv("PATH"); ok {
		paths := strings.Split(path, ":")
		for _, p := range paths {
			s, err := os.Stat(p)
			if err == nil && s.IsDir() {
				files, err := ioutil.ReadDir(p)
				if err != nil {
					panic(err)
				}
				for _, file := range files {
					if strings.HasPrefix(file.Name(), prefix) &&
						!file.IsDir() &&
						file.Mode().IsRegular() &&
						file.Mode().Perm()&0111 != 0 {
						ext_subcommands[file.Name()] = fmt.Sprintf("%s/%s", p, file.Name())
					}
				}
			}
		}
	}

	var cmds []*cli.Command
	for file, path := range ext_subcommands {
		cmd := &cli.Command{
			Use:   fmt.Sprintf("%s --", strings.TrimPrefix(file, prefix)),
			Short: fmt.Sprintf("external command %s", path),
			Args:  cli.ArgsAny(),
		}

		ext_command := exec.Command(path)
		if ex, err := os.Executable(); err == nil {
			ext_command.Env = append(os.Environ(), fmt.Sprintf("EXECUTABLE=%s", ex))
		}
		ext_command.Stdout = os.Stdout
		ext_command.Stderr = os.Stderr

		cmd.Run = func(cmd *cli.Command, args []string) error {
			ext_command.Args = append(ext_command.Args, args...)
			return ext_command.Run()
		}
		cmds = append(cmds, cmd)
	}
	if len(cmds) > 0 {
		return cmds
	}
	return nil
}
