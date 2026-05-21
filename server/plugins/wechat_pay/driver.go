package wechat_pay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ijry/lyshop/core/driver/payment"
)

// WechatPayDriver implements payment.Driver using WeChat Pay v3 API.
// Full production implementation requires wechatpay-go SDK.
// This skeleton wires the interface and config loading.
type WechatPayDriver struct {
	AppID  string
	MchID  string
	APIKey string // v3 API key
	Serial string // certificate serial number
}

func (d *WechatPayDriver) Name() string { return "wechat" }

func (d *WechatPayDriver) CreateOrder(ctx context.Context, p *payment.OrderParams) (*payment.OrderResult, error) {
	if d.AppID == "" {
		return nil, errors.New("微信支付未配置，请在系统设置中填写 AppID/MchID/APIKey")
	}
	// In production: use github.com/wechatpay-apiv3/wechatpay-go to call
	// POST https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi (JSAPI)
	// POST https://api.mch.weixin.qq.com/v3/pay/transactions/h5    (H5)
	// POST https://api.mch.weixin.qq.com/v3/pay/transactions/app   (App)
	return &payment.OrderResult{
		PrepayID: fmt.Sprintf("mock_prepay_%s", p.OrderNo),
		PayParams: map[string]string{
			"appId":     d.AppID,
			"timeStamp": "mock",
			"nonceStr":  "mock",
			"package":   "prepay_id=mock_prepay_" + p.OrderNo,
			"signType":  "RSA",
			"paySign":   "mock_sign",
		},
	}, nil
}

func (d *WechatPayDriver) QueryOrder(ctx context.Context, tradeNo string) (*payment.QueryResult, error) {
	return &payment.QueryResult{TradeNo: tradeNo, Status: "paid"}, nil
}

func (d *WechatPayDriver) Refund(ctx context.Context, p *payment.RefundParams) (*payment.RefundResult, error) {
	return &payment.RefundResult{RefundID: "mock_refund_" + p.RefundNo}, nil
}

func (d *WechatPayDriver) HandleNotify(ctx context.Context, r *http.Request) (*payment.NotifyResult, error) {
	body, _ := io.ReadAll(r.Body)
	var payload map[string]any
	json.Unmarshal(body, &payload)
	orderNo, _ := payload["out_trade_no"].(string)
	return &payment.NotifyResult{OrderNo: orderNo, Paid: true}, nil
}
