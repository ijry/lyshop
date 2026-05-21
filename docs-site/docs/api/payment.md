# 支付接口

## 说明

支付模块用于微信支付、支付宝等第三方支付能力接入。

## 典型接口

- `POST /api/pay/create`
- `POST /api/pay/callback`
- `GET /api/pay/status`

## 说明

- 回调接口必须做签名校验
- 支付通知处理应保持幂等
