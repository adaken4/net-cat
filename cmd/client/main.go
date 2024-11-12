package main

import (
	"log"

	"net-cat/internal/client"
)

func main() {
	cl := client.NewClient("localhost:8080")
	if err := cl.Start(); err != nil {
		log.Fatalf("Client error: %v", err)
	}
}
