<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { createSpecTemplate, updateSpecTemplate } from '@/api/spec-template'
import CategoryTreePicker from '@/components/biz/CategoryTreePicker.vue'

const props = defineProps<{ show: boolean; template?: any }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'saved'): void }>()

const showCatPicker = ref(false)

const form = reactive<{
  name: string
  category_ids: number[]
  attrs: Array<{ name: string; valuesStr: string }>
  status: number
}>({
  name: '',
  category_ids: [],
  attrs: [],
  status: 1,
})

watch(() => props.template, (t) => {
  if (t) {
    form.name = String(t.name || '')
    form.category_ids = Array.isArray(t.category_ids) ? [...t.category_ids] : []
    form.attrs = Array.isArray(t.attrs)
      ? t.attrs.map((a: any) => ({ name: String(a.name || ''), valuesStr: Array.isArray(a.values) ? a.values.join(',') : '' }))
      : []
    form.status = Number(t.status || 0) === 1 ? 1 : 0
  } else {
    form.name = ''; form.category_ids = []; form.attrs = []; form.status = 1
  }
}, { immediate: true })

function addAttr() {
  form.attrs.push({ name: '', valuesStr: '' })
}

function removeAttr(idx: number) {
  form.attrs.splice(idx, 1)
}

function onPickCategory(payload: { id: number; path_name: string }) {
  if (!form.category_ids.includes(payload.id)) {
    form.category_ids.push(payload.id)
  }
  showCatPicker.value = false
}

function removeCategoryId(idx: number) {
  form.category_ids.splice(idx, 1)
}

async function save() {
  if (!form.name.trim()) { uni.showToast({ title: '请输入模板名称', icon: 'none' }); return }
  const payload = {
    name: form.name.trim(),
    category_ids: form.category_ids,
    attrs: form.attrs.map((a) => ({
      name: a.name.trim(),
      values: a.valuesStr.split(',').map((v) => v.trim()).filter(Boolean),
    })).filter((a) => a.name && a.values.length),
    status: form.status,
  }
  if (props.template?.id) {
    await updateSpecTemplate(props.template.id, payload)
  } else {
    await createSpecTemplate(payload)
  }
  emit('saved')
}
</script>

<template>
  <up-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="popup-body">
      <view class="popup-title">{{ template ? '编辑模板' : '新建模板' }}</view>
      <up-input v-model="form.name" placeholder="模板名称" clearable />
      <view class="mt" />
      <view class="label">适用分类</view>
      <view class="cat-row">
        <view v-for="(cid, idx) in form.category_ids" :key="cid" class="cat-tag">
          #{{ cid }}<text class="x" @click="removeCategoryId(idx)">✕</text>
        </view>
        <up-button size="mini" plain @click="showCatPicker = true">+ 添加分类</up-button>
      </view>
      <view class="mt" />
      <view class="label">属性组</view>
      <view v-for="(attr, idx) in form.attrs" :key="idx" class="attr-row">
        <up-input v-model="attr.name" placeholder="属性名" class="attr-name" />
        <up-input v-model="attr.valuesStr" placeholder="值（逗号分隔）" class="attr-vals" />
        <text class="x" @click="removeAttr(idx)">✕</text>
      </view>
      <up-button size="mini" plain @click="addAttr" class="mt">+ 添加属性组</up-button>
      <view class="mt" />
      <view class="row">
        <text>启用</text>
        <switch :checked="form.status === 1" @change="(e: any) => form.status = e.detail.value ? 1 : 0" />
      </view>
      <view class="mt-lg" />
      <up-button type="primary" @click="save">保存</up-button>

      <CategoryTreePicker :show="showCatPicker" :value="0" @close="showCatPicker = false" @pick="onPickCategory" />
    </view>
  </up-popup>
</template>

<style scoped>
.popup-body { padding: 24rpx; box-sizing: border-box; max-height: 80vh; overflow-y: auto; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.label { font-size: 24rpx; color: var(--eapp-text-muted); margin-bottom: 8rpx; }
.cat-row { display: flex; gap: 8rpx; flex-wrap: wrap; align-items: center; }
.cat-tag { background: var(--eapp-bg); padding: 6rpx 14rpx; border-radius: 999rpx; font-size: 22rpx; display: inline-flex; align-items: center; gap: 6rpx; }
.attr-row { display: flex; gap: 8rpx; align-items: center; margin-top: 8rpx; }
.attr-name { width: 200rpx; }
.attr-vals { flex: 1; }
.x { color: var(--eapp-danger); font-size: 22rpx; padding: 0 8rpx; }
.row { display: flex; align-items: center; justify-content: space-between; font-size: 26rpx; }
.mt { margin-top: 12rpx; }
.mt-lg { margin-top: 16rpx; }
</style>
