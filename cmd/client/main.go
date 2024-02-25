package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
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
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()

	log.Println("client: connected to: ", conn.RemoteAddr())

	message := "Hello secure Server\n"
	n, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}

	log.Printf("client: wrote %q (%d bytes)", message, n)

	reply := make([]byte, 256)
	n, err = conn.Read(reply)
	if err != nil {
		log.Fatalf("client: read: %s", err)
	}

	log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)

	log.Print("client: exiting")
}
