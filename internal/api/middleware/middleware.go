package middleware

import (
	"ai-dev-platform/internal/requestid"
	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
	"time"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	xRequestID contextKey = "X-Request-ID"
)

func RequestIdRouter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestId := c.Request.Header.Get(xRequestID.String())

		// Create request id with UUID4
		if requestId == "" {
			u4 := uuid.New()
			requestId = u4.String()
		}

		// Expose it for use in the application
		c.Set(xRequestID.String(), requestId)

		// Set headers
		c.Writer.Header().Set(xRequestID.String(), requestId)

		// set to context with both request_id and trace_id
		ctx := requestid.NewContext(c.Request.Context(), requestId, time.Now())

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
