package agent

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/muesli/reflow/wordwrap"
	"google.golang.org/adk/session"
)

var spinnerFrames = []string{"|", "/", "-", "\\"}

type result struct {
	event *session.Event
	err   error
}

// prettyDisplay is the default non-verbose display. It runs a background
// goroutine that animates a spinner on a single line while the agent works,
// keeping the terminal uncluttered. FinalText shuts the goroutine down.
type prettyDisplay struct {
	out   io.Writer
	tty   bool // true when output is a real TTY; enables ANSI escape sequences
	frame int
	msg   string

	ch        chan result
	wg        sync.WaitGroup
	finalText *strings.Builder
}

func newPrettyDisplay(out io.Writer, tty bool) *prettyDisplay {
	d := &prettyDisplay{
		out:       out,
		tty:       tty,
		ch:        make(chan result),
		finalText: &strings.Builder{},
	}
	d.wg.Add(1)
	go d.run()
	return d
}

func (d *prettyDisplay) run() {
	defer d.wg.Done()

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	d.msg = "tanka-ring..."

	for {
		select {
		case r, ok := <-d.ch:
			if !ok {
				d.clear()
				return
			}
			if r.err != nil {
				d.clear()
				fmt.Fprintf(d.out, "Error: %s\n", r.err)
				continue
			}
			if r.event.Content == nil {
				continue
			}
			for _, part := range r.event.Content.Parts {
				switch {
				case part.FunctionCall != nil:
					d.print("▶ %s%s", part.FunctionCall.Name, formatArgs(part.FunctionCall.Args))
				case part.FunctionResponse != nil:
					if output, ok := part.FunctionResponse.Response["output"].(string); ok {
						d.print("◀ %s - %s", part.FunctionResponse.Name, output)
					} else {
						d.print("◀ %s", part.FunctionResponse.Name)
					}
				case part.Text != "":
					if r.event.IsFinalResponse() {
						d.finalText.WriteString(part.Text)
					} else {
						d.clear()
						wrapped := wordwrap.String(strings.TrimSpace(part.Text), 80)
						colorLLMText.Fprintln(d.out, wrapped)
					}
				}
			}
		case <-ticker.C:
			d.tick()
		}
	}
}

func (d *prettyDisplay) Event(event *session.Event) {
	d.ch <- result{event: event}
}

func (d *prettyDisplay) Error(err error) {
	d.ch <- result{err: err}
}

// FinalText closes the event channel, waits for the background goroutine to
// finish, then returns the accumulated final-response text.
func (d *prettyDisplay) FinalText() string {
	close(d.ch)
	d.wg.Wait()
	return d.finalText.String()
}

// print formats the message, truncates to 120 chars (stripping newlines), and
// writes it with a spinner. Accepts fmt-style format and args.
func (d *prettyDisplay) print(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > 120 {
		s = s[:120] + "..."
	}
	d.msg = s
	d.frame++
	frame := spinnerFrames[d.frame%len(spinnerFrames)]
	if d.tty {
		fmt.Fprintf(d.out, "\r\033[2K%s %s", frame, s)
	} else {
		fmt.Fprintf(d.out, "%s %s\n", frame, s)
	}
}

func (d *prettyDisplay) tick() {
	if d.tty {
		d.frame++
		fmt.Fprintf(d.out, "\r\033[2K%s %s", spinnerFrames[d.frame%len(spinnerFrames)], d.msg)
	}
}

func (d *prettyDisplay) clear() {
	if d.tty {
		fmt.Fprintf(d.out, "\r\033[2K")
	}
}
