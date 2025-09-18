package sender

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

type EventSender interface {
	Send(pack *model.PackEvent) error
}

type Sender struct {
	sarama.SyncProducer
}

func (s Sender) Send(pack *model.PackEvent) error {
	msg, err := json.Marshal(*pack)
	if err != nil {
		return fmt.Errorf("marshaling pack: %w", err)
	}

	_, _, err = s.SendMessage(prepareMessage(pack.Type, msg))
	if err != nil {
		return fmt.Errorf("send message to Kafka: %w", err)
	}
	return nil
}

func NewEventSender(brokers []string) EventSender {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}

	return Sender{producer}
}

func prepareMessage(topic string, message []byte) *sarama.ProducerMessage { //maybe to config too
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.ByteEncoder(message),
	}
	return msg
}
