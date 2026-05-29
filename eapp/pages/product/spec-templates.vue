<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useSpecTemplateList } from '@/composables/useSpecTemplateList'
import { deleteSpecTemplate } from '@/api/spec-template'
import SpecTemplateForm from '@/components/product/SpecTemplateForm.vue'

const h = useSpecTemplateList()
const showForm = ref(false)
const editingTemplate = ref<any>(null)
const keyword = ref('')

function openCreate() { editingTemplate.value = null; showForm.value = true }
function openEdit(item: any) { editingTemplate.value = item; showForm.value = true }

async function onDelete(id: number) {
  await deleteSpecTemplate(id)
  await h.load()
}

async function onSaved() {
  showForm.value = false
  await h.load()
}

function search() {
  h.load({ keyword: keyword.value || undefined })
}

function categoryNames(ids: number[]) {
  if (!Array.isArray(ids) || !ids.length) return '-'
  return ids.map((id) => `#${id}`).join(', ')
}

function attrSummary(attrs: any[]) {
  if (!Array.isArray(attrs) || !attrs.length) return '-'
  return attrs.map((a: any) => `${a.name}(${Array.isArray(a.values) ? a.values.length : 0})`).join(', ')
}

onShow(() => h.load())
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <up-input v-model="keyword" placeholder="搜索模板" clearable class="search" @confirm="search" />
      <up-button size="mini" type="primary" @click="openCreate">新建</up-button>
    </view>
    <view v-if="!h.loading.value && !h.list.value.length" class="empty">暂无规格模板</view>
    <view v-for="item in h.list.value" :key="item.id" class="card">
      <view class="row">
        <text class="name">{{ item.name }}</text>
        <text :class="['badge', item.status === 1 ? 'on' : 'off']">{{ item.status === 1 ? '启用' : '禁用' }}</text>
      </view>
      <view class="desc">分类：{{ categoryNames(item.category_ids) }}</view>
      <view class="desc">属性组：{{ attrSummary(item.attrs) }}</view>
      <view class="actions">
        <up-button size="mini" type="primary" plain @click="openEdit(item)">编辑</up-button>
        <up-button size="mini" type="error" plain @click="onDelete(item.id)">删除</up-button>
      </view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中…</view>

    <SpecTemplateForm :show="showForm" :template="editingTemplate" @close="showForm = false" @saved="onSaved" />
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; display: grid; gap: 14rpx; align-content: start; }
.top-bar { display: flex; gap: 12rpx; align-items: center; }
.search { flex: 1; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.row { display: flex; align-items: center; justify-content: space-between; }
.name { font-size: 28rpx; font-weight: 600; }
.badge { font-size: 22rpx; padding: 4rpx 14rpx; border-radius: 999rpx; }
.badge.on { background: #dcfce7; color: #16a34a; }
.badge.off { background: #fee2e2; color: #dc2626; }
.desc { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 23rpx; }
.actions { margin-top: 12rpx; display: flex; gap: 10rpx; justify-content: flex-end; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>
