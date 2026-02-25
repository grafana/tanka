package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-clix/cli"
	gogit "github.com/go-git/go-git/v5"
	"golang.org/x/term"

	"github.com/grafana/tanka/pkg/agent"
	"github.com/grafana/tanka/pkg/agent/models"
)

func agentCmd(ctx context.Context) *cli.Command {
	cmd := &cli.Command{
		Use:   "agent [query]",
		Short: "AI assistant for Tanka configuration (prepare-only: edits files, validates, diffs, commits — never applies to cluster)",
		Long: `AI assistant for Tanka and Jsonnet-based Kubernetes configuration management.

The agent acts as a senior SRE: it edits Jsonnet files, validates them, shows diffs,
creates branches, and opens pull requests — but never runs 'tanka apply'. All cluster
deploys remain human-controlled.

MODES
  REPL (interactive):   tk agent
  One-shot:             tk agent "create a staging environment"

REPL COMMANDS
  /clear   Reset conversation (start fresh session)
  /exit    Exit (also Ctrl+D)
  Ctrl+C   Cancel current operation, stay in REPL

PROVIDER SELECTION
  Default provider: gemini (gemini-2.0-flash)

  Set provider and model via environment variables:
    TANKA_AGENT_PROVIDER=gemini|anthropic|openai
    TANKA_AGENT_MODEL=<model-id>

  Or via config file (~/.config/tanka/agent.yaml):
    provider: gemini        # gemini | anthropic | openai
    model: gemini-2.0-flash # e.g. claude-opus-4-6, gpt-4o

  Environment variables take priority over the config file.

API KEYS
  Gemini:    GEMINI_API_KEY  or  GOOGLE_API_KEY
  Anthropic: ANTHROPIC_API_KEY
  OpenAI:    OPENAI_API_KEY

  GitHub operations (push, PR creation) also require GITHUB_TOKEN.

REQUIREMENTS
  Must be run from inside a git repository.`,
		Args: cli.ArgsAny(),
	}

	providerFlag := cmd.Flags().String("provider", "",
		`LLM provider: gemini (default), anthropic, or openai.
    Overrides TANKA_AGENT_PROVIDER env var and ~/.config/tanka/agent.yaml.
    API key env vars: GEMINI_API_KEY / GOOGLE_API_KEY, ANTHROPIC_API_KEY, OPENAI_API_KEY`)
	modelFlag := cmd.Flags().String("model", "",
		`Model identifier for the chosen provider.
    Overrides TANKA_AGENT_MODEL env var and ~/.config/tanka/agent.yaml.
    Examples: gemini-2.0-flash (default), claude-opus-4-6, gpt-4o`)
	verboseFlag := cmd.Flags().BoolP("verbose", "v", false,
		"Verbose output: show all LLM messages, tool calls, and full responses with color")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		// 1. Verify we are inside a git repository and find the root
		repoRoot, err := findGitRoot()
		if err != nil {
			return fmt.Errorf("tk agent must be run inside a git repository: %w", err)
		}

		// 2. Load configuration (env vars > config file > defaults)
		cfg, err := agent.LoadConfig()
		if err != nil {
			return fmt.Errorf("loading agent config: %w", err)
		}

		// 3. CLI flags override everything else
		if *providerFlag != "" {
			cfg.Provider = *providerFlag
		}
		if *modelFlag != "" {
			cfg.Model = *modelFlag
		}

		// 4. Validate that an API key is available for the chosen provider
		if err := cfg.Validate(); err != nil {
			return err
		}

		// 5. Construct the model and agent
		llm, err := models.NewModel(ctx, &models.ModelConfig{
			Provider: cfg.Provider,
			Model:    cfg.Model,
			APIKey:   cfg.APIKeyForProvider(),
		})
		if err != nil {
			return fmt.Errorf("initialising model: %w", err)
		}
		a, err := agent.NewAgent(ctx, llm, repoRoot)
		if err != nil {
			return fmt.Errorf("initialising agent: %w", err)
		}

		// 6. Run in one-shot or REPL mode
		if len(args) > 0 {
			// One-shot: tk agent "do something"
			query := args[0]
			tty := term.IsTerminal(int(os.Stdout.Fd()))
			display := agent.NewDisplay(os.Stdout, tty, *verboseFlag)
			err := a.Run(ctx, query, display)
			display.PrintFinalText()
			return err
		}

		// REPL: tk agent (no arguments)
		return agent.RunREPL(ctx, a, os.Stdout, *verboseFlag)
	}

	return cmd
}

// findGitRoot walks up the directory tree to locate the .git directory and
// returns the absolute path to the repository root.
func findGitRoot() (string, error) {
	r, err := gogit.PlainOpenWithOptions(".", &gogit.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		if err == gogit.ErrRepositoryNotExists {
			return "", fmt.Errorf("not inside a git repository")
		}
		return "", err
	}
	wt, err := r.Worktree()
	if err != nil {
		return "", fmt.Errorf("getting worktree: %w", err)
	}
	return wt.Filesystem.Root(), nil
}
