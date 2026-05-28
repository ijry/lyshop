# 会员接口

## 功能说明

会员模块已落地后端最小可用能力，包含：

- 会员档案与开通
- 会员等级与成长值日志
- 会员月领券（自然月配额控制）
- 会员 SKU 价格配置与下单会员价计算（后台在商品管理内按 SKU 配置）

会员价计算器已接入营销管线，优先级为 `15`，且活动价商品（`item.ActivityPrice > 0`）会跳过会员价逻辑。

## 接口变化

### 管理后台（Admin）

- `GET /admin/api/vip/plans`：会员套餐列表
- `GET /admin/api/vip/levels`：会员等级列表
- `GET /admin/api/vip/coupon-rules`：会员月领券规则列表
- `GET /admin/api/vip/sku-prices`：会员 SKU 价列表（支持 `product_id/sku_id/level_id/status/page/size` 过滤）
- `POST /admin/api/vip/sku-prices`：新增会员 SKU 价
- `PUT /admin/api/vip/sku-prices/:id`：更新会员 SKU 价
- `DELETE /admin/api/vip/sku-prices/:id`：删除会员 SKU 价

### C 端（需登录）

- `GET /api/v1/vip/profile`：获取会员档案
- `POST /api/v1/vip/open`：开通会员（入参：`plan_id`）
- `GET /api/v1/vip/coupons/monthly`：查询会员月领券规则与本月领取状态
- `POST /api/v1/vip/coupons/monthly/:rule_id/claim`：领取指定规则会员券
- `GET /api/v1/vip/growth/logs`：分页查询成长值变动日志

## 部署或配置影响

- 需在 `plugins.enabled` 中启用 `vip` 插件。
- 无新增环境变量。
- 插件迁移会新增会员相关数据表（会员档案、等级、套餐、月领券规则、会员价、成长值日志等）。
- 后台菜单不再提供独立“会员 SKU 价”入口，改为在商品编辑页（SKU 区）按会员等级配置，部署无需新增配置。
