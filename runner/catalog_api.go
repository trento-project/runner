package runner

import (
	"github.com/gin-gonic/gin"
)

func CatalogHandler(runnerService RunnerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, runnerService.GetCatalog())
	}
}
