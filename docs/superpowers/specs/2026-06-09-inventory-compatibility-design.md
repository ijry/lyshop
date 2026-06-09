# 库存三态兼容架构设计

## 1. 背景与目标

当前商城库存链路已经具备内置 `wms` 插件，但业务依赖仍然是强绑定关系：

- `order` 插件直接依赖 `wms`
- 下单、支付、取消直接调用 `wms` 预占/确认/释放服务
- 商品创建与新增 SKU 时会直接同步 `wms` 库存行
- 文档默认将 `wms` 视为唯一库存真源

这会带来两个直接问题：

1. `wms` 无法默认关闭，轻量商城部署成本偏高。
2. 商城无法在不启用内置 `wms` 的前提下对接企业已有 WMS 系统。

本次设计目标是建立一套统一库存兼容架构，同时满足以下三种运行模式：

1. `local`：不启用 `wms`，商城使用本地商品库存完成交易。
2. `builtin_wms`：启用内置 `wms` 插件，库存交易走仓储插件。
3. `external_wms`：不使用内置 `wms`，而是对接企业外部 WMS。

其中 `external_wms` 需同时支持两种集成模式：

1. `sync`：下单链路同步调用外部 WMS。
2. `async`：商城落库后异步推送外部 WMS，并提供补偿重试能力。

---

## 2. 设计原则

1. **优先兼容升级现有接口**：订单、支付、取消、商品保存接口尽量不改路径，只升级内部库存实现与返回字段。
2. **库存能力与 WMS 实现解耦**：先抽象库存领域，再把内置 WMS 视为其中一种 provider。
3. **默认轻量部署**：未启用 `wms` 时，商城仍应能完整下单与履约。
4. **外部集成可控**：同步模式追求强一致，异步模式追求系统韧性，并用显式状态与补偿机制兜底。
5. **插件边界清晰**：`wms` 插件负责内置仓储实现与后台仓储能力，不再承担整个平台库存抽象职责。
6. **文档默认写最新架构**：后续对外文档以统一库存架构为主叙述，而非继续把 `wms` 写成唯一真源。

---

## 3. 现状评估

当前代码中的关键耦合点如下：

1. `server/plugins/order/plugin.json`
   - 直接声明依赖 `["product", "wms"]`
2. `server/plugins/order/service/order.go`
   - 直接引用 `wms/service` 的 `ReserveStockTx`、`ConfirmReservationTx`、`ReleaseReservationTx`
3. `server/plugins/product/service/product.go`
   - 创建 SKU 后直接尝试初始化 WMS 库存
4. `docs-site/docs/api/order.md` 与 `docs-site/docs/api/stock-reservation.md`
   - 直接声明 `wms` 为订单库存真源

这说明本次不是简单做插件启停，而是要完成“库存域抽象 + provider 路由 + 交易状态重塑”。

---

## 4. 方案比较

### 4.1 方案 A：业务代码内部分支切换

做法：

- 在 `order/product/points_mall` 里根据配置直接写 `if provider == ...`
- `local` 走商品库存
- `builtin_wms` 走内置 WMS
- `external_wms` 走外部 HTTP 调用

优点：

- 改动快
- 首轮开发量最小

缺点：

- 库存逻辑分散在多个业务模块
- 后续增加第二家外部 WMS、异步补偿、回调重放时复杂度迅速失控
- 无法形成可复用的统一库存边界

结论：

- 不推荐

### 4.2 方案 B：统一库存抽象层

做法：

- 新增统一库存域 `core/inventory`
- 订单、商品、积分商城等模块都只依赖 `inventory` 接口
- `local`、`builtin_wms`、`external_wms` 各自实现 provider

优点：

- 架构边界清晰
- 可同时支持无 WMS、内置 WMS、外部 WMS
- 便于接入多家企业 WMS 与同步/异步双模式

缺点：

- 首轮改造面较大
- 需要迁移现有 `order/product` 直接依赖

结论：

- 推荐

### 4.3 方案 C：仅做插件路由器

做法：

- 参考 `storage_router/logistics_router/delivery_router`
- 只新增一个 `inventory_router` 插件，内部决定调哪个实现

优点：

- 风格接近现有插件体系

缺点：

- 如果没有先抽象库存领域接口，router 只是把复杂度挪位置
- `local` 模式本质不是插件，实现上更适合放在核心域能力里

结论：

- 可作为方案 B 的落地外壳，不适合作为单独方案

### 4.4 最终决策

采用 **方案 B + 方案 C 的外壳**：

1. 先在核心层定义统一库存接口与 provider 注册机制。
2. 再通过 inventory router 选择 `local / builtin_wms / external_wms`。
3. 让 `wms` 从“平台库存唯一真源”降级为“内置库存实现之一”。

---

## 5. 总体架构

### 5.1 目标运行模式

#### 模式一：`local`

- `wms` 插件默认不启用
- 库存真源为 `product_skus.stock`
- 下单、支付、取消的库存交易全部由本地库存 provider 完成

#### 模式二：`builtin_wms`

- 启用 `wms` 插件
- 库存真源为 `inventory_stock`
- 订单库存交易沿用当前 `wms` 预占/确认/释放模型

#### 模式三：`external_wms`

- 不要求启用 `wms`
- 商城库存交易委托企业 WMS
- 由外部 WMS provider 负责签名、调用、幂等、回调、重试

### 5.2 核心模块

建议新增或调整如下模块：

1. `server/core/inventory`
   - 定义统一库存接口、模型、错误、provider 注册表
2. `server/core/inventory/router`
   - 根据配置选择当前库存 provider
3. `server/core/inventory/providers/local`
   - 本地商品库存实现
4. `server/plugins/wms`
   - 继续承担内置仓储后台与库存实现
5. `server/plugins/external_wms` 或 `server/core/inventory/providers/external`
   - 外部 WMS 适配器与任务补偿能力

### 5.3 依赖方向

调整后的依赖应为：

- `order -> core/inventory`
- `product -> core/inventory`
- `points_mall -> core/inventory`
- `marketing -> core/inventory`（仅查询可售库存时）
- `wms -> product`
- `external_wms -> core/inventory`

必须移除的直接依赖：

- `order -> wms service`
- `product -> wms service`

同时，`order` 插件的依赖关系应从 `["product", "wms"]` 调整为只保留 `["product"]`。

---

## 6. 统一库存接口设计

### 6.1 Provider 接口

统一库存 provider 建议提供如下能力：

```go
type Provider interface {
    Name() string

    ReserveTx(tx *gorm.DB, in ReserveInput) error
    ConfirmTx(tx *gorm.DB, bizType, bizNo string) error
    ReleaseTx(tx *gorm.DB, bizType, bizNo, reason string) error

    DeductTx(tx *gorm.DB, in DeductInput) error
    RestoreTx(tx *gorm.DB, in RestoreInput) error

    SyncSkuTx(tx *gorm.DB, in SyncSkuInput) error
    GetSellableStock(ctx context.Context, skuIDs []uint64) ([]SellableStock, error)
}
```

### 6.2 语义边界

1. `Reserve/Confirm/Release`
   - 面向订单式库存交易
   - 对应“下单预占、支付确认、取消释放”
2. `Deduct/Restore`
   - 面向积分商城、独立出入库等非订单场景
3. `SyncSkuTx`
   - 面向商品与 SKU 初始化
   - 不暴露仓储细节给 `product`
4. `GetSellableStock`
   - 对外只暴露“可售库存”语义
   - 不让营销、商品前台直接感知底层库存表结构

### 6.3 Router 职责

router 的职责仅限于：

1. 读取当前库存配置
2. 返回当前 provider
3. 在必要时做 provider 可用性校验

router 不负责堆业务逻辑，也不承担重试补偿。

---

## 7. 配置模型设计

### 7.1 基础配置

建议新增统一库存配置：

```yaml
inventory:
  provider: local          # local | builtin_wms | external_wms
  external_mode: sync      # sync | async，仅 provider=external_wms 时有效
```

### 7.2 插件启用关系

`plugins.enabled` 继续保留，但职责仅用于插件装配，不直接代表库存真源：

```yaml
plugins:
  enabled:
    - product
    - order
    - wms
```

### 7.3 外部 WMS 配置

```yaml
external_wms:
  endpoint: https://wms.example.com/api
  app_key: your-app-key
  app_secret: your-app-secret
  timeout_ms: 3000
  callback_enabled: true
  retry:
    max_attempts: 8
    backoff_seconds: 30
```

### 7.4 启动校验规则

系统启动时应增加配置一致性校验：

1. `provider=local`
   - `wms` 可不启用
2. `provider=builtin_wms`
   - 必须启用 `wms`
3. `provider=external_wms`
   - 默认不要求启用 `wms`
4. `provider=external_wms + async`
   - 必须启用异步任务执行能力

---

## 8. 三类库存模式的交易流

### 8.1 `local`

下单：

1. 锁定 `product_skus.stock`
2. 写本地预占记录
3. 订单提交

支付：

1. 确认预占
2. 正式扣减本地库存

取消：

1. 释放预占

说明：

- 这是默认轻量模式
- 允许商城在完全不启用 `wms` 时独立运行

### 8.2 `builtin_wms`

下单：

1. 调用内置 WMS provider 的 `Reserve`

支付：

1. 调用内置 WMS provider 的 `Confirm`

取消：

1. 调用内置 WMS provider 的 `Release`

说明：

- 保留现有 `wms` 预占模型
- 但业务侧通过统一接口接入，不再直接 import `wms/service`

### 8.3 `external_wms.sync`

下单：

1. 订单服务调用外部 WMS 预占
2. 预占成功后才允许订单事务提交

支付：

1. 调用外部 WMS 确认

取消：

1. 调用外部 WMS 释放

优点：

- 一致性强
- 企业侧库存系统始终是实时真源

风险：

- 外部 WMS 可用性直接影响商城下单成功率

### 8.4 `external_wms.async`

下单：

1. 商城先写订单与库存意图单
2. 订单进入 `inventory_pending`
3. 异步任务投递外部 WMS 预占
4. 预占成功后才允许进入正常待支付状态

支付：

1. 商城写确认任务
2. 异步推送外部 WMS 确认出库

取消：

1. 商城写释放任务
2. 异步推送外部 WMS 释放预占

优点：

- 商城韧性更高
- 能容忍企业 WMS 短时不可用

风险：

- 状态更复杂
- 必须引入显式补偿与失败处理模型

---

## 9. 状态模型设计

### 9.1 新增库存交易状态

建议为订单增加独立库存交易状态字段，例如：

- `inventory_status`
  - `none`
  - `pending`
  - `reserved`
  - `confirmed`
  - `released`
  - `failed`

### 9.2 状态约束

1. `inventory_pending`
   - 不允许支付
2. `inventory_failed`
   - 不允许继续履约
   - 需人工重试或取消
3. 仅 `inventory_reserved`
   - 才允许进入正常支付链路

### 9.3 设计目的

这样可以避免异步模式下出现“订单创建成功但库存尚未锁定”的脏状态，也避免继续把复杂库存语义硬塞进原有订单主状态。

---

## 10. 幂等、补偿与回调设计

### 10.1 幂等键

建议统一库存交易幂等键：

- `biz_type`
- `biz_no`
- `action`
- `sku_id`

示例：

- `order + 202606090001 + reserve + 10001`
- `order + 202606090001 + confirm + 10001`

### 10.2 基本要求

1. 所有外部 WMS 请求都必须携带幂等号
2. 本地必须记录请求与响应日志
3. `reserve/confirm/release` 必须允许重复调用而不重复记账
4. 异步任务采用“至少一次投递”，靠幂等保证最终一致

### 10.3 集成任务表

建议新增集成任务表，例如 `inventory_integration_tasks`：

- `id`
- `provider`
- `biz_type`
- `biz_no`
- `action`
- `payload`
- `status`
- `attempt_count`
- `next_retry_at`
- `last_error`
- `created_at`
- `updated_at`

用途：

1. 出站请求投递
2. 失败重试
3. 死信跟踪
4. 人工补单与审计

### 10.4 外部回调

当外部 WMS 需要回调商城时，可新增明确语义的集成接口：

- `/admin/api/external-wms/callback`
- `/admin/api/inventory/tasks`
- `/admin/api/inventory/tasks/:id/retry`

新增这些接口是合理的，因为现有订单接口并不承载外部集成任务或回调管理语义。

---

## 11. 对现有业务模块的影响

### 11.1 `order`

需要改造：

1. 移除对 `wms/service` 的直接引用
2. 改为依赖统一 `inventory` provider
3. 为异步模式增加 `inventory_status` 约束

保留不变：

1. 下单接口路径
2. 支付接口路径
3. 取消订单接口路径

### 11.2 `product`

需要改造：

1. 新建 SKU 时改为调用 `inventory.SyncSkuTx`
2. 不再直接初始化 WMS 库存行

保留不变：

1. 商品管理接口主体语义

### 11.3 `points_mall`

建议改造：

1. 逐步将积分商品库存扣减升级为统一库存能力
2. 优先使用 `Deduct/Restore` 接口，避免形成第二套库存交易模型

### 11.4 `marketing`

建议改造：

1. 读取统一“可售库存”口径
2. 不直接假设库存来自 `product_skus.stock` 或 `wms`

---

## 12. 兼容升级策略

遵循“优先升级现有接口”的规则，本次设计不建议先新增一套订单主接口。

推荐策略：

1. 订单创建接口保持不变
2. 订单支付接口保持不变
3. 订单取消接口保持不变
4. 商品保存接口保持不变
5. 如需对外暴露新状态，只在返回结构中增量增加字段

只有以下场景新增接口是必要的：

1. 外部 WMS 回调入口
2. 库存集成任务查询与重试
3. 外部集成状态运维入口

---

## 13. 实施顺序建议

建议按四个阶段实施，避免一次性硬切：

### 阶段一：抽象库存域

1. 新增 `core/inventory`
2. 定义 provider 接口、模型、错误与 router

### 阶段二：业务改造解耦

1. `order` 改为依赖 `inventory`
2. `product` 改为依赖 `inventory`
3. `points_mall` 视范围接入统一扣减接口

### 阶段三：包装内置 WMS

1. 将当前 `wms` 实现包装为 `builtin_wms provider`
2. 保持现有仓储后台不变

### 阶段四：接入外部 WMS

1. 实现 `external_wms provider`
2. 先支持 `sync`
3. 再补 `async` 的任务、重试、回调、运维能力

这样做的好处是：

1. 每阶段可独立验证
2. 不会一次性破坏现有 `wms` 功能
3. 风险集中在边界清晰的局部模块

---

## 14. 测试设计

### 14.1 单元测试

1. inventory router 选择逻辑
2. local provider 的预占、确认、释放
3. builtin_wms provider 的适配正确性
4. external_wms provider 的请求签名、幂等、错误映射

### 14.2 集成测试

1. `local` 模式下完整下单、支付、取消链路
2. `builtin_wms` 模式下完整预占交易链路
3. `external_wms.sync` 模式下外部失败回滚链路
4. `external_wms.async` 模式下任务重试与补偿链路

### 14.3 回归测试

1. 现有仓储后台功能不回归
2. 商品创建与 SKU 保存不回归
3. 营销活动库存校验不回归
4. 积分商城库存扣减不出现双轨逻辑冲突

---

## 15. docs-site 同步要求

本次属于系统功能与架构变更，实施落地时必须同步更新 `docs-site`，至少覆盖：

1. 功能说明
   - 统一库存架构
   - `local / builtin_wms / external_wms` 三态模式
   - 外部 WMS 的 `sync / async` 模式
2. 接口变化
   - 订单库存语义升级
   - 新增的集成任务与回调接口
3. 配置影响
   - `inventory.provider`
   - `inventory.external_mode`
   - `external_wms` 相关配置
4. 部署影响
   - 若启用异步模式，需要任务执行与重试能力

对外文档默认应直接描述“当前最新库存兼容架构”，而不是继续以“WMS 是唯一真源”的旧表述作为主线。

---

## 16. 成功标准

本设计落地后的验收标准如下：

1. `wms` 插件可默认不启用，商城仍可正常交易。
2. `order` 不再因为 `wms` 未启用而无法装载。
3. 商城可在不启用内置 `wms` 的情况下对接企业外部 WMS。
4. 外部 WMS 同时支持同步与异步两种对接模式。
5. 订单、商品、积分商城等模块只依赖统一库存能力，不再直接绑定 `wms`。
6. docs-site 对外文档完成统一库存架构更新。
