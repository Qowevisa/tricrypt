package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	com "git.qowevisa.me/Qowevisa/gotell/communication"
	"git.qowevisa.me/Qowevisa/gotell/env"
	"git.qowevisa.me/Qowevisa/gotell/extfyne/widgets"
)

func updateTime(clock *widget.Label) {
	formatted := time.Now().Format("Time: 03:04:05")
	clock.SetText(formatted)
}

var globalApp *MutableApplication
var globalCfg *FyneConfig

var (
	ecdhKeyReceived = false
	ecdhKeySent     = false
)

var (
	globalClipboardChannel = make(chan string)
)

func main() {
	tlepCenter.Init()
	userCenter.Init()
	loadingFileName := env.ServerFullchainFileName
	cert, err := os.ReadFile(loadingFileName)
	if err != nil {
		log.Fatalf("client: load root cert: %s", err)
	}
	log.Printf("Certificate %s loaded successfully!\n", loadingFileName)
	//
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("client: failed to parse root certificate")
	}

	config := &tls.Config{
		RootCAs: roots,
	}
	host := "chat.qowevisa.click"
	port := 3232
	service := fmt.Sprintf("%s:%d", host, port)
	log.Printf("client: connecting to %s", service)
	conn, err := tls.Dial("tcp", service, config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Printf("client: connected to %s", service)
	//

	a := app.New()
	w := a.NewWindow("gotell-fyne")
	// w.Resize(fyne.NewSize(640, 400))

	// dispatch clipboardChannel
	go func(w fyne.Window) {
		for val := range globalClipboardChannel {
			w.Clipboard().SetContent(val)
		}
	}(w)

	clock := widget.NewLabel("")
	go func() {
		for range time.Tick(time.Second) {
			updateTime(clock)
		}
	}()

	btn := widget.NewButtonWithIcon("testBtn", theme.AccountIcon(), func() {
		fmt.Println("I'm pressed!")
	})

	mainContainer := container.New(layout.NewVBoxLayout(), clock, btn)

	w.SetContent(mainContainer)
	// w.SetContent(widget.NewLabel("Hello World!"))
	// w2 := a.NewWindow("Larger")
	// w2.SetContent(widget.NewLabel("More content"))
	// w2.Resize(fyne.NewSize(100, 100))
	// w2.Show()
	intercomFromServer := make(chan *com.Message, 16)
	intercomToServer := make(chan *com.Message, 16)
	actionsChannel := make(chan int)
	cfg := FyneConfig{
		App:    a,
		Window: w,
	}
	globalCfg = &cfg

	go readFromServer(conn, intercomFromServer)
	go analyzeMessages(intercomFromServer, actionsChannel)
	go actOnActions(cfg, actionsChannel, intercomToServer)
	go readFromIntercom(intercomToServer, conn)
	w.ShowAndRun()
}

var r com.RegisteredUser
var tmpLink *com.Link
var tmpNick string

func readFromServer(conn net.Conn, intercom chan *com.Message) {
	buf := make([]byte, 70000)
	for {
		err := conn.SetDeadline(time.Now().Add(1 * time.Minute))
		if err != nil {
			log.Printf("ERROR: conn.SetDeadline: %v", err)
			continue
		}
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("ERROR: client: read: %s", err)
			if errors.Is(err, net.ErrClosed) {
				log.Printf("caught ErrClosed!")
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("caught! Read timeout occurred")
				continue
			}
			return
		}
		msg, err := com.Decode(buf[:n])
		if err != nil {
			log.Printf("ERROR: %#v\n", err)
			continue
		}
		if msg == nil {
			continue
		}
		log.Printf("client: readServer: received message from server: %v", *msg)
		switch msg.ID {
		case com.ID_SERVER_APPROVE_CLIENT_NICKNAME:
			newID := binary.BigEndian.Uint16(msg.Data)
			msg.FromID = newID
			msg.Data = []byte{}
			r.ID = newID
			if tmpNick != "" {
				r.Name = tmpNick
			}
			r.IsRegistered = true
		case com.ID_SERVER_APPROVE_CLIENT_LINK:
			if tmpLink == nil {
				continue
			}
			msg.ToID = tmpLink.UseCount
			msg.Data = tmpLink.Data
		// Crypto stuff
		case com.ID_CLIENT_SEND_CLIENT_ECDH_PUBKEY:
			t, err := tlepCenter.GetTLEP(msg.FromID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			err = t.ECDHApplyOtherKeyBytes(msg.Data)
			if err != nil {
				log.Printf("ERROR: tlep: ECDHApplyOtherKeyBytes: %v\n", err)
				continue
			}
			fromName, err := userCenter.GetName(msg.FromID)
			if err != nil {
				log.Printf("ERROR: userCenter: GetName: %v\n", err)
			} else {
				msg.Data = []byte(fromName)
			}
		case com.ID_CLIENT_SEND_CLIENT_CBES_SPECS:
			t, err := tlepCenter.GetTLEP(msg.FromID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			cbes, err := t.DecryptMessageEA(msg.Data)
			if err != nil {
				log.Printf("ERROR: tlep: DecryptMessageEA: %v\n", err)
				continue
			}
			err = t.CBESSetFromBytes(cbes)
			if err != nil {
				log.Printf("ERROR: tlep: CBESSetFromBytes: %v\n", err)
				continue
			}
			fromName, err := userCenter.GetName(msg.FromID)
			if err != nil {
				log.Printf("ERROR: userCenter: GetName: %v\n", err)
			} else {
				msg.Data = []byte(fromName)
			}
			// message
		case com.ID_CLIENT_SEND_CLIENT_MESSAGE:
			t, err := tlepCenter.GetTLEP(msg.FromID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			decrypedMsg, err := t.DecryptMessageAtMax(msg.Data)
			if err != nil {
				log.Printf("ERROR: tlep: DecryptMessageAtMax: %v\n", err)
				continue
			}
			msg.Data = decrypedMsg
			// switch
		}
		// user stuff
		switch msg.ID {
		case com.ID_CLIENT_ASK_CLIENT_HANDSHAKE,
			com.ID_CLIENT_APPROVE_CLIENT_HANDSHAKE:
			userCenter.AddUser(string(msg.Data), msg.FromID)
		}
		log.Printf("client: readServer: sending message to websocket: %v", *msg)
		intercom <- msg
	}
}

func analyzeMessages(intercomFromServer chan *com.Message, actionsChannel chan int) {
	for msg := range intercomFromServer {
		if msg == nil {
			log.Printf("ERROR: msg is nil")
			continue
		}
		log.Printf("Handling %v from server", *msg)
		switch msg.ID {
		case com.ID_SERVER_ASK_CLIENT_NICKNAME:
			actionsChannel <- ACTION_SET_SCENE_NICKNAME
		case com.ID_SERVER_APPROVE_CLIENT_NICKNAME:
			actionsChannel <- ACTION_SET_SCENE_MAIN
		case com.ID_SERVER_DECLINE_CLIENT_NICKNAME:
			// TODO
		// Link stuff
		case com.ID_SERVER_APPROVE_CLIENT_LINK:
			globalApp.ArrayBundle.LinksMuAr.Add(*tmpLink)
			globalApp.Tabs.LinkTab.Content.Refresh()
		case com.ID_SERVER_SEND_CLIENT_ANOTHER_ID:
			selTab := globalApp.Tabs.AppTabs.Selected()
			id := msg.ToID
			userStr := getStdUserName(id, string(msg.Data))
			globalApp.Widgets.UsersSelect.Options = append(globalApp.Widgets.UsersSelect.Options, userStr)
			globalApp.Tabs.ConnsNoty++
			globalApp.ArrayBundle.UsersOpts[userStr] = &UserOptions{
				Status: USER_STATUS_VISIBLE,
				ToID:   id,
			}
			_, exists := globalApp.ArrayBundle.UserShortcuts[id]
			if !exists {
				globalApp.ArrayBundle.UserShortcuts[id] = userStr
				globalApp.ArrayBundle.UserShortcutsRev[userStr] = id
			}
			if selTab != globalApp.Tabs.ConnectionTab {
				updateTabsBasedOnNotyVals(&globalApp.Tabs)
			}
			actionsChannel <- ACTION_UPDATE_USER_OPTS_CONTS
		case com.ID_CLIENT_ASK_CLIENT_HANDSHAKE:
			selTab := globalApp.Tabs.AppTabs.Selected()
			id := msg.FromID
			userStr := getStdUserName(id, string(msg.Data))
			globalApp.Widgets.UsersSelect.Options = append(globalApp.Widgets.UsersSelect.Options, userStr)
			globalApp.Tabs.ConnsNoty++
			globalApp.ArrayBundle.UsersOpts[userStr] = &UserOptions{
				Status: USER_STATUS_HANDSHAKE_INIT,
				ToID:   id,
			}
			_, exists := globalApp.ArrayBundle.UserShortcuts[id]
			if !exists {
				globalApp.ArrayBundle.UserShortcuts[id] = userStr
				globalApp.ArrayBundle.UserShortcutsRev[userStr] = id
			}
			if selTab != globalApp.Tabs.ConnectionTab {
				updateTabsBasedOnNotyVals(&globalApp.Tabs)
			}
			actionsChannel <- ACTION_UPDATE_USER_OPTS_CONTS
		case com.ID_CLIENT_APPROVE_CLIENT_HANDSHAKE:
			id := msg.FromID
			userStr := getStdUserName(id, string(msg.Data))
			var shortcut string
			val, exists := globalApp.ArrayBundle.UserShortcuts[id]
			if !exists {
				globalApp.ArrayBundle.UserShortcuts[id] = userStr
				globalApp.ArrayBundle.UserShortcutsRev[userStr] = id
				shortcut = userStr
			} else {
				shortcut = val
			}
			u, exists := globalApp.ArrayBundle.UsersOpts[shortcut]
			if !exists {
				log.Printf("WARNING!!: UserOpts[%s] NOT exists!", shortcut)
				continue
			}
			u.Status = USER_STATUS_HANDSHAKE_ACCEPTED
			actionsChannel <- ACTION_UPDATE_USER_OPTS_CONTS
		case com.ID_CLIENT_SEND_CLIENT_ECDH_PUBKEY:
			id := msg.FromID
			if ecdhKeySent {
				val, exists := globalApp.ArrayBundle.UserShortcuts[id]
				if !exists {
					log.Printf("WARNING!!: UserShortcuts[%d] NOT exists!", id)
					continue
				}
				u, exists := globalApp.ArrayBundle.UsersOpts[val]
				if !exists {
					log.Printf("WARNING!!: UserOpts[%s] NOT exists!", val)
					continue
				}
				u.Status = USER_STATUS_ECDH_ESTABLISHED
				actionsChannel <- ACTION_UPDATE_USER_OPTS_CONTS
			}
			ecdhKeyReceived = true
		case com.ID_CLIENT_SEND_CLIENT_CBES_SPECS:
			id := msg.FromID
			val, exists := globalApp.ArrayBundle.UserShortcuts[id]
			if !exists {
				log.Printf("WARNING!!: UserShortcuts[%d] NOT exists!", id)
				continue
			}
			u, exists := globalApp.ArrayBundle.UsersOpts[val]
			if !exists {
				log.Printf("WARNING!!: UserOpts[%s] NOT exists!", val)
				continue
			}
			u.Status = USER_STATUS_ECDH_CBES_ESTABLISHED
			actionsChannel <- ACTION_UPDATE_USER_OPTS_CONTS
			//
			log.Printf("Checkin Messages tab")
			userStr := getStdUserName(id, string(msg.Data))
			var shortcut string
			if !exists {
				globalApp.ArrayBundle.UserShortcuts[id] = userStr
				shortcut = userStr
			} else {
				shortcut = val
			}
			//
			haveUser := false
			for _, v := range globalApp.Widgets.MessagesSelect.Options {
				if v == shortcut {
					haveUser = true
					break
				}
			}
			log.Printf("MessagesSelect.Options1 are %v\n", globalApp.Widgets.MessagesSelect.Options)
			log.Printf("have user is : %t\n", haveUser)
			if !haveUser {
				log.Printf("MessagesSelect.Options1.5 shortcut is '%s'\n", shortcut)
				// LOL wtf
				// globalApp.Widgets.MessagesSelect.Options = append(globalApp.Widgets.UsersSelect.Options, shortcut)
				globalApp.Widgets.MessagesSelect.Options = append(globalApp.Widgets.MessagesSelect.Options, shortcut)
				log.Printf("MessagesSelect.Options2 are %v\n", globalApp.Widgets.MessagesSelect.Options)
				globalApp.Tabs.UsersNoty++
				globalApp.Widgets.MessagesSelect.Refresh()
				updateTabsBasedOnNotyVals(&globalApp.Tabs)
			}
		case com.ID_CLIENT_SEND_CLIENT_MESSAGE:
			id := msg.FromID
			shortcut, exists := globalApp.ArrayBundle.UserShortcuts[id]
			if !exists {
				log.Printf("WARNING: UserShortcuts[%d] NOT exists!", id)
				continue
			}
			msgs, exists := globalApp.ArrayBundle.Messages[shortcut]
			mBoardMsg := widgets.MBoardMessage{
				LeftAlign: true,
				Data:      widget.NewLabel(string(msg.Data)),
			}
			if !exists {
				tmp := CreateMutexArray[widgets.MBoardMessage]()
				tmp.Add(mBoardMsg)
				globalApp.ArrayBundle.Messages[shortcut] = &tmp
				globalApp.Widgets.MBoardCaretacker.Refresh()
				globalApp.Tabs.UsersTab.Content.Refresh()
			} else {
				msgs.Add(mBoardMsg)
				globalApp.Widgets.MBoardCaretacker.Refresh()
				globalApp.Tabs.UsersTab.Content.Refresh()
			}
			haveUser := false
			for _, v := range globalApp.Widgets.MessagesSelect.Options {
				if v == shortcut {
					haveUser = true
					break
				}
			}
			if !haveUser {
				globalApp.Widgets.MessagesSelect.Options = append(globalApp.Widgets.MessagesSelect.Options, shortcut)
				globalApp.Tabs.UsersNoty++
				updateTabsBasedOnNotyVals(&globalApp.Tabs)
				globalApp.Widgets.MessagesSelect.Refresh()
			}
			//
			mBoard, exists := globalApp.Widgets.MBoardMap[shortcut]
			if exists {
				mBoard.Add(mBoardMsg)
			}
			// TODO
		default:
			log.Printf("Can not handle msg with id %d", msg.ID)
			log.Printf("Msg is different : %v", *msg)
		}
	}
}

func actOnActions(cfg FyneConfig, actionsChannel chan int, intercomToServer chan *com.Message) {
	for action := range actionsChannel {
		log.Printf("Receive action %d", action)
		switch action {
		case ACTION_SET_SCENE_NICKNAME:
			// cfg.Window.Content()
			scene := GetNicknameScene(intercomToServer)
			cfg.Window.Resize(fyne.NewSize(640, 400))
			cfg.Window.SetContent(scene)
		case ACTION_SET_SCENE_MAIN:
			scene, app := GetMainScene(intercomToServer, MainSceneConfig{
				UserNickcname: tmpNick,
				UserNameSet:   true,
				UserId:        r.ID,
				UserIdSet:     true,
			})
			globalApp = app
			globalApp.App.App = cfg.App
			globalApp.App.Window = cfg.Window
			// shortcuts
			ctrl1 := &desktop.CustomShortcut{KeyName: fyne.Key1, Modifier: fyne.KeyModifierControl}
			ctrl2 := &desktop.CustomShortcut{KeyName: fyne.Key2, Modifier: fyne.KeyModifierControl}
			ctrl3 := &desktop.CustomShortcut{KeyName: fyne.Key3, Modifier: fyne.KeyModifierControl}
			cfg.Window.Canvas().AddShortcut(ctrl1, func(shortcut fyne.Shortcut) {
				app.Tabs.AppTabs.Select(app.Tabs.LinkTab)
			})
			cfg.Window.Canvas().AddShortcut(ctrl2, func(shortcut fyne.Shortcut) {
				app.Tabs.AppTabs.Select(app.Tabs.ConnectionTab)
			})
			cfg.Window.Canvas().AddShortcut(ctrl3, func(shortcut fyne.Shortcut) {
				app.Tabs.AppTabs.Select(app.Tabs.UsersTab)
			})
			// cfg.Window.Resize(fyne.NewSize(640, 400))
			cfg.Window.SetContent(scene)
			// ACTION_UPDATE_USER_OPTS_CONTS
		case ACTION_UPDATE_USER_OPTS_CONTS:
			for _, u := range globalApp.ArrayBundle.UsersOpts {
				switch u.Status {
				case USER_STATUS_VISIBLE:
					btn := widget.NewButton("Init Handhake", func() {})
					btn.OnTapped = func() {
						intercomToServer <- &com.Message{
							ID:   com.ID_CLIENT_ASK_CLIENT_HANDSHAKE,
							ToID: u.ToID,
						}
						btn.Hide()
					}
					u.OptsContainer = container.NewGridWithRows(1, btn)
					cTaker := globalApp.Widgets.UserOptsCaretaker
					cTaker.RemoveAll()
					cTaker.Add(u.OptsContainer)
				case USER_STATUS_HANDSHAKE_INIT:
					btn := widget.NewButton("Accept Handhake", func() {})
					btn.OnTapped = func() {
						intercomToServer <- &com.Message{
							ID:   com.ID_CLIENT_APPROVE_CLIENT_HANDSHAKE,
							ToID: u.ToID,
						}
						btn.Hide()
						u.Status = USER_STATUS_HANDSHAKE_ACCEPTED
						actionsChannel <- ACTION_UPDATE_USER_OPTS_CONTS
					}
					u.OptsContainer = container.NewGridWithRows(1, btn)
					cTaker := globalApp.Widgets.UserOptsCaretaker
					cTaker.RemoveAll()
					cTaker.Add(u.OptsContainer)
				case USER_STATUS_HANDSHAKE_ACCEPTED:
					btn := widget.NewButton("Send ECDH PubKey", func() {})
					btn.OnTapped = func() {
						intercomToServer <- &com.Message{
							ID:   com.ID_CLIENT_SEND_CLIENT_ECDH_PUBKEY,
							ToID: u.ToID,
						}
						btn.Hide()
						if ecdhKeyReceived {
							u.Status = USER_STATUS_ECDH_ESTABLISHED
							actionsChannel <- ACTION_UPDATE_USER_OPTS_CONTS
						}
						ecdhKeySent = true
					}
					u.OptsContainer = container.NewGridWithRows(1, btn)
					cTaker := globalApp.Widgets.UserOptsCaretaker
					cTaker.RemoveAll()
					cTaker.Add(u.OptsContainer)
				case USER_STATUS_ECDH_ESTABLISHED:
					btn := widget.NewButton("Send CBES Specs", func() {})
					btn.OnTapped = func() {
						intercomToServer <- &com.Message{
							ID:   com.ID_CLIENT_SEND_CLIENT_CBES_SPECS,
							ToID: u.ToID,
						}
						btn.Hide()
					}
					u.OptsContainer = container.NewGridWithRows(1, btn)
					cTaker := globalApp.Widgets.UserOptsCaretaker
					cTaker.RemoveAll()
					cTaker.Add(u.OptsContainer)
				case USER_STATUS_ECDH_CBES_ESTABLISHED:
					btn := widget.NewButton("TBC", func() {})
					btn.OnTapped = func() {
						log.Println("TODO:: 1001")
					}
					u.OptsContainer = container.NewGridWithRows(1, btn)
					cTaker := globalApp.Widgets.UserOptsCaretaker
					cTaker.RemoveAll()
					cTaker.Add(u.OptsContainer)
				default:
				}
			}
			// globalApp.Tabs.ConnectionTab.Content.Refresh()
			// ACTION_UPDATE_USER_OPTS_CONTS
		default:
		}
	}
}

func readFromIntercom(intercom chan *com.Message, conn net.Conn) {
	for msg := range intercom {
		log.Printf("client: received message from Intercom: %v", msg)
		msg.Version = com.V1
		switch msg.ID {
		case com.ID_CLIENT_SEND_SERVER_NICKNAME:
			tmpNick = string(msg.Data)
		case com.ID_CLIENT_SEND_SERVER_LINK:
			if !r.IsRegistered {
				continue
			}
			l, err := r.GenerateLink(msg.ToID)
			if err != nil {
				log.Printf("Error: link: %v", err)
				continue
			}
			tmpLink = &l
			answ, err := com.ClientSendServerLink(r.ID, l)
			if err != nil {
				log.Printf("Error: com: %v", err)
				continue
			}
			log.Printf("client: readWS: sending data to server: %v", answ)
			conn.Write(answ)
			continue
		case com.ID_CLIENT_ASK_CLIENT_HANDSHAKE,
			com.ID_CLIENT_APPROVE_CLIENT_HANDSHAKE,
			com.ID_CLIENT_DECLINE_CLIENT_HANDSHAKE,
			com.ID_CLIENT_SEND_CLIENT_ECDH_PUBKEY,
			com.ID_CLIENT_SEND_CLIENT_CBES_SPECS,
			com.ID_CLIENT_SEND_CLIENT_MKLG_FINGERPRINT,
			com.ID_CLIENT_DECLINE_CLIENT_MKLG_FINGERPRINT,
			com.ID_CLIENT_SEND_CLIENT_MESSAGE:
			if !r.IsRegistered {
				continue
			}
			msg.FromID = r.ID
			// switch
		}
		// user stuff
		switch msg.ID {
		case com.ID_CLIENT_ASK_CLIENT_HANDSHAKE,
			com.ID_CLIENT_APPROVE_CLIENT_HANDSHAKE:
			if r.IsRegistered {
				msg.Data = []byte(r.Name)
			}
		}
		// Crypto stuff
		switch msg.ID {
		case com.ID_CLIENT_ASK_CLIENT_HANDSHAKE,
			com.ID_CLIENT_APPROVE_CLIENT_HANDSHAKE:
			err := tlepCenter.AddTLEP(msg.ToID, fmt.Sprintf("%s-%d", r.Name, msg.ToID))
			if err != nil {
				log.Printf("ERROR: tlepCenter.AddUser: %v\n", err)
			}
		case com.ID_CLIENT_SEND_CLIENT_ECDH_PUBKEY:
			t, err := tlepCenter.GetTLEP(msg.ToID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			key, err := t.ECDHGetPublicKey()
			if err != nil {
				log.Printf("ERROR: tlep: ECDHGetPublicKey: %v\n", err)
				continue
			}
			msg.Data = key
		case com.ID_CLIENT_SEND_CLIENT_CBES_SPECS:
			t, err := tlepCenter.GetTLEP(msg.ToID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			err = t.CBESInitRandom()
			if err != nil {
				log.Printf("ERROR: tlep: CBESInitRandom: %v\n", err)
				continue
			}
			cbes, err := t.CBESGetBytes()
			if err != nil {
				log.Printf("ERROR: tlep: ECDHGetPublicKey: %v\n", err)
				continue
			}
			cbesEAEncr, err := t.EncryptMessageEA(cbes)
			if err != nil {
				log.Printf("ERROR: tlep: EncryptMessageEA: %v\n", err)
				continue
			}
			msg.Data = cbesEAEncr
			// message
		case com.ID_CLIENT_SEND_CLIENT_MESSAGE:
			t, err := tlepCenter.GetTLEP(msg.ToID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			encrypedMsg, err := t.EncryptMessageAtMax(msg.Data)
			if err != nil {
				log.Printf("ERROR: tlep: EncryptMessageAtMax: %v\n", err)
				continue
			}
			msg.Data = encrypedMsg
			// switch
		}
		encodedMsg, err := msg.Bytes()
		if err != nil {
			log.Printf("Encoding error: %s", err)
			continue
		}
		log.Printf("client: readWS: sending data to server: %v", encodedMsg)
		conn.Write(encodedMsg)
	}
}
