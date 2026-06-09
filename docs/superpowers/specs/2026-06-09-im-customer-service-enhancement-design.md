# IM 客服系统渐进增强设计

## 背景

LYShop 已具备完整的 IM/AI 客服主链路：用户会话、WebSocket 实时消息、坐席接入、排队、转接、AI 首接、RAG 知识库、商品信息注入和 AI 反馈评估。对比独立客服系统 AI-CS 后，当前最值得补齐的是运营可观测性、多实例实时消息一致性，以及客服文件/图片消息体验。

本设计采用渐进增强路线，在现有 `im` 插件内扩展能力，不引入独立客服子系统，不改变现有用户、管理员、权限、插件和商城业务边界。

## 目标

第一期实现以下能力：

- IM 报表：统计会话、消息、AI 回复、AI 失败、RAG 命中、转人工、人工接入、关闭、转接、文件消息等指标。
- IM 事件日志：记录客服系统关键事件，支持后台检索和问题排查。
- Redis WebSocket 广播：支持多后端实例下的实时消息跨实例投递。
- 文件/图片消息：用户端和客服端可上传图片/文件，并在会话中展示。

第一期不实现匿名访客小窗。访客体系涉及匿名身份、跨站嵌入、安全域、会话归并和商城用户转化，不纳入本阶段。

## 设计原则

- 优先扩展现有接口和模型，避免引入重复的 conversation/visitor 子系统。
- `ImSession`、`ImMessage`、`ImStaff` 继续作为客服主模型。
- 事件日志用于统计和审计，业务状态仍以现有会话和消息表为准。
- WebSocket 保持当前进程内 Hub 作为本地投递核心，Redis 只负责跨实例转发。
- 文件消息复用现有 storage driver，不新增单独客服存储配置。
- docs-site 文档直接描述最新架构和接口。

## 后端架构

### 模块划分

- `server/plugins/im/model`：新增事件日志模型，复用 `ImMessage.Extra` 承载附件元信息。
- `server/plugins/im/service`：新增日志记录、统计聚合、上传处理和 Redis 广播适配。
- `server/plugins/im/api`：在现有前台和后台 IM API 上增加上传、报表和日志接口。
- `server/plugins/im/service/hub.go`：扩展现有 Hub，增加可选 Redis 发布/订阅，不改变调用方的 `GlobalHub.Send` 用法。
- `docs-site/docs/api/im.md`：同步新增接口、消息格式、配置和部署影响。

### 事件日志

新增 `ImEventLog`：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | uint64 | 主键 |
| `event` | string | 事件编码 |
| `session_id` | uint64 | 会话 ID |
| `user_id` | uint64 | 用户 ID |
| `staff_id` | uint64 | 客服管理员 ID |
| `message_id` | uint64 | 关联消息 ID |
| `source` | string | `user`、`staff`、`ai`、`system` |
| `success` | int8 | 1 成功，0 失败 |
| `latency_ms` | int64 | 耗时毫秒 |
| `extra` | json/text | 扩展信息 |
| `created_at` | time | 创建时间 |

第一期事件编码：

- `session_created`
- `message_sent`
- `ai_reply`
- `ai_failed`
- `rag_hit`
- `to_human`
- `staff_accept`
- `session_close`
- `session_transfer`
- `file_uploaded`

事件记录应由业务流程顺手写入，失败不阻断主流程。写日志失败只记录服务端 logger。

### 报表统计

新增后台接口：

`GET /admin/api/im/analytics`

查询参数：

- `from`：开始日期，格式 `YYYY-MM-DD`。
- `to`：结束日期，格式 `YYYY-MM-DD`。
- `staff_id`：可选，按客服筛选。

响应包含：

- `summary.sessions`：新增会话数。
- `summary.messages`：消息数。
- `summary.ai_replies`：AI 回复数。
- `summary.ai_failed`：AI 失败数。
- `summary.rag_hits`：RAG 命中数。
- `summary.to_human`：转人工次数。
- `summary.accepts`：人工接入次数。
- `summary.closes`：关闭会话次数。
- `summary.transfers`：转接次数。
- `summary.files`：文件消息数。
- `trend`：按日聚合的上述核心指标。

统计以 `ImEventLog` 聚合为主。会话列表和消息列表仍使用现有接口，避免报表接口承担明细查询职责。

### 日志查询

新增后台接口：

`GET /admin/api/im/logs`

查询参数：

- `event`：事件编码。
- `session_id`：会话 ID。
- `user_id`：用户 ID。
- `staff_id`：客服 ID。
- `source`：事件来源。
- `success`：成功状态。
- `page`：页码。
- `size`：每页条数。

返回 `response.PageData`，列表按 `id desc` 排序。

### 文件/图片消息

新增接口：

- `POST /api/v1/im/upload`
- `POST /admin/api/im/upload`

请求为 multipart form：

- `file`：必填。
- `session_id`：必填。

校验规则：

- 文件最大 10MB。
- 支持图片：`jpg`、`jpeg`、`png`、`gif`、`webp`。
- 支持普通文件：`pdf`、`doc`、`docx`、`xls`、`xlsx`、`txt`、`csv`、`md`、`zip`。
- 通过扩展名和 MIME 做基础校验，避免任意可执行文件上传。
- 用户端只能上传自己会话的文件；后台客服端需具备 `im:reply` 权限。

上传成功后返回：

```json
{
  "url": "/uploads/...",
  "path": "im/...",
  "name": "example.png",
  "size": 12345,
  "mime": "image/png",
  "message_type": "image"
}
```

消息发送仍走现有 HTTP/WS 发送逻辑。附件元信息写入 `ImMessage.Extra`：

```json
{
  "file_url": "/uploads/...",
  "file_path": "im/...",
  "file_name": "example.png",
  "file_size": 12345,
  "mime": "image/png"
}
```

`ImMessage.Type` 使用现有类型：

- 图片：`image`
- 普通文件：新增约定 `file`

新增 `file` 类型是必要的，因为现有 `text`、`image`、`product_card`、`order_card`、`system` 不能表达普通附件语义。

### Redis WebSocket 广播

现有 `GlobalHub.Send(targetID, data)` 保持不变。

新增内部广播信封：

```json
{
  "node_id": "hostname-pid-random",
  "target_id": "user_1001",
  "data": "{...frame json...}",
  "created_at": 1780992000000
}
```

机制：

- 本机调用 `Send` 时先进入本地 Hub 投递。
- Redis 可用时，同时发布到 channel `lyshop:im:ws`。
- 每个实例启动时订阅该 channel。
- 收到远端消息时，如果 `node_id` 等于本机，忽略，避免回环。
- 远端消息只投递本机连接，不再二次发布。

Redis 客户端复用 `core/cache.Client`。当前配置未连接外部 Redis 时，单实例仍可正常运行。多实例部署必须配置同一个外部 Redis；嵌入式 Redis 只适合单实例开发环境。

## 前端设计

### Admin

增强 `admin/src/views/im/SessionList.vue`：

- 消息区渲染 `image` 类型为缩略图，点击预览。
- 渲染 `file` 类型为文件卡片，点击打开文件 URL。
- 输入区增加上传按钮。
- 发送附件消息时先调用后台上传接口，再发送 `image` 或 `file` 消息。

新增管理页面：

- `admin/src/views/im/Analytics.vue`：展示 IM 统计汇总和按日趋势。
- `admin/src/views/im/EventLogs.vue`：展示 IM 事件日志列表和筛选条件。

菜单挂在现有「客服中心」下：

- `客服报表`：权限 `im:view`。
- `事件日志`：权限 `im:view`。

### App

增强 `app/pages/im/chat.vue`：

- 输入区增加上传按钮。
- 图片消息显示缩略图，点击预览。
- 文件消息显示文件名、大小和下载/打开入口。
- 继续使用现有 uview-plus 组件，优先使用 `u-*` 或 `up-*` 组件，不回退到 `uni-*` 组件，除非无对应能力。

### Eapp

增强 `eapp/pages/im/chat.vue`：

- 客服端可上传并发送图片/文件。
- 消息气泡组件支持 `image` 和 `file` 类型。

### Web

增强 `web/src/views/Chat.vue`、`web/src/components/ChatDialog.vue` 和 `web/src/stores/chat.ts`：

- PC/H5 Web 客服弹窗支持附件发送。
- 图片消息展示缩略图，文件消息展示链接卡片。

## 权限

沿用现有权限：

- `im:view`：查看报表和事件日志。
- `im:reply`：客服上传并发送附件。

用户端上传要求登录用户身份，且只能作用于自己的会话。

## 错误处理

- 事件日志写入失败不影响主业务。
- 统计查询失败返回标准错误响应。
- 上传文件过大、类型不允许、会话无权限时返回 400 或 403。
- Redis 发布失败只记录日志，不影响本机 WebSocket 投递。
- Redis 订阅断开后应自动重试；重试期间本机实时消息仍可用。

## 测试策略

后端测试：

- `ImEventLog` 事件写入成功。
- 报表接口按事件聚合指标正确。
- 日志查询筛选和分页正确。
- 上传接口拒绝超大文件和不允许类型。
- 用户不能上传到他人会话。
- Redis 广播忽略本机 `node_id`，远端消息只本地投递不回环。

前端测试：

- Admin 会话页可渲染文本、图片、文件消息。
- App/Eapp/Web 发送附件后消息列表展示正确。
- 上传失败时保留输入状态并提示错误。

文档测试：

- `docs-site/docs/api/im.md` 覆盖新增接口、消息类型、Redis 部署说明和统计口径。

## 部署与配置

单实例部署无需新增配置。

多实例部署需要所有后端实例连接同一个外部 Redis。`config.example.yaml` 已包含 Redis 配置示例，文档需明确：嵌入式 Redis 不能用于多实例 WebSocket 广播。

文件存储继续使用现有 storage 插件体系。Docker Compose 已挂载 `./data/uploads:/app/uploads`，本地存储模式无需新增 volume。

## 非目标

- 不实现匿名访客小窗。
- 不引入 AI-CS 的独立用户、conversation、visitor、system_logs 全套模型。
- 不替换现有 Qdrant RAG pipeline。
- 不新增独立客服服务进程。
- 不做客服组织/租户隔离。
