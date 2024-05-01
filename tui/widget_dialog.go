package tui

import (
	"fmt"
	"strings"

	"git.qowevisa.me/Qowevisa/gotell/errors"
)

func createDialog(message, title string) (dialog, error) {
	var buf string
	width, height := UI.getSizes()
	if width == 0 {
		return dialog{}, errors.WrapErr("width", errors.NOT_INIT)
	}
	if height == 0 {
		return dialog{}, errors.WrapErr("height", errors.NOT_INIT)
	}
	maxWidth := width / 3
	maxHeight := 5
	errMsgLen := len(message)
	innerPart := maxWidth - 2
	if errMsgLen <= innerPart {
		maxWidth = errMsgLen + 2
	} else {
		for {
			if errMsgLen <= innerPart {
				break
			}
			maxHeight++
			errMsgLen -= innerPart
		}
	}
	innerPart = maxWidth - 2
	col := (width - maxWidth) / 2
	row := (height - maxHeight) / 2
	startCol := col
	startRow := row

	buf += getBufForMovingCursorTo(row, col)
	buf += strings.Repeat("-", maxWidth)

	row++
	buf += getBufForMovingCursorTo(row, col)
	buf += centerText(maxWidth, title)

	row++
	buf += getBufForMovingCursorTo(row, col)
	buf += strings.Repeat("-", maxWidth)

	startI := 0
	endI := innerPart
	for i := 3; i < maxHeight-1; i++ {
		var tmp string
		if endI > len(message) {
			tmp = message[startI:]
		} else {
			tmp = message[startI:endI]
		}
		row++
		buf += getBufForMovingCursorTo(row, col)
		var spaces string
		if innerPart > len(tmp) {
			spaces = strings.Repeat(" ", innerPart-len(tmp))
		}
		buf += fmt.Sprintf("|%s%s|", tmp, spaces)
		startI += innerPart
		endI += innerPart
	}

	row++
	buf += getBufForMovingCursorTo(row, col)
	buf += strings.Repeat("-", maxWidth)

	row++
	buf += getBufForMovingCursorTo(row, col)
	buf += centerText(maxWidth, "YES | NO")

	row++
	buf += getBufForMovingCursorTo(row, col)
	buf += strings.Repeat("-", maxWidth)
	return dialog{
		Row:    startRow,
		Col:    startCol,
		Width:  maxWidth,
		Height: row - startRow + 1,
		Buf:    buf,
	}, nil
}

const (
	_int_CatcherNone = iota
	_int_Catcher1Arrow
	_int_Catcher2Arrow
	_int_CatcherArrow
)

func dialogRuneCatcher(t *TUI, runes chan (rune)) error {
	state := _int_CatcherNone
	for r := range runes {
		if r == 27 && state == _int_CatcherNone {
			state = _int_Catcher1Arrow
		} else {
			continue
		}
		if r == 91 && state == _int_Catcher1Arrow {
			state = _int_Catcher2Arrow
		} else {
			continue
		}
		if state == _int_Catcher2Arrow {
			switch r {
			case 65:
				t.mySignals <- mySignal{Type: MY_SIGNAL_MOVE_CURSOR_UP}
			case 66:
				t.mySignals <- mySignal{Type: MY_SIGNAL_MOVE_CURSOR_DOWN}
			case 67:
				t.mySignals <- mySignal{Type: MY_SIGNAL_MOVE_CURSOR_RIGHT}
			case 68:
				t.mySignals <- mySignal{Type: MY_SIGNAL_MOVE_CURSOR_LEFT}
			}
		}
	}
	return nil
}
