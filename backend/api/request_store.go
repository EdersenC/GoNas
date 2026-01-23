package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"goevents"
)

// RequestStore attaches a goevents.Store to each request context.
func RequestStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		store := goevents.NewStore()
		requestID := uuid.NewString()
		store.Set("request.id", requestID)
		store.Set("request.method", c.Request.Method)
		store.Set("request.path", c.Request.URL.Path)
		store.Set("client.ip", c.ClientIP())
		store.Set("request.start", time.Now())

		c.Writer.Header().Set("X-Request-Id", requestID)

		ctx := goevents.WithStore(c.Request.Context(), store)
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		if start, ok := store.Get("request.start"); ok {
			if startedAt, ok := start.(time.Time); ok {
				duration := time.Since(startedAt)
				log.Printf("request_id=%s method=%s path=%s status=%d duration=%s",
					requestID,
					c.Request.Method,
					c.Request.URL.Path,
					c.Writer.Status(),
					duration,
				)
			}
		}
	}
}

// StoreFromContext returns the per-request store when available.
func StoreFromContext(c *gin.Context) (*goevents.Store, bool) {
	return goevents.FromContext(c.Request.Context())
}
