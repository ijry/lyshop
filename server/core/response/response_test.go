package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	OK(c, map[string]string{"key": "value"})

	var r R
	json.Unmarshal(w.Body.Bytes(), &r)
	assert.Equal(t, 0, r.Code)
	assert.Equal(t, "success", r.Msg)
}

func TestFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Fail(c, 401, "未登录")

	assert.Equal(t, http.StatusOK, w.Code)
	var r R
	json.Unmarshal(w.Body.Bytes(), &r)
	assert.Equal(t, 401, r.Code)
}
