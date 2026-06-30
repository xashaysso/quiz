package service

import "errors"

var (
	ErrInvalidUsername   = errors.New("username is too short")
	ErrInvalidPassword   = errors.New("password is too short")
	ErrWrongCredentials  = errors.New("wrong username or password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrSessionExpired    = errors.New("user session expired")
)
