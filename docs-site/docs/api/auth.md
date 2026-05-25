# 认证接口

## 说明

认证模块用于前台用户登录、后台管理员登录与会话校验。

## 典型接口

- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/auth/profile`
- `GET /admin/api/menus`

## 管理后台菜单接口（分组导航）

### 请求

- `GET /admin/api/menus`

### 响应示例（当前结构）

```json
{
  "dashboard": { "title": "Dashboard", "path": "/dashboard" },
  "groups": [
    {
      "key": "product",
      "title": "商品",
      "icon": "box",
      "sort": 10,
      "menus": [
        {
          "title": "商品管理",
          "path": "/product",
          "children": [{ "title": "商品列表", "path": "/product/list" }]
        }
      ]
    }
  ]
}
```

### 字段说明

- `dashboard`：固定入口菜单，不属于任何分组。
- `groups[]`：一级分组（如商品、订单、用户、系统、营销、仓储）。
- `groups[].menus[]`：分组下菜单树，兼容原有 `title/path/icon/sort/children` 结构。

### 兼容说明

- 兼容期内，旧版后端可能仍返回“菜单数组”结构。
- 前端应优先识别 `groups` 字段使用新分组渲染；若不存在则降级为旧单列渲染。

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
