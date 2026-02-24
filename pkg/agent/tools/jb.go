package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	jbpkg "github.com/jsonnet-bundler/jsonnet-bundler/pkg"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	v1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const vendorDir = "vendor"

// JBTools provides jsonnet-bundler operations by calling the library directly,
// with no subprocess or exec of the jb binary.
type JBTools struct {
	repoRoot string
}

// NewJBTools creates jb tools anchored to the given repository root.
func NewJBTools(repoRoot string) ([]adktool.Tool, error) {
	jt := &JBTools{repoRoot: repoRoot}
	var tools []adktool.Tool
	for _, mk := range []func() (adktool.Tool, error){
		jt.initTool, jt.installTool, jt.updateTool, jt.listTool,
	} {
		t, err := mk()
		if err != nil {
			return nil, err
		}
		tools = append(tools, t)
	}
	return tools, nil
}

// absDir resolves a relative path (or ".") to an absolute path within the repo.
func (jt *JBTools) absDir(relPath string) (string, error) {
	if relPath == "" || relPath == "." {
		return jt.repoRoot, nil
	}
	p := filepath.Join(jt.repoRoot, filepath.Clean(relPath))
	rel, err := filepath.Rel(jt.repoRoot, p)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path %q would escape the repository root", relPath)
	}
	return p, nil
}

// writeJSONFile serialises v to a pretty-printed JSON file.
func writeJSONFile(path string, v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0644)
}

// writeChangedJSONFile only writes the file if its content differs from the original.
func writeChangedJSONFile(originalBytes []byte, modified *v1.JsonnetFile, path string) error {
	orig, err := jsonnetfile.Unmarshal(originalBytes)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(orig, *modified) {
		return nil
	}
	return writeJSONFile(path, *modified)
}

// ensureWithRecovery wraps jbpkg.Ensure and converts any panic into an error.
// The jb library uses panic() instead of returning errors in some code paths.
func ensureWithRecovery(jf v1.JsonnetFile, vendorDir string, locks *deps.Ordered) (result *deps.Ordered, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			retErr = fmt.Errorf("%v", r)
		}
	}()
	return jbpkg.Ensure(jf, vendorDir, locks)
}

// loadLockFile reads the lock file, returning an empty JsonnetFile on ENOENT.
func loadLockFile(dir string) (v1.JsonnetFile, []byte, error) {
	path := filepath.Join(dir, jsonnetfile.LockFile)
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return v1.New(), nil, nil
	}
	if err != nil {
		return v1.New(), nil, fmt.Errorf("reading lockfile: %w", err)
	}
	lf, err := jsonnetfile.Unmarshal(b)
	if err != nil {
		return v1.New(), b, fmt.Errorf("parsing lockfile: %w", err)
	}
	return lf, b, nil
}

func (jt *JBTools) initTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "jb_init",
			Description: "Initialise a new jsonnetfile.json in the given directory. Equivalent to 'jb init'. Fails if the file already exists.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Directory in which to create jsonnetfile.json (relative to repo root, use '.' for root)"},
					"directory": {"type": "string", "description": "Alias for path — either may be used"}
				}
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			if _, ok := input["path"]; !ok {
				input["path"] = input["directory"]
			}
			var params struct {
				Path string `json:"path"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			if params.Path == "" {
				return nil, fmt.Errorf("path (or directory) is required")
			}
			dir, err := jt.absDir(params.Path)
			if err != nil {
				return nil, err
			}
			jfPath := filepath.Join(dir, jsonnetfile.File)
			exists, err := jsonnetfile.Exists(jfPath)
			if err != nil {
				return nil, fmt.Errorf("checking for existing jsonnetfile: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("jsonnetfile.json already exists at %s", params.Path)
			}
			s := v1.New()
			if err := writeJSONFile(jfPath, s); err != nil {
				return nil, fmt.Errorf("writing jsonnetfile.json: %w", err)
			}
			return map[string]any{"output": fmt.Sprintf("Initialised jsonnetfile.json in %s", params.Path)}, nil
		})
}

func (jt *JBTools) installTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "jb_install",
			Description: "Install jsonnet dependencies. Without packages, installs everything listed in jsonnetfile.json (equivalent to 'jb install'). With packages, adds them and vendors them (equivalent to 'jb install github.com/org/repo/path@version'). Updates jsonnetfile.json and jsonnetfile.lock.json.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Directory containing jsonnetfile.json (relative to repo root, use '.' for root)"},
					"packages": {"type": "array", "items": {"type": "string"}, "description": "Package URIs to install, e.g. ['github.com/grafana/jsonnet-libs/ksonnet-util@main']. Omit or leave empty to install from existing jsonnetfile.json."}
				},
				"required": ["path"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Path     string   `json:"path"`
				Packages []string `json:"packages"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			dir, err := jt.absDir(params.Path)
			if err != nil {
				return nil, err
			}
			jfPath := filepath.Join(dir, jsonnetfile.File)
			jfBytes, err := os.ReadFile(jfPath)
			if err != nil {
				return nil, fmt.Errorf("reading jsonnetfile.json (run jb_init first?): %w", err)
			}
			jf, err := jsonnetfile.Unmarshal(jfBytes)
			if err != nil {
				return nil, fmt.Errorf("parsing jsonnetfile.json: %w", err)
			}
			lf, lfBytes, err := loadLockFile(dir)
			if err != nil {
				return nil, err
			}
			if err := os.MkdirAll(filepath.Join(dir, vendorDir, ".tmp"), os.ModePerm); err != nil {
				return nil, fmt.Errorf("creating vendor directory: %w", err)
			}
			for _, u := range params.Packages {
				d := deps.Parse(dir, u)
				if d == nil {
					return nil, fmt.Errorf("unable to parse package URI %q", u)
				}
				existing, _ := jf.Dependencies.Get(d.Name())
				if !reflect.DeepEqual(existing, *d) {
					jf.Dependencies.Set(d.Name(), *d)
					lf.Dependencies.Delete(d.Name())
				}
			}
			locked, err := ensureWithRecovery(jf, filepath.Join(dir, vendorDir), lf.Dependencies)
			if err != nil {
				return nil, fmt.Errorf("installing packages: %w", err)
			}
			jbpkg.CleanLegacyName(jf.Dependencies)
			if err := writeChangedJSONFile(jfBytes, &jf, jfPath); err != nil {
				return nil, fmt.Errorf("updating jsonnetfile.json: %w", err)
			}
			updatedLock := v1.JsonnetFile{Dependencies: locked}
			if err := writeChangedJSONFile(lfBytes, &updatedLock, filepath.Join(dir, jsonnetfile.LockFile)); err != nil {
				return nil, fmt.Errorf("updating jsonnetfile.lock.json: %w", err)
			}
			var sb strings.Builder
			if len(params.Packages) > 0 {
				fmt.Fprintf(&sb, "Installed %d new package(s) into %s/vendor:\n", len(params.Packages), params.Path)
				for _, u := range params.Packages {
					fmt.Fprintf(&sb, "  • %s\n", u)
				}
			} else {
				count := 0
				if locked != nil {
					count = len(locked.Keys())
				}
				fmt.Fprintf(&sb, "Installed %d package(s) from jsonnetfile.json into %s/vendor", count, params.Path)
			}
			return map[string]any{"output": sb.String()}, nil
		})
}

func (jt *JBTools) updateTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "jb_update",
			Description: "Update jsonnet dependencies to their latest versions. Without packages, updates all dependencies (equivalent to 'jb update'). With packages, updates only those listed.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Directory containing jsonnetfile.json (relative to repo root, use '.' for root)"},
					"packages": {"type": "array", "items": {"type": "string"}, "description": "Package URIs to update. Omit or leave empty to update all packages."}
				},
				"required": ["path"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Path     string   `json:"path"`
				Packages []string `json:"packages"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			dir, err := jt.absDir(params.Path)
			if err != nil {
				return nil, err
			}
			jf, err := jsonnetfile.Load(filepath.Join(dir, jsonnetfile.File))
			if err != nil {
				return nil, fmt.Errorf("loading jsonnetfile.json: %w", err)
			}
			lf, _, err := loadLockFile(dir)
			if err != nil {
				return nil, err
			}
			if err := os.MkdirAll(filepath.Join(dir, vendorDir, ".tmp"), os.ModePerm); err != nil {
				return nil, fmt.Errorf("creating vendor directory: %w", err)
			}
			locks := lf.Dependencies
			if len(params.Packages) == 0 {
				locks = deps.NewOrdered()
			} else {
				for _, u := range params.Packages {
					d := deps.Parse(dir, u)
					if d == nil {
						return nil, fmt.Errorf("unable to parse package URI %q", u)
					}
					locks.Delete(d.Name())
				}
			}
			newLocks, err := ensureWithRecovery(jf, filepath.Join(dir, vendorDir), locks)
			if err != nil {
				return nil, fmt.Errorf("updating packages: %w", err)
			}
			updatedLock := v1.JsonnetFile{Dependencies: newLocks}
			if err := writeJSONFile(filepath.Join(dir, jsonnetfile.LockFile), updatedLock); err != nil {
				return nil, fmt.Errorf("writing jsonnetfile.lock.json: %w", err)
			}
			if len(params.Packages) == 0 {
				return map[string]any{"output": fmt.Sprintf("Updated all packages in %s/vendor", params.Path)}, nil
			}
			return map[string]any{"output": fmt.Sprintf("Updated %s in %s/vendor", strings.Join(params.Packages, ", "), params.Path)}, nil
		})
}

func (jt *JBTools) listTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "jb_list",
			Description: "List the jsonnet dependencies declared in jsonnetfile.json, along with their locked versions from jsonnetfile.lock.json.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Directory containing jsonnetfile.json (relative to repo root, use '.' for root)"}
				},
				"required": ["path"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Path string `json:"path"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			dir, err := jt.absDir(params.Path)
			if err != nil {
				return nil, err
			}
			jf, err := jsonnetfile.Load(filepath.Join(dir, jsonnetfile.File))
			if err != nil {
				return nil, fmt.Errorf("loading jsonnetfile.json: %w", err)
			}
			lf, _, _ := loadLockFile(dir)
			keys := jf.Dependencies.Keys()
			if len(keys) == 0 {
				return map[string]any{"output": "No dependencies in jsonnetfile.json"}, nil
			}
			var sb strings.Builder
			fmt.Fprintf(&sb, "%d dependency(s) in %s:\n\n", len(keys), params.Path)
			for _, k := range keys {
				d, _ := jf.Dependencies.Get(k)
				version := d.Version
				if version == "" {
					version = "(unspecified)"
				}
				lockedVersion := "(not locked)"
				if l, ok := lf.Dependencies.Get(k); ok {
					lockedVersion = l.Version
				}
				fmt.Fprintf(&sb, "  %s\n    declared: %s\n    locked:   %s\n", k, version, lockedVersion)
			}
			return map[string]any{"output": sb.String()}, nil
		})
}
