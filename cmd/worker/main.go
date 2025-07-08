package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hibiken/asynq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"web-scraper.dev/internal/config"
	"web-scraper.dev/internal/utils/logger"
	"web-scraper.dev/internal/workers"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

func main() {
	c := config.NewWorkerConf()
	l := logger.New(true)
	db, err := gormDB(&c.DB)
	if err != nil {
		l.Fatal().Err(err).Msg("DB connection start failure")
	}

	redisClusterAddresses := strings.Split(c.RedisHosts, ",")
	if len(redisClusterAddresses) == 0 {
		l.Fatal().Msg("Redis cluster connection failure")
	}
	redisConnOpt := asynq.RedisClusterClientOpt{
		Addrs: redisClusterAddresses,
	}

	scrapeWorker := workers.NewScrapeWorker(redisConnOpt, db, l)

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		l.Info().Msg("Shutting down worker")

		if err := scrapeWorker.Stop(); err != nil {
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

	l.Info().Msg("Starting worker")
	if err := scrapeWorker.Start(); err != nil {
		l.Fatal().Err(err).Msg("Worker start failure")
	}

	<-closed
	l.Info().Msgf("Worker shutdown successfully")
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
