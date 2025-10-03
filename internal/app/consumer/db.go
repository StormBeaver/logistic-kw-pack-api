package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/StormBeaver/logistic-pack-api/internal/app/eventCounter"
	"github.com/StormBeaver/logistic-pack-api/internal/app/repo"
	"github.com/StormBeaver/logistic-pack-api/internal/model"
)

type Consumer interface {
	Start()
	Close()
}

type consumer struct {
	n      uint64
	events chan<- model.PackEvent

	repo repo.EventRepo

	batchSize uint64
	tick      time.Duration

	wg *sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.PackEvent) *consumer {

	wg := &sync.WaitGroup{}

	return &consumer{
		n:         n,
		batchSize: batchSize,
		tick:      consumeTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
	}
}

func (c *consumer) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.lockDelivery()
	c.mainDelivery()
}

func (c *consumer) Close() {
	c.cancel()
	c.wg.Wait()
}

func (c *consumer) lockDelivery() {
	for range c.n {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.tick)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.PreProcess(c.ctx, c.batchSize)
					if err != nil {
						continue
					}
					eventCounter.EventsCount.Add(float64(len(events)))

					for _, event := range events {
						c.events <- event
					}
					if len(events) == 0 {
						return
					}
				case <-c.ctx.Done():
					return
				}
			}
		}()
	}
	c.wg.Wait()
}

func (c *consumer) mainDelivery() {
	for range c.n {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.tick)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(c.ctx, c.batchSize)
					if err != nil {
						continue
					}
					eventCounter.EventsCount.Add(float64(len(events)))

					for _, event := range events {
						c.events <- event
					}
				case <-c.ctx.Done():
					return
				}
			}
		}()
	}
}
