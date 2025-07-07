package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
)

type Conf struct {
	Server ConfServer
	DB     ConfDB
	Mailer MailerConf
}

type ConfServer struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
}

type ConfDB struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DBName   string `env:"DB_NAME,required"`
	Debug    bool   `env:"DB_DEBUG,required"`
}

type MailerConf struct {
	Host        string `env:"MAILER_HOST,required"`
	Port        int    `env:"MAILER_PORT,required"`
	User        string `env:"MAILER_USER,required"`
	Pass        string `env:"MAILER_PASS,required"`
	FromNoReply string `env:"MAILER_FROM_NO_REPLY,required"`
	WebsiteHost string `env:"MAILER_WEBSITE_HOST,required"`
}

func New() *Conf {
	var c Conf
	if err := env.Parse(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}

func NewDB() *ConfDB {
	var c ConfDB
	if err := env.Parse(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
