package runner

import (
	"github.com/gin-gonic/gin"
)

func ExecutionHandler(runnerService RunnerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r *ExecutionEvent

		if err := c.BindJSON(&r); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(500, gin.H{"status": "nok", "message": err.Error()})
			return
		}

		if err := runnerService.ScheduleExecution(r); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(500, gin.H{"status": "nok", "message": err.Error()})
			return
		}

		c.JSON(202, map[string]string{"status": "ok"})
	}
}
