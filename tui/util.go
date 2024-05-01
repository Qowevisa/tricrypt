package tui

import "git.qowevisa.me/Qowevisa/gotell/errors"

func (t *TUI) SendMessageToServer(title string, minW int) {
	var msg []rune
	h, d := SendMessageToConnectionEasy(&msg)
	err := t.addWidget(widgetConfig{
		Input:           &msg,
		Title:           title,
		MinWidth:        minW,
		HasBorder:       true,
		WidgetPosConfig: widgetPosGeneralCenter,
		CursorPosConfig: cursorPosGeneralCenter,
		DataHandler:     h,
		Data:            d,
		Next:            nil,
		Finale:          nil,
	})
	if err != nil {
		t.errorsChannel <- errors.WrapErr("t.addWidget", err)
	}
	err = t.drawSelectedWidget()
	if err != nil {
		t.errorsChannel <- errors.WrapErr("t.drawSelectedWidget", err)
	}
}
