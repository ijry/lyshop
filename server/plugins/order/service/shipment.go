package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ijry/lyshop/core/db"
	ordermodel "github.com/ijry/lyshop/plugins/order/model"
	"gorm.io/gorm"
)

type CreateShipmentReq struct {
	OrderID         uint64
	AfterSaleCaseID uint64
	ShipType        string
	Direction       string
	Company         string
	TrackingNo      string
	Remark          string
	CreatedByType   string
	CreatedByID     uint64
}

func normalizeShipmentDirection(direction string) string {
	switch strings.ToLower(strings.TrimSpace(direction)) {
	case string(ordermodel.ShipmentDirectionInbound):
		return string(ordermodel.ShipmentDirectionInbound)
	default:
		return string(ordermodel.ShipmentDirectionOutbound)
	}
}

func normalizeShipmentBizType(shipType string, direction string) string {
	if normalizeShipmentDirection(direction) == string(ordermodel.ShipmentDirectionInbound) {
		return string(ordermodel.ShipmentBizTypeReturn)
	}
	switch strings.ToLower(strings.TrimSpace(shipType)) {
	case string(ordermodel.ShipmentBizTypeReship):
		return string(ordermodel.ShipmentBizTypeReship)
	default:
		return string(ordermodel.ShipmentBizTypeInitial)
	}
}

func normalizeOperatorType(op string) string {
	switch strings.ToLower(strings.TrimSpace(op)) {
	case "user":
		return "user"
	case "system":
		return "system"
	default:
		return "admin"
	}
}

func createShipmentTx(tx *gorm.DB, req CreateShipmentReq) (*ordermodel.OrderShipment, error) {
	if req.OrderID == 0 {
		return nil, errors.New("订单不存在")
	}
	trackingNo := strings.TrimSpace(req.TrackingNo)
	if trackingNo == "" {
		return nil, errors.New("请填写快递单号")
	}
	shipDirection := normalizeShipmentDirection(req.Direction)
	shipBizType := normalizeShipmentBizType(req.ShipType, shipDirection)
	now := time.Now()
	row := &ordermodel.OrderShipment{
		OrderID:         req.OrderID,
		AfterSaleCaseID: req.AfterSaleCaseID,
		Direction:       shipDirection,
		BizType:         shipBizType,
		Company:         strings.TrimSpace(req.Company),
		TrackingNo:      trackingNo,
		LogisticsStatus: string(ordermodel.ShipmentStatusShipped),
		Remark:          strings.TrimSpace(req.Remark),
		ShippedAt:       &now,
		CreatedByType:   normalizeOperatorType(req.CreatedByType),
		CreatedByID:     req.CreatedByID,
	}
	if err := tx.Create(row).Error; err != nil {
		return nil, err
	}
	return row, nil
}

func CreateShipment(ctx context.Context, req CreateShipmentReq) (*ordermodel.OrderShipment, error) {
	var row *ordermodel.OrderShipment
	err := db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		row, err = createShipmentTx(tx, req)
		return err
	})
	if err != nil {
		return nil, err
	}
	return row, nil
}

func ListOrderShipments(ctx context.Context, orderID uint64) ([]ordermodel.OrderShipment, error) {
	if orderID == 0 {
		return []ordermodel.OrderShipment{}, nil
	}
	var rows []ordermodel.OrderShipment
	if err := db.DB.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("id desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func listOrderShipmentsMap(ctx context.Context, orderIDs []uint64) (map[uint64][]ordermodel.OrderShipment, error) {
	result := make(map[uint64][]ordermodel.OrderShipment, len(orderIDs))
	if len(orderIDs) == 0 {
		return result, nil
	}
	var rows []ordermodel.OrderShipment
	if err := db.DB.WithContext(ctx).
		Where("order_id IN ?", orderIDs).
		Order("id desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.OrderID] = append(result[row.OrderID], row)
	}
	return result, nil
}
