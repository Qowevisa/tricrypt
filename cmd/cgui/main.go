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
	tlepCenter.Init()
	userCenter.Init()
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

var r com.RegisteredUser
var tmpLink *com.Link
var tmpNick string

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
			r.ID = newID
			if tmpNick != "" {
				r.Name = tmpNick
			}
			r.IsRegistered = true
			break
		case com.ID_SERVER_APPROVE_CLIENT_LINK:
			if tmpLink == nil {
				continue
			}
			msg.ToID = tmpLink.UseCount
			msg.Data = tmpLink.Data
		// Crypto stuff
		case com.ID_CLIENT_SEND_CLIENT_ECDH_PUBKEY:
			t, err := tlepCenter.GetTLEP(msg.FromID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			err = t.ECDHApplyOtherKeyBytes(msg.Data)
			if err != nil {
				log.Printf("ERROR: tlep: ECDHApplyOtherKeyBytes: %v\n", err)
				continue
			}
			fromName, err := userCenter.GetName(msg.FromID)
			if err != nil {
				log.Printf("ERROR: userCenter: GetName: %v\n", err)
			} else {
				msg.Data = []byte(fromName)
			}
		case com.ID_CLIENT_SEND_CLIENT_CBES_SPECS:
			t, err := tlepCenter.GetTLEP(msg.FromID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			cbes, err := t.DecryptMessageEA(msg.Data)
			if err != nil {
				log.Printf("ERROR: tlep: DecryptMessageEA: %v\n", err)
				continue
			}
			err = t.CBESSetFromBytes(cbes)
			if err != nil {
				log.Printf("ERROR: tlep: CBESSetFromBytes: %v\n", err)
				continue
			}
			fromName, err := userCenter.GetName(msg.FromID)
			if err != nil {
				log.Printf("ERROR: userCenter: GetName: %v\n", err)
			} else {
				msg.Data = []byte(fromName)
			}
			// message
		case com.ID_CLIENT_SEND_CLIENT_MESSAGE:
			t, err := tlepCenter.GetTLEP(msg.FromID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			decrypedMsg, err := t.DecryptMessageAtMax(msg.Data)
			if err != nil {
				log.Printf("ERROR: tlep: DecryptMessageAtMax: %v\n", err)
				continue
			}
			msg.Data = decrypedMsg
			// switch
		}
		// user stuff
		switch msg.ID {
		case com.ID_CLIENT_ASK_CLIENT_HANDSHAKE,
			com.ID_CLIENT_APPROVE_CLIENT_HANDSHAKE:
			userCenter.AddUser(string(msg.Data), msg.FromID)
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
		log.Printf("client: readWS: received message from Electron: %v", msg)
		msg.Version = com.V1
		switch msg.ID {
		case com.ID_CLIENT_SEND_SERVER_NICKNAME:
			tmpNick = string(msg.Data)
		case com.ID_CLIENT_SEND_SERVER_LINK:
			if !r.IsRegistered {
				continue
			}
			l, err := r.GenerateLink(msg.ToID)
			if err != nil {
				log.Printf("Error: link: %v", err)
				continue
			}
			tmpLink = &l
			answ, err := com.ClientSendServerLink(r.ID, l)
			if err != nil {
				log.Printf("Error: com: %v", err)
				continue
			}
			log.Printf("client: readWS: sending data to server: %v", answ)
			conn.Write(answ)
			continue
		case com.ID_CLIENT_ASK_CLIENT_HANDSHAKE,
			com.ID_CLIENT_APPROVE_CLIENT_HANDSHAKE,
			com.ID_CLIENT_DECLINE_CLIENT_HANDSHAKE,
			com.ID_CLIENT_SEND_CLIENT_ECDH_PUBKEY,
			com.ID_CLIENT_SEND_CLIENT_CBES_SPECS,
			com.ID_CLIENT_SEND_CLIENT_MKLG_FINGERPRINT,
			com.ID_CLIENT_DECLINE_CLIENT_MKLG_FINGERPRINT,
			com.ID_CLIENT_SEND_CLIENT_MESSAGE:
			if !r.IsRegistered {
				continue
			}
			msg.FromID = r.ID
			// switch
		}
		// user stuff
		switch msg.ID {
		case com.ID_CLIENT_ASK_CLIENT_HANDSHAKE,
			com.ID_CLIENT_APPROVE_CLIENT_HANDSHAKE:
			if r.IsRegistered {
				msg.Data = []byte(r.Name)
			}
		}
		// Crypto stuff
		switch msg.ID {
		case com.ID_CLIENT_ASK_CLIENT_HANDSHAKE,
			com.ID_CLIENT_APPROVE_CLIENT_HANDSHAKE:
			err := tlepCenter.AddTLEP(msg.ToID, fmt.Sprintf("%s-%d", r.Name, msg.ToID))
			if err != nil {
				log.Printf("ERROR: tlepCenter.AddUser: %v\n", err)
			}
		case com.ID_CLIENT_SEND_CLIENT_ECDH_PUBKEY:
			t, err := tlepCenter.GetTLEP(msg.ToID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			key, err := t.ECDHGetPublicKey()
			if err != nil {
				log.Printf("ERROR: tlep: ECDHGetPublicKey: %v\n", err)
				continue
			}
			msg.Data = key
		case com.ID_CLIENT_SEND_CLIENT_CBES_SPECS:
			t, err := tlepCenter.GetTLEP(msg.ToID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			err = t.CBESInitRandom()
			if err != nil {
				log.Printf("ERROR: tlep: CBESInitRandom: %v\n", err)
				continue
			}
			cbes, err := t.CBESGetBytes()
			if err != nil {
				log.Printf("ERROR: tlep: ECDHGetPublicKey: %v\n", err)
				continue
			}
			cbesEAEncr, err := t.EncryptMessageEA(cbes)
			if err != nil {
				log.Printf("ERROR: tlep: EncryptMessageEA: %v\n", err)
				continue
			}
			msg.Data = cbesEAEncr
			// message
		case com.ID_CLIENT_SEND_CLIENT_MESSAGE:
			t, err := tlepCenter.GetTLEP(msg.ToID)
			if err != nil {
				log.Printf("ERROR: tlep: GetTLEP: %v\n", err)
				continue
			}
			encrypedMsg, err := t.EncryptMessageAtMax(msg.Data)
			if err != nil {
				log.Printf("ERROR: tlep: EncryptMessageAtMax: %v\n", err)
				continue
			}
			msg.Data = encrypedMsg
			// switch
		}
		encodedMsg, err := msg.Bytes()
		if err != nil {
			log.Printf("Encoding error: %s", err)
			continue
		}
		log.Printf("client: readWS: sending data to server: %v", encodedMsg)
		conn.Write(encodedMsg)
	}
}
