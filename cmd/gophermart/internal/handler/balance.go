package handler

import (
	"kialkuz/gophermart/internal/service"

	"github.com/gin-gonic/gin"
)

type BalanceHandler struct {
}

func NewBalanceHandler(
	service *service.BalanceService,
) *BalanceHandler {
	return &BalanceHandler{}
}

func (h *BalanceHandler) GetBalance(c *gin.Context) {

}

func (h *BalanceHandler) WithDraw(c *gin.Context) {

}

func (h *BalanceHandler) WithDrawals(c *gin.Context) {

}

func (h *BalanceHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/user/balance", h.GetBalance)
	rg.POST("/user/balance/withdraw", h.WithDraw)
	rg.POST("/user/withdrawals", h.WithDrawals)
}
