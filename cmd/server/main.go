package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	com "git.qowevisa.me/Qowevisa/gotell/communication"
	"git.qowevisa.me/Qowevisa/gotell/env"
)

func main() {
	userCenter.Init()
	host, err := env.GetHost()
	if err != nil {
		log.Fatal(err)
	}
	port, err := env.GetPort()
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.LoadX509KeyPair(
		env.ServerFullchainFileName,
		env.ServerPrivKeyFileName,
	)
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.NoClientCert}
	config.Rand = rand.Reader

	service := fmt.Sprintf("%s:%d", host, port)
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Printf("server: listening on %s", service)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	ask, err := com.ServerAskClientAboutNickname()
	if err != nil {
		log.Printf("ERROR: %#v\n", err)
	} else {
		log.Printf("Trying to send %#v\n", ask)
		_, err = conn.Write(ask)
		if err != nil {
			log.Printf("ERROR: %#v\n", err)
		}
	}
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("server: conn: read: %s", err)
			}
			break
		}
		msg, err := com.Decode(buf[:n])
		if err != nil {
			log.Printf("ERROR: %#v\n", err)
			continue
		}
		if msg == nil {
			log.Printf("ERROR MSG IS NIL\n")
			continue
		}
		log.Printf("server: conn: receive %#v\n", *msg)
		// Handle
		switch msg.ID {
		case com.ID_CLIENT_SEND_SERVER_NICKNAME:
			id, err := userCenter.AddUser(string(msg.Data))
			if err != nil {
				answ, err := com.ServerSendClientDecline()
				if err != nil {
					log.Printf("ERROR: %v\n", err)
					continue
				}
				conn.Write(answ)
			} else {
				idBytes := make([]byte, 4)
				binary.BigEndian.PutUint32(idBytes, uint32(id))
				answ, err := com.ServerSendClientHisID(idBytes)
				if err != nil {
					log.Printf("ERROR: %v\n", err)
					continue
				}
				conn.Write(answ)
			}
		case com.ID_CLIENT_SEND_SERVER_LINK:
		default:
		}
		// Handle
	}
	log.Println("server: conn: closed")
}
