package agent

import (
	"context"
	"fmt"
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
  files, run tanka_diff for each affected environment (if cluster is reachable), then
  summarise the changes for the user before suggesting next steps.
- You cannot and will not apply to any cluster. Your job is to prepare correct,
  validated configuration changes for human review and deployment.
- You do not manage git: no staging, committing, or branch creation. That is the
  user's responsibility. You may read history with git_log and git_show.
- Always finish your response by calling git_diff and including the output so the
  user can see exactly what changed and is ready to review and commit.

Tanka workflow reminders:
- Environments live in subdirectories (often environments/) and have a spec.json
- Shared libraries live in lib/ or vendor/
- Always run tanka_find_environments to discover the repo structure before making changes
- After making changes: validate jsonnet → (optional) tanka_diff → git_diff, then present
  a clear summary of every file changed and what was changed, so the user can review and commit
- Use jb_install / jb_update to manage jsonnet dependencies — never use git_* tools
  to clone or fetch packages manually
- When installing with jb_install, always install packages (e.g. "github.com/jsonnet-libs/k8s-libsonnet/1.35@main"),
  never individual files (e.g. "github.com/jsonnet-libs/k8s-libsonnet/1.30/main.libsonnet")
- Use helm_dependency_build / helm_dependency_update to manage Helm chart dependencies`

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
		tools.NewTankaTools,
		tools.NewJBTools,
		tools.NewHelmTools,
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
func (a *Agent) Reset(ctx context.Context) error {
	a.sessionID = newSessionID()
	return a.createSession(ctx)
}

// Run sends a user message and processes the full agent loop until the LLM
// returns a final text response. Returns the final text response.
func (a *Agent) Run(ctx context.Context, userInput string, display Display) error {
	msg := genai.NewContentFromText(userInput, genai.RoleUser)
	var runErr error
	for event, err := range a.runner.Run(ctx, agentUserID, a.sessionID, msg, agent.RunConfig{}) {
		if err != nil {
			display.Error(err)
			runErr = err
			break
		} else {
			display.Event(event)
		}
	}

	display.PrintFinalText()
	return runErr
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
