# lyshop 2.0 系统设计文档

**日期：** 2026-05-21
**版本：** 1.0.0
**作者：** jry

---

## 目录

1. [项目概述](#1-项目概述)
2. [技术栈](#2-技术栈)
3. [整体架构](#3-整体架构)
4. [插件系统](#4-插件系统)
5. [驱动抽象层](#5-驱动抽象层)
6. [核心插件功能清单](#6-核心插件功能清单)
7. [数据库设计](#7-数据库设计)
8. [API 设计规范](#8-api-设计规范)
9. [IM WebSocket 协议](#9-im-websocket-协议)
10. [AI 生图模块](#10-ai-生图模块)
11. [管理后台 UI 规范](#11-管理后台-ui-规范)
12. [uni-app 前端规范](#12-uni-app-前端规范)
13. [部署方案](#13-部署方案)

---

## 1. 项目概述

lyshop 2.0 是一套开源商城系统，面向中型商户（万级并发）设计，支持私有化部署。

**核心特性：**
- 后端基于 Go（Gin + GORM），管理后台基于 Vue3 + TailwindCSS，前台基于 uni-app + uview-plus 3.x
- 完全插件化架构，每个功能模块以插件形式存在，可按需启用/禁用
- 驱动抽象层统一支付、短信、OAuth登录、文件存储、AI生图接口
- 内置 WebSocket IM 客服系统
- 支持 AI 生成商品轮播图和详情图（多模型聚合）
- uni-app 支持微信小程序 / H5 / App（iOS + Android）三端发布
- Docker Compose 一键部署

**目标规模：** 单机万级并发，数据表预留 `merchant_id` 字段支持未来多商户扩展。

---

## 2. 技术栈

| 层级 | 技术 | 版本 |
|------|------|------|
| 后端框架 | Go + Gin | Go 1.22+ |
| ORM | GORM | v2 |
| 数据库 | MySQL | 8.0 |
| 缓存 | Redis | 7.x |
| 实时通信 | Gorilla WebSocket | v1 |
| 管理后台框架 | Vue3 + Vite | Vue 3.4+ |
| 管理后台样式 | TailwindCSS + shadcn-vue | TailwindCSS 3.x |
| 管理后台状态 | Pinia | v2 |
| 前台框架 | uni-app | Vue3 版本 |
| 前台组件库 | uview-plus | 3.x |
| 前台状态 | Pinia | v2 |
| 容器化 | Docker Compose | v2 |
| 反向代理 | Nginx | 1.25+ |

---

## 3. 整体架构

### 3.1 目录结构

```
lyshop/
├── server/                      # Go 后端
│   ├── core/                    # 框架核心
│   │   ├── app/                 # 应用初始化、生命周期管理
│   │   ├── router/              # 路由注册器（前台/后台分组）
│   │   ├── plugin/              # 插件注册表、加载器
│   │   ├── middleware/          # JWT鉴权、CORS、限流、日志
│   │   ├── db/                  # MySQL连接池（GORM）
│   │   ├── cache/               # Redis客户端
│   │   ├── ws/                  # WebSocket Hub（IM客服）
│   │   └── driver/              # 驱动抽象层
│   │       ├── payment/         # 支付驱动接口
│   │       ├── sms/             # 短信驱动接口
│   │       ├── oauth/           # OAuth登录驱动接口
│   │       ├── storage/         # 文件存储驱动接口
│   │       └── ai_image/        # AI生图驱动接口
│   ├── plugins/                 # 所有插件
│   │   ├── product/             # 商品插件
│   │   ├── order/               # 订单插件
│   │   ├── marketing/           # 营销插件
│   │   ├── im/                  # IM客服插件
│   │   ├── wms/                 # 仓储插件
│   │   ├── ai_image/            # AI生图插件
│   │   ├── wechat_pay/          # 微信支付插件
│   │   ├── alipay/              # 支付宝支付插件
│   │   ├── sms/                 # 短信插件
│   │   ├── wechat_auth/         # 微信登录插件
│   │   ├── decor/               # 装修插件
│   │   ├── storage_local/       # 本地存储插件
│   │   ├── storage_oss/         # 阿里云OSS插件
│   │   ├── storage_cos/         # 腾讯云COS插件
│   │   └── storage_qiniu/       # 七牛云存储插件
│   └── main.go
│
├── admin/                       # Vue3 + TailwindCSS 管理后台
│   ├── src/
│   │   ├── views/               # 核心页面（登录、Dashboard等）
│   │   ├── plugins/             # 各插件对应后台页面（动态加载）
│   │   ├── components/          # 通用业务组件
│   │   └── layouts/             # 布局组件
│   └── vite.config.ts
│
└── app/                         # uni-app 前端
    ├── pages/                   # 核心页面
    ├── plugins/                 # 各插件对应前台页面
    ├── components/              # 基于uview-plus封装的业务组件
    └── manifest.json            # 支持微信小程序/H5/App
```

### 3.2 请求链路

```
uni-app / 浏览器
    │
    ▼
Nginx（反向代理 + SSL终止）
    │
    ├── /api/v1/*     → Go Server :8080（前台接口）
    ├── /admin/api/*  → Go Server :8080（后台接口）
    ├── /ws/im        → Go Server :8080（WebSocket）
    ├── /notify/*     → Go Server :8080（第三方回调）
    ├── /admin/*      → Vue3 Admin 静态文件 :9527
    └── /*            → uni-app H5 静态文件 :9528
```

---

## 4. 插件系统

### 4.1 插件目录结构（以 `product` 为例）

```
plugins/product/
├── plugin.json          # 插件元信息（菜单、权限、版本、依赖）
├── migrate/
│   └── product.sql      # 建表SQL（首次启动自动执行）
├── api/
│   ├── front.go         # 前台接口（注册到 /api/v1/product/...）
│   └── admin.go         # 后台接口（注册到 /admin/api/product/...）
├── model/               # GORM 数据模型
├── service/             # 业务逻辑层
└── plugin.go            # 实现 Plugin 接口（插件注册入口）
```

### 4.2 plugin.json 规范

```json
{
  "name": "product",
  "title": "商品插件",
  "version": "1.0.0",
  "description": "商品管理、分类、SKU、库存",
  "author": "lyshop",
  "depends": [],
  "menus": [
    {
      "title": "商品管理",
      "icon": "box",
      "path": "/product",
      "sort": 10,
      "children": [
        { "title": "商品列表", "path": "/product/list", "permission": "product:view" },
        { "title": "商品分类", "path": "/product/category", "permission": "product:view" },
        { "title": "商品规格", "path": "/product/spec", "permission": "product:view" }
      ]
    }
  ],
  "permissions": [
    "product:view",
    "product:create",
    "product:edit",
    "product:delete"
  ]
}
```

### 4.3 Plugin 接口（Go）

```go
// core/plugin/plugin.go
type Plugin interface {
    Name() string
    RegisterRoutes(front, admin *gin.RouterGroup)
    Migrate(db *gorm.DB) error
    Install() error
    Uninstall() error
}

// 插件注册表
var registry []Plugin

func Register(p Plugin) {
    registry = append(registry, p)
}

func LoadAll(cfg *Config, db *gorm.DB, front, admin *gin.RouterGroup) error {
    for _, name := range cfg.Plugins.Enabled {
        p := findPlugin(name)
        if p == nil {
            return fmt.Errorf("plugin %s not found", name)
        }
        if err := checkDepends(p); err != nil {
            return err
        }
        if err := p.Migrate(db); err != nil {
            return err
        }
        p.RegisterRoutes(front, admin)
        p.Install()
    }
    return nil
}
```

### 4.4 启用/禁用配置（config.yaml）

```yaml
plugins:
  enabled:
    - product
    - order
    - marketing
    - im
    - wms
    - ai_image
    - wechat_pay
    - alipay
    - sms
    - wechat_auth
    - decor
    - storage_local
```

未列出的插件不加载，对应路由、菜单、数据库迁移均不执行。

### 4.5 依赖检查

`plugin.json` 中声明 `depends` 字段，启动时校验依赖插件是否已启用，缺少依赖则启动失败并提示。

```json
// order 插件的 plugin.json
{
  "name": "order",
  "depends": ["product"]
}
```

---

## 5. 驱动抽象层

所有驱动位于 `server/core/driver/`，插件通过实现驱动接口注册到驱动注册表，业务代码通过驱动注册表统一调用。

### 5.1 支付驱动（payment）

```go
type Driver interface {
    Name() string
    CreateOrder(ctx context.Context, p *OrderParams) (*OrderResult, error)
    QueryOrder(ctx context.Context, tradeNo string) (*QueryResult, error)
    Refund(ctx context.Context, p *RefundParams) (*RefundResult, error)
    HandleNotify(ctx context.Context, r *http.Request) (*NotifyResult, error)
}
```

实现插件：`wechat_pay`（JSAPI/H5/App三场景），`alipay`（APP/H5/PC三场景）

### 5.2 短信驱动（sms）

```go
type Driver interface {
    Name() string
    Send(ctx context.Context, phone, templateCode string, params map[string]string) error
}
```

实现插件：`sms`（后台配置阿里云/腾讯云，可切换）

### 5.3 OAuth 登录驱动（oauth）

```go
type Driver interface {
    Name() string
    GetAuthURL(state string) string
    HandleCallback(ctx context.Context, code string) (*UserInfo, error)
    GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)
}
```

实现插件：`wechat_auth`（小程序code2session + H5/App OAuth2）

### 5.4 文件存储驱动（storage）

```go
type Driver interface {
    Name() string
    Upload(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error)
    Delete(ctx context.Context, path string) error
    GetURL(path string) string
}
```

实现插件：`storage_local`，`storage_oss`（阿里云），`storage_cos`（腾讯云），`storage_qiniu`（七牛云）

后台系统设置中选择当前启用的存储驱动，支持一键迁移历史文件。

### 5.5 AI 生图驱动（ai_image）

```go
type Driver interface {
    Name() string
    Generate(ctx context.Context, p *GenerateParams) (*GenerateResult, error)
}

type GenerateParams struct {
    Prompt    string
    NegPrompt string
    Width     int
    Height    int
    Count     int
    Style     string
}
```

实现：`tongyi`（通义万象），`wenxin`（文心一格），`hunyuan`（腾讯混元），`openai`（DALL-E及兼容接口）

### 5.6 驱动汇总

| 驱动 | 接口位置 | 已有实现插件 |
|------|---------|------------|
| payment | core/driver/payment | wechat_pay, alipay |
| sms | core/driver/sms | sms |
| oauth | core/driver/oauth | wechat_auth |
| storage | core/driver/storage | storage_local, storage_oss, storage_cos, storage_qiniu |
| ai_image | core/driver/ai_image | ai_image（内含多模型驱动） |

---

## 6. 核心插件功能清单

### 6.1 商品插件 `product`

| 功能 | 说明 |
|------|------|
| 商品 CRUD | 标题、描述、主图、轮播图、详情富文本 |
| 多规格 SKU | 颜色/尺寸等属性组合，每个SKU独立价格/库存 |
| 商品分类 | 多级分类树，支持图标 |
| 商品相册 | 多图管理，拖拽排序 |
| AI 生图入口 | 调用 `ai_image` 驱动生成轮播图/详情图 |
| 库存预警 | 库存低于阈值时调用 `sms` 驱动通知 |

### 6.2 订单插件 `order`

| 功能 | 说明 |
|------|------|
| 购物车 | Redis 为主存储，支持 SKU，合并登录前后购物车 |
| 下单流程 | 地址选择 → 支付方式选择 → 创建订单 |
| 支付对接 | 调用 `payment.Driver` 统一接口 |
| 订单状态机 | 待付款 → 待发货 → 待收货 → 已完成 → 售后 |
| 发货管理 | 填写快递单号，支持快递查询 |
| 退款/售后 | 退款调用 `payment.Driver.Refund` |

依赖插件：`product`

### 6.3 营销插件 `marketing`

| 功能 | 说明 |
|------|------|
| 优惠券 | 满减券、折扣券、无门槛券，限领次数 |
| 限时秒杀 | 活动时间段内特价，独立库存，Redis原子扣减 |
| 满减活动 | 满X减Y，可叠加优惠券 |
| 积分系统 | 消费积分、积分兑换商品/优惠券 |
| 拼团 | 多人拼团（预留，可选启用） |

依赖插件：`product`，`order`

### 6.4 IM 客服插件 `im`

| 功能 | 说明 |
|------|------|
| 用户发起会话 | uni-app 前台 WebSocket 接入 |
| 客服坐席 | 后台多客服，支持会话分配/转接 |
| 消息类型 | 文字、图片、商品卡片、订单卡片 |
| 离线消息 | 存 MySQL，用户重连后推送未读消息 |
| 会话记录 | 历史消息查询，支持关键词搜索 |
| 自动回复 | 关键词触发预设回复，支持正则 |

### 6.5 仓储插件 `wms`

| 功能 | 说明 |
|------|------|
| 仓库管理 | 多仓库，仓库信息维护 |
| 库存管理 | 商品/SKU 在各仓库的实时库存 |
| 入库单 | 采购入库、退货入库，入库记录 |
| 出库单 | 订单出库、调拨出库，出库记录 |
| 库存调拨 | 仓库间库存转移 |
| 库存盘点 | 盘点单，记录盈亏 |
| 库存流水 | 每笔库存变动可追溯（类型/数量/前后值/关联单号） |
| 低库存预警 | 低于安全库存触发 `sms` 驱动通知 |

依赖插件：`product`

`order` 插件发货时调用 `wms` 扣减库存；`wms` 未启用时降级到 `product` 直接扣减。

### 6.6 AI 生图插件 `ai_image`

| 功能 | 说明 |
|------|------|
| 多模型聚合 | 后台配置多个模型，按场景/商家选择 |
| 商品轮播图生成 | 输入商品名+风格，生成3-5张横版轮播图 |
| 商品详情图生成 | 输入卖点描述，生成竖版长图，一键插入富文本 |
| 生成记录管理 | 历史生成图片管理，一键应用到商品 |
| 模型配置 | 后台配置 API Key、endpoint、模型参数、默认模型 |

支持模型：通义万象 / 文心一格 / 腾讯混元 / OpenAI DALL-E 及兼容接口

### 6.7 支付插件

**微信支付 `wechat_pay`**（实现 `payment.Driver`）

| 场景 | 方式 |
|------|------|
| 微信小程序 | JSAPI 支付 |
| H5 | 微信 H5 支付 |
| App | 微信 App 支付 |
| 回调 | 异步回调签名验证，防重复处理 |
| 退款 | 支持全额/部分退款 |

**支付宝 `alipay`**（实现 `payment.Driver`）

| 场景 | 方式 |
|------|------|
| App | APP 支付 |
| H5 | 手机网站支付 |
| PC | 电脑网站支付 |
| 回调 | 异步通知验签，幂等处理 |
| 退款 | 支持全额/部分退款 |

后台配置：AppID、私钥/公钥、沙箱/生产环境切换。

### 6.8 短信插件 `sms`（实现 `sms.Driver`）

支持阿里云短信、腾讯云短信，后台配置切换。功能：模板管理、发送记录、频率限制。

### 6.9 微信登录插件 `wechat_auth`（实现 `oauth.Driver`）

- 微信小程序：`wx.login` 获取 code → code2session → 绑定/注册用户
- H5/App：微信 OAuth2 授权码模式 → 获取 unionid → 绑定/注册用户

### 6.10 装修插件 `decor`

**核心概念：** 页面由「组件列表」JSON 配置驱动，后台可视化编排，前台动态渲染。

**后台编辑器布局：**

```
┌─────────────────────────────────────────────────┐
│  左侧：组件库    中间：画布预览（手机框）  右侧：属性面板 │
│  拖拽组件到画布  ↕ 拖拽排序              实时编辑属性  │
└─────────────────────────────────────────────────┘
```

**支持的装修组件：**

| 组件类型 | 说明 |
|---------|------|
| `banner` | 轮播图，配置图片列表、跳转链接、高度 |
| `product_grid` | 商品宫格/列表，手动选品或按分类/销量自动拉取 |
| `category_nav` | 分类导航，图标+文字，支持跳转 |
| `notice` | 滚动公告栏 |
| `image_ad` | 单图/多图广告位，支持跳转 |
| `rich_text` | 富文本内容块 |
| `product_recommend` | 猜你喜欢，基于浏览历史 |
| `marketing_zone` | 营销活动入口（依赖 `marketing` 插件） |
| `spacer` | 间距分割块，可配置高度和背景色 |

**数据结构（MySQL JSON 字段存储）：**

```json
{
  "page": "index",
  "components": [
    {
      "type": "banner",
      "id": "c1",
      "props": { "images": [{"url": "...", "link": "..."}], "height": 350 }
    },
    {
      "type": "category_nav",
      "id": "c2",
      "props": { "items": [{"title": "手机", "icon": "...", "link": "/category/1"}] }
    },
    {
      "type": "product_grid",
      "id": "c3",
      "props": { "source": "hot", "limit": 10, "columns": 2 }
    }
  ]
}
```

**前台渲染：** uni-app 通过组件映射表动态渲染，新增装修组件只需注册一个 Vue 组件。

### 6.11 存储插件

| 插件 | 驱动名 | 配置项 |
|------|--------|--------|
| `storage_local` | `local` | 上传目录路径，访问URL前缀 |
| `storage_oss` | `oss` | Endpoint、Bucket、AccessKeyId、AccessKeySecret |
| `storage_cos` | `cos` | Region、Bucket、SecretId、SecretKey |
| `storage_qiniu` | `qiniu` | Zone、Bucket、AccessKey、SecretKey、Domain |

---

## 7. 数据库设计

### 7.1 核心表（core）

```sql
-- 用户表
CREATE TABLE users (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    merchant_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '预留多商户',
    phone       VARCHAR(20) UNIQUE,
    nickname    VARCHAR(64),
    avatar      VARCHAR(255),
    points      INT NOT NULL DEFAULT 0,
    status      TINYINT NOT NULL DEFAULT 1 COMMENT '1正常 0禁用',
    created_at  DATETIME,
    updated_at  DATETIME
);

-- 管理员表
CREATE TABLE admins (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    username    VARCHAR(64) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,
    role_id     BIGINT UNSIGNED NOT NULL,
    status      TINYINT NOT NULL DEFAULT 1,
    created_at  DATETIME,
    updated_at  DATETIME
);

-- 角色表
CREATE TABLE roles (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(64) NOT NULL,
    permissions JSON COMMENT '权限标识列表',
    created_at  DATETIME,
    updated_at  DATETIME
);

-- 系统配置表
CREATE TABLE configs (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    plugin      VARCHAR(64) NOT NULL COMMENT '所属插件',
    key         VARCHAR(128) NOT NULL,
    value       TEXT,
    UNIQUE KEY uk_plugin_key (plugin, key)
);

-- 上传文件记录
CREATE TABLE upload_files (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    driver      VARCHAR(32) NOT NULL COMMENT '存储驱动',
    path        VARCHAR(500) NOT NULL,
    url         VARCHAR(500) NOT NULL,
    size        BIGINT NOT NULL,
    mime_type   VARCHAR(128),
    created_at  DATETIME
);
```

### 7.2 商品插件表（product）

```sql
CREATE TABLE product_categories (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    parent_id   BIGINT UNSIGNED NOT NULL DEFAULT 0,
    name        VARCHAR(64) NOT NULL,
    icon        VARCHAR(255),
    sort        INT NOT NULL DEFAULT 0,
    status      TINYINT NOT NULL DEFAULT 1
);

CREATE TABLE products (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    merchant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    category_id BIGINT UNSIGNED NOT NULL,
    title       VARCHAR(255) NOT NULL,
    subtitle    VARCHAR(255),
    cover       VARCHAR(500),
    price       DECIMAL(10,2) NOT NULL,
    origin_price DECIMAL(10,2),
    stock       INT NOT NULL DEFAULT 0,
    sales       INT NOT NULL DEFAULT 0,
    status      TINYINT NOT NULL DEFAULT 1 COMMENT '1上架 0下架',
    sort        INT NOT NULL DEFAULT 0,
    detail      LONGTEXT COMMENT '富文本详情',
    created_at  DATETIME,
    updated_at  DATETIME
);

CREATE TABLE product_skus (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    product_id  BIGINT UNSIGNED NOT NULL,
    attrs       JSON NOT NULL COMMENT '[{"name":"颜色","value":"红色"}]',
    price       DECIMAL(10,2) NOT NULL,
    stock       INT NOT NULL DEFAULT 0,
    sku_code    VARCHAR(128),
    INDEX idx_product (product_id)
);

CREATE TABLE product_images (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    product_id  BIGINT UNSIGNED NOT NULL,
    url         VARCHAR(500) NOT NULL,
    sort        INT NOT NULL DEFAULT 0,
    INDEX idx_product (product_id)
);
```

### 7.3 订单插件表（order）

```sql
CREATE TABLE carts (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id     BIGINT UNSIGNED NOT NULL,
    sku_id      BIGINT UNSIGNED NOT NULL,
    qty         INT NOT NULL DEFAULT 1,
    created_at  DATETIME,
    updated_at  DATETIME,
    UNIQUE KEY uk_user_sku (user_id, sku_id)
);

CREATE TABLE addresses (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id     BIGINT UNSIGNED NOT NULL,
    name        VARCHAR(64) NOT NULL,
    phone       VARCHAR(20) NOT NULL,
    province    VARCHAR(32),
    city        VARCHAR(32),
    district    VARCHAR(32),
    detail      VARCHAR(255),
    is_default  TINYINT NOT NULL DEFAULT 0,
    INDEX idx_user (user_id)
);

CREATE TABLE orders (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    order_no    VARCHAR(64) UNIQUE NOT NULL,
    user_id     BIGINT UNSIGNED NOT NULL,
    merchant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    status      TINYINT NOT NULL COMMENT '1待付款 2待发货 3待收货 4已完成 5售后',
    payment_method VARCHAR(32) COMMENT 'wechat/alipay',
    goods_amount DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    freight_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(10,2) NOT NULL,
    address_snapshot JSON COMMENT '下单时地址快照',
    remark      VARCHAR(255),
    paid_at     DATETIME,
    created_at  DATETIME,
    updated_at  DATETIME,
    INDEX idx_user (user_id),
    INDEX idx_status (status)
);

CREATE TABLE order_items (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    order_id    BIGINT UNSIGNED NOT NULL,
    product_id  BIGINT UNSIGNED NOT NULL,
    sku_id      BIGINT UNSIGNED NOT NULL,
    title       VARCHAR(255) NOT NULL,
    cover       VARCHAR(500),
    attrs       JSON COMMENT 'SKU属性快照',
    price       DECIMAL(10,2) NOT NULL,
    qty         INT NOT NULL,
    INDEX idx_order (order_id)
);

CREATE TABLE order_payments (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    order_id    BIGINT UNSIGNED NOT NULL,
    driver      VARCHAR(32) NOT NULL,
    trade_no    VARCHAR(128) COMMENT '第三方交易号',
    amount      DECIMAL(10,2) NOT NULL,
    status      TINYINT NOT NULL COMMENT '1待支付 2已支付 3已退款',
    notified_at DATETIME,
    created_at  DATETIME
);

CREATE TABLE order_refunds (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    order_id    BIGINT UNSIGNED NOT NULL,
    reason      VARCHAR(255),
    amount      DECIMAL(10,2) NOT NULL,
    status      TINYINT NOT NULL COMMENT '1申请中 2已退款 3已拒绝',
    refund_no   VARCHAR(128),
    created_at  DATETIME
);
```

### 7.4 仓储插件表（wms）

```sql
CREATE TABLE warehouses (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(64) NOT NULL,
    address     VARCHAR(255),
    contact     VARCHAR(64),
    phone       VARCHAR(20),
    status      TINYINT NOT NULL DEFAULT 1
);

CREATE TABLE wms_stocks (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    warehouse_id BIGINT UNSIGNED NOT NULL,
    sku_id      BIGINT UNSIGNED NOT NULL,
    qty         INT NOT NULL DEFAULT 0,
    safe_qty    INT NOT NULL DEFAULT 0 COMMENT '安全库存阈值',
    UNIQUE KEY uk_warehouse_sku (warehouse_id, sku_id)
);

CREATE TABLE wms_inbound (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    warehouse_id BIGINT UNSIGNED NOT NULL,
    type        TINYINT NOT NULL COMMENT '1采购入库 2退货入库',
    status      TINYINT NOT NULL COMMENT '1待入库 2已完成',
    remark      VARCHAR(255),
    created_at  DATETIME,
    INDEX idx_warehouse (warehouse_id)
);

CREATE TABLE wms_inbound_items (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    inbound_id  BIGINT UNSIGNED NOT NULL,
    sku_id      BIGINT UNSIGNED NOT NULL,
    qty         INT NOT NULL
);

CREATE TABLE wms_outbound (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    warehouse_id BIGINT UNSIGNED NOT NULL,
    type        TINYINT NOT NULL COMMENT '1订单出库 2调拨出库',
    ref_id      BIGINT UNSIGNED COMMENT '关联订单或调拨单ID',
    status      TINYINT NOT NULL,
    created_at  DATETIME
);

CREATE TABLE wms_outbound_items (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    outbound_id BIGINT UNSIGNED NOT NULL,
    sku_id      BIGINT UNSIGNED NOT NULL,
    qty         INT NOT NULL,
    INDEX idx_outbound (outbound_id)
);

CREATE TABLE wms_stock_logs (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    warehouse_id BIGINT UNSIGNED NOT NULL,
    sku_id      BIGINT UNSIGNED NOT NULL,
    type        VARCHAR(32) NOT NULL COMMENT 'inbound/outbound/transfer/inventory',
    qty         INT NOT NULL COMMENT '变动数量（正入负出）',
    before_qty  INT NOT NULL,
    after_qty   INT NOT NULL,
    ref_id      BIGINT UNSIGNED,
    ref_type    VARCHAR(32),
    created_at  DATETIME,
    INDEX idx_sku (sku_id),
    INDEX idx_warehouse (warehouse_id)
);
```

### 7.5 营销插件表（marketing）

```sql
CREATE TABLE coupons (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(64) NOT NULL,
    type        TINYINT NOT NULL COMMENT '1满减 2折扣 3无门槛',
    min_amount  DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '使用门槛',
    discount    DECIMAL(10,2) NOT NULL COMMENT '减免金额或折扣率',
    total_count INT NOT NULL DEFAULT 0 COMMENT '0不限',
    per_limit   INT NOT NULL DEFAULT 1,
    start_at    DATETIME,
    end_at      DATETIME,
    status      TINYINT NOT NULL DEFAULT 1
);

CREATE TABLE coupon_users (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    coupon_id   BIGINT UNSIGNED NOT NULL,
    user_id     BIGINT UNSIGNED NOT NULL,
    status      TINYINT NOT NULL COMMENT '1未使用 2已使用 3已过期',
    used_at     DATETIME,
    order_id    BIGINT UNSIGNED,
    created_at  DATETIME,
    INDEX idx_user (user_id)
);

CREATE TABLE activities (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    type        TINYINT NOT NULL COMMENT '1秒杀 2满减',
    name        VARCHAR(64) NOT NULL,
    config      JSON COMMENT '活动配置（满减规则等）',
    start_at    DATETIME,
    end_at      DATETIME,
    status      TINYINT NOT NULL DEFAULT 1
);

CREATE TABLE activity_products (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    activity_id BIGINT UNSIGNED NOT NULL,
    product_id  BIGINT UNSIGNED NOT NULL,
    sku_id      BIGINT UNSIGNED,
    activity_price DECIMAL(10,2),
    activity_stock INT,
    INDEX idx_activity (activity_id)
);

CREATE TABLE points_logs (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id     BIGINT UNSIGNED NOT NULL,
    type        TINYINT NOT NULL COMMENT '1消费获取 2兑换消耗 3管理员调整',
    points      INT NOT NULL COMMENT '正增负减',
    remark      VARCHAR(128),
    created_at  DATETIME,
    INDEX idx_user (user_id)
);
```

### 7.6 IM 客服插件表（im）

```sql
CREATE TABLE im_sessions (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id     BIGINT UNSIGNED NOT NULL,
    staff_id    BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0未分配',
    status      TINYINT NOT NULL COMMENT '1等待 2进行中 3已关闭',
    last_msg    VARCHAR(255),
    unread_count INT NOT NULL DEFAULT 0,
    created_at  DATETIME,
    updated_at  DATETIME,
    INDEX idx_user (user_id),
    INDEX idx_staff (staff_id)
);

CREATE TABLE im_messages (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    session_id  BIGINT UNSIGNED NOT NULL,
    sender_type TINYINT NOT NULL COMMENT '1用户 2客服',
    sender_id   BIGINT UNSIGNED NOT NULL,
    type        VARCHAR(32) NOT NULL COMMENT 'text/image/product_card/order_card',
    content     TEXT NOT NULL,
    extra       JSON COMMENT '商品卡片/订单卡片扩展数据',
    is_read     TINYINT NOT NULL DEFAULT 0,
    created_at  DATETIME,
    INDEX idx_session (session_id)
);

CREATE TABLE im_auto_replies (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    keyword     VARCHAR(128) NOT NULL,
    match_type  TINYINT NOT NULL COMMENT '1精确 2包含 3正则',
    reply       TEXT NOT NULL,
    sort        INT NOT NULL DEFAULT 0,
    status      TINYINT NOT NULL DEFAULT 1
);
```

### 7.7 AI 生图插件表（ai_image）

```sql
CREATE TABLE ai_models (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(64) NOT NULL COMMENT '显示名称',
    driver      VARCHAR(32) NOT NULL COMMENT 'tongyi/wenxin/hunyuan/openai',
    endpoint    VARCHAR(255),
    api_key     VARCHAR(255),
    params      JSON COMMENT '模型额外参数',
    is_default  TINYINT NOT NULL DEFAULT 0,
    status      TINYINT NOT NULL DEFAULT 1
);

CREATE TABLE ai_image_tasks (
    id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    model_id    BIGINT UNSIGNED NOT NULL,
    scene       VARCHAR(32) NOT NULL COMMENT 'carousel/detail',
    prompt      TEXT NOT NULL,
    neg_prompt  TEXT,
    params      JSON COMMENT '宽高、数量、风格等',
    status      TINYINT NOT NULL COMMENT '1生成中 2完成 3失败',
    result_urls JSON COMMENT '生成图片URL列表',
    error_msg   VARCHAR(255),
    created_at  DATETIME
);
```

### 7.8 装修插件表（decor）

```sql
CREATE TABLE decor_pages (
    id           BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    merchant_id  BIGINT UNSIGNED NOT NULL DEFAULT 0,
    page_key     VARCHAR(64) NOT NULL COMMENT 'index/category等',
    components   JSON NOT NULL COMMENT '组件配置列表',
    published_at DATETIME,
    created_at   DATETIME,
    updated_at   DATETIME,
    UNIQUE KEY uk_merchant_page (merchant_id, page_key)
);
```

---

## 8. API 设计规范

### 8.1 路由分组

| 前缀 | 用途 | 鉴权 |
|------|------|------|
| `/api/v1/` | 前台接口 | JWT 可选（部分需登录） |
| `/admin/api/` | 后台接口 | JWT 必须 + 权限校验 |
| `/ws/im` | WebSocket IM | JWT（query参数传入） |
| `/notify/` | 第三方异步回调 | 无JWT，签名验证 |

### 8.2 统一响应格式

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

### 8.3 错误码规范

| 范围 | 用途 |
|------|------|
| 0 | 成功 |
| 401 | 未登录 |
| 403 | 无权限 |
| 10001-19999 | 商品插件业务错误 |
| 20001-29999 | 订单插件业务错误 |
| 30001-39999 | 营销插件业务错误 |
| 40001-49999 | IM 插件业务错误 |
| 50001-59999 | WMS 插件业务错误 |

### 8.4 鉴权规范

- Token 类型：JWT，有效期 7 天，存 Redis 支持主动失效
- 前台：`Authorization: Bearer <token>`
- 后台：同上，Payload 额外携带 `permissions []string`
- Token 续签：有效期剩余 1 天时自动续签，响应头返回新 Token

### 8.5 核心接口列表（部分）

**前台接口：**
```
GET    /api/v1/index/decor          # 首页装修数据
GET    /api/v1/products             # 商品列表（分页、分类、关键词）
GET    /api/v1/products/:id         # 商品详情（含SKU）
POST   /api/v1/cart/add             # 加入购物车
GET    /api/v1/cart                 # 购物车列表
POST   /api/v1/orders               # 创建订单
POST   /api/v1/orders/:id/pay       # 发起支付
GET    /api/v1/orders               # 我的订单列表
POST   /api/v1/auth/sms/send        # 发送验证码
POST   /api/v1/auth/sms/login       # 手机号登录
POST   /api/v1/auth/wechat/miniapp  # 微信小程序登录
```

**后台接口：**
```
GET    /admin/api/dashboard         # 首页统计数据
GET    /admin/api/products          # 商品列表
POST   /admin/api/products          # 创建商品
PUT    /admin/api/products/:id      # 编辑商品
GET    /admin/api/orders            # 订单列表
PUT    /admin/api/orders/:id/ship   # 发货
GET    /admin/api/wms/stocks        # 库存列表
GET    /admin/api/im/sessions       # 客服会话列表
POST   /admin/api/ai/generate       # AI生图
GET    /admin/api/plugins           # 插件列表及状态
PUT    /admin/api/configs           # 系统配置保存
```

---

## 9. IM WebSocket 协议

### 9.1 连接

```
ws://host/ws/im?token=<jwt_token>
```

连接成功后服务端推送未读离线消息。

### 9.2 消息格式

**通用帧结构：**

```json
{
  "type": "msg | ack | typing | assign | close | ping",
  "session_id": "123",
  "payload": {}
}
```

**type=msg 时 payload：**

```json
{
  "msg_id": "uuid",
  "msg_type": "text | image | product_card | order_card",
  "content": "消息内容",
  "extra": {
    "product_id": 123,
    "product_title": "商品名",
    "product_cover": "https://...",
    "product_price": "99.00"
  }
}
```

**type=ack 时 payload（已读回执）：**

```json
{ "msg_id": "uuid" }
```

**type=assign（客服分配通知）：**

```json
{ "staff_id": 5, "staff_name": "客服小李", "staff_avatar": "https://..." }
```

### 9.3 Hub 并发设计

```go
type Hub struct {
    clients    map[string]*Client  // key: "user_{id}" 或 "staff_{id}"
    broadcast  chan *Message
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}
```

- 每个连接两个 goroutine：readPump（读）+ writePump（写）
- Hub 单 goroutine 管理连接注册/注销，避免竞态
- 心跳：客户端每 30s 发送 `ping`，服务端回 `pong`，60s 无心跳断开

### 9.4 离线消息处理

用户断线时消息存入 `im_messages.is_read=0`；重连时查询未读消息批量推送，推送完毕标记已读。

---

## 10. AI 生图模块

### 10.1 生图入口（商品管理页）

```
商品编辑页
├── 轮播图区域 → [AI生成轮播图] 按钮
│   输入：商品名称 + 风格（电商/写实/插画）→ 生成3-5张 → 勾选应用
└── 详情图区域 → [AI生成详情图] 按钮
    输入：卖点描述（多条）+ 尺寸（竖版750×1000）→ 生成长图 → 一键插入富文本
```

### 10.2 内置模型驱动

| 驱动名 | 服务商 | API |
|--------|--------|-----|
| `tongyi` | 阿里云通义万象 | DashScope 文生图 API |
| `wenxin` | 百度文心一格 | 百度AI开放平台 |
| `hunyuan` | 腾讯混元 | 腾讯云API |
| `openai` | OpenAI兼容 | DALL-E 3 及任意兼容接口 |

### 10.3 生图流程

```
前端提交生图请求
    → 后端创建 ai_image_tasks 记录（status=1）
    → 异步 goroutine 调用对应 Driver.Generate
    → 生成完成后更新 task（status=2，result_urls）
    → 前端轮询 /admin/api/ai/tasks/:id 查询状态
    → 完成后展示图片，用户勾选应用到商品
```

---

## 11. 管理后台 UI 规范

### 11.1 技术栈

Vue3 + Vite + TailwindCSS 3.x + shadcn-vue + Pinia + Vue Router 4

### 11.2 布局

```
┌─────────────────────────────────────────────────────┐
│  顶栏：Logo | 全局搜索 | 消息通知 | 用户头像/退出     │
├──────────────┬──────────────────────────────────────┤
│              │  面包屑导航                            │
│  侧边栏      │────────────────────────────────────── │
│  （深色）    │                                        │
│              │     主内容区                           │
│  菜单由后端  │     （白色卡片，圆角，微阴影）            │
│  动态下发    │                                        │
│  （基于已    │                                        │
│  启用插件）  │                                        │
└──────────────┴──────────────────────────────────────┘
```

### 11.3 设计风格

- 主色：深蓝 `#1e40af`，辅色：`#3b82f6`
- 背景：`#f8fafc`，卡片白色，圆角 `rounded-xl`，阴影 `shadow-sm`
- 侧边栏：深色（`#1e293b`），激活项高亮主色
- 表格：斑马纹，hover 行高亮
- 表单：侧滑 Drawer 而非弹窗，避免遮挡内容

### 11.4 通用组件

| 组件 | 说明 |
|------|------|
| `LyTable` | 带搜索/分页/批量操作/列配置/导出的通用表格 |
| `LyForm` | 基于 JSON Schema 自动渲染的通用表单 |
| `LyUpload` | 图片/文件上传，支持拖拽，接入存储驱动 |
| `LyRichEditor` | 富文本编辑器（商品详情，基于 WangEditor） |
| `LyStatCard` | 首页数据统计卡片（数值+趋势） |
| `LyChart` | 基于 ECharts 的图表组件 |
| `LyChatWindow` | IM 客服坐席窗口 |
| `LyAiImagePicker` | AI生图触发+结果选择组件 |

### 11.5 Dashboard 首屏

- 今日订单数 / 今日销售额 / 待处理售后 / 在线客服会话数（四格统计卡片）
- 7日销售趋势折线图（ECharts）
- 热销商品 TOP10 列表
- 库存预警商品（wms 插件启用时显示）
- 待分配客服会话（im 插件启用时显示）

---

## 12. uni-app 前端规范

### 12.1 技术栈

uni-app（Vue3）+ uview-plus 3.x + Pinia + uni-request

### 12.2 支持平台

微信小程序 / H5 / App（iOS + Android）

### 12.3 页面结构

```
app/
├── pages/
│   ├── index/         # 首页（decor 动态渲染）
│   ├── category/      # 分类页
│   ├── product/       # 商品详情页
│   ├── cart/          # 购物车
│   ├── order/
│   │   ├── confirm/   # 订单确认页
│   │   ├── list/      # 订单列表
│   │   └── detail/    # 订单详情
│   ├── user/          # 个人中心
│   └── login/         # 登录页
├── plugins/           # 各插件前台页面
│   ├── marketing/     # 秒杀、优惠券领取页
│   ├── im/            # 客服聊天页
│   └── wms/           # 物流追踪页
└── components/        # 基于 uview-plus 3.x 封装的业务组件
    ├── LyGoodCard/    # 商品卡片
    ├── LySkuPicker/   # SKU选择器
    ├── LyChatWindow/  # 客服聊天窗口
    └── LyDecorRender/ # 装修组件渲染器
```

### 12.4 核心页面说明

| 页面 | 关键组件 | 说明 |
|------|---------|------|
| 首页 | `u-swiper` + `LyDecorRender` | 从后端拉取装修配置动态渲染 |
| 商品详情 | `u-swiper` + `LySkuPicker` | AI生成轮播图展示，SKU规格选择 |
| 购物车 | `u-swipe-action` | 左滑删除，勾选结算，数量调整 |
| 订单确认 | `u-address` + `u-radio` | 地址选择 + 支付方式动态渲染 |
| IM客服 | `LyChatWindow` | WebSocket，文字/图片/商品卡片 |
| 个人中心 | `u-cell` + `u-avatar` | 订单状态快捷入口，积分，优惠券 |

### 12.5 多平台登录适配

| 平台 | 登录方式 |
|------|---------|
| 微信小程序 | `wechat_auth` 插件：wx.login → code2session |
| H5 | 手机号 + 短信验证码（`sms` 驱动） |
| App | 手机号 + 短信验证码，或微信 OAuth2 |

### 12.6 主题配置

通过 uview-plus `uni.scss` 全局变量配置主色，与管理后台保持一致：

```scss
$u-primary: #1e40af;
$u-primary-light: #3b82f6;
```

---

## 13. 部署方案

### 13.1 Docker Compose 结构

```yaml
version: '3.8'
services:
  nginx:
    image: nginx:1.25
    ports: ['80:80', '443:443']
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on: [server, admin, app_h5]

  server:
    build: ./server
    environment:
      - CONFIG_PATH=/app/config.yaml
    volumes:
      - ./config.yaml:/app/config.yaml
      - ./data/uploads:/app/uploads
    depends_on: [mysql, redis]

  admin:
    build: ./admin
    # Nginx 托管 Vue3 构建产物

  app_h5:
    build: ./app
    # Nginx 托管 uni-app H5 构建产物

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: lyshop
    volumes:
      - ./data/mysql:/var/lib/mysql

  redis:
    image: redis:7
    command: redis-server --appendonly yes
    volumes:
      - ./data/redis:/data
```

### 13.2 数据挂载

```
./data/mysql     → /var/lib/mysql       # MySQL 数据持久化
./data/redis     → /data                # Redis AOF 持久化
./data/uploads   → /app/uploads         # 本地存储插件文件目录
./config.yaml    → /app/config.yaml     # 主配置文件
./ssl/           → /etc/nginx/ssl       # SSL 证书
```

### 13.3 首次启动流程

```
docker compose up -d
    → MySQL 初始化
    → Go Server 启动
    → 执行数据库迁移（按 config.yaml 中 enabled 插件顺序）
    → 检查依赖关系
    → 注册路由和菜单
    → 初始化超级管理员账号（首次启动时）
    → 服务就绪
```

### 13.4 一键启动

```bash
cp config.example.yaml config.yaml
# 编辑 config.yaml：配置数据库密码、启用插件列表等
docker compose up -d
# 访问 http://localhost/admin 进入管理后台
```

---

## 附录：插件清单

| 插件名 | 类型 | 依赖 | 说明 |
|--------|------|------|------|
| `product` | 功能 | - | 商品管理 |
| `order` | 功能 | product | 订单管理 |
| `marketing` | 功能 | product, order | 营销活动 |
| `im` | 功能 | - | IM客服 |
| `wms` | 功能 | product | 仓储管理 |
| `ai_image` | 功能 | product | AI生图 |
| `decor` | 功能 | - | 店铺装修 |
| `wechat_pay` | 驱动实现 | - | 微信支付 |
| `alipay` | 驱动实现 | - | 支付宝支付 |
| `sms` | 驱动实现 | - | 短信服务 |
| `wechat_auth` | 驱动实现 | - | 微信登录 |
| `storage_local` | 驱动实现 | - | 本地存储 |
| `storage_oss` | 驱动实现 | - | 阿里云OSS |
| `storage_cos` | 驱动实现 | - | 腾讯云COS |
| `storage_qiniu` | 驱动实现 | - | 七牛云存储 |
