package sender

import (
	"route255/L-KW-P-API/internal/model"
)

type EventSender interface {
	Send(pack *model.PackEvent) error
}
