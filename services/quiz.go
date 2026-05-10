package services

import (
	"context"
	"errors"
	"quiz/db/repositories"
	entities "quiz/entities/db"
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

func (s *QuizService) DeleteQuiz(ctx context.Context, quizID, userID int)(error) {
	quiz, err := s.QuizRepo.GetQuizByID(ctx, quizID)
	if err != nil {
		return ErrQuizNotFound
	}
	if quiz.CreatorID != userID {
		return ErrNotAnAuthor
	}
	return s.QuizRepo.DeleteQuiz(ctx, quizID)
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

func (s *QuizService) UpdateQuiz(ctx context.Context, quizID int, name, description *string, userID int) (entities.Quiz, error){
	if name == nil && description == nil{
		return entities.Quiz{}, ErrNoRequiredFields
	}
	quiz, err := s.QuizRepo.GetQuizByID(ctx, quizID)
	if err != nil {
		return entities.Quiz{}, ErrQuizNotFound
	}
	if quiz.CreatorID != userID {
		return entities.Quiz{}, ErrNotAnAuthor
	}

	newQuiz, err := s.QuizRepo.UpdateQuiz(ctx, quizID, name, description);
	if err != nil {
		return entities.Quiz{}, err
	}
	return newQuiz, nil
}
