# 积分商城插件 - README

## 简介

积分商城插件是一个完整的积分管理和兑换系统，支持用户通过签到、消费等方式获取积分，并使用积分兑换商品或抵扣订单金额。

## 功能特性

### 核心功能
- ✅ 积分获取（签到、订单完成、管理员调整）
- ✅ 积分消耗（商品兑换、订单抵扣）
- ✅ 积分商品管理（优惠券、实物、虚拟）
- ✅ 兑换流程管理（待发货、已发货、已完成）
- ✅ 积分日志记录
- ✅ 积分统计分析

### 商品类型
1. **优惠券** - 兑换后自动发放到用户账户
2. **实物商品** - 需要填写收货地址，支持物流跟踪
3. **虚拟商品** - 兑换后立即显示虚拟内容（如兑换码）

### 多端支持
- **Admin 端** - 商品管理、兑换记录、积分日志、统计分析
- **Eapp 端** - 商家端商品管理
- **App 端** - 用户端商品浏览和兑换
- **Web 端** - 用户端 Web 版

## 快速开始

### 1. 启用插件

编辑 `server/config.yaml`：
```yaml
plugins:
  enabled:
    - points_mall
```

### 2. 启动服务

```bash
cd server
go run main.go -config config.yaml
```

### 3. 访问管理后台

访问 Admin 端，导航到"积分商城"菜单。

## 技术架构

### 后端
- **语言**: Go
- **框架**: Gin + GORM
- **架构**: 插件化设计

### 前端
- **Admin**: Vue 3 + Vite + UnoCSS
- **Eapp/App**: uni-app + Vue 3
- **Web**: Vue 3 + Vite

## 数据模型

### PointsLog（积分日志）
记录所有积分变动，支持 7 种类型：
1. 签到
2. 订单抵扣
3. 兑换消耗
4. 订单完成
5. 管理员调整
6. 过期扣除
7. 活动奖励

### PointsProduct（积分商品）
支持三种商品类型，包含库存管理、兑换限制等功能。

### PointsExchange（兑换记录）
完整的兑换流程管理，支持发货、确认收货、取消等操作。

## API 文档

### 用户端 API

```
GET  /api/v1/points/products          # 商品列表
GET  /api/v1/points/products/:id      # 商品详情
POST /api/v1/points/products/:id/exchange  # 兑换商品
GET  /api/v1/points/exchanges         # 兑换记录
GET  /api/v1/points/logs              # 积分日志
GET  /api/v1/points/balance           # 积分余额
```

### 管理端 API

```
GET    /admin/api/points/products     # 商品列表
POST   /admin/api/points/products     # 创建商品
PUT    /admin/api/points/products/:id # 更新商品
DELETE /admin/api/points/products/:id # 删除商品
GET    /admin/api/points/exchanges    # 兑换记录
PUT    /admin/api/points/exchanges/:id/ship  # 发货
GET    /admin/api/points/logs         # 积分日志
POST   /admin/api/points/adjust       # 调整积分
GET    /admin/api/points/stats        # 积分统计
```

## 集成指南

### 与 Checkin 插件集成

签到插件已自动集成，签到时会调用积分服务：
```go
pmservice.AddPoints(ctx, userID, points, 1, "每日签到")
```

### 与 Order 插件集成

在订单完成时赠送积分：
```go
import pmservice "github.com/ijry/lyshop/server/plugins/points_mall/service"

// 订单状态变更为已完成时
if newStatus == ordermodel.OrderStatusCompleted {
    pmservice.GrantOrderPoints(ctx, orderID)
}
```

### 与 Marketing 插件集成

- **积分抵扣**: 通过计算器管道自动集成
- **优惠券发放**: 兑换优惠券类商品时自动发放

## 配置说明

### 积分兑换比例
默认：100 积分 = 1 元

### 订单完成赠送
- 开关：可配置
- 比例：默认 1%（消费 100 元送 100 积分）

### 积分过期
- 天数：可配置（0 = 永不过期）

## 开发指南

### 目录结构
```
server/plugins/points_mall/
├── plugin.json          # 插件元数据
├── plugin.go            # 插件入口
├── model/               # 数据模型
├── service/             # 业务逻辑
├── calculator/          # 积分计算器
└── api/                 # API 路由
```

### 添加新功能

1. 在 `service/` 中添加业务逻辑
2. 在 `api/` 中添加 API 端点
3. 在前端添加对应页面
4. 更新文档

## 测试

### 运行单元测试
```bash
cd server/plugins/points_mall/service
go test -v
```

### 测试覆盖率
```bash
go test -cover
```

## 文档

- [完整实现文档](./points_mall_implementation.md)
- [快速启动指南](./points_mall_quickstart.md)

## 许可证

与主项目保持一致

## 贡献

欢迎提交 Issue 和 Pull Request！

## 更新日志

### v1.0.0 (2026-05-31)
- ✅ 初始版本发布
- ✅ 完整的积分管理功能
- ✅ 三种商品类型支持
- ✅ 四端完整实现
- ✅ 与 checkin、marketing、order 插件集成
- ✅ 国际化支持（中英文）
- ✅ 单元测试框架
