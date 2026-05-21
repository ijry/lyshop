package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type R struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// OK writes a successful response.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, R{Code: 0, Msg: "success", Data: data})
}

// Fail writes a business error response (HTTP 200, non-zero code).
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, R{Code: code, Msg: msg, Data: nil})
}

// Err returns a response struct (for use with AbortWithStatusJSON).
func Err(code int, msg string) R {
	return R{Code: code, Msg: msg, Data: nil}
}

// PageData wraps list + total for paginated responses.
type PageData struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}
