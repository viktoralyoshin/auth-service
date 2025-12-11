package config

import (
	"os"
	"time"
)

type Config struct {
	DBUser          string
	DBHost          string
	DBPassword      string
	DBPort          string
	DBName          string
	GRPCPort        string
	JWTSigningKey   string
	Env             string
	RefreshTokenTtl time.Duration
	AccessTokenTtl  time.Duration
}

func Load() (*Config, error) {

	accessTokenTtl, refreshToketTtl := os.Getenv("ACCESS_TOKEN_TTL"), os.Getenv("REFRESH_TOKEN_TTL")

	accessDuration, err := time.ParseDuration(accessTokenTtl)
	if err != nil {
		return nil, err
	}

	refreshDuration, err := time.ParseDuration(refreshToketTtl)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		DBUser:          os.Getenv("DBUser"),
		DBHost:          os.Getenv("DBHost"),
		DBPassword:      os.Getenv("DBPassword"),
		DBPort:          os.Getenv("DBPort"),
		DBName:          os.Getenv("DBName"),
		GRPCPort:        os.Getenv("GRPCPort"),
		JWTSigningKey:   os.Getenv("JWT_SIGNING_KEY"),
		Env:             os.Getenv("ENV"),
		RefreshTokenTtl: refreshDuration,
		AccessTokenTtl:  accessDuration,
	}

	return cfg, nil
}
