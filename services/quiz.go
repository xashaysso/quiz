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

type QuizService struct {
	QuizRepo repositories.QuizRepository
}

func NewQuizService(repo repositories.QuizRepository) QuizServiceInterface {
	return &QuizService{
		QuizRepo: repo,
	}
}

func (s *QuizService) checkIfAuthor(ctx context.Context, quizID int, userID int) error {
	quiz, err := s.QuizRepo.GetQuizByID(ctx, quizID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return ErrQuizNotFound
		}
		return err
	}
	if quiz.CreatorID != userID {
		return ErrNotAnAuthor
	}
	return nil
}

func (s *QuizService) ListQuizzes(ctx context.Context) ([]entities.Quiz, error) {
	quizzes, err := s.QuizRepo.GetQuiz(ctx)
	if err != nil {
		return []entities.Quiz{}, err
	}
	return quizzes, nil
}

func (s *QuizService) DeleteQuiz(ctx context.Context, quizID string, userID int) error {
	qID, err := strconv.Atoi(quizID)
	if err != nil {
		return ErrInvalidIDFormat
	}
	err = s.checkIfAuthor(ctx, qID, userID)
	if err != nil {
		return err
	}
	return s.QuizRepo.DeleteQuiz(ctx, qID)
}

func (s *QuizService) CreateQuiz(ctx context.Context, name, description string, userID int) (entities.Quiz, error) {
	if len(name) < 5 {
		return entities.Quiz{}, ErrInvalidName
	}
	newQuiz, err := s.QuizRepo.CreateQuiz(ctx, name, description, userID)
	if err != nil {
		return entities.Quiz{}, err
	}
	return newQuiz, nil
}

func (s *QuizService) UpdateQuiz(ctx context.Context, quizID string, name, description *string, userID int) (entities.Quiz, error) {
	qID, err := strconv.Atoi(quizID)
	if err != nil {
		return entities.Quiz{}, ErrInvalidIDFormat
	}
	if name == nil && description == nil {
		return entities.Quiz{}, ErrNoRequiredFields
	}
	err = s.checkIfAuthor(ctx, qID, userID)
	if err != nil {
		return entities.Quiz{}, err
	}

	return s.QuizRepo.UpdateQuiz(ctx, qID, name, description)
}
