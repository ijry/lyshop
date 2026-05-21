# 部署文档

本文档描述 LYShop 的常见部署方式。

## 环境要求

- Node.js 18+
- Go 1.21+
- MySQL 8+
- Redis（可选，按插件能力使用）

## 本地开发部署

### 后端服务

```bash
cd server
go mod tidy
go run main.go
```

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

### 文档站

```bash
cd docs-site
npm install
npm run docs:dev
```

## 生产构建

### 管理后台

```bash
cd admin
npm run build
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
