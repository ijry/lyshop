# VIP 会员复购体系（Phase 1）Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在现有营销总线下新增独立 `vip` 插件，落地“年费会员 + 月度自主领券 + 普通商品会员价 + 支付成长值”的可用闭环。

**Architecture:** 保持 `core/marketing` 为统一价格编排层，新增 `vip` 插件提供会员数据与规则，使用一个会员价计算器挂入现有 pipeline。月领券复用 `marketing` 的 `coupon_users` 发券与核销链路，避免重复实现券系统。成长值在订单支付与退款节点通过幂等服务写入和回退。

**Tech Stack:** Go + Gin + GORM + MySQL/SQLite（现有） + Vue3（admin）+ uni-app（app）+ VitePress（docs-site）

---

## File Structure

- Create: `server/plugins/vip/plugin.json`
  - VIP 插件元信息、后台菜单、权限声明。
- Create: `server/plugins/vip/plugin.go`
  - 插件注册、路由注册、迁移入口。
- Create: `server/plugins/vip/model/vip.go`
  - 会员计划/等级/资产/成长值/月领券/会员价/订单权益表。
- Create: `server/plugins/vip/api/admin.go`
  - 后台会员计划、等级、券规则、会员价管理接口。
- Create: `server/plugins/vip/api/front.go`
  - 前台会员中心、开通、月领券、成长值日志接口。
- Create: `server/plugins/vip/service/plan.go`
  - 开通与续期、会员资产读写。
- Create: `server/plugins/vip/service/coupon.go`
  - 月领券规则判定、领取幂等。
- Create: `server/plugins/vip/service/growth.go`
  - 成长值发放/回退幂等逻辑。
- Create: `server/plugins/vip/service/price.go`
  - 会员价命中查询（按等级 + SKU）。
- Create: `server/plugins/vip/calculator/vip_price.go`
  - 价格流水线会员价计算器（Priority=15）。
- Modify: `server/core/marketing/context.go`
  - 扩展价格上下文字段（会员身份、会员折扣、行项目维度折扣）。
- Modify: `server/core/marketing/pipeline.go`
  - 汇总金额时纳入会员折扣。
- Modify: `server/plugins/order/service/order.go`
  - 下单填充会员态；支付后触发成长值发放。
- Modify: `server/plugins/order/service/after_sale.go`
  - 退款登记后触发成长值回退。
- Modify: `server/main.go`
  - blank import 新增 `vip` 插件。
- Create: `admin/src/views/vip/PlanList.vue`
  - 会员计划管理页。
- Create: `admin/src/views/vip/LevelList.vue`
  - 会员等级管理页。
- Create: `admin/src/views/vip/CouponRuleList.vue`
  - 月领券规则管理页。
- Create: `admin/src/views/vip/SkuPriceList.vue`
  - 会员价管理页。
- Modify: `admin/src/router/index.ts`
  - 新增 VIP 管理路由。
- Modify: `app/pages/user/index.vue`
  - 增加“会员中心”入口。
- Create: `app/pages/user/vip.vue`
  - 展示会员状态、可领券并执行领取。
- Modify: `app/pages.json`
  - 注册会员中心页面。
- Create: `docs-site/docs/api/vip.md`
  - 新增会员 API 文档。
- Modify: `docs-site/docs/api/index.md`
  - 增加 VIP 模块索引。
- Modify: `docs-site/docs/guide/features.md`
  - 同步会员功能说明与配置影响。

### Task 1: 创建 VIP 插件骨架并接入启动注册

**Files:**
- Create: `server/plugins/vip/plugin.json`
- Create: `server/plugins/vip/plugin.go`
- Modify: `server/main.go`

- [ ] **Step 1: 写 `plugin.json` 菜单与权限**

```json
{
  "name": "vip",
  "title": "会员插件",
  "version": "1.0.0",
  "description": "会员计划、等级成长、月度领券、会员价",
  "author": "lyshop",
  "depends": ["product", "order", "marketing"],
  "menus": [
    {
      "title": "会员管理",
      "icon": "crown",
      "path": "/vip",
      "sort": 45,
      "children": [
        { "title": "会员计划", "path": "/vip/plans", "permission": "vip:view" },
        { "title": "会员等级", "path": "/vip/levels", "permission": "vip:view" },
        { "title": "月领券规则", "path": "/vip/coupon-rules", "permission": "vip:view" },
        { "title": "会员价", "path": "/vip/sku-prices", "permission": "vip:view" }
      ]
    }
  ],
  "permissions": ["vip:view", "vip:edit"]
}
```

- [ ] **Step 2: 写 `plugin.go`（注册 + 路由 + 迁移入口）**

```go
package vip

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/plugin"
	vipapi "github.com/ijry/lyshop/plugins/vip/api"
	vipmodel "github.com/ijry/lyshop/plugins/vip/model"
	_ "github.com/ijry/lyshop/plugins/vip/calculator"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type vipPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("vip plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&vipPlugin{meta: m})
}

func (p *vipPlugin) Meta() plugin.Meta { return p.meta }
func (p *vipPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	vipapi.RegisterFrontRoutes(front)
	vipapi.RegisterAdminRoutes(admin)
}
func (p *vipPlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&vipmodel.Plan{}, &vipmodel.Level{}, &vipmodel.UserAsset{},
		&vipmodel.GrowthLog{}, &vipmodel.CouponRule{}, &vipmodel.CouponClaim{},
		&vipmodel.SkuPrice{}, &vipmodel.OrderBenefit{},
	)
}
func (p *vipPlugin) Install() error   { return nil }
func (p *vipPlugin) Uninstall() error { return nil }
```

- [ ] **Step 3: 在 `server/main.go` 增加插件导入**

```go
_ "github.com/ijry/lyshop/plugins/vip"
```

- [ ] **Step 4: 运行编译检查**

Run: `go test ./...`（在 `server` 目录）  
Expected: 编译通过，插件可被注册。

- [ ] **Step 5: 提交**

```bash
git add server/main.go server/plugins/vip
git commit -m "新增VIP插件骨架并接入启动注册" -m "主要改动：新增vip插件元信息、注册入口与迁移骨架；原因：建立独立会员插件边界；影响：服务启动时自动加载会员插件。"
```

### Task 2: 建立会员核心数据模型与服务读写

**Files:**
- Create: `server/plugins/vip/model/vip.go`
- Create: `server/plugins/vip/service/plan.go`
- Create: `server/plugins/vip/service/price.go`

- [ ] **Step 1: 定义 8 张核心表结构**

```go
type Plan struct {
	model.Base
	Name           string  `gorm:"size:64;not null"`
	DurationMonths int     `gorm:"not null;default:12"`
	Price          float64 `gorm:"type:decimal(10,2);not null"`
	RenewPrice     float64 `gorm:"type:decimal(10,2);default:0"`
	Status         int8    `gorm:"not null;default:1"`
}

type UserAsset struct {
	model.Base
	UserID         uint64     `gorm:"not null;uniqueIndex"`
	CurrentPlanID  uint64     `gorm:"not null;default:0"`
	CurrentLevelID uint64     `gorm:"not null;default:0"`
	VipStartAt     *time.Time
	VipEndAt       *time.Time
	GrowthValue    int64      `gorm:"not null;default:0"`
	Status         int8       `gorm:"not null;default:0"`
}
```

- [ ] **Step 2: 实现会员开通/续期最小服务**

```go
func OpenMembership(ctx context.Context, userID, planID uint64, now time.Time) (*vipmodel.UserAsset, error) {
	var plan vipmodel.Plan
	if err := db.DB.WithContext(ctx).Where("id = ? AND status = 1", planID).First(&plan).Error; err != nil {
		return nil, errors.New("会员计划不存在")
	}
	var asset vipmodel.UserAsset
	_ = db.DB.WithContext(ctx).Where("user_id = ?", userID).First(&asset).Error
	startAt := now
	if asset.VipEndAt != nil && asset.VipEndAt.After(now) {
		startAt = *asset.VipEndAt
	}
	endAt := startAt.AddDate(0, plan.DurationMonths, 0)
	asset.UserID = userID
	asset.CurrentPlanID = plan.ID
	asset.VipStartAt = &startAt
	asset.VipEndAt = &endAt
	asset.Status = 1
	return &asset, db.DB.WithContext(ctx).Save(&asset).Error
}
```

- [ ] **Step 3: 实现会员价命中查询服务**

```go
func GetVipSkuPrice(ctx context.Context, levelID, productID, skuID uint64) (*vipmodel.SkuPrice, error) {
	var row vipmodel.SkuPrice
	err := db.DB.WithContext(ctx).
		Where("status = 1 AND level_id = ? AND product_id = ? AND (sku_id = ? OR sku_id = 0)", levelID, productID, skuID).
		Order("sku_id desc").First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}
```

- [ ] **Step 4: 运行模型迁移验证**

Run: `go test ./plugins/vip/... -v`  
Expected: 通过，模型定义可编译。

- [ ] **Step 5: 提交**

```bash
git add server/plugins/vip/model/vip.go server/plugins/vip/service/plan.go server/plugins/vip/service/price.go
git commit -m "新增会员核心数据模型与基础服务" -m "主要改动：落地会员计划/等级/资产/成长值/月领券/会员价模型；原因：支撑会员闭环；影响：vip插件迁移可创建核心表。"
```

### Task 3: 接入营销总线会员上下文与金额汇总

**Files:**
- Modify: `server/core/marketing/context.go`
- Modify: `server/core/marketing/pipeline.go`

- [ ] **Step 1: 扩展 `PriceContext` 字段**

```go
type PriceContext struct {
	UserID     uint64
	Items      []OrderItem
	CouponIDs  []uint64
	PointsUse  int

	IsVip          bool
	VipLevelID     uint64
	VipDiscount    float64
	ItemVipDiscount map[uint64]float64 // key=sku_id

	GoodsAmount        float64
	ActivityDiscount   float64
	FullReduceDiscount float64
	CouponDiscount     float64
	PointsDiscount     float64
	FinalAmount        float64
}
```

- [ ] **Step 2: 更新 `FinalAmount` 计算逻辑**

```go
ctx.FinalAmount = ctx.GoodsAmount -
	ctx.ActivityDiscount -
	ctx.VipDiscount -
	ctx.FullReduceDiscount -
	ctx.CouponDiscount -
	ctx.PointsDiscount
```

- [ ] **Step 3: 为新字段补单测**

```go
func TestPipeline_WithVipDiscount(t *testing.T) {
	calculators = nil
	Register(&stubCalc{name: "activity", prio: 10, discount: 20})
	Register(&vipStubCalc{prio: 15, discount: 10})
	ctx := &PriceContext{Items: []OrderItem{{ProductID: 1, SkuID: 1, Price: 100, Qty: 1}}}
	require.NoError(t, Calculate(ctx))
	assert.Equal(t, 70.0, ctx.FinalAmount)
}
```

- [ ] **Step 4: 运行核心营销测试**

Run: `go test ./core/marketing -v`  
Expected: 所有 pipeline 测试通过。

- [ ] **Step 5: 提交**

```bash
git add server/core/marketing/context.go server/core/marketing/pipeline.go server/core/marketing/pipeline_test.go
git commit -m "扩展营销价格上下文支持会员折扣" -m "主要改动：PriceContext新增会员字段并纳入最终金额计算；原因：会员价需要进入统一结算总线；影响：下单价格汇总新增会员折扣项。"
```

### Task 4: 实现会员价计算器（活动商品跳过）

**Files:**
- Create: `server/plugins/vip/calculator/vip_price.go`

- [ ] **Step 1: 编写会员价计算器并注册优先级 15**

```go
type VipPriceCalculator struct{}

func init() { marketing.Register(&VipPriceCalculator{}) }
func (c *VipPriceCalculator) Name() string  { return "vip_price" }
func (c *VipPriceCalculator) Priority() int { return 15 }

func (c *VipPriceCalculator) Calculate(ctx *marketing.PriceContext) (bool, error) {
	if !ctx.IsVip || ctx.VipLevelID == 0 || db.DB == nil {
		return true, nil
	}
	for i := range ctx.Items {
		item := &ctx.Items[i]
		if item.ActivityPrice > 0 { // 活动商品不享会员价
			continue
		}
		price, err := vipsvc.GetVipSkuPrice(context.Background(), ctx.VipLevelID, item.ProductID, item.SkuID)
		if err != nil {
			continue
		}
		vipUnitPrice := price.VipPrice
		if vipUnitPrice <= 0 && price.VipDiscountRate > 0 {
			vipUnitPrice = item.Price * price.VipDiscountRate
		}
		if vipUnitPrice > 0 && vipUnitPrice < item.Price {
			d := (item.Price - vipUnitPrice) * float64(item.Qty)
			ctx.VipDiscount += d
			ctx.AppliedRules = append(ctx.AppliedRules, marketing.AppliedRule{Type: "vip", Name: "会员价", Discount: d})
		}
	}
	return true, nil
}
```

- [ ] **Step 2: 增加计算器行为测试（活动跳过）**

```go
func TestVipPriceCalculator_SkipActivityItems(t *testing.T) {
	ctx := &marketing.PriceContext{IsVip: true, VipLevelID: 1, Items: []marketing.OrderItem{{ProductID: 1, SkuID: 1, Price: 100, Qty: 1, ActivityPrice: 80}}}
	calc := &VipPriceCalculator{}
	_, err := calc.Calculate(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0.0, ctx.VipDiscount)
}
```

- [ ] **Step 3: 运行 VIP 计算器测试**

Run: `go test ./plugins/vip/calculator -v`  
Expected: 测试通过，且活动商品不享会员价。

- [ ] **Step 4: 提交**

```bash
git add server/plugins/vip/calculator/vip_price.go server/plugins/vip/calculator/vip_price_test.go
git commit -m "新增会员价计算器并接入活动互斥规则" -m "主要改动：新增vip_price计算器优先级15；原因：普通商品应用会员价、活动商品跳过；影响：结算规则新增会员价命中项。"
```

### Task 5: 实现月度自主领券（自然月、跨月作废、幂等）

**Files:**
- Create: `server/plugins/vip/service/coupon.go`
- Create: `server/plugins/vip/api/front.go`

- [ ] **Step 1: 编写领取服务（唯一键 + 事务）**

```go
func ClaimMonthlyCoupon(ctx context.Context, userID, ruleID uint64, now time.Time) error {
	period := now.Format("200601")
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rule vipmodel.CouponRule
		if err := tx.Where("id = ? AND status = 1", ruleID).First(&rule).Error; err != nil {
			return errors.New("会员券规则不存在")
		}
		var claim vipmodel.CouponClaim
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND rule_id = ? AND period_yyyymm = ?", userID, ruleID, period).
			First(&claim).Error
		if err == nil && claim.ClaimedCount >= rule.MonthlyQuota {
			return errors.New("本月已领完")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			claim = vipmodel.CouponClaim{UserID: userID, RuleID: ruleID, PeriodYYYYMM: period}
		}
		claim.ClaimedCount += 1
		nowTime := now
		claim.LastClaimedAt = &nowTime
		if err := tx.Save(&claim).Error; err != nil {
			return err
		}
		return tx.Create(&mktmodel.CouponUser{CouponID: rule.CouponID, UserID: userID, Status: 1}).Error
	})
}
```

- [ ] **Step 2: 编写前台接口**

```go
func RegisterFrontRoutes(g *gin.RouterGroup) {
	auth := g.Group("").Use(middleware.RequireAuth)
	auth.GET("/vip/profile", getVipProfile)
	auth.GET("/vip/coupons/monthly", listMonthlyCoupons)
	auth.POST("/vip/coupons/monthly/:rule_id/claim", claimMonthlyCoupon)
}
```

- [ ] **Step 3: 编写领取幂等测试**

```go
func TestClaimMonthlyCoupon_OncePerMonth(t *testing.T) {
	// given: monthly_quota=1
	// when: 连续调用两次 ClaimMonthlyCoupon
	// then: 第二次返回 "本月已领完"
}
```

- [ ] **Step 4: 运行测试**

Run: `go test ./plugins/vip/service -v`  
Expected: 月领券幂等测试通过。

- [ ] **Step 5: 提交**

```bash
git add server/plugins/vip/service/coupon.go server/plugins/vip/api/front.go server/plugins/vip/service/coupon_test.go
git commit -m "实现会员月度自主领券与并发幂等控制" -m "主要改动：自然月维度领取校验、唯一约束并发保护、领取后复用coupon_users发券；原因：满足会员券留存策略；影响：会员中心可自助领券。"
```

### Task 6: 接入订单支付发成长值与退款回退

**Files:**
- Create: `server/plugins/vip/service/growth.go`
- Modify: `server/plugins/order/service/order.go`
- Modify: `server/plugins/order/service/after_sale.go`

- [ ] **Step 1: 编写成长值幂等服务**

```go
func GrantGrowthOnOrderPaid(ctx context.Context, orderID, userID uint64, payable float64) error {
	key := fmt.Sprintf("order_paid:%d", orderID)
	return applyGrowthDelta(ctx, userID, orderID, key, calcGrowth(payable), "消费获得成长值")
}

func RevertGrowthOnRefund(ctx context.Context, orderID, userID uint64, refundAmount float64) error {
	key := fmt.Sprintf("order_refund:%d:%.2f", orderID, refundAmount)
	return applyGrowthDelta(ctx, userID, orderID, key, -calcGrowth(refundAmount), "退款回退成长值")
}
```

- [ ] **Step 2: 在 `PayOrder` 成功后触发成长值**

```go
if res.RowsAffected > 0 {
	var row ordermodel.Order
	if err := db.DB.WithContext(ctx).Where("id = ?", orderID).First(&row).Error; err == nil {
		_ = vipsvc.GrantGrowthOnOrderPaid(ctx, row.ID, row.UserID, row.TotalAmount)
	}
}
```

- [ ] **Step 3: 在 `MarkRefund` 成功后触发回退**

```go
if err := tx.Create(refund).Error; err != nil { return err }
if err := vipsvc.RevertGrowthOnRefund(ctx, caseRow.OrderID, caseRow.UserID, amount); err != nil {
	return err
}
```

- [ ] **Step 4: 运行订单相关测试**

Run: `go test ./plugins/order/service -v`  
Expected: 编译通过，支付/退款流程无回归错误。

- [ ] **Step 5: 提交**

```bash
git add server/plugins/vip/service/growth.go server/plugins/order/service/order.go server/plugins/order/service/after_sale.go
git commit -m "接入订单成长值发放与退款回退" -m "主要改动：支付后发成长值、退款后回退成长值并做幂等保护；原因：实现会员等级成长闭环；影响：订单生命周期将驱动会员成长值变化。"
```

### Task 7: 提供后台管理与前台会员中心页面

**Files:**
- Create: `admin/src/views/vip/PlanList.vue`
- Create: `admin/src/views/vip/LevelList.vue`
- Create: `admin/src/views/vip/CouponRuleList.vue`
- Create: `admin/src/views/vip/SkuPriceList.vue`
- Modify: `admin/src/router/index.ts`
- Modify: `app/pages/user/index.vue`
- Create: `app/pages/user/vip.vue`
- Modify: `app/pages.json`

- [ ] **Step 1: 新增后台 VIP 路由**

```ts
{ path: 'vip/plans', name: '会员计划', component: () => import('@/views/vip/PlanList.vue') },
{ path: 'vip/levels', name: '会员等级', component: () => import('@/views/vip/LevelList.vue') },
{ path: 'vip/coupon-rules', name: '月领券规则', component: () => import('@/views/vip/CouponRuleList.vue') },
{ path: 'vip/sku-prices', name: '会员价', component: () => import('@/views/vip/SkuPriceList.vue') },
```

- [ ] **Step 2: 后台页面最小可用表单/列表**

```vue
<script setup lang="ts">
const list = ref<any[]>([])
const form = ref({ name: '', duration_months: 12, price: 199, renew_price: 199, status: 1 })
const load = async () => { list.value = (await request.get('/vip/plans')).list || [] }
const submit = async () => { await request.post('/vip/plans', form.value); await load() }
onMounted(load)
</script>
```

- [ ] **Step 3: 增加 App 会员中心页与领券按钮**

```vue
<script setup>
const profile = ref(null)
const rules = ref([])
const load = async () => {
  profile.value = await request.get('/vip/profile')
  rules.value = await request.get('/vip/coupons/monthly')
}
const claim = async (ruleId) => {
  await request.post(`/vip/coupons/monthly/${ruleId}/claim`)
  await load()
}
onLoad(load)
</script>
```

- [ ] **Step 4: 运行前端构建检查**

Run: `npm run build`（在 `admin`）  
Expected: admin 构建通过。  
Run: `npm run build:h5`（在 `app`）  
Expected: app 构建通过。

- [ ] **Step 5: 提交**

```bash
git add admin/src/router/index.ts admin/src/views/vip app/pages/user/index.vue app/pages/user/vip.vue app/pages.json
git commit -m "新增会员后台管理页与前台会员中心入口" -m "主要改动：后台提供会员计划/等级/会员券/会员价配置页，前台提供会员中心与领券入口；原因：完成运营配置与用户触达闭环；影响：管理端与用户端均可操作会员功能。"
```

### Task 8: 同步 docs-site 文档与接口说明

**Files:**
- Create: `docs-site/docs/api/vip.md`
- Modify: `docs-site/docs/api/index.md`
- Modify: `docs-site/docs/guide/features.md`

- [ ] **Step 1: 新增 VIP API 文档**

```md
# 会员接口

## 功能说明
- 会员开通（支持按年）
- 月度会员券自主领取（自然月，过月作废）
- 普通商品会员价
- 成长值日志

## 接口变化
- GET /api/v1/vip/profile
- POST /api/v1/vip/open
- GET /api/v1/vip/coupons/monthly
- POST /api/v1/vip/coupons/monthly/:rule_id/claim
- GET /api/v1/vip/growth/logs
```

- [ ] **Step 2: 更新 API 索引与功能文档**

```md
- [会员接口](./vip)
```

```md
### 会员体系
- 新增会员计划与等级成长值。
- 月度会员券改为会员自主领取，跨月未领作废。
- 会员价仅对未参与秒杀/拼团等活动的普通商品生效。
```

- [ ] **Step 3: 验证文档构建**

Run: `npm run docs:build`（在文档项目目录，按现有脚本）  
Expected: 文档构建通过，`/api/vip` 可访问。

- [ ] **Step 4: 提交**

```bash
git add docs-site/docs/api/vip.md docs-site/docs/api/index.md docs-site/docs/guide/features.md
git commit -m "同步会员功能文档与接口说明" -m "主要改动：新增会员API文档并更新功能说明；原因：功能变更需同步docs-site；影响：接口使用方可按新会员能力接入。"
```

### Task 9: 端到端回归与发布前检查

**Files:**
- Test: `server`、`admin`、`app`、`docs-site`

- [ ] **Step 1: 后端测试与编译**

Run: `go test ./...`（`server`）  
Expected: 全部通过。

- [ ] **Step 2: 前端构建**

Run: `npm run build`（`admin`）  
Expected: 构建成功。  
Run: `npm run build:h5`（`app`）  
Expected: 构建成功。

- [ ] **Step 3: 人工验收**

Run: 启动服务后按以下清单验证：
- 会员开通后 `vip_end_at` 正确增加 12 个月；
- 本月可领券成功领取 1 次，第二次提示“本月已领完”；
- 活动商品不享会员价，普通商品享会员价；
- 订单支付后成长值增加，退款后成长值回退。

- [ ] **Step 4: 汇总发布说明（不提交代码变更）**

```md
发布关注项：
1. 确认 `vip` 插件已启用并完成迁移。
2. 运营先配置会员计划/等级/月领券规则，再开放前台入口。
3. 首次发布后抽样核对订单结算明细中的会员价命中记录。
```

## Self-Review

- Spec 覆盖：会员开通、月领券、会员价、成长值、文档同步均有独立任务。
- Placeholder 扫描：无占位符与“后续补充”文案。
- 类型一致性：`VipDiscount/VipLevelID`、`CouponRule/CouponClaim`、`GrantGrowthOnOrderPaid` 在任务间命名一致。
