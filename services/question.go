package services

import (
	"context"
	"database/sql"
	"errors"
	"quiz/db/repositories"
	entities "quiz/entities/api"
	"quiz/entities/dto"
	"strconv"

	"github.com/jackc/pgx/v5"
)

var (
	ErrInvalidIDFormat = errors.New("invalid 'ID' field format")
	ErrNoQuestionAnswers = errors.New("no field 'answers' provided in json")
	ErrQuestionNotFound = errors.New("question not found")
)

type QuestionService struct {
	QuestionRepo repositories.QuestionRepository
}

func (s *QuestionService) CreateQuestion (ctx context.Context, quizID string, body dto.CreateQuestionDTO, userID int) (entities.QuestionAPI, error){
	if len(body.Answers) == 0 {
		return entities.QuestionAPI{}, ErrNoQuestionAnswers
	}
	qID, err := strconv.Atoi(quizID)
	if err != nil {
		return entities.QuestionAPI{}, ErrInvalidIDFormat
	}

	isOwner, err := s.QuestionRepo.CheckIfQuizOwner(ctx, qID, userID)
	if err != nil {
		return entities.QuestionAPI{}, err
	}
	if !isOwner {
		return entities.QuestionAPI{}, ErrNotAnAuthor
	}

	question, err := s.QuestionRepo.CreateQuestion(ctx, qID, body)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return entities.QuestionAPI{}, ErrQuizNotFound
		}
		return entities.QuestionAPI{}, err
	}

	return question, nil
}

func (s *QuestionService) ListQuestions(ctx context.Context, quizID string)([]entities.QuestionAPI, error) {
	qID, err := strconv.Atoi(quizID)
	if err != nil {
		return []entities.QuestionAPI{}, ErrInvalidIDFormat
	}
	questions, err := s.QuestionRepo.GetQuizQuestions(ctx, qID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return []entities.QuestionAPI{}, ErrQuestionNotFound
		}
		return []entities.QuestionAPI{}, err
	}
	return questions, nil
}

func (s *QuestionService) GetQuestion(ctx context.Context, questionID string)(entities.QuestionAPI, error) {
	qID, err := strconv.Atoi(questionID)
	if err != nil {
		return entities.QuestionAPI{}, ErrInvalidIDFormat
	}
	question, err := s.QuestionRepo.GetQuestion(ctx, qID);
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return entities.QuestionAPI{}, ErrQuestionNotFound
		}
		return entities.QuestionAPI{}, err
	}
	return question, nil
}

func (s *QuestionService) UpdateQuestion(ctx context.Context, questionID string, data dto.UpdateQuestionDTO, userID int) (entities.QuestionAPI ,error){
	qID, err := strconv.Atoi(questionID)
	if err != nil {
		return entities.QuestionAPI{}, ErrInvalidIDFormat
	}

	isOwner, err := s.QuestionRepo.CheckIfQuestionOwner(ctx, qID, userID)
	if err != nil {
		return entities.QuestionAPI{}, err
	}
	if !isOwner {
		return entities.QuestionAPI{}, ErrNotAnAuthor
	}

	question, err := s.QuestionRepo.UpdateQuestion(ctx, qID, data)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return entities.QuestionAPI{}, ErrQuestionNotFound
		}
		return entities.QuestionAPI{}, err
	}
	return question, nil
}

func (s *QuestionService) DeleteQuestion(ctx context.Context, questionID string, userID int) (error){
	qID, err := strconv.Atoi(questionID)
	if err != nil {
		return ErrInvalidIDFormat
	}

	isOwner, err := s.QuestionRepo.CheckIfQuestionOwner(ctx, qID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrNotAnAuthor
	}

	err = s.QuestionRepo.DeleteQuestion(ctx, qID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return ErrQuestionNotFound
		}
		return err
	}
	return nil
}