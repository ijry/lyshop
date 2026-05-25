package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	ordermodel "github.com/ijry/lyshop/plugins/order/model"
	vipsvc "github.com/ijry/lyshop/plugins/vip/service"
)

type AfterSaleItemReq struct {
	OrderItemID uint64 `json:"order_item_id"`
	Qty         int    `json:"qty"`
}

type CreateAfterSaleReq struct {
	OrderID      uint64             `json:"order_id"`
	UserID       uint64             `json:"user_id"`
	CaseType     string             `json:"case_type"`
	Reason       string             `json:"reason"`
	ApplyContent string             `json:"apply_content"`
	ApplyImages  []string           `json:"apply_images"`
	Items        []AfterSaleItemReq `json:"items"`
}

type AuditAfterSaleReq struct {
	Approve     bool   `json:"approve"`
	AuditRemark string `json:"audit_remark"`
	AdminID     uint64 `json:"admin_id"`
}

type SubmitReturnShipmentReq struct {
	Company    string `json:"company"`
	TrackingNo string `json:"tracking_no"`
	Remark     string `json:"remark"`
	UserID     uint64 `json:"user_id"`
}

type MarkRefundReq struct {
	Amount   float64 `json:"amount"`
	Reason   string  `json:"reason"`
	RefundNo string  `json:"refund_no"`
	AdminID  uint64  `json:"admin_id"`
}

type CloseAfterSaleReq struct {
	Reason  string `json:"reason"`
	AdminID uint64 `json:"admin_id"`
}

type AfterSaleSummary struct {
	InProgressCount   int64  `json:"in_progress_count"`
	HasOpenCase       bool   `json:"has_open_case"`
	LatestStatus      string `json:"latest_status,omitempty"`
	LatestStatusLabel string `json:"latest_status_label,omitempty"`
	LatestCaseID      uint64 `json:"latest_case_id,omitempty"`
	CanApply          bool   `json:"can_apply"`
}

type AfterSaleCaseListView struct {
	ordermodel.AfterSaleCase
	StatusLabel   string `json:"status_label,omitempty"`
	CaseTypeLabel string `json:"case_type_label,omitempty"`
}

type AfterSaleLogView struct {
	ordermodel.AfterSaleLog
	FromStatusLabel string `json:"from_status_label,omitempty"`
	ToStatusLabel   string `json:"to_status_label,omitempty"`
	ActionLabel     string `json:"action_label,omitempty"`
}

type AfterSaleShipmentView struct {
	ordermodel.OrderShipment
	DirectionLabel       string `json:"direction_label,omitempty"`
	BizTypeLabel         string `json:"biz_type_label,omitempty"`
	LogisticsStatusLabel string `json:"logistics_status_label,omitempty"`
}

type AfterSaleCaseView struct {
	ordermodel.AfterSaleCase
	StatusLabel   string                         `json:"status_label,omitempty"`
	CaseTypeLabel string                         `json:"case_type_label,omitempty"`
	Items         []ordermodel.AfterSaleCaseItem `json:"items"`
	Logs          []AfterSaleLogView             `json:"logs"`
	Shipments     []AfterSaleShipmentView        `json:"shipments"`
}

func generateAfterSaleCaseNo() string {
	return fmt.Sprintf("AS%d%04d", time.Now().UnixMilli(), time.Now().Nanosecond()%10000)
}

func normalizeAfterSaleType(caseType string) string {
	switch strings.ToLower(strings.TrimSpace(caseType)) {
	case string(ordermodel.AfterSaleCaseTypeExchange):
		return string(ordermodel.AfterSaleCaseTypeExchange)
	default:
		return string(ordermodel.AfterSaleCaseTypeReturn)
	}
}

func encodeStringArray(values []string) string {
	out := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, raw := range values {
		v := strings.TrimSpace(raw)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
		if len(out) >= 9 {
			break
		}
	}
	if len(out) == 0 {
		return "[]"
	}
	buf, _ := json.Marshal(out)
	return string(buf)
}

func decodeStringArray(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}
	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return []string{}
	}
	out := make([]string, 0, len(arr))
	seen := map[string]struct{}{}
	for _, item := range arr {
		v := strings.TrimSpace(item)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func isAfterSaleStatusOpen(status string) bool {
	switch status {
	case string(ordermodel.AfterSaleStatusCompleted), string(ordermodel.AfterSaleStatusRejected), string(ordermodel.AfterSaleStatusClosed):
		return false
	default:
		return true
	}
}

func statusAllowClose(status string) bool {
	switch status {
	case string(ordermodel.AfterSaleStatusApplied),
		string(ordermodel.AfterSaleStatusApprovedWaitReturn),
		string(ordermodel.AfterSaleStatusUserReturning),
		string(ordermodel.AfterSaleStatusWarehouseReceived),
		string(ordermodel.AfterSaleStatusRefundPending),
		string(ordermodel.AfterSaleStatusReshipPending):
		return true
	default:
		return false
	}
}

func writeAfterSaleLogTx(tx *gorm.DB, caseID uint64, fromStatus, toStatus, action, operatorType string, operatorID uint64, content string, ext map[string]any) error {
	extJSON := ""
	if len(ext) > 0 {
		buf, _ := json.Marshal(ext)
		extJSON = string(buf)
	}
	row := &ordermodel.AfterSaleLog{
		CaseID:       caseID,
		FromStatus:   fromStatus,
		ToStatus:     toStatus,
		Action:       action,
		OperatorType: operatorType,
		OperatorID:   operatorID,
		Content:      strings.TrimSpace(content),
		ExtJSON:      extJSON,
	}
	return tx.Create(row).Error
}

func lockAfterSaleCase(tx *gorm.DB, caseID uint64) (*ordermodel.AfterSaleCase, error) {
	var row ordermodel.AfterSaleCase
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", caseID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("afterSale.err.notFound")
		}
		return nil, err
	}
	return &row, nil
}

func refreshOrderStatusByAfterSaleTx(tx *gorm.DB, orderID uint64) error {
	var rows []ordermodel.AfterSaleCase
	if err := tx.Where("order_id = ?", orderID).Find(&rows).Error; err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	hasOpen := false
	restoreStatus := int8(ordermodel.OrderStatusCompleted)
	for _, row := range rows {
		if row.OrderStatusSnapshot > 0 && restoreStatus < row.OrderStatusSnapshot {
			restoreStatus = row.OrderStatusSnapshot
		}
		if isAfterSaleStatusOpen(row.Status) {
			hasOpen = true
		}
	}
	if hasOpen {
		return tx.Model(&ordermodel.Order{}).Where("id = ?", orderID).Update("status", ordermodel.OrderStatusAfterSale).Error
	}
	if restoreStatus <= 0 {
		restoreStatus = ordermodel.OrderStatusCompleted
	}
	return tx.Model(&ordermodel.Order{}).Where("id = ?", orderID).Update("status", restoreStatus).Error
}

// CreateAfterSale creates a new after-sale case for an order.
func CreateAfterSale(ctx context.Context, req CreateAfterSaleReq) (uint64, error) {
	if req.OrderID == 0 || req.UserID == 0 {
		return 0, errors.New("afterSale.err.invalidParams")
	}
	caseType := normalizeAfterSaleType(req.CaseType)
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		return 0, errors.New("afterSale.err.reasonRequired")
	}
	if len(req.Items) == 0 {
		return 0, errors.New("afterSale.err.selectProducts")
	}

	var createdID uint64
	err := db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order ordermodel.Order
		if err := tx.Where("id = ? AND user_id = ?", req.OrderID, req.UserID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("afterSale.err.orderNotFound")
			}
			return err
		}
		if order.Status < ordermodel.OrderStatusPaid {
			return errors.New("afterSale.err.orderNotPaid")
		}

		itemQty := map[uint64]int{}
		itemIDs := make([]uint64, 0, len(req.Items))
		for _, item := range req.Items {
			if item.OrderItemID == 0 || item.Qty <= 0 {
				continue
			}
			itemQty[item.OrderItemID] = item.Qty
			itemIDs = append(itemIDs, item.OrderItemID)
		}
		if len(itemQty) == 0 {
			return errors.New("afterSale.err.selectValidProducts")
		}

		var orderItems []ordermodel.OrderItem
		if err := tx.Where("order_id = ? AND id IN ?", req.OrderID, itemIDs).Find(&orderItems).Error; err != nil {
			return err
		}
		if len(orderItems) == 0 {
			return errors.New("afterSale.err.itemNotFound")
		}

		var activeCaseItems []ordermodel.AfterSaleCaseItem
		if err := tx.
			Joins("JOIN after_sale_cases ON after_sale_cases.id = after_sale_case_items.case_id").
			Where("after_sale_cases.order_id = ? AND after_sale_cases.status NOT IN ?", req.OrderID, []string{
				string(ordermodel.AfterSaleStatusCompleted),
				string(ordermodel.AfterSaleStatusClosed),
				string(ordermodel.AfterSaleStatusRejected),
			}).
			Find(&activeCaseItems).Error; err != nil {
			return err
		}
		activeItemMap := map[uint64]bool{}
		for _, row := range activeCaseItems {
			activeItemMap[row.OrderItemID] = true
		}
		for _, item := range orderItems {
			if activeItemMap[item.ID] {
				return fmt.Errorf("afterSale.err.activeCase (item_id=%d)", item.ID)
			}
			if itemQty[item.ID] > item.Qty {
				return fmt.Errorf("afterSale.err.qtyExceeded (item_id=%d)", item.ID)
			}
		}

		caseRow := &ordermodel.AfterSaleCase{
			OrderID:             req.OrderID,
			UserID:              req.UserID,
			MerchantID:          order.MerchantID,
			OrderStatusSnapshot: order.Status,
			CaseNo:              generateAfterSaleCaseNo(),
			CaseType:            caseType,
			Status:              string(ordermodel.AfterSaleStatusApplied),
			Reason:              reason,
			ApplyContent:        strings.TrimSpace(req.ApplyContent),
			ApplyImagesJSON:     encodeStringArray(req.ApplyImages),
			AuditStatus:         string(ordermodel.AfterSaleAuditPending),
		}
		if err := tx.Create(caseRow).Error; err != nil {
			return err
		}
		items := make([]ordermodel.AfterSaleCaseItem, 0, len(orderItems))
		for _, item := range orderItems {
			items = append(items, ordermodel.AfterSaleCaseItem{
				CaseID:      caseRow.ID,
				OrderItemID: item.ID,
				Qty:         itemQty[item.ID],
			})
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		if err := writeAfterSaleLogTx(tx, caseRow.ID, "", caseRow.Status, "apply", "user", req.UserID, "afterSale.log.applied", map[string]any{
			"case_type": caseType,
			"reason":    reason,
		}); err != nil {
			return err
		}
		if err := tx.Model(&ordermodel.Order{}).Where("id = ?", req.OrderID).Update("status", ordermodel.OrderStatusAfterSale).Error; err != nil {
			return err
		}
		createdID = caseRow.ID
		return nil
	})
	if err != nil {
		return 0, err
	}
	return createdID, nil
}

// AuditAfterSale approves or rejects an after-sale case.
func AuditAfterSale(ctx context.Context, caseID uint64, req AuditAfterSaleReq) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		caseRow, err := lockAfterSaleCase(tx, caseID)
		if err != nil {
			return err
		}
		if caseRow.Status != string(ordermodel.AfterSaleStatusApplied) {
			return errors.New("afterSale.err.cannotReview")
		}
		from := caseRow.Status
		updates := map[string]any{
			"audit_remark": strings.TrimSpace(req.AuditRemark),
		}
		if req.Approve {
			updates["status"] = string(ordermodel.AfterSaleStatusApprovedWaitReturn)
			updates["audit_status"] = string(ordermodel.AfterSaleAuditApproved)
		} else {
			updates["status"] = string(ordermodel.AfterSaleStatusRejected)
			updates["audit_status"] = string(ordermodel.AfterSaleAuditRejected)
		}
		if err := tx.Model(&ordermodel.AfterSaleCase{}).Where("id = ?", caseID).Updates(updates).Error; err != nil {
			return err
		}
		to := updates["status"].(string)
		if err := writeAfterSaleLogTx(tx, caseID, from, to, "audit", "admin", req.AdminID, "afterSale.log.reviewed", map[string]any{
			"approve": req.Approve,
		}); err != nil {
			return err
		}
		return refreshOrderStatusByAfterSaleTx(tx, caseRow.OrderID)
	})
}

// SubmitReturnShipment records the user's return shipment.
func SubmitReturnShipment(ctx context.Context, caseID uint64, req SubmitReturnShipmentReq) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		caseRow, err := lockAfterSaleCase(tx, caseID)
		if err != nil {
			return err
		}
		if caseRow.Status != string(ordermodel.AfterSaleStatusApprovedWaitReturn) {
			return errors.New("afterSale.err.cannotReturn")
		}
		if _, err := createShipmentTx(tx, CreateShipmentReq{
			OrderID:         caseRow.OrderID,
			AfterSaleCaseID: caseID,
			Direction:       string(ordermodel.ShipmentDirectionInbound),
			ShipType:        string(ordermodel.ShipmentBizTypeReturn),
			Company:         req.Company,
			TrackingNo:      req.TrackingNo,
			Remark:          req.Remark,
			CreatedByType:   "user",
			CreatedByID:     req.UserID,
		}); err != nil {
			return err
		}
		from := caseRow.Status
		to := string(ordermodel.AfterSaleStatusUserReturning)
		if err := tx.Model(&ordermodel.AfterSaleCase{}).Where("id = ?", caseID).Update("status", to).Error; err != nil {
			return err
		}
		if err := writeAfterSaleLogTx(tx, caseID, from, to, "return_ship", "user", req.UserID, "afterSale.log.returnShipped", map[string]any{
			"company":     strings.TrimSpace(req.Company),
			"tracking_no": strings.TrimSpace(req.TrackingNo),
		}); err != nil {
			return err
		}
		return refreshOrderStatusByAfterSaleTx(tx, caseRow.OrderID)
	})
}

// ReceiveAfterSale confirms warehouse receipt of returned goods.
func ReceiveAfterSale(ctx context.Context, caseID uint64, adminID uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		caseRow, err := lockAfterSaleCase(tx, caseID)
		if err != nil {
			return err
		}
		if caseRow.Status != string(ordermodel.AfterSaleStatusUserReturning) {
			return errors.New("afterSale.err.cannotReceive")
		}
		nextStatus := string(ordermodel.AfterSaleStatusRefundPending)
		if caseRow.CaseType == string(ordermodel.AfterSaleCaseTypeExchange) {
			nextStatus = string(ordermodel.AfterSaleStatusReshipPending)
		}
		if err := tx.Model(&ordermodel.AfterSaleCase{}).Where("id = ?", caseID).Update("status", nextStatus).Error; err != nil {
			return err
		}
		if err := writeAfterSaleLogTx(tx, caseID, caseRow.Status, nextStatus, "receive", "admin", adminID, "afterSale.log.received", nil); err != nil {
			return err
		}
		return refreshOrderStatusByAfterSaleTx(tx, caseRow.OrderID)
	})
}

// MarkRefund records a refund for an after-sale case.
func MarkRefund(ctx context.Context, caseID uint64, req MarkRefundReq) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		caseRow, err := lockAfterSaleCase(tx, caseID)
		if err != nil {
			return err
		}
		if caseRow.Status != string(ordermodel.AfterSaleStatusRefundPending) {
			return errors.New("afterSale.err.cannotRefund")
		}
		amount := req.Amount
		if amount <= 0 {
			if caseRow.RefundAmount > 0 {
				amount = caseRow.RefundAmount
			} else {
				var order ordermodel.Order
				if err := tx.Where("id = ?", caseRow.OrderID).First(&order).Error; err != nil {
					return err
				}
				amount = order.TotalAmount
			}
		}
		if amount <= 0 {
			return errors.New("afterSale.err.invalidRefundAmount")
		}
		refundNo := strings.TrimSpace(req.RefundNo)
		if refundNo == "" {
			refundNo = fmt.Sprintf("RF%d%04d", time.Now().UnixMilli(), time.Now().Nanosecond()%10000)
		}
		refund := &ordermodel.OrderRefund{
			OrderID:         caseRow.OrderID,
			AfterSaleCaseID: caseID,
			Reason:          strings.TrimSpace(req.Reason),
			Amount:          amount,
			Status:          2,
			RefundNo:        refundNo,
		}
		if err := tx.Create(refund).Error; err != nil {
			return err
		}
		if err := tx.Model(&ordermodel.AfterSaleCase{}).Where("id = ?", caseID).Updates(map[string]any{
			"status":        string(ordermodel.AfterSaleStatusRefunded),
			"refund_amount": amount,
		}).Error; err != nil {
			return err
		}
		if err := writeAfterSaleLogTx(tx, caseID, caseRow.Status, string(ordermodel.AfterSaleStatusRefunded), "refund", "admin", req.AdminID, "afterSale.log.refunded", map[string]any{
			"amount":    amount,
			"refund_no": refundNo,
		}); err != nil {
			return err
		}
		if err := vipsvc.RollbackGrowthForRefundTx(tx, caseRow.UserID, caseRow.OrderID, caseRow.ID, amount); err != nil {
			return err
		}
		return refreshOrderStatusByAfterSaleTx(tx, caseRow.OrderID)
	})
}

// CompleteAfterSale marks an after-sale case as completed.
func CompleteAfterSale(ctx context.Context, caseID uint64, adminID uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		caseRow, err := lockAfterSaleCase(tx, caseID)
		if err != nil {
			return err
		}
		if caseRow.Status == string(ordermodel.AfterSaleStatusCompleted) {
			return nil
		}
		allowed := map[string]bool{
			string(ordermodel.AfterSaleStatusRefunded):  true,
			string(ordermodel.AfterSaleStatusReshipped): true,
		}
		if !allowed[caseRow.Status] {
			return errors.New("afterSale.err.cannotComplete")
		}
		now := time.Now()
		if err := tx.Model(&ordermodel.AfterSaleCase{}).Where("id = ?", caseID).Updates(map[string]any{
			"status":       string(ordermodel.AfterSaleStatusCompleted),
			"completed_at": &now,
		}).Error; err != nil {
			return err
		}
		if err := writeAfterSaleLogTx(tx, caseID, caseRow.Status, string(ordermodel.AfterSaleStatusCompleted), "complete", "admin", adminID, "afterSale.log.completed", nil); err != nil {
			return err
		}
		return refreshOrderStatusByAfterSaleTx(tx, caseRow.OrderID)
	})
}

// CloseAfterSale closes an after-sale case.
func CloseAfterSale(ctx context.Context, caseID uint64, req CloseAfterSaleReq) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		caseRow, err := lockAfterSaleCase(tx, caseID)
		if err != nil {
			return err
		}
		if !statusAllowClose(caseRow.Status) {
			return errors.New("afterSale.err.cannotClose")
		}
		closeReason := strings.TrimSpace(req.Reason)
		if closeReason == "" {
			return errors.New("afterSale.err.closeReasonRequired")
		}
		if err := tx.Model(&ordermodel.AfterSaleCase{}).Where("id = ?", caseID).Updates(map[string]any{
			"status":       string(ordermodel.AfterSaleStatusClosed),
			"close_reason": closeReason,
		}).Error; err != nil {
			return err
		}
		if err := writeAfterSaleLogTx(tx, caseID, caseRow.Status, string(ordermodel.AfterSaleStatusClosed), "close", "admin", req.AdminID, "afterSale.log.closed", map[string]any{
			"reason": closeReason,
		}); err != nil {
			return err
		}
		return refreshOrderStatusByAfterSaleTx(tx, caseRow.OrderID)
	})
}

// GetAfterSale returns an after-sale case detail view.
func GetAfterSale(c *gin.Context, caseID uint64) (*AfterSaleCaseView, error) {
	ctx := c.Request.Context()
	var caseRow ordermodel.AfterSaleCase
	if err := db.DB.WithContext(ctx).Where("id = ?", caseID).First(&caseRow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("afterSale.err.notFound")
		}
		return nil, err
	}
	caseRow.ApplyImages = decodeStringArray(caseRow.ApplyImagesJSON)
	var items []ordermodel.AfterSaleCaseItem
	if err := db.DB.WithContext(ctx).Where("case_id = ?", caseID).Order("id asc").Find(&items).Error; err != nil {
		return nil, err
	}
	var logRows []ordermodel.AfterSaleLog
	if err := db.DB.WithContext(ctx).Where("case_id = ?", caseID).Order("id asc").Find(&logRows).Error; err != nil {
		return nil, err
	}
	logs := make([]AfterSaleLogView, 0, len(logRows))
	for _, row := range logRows {
		logs = append(logs, AfterSaleLogView{
			AfterSaleLog:    row,
			FromStatusLabel: afterSaleStatusLabel(c, row.FromStatus),
			ToStatusLabel:   afterSaleStatusLabel(c, row.ToStatus),
			ActionLabel:     afterSaleActionLabel(c, row.Action),
		})
	}
	var shipmentRows []ordermodel.OrderShipment
	if err := db.DB.WithContext(ctx).Where("after_sale_case_id = ?", caseID).Order("id desc").Find(&shipmentRows).Error; err != nil {
		return nil, err
	}
	shipments := make([]AfterSaleShipmentView, 0, len(shipmentRows))
	for _, row := range shipmentRows {
		shipments = append(shipments, AfterSaleShipmentView{
			OrderShipment:        row,
			DirectionLabel:       shipmentDirectionLabel(c, row.Direction),
			BizTypeLabel:         shipmentBizTypeLabel(c, row.BizType),
			LogisticsStatusLabel: shipmentStatusLabel(c, row.LogisticsStatus),
		})
	}
	return &AfterSaleCaseView{
		AfterSaleCase: caseRow,
		StatusLabel:   afterSaleStatusLabel(c, caseRow.Status),
		CaseTypeLabel: afterSaleCaseTypeLabel(c, caseRow.CaseType),
		Items:         items,
		Logs:          logs,
		Shipments:     shipments,
	}, nil
}

// ListAfterSales returns paginated after-sale cases.
func ListAfterSales(c *gin.Context, status string, caseType string, orderID uint64, page int, size int) ([]AfterSaleCaseListView, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	ctx := c.Request.Context()
	tx := db.DB.WithContext(ctx).Model(&ordermodel.AfterSaleCase{})
	status = strings.TrimSpace(status)
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	caseType = strings.TrimSpace(caseType)
	if caseType != "" {
		tx = tx.Where("case_type = ?", caseType)
	}
	if orderID > 0 {
		tx = tx.Where("order_id = ?", orderID)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []ordermodel.AfterSaleCase
	if err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	result := make([]AfterSaleCaseListView, 0, len(list))
	for _, row := range list {
		row.ApplyImages = decodeStringArray(row.ApplyImagesJSON)
		result = append(result, AfterSaleCaseListView{
			AfterSaleCase: row,
			StatusLabel:   afterSaleStatusLabel(c, row.Status),
			CaseTypeLabel: afterSaleCaseTypeLabel(c, row.CaseType),
		})
	}
	return result, total, nil
}

func buildAfterSaleSummaryMap(c *gin.Context, ctx context.Context, orderIDs []uint64) (map[uint64]*AfterSaleSummary, error) {
	result := make(map[uint64]*AfterSaleSummary, len(orderIDs))
	if len(orderIDs) == 0 {
		return result, nil
	}
	var rows []ordermodel.AfterSaleCase
	if err := db.DB.WithContext(ctx).
		Where("order_id IN ?", orderIDs).
		Order("id desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		summary, ok := result[row.OrderID]
		if !ok {
			summary = &AfterSaleSummary{
				CanApply: true,
			}
			result[row.OrderID] = summary
		}
		if summary.LatestCaseID == 0 {
			summary.LatestCaseID = row.ID
			summary.LatestStatus = row.Status
			summary.LatestStatusLabel = afterSaleStatusLabel(c, row.Status)
		}
		if isAfterSaleStatusOpen(row.Status) {
			summary.InProgressCount += 1
			summary.HasOpenCase = true
			summary.CanApply = false
		}
	}
	return result, nil
}
