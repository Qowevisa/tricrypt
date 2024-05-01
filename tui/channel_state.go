package tui

import "git.qowevisa.me/Qowevisa/gotell/errors"

func (t *TUI) launchStateChannel() error {
	if t.stateChannel == nil {
		return errors.WrapErr("t.stateChannel", errors.NOT_INIT)
	}
	go func() {
		for state := range t.stateChannel {
			t.writeMu.Lock()
			oldRow, oldCol := t.getCursorPos()
			t.errorsChannel <- t.moveCursor(t.height, len(footerStart)+1)
			t._clearLine()
			t.errorsChannel <- t.write(state)
			t.errorsChannel <- t.moveCursor(oldRow, oldCol)
			t.writeMu.Unlock()
		}
	}()
	return nil
}
