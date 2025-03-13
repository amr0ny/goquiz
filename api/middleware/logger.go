package middleware

import (
	"api/common"
	common2 "github.com/amr0ny/goquiz/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func LoggerMiddleware(c *gin.Context) {
	log, err := common2.GetLogger()
	if err != nil {
		common.ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	start := time.Now()

	method := c.Request.Method
	path := c.Request.URL.Path
	clientIP := c.ClientIP()

	log.WithFields(logrus.Fields{
		"method": method,
		"path":   path,
		"ip":     clientIP,
	}).Info("Request started")

	c.Next()
	dur := time.Since(start)

	statusCode := c.Writer.Status()
	log.WithFields(logrus.Fields{
		"method":     method,
		"path":       path,
		"ip":         clientIP,
		"statusCode": statusCode,
		"duration":   dur.Seconds(),
	}).Info("Request completed")
}
