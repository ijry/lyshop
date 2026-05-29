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
- [营销接口](./marketing)
- [IM 接口](./im)
- [支付接口](./payment)
- [装修接口](./decor)
- [会员接口](./vip)

## 说明

各模块页面给出当前约定、典型接口和扩展示例，后续可按实际路由继续补全。
