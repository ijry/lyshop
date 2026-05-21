package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/response"
	authsvc "github.com/ijry/lyshop/service/auth"
)

// RegisterAdminRoutes adds admin auth routes (no JWT required for login itself).
func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.POST("/auth/login", adminLogin)
}

func adminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	token, err := authsvc.AdminLogin(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}
	response.OK(c, gin.H{"token": token})
}
