package model

import (
	orderModel "kialkuz/shop-with-loyalty/internal/modules/order/model"
	"time"

	"github.com/google/uuid"
)

type Drawal struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	OrderNumber orderModel.Number
	Sum         int
	ProcessedAt time.Time
}
