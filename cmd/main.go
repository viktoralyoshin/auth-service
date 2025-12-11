package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"

	"github.com/rs/zerolog/log"
	"github.com/viktoralyoshin/utils/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	logger.Setup(cfg.Env)

	app.Start(cfg)
}
