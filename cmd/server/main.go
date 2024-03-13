package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"

	"git.qowevisa.me/Qowevisa/gotell/env"
)

func main() {
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
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("server: conn: read: %s", err)
			}
			break
		}
		log.Printf("server: conn: echo %q\n", string(buf[:n]))
		answer := append([]byte("Hello! I see your message:"), buf[:n]...)
		_, err = conn.Write(answer)
		if err != nil {
			log.Printf("server: conn: write: %s", err)
			break
		}
	}
	log.Println("server: conn: closed")
}
