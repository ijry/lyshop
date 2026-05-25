# docs-site 首页演示截图改造设计

## 目标

将 docs-site 首页从功能清单式介绍升级为更商业化、正规的商城官网首页，并使用真实演示模式截图展示 LYShop 的多端能力。

## 范围

- 自动生成 4 张主要功能截图：H5 首页、H5 商品列表、H5 订单列表、PC 商城首页。
- 截图保存到 `docs-site/docs/public/showcase/`，供 VitePress 首页直接引用。
- 改造 `docs-site/docs/.vitepress/theme/HomePage.vue`，首屏突出品牌价值、演示截图和核心商业能力。
- 更新 docs-site 文档，说明截图命令、资源路径和部署影响。

## 方案

新增 `docs-site/scripts/capture-showcase.js`。脚本优先访问已构建的 `docs-site/docs/public/demo/` 与 `docs-site/docs/public/web-demo/`，通过静态文件服务打开演示页面，再用 Playwright 截图。这样避免依赖后端服务，也能复用现有 mock 演示数据。

首页使用提交到仓库的截图资源，而不是运行时 iframe 截图，确保 GitHub Pages 和静态部署环境都能稳定展示。若截图资源暂未生成，页面仍展示静态占位结构，不影响文档构建。

## 验证

- 运行 `npm run showcase:capture` 生成截图。
- 运行 `npm run docs:build` 确认 VitePress 构建通过。
- 本地预览首页，检查桌面与移动宽度下首屏、截图、CTA 和功能区没有遮挡或横向滚动。
