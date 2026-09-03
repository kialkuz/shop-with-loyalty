package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CurrentUserID(c *gin.Context) uuid.UUID {
	value, _ := c.Get("user_id")
	userID, _ := value.(uuid.UUID)

	return userID
}
