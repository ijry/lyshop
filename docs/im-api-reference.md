# IM 客服插件 API 文档

## 概述

IM 客服插件提供实时客服聊天功能，支持 WebSocket 实时通信、多坐席管理、排队机制、会话转接等企业级客服功能。

**插件标识**: `im`  
**版本**: `1.0.0`  
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
- `sender_type`: 0=系统, 1=用户, 2=客服
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
| `msg` | 双向 | `{msg_type, content, sender_type}` | 消息内容 |
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

## 数据模型

### ImSession (会话表)

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uint64 | 主键 |
| user_id | uint64 | 用户ID (索引) |
| staff_id | uint64 | 客服ID (索引，0表示未分配) |
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
| sender_type | int8 | 发送者类型 (0=系统, 1=用户, 2=客服) |
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

## 权限说明

| 权限 | 说明 |
|---|---|
| `im:view` | 查看客服会话和消息 |
| `im:reply` | 回复消息、接入/结束/转接会话、设置在线状态 |
| `im:staff:manage` | 管理客服坐席（增删改查） |

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

### v1.0.0 (2026-05-31)
- ✅ 基础消息收发功能
- ✅ WebSocket 实时通信
- ✅ 排队机制
- ✅ 会话转接
- ✅ 客服在线状态管理
- ✅ 客服坐席管理
- ✅ 自动回复规则
- ✅ 系统消息通知
