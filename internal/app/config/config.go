package config

import (
	"io"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Pak3n/catalog-service/internal/app/config/section"
)

type Config struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

type LoadArgs struct {
	Output          io.Writer
	EnableSimpleLog bool
}

var Root Config

func createLogger(level zerolog.Level, output io.Writer) zerolog.Logger {
	return zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()
}

func Load(args LoadArgs) {
	// Настройка полей
	zerolog.TimestampFieldName = "timestamp"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC3339

	if args.EnableSimpleLog {
		args.Output = zerolog.ConsoleWriter{Out: args.Output}
	}

	// Глобальный логгер
	log.Logger = createLogger(zerolog.DebugLevel, args.Output)

	log.Debug().Msg("Logger initialized with Debug level")

	// Загрузка .env
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("No .env file found")
	}

	// Парсинг конфига
	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	level, err := zerolog.ParseLevel(Root.Monitor.LogLevel)
	if err != nil {
		log.Warn().Str("log_level", Root.Monitor.LogLevel).Msg("Unknown log level, using debug")
		level = zerolog.DebugLevel
	}

	log.Logger = createLogger(level, args.Output)
	log.Info().Str("log_level", level.String()).Msg("Logger re-initialized with config level")
}
