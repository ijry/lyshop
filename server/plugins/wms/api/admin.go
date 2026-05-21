package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/response"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	wmssvc "github.com/ijry/lyshop/plugins/wms/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/wms/warehouses", listWarehouses)
	g.POST("/wms/warehouses", createWarehouse)
	g.GET("/wms/stocks", listStocks)
	g.POST("/wms/inbound", doInbound)
	g.POST("/wms/outbound", doOutbound)
}

func listWarehouses(c *gin.Context) {
	list, err := wmssvc.ListWarehouses(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func createWarehouse(c *gin.Context) {
	var w wmsmodel.Warehouse
	if err := c.ShouldBindJSON(&w); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := wmssvc.CreateWarehouse(c.Request.Context(), &w); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, w)
}

func listStocks(c *gin.Context) {
	warehouseID, _ := strconv.ParseUint(c.Query("warehouse_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := wmssvc.ListStocks(c.Request.Context(), warehouseID, page, size)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func doInbound(c *gin.Context) {
	var req struct {
		WarehouseID uint64 `json:"warehouse_id"`
		SkuID       uint64 `json:"sku_id"`
		Qty         int    `json:"qty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := wmssvc.Inbound(c.Request.Context(), req.WarehouseID, req.SkuID, req.Qty, 0, "manual"); err != nil {
		response.Fail(c, 50001, err.Error())
		return
	}
	response.OK(c, nil)
}

func doOutbound(c *gin.Context) {
	var req struct {
		WarehouseID uint64 `json:"warehouse_id"`
		SkuID       uint64 `json:"sku_id"`
		Qty         int    `json:"qty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := wmssvc.Outbound(c.Request.Context(), req.WarehouseID, req.SkuID, req.Qty, 0, "manual"); err != nil {
		response.Fail(c, 50002, err.Error())
		return
	}
	response.OK(c, nil)
}
