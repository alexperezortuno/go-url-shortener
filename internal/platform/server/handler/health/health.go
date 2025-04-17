package health

import (
	"github.com/alexperezortuno/go-url-shortner/internal/platform/services/metrics"
	"github.com/gin-gonic/gin"
	"net/http"
)

var metricsService = metrics.GetInstance()

func CheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tags []metrics.MetricTag
		tags = append(tags, metrics.MetricTag{Key: "health_check", Value: "ok"})
		// Increment the health check counter
		metricsService.IncrementCounter("health_check_requests_total", tags)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "everything is ok",
		})
	}
}
