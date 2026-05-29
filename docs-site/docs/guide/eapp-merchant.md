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

## 营销

营销模块覆盖 9 个管理入口：

- **优惠券**：全字段 CRUD + 叠加规则（互斥 / 同类可叠 / 跨类可叠）+ 定向发券（全部用户 / 会员等级 / 新用户）
- **秒杀**：活动 + 活动商品管理
- **拼团**：活动 + 活动商品管理
- **砍价**：活动 + 活动商品管理（起砍价/底价模式）
- **VIP 会员**：套餐 / 等级 / 领券规则 / SKU 专属价 4 页 CRUD
- **店铺装修**：简化版可视化编辑器，8 种组件类型，卡片列表排序 + 属性编辑弹窗 + 多副本管理
- **规格模板**：预定义属性组 + 应用分类，可在商品编辑中一键应用生成 SKU 笛卡尔积
- **分销管理**：分销商管理 / 佣金管理 / 分销配置
- **积分商城**：积分商品 CRUD / 兑换记录 / 签到管理入口

### VIP 会员管理

| 页面 | 功能 |
|------|------|
| 会员套餐 | 月卡/季卡等订阅套餐，设置时长、价格、状态 |
| 会员等级 | 成长等级体系（银卡/金卡），设置成长值门槛与折扣率 |
| 领券规则 | 按等级配置每月自动领券，关联优惠券与限领次数 |
| SKU专属价 | 按商品 SKU + 等级设置专属价格，支持筛选 |

### 移动端装修编辑器

轻量化方案，无拖拽、无 iframe 预览：
- 8 种组件：轮播图 / 分类导航 / 商品网格 / 公告 / 图片广告 / 富文本 / 营销区 / 间距
- 操作：添加组件 → 上下箭头排序 → 点击编辑属性 → 保存/发布
- 副本管理：复制 / 重命名 / 删除

### 优惠券深度

- 全字段 CRUD：名称、类型（满减/折扣/立减）、面额、有效期、发行量等
- 叠加规则 `stack_rule`：`exclusive`（互斥）/ `same_type`（同类可叠）/ `cross_type`（跨类可叠）
- 定向发券 `target_type`：`all`（全部）/ `vip_level`（会员等级）/ `new_user`（新用户）
- 发券操作：选择数量后确认发放，自动累加 `used_count`

### 规格模板

- 模板包含：名称、适用分类（多选）、属性组（动态增删，每组含名称 + 逗号分隔值）
- 在商品编辑页 SKU 区域点击「应用规格模板」→ 选择模板 → 自动生成笛卡尔积填充 SKU 矩阵

## 我的

P1 保留现有实现（店铺设置 / 管理员 / 角色等），P2 之后会与商家工作台进一步集成。

## 仓储管理（WMS）

仓储管理模块提供仓库、库存、出入库单与库存流水的移动端管理能力：

| 页面 | 路径 | 功能 |
|------|------|------|
| 仓储首页 | `/pages/wms/index` | 4 个功能入口卡片 |
| 仓库管理 | `/pages/wms/warehouse-list` | 仓库 CRUD（名称/编码/地址/联系人/状态） |
| 库存台账 | `/pages/wms/stock-ledger` | 仓库筛选 + SKU 搜索 + 安全库存编辑 + 预警标识 |
| 出入库单 | `/pages/wms/doc-list` | 类型 tabs（全部/入库/出库）+ 状态过滤 |
| 单据编辑 | `/pages/wms/doc-editor` | 单据创建/编辑/完成/取消，商品明细行管理 |
| 库存流水 | `/pages/wms/movement-list` | 仓库/SKU/单据号多维搜索，变动量着色 |

## IM 客服

IM 模块提供客服会话管理、实时聊天与自动回复能力：

| 页面 | 路径 | 功能 |
|------|------|------|
| 客服会话 | `/pages/me/im-sessions` | 未读计数、最后消息、状态标签（等待/进行/关闭），点击进入聊天 |
| 聊天详情 | `/pages/im/chat` | 消息气泡（买家/客服），底部输入栏，自动滚底 |
| 自动回复 | `/pages/im/auto-replies` | 规则 CRUD，精确/包含匹配，启用/禁用 |

## 消息分级

消息中心支持按分组和优先级管理：

- **分组 tabs**：全部 / 系统 / 订单 / 营销 / 客服
- **优先级指示**：紧急（红点）/ 重要（琥珀点）/ 普通（无点）
- **标记已读**：单条标记已读操作

## mock 演进

- 所有 mock 路由与示例数据集中维护在 `admin/src/mock/index.ts`，eapp 通过 `import.meta.env.VITE_MOCK==='true'` 触发 `matchMock`
- mock 写入型接口（改价、批量发货、协商消息、商家凭证、商品批量、分类 CRUD）会修改内存源数组，刷新后保留状态
- 批量接口返回 `success_ids[] / fail[{id,reason}]` 结构，含部分失败示例

## 平台兼容

- H5：主跑 demo，`env(safe-area-inset-top)` 在桌面自然降级为 0
- 微信小程序：ly-charts 自带 canvas 兼容；富文本编辑器演示期只读
- App：面单预览用 `rich-text`/`web-view` 渲染
- 旧机型：九宫格不超过 8 项、公告带最多 5 条

## 分销管理

分销管理模块提供二级分销体系的移动端管理能力：

| 页面 | 路径 | 功能 |
|------|------|------|
| 分销总览 | `/pages/distribution/index` | 3 指标卡（总分销商/待结算佣金/已结算总额）+ 3 入口 |
| 分销商管理 | `/pages/distribution/distributor-list` | 关键词搜索 + 状态 tabs（全部/启用/禁用），启用/禁用切换 |
| 佣金管理 | `/pages/distribution/commission-list` | 状态 tabs（全部/待结算/已结算/已退回），结算/退回操作 |
| 分销配置 | `/pages/distribution/config` | 一级/二级佣金比例设置 |

## 积分商城

积分商城模块提供积分商品、兑换记录与签到管理的移动端能力：

| 页面 | 路径 | 功能 |
|------|------|------|
| 积分总览 | `/pages/points/index` | 3 指标卡（发放积分/消耗积分/商品数）+ 3 入口（商品/兑换/签到） |
| 积分商品 | `/pages/points/product-list` | 类型 tabs（全部/实物/虚拟/优惠券）+ CRUD 弹窗 |
| 兑换记录 | `/pages/points/exchange-list` | 状态 tabs（全部/待发货/已完成），发货/完成操作 |

## 签到管理

签到管理依托已有 mock 路由，提供签到规则与日志的移动端管理：

| 页面 | 路径 | 功能 |
|------|------|------|
| 签到规则 | `/pages/checkin/rules` | 动态增删规则行（天数+积分），保存提交 |
| 签到记录 | `/pages/checkin/logs` | 只读列表：用户/日期/连续天数/获得积分 |

## 后续 Phase

- **P5**：客户列表/标签/分群/会员卡 / 群发推送 / 账户/流水/对账/提现 / 盘点/调拨

## 评价管理

- 列表页 `/pages/review/list`：3 状态 tabs（全部/待回复/已回复）+ 关键词搜索
- 评价卡片：商品缩略图 + 星级评分 + 内容摘要 + 图片行 + 回复状态标签 + 追评计数
- 详情弹窗：完整内容/图片（点击预览）/追评列表/已有回复
- 回复弹窗：文本输入 + 提交，支持新增和修改回复
- mock 层：`reply_status` 过滤（pending/replied），`POST /reviews/:id/reply` 写入

## 数据分析

5 个专题页，统一使用 `useAnalytics` composable + `DateRangePicker` 日期切换：

| 页面 | 路径 | 指标卡 | 图表 |
|------|------|--------|------|
| 销售分析 | `/pages/analytics/sales` | GMV / 订单数 / 客单价 | 营收趋势(AreaChart) + 订单趋势(AreaChart) + 支付方式(RingChart) + 时段分布(BarChart) |
| 商品分析 | `/pages/analytics/products` | SKU总数 / 动销率 / 库存周转 | SKU排行(BarChart horizontal) + 分类销售(RingChart) + 价格段(RingChart) + 库存状态(RingChart) |
| 客户分析 | `/pages/analytics/customers` | 总客户 / 复购率 / 客均消费 | 新老客占比(RingChart) + 客单价分布(BarChart) + 购买频次(BarChart) + RFM雷达(RadarChart) |
| 流量分析 | `/pages/analytics/traffic` | PV / UV / 平均停留 / 跳出率 | PV/UV趋势(AreaChart) + 渠道分布(RingChart) + 设备分布(RingChart) + 页面停留(BarChart) |
| 转化分析 | `/pages/analytics/conversion` | 总转化率 / 弃购率 | 转化漏斗(FunnelChart) + 步骤趋势(AreaChart) + 总转化率仪表盘(GaugeChart) |

入口：工作台快捷入口「数据报表」→ 分析首页 5 卡片 → 各专题页
