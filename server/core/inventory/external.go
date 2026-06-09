package inventory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ijry/lyshop/config"
	"gorm.io/gorm"
)

var externalHTTPClient = &http.Client{Timeout: 3 * time.Second}

type externalProvider struct{}

func init() {
	Register(&externalProvider{})
}

func (p *externalProvider) Name() string { return "external_wms" }

func (p *externalProvider) ReserveTx(tx *gorm.DB, in ReserveInput) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", in.BizType, in.BizNo, "reserve", TaskPayload{Items: in.Items})
	}
	return p.postJSON("/reserve", map[string]any{
		"biz_type": in.BizType,
		"biz_no":   in.BizNo,
		"items":    in.Items,
	})
}

func (p *externalProvider) ConfirmTx(tx *gorm.DB, bizType, bizNo string) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", bizType, bizNo, "confirm", TaskPayload{})
	}
	return p.postJSON("/confirm", map[string]any{
		"biz_type": bizType,
		"biz_no":   bizNo,
	})
}

func (p *externalProvider) ReleaseTx(tx *gorm.DB, bizType, bizNo, reason string) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", bizType, bizNo, "release", TaskPayload{Reason: reason})
	}
	return p.postJSON("/release", map[string]any{
		"biz_type": bizType,
		"biz_no":   bizNo,
		"reason":   reason,
	})
}

func (p *externalProvider) DeductTx(tx *gorm.DB, in DeductInput) error {
	return EnqueueInventoryTask(tx, "external_wms", in.BizType, in.BizNo, "deduct", TaskPayload{Items: in.Items})
}

func (p *externalProvider) RestoreTx(tx *gorm.DB, in RestoreInput) error {
	return EnqueueInventoryTask(tx, "external_wms", in.BizType, in.BizNo, "restore", TaskPayload{
		Items:  in.Items,
		Reason: in.Reason,
	})
}

func (p *externalProvider) SyncSkuTx(tx *gorm.DB, in SyncSkuInput) error {
	return EnqueueInventoryTask(tx, "external_wms", in.Source, fmt.Sprintf("sku:%d", in.SkuID), "sync_sku", TaskPayload{
		SkuID: in.SkuID,
		Stock: in.Stock,
	})
}

func (p *externalProvider) GetSellableStock(_ context.Context, _ []uint64) ([]SellableStock, error) {
	return []SellableStock{}, nil
}

func (p *externalProvider) postJSON(path string, payload map[string]any) error {
	raw, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(config.Global.ExternalWMS.Endpoint, "/")+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := externalHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("external wms http status %d", resp.StatusCode)
	}
	return nil
}
