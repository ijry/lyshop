package model

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/model"
)

type AfterSaleCaseType string

const (
	AfterSaleCaseTypeReturn   AfterSaleCaseType = "return"
	AfterSaleCaseTypeExchange AfterSaleCaseType = "exchange"
)

type AfterSaleStatus string

const (
	AfterSaleStatusApplied            AfterSaleStatus = "applied"
	AfterSaleStatusApprovedWaitReturn AfterSaleStatus = "approved_wait_user_return"
	AfterSaleStatusUserReturning      AfterSaleStatus = "user_returning"
	AfterSaleStatusWarehouseReceived  AfterSaleStatus = "warehouse_received"
	AfterSaleStatusRefundPending      AfterSaleStatus = "refund_pending"
	AfterSaleStatusRefunded           AfterSaleStatus = "refunded"
	AfterSaleStatusReshipPending      AfterSaleStatus = "reship_pending"
	AfterSaleStatusReshipped          AfterSaleStatus = "reshipped"
	AfterSaleStatusCompleted          AfterSaleStatus = "completed"
	AfterSaleStatusRejected           AfterSaleStatus = "rejected"
	AfterSaleStatusClosed             AfterSaleStatus = "closed"
)

type AfterSaleAuditStatus string

const (
	AfterSaleAuditPending  AfterSaleAuditStatus = "pending"
	AfterSaleAuditApproved AfterSaleAuditStatus = "approved"
	AfterSaleAuditRejected AfterSaleAuditStatus = "rejected"
)

type AfterSaleCase struct {
	model.Base
	OrderID             uint64     `gorm:"not null;index"                          json:"order_id"`
	UserID              uint64     `gorm:"not null;index"                          json:"user_id"`
	MerchantID          uint64     `gorm:"not null;default:0;index"                json:"merchant_id"`
	OrderStatusSnapshot int8       `gorm:"not null;default:0"                     json:"order_status_snapshot"`
	CaseNo              string     `gorm:"size:64;not null;uniqueIndex"            json:"case_no"`
	CaseType            string     `gorm:"size:16;not null;index"                  json:"case_type"`
	Status              string     `gorm:"size:64;not null;index"                  json:"status"`
	Reason              string     `gorm:"size:255"                                json:"reason"`
	ApplyContent        string     `gorm:"type:text"                               json:"apply_content"`
	ApplyImagesJSON     string     `gorm:"type:text"                               json:"-"`
	AuditStatus         string     `gorm:"size:16;not null;default:'pending'"      json:"audit_status"`
	AuditRemark         string     `gorm:"size:255"                                json:"audit_remark"`
	RefundAmount        float64    `gorm:"type:decimal(10,2);default:0"           json:"refund_amount"`
	CloseReason         string     `gorm:"size:255"                                json:"close_reason"`
	CompletedAt         *time.Time `json:"completed_at"`
	ApplyImages         []string   `gorm:"-"                                        json:"apply_images,omitempty"`
}

func (c *AfterSaleCase) ApplyImagesRaw() string {
	if len(c.ApplyImages) == 0 {
		return "[]"
	}
	buf, _ := json.Marshal(c.ApplyImages)
	return string(buf)
}

type AfterSaleCaseItem struct {
	model.Base
	CaseID      uint64 `gorm:"not null;index"   json:"case_id"`
	OrderItemID uint64 `gorm:"not null;index"   json:"order_item_id"`
	Qty         int    `gorm:"not null"         json:"qty"`
}

type AfterSaleLog struct {
	model.Base
	CaseID       uint64 `gorm:"not null;index"   json:"case_id"`
	FromStatus   string `gorm:"size:64"          json:"from_status"`
	ToStatus     string `gorm:"size:64;not null;index" json:"to_status"`
	Action       string `gorm:"size:32;not null;index" json:"action"`
	OperatorType string `gorm:"size:16;not null" json:"operator_type"`
	OperatorID   uint64 `gorm:"not null;default:0" json:"operator_id"`
	Content      string `gorm:"size:255"        json:"content"`
	ExtJSON      string `gorm:"type:text"       json:"-"`
}
