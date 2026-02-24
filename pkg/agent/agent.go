package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
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

Tanka workflow reminders:
- Environments live in subdirectories (often environments/) and have a spec.json
- Shared libraries live in lib/ or vendor/
- Always run tanka_find_environments to discover the repo structure before making changes
- After making changes: validate jsonnet → (optional) tanka_diff, then present a clear
  summary of every file changed and what was changed, so the user can review and commit
- Use jb_install / jb_update to manage jsonnet dependencies — never use git_* tools
  to clone or fetch packages manually
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
	output     io.Writer
	verbose    bool
	glamour    *glamour.TermRenderer // nil if init failed or not a TTY
}

// NewAgent creates an Agent with all tools registered for the given repository root.
func NewAgent(ctx context.Context, llm model.LLM, repoRoot string, verbose bool) (*Agent, error) {
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
		verbose:    verbose,
	}
	if gr, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(0)); err == nil {
		a.glamour = gr
	}
	if err := a.createSession(ctx); err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}
	return a, nil
}

// renderMarkdown renders text as styled markdown if a glamour renderer is available,
// otherwise returns the text unchanged.
func (a *Agent) renderMarkdown(text string) string {
	if a.glamour == nil {
		return text
	}
	out, err := a.glamour.Render(text)
	if err != nil {
		return text
	}
	return out
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
func (a *Agent) Run(ctx context.Context, userInput string) (string, error) {
	if a.verbose {
		a.logUser(userInput)
	}

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
				if a.verbose {
					a.logToolCall(part.FunctionCall.Name, part.FunctionCall.Args)
				} else {
					argsJSON, _ := json.Marshal(part.FunctionCall.Args)
					fmt.Fprintf(a.output, "[tool: %s] %s\n", part.FunctionCall.Name, summarize(string(argsJSON), 120))
				}
			case part.FunctionResponse != nil:
				if a.verbose {
					a.logToolResponse(part.FunctionResponse.Name, part.FunctionResponse.Response)
				} else {
					if output, ok := part.FunctionResponse.Response["output"].(string); ok {
						fmt.Fprintf(a.output, "[tool: %s] %s\n", part.FunctionResponse.Name, summarize(output, 120))
					}
				}
			case part.Text != "":
				if event.IsFinalResponse() {
					finalText.WriteString(part.Text)
				} else {
					if a.verbose {
						a.logLLMText(part.Text)
					} else {
						// Print text that precedes tool calls
						fmt.Fprintf(a.output, "%s\n\n", part.Text)
					}
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

// PrintContext dumps the full raw session history to out for debugging.
func (a *Agent) PrintContext(ctx context.Context, out io.Writer) error {
	resp, err := a.sessionSvc.Get(ctx, &session.GetRequest{
		AppName:   agentAppName,
		UserID:    agentUserID,
		SessionID: a.sessionID,
	})
	if err != nil {
		return fmt.Errorf("fetching session: %w", err)
	}

	events := resp.Session.Events()
	fmt.Fprintf(out, "=== Session context: %d event(s) ===\n\n", events.Len())

	i := 0
	for event := range events.All() {
		i++
		fmt.Fprintf(out, "--- Event %d | author=%s", i, event.Author)
		if event.Branch != "" {
			fmt.Fprintf(out, " branch=%s", event.Branch)
		}
		fmt.Fprintln(out)

		if event.ErrorMessage != "" {
			fmt.Fprintf(out, "  ERROR: %s (code=%s)\n", event.ErrorMessage, event.ErrorCode)
		}

		if event.Content != nil {
			fmt.Fprintf(out, "  role: %s\n", event.Content.Role)
			for _, part := range event.Content.Parts {
				switch {
				case part.FunctionCall != nil:
					argsJSON, _ := json.MarshalIndent(part.FunctionCall.Args, "    ", "  ")
					fmt.Fprintf(out, "  [function_call] %s\n    %s\n", part.FunctionCall.Name, argsJSON)
				case part.FunctionResponse != nil:
					respJSON, _ := json.MarshalIndent(part.FunctionResponse.Response, "    ", "  ")
					fmt.Fprintf(out, "  [function_response] %s\n    %s\n", part.FunctionResponse.Name, respJSON)
				case part.Text != "":
					fmt.Fprintf(out, "  [text] %s\n", part.Text)
				}
			}
		}

		if m := event.UsageMetadata; m != nil {
			fmt.Fprintf(out, "  tokens: prompt=%d candidates=%d total=%d\n",
				m.PromptTokenCount, m.CandidatesTokenCount, m.TotalTokenCount)
		}
		fmt.Fprintln(out)
	}
	return nil
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
