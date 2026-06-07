<script setup lang="ts">
import { ref, watch, computed } from 'vue'

type SkuRow = { id?: number; attrs: Array<{ name: string; value: string }>; price: number; stock: number }
type SpecTemplateAttr = { name: string; values: string[] }
type SpecTemplate = { id: number; name: string; attrs: SpecTemplateAttr[] }

const props = defineProps<{
  skus: SkuRow[]
  basePrice?: number
  specTemplate?: SpecTemplate | null
  wmsStockMap?: Record<number, { qty: number; reserved_qty: number }>
}>()
const emit = defineEmits<{ (e: 'update', skus: SkuRow[]): void }>()

// ── 自由模式状态 ──────────────────────────────────────────────
const specs = ref<Array<{ name: string; values: string[] }>>([])
const matrix = ref<SkuRow[]>([])
const bulkPrice = ref('')
const bulkStock = ref('')

// ── 模板模式状态 ──────────────────────────────────────────────
const selectedValues = ref<Record<string, Set<string>>>({})

const isTemplateMode = computed(() => !!props.specTemplate)

// 当模板 prop 变化时，重置选中集合
watch(() => props.specTemplate, (tpl) => {
  if (!tpl) {
    rebuildFromSkus()
    return
  }
  const map: Record<string, Set<string>> = {}
  for (const group of tpl.attrs) {
    map[group.name] = new Set()
  }
  // 从已有 skus 反推勾选值
  for (const sku of props.skus || []) {
    for (const attr of sku.attrs || []) {
      if (map[attr.name] !== undefined) map[attr.name].add(attr.value)
    }
  }
  selectedValues.value = map
  rebuildTemplateMatrix()
}, { immediate: false })

watch(() => props.skus, () => {
  if (!isTemplateMode.value) rebuildFromSkus()
}, { immediate: true })

// ── 自由模式逻辑 ──────────────────────────────────────────────
function rebuildFromSkus() {
  const groups: Record<string, Set<string>> = {}
  for (const sku of props.skus || []) {
    for (const a of sku.attrs || []) {
      groups[a.name] = groups[a.name] || new Set()
      groups[a.name].add(a.value)
    }
  }
  specs.value = Object.keys(groups).map((name) => ({ name, values: Array.from(groups[name]) }))
  matrix.value = (props.skus || []).slice()
}

function rebuildMatrix() {
  if (!specs.value.length) { matrix.value = []; emit('update', []); return }
  function cross(idx: number): Array<Array<{ name: string; value: string }>> {
    if (idx >= specs.value.length) return [[]]
    const sub = cross(idx + 1)
    return specs.value[idx].values.flatMap((v) => sub.map((row) => [{ name: specs.value[idx].name, value: v }, ...row]))
  }
  const combos = cross(0)
  const next = combos.map((attrs) => {
    const key = attrs.map((a) => `${a.name}:${a.value}`).join('|')
    const exist = matrix.value.find((row) => row.attrs.map((a) => `${a.name}:${a.value}`).join('|') === key)
    return exist || { attrs, price: Number(props.basePrice || 0), stock: 0 }
  })
  matrix.value = next
  emit('update', matrix.value)
}

function addSpec() {
  specs.value.push({ name: '规格' + (specs.value.length + 1), values: ['默认'] })
  rebuildMatrix()
}
function removeSpec(idx: number) { specs.value.splice(idx, 1); rebuildMatrix() }
function addValue(idx: number) { specs.value[idx].values.push('值' + (specs.value[idx].values.length + 1)); rebuildMatrix() }
function removeValue(idx: number, vIdx: number) { specs.value[idx].values.splice(vIdx, 1); rebuildMatrix() }

// ── 模板模式逻辑 ──────────────────────────────────────────────
function isSelected(attrName: string, value: string) {
  return selectedValues.value[attrName]?.has(value) ?? false
}

function toggleValue(attrName: string, value: string) {
  const map = selectedValues.value
  if (!map[attrName]) map[attrName] = new Set()
  if (map[attrName].has(value)) {
    map[attrName].delete(value)
  } else {
    map[attrName].add(value)
  }
  rebuildTemplateMatrix()
}

function cartesian(groups: string[][]): string[][] {
  return groups.reduce<string[][]>((acc, values) => {
    if (!acc.length) return values.map((v) => [v])
    return acc.flatMap((combo) => values.map((v) => [...combo, v]))
  }, [])
}

function rebuildTemplateMatrix() {
  const tpl = props.specTemplate
  if (!tpl) return
  const activeGroups: { name: string; values: string[] }[] = []
  for (const group of tpl.attrs) {
    const chosen = Array.from(selectedValues.value[group.name] ?? [])
    if (chosen.length) activeGroups.push({ name: group.name, values: chosen })
  }
  if (!activeGroups.length) {
    matrix.value = []
    emit('update', [])
    return
  }
  const combos = cartesian(activeGroups.map((g) => g.values))
  const existing = new Map(matrix.value.map((s) => {
    const key = s.attrs.map((a) => `${a.name}:${a.value}`).sort().join('|')
    return [key, s]
  }))
  matrix.value = combos.map((combo) => {
    const attrs = combo.map((val, i) => ({ name: activeGroups[i].name, value: val }))
    const key = attrs.map((a) => `${a.name}:${a.value}`).sort().join('|')
    const prev = existing.get(key)
    return prev
      ? { ...prev, attrs }
      : { attrs, price: Number(props.basePrice || 0), stock: 0 }
  })
  emit('update', matrix.value)
}

// ── WMS 可售量展示 ────────────────────────────────────────────
function wmsSellable(skuId?: number): string {
  if (!skuId || skuId <= 0) return '保存后初始化'
  const s = props.wmsStockMap?.[skuId]
  if (!s) return '—'
  return String(s.qty - s.reserved_qty)
}

// ── 共用批量操作 ──────────────────────────────────────────────
function applyBulkPrice() {
  if (!bulkPrice.value) return
  for (const row of matrix.value) row.price = Number(bulkPrice.value)
  emit('update', matrix.value)
}
function applyBulkStock() {
  if (!bulkStock.value) return
  for (const row of matrix.value) row.stock = Number(bulkStock.value)
  emit('update', matrix.value)
}
function onCellChange() { emit('update', matrix.value) }
</script>

<template>
  <view class="sku-editor">
    <!-- 模板模式 -->
    <template v-if="isTemplateMode && specTemplate">
      <view v-for="group in specTemplate.attrs" :key="group.name" class="spec">
        <view class="spec-head-label">{{ group.name }}</view>
        <view class="vals">
          <view
            v-for="val in group.values"
            :key="val"
            :class="['chip', isSelected(group.name, val) ? 'chip-on' : '']"
            @click="toggleValue(group.name, val)"
          >{{ val }}</view>
        </view>
      </view>
    </template>

    <!-- 自由模式 -->
    <template v-else>
      <view v-for="(spec, idx) in specs" :key="idx" class="spec">
        <view class="spec-head">
          <up-input v-model="spec.name" class="name-in" />
          <up-button size="mini" type="error" plain @click="removeSpec(idx)">删除规格</up-button>
        </view>
        <view class="vals">
          <view v-for="(v, vIdx) in spec.values" :key="vIdx" class="val">
            <up-input v-model="spec.values[vIdx]" @blur="rebuildMatrix" />
            <text class="x" @click="removeValue(idx, vIdx)">✕</text>
          </view>
          <up-button size="mini" plain @click="addValue(idx)">+ 值</up-button>
        </view>
      </view>
      <up-button size="mini" type="primary" plain @click="addSpec">+ 规格组</up-button>
    </template>

    <!-- 批量操作（两种模式共用） -->
    <view v-if="matrix.length" class="bulk">
      <up-input v-model="bulkPrice" type="digit" placeholder="批量赋价" class="bulk-in" />
      <up-button size="mini" @click="applyBulkPrice">应用</up-button>
      <up-input v-model="bulkStock" type="number" placeholder="批量赋库存" class="bulk-in" />
      <up-button size="mini" @click="applyBulkStock">应用</up-button>
    </view>

    <!-- SKU 矩阵表（两种模式共用） -->
    <view v-if="matrix.length" class="matrix">
      <view class="row head">
        <view class="cell flex">规格组合</view>
        <view class="cell">价格</view>
        <view class="cell">库存</view>
        <view class="cell wms-col">WMS可售</view>
      </view>
      <view v-for="(row, i) in matrix" :key="i" class="row">
        <view class="cell flex">{{ row.attrs.map((a) => a.value).join(' / ') }}</view>
        <up-input v-model="row.price" type="digit" class="cell" @blur="onCellChange" />
        <up-input v-model="row.stock" type="number" class="cell" @blur="onCellChange" />
        <view class="cell wms-col">
          <text :class="row.id && row.id > 0 ? 'wms-val' : 'wms-new'">{{ wmsSellable(row.id) }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.sku-editor { display: grid; gap: 14rpx; }
.spec { background: var(--eapp-bg); padding: 14rpx; border-radius: 14rpx; }
.spec-head { display: flex; gap: 12rpx; align-items: center; }
.spec-head-label { font-size: 26rpx; font-weight: 600; color: var(--eapp-text); margin-bottom: 10rpx; }
.name-in { flex: 1; }
.vals { margin-top: 10rpx; display: flex; gap: 8rpx; flex-wrap: wrap; align-items: center; }
.val { display: flex; align-items: center; gap: 6rpx; background: var(--eapp-card); padding: 4rpx 8rpx; border-radius: 10rpx; }
.x { color: var(--eapp-danger); font-size: 22rpx; padding: 0 6rpx; }
.chip { padding: 8rpx 20rpx; border-radius: 999rpx; font-size: 24rpx; border: 2rpx solid var(--eapp-border); color: var(--eapp-text-muted); background: var(--eapp-card); }
.chip-on { border-color: var(--eapp-primary, #3b82f6); color: var(--eapp-primary, #3b82f6); background: color-mix(in srgb, var(--eapp-primary, #3b82f6) 10%, transparent); }
.bulk { display: flex; gap: 8rpx; align-items: center; flex-wrap: wrap; }
.bulk-in { width: 200rpx; }
.matrix { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 16rpx; padding: 10rpx; }
.row { display: flex; gap: 8rpx; align-items: center; padding: 8rpx 4rpx; border-bottom: 1px dashed var(--eapp-border); }
.row.head { font-weight: 600; color: var(--eapp-text-muted); }
.cell { width: 160rpx; }
.cell.flex { flex: 1; }
.wms-col { width: 130rpx; font-size: 22rpx; }
.wms-val { color: var(--eapp-primary, #3b82f6); font-weight: 600; }
.wms-new { color: var(--eapp-text-muted); font-style: italic; }
</style>
