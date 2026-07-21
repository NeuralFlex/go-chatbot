package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const userIDContextKey = "user_id"

func RequireUser(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}
	c.Set(userIDContextKey, userID)
	c.Next()
}


func UserID(c *gin.Context) string {
	return c.GetString(userIDContextKey)
}
