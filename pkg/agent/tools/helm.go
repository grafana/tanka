package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// HelmTools provides helm CLI operations for managing chart dependencies.
type HelmTools struct {
	repoRoot string
}

// NewHelmTools creates helm tools anchored to the given repository root.
func NewHelmTools(repoRoot string) ([]adktool.Tool, error) {
	ht := &HelmTools{repoRoot: repoRoot}
	var tools []adktool.Tool
	for _, mk := range []func() (adktool.Tool, error){
		ht.dependencyBuildTool, ht.dependencyUpdateTool, ht.dependencyListTool, ht.lintTool,
	} {
		t, err := mk()
		if err != nil {
			return nil, err
		}
		tools = append(tools, t)
	}
	return tools, nil
}

// helmBin returns the path to the helm binary, respecting TANKA_HELM_PATH.
func helmBin() string {
	if env := os.Getenv("TANKA_HELM_PATH"); env != "" {
		return env
	}
	return "helm"
}

// runHelm runs helm with the given arguments inside dir and returns combined output.
func (ht *HelmTools) runHelm(dir string, args ...string) (string, error) {
	cmd := exec.Command(helmBin(), args...)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("helm %s: %w\n%s", args[0], err, strings.TrimSpace(out.String()))
	}
	return strings.TrimSpace(out.String()), nil
}

// absChartDir resolves a chart path within the repo root.
func (ht *HelmTools) absChartDir(relPath string) (string, error) {
	if relPath == "" || relPath == "." {
		return ht.repoRoot, nil
	}
	p := filepath.Join(ht.repoRoot, filepath.Clean(relPath))
	rel, err := filepath.Rel(ht.repoRoot, p)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path %q would escape the repository root", relPath)
	}
	return p, nil
}

func (ht *HelmTools) dependencyBuildTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "helm_dependency_build",
			Description: "Run 'helm dependency build' for a chart directory. Downloads chart dependencies declared in Chart.yaml into the chart's charts/ subdirectory and regenerates Chart.lock.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Path to the chart directory, relative to the repository root"}
				},
				"required": ["path"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Path string `json:"path" aliases:"chart_dir"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			dir, err := ht.absChartDir(params.Path)
			if err != nil {
				return nil, err
			}
			out, err := ht.runHelm(dir, "dependency", "build", ".")
			if err != nil {
				return nil, err
			}
			return map[string]any{"output": out}, nil
		})
}

func (ht *HelmTools) dependencyUpdateTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "helm_dependency_update",
			Description: "Run 'helm dependency update' for a chart directory. Fetches the latest versions of chart dependencies and updates Chart.lock.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Path to the chart directory, relative to the repository root"}
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
			dir, err := ht.absChartDir(params.Path)
			if err != nil {
				return nil, err
			}
			out, err := ht.runHelm(dir, "dependency", "update", ".")
			if err != nil {
				return nil, err
			}
			return map[string]any{"output": out}, nil
		})
}

func (ht *HelmTools) dependencyListTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "helm_dependency_list",
			Description: "Run 'helm dependency list' for a chart directory. Shows the chart's declared dependencies and their status.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Path to the chart directory, relative to the repository root"}
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
			dir, err := ht.absChartDir(params.Path)
			if err != nil {
				return nil, err
			}
			out, err := ht.runHelm(dir, "dependency", "list", ".")
			if err != nil {
				return nil, err
			}
			return map[string]any{"output": out}, nil
		})
}

func (ht *HelmTools) lintTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "helm_lint",
			Description: "Run 'helm lint' against a chart directory to check for errors and best-practice violations.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"path": {"type": "string", "description": "Path to the chart directory, relative to the repository root"}
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
			dir, err := ht.absChartDir(params.Path)
			if err != nil {
				return nil, err
			}
			out, err := ht.runHelm(ht.repoRoot, "lint", dir)
			if err != nil {
				return nil, err
			}
			return map[string]any{"output": out}, nil
		})
}
