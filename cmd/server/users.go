package main

import (
	"errors"
	"math/rand"
	"sync"
)

var (
	userCenter UserCenter
)

type UserCenter struct {
	UsersSTOI map[string]uint16
	UsersITOS map[uint16]string
	Mu        sync.Mutex
}

func (u *UserCenter) Init() {
	u.UsersSTOI = make(map[string]uint16)
	u.UsersITOS = make(map[uint16]string)
}

var (
	ERROR_ALREADY_HAVE = errors.New("Username is already taken")
	ERROR_DONT_HAVE    = errors.New("Username was not found")
)

func (u *UserCenter) AddUser(name string) (uint16, error) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	_, alreadyHave := u.UsersSTOI[name]
	if alreadyHave {
		return 0, ERROR_ALREADY_HAVE
	}
	var ret uint16
	for {
		rInt := rand.Int31()
		_, taken := u.UsersITOS[uint16(rInt)]
		if !taken {
			ret = uint16(rInt)
			break
		}
	}
	u.UsersITOS[ret] = name
	u.UsersSTOI[name] = ret
	return ret, nil
}

func (u *UserCenter) GetID(name string) (uint16, error) {
	id, have := u.UsersSTOI[name]
	if !have {
		return 0, ERROR_DONT_HAVE
	}
	return id, nil
}

func (u *UserCenter) GetName(id uint16) (string, error) {
	name, have := u.UsersITOS[id]
	if !have {
		return "", ERROR_DONT_HAVE
	}
	return name, nil
}
