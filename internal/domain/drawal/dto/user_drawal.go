package drawal

import "time"

type UserDrawal struct {
	OrderNumber string
	Sum         int
	ProcessedAt time.Time
}
