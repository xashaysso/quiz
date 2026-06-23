package repositories

import (
	"context"
	entities "quiz/entities/db"
	"quiz/entities/dto"
	"time"
)

type QuizRepository interface {
	GetQuiz(ctx context.Context) ([]entities.Quiz, error)
	GetQuizByID(ctx context.Context, quizID int) (entities.Quiz, error)
	DeleteQuiz(ctx context.Context, quizID int) error
	CreateQuiz(ctx context.Context, quiz_name string, quiz_description string, userID int) (entities.Quiz, error)
	UpdateQuiz(ctx context.Context, quizID int, name *string, description *string) (entities.Quiz, error)
	SaveAttempt(ctx context.Context, userID int64, quizID int64, score int) error
}

type QuestionRepository interface {
	GetQuizQuestions(ctx context.Context, id int) ([]entities.Question, error)
	GetQuestionIDsByQuizID(ctx context.Context, quizID int64) ([]int64, error)
	CreateQuestion(ctx context.Context, tx any, quizID int, data dto.CreateQuestionDTO) (int, error)
	GetQuestion(ctx context.Context, questionID int) (entities.Question, error)
	UpdateQuestion(ctx context.Context, questionID int, data dto.UpdateQuestionDTO) (entities.Question, error)
	DeleteQuestion(ctx context.Context, questionID int) error
	CheckIfQuizOwner(ctx context.Context, quizID int, userID int) (bool, error)
	CheckIfQuestionOwner(ctx context.Context, questionID int, userID int) (bool, error)
}

type AnswerRepository interface {
	GetQuizAnswers(ctx context.Context, questionID int) ([]entities.Answer, error)
	CheckAnswer(ctx context.Context, questionID int, answerID int) (bool, error)
	CreateAnswer(ctx context.Context, tx any, questionID int, text string, isCorrect bool) (int, error)
	GetAnswer(ctx context.Context, answerID int) (entities.Answer, error)
	DeleteAnswer(ctx context.Context, answerID int) error
	UpdateAnswer(ctx context.Context, answerID int, data dto.UpdateAnswerDTO) (entities.Answer, error)
	CheckIfAnswerOwner(ctx context.Context, answerID int, userID int) (bool, error)
	CheckIfQuestionOwner(ctx context.Context, questionID int, userID int) (bool, error)
	GetAnswersByQuestionIDs(ctx context.Context, questionIDs []int) ([]entities.Answer, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, username string, password_hash string) (entities.User, error)
	GetByUsername(ctx context.Context, username string) (entities.User, error)
}

type SessionRepository interface {
	Get(ctx context.Context, token string) (int, error)
	Set(ctx context.Context, token string, userID int, ttl time.Duration) error
	Delete(ctx context.Context, token string) error
	SaveQuizSession(ctx context.Context, session entities.QuizSession, ttl time.Duration) error
	GetQuizSession(ctx context.Context, sessionID string) (*entities.QuizSession, error)
	DeleteQuizSession(ctx context.Context, sessionID string) error
}

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(tx any) error) error
}

type QuizEventProducer interface {
	PublishQuizPassed(ctx context.Context, userID int64, quizID int64, score int) error
}
