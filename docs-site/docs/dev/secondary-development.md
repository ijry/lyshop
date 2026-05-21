# 二次开发文档

本文档介绍如何在 LYShop 中进行业务扩展与插件开发。

## 扩展原则

- 优先复用现有插件机制
- 新能力按业务域拆分，保持低耦合
- API、Service、Model 分层清晰

## 后端扩展点

- 插件定义：`server/plugins/<plugin-name>/plugin.go`
- 插件配置：`server/plugins/<plugin-name>/plugin.json`
- 业务模块：`api`、`service`、`model`

## 新增插件建议流程

1. 在 `server/plugins` 下创建新插件目录
2. 实现插件注册与初始化逻辑
3. 补充 API 路由与服务实现
4. 编写最小可用测试（模型、插件加载、关键业务）
5. 在管理端与移动端补齐调用与展示

## 前端扩展点

- 管理端路由：`admin/src/router`
- 管理端页面：`admin/src/views`
- 移动端页面：`app/pages`
- 公共请求封装：`admin/src/api`、`app/utils/request.ts`

## 联调建议

- 先定义接口契约，再并行开发
- 使用统一错误码与提示文案
- 对核心流程（登录、下单、支付）做回归用例

## 文档维护规范

- 新增模块需同步更新本章节
- 变更接口需同步更新接口文档
- 发布版本需记录兼容性说明
