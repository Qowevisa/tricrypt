package main

import (
	"log"
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

func (u *UserCenter) AddUser(name string, id uint16) error {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	_, alreadyHave := u.UsersSTOI[name]
	if alreadyHave {
		return ERROR_ALREADY_HAVE
	}
	log.Printf("Users: add %s with %d id\n", name, id)
	u.UsersITOS[id] = name
	u.UsersSTOI[name] = id
	return nil
}

func (u *UserCenter) DeleteIfHaveOne(id uint16) {
	name, found := u.UsersITOS[id]
	if !found {
		log.Printf("User with %d id is not found; Can not delete\n", id)
		return
	}
	delete(u.UsersITOS, id)
	delete(u.UsersSTOI, name)
	log.Printf("User with %s name and %d id was found; User is deleted\n", name, id)
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
