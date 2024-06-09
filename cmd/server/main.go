package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	com "git.qowevisa.me/Qowevisa/gotell/communication"
	"git.qowevisa.me/Qowevisa/gotell/env"
)

func captureHeapProfile() {
	currentTime := time.Now().Format("2006_01_02T15_04")
	filename := fmt.Sprintf("heap_%s.prof", currentTime)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("could not create memory profile: %v", err)
	}
	defer f.Close()
	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatalf("could not write memory profile: %v", err)
	}
}

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	userCenter.Init()
	linkCenter.Init()
	connCenter.Init()
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
	defer listener.Close()
	log.Printf("server: listening on %s", service)

	var wg sync.WaitGroup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("received shutdown signal")
		listener.Close()
		wg.Wait()
		os.Exit(0)
	}()

	ticker := time.NewTicker(3 * time.Minute)
	go func() {
		for range ticker.C {
			captureHeapProfile()
		}
	}()

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
	var registeredID uint16
	var isRegistered bool
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
					log.Printf("ERROR: BYTES: %v\n", err)
					continue
				}
				conn.Write(answ)
			} else {
				idBytes := make([]byte, 2)
				binary.BigEndian.PutUint16(idBytes, id)
				answ, err := com.ServerSendClientHisID(idBytes)
				if err != nil {
					log.Printf("ERROR: BYTES: %v\n", err)
					continue
				}
				conn.Write(answ)
				isRegistered = true
				registeredID = id
				connCenter.AddConn(id, conn)
			}
		case com.ID_CLIENT_SEND_SERVER_LINK:
			l, err := com.DecodeLink(msg.Data)
			if err != nil {
				log.Printf("ERROR: DecodeLink: %v\n", err)
				continue
			}
			err = linkCenter.AddLink(msg.FromID, *l)
			if err != nil {
				log.Printf("ERROR: AddLink: %v\n", err)
				answ, err := com.ServerDeclineClientLink()
				if err != nil {
					log.Printf("ERROR: BYTES: %v\n", err)
					continue
				}
				conn.Write(answ)
				continue
			}
			answ, err := com.ServerApproveClientLink()
			if err != nil {
				log.Printf("ERROR: BYTES: %v\n", err)
				continue
			}
			conn.Write(answ)
		case com.ID_CLIENT_ASK_SERVER_LINK:
			link, err := linkCenter.GetLink(msg.Data)
			if err != nil {
				log.Printf("Error: GetLink: %v\n", err)
				continue
			}
			if link.LeftNum == 0 {
				linkCenter.DeleteLink(msg.Data)
				continue
			}
			// TODO: there can be an error on multi-thread app
			link.LeftNum -= 1
			name, err := userCenter.GetName(link.UserID)
			if err != nil {
				log.Printf("ERROR: userCenter: Getname: %v\n", err)
				continue
			}
			answ, err := com.ServerSendClientIDFromLink(link.UserID, []byte(name))
			if err != nil {
				log.Printf("ERROR: BYTES: %v\n", err)
				continue
			}
			conn.Write(answ)
			// REDIRECTED STUFF
		case com.ID_CLIENT_ASK_CLIENT_HANDSHAKE:
			toConn, err := connCenter.GetConn(msg.ToID)
			if err != nil {
				log.Printf("ERROR: connCenter: GetConn: %v\n", err)
				continue
			}
			log.Printf("Redirecting msg to %d\n", msg.ToID)
			toConn.Write(buf[:n])
		default:
		}
		// Handle
	}
	log.Println("server: conn: closed")
	if isRegistered {
		userCenter.DeleteIfHaveOne(registeredID)
		linkCenter.CleanAfterLeave(registeredID)
		connCenter.DeleteIfHaveOne(registeredID)
	}
}
