package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"

	"github.com/rs/zerolog/log"
	"github.com/viktoralyoshin/utils/pkg/logger"
)

func main() {
	logger.Setup()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Msgf("failed to load config: %v", err)
	}

	app.Start(cfg)
}
