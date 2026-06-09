# 外部 WMS 对接交付说明

## 目的

这份文档面向企业 WMS 对接方，用于快速完成 LYShop 与外部 WMS 的联调。

如果需要完整字段约定、状态机和错误码建议，请再参考：

- [外部 WMS 通用协议说明](./external-wms-generic)

## 对接模式

LYShop 支持两种外部 WMS 对接模式：

1. `sync`
   - 商城下单/确认/释放时，实时调用企业 WMS
2. `async`
   - 商城先落库存任务，再由后台 worker 异步调用企业 WMS
   - 企业 WMS 处理完成后，再回调商城

## 商城需配置

```yaml
inventory:
  provider: external_wms
  external_mode: sync

external_wms:
  endpoint: "https://wms.example.com"
  app_key: "demo-key"
  app_secret: "demo-secret"
  signature_ttl: 300
```

## 商城调用企业 WMS

### 请求头

所有请求统一携带：

- `Content-Type: application/json`
- `X-App-Key`
- `X-Timestamp`
- `X-Nonce`
- `X-Sign`
- `X-Request-Mode: generic`

### 签名规则

签名原文：

```text
app_key + "\n" + timestamp + "\n" + nonce + "\n" + body + "\n" + app_secret
```

签名算法：

- `sha256`
- 十六进制小写输出

### 请求路径

企业 WMS 需提供以下接口：

- `POST /reserve`
- `POST /confirm`
- `POST /release`
- `POST /deduct`
- `POST /restore`
- `POST /sync-sku`
- `POST /stock/sellable`

## 核心请求示例

### 1. 预占库存 `reserve`

```json
{
  "biz_type": "order",
  "biz_no": "ORD-1001",
  "items": [
    { "sku_id": 101, "qty": 2 },
    { "sku_id": 102, "qty": 1 }
  ]
}
```

### 2. 确认扣减 `confirm`

```json
{
  "biz_type": "order",
  "biz_no": "ORD-1001"
}
```

### 3. 释放库存 `release`

```json
{
  "biz_type": "order",
  "biz_no": "ORD-1001",
  "reason": "user_cancel"
}
```

### 4. 查询可售库存 `stock/sellable`

```json
{
  "sku_ids": [101, 102, 103]
}
```

## 企业 WMS 响应商城

统一响应体：

```json
{
  "code": 0,
  "message": "ok",
  "request_id": "WMS-REQ-20260609-0001",
  "data": {}
}
```

约定：

- `code=0` 表示成功
- `code!=0` 表示失败
- `request_id` 建议企业 WMS 每次返回，异步模式下商城会用它关联后续回调

## 企业 WMS 回调商城

异步模式下，企业 WMS 处理完成后回调：

- `POST /admin/api/external-wms/callback`

### 回调外层请求体

```json
{
  "app_key": "demo-key",
  "timestamp": "1717910400",
  "nonce": "nonce-1",
  "sign": "abcdef123456",
  "body": "{\"request_id\":\"WMS-REQ-20260609-0001\",\"callback_id\":\"CALLBACK-1001\",\"status\":\"success\",\"message\":\"ok\"}"
}
```

### `body` 内层 JSON

```json
{
  "request_id": "WMS-REQ-20260609-0001",
  "callback_id": "CALLBACK-1001",
  "status": "success",
  "message": "ok"
}
```

字段说明：

- `request_id`
  - 对应企业 WMS 返回给商城的请求流水号
- `callback_id`
  - 回调唯一号，商城会用它做幂等去重
- `status`
  - `success | failed | processing`
- `message`
  - 成功说明或失败原因

### 回调验签

回调签名规则与请求签名一致：

```text
app_key + "\n" + timestamp + "\n" + nonce + "\n" + body + "\n" + app_secret
```

## 最小联调范围

建议至少联调以下 5 项：

1. `reserve` 成功
2. `reserve` 失败
3. `confirm` 成功
4. `release` 成功
5. `stock/sellable` 查询成功

如果启用异步模式，再补：

1. 异步任务受理成功
2. 成功回调
3. 失败回调
4. 重复回调幂等验证

## 联调通过标准

- 商城能成功调用企业 WMS 预占、确认、释放接口
- 企业 WMS 能正确校验商城签名
- 商城能正确校验企业 WMS 回调签名
- 同一 `callback_id` 重复回调不会重复处理
- 可售库存查询结果能正确映射到商城 SKU 展示

## 备注

- 如果企业现有 WMS 原始协议与本文档不一致，建议企业侧增加一层中间适配服务
- 如果后续要接特定厂商协议，LYShop 可继续在当前 `generic adapter` 之上扩展专用 adapter
