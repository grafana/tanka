//go:build !windows

package agent

import (
	"context"
	"os"
	"syscall"
	"time"

	"golang.org/x/term"
)

// watchForESC polls stdin in raw+nonblocking mode and cancels the context when
// an ESC keystroke (0x1b) is detected. It returns when either ESC is pressed or
// the context is done. Terminal state is always restored before returning.
func watchForESC(ctx context.Context, cancel context.CancelFunc) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return
	}
	defer term.Restore(fd, oldState) //nolint:errcheck

	if err := syscall.SetNonblock(fd, true); err != nil {
		return
	}
	defer syscall.SetNonblock(fd, false) //nolint:errcheck

	buf := make([]byte, 1)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, err := syscall.Read(fd, buf)
		if n == 1 && buf[0] == 0x1b {
			cancel()
			return
		}
		if err != nil && err != syscall.EAGAIN && err != syscall.EWOULDBLOCK {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(50 * time.Millisecond):
		}
	}
}
