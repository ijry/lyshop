package service

import (
	"strconv"
	"strings"
)

var orderStatusLabels = map[int8]string{
	1: "待付款",
	2: "待发货",
	3: "待收货",
	4: "已完成",
	5: "售后",
}

var afterSaleCaseTypeLabels = map[string]string{
	"return":   "退货",
	"exchange": "换货",
}

var afterSaleStatusLabels = map[string]string{
	"applied":                   "已申请",
	"approved_wait_user_return": "待用户回寄",
	"user_returning":            "用户回寄中",
	"warehouse_received":        "仓库已收货",
	"refund_pending":            "待退款",
	"refunded":                  "已退款",
	"reship_pending":            "待补发",
	"reshipped":                 "已补发",
	"completed":                 "已完结",
	"rejected":                  "已拒绝",
	"closed":                    "已关闭",
}

var shipmentStatusLabels = map[string]string{
	"pending":    "待揽收",
	"shipped":    "已发货",
	"in_transit": "运输中",
	"signed":     "已签收",
	"exception":  "物流异常",
}

var shipmentBizTypeLabels = map[string]string{
	"initial": "首发",
	"reship":  "补发",
	"return":  "回寄",
}

var shipmentDirectionLabels = map[string]string{
	"outbound": "寄出",
	"inbound":  "回寄",
}

var afterSaleActionLabels = map[string]string{
	"apply":       "提交申请",
	"audit":       "审核",
	"return_ship": "回寄物流",
	"receive":     "确认收货",
	"refund":      "退款",
	"reship":      "补发",
	"complete":    "完结",
	"close":       "关闭",
}

func orderStatusLabel(status int8) string {
	if label, ok := orderStatusLabels[status]; ok {
		return label
	}
	if status <= 0 {
		return ""
	}
	return strconv.Itoa(int(status))
}

func afterSaleCaseTypeLabel(caseType string) string {
	value := strings.TrimSpace(caseType)
	if label, ok := afterSaleCaseTypeLabels[value]; ok {
		return label
	}
	return value
}

func afterSaleStatusLabel(status string) string {
	value := strings.TrimSpace(status)
	if label, ok := afterSaleStatusLabels[value]; ok {
		return label
	}
	return value
}

func shipmentStatusLabel(status string) string {
	value := strings.TrimSpace(status)
	if label, ok := shipmentStatusLabels[value]; ok {
		return label
	}
	return value
}

func shipmentBizTypeLabel(bizType string) string {
	value := strings.TrimSpace(bizType)
	if label, ok := shipmentBizTypeLabels[value]; ok {
		return label
	}
	return value
}

func shipmentDirectionLabel(direction string) string {
	value := strings.TrimSpace(direction)
	if label, ok := shipmentDirectionLabels[value]; ok {
		return label
	}
	return value
}

func afterSaleActionLabel(action string) string {
	value := strings.TrimSpace(action)
	if label, ok := afterSaleActionLabels[value]; ok {
		return label
	}
	return value
}
