package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"google.golang.org/adk/session"
)

// Display handles all terminal output for a Run() turn.
// Implementations receive the full event stream and are responsible for
// collecting and returning the final LLM response text.
type Display interface {
	// Event processes a single event from the ADK runner event stream.
	Event(event *session.Event)
	// Error displays a terminal error from the runner. The Run loop will call
	// FinalText immediately after, so implementations must remain usable.
	Error(err error)
	// PrintFinalText prints the accumulated final-response text. It must be called
	// exactly once, after all Event/Error calls are complete.
	PrintFinalText()
}

func NewDisplay(out io.Writer, tty bool, verbose bool) Display {
	if verbose {
		return NewVerboseDisplay(out)
	} else {
		return NewPrettyDisplay(out, tty)
	}
}

// formatArgs returns a "(k="v", k2=42)" string with keys sorted for determinism.
func formatArgs(args map[string]any) string {
	if len(args) == 0 {
		return "()"
	}
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf []byte
	buf = append(buf, '(')
	for i, k := range keys {
		if i > 0 {
			buf = append(buf, ',', ' ')
		}
		buf = append(buf, k...)
		buf = append(buf, '=')
		buf = append(buf, formatArgValue(args[k])...)
	}
	buf = append(buf, ')')
	return string(buf)
}

// formatArgValue formats a single argument value: strings are quoted and
// truncated at 80 chars; all other types use %v.
func formatArgValue(v any) string {
	switch s := v.(type) {
	case string:
		q := fmt.Sprintf("%q", s)
		if len(q) > 82 { // 80 chars + two quotes
			q = fmt.Sprintf("%q", s[:80])
			q = q[:len(q)-1] + `..."` // replace trailing " with ..."
		}
		return q
	default:
		return fmt.Sprintf("%v", v)
	}
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
