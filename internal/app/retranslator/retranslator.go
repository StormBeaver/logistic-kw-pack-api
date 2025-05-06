package retranslator

import (
	"time"

	"route255/L-KW-P-API/internal/model"

	"route255/L-KW-P-API/internal/app/consumer"
	"route255/L-KW-P-API/internal/app/producer"
	"route255/L-KW-P-API/internal/app/repo"
	"route255/L-KW-P-API/internal/app/sender"

	"github.com/gammazero/workerpool"
)

type Retranslator interface {
	Start()
	Close()
}

type Config struct {
	ChannelSize uint64

	ConsumerCount uint64

	BatchSize uint64
	Ticker    time.Duration

	ProducerCount uint64
	WorkerCount   int

	Repo   repo.EventRepo
	Sender sender.EventSender
}

type retranslator struct {
	events     chan model.PackEvent
	consumer   consumer.Consumer
	producer   producer.Producer
	workerPool *workerpool.WorkerPool
}

func NewRetranslator(cfg Config) Retranslator {
	events := make(chan model.PackEvent, cfg.ChannelSize)
	workerPool := workerpool.New(cfg.WorkerCount)

	consumer := consumer.NewDbConsumer(
		cfg.ConsumerCount,
		cfg.BatchSize,
		cfg.Ticker,
		cfg.Repo,
		events)

	producer := producer.NewKafkaProducer(
		cfg.ProducerCount,
		cfg.Sender,
		events,
		workerPool,
		cfg.Repo,
		cfg.Ticker,
		cfg.BatchSize)

	return &retranslator{
		events:     events,
		consumer:   consumer,
		producer:   producer,
		workerPool: workerPool,
	}
}

func (r *retranslator) Start() {
	r.producer.Start()
	r.consumer.Start()
}

func (r *retranslator) Close() {
	r.consumer.Close()
	r.producer.Close()
	r.workerPool.StopWait()
}
