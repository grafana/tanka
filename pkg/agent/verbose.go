package agent

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"google.golang.org/adk/session"
)

var (
	colorToolCall = color.New(color.FgCyan, color.Bold)
	colorToolResp = color.New(color.FgGreen, color.Bold)
	colorError    = color.New(color.FgRed, color.Bold)
	colorLLMText  = color.New(color.FgYellow)
)

// verboseDisplay is the --verbose display. It prints every tool call and
// response in full as they arrive, with colour-coded borders and diff colouring.
type verboseDisplay struct {
	out       io.Writer
	finalText *strings.Builder
}

func NewVerboseDisplay(out io.Writer) *verboseDisplay {
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
				fmt.Fprintln(d.out, output)
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
	colorError.Fprintf(d.out, "Error: %s\n", err)
}

func (d *verboseDisplay) PrintFinalText() {
	colorLLMText.Fprintln(d.out, d.finalText.String())
}
