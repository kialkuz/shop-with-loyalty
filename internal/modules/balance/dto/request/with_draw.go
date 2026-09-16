package request

type WithDraw struct {
	Order string  `json:"order" validate:"required"`
	Sum   float64 `json:"sum" validate:"required"`
}
