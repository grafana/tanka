package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

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

// Agent orchestrates the conversation loop between the user and the LLM,
// dispatching tool calls and feeding results back into the conversation.
type Agent struct {
	provider Provider
	allTools []tools.Tool
	messages []Message
	output   io.Writer
}

// NewAgent creates an Agent with all tools registered for the given repository root.
func NewAgent(provider Provider, repoRoot string) *Agent {
	allTools := collectTools(repoRoot)
	return &Agent{
		provider: provider,
		allTools: allTools,
		messages: nil,
		output:   os.Stdout,
	}
}

// collectTools gathers all tool definitions for the given repository root.
func collectTools(repoRoot string) []tools.Tool {
	var all []tools.Tool
	all = append(all, tools.NewFileTools(repoRoot)...)
	all = append(all, tools.NewGitTools(repoRoot)...)
	all = append(all, tools.NewGitHubTools(repoRoot)...)
	all = append(all, tools.NewTankaTools(repoRoot)...)
	all = append(all, tools.NewJBTools(repoRoot)...)
	return all
}

// Reset clears the conversation history for a fresh session.
func (a *Agent) Reset() {
	a.messages = nil
}

// Run sends a user message and processes the full agent loop until the LLM
// returns a final text response (no more tool calls).
// Returns the final text response from the assistant.
func (a *Agent) Run(ctx context.Context, userInput string) (string, error) {
	// Append user message to history
	a.messages = append(a.messages, TextMessage(RoleUser, userInput))

	for {
		// Call the LLM
		response, err := a.provider.Chat(ctx, systemPrompt, a.messages, a.allTools)
		if err != nil {
			return "", fmt.Errorf("LLM error: %w", err)
		}

		// Append assistant response to history
		a.messages = append(a.messages, *response)

		// Check if we have any tool calls to execute
		toolCalls := filterContent(response.Content, ContentTypeToolUse)
		if len(toolCalls) == 0 {
			// No tool calls — return the final text response
			return extractText(response.Content), nil
		}

		// Print any text that preceded the tool calls
		if text := extractText(response.Content); text != "" {
			fmt.Fprintf(a.output, "%s\n\n", text)
		}

		// Execute each tool call and collect results
		toolResultContents := make([]Content, 0, len(toolCalls))
		for _, tc := range toolCalls {
			result, toolErr := a.executeTool(ctx, tc)
			isError := toolErr != nil
			resultText := result
			if toolErr != nil {
				resultText = fmt.Sprintf("Error: %s", toolErr.Error())
			}

			fmt.Fprintf(a.output, "[tool: %s] %s\n", tc.Name, summarize(resultText, 120))

			toolResultContents = append(toolResultContents, Content{
				Type:      ContentTypeToolResult,
				ToolUseID: tc.ID,
				Name:      tc.Name,
				Text:      resultText,
				IsError:   isError,
			})
		}

		// Feed all tool results back as a single user message
		a.messages = append(a.messages, Message{
			Role:    RoleUser,
			Content: toolResultContents,
		})
	}
}

// executeTool finds and runs the tool matching the given tool-use content block.
func (a *Agent) executeTool(ctx context.Context, tc Content) (string, error) {
	for _, t := range a.allTools {
		if t.Name == tc.Name {
			input := tc.Input
			if len(input) == 0 {
				input = []byte("{}")
			}
			return t.Execute(ctx, json.RawMessage(input))
		}
	}
	return "", fmt.Errorf("unknown tool %q", tc.Name)
}

// filterContent returns content blocks matching the given type.
func filterContent(content []Content, ct ContentType) []Content {
	var out []Content
	for _, c := range content {
		if c.Type == ct {
			out = append(out, c)
		}
	}
	return out
}

// extractText concatenates all text blocks in a content slice.
func extractText(content []Content) string {
	var parts []string
	for _, c := range content {
		if c.Type == ContentTypeText && c.Text != "" {
			parts = append(parts, c.Text)
		}
	}
	return strings.Join(parts, "\n")
}

// summarize truncates a string to maxLen characters for display purposes.
func summarize(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	// Replace newlines with spaces for single-line display
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
