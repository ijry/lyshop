package api

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/config"
	"github.com/ijry/lyshop/core/db"
	inventorycore "github.com/ijry/lyshop/core/inventory"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	"gorm.io/gorm"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/inventory/tasks", middleware.RequirePermission("system:config"), listTasks)
	g.POST("/inventory/tasks/:id/retry", middleware.RequirePermission("system:config"), retryTask)
	g.POST("/external-wms/callback", callback)
}

func listTasks(c *gin.Context) {
	var rows []inventorycore.InventoryIntegrationTask
	if err := db.DB.WithContext(c.Request.Context()).Order("id desc").Limit(50).Find(&rows).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, rows)
}

func retryTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	now := time.Now()
	if err := db.DB.WithContext(c.Request.Context()).
		Model(&inventorycore.InventoryIntegrationTask{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":          inventorycore.TaskStatusPending,
			"last_error":      "",
			"lock_owner":      "",
			"lock_expires_at": nil,
			"next_retry_at":   &now,
		}).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, gin.H{"retried": id})
}

type callbackRequest struct {
	AppKey     string `json:"app_key"`
	Timestamp  string `json:"timestamp"`
	Nonce      string `json:"nonce"`
	Sign       string `json:"sign"`
	Body       string `json:"body"`
	RequestID  string `json:"request_id"`
	CallbackID string `json:"callback_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

func callback(c *gin.Context) {
	var req callbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if config.Global.ExternalWMS.AppSecret != "" {
		if err := inventorycore.VerifyCallbackSignature(inventorycore.CallbackEnvelopeFromRequest(req.AppKey, req.Timestamp, req.Nonce, req.Sign, req.Body), config.Global.ExternalWMS.AppSecret, time.Now()); err != nil {
			response.Fail(c, 400, "签名校验失败")
			return
		}
	}
	if strings.TrimSpace(req.Body) != "" {
		payload, err := inventorycore.NewExternalAdapter().ParseCallbackBody(req.Body)
		if err != nil {
			response.Fail(c, 400, "回调体错误")
			return
		}
		req.RequestID = payload.RequestID
		req.CallbackID = payload.CallbackID
		req.Status = payload.Status
		req.Message = payload.Message
	}

	err := db.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		return inventorycore.CompleteTaskByCallback(tx, req.RequestID, req.CallbackID, req.Status, req.Message, time.Now())
	})
	if err != nil {
		response.Fail(c, 500, "回调处理失败")
		return
	}

	response.OK(c, gin.H{"received": true})
}
