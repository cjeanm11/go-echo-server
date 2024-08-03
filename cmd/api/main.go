package main

import (
	"log"
	srv "server-template/internal/server"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env-dev", ".env"); err != nil {
            log.Fatalf("Failed to load env variables: %v", err)
	}
	
	log.Println("Application started")
	server := srv.NewServer(srv.WithPort(8080))
	server.Start()
}
