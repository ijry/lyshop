package alipay

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ijry/lyshop/core/driver/payment"
)

// AlipayDriver implements payment.Driver using Alipay Open API.
// Full production implementation requires go-alipay-sdk.
// This skeleton wires the interface and config loading.
type AlipayDriver struct {
	AppID      string
	PrivateKey string
	PublicKey  string // Alipay public key for signature verification
	Sandbox    bool
}

func (d *AlipayDriver) Name() string { return "alipay" }

func (d *AlipayDriver) CreateOrder(ctx context.Context, p *payment.OrderParams) (*payment.OrderResult, error) {
	if d.AppID == "" {
		return nil, errors.New("支付宝未配置，请在系统设置中填写 AppID/PrivateKey/PublicKey")
	}
	// In production: build alipay.trade.app.pay / alipay.trade.wap.pay request
	// and sign with RSA2 private key
	return &payment.OrderResult{
		PrepayID: fmt.Sprintf("mock_alipay_%s", p.OrderNo),
		PayParams: map[string]string{
			"orderStr": "mock_order_str_for_" + p.OrderNo,
		},
	}, nil
}

func (d *AlipayDriver) QueryOrder(ctx context.Context, tradeNo string) (*payment.QueryResult, error) {
	return &payment.QueryResult{TradeNo: tradeNo, Status: "paid"}, nil
}

func (d *AlipayDriver) Refund(ctx context.Context, p *payment.RefundParams) (*payment.RefundResult, error) {
	return &payment.RefundResult{RefundID: "mock_alipay_refund_" + p.RefundNo}, nil
}

func (d *AlipayDriver) HandleNotify(ctx context.Context, r *http.Request) (*payment.NotifyResult, error) {
	r.ParseForm()
	orderNo := r.FormValue("out_trade_no")
	tradeStatus := r.FormValue("trade_status")
	paid := tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED"
	io.Discard.Write(nil)
	return &payment.NotifyResult{OrderNo: orderNo, Paid: paid}, nil
}
