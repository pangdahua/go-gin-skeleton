// Package gin is a helper package to get a gin compatible middleware.
package gin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/slok/go-http-metrics/middleware"
)

// Handler returns a Gin measuring middleware.
func Handler(paths []string, m middleware.Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if !contains(path, paths) {
			return
		}
		r := &reporter{c: c}
		m.Measure(path, r, func() {
			c.Next()
		})
	}
}

func contains(path string, paths []string) bool {
	for _, p := range paths {
		if path == p {
			return true
		}
	}
	return false
}

type reporter struct {
	c *gin.Context
}

func (r *reporter) Method() string { return r.c.Request.Method }

func (r *reporter) Context() context.Context { return r.c.Request.Context() }

func (r *reporter) URLPath() string { return r.c.Request.URL.Path }

func (r *reporter) StatusCode() int { return r.c.Writer.Status() }

func (r *reporter) BytesWritten() int64 { return int64(r.c.Writer.Size()) }