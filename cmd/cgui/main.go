package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	com "git.qowevisa.me/Qowevisa/gotell/communication"
	"git.qowevisa.me/Qowevisa/gotell/env"
	"github.com/gorilla/websocket"
)

func main() {
	loadingFileName := env.ServerFullchainFileName
	cert, err := os.ReadFile(loadingFileName)
	if err != nil {
		log.Fatalf("client: load root cert: %s", err)
	}
	log.Printf("Certificate %s loaded successfully!\n", loadingFileName)
	//
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("client: failed to parse root certificate")
	}

	config := &tls.Config{
		RootCAs: roots,
	}

	host, err := env.GetHost()
	if err != nil {
		log.Fatal(err)
	}
	port, err := env.GetPort()
	if err != nil {
		log.Fatal(err)
	}

	service := fmt.Sprintf("%s:%d", host, port)
	conn, err := tls.Dial("tcp", service, config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Printf("client: connected to %s", service)

	// Connect to the Electron.js application via WebSocket
	u := url.URL{Scheme: "ws", Host: "localhost:8081", Path: "/ws"}
	var ws *websocket.Conn
	for {
		ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Printf("Error: dial: %v\n", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	defer ws.Close()

	go readFromServer(conn, ws)
	go readFromWebSocket(conn, ws)
	select {}
}

func readFromServer(conn net.Conn, ws *websocket.Conn) {
	buf := make([]byte, 70000)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("client: read: %s", err)
			return
		}
		msg, err := com.Decode(buf[:n])
		if err != nil {
			log.Printf("ERROR: %#v\n", err)
			continue
		}
		if msg == nil {
			continue
		}
		log.Printf("client: readServer: received message from server: %v", *msg)
		switch msg.ID {
		case com.ID_SERVER_APPROVE_CLIENT_NICKNAME:
			newID := binary.BigEndian.Uint16(msg.Data)
			msg.FromID = newID
			msg.Data = []byte{}
		}
		log.Printf("client: readServer: sending message to websocket: %v", *msg)
		ws.WriteJSON(*msg)
	}
}

func readFromWebSocket(conn net.Conn, ws *websocket.Conn) {
	for {
		var msg com.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %s", err)
			return
		}
		msg.Version = com.V1
		log.Printf("client: readWS: received message from Electron: %v", msg)
		encodedMsg, err := msg.Bytes()
		if err != nil {
			log.Printf("Encoding error: %s", err)
			continue
		}
		log.Printf("client: readWS: sending data to server: %v", encodedMsg)
		conn.Write(encodedMsg)
	}
}
