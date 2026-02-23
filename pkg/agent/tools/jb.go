package tools

import (
	"context"
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
)

const vendorDir = "vendor"

// JBTools provides jsonnet-bundler operations by calling the library directly,
// with no subprocess or exec of the jb binary.
type JBTools struct {
	repoRoot string
}

// NewJBTools creates jb tools anchored to the given repository root.
func NewJBTools(repoRoot string) []Tool {
	jt := &JBTools{repoRoot: repoRoot}
	return []Tool{
		jt.initTool(),
		jt.installTool(),
		jt.updateTool(),
		jt.listTool(),
	}
}

// absDir resolves a relative path (or ".") to an absolute path within the repo.
func (jt *JBTools) absDir(relPath string) (string, error) {
	if relPath == "" || relPath == "." {
		return jt.repoRoot, nil
	}
	p := filepath.Join(jt.repoRoot, filepath.Clean(relPath))
	// Ensure it's still inside the repo root
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

// writeChangedJSONFile only writes the file if its content differs from the
// original bytes (same behaviour as jb's writeChangedJsonnetFile).
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

func (jt *JBTools) initTool() Tool {
	return Tool{
		Name:        "jb_init",
		Description: "Initialise a new jsonnetfile.json in the given directory. Equivalent to 'jb init'. Fails if the file already exists.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "Directory in which to create jsonnetfile.json (relative to repo root, use '.' for root)"
				}
			},
			"required": ["path"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				Path string `json:"path"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}
			dir, err := jt.absDir(params.Path)
			if err != nil {
				return "", err
			}

			jfPath := filepath.Join(dir, jsonnetfile.File)
			exists, err := jsonnetfile.Exists(jfPath)
			if err != nil {
				return "", fmt.Errorf("checking for existing jsonnetfile: %w", err)
			}
			if exists {
				return "", fmt.Errorf("jsonnetfile.json already exists at %s", params.Path)
			}

			s := v1.New()
			if err := writeJSONFile(jfPath, s); err != nil {
				return "", fmt.Errorf("writing jsonnetfile.json: %w", err)
			}
			return fmt.Sprintf("Initialised jsonnetfile.json in %s", params.Path), nil
		},
	}
}

func (jt *JBTools) installTool() Tool {
	return Tool{
		Name:        "jb_install",
		Description: "Install jsonnet dependencies. Without packages, installs everything listed in jsonnetfile.json (equivalent to 'jb install'). With packages, adds them and vendors them (equivalent to 'jb install github.com/org/repo/path@version'). Updates jsonnetfile.json and jsonnetfile.lock.json.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "Directory containing jsonnetfile.json (relative to repo root, use '.' for root)"
				},
				"packages": {
					"type": "array",
					"items": {"type": "string"},
					"description": "Package URIs to install, e.g. ['github.com/grafana/jsonnet-libs/ksonnet-util@main']. Omit or leave empty to install from existing jsonnetfile.json."
				}
			},
			"required": ["path"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				Path     string   `json:"path"`
				Packages []string `json:"packages"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}
			dir, err := jt.absDir(params.Path)
			if err != nil {
				return "", err
			}

			// Read jsonnetfile.json
			jfPath := filepath.Join(dir, jsonnetfile.File)
			jfBytes, err := os.ReadFile(jfPath)
			if err != nil {
				return "", fmt.Errorf("reading jsonnetfile.json (run jb_init first?): %w", err)
			}
			jf, err := jsonnetfile.Unmarshal(jfBytes)
			if err != nil {
				return "", fmt.Errorf("parsing jsonnetfile.json: %w", err)
			}

			// Read lock file (missing is OK)
			lf, lfBytes, err := loadLockFile(dir)
			if err != nil {
				return "", err
			}

			// Ensure vendor/.tmp exists
			if err := os.MkdirAll(filepath.Join(dir, vendorDir, ".tmp"), os.ModePerm); err != nil {
				return "", fmt.Errorf("creating vendor directory: %w", err)
			}

			// Add any newly specified packages to the manifest
			for _, u := range params.Packages {
				d := deps.Parse(dir, u)
				if d == nil {
					return "", fmt.Errorf("unable to parse package URI %q", u)
				}
				existing, _ := jf.Dependencies.Get(d.Name())
				if !reflect.DeepEqual(existing, *d) {
					jf.Dependencies.Set(d.Name(), *d)
					// Force re-download by removing any existing lock entry
					lf.Dependencies.Delete(d.Name())
				}
			}

			// Run the actual installation.
			// ensureWithRecovery converts panics from the jb library into errors.
			locked, err := ensureWithRecovery(jf, filepath.Join(dir, vendorDir), lf.Dependencies)
			if err != nil {
				return "", fmt.Errorf("installing packages: %w", err)
			}
			jbpkg.CleanLegacyName(jf.Dependencies)

			// Write back updated files
			if err := writeChangedJSONFile(jfBytes, &jf, jfPath); err != nil {
				return "", fmt.Errorf("updating jsonnetfile.json: %w", err)
			}
			updatedLock := v1.JsonnetFile{Dependencies: locked}
			if err := writeChangedJSONFile(lfBytes, &updatedLock, filepath.Join(dir, jsonnetfile.LockFile)); err != nil {
				return "", fmt.Errorf("updating jsonnetfile.lock.json: %w", err)
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
			return sb.String(), nil
		},
	}
}

func (jt *JBTools) updateTool() Tool {
	return Tool{
		Name:        "jb_update",
		Description: "Update jsonnet dependencies to their latest versions. Without packages, updates all dependencies (equivalent to 'jb update'). With packages, updates only those listed.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "Directory containing jsonnetfile.json (relative to repo root, use '.' for root)"
				},
				"packages": {
					"type": "array",
					"items": {"type": "string"},
					"description": "Package URIs to update. Omit or leave empty to update all packages."
				}
			},
			"required": ["path"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				Path     string   `json:"path"`
				Packages []string `json:"packages"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}
			dir, err := jt.absDir(params.Path)
			if err != nil {
				return "", err
			}

			jf, err := jsonnetfile.Load(filepath.Join(dir, jsonnetfile.File))
			if err != nil {
				return "", fmt.Errorf("loading jsonnetfile.json: %w", err)
			}
			lf, _, err := loadLockFile(dir)
			if err != nil {
				return "", err
			}

			if err := os.MkdirAll(filepath.Join(dir, vendorDir, ".tmp"), os.ModePerm); err != nil {
				return "", fmt.Errorf("creating vendor directory: %w", err)
			}

			locks := lf.Dependencies
			if len(params.Packages) == 0 {
				// Update all: clear the lock
				locks = deps.NewOrdered()
			} else {
				for _, u := range params.Packages {
					d := deps.Parse(dir, u)
					if d == nil {
						return "", fmt.Errorf("unable to parse package URI %q", u)
					}
					locks.Delete(d.Name())
				}
			}

			newLocks, err := ensureWithRecovery(jf, filepath.Join(dir, vendorDir), locks)
			if err != nil {
				return "", fmt.Errorf("updating packages: %w", err)
			}

			updatedLock := v1.JsonnetFile{Dependencies: newLocks}
			if err := writeJSONFile(filepath.Join(dir, jsonnetfile.LockFile), updatedLock); err != nil {
				return "", fmt.Errorf("writing jsonnetfile.lock.json: %w", err)
			}

			if len(params.Packages) == 0 {
				return fmt.Sprintf("Updated all packages in %s/vendor", params.Path), nil
			}
			return fmt.Sprintf("Updated %s in %s/vendor", strings.Join(params.Packages, ", "), params.Path), nil
		},
	}
}

func (jt *JBTools) listTool() Tool {
	return Tool{
		Name:        "jb_list",
		Description: "List the jsonnet dependencies declared in jsonnetfile.json, along with their locked versions from jsonnetfile.lock.json.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "Directory containing jsonnetfile.json (relative to repo root, use '.' for root)"
				}
			},
			"required": ["path"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				Path string `json:"path"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}
			dir, err := jt.absDir(params.Path)
			if err != nil {
				return "", err
			}

			jf, err := jsonnetfile.Load(filepath.Join(dir, jsonnetfile.File))
			if err != nil {
				return "", fmt.Errorf("loading jsonnetfile.json: %w", err)
			}

			lf, _, _ := loadLockFile(dir) // lock is optional

			keys := jf.Dependencies.Keys()
			if len(keys) == 0 {
				return "No dependencies in jsonnetfile.json", nil
			}

			var sb strings.Builder
			fmt.Fprintf(&sb, "%d dependency(s) in %s:\n\n", len(keys), params.Path)
			for _, k := range keys {
				d, _ := jf.Dependencies.Get(k)
				version := d.Version
				if version == "" {
					version = "(unspecified)"
				}

				// Check the locked version
				lockedVersion := "(not locked)"
				if l, ok := lf.Dependencies.Get(k); ok {
					lockedVersion = l.Version
				}

				fmt.Fprintf(&sb, "  %s\n    declared: %s\n    locked:   %s\n", k, version, lockedVersion)
			}
			return sb.String(), nil
		},
	}
}
