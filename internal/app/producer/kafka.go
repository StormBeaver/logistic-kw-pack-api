package producer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/StormBeaver/logistic-pack-api/internal/app/eventCounter"
	"github.com/StormBeaver/logistic-pack-api/internal/app/repo"
	"github.com/StormBeaver/logistic-pack-api/internal/app/sender"
	"github.com/StormBeaver/logistic-pack-api/internal/model"

	"github.com/gammazero/workerpool"
)

type Producer interface {
	Start()
	Close()
}

type producer struct {
	n         uint64
	batchSize uint64
	tick      time.Duration

	sender sender.EventSender
	events <-chan model.PackEvent

	repo repo.EventRepo

	workerPool *workerpool.WorkerPool

	wg *sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

func NewKafkaProducer(
	n uint64,
	sender sender.EventSender,
	events <-chan model.PackEvent,
	workerPool *workerpool.WorkerPool,
	repo repo.EventRepo,
	ticker time.Duration,
	batchSize uint64,
) Producer {
	wg := &sync.WaitGroup{}

	return &producer{
		n:          n,
		sender:     sender,
		events:     events,
		workerPool: workerPool,
		repo:       repo,
		wg:         wg,
		tick:       ticker,
		batchSize:  batchSize,
	}
}

func (p *producer) Start() {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	for range p.n {
		p.wg.Add(1)
		go func() {

			toUnlock := make([]uint64, 0, int(p.batchSize))
			toRemove := make([]uint64, 0, int(p.batchSize))

			defer p.wg.Done()
			ticker := time.NewTicker(p.tick)
			for {
				select {
				case event := <-p.events:
					if err := p.sender.Send(&event); err != nil { //  TODO: compare by event.ID and event.Pack.Id for processing more than 1 type of event
						toUnlock = append(toUnlock, event.ID)
					} else {
						toRemove = append(toRemove, event.ID)
					}
				case <-ticker.C:
					p.delivery(toRemove, toUnlock)
					toUnlock = toUnlock[:0]
					toRemove = toRemove[:0]
				case <-p.ctx.Done():
					p.delivery(toRemove, toUnlock)
					return
				}
			}
		}()
	}
}

func (p *producer) Close() {
	p.cancel()
	p.wg.Wait()
}

func (p *producer) delivery(toRemove, toUnlock []uint64) {
	if len(toUnlock) != 0 {
		tUnlock := append(make([]uint64, 0, len(toUnlock)), toUnlock...)
		p.workerPool.Submit(func() {
			if err := p.repo.Unlock(p.ctx, tUnlock); err != nil {
				log.Println(err)
			}
			eventCounter.EventsCount.Sub(float64(len(tUnlock)))
		})
	}

	if len(toRemove) != 0 {
		tRemove := append(make([]uint64, 0, len(toRemove)), toRemove...)
		p.workerPool.Submit(func() {
			if err := p.repo.Remove(p.ctx, tRemove); err != nil {
				log.Printf("Remove error: %s", err)
				if err := p.repo.Unlock(p.ctx, tRemove); err != nil {
					log.Printf("Unlock error while handling Remove error: %s", err)
				}
			}
			eventCounter.EventsCount.Sub(float64(len(tRemove)))
		})
	}
}
