package services

import "errors"

var (
	ErrInvalidAnswerText = errors.New("field 'text' is required in answer body")
	ErrAnswerNotFound    = errors.New("answer not found")
	ErrNoFieldsToUpdate  = errors.New("fields text and/or correct_id are required in json body")
	ErrInvalidCorrectID  = errors.New("invalid correct answer id")
	ErrInvalidIDFormat   = errors.New("invalid 'ID' field format")
	ErrNoQuestionAnswers = errors.New("no field 'answers' provided in json")
	ErrQuestionNotFound  = errors.New("question not found")
	ErrQuizNotFound      = errors.New("quiz not found")
	ErrNotAnAuthor       = errors.New("you are not an author of this quiz")
	ErrInvalidName       = errors.New("quiz name is too short")
	ErrNoRequiredFields  = errors.New("fields name/description are required in json body")

	// auth
	ErrInvalidUsername   = errors.New("username is too short")
	ErrInvalidPassword   = errors.New("password is too short")
	ErrWrongCredentials  = errors.New("wrong username or password")
	ErrUserAlreadyExists = errors.New("user already exists")
)
