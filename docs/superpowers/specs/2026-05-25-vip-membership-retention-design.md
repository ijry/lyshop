# LYShop 会员体系（复购留存优先）设计

- 日期：2026-05-25
- 目标：建设可持续复购的会员体系，支持年费会员、月度自主领券、会员价、成长值等级。
- 约束确认：
  - 会员券为“自主领取”。
  - 每月未领取券到月底作废。
  - 会员价仅对普通商品生效，参与秒杀/拼团等活动商品不享会员价。

## 1. 设计目标

1. 在不破坏现有营销能力的前提下，新增独立 `vip` 插件。
2. 保持“营销总线”作为统一结算入口，避免未来重复造轮子。
3. 优先落地能提升复购的能力：月领券、等级成长、会员价。

## 2. 架构决策

### 2.1 插件边界

- `marketing` 插件：继续承载优惠券、活动、积分、分销等通用营销能力。
- `vip` 插件：承载会员生命周期、等级成长、会员权益规则。
- 两者通过“价格总线 + 券发放接口复用”协作，不在 `vip` 内重复实现券引擎。

### 2.2 为什么不一次性拆秒杀/拼团

当前代码中秒杀/拼团已稳定运行在营销插件。一次性拆分会引入较高迁移和回归成本，且对当前目标（复购留存）收益有限。

策略：
- 本期先做“可插拔抽象层”和 `vip` 插件接入；
- 秒杀/拼团后续在复杂度达到阈值时再独立拆分。

## 3. 与现有代码的衔接点

- 价格流水线：`server/core/marketing/pipeline.go`
- 价格上下文：`server/core/marketing/context.go`
- 下单入口：`server/plugins/order/service/order.go`
- 活动计算器：`server/plugins/marketing/calculator/activity.go`
- 满减计算器：`server/plugins/marketing/calculator/full_reduce.go`
- 优惠券计算器：`server/plugins/marketing/calculator/coupon.go`

## 4. 结算优先级（本期固定）

1. 活动价（秒杀/拼团等）
2. 会员价（仅普通商品）
3. 满减
4. 优惠券
5. 积分
6. 分销佣金（不影响应付）

说明：
- 行项目若命中活动价（`ActivityPrice > 0`），跳过会员价计算。
- 会员价仅作用于未命中活动价的行项目。

## 5. 数据模型（vip 插件）

### 5.1 会员计划与等级

1. `vip_plans`
   - `id,name,duration_months,price,renew_price,status,created_at,updated_at`
2. `vip_levels`
   - `id,name,growth_threshold,benefit_json,status,sort,created_at,updated_at`

### 5.2 用户会员资产与成长值

3. `vip_user_assets`
   - `user_id,current_plan_id,current_level_id,vip_start_at,vip_end_at,growth_value,status,updated_at`
4. `vip_growth_logs`
   - `id,user_id,order_id,event_type,growth_delta,balance_after,idempotency_key,remark,created_at`
   - 唯一约束：`idempotency_key`（防重复发放/回退）

### 5.3 会员券规则与领取

5. `vip_coupon_rules`
   - `id,plan_id,level_id,coupon_id,monthly_quota,claim_mode,status,created_at,updated_at`
   - `claim_mode` 本期固定 `manual`
6. `vip_coupon_claims`
   - `id,user_id,rule_id,period_yyyymm,claimed_count,last_claimed_at,created_at,updated_at`
   - 唯一约束：`(user_id,rule_id,period_yyyymm)`

### 5.4 会员价

7. `vip_sku_prices`
   - `id,product_id,sku_id,level_id,vip_price,vip_discount_rate,status,created_at,updated_at`
   - `vip_price` 与 `vip_discount_rate` 二选一生效。

### 5.5 订单权益快照（审计/售后回退）

8. `vip_order_benefits`
   - `id,order_id,user_id,vip_discount,growth_granted,growth_reverted,created_at,updated_at`

## 6. API 设计

### 6.1 Admin

- `GET /admin/vip/plans`
- `POST /admin/vip/plans`
- `GET /admin/vip/levels`
- `POST /admin/vip/levels`
- `GET /admin/vip/coupon-rules`
- `POST /admin/vip/coupon-rules`
- `GET /admin/vip/sku-prices`
- `POST /admin/vip/sku-prices`
- `GET /admin/vip/users`

### 6.2 Front

- `GET /api/v1/vip/profile`
- `POST /api/v1/vip/open`（支持按年开通）
- `GET /api/v1/vip/coupons/monthly`
- `POST /api/v1/vip/coupons/monthly/:rule_id/claim`
- `GET /api/v1/vip/growth/logs`

## 7. 核心业务规则

1. 月度会员券：
   - 用户手动领取；
   - 每月按自然月窗口；
   - 未领取额度不滚存，跨月作废；
   - 领取成功后写入 `coupon_users`。
2. 会员价：
   - 活动商品不享会员价；
   - 普通商品按用户当前等级匹配会员价；
   - 行项目会员优惠需写入结算规则明细。
3. 成长值：
   - 订单支付成功发放成长值；
   - 售后退款按规则回退成长值；
   - 全流程幂等，防重复发放/回退。

## 8. 开发分期

### Phase 1（MVP，优先复购）

- `vip` 插件骨架 + 8 张表迁移
- 月领券规则与前台领取接口
- 会员价计算器接入总线
- 用户会员中心基础信息接口

### Phase 2（稳定增长）

- 年费续费、到期任务、降级
- 成长值支付发放与退款回退
- 后台会员用户查询与审计视图

### Phase 3（抽象升级）

- 评估将秒杀/拼团独立插件化
- 营销总线扩展为标准 Provider 注册机制

## 9. 风险与控制

1. 风险：价格规则叠加误配导致客诉。
   - 控制：固化优先级，后台展示“命中规则明细”。
2. 风险：月领券并发重复领取。
   - 控制：唯一索引 + 事务锁 + 幂等键。
3. 风险：成长值重复发放。
   - 控制：`idempotency_key` 与支付状态校验。

## 10. 验收标准

1. 会员用户可在会员中心看到“本月可领券”并成功领取一次。
2. 非活动商品在下单时命中会员价；活动商品不命中会员价。
3. 订单支付后产生成长值日志；退款后可回退成长值且不重复。
4. 后台可配置会员计划、等级、月领券规则、会员价。
