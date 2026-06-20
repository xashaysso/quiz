package services

import (
	"context"
	"errors"
	"log/slog"
	"quiz/db/repositories"
	entities "quiz/entities/db"
	"quiz/entities/dto"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

type AnswerService struct {
	AnswerRepo  repositories.AnswerRepository
	SessionRepo repositories.SessionRepository
	TxManager   repositories.TransactionManager
}

func NewAnswerService(aRepo repositories.AnswerRepository, txm repositories.TransactionManager) AnswerServiceInterface {
	return &AnswerService{
		AnswerRepo: aRepo,
		TxManager:  txm,
	}
}

func (s *AnswerService) CheckAnswer(ctx context.Context, sessionID string, questionID string, answerID int) (bool, error) {
	qID, err := strconv.ParseInt(questionID, 10, 64)
	if err != nil {
		return false, ErrInvalidIDFormat
	}

	session, err := s.SessionRepo.GetQuizSession(ctx, sessionID)
	if err != nil {
		return false, ErrSessionExpired
	}

	isQuestionValid := false
	for _, id := range session.Questions {
		if id == qID {
			isQuestionValid = true
			break
		}
	}
	if !isQuestionValid {
		return false, ErrQuestionDoesNotBelongToQuiz
	}

	for _, id := range session.AnsweredQuestions {
		if id == qID {
			return false, ErrQuestionAlreadyAnswered
		}
	}

	correct, err := s.AnswerRepo.CheckAnswer(ctx, int(qID), answerID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return false, ErrQuestionNotFound
		}
		return false, err
	}

	if correct {
		session.CurrentScore++
	}
	session.AnsweredQuestions = append(session.AnsweredQuestions, qID)

	if len(session.Questions) == len(session.AnsweredQuestions) {
		// TODO:
		// 1. save attempt to quiz_db here
		// 2. push session to kafka

		err = s.SessionRepo.DeleteQuizSession(ctx, sessionID)
		if err != nil {
			return false, err
		}
	} else {
		err = s.SessionRepo.SaveQuizSession(ctx, *session, time.Hour)
		if err != nil {
			return false, err
		}
	}

	return correct, nil
}

func (s *AnswerService) ListAnswers(ctx context.Context, questionID string) ([]entities.Answer, error) {
	qID, err := strconv.Atoi(questionID)
	if err != nil {
		return []entities.Answer{}, ErrInvalidIDFormat
	}

	answers, err := s.AnswerRepo.GetQuizAnswers(ctx, qID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return []entities.Answer{}, ErrQuestionNotFound
		}
		return []entities.Answer{}, err
	}

	return answers, nil
}

func (s *AnswerService) CreateAnswer(ctx context.Context, questionID string, data dto.CreateAnswerDTO, userID int) (dto.AnswerPublicResponse, error) {
	qID, err := strconv.Atoi(questionID)
	if err != nil {
		return dto.AnswerPublicResponse{}, ErrInvalidIDFormat
	}

	if data.Text == "" {
		return dto.AnswerPublicResponse{}, ErrInvalidAnswerText
	}

	isOwner, err := s.AnswerRepo.CheckIfQuestionOwner(ctx, qID, userID)
	if err != nil {
		return dto.AnswerPublicResponse{}, err
	}
	if !isOwner {
		return dto.AnswerPublicResponse{}, ErrNotAnAuthor
	}

	var answer dto.AnswerPublicResponse

	err = s.TxManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		ansID, err := s.AnswerRepo.CreateAnswer(ctx, tx, qID, data.Text, data.IsCorrect)
		if err != nil {
			return err
		}

		answer = dto.AnswerPublicResponse{
			ID:   ansID,
			Text: data.Text,
		}
		return nil
	})
	if err != nil {
		return dto.AnswerPublicResponse{}, err
	}

	slog.Info("answer created successfully", slog.Int("question_id", qID), slog.Int("creator_id", userID))

	return answer, nil
}

func (s *AnswerService) GetAnswer(ctx context.Context, answerID string) (dto.AnswerPublicResponse, error) {
	aID, err := strconv.Atoi(answerID)
	if err != nil {
		return dto.AnswerPublicResponse{}, err
	}

	answer, err := s.AnswerRepo.GetAnswer(ctx, aID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return dto.AnswerPublicResponse{}, ErrAnswerNotFound
		}
		return dto.AnswerPublicResponse{}, err
	}

	return dto.NewAnswerResponse(answer), nil
}

func (s *AnswerService) DeleteAnswer(ctx context.Context, answerID string, userID int) error {
	aID, err := strconv.Atoi(answerID)
	if err != nil {
		return ErrInvalidIDFormat
	}
	isOwner, err := s.AnswerRepo.CheckIfAnswerOwner(ctx, aID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		slog.Warn("unauthorized mutation attempt", slog.Int("user_id", userID), slog.Int("target_answer_id", aID))
		return ErrNotAnAuthor
	}
	err = s.AnswerRepo.DeleteAnswer(ctx, aID)

	slog.Info("answer deleted successfully", slog.Int("answer_id", aID), slog.Int("deleted_by", userID))

	return err
}

func (s *AnswerService) UpdateAnswer(ctx context.Context, answerID string, data dto.UpdateAnswerDTO, userID int) (dto.AnswerPublicResponse, error) {
	aID, err := strconv.Atoi(answerID)
	if err != nil {
		return dto.AnswerPublicResponse{}, ErrInvalidIDFormat
	}
	if data.Text == nil && data.NewCorrectID == nil {
		return dto.AnswerPublicResponse{}, ErrNoFieldsToUpdate
	}

	isOwner, err := s.AnswerRepo.CheckIfAnswerOwner(ctx, aID, userID)
	if err != nil {
		return dto.AnswerPublicResponse{}, err
	}
	if !isOwner {
		slog.Warn("unauthorized mutation attempt", slog.Int("user_id", userID), slog.Int("target_answer_id", aID))
		return dto.AnswerPublicResponse{}, ErrNotAnAuthor
	}

	answer, err := s.AnswerRepo.UpdateAnswer(ctx, aID, data)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return dto.AnswerPublicResponse{}, ErrAnswerNotFound
		}
		if errors.Is(err, repositories.ErrInvalidCorrectID) {
			return dto.AnswerPublicResponse{}, ErrInvalidCorrectID
		}
		return dto.AnswerPublicResponse{}, err
	}

	slog.Info("answer updated successfully", slog.Int("answer_id", aID), slog.Int("updated_by", userID))

	return dto.NewAnswerResponse(answer), nil
}
