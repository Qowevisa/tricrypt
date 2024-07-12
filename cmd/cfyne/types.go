package main

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"git.qowevisa.me/Qowevisa/gotell/extfyne/widgets"
	// com "git.qowevisa.me/Qowevisa/gotell/communication"
)

type FyneConfig struct {
	App    fyne.App
	Window fyne.Window
}

type MainSceneConfig struct {
	UserNickcname string
	// Just for my calmness
	UserNameSet bool
	UserId      uint16
	// Just for my calmness
	UserIdSet bool
}

type MutableStructAboutTabs struct {
	AppTabs       *container.AppTabs
	LinkTab       *container.TabItem
	LinksBT       string
	LinksNoty     uint
	ConnectionTab *container.TabItem
	ConnsBT       string
	ConnsNoty     uint
	UsersTab      *container.TabItem
	UsersBT       string
	UsersNoty     uint
}

type MutableStructAboutWidgets struct {
	UsersSelect       *widget.Select
	UserOptsCaretaker *fyne.Container
	MessagesSelect    *widget.Select
	MBoardCaretacker  *fyne.Container
	MsgButton         *widget.Button
	MBoardMap         map[string]*widgets.MessageBoard
}

type UserOptions struct {
	OptsContainer *fyne.Container
	Status        int
	ToID          uint16
	Name          string
}

type BundleOfMutexArrays struct {
	LinksMuAr        *MutexLinksArray
	UsersMuAr        *MutexStringArray
	UsersOpts        map[string]*UserOptions
	UserShortcuts    map[uint16]string
	UserShortcutsRev map[string]uint16
	Messages         map[string]*MutexArray[widgets.MBoardMessage]
}

type MutableApp struct {
	App    fyne.App
	Window fyne.Window
}

type MutableApplication struct {
	Tabs        MutableStructAboutTabs
	Widgets     MutableStructAboutWidgets
	ArrayBundle BundleOfMutexArrays
	App         MutableApp
}

// I don't know if it is a good idea
type MutexArray[T any] struct {
	Ar []T
	Mu sync.RWMutex
}

func CreateMutexArray[T any]() MutexArray[T] {
	var tmp []T
	return MutexArray[T]{
		Ar: tmp,
	}
}

func (ma *MutexArray[T]) Add(v T) {
	ma.Mu.Lock()
	ma.Ar = append(ma.Ar, v)
	ma.Mu.Unlock()
}

func (ma *MutexArray[T]) GetArray() []T {
	ma.Mu.RLock()
	defer ma.Mu.RUnlock()
	return ma.Ar
}

// Non-Generic variant
type MutexStringArray struct {
	Ar []string
	Mu sync.RWMutex
}

func CreateMutexStringArray() *MutexStringArray {
	var tmp []string
	return &MutexStringArray{
		Ar: tmp,
	}
}

func (ma *MutexStringArray) Add(v string) {
	ma.Mu.Lock()
	ma.Ar = append(ma.Ar, v)
	ma.Mu.Unlock()
}

func (ma *MutexStringArray) GetArray() []string {
	ma.Mu.RLock()
	defer ma.Mu.RUnlock()
	return ma.Ar
}
