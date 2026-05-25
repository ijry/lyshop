# 认证接口

## 说明

认证模块用于前台用户登录、后台管理员登录与会话校验。

## 典型接口

- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/auth/profile`

## 请求示例

```json
{
  "username": "demo",
  "password": "******"
}
```

## 响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "jwt-token"
  }
}
```

## 最近前端交互调整（2026-05-25）

- 功能说明：H5 个人中心“注销账号”入口已迁移到“账号与安全”二级菜单。
- 接口变化：无，继续使用既有认证与销户相关接口。
- 部署/配置影响：无，仅前端页面结构与图标展示调整。
