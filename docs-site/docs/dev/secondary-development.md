# 二次开发文档

本文档面向需要在 LYShop 基础上进行业务扩展、插件开发、定制化改造的开发者，覆盖后端插件体系、驱动接口、数据模型、API 规范、营销管线、前端开发、配置权限、测试部署等全部环节。每个章节均包含架构说明、完整代码示例和最佳实践。

---

## 目录

1. [概述与环境准备](#一、概述与环境准备)
2. [后端插件开发](#二、后端插件开发)
3. [后端驱动开发](#三、后端驱动开发)
4. [数据模型与迁移](#四、数据模型与迁移)
5. [API 开发规范](#五、api-开发规范)
6. [营销管线扩展](#六、营销管线扩展)
7. [Admin 后台前端开发](#七、admin-后台前端开发)
8. [PC 商城与 H5 前端开发](#八、pc-商城与-h5-前端开发)
9. [配置中心与权限系统](#九、配置中心与权限系统)
10. [测试、调试与部署](#十、测试调试与部署)

---

## 一、概述与环境准备

### 1.1 项目架构总览

LYShop 采用前后端分离架构，后端为 Go 单体服务，前端分为三个独立应用：

```
lyshop/
├── server/            # Go 后端（Gin + GORM + MySQL + Redis）
│   ├── main.go        # 入口，空导入启用的插件
│   ├── core/          # 框架核心
│   │   ├── app/       # 应用初始化与启动
│   │   ├── plugin/    # 插件接口、注册、加载
│   │   ├── driver/    # 驱动接口（存储/支付/短信/物流/AI/OAuth/配送）
│   │   ├── middleware/ # 中间件（CORS/日志/JWT认证）
│   │   ├── marketing/ # 营销价格管线
│   │   ├── response/  # 统一响应格式
│   │   ├── db/        # 数据库初始化
│   │   └── cache/     # Redis 初始化
│   ├── model/         # 核心数据模型
│   ├── plugins/       # 25+ 业务插件
│   ├── api/           # 核心路由（认证、管理员）
│   └── service/       # 核心业务逻辑
├── admin/             # Vue3 管理后台（TailwindCSS + Vite）
├── web/               # Vue3 PC 端商城（UnoCSS + Vite）
├── app/               # uni-app 移动端/小程序（UnoCSS + uview-plus）
└── docs-site/         # VitePress 文档站点
```

### 1.2 技术栈

| 层级 | 技术 | 版本 |
|------|------|------|
| **后端框架** | Go + Gin | Go 1.26+ |
| **ORM** | GORM | v1.31+ |
| **数据库** | MySQL | 8.0+ |
| **缓存** | Redis | 6.0+ |
| **认证** | JWT (HS256) | golang-jwt/v5 |
| **前端框架** | Vue 3 (Composition API) | 3.4+ |
| **构建工具** | Vite | 5.x |
| **状态管理** | Pinia | 2.1+ |
| **路由** | Vue Router | 4.3+ |
| **HTTP 客户端** | Axios / uni.request | 1.7+ |
| **Admin 样式** | TailwindCSS | 3.4+ |
| **Web 样式** | UnoCSS | - |
| **App 框架** | uni-app + uview-plus | 3.0+ |
| **文档** | VitePress | - |

### 1.3 开发环境搭建

#### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis（可选，不配则部分缓存功能不可用）

#### 后端启动

```bash
# 1. 进入后端目录
cd server

# 2. 复制配置文件并修改数据库/Redis 连接信息
cp ../config.example.yaml config.yaml

# 3. 安装依赖并启动
go mod tidy
go run main.go -config config.yaml
# 服务默认监听 :8080
```

#### 前端启动

```bash
# Admin 管理后台（端口 9527）
cd admin && npm install && npm run dev

# PC 商城（端口 9529）
cd web && npm install && npm run dev

# H5 移动端
cd app && npm install && npm run dev:h5

# 微信小程序
cd app && npm run dev:mp-weixin

# 演示模式（使用 Mock 数据，无需后端）
npm run dev:demo
```

#### 目录约定

开发新功能时，请遵循以下目录约定：

| 目录 | 说明 |
|------|------|
| `server/plugins/<name>/plugin.go` | 插件入口 |
| `server/plugins/<name>/plugin.json` | 插件元数据 |
| `server/plugins/<name>/api/` | HTTP 路由处理 |
| `server/plugins/<name>/service/` | 业务逻辑 |
| `server/plugins/<name>/model/` | 数据模型 |
| `admin/src/views/<module>/` | 管理端页面 |
| `admin/src/router/index.ts` | 管理端路由 |
| `web/src/views/` | PC 端页面 |
| `app/pages/<module>/` | H5/小程序页面 |

---

## 二、后端插件开发

### 2.1 插件体系架构

LYShop 的全部业务功能都以插件形式实现。每个插件是一个独立的 Go 包，通过 `init()` 函数自注册到全局注册表，由框架在启动时统一加载。

**插件生命周期：**

```
init() 注册  →  Load() 校验依赖  →  Migrate() 建表  →  RegisterRoutes() 注册路由  →  Install() 初始化
```

### 2.2 Plugin 接口

所有插件必须实现 `Plugin` 接口（定义在 `server/core/plugin/plugin.go`）：

```go
// Plugin 是所有插件必须实现的接口
type Plugin interface {
    // Meta 返回插件元数据（从嵌入的 plugin.json 解析）
    Meta() Meta
    // RegisterRoutes 注册前台和管理端 API 路由
    RegisterRoutes(front, admin *gin.RouterGroup)
    // Migrate 执行数据库建表/迁移（必须幂等）
    Migrate(db *gorm.DB) error
    // Install 在 Migrate 和 RegisterRoutes 之后调用一次
    Install() error
    // Uninstall 插件被禁用时调用
    Uninstall() error
}
```

### 2.3 Meta 元数据结构

```go
type Meta struct {
    Name        string        `json:"name"`         // 唯一标识符，如 "product"
    Title       string        `json:"title"`        // 显示名称，如 "商品插件"
    Version     string        `json:"version"`      // 语义化版本
    Description string        `json:"description"`  // 简述
    Author      string        `json:"author"`       // 作者
    Depends     []string      `json:"depends"`      // 依赖的插件列表
    Menus       []MenuItem    `json:"menus"`         // 管理端菜单
    Permissions []string      `json:"permissions"`   // 声明的权限标识
    ConfigItems []ConfigField `json:"config_items"`  // 配置中心字段声明
}
```

**MenuItem（管理端菜单项）：**

```go
type MenuItem struct {
    Title      string     `json:"title"`               // 菜单标题
    Icon       string     `json:"icon"`                // 图标名称
    Path       string     `json:"path"`                // 路由路径
    Sort       int        `json:"sort"`                // 排序权重（越小越靠前）
    Permission string     `json:"permission,omitempty"` // 所需权限（空 = 无需权限）
    Children   []MenuItem `json:"children,omitempty"`   // 子菜单
}
```

**ConfigField（配置中心字段）：**

```go
type ConfigField struct {
    Key         string `json:"key"`           // 配置键名
    Label       string `json:"label"`         // 显示标签
    Type        string `json:"type"`          // text|password|textarea|number|select|switch
    Placeholder string `json:"placeholder"`   // 输入提示
    Required    bool   `json:"required"`      // 是否必填
    Options     []struct {                    // type=select 时的选项
        Label string `json:"label"`
        Value string `json:"value"`
    } `json:"options,omitempty"`
}
```

### 2.4 插件注册与加载

**注册机制（`server/core/plugin/registry.go`）：**

```go
var registry []Plugin

// Register 将插件加入全局注册表，在 init() 中调用
func Register(p Plugin)

// Find 按名称查找插件
func Find(name string) Plugin

// All 返回所有已注册插件的快照
func All() []Plugin

// EnabledMenus 根据启用的插件和权限，返回过滤后的菜单树
func EnabledMenus(enabled []string, perms []string) []MenuItem
```

**加载流程（`server/core/plugin/loader.go`）：**

```go
func Load(enabled []string, db *gorm.DB, front, admin *gin.RouterGroup) error {
    // 1. 校验所有 enabled 插件已注册（否则报错提示加 blank import）
    // 2. 校验依赖链（Depends 中的插件必须也在 enabled 中）
    // 3. 按 config 中的顺序依次执行：
    //    - p.Migrate(db)         → 建表
    //    - p.RegisterRoutes(...)  → 注册路由
    //    - p.Install()           → 初始化数据
}
```

**启用方式（`server/main.go`）：**

```go
import (
    // 空导入使插件的 init() 被执行
    _ "github.com/ijry/lyshop/plugins/product"
    _ "github.com/ijry/lyshop/plugins/order"
    // 新增插件在此添加空导入
    _ "github.com/ijry/lyshop/plugins/my_plugin"
)
```

同时在 `config.yaml` 的 `plugins.enabled` 列表中添加插件名称。

### 2.5 完整示例：创建一个新插件

下面以创建一个「公告管理」插件为例，演示完整的插件开发流程。

#### 步骤 1：创建目录结构

```
server/plugins/announcement/
├── plugin.go          # 插件入口
├── plugin.json        # 插件元数据
├── api/
│   ├── front.go       # 前台路由
│   └── admin.go       # 管理端路由
├── model/
│   └── announcement.go # 数据模型
└── service/
    └── announcement.go # 业务逻辑
```

#### 步骤 2：编写 plugin.json

```json
{
  "name": "announcement",
  "title": "公告插件",
  "version": "1.0.0",
  "description": "站点公告管理，支持前台展示和后台 CRUD",
  "author": "developer",
  "depends": [],
  "menus": [
    {
      "title": "内容管理",
      "icon": "megaphone",
      "path": "/announcement",
      "sort": 50,
      "children": [
        {
          "title": "公告列表",
          "path": "/announcement/list",
          "permission": "announcement:view"
        }
      ]
    }
  ],
  "permissions": [
    "announcement:view",
    "announcement:create",
    "announcement:edit",
    "announcement:delete"
  ]
}
```

#### 步骤 3：定义数据模型

```go
// server/plugins/announcement/model/announcement.go
package model

import "github.com/ijry/lyshop/model"

type Announcement struct {
    model.Base
    Title   string `gorm:"size:255;not null"    json:"title"`
    Content string `gorm:"type:text"            json:"content"`
    Sort    int    `gorm:"not null;default:0"   json:"sort"`
    Status  int8   `gorm:"not null;default:1"   json:"status"` // 1=显示 0=隐藏
}
```

#### 步骤 4：编写业务逻辑

```go
// server/plugins/announcement/service/announcement.go
package service

import (
    "context"

    "github.com/ijry/lyshop/core/db"
    annomodel "github.com/ijry/lyshop/plugins/announcement/model"
)

// ListPublished 查询已发布的公告
func ListPublished(ctx context.Context) ([]annomodel.Announcement, error) {
    var list []annomodel.Announcement
    err := db.DB.WithContext(ctx).
        Where("status = ?", 1).
        Order("sort DESC, id DESC").
        Find(&list).Error
    return list, err
}

// ListAll 查询全部公告（管理端）
func ListAll(ctx context.Context, page, size int) ([]annomodel.Announcement, int64, error) {
    var list []annomodel.Announcement
    var total int64
    q := db.DB.WithContext(ctx).Model(&annomodel.Announcement{})
    q.Count(&total)
    err := q.Order("id DESC").
        Offset((page - 1) * size).Limit(size).
        Find(&list).Error
    return list, total, err
}

// Create 创建公告
func Create(ctx context.Context, a *annomodel.Announcement) error {
    return db.DB.WithContext(ctx).Create(a).Error
}

// Update 更新公告
func Update(ctx context.Context, id uint64, updates map[string]any) error {
    return db.DB.WithContext(ctx).
        Model(&annomodel.Announcement{}).
        Where("id = ?", id).
        Updates(updates).Error
}

// Delete 删除公告
func Delete(ctx context.Context, id uint64) error {
    return db.DB.WithContext(ctx).
        Where("id = ?", id).
        Delete(&annomodel.Announcement{}).Error
}
```

#### 步骤 5：注册路由

**前台路由：**

```go
// server/plugins/announcement/api/front.go
package api

import (
    "github.com/gin-gonic/gin"
    "github.com/ijry/lyshop/core/response"
    annosvc "github.com/ijry/lyshop/plugins/announcement/service"
)

func RegisterFrontRoutes(g *gin.RouterGroup) {
    g.GET("/announcements", listPublished)
}

func listPublished(c *gin.Context) {
    list, err := annosvc.ListPublished(c.Request.Context())
    if err != nil {
        response.Fail(c, 500, err.Error())
        return
    }
    response.OK(c, list)
}
```

**管理端路由：**

```go
// server/plugins/announcement/api/admin.go
package api

import (
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/ijry/lyshop/core/middleware"
    "github.com/ijry/lyshop/core/response"
    annomodel "github.com/ijry/lyshop/plugins/announcement/model"
    annosvc "github.com/ijry/lyshop/plugins/announcement/service"
)

func RegisterAdminRoutes(g *gin.RouterGroup) {
    g.GET("/announcements",
        middleware.RequirePermission("announcement:view"), adminList)
    g.POST("/announcements",
        middleware.RequirePermission("announcement:create"), adminCreate)
    g.PUT("/announcements/:id",
        middleware.RequirePermission("announcement:edit"), adminUpdate)
    g.DELETE("/announcements/:id",
        middleware.RequirePermission("announcement:delete"), adminDelete)
}

func adminList(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
    list, total, err := annosvc.ListAll(c.Request.Context(), page, size)
    if err != nil {
        response.Fail(c, 500, err.Error())
        return
    }
    response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}

func adminCreate(c *gin.Context) {
    var anno annomodel.Announcement
    if err := c.ShouldBindJSON(&anno); err != nil {
        response.Fail(c, 400, err.Error())
        return
    }
    if err := annosvc.Create(c.Request.Context(), &anno); err != nil {
        response.Fail(c, 500, err.Error())
        return
    }
    response.OK(c, anno)
}

func adminUpdate(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    var updates map[string]any
    c.ShouldBindJSON(&updates)
    if err := annosvc.Update(c.Request.Context(), id, updates); err != nil {
        response.Fail(c, 500, err.Error())
        return
    }
    response.OK(c, nil)
}

func adminDelete(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    if err := annosvc.Delete(c.Request.Context(), id); err != nil {
        response.Fail(c, 500, err.Error())
        return
    }
    response.OK(c, nil)
}
```

#### 步骤 6：编写插件入口

```go
// server/plugins/announcement/plugin.go
package announcement

import (
    _ "embed"
    "encoding/json"

    "github.com/gin-gonic/gin"
    "github.com/ijry/lyshop/core/plugin"
    annoapi "github.com/ijry/lyshop/plugins/announcement/api"
    annomodel "github.com/ijry/lyshop/plugins/announcement/model"
    "gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type annoPlugin struct {
    meta plugin.Meta
}

func init() {
    var m plugin.Meta
    if err := json.Unmarshal(metaJSON, &m); err != nil {
        panic("announcement plugin: invalid plugin.json: " + err.Error())
    }
    plugin.Register(&annoPlugin{meta: m})
}

func (p *annoPlugin) Meta() plugin.Meta { return p.meta }

func (p *annoPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
    annoapi.RegisterFrontRoutes(front)
    annoapi.RegisterAdminRoutes(admin)
}

func (p *annoPlugin) Migrate(db *gorm.DB) error {
    return db.AutoMigrate(&annomodel.Announcement{})
}

func (p *annoPlugin) Install() error   { return nil }
func (p *annoPlugin) Uninstall() error { return nil }
```

#### 步骤 7：注册并启用

**在 `server/main.go` 添加空导入：**

```go
_ "github.com/ijry/lyshop/plugins/announcement"
```

**在 `config.yaml` 添加启用项：**

```yaml
plugins:
  enabled:
    - product
    - order
    # ...
    - announcement   # 新增
```

重启后端服务，插件自动建表、注册路由、挂载菜单。

### 2.6 最佳实践

- **Migrate 必须幂等**：使用 `db.AutoMigrate()`，它只添加不存在的列和索引，不会破坏已有数据
- **依赖声明要准确**：如果你的插件用到了 product 的 model 或 service，务必在 `depends` 中声明
- **按 config 顺序加载**：`plugins.enabled` 列表的顺序即加载顺序，被依赖的插件需排在前面
- **权限粒度明确**：`permission` 字段使用 `模块:操作` 格式，如 `announcement:view`、`announcement:create`
- **每个插件一个 Go 包**：避免跨插件直接引用 model/service，需要时通过接口或事件解耦

---

## 三、后端驱动开发

### 3.1 驱动体系概述

LYShop 的驱动（Driver）系统提供了对外部服务的可插拔抽象。每种驱动类型定义统一接口，不同实现以插件形式注册。系统启动时自动选择配置的驱动。

当前支持 7 类驱动：

| 驱动类型 | 接口目录 | 用途 | 已有实现 |
|----------|----------|------|----------|
| Storage | `core/driver/storage/` | 文件上传 | local, oss, cos, qiniu |
| Payment | `core/driver/payment/` | 支付网关 | wechat_pay, alipay |
| SMS | `core/driver/sms/` | 短信发送 | aliyun, tencent |
| Logistics | `core/driver/logistics/` | 物流追踪 | kuaidi100, kdniao |
| AI | `core/driver/ai/` | AI 图片生成 | tongyi, wenxin, dalle |
| OAuth | `core/driver/oauth/` | 第三方登录 | wechat |
| Delivery | `core/driver/delivery/` | 配送方式 | express, local |

### 3.2 存储驱动

**接口定义（`server/core/driver/storage/storage.go`）：**

```go
type UploadResult struct {
    Path string `json:"path"` // 存储路径
    URL  string `json:"url"`  // 公开访问 URL
    Size int64  `json:"size"` // 文件大小（字节）
    Mime string `json:"mime"` // MIME 类型
}

type Driver interface {
    Name() string
    Upload(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error)
    Delete(ctx context.Context, path string) error
    GetURL(path string) string
}
```

**注册表特点：**
- 支持名称别名（`storage_local` 和 `local` 等价）
- 支持默认驱动选择（`SetDefault` / `GetByName("")`）
- 通过 `storage_router` 插件可实现多驱动路由

**示例：实现一个自定义存储驱动**

```go
// server/plugins/storage_minio/driver.go
package storage_minio

import (
    "context"
    "fmt"
    "mime/multipart"
    "path/filepath"
    "time"

    "github.com/ijry/lyshop/core/driver/storage"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

type minioDriver struct {
    client   *minio.Client
    bucket   string
    endpoint string
}

func NewDriver(endpoint, accessKey, secretKey, bucket string) (*minioDriver, error) {
    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: true,
    })
    if err != nil {
        return nil, err
    }
    return &minioDriver{client: client, bucket: bucket, endpoint: endpoint}, nil
}

func (d *minioDriver) Name() string { return "minio" }

func (d *minioDriver) Upload(ctx context.Context, fh *multipart.FileHeader) (*storage.UploadResult, error) {
    file, err := fh.Open()
    if err != nil {
        return nil, err
    }
    defer file.Close()

    key := fmt.Sprintf("uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fh.Filename))
    _, err = d.client.PutObject(ctx, d.bucket, key, file, fh.Size, minio.PutObjectOptions{
        ContentType: fh.Header.Get("Content-Type"),
    })
    if err != nil {
        return nil, err
    }
    return &storage.UploadResult{
        Path: key,
        URL:  d.GetURL(key),
        Size: fh.Size,
        Mime: fh.Header.Get("Content-Type"),
    }, nil
}

func (d *minioDriver) Delete(ctx context.Context, path string) error {
    return d.client.RemoveObject(ctx, d.bucket, path, minio.RemoveObjectOptions{})
}

func (d *minioDriver) GetURL(path string) string {
    return fmt.Sprintf("https://%s/%s/%s", d.endpoint, d.bucket, path)
}
```

然后在插件的 `Install()` 中读取配置并注册：

```go
func (p *minioPlugin) Install() error {
    // 从 ConfigKV 读取配置
    cfg := readConfig("storage_minio")
    driver, err := NewDriver(cfg["endpoint"], cfg["access_key"], cfg["secret_key"], cfg["bucket"])
    if err != nil {
        return err
    }
    storage.Register(driver)
    return nil
}
```

### 3.3 支付驱动

**接口定义（`server/core/driver/payment/payment.go`）：**

```go
type OrderParams struct {
    OrderNo     string // 商户订单号
    Amount      int64  // 金额（分），如 9900 = ¥99.00
    Description string // 订单描述
    NotifyURL   string // 回调通知地址
    OpenID      string // 微信小程序 JSAPI 需要
    ClientIP    string // H5/App 需要
}

type OrderResult struct {
    PrepayID  string            // 预支付 ID
    PayParams map[string]string // 传给前端 SDK 的参数
}

type QueryResult struct {
    OutTradeNo string // 商户订单号
    TradeNo    string // 平台交易号
    Status     string // "paid" | "unpaid" | "refunded"
    Amount     int64  // 金额（分）
}

type RefundParams struct {
    OrderNo     string // 原订单号
    RefundNo    string // 退款单号
    Amount      int64  // 退款金额
    TotalAmount int64  // 原订单总额
    Reason      string // 退款原因
}

type Driver interface {
    Name() string
    CreateOrder(ctx context.Context, p *OrderParams) (*OrderResult, error)
    QueryOrder(ctx context.Context, tradeNo string) (*QueryResult, error)
    Refund(ctx context.Context, p *RefundParams) (*RefundResult, error)
    HandleNotify(ctx context.Context, r *http.Request) (*NotifyResult, error)
}
```

**注册方式：**

```go
// 在插件的 init() 或 Install() 中
payment.Register(myPaymentDriver)

// 使用时按名称获取
driver, err := payment.Get("wechat_pay")
```

### 3.4 短信驱动

**接口定义（`server/core/driver/sms/sms.go`）：**

```go
type Driver interface {
    Name() string
    Send(ctx context.Context, phone, templateCode string, params map[string]string) error
}
```

短信驱动为单实例注册（最后注册的生效）：

```go
sms.Register(myDriver)

// 使用
driver, err := sms.Get()
driver.Send(ctx, "13800138000", "SMS_CODE_TPL", map[string]string{"code": "1234"})
```

### 3.5 物流驱动

**接口定义（`server/core/driver/logistics/logistics.go`）：**

```go
type QueryReq struct {
    CompanyCode string // 快递公司编码
    TrackingNo  string // 运单号
    Phone       string // 收件人手机尾号（部分平台需要）
}

type TrackNode struct {
    Time       string          // 时间
    Location   string          // 地点
    StatusCode string          // 状态码
    StatusText string          // 状态文本
    RawPayload json.RawMessage // 原始数据
}

type TrackResult struct {
    Provider   string      // 驱动名称
    StatusCode string      // 总状态码
    StatusText string      // 总状态文本
    SignedAt   *time.Time  // 签收时间
    Nodes      []TrackNode // 轨迹节点
}

type Driver interface {
    Name() string
    Query(ctx context.Context, req QueryReq) (*TrackResult, error)
}
```

**物流驱动支持主备切换：**

```go
// 设置主备驱动
logistics.SetDefaultDrivers("kuaidi100", "kdniao")

// 自动 fallback
driver, name, err := logistics.ResolveByPinnedOrFallback(pinnedDriver)
```

### 3.6 自定义驱动开发流程

1. **定义驱动实现**：在 `server/plugins/<driver_name>/` 下实现对应 Driver 接口
2. **编写 plugin.json**：声明 config_items 让管理员在配置中心填写密钥等信息
3. **在 Install() 中注册**：读取 ConfigKV，创建驱动实例，调用 `driver.Register()`
4. **在 main.go 空导入**：添加 `_ "github.com/ijry/lyshop/plugins/<driver_name>"`
5. **在 config.yaml 启用**：添加到 `plugins.enabled` 列表

### 3.7 最佳实践

- **驱动接口保持精简**：方法不超过 5 个，只抽象核心操作
- **配置从 ConfigKV 读取**：不要硬编码密钥或地址，利用配置中心
- **错误信息要脱敏**：返回给前端的错误不应包含第三方 SDK 的敏感信息
- **Router 插件模式**：多个同类驱动通过 router 插件统一管理（参考 `storage_router`、`logistics_router`）

---

## 四、数据模型与迁移

### 4.1 Base 模型

所有模型嵌入 `model.Base`（定义在 `server/model/base.go`），自动提供主键和时间戳：

```go
type Base struct {
    ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 4.2 核心模型参考

**User（用户）：**

```go
type User struct {
    Base
    MerchantID uint64 `gorm:"not null;default:0;index" json:"merchant_id"` // 多商户预留
    Phone      string `gorm:"size:20;uniqueIndex"      json:"phone"`
    Nickname   string `gorm:"size:64"                  json:"nickname"`
    Avatar     string `gorm:"size:500"                 json:"avatar"`
    Points     int    `gorm:"not null;default:0"       json:"points"`
    Status     int8   `gorm:"not null;default:1"       json:"status"` // 1=活跃 0=禁用
}
```

**Admin（管理员）：**

```go
type Admin struct {
    Base
    Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
    Password string `gorm:"size:255;not null"            json:"-"` // bcrypt 哈希，json 忽略
    RoleID   uint64 `gorm:"not null"                     json:"role_id"`
    Status   int8   `gorm:"not null;default:1"           json:"status"`
}
```

**Role（角色）：**

```go
type Role struct {
    Base
    Name        string          `gorm:"size:64;not null" json:"name"`
    Permissions json.RawMessage `gorm:"type:json"        json:"permissions"` // []string
}
```

**ConfigKV（配置键值对）：**

```go
type ConfigKV struct {
    ID     uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    Plugin string `gorm:"size:64;not null;uniqueIndex:uk_plugin_key"  json:"plugin"`
    Key    string `gorm:"size:128;not null;uniqueIndex:uk_plugin_key" json:"key"`
    Value  string `gorm:"type:text"                                  json:"value"`
}

func (ConfigKV) TableName() string { return "configs" }
```

### 4.3 插件模型示例

以商品模型为例（`server/plugins/product/model/product.go`）：

```go
type Product struct {
    model.Base
    MerchantID    uint64          `gorm:"not null;default:0;index"    json:"merchant_id"`
    CategoryID    uint64          `gorm:"not null;index"               json:"category_id"`
    Title         string          `gorm:"size:255;not null"            json:"title"`
    Subtitle      string          `gorm:"size:255"                     json:"subtitle"`
    Cover         string          `gorm:"size:500"                     json:"cover"`
    Price         float64         `gorm:"type:decimal(10,2);not null"  json:"price"`
    OriginPrice   float64         `gorm:"type:decimal(10,2)"           json:"origin_price"`
    Stock         int             `gorm:"not null;default:0"           json:"stock"`
    Sales         int             `gorm:"not null;default:0"           json:"sales"`
    FavoriteCount int             `gorm:"not null;default:0"           json:"favorite_count"`
    Status        int8            `gorm:"not null;default:1"           json:"status"`
    Sort          int             `gorm:"not null;default:0"           json:"sort"`
    Detail        json.RawMessage `gorm:"type:json"                    json:"detail"`
}
```

### 4.4 GORM Tag 规范

| Tag | 说明 | 示例 |
|-----|------|------|
| `gorm:"primaryKey"` | 主键 | ID 字段 |
| `gorm:"autoIncrement"` | 自增 | ID 字段 |
| `gorm:"size:N"` | varchar 长度 | `size:255` |
| `gorm:"type:xxx"` | 指定列类型 | `type:decimal(10,2)`、`type:text`、`type:json` |
| `gorm:"not null"` | 非空约束 | |
| `gorm:"default:N"` | 默认值 | `default:0`、`default:1` |
| `gorm:"uniqueIndex"` | 唯一索引 | 单字段唯一 |
| `gorm:"uniqueIndex:name"` | 复合唯一索引 | `uniqueIndex:uk_plugin_key` |
| `gorm:"index"` | 普通索引 | |
| `json:"-"` | JSON 序列化忽略 | Password 字段 |
| `json:"field_name"` | JSON 字段名 | 使用 snake_case |

### 4.5 迁移方式

LYShop 使用 GORM 的 `AutoMigrate` 进行迁移，在插件的 `Migrate()` 方法中调用：

```go
func (p *myPlugin) Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &mymodel.ModelA{},
        &mymodel.ModelB{},
    )
}
```

**AutoMigrate 行为：**

- 表不存在 → 创建表
- 表存在但缺少列 → 添加列
- 列存在但类型变更 → 尝试修改（不删除列）
- 索引不存在 → 创建索引
- **不会删除已有列或数据**

### 4.6 关联关系

**一对多：**

```go
// 一个分类有多个商品
type Category struct {
    model.Base
    Name     string    `gorm:"size:64;not null" json:"name"`
    Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

type Product struct {
    model.Base
    CategoryID uint64 `gorm:"not null;index" json:"category_id"`
    // ...
}

// 查询时预加载
db.Preload("Products").Find(&categories)
```

**多对多（使用中间表）：**

```go
// 商品收藏关系
type ProductFavorite struct {
    model.Base
    UserID    uint64 `gorm:"not null;uniqueIndex:uk_user_product" json:"user_id"`
    ProductID uint64 `gorm:"not null;uniqueIndex:uk_user_product" json:"product_id"`
}
```

### 4.7 最佳实践

- **所有模型嵌入 `model.Base`**：保持 ID、时间戳的一致性
- **字段命名使用 snake_case**：与 JSON 输出保持一致
- **字符串字段指定 size**：避免默认 varchar(255) 浪费空间
- **状态字段使用 int8**：1=正常/启用，0=禁用/隐藏，预留扩展空间
- **金额使用 `decimal(10,2)`**：避免浮点精度问题
- **JSON 字段使用 `json.RawMessage`**：灵活存储结构化数据
- **敏感字段加 `json:"-"`**：如密码不输出到 API 响应

---

## 五、API 开发规范

### 5.1 统一响应格式

所有 API 使用统一的 JSON 信封格式（`server/core/response/response.go`）：

```go
type R struct {
    Code int    `json:"code"` // 0=成功，非0=业务错误
    Msg  string `json:"msg"`  // 提示信息
    Data any    `json:"data"` // 业务数据
}
```

**响应函数：**

```go
// 成功响应
response.OK(c, data)
// 输出: {"code": 0, "msg": "success", "data": {...}}

// 业务错误
response.Fail(c, 10001, "商品不存在")
// 输出: {"code": 10001, "msg": "商品不存在", "data": null}

// 用于 AbortWithStatusJSON 的错误
response.Err(401, "请先登录")
```

**分页数据：**

```go
response.OK(c, response.PageData{
    List:  list,
    Total: total,
    Page:  page,
    Size:  size,
})
// 输出: {"code": 0, "msg": "success", "data": {"list": [...], "total": 100, "page": 1, "size": 20}}
```

### 5.2 路由组织

路由分为三个组，在 `server/core/app/app.go` 中创建：

```go
// 前台公开路由（用户端 API）
front := r.Group("/api/v1")

// 管理端公开路由（如登录）
adminPublic := r.Group("/admin/api")

// 管理端认证路由（需要 JWT + admin 角色）
adminAuth := r.Group("/admin/api")
adminAuth.Use(middleware.RequireAdmin)
```

**插件中注册路由：**

```go
func (p *myPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
    // front 是 /api/v1 组
    // admin 是 /admin/api 组（已包含 RequireAdmin 中间件）
    myapi.RegisterFrontRoutes(front)
    myapi.RegisterAdminRoutes(admin)
}
```

### 5.3 前台路由示例

参考商品插件的前台路由（`server/plugins/product/api/front.go`）：

```go
func RegisterFrontRoutes(g *gin.RouterGroup) {
    // 公开路由 — 无需登录
    g.GET("/categories", listCategories)
    g.GET("/products", listProducts)
    g.GET("/products/:id", getProduct)

    // 认证路由 — 需要用户登录
    auth := g.Group("")
    auth.Use(middleware.RequireAuth)
    auth.POST("/products/:id/favorite", favoriteProduct)
    auth.DELETE("/products/:id/favorite", unfavoriteProduct)
    auth.GET("/user/favorites", listUserFavorites)
}
```

### 5.4 管理端路由示例

参考商品插件的管理端路由（`server/plugins/product/api/admin.go`）：

```go
func RegisterAdminRoutes(g *gin.RouterGroup) {
    // 每个接口单独声明权限
    g.GET("/products",     middleware.RequirePermission("product:view"),   adminListProducts)
    g.POST("/products",    middleware.RequirePermission("product:create"), adminCreateProduct)
    g.PUT("/products/:id", middleware.RequirePermission("product:edit"),   adminUpdateProduct)
    g.DELETE("/products/:id", middleware.RequirePermission("product:delete"), adminDeleteProduct)
}
```

### 5.5 认证中间件

JWT 认证定义在 `server/core/middleware/auth.go`：

```go
// JWT Claims 结构
type Claims struct {
    UserID uint64   `json:"user_id"`
    Role   string   `json:"role"`  // "user" | "admin"
    Perms  []string `json:"perms"` // 权限列表
    jwt.RegisteredClaims
}

// 生成 Token
token, err := middleware.GenerateToken(userID, "admin", []string{"product:view", "product:create"})

// 解析 Token
claims, err := middleware.ParseToken(tokenString)
```

**三层认证中间件：**

```go
// RequireAuth — 校验 JWT，设置 user_id/role/perms 到 Context
middleware.RequireAuth

// RequireAdmin — 在 RequireAuth 基础上检查 role == "admin"
middleware.RequireAdmin

// RequirePermission — 检查具体权限（支持 "*" 通配符）
middleware.RequirePermission("product:view")
```

**在 Handler 中获取用户信息：**

```go
func myHandler(c *gin.Context) {
    // 获取当前用户 ID
    userID, _ := c.Get("user_id")
    uid := userID.(uint64)

    // 获取角色
    role, _ := c.Get("role")

    // 获取权限列表
    perms, _ := c.Get("perms")
    permList := perms.([]string)
}
```

### 5.6 请求参数处理

**Query 参数（GET）：**

```go
// 定义查询结构
type ProductListQuery struct {
    Page       int    `form:"page"`
    Size       int    `form:"size"`
    Keyword    string `form:"keyword"`
    CategoryID uint64 `form:"category_id"`
}

func listProducts(c *gin.Context) {
    var q ProductListQuery
    c.ShouldBindQuery(&q)
    // 使用 q.Page, q.Keyword 等
}
```

**Path 参数：**

```go
id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
```

**JSON Body（POST/PUT）：**

```go
// 绑定到结构体
var req CreateRequest
if err := c.ShouldBindJSON(&req); err != nil {
    response.Fail(c, 400, err.Error())
    return
}

// 绑定到 map（灵活更新）
var updates map[string]any
c.ShouldBindJSON(&updates)
```

**文件上传：**

```go
fh, err := c.FormFile("file")
if err != nil {
    response.Fail(c, 400, "请选择文件")
    return
}
driver, _ := storage.Get()
result, err := driver.Upload(c.Request.Context(), fh)
```

### 5.7 错误码约定

| 范围 | 说明 |
|------|------|
| 0 | 成功 |
| 400 | 参数校验错误 |
| 401 | 未登录 / Token 无效 |
| 403 | 无权限 |
| 500 | 服务器内部错误 |
| 10001-19999 | 业务错误（各插件自定义） |

### 5.8 最佳实践

- **HTTP 状态码始终 200**：业务错误通过 `code` 字段区分，前端统一处理
- **前台和管理端路由分离**：分别注册到 `front` 和 `admin` 路由组
- **权限逐接口声明**：每个管理端接口单独使用 `RequirePermission`
- **Service 层不感知 HTTP**：Handler 解析参数、调用 Service、返回响应，Service 只接受纯数据
- **分页接口返回 PageData**：统一 list/total/page/size 结构
- **Context 传递**：使用 `c.Request.Context()` 传递给 Service 和 DB 操作

---

## 六、营销管线扩展

### 6.1 管线架构

LYShop 的营销价格计算使用管线（Pipeline）模式。多个 `PriceCalculator` 按优先级排序执行，依次对 `PriceContext` 中的价格进行修改。

```
商品原价 → 活动折扣 → VIP 折扣 → 满减折扣 → 优惠券折扣 → 积分抵扣 → 最终价格
```

### 6.2 PriceCalculator 接口

定义在 `server/core/marketing/pipeline.go`：

```go
type PriceCalculator interface {
    // Name 返回可读标识
    Name() string
    // Priority 控制执行顺序（数值越小越先执行）
    Priority() int
    // Calculate 修改 ctx 中的价格数据
    // 返回 continueNext=false 可中断管线（如排他性活动）
    Calculate(ctx *PriceContext) (continueNext bool, err error)
}
```

**注册方式：**

```go
// 在 init() 中注册，自动按 Priority 排序
marketing.Register(myCalculator)
```

**管线执行流程（`Calculate` 函数）：**

```go
func Calculate(ctx *PriceContext) error {
    // 1. 计算商品小计
    ctx.GoodsAmount = sum(item.Price * item.Qty)
    ctx.FinalAmount = ctx.GoodsAmount

    // 2. 按优先级执行各 Calculator
    for _, calc := range calculators {
        cont, err := calc.Calculate(ctx)
        // 每步重新计算最终价格
        ctx.FinalAmount = ctx.GoodsAmount
            - ctx.ActivityDiscount
            - ctx.VipDiscount
            - ctx.FullReduceDiscount
            - ctx.CouponDiscount
            - ctx.PointsDiscount
        if !cont { break }
    }

    // 3. 最终价格不低于 0
    if ctx.FinalAmount < 0 { ctx.FinalAmount = 0 }
}
```

### 6.3 PriceContext 数据结构

定义在 `server/core/marketing/context.go`：

```go
type PriceContext struct {
    // === 输入 ===
    UserID          uint64
    Items           []OrderItem            // 订单商品行
    CouponIDs       []uint64               // 用户选择的优惠券
    ActivityID      uint64                 // 指定活动（0=自动检测）
    PointsUse       int                    // 用户选择使用的积分数
    IsVIP           bool                   // 是否 VIP
    VIPLevelID      uint64                 // VIP 等级 ID
    ItemVIPDiscount map[uint64]float64     // sku_id -> VIP 折扣金额

    // === 计算结果（各步骤逐步填充）===
    GoodsAmount        float64             // 商品小计
    ActivityDiscount   float64             // 活动折扣
    VipDiscount        float64             // VIP 折扣
    FullReduceDiscount float64             // 满减折扣
    CouponDiscount     float64             // 优惠券折扣
    PointsDiscount     float64             // 积分抵扣
    FinalAmount        float64             // 最终价格

    AppliedRules []AppliedRule              // 已应用的规则记录
    Commissions  []Commission              // 分销返佣（不影响 FinalAmount）
}

type OrderItem struct {
    ProductID     uint64
    SkuID         uint64
    Title         string
    Price         float64 // 原始单价
    Qty           int
    ActivityPrice float64 // 活动价（0 表示无）
    VipPrice      float64 // VIP 价（0 表示无）
}

type AppliedRule struct {
    Type     string  `json:"type"`     // activity|vip|coupon|full_reduce|points
    Name     string  `json:"name"`     // 可读标签
    Discount float64 `json:"discount"` // 折扣金额（正数）
}
```

### 6.4 完整示例：实现满减计算器

以下实现「满 200 减 30」的满减规则：

```go
// server/plugins/marketing/calculator/full_reduce.go
package calculator

import (
    "github.com/ijry/lyshop/core/marketing"
)

type fullReduceCalc struct{}

func init() {
    marketing.Register(&fullReduceCalc{})
}

func (c *fullReduceCalc) Name() string     { return "满减折扣" }
func (c *fullReduceCalc) Priority() int    { return 300 } // 在活动(100)、VIP(200) 之后

func (c *fullReduceCalc) Calculate(ctx *marketing.PriceContext) (bool, error) {
    // 计算当前已扣折扣后的实付金额
    currentAmount := ctx.GoodsAmount - ctx.ActivityDiscount - ctx.VipDiscount

    // 示例规则：满 200 减 30，满 500 减 100
    rules := []struct {
        Threshold float64
        Discount  float64
    }{
        {500, 100},
        {200, 30},
    }

    for _, rule := range rules {
        if currentAmount >= rule.Threshold {
            ctx.FullReduceDiscount = rule.Discount
            ctx.AppliedRules = append(ctx.AppliedRules, marketing.AppliedRule{
                Type:     "full_reduce",
                Name:     fmt.Sprintf("满%.0f减%.0f", rule.Threshold, rule.Discount),
                Discount: rule.Discount,
            })
            break // 只命中最高一档
        }
    }

    return true, nil // continueNext=true，允许后续折扣继续
}
```

### 6.5 优先级推荐

| 优先级 | 阶段 | 说明 |
|--------|------|------|
| 100 | 活动折扣 | 秒杀、团购等排他活动 |
| 200 | VIP 折扣 | 会员专属价 |
| 300 | 满减折扣 | 满减规则 |
| 400 | 优惠券折扣 | 用户手动选择使用 |
| 500 | 积分抵扣 | 积分兑换 |
| 900 | 分销佣金 | 计算返佣（不影响支付价格） |

### 6.6 最佳实践

- **排他性活动**：返回 `continueNext=false` 可阻止后续折扣（如秒杀不可叠加优惠券）
- **AppliedRules 必须记录**：所有折扣都需追加到 `AppliedRules`，用于订单详情展示
- **折扣字段互不干扰**：每个 Calculator 只修改自己对应的字段（ActivityDiscount / VipDiscount / ...）
- **管线中不做数据库写操作**：管线是纯计算，数据持久化由订单模块负责
- **在实际项目中应从数据库读取规则**：示例中硬编码仅作演示

---

## 七、Admin 后台前端开发

### 7.1 技术架构

| 项目 | 说明 |
|------|------|
| 框架 | Vue 3 + TypeScript |
| 构建 | Vite 5 |
| 样式 | TailwindCSS 3.4 |
| 状态 | Pinia |
| 路由 | Vue Router 4 (hash/history) |
| 图标 | lucide-vue-next |
| 图表 | ECharts 5 |
| API | Axios |

**目录结构：**

```
admin/src/
├── api/              # API 模块
│   ├── request.ts    # Axios 封装
│   ├── auth.ts       # 认证相关 API
│   └── plugins.ts    # 插件相关 API
├── components/       # 通用组件
├── layouts/
│   └── AdminLayout.vue # 管理后台主布局
├── mock/             # Mock 数据
├── router/
│   └── index.ts      # 路由配置
├── stores/
│   └── auth.ts       # 认证状态
├── types/            # TypeScript 类型定义
├── utils/            # 工具函数
│   ├── notify.ts     # 通知处理
│   ├── toast.ts      # Toast 提示
│   └── dialog.ts     # 对话框
└── views/            # 页面组件（按模块分目录）
```

### 7.2 路由注册

在 `admin/src/router/index.ts` 中注册新页面：

```typescript
import { createRouter, createWebHistory, createWebHashHistory } from 'vue-router'

const isMock = import.meta.env.VITE_MOCK === 'true'

const router = createRouter({
  // Mock 模式用 hash，生产用 history + /admin 前缀
  history: isMock ? createWebHashHistory() : createWebHistory('/admin'),
  routes: [
    { path: '/login', component: () => import('@/views/Login.vue') },
    {
      path: '/',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', name: '首页', component: () => import('@/views/Dashboard.vue') },

        // 示例：为公告插件添加管理端页面
        { path: 'announcement/list', name: '公告列表',
          component: () => import('@/views/announcement/AnnouncementList.vue') },
        { path: 'announcement/form', name: '新增公告',
          component: () => import('@/views/announcement/AnnouncementForm.vue') },
        { path: 'announcement/form/:id', name: '编辑公告',
          component: () => import('@/views/announcement/AnnouncementForm.vue') },
      ]
    }
  ]
})

// 路由守卫：检查 admin_token
router.beforeEach(to => {
  const token = localStorage.getItem('admin_token')
  if (to.meta.requiresAuth && !token) return '/login'
})
```

**要点：**
- 路由的 `name` 会显示在顶部导航栏
- 所有页面放在 AdminLayout 的 `children` 中
- 使用动态导入实现代码分割

### 7.3 API 调用

**请求封装（`admin/src/api/request.ts`）：**

```typescript
const http = axios.create({ baseURL: '/admin/api', timeout: 30000 })

// 请求拦截器：注入 Token
http.interceptors.request.use(config => {
  const token = localStorage.getItem('admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// 响应拦截器：解包 { code, msg, data }
http.interceptors.response.use(res => {
  const { code, msg, data } = res.data
  if (code !== 0) return Promise.reject(new Error(msg || '请求失败'))
  return data  // 直接返回 data 字段
})
```

**编写 API 模块：**

```typescript
// admin/src/api/announcement.ts
import { get, post, put, del } from './request'

export interface Announcement {
  id: number
  title: string
  content: string
  sort: number
  status: number
  created_at: string
}

export function listAnnouncements(params: { page: number; size: number }) {
  return get<{ list: Announcement[]; total: number }>('/announcements', params)
}

export function createAnnouncement(data: Partial<Announcement>) {
  return post<Announcement>('/announcements', data)
}

export function updateAnnouncement(id: number, data: Partial<Announcement>) {
  return put(`/announcements/${id}`, data)
}

export function deleteAnnouncement(id: number) {
  return del(`/announcements/${id}`)
}
```

### 7.4 列表页开发

以公告列表为例：

```vue
<!-- admin/src/views/announcement/AnnouncementList.vue -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { listAnnouncements, deleteAnnouncement } from '@/api/announcement'
import type { Announcement } from '@/api/announcement'
import { useRouter } from 'vue-router'
import { showToast } from '@/utils/toast'
import { showDialog } from '@/utils/dialog'

const router = useRouter()
const list = ref<Announcement[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const loading = ref(false)

async function fetchList() {
  loading.value = true
  try {
    const data = await listAnnouncements({ page: page.value, size: size.value })
    list.value = data.list || []
    total.value = data.total
  } finally {
    loading.value = false
  }
}

async function handleDelete(id: number) {
  const ok = await showDialog('确认删除该公告？')
  if (!ok) return
  await deleteAnnouncement(id)
  showToast('删除成功')
  fetchList()
}

onMounted(fetchList)
</script>

<template>
  <div class="p-6">
    <!-- 头部 -->
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-bold">公告列表</h2>
      <button
        class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
        @click="router.push('/announcement/form')"
      >
        新增公告
      </button>
    </div>

    <!-- 表格 -->
    <div class="bg-white rounded-lg shadow overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left">
          <tr>
            <th class="px-4 py-3 font-medium">ID</th>
            <th class="px-4 py-3 font-medium">标题</th>
            <th class="px-4 py-3 font-medium">排序</th>
            <th class="px-4 py-3 font-medium">状态</th>
            <th class="px-4 py-3 font-medium">创建时间</th>
            <th class="px-4 py-3 font-medium">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y">
          <tr v-for="item in list" :key="item.id" class="hover:bg-gray-50">
            <td class="px-4 py-3">{{ item.id }}</td>
            <td class="px-4 py-3">{{ item.title }}</td>
            <td class="px-4 py-3">{{ item.sort }}</td>
            <td class="px-4 py-3">
              <span :class="item.status === 1 ? 'text-green-600' : 'text-gray-400'">
                {{ item.status === 1 ? '显示' : '隐藏' }}
              </span>
            </td>
            <td class="px-4 py-3">{{ item.created_at }}</td>
            <td class="px-4 py-3 space-x-2">
              <button
                class="text-blue-600 hover:underline"
                @click="router.push(`/announcement/form/${item.id}`)"
              >
                编辑
              </button>
              <button
                class="text-red-600 hover:underline"
                @click="handleDelete(item.id)"
              >
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 分页 -->
    <div class="flex justify-end mt-4 gap-2">
      <button
        class="px-3 py-1 border rounded"
        :disabled="page <= 1"
        @click="page--; fetchList()"
      >上一页</button>
      <span class="px-3 py-1">{{ page }} / {{ Math.ceil(total / size) }}</span>
      <button
        class="px-3 py-1 border rounded"
        :disabled="page * size >= total"
        @click="page++; fetchList()"
      >下一页</button>
    </div>
  </div>
</template>
```

### 7.5 表单页开发

```vue
<!-- admin/src/views/announcement/AnnouncementForm.vue -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createAnnouncement, updateAnnouncement } from '@/api/announcement'
import { get } from '@/api/request'
import { showToast } from '@/utils/toast'

const route = useRoute()
const router = useRouter()
const id = route.params.id ? Number(route.params.id) : null
const isEdit = !!id

const form = ref({
  title: '',
  content: '',
  sort: 0,
  status: 1,
})
const submitting = ref(false)

onMounted(async () => {
  if (isEdit) {
    const data = await get(`/announcements/${id}`)
    Object.assign(form.value, data)
  }
})

async function handleSubmit() {
  if (!form.value.title.trim()) {
    showToast('请输入标题')
    return
  }
  submitting.value = true
  try {
    if (isEdit) {
      await updateAnnouncement(id!, form.value)
      showToast('更新成功')
    } else {
      await createAnnouncement(form.value)
      showToast('创建成功')
    }
    router.push('/announcement/list')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="p-6 max-w-2xl">
    <h2 class="text-xl font-bold mb-6">{{ isEdit ? '编辑公告' : '新增公告' }}</h2>

    <div class="space-y-4">
      <!-- 标题 -->
      <div>
        <label class="block text-sm font-medium mb-1">标题 <span class="text-red-500">*</span></label>
        <input
          v-model="form.title"
          class="w-full border rounded-lg px-3 py-2"
          placeholder="请输入公告标题"
        />
      </div>

      <!-- 内容 -->
      <div>
        <label class="block text-sm font-medium mb-1">内容</label>
        <textarea
          v-model="form.content"
          class="w-full border rounded-lg px-3 py-2"
          rows="6"
          placeholder="请输入公告内容"
        />
      </div>

      <!-- 排序 -->
      <div>
        <label class="block text-sm font-medium mb-1">排序</label>
        <input
          v-model.number="form.sort"
          type="number"
          class="w-full border rounded-lg px-3 py-2"
        />
      </div>

      <!-- 状态 -->
      <div>
        <label class="block text-sm font-medium mb-1">状态</label>
        <select v-model="form.status" class="w-full border rounded-lg px-3 py-2">
          <option :value="1">显示</option>
          <option :value="0">隐藏</option>
        </select>
      </div>

      <!-- 操作按钮 -->
      <div class="flex gap-3 pt-4">
        <button
          class="px-6 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50"
          :disabled="submitting"
          @click="handleSubmit"
        >
          {{ isEdit ? '保存' : '创建' }}
        </button>
        <button
          class="px-6 py-2 border rounded-lg hover:bg-gray-50"
          @click="router.back()"
        >
          取消
        </button>
      </div>
    </div>
  </div>
</template>
```

### 7.6 认证 Store

管理端认证状态（`admin/src/stores/auth.ts`）：

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const isLoggedIn = computed(() => !!token.value)

  // 从 JWT payload 解析权限
  const permissions = computed<string[]>(() => {
    if (!token.value) return []
    try {
      const payload = token.value.split('.')[1]
      // 标准 base64 解码
      const normalized = payload.replace(/-/g, '+').replace(/_/g, '/')
      const padded = normalized + '='.repeat((4 - normalized.length % 4) % 4)
      const data = JSON.parse(atob(padded))
      return data.perms || []
    } catch { return [] }
  })

  function loginAction(newToken: string) {
    token.value = newToken
    localStorage.setItem('admin_token', newToken)
  }

  function hasPermission(perm: string): boolean {
    return permissions.value.includes('*') || permissions.value.includes(perm)
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('admin_token')
  }

  return { token, isLoggedIn, permissions, loginAction, hasPermission, logout }
})
```

### 7.7 工具函数

**Toast 提示：**

```typescript
// admin/src/utils/toast.ts
import { ref } from 'vue'

export interface ToastItem {
  id: number
  message: string
  type: 'success' | 'error' | 'info'
}

export const toasts = ref<ToastItem[]>([])
let seq = 0

export function showToast(message: string, type: 'success' | 'error' | 'info' = 'success') {
  const id = ++seq
  toasts.value.push({ id, message, type })
  setTimeout(() => {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }, 3000)
}
```

**对话框：**

```typescript
// admin/src/utils/dialog.ts
export function showDialog(message: string): Promise<boolean> {
  return Promise.resolve(confirm(message))
}
```

### 7.8 最佳实践

- **页面按模块分目录**：`views/announcement/`、`views/product/` 等
- **路由名称即页面标题**：`name` 字段会显示在 AdminLayout 的顶部栏
- **使用 TailwindCSS 原子类**：保持样式一致性，无需编写 `<style>` 块
- **API 模块独立封装**：每个业务模块一个 API 文件，导出类型和函数
- **请求错误统一处理**：拦截器自动解包响应并抛出错误

---

## 八、PC 商城与 H5 前端开发

### 8.1 Web 端（PC 商城）

#### 技术架构

```
web/src/
├── api/request.ts     # Axios 封装 + Mock 支持
├── components/        # 通用组件
├── mock/              # Mock 数据
├── router/index.ts    # 路由配置
├── stores/            # Pinia 状态管理
│   ├── auth.ts        # 认证
│   ├── cart.ts        # 购物车
│   ├── chat.ts        # 聊天
│   └── site.ts        # 站点配置
├── utils/             # 工具函数
└── views/             # 页面
```

#### API 请求封装

Web 端的请求封装（`web/src/api/request.ts`）支持 Mock 和真实 API 两种模式：

```typescript
import axios from 'axios'

const MOCK_ENABLED = import.meta.env.VITE_MOCK === 'true'
const BASE_URL = MOCK_ENABLED ? '' : (import.meta.env.VITE_API_URL || '')

const http = axios.create({ baseURL: BASE_URL, timeout: 30000 })

// Token 注入
http.interceptors.request.use(config => {
  const token = localStorage.getItem('user_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// 响应解包
http.interceptors.response.use(
  res => {
    const { code, msg, data } = res.data
    if (code !== 0) return Promise.reject(new Error(msg || '请求失败'))
    return data
  },
  err => Promise.reject(err)
)

// Mock 模式下的请求处理
async function mockRequest<T>(method: string, url: string, params?: any): Promise<T> {
  const { matchMock } = await import('@/mock/index')
  const result = matchMock(method, url, params)
  await new Promise(r => setTimeout(r, 100 + Math.random() * 200))
  if (result.matched) return (result.data ?? null) as T
  console.warn(`[Mock] No data for: ${method} ${url}`)
  return null as T
}

// 导出请求方法
export async function get<T = any>(url: string, params?: any): Promise<T> {
  if (MOCK_ENABLED) return mockRequest<T>('GET', url, params)
  return http.get(url, { params }) as Promise<T>
}

export async function post<T = any>(url: string, data?: any): Promise<T> {
  if (MOCK_ENABLED) return mockRequest<T>('POST', url, data)
  return http.post(url, data) as Promise<T>
}

export async function put<T = any>(url: string, data?: any): Promise<T> {
  if (MOCK_ENABLED) return mockRequest<T>('PUT', url, data)
  return http.put(url, data) as Promise<T>
}

export async function del<T = any>(url: string, data?: any): Promise<T> {
  if (MOCK_ENABLED) return mockRequest<T>('DELETE', url, data)
  return http.delete(url, { data }) as Promise<T>
}

export async function upload<T = any>(url: string, file: File): Promise<T> {
  if (MOCK_ENABLED) return mockRequest<T>('POST', url, { name: file.name, size: file.size })
  const form = new FormData()
  form.append('file', file)
  return http.post(url, form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }) as Promise<T>
}
```

#### Pinia Store 模式

**认证 Store：**

```typescript
// web/src/stores/auth.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('user_token') || '')
  const isLoggedIn = computed(() => !!token.value)

  function setToken(t: string) {
    token.value = t
    localStorage.setItem('user_token', t)
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('user_token')
  }

  return { token, isLoggedIn, setToken, logout }
})
```

**购物车 Store：**

```typescript
// web/src/stores/cart.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

interface CartItem {
  sku_id: number
  qty: number
  product: { id: number; title: string; cover: string; price: number }
  sku: { id: number; attrs: string; price: number; stock: number }
}

export const useCartStore = defineStore('cart', () => {
  const items = ref<CartItem[]>([])
  const total = computed(() => items.value.reduce((s, i) => s + i.sku.price * i.qty, 0))
  const count = computed(() => items.value.length)

  function setItems(list: CartItem[]) { items.value = list }
  function removeItem(skuId: number) {
    items.value = items.value.filter(i => i.sku_id !== skuId)
  }
  function updateQty(skuId: number, qty: number) {
    const item = items.value.find(i => i.sku_id === skuId)
    if (item) item.qty = qty
  }

  return { items, total, count, setItems, removeItem, updateQty }
})
```

#### 主题色系统

站点配置 Store（`web/src/stores/site.ts`）负责动态主题：

```typescript
export interface SiteSettings {
  site_name: string
  site_logo: string
  seo_title: string
  seo_keywords: string
  seo_description: string
  icp: string
  // Hero 区域
  hero_badge: string
  hero_title: string
  hero_subtitle: string
  hero_btn_text: string
  hero_btn_link: string
  // 主题色
  color_primary: string       // 主色，如 #dc2626
  color_primary_light: string // 亮主色
  color_primary_dark: string  // 暗主色
  color_bg_page: string       // 页面背景
  color_bg_header: string     // 头部背景
  color_bg_footer: string     // 底部背景
  color_price: string         // 价格颜色
  color_hero_from: string     // Hero 渐变起始
  color_hero_to: string       // Hero 渐变结束
}
```

**主题应用机制：**

```typescript
function applyTheme(s: SiteSettings) {
  const root = document.documentElement.style
  root.setProperty('--color-primary', s.color_primary)
  root.setProperty('--color-primary-light', s.color_primary_light)
  root.setProperty('--color-primary-dark', s.color_primary_dark)
  root.setProperty('--color-bg-page', s.color_bg_page)
  root.setProperty('--color-bg-header', s.color_bg_header)
  root.setProperty('--color-bg-footer', s.color_bg_footer)
  root.setProperty('--color-price', s.color_price)
  root.setProperty('--color-hero-from', s.color_hero_from)
  root.setProperty('--color-hero-to', s.color_hero_to)

  // 同步 SEO 元信息
  document.title = s.seo_title || s.site_name
}
```

**在组件中使用 CSS 变量：**

```vue
<template>
  <!-- 使用 CSS 变量 -->
  <button :style="{ background: 'var(--color-primary)' }">购买</button>
  <span :style="{ color: 'var(--color-price)' }">¥{{ price }}</span>

  <!-- 在 UnoCSS 中使用 -->
  <div class="bg-[var(--color-bg-page)]">...</div>
</template>
```

#### 路由配置

```typescript
// web/src/router/index.ts
import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  scrollBehavior: () => ({ top: 0 }),
  routes: [
    { path: '/', component: () => import('@/views/Home.vue') },
    { path: '/products', component: () => import('@/views/ProductList.vue') },
    { path: '/product/:id', component: () => import('@/views/ProductDetail.vue') },
    { path: '/cart', component: () => import('@/views/Cart.vue') },
    { path: '/orders', component: () => import('@/views/OrderList.vue') },
    { path: '/orders/:id', component: () => import('@/views/OrderDetail.vue') },
    { path: '/login', component: () => import('@/views/Login.vue') },
    { path: '/user', component: () => import('@/views/UserCenter.vue') },
    // 添加新页面...
  ]
})
```

#### 新增页面示例

在 Web 端添加一个公告展示页面：

```vue
<!-- web/src/views/Announcements.vue -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/api/request'

interface Announcement {
  id: number
  title: string
  content: string
  created_at: string
}

const list = ref<Announcement[]>([])

onMounted(async () => {
  list.value = await get<Announcement[]>('/api/v1/announcements') || []
})
</script>

<template>
  <div class="max-w-4xl mx-auto py-8 px-4">
    <h1 class="text-2xl font-bold mb-6">公告通知</h1>
    <div class="space-y-4">
      <div
        v-for="item in list"
        :key="item.id"
        class="bg-white rounded-lg shadow p-6"
      >
        <h2 class="text-lg font-semibold mb-2">{{ item.title }}</h2>
        <p class="text-gray-600 text-sm mb-3">{{ item.created_at }}</p>
        <div class="text-gray-700 leading-relaxed">{{ item.content }}</div>
      </div>
      <div v-if="!list.length" class="text-center text-gray-400 py-12">
        暂无公告
      </div>
    </div>
  </div>
</template>
```

然后在 `web/src/router/index.ts` 添加路由：

```typescript
{ path: '/announcements', component: () => import('@/views/Announcements.vue') },
```

### 8.2 App 端（H5/小程序）

#### 技术架构

App 端使用 uni-app 框架，同时支持 H5 和微信小程序：

```
app/
├── composables/       # 组合式函数
│   └── useTheme.ts    # 主题管理
├── components/        # 组件
├── mock/              # Mock 数据 + 行业预设
├── pages/             # 页面（按 pages.json 配置）
├── utils/
│   └── request.ts     # 请求封装（uni.request）
├── pages.json         # 路由配置
├── App.vue            # 根组件
└── main.ts            # 入口
```

#### 路由配置（pages.json）

uni-app 使用 `pages.json` 代替 Vue Router：

```json
{
  "pages": [
    { "path": "pages/index/index", "style": { "navigationBarTitleText": "首页", "navigationStyle": "custom" } },
    { "path": "pages/product/list", "style": { "navigationBarTitleText": "商品" } },
    { "path": "pages/product/detail", "style": { "navigationBarTitleText": "商品详情" } },
    { "path": "pages/cart/index", "style": { "navigationBarTitleText": "购物车" } },
    { "path": "pages/order/list", "style": { "navigationBarTitleText": "订单" } }
  ],
  "tabBar": {
    "color": "#999",
    "selectedColor": "#dc2626",
    "list": [
      { "pagePath": "pages/index/index", "text": "首页" },
      { "pagePath": "pages/product/list", "text": "商品" },
      { "pagePath": "pages/cart/index", "text": "购物车" },
      { "pagePath": "pages/order/list", "text": "订单" },
      { "pagePath": "pages/user/index", "text": "我的" }
    ]
  }
}
```

**添加新页面：**

```json
// 在 pages 数组中添加
{ "path": "pages/announcement/index", "style": { "navigationBarTitleText": "公告" } }
```

然后创建对应文件 `app/pages/announcement/index.vue`。

#### 请求封装

App 端使用 `uni.request` 封装（`app/utils/request.ts`）：

```typescript
const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export function get<T>(url: string, params?: any): Promise<T> {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('user_token')
    uni.request({
      url: `${BASE_URL}${url}`,
      method: 'GET',
      data: params,
      header: { Authorization: token ? `Bearer ${token}` : '' },
      success(res: any) {
        const { code, msg, data } = res.data
        if (code !== 0) return reject(new Error(msg))
        resolve(data as T)
      },
      fail: reject,
    })
  })
}
```

#### 暗黑模式 Composable

主题管理（`app/composables/useTheme.ts`）：

```typescript
import { ref, computed } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'auto'

const STORAGE_KEY = 'app_theme_mode'
const themeMode = ref<ThemeMode>(uni.getStorageSync(STORAGE_KEY) || 'auto')

function getSystemTheme(): 'light' | 'dark' {
  try {
    const info = uni.getSystemInfoSync()
    if ((info as any).theme === 'dark') return 'dark'
  } catch {}
  if (typeof window !== 'undefined' && window.matchMedia?.('(prefers-color-scheme: dark)').matches) {
    return 'dark'
  }
  return 'light'
}

const effectiveTheme = computed(() => {
  if (themeMode.value === 'auto') return getSystemTheme()
  return themeMode.value
})

function applyTheme() {
  const theme = effectiveTheme.value
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('data-up-theme', theme)
  }
}

function setTheme(mode: ThemeMode) {
  themeMode.value = mode
  uni.setStorageSync(STORAGE_KEY, mode)
  applyTheme()
}

function toggleTheme() {
  setTheme(effectiveTheme.value === 'light' ? 'dark' : 'light')
}

function initTheme() {
  applyTheme()
  if (typeof window !== 'undefined') {
    window.matchMedia?.('(prefers-color-scheme: dark)')
      .addEventListener('change', () => { if (themeMode.value === 'auto') applyTheme() })
  }
}

export function useTheme() {
  return { themeMode, effectiveTheme, setTheme, toggleTheme, applyTheme, initTheme }
}
```

**在页面中使用：**

```vue
<script setup lang="ts">
import { useTheme } from '@/composables/useTheme'
const { effectiveTheme, toggleTheme } = useTheme()
</script>

<template>
  <view @tap="toggleTheme">
    当前主题：{{ effectiveTheme }}
  </view>
</template>
```

### 8.3 跨端开发注意事项

| 差异点 | Web 端 | App 端 |
|--------|--------|--------|
| HTTP 客户端 | Axios | uni.request |
| 路由方式 | Vue Router | pages.json + uni.navigateTo |
| 存储 | localStorage | uni.getStorageSync |
| 样式框架 | UnoCSS | UnoCSS + uview-plus |
| 条件编译 | 不支持 | 支持 `#ifdef H5` / `#ifdef MP-WEIXIN` |
| 页面导航 | `router.push()` | `uni.navigateTo({ url: '/pages/xxx' })` |
| 文件上传 | `FormData` | `uni.uploadFile` |

---

## 九、配置中心与权限系统

### 9.1 配置中心

#### 配置存储模型

配置数据存储在 `configs` 表中，以插件为命名空间：

```go
type ConfigKV struct {
    ID     uint64 `gorm:"primaryKey;autoIncrement"`
    Plugin string `gorm:"size:64;not null;uniqueIndex:uk_plugin_key"`  // 插件名
    Key    string `gorm:"size:128;not null;uniqueIndex:uk_plugin_key"` // 配置键
    Value  string `gorm:"type:text"`                                   // 配置值
}
```

#### 配置声明（plugin.json）

插件通过 `config_items` 声明配置字段，管理端自动渲染配置表单：

```json
{
  "name": "sms",
  "config_items": [
    {
      "key": "provider",
      "label": "短信服务商",
      "type": "select",
      "required": true,
      "options": [
        { "label": "阿里云", "value": "aliyun" },
        { "label": "腾讯云", "value": "tencent" }
      ]
    },
    {
      "key": "access_key",
      "label": "AccessKey",
      "type": "text",
      "required": true
    },
    {
      "key": "secret_key",
      "label": "SecretKey",
      "type": "password",
      "required": true
    },
    {
      "key": "sign_name",
      "label": "签名名称",
      "type": "text",
      "required": true,
      "placeholder": "如：LYShop"
    }
  ]
}
```

**支持的字段类型：**

| type | 渲染方式 | 适用场景 |
|------|----------|----------|
| `text` | 单行文本输入 | AccessKey、域名 |
| `password` | 密码输入框 | SecretKey |
| `textarea` | 多行文本输入 | 证书内容 |
| `number` | 数字输入 | 端口号、限额 |
| `select` | 下拉选择 | 服务商选择 |
| `switch` | 开关 | 功能启用/禁用 |

#### 配置中心 API

```
GET  /admin/api/config/schemas     → 获取所有插件的配置声明
GET  /admin/api/config/:plugin     → 获取某插件的配置值
PUT  /admin/api/config/:plugin     → 保存某插件的配置值
```

**示例请求：**

```bash
# 获取短信插件配置
GET /admin/api/config/sms
# 响应: {"code":0, "data": {"provider":"aliyun", "access_key":"xxx", ...}}

# 保存配置
PUT /admin/api/config/sms
Body: {"provider":"aliyun", "access_key":"LTAI...", "secret_key":"xxx", "sign_name":"LYShop"}
```

#### 在插件中读取配置

```go
func readConfig(pluginName string) map[string]string {
    var kvs []model.ConfigKV
    db.DB.Where("plugin = ?", pluginName).Find(&kvs)
    result := make(map[string]string)
    for _, kv := range kvs {
        result[kv.Key] = kv.Value
    }
    return result
}

// 使用
cfg := readConfig("sms")
provider := cfg["provider"]  // "aliyun"
accessKey := cfg["access_key"]
```

### 9.2 权限系统（RBAC）

#### 架构设计

LYShop 使用基于角色的访问控制（RBAC）：

```
Admin —— belongsTo ——→ Role —— has ——→ []Permission
```

- 每个 Admin 关联一个 Role
- 每个 Role 包含一组权限字符串（`json.RawMessage`）
- 权限格式：`模块:操作`，如 `product:view`、`order:edit`
- `"*"` 通配符表示超级管理员（拥有所有权限）

#### 权限声明

每个插件在 `plugin.json` 中声明自己的权限：

```json
{
  "permissions": [
    "product:view",
    "product:create",
    "product:edit",
    "product:delete"
  ]
}
```

框架自动聚合所有启用插件的权限：

```go
// 获取所有已声明的权限（用于角色管理界面）
allPerms := plugin.AllPermissions(config.Global.Plugins.Enabled)
```

#### 权限校验

**后端校验（中间件）：**

```go
// 单个接口声明所需权限
g.GET("/products", middleware.RequirePermission("product:view"), handler)

// RequirePermission 的实现
func RequirePermission(perm string) gin.HandlerFunc {
    return func(c *gin.Context) {
        perms := c.GetStringSlice("perms") // 从 JWT 解析
        for _, p := range perms {
            if p == "*" || p == perm {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(200, response.Err(403, "无权限: "+perm))
    }
}
```

**前端校验（Store）：**

```typescript
const authStore = useAuthStore()

// 在页面或组件中检查权限
if (authStore.hasPermission('product:create')) {
  // 显示「新增商品」按钮
}
```

### 9.3 菜单系统

管理端菜单由插件动态声明，框架按权限过滤：

```go
// 管理端菜单 API（app.go 中注册）
adminAuth.GET("/menus", func(c *gin.Context) {
    perms, _ := c.Get("perms")
    permList, _ := perms.([]string)
    // 只返回当前管理员有权限的菜单项
    menus := plugin.EnabledMenus(config.Global.Plugins.Enabled, permList)
    c.JSON(200, menus)
})
```

**菜单过滤逻辑：**
1. 遍历 `plugins.enabled` 列表中的每个插件
2. 获取插件的 `Menus` 声明
3. 递归过滤：如果菜单项有 `permission` 字段且用户无此权限，则隐藏
4. 如果父菜单的所有子菜单都被过滤，则父菜单也隐藏

### 9.4 完整的权限配置流程

1. **插件声明权限**：在 `plugin.json` 的 `permissions` 字段中列出
2. **创建角色**：在管理端创建角色，勾选所需权限
3. **分配角色**：将角色分配给管理员账户
4. **后端校验**：每个 API 通过 `RequirePermission` 中间件校验
5. **前端校验**：菜单 API 按权限过滤，页面内按需检查

### 9.5 最佳实践

- **权限命名规范**：`模块:操作`，如 `product:view`、`order:create`
- **配置键命名**：使用 snake_case，与 Go 和 JSON 风格一致
- **敏感配置用 password 类型**：在前端以掩码显示
- **`"*"` 权限仅给超级管理员**：不要在普通角色中使用通配符
- **配置变更后需重新初始化驱动**：部分驱动在 Install 时读取配置，修改后可能需要重启

---

## 十、测试、调试与部署

### 10.1 Mock 系统

所有前端应用都支持 Mock 模式，通过 `VITE_MOCK=true` 环境变量启用。

**启动命令：**

```bash
# Web 端演示模式
cd web && npm run dev:demo

# Admin 端演示模式
cd admin && npm run dev:demo

# App H5 演示模式
cd app && npm run dev:h5:demo
```

**Mock 数据组织：**

```
web/src/mock/index.ts      # Web 端 Mock 路由匹配
admin/src/mock/index.ts    # Admin 端 Mock
app/mock/
├── index.ts               # App 端 Mock 路由匹配
├── data/                  # 原始数据（JSON）
└── presets/               # 行业预设（商城/超市/生鲜/珠宝 等 7 个）
    ├── types.ts           # MockPreset 类型定义
    ├── index.ts           # 预设加载器
    ├── mall.ts            # 综合商城预设
    └── ...
```

**App 商品详情演示补全：**

- `app/mock/index.ts` 继续复用现有 `/api/v1/products/:id` 和 `/api/v1/products/:id/reviews`，不新增演示专用接口。
- 当预设仅提供商品列表基础字段时，Mock 层会自动补齐商品详情所需的轮播图、SKU 规格、详情图文 blocks、销量/库存/收藏数。
- 对没有真实订单评价沉淀的商品，Mock 层会补齐评价列表、追评和商家回复，确保列表页、推荐位和营销页跳转到详情页时都有完整演示数据。
- 商品详情轮播图、详情图和评价晒图会优先复用商品自身 `cover`、`images` 和 `detail.blocks` 中已有图片，减少随机图与商品语义不匹配的问题。

**Mock 路由匹配原理：**

```typescript
// mock/index.ts 中的 matchMock 函数
export function matchMock(method: string, url: string, params?: any) {
  // 路由表以 "METHOD /path" 为 key
  const routes: Record<string, Function> = {
    'GET /api/v1/products': (params) => {
      // 返回 Mock 数据
      return { matched: true, data: { list: [...], total: 100 } }
    },
    'GET /api/v1/products/:id': (params) => {
      const id = extractParam(url, '/api/v1/products/:id')
      return { matched: true, data: productDetail }
    },
    // ...
  }
}
```

**添加新接口的 Mock 数据：**

在对应的 `mock/index.ts` 中添加路由匹配规则即可：

```typescript
// 为公告接口添加 Mock
'GET /api/v1/announcements': () => ({
  matched: true,
  data: [
    { id: 1, title: '618 大促开启', content: '...', created_at: '2026-06-01' },
    { id: 2, title: '新用户注册送礼', content: '...', created_at: '2026-05-20' },
  ]
}),
```



### 10.1.1 PC 装修高级样式结构（PcDecorPage）

PC 装修接口路径为（`/admin/api/decor/pc`、`/api/v1/pc/decor`），`components` 字段使用页面对象：

```json
{
  "components": {
    "pageStyle": {
      "background": {
        "mode": "solid",
        "solidColor": "#f8fafc",
        "gradient": {
          "angle": 135,
          "stops": [
            { "color": "#f8fafc", "position": 0 },
            { "color": "#eef2ff", "position": 100 }
          ]
        },
        "image": {
          "url": "",
          "size": "cover",
          "customSize": "100% auto",
          "position": "center top",
          "repeat": "no-repeat",
          "attachment": "scroll"
        },
        "overlay": { "enabled": false, "color": "#000000", "opacity": 0.2 }
      },
      "content": { "maxWidth": 1280, "gutterX": 24, "sectionGap": 24 },
      "surface": { "radius": 12, "shadow": "none" }
    },
    "components": [
      {
        "id": "pc_hero",
        "type": "hero",
        "props": {},
        "style": {
          "marginTop": 24,
          "marginBottom": 0,
          "paddingX": 24,
          "paddingY": 0,
          "backgroundColor": "transparent",
          "borderRadius": 12,
          "borderWidth": 0,
          "borderColor": "#e5e7eb",
          "shadow": "none"
        }
      }
    ]
  }
}
```

说明：

- `pageStyle` 管页面级背景、遮罩、内容宽度和默认表面样式。
- `components[].style` 是组件级覆盖，优先级高于 `pageStyle` 的默认值。
- 背景图支持 URL + 上传配置，字段由 Admin 装修页直接产出。

### 10.1.2 Mock 同步规则（admin/web/app）

PC 装修字段要求三端 Mock 同步：

- `admin/src/mock/index.ts`：`GET/PUT /admin/api/decor/pc` 使用 `PcDecorPage` 结构。
- `web/src/mock/index.ts`：`GET /api/v1/pc/decor` 返回同结构对象。
- `app/mock/presets/types.ts` 与各行业预设：`pcDecor` 使用 `{ pageStyle, components }`。

若仅改其中一端，会出现「后台预览正常、前台演示异常」或类型检查失败。

### 10.2 本地调试

#### 后端调试

```bash
# 开发模式（debug 级日志）
# config.yaml 中设置 server.mode: debug
go run main.go -config config.yaml

# 只编译不运行（检查编译错误）
go build ./...

# 运行测试
go test ./...
```

**代理配置：**

前端开发服务器通过 Vite 代理转发 API 请求到后端：

```typescript
// web/vite.config.ts
export default defineConfig({
  server: {
    port: 9529,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})

// admin/vite.config.ts
export default defineConfig({
  server: {
    port: 9527,
    proxy: {
      '/admin/api': { target: 'http://localhost:8080', changeOrigin: true },
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
    },
  },
})
```

#### 前端调试

```bash
# 真实 API 模式（需要后端运行）
npm run dev

# Mock 模式（无需后端）
npm run dev:demo
```

#### 联调流程

1. 先用 Mock 模式开发前端页面，确保 UI 正确
2. 后端实现 API 后，切换到真实模式验证数据交互
3. 使用浏览器开发者工具检查请求和响应
4. 对核心流程（登录、下单、支付）做端到端回归

### 10.3 Go 测试

```go
// server/plugins/announcement/service/announcement_test.go
package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    annosvc "github.com/ijry/lyshop/plugins/announcement/service"
    annomodel "github.com/ijry/lyshop/plugins/announcement/model"
)

func TestCreateAndList(t *testing.T) {
    ctx := context.Background()

    // 创建公告
    anno := &annomodel.Announcement{
        Title:   "测试公告",
        Content: "这是一条测试公告",
        Status:  1,
    }
    err := annosvc.Create(ctx, anno)
    assert.NoError(t, err)
    assert.NotZero(t, anno.ID)

    // 查询已发布公告
    list, err := annosvc.ListPublished(ctx)
    assert.NoError(t, err)
    assert.NotEmpty(t, list)
}
```

### 10.4 CI/CD

项目使用 GitHub Actions 进行 CI（`.github/workflows/ci.yml`）：

```yaml
name: CI
on:
  push:
    branches: [master]
  pull_request:
    branches: [master]
jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: go test ./...
      - run: go build ./...

  admin:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '20' }
      - run: cd admin && npm install && npm run build
```

### 10.5 生产构建

```bash
# 后端构建
cd server && go build -o lyshop main.go

# Admin 构建（输出到 admin/dist/）
cd admin && npm run build

# PC 商城构建（输出到 web/dist/）
cd web && npm run build

# H5 构建（输出到 app/dist/build/h5/）
cd app && npm run build:h5

# 微信小程序构建（输出到 app/dist/build/mp-weixin/）
cd app && npm run build:mp-weixin
```

### 10.6 Docker 部署

项目提供完整的 Docker Compose 配置：

```bash
# 一键部署
docker-compose up -d
```

**docker-compose.yml 架构：**

```
┌─────────┐     ┌──────────────┐     ┌──────────┐
│  Nginx  │────▶│  Go Server   │────▶│  MySQL   │
│ (反向代理) │     │  (:8080)     │     │          │
└────┬────┘     └──────┬───────┘     └──────────┘
     │                 │
     │                 └────▶┌──────────┐
     │                       │  Redis   │
     ├──▶ /admin → admin/dist│          │
     ├──▶ /h5   → app/dist  └──────────┘
     └──▶ /     → web/dist
```

**环境配置清单：**

| 配置项 | 说明 | 示例 |
|--------|------|------|
| `database.dsn` | MySQL 连接串 | `root:pass@tcp(mysql:3306)/lyshop?...` |
| `redis.addr` | Redis 地址 | `redis:6379` |
| `jwt.secret` | JWT 密钥（**务必修改**） | 随机字符串 |
| `server.mode` | 运行模式 | `release` |
| `plugins.enabled` | 启用的插件列表 | 见 config.example.yaml |

### 10.7 Nginx 配置示例

```nginx
server {
    listen 80;
    server_name shop.example.com;

    # PC 商城
    location / {
        root /usr/share/nginx/html/web;
        try_files $uri $uri/ /index.html;
    }

    # Admin 后台
    location /admin {
        alias /usr/share/nginx/html/admin;
        try_files $uri $uri/ /admin/index.html;
    }

    # H5 移动端
    location /h5 {
        alias /usr/share/nginx/html/h5;
        try_files $uri $uri/ /h5/index.html;
    }

    # API 反向代理
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /admin/api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # WebSocket
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

### 10.8 最佳实践

- **开发时先用 Mock 模式**：快速迭代 UI，避免依赖后端
- **为新接口同步添加 Mock 数据**：保证演示模式可用
- **生产环境必须修改 JWT Secret**：使用随机字符串
- **使用 `server.mode: release`**：生产环境关闭 Gin 的 debug 日志
- **WebSocket 需要 Nginx 特殊配置**：升级 HTTP 协议为 WebSocket
- **数据库备份**：生产环境定期备份 MySQL 数据

---

## 附录：常用开发速查

### API 响应格式

```json
// 成功
{"code": 0, "msg": "success", "data": {...}}

// 分页
{"code": 0, "msg": "success", "data": {"list": [...], "total": 100, "page": 1, "size": 20}}

// 失败
{"code": 10001, "msg": "商品不存在", "data": null}
```

### 新增插件 Checklist

- [ ] 创建 `server/plugins/<name>/` 目录
- [ ] 编写 `plugin.json`（name, menus, permissions, config_items）
- [ ] 定义 Model（嵌入 `model.Base`）
- [ ] 编写 Service（业务逻辑）
- [ ] 编写 API（front.go + admin.go）
- [ ] 编写 `plugin.go`（init 注册 + Migrate + RegisterRoutes）
- [ ] 在 `main.go` 添加空导入
- [ ] 在 `config.yaml` 的 `plugins.enabled` 中启用
- [ ] 在 `admin/src/router/index.ts` 添加管理端路由
- [ ] 在 `admin/src/views/<name>/` 创建管理端页面
- [ ] 为新接口添加 Mock 数据
- [ ] 更新文档

### 文件命名约定

| 位置 | 命名规则 | 示例 |
|------|----------|------|
| Go 包 | snake_case | `announcement`, `wechat_pay` |
| Go 文件 | snake_case | `announcement.go` |
| Vue 组件 | PascalCase | `AnnouncementList.vue` |
| TypeScript | camelCase | `request.ts`, `orderStatus.ts` |
| CSS 类 | kebab-case / TailwindCSS | `text-red-600` |
| 数据库表 | snake_case (复数) | `announcements` |
| JSON 字段 | snake_case | `created_at`, `category_id` |
