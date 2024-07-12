package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"

	// "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	com "git.qowevisa.me/Qowevisa/gotell/communication"
	"git.qowevisa.me/Qowevisa/gotell/extfyne/layouts"
	"git.qowevisa.me/Qowevisa/gotell/extfyne/widgets"
)

func GetNicknameScene(intercom chan *com.Message) *fyne.Container {
	// button := widget.NewButtonWithIcon(
	// 	"Send",
	// 	theme.ConfirmIcon(),
	// 	getSendMessageFuncToIntercom(intercom, com.ID_CLIENT_SEND_SERVER_NICKNAME, []byte("test")))
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter Nickname:")
	btn := widget.NewButtonWithIcon("Send", theme.MailSendIcon(), func() {
		log.Println("Nickname submitted:", entry.Text)
		intercom <- &com.Message{
			ID:   com.ID_CLIENT_SEND_SERVER_NICKNAME,
			Data: []byte(entry.Text),
		}
	})
	contentInner := container.New(layout.NewVBoxLayout(), entry, btn)
	return container.NewGridWithRows(3,
		layout.NewSpacer(),
		container.NewGridWithColumns(3,
			layout.NewSpacer(),
			contentInner,
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
}

func GetMainScene(intercom chan *com.Message, cfg MainSceneConfig) (*fyne.Container, *MutableApplication) {
	headerWidget := widget.NewLabel(fmt.Sprintf("Hello, %s! Your ID: %d", cfg.UserNickcname, cfg.UserId))
	header := container.NewGridWithColumns(3,
		layout.NewSpacer(),
		headerWidget,
		layout.NewSpacer(),
	)

	linksSS, linksArr := GetLinksSubScene(intercom, header)
	linksTab := container.NewTabItemWithIcon("Links", theme.MailComposeIcon(), linksSS)
	conSS, conBundle, conWidgets := GetConnectionsSubScene(intercom, header)
	connsTab := container.NewTabItemWithIcon("Connections", theme.ComputerIcon(), conSS)
	userSS, userBundle, userWidgets := GetUsersSubScene(intercom, header, conBundle.UserShortcutsRev)
	usersTab := container.NewTabItemWithIcon("Users", theme.AccountIcon(), userSS)
	tabs := container.NewAppTabs(
		linksTab,
		connsTab,
		usersTab,
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	return container.NewGridWithColumns(1,
			tabs),
		&MutableApplication{
			Tabs: MutableStructAboutTabs{
				AppTabs:       tabs,
				LinkTab:       linksTab,
				LinksBT:       "Links",
				LinksNoty:     0,
				ConnectionTab: connsTab,
				ConnsBT:       "Connections",
				ConnsNoty:     0,
				UsersTab:      usersTab,
				UsersBT:       "Users",
				UsersNoty:     0,
			},
			ArrayBundle: BundleOfMutexArrays{
				LinksMuAr:        linksArr,
				UsersMuAr:        conBundle.UsersMuAr,
				UsersOpts:        conBundle.UsersOpts,
				UserShortcuts:    conBundle.UserShortcuts,
				UserShortcutsRev: conBundle.UserShortcutsRev,
				Messages:         userBundle.Messages,
			},
			Widgets: MutableStructAboutWidgets{
				UsersSelect:       conWidgets.UsersSelect,
				UserOptsCaretaker: conWidgets.UserOptsCaretaker,
				MessagesSelect:    userWidgets.MessagesSelect,
				MBoardCaretacker:  userWidgets.MBoardCaretacker,
				MsgButton:         userWidgets.MsgButton,
				MBoardMap:         userWidgets.MBoardMap,
			},
		}

}

func GetLinksSubScene(intercom chan *com.Message, header *fyne.Container) (*fyne.Container, *MutexLinksArray) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Count:")
	entry.Validator = func(s string) error {
		val, err := strconv.ParseUint(s, 10, 16)
		if err == nil && val == 0 {
			return errors.New("Use count can't be 0")
		}
		return err
	}
	btn := widget.NewButtonWithIcon("Generate", theme.MailSendIcon(), func() {
		count, err := strconv.ParseUint(entry.Text, 10, 16)
		if err != nil || count == 0 {
			log.Printf("dafuq: GetLinksSubScene: %v ; %d\n", err, count)
			return
		}
		intercom <- &com.Message{
			ID:   com.ID_CLIENT_SEND_SERVER_LINK,
			ToID: uint16(count),
		}
	})
	entry.SetOnValidationChanged(func(err error) {
		if err != nil {
			btn.Hide()
		} else {
			if btn.Hidden {
				btn.Show()
			}
		}
	})
	contentInner := container.NewGridWithColumns(2, entry, btn)
	// /
	linksMutAt := CreateMutexLinksArray()
	linksList := widget.NewList(func() int {
		return len(linksMutAt.GetArray())
	}, func() fyne.CanvasObject {
		return widget.NewLabel("template")
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
		link := linksMutAt.Ar[lii]
		str := fmt.Sprintf("Count: %d ; %s", link.UseCount, link.Data)
		co.(*widget.Label).SetText(str)
	},
	)
	linksList.OnSelected = func(id widget.ListItemID) {
		link := linksMutAt.Ar[id]
		globalClipboardChannel <- string(link.Data)
	}
	// //
	linksGetEntry := widget.NewEntry()
	linksGetEntry.SetPlaceHolder("Link:")
	linksGetEntry.Validator = func(s string) error {
		validated, err := com.IsThisALinkData(s)
		if err != nil {
			return err
		}
		if !validated {
			return errors.New("Link is not validated")
		}
		return nil
	}
	linksGetBtn := widget.NewButtonWithIcon("Get", theme.MailReplyIcon(), func() {
		link := linksGetEntry.Text
		validated, err := com.IsThisALinkData(link)
		if err != nil {
			log.Printf("linksGetBtn: Error: %v\n", err)
			return
		}
		if !validated {
			log.Printf("linksGetBtn: Error: link %v not validated\n", link)
			return
		}
		intercom <- &com.Message{
			ID:   com.ID_CLIENT_ASK_SERVER_LINK,
			Data: []byte(link),
		}
	})
	linksGetEntry.SetOnValidationChanged(func(err error) {
		if err != nil {
			linksGetBtn.Hide()
		} else {
			if linksGetBtn.Hidden {
				linksGetBtn.Show()
			}
		}
	})
	contentInner2 := container.NewGridWithColumns(2, linksGetEntry, linksGetBtn)
	// /
	// //
	return container.NewGridWithRows(3,
		header,
		container.New(
			layouts.NewVariableGridWithColumns(
				6, []int{2, 1, 3}),
			container.NewGridWithRows(3,
				container.NewGridWithColumns(3,
					layout.NewSpacer(),
					widget.NewLabel("Create link:"),
					layout.NewSpacer()),
				contentInner,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
			linksList,
		),
		container.New(
			layouts.NewVariableGridWithColumns(
				6, []int{2, 4}),
			container.NewGridWithRows(3,
				container.NewGridWithColumns(3,
					layout.NewSpacer(),
					widget.NewLabel("Get user from link:"),
					layout.NewSpacer()),
				contentInner2,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		),
	), linksMutAt
}

func GetConnectionsSubScene(intercom chan *com.Message, header *fyne.Container) (*fyne.Container, BundleOfMutexArrays, MutableStructAboutWidgets) {
	optionsHeaderLabel := widget.NewLabel("Options for user: ")
	usersMuAr := CreateMutexStringArray()
	usersOptions := make(map[string]*UserOptions)
	userShortcuts := make(map[uint16]string)
	userShortcutsRev := make(map[string]uint16)
	userOptsCaretaker := container.NewGridWithRows(1)
	optsContainer := container.NewGridWithRows(2,
		optionsHeaderLabel,
		userOptsCaretaker,
	)
	userSelect := widget.NewSelect(usersMuAr.Ar, func(s string) {
		optionsHeaderLabel.SetText(fmt.Sprintf("Options for user: %s", s))
		uOpts, exists := usersOptions[s]
		if !exists {
			log.Printf("GetConnectionsSubScene::1 : TODO!!\n")
			// TODO
			return
		}
		userOptsCaretaker.RemoveAll()
		userOptsCaretaker.Add(uOpts.OptsContainer)
		optsContainer.Refresh()
	})

	contentInner := container.New(layout.NewVBoxLayout(), userSelect, optsContainer)
	return container.New(
			layout.NewVBoxLayout(),
			header,
			container.New(
				layout.NewHBoxLayout(),
				layout.NewSpacer(),
				contentInner,
				layout.NewSpacer(),
			),
		), BundleOfMutexArrays{
			UsersMuAr:        usersMuAr,
			UsersOpts:        usersOptions,
			UserShortcuts:    userShortcuts,
			UserShortcutsRev: userShortcutsRev,
		}, MutableStructAboutWidgets{
			UsersSelect:       userSelect,
			UserOptsCaretaker: userOptsCaretaker,
		}
}

func GetUsersSubScene(intercom chan *com.Message, header *fyne.Container, userShortcutsRev map[string]uint16) (*fyne.Container, BundleOfMutexArrays, MutableStructAboutWidgets) {
	msgHeader := widget.NewLabel("Messages With User: ")
	usersMessages := make(map[string]*MutexArray[widgets.MBoardMessage])
	usersMBoards := make(map[string]*widgets.MessageBoard)
	msgEntry := widget.NewEntry()
	msgEntry.SetPlaceHolder("Message...")
	userMsgSendBtn := widget.NewButtonWithIcon("Send", theme.MailSendIcon(), func() {})
	mBoardCareTaker := container.NewGridWithRows(1)
	userMsgsSelect := widget.NewSelect([]string{}, func(s string) {
		msgHeader.SetText(fmt.Sprintf("Messages With User: %s", s))
		msgsMuAr, exists := usersMessages[s]
		if !exists {
			tmp := CreateMutexArray[widgets.MBoardMessage]()
			usersMessages[s] = &tmp
			msgsMuAr = &tmp
			return
		}
		mBoardCareTaker.RemoveAll()
		mmBoard := widgets.NewMessageBoard(msgsMuAr.Ar)
		usersMBoards[s] = mmBoard
		mBoardCareTaker.Add(mmBoard)
		mBoardCareTaker.Refresh()
		id, exists := userShortcutsRev[s]
		if !exists {
			log.Printf("GetUsersSubScene::2 : TODO!!\n")
			// TODO
			return
		}
		userMsgSendBtn.Text = fmt.Sprintf("%s", s)
		userMsgSendBtn.OnTapped = func() {
			intercom <- &com.Message{
				ID:   com.ID_CLIENT_SEND_CLIENT_MESSAGE,
				ToID: id,
				Data: []byte(msgEntry.Text),
			}
			msg := widgets.MBoardMessage{
				LeftAlign: false,
				Data:      widget.NewLabel(msgEntry.Text),
			}
			msgsMuAr.Add(msg)
			mmBoard.Add(msg)
			msgEntry.Text = ""
			msgEntry.Refresh()
			mBoardCareTaker.Refresh()
			mmBoard.Refresh()
		}
		userMsgSendBtn.Refresh()
	})

	// entryWithButton := container.NewHBox(
	// 	layout.NewSpacer(),
	// 	container.NewPadded(msgEntry),
	// 	userMsgSendBtn,
	// )
	entryWithButton := container.New(
		layouts.NewEntryBtn7030(),
		msgEntry,
		userMsgSendBtn,
	)
	fs := globalCfg.Window.Canvas().Size()

	scene := container.New(
		layout.NewVBoxLayout(),
		header,
		container.NewGridWithColumns(2,
			container.New(
				layout.NewVBoxLayout(),
				layout.NewSpacer(),
				container.New(
					layout.NewHBoxLayout(),
					userMsgsSelect,
					msgHeader,
				),
			),
			container.New(
				layout.NewVBoxLayout(),
				container.New(
					layouts.NewFullWidthWithSize(fyne.NewSize(0, fs.Height*0.8)),
					mBoardCareTaker,
				),
				layout.NewSpacer(),
				entryWithButton,
			),
		),
	)
	// container.New(
	// 	layouts.NewVariableGridWithRows(3, []int{1, 2}),
	// 	header,
	// 	container.New(
	// 		layouts.NewVariableGridWithColumns(4, []int{1, 3}),
	// 		container.NewGridWithColumns(3,
	// 			layout.NewSpacer(),
	// 			container.NewVBox(userMsgsSelect, msgHeader),
	// 			layout.NewSpacer()),
	// 		container.New(
	// 			layouts.NewVariableGridWithRows(6, []int{5, 1}),
	// 			mBoardCareTaker,
	// 			container.NewHBox(msgEntry, userMsgSendBtn),
	// 		),
	// 	),
	// )
	return scene, BundleOfMutexArrays{
			Messages: usersMessages,
		}, MutableStructAboutWidgets{
			MessagesSelect:   userMsgsSelect,
			MBoardCaretacker: mBoardCareTaker,
			MsgButton:        userMsgSendBtn,
			MBoardMap:        usersMBoards,
		}
}
