# IM 接口

## 说明

IM 模块提供会话、消息、附件、事件审计、报表统计与通知能力，并内置基于本地大模型的 AI 智能客服（RAG 知识库 + 商品信息分析）。新会话默认由 AI 接待，用户可随时转人工排队。多实例部署时，WebSocket Hub 通过 Redis Pub/Sub 做跨实例消息扇出。

接口默认返回系统统一响应结构，业务数据位于 `data`。分页接口统一返回 `{list,total,page,size}`，未显式指定时 `page=1`、`size=20` 或接口默认值。

## 用户端接口

| 接口 | 权限 | 说明 |
|---|---|---|
| `GET /api/v1/im/session` | 用户登录 | 获取或创建当前用户开放会话，返回 `mode`：`ai`/`human` |
| `GET /api/v1/im/messages?session_id=1&page=1&size=50` | 用户登录 | 拉取会话历史消息 |
| `POST /api/v1/im/upload` | 用户登录 | 上传当前用户所属会话附件，表单字段：`session_id`、`file` |
| `POST /api/v1/im/feedback` | 用户登录 | 提交 AI 回答反馈，`rating` 为 `1` 或 `-1` |
| `GET /ws/im?token=...` | JWT | WebSocket 长连接，用户 token 自动绑定当前会话 |

## 管理端接口

| 接口 | 权限 | 说明 |
|---|---|---|
| `GET /admin/api/im/sessions?staff_id=&status=` | `im:view` | 会话列表，状态：1等待、2服务中、3已关闭 |
| `GET /admin/api/im/sessions/:id/messages?page=1&size=50` | `im:view` | 会话消息 |
| `POST /admin/api/im/sessions/:id/reply` | `im:reply` | 人工回复，支持 `type=text/image/file/system` 与 `extra` |
| `POST /admin/api/im/sessions/:id/accept` | `im:reply` | 接入等待会话 |
| `POST /admin/api/im/sessions/:id/close` | `im:reply` | 结束会话并释放客服负载 |
| `POST /admin/api/im/sessions/:id/transfer` | `im:reply` | 转接会话，参数 `to_staff_id`、`remark` |
| `POST /admin/api/im/upload` | `im:reply` | 客服上传会话附件，表单字段：`session_id`、`file` |
| `GET /admin/api/im/analytics?from=&to=&staff_id=` | `im:view` | 客服报表汇总与日趋势 |
| `GET /admin/api/im/logs?event=&session_id=&source=&success=` | `im:view` | 事件日志查询 |
| `GET /admin/api/im/feedback` | `im:view` | AI 反馈列表 |
| `GET /admin/api/im/feedback/stats` | `im:view` | AI 反馈统计 |
- 知识库（需 `im:knowledge` 权限）：
  - `GET /admin/api/im/knowledge` — 列表（支持 `keyword` 检索）
  - `POST /admin/api/im/knowledge` — 新增
  - `PUT /admin/api/im/knowledge/:id` — 更新
  - `DELETE /admin/api/im/knowledge/:id` — 删除
  - `POST /admin/api/im/knowledge/reindex` — 重建向量索引
  - `POST /admin/api/im/knowledge/import` — 上传多格式文档，自动切片入库
  - `POST /admin/api/im/ai/test` — 测试本地大模型连通性

## WebSocket 帧类型

| 类型 | 方向 | 说明 |
|---|---|---|
| `msg` | 双向 | 消息内容（`sender_type`：0系统/1用户/2人工/3AI） |
| `typing` | 服务端→客户端 | AI 正在生成回复 |
| `to_human` | 客户端→服务端 | 用户请求转人工 |
| `queue` | 服务端→客户端 | 排队位置更新 |
| `assign` | 服务端→客户端 | 接入/转接通知 |
| `close` | 服务端→客户端 | 会话结束通知 |
| `ping/pong` | 双向 | 心跳保活 |

`msg` 帧的 `payload` 支持：

```json
{
  "msg_type": "text|image|file|product_card|order_card|system",
  "content": "消息文本或附件名称",
  "sender_type": 1,
  "extra": {
    "file_url": "https://example.com/uploads/im/a.png",
    "file_path": "uploads/im/a.png",
    "file_name": "a.png",
    "file_size": 1024,
    "mime": "image/png"
  }
}
```

客户端发送附件消息时，先调用上传接口得到 `url/path/name/size/mime/message_type`，再发送 `msg` 帧。服务端保存后的消息会回推给会话双方，`extra` 建议按 JSON 字符串或对象保存同一组附件字段。

## 附件消息

- 用户端上传：`POST /api/v1/im/upload`
- 管理端上传：`POST /admin/api/im/upload`
- 表单字段：`session_id`、`file`
- 会话校验：用户端必须上传到自己的会话；管理端校验会话存在并要求 `im:reply`
- 文件大小：单文件最大 10MB
- 图片类型：`.jpg`、`.jpeg`、`.png`、`.gif`、`.webp`
- 文件类型：`.pdf`、`.doc`、`.docx`、`.xls`、`.xlsx`、`.txt`、`.csv`、`.md`、`.zip`
- 存储：复用系统当前启用的 storage driver，本地 Docker 部署沿用 `./data/uploads:/app/uploads`
- 消息类型：图片保存为 `image`，普通文件保存为 `file`，附件元数据写入 `ImMessage.extra`

## 报表与事件日志

后台客服报表基于 `ImEventLog` 聚合生成，统计项包含：

- `sessions`：新会话数
- `messages`：消息数
- `ai_replies` / `ai_failed`：AI 回复成功/失败
- `rag_hits`：RAG 或商品上下文命中
- `to_human`：转人工
- `accepts`：人工接入
- `closes`：会话关闭
- `transfers`：会话转接
- `files`：附件上传

事件日志记录 `event`、`session_id`、`user_id`、`staff_id`、`message_id`、`source`、`success`、`latency_ms`、`extra`，用于审计、排障和报表统计。

当前事件类型：

| event | 计入口径 | 典型来源 |
|---|---|---|
| `session_created` | `sessions` | 用户创建会话 |
| `message_sent` | `messages` | 用户、客服、AI 消息落库 |
| `ai_reply` | `ai_replies` | AI 成功生成回复 |
| `ai_failed` | `ai_failed` | AI 调用失败 |
| `rag_hit` | `rag_hits` | 知识库或商品上下文命中 |
| `to_human` | `to_human` | 用户关键词或按钮转人工 |
| `staff_accept` | `accepts` | 客服接入 |
| `session_close` | `closes` | 会话关闭 |
| `session_transfer` | `transfers` | 会话转接 |
| `file_uploaded` | `files` | 用户或客服上传附件 |

`/admin/api/im/analytics` 的 `from/to` 使用 `YYYY-MM-DD`。`to` 按结束日期当天包含处理，后端实际使用次日零点开区间过滤。事件日志分页 `size` 最大 100。

## AI 智能客服

- **接待流程**：会话创建即进入 AI 模式，AI 基于知识库与在售商品信息作答；用户输入“人工/转人工”等关键词或点击「转人工」按钮后进入人工排队。
- **本地大模型**：通过 OpenAI 兼容接口（`/chat/completions`、`/embeddings`）对接本地推理服务，可在「配置中心 → IM客服」配置服务地址、对话模型、向量模型、系统提示词、转人工关键词、召回条数、温度、商品分析开关与超时时间。
- **完整 RAG pipeline**（与行业标准对齐）：查询改写（rewrite/HyDE/multi-query）→ 双路召回（Qdrant ANN + 关键词）→ RRF 融合（`ai_hybrid`）→ cross-encoder 重排（`ai_rerank_url`，兼容 Cohere/Jina/TEI）。各阶段独立可关，未配置时依次退化为内存余弦 → 关键词兜底。知识条目与 Qdrant 双写同步；`docker-compose.yml` 内置 `qdrant` 服务（容器内 `http://qdrant:6333`）。
- **评估闭环**：用户可👍/👎 AI 回答（`POST /api/v1/im/feedback`）；开启 `ai_auto_eval` 后 LLM-as-Judge 自动评估忠实度和相关性（0-5，异步存入 `ImFeedback`）；管理后台可查列表与聚合统计（`GET /admin/api/im/feedback/stats`）。
- **文档切片入库**：支持上传企业多格式文档（TXT/Markdown/CSV/TSV/JSON/XML/HTML/DOCX/PDF/XLSX，≤20MB），后端提取纯文本并按段落/句子边界切片（超长段落硬切、片间可重叠），每片落库为一条知识并自动向量化（`POST /admin/api/im/knowledge/import`）。

## 集成建议

- 会话应绑定用户身份和业务上下文。
- H5 端在 WebSocket 不可用时应提供本地发送与兜底回复，保障可对话性。
- PC 端推荐统一使用站内客服弹窗，而非新页面跳转。
- 客服入口应保持可输入可发送的会话状态，避免仅打开页面但无法开始对话。
- 未配置或关闭 AI（`ai_enabled` 关闭）时，新会话自动回退到传统人工分配/排队流程。
- 多个后端副本必须共享同一个外部 Redis，才能保证 WebSocket 跨实例投递；嵌入式 Redis 只适合单实例运行。

## 部署检查清单

- 已启用 IM 插件并完成数据库迁移，包含 `im_sessions`、`im_messages`、`im_event_logs`、`im_feedbacks`、`im_knowledges` 等表。
- 后台角色已授予 `im:view`、`im:reply`，需要坐席管理时授予 `im:staff:manage`，需要 AI 知识库时授予 `im:knowledge`。
- 附件上传目录或对象存储已配置持久化，Nginx/CDN 可访问返回的 `url`。
- 多副本后端连接同一个外部 Redis；单实例可以使用嵌入式 Redis。
- AI 客服需要配置 OpenAI 兼容对话服务；需要向量召回时配置 embedding 模型和 Qdrant。
