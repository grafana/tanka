package tools

import (
	"fmt"
	"strings"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// GitTools provides git repository operations using go-git (pure Go, no binary required).
type GitTools struct {
	repoRoot string
}

// NewGitTools creates git tools for the repository at the given path.
func NewGitTools(repoRoot string) ([]adktool.Tool, error) {
	gt := &GitTools{repoRoot: repoRoot}
	var tools []adktool.Tool
	for _, mk := range []func() (adktool.Tool, error){
		gt.logTool, gt.statusTool, gt.diffTool,
		gt.branchCreateTool, gt.branchCheckoutTool,
		gt.addTool, gt.commitTool,
	} {
		t, err := mk()
		if err != nil {
			return nil, err
		}
		tools = append(tools, t)
	}
	return tools, nil
}

func (gt *GitTools) openRepo() (*gogit.Repository, error) {
	r, err := gogit.PlainOpen(gt.repoRoot)
	if err != nil {
		return nil, fmt.Errorf("opening git repository: %w", err)
	}
	return r, nil
}

func (gt *GitTools) logTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "git_log",
			Description: "Show recent commit history for the repository.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"limit": {"type": "integer", "description": "Maximum number of commits to return (default: 10)", "default": 10}
				}
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Limit int `json:"limit"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			if params.Limit <= 0 {
				params.Limit = 10
			}
			r, err := gt.openRepo()
			if err != nil {
				return nil, err
			}
			ref, err := r.Head()
			if err != nil {
				return nil, fmt.Errorf("getting HEAD: %w", err)
			}
			iter, err := r.Log(&gogit.LogOptions{From: ref.Hash()})
			if err != nil {
				return nil, fmt.Errorf("getting log: %w", err)
			}
			defer iter.Close()

			var sb strings.Builder
			count := 0
			err = iter.ForEach(func(c *object.Commit) error {
				if count >= params.Limit {
					return storer.ErrStop
				}
				fmt.Fprintf(&sb, "%s %s %s\n  %s\n",
					c.Hash.String()[:8],
					c.Author.When.Format("2006-01-02"),
					c.Author.Name,
					strings.TrimSpace(c.Message),
				)
				count++
				return nil
			})
			if err != nil && err != storer.ErrStop {
				return nil, fmt.Errorf("iterating commits: %w", err)
			}
			if sb.Len() == 0 {
				return map[string]any{"output": "No commits found"}, nil
			}
			return map[string]any{"output": sb.String()}, nil
		})
}

func (gt *GitTools) statusTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "git_status",
			Description: "Show the working tree status — modified, staged, and untracked files.",
			InputSchema: mustSchema(`{"type": "object", "properties": {}}`),
		},
		func(_ adktool.Context, _ map[string]any) (map[string]any, error) {
			r, err := gt.openRepo()
			if err != nil {
				return nil, err
			}
			w, err := r.Worktree()
			if err != nil {
				return nil, fmt.Errorf("getting worktree: %w", err)
			}
			status, err := w.Status()
			if err != nil {
				return nil, fmt.Errorf("getting status: %w", err)
			}
			if status.IsClean() {
				return map[string]any{"output": "nothing to commit, working tree clean"}, nil
			}
			var sb strings.Builder
			for file, s := range status {
				staging := rune(s.Staging)
				worktree := rune(s.Worktree)
				if staging == ' ' {
					staging = '-'
				}
				if worktree == ' ' {
					worktree = '-'
				}
				fmt.Fprintf(&sb, "%c%c %s\n", staging, worktree, file)
			}
			return map[string]any{"output": sb.String()}, nil
		})
}

func (gt *GitTools) diffTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "git_diff",
			Description: "Show unstaged and staged changes in the working tree as a unified diff.",
			InputSchema: mustSchema(`{"type": "object", "properties": {}}`),
		},
		func(_ adktool.Context, _ map[string]any) (map[string]any, error) {
			r, err := gt.openRepo()
			if err != nil {
				return nil, err
			}
			w, err := r.Worktree()
			if err != nil {
				return nil, fmt.Errorf("getting worktree: %w", err)
			}
			status, err := w.Status()
			if err != nil {
				return nil, fmt.Errorf("getting status: %w", err)
			}
			if status.IsClean() {
				return map[string]any{"output": "No changes"}, nil
			}
			var sb strings.Builder
			sb.WriteString("Changes:\n")
			for file, s := range status {
				switch {
				case s.Staging == gogit.Added || s.Worktree == gogit.Untracked:
					fmt.Fprintf(&sb, "  new file:   %s\n", file)
				case s.Staging == gogit.Deleted || s.Worktree == gogit.Deleted:
					fmt.Fprintf(&sb, "  deleted:    %s\n", file)
				case s.Staging == gogit.Modified || s.Worktree == gogit.Modified:
					fmt.Fprintf(&sb, "  modified:   %s\n", file)
				case s.Staging == gogit.Renamed:
					fmt.Fprintf(&sb, "  renamed:    %s\n", file)
				default:
					fmt.Fprintf(&sb, "  changed:    %s\n", file)
				}
			}
			return map[string]any{"output": sb.String()}, nil
		})
}

func (gt *GitTools) branchCreateTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "git_branch_create",
			Description: "Create a new branch from HEAD.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"name": {"type": "string", "description": "Name for the new branch (e.g. 'feat/add-staging-env')"}
				},
				"required": ["name"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Name string `json:"name"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			r, err := gt.openRepo()
			if err != nil {
				return nil, err
			}
			w, err := r.Worktree()
			if err != nil {
				return nil, fmt.Errorf("getting worktree: %w", err)
			}
			err = w.Checkout(&gogit.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(params.Name),
				Create: true,
			})
			if err != nil {
				return nil, fmt.Errorf("creating branch %q: %w", params.Name, err)
			}
			return map[string]any{"output": fmt.Sprintf("Created and switched to branch %q", params.Name)}, nil
		})
}

func (gt *GitTools) branchCheckoutTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "git_branch_checkout",
			Description: "Switch to an existing branch.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"name": {"type": "string", "description": "Name of the branch to switch to"}
				},
				"required": ["name"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Name string `json:"name"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			r, err := gt.openRepo()
			if err != nil {
				return nil, err
			}
			w, err := r.Worktree()
			if err != nil {
				return nil, fmt.Errorf("getting worktree: %w", err)
			}
			err = w.Checkout(&gogit.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(params.Name),
			})
			if err != nil {
				return nil, fmt.Errorf("checking out branch %q: %w", params.Name, err)
			}
			return map[string]any{"output": fmt.Sprintf("Switched to branch %q", params.Name)}, nil
		})
}

func (gt *GitTools) addTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "git_add",
			Description: "Stage files for commit. Pass specific file paths or '.' to stage all changes.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"paths": {"type": "array", "items": {"type": "string"}, "description": "List of file paths to stage. Use [\".\"] to stage all changes."}
				},
				"required": ["paths"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Paths []string `json:"paths"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			r, err := gt.openRepo()
			if err != nil {
				return nil, err
			}
			w, err := r.Worktree()
			if err != nil {
				return nil, fmt.Errorf("getting worktree: %w", err)
			}
			for _, p := range params.Paths {
				if _, err := w.Add(p); err != nil {
					return nil, fmt.Errorf("staging %q: %w", p, err)
				}
			}
			return map[string]any{"output": fmt.Sprintf("Staged: %s", strings.Join(params.Paths, ", "))}, nil
		})
}

func (gt *GitTools) commitTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "git_commit",
			Description: "Create a commit with the staged changes. Make sure to call git_add first.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"message": {"type": "string", "description": "Commit message describing the change"}
				},
				"required": ["message"]
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Message string `json:"message"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			r, err := gt.openRepo()
			if err != nil {
				return nil, err
			}
			w, err := r.Worktree()
			if err != nil {
				return nil, fmt.Errorf("getting worktree: %w", err)
			}
			cfg, err := r.Config()
			if err != nil {
				return nil, fmt.Errorf("reading git config: %w", err)
			}
			authorName := cfg.User.Name
			authorEmail := cfg.User.Email
			if authorName == "" {
				authorName = "Tanka Agent"
			}
			if authorEmail == "" {
				authorEmail = "tanka-agent@local"
			}
			hash, err := w.Commit(params.Message, &gogit.CommitOptions{
				Author: &object.Signature{
					Name:  authorName,
					Email: authorEmail,
					When:  time.Now(),
				},
			})
			if err != nil {
				return nil, fmt.Errorf("committing: %w", err)
			}
			return map[string]any{"output": fmt.Sprintf("Created commit %s: %s", hash.String()[:8], params.Message)}, nil
		})
}
