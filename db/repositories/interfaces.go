package repositories

import (
	"context"
	APIentities "quiz/entities/api"
	entities "quiz/entities/db"
	"quiz/entities/dto"
)

type QuizRepository interface {
	GetQuiz(ctx context.Context)([]entities.Quiz, error)
	DeleteQuiz(ctx context.Context, quizID string)(error)
	CreateQuiz(ctx context.Context, quiz_name string, quiz_description string)(entities.Quiz, error)
	UpdateQuiz(ctx context.Context, quizID string, name *string, description *string)(entities.Quiz, error)
}

type QuestionRepository interface {
	GetQuizQuestions(ctx context.Context, id string) ([]APIentities.QuestionAPI, error)
	CreateQuestion(ctx context.Context, quizID int, data dto.CreateQuestionDTO)(APIentities.QuestionAPI, error)
	GetQuestion(ctx context.Context, questionID string)(APIentities.QuestionAPI, error)
	UpdateQuestion(ctx context.Context, questionID string, data dto.UpdateQuestionDTO)(APIentities.QuestionAPI, error)
	DeleteQuestion(ctx context.Context, questionID string)(error)
}

type AnswerRepository interface {
	GetQuizAnswers(ctx context.Context, question_id string) ([]entities.Answer, error)
	CheckAnswer(ctx context.Context, questionID string, answerID int) (bool, error)
	CreateAnswer(ctx context.Context, questionID int, data dto.CreateAnswerDTO)(APIentities.AnswerAPI, error)
	GetAnswer(ctx context.Context, answerID int)(APIentities.AnswerAPI, error)
	DeleteAnswer(ctx context.Context, answerID int)(error)
	UpdateAnswer(ctx context.Context, answerID int, data dto.UpdateAnswerDTO)(APIentities.AnswerAPI, error)
}