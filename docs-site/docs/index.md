---
layout: home

hero:
  name: "LYShop"
  text: "开箱即用的电商全栈项目"
  tagline: "官网文档包含功能介绍、部署说明、接口规范与二次开发指南"
  actions:
    - theme: brand
      text: 立即开始
      link: /guide/features
    - theme: alt
      text: 部署文档
      link: /deploy/

features:
  - title: 全栈能力
    details: 后端、管理端、移动端一体化，覆盖电商核心业务流程。
  - title: 插件化架构
    details: 支持按业务域扩展插件，低耦合、易维护。
  - title: 文档体系
    details: 提供从部署到开发的完整文档，降低接入和维护成本。
---

## 快速导航

- [功能介绍](/guide/features)
- [部署文档](/deploy/)
- [接口文档](/api/)
- [二次开发](/dev/secondary-development)

## 项目架构

```text
Client(App/管理后台) -> API(Server) -> 插件层(商品/订单/营销/IM...) -> 基础设施(DB/缓存/存储/支付)
```
