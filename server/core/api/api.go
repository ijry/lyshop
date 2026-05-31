// Package api provides convenience helpers for Gin handler responses and middleware.
// It wraps core/response and core/middleware for simpler imports in plugin API files.
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
)

// RequirePermission is an alias for middleware.RequirePermission.
func RequirePermission(perm string) gin.HandlerFunc {
	return middleware.RequirePermission(perm)
}

// Success writes a successful JSON response.
func Success(c *gin.Context, data any) {
	response.OK(c, data)
}

// Error writes an error JSON response from an error value.
func Error(c *gin.Context, err error) {
	response.Fail(c, 1, err.Error())
}

// BadRequest writes a 400-style error response.
func BadRequest(c *gin.Context, msg string) {
	response.Fail(c, 400, msg)
}

// Unauthorized writes a 401-style error response.
func Unauthorized(c *gin.Context, msg string) {
	response.Fail(c, 401, msg)
}
