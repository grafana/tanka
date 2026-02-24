package tools

import (
	"fmt"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// GitTools provides read-only git repository operations using go-git (pure Go, no binary required).
type GitTools struct {
	repoRoot string
}

// NewGitTools creates git tools for the repository at the given path.
func NewGitTools(repoRoot string) ([]adktool.Tool, error) {
	gt := &GitTools{repoRoot: repoRoot}
	var tools []adktool.Tool
	for _, mk := range []func() (adktool.Tool, error){
		gt.logTool, gt.showTool,
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

func (gt *GitTools) showTool() (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "git_show",
			Description: "Show the full details of a commit: author, date, message, and a unified diff of all changed files.",
			InputSchema: mustSchema(`{
				"type": "object",
				"properties": {
					"ref": {"type": "string", "description": "Commit ref: full or short hash, HEAD, HEAD~1, branch name, etc. Defaults to HEAD."}
				}
			}`),
		},
		func(_ adktool.Context, input map[string]any) (map[string]any, error) {
			var params struct {
				Ref string `json:"ref"`
			}
			if err := bind(input, &params); err != nil {
				return nil, err
			}
			if params.Ref == "" {
				params.Ref = "HEAD"
			}

			r, err := gt.openRepo()
			if err != nil {
				return nil, err
			}

			hash, err := r.ResolveRevision(plumbing.Revision(params.Ref))
			if err != nil {
				return nil, fmt.Errorf("resolving %q: %w", params.Ref, err)
			}

			commit, err := r.CommitObject(*hash)
			if err != nil {
				return nil, fmt.Errorf("getting commit: %w", err)
			}

			var sb strings.Builder
			fmt.Fprintf(&sb, "commit %s\nAuthor: %s <%s>\nDate:   %s\n\n    %s\n\n",
				commit.Hash.String(),
				commit.Author.Name,
				commit.Author.Email,
				commit.Author.When.Format("Mon Jan 2 15:04:05 2006 -0700"),
				strings.ReplaceAll(strings.TrimSpace(commit.Message), "\n", "\n    "),
			)

			if commit.NumParents() == 0 {
				// Initial commit — list files added
				tree, err := commit.Tree()
				if err != nil {
					return nil, fmt.Errorf("getting tree: %w", err)
				}
				err = tree.Files().ForEach(func(f *object.File) error {
					fmt.Fprintf(&sb, "new file: %s\n", f.Name)
					return nil
				})
				if err != nil {
					return nil, fmt.Errorf("listing files: %w", err)
				}
			} else {
				parent, err := commit.Parent(0)
				if err != nil {
					return nil, fmt.Errorf("getting parent commit: %w", err)
				}
				patch, err := parent.Patch(commit)
				if err != nil {
					return nil, fmt.Errorf("computing diff: %w", err)
				}
				sb.WriteString(patch.String())
			}

			return map[string]any{"output": sb.String()}, nil
		})
}
