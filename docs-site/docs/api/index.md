# 接口文档

当前章节提供 LYShop API 使用约定与模块索引。

## 基础约定

- 基础路径：`/api`
- 请求格式：`application/json`
- 响应格式：统一 JSON 包装
- 认证方式：`Authorization: Bearer <token>`

## 统一响应结构

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

## 模块索引

- [认证接口](./auth)
- [后台概览接口](./admin)
- [商品接口](./product)
- [订单接口](./order)
- [仓储接口](./wms)
- [库存预占规则](./stock-reservation)
- [外部 WMS 通用协议说明](./external-wms-generic)
- [外部 WMS 对接交付说明](./external-wms-handover)
- [营销接口](./marketing)
- [IM 接口](./im)
- [支付接口](./payment)
- [装修接口](./decor)
- [会员接口](./vip)

## 说明

各模块页面给出当前约定、典型接口和扩展示例，后续可按实际路由继续补全。

库存相关文档建议按以下顺序阅读：

1. [商品接口](./product)
   - 理解 SPU / SKU 与商品展示模型
2. [订单接口](./order)
   - 理解订单库存状态与库存交易入口
3. [库存预占规则](./stock-reservation)
   - 理解统一 `inventory` 架构与 provider 模式
4. [仓储接口](./wms)
   - 理解内置 `builtin_wms` provider 能力边界
5. [外部 WMS 通用协议说明](./external-wms-generic)
   - 理解 `external_wms` 默认企业协议
