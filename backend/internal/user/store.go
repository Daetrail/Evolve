package user

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, username, email, passwordHash string) (string, error) {
	id := uuid.New().String()
	_, err := s.db.ExecContext(ctx, 
		"INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)",
		id, username, email, passwordHash,
	)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *UserStore) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool 
	err := s.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)",
		strings.ToLower(email),
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}