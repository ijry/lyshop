# eapp 商家移动端 P1 实施计划 — Mock 补齐 / 文档 / i18n（Part 3）

> 本文件接续 Part 1（A+B）与 Part 2（C+D+E）。最后 3 个 Phase 与 commit 6/7/8。

## Phase F — mock 路由与示例数据增量补齐（→ commit 6）

> 前序 Phase 已写入 25 条接口与初步数据。本阶段对**示例数据丰富度**做最后一次扫尾，让演示更真实。

### Task F1：扩订单到 40 单覆盖全状态

**Files:** 修改 `admin/src/mock/index.ts`

- [ ] **Step 1：生成订单数据**

在 `orderListSource` 定义之后（或文件顶部 source 区追加）插入扩展逻辑：

```ts
function expandOrderListTo40() {
  if (orderListSource.length >= 40) return
  const baseCount = orderListSource.length
  const need = 40 - baseCount
  const statuses = ['1', '2', '3', '4', '5']
  const provinces = ['上海市浦东新区', '北京市朝阳区', '广东省广州市', '浙江省杭州市', '四川省成都市', '湖北省武汉市', '辽宁省沈阳市']
  const companies = ['SF', 'ZTO', 'YTO', 'STO', 'YD', 'JD', 'EMS']
  const pickFromProducts = (i: number) => {
    const arr = (productListSource as any[]).slice(0, 6)
    return arr[i % arr.length]
  }
  for (let i = 0; i < need; i++) {
    const idx = baseCount + i + 1
    const id = 10000 + idx
    const status = statuses[i % statuses.length]
    const product = pickFromProducts(i)
    const qty = 1 + (i % 3)
    const price = Number(product?.price || 199)
    const goods = +(price * qty).toFixed(2)
    const discount = +((goods * (i % 3 === 0 ? 0.1 : 0)).toFixed(2))
    const pay = +(goods - discount).toFixed(2)
    const dateOffset = 86400000 * (i % 14)
    const created = new Date(Date.parse('2026-05-28T08:00:00Z') - dateOffset).toISOString()
    const order: any = {
      id,
      status,
      status_label: { '1': '待付款', '2': '待发货', '3': '已发货', '4': '已完成', '5': '已关闭' }[status],
      pay_method: ['wechat', 'alipay', 'wechat'][i % 3],
      user_nickname: ['张三', '李四', '王五', '赵六', '陈七', '何八'][i % 6],
      receiver_name: ['张先生', '李女士', '王先生', '赵女士'][i % 4],
      receiver_phone: '138****' + String(1000 + i).slice(-4),
      receiver_address: provinces[i % provinces.length] + ` ${i + 1}号`,
      goods_amount: goods,
      discount_amount: discount,
      pay_amount: pay,
      total_amount: pay,
      amount_breakdown: { goods_amount: goods, discount_amount: discount, payable_amount: pay },
      created_at: created,
      items: [{ id: idx * 100 + 1, product_id: product?.id, title: product?.title, cover: product?.cover || `https://picsum.photos/120/120?random=${1000 + i}`, qty, price }],
      shipments: status === '3' || status === '4' ? [{ id: idx * 10, company: companies[i % companies.length], tracking_no: `SF${1000000 + idx}`, delivery_type: 'express', logistics_status: status === '4' ? 'signed' : 'in_transit', logistics_status_label: status === '4' ? '已签收' : '运输中', created_at: created }] : [],
      logs: [],
      notes: [],
      has_after_sale: i % 7 === 0,
    }
    orderListSource.push(order)
  }
}
expandOrderListTo40()
```

放在文件顶部 helper 之后即可（无需修改 matchMock）。

### Task F2：补 5 条公告 + 协商消息 + 面单模板

- [ ] **Step 1：在 `announcementsSource` 中再加 2 条**（合计 5）：

```ts
announcementsSource.push(
  { id: 4, title: '物流时效升级', content: '部分城市顺丰提供「夜间寄」服务', type: 'normal', created_at: '2026-05-24T10:00:00Z' },
  { id: 5, title: '商家服务热线变更', content: '客服热线已升级 7×24，新号码 400-000-XXXX', type: 'urgent', created_at: '2026-05-23T18:00:00Z' },
)
```

- [ ] **Step 2：dev 验证 + commit 6**

```bash
cd eapp && npm run dev:h5 -- --mode demo
# 验证 dashboard 趋势图与销量榜数据丰满；订单列表至少 40 条；售后/分类样本齐备
cd eapp && npx vitest run
git -C 'D:\Repos\xyito\open\lyshop' add admin/src/mock/index.ts
git -C 'D:\Repos\xyito\open\lyshop' commit -m "eapp: mock 路由与示例数据增量补齐"
```

commit body：

```
- 订单 mock 数据扩展到 40 单，覆盖 5 个状态、含 1-3 个 SKU、含物流轨迹与跨期日期
- 公告补充至 5 条，售后协商示例完整、分类树预置 20 节点、6 商品多规格 SKU 完整
- 至此 P1 所有 mock 路由与示例数据可独立支撑完整演示
```

---

## Phase G — 文档同步（→ commit 7）

### Task G1：重写 docs-site/docs/guide/eapp-merchant.md

**Files:** 修改 `docs-site/docs/guide/eapp-merchant.md`

- [ ] **Step 1：重写文档（覆盖整文件）**

```markdown
# 商家移动端 eapp（最新架构）

`eapp/` 是 lyshop 的商家专用移动端，覆盖 `H5 + 微信小程序 + App` 三端，承担「订单履约、商品运营、营销活动、客户与店铺管理」的成熟移动端商家工作台能力。

## 在线演示

`npm run dev:h5 -- --mode demo` 启动后即为 mock 演示模式，账号 admin / admin123 自动登录。生产模式连接后端 `/admin/api/*`。

## 顶层结构

- **5 个底部 Tab**：工作台 / 订单 / 商品 / 营销 / 我的
- **图表能力**：基于 `ly-charts (qiun-data-charts)` 跨端集成，业务方用 `ChartPanel + AreaChart / RingChart / BarChart` 统一调用
- **商家高阶组件库**：`components/biz/` 提供 PageHeader / MetricCard / ActionGrid / TodoCenter / AnnouncementBar / FilterDrawer / BatchBar / BatchResultPopup / Timeline / EmptyState / OrderCard / ProductCard / AfterSaleCard / SkuMatrixEditor / CategoryTreePicker / RichTextEditor / ShipPopup / RepricingPopup / RemarkPopup
- **composables**：`useDashboard / useOrderList / useProductList / useBatchSelection / useFilter / useRequest`

## 工作台（Dashboard）

- 顶部：店铺名（来自 `/shops/current`）+ 扫码 + 消息入口
- 营业概览：今日营收 / 今日订单 / 客单价 + 同环比小标签
- 营收趋势：近 7 / 30 日切换的 `AreaChart`
- 订单状态分布：`RingChart` + 图例
- 公告通知带（横滑）
- 待办中心 6 项：待发货 / 待审售后 / 待回消息 / 库存预警 / 待开发票 / 待处理退款（点击直跳到对应业务页并自动应用筛选）
- 快捷九宫格 8 项：扫码发货 / 新建商品 / 优惠券 / 装修 / 数据 / 客户 / 财务 / WMS（后 4 项 P2 阶段实现）
- 商品销量 Top5 `BarChart`
- 库存预警列表

## 订单

### 列表
- 顶部 9 状态 tabs：全部 / 待付款 / 待发货 / 已发货 / 已完成 / 已关闭 / 售后中 / 待评价 / 待开票
- 高级筛选抽屉：关键词、金额区间、物流公司、收货省份、支付方式、是否含售后、时间
- 卡片操作菜单：详情 / 改价 / 催付 / 发货 / 备注 / 打单
- 长按或勾选进入批量模式：批量发货 / 批量备注 / 批量关闭

### 详情
- 顶部订单进度 Timeline（下单→支付→发货→签收→完成）
- 操作 ActionSheet（改价 / 备注 / 打单 / 催付）
- 物流卡片、商品明细、操作日志

### 批量操作集中页 `/pages/order/batch`
- 由列表「批量发货」或工作台「扫码发货」进入
- 支持统一选择快递公司，按行填写运单号
- 提交后展示成功/失败结果

### 面单预览 `/pages/order/print-preview`
- 服务端模板 HTML 通过 `rich-text` 渲染（演示模式 mock 静态模板）

## 售后

- 列表：6 状态 tabs + 售后类型筛选
- 详情：5 步进度条 / 协商沟通时间线（双向带头像气泡）/ 商家凭证上传 / 双向物流卡片 / 操作按钮按状态机启用

## 商品

### 列表
- 4 状态 tabs：在售 / 仓库 / 预警 / 全部
- 筛选：关键词 + 排序（销量 / 库存 / 价格 / 最新）
- 卡片：缩略图 + 多指标 + 状态徽标 + 操作
- 批量：批量上下架 / 批量分类 / 批量调价（百分比）
- FAB：新建 / AI 提示

### 编辑
- 分段表单：基础信息 / 主图轮播 / 详情（演示期只读 RichTextEditor）/ 价格库存 / 规格 SKU / 分类与标签 / 物流与营销 / 状态控制
- 多规格 SKU 矩阵编辑器：动态规格组 + 全组合矩阵 + 批量赋价/赋库存

### 分类管理 `/pages/product/category-tree`
- 三级树展开/折叠 + 改名 / 删除 / 增子

## 营销 / 我的

P1 保留现有实现（优惠券 / 秒杀 / 拼团 / 砍价 / 店铺设置 / 管理员 / 角色等），P2 之后会与商家工作台进一步集成。

## mock 演进

- 所有 mock 路由与示例数据集中维护在 `admin/src/mock/index.ts`，eapp 通过 `import.meta.env.VITE_MOCK==='true'` 触发 `matchMock`
- mock 写入型接口（改价、批量发货、协商消息、商家凭证、商品批量、分类 CRUD）会修改内存源数组，刷新后保留状态
- 批量接口返回 `success_ids[] / fail[{id,reason}]` 结构，含部分失败示例

## 平台兼容

- H5：主跑 demo，`env(safe-area-inset-top)` 在桌面自然降级为 0
- 微信小程序：ly-charts 自带 canvas 兼容；富文本编辑器演示期只读
- App：面单预览用 `rich-text`/`web-view` 渲染
- 旧机型：九宫格不超过 8 项、公告带最多 5 条

## 后续 Phase

- **P2**：商品分类树管理强化 / 规格模板 / 营销补会员价&分销&积分商城 / 移动端装修编辑器 / 优惠券深度
- **P3**：客户列表/标签/分群/会员卡 / 数据分析专题 / 群发推送 / 评价管理
- **P4**：账户/流水/对账/提现 / WMS 出入库/盘点/调拨 / IM 多坐席与工单 / 系统消息分级
```

### Task G2：追加接口文档

- [ ] **Step 1：在 `docs-site/docs/api/order.md` 末尾追加**

```markdown
## 商家移动端订单增强接口

> 用于商家移动端 eapp 的批量与高频操作，统一前缀 `/admin/api`。

### POST /orders/{id}/repricing

订单改价（仅 status=1 可用）。

请求：`{ items: [{ item_id, price }], remark }`
返回：`{ id, amount_breakdown: { goods_amount, discount_amount, payable_amount } }`

### POST /orders/{id}/notes

订单备注。

请求：`{ content, visible_to?: 'merchant_only' }`
返回：`{ id, notes: [...] }`

### POST /orders/{id}/remind-pay

催付款（短信 / 微信）。

请求：`{ channel: 'sms' | 'wx' }`
返回：`{ sent_at, channel }`

### GET /orders/{id}/print-template

电子面单模板。

返回：`{ template: '<html>...</html>' }`

### GET /orders/{id}/timeline

订单时间线。

返回：`[{ stage, status, time, content }]`

### POST /orders/batch/ship | notes | repricing | close

批量操作系列。`ship` 接收数组 `[{order_id,company,tracking_no}]`；其他统一 `{ ids[], ... }`。

返回：`{ success_ids: number[], fail: [{ id, reason }] }`

### GET /orders?keyword=&time_start=&time_end=&amount_min=&amount_max=&logistics_company=&province=&pay_method=&has_after_sale=

列表接口扩展查询参数。其它字段与原接口一致。
```

- [ ] **Step 2：在 `docs-site/docs/api/product.md` 末尾追加**

```markdown
## 商家移动端商品增强接口

### GET /categories/tree

返回 3 级分类树：`[{ id, parent_id, name, sort, product_count, children: [...] }]`

### POST /categories | PUT /categories/{id} | DELETE /categories/{id}

CRUD：POST `{ name, parent_id?, sort? }`、PUT `{ name?, sort? }`、DELETE 无 body。

### PUT /products/batch/status | category | price

批量上下架 / 批量分类 / 批量调价。

请求：
- status: `{ ids: number[], status: 0|1 }`
- category: `{ ids: number[], category_id }`
- price: `{ ids: number[], adjustment: { type: 'set'|'percent'|'amount', value, scope?: 'all'|'main_sku' } }`

返回：`{ success_ids[], fail: [{id, reason}] }`

### GET /products?status=&category_id=&sort_by=&low_stock=

列表扩展查询参数。`sort_by` 取值 `sales|stock|price_asc|price_desc|created`。
```

- [ ] **Step 3：在 `docs-site/docs/api/admin.md` 末尾追加**

```markdown
## 商家移动端工作台与基础接口

### GET /dashboard（升级）

在原有 today_* 与待办字段基础上，追加返回：

- `today_avg_price`
- `compare`：`{ revenue_yoy, revenue_mom, order_yoy, order_mom }` （0.18 表示 +18%）
- `trend`：`{ revenue_7d, revenue_30d, order_7d }`，每项为 `{ categories: string[], series: [{ name, data }] }`
- `status_distribution`：`[{ name, value }]`
- `hot_products`：`[{ id, title, cover, sold_qty }]`
- `announcements`：`[{ id, title, content, type, created_at }]`
- `stock_warning_list`：`[{ product_id, sku_id, title, stock, threshold }]`

### GET /shops/current

当前店铺：`{ id, name, logo, owner, decor_status }`

### GET /announcements

平台公告列表：标准 `{ list, total, page, size }` 分页结构。

### POST /after-sales/{id}/messages

售后协商消息：`{ from: 'merchant'|'user', content, images? }` → `{ id, messages: [...] }`

### POST /after-sales/{id}/evidences

商家凭证：`{ images: string[], remark? }` → `{ id, evidences: [...] }`
```

- [ ] **Step 4：在 `docs-site/docs/guide/features.md` 中「移动端商家工作台」段落（如无则新建一段）替换/插入**

```markdown
## 移动端商家工作台

`eapp/` 为商家专用移动端，覆盖 H5 / 微信小程序 / App，提供完整的订单履约、商品运营、营销与店铺管理：

- 工作台首屏：多指标卡 + 营收趋势图 + 订单状态环形 + 待办中心 + 9 状态订单 tabs + 销量榜 + 库存预警
- 订单：高级筛选抽屉 / 批量发货 / 改价 / 备注 / 面单预览 / 操作日志
- 售后：5 步进度条 + 协商时间线 + 凭证上传 + 双向物流
- 商品：多规格 SKU 矩阵 + 分类树 + 批量上下架 / 批量分类 / 批量调价 + 三级分类管理
- 图表：基于 `ly-charts` 跨端实现
- mock 模式：`npm run dev:h5 -- --mode demo` 即开即用
```

- [ ] **Step 5：commit 7**

```bash
git -C 'D:\Repos\xyito\open\lyshop' add docs-site
git -C 'D:\Repos\xyito\open\lyshop' commit -m "docs: 同步 eapp 商家端最新架构与接口"
```

commit body：

```
- 重写 docs-site/docs/guide/eapp-merchant.md 体现新工作台 / 订单 / 售后 / 商品全套
- 在 docs-site/docs/api/order.md、product.md、admin.md 追加新接口契约
- features.md 移动端商家工作台一节同步升级
- 文档使用中性表述，未出现第三方品牌名
```

---

## Phase H — 国际化补全 + 最终单测（→ commit 8）

### Task H1：补 i18n 键

**Files:** 修改 `eapp/locales/zh-CN.ts`、`eapp/locales/en.ts`

- [ ] **Step 1：补 zh-CN.ts**

```ts
// eapp/locales/zh-CN.ts 完整替换以下命名空间
export default {
  common: {
    refresh: '刷新', loading: '加载中…', empty: '暂无数据', save: '保存', cancel: '取消', confirm: '确认', reset: '重置',
  },
  login: {
    title: '商家工作台登录',
    username: '账号', password: '密码', submit: '登录',
    usernameRequired: '请输入账号', passwordRequired: '请输入密码',
  },
  dashboard: {
    title: '工作台',
    todayOrders: '今日订单', todaySales: '今日营收', todayAvgPrice: '客单价',
    pendingShip: '待发货', pendingAfterSale: '待审售后', unreadMessage: '未读消息',
    stockWarning: '库存预警', pendingInvoice: '待开发票', pendingRefund: '待处理退款',
    trend7: '近 7 日', trend30: '近 30 日',
    statusDistribution: '订单状态分布', hotProducts: '商品销量 Top5',
    quickEntries: '快捷入口', announcements: '公告',
  },
  order: {
    all: '全部', pendingPay: '待付款', pendingShip: '待发货', shipped: '已发货',
    completed: '已完成', closed: '已关闭', hasAfterSale: '售后中', pendingReview: '待评价', pendingInvoice: '待开票',
    filterTitle: '订单筛选',
    actions: { detail: '详情', reprice: '改价', remindPay: '催付', ship: '发货', note: '备注', print: '打单' },
    batch: { ship: '批量发货', notes: '批量备注', close: '批量关闭' },
  },
  afterSale: {
    all: '全部', applied: '待审核', returning: '退货中', refunding: '退款中', refunded: '已完成', closed: '已关闭',
    types: { all: '全部', refundOnly: '仅退款', returnRefund: '退货退款', exchange: '换货' },
    progress: { applied: '申请', approved: '审核', returning: '寄回', received: '收货', refunded: '退款' },
    chatPlaceholder: '回复买家',
    evidenceUpload: '上传凭证',
  },
  product: {
    search: '搜索商品', edit: '编辑', onSale: '上架', offSale: '下架',
    tabs: { all: '全部', onSale: '在售', off: '仓库', warning: '预警' },
    sortBy: { default: '默认', sales: '销量', stock: '库存', priceAsc: '价格升', priceDesc: '价格降', created: '最新' },
    batch: { shelfOn: '批量上架', shelfOff: '批量下架', category: '批量分类', price: '批量调价' },
    edit: {
      base: '基础信息', covers: '主图轮播', detail: '商品详情', pricing: '价格库存',
      sku: '规格 SKU', category: '分类与标签', shipping: '物流与营销', status: '状态控制',
    },
  },
  marketing: { coupon: '优惠券', seckill: '秒杀', groupBuy: '拼团', bargain: '砍价' },
  me: {
    title: '我的', messages: '消息中心', sessions: '客服会话',
    siteSettings: '店铺设置', admins: '管理员', roles: '角色权限', logout: '退出登录',
  },
  biz: {
    selected: '已选 {count} 项',
    batchResult: '批量操作结果',
    selectFirst: '请先勾选',
    soon: '即将上线',
  },
}
```

- [ ] **Step 2：英文版 `eapp/locales/en.ts` 同步键，文本简译即可（保持结构完全一致）**

```ts
export default {
  common: { refresh: 'Refresh', loading: 'Loading…', empty: 'No data', save: 'Save', cancel: 'Cancel', confirm: 'OK', reset: 'Reset' },
  login: { title: 'Merchant Console Login', username: 'Account', password: 'Password', submit: 'Sign In', usernameRequired: 'Account required', passwordRequired: 'Password required' },
  dashboard: { title: 'Dashboard', todayOrders: "Today's Orders", todaySales: "Today's Sales", todayAvgPrice: 'Avg Order', pendingShip: 'To Ship', pendingAfterSale: 'After-sale', unreadMessage: 'Messages', stockWarning: 'Low Stock', pendingInvoice: 'Invoice', pendingRefund: 'Refund', trend7: '7 days', trend30: '30 days', statusDistribution: 'Order Status', hotProducts: 'Top Products', quickEntries: 'Shortcuts', announcements: 'Notices' },
  order: { all: 'All', pendingPay: 'Unpaid', pendingShip: 'To Ship', shipped: 'Shipped', completed: 'Completed', closed: 'Closed', hasAfterSale: 'After-sale', pendingReview: 'Review', pendingInvoice: 'Invoice', filterTitle: 'Order Filter', actions: { detail: 'Detail', reprice: 'Reprice', remindPay: 'Remind', ship: 'Ship', note: 'Note', print: 'Print' }, batch: { ship: 'Batch Ship', notes: 'Batch Note', close: 'Batch Close' } },
  afterSale: { all: 'All', applied: 'Applied', returning: 'Returning', refunding: 'Refunding', refunded: 'Refunded', closed: 'Closed', types: { all: 'All', refundOnly: 'Refund only', returnRefund: 'Return & refund', exchange: 'Exchange' }, progress: { applied: 'Apply', approved: 'Review', returning: 'Return', received: 'Receive', refunded: 'Refund' }, chatPlaceholder: 'Reply…', evidenceUpload: 'Upload evidence' },
  product: { search: 'Search', edit: 'Edit', onSale: 'List', offSale: 'Unlist', tabs: { all: 'All', onSale: 'On Sale', off: 'Stock', warning: 'Low' }, sortBy: { default: 'Default', sales: 'Sales', stock: 'Stock', priceAsc: 'Price↑', priceDesc: 'Price↓', created: 'Newest' }, batch: { shelfOn: 'List', shelfOff: 'Unlist', category: 'Category', price: 'Price' }, edit: { base: 'Basics', covers: 'Covers', detail: 'Detail', pricing: 'Pricing', sku: 'SKU', category: 'Category', shipping: 'Shipping', status: 'Status' } },
  marketing: { coupon: 'Coupon', seckill: 'Seckill', groupBuy: 'Group Buy', bargain: 'Bargain' },
  me: { title: 'Me', messages: 'Messages', sessions: 'Chats', siteSettings: 'Shop', admins: 'Admins', roles: 'Roles', logout: 'Sign Out' },
  biz: { selected: 'Selected {count}', batchResult: 'Batch Result', selectFirst: 'Select first', soon: 'Coming soon' },
}
```

### Task H2：跑全部测试 + dev 烟囱测试

- [ ] **Step 1：跑全部 vitest**

```bash
cd eapp && npx vitest run
```

预期：全绿。任何失败按错误信息修复，不要继续往下走。

- [ ] **Step 2：dev:h5 演示模式人工烟囱测试**

```bash
cd eapp && npm run dev:h5 -- --mode demo
```

人工核对：
- 工作台首屏：店铺名、3 指标卡、营收图（7/30 切换）、环形图、公告带、6 待办（点击跳转）、9 快捷入口、销量榜、库存预警
- 订单：9 tabs、筛选抽屉、勾选批量发货可走完批量页流程
- 订单详情：进度时间线、操作 ActionSheet 4 项可触发
- 售后详情：进度条、协商可发送可显示、凭证上传后可见
- 商品列表：4 tabs、批量上下架/调价/分类、FAB
- 商品编辑：分段表单完整渲染、SKU 矩阵增删规格不丢数据
- 分类树：增删改可操作
- 控制台：无 ly-charts 与 uview-plus 报错

- [ ] **Step 3：commit 8**

```bash
git -C 'D:\Repos\xyito\open\lyshop' add eapp
git -C 'D:\Repos\xyito\open\lyshop' commit -m "eapp: 国际化补全与单测"
```

commit body：

```
- locales/zh-CN.ts 与 en.ts 全量补齐 dashboard / order / afterSale / product / biz 等命名空间
- 新增 useDashboard、useOrderList、useProductList、useBatchSelection、useFilter 单测；mock 路由 5 个 spec
- 烟囱测试覆盖工作台 / 订单 / 售后 / 商品全链路，演示模式可独立运行
```

---

## 完结自审

- [ ] **Step 1：扫一遍 spec 覆盖**

确认 spec `docs/superpowers/specs/2026-05-28-eapp-merchant-overhaul-design.md` 第 4 章 P1 各模块全部命中到本 plan 的 Task。

- [ ] **Step 2：扫一遍占位符**

确认本 plan 三个文件无 "TBD / TODO / 见上文 / 类似" 等占位（Task C4 Step 2 已用「完整代码」段落直接给出代码）。

- [ ] **Step 3：type 一致性**

确认 `useOrderList` / `useProductList` / `useDashboard` 暴露的 ref 字段在页面里访问方式一致（统一用 `h.list.value` / `h.filter.value` 等）。

- [ ] **Step 4：commit 顺序**

提交从 commit 1 到 commit 8 严格按 Phase A→H 顺序；每个 commit 都能让仓库处于「编译通过 + vitest 通过」的状态。

---

## 执行方式

`writing-plans` 推荐两种执行方式：

**1. Subagent-Driven（推荐）** — 我每个 Task 派一个独立 subagent 执行，两阶段 review，最适合本 plan 的颗粒度。

**2. Inline Execution** — 在当前会话里按批次跑 + 检查点。

请你告诉我选哪种，开始动手。
