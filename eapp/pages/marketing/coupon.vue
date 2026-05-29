<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useCouponList } from '@/composables/useCouponList'
import { createCoupon, updateCoupon, deleteCoupon, sendCoupon } from '@/api/marketing'
import CouponForm from '@/components/marketing/CouponForm.vue'
import CouponSendPopup from '@/components/marketing/CouponSendPopup.vue'

const h = useCouponList()
const showForm = ref(false)
const showSend = ref(false)
const editingCoupon = ref<any>(null)
const activeCoupon = ref<any>(null)
const keyword = ref('')

const stackLabels: Record<string, string> = { exclusive: '互斥', same_type: '同类可叠', cross_type: '跨类可叠' }
const typeLabels: Record<number, string> = { 1: '满减', 2: '折扣', 3: '立减' }

function openCreate() { editingCoupon.value = null; showForm.value = true }
function openEdit(item: any) { editingCoupon.value = item; showForm.value = true }
function openSend(item: any) { activeCoupon.value = item; showSend.value = true }

async function onSubmitForm(payload: any) {
  if (editingCoupon.value) {
    await updateCoupon(editingCoupon.value.id, payload)
  } else {
    await createCoupon(payload)
  }
  showForm.value = false
  await h.load()
}

async function onDelete(id: number) {
  await deleteCoupon(id)
  await h.load()
}

async function onSend(payload: { count: number }) {
  if (!activeCoupon.value) return
  await sendCoupon(activeCoupon.value.id, payload)
  uni.showToast({ title: `已发送 ${payload.count} 张`, icon: 'success' })
  showSend.value = false
  await h.load()
}

function search() {
  h.applyFilter({ keyword: keyword.value || undefined })
  h.load()
}

onShow(() => h.load())
</script>

<template>
  <view class="page">
    <view class="top-bar">
      <up-input v-model="keyword" placeholder="搜索优惠券" clearable class="search" @confirm="search" />
      <up-button size="mini" type="primary" @click="openCreate">新建</up-button>
    </view>
    <view v-if="!h.loading.value && !h.list.value.length" class="empty">暂无优惠券</view>
    <view v-for="item in h.list.value" :key="item.id" class="card">
      <view class="row">
        <text class="name">{{ item.name }}</text>
        <text :class="['badge', item.status === 1 ? 'on' : 'off']">{{ item.status === 1 ? '启用' : '禁用' }}</text>
      </view>
      <view class="desc">
        类型：{{ typeLabels[item.type] || '-' }} ·
        面额：{{ item.discount }} ·
        已用：{{ item.used_count || 0 }}
      </view>
      <view class="desc" v-if="item.stack_rule">
        叠加：<text class="stack-badge">{{ stackLabels[item.stack_rule] || item.stack_rule }}</text>
      </view>
      <view class="actions">
        <up-button size="mini" type="primary" plain @click="openEdit(item)">编辑</up-button>
        <up-button size="mini" type="warning" plain @click="openSend(item)">发券</up-button>
        <up-button size="mini" type="error" plain @click="onDelete(item.id)">删除</up-button>
      </view>
    </view>
    <view v-if="h.loading.value" class="loading">加载中…</view>

    <CouponForm :show="showForm" :coupon="editingCoupon" @close="showForm = false" @submit="onSubmitForm" />
    <CouponSendPopup :show="showSend" :coupon="activeCoupon" @close="showSend = false" @send="onSend" />
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
.stack-badge { display: inline-block; background: #eff6ff; color: #2563eb; padding: 2rpx 10rpx; border-radius: 6rpx; font-size: 22rpx; }
.actions { margin-top: 12rpx; display: flex; gap: 10rpx; justify-content: flex-end; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.loading { text-align: center; color: var(--eapp-text-muted); padding: 20rpx 0; }
</style>
