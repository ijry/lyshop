---
layout: home

hero:
  name: "LYShop 零云商城"
  text: "开源插件化多端商城系统"
  tagline: "基于 Go + Vue3 + uni-app，支持 PC / H5 / 小程序 / App 四端，插件化架构，AI生图，IM客服，Docker一键部署"
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/features
    - theme: alt
      text: 在线演示
      link: /web-demo/index.html

features:
  - title: 全栈多端
    details: Go 后端 + Vue3 管理后台 + uni-app 移动端 + PC Web 商城，一套代码覆盖全平台。
  - title: 完全插件化
    details: 商品、订单、营销、IM、支付、存储等均为独立插件，按需启用，零耦合扩展。
  - title: AI 生图
    details: 集成通义万象、文心、DALL-E 等多模型，一键生成商品轮播图和详情图。
  - title: IM 客服
    details: WebSocket 实时通信，多坐席会话分配，离线消息，自动回复。
  - title: 驱动抽象层
    details: 支付、短信、OAuth、存储、AI 统一驱动接口，一行代码切换服务商。
  - title: 一键部署
    details: Docker Compose 一键启动，MySQL + Redis + Nginx 全容器化，开箱即用。
---

## 为什么选择 LYShop？

LYShop 是一套面向中小商户的**开源免费**商城系统，采用现代技术栈和插件化架构设计，无需授权费用，支持私有化部署，适合独立开发者和团队二次开发。

### 核心优势

- **真正的插件化** — 每个功能模块是独立插件，包含自己的模型、接口、菜单、迁移脚本，通过配置文件一行开关
- **多端覆盖** — PC Web + H5 + 微信小程序 + iOS/Android App，共享同一套后端 API
- **AI 赋能电商** — 多模型聚合生成商品图片，降低运营成本
- **开发者友好** — 清晰的分层架构，完善的文档，TypeScript 全覆盖

## 快速导航

- [功能介绍](/guide/features) — 完整功能清单
- [部署文档](/deploy/) — 本地开发 & Docker 部署
- [接口文档](/api/) — RESTful API 参考
- [二次开发](/dev/secondary-development) — 插件开发指南
- [PC 商城演示](/web-demo/index.html) — 在线体验

## 技术架构

```text
                    ┌─────────────┐
                    │   Nginx     │
                    └──────┬──────┘
           ┌───────────────┼───────────────┐
           ▼               ▼               ▼
    ┌──────────┐    ┌──────────┐    ┌──────────┐
    │ PC Web   │    │ H5/App   │    │  Admin   │
    │ Vue3     │    │ uni-app  │    │  Vue3    │
    └────┬─────┘    └────┬─────┘    └────┬─────┘
         └───────────────┼───────────────┘
                         ▼
              ┌─────────────────────┐
              │    Go Server (Gin)  │
              │  ┌───────────────┐  │
              │  │  Plugin System │  │
              │  │ product│order │  │
              │  │ im │marketing │  │
              │  │ wms│ai_image  │  │
              │  └───────────────┘  │
              │  ┌───────────────┐  │
              │  │ Driver Layer  │  │
              │  │payment│sms    │  │
              │  │oauth │storage │  │
              │  └───────────────┘  │
              └──────────┬──────────┘
                         ▼
              ┌──────┐ ┌───────┐
              │MySQL │ │ Redis │
              └──────┘ └───────┘
```

## 开源协议

LYShop 基于 [MIT License](https://github.com/ijry/lyshop/blob/master/LICENSE) 开源，可免费用于商业项目。
