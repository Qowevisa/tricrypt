package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"git.qowevisa.me/Qowevisa/gotell/env"
)

func main() {
	url := fmt.Sprintf("127.0.0.1:%d", env.Port)
	listener, err := net.Listen("tcp", url)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("Server is listening on %s\n", url)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Client connected: %v\n", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("Received: %s\n", text)
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
