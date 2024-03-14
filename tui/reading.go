package tui

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"git.qowevisa.me/Qowevisa/gotell/errors"
)

// NOTE: should be launched as goroutine
func (t *TUI) launchReadingMessagesFromConnection(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this goroutine as done when it exits

	if t.messageChannel == nil {
		t.errors <- errors.WrapErr("t.messageChannel", errors.NOT_INIT)
		return
	}
	buf := make([]byte, CONST_MESSAGE_LEN)
	for {
		select {
		case <-ctx.Done(): // Check if context cancellation has been requested
			return
		default:
			timeoutDuration := 5 * time.Second
			err := t.tlsConnection.SetReadDeadline(time.Now().Add(timeoutDuration))
			if err != nil {
				t.errors <- errors.WrapErr("SetReadDeadline", err)
				return
			}
			n, err := t.tlsConnection.Read(buf)
			if err != nil {
				if err != io.EOF {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue
					}
					t.errors <- errors.WrapErr("t.tlsConnection.Read", err)
				}
				return
			}
			t.messageChannel <- buf[:n]
		}
	}
}
