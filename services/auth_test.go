package services_test

import (
	"context"
	"quiz/db/repositories"
	mock_repositories "quiz/db/repositories/mocks"
	entities "quiz/entities/db"
	"quiz/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {

	type testCase struct {
		name        string
		username    string
		password    string
		prepareMock func(m *mock_repositories.MockSessionRepository, mu *mock_repositories.MockUserRepository)

		wantErr error
	}

	cases := []testCase{
		{
			name: "err_invalid_username",

			username:    "ab",
			password:    "12345",
			prepareMock: func(m *mock_repositories.MockSessionRepository, mu *mock_repositories.MockUserRepository) {},

			wantErr: services.ErrInvalidUsername,
		},
		{
			name: "err_invalid_password",

			username:    "abcdef",
			password:    "1",
			prepareMock: func(m *mock_repositories.MockSessionRepository, mu *mock_repositories.MockUserRepository) {},

			wantErr: services.ErrInvalidPassword,
		},
		{
			name: "err_user_already_exists",

			username: "abcdef",
			password: "12345",
			prepareMock: func(m *mock_repositories.MockSessionRepository, mu *mock_repositories.MockUserRepository) {
				mu.EXPECT().CreateUser(gomock.Any(), "abcdef", gomock.Any()).Return(entities.User{}, repositories.ErrUserAlreadyExists)
			},

			wantErr: services.ErrUserAlreadyExists,
		},
		{
			name: "success_register",

			username: "abcdef",
			password: "12345",
			prepareMock: func(m *mock_repositories.MockSessionRepository, mu *mock_repositories.MockUserRepository) {
				validHash, _ := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)

				fakeUser := entities.User{
					ID:           1,
					Username:     "abcdef",
					PasswordHash: string(validHash),
				}

				mu.EXPECT().CreateUser(gomock.Any(), "abcdef", gomock.Any()).Return(fakeUser, nil)
				mu.EXPECT().GetByUsername(gomock.Any(), "abcdef").Return(fakeUser, nil)
				m.EXPECT().Set(gomock.Any(), gomock.Any(), 1, 24*time.Hour).Return(nil)
			},

			wantErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSessionRepository := mock_repositories.NewMockSessionRepository(ctrl)
			mockUserRepository := mock_repositories.NewMockUserRepository(ctrl)
			tc.prepareMock(mockSessionRepository, mockUserRepository)

			service := services.NewAuthService(mockUserRepository, mockSessionRepository)
			ctx := context.Background()

			_, _, err := service.Register(ctx, tc.username, tc.password)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestAuthService_Login(t *testing.T) {

	type testCase struct {
		name        string
		username    string
		password    string
		prepareMock func(m *mock_repositories.MockSessionRepository, mu *mock_repositories.MockUserRepository)

		wantErr error
	}

	cases := []testCase{
		{
			name: "err_invalid_credentials",

			username: "abcdef",
			password: "12345",
			prepareMock: func(m *mock_repositories.MockSessionRepository, mu *mock_repositories.MockUserRepository) {
				mu.EXPECT().GetByUsername(gomock.Any(), "abcdef").Return(entities.User{}, repositories.ErrRecordNotFound)
			},

			wantErr: services.ErrWrongCredentials,
		},
		{
			name: "success_login",

			username: "abcdef",
			password: "12345",
			prepareMock: func(m *mock_repositories.MockSessionRepository, mu *mock_repositories.MockUserRepository) {
				validHash, _ := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)

				fakeUser := entities.User{
					ID:           1,
					Username:     "abcdef",
					PasswordHash: string(validHash),
				}

				mu.EXPECT().GetByUsername(gomock.Any(), "abcdef").Return(fakeUser, nil)
				m.EXPECT().Set(gomock.Any(), gomock.Any(), 1, 24*time.Hour).Return(nil)
			},

			wantErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSessionRepository := mock_repositories.NewMockSessionRepository(ctrl)
			mockUserRepository := mock_repositories.NewMockUserRepository(ctrl)
			tc.prepareMock(mockSessionRepository, mockUserRepository)

			service := services.NewAuthService(mockUserRepository, mockSessionRepository)
			ctx := context.Background()

			_, err := service.Login(ctx, tc.username, tc.password)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestAuthService_Logout(t *testing.T) {

	type testCase struct {
		name        string
		token       string
		prepareMock func(m *mock_repositories.MockSessionRepository, mu *mock_repositories.MockUserRepository)

		wantErr error
	}

	cases := []testCase{
		{
			name:  "success_logout",
			token: "123",

			prepareMock: func(m *mock_repositories.MockSessionRepository, mu *mock_repositories.MockUserRepository) {
				m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},

			wantErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSessionRepository := mock_repositories.NewMockSessionRepository(ctrl)
			mockUserRepository := mock_repositories.NewMockUserRepository(ctrl)
			tc.prepareMock(mockSessionRepository, mockUserRepository)

			service := services.NewAuthService(mockUserRepository, mockSessionRepository)
			ctx := context.Background()

			err := service.Logout(ctx, tc.token)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
