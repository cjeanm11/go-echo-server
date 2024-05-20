package database

import (
	"time"
	"fmt"
	"database/sql"
	"context"
)

type UserRepository interface {
    GetUserByEmail(email string) (map[string]string, error)
    AddUser(username string, email string, password string) map[string]string
}

func (s *Service) AddUser(username string, email string, password string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return map[string]string{"error": err.Error()}
	}

	var userID int
	err = stmt.QueryRowContext(ctx, username, email, password).Scan(&userID)
	if err != nil {
		return map[string]string{"error": err.Error()}
	}

	defer stmt.Close()

	return map[string]string{
		"message": "User added successfully",
		"user_id": fmt.Sprintf("%d", userID),
	}
}

func (s *Service) GetUserByEmail(email string) (map[string]string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var userID int
    var username string
    err := s.db.QueryRowContext(ctx, "SELECT id, username FROM users WHERE email = $1", email).Scan(&userID, &username)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil 
        }
        return nil, err
    }

    return map[string]string{
        "user_id":  fmt.Sprintf("%d", userID),
        "username": username,
    }, nil
}
 