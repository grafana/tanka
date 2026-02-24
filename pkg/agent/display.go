package agent

import (
	"fmt"
	"sort"

	"google.golang.org/adk/session"
)

// display handles all terminal output for a Run() turn.
// Implementations receive the full event stream and are responsible for
// collecting and returning the final LLM response text.
type display interface {
	// Event processes a single event from the ADK runner event stream.
	Event(event *session.Event)
	// Error displays a terminal error from the runner. The Run loop will call
	// FinalText immediately after, so implementations must remain usable.
	Error(err error)
	// FinalText returns the accumulated final-response text. It must be called
	// exactly once, after all Event/Error calls are complete.
	FinalText() string
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
