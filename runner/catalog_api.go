package runner

import (
	"github.com/gin-gonic/gin"
)

func CatalogHandler(runnerService RunnerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if runnerService.IsCatalogReady() {
			c.JSON(200, runnerService.GetCatalog())
		} else {
			c.JSON(204, nil)
		}
	}
}
