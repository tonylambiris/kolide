package requestlogger

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// New builds a new custom logrus format for
// gin route/request information.
func New(logger *logrus.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		status := c.Writer.Status()

		entry := logger.WithFields(logrus.Fields{
			"path":       path,
			"status":     status,
			"method":     c.Request.Method,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
			"time":       end.Format(timeFormat),
		})

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			entry.Error(c.Errors.String())
		} else {
			switch status {
			case 200:
				entry.Debug(fmt.Sprintf("[%d] %s", status, path))
			case 404:
				entry.Warn(fmt.Sprintf("[%d] %s", status, path))
			case 500:
			case 401:
				entry.Error(fmt.Sprintf("[%d] %s", status, path))
			default:
				entry.Debug(fmt.Sprintf("[%d] %s", status, path))
			}
		}
	}
}
