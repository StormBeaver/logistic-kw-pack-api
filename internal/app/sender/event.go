package sender

import (
	"route255/logistic-kw-pack-api/internal/model"
)

type EventSender interface {
	Send(pack *model.PackEvent) error
}
