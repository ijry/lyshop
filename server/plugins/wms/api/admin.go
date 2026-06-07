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
	g.GET("/wms/stocks/by-skus", middleware.RequirePermission("wms:view"), listStocksBySkuIDs)
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
	page, size, err := parsePageSize(c)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
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
		writeServiceError(c, err)
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
		writeServiceError(c, err)
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
		writeServiceError(c, err)
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
		Status *int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if req.Status == nil {
		response.Fail(c, 400, "status 必填")
		return
	}
	if err := wmssvc.UpdateWarehouseStatus(c.Request.Context(), id, *req.Status); err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, nil)
}

func listStocks(c *gin.Context) {
	page, size, err := parsePageSize(c)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	warehouseID, err := parseOptionalUint64Query(c.Query("warehouse_id"), "warehouse_id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	skuID, err := parseOptionalUint64Query(c.Query("sku_id"), "sku_id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	warningOnly, err := parseOptionalBoolQuery(c.Query("warning_only"), "warning_only")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	list, total, err := wmssvc.ListStocks(c.Request.Context(), wmssvc.StockListQuery{
		Page:        page,
		Size:        size,
		WarehouseID: warehouseID,
		SkuID:       skuID,
		WarningOnly: warningOnly,
	})
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func listStocksBySkuIDs(c *gin.Context) {
	rawIDs := strings.TrimSpace(c.Query("sku_ids"))
	if rawIDs == "" {
		response.OK(c, []wmsmodel.InventoryStock{})
		return
	}
	parts := strings.Split(rawIDs, ",")
	skuIDs := make([]uint64, 0, len(parts))
	for _, p := range parts {
		id, err := strconv.ParseUint(strings.TrimSpace(p), 10, 64)
		if err != nil || id == 0 {
			continue
		}
		skuIDs = append(skuIDs, id)
	}
	if len(skuIDs) == 0 {
		response.OK(c, []wmsmodel.InventoryStock{})
		return
	}
	list, err := wmssvc.ListStocksBySkuIDs(c.Request.Context(), skuIDs)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, list)
}

func updateStockSafety(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	var req struct {
		SafeQty *int `json:"safe_qty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if req.SafeQty == nil {
		response.Fail(c, 400, "safe_qty 必填")
		return
	}
	if err := wmssvc.UpdateStockSafety(c.Request.Context(), id, *req.SafeQty); err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, nil)
}

func listDocs(c *gin.Context) {
	page, size, err := parsePageSize(c)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	warehouseID, err := parseOptionalUint64Query(c.Query("warehouse_id"), "warehouse_id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
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
		writeServiceError(c, err)
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
		writeServiceError(c, err)
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
		writeServiceError(c, err)
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
		writeServiceError(c, err)
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
		writeServiceError(c, err)
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
		writeServiceError(c, err)
		return
	}
	response.OK(c, nil)
}

func listMovements(c *gin.Context) {
	page, size, err := parsePageSize(c)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	warehouseID, err := parseOptionalUint64Query(c.Query("warehouse_id"), "warehouse_id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	skuID, err := parseOptionalUint64Query(c.Query("sku_id"), "sku_id")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
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
		writeServiceError(c, err)
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func writeServiceError(c *gin.Context, err error) {
	switch wmssvc.ErrorKindOf(err) {
	case wmssvc.ErrorKindInvalid:
		response.Fail(c, 400, err.Error())
	case wmssvc.ErrorKindNotFound:
		response.Fail(c, 404, err.Error())
	case wmssvc.ErrorKindConflict:
		response.Fail(c, 409, err.Error())
	case wmssvc.ErrorKindForbidden:
		response.Fail(c, 403, err.Error())
	default:
		response.Fail(c, 500, err.Error())
	}
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

func parseOptionalUint64Query(raw, key string) (uint64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	val, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s 参数非法", key)
	}
	return val, nil
}

func parsePageSize(c *gin.Context) (int, int, error) {
	page, err := parsePositiveIntQuery(c.Query("page"), "page", 1, 1000000)
	if err != nil {
		return 0, 0, err
	}
	size, err := parsePositiveIntQuery(c.Query("size"), "size", 20, 100)
	if err != nil {
		return 0, 0, err
	}
	return page, size, nil
}

func parsePositiveIntQuery(raw, key string, defaultVal, maxVal int) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultVal, nil
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s 参数非法", key)
	}
	if val <= 0 || (maxVal > 0 && val > maxVal) {
		return 0, fmt.Errorf("%s 参数非法", key)
	}
	return val, nil
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

func parseOptionalBoolQuery(raw, key string) (bool, error) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return false, nil
	}
	if raw == "1" || raw == "true" {
		return true, nil
	}
	if raw == "0" || raw == "false" {
		return false, nil
	}
	return false, fmt.Errorf("%s 参数非法", key)
}
