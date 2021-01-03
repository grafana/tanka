package tanka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/process"
)

// BelRune is a string of the Ascii character BEL which made computers ring in ancient times
// We use it as "magic" char for the subfolder creation as it is a non printable character and thereby will never be
// in a valid filepath by accident. Only when we include it.
const BelRune = string(rune(7))

// When exporting manifests to files, it becomes increasingly hard to map manifests back to its environment, this file
// can be used to map the files back to their environment. This is aimed to be used by CI/CD but can also be used for
// debugging purposes.
const manifestFile = "manifest.json"

// ExportEnvOpts specify options on how to export environments
type ExportEnvOpts struct {
	// formatting the filename based on the exported Kubernetes manifest
	Format string
	// extension of the filename
	Extension string
	// merge export with existing directory
	Merge bool
	// optional: only export specified Kubernetes manifests
	Targets []string
	// optional: options for parsing Environments
	ParallelOpts ParallelOpts
}

func ExportEnvironments(paths []string, to string, opts *ExportEnvOpts) error {
	// Keep track of which file maps to which environment
	fileToEnv := map[string]string{}

	// dir must be empty
	empty, err := dirEmpty(to)
	if err != nil {
		return fmt.Errorf("Checking target dir: %s", err)
	}
	if !empty && !opts.Merge {
		return fmt.Errorf("Output dir `%s` not empty. Pass --merge to ignore this", to)
	}

	for _, path := range paths {
		// select targets to export
		filter, err := process.StrExps(opts.Targets...)
		if err != nil {
			return err
		}

		// get the manifests
		loaded, err := Load(path, Opts{
			Filters: filter,
		})
		if err != nil {
			return err
		}

		env := loaded.Env
		res := loaded.Resources

		// create raw manifest version of env for templating
		env.Data = nil
		raw, err := json.Marshal(env)
		if err != nil {
			return err
		}
		var menv manifest.Manifest
		if err := json.Unmarshal(raw, &menv); err != nil {
			return err
		}

		// create template
		manifestTemplate, err := createTemplate(opts.Format, menv)
		if err != nil {
			return fmt.Errorf("Parsing format: %s", err)
		}

		// write each to a file
		for _, m := range res {
			// apply template
			name, err := applyTemplate(manifestTemplate, m)
			if err != nil {
				return fmt.Errorf("executing name template: %w", err)
			}

			// Create all subfolders in path
			path := filepath.Join(to, name+"."+opts.Extension)

			fileToEnv[path] = env.Metadata.Namespace

			// Abort if already exists
			if exists, err := fileExists(path); err != nil {
				return err
			} else if exists {
				return fmt.Errorf("File '%s' already exists. Aborting", path)
			}

			// Write manifest
			data := m.String()
			if err := writeExportFile(path, []byte(data)); err != nil {
				return err
			}
		}
	}

	// Write manifest file
	if len(fileToEnv) != 0 {
		data, err := json.MarshalIndent(fileToEnv, "", "    ")
		if err != nil {
			return err
		}
		path := filepath.Join(to, manifestFile)
		if err := writeExportFile(path, data); err != nil {
			return err
		}
	}

	return nil
}

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func dirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if os.IsNotExist(err) {
		return true, os.MkdirAll(dir, os.ModePerm)
	} else if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func writeExportFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating filepath '%s': %s", filepath.Dir(path), err)
	}

	return ioutil.WriteFile(path, data, 0644)
}

func createTemplate(format string, env manifest.Manifest) (*template.Template, error) {
	// Replace all os.path separators in string with BelRune for creating subfolders
	replaceFormat := strings.Replace(format, string(os.PathSeparator), BelRune, -1)

	envMap := template.FuncMap{"env": func() manifest.Manifest { return env }}

	template, err := template.New("").
		Funcs(sprig.TxtFuncMap()). // register Masterminds/sprig
		Funcs(envMap).             // register environment mapping
		Parse(replaceFormat)       // parse template
	if err != nil {
		return nil, err
	}
	return template, nil
}

func applyTemplate(template *template.Template, m manifest.Manifest) (path string, err error) {
	buf := bytes.Buffer{}
	if err := template.Execute(&buf, m); err != nil {
		return "", err
	}

	// Replace all os.path separators in string in order to not accidentally create subfolders
	path = strings.Replace(buf.String(), string(os.PathSeparator), "-", -1)
	// Replace the BEL character inserted with a path separator again in order to create a subfolder
	path = strings.Replace(path, BelRune, string(os.PathSeparator), -1)

	return path, nil
}
