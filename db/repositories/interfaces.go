package repositories

import (
	"context"
	APIentities "quiz/entities/api"
	entities "quiz/entities/db"
	"quiz/entities/dto"
	"time"
)

type QuizRepository interface {
	GetQuiz(ctx context.Context)([]entities.Quiz, error)
	GetQuizByID(ctx context.Context, quizID int) (entities.Quiz, error)
	DeleteQuiz(ctx context.Context, quizID int)(error)
	CreateQuiz(ctx context.Context, quiz_name string, quiz_description string, userID int)(entities.Quiz, error)
	UpdateQuiz(ctx context.Context, quizID int, name *string, description *string)(entities.Quiz, error)
}

type QuestionRepository interface {
	GetQuizQuestions(ctx context.Context, id int) ([]APIentities.QuestionAPI, error)
	CreateQuestion(ctx context.Context, quizID int, data dto.CreateQuestionDTO)(APIentities.QuestionAPI, error)
	GetQuestion(ctx context.Context, questionID int)(APIentities.QuestionAPI, error)
	UpdateQuestion(ctx context.Context, questionID int, data dto.UpdateQuestionDTO)(APIentities.QuestionAPI, error)
	DeleteQuestion(ctx context.Context, questionID int)(error)
}

type AnswerRepository interface {
	GetQuizAnswers(ctx context.Context, questionID int) ([]entities.Answer, error)
	CheckAnswer(ctx context.Context, questionID int, answerID int) (bool, error)
	CreateAnswer(ctx context.Context, questionID int, data dto.CreateAnswerDTO)(APIentities.AnswerAPI, error)
	GetAnswer(ctx context.Context, answerID int)(APIentities.AnswerAPI, error)
	DeleteAnswer(ctx context.Context, answerID int)(error)
	UpdateAnswer(ctx context.Context, answerID int, data dto.UpdateAnswerDTO)(APIentities.AnswerAPI, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, username string, password_hash string) (entities.User, error)
	GetByUsername(ctx context.Context, username string) (entities.User, error)
}

type SessionRepository interface {
	Get(ctx context.Context, token string)(int, error)
	Set(ctx context.Context, token string, userID int, ttl time.Duration)(error)
	Delete(ctx context.Context, token string) (error)
}