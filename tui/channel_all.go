package tui

func (t *TUI) launchAllChannels() error {
	var err error
	err = t.launchErrorsChannel()
	if err != nil {
		return err
	}
	t.errorsChannel <- t.launchInputChannel()
	t.errorsChannel <- t.launchInputChannel()
	t.errorsChannel <- t.launchMessageChannel()
	t.errorsChannel <- t.launchSignalsChannel()
	t.errorsChannel <- t.launchStateChannel()
	return nil
}
