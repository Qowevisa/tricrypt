package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"git.qowevisa.me/Qowevisa/gotell/env"
)

func main() {
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("Reading CA cert file: %s", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	url := fmt.Sprintf("https://chat.qowevisa.me:%d", env.ConnectPort)
	response, err := client.Get(url)
	if err != nil {
		log.Fatalf("Failed to request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %s", err)
	}

	log.Printf("Server response: %s", body)
}

