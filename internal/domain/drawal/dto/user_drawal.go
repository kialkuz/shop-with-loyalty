package drawal

import (
	modelOrder "kialkuz/shop-with-loyalty/internal/domain/order/model"
	"time"
)

type UserDrawal struct {
	OrderNumber modelOrder.Number
	Sum         int
	ProcessedAt time.Time
}
