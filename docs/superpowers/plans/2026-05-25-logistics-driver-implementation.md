# 物流驱动化（快递100 + 快递鸟）Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为订单物流新增可插拔驱动体系（快递100/快递鸟），支持主备路由、渠道绑定、自动轮询、手动刷新，并在前后台展示完整轨迹。

**Architecture:** 新增 `core/driver/logistics` 作为统一抽象；`logistics_kuaidi100` 与 `logistics_kdniao` 插件注册驱动；`logistics_router` 插件负责主备选择、渠道绑定与轮询调度；`order` 插件负责轨迹落库与对外 API。

**Tech Stack:** Go 1.26 · Gin · GORM · Vue3 + Tailwind（admin/web）· uni-app（app）· VitePress（docs-site）

---

## 文件结构与职责边界

### 后端核心

- Create: `server/core/driver/logistics/logistics.go`
  - 物流驱动接口、标准响应结构、驱动注册/获取
- Create: `server/core/driver/logistics/logistics_test.go`
  - 驱动注册、主备获取、空驱动错误测试

### 物流驱动插件

- Create: `server/plugins/logistics_kuaidi100/plugin.json`
- Create: `server/plugins/logistics_kuaidi100/plugin.go`
- Create: `server/plugins/logistics_kuaidi100/driver.go`
- Create: `server/plugins/logistics_kdniao/plugin.json`
- Create: `server/plugins/logistics_kdniao/plugin.go`
- Create: `server/plugins/logistics_kdniao/driver.go`
- Modify: `server/main.go`
  - blank import 新增两个物流插件

### 路由插件

- Create: `server/plugins/logistics_router/plugin.json`
- Create: `server/plugins/logistics_router/plugin.go`
- Create: `server/plugins/logistics_router/service/router.go`
  - 主备路由、首查绑定、轮询入口
- Create: `server/plugins/logistics_router/service/polling.go`
  - 定时任务循环（开关+频率）

### 订单域改造

- Modify: `server/plugins/order/model/shipment.go`
  - `OrderShipment` 增加渠道与同步元数据
- Create: `server/plugins/order/model/shipment_track.go`
  - 轨迹明细、同步日志模型
- Modify: `server/plugins/order/plugin.go`
  - 迁移注册新表
- Modify: `server/plugins/order/service/shipment.go`
  - 同步轨迹、状态映射、去重落库
- Modify: `server/plugins/order/service/order.go`
  - 发货后可触发首次拉取（可配置）
- Modify: `server/plugins/order/api/admin.go`
  - 新增手动刷新与后台轨迹接口
- Modify: `server/plugins/order/api/front.go`
  - 新增前台轨迹查询接口

### 前端改造

- Modify: `admin/src/api/plugins.ts`
  - 新增后台手动刷新与轨迹查询 API
- Modify: `admin/src/views/order/OrderDetail.vue`
  - 轨迹节点展示 + 手动刷新
- Modify: `web/src/views/OrderDetail.vue`
  - 轨迹节点展示
- Modify: `app/pages/order/detail.vue`
  - 轨迹节点展示

### 文档

- Modify: `docs-site/docs/api/order.md`
- Modify: `docs-site/docs/guide/features.md`

---

### Task 1: 新增物流驱动核心抽象

**Files:**
- Create: `server/core/driver/logistics/logistics.go`
- Create: `server/core/driver/logistics/logistics_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestGetByName_NotFound(t *testing.T) {
	_, err := GetByName("missing")
	require.EqualError(t, err, `logistics driver "missing" not registered`)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server; go test ./core/driver/logistics -run TestGetByName_NotFound -v`  
Expected: FAIL（包尚未创建）

- [ ] **Step 3: Write minimal implementation**

```go
package logistics

type Driver interface {
	Name() string
	Query(ctx context.Context, req QueryReq) (*TrackResult, error)
}

type QueryReq struct {
	CompanyCode string
	TrackingNo  string
}

type TrackNode struct {
	Time       time.Time       `json:"time"`
	Location   string          `json:"location"`
	StatusCode string          `json:"status_code"`
	StatusText string          `json:"status_text"`
	RawPayload json.RawMessage `json:"raw_payload"`
}

type TrackResult struct {
	Provider    string      `json:"provider"`
	StatusCode  string      `json:"status_code"`
	StatusText  string      `json:"status_text"`
	SignedAt    *time.Time  `json:"signed_at"`
	Nodes       []TrackNode `json:"nodes"`
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server; go test ./core/driver/logistics -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/core/driver/logistics/logistics.go server/core/driver/logistics/logistics_test.go
git commit -m "新增物流驱动核心抽象

建立 logistics driver 的统一注册与查询接口，定义标准轨迹返回结构，作为快递100与快递鸟插件的公共契约。影响范围为后端驱动基础层。"
```

---

### Task 2: 接入快递100与快递鸟驱动插件

**Files:**
- Create: `server/plugins/logistics_kuaidi100/plugin.json`
- Create: `server/plugins/logistics_kuaidi100/plugin.go`
- Create: `server/plugins/logistics_kuaidi100/driver.go`
- Create: `server/plugins/logistics_kdniao/plugin.json`
- Create: `server/plugins/logistics_kdniao/plugin.go`
- Create: `server/plugins/logistics_kdniao/driver.go`
- Modify: `server/main.go`

- [ ] **Step 1: Write the failing test**

```go
func TestDrivers_RegisterSuccessfully(t *testing.T) {
	_, err1 := logistics.GetByName("kuaidi100")
	_, err2 := logistics.GetByName("kdniao")
	require.NoError(t, err1)
	require.NoError(t, err2)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server; go test ./plugins/logistics_kuaidi100 ./plugins/logistics_kdniao -v`  
Expected: FAIL（插件未创建）

- [ ] **Step 3: Write minimal implementation**

```go
// plugin.go
func (p *logisticsKuaidi100Plugin) Install() error {
	logistics.Register(&kuaidi100Driver{})
	return nil
}
```

```go
// driver.go (mapping excerpt)
func mapStatus(code string) string {
	switch strings.ToLower(strings.TrimSpace(code)) {
	case "signed", "delivered":
		return "signed"
	case "in_transit", "transport":
		return "in_transit"
	case "exception", "problem":
		return "exception"
	default:
		return "shipped"
	}
}
```

```go
// server/main.go
_ "github.com/ijry/lyshop/plugins/logistics_kuaidi100"
_ "github.com/ijry/lyshop/plugins/logistics_kdniao"
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server; go test ./plugins/logistics_kuaidi100 ./plugins/logistics_kdniao -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/plugins/logistics_kuaidi100 server/plugins/logistics_kdniao server/main.go
git commit -m "新增快递100与快递鸟物流驱动插件

实现两家物流服务商驱动与标准状态映射，并通过插件安装自动注册到 logistics driver 注册中心。影响范围为物流驱动插件与服务启动入口。"
```

---

### Task 3: 新增 logistics_router 路由插件（主备+轮询配置）

**Files:**
- Create: `server/plugins/logistics_router/plugin.json`
- Create: `server/plugins/logistics_router/plugin.go`
- Create: `server/plugins/logistics_router/service/router.go`
- Create: `server/plugins/logistics_router/service/polling.go`
- Modify: `server/main.go`

- [ ] **Step 1: Write the failing test**

```go
func TestSelectDriver_PrimaryFailFallbackSecondary(t *testing.T) {
	driver, err := SelectDriverForFirstQuery(ctx, "xxx")
	require.NoError(t, err)
	require.Equal(t, "kdniao", driver.Name())
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server; go test ./plugins/logistics_router/service -run TestSelectDriver_PrimaryFailFallbackSecondary -v`  
Expected: FAIL（路由服务未实现）

- [ ] **Step 3: Write minimal implementation**

```go
func SelectDriverForFirstQuery(ctx context.Context, companyCode string) (logistics.Driver, error) {
	primary := strings.TrimSpace(loadCfg("primary_driver", "kuaidi100"))
	secondary := strings.TrimSpace(loadCfg("secondary_driver", "kdniao"))

	pDrv, pErr := logistics.GetByName(primary)
	if pErr == nil {
		return pDrv, nil
	}
	sDrv, sErr := logistics.GetByName(secondary)
	if sErr == nil {
		return sDrv, nil
	}
	return nil, fmt.Errorf("主备物流驱动均不可用: primary=%v secondary=%v", pErr, sErr)
}
```

```go
func StartPollingLoop() {
	if !loadBool("polling_enabled", false) {
		return
	}
	interval := loadInt("polling_interval_seconds", 300)
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			_ = PollAndSyncShipments(context.Background(), 100)
		}
	}()
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server; go test ./plugins/logistics_router/service -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/plugins/logistics_router server/main.go
git commit -m "新增物流路由插件与轮询调度

引入 logistics_router 负责主备驱动选择、轮询开关与频率配置，为运单首次绑定与自动同步提供统一路由能力。影响范围为物流路由层与插件启动流程。"
```

---

### Task 4: 订单模型扩展与轨迹落库

**Files:**
- Modify: `server/plugins/order/model/shipment.go`
- Create: `server/plugins/order/model/shipment_track.go`
- Modify: `server/plugins/order/plugin.go`
- Modify: `server/plugins/order/service/shipment.go`
- Test: `server/plugins/order/service/shipment_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestSyncShipmentTracks_BindProviderAndDedupNodes(t *testing.T) {
	ctx, env := setupShipmentSyncEnv(t)
	err := SyncShipmentTracks(ctx, env.ShipmentID, SyncShipmentReq{Manual: true})
	require.NoError(t, err)
	err = SyncShipmentTracks(ctx, env.ShipmentID, SyncShipmentReq{Manual: true})
	require.NoError(t, err)

	ship := loadShipment(t, env.DB, env.ShipmentID)
	require.NotEmpty(t, ship.ChannelProvider)
	require.Equal(t, 0, ship.SyncFailCount)
	require.Equal(t, 3, countTracks(t, env.DB, env.ShipmentID))
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server; go test ./plugins/order/service -run TestSyncShipmentTracks_BindProviderAndDedupNodes -v`  
Expected: FAIL（模型与同步逻辑缺失）

- [ ] **Step 3: Write minimal implementation**

```go
// shipment.go add fields
ChannelProvider string     `gorm:"size:32;index" json:"channel_provider"`
LastQueryAt     *time.Time `json:"last_query_at"`
LastSyncOKAt    *time.Time `json:"last_sync_ok_at"`
SyncFailCount   int        `gorm:"not null;default:0" json:"sync_fail_count"`
```

```go
// shipment_track.go
type OrderShipmentTrack struct {
	model.Base
	ShipmentID  uint64          `gorm:"not null;index:idx_shipment_time" json:"shipment_id"`
	Provider    string          `gorm:"size:32;not null" json:"provider"`
	TrackHash   string          `gorm:"size:64;not null;uniqueIndex:uk_shipment_track" json:"track_hash"`
	StatusCode  string          `gorm:"size:32;not null" json:"status_code"`
	StatusText  string          `gorm:"size:255;not null" json:"status_text"`
	EventTime   time.Time       `gorm:"index:idx_shipment_time" json:"event_time"`
	Location    string          `gorm:"size:255" json:"location"`
	RawPayload  json.RawMessage `gorm:"type:json" json:"raw_payload"`
}
```

```go
func SyncShipmentTracks(ctx context.Context, shipmentID uint64, req SyncShipmentReq) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row, driver, err := lockAndResolveDriver(ctx, tx, shipmentID)
		if err != nil {
			return err
		}
		result, err := driver.Query(ctx, logistics.QueryReq{
			CompanyCode: row.Company,
			TrackingNo:  row.TrackingNo,
		})
		if err != nil {
			return markShipmentSyncFail(ctx, tx, row.ID, err)
		}
		return applyShipmentTrackResult(ctx, tx, row, result)
	})
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server; go test ./plugins/order/service -run TestSyncShipmentTracks_BindProviderAndDedupNodes -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/plugins/order/model/shipment.go server/plugins/order/model/shipment_track.go server/plugins/order/plugin.go server/plugins/order/service/shipment.go server/plugins/order/service/shipment_test.go
git commit -m "订单物流新增渠道绑定与轨迹明细落库

扩展 order_shipments 同步元数据并新增轨迹/同步日志模型，落地运单渠道绑定、轨迹去重与状态更新逻辑。影响范围为订单物流模型与服务层。"
```

---

### Task 5: 新增物流查询 API（前后台）与手动刷新

**Files:**
- Modify: `server/plugins/order/api/admin.go`
- Modify: `server/plugins/order/api/front.go`
- Modify: `server/plugins/order/service/order.go`
- Test: `server/plugins/order/api/order_logistics_api_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestAdminSyncShipment_API(t *testing.T) {
	r := setupAdminRouter(t)
	w := performPost(t, r, "/admin/api/orders/1/shipments/10/sync", adminToken, nil)
	require.Equal(t, 200, w.Code)
	body := decodeJSON(t, w.Body)
	require.Equal(t, float64(0), body["code"])
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server; go test ./plugins/order/api -run TestAdminSyncShipment_API -v`  
Expected: FAIL（接口未注册）

- [ ] **Step 3: Write minimal implementation**

```go
// admin.go register
g.POST("/orders/:id/shipments/:shipment_id/sync", middleware.RequirePermission("order:ship"), adminSyncShipment)
g.GET("/orders/:id/shipments/:shipment_id/tracks", middleware.RequirePermission("order:view"), adminGetShipmentTracks)
```

```go
// front.go register
auth.GET("/orders/:id/shipments/:shipment_id/tracks", myShipmentTracks)
```

```go
func adminSyncShipment(c *gin.Context) {
	shipmentID, _ := strconv.ParseUint(c.Param("shipment_id"), 10, 64)
	if err := ordersvc.SyncShipmentTracks(c.Request.Context(), shipmentID, ordersvc.SyncShipmentReq{Manual: true}); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server; go test ./plugins/order/api -run TestAdminSyncShipment_API -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/plugins/order/api/admin.go server/plugins/order/api/front.go server/plugins/order/service/order.go server/plugins/order/api/order_logistics_api_test.go
git commit -m "新增前后台物流轨迹查询与手动刷新接口

开放订单运单轨迹查询与后台手动同步能力，前后台可按 shipment 维度查看完整轨迹并触发即时刷新。影响范围为订单 API 层与物流查询入口。"
```

---

### Task 6: 前端接入轨迹时间线与后台刷新按钮

**Files:**
- Modify: `admin/src/api/plugins.ts`
- Modify: `admin/src/views/order/OrderDetail.vue`
- Modify: `web/src/views/OrderDetail.vue`
- Modify: `app/pages/order/detail.vue`

- [ ] **Step 1: Write the failing test**

```ts
it('admin order detail can trigger shipment sync', async () => {
  const wrapper = mount(OrderDetail, { global: mockOrderDetailDeps() })
  await wrapper.find('[data-test="sync-shipment"]').trigger('click')
  expect(mockSyncShipment).toHaveBeenCalled()
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd admin; npm run test -- OrderDetail.spec.ts`  
Expected: FAIL（缺少刷新按钮与 API）

- [ ] **Step 3: Write minimal implementation**

```ts
// admin/src/api/plugins.ts
export const syncShipment = (orderID: number, shipmentID: number) =>
  request.post(`/orders/${orderID}/shipments/${shipmentID}/sync`)
export const getShipmentTracks = (orderID: number, shipmentID: number) =>
  request.get(`/orders/${orderID}/shipments/${shipmentID}/tracks`)
```

```vue
<!-- admin OrderDetail.vue -->
<button data-test="sync-shipment" class="text-xs text-blue-600 hover:underline" @click="syncShip(ship.id)">刷新物流</button>
```

```ts
async function syncShip(shipmentID: number) {
  await syncShipment(Number(route.params.id), shipmentID)
  await loadDetail()
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd admin; npm run test -- OrderDetail.spec.ts`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add admin/src/api/plugins.ts admin/src/views/order/OrderDetail.vue web/src/views/OrderDetail.vue app/pages/order/detail.vue
git commit -m "前端接入物流轨迹刷新与时间线展示

后台订单详情新增手动刷新物流入口，web 与 app 订单详情补充轨迹明细展示，统一使用 shipment 轨迹接口。影响范围为三端订单详情页面。"
```

---

### Task 7: 同步 docs-site 文档与全链路回归

**Files:**
- Modify: `docs-site/docs/api/order.md`
- Modify: `docs-site/docs/guide/features.md`

- [ ] **Step 1: Write the failing doc check**

Run:  
`rg -n "快递100|快递鸟|channel_provider|轮询|shipments/.*/sync|tracks" docs-site/docs/api/order.md docs-site/docs/guide/features.md`

Expected: 关键内容缺失或不完整

- [ ] **Step 2: Update docs with concrete changes**

```md
### 物流驱动化升级

- 系统支持 `kuaidi100` / `kdniao` 双驱动。
- 运单首次同步按主备选择，成功后写入 `channel_provider` 并固定后续渠道。
- 后台可配置自动轮询开关与频率。

### 新增接口

- `POST /admin/api/orders/:id/shipments/:shipment_id/sync`
- `GET /admin/api/orders/:id/shipments/:shipment_id/tracks`
- `GET /api/v1/orders/:id/shipments/:shipment_id/tracks`
```

- [ ] **Step 3: Run targeted verification**

Run:  
`rg -n "快递100|快递鸟|channel_provider|轮询|/shipments/:shipment_id/sync|/shipments/:shipment_id/tracks" docs-site/docs/api/order.md docs-site/docs/guide/features.md`  
Expected: 全部命中

- [ ] **Step 4: Run core regression suites**

Run:
- `cd server; go test ./plugins/order/...`
- `cd server; go test ./plugins/logistics_router/...`
- `cd admin; npm run build`
- `cd web; npm run build`

Expected: 全部通过；如有历史失败，单列说明且不在本次顺带修复。

- [ ] **Step 5: Commit**

```bash
git add docs-site/docs/api/order.md docs-site/docs/guide/features.md
git commit -m "同步物流驱动化能力文档

更新订单 API 与功能说明，覆盖双驱动路由、渠道绑定、自动轮询配置及新增轨迹查询/刷新接口，确保功能变更与部署影响可追溯。影响范围为 docs-site 文档站。"
```

---

## 自检结果（计划层）

- Spec coverage：已覆盖驱动抽象、双插件实现、主备路由、渠道绑定、轮询配置、前后台接口、三端展示、文档同步。
- Placeholder scan：无 `TODO/TBD/待补` 占位内容。
- Type consistency：统一使用 `channel_provider`、`sync_fail_count`、`shipment_id`、`tracks` 命名。

## 交付检查清单

- [ ] 运单首次查询按主备路由并成功绑定渠道
- [ ] 已绑定运单后续刷新固定同一渠道
- [ ] 自动轮询支持后台开关与频率配置
- [ ] 同步保存最新状态与完整轨迹明细
- [ ] 前后台可查询轨迹，后台可手动刷新
- [ ] `docs-site` 完成功能/接口/配置影响同步

