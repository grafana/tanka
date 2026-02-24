package tools

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// FileTools provides sandboxed filesystem operations anchored to the repo root.
type FileTools struct {
	repoRoot string
}

// NewFileTools creates file tools sandboxed to the given repository root.
func NewFileTools(repoRoot string) ([]adktool.Tool, error) {
	ft := &FileTools{repoRoot: repoRoot}
	var tools []adktool.Tool
	for _, mk := range []func() (adktool.Tool, error){
		ft.readTool, ft.writeTool, ft.listTool, ft.searchTool, ft.deleteTool,
	} {
		t, err := mk()
		if err != nil {
			return nil, err
		}
		tools = append(tools, t)
	}
	return tools, nil
}

// safePath validates and resolves a relative path within the repo root.
func (ft *FileTools) safePath(relPath string) (string, error) {
	cleaned := filepath.Clean(relPath)
	if filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("absolute paths not allowed; use paths relative to the repository root")
	}
	full := filepath.Join(ft.repoRoot, cleaned)
	rel, err := filepath.Rel(ft.repoRoot, full)
	if err != nil {
		return "", fmt.Errorf("invalid path %q: %w", relPath, err)
	}
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path %q would escape the repository root", relPath)
	}
	return full, nil
}

func (ft *FileTools) readTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "file_read",
			Description: "Read the contents of a file. Supports pagination via offset/limit (in lines). Use offset+limit to page through large files.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Relative path to the file from the repository root"},
					"offset": {"type": "integer", "description": "Number of lines to skip from the start (default: 0)"},
					"limit": {"type": "integer", "description": "Maximum number of lines to return (default: 500, max: 500)"}
				},
				"required": ["path"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Path   string `json:"path"`
				Offset int    `json:"offset"`
				Limit  int    `json:"limit"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			if params.Limit <= 0 || params.Limit > 500 {
				params.Limit = 500
			}
			fullPath, err := ft.safePath(params.Path)
			if err != nil {
				return nil, err
			}
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("reading %s: %w", params.Path, err)
			}
			lines := strings.Split(string(content), "\n")
			total := len(lines)
			if params.Offset >= total {
				return map[string]any{"output": fmt.Sprintf("[%s: %d lines total, offset %d is past end of file]", params.Path, total, params.Offset)}, nil
			}
			end := params.Offset + params.Limit
			if end > total {
				end = total
			}
			page := strings.Join(lines[params.Offset:end], "\n")
			if total <= params.Limit && params.Offset == 0 {
				return map[string]any{"output": page}, nil
			}
			header := fmt.Sprintf("[%s: lines %d–%d of %d]", params.Path, params.Offset+1, end, total)
			if end < total {
				header += fmt.Sprintf(" (%d more lines, use offset=%d to continue)", total-end, end)
			}
			return map[string]any{"output": header + "\n" + page}, nil
		})
}

func (ft *FileTools) writeTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "file_write",
			Description: "Write content to a file in the repository, creating it if it doesn't exist. Use relative paths from the repository root.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Relative path to the file from the repository root"},
					"content": {"type": "string", "description": "Content to write to the file"}
				},
				"required": ["path", "content"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Path    string `json:"path"`
				Content string `json:"content"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			fullPath, err := ft.safePath(params.Path)
			if err != nil {
				return nil, err
			}
			if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
				return nil, fmt.Errorf("creating parent directories for %s: %w", params.Path, err)
			}
			if err := os.WriteFile(fullPath, []byte(params.Content), 0644); err != nil {
				return nil, fmt.Errorf("writing %s: %w", params.Path, err)
			}
			return map[string]any{"output": fmt.Sprintf("Successfully wrote %d bytes to %s", len(params.Content), params.Path)}, nil
		})
}

func (ft *FileTools) listTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "file_list",
			Description: "List files matching a glob pattern in the repository. Use relative paths. Examples: '**/*.jsonnet', 'environments/**', 'lib/*.libsonnet'. Supports pagination via offset/limit.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"glob_pattern": {"type": "string", "description": "Glob pattern to match files, relative to repository root"},
					"offset": {"type": "integer", "description": "Number of results to skip (default: 0)"},
					"limit": {"type": "integer", "description": "Maximum number of results to return (default: 200)"}
				},
				"required": ["glob_pattern"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				GlobPattern string `json:"glob_pattern"`
				Offset      int    `json:"offset"`
				Limit       int    `json:"limit"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			if params.Limit <= 0 {
				params.Limit = 200
			}
			var matches []string
			err := filepath.WalkDir(ft.repoRoot, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				rel, relErr := filepath.Rel(ft.repoRoot, path)
				if relErr != nil {
					return nil
				}
				matched, matchErr := filepath.Match(params.GlobPattern, rel)
				if matchErr != nil {
					return matchErr
				}
				if !matched {
					matched = matchDoubleGlob(rel, params.GlobPattern)
				}
				if matched && !d.IsDir() {
					matches = append(matches, rel)
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("listing files: %w", err)
			}
			if len(matches) == 0 {
				return map[string]any{"output": fmt.Sprintf("No files found matching %q", params.GlobPattern)}, nil
			}
			total := len(matches)
			if params.Offset >= total {
				return map[string]any{"output": fmt.Sprintf("[%d files matched %q, offset %d is past end]", total, params.GlobPattern, params.Offset)}, nil
			}
			end := params.Offset + params.Limit
			if end > total {
				end = total
			}
			page := strings.Join(matches[params.Offset:end], "\n")
			if total <= params.Limit && params.Offset == 0 {
				return map[string]any{"output": page}, nil
			}
			header := fmt.Sprintf("[%d–%d of %d files matching %q]", params.Offset+1, end, total, params.GlobPattern)
			if end < total {
				header += fmt.Sprintf(" (%d more, use offset=%d to continue)", total-end, end)
			}
			return map[string]any{"output": header + "\n" + page}, nil
		})
}

// matchDoubleGlob handles ** patterns by checking if any suffix of the path
// matches the pattern suffix after the **.
func matchDoubleGlob(path, pattern string) bool {
	if !strings.Contains(pattern, "**") {
		return false
	}
	parts := strings.SplitN(pattern, "**", 2)
	prefix := parts[0]
	suffix := parts[1]
	if suffix != "" && strings.HasPrefix(suffix, "/") {
		suffix = suffix[1:]
	}
	if prefix != "" && !strings.HasPrefix(path, prefix) {
		return false
	}
	if suffix == "" {
		return true
	}
	matched, _ := filepath.Match(suffix, filepath.Base(path))
	return matched
}

func (ft *FileTools) searchTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "file_search",
			Description: "Search for text within files matching a glob pattern. Returns matching lines with file paths and line numbers. Supports pagination via offset/limit.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"glob_pattern": {"type": "string", "description": "Glob pattern to match files to search within"},
					"text_query": {"type": "string", "description": "Text to search for within the matched files"},
					"offset": {"type": "integer", "description": "Number of results to skip (default: 0)"},
					"limit": {"type": "integer", "description": "Maximum number of results to return (default: 200)"}
				},
				"required": ["glob_pattern", "text_query"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				GlobPattern string `json:"glob_pattern"`
				TextQuery   string `json:"text_query"`
				Offset      int    `json:"offset"`
				Limit       int    `json:"limit"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			if params.Limit <= 0 {
				params.Limit = 200
			}
			var matches []string
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
						matches = append(matches, fmt.Sprintf("%s:%d: %s", rel, i+1, line))
					}
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("searching files: %w", err)
			}
			if len(matches) == 0 {
				return map[string]any{"output": fmt.Sprintf("No matches found for %q in files matching %q", params.TextQuery, params.GlobPattern)}, nil
			}
			total := len(matches)
			if params.Offset >= total {
				return map[string]any{"output": fmt.Sprintf("[%d matches for %q, offset %d is past end]\n", total, params.TextQuery, params.Offset)}, nil
			}
			end := params.Offset + params.Limit
			if end > total {
				end = total
			}
			page := strings.Join(matches[params.Offset:end], "\n") + "\n"
			if total <= params.Limit && params.Offset == 0 {
				return map[string]any{"output": page}, nil
			}
			header := fmt.Sprintf("[%d–%d of %d matches for %q]", params.Offset+1, end, total, params.TextQuery)
			if end < total {
				header += fmt.Sprintf(" (%d more, use offset=%d to continue)", total-end, end)
			}
			return map[string]any{"output": header + "\n" + page}, nil
		})
}

func (ft *FileTools) deleteTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "file_delete",
			Description: "Delete a file from the repository. Use relative paths from the repository root.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Relative path to the file to delete from the repository root"}
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
			fullPath, err := ft.safePath(params.Path)
			if err != nil {
				return nil, err
			}
			if err := os.Remove(fullPath); err != nil {
				return nil, fmt.Errorf("deleting %s: %w", params.Path, err)
			}
			return map[string]any{"output": fmt.Sprintf("Successfully deleted %s", params.Path)}, nil
		})
}
