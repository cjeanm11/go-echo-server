package main

import (
	"log"
	srv "server-template/internal/server"
	"github.com/joho/godotenv"
	"os"
)

func main() {

	 if err := godotenv.Load(".env-dev", ".env"); err != nil {
        log.Fatalf("Failed to load env variables: %v", err)
    }
	logFile, err := os.OpenFile("./app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	
	log.Println("Application started")
	server := srv.NewServer(srv.WithPort(8080))
	server.Start()
}
