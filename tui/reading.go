package tui

import (
	"io"

	"git.qowevisa.me/Qowevisa/gotell/errors"
)

// NOTE: should be launched as goroutine
func (t *TUI) launchReadingMessagesFromConnection() {
	if t.messageChannel == nil {
		t.errors <- errors.WrapErr("t.messageChannel", errors.NOT_INIT)
		return
	}
	buf := make([]byte, CONST_MESSAGE_LEN)
	for {
		n, err := t.tlsConnection.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.errors <- errors.WrapErr("t.tlsConnection.Read", err)
			}
			break
		}
		t.messageChannel <- buf[:n]
	}
}
