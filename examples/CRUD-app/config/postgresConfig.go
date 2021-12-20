package config

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST,required"`
	PORT     string `env:"POSTGRES_PORT,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	User     string `env:"POSTGRES_USER,required"`
	DbName   string `env:"POSTGRES_DBNAME,required"`
}

func NewPostgresConfig() *PostgresConfig {
	cfg := &PostgresConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Couldn't parse postgres config: %v", err)
		return nil
	}

	return cfg
}
