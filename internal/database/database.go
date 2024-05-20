package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8" 
	_ "github.com/joho/godotenv/autoload"
	"github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/pgx"
    _ "github.com/golang-migrate/migrate/v4/database/postgres" 
    _ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"	
	_ "database/sql"
	"server-template/pkg/util"
)

var (
	database  = util.GetEnvOrDefault("DB_DATABASE","postgres")
	password  = util.GetEnvOrDefault("DB_PASSWORD", "postgres")
	username  = util.GetEnvOrDefault("DB_USERNAME", "")
	port      = util.GetEnvOrDefault("DB_PORT","5432")
	host      = util.GetEnvOrDefault("DB_HOST","localhost")
)

type Service struct {
	db *sql.DB
	redis *redis.Client 
}

func MigrateDatabase(connStr string) error {
    m, err := migrate.New("file://db/migrations", connStr)
    if err != nil {
        return fmt.Errorf("failed to create migration instance: %w", err)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("failed to run migrations: %w", err)
    }
    return nil
}

func New() *Service {

	
	// Postgres Connection
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	
    // Redis Connection
    redisClient := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", 
        Password: "",               
        DB:       0,                
    })

    // database Migration
    if err := MigrateDatabase(connStr); err != nil {
 	  	log.Fatalf("Failed database migration: %v", err)
    }

	return &Service{db: db, redis: redisClient}
}

func (s *Service) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}


func (s *Service) Health() map[string]string {
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
