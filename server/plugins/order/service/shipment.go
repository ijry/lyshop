package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ijry/lyshop/core/db"
	logisticsDriver "github.com/ijry/lyshop/core/driver/logistics"
	ordermodel "github.com/ijry/lyshop/plugins/order/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	company := strings.ToUpper(strings.TrimSpace(req.Company))
	if company == "" {
		return nil, errors.New("请选择快递公司")
	}
	shipDirection := normalizeShipmentDirection(req.Direction)
	shipBizType := normalizeShipmentBizType(req.ShipType, shipDirection)
	now := time.Now()
	row := &ordermodel.OrderShipment{
		OrderID:         req.OrderID,
		AfterSaleCaseID: req.AfterSaleCaseID,
		Direction:       shipDirection,
		BizType:         shipBizType,
		Company:         company,
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

type SyncShipmentReq struct {
	Manual bool `json:"manual"`
}

func SyncShipmentTracks(ctx context.Context, shipmentID uint64, _ SyncShipmentReq) error {
	if shipmentID == 0 {
		return errors.New("运单不存在")
	}
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var shipment ordermodel.OrderShipment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", shipmentID).
			First(&shipment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("运单不存在")
			}
			return err
		}

		driver, provider, err := logisticsDriver.ResolveByPinnedOrFallback(shipment.ChannelProvider)
		if err != nil {
			return markShipmentSyncFail(ctx, tx, &shipment, provider, err, 0)
		}

		start := time.Now()
		result, err := driver.Query(ctx, logisticsDriver.QueryReq{
			CompanyCode: shipment.Company,
			TrackingNo:  shipment.TrackingNo,
		})
		cost := time.Since(start).Milliseconds()
		if err != nil {
			return markShipmentSyncFail(ctx, tx, &shipment, provider, err, cost)
		}
		return applyShipmentTrackResult(ctx, tx, &shipment, provider, result, cost)
	})
}

func ListShipmentTracks(ctx context.Context, orderID, shipmentID uint64) ([]ordermodel.OrderShipmentTrack, error) {
	if orderID == 0 || shipmentID == 0 {
		return []ordermodel.OrderShipmentTrack{}, nil
	}
	var shipment ordermodel.OrderShipment
	if err := db.DB.WithContext(ctx).Where("id = ? AND order_id = ?", shipmentID, orderID).First(&shipment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("运单不存在")
		}
		return nil, err
	}
	var rows []ordermodel.OrderShipmentTrack
	if err := db.DB.WithContext(ctx).
		Where("shipment_id = ?", shipmentID).
		Order("event_time desc, id desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func PollAndSyncShipments(ctx context.Context, limit int) (int, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	var rows []ordermodel.OrderShipment
	if err := db.DB.WithContext(ctx).
		Where("tracking_no <> ''").
		Where("logistics_status <> ?", string(ordermodel.ShipmentStatusSigned)).
		Order("id desc").
		Limit(limit).
		Find(&rows).Error; err != nil {
		return 0, err
	}
	success := 0
	for _, row := range rows {
		if err := SyncShipmentTracks(ctx, row.ID, SyncShipmentReq{Manual: false}); err == nil {
			success++
		}
	}
	return success, nil
}

func markShipmentSyncFail(ctx context.Context, tx *gorm.DB, shipment *ordermodel.OrderShipment, provider string, err error, costMS int64) error {
	now := time.Now()
	if provider == "" {
		provider = strings.TrimSpace(shipment.ChannelProvider)
	}
	if err := tx.WithContext(ctx).Model(&ordermodel.OrderShipment{}).
		Where("id = ?", shipment.ID).
		Updates(map[string]any{
			"last_query_at":   &now,
			"sync_fail_count": gorm.Expr("sync_fail_count + 1"),
		}).Error; err != nil {
		return err
	}
	return tx.WithContext(ctx).Create(&ordermodel.OrderShipmentSyncLog{
		ShipmentID:     shipment.ID,
		Provider:       provider,
		Success:        false,
		ErrorCode:      "query_failed",
		ErrorMessage:   strings.TrimSpace(err.Error()),
		CostMS:         costMS,
		ResponseDigest: "",
	}).Error
}

func applyShipmentTrackResult(ctx context.Context, tx *gorm.DB, shipment *ordermodel.OrderShipment, provider string, result *logisticsDriver.TrackResult, costMS int64) error {
	now := time.Now()
	channelProvider := strings.TrimSpace(shipment.ChannelProvider)
	if channelProvider == "" {
		channelProvider = strings.TrimSpace(provider)
	}
	statusCode := normalizeShipmentStatusCode(result.StatusCode)
	updates := map[string]any{
		"channel_provider": channelProvider,
		"logistics_status": statusCode,
		"last_query_at":    &now,
		"last_sync_ok_at":  &now,
		"sync_fail_count":  0,
	}
	if result.SignedAt != nil {
		updates["signed_at"] = result.SignedAt
	} else if statusCode == string(ordermodel.ShipmentStatusSigned) {
		updates["signed_at"] = &now
	}
	if err := tx.WithContext(ctx).Model(&ordermodel.OrderShipment{}).Where("id = ?", shipment.ID).Updates(updates).Error; err != nil {
		return err
	}

	for _, node := range result.Nodes {
		hash := trackNodeHash(shipment.ID, node)
		row := ordermodel.OrderShipmentTrack{
			ShipmentID: shipment.ID,
			Provider:   channelProvider,
			TrackHash:  hash,
			StatusCode: normalizeShipmentStatusCode(node.StatusCode),
			StatusText: strings.TrimSpace(node.StatusText),
			EventTime:  node.Time,
			Location:   strings.TrimSpace(node.Location),
			RawPayload: node.RawPayload,
		}
		if row.StatusCode == "" {
			row.StatusCode = statusCode
		}
		if row.EventTime.IsZero() {
			row.EventTime = now
		}
		if row.StatusText == "" {
			row.StatusText = strings.TrimSpace(result.StatusText)
		}
		if row.RawPayload == nil {
			raw, _ := json.Marshal(node)
			row.RawPayload = raw
		}
		var existing ordermodel.OrderShipmentTrack
		if err := tx.WithContext(ctx).
			Where("shipment_id = ? AND track_hash = ?", shipment.ID, hash).
			First(&existing).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err := tx.WithContext(ctx).Create(&row).Error; err != nil {
				return err
			}
		}
	}

	return tx.WithContext(ctx).Create(&ordermodel.OrderShipmentSyncLog{
		ShipmentID:     shipment.ID,
		Provider:       channelProvider,
		Success:        true,
		ErrorCode:      "",
		ErrorMessage:   "",
		CostMS:         costMS,
		ResponseDigest: fmt.Sprintf("status=%s,nodes=%d", statusCode, len(result.Nodes)),
	}).Error
}

func normalizeShipmentStatusCode(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case string(ordermodel.ShipmentStatusPending):
		return string(ordermodel.ShipmentStatusPending)
	case string(ordermodel.ShipmentStatusInTransit), "transit":
		return string(ordermodel.ShipmentStatusInTransit)
	case string(ordermodel.ShipmentStatusSigned), "delivered":
		return string(ordermodel.ShipmentStatusSigned)
	case string(ordermodel.ShipmentStatusException), "problem":
		return string(ordermodel.ShipmentStatusException)
	default:
		return string(ordermodel.ShipmentStatusShipped)
	}
}

func trackNodeHash(shipmentID uint64, node logisticsDriver.TrackNode) string {
	source := fmt.Sprintf("%d|%d|%s|%s|%s",
		shipmentID,
		node.Time.Unix(),
		strings.TrimSpace(node.Location),
		strings.TrimSpace(node.StatusCode),
		strings.TrimSpace(node.StatusText),
	)
	sum := sha1.Sum([]byte(source))
	return hex.EncodeToString(sum[:])
}
