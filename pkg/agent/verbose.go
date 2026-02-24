package agent

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"google.golang.org/adk/session"

	"github.com/grafana/tanka/pkg/term"
)

var (
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

// verboseDisplay is the --verbose display. It prints every tool call and
// response in full as they arrive, with colour-coded borders and diff colouring.
type verboseDisplay struct {
	out       io.Writer
	finalText *strings.Builder
}

func newVerboseDisplay(out io.Writer) *verboseDisplay {
	return &verboseDisplay{
		out:       out,
		finalText: &strings.Builder{},
	}
}

func (d *verboseDisplay) Event(event *session.Event) {
	if event.Content == nil {
		return
	}
	for _, part := range event.Content.Parts {
		switch {
		case part.FunctionCall != nil:
			colorToolCall.Fprintf(d.out, "▶ %s%s\n", part.FunctionCall.Name, formatArgs(part.FunctionCall.Args))
		case part.FunctionResponse != nil:
			colorToolResp.Fprintf(d.out, "◀ %s\n", part.FunctionResponse.Name)

			if output, ok := part.FunctionResponse.Response["output"].(string); ok {
				var content string
				if isDiff(output) {
					content = term.Colordiff(output).String()
				} else {
					content = colorDimGrey.Sprint(output)
				}
				fmt.Fprintln(d.out, styleToolBorder.Render(content))
			}

		case part.Text != "":
			if event.IsFinalResponse() {
				d.finalText.WriteString(part.Text)
			} else {
				colorLLMText.Fprintln(d.out, strings.TrimSpace(part.Text))
			}
		}
	}
}

func (d *verboseDisplay) Error(err error) {
	color.New(color.FgRed, color.Bold).Fprintf(d.out, "Error: %s\n", err)
}

func (d *verboseDisplay) FinalText() string {
	return d.finalText.String()
}

// isDiff returns true if s looks like a unified or git diff.
func isDiff(s string) bool {
	return strings.HasPrefix(s, "diff ") ||
		strings.HasPrefix(s, "---") ||
		strings.Contains(s, "\n+++ ")
}
