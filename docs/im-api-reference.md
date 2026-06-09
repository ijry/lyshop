# IM 客服插件 API 文档

## 概述

IM 客服插件提供实时客服聊天功能，支持 WebSocket 实时通信、多坐席管理、排队机制、会话转接、附件消息、事件审计、报表统计等企业级客服功能，并内置基于本地大模型的 AI 智能客服（RAG 知识库 + 商品信息分析）。新会话默认由 AI 接待，用户可随时转人工排队。

**插件标识**: `im`  
**版本**: `1.1.0`  
**依赖**: 无

### 通用约定

- 用户端接口前缀：`/api/v1`
- 管理端接口前缀：`/admin/api`
- WebSocket 地址：`/ws/im?token={jwt_token}`
- HTTP 接口使用系统统一响应结构，示例中的 JSON 为 `data` 内容。
- 分页接口返回 `{list,total,page,size}`，`page` 从 1 开始。
- 时间字段使用 ISO 8601 字符串，日期筛选参数使用 `YYYY-MM-DD`。

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
  "mode": "ai",
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
- `extra`: 附件或卡片扩展信息，附件消息使用 `{file_url,file_path,file_name,file_size,mime}`

---

### 3. 上传会话附件

**接口**: `POST /api/v1/im/upload`  
**权限**: 需要用户登录  
**请求**: `multipart/form-data`

**表单字段**:
- `session_id` (必填): 当前用户所属会话 ID
- `file` (必填): 图片或文件

**响应示例**:
```json
{
  "url": "https://example.com/uploads/im/photo.png",
  "path": "uploads/im/photo.png",
  "name": "photo.png",
  "size": 1024,
  "mime": "image/png",
  "message_type": "image"
}
```

**说明**: 单文件最大 10MB；图片类型为 `.jpg/.jpeg/.png/.gif/.webp`，文件类型为 `.pdf/.doc/.docx/.xls/.xlsx/.txt/.csv/.md/.zip`。上传结果通过 WebSocket `msg` 帧发送，`msg_type` 为 `image` 或 `file`。

**会话校验**: 用户端只能上传到自己的会话，否则返回 403。

**发送流程**:
1. 调用 `POST /api/v1/im/upload` 上传文件。
2. 将响应中的 `message_type` 作为 WebSocket `payload.msg_type`。
3. 将 `name` 作为 `content`，并把 `url/path/name/size/mime` 写入 `payload.extra`。
4. 服务端落库后推送 `msg` 帧给会话双方。

---

### 4. WebSocket 连接

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
  "type": "msg|typing|to_human|queue|assign|close|ping|pong",
  "session_id": 1,
  "payload": {}
}
```

**帧类型说明**:

| 类型 | 方向 | Payload | 说明 |
|---|---|---|---|
| `msg` | 双向 | `{msg_type, content, sender_type, extra}` | 消息内容（`sender_type`：0系统/1用户/2人工/3AI） |
| `typing` | 服务端→客户端 | `{sender_type}` | AI 正在生成回复的输入指示 |
| `to_human` | 客户端→服务端 | `{}` | 用户请求转人工（等价于发送转人工关键词） |
| `queue` | 服务端→客户端 | `{position}` | 排队位置更新 |
| `assign` | 服务端→客户端 | `{action}` | 接入/转接通知 |
| `close` | 服务端→客户端 | `{}` | 会话结束通知 |
| `ping` | 客户端→服务端 | `{}` | 心跳请求 |
| `pong` | 服务端→客户端 | `{}` | 心跳响应 |

**连接身份**:
- 用户 token：服务端自动获取或创建当前用户会话，并绑定 `user_{user_id}` 客户端。
- 管理员 token：服务端绑定 `staff_{admin_id}` 客户端，用于接收接入、转接和会话消息。

**可靠性边界**:
- WebSocket 用于实时投递，消息最终以数据库 `im_messages` 为准。
- 客户端重连后应调用消息历史接口补拉，避免网络中断期间丢失展示。
- `ping/pong` 用于连接保活，不替代消息确认机制。

**发送消息示例**:
```json
{
  "type": "msg",
  "session_id": 1,
  "payload": {
    "msg_type": "image",
    "content": "photo.png",
    "extra": {
      "file_url": "https://example.com/uploads/im/photo.png",
      "file_path": "uploads/im/photo.png",
      "file_name": "photo.png",
      "file_size": 1024,
      "mime": "image/png"
    }
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
    "msg_type": "file",
    "content": "policy.pdf",
    "sender_type": 2,
    "extra": {
      "file_url": "https://example.com/uploads/im/policy.pdf",
      "file_name": "policy.pdf"
    }
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
  "type": "text",
  "extra": ""
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

### 22. 管理端上传会话附件

**接口**: `POST /admin/api/im/upload`  
**权限**: `im:reply`  
**请求**: `multipart/form-data`

**表单字段**:
- `session_id` (必填): 会话 ID
- `file` (必填): 图片或文件

**响应**: 同用户端上传附件接口。

**说明**: 管理端校验会话存在并记录 `file_uploaded` 事件，事件来源为 `staff`。业务端发送附件消息仍需调用 `POST /admin/api/im/sessions/:id/reply` 或通过 WebSocket 发送消息，将上传结果写入消息 `extra`。

---

### 23. 客服报表

**接口**: `GET /admin/api/im/analytics`  
**权限**: `im:view`

**参数**:
- `from` (可选): 起始日期，格式 `YYYY-MM-DD`
- `to` (可选): 结束日期，格式 `YYYY-MM-DD`
- `staff_id` (可选): 客服 ID

**时间范围**: `from` 按当天 00:00:00 起算；`to` 包含当天，后端实际转换为次日 00:00:00 开区间。

**响应示例**:
```json
{
  "summary": {
    "sessions": 12,
    "messages": 88,
    "ai_replies": 31,
    "ai_failed": 1,
    "rag_hits": 20,
    "to_human": 6,
    "accepts": 5,
    "closes": 4,
    "transfers": 1,
    "files": 3
  },
  "trend": [
    {
      "date": "2026-06-09",
      "sessions": 12,
      "messages": 88,
      "files": 3
    }
  ]
}
```

**统计口径**:

| 字段 | 来源事件 | 说明 |
|---|---|---|
| `sessions` | `session_created` | 新会话数 |
| `messages` | `message_sent` | 消息落库数，包含用户、客服、AI 和系统消息 |
| `ai_replies` | `ai_reply` | AI 成功回复数 |
| `ai_failed` | `ai_failed` | AI 调用或生成失败数 |
| `rag_hits` | `rag_hit` | 知识库或商品上下文命中次数 |
| `to_human` | `to_human` | 转人工次数 |
| `accepts` | `staff_accept` | 客服接入次数 |
| `closes` | `session_close` | 会话关闭次数 |
| `transfers` | `session_transfer` | 会话转接次数 |
| `files` | `file_uploaded` | 附件上传次数 |

---

### 24. 事件日志

**接口**: `GET /admin/api/im/logs`  
**权限**: `im:view`

**参数**:
- `event`、`session_id`、`user_id`、`staff_id`、`source`、`success`
- `page`、`size`

**参数说明**:
- `event`: 事件类型，如 `message_sent`
- `source`: `user`、`staff`、`ai`、`system`
- `success`: `1` 成功，`0` 失败
- `size`: 默认 20，最大 100

**响应示例**:
```json
{
  "list": [
    {
      "id": 101,
      "event": "file_uploaded",
      "session_id": 1,
      "user_id": 1001,
      "staff_id": 0,
      "message_id": 0,
      "source": "user",
      "success": 1,
      "latency_ms": 0,
      "extra": "{\"name\":\"photo.png\",\"size\":1024,\"type\":\"image\"}",
      "created_at": "2026-06-09T10:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "size": 20
}
```

**事件类型**:

| event | 触发场景 |
|---|---|
| `session_created` | 创建用户会话 |
| `message_sent` | 消息保存成功 |
| `ai_reply` | AI 回复保存成功 |
| `ai_failed` | AI 回复生成失败 |
| `rag_hit` | RAG 或商品上下文召回命中 |
| `to_human` | 用户请求转人工 |
| `staff_accept` | 客服接入会话 |
| `session_close` | 会话关闭 |
| `session_transfer` | 会话转接 |
| `file_uploaded` | 用户或客服上传附件 |

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
| `ai_query_rewrite` | 查询改写模式：`""` 关闭 / `rewrite` LLM改写 / `hyde` 生成假设回答 / `multi` 多路变体+RRF |
| `ai_query_rewrite_n` | `multi` 模式生成的查询变体数（默认 3） |
| `ai_auto_eval` | 开启 LLM-as-Judge 自动评估（忠实度 + 相关性，异步存入 `ImFeedback`） |

### 召回 → 融合 → 重排 Pipeline

知识库召回完整流程：

```
用户问题 (userText)
  │
  ├─ QueryRewrite ─→ retrievalQuery
  │   rewrite: LLM 扩写     hyde: 生成假设回答     multi: N 变体各自检索→RRF
  │
  ├─ recallVector(retrievalQuery, RecallK)  ← Qdrant ANN / 内存余弦 / nil
  └─ recallKeyword(retrievalQuery, RecallK) ← token-overlap (Hybrid=true)
       │
       └─ RRF 融合 (Hybrid=true) ─→ 去重候选池
              │
              └─ cross-encoder Rerank (RerankURL 非空) ─→ TopK 精排
                     │
                     └─ 注入 Prompt → LLM 生成 → reply
                                              │
                                              └─ AutoEval (异步) → ImFeedback
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

**说明**: 仅测试当前配置的 `ai_rerank_url`、`ai_rerank_api_key`、`ai_rerank_model` 是否可用；未配置重排服务时应保持为空，不影响基础 RAG。

---

### 用户提交反馈

**接口**: `POST /api/v1/im/feedback`  
**权限**: 需用户登录  
**请求体**: `{ "session_id": 1, "rating": 1, "comment": "很有帮助", "query": "...", "answer": "..." }`  
**说明**: `rating` 1=👍 -1=👎；`query`/`answer` 可选，存入用于后续分析

---

### 管理端反馈列表

**接口**: `GET /admin/api/im/feedback`  
**权限**: `im:view`  
**参数**: `session_id`（可选）、`page`、`size`  
**响应**: `{ list: [...], total, page, size }`  
**字段**:
- `source`: `user`（用户提交）/ `auto`（LLM-as-Judge 自动评估）
- `rating`: 用户评分 1=👍 -1=👎 0=未评
- `faithfulness` / `relevance`: 自动评估分数（0-5），`source=user` 时为 0

---

### 管理端反馈统计

**接口**: `GET /admin/api/im/feedback/stats`  
**权限**: `im:view`  
**响应**: `{ "auto": { count, avg_faith, avg_relevance, avg_rating }, "user": { ... } }`

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
| type | string | 消息类型 (text/image/file/product_card/order_card/system) |
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

### ImFeedback (AI 答案评估反馈表)

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uint64 | 主键 |
| session_id | uint64 | 会话ID (索引) |
| source | string | `user`=用户提交，`auto`=LLM-as-Judge |
| rating | int8 | 用户评分 1=👍 -1=👎 0=未评 |
| comment | string | 用户评语 (512字符) |
| faithfulness | float64 | 忠实度（0-5，auto 填写） |
| relevance | float64 | 相关性（0-5，auto 填写） |
| query | text | 用户问题 |
| answer | text | AI 回答 |
| created_at | time | 创建时间 |

---

### ImEventLog (事件日志表)

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uint64 | 主键 |
| event | string | 事件类型，如 `session_created`、`message_sent`、`file_uploaded` |
| session_id | uint64 | 会话 ID |
| user_id | uint64 | 用户 ID |
| staff_id | uint64 | 客服 ID |
| message_id | uint64 | 消息 ID |
| source | string | 来源：`user`、`staff`、`ai`、`system` |
| success | int8 | 是否成功 |
| latency_ms | int64 | 延迟毫秒 |
| extra | json | 扩展信息 |
| created_at | time | 创建时间 |

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

### 常见错误

| 场景 | 返回 | 处理建议 |
|---|---|---|
| 上传缺少 `session_id` | 400 | 表单补充会话 ID |
| 用户上传到非本人会话 | 403 | 重新获取当前用户会话 |
| 上传文件超过 10MB | 400 | 压缩文件或改用外部文件系统 |
| 上传扩展名不支持 | 400 | 使用允许的图片或文件类型 |
| 管理端缺少 `im:reply` | 403 | 为客服角色授予回复权限 |
| AI 服务不可用 | 400/500 | 检查 `ai_base_url`、模型名、超时和本地推理服务 |
| Qdrant 未配置或不可用 | 不影响基础问答 | 系统回退内存向量或关键词召回 |

---

## 部署与存储

- 附件上传复用系统当前启用的 storage driver；Docker 本地存储沿用 `./data/uploads:/app/uploads`。
- 多个后端副本必须共享同一个外部 Redis，WebSocket Hub 通过 `lyshop:im:ws` Pub/Sub 频道跨实例扇出消息，并通过节点 ID 避免回环。
- 嵌入式 Redis 适合单实例运行，不提供跨副本投递能力。
- 对象存储或本地上传目录必须对 Web、App、Eapp 和 Admin 可访问；返回的 `url` 是前端预览和下载的唯一入口。
- 多副本部署时不要把上传目录放在单个 Pod/容器临时目录，除非所有副本共享同一挂载或使用对象存储。
- 使用本地大模型时，后端服务需要能访问 `ai_base_url`、`ai_qdrant_url` 和 `ai_rerank_url` 对应网络地址；容器内地址应使用 compose/service 名称。
