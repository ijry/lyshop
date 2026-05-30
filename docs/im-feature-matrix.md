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

## 权限说明

| 权限 | 说明 | 默认角色 |
|---|---|---|
| `im:view` | 查看客服会话和消息 | 客服、管理员 |
| `im:reply` | 回复消息、接入/结束/转接会话 | 客服、管理员 |
| `im:staff:manage` | 管理客服坐席（增删改查） | 管理员 |

超级管理员默认拥有所有权限。

## 技术架构

### 后端
- **语言**: Go
- **WebSocket**: gorilla/websocket
- **消息总线**: Hub 模式（单播/广播）
- **数据库表**:
  - `im_session` - 会话表（用户、客服、状态、排队位置）
  - `im_message` - 消息表（会话ID、发送者、内容）
  - `im_staff` - 客服表（管理员ID、在线状态、负载）
  - `im_transfer_log` - 转接记录表
  - `im_auto_reply` - 自动回复规则表

### 前端
- **Admin**: Vue 3 + TypeScript
- **Eapp**: UniApp (Vue 3)
- **App**: UniApp (Vue 3)
- **Web**: Vue 3 + TypeScript

### WebSocket 帧类型

| 类型 | 方向 | 说明 |
|---|---|---|
| `msg` | 双向 | 消息内容 |
| `queue` | 服务端→客户端 | 排队位置更新 |
| `assign` | 服务端→客户端 | 接入/转接通知 |
| `close` | 服务端→客户端 | 会话结束通知 |
| `ping/pong` | 双向 | 心跳保活 |

## 主要功能流程

### 1. 用户发起咨询
1. 用户打开聊天页面
2. 后端创建会话（status: WAITING）
3. 查找在线且未满载的客服
   - 有 → 直接分配（status: ONGOING）
   - 无 → 进入排队（queue_position: N）
4. 推送排队位置给用户

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
- `server/plugins/im/service/session.go` - 业务逻辑
- `server/plugins/im/service/hub.go` - WebSocket Hub
- `server/plugins/im/api/admin.go` - Admin API
- `server/plugins/im/api/front.go` - 用户端 API
- `server/plugins/im/plugin.json` - 插件配置

### 前端
- `admin/src/views/im/SessionList.vue` - 客服会话页面
- `admin/src/views/im/StaffManage.vue` - 坐席管理页面
- `eapp/pages/me/im-sessions.vue` - 会话列表
- `eapp/pages/im/chat.vue` - 聊天页面
- `eapp/pages/im/staff-manage.vue` - 坐席管理
- `app/pages/im/chat.vue` - 用户聊天页面
- `web/src/views/Chat.vue` - Web 聊天页面
- `web/src/components/ChatDialog.vue` - 聊天弹窗
- `web/src/stores/chat.ts` - 聊天状态管理

## 更新日志

### 2026-05-31
- ✅ 完成排队机制（排队位置显示、自动分配）
- ✅ 完成转接功能（客服间转接、转接通知）
- ✅ 完成客服在线状态管理
- ✅ 完成客服坐席管理（CRUD）
- ✅ 四端 WebSocket 实时通信全部实现
- ✅ 系统消息支持（接入/转接/结束通知）
