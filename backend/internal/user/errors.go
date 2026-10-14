package user

import "errors"

var ErrEmailTaken = errors.New("email already taken")
var ErrUsernameTaken = errors.New("username already taken")
var ErrWeakPassword = errors.New("password is too weak")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserDoesNotExist = errors.New("user does not exist")
