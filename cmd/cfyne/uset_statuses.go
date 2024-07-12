package main

import (
// "log"
//
// "fyne.io/fyne/v2"
// "fyne.io/fyne/v2/container"
// "fyne.io/fyne/v2/widget"
// com "git.qowevisa.me/Qowevisa/gotell/communication"
)

const (
	USER_STATUS_NOT_VISIBLE = 0 + iota
	USER_STATUS_VISIBLE
	USER_STATUS_HANDSHAKE_INIT
	USER_STATUS_HANDSHAKE_ACCEPTED
	USER_STATUS_ECDH_ESTABLISHED
	USER_STATUS_ECDH_CBES_ESTABLISHED
	USER_STATUS_ECDH_CBES_MKLG_ESTABLISHED
)

// func getUserOptsContainer(userStatus int, intercom chan *com.Message) *fyne.Container {
// 	switch {
//   case USER_STATUS_HANDSHAKE_INIT:
//     btn := widget.NewButton("Init Hadnshake", func() {
//       intercom <- &com.Message{
//         ID: ID_CLIENT_ASK_CLIENT_HANDSHAKE,
//         ToID: ,
//       }
//     })
//     return container.NewGridWithRows(1, )
// 	default:
// 		log.Printf("ERROR: getUserOptsContainer: status %d was not handled!", userStatus)
// 		return container.NewGridWithRows(1, widget.NewLabel("ERROR: 500"))
// 	}
// }
