package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/tanka"
	"sigs.k8s.io/yaml"
)

// TankaTools provides Tanka and Jsonnet operations using the in-process Go API.
type TankaTools struct {
	repoRoot string
}

// NewTankaTools creates Tanka tools for the repository at the given path.
func NewTankaTools(repoRoot string) []Tool {
	tt := &TankaTools{repoRoot: repoRoot}
	return []Tool{
		tt.findEnvironmentsTool(),
		tt.showTool(),
		tt.diffTool(),
		tt.validateJsonnetTool(),
		tt.formatJsonnetTool(),
	}
}

func (tt *TankaTools) findEnvironmentsTool() Tool {
	return Tool{
		Name:        "tanka_find_environments",
		Description: "Find all Tanka environments recursively within a path. Returns a list of environments with their metadata including name, API server, and namespace.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"path": {
					"type": "string",
					"description": "Relative path to search from (use '.' for the entire repository)"
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

			searchPath := tt.repoRoot
			if params.Path != "" && params.Path != "." {
				searchPath = strings.TrimSuffix(tt.repoRoot, "/") + "/" + strings.TrimPrefix(params.Path, "/")
			}

			envs, err := tanka.FindEnvs(ctx, searchPath, tanka.FindOpts{})
			if err != nil {
				return "", fmt.Errorf("finding environments: %w", err)
			}
			if len(envs) == 0 {
				return "No Tanka environments found", nil
			}

			var sb strings.Builder
			fmt.Fprintf(&sb, "Found %d environment(s):\n\n", len(envs))
			for _, env := range envs {
				fmt.Fprintf(&sb, "Name:       %s\n", env.Metadata.Name)
				fmt.Fprintf(&sb, "Namespace:  %s\n", env.Spec.Namespace)
				fmt.Fprintf(&sb, "API Server: %s\n", env.Spec.APIServer)
				if len(env.Metadata.Labels) > 0 {
					fmt.Fprintf(&sb, "Labels:     ")
					for k, v := range env.Metadata.Labels {
						fmt.Fprintf(&sb, "%s=%s ", k, v)
					}
					sb.WriteString("\n")
				}
				sb.WriteString("\n")
			}
			return sb.String(), nil
		},
	}
}

func (tt *TankaTools) showTool() Tool {
	return Tool{
		Name:        "tanka_show",
		Description: "Render and display the Kubernetes manifests for a Tanka environment without connecting to a cluster. Shows what would be applied.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"env_path": {
					"type": "string",
					"description": "Relative path to the Tanka environment directory (the one containing spec.json or main.jsonnet)"
				}
			},
			"required": ["env_path"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				EnvPath string `json:"env_path"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}

			envPath := strings.TrimSuffix(tt.repoRoot, "/") + "/" + strings.TrimPrefix(params.EnvPath, "/")

			manifests, err := tanka.Show(ctx, envPath, tanka.Opts{})
			if err != nil {
				return "", fmt.Errorf("rendering environment %s: %w", params.EnvPath, err)
			}

			if len(manifests) == 0 {
				return "No resources rendered", nil
			}

			var sb strings.Builder
			fmt.Fprintf(&sb, "Rendered %d resource(s) for environment %s:\n\n", len(manifests), params.EnvPath)
			for _, m := range manifests {
				out, err := yaml.Marshal(map[string]interface{}(m))
				if err != nil {
					continue
				}
				sb.WriteString("---\n")
				sb.Write(out)
			}
			return sb.String(), nil
		},
	}
}

func (tt *TankaTools) diffTool() Tool {
	return Tool{
		Name:        "tanka_diff",
		Description: "Show the diff between the local Tanka environment configuration and the live cluster state. Requires cluster connectivity (kubectl configured).",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"env_path": {
					"type": "string",
					"description": "Relative path to the Tanka environment directory"
				}
			},
			"required": ["env_path"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				EnvPath string `json:"env_path"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}

			envPath := strings.TrimSuffix(tt.repoRoot, "/") + "/" + strings.TrimPrefix(params.EnvPath, "/")

			diff, err := tanka.Diff(ctx, envPath, tanka.DiffOpts{})
			if err != nil {
				return "", fmt.Errorf("diffing environment %s (note: cluster connectivity is required): %w", params.EnvPath, err)
			}
			if diff == nil || *diff == "" {
				return fmt.Sprintf("No differences found for environment %s — cluster matches local config", params.EnvPath), nil
			}
			return *diff, nil
		},
	}
}

func (tt *TankaTools) validateJsonnetTool() Tool {
	return Tool{
		Name:        "tanka_validate_jsonnet",
		Description: "Validate and lint a Jsonnet or Libsonnet file, reporting any syntax errors or lint warnings.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"file_path": {
					"type": "string",
					"description": "Relative path to the .jsonnet or .libsonnet file to validate"
				}
			},
			"required": ["file_path"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				FilePath string `json:"file_path"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}

			fullPath := strings.TrimSuffix(tt.repoRoot, "/") + "/" + strings.TrimPrefix(params.FilePath, "/")

			// Verify the file exists
			if _, err := os.Stat(fullPath); err != nil {
				return "", fmt.Errorf("file not found: %s", params.FilePath)
			}

			var lintOut bytes.Buffer
			err := jsonnet.Lint([]string{fullPath}, &jsonnet.LintOpts{
				Parallelism: 1,
				Out:         &lintOut,
			})
			if err != nil {
				output := lintOut.String()
				if output != "" {
					return "", fmt.Errorf("validation failed for %s:\n%s", params.FilePath, output)
				}
				return "", fmt.Errorf("validation failed for %s: %w", params.FilePath, err)
			}

			return fmt.Sprintf("✓ %s is valid (no errors or warnings)", params.FilePath), nil
		},
	}
}

func (tt *TankaTools) formatJsonnetTool() Tool {
	return Tool{
		Name:        "tanka_format_jsonnet",
		Description: "Format a Jsonnet or Libsonnet file according to standard Jsonnet formatting rules. Returns the formatted content without modifying the file.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"file_path": {
					"type": "string",
					"description": "Relative path to the .jsonnet or .libsonnet file to format"
				}
			},
			"required": ["file_path"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				FilePath string `json:"file_path"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}

			fullPath := strings.TrimSuffix(tt.repoRoot, "/") + "/" + strings.TrimPrefix(params.FilePath, "/")

			content, err := os.ReadFile(fullPath)
			if err != nil {
				return "", fmt.Errorf("reading %s: %w", params.FilePath, err)
			}

			formatted, err := tanka.Format(fullPath, string(content))
			if err != nil {
				return "", fmt.Errorf("formatting %s: %w", params.FilePath, err)
			}

			if string(content) == formatted {
				return fmt.Sprintf("File %s is already properly formatted", params.FilePath), nil
			}
			return fmt.Sprintf("Formatted content for %s:\n\n%s", params.FilePath, formatted), nil
		},
	}
}
