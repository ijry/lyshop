# 部署文档

本文档描述 LYShop 的常见部署方式。

## 环境要求

- Node.js 18+
- Go 1.21+
- SQLite（默认）或 MySQL 8+
- Redis（可选，按插件能力使用）

## 本地开发部署

### 后端服务

```bash
cd server
cp ../config.example.yaml ../config.yaml
go mod tidy
go run main.go -config ../config.yaml
```

`config.example.yaml` 默认使用 SQLite：

```yaml
database:
  dsn: "lyshop.db"
```

如需改为 MySQL，只需将 `database.dsn` 替换为 MySQL 连接串。

`redis.addr` 留空时，后端会自动启用进程内嵌 Redis，适合个人单机运行。

### 管理后台

```bash
cd admin
npm install
npm run dev
```

### 移动端 H5

```bash
cd app
npm install
npm run dev:h5
```

### 商家移动端 eapp

```bash
cd eapp
npm install --legacy-peer-deps
npm run dev:h5
```

### 文档站

```bash
cd docs-site
npm install
npm run docs:dev
```

## 单机内嵌部署（个人用户推荐）

该模式会把 `web/admin/app(H5)` 构建产物嵌入后端二进制，并提供 Windows GUI 启动器。

```bash
pwsh -File scripts/build-standalone.ps1
```

产物目录：`dist/`

- `lyshop.exe`：命令行后端
- `lyshop-gui.exe`：Windows 无控制台 GUI 启动器

内嵌模式路由：

- 商城：`/`
- 管理后台：`/admin/`（接口仍为 `/admin/api/*`）
- H5：`/h5/`

## 生产构建

### 管理后台

```bash
cd admin
npm run build
```

管理后台构建同时执行 TypeScript 检查和 Vite 打包，产物位于 `admin/dist`。Element Plus 的组件、中文语言包及样式随前端产物交付，无需新增环境变量或后端配置。替换部署时应同时更新入口 HTML 与带哈希的静态资源；使用内嵌部署时需重新构建后端可执行文件。

分销管理接口统一使用 `/admin/api/distribution/*`，路由前缀由后端应用和前端请求封装各自统一配置；分销用户端为 `/api/v1/distribution/*`。部署需使用对应的前后端构建产物，无需数据迁移。

演示管理后台使用 `npm run build:demo` 生成带 Mock 数据的产物，分销筛选、状态更新、订单、提现和配置读写均在浏览器内演示，不连接真实业务接口。验证可用 `npm run dev:demo`，并检查分销、拼团、砍价表格及首页日期选择器。

### 商家移动端 eapp（H5）

```bash
cd eapp
npm run build:h5
```

若需要构建微信小程序：

```bash
cd eapp
npm run build:mp-weixin
```

### 文档站

```bash
cd docs-site
npm run docs:build
```

生成目录为 `docs-site/docs/.vitepress/dist`，可作为静态站点发布。

若部署到 GitHub Pages 项目页（例如 `https://ijry.github.io/lyshop/`），需在
`docs-site/docs/.vitepress/config.mts` 中配置：

```ts
base: "/lyshop/"
```

## Nginx 托管示例

```nginx
server {
  listen 80;
  server_name docs.example.com;

  root /var/www/lyshop-docs;
  index index.html;

  location / {
    try_files $uri $uri/ /index.html;
  }
}
```

## Docker Compose

根目录 `docker-compose.yml` 可用于整体容器化部署。建议在生产环境补充：

- 环境变量与密钥管理
- 数据卷持久化
- 反向代理与 HTTPS
- 若使用 MySQL 容器，请在 `config.yaml` 中将 `database.dsn` 设置为 MySQL 连接串（例如 `root:password@tcp(mysql:3306)/lyshop?...`）。

## eapp 部署影响

- 后端无新增配置项，继续复用现有 `/admin/api/*`。
- Nginx 需新增 eapp H5 静态目录映射（如 `/eapp/`）。
- 微信小程序与 App 仍按 uni-app 常规流程在对应平台发布，不依赖额外后端开关。
