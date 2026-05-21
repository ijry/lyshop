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
