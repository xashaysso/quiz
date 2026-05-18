package services

import (
	"context"
	"database/sql"
	"errors"
	"quiz/db/repositories"
	APIentities "quiz/entities/api"
	entities "quiz/entities/db"
	"quiz/entities/dto"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

var (
	ErrInvalidAnswerText = errors.New("field 'text' is required in answer body")
	ErrAnswerNotFound = errors.New("answer not found")
	ErrNoFieldsToUpdate = errors.New("fields text and/or correct_id are required in json body")
	ErrInvalidCorrectID = errors.New("invalid correct answer id")
)

type AnswerService struct {
	AnswerRepo repositories.AnswerRepository
}

func (s *AnswerService) CheckAnswer(ctx context.Context, questionID string, answerID int) (bool, error){
	qID, err := strconv.Atoi(questionID)
	if err != nil {
		return false, ErrInvalidIDFormat
	}

	correct, err := s.AnswerRepo.CheckAnswer(ctx, qID, answerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return false, ErrQuestionNotFound
		}
		return false, err
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
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return []entities.Answer{}, ErrQuestionNotFound
		}
		return []entities.Answer{}, err
	}

	return answers, nil
}

func (s *AnswerService) CreateAnswer(ctx context.Context, questionID string, data dto.CreateAnswerDTO) (APIentities.AnswerAPI, error) {
	qID, err := strconv.Atoi(questionID);
	if err != nil{
		return APIentities.AnswerAPI{}, ErrInvalidIDFormat;
	}

	if data.Text == "" {
		return APIentities.AnswerAPI{}, ErrInvalidAnswerText;
	}

	return s.AnswerRepo.CreateAnswer(ctx, qID, data)
}

func (s *AnswerService) GetAnswer(ctx context.Context, answerID string) (APIentities.AnswerAPI, error){
	aID, err := strconv.Atoi(answerID)
	if err != nil{
		return APIentities.AnswerAPI{}, err;
	}

	answer, err := s.AnswerRepo.GetAnswer(ctx, aID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return APIentities.AnswerAPI{}, ErrAnswerNotFound
		}
		return APIentities.AnswerAPI{}, err;
	}

	return answer, nil
}

func(s *AnswerService) DeleteAnswer(ctx context.Context, answerID string)(error) {
	aID, err := strconv.Atoi(answerID)
	if err != nil {
		return ErrInvalidIDFormat
	}
	err = s.AnswerRepo.DeleteAnswer(ctx, aID);
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		return ErrAnswerNotFound;
	}

	return nil
}

func (s *AnswerService) UpdateAnswer(ctx context.Context, answerID string, data dto.UpdateAnswerDTO)(APIentities.AnswerAPI, error) {
	aID, err := strconv.Atoi(answerID)
	if err != nil {
		return APIentities.AnswerAPI{}, ErrInvalidIDFormat
	}
	if data.Text == nil && data.NewCorrectID == nil{
		return APIentities.AnswerAPI{}, ErrNoFieldsToUpdate;
	}
	answer, err := s.AnswerRepo.UpdateAnswer(ctx, aID, data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return APIentities.AnswerAPI{}, ErrAnswerNotFound
		}
		errMsg := err.Error();
			if strings.Contains(errMsg, "not found"){
				return APIentities.AnswerAPI{}, ErrAnswerNotFound;
			}
			if strings.Contains(errMsg, "new correct answer id"){
				return APIentities.AnswerAPI{}, ErrInvalidCorrectID;
			}
		return APIentities.AnswerAPI{}, err
	}

	return answer, nil
}