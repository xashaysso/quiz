package services

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"quiz/db/repositories"
	entities "quiz/entities/db"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type QuizService struct {
	QuizRepo     repositories.QuizRepository
	QuestionRepo repositories.QuestionRepository
	SessionRepo  repositories.SessionRepository
}

func NewQuizService(quizRepo repositories.QuizRepository, questionRepo repositories.QuestionRepository, sessionRepo repositories.SessionRepository) QuizServiceInterface {
	return &QuizService{
		QuizRepo:     quizRepo,
		QuestionRepo: questionRepo,
		SessionRepo:  sessionRepo,
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

	slog.Info("quiz deleted successfully", slog.Int("quiz_id", qID), slog.Int("deleted_by", userID))

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

	slog.Info("quiz created successfully", slog.Int("quiz_id", newQuiz.ID), slog.Int("creator_id", userID))

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

	slog.Info("quiz updated successfully", slog.Int("quiz_id", qID), slog.Int("updated_by", userID))

	return s.QuizRepo.UpdateQuiz(ctx, qID, name, description)
}

func (s *QuizService) StartQuiz(ctx context.Context, userID int64, quizID string) (string, error) {
	qID, err := strconv.ParseInt(quizID, 10, 64)
	if err != nil {
		return "", ErrInvalidIDFormat
	}

	questionIDs, err := s.QuestionRepo.GetQuestionIDsByQuizID(ctx, qID)
	if err != nil {
		return "", err
	}
	if len(questionIDs) == 0 {
		return "", ErrQuizHasNoQuestions
	}

	sessionID := uuid.New().String()

	session := entities.QuizSession{
		SessionID:         sessionID,
		UserID:            userID,
		QuizID:            qID,
		CurrentScore:      0,
		Questions:         questionIDs,
		AnsweredQuestions: make([]int64, 0),
	}

	err = s.SessionRepo.SaveQuizSession(ctx, session, time.Hour)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}
