package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/grafana/tanka/pkg/agent/tools"
)

const systemPrompt = `You are a senior Site Reliability Engineer specializing in Tanka and Jsonnet-based
Kubernetes configuration management.

Principles you always follow:
- Minimize blast radius: prefer the smallest possible scope of change. Target one
  environment at a time; don't touch shared libraries unless essential.
- Staged rollout: change one environment first, validate it, then propagate to others.
- Minimal, reviewable changes: keep diffs small and easy for humans to review.
  Avoid reformatting unrelated code.
- Always validate: after every file change, call tanka_validate_jsonnet on modified
  files, run tanka_diff for each affected environment (if cluster is reachable), and
  show the user the git diff and tanka diff before suggesting next steps.
- You cannot and will not apply to any cluster. Your job is to prepare correct,
  validated configuration changes for human review and deployment.

Tanka workflow reminders:
- Environments live in subdirectories (often environments/) and have a spec.json
- Shared libraries live in lib/ or vendor/
- Always run tanka_find_environments to discover the repo structure before making changes
- After making changes: validate jsonnet → git_diff → (optional) tanka_diff → git_add → git_commit
- Prefer creating a new branch (git_branch_create) before making changes
- Suggest a pull request (github_push + github_pr_create) when the work is complete`

const (
	agentAppName = "tanka"
	agentUserID  = "user"
)

// Agent orchestrates the conversation loop between the user and the LLM via
// the ADK runner, dispatching tool calls and maintaining session history.
type Agent struct {
	runner     *runner.Runner
	sessionSvc session.Service
	sessionID  string
	output     io.Writer
}

// NewAgent creates an Agent with all tools registered for the given repository root.
func NewAgent(ctx context.Context, llm model.LLM, repoRoot string) (*Agent, error) {
	adkTools, err := collectTools(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("registering tools: %w", err)
	}

	ag, err := llmagent.New(llmagent.Config{
		Name:        "tanka-agent",
		Model:       llm,
		Instruction: systemPrompt,
		Tools:       adkTools,
	})
	if err != nil {
		return nil, fmt.Errorf("creating agent: %w", err)
	}

	sessionSvc := session.InMemoryService()

	r, err := runner.New(runner.Config{
		AppName:        agentAppName,
		Agent:          ag,
		SessionService: sessionSvc,
	})
	if err != nil {
		return nil, fmt.Errorf("creating runner: %w", err)
	}

	a := &Agent{
		runner:     r,
		sessionSvc: sessionSvc,
		sessionID:  newSessionID(),
		output:     os.Stdout,
	}
	if err := a.createSession(ctx); err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}
	return a, nil
}

// collectTools gathers all ADK tools for the given repository root.
func collectTools(repoRoot string) ([]adktool.Tool, error) {
	var all []adktool.Tool
	for _, fn := range []func(string) ([]adktool.Tool, error){
		tools.NewFileTools,
		tools.NewGitTools,
		tools.NewGitHubTools,
		tools.NewTankaTools,
		tools.NewJBTools,
	} {
		ts, err := fn(repoRoot)
		if err != nil {
			return nil, err
		}
		all = append(all, ts...)
	}
	return all, nil
}

// Reset clears the conversation history by starting a fresh session.
func (a *Agent) Reset(ctx context.Context) {
	a.sessionID = newSessionID()
	_ = a.createSession(ctx)
}

// Run sends a user message and processes the full agent loop until the LLM
// returns a final text response. Returns the final text response.
func (a *Agent) Run(ctx context.Context, userInput string) (string, error) {
	msg := genai.NewContentFromText(userInput, genai.RoleUser)

	var finalText strings.Builder
	for event, err := range a.runner.Run(ctx, agentUserID, a.sessionID, msg, agent.RunConfig{}) {
		if err != nil {
			return "", fmt.Errorf("agent error: %w", err)
		}
		if event.Content == nil {
			continue
		}

		for _, part := range event.Content.Parts {
			switch {
			case part.FunctionCall != nil:
				argsJSON, _ := json.Marshal(part.FunctionCall.Args)
				fmt.Fprintf(a.output, "[tool: %s] %s\n", part.FunctionCall.Name, summarize(string(argsJSON), 120))
			case part.FunctionResponse != nil:
				if output, ok := part.FunctionResponse.Response["output"].(string); ok {
					fmt.Fprintf(a.output, "[tool: %s] %s\n", part.FunctionResponse.Name, summarize(output, 120))
				}
			case part.Text != "":
				if event.IsFinalResponse() {
					finalText.WriteString(part.Text)
				} else {
					// Print text that precedes tool calls
					fmt.Fprintf(a.output, "%s\n\n", part.Text)
				}
			}
		}
	}

	return finalText.String(), nil
}

func (a *Agent) createSession(ctx context.Context) error {
	_, err := a.sessionSvc.Create(ctx, &session.CreateRequest{
		AppName:   agentAppName,
		UserID:    agentUserID,
		SessionID: a.sessionID,
	})
	return err
}

func newSessionID() string {
	return fmt.Sprintf("s%d", time.Now().UnixNano())
}

// summarize truncates a string to maxLen characters for display purposes.
func summarize(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
