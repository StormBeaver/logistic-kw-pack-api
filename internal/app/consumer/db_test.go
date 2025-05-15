package consumer

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"route255/logistic-kw-pack-api/internal/app/repo"
	"route255/logistic-kw-pack-api/internal/mocks"
	"route255/logistic-kw-pack-api/internal/model"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

type Config struct {
	n         uint64
	events    chan<- model.PackEvent
	repo      repo.EventRepo
	batchSize uint64
	timeout   time.Duration
}

func TestMainDelivery(t *testing.T) {
	var (
		Id       uint64
		mu       sync.Mutex
		events   = make(chan model.PackEvent, 512)
		ctrl     = gomock.NewController(t)
		repo     = mocks.NewMockEventRepo(ctrl)
		err      = errors.New("planned Error")
		countErr int
	)

	defer ctrl.Finish()

	cfg := Config{
		n:         2,
		events:    events,
		repo:      repo,
		batchSize: 5,
		timeout:   2 * time.Second,
	}

	firstCase := repo.EXPECT().Lock(cfg.batchSize).Times(int(cfg.n) * 2).DoAndReturn(func(n uint64) ([]model.PackEvent, error) {
		output := make([]model.PackEvent, 2)
		for i := range output {
			output[i] = model.PackEvent{ID: Id}
			mu.Lock()
			Id++
			mu.Unlock()
		}
		return output, nil
	})

	secondCase := repo.EXPECT().Lock(cfg.batchSize).Times(int(cfg.n) * 2).DoAndReturn(func(n uint64) ([]model.PackEvent, error) {
		mu.Lock()
		countErr++
		mu.Unlock()
		fmt.Println("error #:", countErr)
		return []model.PackEvent{}, err
	}).After(firstCase)

	repo.EXPECT().Lock(cfg.batchSize).Times(int(cfg.n) * 2).DoAndReturn(func(n uint64) ([]model.PackEvent, error) {
		output := make([]model.PackEvent, rand.N(cfg.batchSize))
		for i := range output {
			output[i] = model.PackEvent{ID: Id}
			mu.Lock()
			Id++
			mu.Unlock()
		}
		return output, nil
	}).After(secondCase)

	c := NewDbConsumer(cfg.n,
		cfg.batchSize,
		cfg.timeout,
		cfg.repo,
		cfg.events)

	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.mainDelivery()

	go func() {
		for v := range events {
			fmt.Println(v.ID)
		}
	}()
	time.Sleep(13 * time.Second)
	c.Close()
}

func TestMainDeliveryV2(t *testing.T) {
	var (
		Id, idCounter uint64
		mu            sync.Mutex
		events        = make(chan model.PackEvent, 512)
		ctrl          = gomock.NewController(t)
		repo          = mocks.NewMockEventRepo(ctrl)
	)

	defer ctrl.Finish()

	cfg := Config{
		n:         1,
		events:    events,
		repo:      repo,
		batchSize: 5,
		timeout:   2 * time.Second,
	}

	repo.EXPECT().Lock(cfg.batchSize).Times(int(cfg.n) * 2).DoAndReturn(func(n uint64) ([]model.PackEvent, error) {
		output := make([]model.PackEvent, 2)
		for i := range output {
			output[i] = model.PackEvent{ID: Id}
			mu.Lock()
			Id++
			mu.Unlock()
		}
		return output, nil
	})

	c := NewDbConsumer(cfg.n,
		cfg.batchSize,
		cfg.timeout,
		cfg.repo,
		cfg.events)

	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.mainDelivery()
	go func() {
		for v := range events {
			if v.ID != uint64(idCounter) {
				t.Fail()
			}
			idCounter++
		}
	}()
	time.Sleep(5 * time.Second)
	c.Close()
}

func TestLockDelivery(t *testing.T) {
	var (
		events = make(chan model.PackEvent, 512)
		ctrl   = gomock.NewController(t)
		repo   = mocks.NewMockEventRepo(ctrl)
	)

	defer ctrl.Finish()

	cfg := Config{
		n:         2,
		events:    events,
		repo:      repo,
		batchSize: 5,
		timeout:   2 * time.Second,
	}

	c := NewDbConsumer(cfg.n,
		cfg.batchSize,
		cfg.timeout,
		cfg.repo,
		cfg.events)

	c.ctx, c.cancel = context.WithCancel(context.Background())

	repo.EXPECT().PreProcess(cfg.batchSize).Times(int(cfg.n)).Return([]model.PackEvent{}, nil)

	go c.lockDelivery()
	time.Sleep(3 * time.Second)

}
