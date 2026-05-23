package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	checkinsvc "github.com/ijry/lyshop/plugins/checkin/service"
)

func RegisterFrontRoutes(g *gin.RouterGroup) {
	auth := g.Group("").Use(middleware.RequireAuth)
	auth.POST("/checkin", doCheckin)
	auth.GET("/checkin/status", getStatus)
	auth.GET("/checkin/rules", getRules)
}

func doCheckin(c *gin.Context) {
	userID, _ := c.Get("user_id")
	points, consecutive, err := checkinsvc.Checkin(c.Request.Context(), userID.(uint64))
	if err != nil {
		response.Fail(c, 30001, err.Error())
		return
	}
	response.OK(c, gin.H{"points": points, "consecutive_days": consecutive})
}

func getStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	status, err := checkinsvc.GetStatus(c.Request.Context(), userID.(uint64))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, status)
}

func getRules(c *gin.Context) {
	rules, err := checkinsvc.GetRules(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, rules)
}
