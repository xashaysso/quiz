package services

import (
	"context"
	"database/sql"
	"errors"
	"quiz/db/repositories"
	entities "quiz/entities/db"
	"strconv"

	"github.com/jackc/pgx/v5"
)

var (
	ErrQuizNotFound = errors.New("quiz not found")
	ErrNotAnAuthor = errors.New("you are not an author of this quiz")
	ErrInvalidName = errors.New("quiz name is too short")
	ErrNoRequiredFields = errors.New("fields name/description are required in json body")
)

type QuizService struct {
	QuizRepo repositories.QuizRepository
}

func (s *QuizService) ListQuizzes(ctx context.Context) ([]entities.Quiz, error) {
	quizzes, err := s.QuizRepo.GetQuiz(ctx)
	if err != nil {
		return []entities.Quiz{}, err
	}
	return quizzes, nil
}

func (s *QuizService) DeleteQuiz(ctx context.Context, quizID string, userID int)(error) {
	qID, err := strconv.Atoi(quizID)
	if err != nil {
		return ErrInvalidIDFormat
	}
	quiz, err := s.QuizRepo.GetQuizByID(ctx, qID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return ErrQuizNotFound
		}
		return err
	}
	if quiz.CreatorID != userID {
		return ErrNotAnAuthor
	}
	return s.QuizRepo.DeleteQuiz(ctx, qID)
}

func (s *QuizService) CreateQuiz(ctx context.Context, name, description string, userID int) (entities.Quiz, error){
	if len(name) < 5 {
		return entities.Quiz{}, ErrInvalidName
	}
	newQuiz, err := s.QuizRepo.CreateQuiz(ctx, name, description, userID);
	if err != nil {
		return entities.Quiz{}, err
	}
	return newQuiz, nil
}

func (s *QuizService) UpdateQuiz(ctx context.Context, quizID string, name, description *string, userID int) (entities.Quiz, error){
	qID, err := strconv.Atoi(quizID)
	if err != nil {
		return entities.Quiz{}, ErrInvalidIDFormat
	}
	if name == nil && description == nil{
		return entities.Quiz{}, ErrNoRequiredFields
	}
	quiz, err := s.QuizRepo.GetQuizByID(ctx, qID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return entities.Quiz{}, ErrQuizNotFound
		}
		return entities.Quiz{}, err
	}
	if quiz.CreatorID != userID {
		return entities.Quiz{}, ErrNotAnAuthor
	}

	newQuiz, err := s.QuizRepo.UpdateQuiz(ctx, qID, name, description);
	if err != nil {
		return entities.Quiz{}, err
	}
	return newQuiz, nil
}
