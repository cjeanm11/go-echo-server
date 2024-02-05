package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
	AddUser(username string, email string) map[string]string
	Close() error
}

type service struct {
	db *sql.DB
}

func New() Service {
	database := os.Getenv("DB_DATABASE")
	password := os.Getenv("DB_PASSWORD")
	username := os.Getenv("DB_USERNAME")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	return &service{db: db}
}

func (s *service) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *service) AddUser(username string, email string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return map[string]string{"error": err.Error()}
	}

	var userID int
	err = stmt.QueryRowContext(ctx, username, email).Scan(&userID)
	if err != nil {
		return map[string]string{"error": err.Error()}
	}

	defer stmt.Close()

	return map[string]string{
		"message": "User added successfully",
		"user_id": fmt.Sprintf("%d", userID),
	}
}
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
