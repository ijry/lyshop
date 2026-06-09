# IM 客服系统功能矩阵

本文档记录 LYShop IM 客服系统在四端（Admin、Eapp、App、Web）的完整功能实现情况。

## 功能对比表

| 功能 | Admin | Eapp | App | Web | 说明 |
|---|:---:|:---:|:---:|:---:|---|
| **基础功能** |
| 消息收发 | ✅ | ✅ | ✅ | ✅ | 实时文本消息 |
| 图片/文件消息 | ✅ | ✅ | ✅ | ✅ | 上传、发送、预览图片，文件以卡片展示 |
| WebSocket | ✅ | ✅ | ✅ | ✅ | 实时通信，心跳保活，自动重连 |
| 跨实例 WebSocket | ✅ | ✅ | ✅ | ✅ | 后端通过 Redis Pub/Sub 扇出到其他副本 |
| 实时输入草稿 | ✅ | ✅ | - | ✅ | `typing_draft`/`typing_stop` 实时转发，客服可看到用户当前输入 |
| 访客上下文 | ✅ | - | - | ✅ | 会话记录访客 ID、IP、来源页面、浏览器、系统、语言、设备等 |
| Web 嵌入脚本 | - | - | - | ✅ | `/im-widget.js` 以 iframe 方式嵌入现有 `/chat?embed=1` |
| 系统消息 | ✅ | ✅ | ✅ | ✅ | 接入/转接/结束通知 |
| **排队机制** |
| 排队位置显示 | ✅ | ✅ | ✅ | ✅ | 显示"排队第N位" |
| 排队状态更新 | ✅ | ✅ | ✅ | ✅ | WebSocket 实时推送 |
| **客服操作** |
| 接入按钮 | ✅ | ✅ | - | - | 手动接入等待会话 |
| 结束按钮 | ✅ | ✅ | - | - | 结束当前会话 |
| 转接按钮 | ✅ | ✅ | - | - | 转接给其他客服 |
| 转接通知 | ✅ | ✅ | ✅ | ✅ | 用户收到转接系统消息 |
| **客服状态** |
| 在线/离线切换 | ✅ | ✅ | - | - | 客服手动上线/下线 |
| 客服负载显示 | ✅ | ✅ | - | - | 当前接待数/最大负载 |
| 连接状态指示 | ✅ | ✅ | ✅ | ✅ | WebSocket 连接状态 |
| **管理功能** |
| 坐席管理 | ✅ | ✅ | - | - | 客服人员 CRUD |
| 自动回复规则 | ✅ | ✅ | - | - | 关键词匹配自动回复 |
| 会话列表 | ✅ | ✅ | - | - | 查看所有会话 |
| 会话历史 | ✅ | ✅ | ✅ | ✅ | 消息历史记录 |
| 客服报表 | ✅ | - | - | - | 会话、消息、AI、RAG、转人工、附件等统计 |
| 事件日志 | ✅ | - | - | - | 查询会话、消息、AI、转接、上传等事件 |
| **AI 智能客服** |
| 本地大模型应答 | ✅ | - | ✅ | ✅ | 新会话默认由本地大模型 AI 接待 |
| RAG 知识库召回 | ✅ | - | ✅ | ✅ | 向量召回，无向量模型时退化为关键词召回 |
| 商品信息分析 | ✅ | - | ✅ | ✅ | 回答时检索在售商品价格/库存/销量 |
| 可选联网搜索 | ✅ | - | ✅ | ✅ | 默认关闭；配置 Serper 后注入外部搜索摘要并记录日志 |
| 转人工 | ✅ | - | ✅ | ✅ | 输入“人工”关键词或点击转人工按钮进入排队 |
| 知识库管理 | ✅ | - | - | - | 知识条目 CRUD + 重建向量索引 + 连通性测试 |
| 文档切片入库 | ✅ | - | - | - | 上传 TXT/MD/CSV/JSON/XML/HTML/DOCX/PDF/XLSX 自动切片为多条知识 |
| 大模型配置 | ✅ | - | - | - | 配置中心维护服务地址/模型/提示词等 |
| AI 反馈统计 | ✅ | - | ✅ | ✅ | 用户端提交反馈，后台查看列表和统计 |

> AI 智能客服为可选能力，由插件配置 `ai_enabled` 控制。关闭时新会话回退到传统“分配/排队”人工流程。

## 端侧能力边界

| 端 | 主要角色 | 已覆盖能力 | 边界说明 |
|---|---|---|---|
| Admin | 平台/商家后台客服与管理员 | 会话处理、坐席管理、知识库、报表、日志、反馈统计、附件发送 | 需要按角色授予 `im:*` 权限；报表和日志只在 Admin 提供 |
| Eapp | 商家移动端客服 | 会话列表、接入/回复/转接/结束、附件发送、连接状态 | 面向移动客服处理，不承载知识库和报表管理 |
| App | 用户移动端 | 创建会话、AI 接待、转人工、附件消息、反馈 | 附件上传仅允许当前用户自己的会话 |
| Web | PC/H5 用户端 | 站内聊天弹窗、嵌入 iframe、AI 接待、转人工、附件消息、反馈、草稿同步 | 推荐使用弹窗或 `/im-widget.js` 保持购物链路，不建议跳转独立客服页 |

## 消息类型与展示

| 消息类型 | 说明 | 端侧展示 |
|---|---|---|
| `text` | 普通文本消息 | 四端文本气泡 |
| `image` | 图片附件 | 四端预览图片，点击查看或打开原图 |
| `file` | 普通文件附件 | 文件卡片，展示文件名、大小和下载入口 |
| `product_card` | 商品卡片 | 预留业务卡片类型，可携带商品扩展信息 |
| `order_card` | 订单卡片 | 预留业务卡片类型，可携带订单扩展信息 |
| `system` | 系统通知 | 接入、转接、关闭、排队等状态提示 |

附件消息的 `extra` 使用 `{file_url,file_path,file_name,file_size,mime}`。前端展示应优先使用 `file_url`，缺失时再按业务配置拼接 `file_path`。

## 权限说明

| 权限 | 说明 | 默认角色 |
|---|---|---|
| `im:view` | 查看客服会话和消息 | 客服、管理员 |
| `im:reply` | 回复消息、接入/结束/转接会话 | 客服、管理员 |
| `im:staff:manage` | 管理客服坐席（增删改查） | 管理员 |
| `im:knowledge` | 管理 AI 知识库、重建索引、测试大模型连通 | 管理员 |

超级管理员默认拥有所有权限。

## 技术架构

### 后端
- **语言**: Go
- **WebSocket**: gorilla/websocket
- **消息总线**: Hub 模式（单播/广播）
- **跨实例投递**: Redis Pub/Sub（频道 `lyshop:im:ws`，节点 ID 防回环）
- **附件存储**: 复用系统 storage driver，本地 Docker 使用 `./data/uploads:/app/uploads`
- **事件统计**: `im_event_logs` 统一记录会话、消息、AI、转人工、转接、上传事件，报表按事件聚合生成
- **日志中心**: `im_event_logs` 包含 `level/category/trace_id/message/meta`，支持事件检索、关键字检索和 AI 联网搜索诊断
- **数据库表**:
  - `im_sessions` - 会话表（用户、客服、**接待模式 mode(ai/human)**、状态、排队位置、访客上下文）
  - `im_messages` - 消息表（会话ID、发送者类型含 **AI=3**、内容）
  - `im_event_logs` - 事件日志表（会话、消息、AI、转人工、上传等审计和报表来源）
  - `im_staffs` - 客服表（管理员ID、在线状态、负载）
  - `im_transfer_logs` - 转接记录表
  - `im_auto_replies` - 自动回复规则表
  - `im_knowledges` - AI 知识库表（标题、内容、标签、向量 embedding、是否已索引、状态）
  - `im_feedbacks` - 用户反馈与 LLM-as-Judge 自动评估结果

### AI 客服（本地大模型 + RAG）
- **接口协议**: OpenAI 兼容 `/chat/completions` 与 `/embeddings`，可对接本地推理服务（如 Ollama / vLLM）。
- **检索增强**:
  - 知识库召回（按优先级）：① Qdrant 向量库 ANN 检索（配置 `ai_qdrant_url` + 向量模型，`status=1` 过滤、可选相似度阈值，按命中 ID 回查 DB 并保序，可扩展到大规模知识库）；② 仅向量模型时全量内存余弦（适合小库）；③ 均未配置时标题/内容/标签关键词召回（含 CJK 二元切分）兜底。
  - 向量数据同步：知识 CRUD / 文档导入异步 upsert 到 Qdrant，删除同步删点，停用经 payload `status` 失效；`reindex` 重建集合并全量重灌；DB `embedding` 列作本地缓存/回退。
  - 商品信息分析：按用户问题检索在售商品（标题/副标题 LIKE），注入价格、库存、销量供模型参考。
- **混合检索（Hybrid + RRF）**：`ai_hybrid=on` 时向量召回与关键词召回并行，结果经 Reciprocal Rank Fusion（k=60）融合，召回长尾更稳。
- **重排（Rerank）**：配置 `ai_rerank_url` 后，召回候选池（RecallK）送 cross-encoder 精排，支持 Cohere / Jina / TEI 兼容 `/rerank` 接口；不配置则保持召回顺序直接取 TopK。
- **查询改写**：`ai_query_rewrite` 可选 `rewrite`（LLM 扩写口语化问题）、`hyde`（生成假设回答作为检索向量）、`multi`（生成 N 个变体各自检索再 RRF 融合）；改写仅作用于检索，不影响发给 LLM 的原始问题。
- **联网搜索**：`ai_web_search_enabled` 默认关闭；开启后通过 `ai_web_search_provider=serper`、`ai_web_search_api_key`、`ai_web_search_endpoint`、`ai_web_search_top_k` 获取搜索摘要并注入 `【联网搜索】` 上下文，失败只记录 `web_search` 事件，不阻断回答。
- **评估闭环**：用户可在聊天页对 AI 回答👍/👎，结果存入 `ImFeedback`（source=user）；开启 `ai_auto_eval` 后 AI 回答完成后自动用 LLM-as-Judge 评估忠实度和相关性（0-5，source=auto）；管理后台可查看列表和聚合统计（`/admin/api/im/feedback`）。
- **向量库部署**: `docker-compose.yml` 内置 `qdrant` 服务（`qdrant/qdrant`），容器内地址 `http://qdrant:6333`，默认 Qdrant collection 为 `im_knowledge`。
- **配置项**（配置中心 → IM客服）：`ai_enabled`、`ai_base_url`、`ai_api_key`、`ai_chat_model`、`ai_embed_model`、`ai_system_prompt`、`ai_human_keywords`、`ai_top_k`、`ai_temperature`、`ai_product_search`、`ai_timeout_sec`、`ai_qdrant_url`、`ai_qdrant_api_key`、`ai_qdrant_collection`、`ai_score_threshold`、`ai_hybrid`、`ai_recall_k`、`ai_rerank_url`、`ai_rerank_api_key`、`ai_rerank_model`。

### 事件与报表口径

| 指标 | 事件 | 说明 |
|---|---|---|
| `sessions` | `session_created` | 新建会话 |
| `messages` | `message_sent` | 保存消息 |
| `ai_replies` | `ai_reply` | AI 回复成功 |
| `ai_failed` | `ai_failed` | AI 回复失败 |
| `rag_hits` | `rag_hit` | 知识库或商品上下文命中 |
| `to_human` | `to_human` | 用户请求转人工 |
| `accepts` | `staff_accept` | 客服接入 |
| `closes` | `session_close` | 关闭会话 |
| `transfers` | `session_transfer` | 转接会话 |
| `files` | `file_uploaded` | 上传附件 |

报表查询支持 `from`、`to`、`staff_id`。日期按本地时区解析，`to` 包含当天；事件日志支持按 `event`、`session_id`、`user_id`、`staff_id`、`source`、`success` 过滤。

### 前端
- **Admin**: Vue 3 + TypeScript
- **Eapp**: UniApp (Vue 3)
- **App**: UniApp (Vue 3)
- **Web**: Vue 3 + TypeScript

### WebSocket 帧类型

| 类型 | 方向 | 说明 |
|---|---|---|
| `msg` | 双向 | 消息内容（`sender_type`：0系统/1用户/2人工客服/3AI） |
| `typing` | 服务端→客户端 | AI 正在生成回复的输入指示 |
| `typing_draft` | 双向转发 | 实时输入草稿，payload 含 `draft/sender_type/sender_id/updated_at` |
| `typing_stop` | 双向转发 | 输入停止或消息已发送 |
| `to_human` | 客户端→服务端 | 用户点击“转人工”，请求转接人工 |
| `queue` | 服务端→客户端 | 排队位置更新 |
| `assign` | 服务端→客户端 | 接入/转接通知 |
| `close` | 服务端→客户端 | 会话结束通知 |
| `ping/pong` | 双向 | 心跳保活 |

## 主要功能流程

### 1. 用户发起咨询
1. 用户打开聊天页面，后端 `GetOrCreateSession` 创建会话
2. **AI 已启用**：会话以 `mode=ai`、`status=ONGOING` 创建，由本地大模型直接接待，无需排队
   - 用户每条消息：先落库 → 命中转人工关键词则转人工，否则推送 `typing` 并异步生成 RAG 回复（知识库 + 商品信息）→ 推送 `sender_type=3` 的 AI 消息
3. **AI 未启用**：会话以 `mode=human`、`status=WAITING` 创建，走传统人工流程
   - 查找在线且未满载的客服：有 → 直接分配（ONGOING）；无 → 进入排队（queue_position: N）并推送排队位置

### 1.1 转人工
1. 用户输入“人工/转人工/...”关键词，或点击“转人工”按钮（发送 `to_human` 帧）
2. 后端 `SwitchToHuman`：写入系统消息、将 `mode` 置为 `human`
3. 有空闲客服 → 直接分配并推送 `assign:accepted`；否则置 `WAITING` 入队并推送 `queue` 位置

### 2. 客服接入
1. 客服点击"接入"按钮
2. 更新会话状态为 ONGOING
3. 客服负载 +1
4. 推送接入通知给用户
5. 清除用户排队位置

### 3. 会话转接
1. 客服A点击"转接"按钮，选择客服B
2. 更新 session.staff_id = B
3. 客服A负载 -1，客服B负载 +1
4. 插入系统消息到会话
5. 推送通知：
   - 客服A: `transfer_out`（从列表移除）
   - 客服B: `transfer_in`（添加到列表）
   - 用户: `transfer`（显示转接通知）

### 4. 结束会话
1. 客服点击"结束"按钮
2. 更新会话状态为 CLOSED
3. 客服负载 -1
4. 推送结束通知给用户
5. 从队列分配下一个等待会话

### 5. 附件消息
1. 用户或客服选择图片/文件
2. 调用对应上传接口，后端校验文件大小、扩展名、MIME 和会话权限
3. 上传服务复用系统 storage driver，返回 `url/path/name/size/mime/message_type`
4. 前端发送 `msg` 帧或调用回复接口，消息类型使用 `image` 或 `file`
5. 服务端保存消息并推送给会话双方，同时记录 `file_uploaded` 事件

### 6. 多副本实时投递
1. 当前副本优先向本机连接投递 WebSocket 消息
2. 同一消息发布到 Redis Pub/Sub 频道 `lyshop:im:ws`
3. 其他副本收到广播后按会话和客户端 ID 投递给本机连接
4. 节点 ID 用于忽略自身广播，避免回环重复投递

如果未连接外部 Redis，系统仍可单实例工作；多副本场景会出现用户和客服连接到不同副本时无法实时互通的问题。

### 7. Web 嵌入与访客上下文
1. 宿主站引入 `/im-widget.js` 并调用 `LYShopIMWidget.init({baseUrl, token, context})`
2. 脚本生成右下角按钮和 iframe，打开 `/chat?embed=1`
3. iframe 通过 `postMessage` 接收访客上下文，调用 `/api/v1/im/session` 时写入 `ImSession`
4. 客服接入或查看会话时，可看到访客来源、页面、语言、浏览器、设备等信息

## 文件位置

### 后端
- `server/plugins/im/model/im.go` - 数据模型
- `server/plugins/im/service/session.go` - 业务逻辑（会话/排队/转接/转人工/WS）
- `server/plugins/im/service/ai.go` - 本地大模型调用、RAG 召回、商品信息分析
- `server/plugins/im/service/web_search.go` - 可选联网搜索 Provider
- `server/plugins/im/service/vectorstore.go` - Qdrant 向量库客户端（REST）
- `server/plugins/im/service/rerank.go` - 混合检索 RRF 融合 + cross-encoder 重排客户端
- `server/plugins/im/service/query_rewrite.go` - 查询改写（rewrite/HyDE/multi）
- `server/plugins/im/service/eval.go` - LLM-as-Judge 评估 + 用户反馈 CRUD
- `server/plugins/im/service/knowledge.go` - 知识库 CRUD、文档导入切片
- `server/plugins/im/service/document.go` - 多格式文本提取与切片
- `server/plugins/im/service/hub.go` - WebSocket Hub
- `server/plugins/im/service/event.go` - 事件日志与报表聚合
- `server/plugins/im/service/upload.go` - 附件校验与上传
- `server/plugins/im/api/admin.go` - Admin API（含知识库/AI 测试）
- `server/plugins/im/api/front.go` - 用户端 API
- `server/plugins/im/plugin.json` - 插件配置（菜单/权限/config_items）

### 前端
- `admin/src/views/im/SessionList.vue` - 客服会话页面
- `admin/src/views/im/Analytics.vue` - 客服报表页面
- `admin/src/views/im/EventLogs.vue` - 事件日志页面
- `admin/src/views/im/StaffManage.vue` - 坐席管理页面
- `admin/src/views/im/KnowledgeManage.vue` - AI 知识库管理页面
- `eapp/pages/me/im-sessions.vue` - 会话列表
- `eapp/pages/im/chat.vue` - 聊天页面
- `eapp/pages/im/staff-manage.vue` - 坐席管理
- `app/pages/im/chat.vue` - 用户聊天页面（AI 气泡 + 转人工）
- `web/src/views/Chat.vue` - Web 聊天页面
- `web/src/components/ChatDialog.vue` - 聊天弹窗
- `web/src/stores/chat.ts` - 聊天状态管理
- `web/public/im-widget.js` - 可嵌入站点的客服 iframe 脚本

## 部署影响

- 多副本后端必须连接同一个外部 Redis，才能保证用户和客服 WebSocket 分布在不同副本时仍可互通。
- 嵌入式 Redis 只适合单实例运行。
- 附件消息使用当前启用的 storage driver；本地存储部署需确保 uploads 目录持久化。

## 上线检查清单

- 数据库迁移已执行，IM 相关表存在且插件已启用。
- 后台角色已配置 `im:view`、`im:reply`；管理员额外配置 `im:staff:manage`、`im:knowledge`。
- 至少创建一个客服坐席，并设置合理 `max_load`。
- 用户端、Web 端和客服端都能访问同一后端域名或网关下的 `/ws/im`。
- 上传目录或对象存储已持久化，返回 URL 可被四端访问。
- 多副本部署使用外部 Redis，并确认所有副本连接同一个 Redis。
- AI 客服启用前已验证 `ai_base_url` 与 `ai_chat_model`；需要 RAG 时验证 embedding/Qdrant；需要精排时验证 rerank 服务。
- 联网搜索启用前已配置 `ai_web_search_api_key`，并确认后端可以访问 `ai_web_search_endpoint`。
- 前端重连后会补拉历史消息，避免 WebSocket 断线期间漏展示。
