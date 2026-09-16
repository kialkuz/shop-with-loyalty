package scheduler

type GetAccrual struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
	Err     error
}
