package responce

type ViewBalance struct {
	Current   int `json:"current"`
	WithDrawn int `json:"withdrawn"`
}
