package service

import (
	"auth-service/internal/model"
	"auth-service/internal/storage"
	"context"

	authpb "github.com/viktoralyoshin/playhub-proto/gen/go/auth"
	"github.com/viktoralyoshin/utils/pkg/errs"
	"github.com/viktoralyoshin/utils/pkg/hasher"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	repo *storage.UserRepo
}

func NewAuthService(repo *storage.UserRepo) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Register(ctx context.Context, user *model.User) (*model.User, error) {
	existingEmail, err := s.repo.GetUserByLogin(ctx, user.Email)
	if err != nil {
		return &model.User{}, err
	}

	if existingEmail != nil {
		return nil, errs.ErrUserEmailExists
	}

	existingUsername, err := s.repo.GetUserByLogin(ctx, user.Username)
	if err != nil {
		return &model.User{}, err
	}

	if existingUsername != nil {
		return nil, errs.ErrUserUsernameExists
	}

	passwordHash, err := hasher.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = passwordHash

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
