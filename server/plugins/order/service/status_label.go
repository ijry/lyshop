package service

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/i18n"
)

var orderStatusI18nKeys = map[int8]string{
	1: "order.status.pending",
	2: "order.status.shipped",
	3: "order.status.delivering",
	4: "order.status.completed",
	5: "order.status.afterSale",
}

var afterSaleCaseTypeI18nKeys = map[string]string{
	"return":   "afterSale.type.return",
	"exchange": "afterSale.type.exchange",
}

var afterSaleStatusI18nKeys = map[string]string{
	"applied":                   "afterSale.status.applied",
	"approved_wait_user_return": "afterSale.status.waitReturn",
	"user_returning":            "afterSale.status.returning",
	"warehouse_received":        "afterSale.status.received",
	"refund_pending":            "afterSale.status.refundPending",
	"refunded":                  "afterSale.status.refunded",
	"reship_pending":            "afterSale.status.reshipPending",
	"reshipped":                 "afterSale.status.reshipped",
	"completed":                 "afterSale.status.completed",
	"rejected":                  "afterSale.status.rejected",
	"closed":                    "afterSale.status.closed",
}

var deliveryTypeI18nKeys = map[string]string{
	"express": "delivery.type.express",
	"local":   "delivery.type.local",
}

var shipmentStatusI18nKeys = map[string]string{
	"pending":    "shipment.status.pending",
	"shipped":    "shipment.status.shipped",
	"in_transit": "shipment.status.inTransit",
	"signed":     "shipment.status.signed",
	"exception":  "shipment.status.exception",
}

var shipmentBizTypeI18nKeys = map[string]string{
	"initial": "shipment.bizType.initial",
	"reship":  "shipment.bizType.reship",
	"return":  "shipment.bizType.return",
}

var shipmentDirectionI18nKeys = map[string]string{
	"outbound": "shipment.direction.outbound",
	"inbound":  "shipment.direction.inbound",
}

var afterSaleActionI18nKeys = map[string]string{
	"apply":       "afterSale.action.apply",
	"audit":       "afterSale.action.review",
	"return_ship": "afterSale.action.return",
	"receive":     "afterSale.action.receive",
	"refund":      "afterSale.action.refund",
	"reship":      "afterSale.action.reship",
	"complete":    "afterSale.action.complete",
	"close":       "afterSale.action.close",
}

func deliveryTypeLabel(c *gin.Context, dt string) string {
	value := strings.TrimSpace(dt)
	if key, ok := deliveryTypeI18nKeys[value]; ok {
		return i18n.T(c, key)
	}
	if value == "" {
		return i18n.T(c, "delivery.type.express")
	}
	return value
}

func orderStatusLabel(c *gin.Context, status int8) string {
	if key, ok := orderStatusI18nKeys[status]; ok {
		return i18n.T(c, key)
	}
	if status <= 0 {
		return ""
	}
	return strconv.Itoa(int(status))
}

func afterSaleCaseTypeLabel(c *gin.Context, caseType string) string {
	value := strings.TrimSpace(caseType)
	if key, ok := afterSaleCaseTypeI18nKeys[value]; ok {
		return i18n.T(c, key)
	}
	return value
}

func afterSaleStatusLabel(c *gin.Context, status string) string {
	value := strings.TrimSpace(status)
	if key, ok := afterSaleStatusI18nKeys[value]; ok {
		return i18n.T(c, key)
	}
	return value
}

func shipmentStatusLabel(c *gin.Context, status string) string {
	value := strings.TrimSpace(status)
	if key, ok := shipmentStatusI18nKeys[value]; ok {
		return i18n.T(c, key)
	}
	return value
}

func shipmentBizTypeLabel(c *gin.Context, bizType string) string {
	value := strings.TrimSpace(bizType)
	if key, ok := shipmentBizTypeI18nKeys[value]; ok {
		return i18n.T(c, key)
	}
	return value
}

func shipmentDirectionLabel(c *gin.Context, direction string) string {
	value := strings.TrimSpace(direction)
	if key, ok := shipmentDirectionI18nKeys[value]; ok {
		return i18n.T(c, key)
	}
	return value
}

func afterSaleActionLabel(c *gin.Context, action string) string {
	value := strings.TrimSpace(action)
	if key, ok := afterSaleActionI18nKeys[value]; ok {
		return i18n.T(c, key)
	}
	return value
}
