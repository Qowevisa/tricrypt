package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"git.qowevisa.me/Qowevisa/gotell/env"
)

func main() {
	// Listen on TCP port 8080 on all available unicast and anycast IP addresses of the local system.
	cert, err := tls.LoadX509KeyPair("tls.crt", "tls.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	url := fmt.Sprintf(":%d", env.Port)
	listener, err := tls.Listen("tcp", url, &config)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleConnection(conn)
	}
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
