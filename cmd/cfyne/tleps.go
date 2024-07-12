package main

import (
	"errors"
	"log"
	"sync"

	"git.qowevisa.me/Qowevisa/gotell/gmyerr"
	"git.qowevisa.me/Qowevisa/gotell/tlep"
)

var (
	tlepCenter TlepCenter
)

type TlepCenter struct {
	TLEPs map[uint16]*tlep.TLEP
	Mu    sync.Mutex
}

func (t *TlepCenter) Init() {
	t.TLEPs = make(map[uint16]*tlep.TLEP)
}

var (
	ERROR_ALREADY_HAVE = errors.New("Already taken")
	ERROR_DONT_HAVE    = errors.New("Not found")
)

func (t *TlepCenter) AddTLEP(id uint16, name string) error {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	_, alreadyHave := t.TLEPs[id]
	if alreadyHave {
		return ERROR_ALREADY_HAVE
	}
	val, err := tlep.InitTLEP(name)
	if err != nil {
		return gmyerr.WrapPrefix("tlep.InitTLEP", err)
	}
	val.Debug = true
	t.TLEPs[id] = val
	log.Printf("TLEPs: add %p for %d id\n", val, id)
	return nil
}

func (t *TlepCenter) DeleteIfHaveOne(id uint16) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	val, found := t.TLEPs[id]
	if !found {
		log.Printf("TLEP with %d id is not found; Can not delete\n", id)
		return
	}
	delete(t.TLEPs, id)
	log.Printf("TLEP with %v val and %d id was found; TLEP is deleted\n", val, id)
}

func (t *TlepCenter) GetTLEP(id uint16) (*tlep.TLEP, error) {
	log.Printf("Getting tlep by id = %d\n", id)
	name, have := t.TLEPs[id]
	if !have {
		return nil, ERROR_DONT_HAVE
	}
	return name, nil
}
