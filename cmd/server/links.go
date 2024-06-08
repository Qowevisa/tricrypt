package main

import (
	"log"
	"sync"

	com "git.qowevisa.me/Qowevisa/gotell/communication"
)

type UserLink struct {
	LeftNum uint16
	UserID  uint16
}

type LinkArray struct {
	Array []string
}

type LinkCenter struct {
	Links      map[string]*UserLink
	SavedLinks map[uint16]*LinkArray
	Mu         sync.Mutex
}

var (
	linkCenter LinkCenter
)

func (l *LinkCenter) Init() {
	l.Links = make(map[string]*UserLink)
	l.SavedLinks = make(map[uint16]*LinkArray)
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
	val, found := l.SavedLinks[id]
	if !found {
		var tmpAr []string
		tmpAr = append(tmpAr, string(link.Data))
		l.SavedLinks[id] = &LinkArray{
			Array: tmpAr,
		}
	} else {
		val.Array = append(val.Array, string(link.Data))
	}
	log.Printf("Added link by %s\n", string(link.Data))
	log.Printf("SavedLinks[%d] is now %v\n", id, l.SavedLinks[id])
	l.Mu.Unlock()
	return nil
}

func (l *LinkCenter) debug() {
	for val, key := range l.Links {
		if key == nil {
			log.Printf("DEBUG: LINKCENTER: VAL = %s LINK = NIL\n", val)
			continue
		}
		log.Printf("DEBUG: LINKCENTER: VAL = %s LINK = %v\n", val, *key)
	}
}

func (l *LinkCenter) DeleteLink(data []byte) error {
	l.Mu.Lock()
	delete(l.Links, string(data))
	l.Mu.Unlock()
	return nil
}

func (l *LinkCenter) GetLink(data []byte) (*UserLink, error) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	val, found := l.Links[string(data)]
	if !found {
		return nil, ERROR_DONT_HAVE
	}
	return val, nil
}

func (l *LinkCenter) CleanAfterLeave(id uint16) {
	log.Printf("Cleaning after id=%d left;\n", id)
	ar, found := l.SavedLinks[id]
	if !found {
		log.Printf("Cleaning: Id=%d not found;\n", id)
		return
	}
	l.Mu.Lock()
	for _, link := range ar.Array {
		log.Printf("Cleaning: Id=%d deleting %s;\n", id, link)
		delete(l.Links, link)
	}
	l.debug()
	l.Mu.Unlock()
}
