package agent

import (
	"context"
	"fmt"
	"io"
)

// RunOneShot sends a single request to the agent and prints the response.
// This is the non-interactive mode for queries like: tk agent "list environments"
func RunOneShot(ctx context.Context, a *Agent, query string, out io.Writer) error {
	response, err := a.Run(ctx, query)
	if err != nil {
		return fmt.Errorf("agent error: %w", err)
	}
	fmt.Fprintln(out, response)
	return nil
}
