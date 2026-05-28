<script setup lang="ts">
import { ref, watch } from 'vue'
const props = defineProps<{ skus: Array<{ id?: number; attrs: Array<{ name: string; value: string }>; price: number; stock: number }>; basePrice?: number }>()
const emit = defineEmits<{ (e: 'update', skus: any[]): void }>()

const specs = ref<Array<{ name: string; values: string[] }>>([])
const matrix = ref<Array<{ id?: number; attrs: Array<{ name: string; value: string }>; price: number; stock: number }>>([])
const bulkPrice = ref('')
const bulkStock = ref('')

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
  if (!specs.value.length) { matrix.value = []; return }
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
function removeSpec(idx: number) {
  specs.value.splice(idx, 1)
  rebuildMatrix()
}
function addValue(idx: number) {
  specs.value[idx].values.push('值' + (specs.value[idx].values.length + 1))
  rebuildMatrix()
}
function removeValue(idx: number, vIdx: number) {
  specs.value[idx].values.splice(vIdx, 1)
  rebuildMatrix()
}
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

watch(() => props.skus, () => rebuildFromSkus(), { immediate: true })
</script>

<template>
  <view class="sku-editor">
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

    <view v-if="matrix.length" class="bulk">
      <up-input v-model="bulkPrice" type="digit" placeholder="批量赋价" class="bulk-in" />
      <up-button size="mini" @click="applyBulkPrice">应用</up-button>
      <up-input v-model="bulkStock" type="number" placeholder="批量赋库存" class="bulk-in" />
      <up-button size="mini" @click="applyBulkStock">应用</up-button>
    </view>

    <view v-if="matrix.length" class="matrix">
      <view class="row head">
        <view class="cell flex">规格组合</view>
        <view class="cell">价格</view>
        <view class="cell">库存</view>
      </view>
      <view v-for="(row, i) in matrix" :key="i" class="row">
        <view class="cell flex">{{ row.attrs.map((a) => `${a.value}`).join(' / ') }}</view>
        <up-input v-model="row.price" type="digit" class="cell" @blur="onCellChange" />
        <up-input v-model="row.stock" type="number" class="cell" @blur="onCellChange" />
      </view>
    </view>
  </view>
</template>

<style scoped>
.sku-editor { display: grid; gap: 14rpx; }
.spec { background: var(--eapp-bg); padding: 14rpx; border-radius: 14rpx; }
.spec-head { display: flex; gap: 12rpx; align-items: center; }
.name-in { flex: 1; }
.vals { margin-top: 10rpx; display: flex; gap: 8rpx; flex-wrap: wrap; align-items: center; }
.val { display: flex; align-items: center; gap: 6rpx; background: var(--eapp-card); padding: 4rpx 8rpx; border-radius: 10rpx; }
.x { color: var(--eapp-danger); font-size: 22rpx; padding: 0 6rpx; }
.bulk { display: flex; gap: 8rpx; align-items: center; flex-wrap: wrap; }
.bulk-in { width: 200rpx; }
.matrix { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 16rpx; padding: 10rpx; }
.row { display: flex; gap: 8rpx; align-items: center; padding: 8rpx 4rpx; border-bottom: 1px dashed var(--eapp-border); }
.row.head { font-weight: 600; color: var(--eapp-text-muted); }
.cell { width: 160rpx; }
.cell.flex { flex: 1; }
</style>
