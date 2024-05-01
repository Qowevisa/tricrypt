package tui

const (
	MY_SIGNAL_EXIT = iota
	MY_SIGNAL_MESSAGE
	MY_SIGNAL_CONNECT
	MY_SIGNAL_CLOSE
	MY_SIGNAL_MOVE_CURSOR_UP
	MY_SIGNAL_MOVE_CURSOR_DOWN
	MY_SIGNAL_MOVE_CURSOR_LEFT
	MY_SIGNAL_MOVE_CURSOR_RIGHT
)

const (
	cursorPosGeneralCenter cursorPosConfigValue = -1
	cursorPosGeneralLeft   cursorPosConfigValue = -2
	cursorPosGeneralRight  cursorPosConfigValue = -3
)

const (
	CONST_MESSAGE_LEN   = 1024
	CONST_ERROR_N_TITLE = "ERROR"
)

const (
	footerStart = "State: "
)

func (c cursorPosConfigValue) isGeneral() bool {
	switch c {
	case cursorPosGeneralCenter:
		return true
	case cursorPosGeneralLeft:
		return true
	case cursorPosGeneralRight:
		return true
	}
	return false
}

const (
	widgetPosGeneralCenter      widgetPosConfigValue = -1
	widgetPosGeneralLeftCenter  widgetPosConfigValue = -2
	widgetPosGeneralRightCenter widgetPosConfigValue = -3
)

func (w widgetPosConfigValue) isGeneral() bool {
	switch w {
	case widgetPosGeneralCenter:
		return true
	case widgetPosGeneralLeftCenter:
		return true
	case widgetPosGeneralRightCenter:
		return true
	}
	return false
}
