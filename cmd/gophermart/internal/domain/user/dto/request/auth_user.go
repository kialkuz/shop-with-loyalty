package request

type AuthUser struct {
	Login    string `json:"login" validate:"required,gte=4,lte=50"`
	Password string `json:"password" validate:"required,gte=5,lte=20"`
}
