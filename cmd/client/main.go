package main

import (
	"bufio"
	"crypto/tls"
	"fmt"

	// "fmt"
	"io"
	"log"
	"os"

	"git.qowevisa.me/Qowevisa/gotell/env"
)

func main() {
	cert, err := tls.LoadX509KeyPair("tls.crt", "tls.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	// config.Rand = rand.Reader

	url := fmt.Sprintf("chat.qowevisa.me:%d", env.ConnectPort)
	// Dial a TLS connection
	conn, err := tls.Dial("tcp", url, &config)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Buffer to hold data read from the connection
	// buf := make([]byte, 1024) // Adjust size as needed
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		text := reader.Text()
		// Read from the connection
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %v", err)
			}
			break
		}
		// fmt.Printf("Received: %s\n", string(buf[:n]))
	}
}
