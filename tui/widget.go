package tui

import (
	"fmt"
	"strings"

	"git.qowevisa.me/Qowevisa/gotell/errors"
)

// func (w *widget) init(startX, startY, endX, endY int) error {
// 	width, height := UI.getSizes()
//
// 	return nil
// }

func (w *widget) init(config widgetConfig) error {
	width, height := UI.getSizes()
	if width == 0 {
		return errors.WrapErr("width", errors.NOT_INIT)
	}
	if height == 0 {
		return errors.WrapErr("height", errors.NOT_INIT)
	}
	var startRow, startCol int
	// I guess I really has to do it that way bc go doesn't like my C style
	//   of writting `-2 * config.HasBorder` or `-2 * (config.HasBorder == true)`
	var factorIfHasBorder int
	if config.HasBorder {
		factorIfHasBorder = 1
	} else {
		factorIfHasBorder = 0
	}
	if config.WidgetPosConfig.isGeneral() {
		switch config.WidgetPosConfig {
		case widgetPosGeneralCenter:
			startCol = (width - len(config.Title)) / 2
			startRow = (height - 3 - 2*factorIfHasBorder) / 2
		case widgetPosGeneralLeftCenter:
			startCol = 0
			startRow = (height - 3 - 2*factorIfHasBorder) / 2
		case widgetPosGeneralRightCenter:
			startCol = (width - len(config.Title) - 2*factorIfHasBorder)
			startRow = (height - 3 - 2*factorIfHasBorder) / 2
		default:
			return errors.WrapErr(fmt.Sprintf("config.WidgetPosConfig: %d :", config.WidgetPosConfig), errors.NOT_HANDLED)
		}
	}
	snappedPair := tuiPointPair{
		startCol: startCol,
		startRow: startRow,
		endCol:   startCol + len(config.Title) + 2,
		endRow:   startRow + 5,
	}
	percentPair := tuiPointPair{
		startCol: GetIntPercentFromData(snappedPair.startCol, width),
		startRow: GetIntPercentFromData(snappedPair.startRow, height),
		endCol:   GetIntPercentFromData(snappedPair.endCol, width),
		endRow:   GetIntPercentFromData(snappedPair.endRow, height),
	}
	w.percentPair = percentPair
	w.snappedPair = snappedPair
	w.startupConfig = config
	w.Input = config.Input
	w.Handler = config.DataHandler
	w.Data = config.Data
	w.Title = config.Title
	w.MinWidth = config.MinWidth
	w.MinHeight = config.MinHeight
	if w.MinWidth > len(config.Title) {
		w.Width = w.MinWidth
	} else {
		if w.startupConfig.HasBorder {
			w.Width = len(config.Title) + 2
		} else {
			w.Width = len(config.Title)
		}
	}
	if config.HasBorder {
		if config.MinHeight > 5 {
			w.Height = config.MinHeight
		} else {
			w.Height = 5
		}
	}
	w.Next = config.Next
	w.Finale = config.Finale
	w.FinaleData = config.FinaleData
	return nil
}

func getBufForMovingCursorTo(row, col int) string {
	return fmt.Sprintf("\033[%d;%dH", row, col)
}

func (w *widget) moveToNextLine() string {
	w.row++
	return getBufForMovingCursorTo(w.row, w.col)
}

func (w *widget) Draw() (widgetDraw, error) {
	w.row = w.snappedPair.startRow
	w.col = w.snappedPair.startCol
	var buf string
	title := w.startupConfig.Title
	buf += getBufForMovingCursorTo(w.row, w.col)
	if w.startupConfig.HasBorder {
		buf += strings.Repeat("-", w.Width)
		buf += w.moveToNextLine()
		emptyLen := w.Width - len(title) - 2
		firstHalf := (emptyLen) / 2
		buf += fmt.Sprintf("|%s%s%s|",
			strings.Repeat(" ", firstHalf),
			title, strings.Repeat(" ",
				emptyLen-firstHalf))
		buf += w.moveToNextLine()
		buf += strings.Repeat("-", w.Width)
		buf += w.moveToNextLine()
		buf += fmt.Sprintf("|%s%s|", string(*w.Input), strings.Repeat(" ", w.Width-2-len(*w.Input)))
		buf += w.moveToNextLine()
		buf += strings.Repeat("-", w.Width)
		w.col++
		w.row--
		buf += getBufForMovingCursorTo(w.row, w.col)
	} else {
		buf += fmt.Sprintf("%s", title)
		buf += w.moveToNextLine()
		buf += strings.Repeat("-", w.Width)
		buf += w.moveToNextLine()
		buf += fmt.Sprintf("%s%s", string(*w.Input), strings.Repeat(" ", len(title)-len(*w.Input)))
		buf += getBufForMovingCursorTo(w.row, w.col)
	}
	return widgetDraw{
		Buf: buf,
		Row: w.row,
		Col: w.col,
	}, nil
}

func (w *widget) Clear() string {
	var buf string
	w.row = w.snappedPair.startRow
	w.col = w.snappedPair.startCol
	title := w.startupConfig.Title
	buf += getBufForMovingCursorTo(w.row, w.col)
	if w.startupConfig.HasBorder {
		buf += strings.Repeat(" ", w.Width)
		buf += w.moveToNextLine()
		buf += strings.Repeat(" ", w.Width)
		buf += w.moveToNextLine()
		buf += strings.Repeat(" ", w.Width)
		buf += w.moveToNextLine()
		maxClearInpLen := len(*w.Input) + 2
		if w.Width > maxClearInpLen {
			maxClearInpLen = w.Width
		}
		buf += fmt.Sprintf("%s", strings.Repeat(" ", maxClearInpLen))
		buf += w.moveToNextLine()
		buf += strings.Repeat(" ", w.Width)
		w.row++
		w.col--
		buf += getBufForMovingCursorTo(0, 0)
	} else {
		buf += fmt.Sprintf("%s", strings.Repeat(" ", len(title)))
		buf += w.moveToNextLine()
		buf += strings.Repeat(" ", w.Width)
		buf += w.moveToNextLine()
		maxClearInpLen := len(*w.Input)
		if w.Width > maxClearInpLen {
			maxClearInpLen = w.Width
		}
		buf += fmt.Sprintf("%s", strings.Repeat(" ", maxClearInpLen))
		buf += getBufForMovingCursorTo(0, 0)
	}
	return buf
}
