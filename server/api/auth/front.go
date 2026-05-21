package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/response"
	authsvc "github.com/ijry/lyshop/service/auth"
)

// RegisterFrontRoutes adds user auth routes to the front-end router group.
func RegisterFrontRoutes(g *gin.RouterGroup) {
	g.POST("/auth/sms/send", sendSMSCode)
	g.POST("/auth/sms/login", smsLogin)
}

func sendSMSCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "手机号不能为空")
		return
	}
	code, err := authsvc.SendSMSCode(c.Request.Context(), req.Phone)
	if err != nil {
		response.Fail(c, 500, "发送失败: "+err.Error())
		return
	}
	// dev_code is included for development/testing; remove or gate behind debug mode in production
	response.OK(c, gin.H{"dev_code": code})
}

func smsLogin(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	token, err := authsvc.SMSLogin(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		response.Fail(c, 10001, err.Error())
		return
	}
	response.OK(c, gin.H{"token": token})
}
