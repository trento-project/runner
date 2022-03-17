package runner

import (
	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	c.JSON(200, map[string]string{"status": "ok"})
}
