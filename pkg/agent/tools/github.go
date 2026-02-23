package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	gogitconfig "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v67/github"
)

// GitHubTools provides GitHub API operations.
type GitHubTools struct {
	repoRoot string
}

// NewGitHubTools creates GitHub tools for the repository at the given path.
// Requires GITHUB_TOKEN environment variable for authentication.
func NewGitHubTools(repoRoot string) []Tool {
	ght := &GitHubTools{repoRoot: repoRoot}
	return []Tool{
		ght.pushTool(),
		ght.prCreateTool(),
		ght.prListTool(),
	}
}

// githubClient creates an authenticated GitHub client using GITHUB_TOKEN.
func githubClient() (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is required for GitHub operations")
	}
	return github.NewClient(nil).WithAuthToken(token), nil
}

// detectOwnerRepo parses the GitHub owner and repo name from the git remote URL.
func (ght *GitHubTools) detectOwnerRepo() (owner, repo string, err error) {
	r, err := gogit.PlainOpen(ght.repoRoot)
	if err != nil {
		return "", "", fmt.Errorf("opening repository: %w", err)
	}
	remote, err := r.Remote("origin")
	if err != nil {
		return "", "", fmt.Errorf("getting origin remote: %w", err)
	}
	if len(remote.Config().URLs) == 0 {
		return "", "", fmt.Errorf("origin remote has no URLs")
	}
	return parseGitHubURL(remote.Config().URLs[0])
}

// parseGitHubURL extracts owner and repo from HTTPS or SSH GitHub URLs.
func parseGitHubURL(rawURL string) (owner, repo string, err error) {
	// SSH: git@github.com:owner/repo.git
	sshRe := regexp.MustCompile(`git@github\.com:([^/]+)/(.+?)(?:\.git)?$`)
	if m := sshRe.FindStringSubmatch(rawURL); m != nil {
		return m[1], m[2], nil
	}
	// HTTPS: https://github.com/owner/repo.git
	u, parseErr := url.Parse(rawURL)
	if parseErr != nil {
		return "", "", fmt.Errorf("parsing remote URL %q: %w", rawURL, parseErr)
	}
	if !strings.Contains(u.Host, "github.com") {
		return "", "", fmt.Errorf("remote URL %q does not appear to be a GitHub URL", rawURL)
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("cannot extract owner/repo from URL %q", rawURL)
	}
	return parts[0], strings.TrimSuffix(parts[1], ".git"), nil
}

func (ght *GitHubTools) pushTool() Tool {
	return Tool{
		Name:        "github_push",
		Description: "Push a local branch to the origin remote on GitHub. Requires GITHUB_TOKEN.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"branch": {
					"type": "string",
					"description": "Name of the local branch to push"
				}
			},
			"required": ["branch"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				Branch string `json:"branch"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}
			token := os.Getenv("GITHUB_TOKEN")
			if token == "" {
				return "", fmt.Errorf("GITHUB_TOKEN environment variable is required to push")
			}

			r, err := gogit.PlainOpen(ght.repoRoot)
			if err != nil {
				return "", fmt.Errorf("opening repository: %w", err)
			}

			err = r.PushContext(ctx, &gogit.PushOptions{
				RemoteName: "origin",
				Auth: &gogitconfig.BasicAuth{
					Username: "x-token",
					Password: token,
				},
			})
			if err == gogit.NoErrAlreadyUpToDate {
				return fmt.Sprintf("Branch %q is already up to date on origin", params.Branch), nil
			}
			if err != nil {
				return "", fmt.Errorf("pushing branch %q: %w", params.Branch, err)
			}
			return fmt.Sprintf("Successfully pushed branch %q to origin", params.Branch), nil
		},
	}
}

func (ght *GitHubTools) prCreateTool() Tool {
	return Tool{
		Name:        "github_pr_create",
		Description: "Create a pull request on GitHub. Requires GITHUB_TOKEN. The head branch should already be pushed.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"title": {
					"type": "string",
					"description": "Pull request title (keep it under 70 characters)"
				},
				"body": {
					"type": "string",
					"description": "Pull request description in Markdown"
				},
				"head_branch": {
					"type": "string",
					"description": "The branch containing the changes"
				},
				"base_branch": {
					"type": "string",
					"description": "The branch to merge into (e.g. 'main')",
					"default": "main"
				}
			},
			"required": ["title", "body", "head_branch"]
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			var params struct {
				Title      string `json:"title"`
				Body       string `json:"body"`
				HeadBranch string `json:"head_branch"`
				BaseBranch string `json:"base_branch"`
			}
			if err := json.Unmarshal(input, &params); err != nil {
				return "", fmt.Errorf("invalid parameters: %w", err)
			}
			if params.BaseBranch == "" {
				params.BaseBranch = "main"
			}

			client, err := githubClient()
			if err != nil {
				return "", err
			}
			owner, repo, err := ght.detectOwnerRepo()
			if err != nil {
				return "", err
			}

			pr, _, err := client.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
				Title: github.String(params.Title),
				Body:  github.String(params.Body),
				Head:  github.String(params.HeadBranch),
				Base:  github.String(params.BaseBranch),
			})
			if err != nil {
				return "", fmt.Errorf("creating pull request: %w", err)
			}
			return fmt.Sprintf("Created PR #%d: %s\nURL: %s", pr.GetNumber(), pr.GetTitle(), pr.GetHTMLURL()), nil
		},
	}
}

func (ght *GitHubTools) prListTool() Tool {
	return Tool{
		Name:        "github_pr_list",
		Description: "List open pull requests on GitHub for this repository. Requires GITHUB_TOKEN.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {}
		}`),
		Execute: func(ctx context.Context, input json.RawMessage) (string, error) {
			client, err := githubClient()
			if err != nil {
				return "", err
			}
			owner, repo, err := ght.detectOwnerRepo()
			if err != nil {
				return "", err
			}

			prs, _, err := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
				State: "open",
			})
			if err != nil {
				return "", fmt.Errorf("listing pull requests: %w", err)
			}
			if len(prs) == 0 {
				return "No open pull requests", nil
			}

			var sb strings.Builder
			for _, pr := range prs {
				fmt.Fprintf(&sb, "#%d %s\n  Author: %s\n  Branch: %s → %s\n  URL: %s\n\n",
					pr.GetNumber(),
					pr.GetTitle(),
					pr.GetUser().GetLogin(),
					pr.GetHead().GetRef(),
					pr.GetBase().GetRef(),
					pr.GetHTMLURL(),
				)
			}
			return sb.String(), nil
		},
	}
}
