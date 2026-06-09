package inventory

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ijry/lyshop/config"
	"gorm.io/gorm"
)

var externalHTTPClient = &http.Client{Timeout: 3 * time.Second}

type signedPayload struct {
	AppKey    string
	Timestamp string
	Nonce     string
	Body      string
}

type callbackEnvelope struct {
	AppKey    string `json:"app_key"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Sign      string `json:"sign"`
	Body      string `json:"body"`
}

func CallbackEnvelopeFromRequest(appKey, timestamp, nonce, sign, body string) callbackEnvelope {
	return callbackEnvelope{
		AppKey:    appKey,
		Timestamp: timestamp,
		Nonce:     nonce,
		Sign:      sign,
		Body:      body,
	}
}

type externalProvider struct{}

func init() {
	Register(&externalProvider{})
}

func (p *externalProvider) Name() string { return "external_wms" }

func (p *externalProvider) ReserveTx(tx *gorm.DB, in ReserveInput) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", in.BizType, in.BizNo, "reserve", TaskPayload{Items: in.Items})
	}
	return p.postJSON(context.Background(), "/reserve", map[string]any{
		"biz_type": in.BizType,
		"biz_no":   in.BizNo,
		"items":    in.Items,
	})
}

func (p *externalProvider) ConfirmTx(tx *gorm.DB, bizType, bizNo string) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", bizType, bizNo, "confirm", TaskPayload{})
	}
	return p.postJSON(context.Background(), "/confirm", map[string]any{
		"biz_type": bizType,
		"biz_no":   bizNo,
	})
}

func (p *externalProvider) ReleaseTx(tx *gorm.DB, bizType, bizNo, reason string) error {
	if strings.ToLower(strings.TrimSpace(config.Global.Inventory.ExternalMode)) == "async" {
		return EnqueueInventoryTask(tx, "external_wms", bizType, bizNo, "release", TaskPayload{Reason: reason})
	}
	return p.postJSON(context.Background(), "/release", map[string]any{
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

func (p *externalProvider) postJSON(ctx context.Context, path string, payload map[string]any) error {
	raw, _ := json.Marshal(payload)
	_, err := p.doSignedRequest(ctx, http.MethodPost, path, raw)
	if err != nil {
		return err
	}
	return nil
}

func (p *externalProvider) doSignedRequest(ctx context.Context, method, path string, raw []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(config.Global.ExternalWMS.Endpoint, "/")+path, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	sign := BuildSignature(signedPayload{
		AppKey:    config.Global.ExternalWMS.AppKey,
		Timestamp: timestamp,
		Nonce:     nonce,
		Body:      string(raw),
	}, config.Global.ExternalWMS.AppSecret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Key", config.Global.ExternalWMS.AppKey)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Nonce", nonce)
	req.Header.Set("X-Sign", sign)
	resp, err := externalHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("external wms http status %d", resp.StatusCode)
	}
	return body, nil
}

func BuildSignature(req signedPayload, secret string) string {
	raw := strings.Join([]string{
		req.AppKey,
		req.Timestamp,
		req.Nonce,
		req.Body,
		secret,
	}, "\n")
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func VerifyCallbackSignature(req callbackEnvelope, secret string, now time.Time) error {
	ts, err := strconv.ParseInt(req.Timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %w", err)
	}
	ttl := config.Global.ExternalWMS.SignatureTTL
	if ttl <= 0 {
		ttl = 300
	}
	if absDuration(now.Sub(time.Unix(ts, 0))) > time.Duration(ttl)*time.Second {
		return fmt.Errorf("signature expired")
	}

	expected := BuildSignature(signedPayload{
		AppKey:    req.AppKey,
		Timestamp: req.Timestamp,
		Nonce:     req.Nonce,
		Body:      req.Body,
	}, secret)
	if !strings.EqualFold(expected, req.Sign) {
		return fmt.Errorf("signature mismatch")
	}
	return nil
}

func absDuration(v time.Duration) time.Duration {
	if v < 0 {
		return -v
	}
	return v
}
