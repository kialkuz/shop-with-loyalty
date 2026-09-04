package responce

type ViewBalance struct {
	Current   float64 `json:"current"`
	WithDrawn float64 `json:"withdrawn"`
}
