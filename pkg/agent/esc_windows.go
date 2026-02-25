//go:build windows

package agent

import "context"

// watchForESC is a no-op on Windows.
func watchForESC(_ context.Context, _ context.CancelFunc) {}
