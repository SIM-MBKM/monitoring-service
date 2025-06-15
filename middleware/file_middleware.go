package middleware

import (
	"monitoring-service/dto"
	"monitoring-service/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func FileUploadSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set maximum request size
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, helper.MaxFileSize+helper.MaxContentLength)

		// Check Content-Type for multipart
		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "multipart/form-data") {
			c.JSON(http.StatusBadRequest, dto.Response{
				Status:  dto.STATUS_ERROR,
				Message: "Invalid content type for file upload",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
