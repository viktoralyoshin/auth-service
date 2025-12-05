package grpc

import (
	"auth-service/internal/handler"
	"auth-service/internal/service"
	"auth-service/internal/storage"
	"database/sql"

	authpb "github.com/viktoralyoshin/playhub-proto/gen/go/auth"
	"github.com/viktoralyoshin/utils/pkg/jwt"
	"google.golang.org/grpc"
)

func Init(db *sql.DB, tokenManager jwt.TokenManager) *grpc.Server {
	s := grpc.NewServer()

	userRepo := storage.NewUserRepo(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, tokenManager)

	authpb.RegisterAuthServiceServer(s, authHandler)

	return s
}
