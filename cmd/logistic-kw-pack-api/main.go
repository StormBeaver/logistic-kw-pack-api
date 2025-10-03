package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/StormBeaver/logistic-pack-api/internal/app/repo"
	"github.com/StormBeaver/logistic-pack-api/internal/app/retranslator"
	"github.com/StormBeaver/logistic-pack-api/internal/app/sender"
	"github.com/StormBeaver/logistic-pack-api/internal/config"
	"github.com/StormBeaver/logistic-pack-api/internal/database"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func main() {

	sigs := make(chan os.Signal, 1)

	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}

	cfg := config.GetConfigInstance()

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
		log.Fatal().Err(err).Msg("Failed init postgres")
	}
	defer db.Close()

	cfgR := retranslator.Config{
		ChannelSize:   cfg.Kafka.Capacity,
		ConsumerCount: cfg.Retranslator.ConsumerCount,
		BatchSize:     cfg.Retranslator.BatchSize,
		ProducerCount: cfg.Retranslator.ProducerCount,
		WorkerCount:   cfg.Retranslator.WorkerCount,
		Repo:          repo.NewEventRepo(db),
		Sender:        sender.NewEventSender(cfg.Kafka.Brokers),
	}

	retranslator := retranslator.NewRetranslator(cfgR)
	retranslator.Start()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	retranslator.Close()
}
