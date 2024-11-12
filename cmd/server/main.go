package main

import (
	"log"
	"net-cat/internal/server"
)

func main() {
	srv := server.NewServer("localhost:8080")
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
