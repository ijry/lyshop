# PC 装修高级样式系统 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 完成 PC 装修高级样式系统落地，支持页面级背景与组件级样式覆盖，并确保 admin/web/app Mock 数据结构一致。

**Architecture:** 继续复用现有 `/decor/pc` 接口路径与后端 JSON 存储字段，将 PC 页面配置统一为 `PcDecorPage = { pageStyle, components }`。Admin 编辑器维护该结构并实时驱动预览，Web 前台按同一结构渲染。Mock（admin/web/app）统一升级，避免数据分叉。

**Tech Stack:** Vue 3 + TypeScript（admin/web）、Gin + Gorm（server）、Mock preset（app/admin/web）、docs-site Markdown

---

## File Structure

- Modify: `admin/src/types/decor.ts`
  - 新增 `PcDecorPageStyle`、`PcDecorComponentStyle`、`PcDecorPagePayload` 类型与默认值工厂。
- Create: `admin/src/views/decor/widgets/PcPageStyleEditor.vue`
  - 页面样式配置面板（背景模式、渐变、图片、遮罩、内容宽度、默认圆角阴影）。
- Create: `admin/src/views/decor/widgets/PcComponentStyleEditor.vue`
  - 组件样式覆盖面板（margin/padding/background/border/radius/shadow）。
- Modify: `admin/src/views/decor/PcDecorEditor.vue`
  - 编辑状态从 `components[]` 升级为 `pagePayload`，保存/加载逻辑改造，挂接两个新面板。
- Modify: `admin/src/views/decor/PcDecorPreview.vue`
  - 接收 `pageStyle`，应用页面背景与组件样式覆盖渲染规则。
- Modify: `admin/src/mock/index.ts`
  - `pcDecorSource` 与 `PUT /admin/api/decor/pc` 改为完整 `PcDecorPage`。
- Modify: `app/mock/presets/types.ts`
  - `pcDecor` 类型从 `{ components: any[] }` 升级为 `PcDecorPagePayload`。
- Modify: `app/mock/presets/mall.ts`
- Modify: `app/mock/presets/cake.ts`
- Modify: `app/mock/presets/farm.ts`
- Modify: `app/mock/presets/fresh.ts`
- Modify: `app/mock/presets/jewelry.ts`
- Modify: `app/mock/presets/mother.ts`
- Modify: `app/mock/presets/supermarket.ts`
  - 各预设 `pcDecor` 补齐 `pageStyle` 并为组件增加 `style` 空对象（或按主题给默认值）。
- Modify: `web/src/views/Home.vue`
  - 读取并传递 `pageStyle + components`。
- Modify: `web/src/components/decor/DecorRenderer.vue`
  - 页面背景渲染与组件样式覆盖渲染。
- Modify: `server/plugins/decor/service/decor.go`
  - 为 `page_key=pc` 提供默认空 payload（非数组）。
- Create: `server/plugins/decor/service/decor_pc_payload_test.go`
  - 测试 PC 默认 payload 与校验回退逻辑。
- Modify: `docs-site/docs/dev/secondary-development.md`
  - 文档补充 PC 装修样式结构、接口载荷变化、Mock 说明。

---

### Task 1: 定义 PC 样式 Schema 与默认值（Admin 类型层）

**Files:**
- Modify: `admin/src/types/decor.ts`

- [ ] **Step 1: 先写类型断言（“失败用例”）注释块，明确目标结构**

```ts
// 目标：PcDecorPagePayload 必须包含 pageStyle 与 components
// 目标：每个组件可选 style 覆盖字段（margin/padding/background/border/radius/shadow）
```

- [ ] **Step 2: 增加页面样式与组件样式类型**

```ts
export interface PcDecorPageStyle {
  background: {
    mode: 'solid' | 'gradient' | 'image'
    solidColor?: string
    gradient?: { angle: number; stops: Array<{ color: string; position: number }> }
    image?: {
      url: string
      size: 'cover' | 'contain' | 'auto' | 'custom'
      customSize?: string
      position: string
      repeat: 'no-repeat' | 'repeat' | 'repeat-x' | 'repeat-y'
      attachment: 'scroll' | 'fixed'
    }
    overlay?: { enabled: boolean; color: string; opacity: number }
  }
  content: { maxWidth: number; gutterX: number; sectionGap: number }
  surface: { radius: number; shadow: 'none' | 'sm' | 'md' | 'lg' }
}
```

- [ ] **Step 3: 增加 payload 与默认值工厂**

```ts
export interface PcDecorPagePayload {
  pageStyle: PcDecorPageStyle
  components: Array<PcDecorComponent & { style?: PcDecorComponentStyle }>
}

export function createDefaultPcPageStyle(): PcDecorPageStyle {
  return {
    background: { mode: 'solid', solidColor: '#f8fafc', overlay: { enabled: false, color: '#000000', opacity: 0.2 } },
    content: { maxWidth: 1280, gutterX: 24, sectionGap: 24 },
    surface: { radius: 12, shadow: 'none' },
  }
}
```

- [ ] **Step 4: 给 createPcDefaultProps 产物预留 style 字段挂载约定（不在 props 内）**

```ts
// append 时统一在组件对象根层增加 style?: PcDecorComponentStyle
```

- [ ] **Step 5: 运行 admin 类型检查**

Run: `npm run build`（工作目录 `admin`）  
Expected: `vue-tsc` 通过，无类型错误。

- [ ] **Step 6: Commit**

```bash
git add admin/src/types/decor.ts
git commit -m "定义PC装修样式Schema与默认值" -m "新增页面级与组件级样式类型，并补充PC页面payload默认值工厂。"
```

---

### Task 2: 后端 PC 默认载荷改造与测试

**Files:**
- Modify: `server/plugins/decor/service/decor.go`
- Create: `server/plugins/decor/service/decor_pc_payload_test.go`

- [ ] **Step 1: 写失败测试，断言 page_key=pc 返回对象结构而非数组**

```go
func TestDefaultPayloadForPCPage(t *testing.T) {
  got := defaultComponentsForPage("pc")
  require.JSONEq(t, `{"pageStyle":{"background":{"mode":"solid","solidColor":"#f8fafc","overlay":{"enabled":false,"color":"#000000","opacity":0.2}},"content":{"maxWidth":1280,"gutterX":24,"sectionGap":24},"surface":{"radius":12,"shadow":"none"}},"components":[]}`, string(got))
}
```

- [ ] **Step 2: 运行失败测试**

Run: `go test ./plugins/decor/service -run TestDefaultPayloadForPCPage -v`（工作目录 `server`）  
Expected: FAIL（`defaultComponentsForPage` 未定义）。

- [ ] **Step 3: 实现默认 payload 分发函数**

```go
func defaultComponentsForPage(pageKey string) json.RawMessage {
  if strings.EqualFold(strings.TrimSpace(pageKey), "pc") {
    return json.RawMessage(`{"pageStyle":{"background":{"mode":"solid","solidColor":"#f8fafc","overlay":{"enabled":false,"color":"#000000","opacity":0.2}},"content":{"maxWidth":1280,"gutterX":24,"sectionGap":24},"surface":{"radius":12,"shadow":"none"}},"components":[]}`)
  }
  return emptyComponents()
}
```

- [ ] **Step 4: 在 GetPage/ListVariants/CreateVariantCopy 默认分支中使用新函数**

```go
Components: defaultComponentsForPage(pageKey),
```

- [ ] **Step 5: 运行通过测试**

Run: `go test ./plugins/decor/service -run TestDefaultPayloadForPCPage -v`  
Expected: PASS。

- [ ] **Step 6: 回归装饰服务包测试**

Run: `go test ./plugins/decor/service -v`  
Expected: PASS。

- [ ] **Step 7: Commit**

```bash
git add server/plugins/decor/service/decor.go server/plugins/decor/service/decor_pc_payload_test.go
git commit -m "调整PC装修默认载荷并补充测试" -m "按page_key为pc返回页面级样式对象默认结构，避免空数组导致前端解析分歧。"
```

---

### Task 3: Admin Mock 升级为完整 PcDecorPage

**Files:**
- Modify: `admin/src/mock/index.ts`

- [ ] **Step 1: 更新 pcDecorSource 初始结构**

```ts
const pcDecorSource: any = {
  id: 1,
  components: JSON.stringify({
    pageStyle: { /* 默认 pageStyle */ },
    components: [ /* 原有 pc 组件 */ ],
  }),
  is_published: true,
}
```

- [ ] **Step 2: 更新 PUT /admin/api/decor/pc 保存逻辑**

```ts
if (key === 'PUT /admin/api/decor/pc') {
  const payload = (params as any)?.components
  pcDecorSource.components = JSON.stringify(payload || { pageStyle: defaultPageStyle, components: [] })
  return { matched: true, data: clone(pcDecorSource) }
}
```

- [ ] **Step 3: 校验 GET 返回结构与保存一致**

Run: `npm run build:demo`（工作目录 `admin`）  
Expected: 构建通过，无 mock 语法错误。

- [ ] **Step 4: Commit**

```bash
git add admin/src/mock/index.ts
git commit -m "升级Admin端PC装修Mock结构" -m "pcDecor改为页面级样式+组件数组的完整对象，PUT接口按新结构整体保存。"
```

---

### Task 4: App 预设类型与所有行业预设同步

**Files:**
- Modify: `app/mock/presets/types.ts`
- Modify: `app/mock/presets/mall.ts`
- Modify: `app/mock/presets/cake.ts`
- Modify: `app/mock/presets/farm.ts`
- Modify: `app/mock/presets/fresh.ts`
- Modify: `app/mock/presets/jewelry.ts`
- Modify: `app/mock/presets/mother.ts`
- Modify: `app/mock/presets/supermarket.ts`

- [ ] **Step 1: 升级预设类型声明**

```ts
export interface PcDecorPagePayload {
  pageStyle: { /* 与 admin 对齐 */ }
  components: any[]
}

export interface MockPreset {
  // ...
  pcDecor: PcDecorPagePayload
}
```

- [ ] **Step 2: 批量为各预设 pcDecor 添加 pageStyle**

```ts
pcDecor: {
  pageStyle: {
    background: { mode: 'solid', solidColor: '#f8fafc', overlay: { enabled: false, color: '#000000', opacity: 0.2 } },
    content: { maxWidth: 1280, gutterX: 24, sectionGap: 24 },
    surface: { radius: 12, shadow: 'none' },
  },
  components: [ /* 原有数组 */ ],
},
```

- [ ] **Step 3: 运行 app 构建检查**

Run: `npm run build:h5:demo`（工作目录 `app`）  
Expected: 通过，无预设类型错误。

- [ ] **Step 4: Commit**

```bash
git add app/mock/presets/types.ts app/mock/presets/mall.ts app/mock/presets/cake.ts app/mock/presets/farm.ts app/mock/presets/fresh.ts app/mock/presets/jewelry.ts app/mock/presets/mother.ts app/mock/presets/supermarket.ts
git commit -m "统一App端预设PC装修数据结构" -m "所有行业预设pcDecor补充pageStyle，确保与admin/web mock结构一致。"
```

---

### Task 5: Admin 编辑器切换为 PcDecorPage 状态模型

**Files:**
- Modify: `admin/src/views/decor/PcDecorEditor.vue`
- Create: `admin/src/views/decor/widgets/PcPageStyleEditor.vue`
- Create: `admin/src/views/decor/widgets/PcComponentStyleEditor.vue`

- [ ] **Step 1: 在编辑器中引入 pagePayload 状态并替换 components 顶层状态**

```ts
const pagePayload = ref<PcDecorPagePayload>({
  pageStyle: createDefaultPcPageStyle(),
  components: [],
})
const components = computed({
  get: () => pagePayload.value.components,
  set: (list) => { pagePayload.value.components = list as any[] },
})
```

- [ ] **Step 2: 新增页面样式面板组件并接入右侧配置区**

```vue
<PcPageStyleEditor v-model="pagePayload.pageStyle" />
```

- [ ] **Step 3: 新增组件样式面板并接入选中组件**

```vue
<PcComponentStyleEditor v-if="selectedComp" v-model="selectedComp.style" />
```

- [ ] **Step 4: 改造保存/加载逻辑为整包读写**

```ts
await request.put('/decor/pc', { components: pagePayload.value })

const data: any = await request.get('/decor/pc')
const raw = typeof data?.components === 'string' ? JSON.parse(data.components) : data?.components
pagePayload.value = normalizePcDecorPayload(raw)
```

- [ ] **Step 5: 运行 admin 构建**

Run: `npm run build`（工作目录 `admin`）  
Expected: PASS。

- [ ] **Step 6: Commit**

```bash
git add admin/src/views/decor/PcDecorEditor.vue admin/src/views/decor/widgets/PcPageStyleEditor.vue admin/src/views/decor/widgets/PcComponentStyleEditor.vue
git commit -m "改造PC装修编辑器为页面级样式模型" -m "编辑状态升级为PcDecorPage，新增页面样式与组件样式覆盖配置面板。"
```

---

### Task 6: Admin 预览渲染实现页面背景与组件样式覆盖

**Files:**
- Modify: `admin/src/views/decor/PcDecorPreview.vue`

- [ ] **Step 1: 扩展组件入参**

```ts
const props = defineProps<{ components: any[]; pageStyle: any }>()
```

- [ ] **Step 2: 实现页面容器背景样式计算**

```ts
const pageBackgroundStyle = computed(() => {
  // 根据 mode 生成 background / backgroundImage / backgroundSize / backgroundPosition / backgroundRepeat / backgroundAttachment
})
```

- [ ] **Step 3: 实现组件级 style 合并函数**

```ts
function sectionStyle(comp: any) {
  return {
    marginTop: px(comp?.style?.marginTop),
    marginBottom: px(comp?.style?.marginBottom ?? props.pageStyle?.content?.sectionGap),
    paddingLeft: px(comp?.style?.paddingX),
    paddingRight: px(comp?.style?.paddingX),
    borderRadius: px(comp?.style?.borderRadius ?? props.pageStyle?.surface?.radius),
    // shadow / border / bg
  }
}
```

- [ ] **Step 4: 在模板根层与各 section 容器应用计算样式**

```vue
<div class="min-h-full relative" :style="pageBackgroundStyle">
  <div class="absolute inset-0 pointer-events-none" :style="overlayStyle" />
  <section :style="sectionStyle(comp)">...</section>
</div>
```

- [ ] **Step 5: 运行 admin 构建**

Run: `npm run build`（工作目录 `admin`）  
Expected: PASS。

- [ ] **Step 6: Commit**

```bash
git add admin/src/views/decor/PcDecorPreview.vue
git commit -m "增强PC装修预览样式渲染能力" -m "支持页面背景模式、遮罩与组件样式覆盖，预览效果与前台语义对齐。"
```

---

### Task 7: Web 前台消费 PcDecorPage 并渲染样式

**Files:**
- Modify: `web/src/views/Home.vue`
- Modify: `web/src/components/decor/DecorRenderer.vue`

- [ ] **Step 1: Home 页改为接收 pageStyle + components**

```ts
const pageStyle = ref<any>(null)
const components = ref<any[]>([])
const data = await get<any>('/api/v1/pc/decor')
const raw = data?.components || {}
pageStyle.value = raw?.pageStyle || createFallbackPageStyle()
components.value = Array.isArray(raw?.components) ? raw.components : []
```

- [ ] **Step 2: Renderer 新增 pageStyle 入参与背景/遮罩渲染**

```ts
const props = defineProps<{ components: any[]; pageStyle?: any }>()
```

- [ ] **Step 3: 将组件 section 样式应用逻辑与 Admin 预览保持同构**

```ts
function resolveSectionStyle(comp: any) { /* 同 admin 规则 */ }
```

- [ ] **Step 4: 运行 web 构建（含 mock 模式）**

Run: `npm run build`（工作目录 `web`）  
Expected: PASS。  

Run: `npm run build:demo`（工作目录 `web`）  
Expected: PASS。

- [ ] **Step 5: Commit**

```bash
git add web/src/views/Home.vue web/src/components/decor/DecorRenderer.vue
git commit -m "前台支持PC装修页面级与组件级样式渲染" -m "首页读取PcDecorPage结构，渲染页面背景与组件覆写样式，mock与真实接口语义一致。"
```

---

### Task 8: docs-site 文档同步

**Files:**
- Modify: `docs-site/docs/dev/secondary-development.md`

- [ ] **Step 1: 补充 PC 装修新结构说明**

```md
### PC 装修页面配置结构（PcDecorPage）
- pageStyle.background / content / surface
- components[].style（可选覆盖）
```

- [ ] **Step 2: 补充接口载荷变化示例**

```json
{
  "components": {
    "pageStyle": { "...": "..." },
    "components": [{ "id": "pc_hero", "type": "hero", "props": {}, "style": {} }]
  }
}
```

- [ ] **Step 3: 补充 Mock 同步说明（admin/web/app）**

```md
三端 `pcDecor` 结构统一为 PcDecorPage；新增字段必须三端同时更新。
```

- [ ] **Step 4: Commit**

```bash
git add docs-site/docs/dev/secondary-development.md
git commit -m "更新PC装修高级样式开发文档" -m "补充PcDecorPage结构、接口示例与三端Mock同步规范。"
```

---

### Task 9: 端到端回归与打包验证

**Files:**
- Modify: 无（验证任务）

- [ ] **Step 1: Server 测试回归**

Run: `go test ./...`（工作目录 `server`）  
Expected: PASS。

- [ ] **Step 2: Admin 构建回归**

Run: `npm run build`（工作目录 `admin`）  
Expected: PASS。

- [ ] **Step 3: Web 构建回归（普通+demo）**

Run: `npm run build && npm run build:demo`（工作目录 `web`）  
Expected: PASS。

- [ ] **Step 4: 手工验收 checklist（Mock 模式）**

```md
1) admin dev:demo 打开 /decor/pc
2) 调整页面背景（纯色→渐变→图片）后中间预览即时变化
3) 给任意组件设置圆角/阴影/边距，预览即时变化
4) 点击保存后刷新页面，配置完整回显
5) web dev:demo 打开首页，背景与组件样式一致展示
```

- [ ] **Step 5: 汇总提交（若需）**

```bash
git status
```

---

## Self-Review

### 1) Spec coverage
- 页面级背景与遮罩：Task 5/6/7
- 组件级样式覆盖：Task 5/6/7
- Mock 三端同步：Task 3/4/7
- 接口升级但路径不变：Task 2/5/7
- docs-site 同步：Task 8

### 2) Placeholder scan
- 已检查无 TBD/TODO/“后续补充”等占位语句。

### 3) Type consistency
- 统一使用 `PcDecorPagePayload` / `pageStyle` / `components[].style` 命名，前后端与 mock 同名字段。

