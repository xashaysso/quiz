package repository

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrUserAlreadyExists = errors.New("this user already exists")
	ErrSessionExpired    = errors.New("quiz session expired")
)
