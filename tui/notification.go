package tui

import (
	"fmt"
	"strings"

	"git.qowevisa.me/Qowevisa/gotell/errors"
)

func centerText(width int, text string) string {
	emptyLen := width - len(text) - 2
	leftEmptyLen := emptyLen / 2
	return fmt.Sprintf(
		"|%s%s%s|",
		strings.Repeat(" ", leftEmptyLen),
		text,
		strings.Repeat(" ", emptyLen-leftEmptyLen),
	)
}

func createNotification(message string) (notifier, error) {
	title := "ERROR"
	var buf string
	width, height := UI.getSizes()
	if width == 0 {
		return notifier{}, errors.WrapErr("width", errors.NOT_INIT)
	}
	if height == 0 {
		return notifier{}, errors.WrapErr("height", errors.NOT_INIT)
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
	buf += centerText(maxWidth, "OK")

	row++
	buf += getBufForMovingCursorTo(row, col)
	buf += strings.Repeat("-", maxWidth)
	return notifier{
		Row:    startRow,
		Col:    startCol,
		Width:  maxWidth,
		Height: row - startRow + 1,
		Buf:    buf,
	}, nil
}

func (n *notifier) Clear() string {
	var buf string
	for i := 0; i < n.Height; i++ {
		buf += getBufForMovingCursorTo(n.Row, n.Col)
		buf += strings.Repeat(" ", n.Width)
		n.Row++
	}
	return buf
}
