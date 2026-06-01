# IM 客服系统功能矩阵

本文档记录 LYShop IM 客服系统在四端（Admin、Eapp、App、Web）的完整功能实现情况。

## 功能对比表

| 功能 | Admin | Eapp | App | Web | 说明 |
|---|:---:|:---:|:---:|:---:|---|
| **基础功能** |
| 消息收发 | ✅ | ✅ | ✅ | ✅ | 实时文本消息 |
| WebSocket | ✅ | ✅ | ✅ | ✅ | 实时通信，心跳保活，自动重连 |
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
| **AI 智能客服** |
| 本地大模型应答 | ✅ | - | ✅ | ✅ | 新会话默认由本地大模型 AI 接待 |
| RAG 知识库召回 | ✅ | - | ✅ | ✅ | 向量召回，无向量模型时退化为关键词召回 |
| 商品信息分析 | ✅ | - | ✅ | ✅ | 回答时检索在售商品价格/库存/销量 |
| 转人工 | ✅ | - | ✅ | ✅ | 输入“人工”关键词或点击转人工按钮进入排队 |
| 知识库管理 | ✅ | - | - | - | 知识条目 CRUD + 重建向量索引 + 连通性测试 |
| 文档切片入库 | ✅ | - | - | - | 上传 TXT/MD/CSV/JSON/XML/HTML/DOCX/PDF/XLSX 自动切片为多条知识 |
| 大模型配置 | ✅ | - | - | - | 配置中心维护服务地址/模型/提示词等 |

> AI 智能客服为可选能力，由插件配置 `ai_enabled` 控制。关闭时新会话回退到传统“分配/排队”人工流程。

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
- **数据库表**:
  - `im_session` - 会话表（用户、客服、**接待模式 mode(ai/human)**、状态、排队位置）
  - `im_message` - 消息表（会话ID、发送者类型含 **AI=3**、内容）
  - `im_staff` - 客服表（管理员ID、在线状态、负载）
  - `im_transfer_log` - 转接记录表
  - `im_auto_reply` - 自动回复规则表
  - `im_knowledge` - AI 知识库表（标题、内容、标签、向量 embedding、是否已索引、状态）

### AI 客服（本地大模型 + RAG）
- **接口协议**: OpenAI 兼容 `/chat/completions` 与 `/embeddings`，可对接本地推理服务（如 Ollama / vLLM）。
- **检索增强**:
  - 知识库召回（按优先级）：① Qdrant 向量库 ANN 检索（配置 `ai_qdrant_url` + 向量模型，`status=1` 过滤、可选相似度阈值，按命中 ID 回查 DB 并保序，可扩展到大规模知识库）；② 仅向量模型时全量内存余弦（适合小库）；③ 均未配置时标题/内容/标签关键词召回（含 CJK 二元切分）兜底。
  - 向量数据同步：知识 CRUD / 文档导入异步 upsert 到 Qdrant，删除同步删点，停用经 payload `status` 失效；`reindex` 重建集合并全量重灌；DB `embedding` 列作本地缓存/回退。
  - 商品信息分析：按用户问题检索在售商品（标题/副标题 LIKE），注入价格、库存、销量供模型参考。
- **向量库部署**: `docker-compose.yml` 内置 `qdrant` 服务（`qdrant/qdrant`），容器内地址 `http://qdrant:6333`。
- **配置项**（配置中心 → IM客服）：`ai_enabled`、`ai_base_url`、`ai_api_key`、`ai_chat_model`、`ai_embed_model`、`ai_system_prompt`、`ai_human_keywords`、`ai_top_k`、`ai_temperature`、`ai_product_search`、`ai_timeout_sec`、`ai_qdrant_url`、`ai_qdrant_api_key`、`ai_qdrant_collection`、`ai_score_threshold`。

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

## 文件位置

### 后端
- `server/plugins/im/model/im.go` - 数据模型
- `server/plugins/im/service/session.go` - 业务逻辑（会话/排队/转接/转人工/WS）
- `server/plugins/im/service/ai.go` - 本地大模型调用、RAG 召回、商品信息分析
- `server/plugins/im/service/vectorstore.go` - Qdrant 向量库客户端（REST）
- `server/plugins/im/service/knowledge.go` - 知识库 CRUD、文档导入切片
- `server/plugins/im/service/document.go` - 多格式文本提取与切片
- `server/plugins/im/service/hub.go` - WebSocket Hub
- `server/plugins/im/api/admin.go` - Admin API（含知识库/AI 测试）
- `server/plugins/im/api/front.go` - 用户端 API
- `server/plugins/im/plugin.json` - 插件配置（菜单/权限/config_items）

### 前端
- `admin/src/views/im/SessionList.vue` - 客服会话页面
- `admin/src/views/im/StaffManage.vue` - 坐席管理页面
- `admin/src/views/im/KnowledgeManage.vue` - AI 知识库管理页面
- `eapp/pages/me/im-sessions.vue` - 会话列表
- `eapp/pages/im/chat.vue` - 聊天页面
- `eapp/pages/im/staff-manage.vue` - 坐席管理
- `app/pages/im/chat.vue` - 用户聊天页面（AI 气泡 + 转人工）
- `web/src/views/Chat.vue` - Web 聊天页面
- `web/src/components/ChatDialog.vue` - 聊天弹窗
- `web/src/stores/chat.ts` - 聊天状态管理

## 更新日志

### 2026-06-01
- ✅ 接入本地大模型 AI 客服：新会话默认 AI 接待（`mode=ai`），可一键/关键词转人工
- ✅ RAG 知识库：`im_knowledge` 表 + 向量召回（无向量模型时关键词召回兜底）
- ✅ 文档切片入库：上传企业多格式文档（TXT/MD/CSV/TSV/JSON/XML/HTML/DOCX/PDF/XLSX）自动提取并切片为多条知识
- ✅ 商品信息分析：回答时检索在售商品价格/库存/销量并注入提示
- ✅ Admin 新增「AI知识库」页面（CRUD + 重建索引 + 连通性测试）与 `im:knowledge` 权限
- ✅ 插件 `config_items`：在配置中心维护大模型服务地址/模型/提示词等
- ✅ Qdrant 向量库检索：docker-compose 内置 qdrant，CRUD/导入/重建双写同步、按 ID 回查保序、未配置时回退内存余弦/关键词
- ✅ 新增 WS 帧：`typing`（AI 输入指示）、`to_human`（转人工请求）

### 2026-05-31
- ✅ 完成排队机制（排队位置显示、自动分配）
- ✅ 完成转接功能（客服间转接、转接通知）
- ✅ 完成客服在线状态管理
- ✅ 完成客服坐席管理（CRUD）
- ✅ 四端 WebSocket 实时通信全部实现
- ✅ 系统消息支持（接入/转接/结束通知）
