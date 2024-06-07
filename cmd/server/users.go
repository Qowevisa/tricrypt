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
	UsersSTOI map[string]int32
	UsersITOS map[int32]string
	Mu        sync.Mutex
}

func (u *UserCenter) Init() {
	u.UsersSTOI = make(map[string]int32)
	u.UsersITOS = make(map[int32]string)
}

var (
	ERROR_ALREADY_HAVE = errors.New("Username is already taken")
	ERROR_DONT_HAVE    = errors.New("Username was not found")
)

func (u *UserCenter) AddUser(name string) (int32, error) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	_, alreadyHave := u.UsersSTOI[name]
	if alreadyHave {
		return -1, ERROR_ALREADY_HAVE
	}
	var ret int32
	for {
		rInt := rand.Int31()
		_, taken := u.UsersITOS[rInt]
		if !taken {
			ret = rInt
			break
		}
	}
	u.UsersITOS[ret] = name
	u.UsersSTOI[name] = ret
	return ret, nil
}

func (u *UserCenter) GetID(name string) (int32, error) {
	id, have := u.UsersSTOI[name]
	if !have {
		return -1, ERROR_DONT_HAVE
	}
	return id, nil
}

func (u *UserCenter) GetName(id int32) (string, error) {
	name, have := u.UsersITOS[id]
	if !have {
		return "", ERROR_DONT_HAVE
	}
	return name, nil
}
