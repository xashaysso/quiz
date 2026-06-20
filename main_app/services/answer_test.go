package services_test

import (
	"context"
	"quiz/db/repositories"
	mock_repositories "quiz/db/repositories/mocks"
	entities "quiz/entities/db"
	"quiz/entities/dto"
	"quiz/services"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAnswerService_CheckAnswer(t *testing.T) {

	type testCase struct {
		name        string
		questionID  string
		answerID    int
		prepareMock func(m *mock_repositories.MockAnswerRepository)

		wantResult bool
		wantErr    error
	}

	cases := []testCase{
		{
			name:       "success_correct_answer",
			questionID: "10",
			answerID:   2,
			prepareMock: func(m *mock_repositories.MockAnswerRepository) {
				m.EXPECT().CheckAnswer(gomock.Any(), 10, 2).Return(true, nil)
			},

			wantResult: true,
			wantErr:    nil,
		},
		{
			name:       "question_not_found",
			questionID: "999",
			answerID:   1,
			prepareMock: func(m *mock_repositories.MockAnswerRepository) {
				m.EXPECT().CheckAnswer(gomock.Any(), 999, 1).Return(false, repositories.ErrRecordNotFound)
			},

			wantResult: false,
			wantErr:    services.ErrQuestionNotFound,
		},
		{
			name:        "invalid_id_format",
			questionID:  "abc",
			answerID:    1,
			prepareMock: func(m *mock_repositories.MockAnswerRepository) {},

			wantResult: false,
			wantErr:    services.ErrInvalidIDFormat,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAnswerRepo := mock_repositories.NewMockAnswerRepository(ctrl)
			tc.prepareMock(mockAnswerRepo)

			service := services.NewAnswerService(mockAnswerRepo, nil)
			ctx := context.Background()

			res, err := service.CheckAnswer(ctx, "random_session", tc.questionID, tc.answerID)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantResult, res)
		})
	}
}

func TestAnswerService_CreateAnswer(t *testing.T) {
	type testCase struct {
		name        string
		questionID  string
		dto         dto.CreateAnswerDTO
		userID      int
		prepareMock func(m *mock_repositories.MockAnswerRepository, tx *mock_repositories.MockTransactionManager)
		wantErr     error
	}

	cases := []testCase{
		{
			name:       "error_not_an_author",
			questionID: "10",
			dto:        dto.CreateAnswerDTO{Text: "Correct answer", IsCorrect: true},
			userID:     5,
			prepareMock: func(m *mock_repositories.MockAnswerRepository, tx *mock_repositories.MockTransactionManager) {
				m.EXPECT().CheckIfQuestionOwner(gomock.Any(), 10, 5).Return(false, nil)
			},

			wantErr: services.ErrNotAnAuthor,
		},
		{
			name:       "success_create_answer",
			questionID: "10",
			dto:        dto.CreateAnswerDTO{Text: "Cool answer", IsCorrect: true},
			userID:     5,
			prepareMock: func(m *mock_repositories.MockAnswerRepository, tx *mock_repositories.MockTransactionManager) {
				m.EXPECT().CheckIfQuestionOwner(gomock.Any(), 10, 5).Return(true, nil)
				tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx pgx.Tx) error) error {
					return fn(nil)
				})
				m.EXPECT().CreateAnswer(gomock.Any(), gomock.Any(), 10, "Cool answer", true).Return(1, nil)
			},

			wantErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAnswerRepo := mock_repositories.NewMockAnswerRepository(ctrl)
			mockTx := mock_repositories.NewMockTransactionManager(ctrl)
			tc.prepareMock(mockAnswerRepo, mockTx)

			service := services.NewAnswerService(mockAnswerRepo, mockTx)
			ctx := context.Background()

			_, err := service.CreateAnswer(ctx, tc.questionID, tc.dto, tc.userID)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestAnswerService_UpdateAnswer(t *testing.T) {
	type testCase struct {
		name        string
		answerID    string
		dto         dto.UpdateAnswerDTO
		userID      int
		prepareMock func(m *mock_repositories.MockAnswerRepository)
		wantErr     error
	}

	strPtr := func(s string) *string { return &s }

	cases := []testCase{
		{
			name:     "error_not_an_author",
			answerID: "10",
			dto:      dto.UpdateAnswerDTO{Text: strPtr("Updated answer")},
			userID:   5,
			prepareMock: func(m *mock_repositories.MockAnswerRepository) {
				m.EXPECT().CheckIfAnswerOwner(gomock.Any(), 10, 5).Return(false, nil)
			},

			wantErr: services.ErrNotAnAuthor,
		},
		{
			name:        "error_no_fields_to_update",
			answerID:    "10",
			dto:         dto.UpdateAnswerDTO{Text: nil, NewCorrectID: nil},
			userID:      5,
			prepareMock: func(m *mock_repositories.MockAnswerRepository) {},

			wantErr: services.ErrNoFieldsToUpdate,
		},
		{
			name:     "success_update_answer",
			answerID: "10",
			dto:      dto.UpdateAnswerDTO{Text: strPtr("New text")},
			userID:   5,
			prepareMock: func(m *mock_repositories.MockAnswerRepository) {
				m.EXPECT().CheckIfAnswerOwner(gomock.Any(), 10, 5).Return(true, nil)
				m.EXPECT().UpdateAnswer(gomock.Any(), 10, gomock.Any()).Return(entities.Answer{ID: 10, Text: "New text"}, nil)
			},

			wantErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAnswerRepo := mock_repositories.NewMockAnswerRepository(ctrl)
			mockTx := mock_repositories.NewMockTransactionManager(ctrl)
			tc.prepareMock(mockAnswerRepo)

			service := services.NewAnswerService(mockAnswerRepo, mockTx)
			ctx := context.Background()

			_, err := service.UpdateAnswer(ctx, tc.answerID, tc.dto, tc.userID)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestAnswerService_DeleteAnswer(t *testing.T) {
	type testCase struct {
		name        string
		answerID    string
		userID      int
		prepareMock func(m *mock_repositories.MockAnswerRepository)
		wantErr     error
	}

	cases := []testCase{
		{
			name:     "error_not_an_author",
			answerID: "10",
			userID:   5,
			prepareMock: func(m *mock_repositories.MockAnswerRepository) {
				m.EXPECT().CheckIfAnswerOwner(gomock.Any(), 10, 5).Return(false, nil)
			},

			wantErr: services.ErrNotAnAuthor,
		},
		{
			name:        "error_invalid_id_format",
			answerID:    "abc",
			userID:      5,
			prepareMock: func(m *mock_repositories.MockAnswerRepository) {},

			wantErr: services.ErrInvalidIDFormat,
		},
		{
			name:     "success_delete_answer",
			answerID: "10",
			userID:   5,
			prepareMock: func(m *mock_repositories.MockAnswerRepository) {
				m.EXPECT().CheckIfAnswerOwner(gomock.Any(), 10, 5).Return(true, nil)
				m.EXPECT().DeleteAnswer(gomock.Any(), 10).Return(nil)
			},

			wantErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAnswerRepo := mock_repositories.NewMockAnswerRepository(ctrl)
			mockTx := mock_repositories.NewMockTransactionManager(ctrl)
			tc.prepareMock(mockAnswerRepo)

			service := services.NewAnswerService(mockAnswerRepo, mockTx)
			ctx := context.Background()

			err := service.DeleteAnswer(ctx, tc.answerID, tc.userID)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
