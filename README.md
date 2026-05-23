# LYShop 零云商城

开源插件化多端商城系统，基于 Go + Vue3 + uni-app。

[![CI](https://github.com/ijry/lyshop/actions/workflows/ci.yml/badge.svg)](https://github.com/ijry/lyshop/actions/workflows/ci.yml)
[![Docs](https://github.com/ijry/lyshop/actions/workflows/docs.yml/badge.svg)](https://ijry.github.io/lyshop/)
[![License: MIT](https://img.shields.io/badge/License-MIT-red.svg)](LICENSE)

## 在线演示

- [PC 商城演示](https://ijry.github.io/lyshop/web-demo/)
- [管理后台演示](https://ijry.github.io/lyshop/admin-demo/)（账号 admin / admin123）
- [H5 移动端演示](https://ijry.github.io/lyshop/)（右下角浮窗）
- [项目文档](https://ijry.github.io/lyshop/)

## 特性

- **完全插件化** — 商品、订单、营销、IM、支付等均为独立插件，config.yaml 一行开关
- **多端覆盖** — PC Web + H5 + 微信小程序 + App（iOS/Android）+ 管理后台
- **驱动抽象层** — 支付、短信、OAuth、存储、AI 统一接口，一行代码切换服务商
- **AI 生图** — 通义万象 / 文心 / DALL-E 多模型聚合生成商品图
- **IM 客服** — WebSocket 实时通信，多坐席，断线重连，声音提醒
- **营销引擎** — 价格计算管线，秒杀/拼团/砍价/优惠券（可叠加）/积分/分销（2级返利）
- **RBAC 权限** — 角色 + 细粒度权限，菜单按权限动态过滤
- **一键部署** — Docker Compose，MySQL + Redis + Nginx 全容器化

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go + Gin + GORM + MySQL + Redis |
| 管理后台 | Vue3 + Vite + TailwindCSS + shadcn-vue |
| PC 商城 | Vue3 + Vite + UnoCSS |
| 移动端 | uni-app + uview-plus 3.x + UnoCSS |
| 文档站 | VitePress |
| 部署 | Docker Compose + Nginx |

## 目录结构

```
lyshop/
├── server/        # Go 后端（插件化架构）
├── admin/         # Vue3 管理后台
├── web/           # Vue3 PC 商城前端
├── app/           # uni-app 移动端（H5/小程序/App）
├── docs-site/     # VitePress 文档站
└── docker-compose.yml
```

## 插件列表

| 插件 | 类型 | 说明 |
|------|------|------|
| `product` | 功能 | 商品管理、多规格 SKU、分类 |
| `order` | 功能 | 购物车、订单、支付、售后 |
| `marketing` | 功能 | 优惠券、秒杀、拼团、砍价、满减、积分、分销 |
| `im` | 功能 | WebSocket IM 客服 |
| `wms` | 功能 | 仓储管理 |
| `ai_image` | 功能 | AI 生成商品图 |
| `decor` | 功能 | 店铺装修 |
| `checkin` | 功能 | 每日签到 |
| `message` | 功能 | 消息中心 |
| `wechat_pay` | 驱动 | 微信支付 |
| `alipay` | 驱动 | 支付宝支付 |
| `sms` | 驱动 | 短信（阿里云/腾讯云） |
| `wechat_auth` | 驱动 | 微信登录 |
| `storage_local` | 驱动 | 本地存储 |
| `storage_oss` | 驱动 | 阿里云 OSS |
| `storage_cos` | 驱动 | 腾讯云 COS |
| `storage_qiniu` | 驱动 | 七牛云 |

## 快速开始

### 1. 后端

```bash
cd server
cp ../config.example.yaml ../config.yaml  # 编辑配置
go mod tidy
go run main.go
# 首次启动自动建表 + 创建超管 admin/admin123
```

### 2. 管理后台

```bash
cd admin
npm install
npm run dev        # http://localhost:9527
npm run dev:demo   # mock 演示模式
```

### 3. PC 商城

```bash
cd web
npm install
npm run dev        # http://localhost:9529
npm run dev:demo   # mock 演示模式
```

### 4. 移动端

```bash
cd app
npm install --legacy-peer-deps
npm run dev:h5       # H5 开发
npm run dev:h5:demo  # mock 演示模式
```

### 5. Docker 一键部署

```bash
cp config.example.yaml config.yaml
docker compose up -d
```

## 文档

- [功能介绍](https://ijry.github.io/lyshop/guide/features)
- [部署文档](https://ijry.github.io/lyshop/deploy/)
- [接口文档](https://ijry.github.io/lyshop/api/)
- [二次开发](https://ijry.github.io/lyshop/dev/secondary-development)

## License

[MIT](LICENSE)
