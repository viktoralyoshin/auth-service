package handler

import (
	"auth-service/internal/model"
	"auth-service/internal/service"
	"context"

	"github.com/rs/zerolog/log"
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

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	user, err := h.service.Login(ctx, req)
	if err != nil {
		log.Error().Msgf("Login, failed to login user: %v", err)
		return &authpb.LoginResponse{}, err
	}

	accessToken, refreshToken, err := h.tokenManager.GenerateTokens(user.Id.String(), string(user.Role))
	if err != nil {
		log.Error().Msgf("Login, failed to generate tokens: %v", err)
		return nil, err
	}

	return &authpb.LoginResponse{
		UserId:       user.Id.String(),
		Username:     user.Username,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user := &model.User{Email: req.Email, Username: req.Username, Password: req.Password}

	createdUser, err := h.service.Register(ctx, user)
	if err != nil {
		log.Error().Msgf("Register, failed to create user: %v", err)
		return &authpb.RegisterResponse{}, err
	}

	accessToken, refreshToken, err := h.tokenManager.GenerateTokens(createdUser.Id.String(), string(createdUser.Role))
	if err != nil {
		log.Error().Msgf("Register, failed to generate tokens: %v", err)
		return nil, err
	}

	return &authpb.RegisterResponse{
		UserId:       createdUser.Id.String(),
		Username:     createdUser.Username,
		Email:        createdUser.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *authpb.TokenRequest) (*authpb.TokenResponse, error) {
	userClaims, err := h.tokenManager.ParseToken(req.TokenStr)
	if err != nil {
		return nil, err
	}

	return &authpb.TokenResponse{
		UserId:   userClaims.UserId,
		UserRole: userClaims.UserRole,
	}, nil
}
