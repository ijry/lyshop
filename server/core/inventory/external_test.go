package inventory

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ijry/lyshop/config"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestExternalProviderReserveSyncCallsRemoteAPI(t *testing.T) {
	originalClient := externalHTTPClient
	externalHTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, "POST", req.Method)
		require.Equal(t, "/reserve", req.URL.Path)
		body, _ := io.ReadAll(req.Body)
		require.Contains(t, string(body), `"biz_no":"ORD-9"`)
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"code":0}`)),
			Header:     make(http.Header),
		}, nil
	})}
	t.Cleanup(func() { externalHTTPClient = originalClient })

	original := config.Global
	t.Cleanup(func() { config.Global = original })
	config.Global.Inventory.Provider = "external_wms"
	config.Global.Inventory.ExternalMode = "sync"
	config.Global.ExternalWMS.Endpoint = "https://wms.example.com"

	p := &externalProvider{}
	err := p.ReserveTx(&gorm.DB{}, ReserveInput{
		BizType: "order",
		BizNo:   "ORD-9",
		Items:   []ReserveItem{{SkuID: 1, Qty: 2}},
	})
	require.NoError(t, err)
}

func TestBuildExternalSignature(t *testing.T) {
	req := signedPayload{
		AppKey:    "demo-key",
		Timestamp: "1717910400",
		Nonce:     "nonce-1",
		Body:      `{"biz_no":"O1001"}`,
	}

	sign := BuildSignature(req, "demo-secret")
	require.Equal(t, "e7593e961fe2c512caee1919359e5eee2c85003beded98875e35dd203fc8b625", sign)
}

func TestVerifyCallbackSignatureRejectsInvalidSign(t *testing.T) {
	err := VerifyCallbackSignature(callbackEnvelope{
		AppKey:    "demo-key",
		Timestamp: "1717910400",
		Nonce:     "nonce-1",
		Sign:      "bad-sign",
		Body:      `{"request_id":"REQ-1"}`,
	}, "demo-secret", time.Unix(1717910400, 0))
	require.ErrorContains(t, err, "signature")
}

func TestExternalProviderGetSellableStock(t *testing.T) {
	originalClient := externalHTTPClient
	externalHTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, "POST", req.Method)
		require.Equal(t, "/stock/sellable", req.URL.Path)
		return &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(bytes.NewBufferString(`{"code":0,"data":[{"sku_id":11,"sellable_stock":7},{"sku_id":12,"sellable_stock":0}]}`)),
			Header: make(http.Header),
		}, nil
	})}
	t.Cleanup(func() { externalHTTPClient = originalClient })

	original := config.Global
	t.Cleanup(func() { config.Global = original })
	config.Global.ExternalWMS.Endpoint = "https://wms.example.com"

	provider := &externalProvider{}
	stocks, err := provider.GetSellableStock(context.Background(), []uint64{11, 12})
	require.NoError(t, err)
	require.Len(t, stocks, 2)
	require.Equal(t, uint64(11), stocks[0].SkuID)
	require.Equal(t, 7, stocks[0].SellableStock)
}

func TestExternalProviderGetSellableStockMapsRemoteError(t *testing.T) {
	originalClient := externalHTTPClient
	externalHTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(bytes.NewBufferString(`{"code":5001,"msg":"remote error"}`)),
			Header: make(http.Header),
		}, nil
	})}
	t.Cleanup(func() { externalHTTPClient = originalClient })

	original := config.Global
	t.Cleanup(func() { config.Global = original })
	config.Global.ExternalWMS.Endpoint = "https://wms.example.com"

	provider := &externalProvider{}
	_, err := provider.GetSellableStock(context.Background(), []uint64{11})
	require.ErrorContains(t, err, "remote error")
}
