package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags Health
// @Description Check the status of the server
// @Success 200
// @Router /status [get]
func GetStatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "up",
	})
}
