package transport

import (
	"auth/pkg/authv1"
	"auth/service"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	authService service.AuthServiceInterface
}

func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user, token, err := h.authService.Register(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, service.ErrInvalidUsername) || errors.Is(err, service.ErrInvalidPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.RegisterResponse{
		Token:  token,
		UserId: int64(user.ID),
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	token, err := h.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, service.ErrWrongCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (h *AuthHandler) CheckSession(ctx context.Context, req *authv1.CheckSessionRequest) (*authv1.CheckSessionResponse, error) {
	userID, err := h.authService.CheckSession(ctx, req.GetSessionId())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "session is invalid or expired")
	}

	return &authv1.CheckSessionResponse{
		UserId: int64(userID),
	}, nil
}

func (h *AuthHandler) DeleteSession(ctx context.Context, req *authv1.DeleteSessionRequest) (*authv1.DeleteSessionResponse, error) {
	err := h.authService.Logout(ctx, req.GetSessionId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete session")
	}

	return &authv1.DeleteSessionResponse{
		Success: true,
	}, nil
}
