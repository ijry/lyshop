package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ijry/lyshop/config"
	"gorm.io/gorm"
)

var externalHTTPClient = &http.Client{}

type externalProvider struct {
	adapter ExternalAdapter
}

func init() {
	Register(&externalProvider{adapter: NewExternalAdapter()})
}

func (p *externalProvider) Name() string { return "external_wms" }

func (p *externalProvider) ReserveTx(tx *gorm.DB, in ReserveInput) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", in.BizType, in.BizNo, "reserve", TaskPayload{Items: in.Items})
	}
	_, err := p.callRemote(context.Background(), "/reserve", map[string]any{
		"biz_type": in.BizType,
		"biz_no":   in.BizNo,
		"items":    in.Items,
	})
	return err
}

func (p *externalProvider) ConfirmTx(tx *gorm.DB, bizType, bizNo string) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", bizType, bizNo, "confirm", TaskPayload{})
	}
	_, err := p.callRemote(context.Background(), "/confirm", map[string]any{
		"biz_type": bizType,
		"biz_no":   bizNo,
	})
	return err
}

func (p *externalProvider) ReleaseTx(tx *gorm.DB, bizType, bizNo, reason string) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", bizType, bizNo, "release", TaskPayload{Reason: reason})
	}
	_, err := p.callRemote(context.Background(), "/release", map[string]any{
		"biz_type": bizType,
		"biz_no":   bizNo,
		"reason":   reason,
	})
	return err
}

func (p *externalProvider) DeductTx(tx *gorm.DB, in DeductInput) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", in.BizType, in.BizNo, "deduct", TaskPayload{Items: in.Items})
	}
	_, err := p.callRemote(context.Background(), "/deduct", map[string]any{
		"biz_type": in.BizType,
		"biz_no":   in.BizNo,
		"items":    in.Items,
	})
	return err
}

func (p *externalProvider) RestoreTx(tx *gorm.DB, in RestoreInput) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", in.BizType, in.BizNo, "restore", TaskPayload{
			Items:  in.Items,
			Reason: in.Reason,
		})
	}
	_, err := p.callRemote(context.Background(), "/restore", map[string]any{
		"biz_type": in.BizType,
		"biz_no":   in.BizNo,
		"items":    in.Items,
		"reason":   in.Reason,
	})
	return err
}

func (p *externalProvider) SyncSkuTx(tx *gorm.DB, in SyncSkuInput) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", in.Source, fmt.Sprintf("sku:%d", in.SkuID), "sync_sku", TaskPayload{
			SkuID: in.SkuID,
			Stock: in.Stock,
		})
	}
	_, err := p.callRemote(context.Background(), "/sync-sku", map[string]any{
		"source": in.Source,
		"sku_id": in.SkuID,
		"stock":  in.Stock,
	})
	return err
}

func (p *externalProvider) GetSellableStock(ctx context.Context, skuIDs []uint64) ([]SellableStock, error) {
	resp, err := p.callRemote(ctx, "/stock/sellable", map[string]any{"sku_ids": skuIDs})
	if err != nil {
		return nil, err
	}
	var stocks []SellableStock
	if len(resp.Data) > 0 {
		if err := json.Unmarshal(resp.Data, &stocks); err != nil {
			return nil, err
		}
	}
	return stocks, nil
}

func (p *externalProvider) ProcessTask(_ *gorm.DB, task *InventoryIntegrationTask, _ time.Time) error {
	payload, err := decodeTaskPayload(task.Payload)
	if err != nil {
		return err
	}
	resp, err := p.callRemote(context.Background(), externalActionPath(task.Action), buildTaskRequest(task, payload))
	if err != nil {
		return err
	}
	if strings.TrimSpace(resp.RequestID) != "" {
		task.RequestID = resp.RequestID
	}
	return nil
}

func (p *externalProvider) callRemote(ctx context.Context, path string, payload any) (*genericExternalResponse, error) {
	if p.adapter == nil {
		p.adapter = NewExternalAdapter()
	}
	req, err := p.adapter.BuildSignedRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return nil, err
	}
	resp, err := externalHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parsed, err := p.adapter.ParseResponse(resp)
	if err != nil {
		return nil, err
	}
	if !p.adapter.IsSuccess(parsed.Code) {
		if p.adapter.IsRetryable(parsed.Code) {
			return nil, ErrInventoryBusy
		}
		return nil, &externalWMSError{Code: parsed.Code, Message: p.adapter.ErrorMessage(parsed)}
	}
	return parsed, nil
}

func decodeTaskPayload(raw string) (TaskPayload, error) {
	var payload TaskPayload
	if strings.TrimSpace(raw) == "" {
		return payload, nil
	}
	err := json.Unmarshal([]byte(raw), &payload)
	return payload, err
}

func externalActionPath(action string) string {
	switch action {
	case "reserve":
		return "/reserve"
	case "confirm":
		return "/confirm"
	case "release":
		return "/release"
	case "deduct":
		return "/deduct"
	case "restore":
		return "/restore"
	case "sync_sku":
		return "/sync-sku"
	default:
		return "/" + action
	}
}

func buildTaskRequest(task *InventoryIntegrationTask, payload TaskPayload) map[string]any {
	req := map[string]any{
		"biz_type": task.BizType,
		"biz_no":   task.BizNo,
	}
	if len(payload.Items) > 0 {
		req["items"] = payload.Items
	}
	if payload.Reason != "" {
		req["reason"] = payload.Reason
	}
	if payload.SkuID > 0 {
		req["sku_id"] = payload.SkuID
	}
	if payload.Stock != 0 {
		req["stock"] = payload.Stock
	}
	if task.RequestID != "" {
		req["request_id"] = task.RequestID
	}
	return req
}
