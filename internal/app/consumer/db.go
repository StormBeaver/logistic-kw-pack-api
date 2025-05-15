package consumer

import (
	"context"
	"sync"
	"time"

	"route255/logistic-kw-pack-api/internal/app/repo"
	"route255/logistic-kw-pack-api/internal/model"
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

func (c *consumer) mainDelivery() {
	for range c.n {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.tick)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(c.batchSize)
					if err != nil {
						continue
					}
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

func (c *consumer) lockDelivery() {
	for range c.n {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.tick)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.PreProcess(c.batchSize)
					if err != nil {
						continue
					}
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
