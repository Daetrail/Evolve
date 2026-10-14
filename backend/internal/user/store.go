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

func (s *UserStore) GetByID(ctx context.Context, id string) (User, error) {
	var user User
	err := s.db.QueryRowContext(ctx,
		"SELECT id, username, email, password FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := s.db.QueryRowContext(ctx,
		"SELECT id, username, email, password FROM users WHERE email = ?",
		strings.ToLower(email),
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
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

func (s *UserStore) UsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)",
		strings.ToLower(username),
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

func (s *UserStore) IDExists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)",
		id,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}