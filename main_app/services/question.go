package services

import (
	"context"
	"errors"
	"log/slog"
	"quiz/db/repositories"
	entities "quiz/entities/db"
	"quiz/entities/dto"
	"strconv"
)

type QuestionService struct {
	QuestionRepo repositories.QuestionRepository
	AnswerRepo   repositories.AnswerRepository
	TxManager    repositories.TransactionManager
}

func NewQuestionService(qRepo repositories.QuestionRepository, aRepo repositories.AnswerRepository, txm repositories.TransactionManager) QuestionServiceInterface {
	return &QuestionService{
		QuestionRepo: qRepo,
		AnswerRepo:   aRepo,
		TxManager:    txm,
	}
}

func (s *QuestionService) CreateQuestion(ctx context.Context, quizID string, body dto.CreateQuestionDTO, userID int) (dto.QuestionResponse, error) {
	if len(body.Answers) == 0 {
		return dto.QuestionResponse{}, ErrNoQuestionAnswers
	}
	qID, err := strconv.Atoi(quizID)
	if err != nil {
		return dto.QuestionResponse{}, ErrInvalidIDFormat
	}

	isOwner, err := s.QuestionRepo.CheckIfQuizOwner(ctx, qID, userID)
	if err != nil {
		return dto.QuestionResponse{}, err
	}
	if !isOwner {
		slog.Warn("unauthorized mutation attempt", slog.Int("user_id", userID), slog.Int("target_quiz_id", qID))
		return dto.QuestionResponse{}, ErrNotAnAuthor
	}

	var createdQuestionID int
	var answers []entities.Answer

	err = s.TxManager.WithinTransaction(ctx, func(tx any) error {
		id, err := s.QuestionRepo.CreateQuestion(ctx, tx, qID, body)
		if err != nil {
			if errors.Is(err, repositories.ErrRecordNotFound) {
				return ErrQuizNotFound
			}
			return err
		}

		createdQuestionID = id

		answers = make([]entities.Answer, len(body.Answers))

		for i, ansDTO := range body.Answers {
			ansID, err := s.AnswerRepo.CreateAnswer(ctx, tx, createdQuestionID, ansDTO.Text, ansDTO.IsCorrect)
			if err != nil {
				return err
			}

			answers[i] = entities.Answer{
				ID:         ansID,
				QuestionID: createdQuestionID,
				Text:       ansDTO.Text,
				IsCorrect:  ansDTO.IsCorrect,
			}
		}

		return nil
	})
	if err != nil {
		return dto.QuestionResponse{}, err
	}

	question := entities.Question{
		ID:     createdQuestionID,
		Text:   body.Text,
		QuizID: qID,
	}

	slog.Info("question with answers created successfully", slog.Int("quiz_id", qID), slog.Int("answers_count", len(body.Answers)))

	return dto.NewQuestionResponse(question, answers), nil
}

func (s *QuestionService) ListQuestions(ctx context.Context, quizID string) ([]dto.QuestionResponse, error) {
	qID, err := strconv.Atoi(quizID)
	if err != nil {
		return []dto.QuestionResponse{}, ErrInvalidIDFormat
	}
	questions, err := s.QuestionRepo.GetQuizQuestions(ctx, qID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return []dto.QuestionResponse{}, ErrQuestionNotFound
		}
		return []dto.QuestionResponse{}, err
	}
	qIDs := make([]int, len(questions))
	for i, q := range questions {
		qIDs[i] = q.ID
	}

	answers, err := s.AnswerRepo.GetAnswersByQuestionIDs(ctx, qIDs)
	answerMap := make(map[int][]entities.Answer)
	for _, ans := range answers {
		answerMap[ans.QuestionID] = append(answerMap[ans.QuestionID], ans)
	}

	response := make([]dto.QuestionResponse, len(questions))
	for i, q := range questions {
		currentAnswers := answerMap[q.ID]
		response[i] = dto.NewQuestionResponse(q, currentAnswers)
	}

	return response, nil
}

func (s *QuestionService) GetQuestion(ctx context.Context, questionID string) (dto.QuestionResponse, error) {
	qID, err := strconv.Atoi(questionID)
	if err != nil {
		return dto.QuestionResponse{}, ErrInvalidIDFormat
	}
	question, err := s.QuestionRepo.GetQuestion(ctx, qID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return dto.QuestionResponse{}, ErrQuestionNotFound
		}
		return dto.QuestionResponse{}, err
	}

	answers, err := s.AnswerRepo.GetQuizAnswers(ctx, qID)
	if err != nil {
		return dto.QuestionResponse{}, ErrAnswerNotFound
	}

	return dto.NewQuestionResponse(question, answers), nil
}

func (s *QuestionService) UpdateQuestion(ctx context.Context, questionID string, data dto.UpdateQuestionDTO, userID int) (dto.QuestionResponse, error) {
	qID, err := strconv.Atoi(questionID)
	if err != nil {
		return dto.QuestionResponse{}, ErrInvalidIDFormat
	}

	isOwner, err := s.QuestionRepo.CheckIfQuestionOwner(ctx, qID, userID)
	if err != nil {
		return dto.QuestionResponse{}, err
	}
	if !isOwner {
		slog.Warn("unauthorized mutation attempt", slog.Int("user_id", userID), slog.Int("target_question_id", qID))
		return dto.QuestionResponse{}, ErrNotAnAuthor
	}

	question, err := s.QuestionRepo.UpdateQuestion(ctx, qID, data)

	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return dto.QuestionResponse{}, ErrQuestionNotFound
		}
		return dto.QuestionResponse{}, err
	}

	answers, err := s.AnswerRepo.GetQuizAnswers(ctx, qID)
	if err != nil {
		return dto.QuestionResponse{}, ErrAnswerNotFound
	}

	slog.Info("question updated successfully", slog.Int("question_id", qID), slog.Int("updated_by", userID))

	return dto.NewQuestionResponse(question, answers), nil
}

func (s *QuestionService) DeleteQuestion(ctx context.Context, questionID string, userID int) error {
	qID, err := strconv.Atoi(questionID)
	if err != nil {
		return ErrInvalidIDFormat
	}

	isOwner, err := s.QuestionRepo.CheckIfQuestionOwner(ctx, qID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		slog.Warn("unauthorized mutation attempt", slog.Int("user_id", userID), slog.Int("target_question_id", qID))
		return ErrNotAnAuthor
	}

	err = s.QuestionRepo.DeleteQuestion(ctx, qID)

	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return ErrQuestionNotFound
		}
		return err
	}

	slog.Info("question deleted successfully", slog.Int("question_id", qID), slog.Int("deleted_by", userID))

	return nil
}
