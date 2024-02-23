package main

import (
	"crypto/tls"
	"crypto/x509"
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
		panic(err)
	}
	port, err := env.GetPort()
	if err != nil {
		panic(err)
	}
	//
	rootCert, err := os.ReadFile("./server.pem")
	if err != nil {
		panic(err)
	}
	//

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootCert)
	if !ok {
		log.Fatal("failed to parse root certificate")
	}
	config := &tls.Config{RootCAs: roots, ServerName: "my-server"}

	log.Printf("Trying to dial %s:%d\n", host, port)
	connp, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
	}

	conn := tls.Client(connp, config)
	io.WriteString(conn, "Hello secure Server")
	conn.Close()
}
