package repositories

import "errors"

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrInvalidCorrectID = errors.New("new correct answer id belongs to another question")
)
