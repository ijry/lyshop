package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Provider constants
const (
	ProviderAliyun  = "aliyun"
	ProviderTencent = "tencent"
)

// SmsDriver implements sms.Driver, supporting aliyun and tencent.
type SmsDriver struct {
	Provider  string // "aliyun" | "tencent"
	AccessKey string
	SecretKey string
	SignName  string
}

func (d *SmsDriver) Name() string { return "sms" }

func (d *SmsDriver) Send(ctx context.Context, phone, templateCode string, params map[string]string) error {
	switch d.Provider {
	case ProviderAliyun:
		return d.sendAliyun(ctx, phone, templateCode, params)
	case ProviderTencent:
		return d.sendTencent(ctx, phone, templateCode, params)
	default:
		return fmt.Errorf("unknown SMS provider: %s", d.Provider)
	}
}

func (d *SmsDriver) sendAliyun(_ context.Context, phone, templateCode string, params map[string]string) error {
	// Production: use https://dysmsapi.aliyuncs.com with HMAC-SHA1 signature
	// Simplified skeleton — real impl needs proper HMAC signing
	paramsJSON, _ := json.Marshal(params)
	fmt.Printf("[SMS Aliyun] phone=%s tpl=%s sign=%s params=%s\n",
		phone, templateCode, d.SignName, paramsJSON)
	_ = http.DefaultClient // placeholder
	return nil
}

func (d *SmsDriver) sendTencent(_ context.Context, phone, templateCode string, params map[string]string) error {
	// Production: use https://sms.tencentcloudapi.com with TC3-HMAC-SHA256 signature
	vals := make([]string, 0, len(params))
	for _, v := range params { vals = append(vals, v) }
	fmt.Printf("[SMS Tencent] phone=%s tpl=%s params=%v\n",
		phone, templateCode, strings.Join(vals, ","))
	return nil
}
