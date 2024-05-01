package tui

import "git.qowevisa.me/Qowevisa/gotell/errors"

func (t *TUI) launchErrorsChannel() error {
	if t.errorsChannel == nil {
		return errors.WrapErr("t.errors", errors.NOT_INIT)
	}
	go func() {
		for err := range t.errorsChannel {
			if err != nil {
				t.createNotification(err.Error(), CONST_ERROR_N_TITLE)
			}
		}
	}()
	return nil
}
