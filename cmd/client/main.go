package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"git.qowevisa.me/Qowevisa/gotell/communication"
	"git.qowevisa.me/Qowevisa/gotell/env"
)

var messageChannel chan ([]byte)

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
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done() // Mark this goroutine as done when it exits
		buf := make([]byte, 1024)
		for {
			select {
			case <-ctx.Done(): // Check if context cancellation has been requested
				return
			default:
				timeoutDuration := 5 * time.Second
				err := conn.SetReadDeadline(time.Now().Add(timeoutDuration))
				if err != nil {
					panic(err)
				}
				n, err := conn.Read(buf)
				if err != nil {
					if err != io.EOF {
						if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
							continue
						}
						panic(err)
					}
					return
				}
				messageChannel <- buf[:n]
			}
		}
	}(ctx, wg)
	go func() {
		for message := range messageChannel {
			msg, err := communication.Decode(message)
			if err != nil {
				log.Printf("ERROR: %v\n", err)
			}
			switch msg.From {
			case communication.FROM_SERVER:

			case communication.FROM_CLIENT:
			}
		}
	}()
	cancel()
	wg.Wait()
}
