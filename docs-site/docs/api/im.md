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
