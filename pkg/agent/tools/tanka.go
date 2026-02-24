package tools

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/tanka"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"sigs.k8s.io/yaml"
)

// TankaTools provides Tanka and Jsonnet operations using the in-process Go API.
type TankaTools struct {
	repoRoot string
}

// NewTankaTools creates Tanka tools for the repository at the given path.
func NewTankaTools(repoRoot string) ([]adktool.Tool, error) {
	tt := &TankaTools{repoRoot: repoRoot}
	var tools []adktool.Tool
	for _, mk := range []func() (adktool.Tool, error){
		tt.findEnvironmentsTool, tt.showTool, tt.diffTool,
		tt.validateJsonnetTool, tt.formatJsonnetTool,
	} {
		t, err := mk()
		if err != nil {
			return nil, err
		}
		tools = append(tools, t)
	}
	return tools, nil
}

func (tt *TankaTools) findEnvironmentsTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "tanka_find_environments",
			Description: "Find all Tanka environments recursively within a path. Returns a list of environments with their metadata including name, API server, and namespace.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Relative path to search from (default: '.' — the entire repository)"}
				}
			}`),
		},
		func(ctx adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Path string `json:"path"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			if params.Path == "" {
				params.Path = "."
			}
			searchPath := tt.repoRoot
			if params.Path != "." {
				searchPath = strings.TrimSuffix(tt.repoRoot, "/") + "/" + strings.TrimPrefix(params.Path, "/")
			}
			envs, err := tanka.FindEnvs(ctx, searchPath, tanka.FindOpts{})
			if err != nil {
				return nil, fmt.Errorf("finding environments: %w", err)
			}
			if len(envs) == 0 {
				return map[string]any{"output": "No Tanka environments found"}, nil
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
			return map[string]any{"output": sb.String()}, nil
		})
}

func (tt *TankaTools) showTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "tanka_show",
			Description: "Render and display the Kubernetes manifests for a Tanka environment without connecting to a cluster. Shows what would be applied.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"env_path": {"type": "string", "description": "Relative path to the Tanka environment directory (the one containing spec.json or main.jsonnet)"},
					"path": {"type": "string", "description": "Alias for env_path"},
					"env": {"type": "string", "description": "Alias for env_path"}
				},
				"required": ["env_path"]
			}`),
		},
		func(ctx adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				EnvPath string `json:"env_path" aliases:"path,env"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			envPath := strings.TrimSuffix(tt.repoRoot, "/") + "/" + strings.TrimPrefix(params.EnvPath, "/")
			manifests, err := tanka.Show(ctx, envPath, tanka.Opts{})
			if err != nil {
				return nil, fmt.Errorf("rendering environment %s: %w", params.EnvPath, err)
			}
			if len(manifests) == 0 {
				return map[string]any{"output": "No resources rendered"}, nil
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
			return map[string]any{"output": sb.String()}, nil
		})
}

func (tt *TankaTools) diffTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "tanka_diff",
			Description: "Show the diff between the local Tanka environment configuration and the live cluster state. Requires cluster connectivity (kubectl configured).",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"env_path": {"type": "string", "description": "Relative path to the Tanka environment directory"}
				},
				"required": ["env_path"]
			}`),
		},
		func(ctx adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				EnvPath string `json:"env_path"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			envPath := strings.TrimSuffix(tt.repoRoot, "/") + "/" + strings.TrimPrefix(params.EnvPath, "/")
			diff, err := tanka.Diff(ctx, envPath, tanka.DiffOpts{})
			if err != nil {
				return nil, fmt.Errorf("diffing environment %s (note: cluster connectivity is required): %w", params.EnvPath, err)
			}
			if diff == nil || *diff == "" {
				return map[string]any{"output": fmt.Sprintf("No differences found for environment %s — cluster matches local config", params.EnvPath)}, nil
			}
			return map[string]any{"output": *diff}, nil
		})
}

func (tt *TankaTools) validateJsonnetTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "tanka_validate_jsonnet",
			Description: "Validate and lint a Jsonnet or Libsonnet file, reporting any syntax errors or lint warnings.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"file_path": {"type": "string", "description": "Relative path to the .jsonnet or .libsonnet file to validate"}
				},
				"required": ["file_path"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				FilePath string `json:"file_path"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			fullPath := strings.TrimSuffix(tt.repoRoot, "/") + "/" + strings.TrimPrefix(params.FilePath, "/")
			if _, err := os.Stat(fullPath); err != nil {
				return nil, fmt.Errorf("file not found: %s", params.FilePath)
			}
			var lintOut bytes.Buffer
			err := jsonnet.Lint([]string{fullPath}, &jsonnet.LintOpts{
				Parallelism: 1,
				Out:         &lintOut,
			})
			if err != nil {
				output := lintOut.String()
				if output != "" {
					return nil, fmt.Errorf("validation failed for %s:\n%s", params.FilePath, output)
				}
				return nil, fmt.Errorf("validation failed for %s: %w", params.FilePath, err)
			}
			return map[string]any{"output": fmt.Sprintf("✓ %s is valid (no errors or warnings)", params.FilePath)}, nil
		})
}

func (tt *TankaTools) formatJsonnetTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "tanka_format_jsonnet",
			Description: "Format a Jsonnet or Libsonnet file according to standard Jsonnet formatting rules. Returns the formatted content without modifying the file.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"file_path": {"type": "string", "description": "Relative path to the .jsonnet or .libsonnet file to format"}
				},
				"required": ["file_path"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				FilePath string `json:"file_path"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			fullPath := strings.TrimSuffix(tt.repoRoot, "/") + "/" + strings.TrimPrefix(params.FilePath, "/")
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("reading %s: %w", params.FilePath, err)
			}
			formatted, err := tanka.Format(fullPath, string(content))
			if err != nil {
				return nil, fmt.Errorf("formatting %s: %w", params.FilePath, err)
			}
			if string(content) == formatted {
				return map[string]any{"output": fmt.Sprintf("File %s is already properly formatted", params.FilePath)}, nil
			}
			return map[string]any{"output": fmt.Sprintf("Formatted content for %s:\n\n%s", params.FilePath, formatted)}, nil
		})
}
