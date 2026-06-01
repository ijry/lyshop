# IM 客服插件 API 文档

## 概述

IM 客服插件提供实时客服聊天功能，支持 WebSocket 实时通信、多坐席管理、排队机制、会话转接等企业级客服功能，并内置基于本地大模型的 AI 智能客服（RAG 知识库 + 商品信息分析）。新会话默认由 AI 接待，用户可随时转人工排队。

**插件标识**: `im`  
**版本**: `1.1.0`  
**依赖**: 无

---

## 用户端 API

### 1. 获取或创建会话

**接口**: `GET /api/v1/im/session`  
**权限**: 需要用户登录  
**说明**: 获取当前用户的开放会话，不存在则自动创建

**响应示例**:
```json
{
  "id": 1,
  "user_id": 1001,
  "staff_id": 2,
  "status": 2,
  "queue_position": 0,
  "last_msg": "您好，有什么可以帮您？",
  "unread_count": 0,
  "created_at": "2026-05-31T10:00:00Z",
  "updated_at": "2026-05-31T10:05:00Z"
}
```

**字段说明**:
- `status`: 1=等待接入, 2=服务中, 3=已关闭
- `mode`: `ai`=AI 接待中, `human`=人工接待/排队
- `queue_position`: 排队位置，0表示未排队或已接入

---

### 2. 获取消息历史

**接口**: `GET /api/v1/im/messages`  
**权限**: 需要用户登录  
**参数**:
- `session_id` (必填): 会话ID
- `page` (可选): 页码，默认1
- `size` (可选): 每页数量，默认50，最大100

**响应示例**:
```json
{
  "list": [
    {
      "id": 1,
      "session_id": 1,
      "sender_type": 1,
      "sender_id": 1001,
      "type": "text",
      "content": "你好，我想咨询一下",
      "created_at": "2026-05-31T10:00:00Z"
    },
    {
      "id": 2,
      "session_id": 1,
      "sender_type": 2,
      "sender_id": 2,
      "type": "text",
      "content": "您好，有什么可以帮您？",
      "created_at": "2026-05-31T10:01:00Z"
    }
  ],
  "total": 2,
  "page": 1,
  "size": 50
}
```

**字段说明**:
- `sender_type`: 0=系统, 1=用户, 2=人工客服, 3=AI 客服
- `type`: text=文本, image=图片, product_card=商品卡片, order_card=订单卡片, system=系统消息

---

### 3. WebSocket 连接

**接口**: `GET /ws/im?token={jwt_token}`  
**协议**: WebSocket  
**权限**: 需要 JWT token

**连接示例**:
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/im?token=xxx')
```

**消息帧格式**:
```json
{
  "type": "msg|queue|assign|close|ping|pong",
  "session_id": 1,
  "payload": {}
}
```

**帧类型说明**:

| 类型 | 方向 | Payload | 说明 |
|---|---|---|---|
| `msg` | 双向 | `{msg_type, content, sender_type}` | 消息内容（`sender_type`：0系统/1用户/2人工/3AI） |
| `typing` | 服务端→客户端 | `{sender_type}` | AI 正在生成回复的输入指示 |
| `to_human` | 客户端→服务端 | `{}` | 用户请求转人工（等价于发送转人工关键词） |
| `queue` | 服务端→客户端 | `{position}` | 排队位置更新 |
| `assign` | 服务端→客户端 | `{action}` | 接入/转接通知 |
| `close` | 服务端→客户端 | `{}` | 会话结束通知 |
| `ping` | 客户端→服务端 | `{}` | 心跳请求 |
| `pong` | 服务端→客户端 | `{}` | 心跳响应 |

**发送消息示例**:
```json
{
  "type": "msg",
  "session_id": 1,
  "payload": {
    "msg_type": "text",
    "content": "你好"
  }
}
```

**接收消息示例**:
```json
{
  "type": "msg",
  "session_id": 1,
  "payload": {
    "msg_id": "123",
    "msg_type": "text",
    "content": "您好，有什么可以帮您？",
    "sender_type": 2
  }
}
```

---

## 管理端 API

### 1. 获取会话列表

**接口**: `GET /admin/api/im/sessions`  
**权限**: `im:view`  
**参数**:
- `staff_id` (可选): 按客服ID筛选
- `status` (可选): 按状态筛选 (1=等待, 2=服务中, 3=已关闭)

**响应示例**:
```json
[
  {
    "id": 1,
    "user_id": 1001,
    "staff_id": 2,
    "status": 2,
    "queue_position": 0,
    "last_msg": "您好，有什么可以帮您？",
    "unread_count": 0,
    "updated_at": "2026-05-31T10:05:00Z"
  }
]
```

---

### 2. 获取会话消息

**接口**: `GET /admin/api/im/sessions/:id/messages`  
**权限**: `im:view`  
**参数**:
- `page` (可选): 页码，默认1
- `size` (可选): 每页数量，默认50

**响应**: 同用户端消息历史接口

---

### 3. 回复消息

**接口**: `POST /admin/api/im/sessions/:id/reply`  
**权限**: `im:reply`  
**请求体**:
```json
{
  "content": "您好，有什么可以帮您？",
  "type": "text"
}
```

**响应**: 返回创建的消息对象

---

### 4. 接入会话

**接口**: `POST /admin/api/im/sessions/:id/accept`  
**权限**: `im:reply`  
**说明**: 客服手动接入等待中的会话

**响应**:
```json
{
  "success": true
}
```

---

### 5. 结束会话

**接口**: `POST /admin/api/im/sessions/:id/close`  
**权限**: `im:reply`  
**说明**: 结束当前会话，释放客服负载

**响应**:
```json
{
  "success": true
}
```

---

### 6. 转接会话

**接口**: `POST /admin/api/im/sessions/:id/transfer`  
**权限**: `im:reply`  
**请求体**:
```json
{
  "to_staff_id": 3,
  "remark": "专业问题转技术客服"
}
```

**说明**: 将会话转接给其他客服

**响应**:
```json
{
  "success": true
}
```

---

### 7. 获取客服状态

**接口**: `GET /admin/api/im/staff/status`  
**权限**: `im:view`  
**说明**: 获取当前客服的在线状态和负载

**响应示例**:
```json
{
  "admin_id": 1,
  "is_online": 1,
  "max_load": 5,
  "current_load": 3
}
```

---

### 8. 设置在线状态

**接口**: `POST /admin/api/im/staff/online`  
**权限**: `im:reply`  
**请求体**:
```json
{
  "online": true
}
```

**说明**: 切换客服上线/下线状态

**响应**:
```json
{
  "success": true
}
```

---

### 9. 客服列表

**接口**: `GET /admin/api/im/staff`  
**权限**: `im:staff:manage`  
**说明**: 获取所有客服人员列表

**响应示例**:
```json
[
  {
    "id": 1,
    "admin_id": 1,
    "is_online": 1,
    "max_load": 5,
    "current_load": 3,
    "created_at": "2026-05-31T09:00:00Z",
    "updated_at": "2026-05-31T10:00:00Z"
  }
]
```

---

### 10. 创建客服

**接口**: `POST /admin/api/im/staff`  
**权限**: `im:staff:manage`  
**请求体**:
```json
{
  "admin_id": 2,
  "max_load": 5
}
```

**说明**: 添加新的客服人员

**响应**: 返回创建的客服对象

---

### 11. 更新客服

**接口**: `PUT /admin/api/im/staff/:id`  
**权限**: `im:staff:manage`  
**请求体**:
```json
{
  "max_load": 8
}
```

**说明**: 更新客服最大负载

**响应**:
```json
{
  "success": true
}
```

---

### 12. 删除客服

**接口**: `DELETE /admin/api/im/staff/:id`  
**权限**: `im:staff:manage`  
**说明**: 删除客服人员记录

**响应**:
```json
{
  "success": true
}
```

---

### 13. 自动回复列表

**接口**: `GET /admin/api/im/auto-replies`  
**权限**: 无  
**说明**: 获取自动回复规则列表（当前为 stub）

---

### 14. 创建自动回复

**接口**: `POST /admin/api/im/auto-replies`  
**权限**: 无  
**说明**: 创建自动回复规则（当前为 stub）

---

### 15. AI 知识库列表

**接口**: `GET /admin/api/im/knowledge`  
**权限**: `im:knowledge`  
**参数**: `keyword`（可选，按标题/内容/标签检索）、`page`、`size`  
**说明**: 分页返回知识库条目（`{list,total,page,size}`）

---

### 16. 新增知识条目

**接口**: `POST /admin/api/im/knowledge`  
**权限**: `im:knowledge`  
**请求体**: `{ "title": "退货政策", "content": "7天无理由退货...", "tags": "退货,售后", "sort": 0, "status": 1 }`  
**说明**: 创建后异步进行向量化（配置了向量模型时）

---

### 17. 更新知识条目

**接口**: `PUT /admin/api/im/knowledge/:id`  
**权限**: `im:knowledge`  
**请求体**: 任意可选字段 `{title, content, tags, sort, status}`  
**说明**: 内容变更后会重新向量化

---

### 18. 删除知识条目

**接口**: `DELETE /admin/api/im/knowledge/:id`  
**权限**: `im:knowledge`

---

### 19. 重建向量索引

**接口**: `POST /admin/api/im/knowledge/reindex`  
**权限**: `im:knowledge`  
**响应**: `{ "indexed": 12 }`  
**说明**: 对全部知识条目重新向量化；未配置向量模型时返回错误提示并退化为关键词召回

---

### 20. 导入文档（多格式切片）

**接口**: `POST /admin/api/im/knowledge/import`  
**权限**: `im:knowledge`  
**请求**: `multipart/form-data`  
**表单字段**:
- `file` (必填): 上传的文档，单文件 ≤ 20MB
- `title` (可选): 标题前缀，留空使用文件名
- `tags` (可选): 共享标签，逗号分隔
- `chunk_size` (可选): 每片字数（按 rune 计，默认 500）
- `overlap` (可选): 相邻片重叠字数（默认 50）

**支持格式**: `.txt` `.md` `.markdown` `.text` `.log` `.csv` `.tsv` `.json` `.xml` `.html` `.htm` `.docx` `.pdf` `.xlsx`

**响应**: `{ "filename": "faq.docx", "chunks": 8 }`  
**说明**: 提取文本 → 按段落/句子边界切片（超长段落硬切）→ 每片生成一条 `ImKnowledge`（标题形如 `标题 (1/8)`）→ 后台逐条向量化

---

### 21. 测试大模型连通性

**接口**: `POST /admin/api/im/ai/test`  
**权限**: `im:knowledge`  
**响应**: `{ "reply": "连接正常" }`  
**说明**: 使用当前配置发起一次最小对话以校验服务地址与对话模型

---

## AI 智能客服

### 接待流程

1. 用户进入会话，`GetOrCreateSession` 创建 `mode=ai`、`status=2` 的会话，由本地大模型接待（无需排队）。
2. 用户每条消息先落库，命中转人工关键词（或收到 `to_human` 帧）则调用 `SwitchToHuman` 进入人工流程；否则推送 `typing` 帧并异步生成回复。
3. 回复生成：检索知识库（RAG）与在售商品信息，连同最近若干轮对话与系统提示词一并请求大模型，回复以 `sender_type=3` 推送给用户。
4. 关闭 AI（`ai_enabled=false`）时，新会话回退到传统人工分配/排队流程。

### 配置项（配置中心 → IM客服）

| Key | 说明 |
|---|---|
| `ai_enabled` | 是否启用 AI 客服 |
| `ai_base_url` | OpenAI 兼容服务地址，如 `http://localhost:11434/v1` |
| `ai_api_key` | API Key（本地服务可留空） |
| `ai_chat_model` | 对话模型，如 `qwen2.5:7b` |
| `ai_embed_model` | 向量模型，如 `bge-m3`；留空则关键词召回 |
| `ai_system_prompt` | 系统提示词（人设与回答约束） |
| `ai_human_keywords` | 转人工关键词，逗号分隔 |
| `ai_top_k` | 知识库召回条数（默认 3） |
| `ai_temperature` | 采样温度（默认 0.3） |
| `ai_product_search` | 是否启用商品信息分析 |
| `ai_timeout_sec` | 大模型请求超时秒数（默认 30） |
| `ai_qdrant_url` | Qdrant 向量库地址，如 `http://localhost:6333`；留空则用内存余弦/关键词召回 |
| `ai_qdrant_api_key` | Qdrant API Key（自建无鉴权可留空） |
| `ai_qdrant_collection` | Qdrant 集合名（默认 `im_knowledge`） |
| `ai_score_threshold` | 相似度阈值（0-1），低于该分数的召回结果丢弃，默认 0 不过滤 |
| `ai_hybrid` | 开启混合检索：向量召回 + 关键词召回经 RRF 融合（长尾召回更稳，推荐开启） |
| `ai_recall_k` | 重排前每路候选数，默认 4×`ai_top_k` |
| `ai_rerank_url` | 重排服务地址（Cohere/Jina/TEI 兼容 `/rerank`，留空则不重排） |
| `ai_rerank_api_key` | 重排 API Key |
| `ai_rerank_model` | 重排模型，如 `bge-reranker-v2-m3` |

### 召回 → 融合 → 重排 Pipeline

知识库召回完整流程：

```
query
  ├─ recallVector (Qdrant ANN / in-memory cosine)  ── RecallK 候选
  └─ recallKeyword (token-overlap, Hybrid=true 时)  ── RecallK 候选
       │
       └─ RRF 融合 (Hybrid=true) ─→ 去重排序候选池
              │
              └─ cross-encoder Rerank (RerankURL 非空) ─→ TopK 精排结果
```

各阶段均可独立关闭（不配置即跳过），后向兼容现有关键词兜底：

| 配置 | 召回 | 是否融合 | 是否重排 |
|---|---|---|---|
| 无 embed/无 Qdrant | 关键词 | - | - |
| embed only | 内存余弦 | - | - |
| Qdrant + embed | Qdrant ANN | - | - |
| 上述任一 + `ai_hybrid=on` | 向量 + 关键词 | ✅ RRF | - |
| 上述任一 + `ai_rerank_url` | 视配置 | 视配置 | ✅ cross-encoder |

**数据同步**：知识条目的新增/编辑会异步向量化并 upsert 到 Qdrant；删除会同步删除向量点；状态停用通过 upsert 更新 `status` payload，使其不再被检索命中。`POST /im/knowledge/reindex` 会重建 Qdrant 集合并全量重灌。DB 的 `embedding` 列作为本地缓存与回退保留。

> 部署：`docker-compose.yml` 已内置 `qdrant` 服务，容器内将地址配置为 `http://qdrant:6333` 即可。

### 重排服务连通性测试

**接口**: `POST /admin/api/im/ai/rerank-test`  
**权限**: `im:knowledge`  
**响应**: `{ "reply": "连接正常" }`

---

## 数据模型

### ImSession (会话表)

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uint64 | 主键 |
| user_id | uint64 | 用户ID (索引) |
| staff_id | uint64 | 客服ID (索引，0表示未分配) |
| mode | string | 接待模式 (`ai`=AI接待, `human`=人工) |
| status | int8 | 状态 (1=等待, 2=服务中, 3=已关闭) |
| queue_position | int | 排队位置 (0=未排队) |
| last_msg | string | 最后一条消息 (255字符) |
| unread_count | int | 未读消息数 |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |

---

### ImMessage (消息表)

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uint64 | 主键 |
| session_id | uint64 | 会话ID (索引) |
| sender_type | int8 | 发送者类型 (0=系统, 1=用户, 2=人工客服, 3=AI客服) |
| sender_id | uint64 | 发送者ID |
| type | string | 消息类型 (text/image/product_card/order_card/system) |
| content | text | 消息内容 |
| extra | json | 扩展字段 |
| is_read | int8 | 是否已读 (0=未读, 1=已读) |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |

---

### ImStaff (客服表)

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uint64 | 主键 |
| admin_id | uint64 | 管理员ID (索引) |
| is_online | int8 | 在线状态 (0=离线, 1=在线) |
| max_load | int | 最大负载 (默认5) |
| current_load | int | 当前负载 |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |

---

### ImTransferLog (转接记录表)

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uint64 | 主键 |
| session_id | uint64 | 会话ID (索引) |
| from_staff_id | uint64 | 原客服ID |
| to_staff_id | uint64 | 目标客服ID |
| remark | string | 转接备注 (255字符) |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |

---

### ImAutoReply (自动回复表)

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uint64 | 主键 |
| keyword | string | 触发关键词 (128字符) |
| match_type | int8 | 匹配方式 (1=精确, 2=包含, 3=正则) |
| reply | text | 回复内容 |
| sort | int | 排序 |
| status | int8 | 状态 (0=禁用, 1=启用) |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |

---

### ImKnowledge (AI 知识库表)

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uint64 | 主键 |
| title | string | 标题 (255字符) |
| content | text | 内容 |
| tags | string | 标签，逗号分隔 (255字符) |
| embedding | json | 内容向量（[]float64），未配置向量模型时为空 |
| indexed | int8 | 是否已向量化 (0/1) |
| sort | int | 排序 |
| status | int8 | 状态 (0=停用, 1=启用) |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |

---

## 权限说明

| 权限 | 说明 |
|---|---|
| `im:view` | 查看客服会话和消息 |
| `im:reply` | 回复消息、接入/结束/转接会话、设置在线状态 |
| `im:staff:manage` | 管理客服坐席（增删改查） |
| `im:knowledge` | 管理 AI 知识库、重建索引、测试大模型连通 |

---

## 错误码

| HTTP 状态码 | 说明 |
|---|---|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（token 无效或过期） |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 更新日志

### v1.1.0 (2026-06-01)
- ✅ 本地大模型 AI 智能客服：新会话默认 AI 接待，输入“人工”或点击转人工进入排队
- ✅ RAG 知识库：`im_knowledge` 表 + 向量召回（无向量模型时关键词召回兜底）
- ✅ 文档切片入库：上传 TXT/MD/CSV/TSV/JSON/XML/HTML/DOCX/PDF/XLSX，自动提取并切片为多条知识
- ✅ 商品信息分析：回答时检索在售商品价格/库存/销量
- ✅ 知识库管理接口与 `im:knowledge` 权限、配置中心 `config_items`
- ✅ Qdrant 向量库检索：CRUD/导入/重建双写同步，按 ID 回查并保序，未配置时回退内存余弦/关键词
- ✅ 混合检索（RRF 融合）：向量召回 + 关键词召回按 Reciprocal Rank Fusion 融合（`ai_hybrid=on`）
- ✅ 重排（Rerank）：cross-encoder 精排候选池至 TopK，兼容 Cohere/Jina/TEI `/rerank` 接口
- ✅ 新增 WS 帧 `typing` / `to_human`，`sender_type` 扩展 AI=3

### v1.0.0 (2026-05-31)
- ✅ 基础消息收发功能
- ✅ WebSocket 实时通信
- ✅ 排队机制
- ✅ 会话转接
- ✅ 客服在线状态管理
- ✅ 客服坐席管理
- ✅ 自动回复规则
- ✅ 系统消息通知
