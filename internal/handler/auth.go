package handler

import (
	"auth-service/internal/model"
	"auth-service/internal/service"
	"context"
	"errors"

	jwtp "github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	authpb "github.com/viktoralyoshin/playhub-proto/gen/go/auth"
	"github.com/viktoralyoshin/utils/pkg/errs"
	"github.com/viktoralyoshin/utils/pkg/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	log.Info().Str("login", req.Login).Msg("login attempt")

	user, err := h.service.Login(ctx, req)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			log.Warn().Str("email", req.Login).Msg("login failed: user not found")
			return nil, status.Error(codes.NotFound, "user not found")
		}

		log.Warn().Err(err).Str("email", req.Login).Msg("login failed")
		return nil, status.Error(codes.InvalidArgument, "invalid email or password")
	}

	accessToken, refreshToken, err := h.tokenManager.GenerateTokens(user.Id.String(), string(user.Role))
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", user.Id.String()).
			Msg("failed to generate tokens during login")

		return nil, status.Error(codes.Internal, "internal error during token generation")
	}

	log.Info().
		Str("user_id", user.Id.String()).
		Msg("user logged in successfully")

	return &authpb.LoginResponse{
		UserId:       user.Id.String(),
		Username:     user.Username,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	log.Info().
		Str("email", req.Email).
		Str("username", req.Username).
		Msg("registration attempt")

	user := &model.User{Email: req.Email, Username: req.Username, Password: req.Password}

	createdUser, err := h.service.Register(ctx, user)
	if err != nil {
		if errors.Is(err, errs.ErrUserEmailExists) {
			log.Warn().Str("email", req.Email).Msg("registration failed: email already exists")
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		} else if errors.Is(err, errs.ErrUserUsernameExists) {
			log.Warn().Str("username", req.Username).Msg("registration failed: username already exists")
			return nil, status.Error(codes.AlreadyExists, "username already exists")
		}

		log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("failed to register user")

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	accessToken, refreshToken, err := h.tokenManager.GenerateTokens(createdUser.Id.String(), string(createdUser.Role))
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", createdUser.Id.String()).
			Msg("failed to generate tokens during registration")

		return nil, status.Error(codes.Internal, "internal error during token generation")
	}

	log.Info().
		Str("user_id", createdUser.Id.String()).
		Msg("user registered successfully")

	return &authpb.RegisterResponse{
		UserId:       createdUser.Id.String(),
		Username:     createdUser.Username,
		Email:        createdUser.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *authpb.TokenRequest) (*authpb.TokenResponse, error) {
	log.Debug().Msg("validating token")

	userClaims, err := h.tokenManager.ParseToken(req.TokenStr)
	if err != nil {
		log.Debug().Err(err).Msg("token validation failed")

		if errors.Is(err, jwtp.ErrTokenExpired) {
			return nil, status.Error(codes.Unauthenticated, "token expired")
		}

		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	log.Debug().Str("user_id", userClaims.UserId).Msg("token valid")

	return &authpb.TokenResponse{
		UserId:   userClaims.UserId,
		UserRole: userClaims.UserRole,
	}, nil
}
