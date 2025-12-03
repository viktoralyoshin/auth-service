package config

import (
	"os"
)

type Config struct {
	DBUser     string
	DBHost     string
	DBPassword string
	DBPort     string
	DBName     string
	GRPCPort   string
}

func Load() (*Config, error) {

	cfg := &Config{
		DBUser:     os.Getenv("DBUser"),
		DBHost:     os.Getenv("DBHost"),
		DBPassword: os.Getenv("DBPassword"),
		DBPort:     os.Getenv("DBPort"),
		DBName:     os.Getenv("DBName"),
		GRPCPort:   os.Getenv("GRPCPort"),
	}

	return cfg, nil
}
