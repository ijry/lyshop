# Marketing 插件拆分设计方案

## 概述

将 marketing 插件拆分为 5 个独立插件：
1. **marketing** - 保留优惠券和满减功能
2. **seckill** - 秒杀功能
3. **group_buy** - 团购功能
4. **bargain** - 砍价功能
5. **distribution** - 分销功能

## 拆分原则

### 1. 保持价格计算管道统一
所有活动类插件的价格计算器仍然注册到 `core/marketing` 管道，确保：
- 统一的价格计算流程
- 明确的优先级顺序
- 可插拔的计算器架构

### 2. 插件依赖关系
```
marketing (核心)
  ├── 优惠券
  └── 满减活动

seckill (秒杀)
  └── depends: ["product", "marketing"]

group_buy (团购)
  └── depends: ["product", "marketing"]

bargain (砍价)
  └── depends: ["product", "marketing"]

distribution (分销)
  └── depends: ["order", "marketing"]
```

### 3. 数据模型迁移策略

#### Marketing 插件保留
- `Coupon` - 优惠券
- `CouponUser` - 用户优惠券
- `Activity` (type='full_save') - 满减活动
- `ActivityProduct` (满减商品)

#### Seckill 插件
- `Activity` (type='seckill') → `SeckillActivity`
- `ActivityProduct` → `SeckillProduct`

#### GroupBuy 插件
- `Activity` (type='group_buy') → `GroupBuyActivity`
- `ActivityProduct` → `GroupBuyProduct`
- `GroupBuyOrder` - 拼团订单
- `GroupBuyMember` - 拼团成员

#### Bargain 插件
- `Activity` (type='bargain') → `BargainActivity`
- `ActivityProduct` → `BargainProduct`
- `BargainOrder` - 砍价订单
- `BargainHelper` - 砍价助手

#### Distribution 插件
- `Distributor` - 分销商
- `DistributionCommission` - 分销佣金
- `DistributionConfig` - 分销配置

## 详细设计

### 1. Marketing 插件（重构后）

**保留功能**：
- 优惠券系统
- 满减活动
- 价格计算管道核心

**数据模型**：
```go
// marketing/model/marketing.go
type Coupon struct { ... }
type CouponUser struct { ... }
type Activity struct { ... }  // 只保留 type='full_save'
type ActivityProduct struct { ... }
```

**服务层**：
- `coupon.go` - 优惠券服务
- `activity.go` - 满减活动服务（简化版）

**计算器**：
- `calculator/coupon.go` - 优惠券计算器（优先级30）
- `calculator/full_reduce.go` - 满减计算器（优先级20）

**API**：
- 优惠券 CRUD
- 满减活动 CRUD

---

### 2. Seckill 插件（秒杀）

**插件结构**：
```
server/plugins/seckill/
├── plugin.json
├── plugin.go
├── model/
│   └── seckill.go
├── service/
│   ├── activity.go
│   └── product.go
├── calculator/
│   └── seckill.go
└── api/
    ├── front.go
    └── admin.go
```

**数据模型**：
```go
type SeckillActivity struct {
    ID        uint64
    Name      string
    StartAt   *time.Time
    EndAt     *time.Time
    Status    int8
}

type SeckillProduct struct {
    ID              uint64
    ActivityID      uint64
    ProductID       uint64
    SkuID           uint64
    SeckillPrice    float64
    LimitPerOrder   int
    TotalStock      int
    SoldQty         int
}
```

**价格计算器**：
```go
// calculator/seckill.go
type SeckillCalculator struct{}

func (c *SeckillCalculator) Priority() int { return 10 }

func (c *SeckillCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
    // 查询秒杀商品
    // 应用秒杀价格
    // 标记为排他性活动
}
```

**前端页面**：
- Admin: `SeckillActivityList.vue`, `SeckillProductManage.vue`
- Eapp: `seckill/index.vue`, `seckill/detail.vue`
- App: `seckill/index.vue`, `seckill/detail.vue`
- Web: `Seckill.vue`

---

### 3. GroupBuy 插件（团购）

**插件结构**：
```
server/plugins/group_buy/
├── plugin.json
├── plugin.go
├── model/
│   └── group_buy.go
├── service/
│   ├── activity.go
│   ├── product.go
│   └── order.go
├── calculator/
│   └── group_buy.go
└── api/
    ├── front.go
    └── admin.go
```

**数据模型**：
```go
type GroupBuyActivity struct {
    ID          uint64
    Name        string
    GroupSize   int     // 成团人数
    ExpireHours int     // 过期小时数
    StartAt     *time.Time
    EndAt       *time.Time
    Status      int8
}

type GroupBuyProduct struct {
    ID            uint64
    ActivityID    uint64
    ProductID     uint64
    SkuID         uint64
    GroupPrice    float64
    LimitPerOrder int
}

type GroupBuyOrder struct {
    ID          uint64
    ActivityID  uint64
    ProductID   uint64
    LeaderID    uint64
    GroupSize   int
    JoinedCount int
    Status      string  // pending|success|failed
    ExpireAt    *time.Time
}

type GroupBuyMember struct {
    ID           uint64
    GroupOrderID uint64
    UserID       uint64
    OrderID      uint64
}
```

**核心服务**：
```go
// service/order.go
func CreateGroupBuyOrder(ctx, activityID, productID, userID) (*GroupBuyOrder, error)
func JoinGroupBuy(ctx, groupOrderID, userID, orderID) error
func CheckAndCompleteGroupBuy(ctx, groupOrderID) error
func ExpireGroupBuy(ctx, groupOrderID) error
```

**价格计算器**：
```go
type GroupBuyCalculator struct{}
func (c *GroupBuyCalculator) Priority() int { return 10 }
```

**前端页面**：
- Admin: 活动管理、商品管理、拼团订单管理
- Eapp: 拼团列表、创建拼团、加入拼团、拼团详情
- App: 拼团列表、拼团详情、邀请好友
- Web: 拼团活动页

---

### 4. Bargain 插件（砍价）

**插件结构**：
```
server/plugins/bargain/
├── plugin.json
├── plugin.go
├── model/
│   └── bargain.go
├── service/
│   ├── activity.go
│   ├── product.go
│   └── order.go
├── calculator/
│   └── bargain.go
└── api/
    ├── front.go
    └── admin.go
```

**数据模型**：
```go
type BargainActivity struct {
    ID          uint64
    Name        string
    MaxHelpers  int     // 最多砍价人数
    ExpireHours int
    StartAt     *time.Time
    EndAt       *time.Time
    Status      int8
}

type BargainProduct struct {
    ID          uint64
    ActivityID  uint64
    ProductID   uint64
    SkuID       uint64
    StartPrice  float64
    FloorPrice  float64
}

type BargainOrder struct {
    ID            uint64
    ActivityID    uint64
    ProductID     uint64
    UserID        uint64
    OriginalPrice float64
    CurrentPrice  float64
    FloorPrice    float64
    HelperCount   int
    MaxHelpers    int
    Status        string  // pending|completed|expired
    ExpireAt      *time.Time
}

type BargainHelper struct {
    ID             uint64
    BargainOrderID uint64
    UserID         uint64
    CutAmount      float64
}
```

**核心服务**：
```go
// service/order.go
func CreateBargainOrder(ctx, activityID, productID, userID) (*BargainOrder, error)
func HelpBargain(ctx, bargainOrderID, userID) (cutAmount float64, err error)
func CompleteBargain(ctx, bargainOrderID) error
func ExpireBargain(ctx, bargainOrderID) error
```

**价格计算器**：
```go
type BargainCalculator struct{}
func (c *BargainCalculator) Priority() int { return 10 }
```

**前端页面**：
- Admin: 活动管理、商品管理、砍价订单管理
- Eapp: 砍价列表、发起砍价、砍价详情
- App: 砍价列表、砍价详情、邀请好友砍价
- Web: 砍价活动页

---

### 5. Distribution 插件（分销）

**插件结构**：
```
server/plugins/distribution/
├── plugin.json
├── plugin.go
├── model/
│   └── distribution.go
├── service/
│   ├── distributor.go
│   ├── commission.go
│   └── withdraw.go
├── calculator/
│   └── distribution.go
└── api/
    ├── front.go
    └── admin.go
```

**数据模型**：
```go
type Distributor struct {
    ID        uint64
    UserID    uint64
    ParentID  uint64
    Level     int
    TotalEarn float64
    Balance   float64
    Status    int8
}

type DistributionCommission struct {
    ID            uint64
    OrderID       uint64
    DistributorID uint64
    Level         int
    Amount        float64
    Status        string  // pending|settled|refunded
}

type DistributionConfig struct {
    ID         uint64
    Level1Rate float64
    Level2Rate float64
    MinWithdraw float64
}

type WithdrawRecord struct {
    ID            uint64
    DistributorID uint64
    Amount        float64
    Status        string  // pending|approved|rejected|completed
    ApprovedAt    *time.Time
}
```

**核心服务**：
```go
// service/distributor.go
func RegisterDistributor(ctx, userID, parentID) error
func GetDistributorInfo(ctx, userID) (*Distributor, error)
func GetDistributorChain(ctx, userID) ([]uint64, error)

// service/commission.go
func CalculateCommission(ctx, orderID, buyerID, orderAmount) ([]Commission, error)
func SettleCommissions(ctx, orderID) error
func RefundCommissions(ctx, orderID) error

// service/withdraw.go
func CreateWithdrawRequest(ctx, distributorID, amount) error
func ApproveWithdraw(ctx, withdrawID) error
func RejectWithdraw(ctx, withdrawID, reason) error
func CompleteWithdraw(ctx, withdrawID) error
```

**价格计算器**：
```go
type DistributionCalculator struct{}
func (c *DistributionCalculator) Priority() int { return 50 }
// 不影响价格，只记录佣金
```

**前端页面**：
- Admin: 分销商管理、佣金管理、提现审核、分销配置
- Eapp: 分销商中心、我的佣金、提现申请、推广海报
- App: 成为分销商、我的团队、佣金明细、提现
- Web: 分销中心

---

## 价格计算管道集成

### 计算器优先级分配

| 优先级 | 计算器 | 插件 | 说明 |
|--------|--------|------|------|
| 10 | SeckillCalculator | seckill | 秒杀价格 |
| 10 | GroupBuyCalculator | group_buy | 团购价格 |
| 10 | BargainCalculator | bargain | 砍价价格 |
| 15 | VipPriceCalculator | vip | 会员价格 |
| 20 | FullReduceCalculator | marketing | 满减活动 |
| 30 | CouponCalculator | marketing | 优惠券 |
| 40 | PointsCalculator | points_mall | 积分抵扣 |
| 50 | DistributionCalculator | distribution | 分销佣金 |

### 排他性规则

- 秒杀、团购、砍价活动是排他性的（Exclusive=true）
- 当应用排他性活动后，停止后续计算器（除了分销）
- 会员价、满减、优惠券、积分可以叠加

### 计算器注册

所有计算器仍然注册到 `core/marketing` 管道：

```go
// seckill/calculator/seckill.go
func init() {
    marketing.Register(&SeckillCalculator{})
}

// group_buy/calculator/group_buy.go
func init() {
    marketing.Register(&GroupBuyCalculator{})
}

// bargain/calculator/bargain.go
func init() {
    marketing.Register(&BargainCalculator{})
}

// distribution/calculator/distribution.go
func init() {
    marketing.Register(&DistributionCalculator{})
}
```

---

## 迁移策略

### 阶段 1: 创建新插件（不影响现有功能）
1. 创建 seckill 插件目录和基础结构
2. 创建 group_buy 插件目录和基础结构
3. 创建 bargain 插件目录和基础结构
4. 创建 distribution 插件目录和基础结构

### 阶段 2: 迁移数据模型
1. 复制相关模型到新插件
2. 重命名模型（Activity → SeckillActivity）
3. 调整字段和关系

### 阶段 3: 迁移服务层
1. 复制服务代码到新插件
2. 调整导入路径
3. 完善业务逻辑

### 阶段 4: 迁移计算器
1. 复制计算器到新插件
2. 调整模型引用
3. 保持注册到 core/marketing

### 阶段 5: 迁移 API
1. 复制 API 端点到新插件
2. 调整路由路径
3. 更新服务调用

### 阶段 6: 迁移前端
1. 复制前端页面到新目录
2. 调整 API 调用路径
3. 更新路由配置

### 阶段 7: 数据迁移
1. 编写数据迁移脚本
2. 迁移 activities 表数据到新表
3. 迁移 activity_products 表数据到新表

### 阶段 8: 清理 marketing 插件
1. 删除已迁移的代码
2. 保留优惠券和满减功能
3. 更新文档

---

## 前端路由规划

### Admin 端
```
/marketing/coupons              → marketing 插件
/marketing/full-reduce          → marketing 插件
/seckill/activities             → seckill 插件
/seckill/products               → seckill 插件
/group-buy/activities           → group_buy 插件
/group-buy/products             → group_buy 插件
/group-buy/orders               → group_buy 插件
/bargain/activities             → bargain 插件
/bargain/products               → bargain 插件
/bargain/orders                 → bargain 插件
/distribution/distributors      → distribution 插件
/distribution/commissions       → distribution 插件
/distribution/withdraws         → distribution 插件
/distribution/config            → distribution 插件
```

### Eapp/App/Web 端
```
/marketing/coupons              → marketing 插件
/seckill/list                   → seckill 插件
/seckill/detail/:id             → seckill 插件
/group-buy/list                 → group_buy 插件
/group-buy/detail/:id           → group_buy 插件
/group-buy/my-groups            → group_buy 插件
/bargain/list                   → bargain 插件
/bargain/detail/:id             → bargain 插件
/bargain/my-bargains            → bargain 插件
/distribution/center            → distribution 插件
/distribution/team              → distribution 插件
/distribution/commission        → distribution 插件
/distribution/withdraw          → distribution 插件
```

---

## 优势

1. **职责清晰** - 每个插件专注于一个营销功能
2. **独立开发** - 各插件可以独立开发和测试
3. **按需启用** - 商家可以选择启用需要的营销功能
4. **易于维护** - 代码结构清晰，便于维护和扩展
5. **统一计算** - 价格计算管道保持统一，确保一致性

## 风险和注意事项

1. **数据迁移** - 需要谨慎处理现有数据的迁移
2. **API 兼容** - 需要保持 API 向后兼容或提供迁移指南
3. **前端更新** - 前端需要同步更新 API 调用路径
4. **测试覆盖** - 需要全面测试各插件的功能和集成
5. **文档更新** - 需要更新所有相关文档

---

## 实施计划

**预计时间**: 2-3 天

**Day 1**: 创建 seckill 和 group_buy 插件
**Day 2**: 创建 bargain 和 distribution 插件
**Day 3**: 完善前端、测试、文档

**优先级**:
1. Seckill（最常用）
2. GroupBuy（次常用）
3. Distribution（重要）
4. Bargain（可选）
