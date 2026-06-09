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
)

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

type externalCallbackPayload struct {
	RequestID  string `json:"request_id"`
	CallbackID string `json:"callback_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

type genericExternalResponse struct {
	Code      int             `json:"code"`
	Message   string          `json:"message"`
	Msg       string          `json:"msg"`
	RequestID string          `json:"request_id"`
	Data      json.RawMessage `json:"data"`
}

type externalWMSError struct {
	Code    int
	Message string
}

func (e *externalWMSError) Error() string {
	return fmt.Sprintf("external wms error code=%d message=%s", e.Code, e.Message)
}

type ExternalAdapter interface {
	BuildSignedRequest(ctx context.Context, method, path string, payload any) (*http.Request, error)
	VerifyCallback(req callbackEnvelope, now time.Time) error
	ParseCallbackBody(body string) (*externalCallbackPayload, error)
	ParseResponse(resp *http.Response) (*genericExternalResponse, error)
	IsSuccess(code int) bool
	IsRetryable(code int) bool
	ErrorMessage(resp *genericExternalResponse) string
}

type genericExternalAdapter struct{}

func NewExternalAdapter() ExternalAdapter {
	return &genericExternalAdapter{}
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

func (a *genericExternalAdapter) BuildSignedRequest(ctx context.Context, method, path string, payload any) (*http.Request, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
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
	req.Header.Set("X-Request-Mode", "generic")
	return req, nil
}

func (a *genericExternalAdapter) VerifyCallback(req callbackEnvelope, now time.Time) error {
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
	}, config.Global.ExternalWMS.AppSecret)
	if !strings.EqualFold(expected, req.Sign) {
		return fmt.Errorf("signature mismatch")
	}
	return nil
}

func (a *genericExternalAdapter) ParseCallbackBody(body string) (*externalCallbackPayload, error) {
	var payload externalCallbackPayload
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return nil, err
	}
	if strings.TrimSpace(payload.RequestID) == "" {
		return nil, fmt.Errorf("callback request_id is required")
	}
	return &payload, nil
}

func (a *genericExternalAdapter) ParseResponse(resp *http.Response) (*genericExternalResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("external wms http status %d", resp.StatusCode)
	}
	var out genericExternalResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (a *genericExternalAdapter) IsSuccess(code int) bool {
	return code == 0
}

func (a *genericExternalAdapter) IsRetryable(code int) bool {
	switch code {
	case 1001, 1002, 1003, 2001, 5000, 5001, 5002:
		return true
	default:
		return false
	}
}

func (a *genericExternalAdapter) ErrorMessage(resp *genericExternalResponse) string {
	if strings.TrimSpace(resp.Message) != "" {
		return resp.Message
	}
	return resp.Msg
}

func VerifyCallbackSignature(req callbackEnvelope, secret string, now time.Time) error {
	adapter := NewExternalAdapter()
	return adapter.VerifyCallback(req, now)
}

func absDuration(v time.Duration) time.Duration {
	if v < 0 {
		return -v
	}
	return v
}
