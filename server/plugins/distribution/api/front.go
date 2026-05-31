package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/api"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/plugins/distribution/model"
	"github.com/ijry/lyshop/plugins/distribution/service"
)

// RegisterFrontRoutes 注册前端路由
func RegisterFrontRoutes(r *gin.RouterGroup) {
	g := r.Group("/api/v1/distribution").Use(middleware.RequireAuth)
	{
		// 分销商信息
		g.GET("/info", getDistributorInfo)
		g.POST("/apply", applyDistributor)

		// 我的团队
		g.GET("/team", getTeam)

		// 佣金订单
		g.GET("/orders", listMyOrders)

		// 提现
		g.POST("/withdrawals", createWithdrawal)
		g.GET("/withdrawals", listMyWithdrawals)
		g.GET("/withdrawals/:id", getMyWithdrawal)
	}
}

func getDistributorInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		api.Unauthorized(c, "未登录")
		return
	}

	distributor, err := service.GetDistributor(c.Request.Context(), userID.(uint64))
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, distributor)
}

func applyDistributor(c *gin.Context) {
	var req struct {
		ParentID uint64 `json:"parent_id"`
		RealName string `json:"real_name"`
		IDCard   string `json:"id_card"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		api.Unauthorized(c, "未登录")
		return
	}

	distributor := &model.Distributor{
		UserID:   userID.(uint64),
		ParentID: req.ParentID,
		Status:   "pending",
		RealName: req.RealName,
		IDCard:   req.IDCard,
		Phone:    req.Phone,
	}

	if err := service.CreateDistributor(c.Request.Context(), distributor); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, distributor)
}

func getTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		api.Unauthorized(c, "未登录")
		return
	}

	distributor, err := service.GetDistributor(c.Request.Context(), userID.(uint64))
	if err != nil {
		api.Error(c, err)
		return
	}

	// 获取下级分销商
	var team []model.Distributor
	// TODO: 实现获取下级分销商逻辑

	api.Success(c, gin.H{
		"distributor": distributor,
		"team":        team,
	})
}

func listMyOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := c.Query("status")

	userID, exists := c.Get("user_id")
	if !exists {
		api.Unauthorized(c, "未登录")
		return
	}

	distributor, err := service.GetDistributor(c.Request.Context(), userID.(uint64))
	if err != nil {
		api.Error(c, err)
		return
	}

	orders, total, err := service.ListDistributionOrders(c.Request.Context(), page, size, status, distributor.ID)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

func createWithdrawal(c *gin.Context) {
	var req struct {
		Amount      float64 `json:"amount" binding:"required"`
		BankName    string  `json:"bank_name" binding:"required"`
		BankAccount string  `json:"bank_account" binding:"required"`
		AccountName string  `json:"account_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		api.Unauthorized(c, "未登录")
		return
	}

	distributor, err := service.GetDistributor(c.Request.Context(), userID.(uint64))
	if err != nil {
		api.Error(c, err)
		return
	}

	withdrawal := &model.DistributionWithdrawal{
		DistributorID: distributor.ID,
		Amount:        req.Amount,
		BankName:      req.BankName,
		BankAccount:   req.BankAccount,
		AccountName:   req.AccountName,
	}

	if err := service.CreateWithdrawal(c.Request.Context(), withdrawal); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, withdrawal)
}

func listMyWithdrawals(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := c.Query("status")

	userID, exists := c.Get("user_id")
	if !exists {
		api.Unauthorized(c, "未登录")
		return
	}

	distributor, err := service.GetDistributor(c.Request.Context(), userID.(uint64))
	if err != nil {
		api.Error(c, err)
		return
	}

	withdrawals, total, err := service.ListWithdrawals(c.Request.Context(), page, size, status, distributor.ID)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, gin.H{
		"list":  withdrawals,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

func getMyWithdrawal(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	withdrawal, err := service.GetWithdrawal(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, withdrawal)
}
