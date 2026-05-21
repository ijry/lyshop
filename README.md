# LYShop

LYShop 是一个面向电商场景的全栈项目，包含后端服务、管理后台、移动端应用与插件化业务能力。

## 项目简介

- `server`：Go 后端服务，提供认证、商品、订单、营销、IM、支付等能力
- `admin`：基于 Vue 3 + Vite 的管理后台
- `app`：基于 uni-app 的移动端应用
- `docs-site`：基于 VitePress 的项目官网与文档站

## 技术栈

- 后端：Go、GORM、插件化模块机制
- 管理端：Vue 3、TypeScript、Vite、Pinia
- 移动端：uni-app、TypeScript、Pinia
- 文档站：VitePress

## 目录结构

```text
lyshop/
├─ server/              # 后端服务
├─ admin/               # 管理后台
├─ app/                 # 移动端应用
├─ docs/                # 研发规划与设计文档
└─ docs-site/           # 官网文档站（VitePress）
```

## 快速开始

### 1. 启动后端

```bash
cd server
go mod tidy
go run main.go
```

### 2. 启动管理后台

```bash
cd admin
npm install
npm run dev
```

### 3. 启动移动端（H5）

```bash
cd app
npm install
npm run dev:h5
```

### 4. 启动文档站

```bash
cd docs-site
npm install
npm run docs:dev
```

## Docker 部署

项目根目录已提供 `docker-compose.yml`，可按实际环境补充配置后执行：

```bash
docker compose up -d --build
```

## 文档入口

- 官网首页：`docs-site/docs/index.md`
- 功能介绍：`docs-site/docs/guide/features.md`
- 部署文档：`docs-site/docs/deploy/index.md`
- 接口文档：`docs-site/docs/api/index.md`
- 二次开发：`docs-site/docs/dev/secondary-development.md`
- 在线文档：`https://ijry.github.io/lyshop/`

## 后续建议

- 接口章节可后续接入 OpenAPI 自动生成
- 发布时可将 `docs-site/docs/.vitepress/dist` 托管到 Nginx 或对象存储
