package runner

import (
	"github.com/gin-gonic/gin"
)

func setupTestDependencies() Dependencies {
	return Dependencies{
		webEngine: gin.Default(),
	}
}
