package tui

import (
	"log"

	"git.qowevisa.me/Qowevisa/gotell/errors"
)

func (t *TUI) launchInputChannel() error {
	if t.inputChannel == nil {
		return errors.WrapErr("t.inputChannel", errors.NOT_INIT)
	}
	go func() {
		for r := range t.inputChannel {
			log.Printf("Read rune: %#v from t.input\n", r)
			selWidget := t.selectedWidget
			// Enter
			if r == 13 {
				log.Printf("Debug: selWidget: %#v ; selWidget.Input: %#v", selWidget, selWidget.Input)
				if selWidget != nil && selWidget.Input != nil {
					t.errorsChannel <- selWidget.Handler(t, selWidget.Data)
					buf := selWidget.Clear()
					t.writeMu.Lock()
					err := t.write(buf)
					t.writeMu.Unlock()
					if err != nil {
						t.errorsChannel <- errors.WrapErr("t.write", err)
					}
					if selWidget.Next != nil {
						log.Printf("Seeing that widget.Next is not nil")
						t.errorsChannel <- t.addWidget(*selWidget.Next)
						t.errorsChannel <- t.drawSelectedWidget()
					} else {
						t.readInputState <- false
					}
					if selWidget.Finale != nil {
						log.Printf("Seeing that widget.Finale is not nil")
						t.errorsChannel <- selWidget.Finale(t, selWidget.FinaleData)
					}
				}
				continue
			} else if r == 127 {
				log.Printf("seeing r = 127")
				sliceLen := len(*selWidget.Input)
				log.Printf("sliceLen = %d\n", sliceLen)
				if sliceLen > 0 {
					log.Printf("sliceLen > 0")
					*selWidget.Input = (*selWidget.Input)[:sliceLen-1]
					t.writeMu.Lock()
					t.errorsChannel <- t.moveCursor(t.cursorPosRow, t.cursorPosCol-1)
					t.errorsChannel <- t.writeRune(' ')
					t.errorsChannel <- t.moveCursor(t.cursorPosRow, t.cursorPosCol-1)
					t.writeMu.Unlock()
				}
				continue
			}
			if selWidget != nil && selWidget.Input != nil {
				log.Printf("t.input: append %#v to widget input\n", r)
				*selWidget.Input = append(*selWidget.Input, r)
				log.Printf("t.input: trying to write %#v", r)
				t.writeMu.Lock()
				err := t.writeRune(r)
				t.writeMu.Unlock()
				if err != nil {
					t.errorsChannel <- err
				}
			}
		}
	}()
	return nil
}
