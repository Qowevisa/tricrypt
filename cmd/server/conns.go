package main

import (
	"log"
	"net"
	"sync"
)

var (
	connCenter ConnCenter
)

type ConnCenter struct {
	Conns map[uint16]net.Conn
	Mu    sync.Mutex
}

func (c *ConnCenter) Init() {
	c.Conns = make(map[uint16]net.Conn)
}

func (c *ConnCenter) AddConn(id uint16, con net.Conn) (uint16, error) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	log.Printf("Conns: add %v with %d id\n", con, id)
	c.Conns[id] = con
	return id, nil
}

func (c *ConnCenter) DeleteIfHaveOne(id uint16) {
	name, found := c.Conns[id]
	if !found {
		log.Printf("Conn with %d id is not found; Can not delete\n", id)
		return
	}
	delete(c.Conns, id)
	log.Printf("Conn with %v con and %d id was found; Conn is deleted\n", name, id)
}

func (c *ConnCenter) GetConn(id uint16) (net.Conn, error) {
	con, have := c.Conns[id]
	if !have {
		return nil, ERROR_DONT_HAVE
	}
	return con, nil
}
