package main

import (
	"log"
	"os"
	"sample/server"
)

func main() {
	log.Println("Starting gRPC servers...")
	err := server.StartServers()
	if err != nil {
		log.Fatalf("Failed to start servers: %v", err)
		os.Exit(1)
	}
	log.Println("Servers are running. Press Ctrl+C to stop.")
	select {}
}
