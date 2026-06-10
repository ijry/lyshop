# PC Decor Variants Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 让后台 PC 首页装修支持与移动端首页一致的多副本创建、切换、保存、发布、重命名和删除能力。

**Architecture:** 后端复用现有 `decor/:page_key` 多副本服务，PC 通过 `page_key=pc` 进入同一套 `variant_key` 流程。Admin PC 编辑器补齐版本状态与操作按钮，Admin Mock 从单对象改为 PC 副本数组，docs-site 更新为当前通用装修多副本架构。

**Tech Stack:** Go + Gin + Gorm + SQLite 测试；Vue 3 + TypeScript + Vite；Admin Mock；VitePress docs-site。

---

## File Structure

- Modify: `server/plugins/decor/service/decor_pc_payload_test.go`
  - 增加内存数据库测试，证明 `page_key=pc` 可列出默认副本、复制 PC payload、发布 PC 副本且不影响 `index`。
- Modify: `admin/src/views/decor/PcDecorEditor.vue`
  - 引入 `useI18n`、`confirmAction`、`promptText`，增加版本列表、当前版本、复制、重命名、删除、按版本保存和发布。
- Modify: `admin/src/mock/index.ts`
  - 将 `pcDecorSource` 改为 `pcDecorVariantsSource`，补齐 PC versions/copies/variant CRUD/publish Mock 分支。
- Modify: `docs-site/docs/api/decor.md`
  - 将装修多副本描述从“首页装修”升级为“任意页面装修”，补充 PC 首页装修多版本语义。

---

### Task 1: 后端补充 PC 多副本服务测试

**Files:**
- Modify: `server/plugins/decor/service/decor_pc_payload_test.go`

- [ ] **Step 1: 增加测试依赖 import**

在 `server/plugins/decor/service/decor_pc_payload_test.go` 中把 import 改为：

```go
import (
	"context"
	"encoding/json"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/core/db"
	decormodel "github.com/ijry/lyshop/plugins/decor/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)
```

- [ ] **Step 2: 增加内存数据库 helper**

在现有 `TestDefaultComponentsForPageNonPC` 后追加：

```go
func setupDecorServiceTestDB(t *testing.T) {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open("file:decor-service-test?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, gdb.AutoMigrate(&decormodel.DecorPage{}))
	db.DB = gdb
}
```

- [ ] **Step 3: 写 PC 默认副本列表测试**

继续追加：

```go
func TestListVariantsForPCReturnsDefaultPayload(t *testing.T) {
	setupDecorServiceTestDB(t)

	rows, err := ListVariants(context.Background(), 0, "pc")

	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, "pc", rows[0].PageKey)
	require.Equal(t, DefaultVariantKey, rows[0].VariantKey)
	require.JSONEq(t, `{"pageStyle":{"background":{"mode":"solid","solidColor":"#f8fafc","overlay":{"enabled":false,"color":"#000000","opacity":0.2}},"content":{"maxWidth":1280,"gutterX":24,"sectionGap":24},"surface":{"radius":12,"shadow":"none"}},"components":[]}`, string(rows[0].Components))
}
```

- [ ] **Step 4: 写 PC 副本复制测试**

继续追加：

```go
func TestCreateVariantCopyForPCKeepsPagePayload(t *testing.T) {
	setupDecorServiceTestDB(t)
	ctx := context.Background()
	payload := json.RawMessage(`{"pageStyle":{"background":{"mode":"solid","solidColor":"#111111"},"content":{"maxWidth":1180,"gutterX":20,"sectionGap":18},"surface":{"radius":8,"shadow":"sm"}},"components":[{"id":"pc_hero","type":"hero","props":{"title":"A"}}]}`)
	_, err := SavePage(ctx, 0, "pc", payload, "default")
	require.NoError(t, err)

	row, err := CreateVariantCopy(ctx, 0, "pc", "default", "summer", "夏季版")

	require.NoError(t, err)
	require.Equal(t, "pc", row.PageKey)
	require.Equal(t, "summer", row.VariantKey)
	require.Equal(t, "夏季版", row.VariantName)
	require.False(t, row.IsPublished)
	require.JSONEq(t, string(payload), string(row.Components))
}
```

- [ ] **Step 5: 写 PC 发布单活隔离测试**

继续追加：

```go
func TestPublishPageForPCIsSingleActiveWithinPCOnly(t *testing.T) {
	setupDecorServiceTestDB(t)
	ctx := context.Background()
	pcPayload := json.RawMessage(`{"pageStyle":{"background":{"mode":"solid","solidColor":"#f8fafc"},"content":{"maxWidth":1280,"gutterX":24,"sectionGap":24},"surface":{"radius":12,"shadow":"none"}},"components":[]}`)
	indexPayload := json.RawMessage(`[{"id":"m_1","type":"banner","props":{}}]`)
	_, err := SavePage(ctx, 0, "pc", pcPayload, "default")
	require.NoError(t, err)
	_, err = CreateVariantCopy(ctx, 0, "pc", "default", "festival", "节日版")
	require.NoError(t, err)
	_, err = SavePage(ctx, 0, "index", indexPayload, "default")
	require.NoError(t, err)
	require.NoError(t, PublishPage(ctx, 0, "index", "default"))
	require.NoError(t, PublishPage(ctx, 0, "pc", "default"))

	require.NoError(t, PublishPage(ctx, 0, "pc", "festival"))

	var pcRows []decormodel.DecorPage
	require.NoError(t, db.DB.Where("page_key = ?", "pc").Order("variant_key asc").Find(&pcRows).Error)
	require.Len(t, pcRows, 2)
	require.False(t, pcRows[0].IsPublished)
	require.Equal(t, "default", pcRows[0].VariantKey)
	require.True(t, pcRows[1].IsPublished)
	require.Equal(t, "festival", pcRows[1].VariantKey)

	indexPublished, err := GetPublishedPage(ctx, 0, "index")
	require.NoError(t, err)
	require.Equal(t, "index", indexPublished.PageKey)
	require.Equal(t, "default", indexPublished.VariantKey)
	require.True(t, indexPublished.IsPublished)
}
```

- [ ] **Step 6: 运行服务层测试**

Run: `go test ./plugins/decor/service -v`（工作目录 `server`）

Expected: PASS。

- [ ] **Step 7: Commit**

```bash
git add server/plugins/decor/service/decor_pc_payload_test.go
git commit -m "补充PC装修多版本服务测试" -m "覆盖PC默认副本、复制副本和发布单活隔离，确认复用decor通用多副本能力。"
```

---

### Task 2: Admin PC 编辑器接入版本管理

**Files:**
- Modify: `admin/src/views/decor/PcDecorEditor.vue`

- [ ] **Step 1: 修改顶部模板为版本工具栏**

将 `<h2 class="text-xl font-semibold text-slate-800">PC 首页装修</h2>` 所在区域替换为：

```vue
<div class="flex items-center gap-3">
  <h2 class="text-xl font-semibold text-slate-800">PC 首页装修</h2>
  <select v-model="currentVariantKey" @change="changeVariant"
    class="border border-slate-200 rounded-lg px-3 py-1.5 text-sm text-slate-700 bg-white">
    <option v-for="v in variants" :key="v.variant_key" :value="v.variant_key">
      {{ v.variant_name }}（{{ v.variant_key }}）{{ v.is_published ? ' · ' + t('decor.published') : '' }}
    </option>
  </select>
  <button @click="copyVariant"
    class="px-3 py-1.5 bg-slate-100 text-slate-700 rounded-lg text-xs hover:bg-slate-200 transition">
    {{ t('decor.copyVariant') }}
  </button>
  <button @click="renameVariant"
    class="px-3 py-1.5 bg-slate-100 text-slate-700 rounded-lg text-xs hover:bg-slate-200 transition">
    {{ t('decor.rename') }}
  </button>
  <button @click="deleteVariant"
    class="px-3 py-1.5 bg-red-50 text-red-600 rounded-lg text-xs hover:bg-red-100 transition">
    {{ t('decor.deleteVariant') }}
  </button>
</div>
```

- [ ] **Step 2: 改保存与发布按钮文案**

把保存按钮内容替换为：

```vue
{{ saving ? '保存中...' : t('decor.saveDraft') }}
```

把发布按钮内容替换为：

```vue
{{ t('decor.publish') }}
```

- [ ] **Step 3: 增加 imports**

在 `<script setup lang="ts">` 顶部 import 区域改为：

```ts
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'
import { notify } from '@/utils/notify'
import { confirmAction, promptText } from '@/utils/dialog'
```

- [ ] **Step 4: 增加 i18n 与版本状态**

在类型 import 后、`pagePayload` 前增加：

```ts
const { t } = useI18n()
```

在 `const pagePayload = ...` 后增加：

```ts
const variants = ref<any[]>([])
const currentVariantKey = ref('default')
```

- [ ] **Step 5: 修改 save/publish/load 流程**

替换现有 `save`、`publish`、`onMounted` 三段为：

```ts
async function save() {
  saving.value = true
  try {
    await request.put(`/decor/pc?variant=${encodeURIComponent(currentVariantKey.value)}`, { components: pagePayload.value })
    await loadVariants()
    notify('保存成功')
  } finally { saving.value = false }
}

async function publish() {
  await save()
  await request.post(`/decor/pc/publish?variant=${encodeURIComponent(currentVariantKey.value)}`)
  await loadVariants()
  notify(t('decor.publishedNote'))
}

async function loadVariants() {
  const data: any = await request.get('/decor/pc/variants')
  variants.value = Array.isArray(data) ? data : []
  if (!variants.value.length) {
    variants.value = [{
      variant_key: 'default',
      variant_name: t('decor.variantName', { key: 'default' }),
      is_published: false,
    }]
  }
  const currentExists = variants.value.some(v => v.variant_key === currentVariantKey.value)
  if (!currentExists) {
    const published = variants.value.find(v => v.is_published)
    currentVariantKey.value = published?.variant_key || variants.value[0].variant_key || 'default'
  }
}

async function loadCurrentVariant() {
  const data: any = await request.get(`/decor/pc?variant=${encodeURIComponent(currentVariantKey.value)}`)
  pagePayload.value = normalizePayload(data?.components)
  selectedIndex.value = null
}

async function changeVariant() {
  await loadCurrentVariant()
}

function toVariantKey(raw: string) {
  return raw.trim().toLowerCase().replace(/\s+/g, '_').replace(/[^a-z0-9_-]/g, '')
}

async function copyVariant() {
  const keyRaw = promptText(t('decor.promptKey'))
  if (!keyRaw) return
  const newVariantKey = toVariantKey(keyRaw)
  if (!newVariantKey) {
    notify(t('decor.invalidKey'))
    return
  }
  const defaultName = t('decor.variantName', { key: newVariantKey })
  const newVariantName = promptText(t('decor.promptName'), defaultName) || defaultName
  await request.post('/decor/pc/copies', {
    from_variant_key: currentVariantKey.value,
    new_variant_key: newVariantKey,
    new_variant_name: newVariantName,
  })
  await loadVariants()
  currentVariantKey.value = newVariantKey
  await loadCurrentVariant()
}

async function renameVariant() {
  const current = variants.value.find(v => v.variant_key === currentVariantKey.value)
  const next = promptText(t('decor.promptName'), current?.variant_name || '')
  if (!next) return
  await request.put(`/decor/pc/variants/${encodeURIComponent(currentVariantKey.value)}`, {
    variant_name: next,
  })
  await loadVariants()
}

async function deleteVariant() {
  if (currentVariantKey.value === 'default') {
    notify(t('decor.defaultNoDelete'))
    return
  }
  if (!confirmAction(t('decor.confirmDeleteVariant', { key: currentVariantKey.value }))) return
  await request.delete(`/decor/pc/variants/${encodeURIComponent(currentVariantKey.value)}`)
  await loadVariants()
  await loadCurrentVariant()
}

onMounted(async () => {
  await loadVariants()
  await loadCurrentVariant()
})
```

- [ ] **Step 6: 运行 Admin 构建**

Run: `npm run build`（工作目录 `admin`）

Expected: PASS。

- [ ] **Step 7: Commit**

```bash
git add admin/src/views/decor/PcDecorEditor.vue
git commit -m "接入PC装修多版本编辑管理" -m "PC装修编辑器复用decor通用副本接口，支持切换、复制、重命名、删除、保存和发布指定副本。"
```

---

### Task 3: Admin Mock 支持 PC 多副本接口

**Files:**
- Modify: `admin/src/mock/index.ts`

- [ ] **Step 1: 将单对象 PC source 改为数组**

找到 `const pcDecorSource` 定义，将其替换为：

```ts
const pcDecorVariantsSource: any[] = [{
  id: 101,
  page_key: 'pc',
  variant_key: 'default',
  variant_name: '默认副本',
  components: JSON.stringify(createDefaultPcDecorPayload()),
  is_published: true,
  published_at: '2026-05-28T10:00:00Z',
}]
```

如果 `createDefaultPcDecorPayload` 不在当前文件可用，则在现有 `pcDecorSource` 组件结构基础上包成同样数组，保持 `components` 为完整 PC 页面载荷 JSON 字符串。

- [ ] **Step 2: 替换 PC Decor Mock 分支**

把当前 `// PC Decor` 下三段 `GET/PUT/POST /admin/api/decor/pc` 分支替换为：

```ts
  // PC Decor
  if (key === 'GET /admin/api/decor/pc/variants') {
    return { matched: true, data: clone(pcDecorVariantsSource) }
  }
  if (key.startsWith('GET /admin/api/decor/pc')) {
    const parsed = new URL(url, 'https://mock.local')
    const variantKey = String(parsed.searchParams.get('variant') || 'default')
    const item = pcDecorVariantsSource.find((row: any) => String(row.variant_key) === variantKey) || pcDecorVariantsSource[0]
    return { matched: true, data: clone(item) }
  }
  if (key.startsWith('PUT /admin/api/decor/pc')) {
    const parsed = new URL(url, 'https://mock.local')
    const variantKey = String(parsed.searchParams.get('variant') || 'default')
    const target = pcDecorVariantsSource.find((row: any) => String(row.variant_key) === variantKey)
    if (target) {
      target.components = JSON.stringify((params as any)?.components || { pageStyle: null, components: [] })
    }
    return { matched: true, data: target ? clone(target) : null }
  }
  if (key.startsWith('POST /admin/api/decor/pc/publish')) {
    const parsed = new URL(url, 'https://mock.local')
    const variantKey = String(parsed.searchParams.get('variant') || 'default')
    const now = new Date().toISOString()
    for (const row of pcDecorVariantsSource) {
      row.is_published = false
      row.published_at = null
    }
    const target = pcDecorVariantsSource.find((row: any) => String(row.variant_key) === variantKey)
    if (target) {
      target.is_published = true
      target.published_at = now
    }
    return { matched: true, data: null }
  }
  if (key === 'POST /admin/api/decor/pc/copies') {
    const payload: any = params || {}
    const fromVariantKey = String(payload.from_variant_key || 'default')
    const newVariantKey = String(payload.new_variant_key || '').trim()
    const newVariantName = String(payload.new_variant_name || '').trim() || `副本 ${newVariantKey}`
    const source = pcDecorVariantsSource.find((row: any) => String(row.variant_key) === fromVariantKey)
    if (!source || !newVariantKey) return { matched: true, data: null }
    if (!pcDecorVariantsSource.find((row: any) => String(row.variant_key) === newVariantKey)) {
      pcDecorVariantsSource.push({
        ...clone(source),
        id: Math.max(...pcDecorVariantsSource.map((row: any) => Number(row.id || 0))) + 1,
        variant_key: newVariantKey,
        variant_name: newVariantName,
        is_published: false,
        published_at: null,
      })
    }
    return { matched: true, data: null }
  }
  if (key.startsWith('PUT /admin/api/decor/pc/variants/')) {
    const variantKey = decodeURIComponent(url.split('/').pop() || '')
    const target = pcDecorVariantsSource.find((row: any) => String(row.variant_key) === variantKey)
    if (target) {
      target.variant_name = String((params as any)?.variant_name || target.variant_name)
    }
    return { matched: true, data: null }
  }
  if (key.startsWith('DELETE /admin/api/decor/pc/variants/')) {
    const variantKey = decodeURIComponent(url.split('/').pop() || '')
    if (variantKey !== 'default') {
      const idx = pcDecorVariantsSource.findIndex((row: any) => String(row.variant_key) === variantKey && !row.is_published)
      if (idx >= 0) pcDecorVariantsSource.splice(idx, 1)
    }
    return { matched: true, data: null }
  }
```

- [ ] **Step 3: 运行 Admin demo 构建**

Run: `npm run build:demo`（工作目录 `admin`）

Expected: PASS。

- [ ] **Step 4: Commit**

```bash
git add admin/src/mock/index.ts
git commit -m "补齐PC装修多版本Mock接口" -m "Admin Mock为PC装修维护独立副本列表，覆盖版本列表、复制、重命名、删除、保存与单活发布。"
```

---

### Task 4: docs-site 更新装修接口文档

**Files:**
- Modify: `docs-site/docs/api/decor.md`

- [ ] **Step 1: 更新说明段落**

将“装修模块支持页面组件化配置、草稿保存、发布上线，以及首页多副本管理（单活发布）。”替换为：

```md
装修模块支持页面组件化配置、草稿保存、发布上线，以及任意页面多副本管理（单活发布）。移动端首页通常使用 `page_key=index`，PC 首页使用 `page_key=pc`。
```

- [ ] **Step 2: 更新多副本标题与说明**

将 `## 多副本（首页装修）` 改为：

```md
## 多副本（任意页面装修）
```

将第一条说明替换为：

```md
- 同一 `page_key`（如 `index`、`pc`）下可存在多个副本（`variant_key`）。
```

- [ ] **Step 3: 更新发布规则**

将“前台 `GET /api/v1/index/decor` 仅返回当前已发布副本。”替换为：

```md
- 前台 `GET /api/v1/index/decor` 返回 `index` 当前已发布副本；`GET /api/v1/decor/pc` 返回 PC 首页当前已发布副本。
```

- [ ] **Step 4: 补充 PC 装修样式系统版本语义**

在“PC 装修样式系统”功能说明列表中追加：

```md
  - PC 首页装修支持多副本管理，每个副本保存完整 `PcDecorPage` 载荷。
```

在该章节接口列表中追加：

```md
  - `GET /admin/api/decor/pc/variants`
  - `GET /admin/api/decor/pc?variant=<variant_key>`
  - `PUT /admin/api/decor/pc?variant=<variant_key>`
  - `POST /admin/api/decor/pc/publish?variant=<variant_key>`
  - `POST /admin/api/decor/pc/copies`
  - `PUT /admin/api/decor/pc/variants/:variant_key`
  - `DELETE /admin/api/decor/pc/variants/:variant_key`
```

- [ ] **Step 5: 运行 docs-site 构建**

Run: `npm run build`（工作目录 `docs-site`）

Expected: PASS。

- [ ] **Step 6: Commit**

```bash
git add docs-site/docs/api/decor.md
git commit -m "更新PC装修多版本接口文档" -m "装修接口文档改为任意页面多副本语义，并补充PC首页装修版本管理接口与发布规则。"
```

---

### Task 5: 回归验证

**Files:**
- Modify: none

- [ ] **Step 1: 运行后端 Decor 服务测试**

Run: `go test ./plugins/decor/service -v`（工作目录 `server`）

Expected: PASS。

- [ ] **Step 2: 运行 Admin 普通构建**

Run: `npm run build`（工作目录 `admin`）

Expected: PASS。

- [ ] **Step 3: 运行 Admin Mock 构建**

Run: `npm run build:demo`（工作目录 `admin`）

Expected: PASS。

- [ ] **Step 4: 运行文档站构建**

Run: `npm run build`（工作目录 `docs-site`）

Expected: PASS。

- [ ] **Step 5: 查看最终变更**

Run: `git status --short`

Expected: clean working tree。

---

## Self-Review

### Spec coverage

- PC 后台版本下拉、复制、重命名、删除、保存、发布：Task 2。
- 后端复用 `page_key=pc` 通用多副本能力：Task 1。
- Admin Mock 与真实接口一致：Task 3。
- docs-site 同步更新功能说明、接口变化、部署影响：Task 4。
- 验证命令覆盖服务层、Admin、Mock、文档站：Task 5。

### Placeholder scan

计划不包含 TBD、TODO、待定或未定义步骤。每个代码变更步骤包含目标代码或明确替换内容。

### Type consistency

计划统一使用 `variant_key`、`variant_name`、`is_published`、`pagePayload`、`PcDecorPagePayload`、`page_key=pc`，与现有服务层和编辑器命名一致。
