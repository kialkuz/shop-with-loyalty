package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CurrentUserId(c *gin.Context) uuid.UUID {
	value, _ := c.Get("user_id")
	userId, _ := value.(uuid.UUID)

	return userId
}
