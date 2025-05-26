package sender

import (
	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

type EventSender interface {
	Send(pack *model.PackEvent) error
}
