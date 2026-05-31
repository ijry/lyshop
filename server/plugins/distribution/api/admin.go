package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/api"
	"github.com/ijry/lyshop/server/plugins/distribution/model"
	"github.com/ijry/lyshop/server/plugins/distribution/service"
)

// RegisterAdminRoutes 注册管理端路由
func RegisterAdminRoutes(r *gin.RouterGroup) {
	g := r.Group("/admin/api/distribution")
	{
		// 配置管理
		g.GET("/config", api.RequirePermission("distribution:view"), getConfig)
		g.PUT("/config", api.RequirePermission("distribution:edit"), updateConfig)

		// 分销商管理
		g.GET("/distributors", api.RequirePermission("distribution:view"), listDistributors)
		g.GET("/distributors/:id", api.RequirePermission("distribution:view"), getDistributor)
		g.PUT("/distributors/:id", api.RequirePermission("distribution:edit"), updateDistributor)

		// 分销订单
		g.GET("/orders", api.RequirePermission("distribution:view"), listOrders)
		g.GET("/orders/:id", api.RequirePermission("distribution:view"), getOrder)

		// 提现管理
		g.GET("/withdrawals", api.RequirePermission("distribution:view"), listWithdrawals)
		g.GET("/withdrawals/:id", api.RequirePermission("distribution:view"), getWithdrawal)
		g.POST("/withdrawals/:id/approve", api.RequirePermission("distribution:edit"), approveWithdrawal)
		g.POST("/withdrawals/:id/reject", api.RequirePermission("distribution:edit"), rejectWithdrawal)
		g.POST("/withdrawals/:id/complete", api.RequirePermission("distribution:edit"), completeWithdrawal)
	}
}

func getConfig(c *gin.Context) {
	config, err := service.GetConfig(c.Request.Context())
	if err != nil {
		api.Error(c, err)
		return
	}
	api.Success(c, config)
}

func updateConfig(c *gin.Context) {
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	if err := service.UpdateConfig(c.Request.Context(), updates); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, nil)
}

func listDistributors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := c.Query("status")

	distributors, total, err := service.ListDistributors(c.Request.Context(), page, size, status)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, gin.H{
		"list":  distributors,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

func getDistributor(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	distributor, err := service.GetDistributorByID(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, distributor)
}

func updateDistributor(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	if err := service.UpdateDistributor(c.Request.Context(), id, updates); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, nil)
}

func listOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := c.Query("status")
	distributorID, _ := strconv.ParseUint(c.Query("distributor_id"), 10, 64)

	orders, total, err := service.ListDistributionOrders(c.Request.Context(), page, size, status, distributorID)
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

func getOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	order, err := service.GetDistributionOrder(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, order)
}

func listWithdrawals(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := c.Query("status")
	distributorID, _ := strconv.ParseUint(c.Query("distributor_id"), 10, 64)

	withdrawals, total, err := service.ListWithdrawals(c.Request.Context(), page, size, status, distributorID)
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

func getWithdrawal(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	withdrawal, err := service.GetWithdrawal(c.Request.Context(), id)
	if err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, withdrawal)
}

func approveWithdrawal(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := service.ApproveWithdrawal(c.Request.Context(), id); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, nil)
}

func rejectWithdrawal(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	if err := service.RejectWithdrawal(c.Request.Context(), id, req.Reason); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, nil)
}

func completeWithdrawal(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := service.CompleteWithdrawal(c.Request.Context(), id); err != nil {
		api.Error(c, err)
		return
	}

	api.Success(c, nil)
}
