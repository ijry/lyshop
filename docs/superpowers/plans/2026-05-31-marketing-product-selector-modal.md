# Marketing Product Selector Modal Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the product `<select>` dropdown in `ActivityProductManageBase.vue` with a modal dialog that supports keyword search, pagination, and displays product cover, title, price, stock, favorites, and status.

**Architecture:** Create a standalone `ProductSelectorModal.vue` component under `admin/src/components/common/` that fetches products via the existing `getProducts` API, renders a searchable paginated list, and emits the selected product on row click. Update `ActivityProductManageBase.vue` to open this modal instead of the `<select>` dropdown. Add two new i18n keys to both locale files.

**Tech Stack:** Vue 3 Composition API, TypeScript, Tailwind CSS, vue-i18n, existing `getProducts` / `getCategories` from `@/api/plugins`

---

## File Map

| Action | File | Purpose |
|--------|------|---------|
| Create | `admin/src/components/common/ProductSelectorModal.vue` | Reusable modal: search, pagination, product list |
| Modify | `admin/src/views/marketing/ActivityProductManageBase.vue` | Replace select dropdown with modal trigger |
| Modify | `admin/src/locales/zh-CN.ts` | Add 2 new i18n keys |
| Modify | `admin/src/locales/en.ts` | Add 2 new i18n keys |

---

## Task 1: Add i18n keys

**Files:**
- Modify: `admin/src/locales/zh-CN.ts:506-512`
- Modify: `admin/src/locales/en.ts:506-512`

- [ ] **Step 1: Add keys to zh-CN.ts**

In `zh-CN.ts`, after line `'marketingProduct.selectProduct': '选择商品',` add:

```ts
  'marketingProduct.selectProductBtn': '选择商品',
  'marketingProduct.changeProduct': '更换商品',
```

- [ ] **Step 2: Add keys to en.ts**

In `en.ts`, after line `'marketingProduct.selectProduct': 'Select Product',` add:

```ts
  'marketingProduct.selectProductBtn': 'Select Product',
  'marketingProduct.changeProduct': 'Change Product',
```

- [ ] **Step 3: Commit**

```bash
git add admin/src/locales/zh-CN.ts admin/src/locales/en.ts
git commit -m "i18n: add product selector modal keys"
```

---

## Task 2: Create ProductSelectorModal.vue

**Files:**
- Create: `admin/src/components/common/ProductSelectorModal.vue`

- [ ] **Step 1: Create the component**

```vue
<template>
  <div v-if="visible" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50" @click.self="close">
    <div class="bg-white rounded-xl shadow-2xl w-[760px] max-w-[96vw] max-h-[88vh] flex flex-col">
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-100">
        <h3 class="text-base font-semibold text-slate-800">{{ $t('marketingProduct.selectProduct') }}</h3>
        <button @click="close" class="text-slate-400 hover:text-slate-600 text-lg leading-none">✕</button>
      </div>

      <!-- Search bar -->
      <div class="px-6 py-3 border-b border-slate-100 flex gap-3">
        <input
          v-model="keyword"
          :placeholder="$t('common.search')"
          class="border border-slate-200 rounded-lg px-3 py-2 text-sm flex-1 focus:outline-none focus:border-blue-400"
          @keyup.enter="search"
        />
        <select
          v-model="categoryId"
          class="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none"
          @change="search"
        >
          <option value="">{{ $t('product.list.allCategory') }}</option>
          <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>
        <button @click="search" class="px-4 py-2 bg-slate-100 rounded-lg text-sm hover:bg-slate-200">
          {{ $t('common.search') }}
        </button>
      </div>

      <!-- Product list -->
      <div class="flex-1 overflow-y-auto">
        <table class="w-full text-sm">
          <thead class="bg-slate-50 border-b border-slate-100 sticky top-0">
            <tr>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.id') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('product.list.product') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('product.list.price') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('product.list.stock') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('product.list.favorites') }}</th>
              <th class="px-4 py-3 text-left text-slate-500 font-medium">{{ $t('common.status') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-50">
            <tr
              v-for="p in products"
              :key="p.id"
              class="hover:bg-blue-50 cursor-pointer"
              @click="select(p)"
            >
              <td class="px-4 py-3 text-slate-400">{{ p.id }}</td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <img v-if="p.cover" :src="p.cover" class="w-10 h-10 rounded-lg object-cover flex-shrink-0" />
                  <div v-else class="w-10 h-10 rounded-lg bg-slate-100 flex items-center justify-center text-slate-400 text-xs flex-shrink-0">
                    {{ $t('product.list.image') }}
                  </div>
                  <span class="font-medium text-slate-700 truncate max-w-[200px]">{{ p.title }}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-slate-700">¥{{ p.price }}</td>
              <td class="px-4 py-3 text-slate-700">{{ p.stock }}</td>
              <td class="px-4 py-3 text-slate-700">{{ p.favorite_count || 0 }}</td>
              <td class="px-4 py-3">
                <span
                  :class="p.status === 1 ? 'text-green-600 bg-green-50' : 'text-slate-400 bg-slate-100'"
                  class="px-2 py-1 rounded-full text-xs"
                >
                  {{ p.status === 1 ? $t('product.list.onSale') : $t('product.list.offSale') }}
                </span>
              </td>
            </tr>
            <tr v-if="!loading && !products.length">
              <td colspan="6" class="px-4 py-12 text-center text-slate-400">{{ $t('common.noData') }}</td>
            </tr>
            <tr v-if="loading">
              <td colspan="6" class="px-4 py-12 text-center text-slate-400">{{ $t('common.loading') }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="px-6 py-3 border-t border-slate-100 flex items-center justify-between text-sm text-slate-500">
        <span>{{ $t('common.totalCount', { total }) }}</span>
        <div class="flex gap-2">
          <button
            :disabled="page <= 1"
            @click="page--; load()"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40"
          >{{ $t('common.prevPage') }}</button>
          <button
            :disabled="page * pageSize >= total"
            @click="page++; load()"
            class="px-3 py-1 rounded-lg border hover:bg-slate-50 disabled:opacity-40"
          >{{ $t('common.nextPage') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { getProducts, getCategories } from '@/api/plugins'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'select', product: any): void
}>()

const keyword = ref('')
const categoryId = ref('')
const page = ref(1)
const pageSize = 10
const total = ref(0)
const products = ref<any[]>([])
const categories = ref<any[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const data: any = await getProducts({
      keyword: keyword.value || undefined,
      category_id: categoryId.value || undefined,
      page: page.value,
      size: pageSize,
    })
    products.value = data?.list || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  load()
}

function select(product: any) {
  emit('select', product)
  emit('close')
}

function close() {
  emit('close')
}

watch(() => props.visible, (val) => {
  if (val) {
    keyword.value = ''
    categoryId.value = ''
    page.value = 1
    if (!categories.value.length) {
      getCategories().then((d: any) => { categories.value = d || [] })
    }
    load()
  }
})
</script>
```

- [ ] **Step 2: Commit**

```bash
git add admin/src/components/common/ProductSelectorModal.vue
git commit -m "feat: add ProductSelectorModal component with search and pagination"
```

---

## Task 3: Update ActivityProductManageBase.vue

**Files:**
- Modify: `admin/src/views/marketing/ActivityProductManageBase.vue`

Replace the product `<select>` in the form panel (lines 77–87) with a button that opens the modal, and replace the `loadProductOptions` / `productOptions` approach with the modal-based selection.

- [ ] **Step 1: Update the template — replace product select with modal trigger**

In the `<template>` section, replace the product select block (lines 76–87):

```html
<!-- OLD -->
<div>
  <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('marketingProduct.selectProduct') }}</label>
  <select
    v-model.number="selectedProductID"
    class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm"
    @change="onSelectProduct"
    :disabled="editingID > 0"
  >
    <option :value="0">{{ $t('marketingProduct.selectProduct') }}</option>
    <option v-for="item in productOptions" :key="item.id" :value="item.id">{{ item.title }} (ID: {{ item.id }})</option>
  </select>
</div>
```

With:

```html
<div>
  <label class="block text-sm font-medium text-slate-700 mb-1">{{ $t('marketingProduct.selectProduct') }}</label>
  <div v-if="selectedProductID > 0" class="flex items-center gap-3 p-3 border border-slate-200 rounded-xl mb-2">
    <img v-if="selectedProductCover" :src="selectedProductCover" class="w-10 h-10 rounded-lg object-cover flex-shrink-0" />
    <div v-else class="w-10 h-10 rounded-lg bg-slate-100 flex-shrink-0" />
    <span class="text-sm text-slate-700 flex-1 truncate">{{ selectedProductTitle }}</span>
  </div>
  <button
    type="button"
    :disabled="editingID > 0"
    @click="showProductSelector = true"
    class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm text-left text-slate-500 hover:border-blue-400 hover:text-blue-600 transition disabled:opacity-50 disabled:cursor-not-allowed"
  >
    {{ selectedProductID > 0 ? $t('marketingProduct.changeProduct') : $t('marketingProduct.selectProductBtn') }}
  </button>
</div>
```

- [ ] **Step 2: Add ProductSelectorModal to the template**

At the end of the template, just before the closing `</div>` of the outer wrapper, add:

```html
<ProductSelectorModal
  :visible="showProductSelector"
  @close="showProductSelector = false"
  @select="onProductSelected"
/>
```

- [ ] **Step 3: Update the script — imports and new refs**

Add to the imports at the top of `<script setup>`:

```ts
import ProductSelectorModal from '@/components/common/ProductSelectorModal.vue'
```

Add new refs after the existing `const skuOptions = ref<SkuOption[]>([])` line:

```ts
const showProductSelector = ref(false)
const selectedProductTitle = ref('')
const selectedProductCover = ref('')
```

- [ ] **Step 4: Add onProductSelected handler**

Replace the existing `onSelectProduct` function:

```ts
// OLD
async function onSelectProduct() {
  form.product_id = Number(selectedProductID.value || 0)
  form.sku_id = 0
  await loadSkuOptions(form.product_id)
}
```

With:

```ts
async function onProductSelected(product: any) {
  selectedProductID.value = Number(product?.id || 0)
  selectedProductTitle.value = String(product?.title || '')
  selectedProductCover.value = String(product?.cover || '')
  form.product_id = selectedProductID.value
  form.sku_id = 0
  await loadSkuOptions(form.product_id)
}
```

- [ ] **Step 5: Update resetForm to clear new refs**

In `resetForm()`, add after `selectedProductID.value = 0`:

```ts
selectedProductTitle.value = ''
selectedProductCover.value = ''
```

- [ ] **Step 6: Update openEdit to populate new refs**

In `openEdit()`, after `selectedProductID.value = form.product_id`, add:

```ts
selectedProductTitle.value = String(row?.product_title || '')
selectedProductCover.value = String(row?.product_cover || '')
```

- [ ] **Step 7: Remove loadProductOptions and its call**

Delete the `loadProductOptions` function entirely:

```ts
// DELETE this entire function:
async function loadProductOptions() {
  const data: any = await request.get('/products', { params: { page: 1, size: 200 } })
  productOptions.value = data?.list || []
}
```

Delete `const productOptions = ref<any[]>([])`.

Remove the `loadProductOptions()` call from `onMounted`.

Also remove the filter product `<select>` in the search bar (lines 22–25) since it relied on `productOptions` — replace it with a plain keyword search (the filter-by-product dropdown is no longer needed as the modal handles product selection):

```html
<!-- REMOVE this block from the search bar -->
<select v-model.number="filterProductID" class="border border-slate-200 rounded-lg px-3 py-2 text-sm">
  <option :value="0">{{ $t('marketingProduct.filterProduct') }}</option>
  <option v-for="item in productOptions" :key="item.id" :value="item.id">{{ item.title }} (ID: {{ item.id }})</option>
</select>
```

Also remove `const filterProductID = ref(0)` and the `product_id: filterProductID.value || undefined` param from `loadProducts`.

- [ ] **Step 8: Commit**

```bash
git add admin/src/views/marketing/ActivityProductManageBase.vue
git commit -m "feat: replace product select dropdown with ProductSelectorModal in marketing activity form"
```

---

## Self-Review

**Spec coverage:**
- [x] 弹出框替换下拉框 → Task 3 replaces `<select>` with modal trigger
- [x] 分页 → Task 2 implements pagination in modal
- [x] 关键字搜索 → Task 2 implements keyword + category search
- [x] 显示封面、标题、价格、库存、收藏数、状态 → Task 2 table columns
- [x] 复用商品列表组件 → `ProductSelectorModal` is standalone reusable component
- [x] 单选模式 → click row emits select and closes modal

**Placeholder scan:** None found.

**Type consistency:** `product: any` used consistently across Task 2 and Task 3. `selectedProductID`, `selectedProductTitle`, `selectedProductCover` defined in Task 3 Step 3 and used in Steps 4, 5, 6.
