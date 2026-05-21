package payment

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

type OrderParams struct {
	OrderNo     string
	Amount      int64  // cents (e.g. 9900 = ¥99.00)
	Description string
	NotifyURL   string
	OpenID      string // JSAPI (WeChat mini-program)
	ClientIP    string // H5 / App
}

type OrderResult struct {
	PrepayID  string
	PayParams map[string]string // passed directly to frontend SDK
}

type QueryResult struct {
	OutTradeNo string
	TradeNo    string // platform trade number
	Status     string // "paid" | "unpaid" | "refunded"
	Amount     int64
}

type RefundParams struct {
	OrderNo     string
	RefundNo    string
	Amount      int64
	TotalAmount int64
	Reason      string
}

type RefundResult struct {
	RefundID string
}

type NotifyResult struct {
	OrderNo string
	Amount  int64
	Paid    bool
}

// Driver is the interface all payment plugins must implement.
type Driver interface {
	Name() string
	CreateOrder(ctx context.Context, p *OrderParams) (*OrderResult, error)
	QueryOrder(ctx context.Context, tradeNo string) (*QueryResult, error)
	Refund(ctx context.Context, p *RefundParams) (*RefundResult, error)
	HandleNotify(ctx context.Context, r *http.Request) (*NotifyResult, error)
}

var (
	mu      sync.RWMutex
	drivers = map[string]Driver{}
)

func Register(d Driver) { mu.Lock(); drivers[d.Name()] = d; mu.Unlock() }

func Get(name string) (Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	d, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("payment driver %q not registered", name)
	}
	return d, nil
}

// Names returns the names of all registered payment drivers.
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(drivers))
	for n := range drivers {
		out = append(out, n)
	}
	return out
}
