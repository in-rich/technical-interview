package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

func Logger(logger zerolog.Logger, projectID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		status := c.Writer.Status()
		errs := c.Errors.Errors()

		logLevel := zerolog.TraceLevel
		severity := "INFO" // For GCP.
		if status > 499 {
			logLevel = zerolog.ErrorLevel
			severity = "ERROR"
		} else if status > 399 || len(errs) > 0 {
			logLevel = zerolog.WarnLevel
			severity = "WARNING"
		}

		parserQuery := zerolog.Dict()
		for k, v := range c.Request.URL.Query() {
			parserQuery.Strs(k, v)
		}

		// Allow logs to be grouped in log explorer.
		// https://cloud.google.com/run/docs/logging#run_manual_logging-go
		var trace string
		if projectID != "" {
			traceHeader := c.GetHeader("X-Cloud-Trace-Context")
			traceParts := strings.Split(traceHeader, "/")
			if len(traceParts) > 0 && len(traceParts[0]) > 0 {
				trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
			}
		}

		ll := logger.WithLevel(logLevel).
			Dict(
				"httpRequest", zerolog.Dict().
					Str("requestMethod", c.Request.Method).
					Str("requestUrl", c.FullPath()).
					Int("status", status).
					Str("userAgent", c.Request.UserAgent()).
					Str("remoteIp", c.ClientIP()).
					Str("protocol", c.Request.Proto).
					Str("latency", end.Sub(start).String()),
			).
			Time("start", start).
			Dur("postProcessingLatency", time.Now().Sub(end)).
			Int64("contentLength", c.Request.ContentLength).
			Str("ip", c.ClientIP()).
			Str("contentType", c.ContentType()).
			Str("auth", c.GetHeader("Authorization")).
			Strs("errors", errs).
			Str("severity", severity)

		if len(trace) > 0 {
			ll = ll.Str("logging.googleapis.com/trace", trace)
		}

		ll.Msg(c.Request.URL.String())
	}
}
