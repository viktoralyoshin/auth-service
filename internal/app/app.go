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

	db, err := database.Init(cfg)
	if err != nil {
		log.Fatal().Msgf("unable to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal().Msgf("failed to close database connection: %v", err)
		}
	}()

	if err := database.Migrate(db); err != nil {
		log.Fatal().Msgf("migration failed: %v", err)
	}

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	tokenManager, err := jwt.NewManager(cfg.JWTSigningKey, cfg.AccessTokenTtl, cfg.RefreshTokenTtl)
	if err != nil {
		log.Fatal().Msgf("initialization jwt token manager failed: %v", err)
	}

	s := grpc.Init(db, tokenManager)

	log.Info().Msgf("Auth Service running on :%s", cfg.GRPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}
