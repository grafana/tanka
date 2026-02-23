package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FileTools provides sandboxed filesystem operations anchored to the repo root.
type FileTools struct {
	repoRoot string
}

// NewFileTools creates file tools sandboxed to the given repository root.
func NewFileTools(repoRoot string) []Tool {
	ft := &FileTools{repoRoot: repoRoot}
	return []Tool{
		ft.readTool(),
		ft.writeTool(),
		ft.listTool(),
		ft.searchTool(),
		ft.deleteTool(),
	}
}

// safePath validates and resolves a relative path within the repo root.
// Returns an error if the path would escape the sandbox.
func (ft *FileTools) safePath(relPath string) (string, error) {
	// Clean the path to remove . and .. components
	cleaned := filepath.Clean(relPath)
	// Reject absolute paths
	if filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("absolute paths not allowed; use paths relative to the repository root")
	}
	// Build full path
	full := filepath.Join(ft.repoRoot, cleaned)
	// Verify the result is still within the repo root
	rel, err := filepath.Rel(ft.repoRoot, full)
	if err != nil {
		return "", fmt.Errorf("invalid path %q: %w", relPath, err)
	}
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path %q would escape the repository root", relPath)
	}
	return full, nil
}

func (ft *FileTools) readTool() Tool {
	return Tool{
		Name:        "file_read",
		Description: "Read the contents of a file in the repository. Use relative paths from the repository root.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "Relative path to the file from the repository root"
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
			fullPath, err := ft.safePath(params.Path)
			if err != nil {
				return "", err
			}
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return "", fmt.Errorf("reading %s: %w", params.Path, err)
			}
			return string(content), nil
		},
	}
}

func (ft *FileTools) writeTool() Tool {
	return Tool{
		Name:        "file_write",
		Description: "Write content to a file in the repository, creating it if it doesn't exist. Use relative paths from the repository root.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "Relative path to the file from the repository root"
				},
				"content": {
					"type": "string",
					"description": "Content to write to the file"
				}
			},
			"required": ["path", "content"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				Path    string `json:"path"`
				Content string `json:"content"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}
			fullPath, err := ft.safePath(params.Path)
			if err != nil {
				return "", err
			}
			// Create parent directories if needed
			if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
				return "", fmt.Errorf("creating parent directories for %s: %w", params.Path, err)
			}
			if err := os.WriteFile(fullPath, []byte(params.Content), 0644); err != nil {
				return "", fmt.Errorf("writing %s: %w", params.Path, err)
			}
			return fmt.Sprintf("Successfully wrote %d bytes to %s", len(params.Content), params.Path), nil
		},
	}
}

func (ft *FileTools) listTool() Tool {
	return Tool{
		Name:        "file_list",
		Description: "List files matching a glob pattern in the repository. Use relative paths. Examples: '**/*.jsonnet', 'environments/**', 'lib/*.libsonnet'",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"glob_pattern": {
					"type": "string",
					"description": "Glob pattern to match files, relative to repository root"
				}
			},
			"required": ["glob_pattern"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				GlobPattern string `json:"glob_pattern"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}
			// Walk from repo root matching pattern
			var matches []string
			err := filepath.WalkDir(ft.repoRoot, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil // skip errors
				}
				rel, relErr := filepath.Rel(ft.repoRoot, path)
				if relErr != nil {
					return nil
				}
				matched, matchErr := filepath.Match(params.GlobPattern, rel)
				if matchErr != nil {
					return matchErr
				}
				// Also try double-star matching by checking path components
				if !matched {
					matched = matchDoubleGlob(rel, params.GlobPattern)
				}
				if matched && !d.IsDir() {
					matches = append(matches, rel)
				}
				return nil
			})
			if err != nil {
				return "", fmt.Errorf("listing files: %w", err)
			}
			if len(matches) == 0 {
				return fmt.Sprintf("No files found matching %q", params.GlobPattern), nil
			}
			return strings.Join(matches, "\n"), nil
		},
	}
}

// matchDoubleGlob handles ** patterns by checking if any suffix of the path
// matches the pattern suffix after the **.
func matchDoubleGlob(path, pattern string) bool {
	// Simple ** handling: replace ** with a check against path suffix
	if !strings.Contains(pattern, "**") {
		return false
	}
	parts := strings.SplitN(pattern, "**", 2)
	prefix := parts[0]
	suffix := parts[1]
	if suffix != "" && strings.HasPrefix(suffix, "/") {
		suffix = suffix[1:]
	}
	// Check prefix
	if prefix != "" && !strings.HasPrefix(path, prefix) {
		return false
	}
	// Check suffix against path or any path component
	if suffix == "" {
		return true
	}
	matched, _ := filepath.Match(suffix, filepath.Base(path))
	return matched
}

func (ft *FileTools) searchTool() Tool {
	return Tool{
		Name:        "file_search",
		Description: "Search for text within files matching a glob pattern. Returns matching lines with file paths and line numbers.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"glob_pattern": {
					"type": "string",
					"description": "Glob pattern to match files to search within"
				},
				"text_query": {
					"type": "string",
					"description": "Text to search for within the matched files"
				}
			},
			"required": ["glob_pattern", "text_query"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				GlobPattern string `json:"glob_pattern"`
				TextQuery   string `json:"text_query"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}

			var results bytes.Buffer
			matchCount := 0

			err := filepath.WalkDir(ft.repoRoot, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				rel, relErr := filepath.Rel(ft.repoRoot, path)
				if relErr != nil {
					return nil
				}
				matched, _ := filepath.Match(params.GlobPattern, rel)
				if !matched {
					matched = matchDoubleGlob(rel, params.GlobPattern)
				}
				if !matched {
					return nil
				}

				content, readErr := os.ReadFile(path)
				if readErr != nil {
					return nil
				}
				lines := strings.Split(string(content), "\n")
				for i, line := range lines {
					if strings.Contains(line, params.TextQuery) {
						fmt.Fprintf(&results, "%s:%d: %s\n", rel, i+1, line)
						matchCount++
					}
				}
				return nil
			})
			if err != nil {
				return "", fmt.Errorf("searching files: %w", err)
			}
			if matchCount == 0 {
				return fmt.Sprintf("No matches found for %q in files matching %q", params.TextQuery, params.GlobPattern), nil
			}
			return results.String(), nil
		},
	}
}

func (ft *FileTools) deleteTool() Tool {
	return Tool{
		Name:        "file_delete",
		Description: "Delete a file from the repository. Use relative paths from the repository root.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "Relative path to the file to delete from the repository root"
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
			fullPath, err := ft.safePath(params.Path)
			if err != nil {
				return "", err
			}
			if err := os.Remove(fullPath); err != nil {
				return "", fmt.Errorf("deleting %s: %w", params.Path, err)
			}
			return fmt.Sprintf("Successfully deleted %s", params.Path), nil
		},
	}
}
