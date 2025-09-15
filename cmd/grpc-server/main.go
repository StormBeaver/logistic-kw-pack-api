package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	loggerApp "github.com/stormbeaver/logistic-pack-api/internal/logger"
	"github.com/stormbeaver/logistic-pack-api/internal/server"

	"github.com/stormbeaver/logistic-pack-api/internal/tracer"

	"github.com/stormbeaver/logistic-pack-api/internal/database"

	"github.com/stormbeaver/logistic-pack-api/internal/config"
)

var (
	batchSize uint = 2
	logger    zerolog.Logger
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}

	cfg := config.GetConfigInstance()

	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()

	log.Info().
		Str("version", cfg.Project.Version).
		Str("commitHash", cfg.Project.CommitHash).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)

	logger = loggerApp.LogInit(cfg.Project.Debug)

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	initCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewPostgres(initCtx, dsn, cfg.Database.Driver, &cfg.Database.Connections)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed init postgres")
	}
	defer db.Close()

	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			logger.Error().Err(err).Msg("Migration failed")

			return
		}
	}

	tracing, err := tracer.NewTracer(&cfg, logger)
	if err != nil {
		logger.Error().Err(err).Msg("Failed init tracing")

		return
	}
	defer tracing.Close()

	if err := server.NewGrpcServer(db, batchSize).Start(&cfg, logger); err != nil {
		logger.Error().Err(err).Msg("Failed creating gRPC server")

		return
	}
}
