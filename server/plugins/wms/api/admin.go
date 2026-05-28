package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/response"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	wmssvc "github.com/ijry/lyshop/plugins/wms/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
	g.GET("/wms/warehouses", middleware.RequirePermission("wms:view"), listWarehouses)
	g.POST("/wms/warehouses", middleware.RequirePermission("wms:edit"), createWarehouse)
	g.PUT("/wms/warehouses/:id", middleware.RequirePermission("wms:edit"), updateWarehouse)
	g.PUT("/wms/warehouses/:id/status", middleware.RequirePermission("wms:edit"), updateWarehouseStatus)

	g.GET("/wms/stocks", middleware.RequirePermission("wms:view"), listStocks)
	g.PUT("/wms/stocks/:id/safety", middleware.RequirePermission("wms:edit"), updateStockSafety)

	g.GET("/wms/docs", middleware.RequirePermission("wms:view"), listDocs)
	g.POST("/wms/docs", middleware.RequirePermission("wms:edit"), createDoc)
	g.GET("/wms/docs/:id", middleware.RequirePermission("wms:view"), getDocDetail)
	g.PUT("/wms/docs/:id", middleware.RequirePermission("wms:edit"), updateDoc)
	g.POST("/wms/docs/:id/complete", middleware.RequirePermission("wms:edit"), completeDoc)
	g.POST("/wms/docs/:id/cancel", middleware.RequirePermission("wms:edit"), cancelDoc)

	g.GET("/wms/movements", middleware.RequirePermission("wms:view"), listMovements)
}

func listWarehouses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	statusPtr, err := parseOptionalInt8(c.Query("status"))
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	list, total, err := wmssvc.ListWarehouses(c.Request.Context(), wmssvc.WarehouseListQuery{
		Page:    page,
		Size:    size,
		Keyword: c.Query("keyword"),
		Status:  statusPtr,
	})
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
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

func updateWarehouse(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	var req wmsmodel.Warehouse
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := wmssvc.UpdateWarehouse(c.Request.Context(), id, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func updateWarehouseStatus(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := wmssvc.UpdateWarehouseStatus(c.Request.Context(), id, req.Status); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func listStocks(c *gin.Context) {
	warehouseID, _ := strconv.ParseUint(c.Query("warehouse_id"), 10, 64)
	skuID, _ := strconv.ParseUint(c.Query("sku_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	warningOnly := parseBool(c.Query("warning_only"))
	list, total, err := wmssvc.ListStocks(c.Request.Context(), wmssvc.StockListQuery{
		Page:        page,
		Size:        size,
		WarehouseID: warehouseID,
		SkuID:       skuID,
		WarningOnly: warningOnly,
	})
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func updateStockSafety(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	var req struct {
		SafeQty int `json:"safe_qty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := wmssvc.UpdateStockSafety(c.Request.Context(), id, req.SafeQty); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func listDocs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	warehouseID, _ := strconv.ParseUint(c.Query("warehouse_id"), 10, 64)
	startAt, err := parseOptionalTime(c.Query("start_at"))
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	endAt, err := parseOptionalTime(c.Query("end_at"))
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	list, total, err := wmssvc.ListDocs(c.Request.Context(), wmssvc.DocListQuery{
		Page:        page,
		Size:        size,
		WarehouseID: warehouseID,
		DocType:     c.Query("doc_type"),
		Status:      c.Query("status"),
		DocNo:       c.Query("doc_no"),
		StartAt:     startAt,
		EndAt:       endAt,
	})
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func createDoc(c *gin.Context) {
	var req wmssvc.CreateDocInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	doc, err := wmssvc.CreateDraftDoc(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, doc)
}

func getDocDetail(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	detail, err := wmssvc.GetDocDetail(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, detail)
}

func updateDoc(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	var req wmssvc.UpdateDocInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := wmssvc.UpdateDraftDoc(c.Request.Context(), id, req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func completeDoc(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := wmssvc.CompleteDraftDoc(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func cancelDoc(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := wmssvc.CancelDraftDoc(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func listMovements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	warehouseID, _ := strconv.ParseUint(c.Query("warehouse_id"), 10, 64)
	skuID, _ := strconv.ParseUint(c.Query("sku_id"), 10, 64)
	startAt, err := parseOptionalTime(c.Query("start_at"))
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	endAt, err := parseOptionalTime(c.Query("end_at"))
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	list, total, err := wmssvc.ListMovements(c.Request.Context(), wmssvc.MovementListQuery{
		Page:        page,
		Size:        size,
		WarehouseID: warehouseID,
		SkuID:       skuID,
		BizType:     c.Query("biz_type"),
		DocNo:       c.Query("doc_no"),
		StartAt:     startAt,
		EndAt:       endAt,
	})
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func parseIDParam(c *gin.Context, key string) (uint64, error) {
	v := strings.TrimSpace(c.Param(key))
	id, err := strconv.ParseUint(v, 10, 64)
	if err != nil || id == 0 {
		return 0, fmt.Errorf("参数 %s 非法", key)
	}
	return id, nil
}

func parseOptionalInt8(raw string) (*int8, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.ParseInt(raw, 10, 8)
	if err != nil {
		return nil, fmt.Errorf("status 参数非法")
	}
	val := int8(v)
	return &val, nil
}

func parseOptionalTime(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, raw)
		if err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("时间参数格式非法")
}

func parseBool(raw string) bool {
	raw = strings.TrimSpace(strings.ToLower(raw))
	return raw == "1" || raw == "true" || raw == "yes"
}
