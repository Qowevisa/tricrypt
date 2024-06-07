package main

import (
	"sync"

	com "git.qowevisa.me/Qowevisa/gotell/communication"
)

type UserLink struct {
	LeftNum uint16
	UserID  uint16
}

type LinkCenter struct {
	Links map[string]*UserLink
	Mu    sync.Mutex
}

var (
	linkCenter LinkCenter
)

func (l *LinkCenter) Init() {
	l.Links = make(map[string]*UserLink)
}

func (l *LinkCenter) AddLink(id uint16, link com.Link) error {
	_, found := l.Links[string(link.Data)]
	if found {
		return ERROR_ALREADY_HAVE
	}
	l.Mu.Lock()
	l.Links[string(link.Data)] = &UserLink{
		LeftNum: link.UseCount,
		UserID:  id,
	}
	l.Mu.Unlock()
	return nil
}

func (l *LinkCenter) DeleteLink(data []byte) error {
	l.Mu.Lock()
	delete(l.Links, string(data))
	l.Mu.Unlock()
	return nil
}

func (l *LinkCenter) GetLink(data []byte) (*UserLink, error) {
	val, found := l.Links[string(data)]
	if !found {
		return nil, ERROR_DONT_HAVE
	}
	return val, nil
}
