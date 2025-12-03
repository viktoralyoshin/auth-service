package grpc

import (
	"auth-service/internal/service"

	authpb "github.com/viktoralyoshin/playhub-proto/gen/go/auth"
	"google.golang.org/grpc"
)

func Init() *grpc.Server {
	s := grpc.NewServer()

	authService := service.NewAuthService()

	authpb.RegisterAuthServiceServer(s, authService)

	return s
}
