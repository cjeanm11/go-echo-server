package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"server-template/internal/database"
)

type Server struct {
	port       int
	db         database.Service
	httpServer *http.Server
}

type Option func(*Server)

func NewServer(options ...Option) *Server {
	s := &Server{
		port: loadPortFromEnv(), 
		db:   database.New(),
	}

	for _, option := range options {
		option(s)
	}

	s.initHTTPServer()
	return s
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithDatabaseService(db database.Service) Option {
	return func(s *Server) {
		s.db = db
	}
}

func (s *Server) initHTTPServer() {
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(), // Implement this method to setup routes
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) Start() {
	s.db = database.New()
	defer func() {
		// graceful shutdown
		if err := s.db.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
	}()

	log.Printf("Server starting on port %d\n", s.port)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

}

func loadPortFromEnv() int {
	portStr, exists := os.LookupEnv("PORT")
	if !exists {
		return 8080 
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Warning: Invalid PORT environment variable '%s', falling back to default port 8080.\n", portStr)
		return 8080
	}

	return port
}
