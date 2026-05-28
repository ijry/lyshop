# 商家移动端 eapp 全面成熟化设计（最新架构）

## 1. 背景与目标

`eapp/` 已完成初始骨架（参见 `2026-05-28-eapp-merchant-mobile-design.md`），目前 5 个底部 Tab 与登录态、权限、mock 复用、i18n 都已就绪，但页面厚度普遍偏薄，难以与成熟移动端商家工作台正面对比。本设计在「不推翻已有架构」的前提下，将 `eapp` 升级为可演示、可上线的成熟商家移动端，分 4 个 Phase 落地。

本次目标：

1. 总体设计覆盖 P1-P4 全范围模块；首期实施 P1（核心交易闭环 + Dashboard + 商品 + 售后 + 高阶组件基础设施）。
2. 引入 `ly-charts` 图表能力（lyCharts，作者 ultraUI，跨 H5 / 微信小程序 / App / 鸿蒙）。
3. 沉淀一套商家专用高阶组件库（`components/biz/`），后续 Phase 直接复用。
4. mock 在 `admin/src/mock/index.ts` 内统一增量，eapp 不分裂 mock；写入型接口要修改内存源数据，刷新可见。
5. 文档同步 `docs-site/docs/guide/eapp-merchant.md` 与相关接口文档（写最新架构）。
6. 公开文档全部中性表述，**不出现第三方商城/SaaS 商标**；参考类备忘集中于项目根 `eapp-ui-reference-private.md`（不入库）。

---

## 2. 约束与边界

1. 工程：继续使用现有 `eapp/`，不另起工程。
2. 接口前缀：保持 `/admin/api/*`，不新增 `/merchant/*`。
3. 组件库：保留 `uview-plus`，新增的高阶组件统一在 `components/biz/`、图表在 `components/charts/`。
4. mock 路径：继续复用 `admin/src/mock/index.ts`，所有新增路由与示例数据写在此文件；批量写接口必须改动内存源数组以便刷新可见。
5. 后端真实接口：本期只交付接口契约，server 端实现按需评估；mock 模式下不阻塞演示。
6. 不破坏现有页面与接口语义，按 AGENTS 规则 3 优先升级现有接口而非新增。
7. **公开文档与提交信息禁出现第三方品牌名**；所有「对标」表述统一改为「成熟移动端商家工作台」「行业常规商家工作台」。

---

## 3. Phase 总览

| Phase | 主题 | 模块覆盖 | 演示价值 | 本次是否落地 |
|---|---|---|---|---|
| **P1** | 核心交易闭环 | Dashboard 工作台 / 订单（含批量&改价&打单&轨迹）/ 售后流程 / 商品（多规格 SKU + 分类 + 富文本占位 + 库存预警）/ ly-charts 基础设施 / `components/biz/` 高阶组件库 | 演示首屏震撼 + 完整交易闭环 | ✅ |
| **P2** | 商品+营销+店铺装修 | 商品分类树管理 / 规格模板 / 营销补会员价&分销&积分商城 / 移动端装修编辑器 / 优惠券深度（叠加规则、定向发券） | 运营特色与多元变现 | 仅写设计，不实施 |
| **P3** | 客户 + 数据分析 | 客户列表/标签/分群/会员卡 / 数据分析专题（销售/商品/客户/流量/转化）/ 群发推送 / 评价管理 | 数据驱动经营 | 仅写设计，不实施 |
| **P4** | 财务 + WMS + IM 工单 | 账户/流水/对账/提现 / WMS 出入库/盘点/调拨 / IM 客服多坐席与工单 / 系统消息分级 | 财务结算与履约后台 | 仅写设计，不实施 |

---

## 4. P1 详细模块设计

### 4.1 Dashboard 工作台（重写 `pages/dashboard/index.vue`）

页面节奏自上而下：

1. **顶部栏**（自定义状态栏 + 安全区）：店铺名（来自 `useShopStore`）+ 扫码（`uni.scanCode` 触发，识别到订单号跳详情、识别到商品 SKU 跳商品编辑）+ 消息（带未读徽标）。多店铺切换属 P2+ 范围，本期不出现入口。
2. **营业概览卡**：今日营收 / 今日订单 / 今日客单价，每项含同比/环比小标签（涨/跌色块）
3. **营收趋势图**：近 7 / 30 日切换的 `AreaChart`
4. **订单状态分布**：`RingChart`（待付/待发/已发/已完成/已关闭/售后中）+ 右侧文字列表
5. **待办中心**：6 个待办（待发货 / 待审核售后 / 待回复消息 / 库存预警 / 待开发票 / 退款待处理），点击进入对应业务页（带筛选参数）
6. **快捷九宫格**：扫描发货（`uni.scanCode` → 识别物流单号填入批量发货页 `batch.vue`）/ 新建商品（跳 `pages/product/edit?id=0`）/ 优惠券（跳 `pages/marketing/coupon`）/ 装修 / 数据 / 客户 / 财务 / WMS（后 4 项为 P2+ 占位，跳「即将上线」轻量提示页）
7. **公告通知带**：横滑展示 3 条平台公告
8. **商品销量排行 Top5**：迷你 `BarChart`

依赖接口：`GET /dashboard` 升级返回字段
```ts
{
  // 原有
  today_orders, today_sales, pending_ship, pending_after_sale, unread_message, stock_warning,
  // 新增
  today_avg_price: number,
  compare: { revenue_yoy, revenue_mom, order_yoy, order_mom },          // 同环比百分比
  trend: {
    revenue_7d: { categories: string[]; series: [{ name: '营收'; data: number[] }] },
    revenue_30d: { categories: string[]; series: [{ name: '营收'; data: number[] }] },
    order_7d:   { categories: string[]; series: [{ name: '订单'; data: number[] }] }
  },
  status_distribution: [{ name: string; value: number }],
  hot_products: [{ id, title, cover, sold_qty }],                       // Top5
  announcements: [{ id, title, content, type, created_at }],
  stock_warning_list: [{ product_id, sku_id, title, stock, threshold }]
}
```

### 4.2 订单模块

#### 4.2.1 列表页 `pages/order/list.vue`（重写）

- 顶部粘性 tabs：全部 / 待付款 / 待发货 / 已发货 / 已完成 / 已关闭 / 售后中 / 待评价 / 待开票（9 个，可横滑）
- 顶部搜索 + 「筛选」按钮，点开 `FilterDrawer`：
  - 时间范围（今日/昨日/近 7 日/近 30 日/自定义）
  - 金额区间（min/max）
  - 物流公司（多选）
  - 收货省份（picker）
  - 支付方式（微信/支付宝/余额）
  - 来源渠道（H5/小程序/PC）
  - 是否含售后
  - SKU 名称关键词
- 卡片采用 `OrderCard.vue`：商品缩略图（最多 3 张，超出 +N）+ 状态徽标 + 收货人 + 金额拆解 + 操作菜单（详情/改价/备注/打单/复制订单号/催发货）
- 长按或勾选进入批量模式：底部 `BatchBar` 显示「批量发货 / 批量备注 / 批量改价 / 批量打单 / 批量关闭」
- 下拉刷新 + 触底加载更多（接管 onPullDownRefresh / onReachBottom）

#### 4.2.2 详情页 `pages/order/detail.vue`（增强）

- 顶部 `Timeline`：下单 → 支付 → 发货 → 签收 → 评价（带时间戳，未发生节点灰显）
- 商品明细：缩略图 + SKU 属性 + 优惠拆分（每行可见 `original / discount / final`）
- 物流卡片改为 `Timeline` 视觉（合并 `getShipmentTracks`）
- 顶部操作菜单（ActionSheet）：发货 / 改价 / 修改地址 / 备注 / 关闭订单 / 打印面单 / 复制 / 申请发票
- Tab 切换「订单详情 / 操作日志」

#### 4.2.3 新增 `pages/order/batch.vue`

集中处理批量发货：
- 模式选择：批量发货 / 批量备注 / 批量改价 / 批量关闭
- 列表展示已选订单（来自上一页传参或本页扫码加入）
- 表单按模式动态切换；提交后展示 `BatchResultPopup`（成功 ids / 失败 ids[{id,reason}]）

#### 4.2.4 新增 `pages/order/print-preview.vue`

- mock 模式：`GET /orders/{id}/print-template` 返回 HTML 字符串，使用 `rich-text`（小程序）或 `web-view`（App/H5）渲染
- 顶部「发送至打印机」按钮（mock 直接 toast「已发送至默认打印机」）

#### 4.2.5 新接口与契约

```text
POST /orders/{id}/repricing          body: { items:[{item_id,price}], remark }       -> 返回更新后的 amount_breakdown
POST /orders/{id}/notes              body: { content, visible_to:'merchant_only' }    -> 推入 notes[]
POST /orders/{id}/remind-pay         body: { channel:'sms'|'wx' }                     -> { sent_at, channel }
GET  /orders/{id}/print-template                                                       -> { template:'<html>...</html>' }
GET  /orders/{id}/timeline                                                             -> [{ stage, status, time, content }]
POST /orders/batch/ship              body: [{order_id,company,tracking_no}]           -> { success_ids[], fail:[{id,reason}] }
POST /orders/batch/notes             body: { ids[], content }                          -> 同上风格
POST /orders/batch/repricing         body: { ids[], adjustment:{type:'percent'|'amount', value} } -> 同上风格
POST /orders/batch/close             body: { ids[], reason }                           -> 同上风格
GET  /orders?keyword=&time_start=&time_end=&amount_min=&amount_max=&logistics_company=&province=&pay_method=&has_after_sale=  // 在现有基础上扩展查询参数
```

### 4.3 售后模块

#### 4.3.1 列表 `pages/order/after-sale-list.vue`（增强）

- 顶部 tabs：全部 / 待审核 / 退货中 / 退款中 / 已完成 / 已关闭
- 时间筛选 + 售后类型筛选（仅退款/退货退款/换货）
- 卡片 `AfterSaleCard.vue`：售后号 + 类型 + 关联订单 + 状态 + 申请时间

#### 4.3.2 详情 `pages/order/after-sale-detail.vue`（增强）

- 顶部进度条（5 步：申请 → 审核 → 寄回 → 收货 → 退款，当前步骤高亮）
- 协商沟通时间线（商家/买家头像、文字、图片附件，输入框直接发消息）
- 上传凭证入口（选择本地图片，mock 阶段仅接收文件名）
- 物流卡片支持双向（用户寄回 / 商家补发），并接入轨迹
- 操作按钮按状态机启用/禁用（原有逻辑保留）

#### 4.3.3 新接口

```text
POST /after-sales/{id}/messages      body: { from:'merchant', content, images?[] }    -> 追加到 messages[]
POST /after-sales/{id}/evidences     body: { images:[url], remark? }                  -> 追加到 evidences[]
GET  /after-sales?type=&time_start=&time_end=                                         // 扩展查询参数
```

### 4.4 商品模块

#### 4.4.1 列表 `pages/product/list.vue`（重写）

- 顶部搜索 + 状态 tabs（在售 / 仓库 / 预警 / 全部）+ 分类筛选 + 排序（销量 / 库存 / 价格 / 创建时间）
- 卡片 `ProductCard.vue`：主图 + 标题 + 价格 + 库存 + 销量 + 状态徽标
- 批量模式：批量上下架 / 批量分类 / 批量删除 / 批量调价
- 右下角 FAB（单按钮 + 展开菜单）：「新建商品」直接跳 `edit.vue?id=0`；「AI 生成商品」为占位入口，弹层提示「请到管理后台 AI 工作流操作」并附复制链接

#### 4.4.2 编辑 `pages/product/edit.vue`（重写）

分段式表单（折叠/锚点跳转）：

1. **基础信息**：标题、副标题、卖点（多 tag）、商品编号
2. **主图轮播**：上传/排序（mock 阶段允许填 URL）
3. **详情**：`RichTextEditor.vue`（演示期只读 `rich-text`，编辑入口提示「请到管理后台编辑长详情」）
4. **价格库存**：主价 + 主库存 + 单位 + 重量 + 体积
5. **规格 SKU**：`SkuMatrixEditor.vue` —— 动态规格组（颜色、尺码…）+ SKU 矩阵（批量赋价 / 批量赋库存 / 单格编辑）
6. **分类与标签**：`CategoryTreePicker.vue` 多选 + 标签数组
7. **物流模板**：picker（默认/包邮/到付/同城）
8. **营销限制**：限购、不参与活动开关
9. **状态控制**：上下架、上架时间、下架时间

#### 4.4.3 新增 `pages/product/category-tree.vue`

- 三级分类树展示（手风琴）
- 增 / 改 / 删 / 上下移（mock 写入内存）

#### 4.4.4 新接口

```text
GET  /categories/tree                                          -> [{id,parent_id,name,path,product_count,children[]}]
POST /categories                       body: { name, parent_id, sort? }
PUT  /categories/{id}                  body: { name, sort? }
DELETE /categories/{id}
PUT  /products/batch/status            body: { ids[], status }
PUT  /products/batch/category          body: { ids[], category_id }
PUT  /products/batch/price             body: { ids[], adjustment:{type:'set'|'percent'|'amount', value, scope:'all'|'main_sku'} }
GET  /products?status=&category_id=&sort_by=&low_stock=          // 扩展查询参数
PUT  /products/{id}                    // payload 兼容 skus 数组与 multi-spec 编辑
```

### 4.5 ly-charts 基础设施 + 商家高阶组件库

#### 4.5.1 ly-charts 集成

- 通过 `uni_modules` 引入 `qiun-data-charts`（ly-charts 实际暴露的组件名），放入 `eapp/uni_modules/qiun-data-charts/`
- 跨端策略：H5 用 canvas 2D，小程序用原生 canvas-2d，App 用 cover-view + canvas
- 主题色由封装层注入，对齐 `--eapp-primary` (#2563EB) / `--eapp-accent` (#F97316) / `--eapp-success` (#16A34A) / `--eapp-warning` (#F59E0B)

#### 4.5.2 `components/charts/`

| 组件 | 作用 | 关键 props |
|---|---|---|
| `ChartPanel.vue` | 通用图表壳（标题/工具栏/loading/empty/error） | `title`, `type`, `data`, `height?`, `loading?`, `empty?`, `extra?` slot |
| `AreaChart.vue` | 折线/区域图预设 | `data`（categories + series）|
| `RingChart.vue` | 环形图预设 | `data`（[{name,value}]）|
| `BarChart.vue` | 柱状图预设（含水平） | `data`，`horizontal?` |

#### 4.5.3 `components/biz/`（商家高阶组件）

| 组件 | 作用 |
|---|---|
| `PageHeader.vue` | 自定义状态栏 + safe-area + 标题 + 右侧操作槽 |
| `MetricCard.vue` | 单/多指标卡（含同环比、点击跳转） |
| `ActionGrid.vue` | 九宫格快捷入口（icon + label + 跳转） |
| `TodoCenter.vue` | 待办中心列表 |
| `AnnouncementBar.vue` | 横滑公告通知带 |
| `FilterDrawer.vue` | 通用筛选抽屉（slot 自定义字段） |
| `BatchBar.vue` | 底部批量操作条（操作按钮组） |
| `BatchResultPopup.vue` | 批量操作结果展示（成功/失败列表） |
| `Timeline.vue` | 通用时间线（订单状态、操作日志、协商记录） |
| `EmptyState.vue` | 统一空态（图标 + 文案 + 按钮） |
| `OrderCard.vue` | 订单列表卡片 |
| `ProductCard.vue` | 商品列表卡片 |
| `AfterSaleCard.vue` | 售后列表卡片 |
| `SkuMatrixEditor.vue` | 规格矩阵编辑器（动态规格组 + 矩阵） |
| `CategoryTreePicker.vue` | 分类树多选弹层 |
| `RichTextEditor.vue` | 详情富文本（演示期只读） |
| `ShipPopup.vue` | 发货表单弹层（抽自现有 detail.vue） |
| `RepricingPopup.vue` | 改价表单弹层 |
| `RemarkPopup.vue` | 备注表单弹层 |

---

## 5. 架构与目录布局

```text
eapp/
├── pages/
│   ├── dashboard/index.vue           ★ 重写
│   ├── login/index.vue               • 保留
│   ├── order/
│   │   ├── list.vue                  ★ 重写
│   │   ├── detail.vue                ★ 增强
│   │   ├── batch.vue                 + 新增 批量集中页
│   │   ├── print-preview.vue         + 新增 面单预览
│   │   ├── after-sale-list.vue       ★ 增强
│   │   └── after-sale-detail.vue     ★ 增强
│   ├── product/
│   │   ├── list.vue                  ★ 重写
│   │   ├── edit.vue                  ★ 重写
│   │   └── category-tree.vue         + 新增 分类树管理
│   ├── marketing/...                 • 保留（P2 才动）
│   └── me/...                        • 保留（P2 才动）
├── components/
│   ├── charts/                       + ly-charts 封装
│   ├── biz/                          + 商家高阶组件库
│   ├── common/                       ★ StatusTag 扩展
│   └── layout/EappShell.vue          ★ 增强 safe-area
├── composables/
│   ├── useDashboard.ts               + 工作台数据聚合
│   ├── useOrderList.ts               + 列表分页/筛选/批量
│   ├── useBatchSelection.ts          + 通用批量勾选
│   ├── useFilter.ts                  + 筛选状态持久化（按 page key 存 storage）
│   └── useRequest.ts                 + loading/error 包装
├── api/
│   ├── dashboard.ts                  ★ 扩展返回字段
│   ├── order.ts                      ★ 增 batch/reprice/note/print/timeline/remind-pay
│   ├── product.ts                    ★ 增 batch/categories
│   ├── after-sale.ts                 + 抽出（原在 order.ts 内）
│   └── category.ts                   + 新增
├── stores/
│   ├── auth.ts                       • 保留
│   ├── badge.ts                      ★ 接入新待办字段
│   └── shop.ts                       + 当前店铺信息缓存（含切换店铺占位）
└── utils/
    └── ly-charts.ts                  + ly-charts 适配封装
```

**注**：`api/after-sale.ts` 从 `api/order.ts` 抽出，是为模块边界清晰，原 `order.ts` 中 `getAfterSales` 等导出保留兼容（re-export）。

### 5.1 数据流

```
页面 (pages/*.vue)
  ↓ 调 composable
composables/useXxx.ts (loading/error/分页/筛选/批量统一)
  ↓ 调 api
api/*.ts （只做 url 与方法）
  ↓ 调 utils/request
utils/request.ts
  ├─ MOCK=true → admin/src/mock/index.ts.matchMock()
  └─ MOCK=false → uni.request(BASE_URL + /admin/api/*)
```

`useOrderList` 暴露：`page/size/total/list/loading/refreshing/filter/selectedIds`、方法 `load/refresh/loadMore/toggleSelect/clearSelect/applyFilter`。同款思路用于 `useProductList`、`useAfterSaleList`。

### 5.2 国际化与权限

- 所有新文案进 `locales/zh-CN.ts` + `en.ts`，按模块前缀分组（`dashboard.*` / `order.*` / `product.*` / `afterSale.*` / `biz.*`）。
- 所有操作（批量发货、批量改价、面单、分类管理、商品批量上下架）通过 `hasPermission(authStore.perms, 'xxx.yyy')` 判断；mock 模式下 token 内 `perms='*'`，全部可见。

---

## 6. mock 演进策略（统一在 `admin/src/mock/index.ts`）

### 6.1 路由清单（新增/升级约 25 条）

| 接口 | 方法 | 说明 |
|---|---|---|
| `/dashboard` | GET ★升级 | 追加 `trend/status_distribution/hot_products/announcements/stock_warning_list/compare/today_avg_price` |
| `/categories/tree` | GET + | 三级分类树 |
| `/categories` | POST + | 新增分类 |
| `/categories/{id}` | PUT / DELETE + | 改/删分类 |
| `/orders` | GET ★扩展 | 新增多条 query 参数（见 4.2.5） |
| `/orders/{id}/repricing` | POST + | 改价 |
| `/orders/{id}/notes` | POST + | 备注 |
| `/orders/{id}/remind-pay` | POST + | 催付款 |
| `/orders/{id}/print-template` | GET + | 面单 HTML |
| `/orders/{id}/timeline` | GET + | 订单时间线 |
| `/orders/batch/ship` | POST + | 批量发货（含 partial fail） |
| `/orders/batch/notes` | POST + | 批量备注 |
| `/orders/batch/repricing` | POST + | 批量改价 |
| `/orders/batch/close` | POST + | 批量关闭 |
| `/after-sales` | GET ★扩展 | 扩展 type / time 筛选 |
| `/after-sales/{id}/messages` | POST + | 协商消息 |
| `/after-sales/{id}/evidences` | POST + | 上传凭证 |
| `/products` | GET ★扩展 | status / category_id / sort_by / low_stock 筛选 |
| `/products/batch/status` | PUT + | 批量上下架 |
| `/products/batch/category` | PUT + | 批量分类 |
| `/products/batch/price` | PUT + | 批量调价 |
| `/products/{id}` | PUT ★兼容 | 兼容多规格 SKU 矩阵 |
| `/shops/current` | GET + | 当前店铺 |
| `/announcements` | GET + | 平台公告（3 条） |
| `/print/templates` | GET + | 面单模板列表（占位） |

### 6.2 mock 数据扩充

为达到「演示看起来真实」，需补：

1. **订单**：mock orders 从 ~10 单扩到 40 单，覆盖所有 9 个状态、每单 1-3 个 SKU、含售后/物流轨迹/操作日志。
2. **商品**：现有商品列表中挑 6 个补全多规格 SKU 矩阵（颜色 × 尺码 × 版本），并对每个商品补 `category_id`、`category_path`、`category_path_name`、`sales_count`、`low_stock_threshold`。
3. **分类树**：3 级共 ~20 节点（电子/服饰/生活/食品/家居/美妆…）。
4. **dashboard 趋势**：30 日营收按确定性公式生成（用日期作 seed，正弦 + 抖动 + 周末微涨），保证刷新数据稳定。
5. **公告**：3 条（系统升级、政策提醒、营销活动预告）。
6. **协商消息**：示例售后单中预置 5 条来回沟通。
7. **打印模板**：1 段简洁面单 HTML 字符串（无外部资源依赖）。

### 6.3 mock 实现注意

- `matchMock(method, url, data)` 要支持路径参数（如 `/orders/{id}/repricing`）。如现有实现仅做严格匹配，实施阶段需要扩展为正则/参数提取。
- 写入型接口（如 `POST /orders/{id}/repricing`、`POST /orders/batch/ship`、`POST /after-sales/{id}/messages`）必须**修改 mock 内存中的 source 数组**，使刷新后能看到效果（参考现有 `shipOrder` / `auditAfterSale` 逻辑）。
- 批量接口必须返回部分失败示例（如 ids 中第 3 个总是失败，提示「订单已关闭」），便于演示错误提示 UI。
- 30 日趋势用 `seedRandom(date)` 生成，避免每次刷新数据跳变。

### 6.4 真实接口同步策略（非 mock 模式）

- 新增接口要求 server 端实现的功能，本期只规划契约。若 server 端实现成本过高，对应功能在非 mock 模式下显示「即将上线」占位（按钮 disable + tooltip），不阻塞 mock 演示。
- 已有接口尽量「升级」（按 AGENTS 规则 3），如 `/dashboard` 在原响应基础上追加字段，老消费者不受影响。
- `docs-site/docs/api/order.md` / `product.md` / `admin.md` 同步刷新（按 AGENTS 规则 2、4）。

---

## 7. 兼容性、错误处理、测试与文档

### 7.1 平台与版本兼容

| 平台 | 要求 | 风险点 | 应对 |
|---|---|---|---|
| H5 | 主跑 demo | 自定义状态栏需根据浏览器，`env(safe-area-inset-top)` 在桌面 = 0 自然降级 | `PageHeader` 用 `env()`，无须分支判断 |
| 微信小程序 | 支持 | ly-charts 自带 canvas 兼容；富文本用 `editor` 需在 `pages.json` 注册，演示期只读 | `RichTextEditor` 演示期只读 `rich-text` |
| App | 兼容 | 批量打印面单 WebView 渲染 | `print-preview.vue` 用 `web-view` 显示模板 HTML |
| 旧机型 | iOS 12+ / Android 8+ | `scroll-view scroll-x` 横滑性能 | 九宫格不超过 8 项；公告带最多 5 条 |

### 7.2 错误处理与降级

- 现状 `utils/request.ts` 已统一 401 + `code !== 0` 提示，保留。新增：
  - 网络错误：toast「网络异常，已切换离线展示数据」+ 列表显示空态卡片（不抛错）
  - 批量操作：成功部分 toast「成功 X / 失败 Y」，失败列表展示在 `BatchResultPopup`
  - 图表无数据：`ChartPanel` 内置 `EmptyState`（图标 + 文案 + 重试按钮）
  - 富文本/打印模板加载失败：占位卡片「内容暂不可用」
- 表单弹层（改价/备注/发货/批量改价）：
  - 提交前本地校验（非空 / 数值范围 / SKU 选择）
  - 提交时按钮 `loading`，避免重复点击
  - 网络失败保留弹层与已填字段

### 7.3 测试范围（vitest）

- **单元测试新增**：
  - `composables/useOrderList.test.ts`：分页、筛选合并、批量选择幂等
  - `composables/useBatchSelection.test.ts`：全选/反选/最大选中限制
  - `components/biz/SkuMatrixEditor.test.ts`：规格组动态变化下 SKU 矩阵的稳态（增删规格不丢已编辑数据）
  - `components/charts/ChartPanel.test.ts`：props 透传 + 空态分支
- **mock 路由测试新增**：
  - `admin/src/mock/index.test.ts`：校验新增 25 条路由 matchMock 命中并返回结构正确（路径参数、批量返回 partial fail）
- **现有测试**：保持通过；若新字段影响订单详情既有断言，同步更新。

### 7.4 文档同步（按 AGENTS 规则 2、4）

- `docs-site/docs/guide/eapp-merchant.md`：主指南重写「页面总览 / 操作流程 / 演示路径」
- `docs-site/docs/api/order.md` / `product.md` / `admin.md`：追加新接口契约（请求 / 响应 / 示例），写最新架构语气
- `docs-site/docs/guide/features.md`：「移动端商家工作台」一节升级
- 公开文档**不出现第三方品牌名**；如需「灵感参考」，集中写到项目根 `eapp-ui-reference-private.md`，并在 `.gitignore` 加入（不入库）

### 7.5 Git 提交划分（按 AGENTS 规则 1）

P1 实施阶段按以下划分 commit（每次都中文，head + body，不带 Co-Authored-By）：

1. `eapp: 引入 ly-charts 与商家高阶组件骨架`
2. `eapp: 重写工作台 dashboard 与待办中心`
3. `eapp: 订单列表与详情升级，新增批量与改价流程`
4. `eapp: 售后协商时间线与凭证上传`
5. `eapp: 商品列表与编辑（多规格 SKU + 分类）`
6. `eapp: mock 路由与示例数据增量`
7. `docs: 同步 eapp 商家端最新架构与接口`
8. `eapp: 国际化补全与单测`

### 7.6 范围外（明确不做）

- P2-P4 模块（装修编辑器、客户 CRM、数据分析专题、财务、WMS、IM 工单、群发推送）：仅写入本 spec 第 3 节，不在 P1 实施
- 真实后端接口的服务端实现：本期 spec 给出契约，server 端实现由后续 PR 评估（mock 模式下不阻塞）
- 富文本编辑器的完整 WYSIWYG：演示阶段只读，编辑跳 admin

---

## 8. 验收标准（P1）

1. 在 `npm run dev:h5` mock 模式下，5 个底部 Tab 主流程闭环：
   - 工作台首屏含 2 个图表 + 9 个快捷入口 + 6 个待办 + 公告带 + 销量榜
   - 订单列表 9 状态 tabs + 高级筛选 + 批量发货/改价/备注/打单/关闭可演示
   - 售后详情含协商时间线 + 凭证上传 + 双向物流
   - 商品编辑支持多规格 SKU 矩阵 + 分类树选择 + 富文本占位
2. mock 接口刷新数据稳定（同一日期产生一致趋势）；批量接口返回 partial fail 可触发结果弹层
3. `vitest run` 全绿；新增测试覆盖 composables 与 mock 路由
4. `docs-site` 三处文档（guide/eapp-merchant、api/order/product/admin、guide/features）已更新且不出现第三方品牌名
5. `eapp-ui-reference-private.md`（如存在）已加入 `.gitignore`
6. 微信小程序 / App 端首屏渲染无明显报错（不要求三端 100% 等价，但 dashboard / 订单列表 / 商品列表三页可正常打开）

---

## 9. 风险与对接

- **ly-charts 集成体积**：包体 ~78KB，对 H5 影响可控；小程序需注意分包，必要时将 `pages/dashboard` 移入主包，其他低频图表页移入子包。
- **mock 路径参数支持**：若 `matchMock` 当前不支持，实施前先小幅重构匹配函数（独立 commit），避免大改逻辑。
- **server 端跟进**：批量接口、改价、面单模板涉及后端实现工作量，实施阶段需与 server 端约定 Phase 2 真实接口对接节奏。
- **demo 与生产差异**：mock 模式下所有批量、改价、协商均「乐观提交」，生产模式需要后端真实校验，弹层 loading 已为此预留。

---

## 10. 里程碑

- **M1（本期）**：P1 全套交付 + mock 演示路径打通 + 文档同步
- **M2**：进入 P2 实施（先 brainstorm 二级 spec，再 writing-plans）
- **M3-M4**：P3 / P4 同模式推进
