package handlers

import "github.com/gin-gonic/gin"

// PingHandler godoc
// @Summary      Ping the server
// @Description  Returns pong to confirm the server is alive
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
