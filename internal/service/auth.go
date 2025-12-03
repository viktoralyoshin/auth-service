package service

import (
	"context"

	authpb "github.com/viktoralyoshin/playhub-proto/gen/go/auth"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	return &authpb.RegisterResponse{
		UserId: "313",
	}, nil
}
