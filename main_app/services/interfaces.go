package services

import (
	"context"
	entities "quiz/entities/db"
	"quiz/entities/dto"
)

type QuizServiceInterface interface {
	checkIfAuthor(ctx context.Context, quizID int, userID int) error
	ListQuizzes(ctx context.Context) ([]entities.Quiz, error)
	DeleteQuiz(ctx context.Context, quizID string, userID int) error
	CreateQuiz(ctx context.Context, name, description string, userID int) (entities.Quiz, error)
	UpdateQuiz(ctx context.Context, quizID string, name, description *string, userID int) (entities.Quiz, error)
}

type QuestionServiceInterface interface {
	CreateQuestion(ctx context.Context, quizID string, body dto.CreateQuestionDTO, userID int) (dto.QuestionResponse, error)
	ListQuestions(ctx context.Context, quizID string) ([]dto.QuestionResponse, error)
	GetQuestion(ctx context.Context, questionID string) (dto.QuestionResponse, error)
	UpdateQuestion(ctx context.Context, questionID string, data dto.UpdateQuestionDTO, userID int) (dto.QuestionResponse, error)
	DeleteQuestion(ctx context.Context, questionID string, userID int) error
}

type AnswerServiceInterface interface {
	CheckAnswer(ctx context.Context, questionID string, answerID int) (bool, error)
	ListAnswers(ctx context.Context, questionID string) ([]entities.Answer, error)
	CreateAnswer(ctx context.Context, questionID string, data dto.CreateAnswerDTO, userID int) (dto.AnswerPublicResponse, error)
	GetAnswer(ctx context.Context, answerID string) (dto.AnswerPublicResponse, error)
	DeleteAnswer(ctx context.Context, answerID string, userID int) error
	UpdateAnswer(ctx context.Context, answerID string, data dto.UpdateAnswerDTO, userID int) (dto.AnswerPublicResponse, error)
}

type AuthServiceInterface interface {
	Register(ctx context.Context, username, password string) (entities.User, string, error)
	Login(ctx context.Context, username, password string) (string, error)
	Logout(ctx context.Context, token string) error
}
