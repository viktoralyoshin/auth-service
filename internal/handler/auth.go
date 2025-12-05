package handler

import (
	"auth-service/internal/model"
	"auth-service/internal/service"
	"context"

	authpb "github.com/viktoralyoshin/playhub-proto/gen/go/auth"
	"github.com/viktoralyoshin/utils/pkg/jwt"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	service      *service.AuthService
	tokenManager jwt.TokenManager
}

func NewAuthHandler(service *service.AuthService, tokenManager jwt.TokenManager) *AuthHandler {
	return &AuthHandler{
		service:      service,
		tokenManager: tokenManager,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user := &model.User{Email: req.Email, Username: req.Username, Password: req.Password}

	createdUser, err := h.service.Register(ctx, user)
	if err != nil {
		return &authpb.RegisterResponse{}, err
	}

	accessToken, refreshToken, err := h.tokenManager.GenerateTokens(createdUser.Id.String(), string(createdUser.Role))
	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{
		UserId:       createdUser.Id.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
