package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"web-scraper.dev/internal/api/router"
	"web-scraper.dev/internal/config"
	"web-scraper.dev/internal/utils/logger"
	"web-scraper.dev/internal/utils/validator"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

func main() {
	c := config.New()
	l := logger.New(c.Server.Debug)
	v := validator.New()
	db, err := gormDB(&c.DB)
	if err != nil {
		l.Fatal().Err(err).Msg("DB connection start failure")
	}

	r := router.New(c.Server.TimeoutRead, c.Server.TimeoutWrite, db, l, v)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		l.Info().Msgf("Shutting down server %v", s.Addr)

		ctx, cancel := context.WithTimeout(context.Background(), s.IdleTimeout)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			l.Error().Err(err).Msg("Server shutdown failure")
		}

		sqlDB, err := db.DB()
		if err == nil {
			if err = sqlDB.Close(); err != nil {
				l.Error().Err(err).Msg("DB connection closing failure")
			}
		}

		close(closed)
	}()

	l.Info().Msgf("Starting server %v", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
	l.Info().Msgf("Server shutdown successfully")
}

func gormDB(conf *config.ConfDB) (*gorm.DB, error) {
	var logLevel gormlogger.LogLevel
	if conf.Debug {
		logLevel = gormlogger.Info
	} else {
		logLevel = gormlogger.Error
	}

	dbString := fmt.Sprintf(fmtDBString, conf.Host, conf.Username, conf.Password, conf.DBName, conf.Port)
	return gorm.Open(postgres.Open(dbString), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})
}
