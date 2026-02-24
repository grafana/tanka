package agent

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"

	"github.com/grafana/tanka/pkg/term"
)

var (
	colorUser     = color.New(color.FgBlue, color.Bold)
	colorToolCall = color.New(color.FgCyan, color.Bold)
	colorToolResp = color.New(color.FgGreen, color.Bold)
	colorDimGrey  = color.New(color.FgHiBlack)
	colorLLMText  = color.New(color.FgYellow)

	styleToolBorder = lipgloss.NewStyle().
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			PaddingLeft(1)
)

// logUser prints the user message in blue bold.
func (a *Agent) logUser(msg string) {
	colorUser.Fprintf(a.output, "→ user: %s\n\n", msg)
}

// logToolCall prints a cyan function-call line: ▶ name(k="v", k2=42)
func (a *Agent) logToolCall(name string, args map[string]any) {
	colorToolCall.Fprintf(a.output, "▶ %s%s\n", name, formatArgs(args))
}

// logToolResponse prints a green header and the full output in a left-bordered block.
// Diff output is colourised; other output is printed in dim grey.
func (a *Agent) logToolResponse(name string, response map[string]any) {
	colorToolResp.Fprintf(a.output, "◀ %s\n", name)
	if output, ok := response["output"].(string); ok {
		var content string
		if isDiff(output) {
			content = term.Colordiff(output).String()
		} else {
			content = colorDimGrey.Sprint(output)
		}
		fmt.Fprintln(a.output, styleToolBorder.Render(content))
	}
}

// isDiff returns true if s looks like a unified or git diff.
func isDiff(s string) bool {
	return strings.HasPrefix(s, "diff ") ||
		strings.HasPrefix(s, "---") ||
		strings.Contains(s, "\n+++ ")
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
