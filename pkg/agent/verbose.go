package agent

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
)

var (
	colorUser     = color.New(color.FgBlue, color.Bold)
	colorToolCall = color.New(color.FgCyan, color.Bold)
	colorToolResp = color.New(color.FgGreen, color.Bold)
	colorDimGrey  = color.New(color.FgHiBlack)
	colorLLMText  = color.New(color.FgYellow)
)

// logUser prints the user message in blue bold.
func (a *Agent) logUser(msg string) {
	colorUser.Fprintf(a.output, "→ user: %s\n\n", msg)
}

// logToolCall prints a cyan function-call line: ▶ name(k="v", k2=42)
func (a *Agent) logToolCall(name string, args map[string]any) {
	colorToolCall.Fprintf(a.output, "▶ %s%s\n", name, formatArgs(args))
}

// logToolResponse prints a green header and the full output indented in dim grey.
func (a *Agent) logToolResponse(name string, response map[string]any) {
	colorToolResp.Fprintf(a.output, "◀ %s\n", name)
	if output, ok := response["output"].(string); ok {
		for _, line := range splitLines(output) {
			colorDimGrey.Fprintf(a.output, "  %s\n", line)
		}
	}
	fmt.Fprintln(a.output)
}

// logLLMText prints intermediate LLM text (before tool calls) in yellow.
func (a *Agent) logLLMText(text string) {
	colorLLMText.Fprintln(a.output, text)
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
			q = fmt.Sprintf("%q", s[:80]) + `..."`
		}
		return q
	default:
		return fmt.Sprintf("%v", v)
	}
}

// splitLines splits text into lines without importing strings in this file.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
