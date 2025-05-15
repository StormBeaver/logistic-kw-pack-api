package main

import (
	"os"
	"os/signal"
	"syscall"

	"route255/logistic-kw-pack-api/internal/app/retranslator"
)

func main() {

	sigs := make(chan os.Signal, 1)

	cfg := retranslator.Config{
		ChannelSize:   512,
		ConsumerCount: 2,
		BatchSize:     10,
		ProducerCount: 28,
		WorkerCount:   2,
	}

	retranslator := retranslator.NewRetranslator(cfg)
	retranslator.Start()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	retranslator.Close()
}
