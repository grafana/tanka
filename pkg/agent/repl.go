package agent

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

const replPrompt = "❯ "

// RunREPL starts an interactive REPL session with the agent.
// Supports:
//   - /exit or Ctrl+D: exit cleanly
//   - /clear: reset the conversation (start fresh)
//   - Ctrl+C: interrupt the current operation but stay in the REPL
func RunREPL(ctx context.Context, a *Agent, out io.Writer) error {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          replPrompt,
		HistoryFile:     historyFilePath(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		Stdin:           os.Stdin,
		Stdout:          out,
		Stderr:          os.Stderr,
	})
	if err != nil {
		return fmt.Errorf("initialising REPL: %w", err)
	}
	defer rl.Close()

	fmt.Fprintln(out, "Tanka Agent — type your request, /clear to reset, /exit or Ctrl+D to quit.")
	fmt.Fprintln(out)

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				// Ctrl+C — just show a fresh prompt
				continue
			}
			// EOF (Ctrl+D) or other read error — exit cleanly
			if err == io.EOF {
				fmt.Fprintln(out, "\nGoodbye!")
				return nil
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle REPL commands
		switch strings.ToLower(line) {
		case "/exit", "/quit":
			fmt.Fprintln(out, "Goodbye!")
			return nil
		case "/clear":
			a.Reset()
			fmt.Fprintln(out, "Conversation cleared.")
			continue
		case "/help":
			printHelp(out)
			continue
		}

		// Send to agent
		response, err := a.Run(ctx, line)
		if err != nil {
			fmt.Fprintf(out, "Error: %s\n\n", err)
			continue
		}
		if response != "" {
			fmt.Fprintln(out, response)
			fmt.Fprintln(out)
		}
	}
}

func printHelp(out io.Writer) {
	fmt.Fprintln(out, `Available commands:
  /clear   Reset the conversation (start a fresh session)
  /exit    Exit the REPL (also: Ctrl+D)
  /help    Show this help

Keyboard shortcuts:
  Ctrl+C   Interrupt the current operation; stay in REPL
  Ctrl+D   Exit
  Up/Down  Navigate history`)
	fmt.Fprintln(out)
}

// historyFilePath returns the path to the readline history file.
func historyFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// Ensure the directory exists
	dir := home + "/.config/tanka"
	_ = os.MkdirAll(dir, 0700)
	return dir + "/agent_history"
}
