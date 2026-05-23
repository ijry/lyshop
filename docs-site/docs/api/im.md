# IM 接口

## 说明

IM 模块提供会话、消息与通知能力。

## 典型接口

- `GET /api/im/session/list`
- `POST /api/im/session/create`
- `GET /api/im/message/list`
- `POST /api/im/message/send`

## 说明

- 会话应绑定用户身份和业务上下文
- 消息投递建议与前端轮询或长连接方案配合使用
- H5 端在 WebSocket 不可用时应提供本地发送与自动回复兜底，保障可对话性。
- PC 端推荐统一使用站内客服弹窗，而非新页面跳转。
- 客服入口应保持可输入可发送的会话状态，避免仅打开页面但无法开始对话。
