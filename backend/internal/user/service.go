package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	user   *UserStore
	secret string
}

func NewUserService(user *UserStore, secret string) *UserService {
	return &UserService{
		user: user,
		secret: secret,
	}
}

func (s *UserService) SignUp(ctx context.Context, username, email, password string) (string, error) {
	// Check if username already exists
	exists, err := s.user.UsernameExists(ctx, username)
	if err != nil {
		return "", err
	}
	if exists {
		return "", ErrUsernameTaken
	}

	// Check if email already exists
	exists, err = s.user.EmailExists(ctx, email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", ErrUsernameTaken
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Save user to databse
	userID, err := s.user.Create(ctx, username, email, string(hash))
	if err != nil {
		return "", err
	}

	// Return JWT
	return s.generateToken(userID)
}

func (s *UserService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := s.user.GetByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	return s.generateToken(user.ID)
}

func (s *UserService) Me(ctx context.Context, id string) (User, error) {
	user, err := s.user.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserDoesNotExist
		}
		return User{}, err
	}

	return user, nil
}

func (s *UserService) generateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
	})
	return token.SignedString([]byte(s.secret))
}
