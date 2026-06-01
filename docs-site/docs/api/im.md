# IM 接口

## 说明

IM 模块提供会话、消息与通知能力，并内置基于本地大模型的 AI 智能客服（RAG 知识库 + 商品信息分析）。新会话默认由 AI 接待，用户可随时转人工排队。

## 用户端接口

- `GET /api/v1/im/session` — 获取或创建当前用户会话（返回 `mode`：`ai`/`human`）
- `GET /api/v1/im/messages` — 拉取会话历史消息
- `GET /ws/im?token=...` — WebSocket 长连接（实时消息、排队、转接、AI 输入指示）

## 管理端接口

- `GET /admin/api/im/sessions` — 会话列表
- `GET /admin/api/im/sessions/:id/messages` — 会话消息
- `POST /admin/api/im/sessions/:id/{reply,accept,close,transfer}` — 人工回复/接入/结束/转接
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
