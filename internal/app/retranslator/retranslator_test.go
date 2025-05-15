package retranslator

import (
	"testing"
	"time"

	"route255/logistic-kw-pack-api/internal/mocks"

	"github.com/golang/mock/gomock"
)

func TestStart(t *testing.T) {

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	repo.EXPECT().PreProcess(gomock.Any()).AnyTimes()
	repo.EXPECT().Lock(gomock.Any()).AnyTimes()

	cfg := Config{
		ChannelSize:   512,
		ConsumerCount: 2,
		BatchSize:     10,
		Ticker:        10 * time.Second,
		ProducerCount: 2,
		WorkerCount:   2,
		Repo:          repo,
		Sender:        sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start()

	retranslator.Close()
}
