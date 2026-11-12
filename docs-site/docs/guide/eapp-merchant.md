# 商家移动端 eapp（最新架构）

`eapp/` 是 lyshop 的商家专用移动端，覆盖 `H5 + 微信小程序 + App` 三端，承担「订单履约、商品运营、营销活动、客户与店铺管理」的成熟移动端商家工作台能力。

## 在线演示

`npm run dev:h5 -- --mode demo` 启动后即为 mock 演示模式，账号 admin / admin123 自动登录。生产模式连接后端 `/admin/api/*`。

## 顶层结构

- **4 个底部 Tab**：工作台 / 订单 / 客服 / 我的，各带常态与选中态线性 PNG 图标
- **自定义导航**：登录、工作台、订单、客服、我的、商品列表使用 `navigationStyle: custom`，由 `PageHeader` 统一渲染标题、辅助文案与右侧图标按钮
- **图表能力**：基于 `ly-charts (qiun-data-charts)` 跨端集成，业务方用 `ChartPanel + AreaChart / RingChart / BarChart` 统一调用
- **商家高阶组件库**：`components/biz/` 提供 PageHeader / SectionHeader / EntryList / MetricCard / ActionGrid / TodoCenter / AnnouncementBar / FilterDrawer / BatchBar / BatchResultPopup / Timeline / EmptyState / OrderCard / ProductCard / AfterSaleCard / SkuMatrixEditor / CategoryTreePicker / RichTextEditor / ShipPopup / RepricingPopup / RemarkPopup
- **通用组件**：`components/common/` 提供 AppIcon / SkeletonCard / StatusTag / StatCard / TodoCard
- **composables**：`useDashboard / useOrderList / useProductList / useBatchSelection / useFilter / useRequest`

## 视觉与交互规范

界面采用沉稳商务蓝 + 中性石板灰，所有颜色、圆角、阴影、字阶、间距都从 `App.vue` 的 CSS 变量取值，页面与组件不硬编码颜色。

| 令牌组 | 变量 | 用途 |
|--------|------|------|
| 品牌 | `--eapp-primary` / `--eapp-primary-deep` / `--eapp-primary-soft` | 主操作、品牌区、图标底 |
| 语义 | `--eapp-success` / `--eapp-warning` / `--eapp-danger` / `--eapp-accent` | 状态、趋势、预警（低饱和） |
| 中性 | `--eapp-bg` / `--eapp-card` / `--eapp-text*` / `--eapp-border*` | 背景、表面、文字、描边 |
| 形状 | `--eapp-radius*` / `--eapp-shadow*` / `--eapp-space-*` | 圆角、阴影、间距 |

- **uview-plus 桥接**：`--up-*` / `--u-*` 变量在 `App.vue` 统一定义，`up-*` 组件自动继承同一主题；`uni.scss` 只引入 `uview-plus/theme.scss`，不重复声明颜色。
- **图标系统**：`utils/icons.ts` 维护一套 24×24、1.7 描边的线性图标，运行时生成 base64 SVG data URI，由 `AppIcon.vue` 以 `<image>` 渲染，兼容 H5 / 微信小程序 / App；不使用 emoji，也不使用内联 `<svg>` 标签。
- **统一组件**：卡片用 1rpx 描边 + 极淡阴影；状态用带圆点的低饱和标签；区块标题用 `SectionHeader`（左侧主色竖条 + 可选「查看全部」）；入口列表用 `EntryList`；加载态用 `SkeletonCard`；点击反馈统一 `hover-class="eapp-hover"`。
- **图表色板**：`utils/ly-charts.ts` 的 `EAPP_CHART_COLORS` 与主题同源（品牌蓝 → 青 → 绿 → 琥珀 → 红 → 石板灰），坐标轴与网格线使用中性灰，避免高饱和彩虹配色。
  > 实现约束：本项目内置的 `uni_modules/ly-charts` 是 **26.1.4 旧副本**，其 line / bar / pie 组件**不读取** ECharts 顶层的 `option.color`，只按 `series.color`、`series.itemStyle.color`、饼图数据项的 `color` / `itemStyle.color` 取色，未指定时会回退到组件内置的 ECharts 高饱和默认色板。因此所有 option 构造函数除了设置顶层 `color`，还必须通过 `paletteColor(index)` 给每条 series、每个扇区显式着色，否则截图里仍会出现淡紫面积图和彩虹环形图。新增图表类型时同样处理。
  >
  > 上游组件库 `D:\Repos\xyito\open\lyCharts` 已修复该缺陷（Vue / uvue / uni-app x 全端支持顶层 `option.color`，取色优先级为显式色 → `itemStyle.color` → `option.color[index]` → 内置色板）。待本项目把 `uni_modules/ly-charts` 升级到含该修复的版本后，可保留顶层 `color` 并移除逐项显式着色；在升级之前不要删除 `paletteColor` 逻辑，也不要直接改项目内的第三方组件副本。
- **数字排版**：金额、指标、编号使用 `.eapp-num` 等宽数字对齐；`MetricCard` / `StatCard` 按数值字符数自适应字号，避免 393px 手机三列布局下长金额（如 `¥28960.50`）被省略号截断。

## 构建与静态资源

- `scripts/gen-tabbar-icons.cjs` 用纯 Node 光栅化生成 `static/tabbar/` 下 8 张 tabBar PNG（4 个入口 × 常态/选中态），修改图标后重新执行 `node scripts/gen-tabbar-icons.cjs`。
- `scripts/normalize-theme.cjs` 用于把页面里遗留的硬编码颜色批量收敛到设计令牌，新增页面后可复跑。
- `vite.config.ts` 的 `copyStaticAssets` 插件在 H5 构建时把 `static/` 复制到产物根目录：uni 的 Vite 构建默认不复制该目录，而 tabBar 的 `iconPath` 在 H5 端会解析为 `/static/...` 绝对路径，缺少这一步会导致图标 404。小程序与 App 端由 uni 自行处理静态资源。
- 部署 H5 时按 `manifest.json` 的 `h5.router.base`（默认 `/eapp/`）配置静态站点与 history 回退。

## 接口与配置

- 数据交互统一走 `/admin/api/*`，界面层不持有独立接口，mock 模式由 `VITE_MOCK=true` 切换。
- 运行与部署无需额外环境变量；构建产物包含 `static/tabbar/`，需随站点一起发布。

## 工作台（Dashboard）

页面按「店铺状态 → 核心指标 → 待办 → 公告 → 入口 → 分析 → 预警」组织，区块标题统一使用 `SectionHeader`。

- 顶部 `PageHeader`：店铺名（来自 `/shops/current`）+ 扫码 + 消息入口（带未读红点）
- 店铺概览条：店铺标识、营业状态、当前日期、经营报表入口
- 今日营业概览：今日营收 / 今日订单 / 客单价三列 `MetricCard`，附环比趋势
- 待办中心：待发货 / 待审售后 / 未读消息 / 库存预警四项，点击直达对应业务页并自动应用筛选
- 公告通知：`AnnouncementBar` 按通知 / 提醒 / 紧急分类，点击查看详情
- 快捷入口：订单 / 商品 / 售后 / 评价 / 装修 / 数据报表 / 仓储，统一线性图标
- 营销中心：优惠券 / 秒杀 / 拼团 / 砍价 / 会员 / 消息群发 / 装修 / 规格模板 / 分销 / 积分商城
- 经营分析：营收趋势（近 7 / 30 日 `AreaChart`）、订单状态分布（`RingChart`）、商品销量 Top5（`BarChart`）
- 库存预警列表：显示预警阈值与当前库存，点击跳转商品列表「预警」页签
- 扫码结果写入订单关键词，跳转到订单列表后自动作为搜索条件生效

## 订单

### 列表
- 顶部 7 状态 tabs：全部 / 待付款 / 待发货 / 已发货 / 已完成 / 已取消 / 售后中
- 快捷搜索：订单号 / 买家 / 商品，支持回车提交与一键清空；工作台扫码结果可直接作为关键词
- 高级筛选抽屉：关键词、金额区间、物流公司、收货省份、支付方式
- 卡片操作：详情 / 改价 / 催付 / 发货 / 备注 / 打单
- 勾选进入批量模式：批量发货 / 批量备注 / 批量关闭；结果弹窗分别列出成功与失败项
- 首次加载用 `SkeletonCard` 骨架屏，空态提供「清除筛选」快捷操作

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

营销中心按「活动与优惠 / 会员与运营 / 内容与工具」三组组织 10 个管理入口，统一使用 `EntryList`：

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

- 账号卡：头像、登录账号、认证状态
- 店铺管理：店铺设置 / 管理员 / 角色权限（按 RBAC 权限动态过滤）
- 客服与消息：自动回复 / 消息群发 / 客服会话（tabBar 页面用 `switchTab` 跳转）
- 经营工具：数据分析 / 店铺装修 / 仓储管理
- 退出登录带二次确认

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

- H5：主跑 demo，`env(safe-area-inset-top)` 在桌面自然降级为 0；`vite.config.ts` 负责把 `static/` 复制到产物
- 微信小程序：`AppIcon` 以 base64 SVG data URI 渲染 `<image>`，不依赖内联 `<svg>`；ly-charts 自带 canvas 兼容；富文本编辑器演示期只读
- App：面单预览用 `rich-text`/`web-view` 渲染
- 安全区：抽屉、批量操作条、页面底部留白统一用 `env(safe-area-inset-bottom)` 处理
- 旧机型：快捷入口 4 列、营销中心 5 列，公告按条展示不做横滑

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
