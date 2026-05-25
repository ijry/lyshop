# docs-site AI花絮导航与文档整理 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在 docs-site 顶部导航最后新增「AI花絮」，并基于现有文档整理出索引页与三张子页（能力、接口、演进）。

**Architecture:** 仅调整 VitePress 配置与文档 Markdown 内容，不新增业务接口。导航增加 `/ai-notes/` 专区入口，专区内部通过 sidebar 管理章节跳转，内容严格引用现有 docs 来源。

**Tech Stack:** VitePress, Markdown, Node.js scripts（docs-site build）

---

## File Structure

- Modify: `docs-site/docs/.vitepress/config.mts`
  - 顶部导航新增「AI花絮」（最后一项）
  - 增加 `/ai-notes/` sidebar
- Create: `docs-site/docs/ai-notes/index.md`
  - AI 花絮索引页
- Create: `docs-site/docs/ai-notes/capabilities.md`
  - AI 能力总览（来源 features/README）
- Create: `docs-site/docs/ai-notes/api.md`
  - AI 接口速查（仅现有文档可见信息）
- Create: `docs-site/docs/ai-notes/evolution.md`
  - AI 演进记录与实践要点

### Task 1: 更新导航与侧边栏配置

**Files:**
- Modify: `docs-site/docs/.vitepress/config.mts`

- [ ] **Step 1: 在 nav 最后一项新增 AI花絮入口**

```ts
{ text: "AI花絮", link: "/ai-notes/" }
```

- [ ] **Step 2: 增加 /ai-notes/ 的 sidebar 分组与 4 个页面链接**

```ts
"/ai-notes/": [
  {
    text: "AI花絮",
    items: [
      { text: "导读", link: "/ai-notes/" },
      { text: "能力总览", link: "/ai-notes/capabilities" },
      { text: "接口速查", link: "/ai-notes/api" },
      { text: "演进记录", link: "/ai-notes/evolution" }
    ]
  }
]
```

### Task 2: 新增 AI 花絮索引页

**Files:**
- Create: `docs-site/docs/ai-notes/index.md`

- [ ] **Step 1: 编写导读结构与阅读路径**

```md
# AI 花絮

这里整理 LYShop AI 生图能力在文档中的关键信息，便于产品、研发、运营快速对齐。

## 阅读顺序

1. 能力总览：先看支持范围与业务定位
2. 接口速查：再看现有文档中的 API 可见信息
3. 演进记录：最后看架构与交互演进
```

- [ ] **Step 2: 添加“内容边界说明”**

```md
> 本专区仅整理现有文档事实，不新增未落地接口定义。
```

### Task 3: 新增能力总览页

**Files:**
- Create: `docs-site/docs/ai-notes/capabilities.md`

- [ ] **Step 1: 从 features/README 提炼能力要点**

```md
## 核心能力
- 多模型聚合：通义万象 / 文心一格 / 腾讯混元 / OpenAI DALL-E
- 生成目标：封面图、轮播图、详情图
- 参考图能力：按模型能力控制可用性
```

- [ ] **Step 2: 补充业务集成定位**

```md
## 业务集成位置
- AI 生图能力已内嵌商品编辑页
- 后台不再保留独立 AI 生图菜单入口
```

### Task 4: 新增接口速查页

**Files:**
- Create: `docs-site/docs/ai-notes/api.md`

- [ ] **Step 1: 编写“现有文档可见接口”与说明**

```md
## 文档可见接口
- GET /admin/api/ai/models
- GET /admin/api/ai/tasks
- POST /admin/api/ai/generate
- GET /admin/api/ai/tasks/:id
```

- [ ] **Step 2: 标注信息来源与一致性约束**

```md
> 接口字段以 docs-site 当前 API 文档与已公开示例为准。
> 若后端能力先行变更，请同步更新 api 文档后再更新本页。
```

### Task 5: 新增演进记录页

**Files:**
- Create: `docs-site/docs/ai-notes/evolution.md`

- [ ] **Step 1: 编写演进主线**

```md
## 演进主线
- 初期：独立 AI 能力展示
- 当前：能力下沉到商品编辑工作流
- 结果：减少跨菜单跳转，提高素材生产闭环效率
```

- [ ] **Step 2: 增加实践要点**

```md
## 实践建议
- 先选模型再写 Prompt，避免模型能力与目标不匹配
- 明确图片用途（封面/轮播/详情）再生成，降低返工
- 关键素材保留可追溯记录，便于复用与审计
```

### Task 6: 构建验证

**Files:**
- Test: `docs-site` build

- [ ] **Step 1: 运行文档构建**

Run: `npm run docs:build`（`docs-site` 目录）  
Expected: 构建通过，`/ai-notes/` 页面被产出。

- [ ] **Step 2: 快速路径检查**

Run: 在构建输出中确认以下页面存在：
- `docs/.vitepress/dist/ai-notes/index.html`
- `docs/.vitepress/dist/ai-notes/capabilities.html`
- `docs/.vitepress/dist/ai-notes/api.html`
- `docs/.vitepress/dist/ai-notes/evolution.html`

## Self-Review

- Spec 覆盖：导航最后项、索引+子页、内容来源约束均有任务对应。
- Placeholder 扫描：无 TBD/TODO/后续补充等占位语句。
- 一致性：路径与命名统一使用 `/ai-notes/`。
