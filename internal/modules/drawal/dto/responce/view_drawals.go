package responce

import "time"

type ViewDrawals struct {
	Number      string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
