package inventory

import (
	"bytes"
	"io"
	"net/http"
	"testing"

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
