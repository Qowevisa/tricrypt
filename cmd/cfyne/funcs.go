package main

import (
	// "fyne.io/fyne/v2/widget"
	com "git.qowevisa.me/Qowevisa/gotell/communication"
)

func getSendMessageFuncToIntercom(intercom chan *com.Message, id uint8, data []byte) func() {
	return func() {
		intercom <- &com.Message{
			ID:   id,
			Data: data,
		}
	}
}
