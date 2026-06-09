# 外部 WMS 通用协议说明

## 说明

LYShop 在 `inventory.provider=external_wms` 模式下，默认内置一个 `generic adapter`，用于统一企业外部 WMS 的基础对接协议。

这份协议文档描述的是当前系统默认支持的“通用企业协议”，目标是先统一商城与企业 WMS 的集成边界：

- 统一请求头
- 统一签名方式
- 统一同步响应结构
- 统一异步回调结构
- 统一任务重试与状态语义

如果后续要接入某个特定厂商 WMS，建议在此基础上新增独立 adapter，而不是继续把厂商差异直接写进业务逻辑。

## 适用范围

- 同步模式：`inventory.external_mode=sync`
- 异步模式：`inventory.external_mode=async`

统一库存三段式交易动作：

1. `reserve`
2. `confirm`
3. `release`

以及扩展动作：

- `deduct`
- `restore`
- `sync_sku`
- `stock/sellable`

## 配置项

至少需要配置：

```yaml
inventory:
  provider: external_wms
  external_mode: sync

external_wms:
  endpoint: "https://wms.example.com"
  app_key: "demo-key"
  app_secret: "demo-secret"
  timeout_ms: 3000
  callback_enabled: true
  signature_ttl: 300
  worker_interval_sec: 5
  retry:
    max_attempts: 8
    backoff_seconds: 30
```

字段说明：

- `endpoint`
  - 外部 WMS 根地址
- `app_key`
  - 调用方身份标识
- `app_secret`
  - 请求签名与回调验签密钥
- `signature_ttl`
  - 回调签名有效期，单位秒
- `worker_interval_sec`
  - 异步 worker 轮询频率，单位秒

## 请求协议

### 请求头

所有请求统一携带以下请求头：

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
- 输出为十六进制小写字符串

### 幂等约定

请求体可携带：

- `biz_type`
- `biz_no`
- `request_id`

推荐外部 WMS 以 `biz_type + biz_no + action` 或 `request_id` 作为幂等键，避免重复执行库存动作。

## 动作与路径

当前默认动作与路径约定：

- `POST /reserve`
- `POST /confirm`
- `POST /release`
- `POST /deduct`
- `POST /restore`
- `POST /sync-sku`
- `POST /stock/sellable`

## 请求体示例

### reserve

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

### confirm

```json
{
  "biz_type": "order",
  "biz_no": "ORD-1001"
}
```

### release

```json
{
  "biz_type": "order",
  "biz_no": "ORD-1001",
  "reason": "user_cancel"
}
```

### deduct

```json
{
  "biz_type": "order",
  "biz_no": "ORD-1001",
  "items": [
    { "sku_id": 101, "qty": 2 }
  ]
}
```

### restore

```json
{
  "biz_type": "order",
  "biz_no": "ORD-1001",
  "reason": "refund",
  "items": [
    { "sku_id": 101, "qty": 2 }
  ]
}
```

### sync-sku

```json
{
  "biz_type": "product",
  "biz_no": "sku:101",
  "sku_id": 101,
  "stock": 80
}
```

### stock/sellable

```json
{
  "sku_ids": [101, 102, 103]
}
```

## 同步响应协议

统一响应体：

```json
{
  "code": 0,
  "message": "ok",
  "request_id": "WMS-REQ-20260609-0001",
  "data": {}
}
```

兼容字段：

- `message`
- `msg`

成功判定：

- `code = 0`

失败判定：

- `code != 0`

其中部分错误码会被系统视为可重试错误。

### 建议错误码

- `0`
  - 成功
- `1001`
  - 外部系统繁忙，可重试
- `1002`
  - 网络或依赖超时，可重试
- `1003`
  - 锁冲突，可重试
- `2001`
  - 下游处理中，可重试
- `4001`
  - 参数错误，不重试
- `4002`
  - 签名错误，不重试
- `4003`
  - 库存不足，不重试
- `4004`
  - 业务失败，不重试
- `5000`
  - 外部服务异常，可重试
- `5001`
  - 外部网关异常，可重试
- `5002`
  - 外部临时不可用，可重试

## 异步回调协议

异步模式下，外部 WMS 可以通过：

- `POST /admin/api/external-wms/callback`

回调商城。

### 外层请求体

```json
{
  "app_key": "demo-key",
  "timestamp": "1717910400",
  "nonce": "nonce-1",
  "sign": "abcdef123456",
  "body": "{\"request_id\":\"WMS-REQ-20260609-0001\",\"callback_id\":\"CALLBACK-1001\",\"status\":\"success\",\"message\":\"ok\"}"
}
```

### body 内层结构

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
  - 外部 WMS 请求流水号，应与同步响应或异步受理返回值对应
- `callback_id`
  - 回调唯一标识，用于幂等去重
- `status`
  - `success | failed | processing`
- `message`
  - 补充说明或失败原因

### 回调验签

回调验签规则与请求签名一致：

```text
app_key + "\n" + timestamp + "\n" + nonce + "\n" + body + "\n" + app_secret
```

当满足以下任一条件时，商城会拒绝回调：

- 签名不匹配
- `timestamp` 超过 `external_wms.signature_ttl`
- `body` 无法解析

## 异步任务与订单状态

异步模式下，库存动作先写入 `inventory_integration_tasks`，再由后台 worker 执行。

任务状态：

- `pending`
- `processing`
- `success`
- `failed`

订单 `inventory_status` 收口规则：

- `reserve` 成功 -> `reserved`
- `confirm` / `deduct` 成功 -> `confirmed`
- `release` / `restore` 成功 -> `released`
- 任意动作失败 -> `failed`

在异步执行完成前，订单库存状态保持：

- `pending`

## 重试与幂等

系统会对可重试错误码进行自动重试，具体受以下配置控制：

- `external_wms.retry.max_attempts`
- `external_wms.retry.backoff_seconds`

幂等规则：

- 同一 `callback_id` 的重复回调会被忽略
- 已进入最终态（`success / failed`）的任务不会被后续回调覆盖
- 建议外部 WMS 也对 `request_id` 做去重

## 对接建议

- 先按本文档实现一个企业侧中间适配层，而不是直接把企业 WMS 原始字段暴露给商城
- 如果企业现有 WMS 无法直接满足本协议，建议企业侧做一次协议转换
- 当需要接入多家 WMS 或某家 WMS 有特殊签名算法时，再新增专用 adapter

## 关联文档

- [库存预占交易规则](./stock-reservation)
- [仓储接口](./wms)
- [订单接口](./order)
