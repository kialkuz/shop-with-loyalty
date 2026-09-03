package model

import (
	orderModel "kialkuz/shop-with-loyalty/internal/domain/order/model"
	"time"

	"github.com/google/uuid"
)

type Drawal struct {
	ID          uuid.UUID
	UserId      uuid.UUID
	OrderNumber orderModel.Number
	Sum         int
	ProcessedAt time.Time
}
