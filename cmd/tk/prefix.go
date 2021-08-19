package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-clix/cli"
)

func prefixCommands(prefix string) (cmds []*cli.Command) {
	externalCommands, err := executablesOnPath(prefix)
	if err != nil {
		// soft fail if no commands found
		return nil
	}

	for file, path := range externalCommands {
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

func executablesOnPath(prefix string) (map[string]string, error) {
	path, ok := os.LookupEnv("PATH")
	if !ok {
		// if PATH not set, soft fail
		return nil, fmt.Errorf("PATH not set")
	}

	executables := make(map[string]string)
	paths := strings.Split(path, ":")
	for _, p := range paths {
		s, err := os.Stat(p)
		if err != nil && os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !s.IsDir() {
			continue
		}

		files, err := filepath.Glob(fmt.Sprintf("%s/%s*", p, prefix))
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			base := filepath.Base(file)
			// guarding against a glob character in the prefix or path
			if !strings.HasPrefix(base, prefix) {
				continue
			}
			info, err := os.Stat(file)
			if err != nil {
				return nil, err
			}
			if !info.Mode().IsRegular() {
				continue
			}
			if info.Mode().Perm()&0111 == 0 {
				continue
			}
			executables[base] = file
		}
	}
	return executables, nil
}
