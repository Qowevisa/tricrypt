package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"git.qowevisa.me/Qowevisa/gotell/env"
)

func main() {
	url := fmt.Sprintf("chat.qowevisa.me:%d", env.ConnectPort)
	conn, err := tls.Dial("tcp", url, &tls.Config{
		InsecureSkipVerify: false, // Set to true if using self-signed certificates
	})
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		text := reader.Text()
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}
