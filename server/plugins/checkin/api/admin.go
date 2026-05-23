package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	checkinmodel "github.com/ijry/lyshop/plugins/checkin/model"
	checkinsvc "github.com/ijry/lyshop/plugins/checkin/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/checkin/rules", middleware.RequirePermission("checkin:config"), adminGetRules)
	g.PUT("/checkin/rules", middleware.RequirePermission("checkin:config"), adminSaveRules)
	g.GET("/checkin/logs", middleware.RequirePermission("checkin:view"), adminLogs)
}

func adminGetRules(c *gin.Context) {
	rules, err := checkinsvc.GetRules(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, rules)
}

func adminSaveRules(c *gin.Context) {
	var rules []checkinmodel.CheckinRule
	if err := c.ShouldBindJSON(&rules); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := checkinsvc.SaveRules(c.Request.Context(), rules); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func adminLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := checkinsvc.AdminLogs(c.Request.Context(), page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}
