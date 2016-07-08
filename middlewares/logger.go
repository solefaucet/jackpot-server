package middlewares

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/solefaucet/jackpot-server/models"
)

// Logger returns a middleware that logs all request
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// mark start time
		start := time.Now()

		// process request
		c.Next()

		// log after processing
		status := c.Writer.Status()
		ip := c.ClientIP()

		fields := logrus.Fields{
			"event":            models.LogEventHTTPRequest,
			"method":           c.Request.Method,
			"path":             c.Request.URL.Path,
			"query":            c.Request.URL.Query().Encode(),
			"http_status_code": status,
			"response_time":    float64(time.Since(start).Nanoseconds()) / 1e6,
			"ip":               ip,
		}
		if err := c.Errors.ByType(gin.ErrorTypeAny).Last(); err != nil {
			fields["error"] = err.Error()
		}

		entry := logrus.WithFields(fields)
		if status == http.StatusInternalServerError {
			entry.Error("internal server error")
			return
		}

		entry.Info("succeed to process http request")
	}
}
