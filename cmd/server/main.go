package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"

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
	//
	serverCert, err := os.ReadFile("./server.pem")
	if err != nil {
		log.Fatal(err)
	}
	serverKey, err := os.ReadFile("./server.key")
	if err != nil {
		log.Fatal(err)
	}
	cer, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	//

	log.Printf("Serving on %s:%d\n", host, port)
	l, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			io.Copy(os.Stdout, c)
			fmt.Println()
			c.Close()
		}(conn)
	}
}
