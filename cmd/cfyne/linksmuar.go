package main

import (
	com "git.qowevisa.me/Qowevisa/gotell/communication"
	"sync"
)

// Non-Generic variant
type MutexLinksArray struct {
	Ar []com.Link
	Mu sync.RWMutex
}

func CreateMutexLinksArray() *MutexLinksArray {
	var tmp []com.Link
	return &MutexLinksArray{
		Ar: tmp,
	}
}

func (ma *MutexLinksArray) Add(v com.Link) {
	ma.Mu.Lock()
	ma.Ar = append(ma.Ar, v)
	ma.Mu.Unlock()
}

func (ma *MutexLinksArray) GetArray() []com.Link {
	ma.Mu.RLock()
	defer ma.Mu.RUnlock()
	return ma.Ar
}
