package tui

import (
	"log"
	"syscall"

	"git.qowevisa.me/Qowevisa/gotell/errors"
)

func (t *TUI) launchSignalsChannel() error {
	if t.osSignals == nil {
		return errors.WrapErr("t.osSignals", errors.NOT_INIT)
	}
	go func() {
		for sig := range t.osSignals {
			log.Printf("Receive OS.signal: %#v\n", sig)
			switch sig {
			case syscall.SIGWINCH:
				t.errorsChannel <- t.redraw()
			}
		}
	}()
	return nil
}
