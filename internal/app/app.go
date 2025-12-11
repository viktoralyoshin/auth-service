package app

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/grpc"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/viktoralyoshin/utils/pkg/jwt"
)

func Start(cfg *config.Config) {
	// Инициализация БД
	db, err := database.Init(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to database")
	}

	// Graceful shutdown для базы
	defer func() {
		if err := db.Close(); err != nil {
			// В defer используем Error, а не Fatal
			log.Error().Err(err).Msg("failed to close database connection")
		}
	}()

	// Миграции
	if err := database.Migrate(db); err != nil {
		log.Fatal().Err(err).Msg("migration failed")
	}

	// Token Manager
	tokenManager, err := jwt.NewManager(cfg.JWTSigningKey, cfg.AccessTokenTtl, cfg.RefreshTokenTtl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize jwt token manager")
	}

	// Listener
	// Собираем адрес явно для лога
	addr := ":" + cfg.GRPCPort
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Str("addr", addr).Msg("failed to listen tcp")
	}

	// gRPC Server Init
	s := grpc.Init(db, tokenManager)

	// Инфо о запуске
	log.Info().
		Str("port", cfg.GRPCPort).
		Str("service", "auth-service").
		Msg("gRPC server started")

	// Блокирующий запуск
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("gRPC server stopped unexpectedly")
	}
}
