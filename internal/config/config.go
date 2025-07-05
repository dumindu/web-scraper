package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type ConfDB struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DBName   string `env:"DB_NAME,required"`
	Debug    bool   `env:"DB_DEBUG,required"`
}

func NewDB() *ConfDB {
	var c ConfDB
	if err := env.Parse(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
