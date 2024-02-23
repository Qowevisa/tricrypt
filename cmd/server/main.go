package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"git.qowevisa.me/Qowevisa/gotell/env"
)

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've connected to the server!")
	log.Printf("w: %#v", w)
	log.Printf("r: %#v", r)
}

func main() {
	// Listen on TCP port 8080 on all available unicast and anycast IP addresses of the local system.
	cert, err := tls.LoadX509KeyPair("tls.crt", "tls.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	host, err := env.GetHost()
	if err != nil {
		panic(err)
	}
	port, err := env.GetPort()
	if err != nil {
		panic(err)
	}
	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		TLSConfig:    &config,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		Handler:      http.HandlerFunc(handle),
	}
	log.Printf("Start http server on %s:%d\n", host, port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
	defer srv.Close()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Client connected: %v\n", conn.RemoteAddr())

	// Create a new reader for each client.
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// Read the incoming connection into the buffer.
		text := scanner.Text()
		fmt.Printf("Received: %s\n", text)

		// Send a response back to client.
		_, err := conn.Write([]byte("Message received: " + text + "\n"))
		if err != nil {
			log.Printf("Failed to write to connection: %v", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from connection: %s\n", err)
	}

	fmt.Printf("Client disconnected: %v\n", conn.RemoteAddr())
}
