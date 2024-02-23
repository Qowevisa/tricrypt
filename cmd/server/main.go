package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've reached the server!")
	})

	log.Println("Starting server on :1993...")
	err := http.ListenAndServeTLS(":1993", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

