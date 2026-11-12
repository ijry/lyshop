<script setup lang="ts">
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import PageHeader from '@/components/biz/PageHeader.vue'
import SectionHeader from '@/components/biz/SectionHeader.vue'
import EmptyState from '@/components/biz/EmptyState.vue'
import AppIcon from '@/components/common/AppIcon.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import SkeletonCard from '@/components/common/SkeletonCard.vue'
import CouponForm from '@/components/marketing/CouponForm.vue'
import CouponSendPopup from '@/components/marketing/CouponSendPopup.vue'
import { useCouponList } from '@/composables/useCouponList'
import { createCoupon, updateCoupon, deleteCoupon, sendCoupon } from '@/api/marketing'
import { formatActivityTime } from '@/utils/marketing'

const h = useCouponList()
const { list, loading, applyFilter, load } = h

const showForm = ref(false)
const showSend = ref(false)
const editingCoupon = ref<any>(null)
const activeCoupon = ref<any>(null)
const keyword = ref('')
const statusFilter = ref<'all' | 1 | 0>('all')
const formSubmitting = ref(false)
const sendSubmitting = ref(false)

const TYPE_LABELS: Record<number, string> = { 1: '满减', 2: '折扣', 3: '立减' }
const STACK_LABELS: Record<string, string> = {
  exclusive: '互斥',
  same_type: '同类可叠',
  cross_type: '跨类可叠',
}
const TARGET_LABELS: Record<string, string> = {
  all: '全部用户',
  vip_level: '会员等级',
  new_user: '新用户',
}
const FILTERS: Array<{ key: 'all' | 1 | 0; label: string }> = [
  { key: 'all', label: '全部' },
  { key: 1, label: '启用中' },
  { key: 0, label: '已停用' },
]

const stats = computed(() => ({
  total: list.value.length,
  enabled: list.value.filter((item) => Number(item.status || 0) === 1).length,
  disabled: list.value.filter((item) => Number(item.status || 0) !== 1).length,
}))

const filtered = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  return list.value.filter((item) => {
    if (statusFilter.value !== 'all' && Number(item.status || 0) !== statusFilter.value) return false
    if (kw && !String(item.name || '').toLowerCase().includes(kw)) return false
    return true
  })
})

function faceValue(item: any) {
  const value = Number(item.discount || 0)
  return Number(item.type || 1) === 2 ? `${value} 折` : `¥${value}`
}

function period(item: any) {
  if (!item.start_at && !item.end_at) return '长期有效'
  return `${formatActivityTime(item.start_at)} — ${formatActivityTime(item.end_at)}`
}

function openCreate() {
  editingCoupon.value = null
  showForm.value = true
}

function openEdit(item: any) {
  editingCoupon.value = item
  showForm.value = true
}

function openSend(item: any) {
  activeCoupon.value = item
  showSend.value = true
}

async function onSubmitForm(payload: any) {
  if (formSubmitting.value) return
  formSubmitting.value = true
  try {
    if (editingCoupon.value) await updateCoupon(editingCoupon.value.id, payload)
    else await createCoupon(payload)
    showForm.value = false
    uni.showToast({ title: editingCoupon.value ? '优惠券已更新' : '优惠券已创建', icon: 'success' })
    await load()
  } catch (error: any) {
    if (!error?.__toasted) uni.showToast({ title: String(error?.message || '保存失败').slice(0, 20), icon: 'none' })
  } finally {
    formSubmitting.value = false
  }
}

async function onDelete(item: any) {
  uni.showModal({
    title: '删除优惠券',
    content: `确认删除「${item.name}」？已发放的券不受影响。`,
    confirmText: '删除',
    confirmColor: '#b42318',
    success: async (res) => {
      if (!res.confirm) return
      try {
        await deleteCoupon(item.id)
        uni.showToast({ title: '优惠券已删除', icon: 'success' })
        await load()
      } catch (error: any) {
        if (!error?.__toasted) uni.showToast({ title: '删除失败', icon: 'none' })
      }
    },
  })
}

async function onSend(payload: { count: number }) {
  if (!activeCoupon.value || sendSubmitting.value) return
  sendSubmitting.value = true
  try {
    await sendCoupon(activeCoupon.value.id, payload)
    uni.showToast({ title: `已发送 ${payload.count} 张`, icon: 'success' })
    showSend.value = false
    await load()
  } catch (error: any) {
    if (!error?.__toasted) uni.showToast({ title: '发送失败', icon: 'none' })
  } finally {
    sendSubmitting.value = false
  }
}

function onSearch() {
  applyFilter({ keyword: keyword.value || undefined })
  load()
}

onShow(load)
</script>

<template>
  <view class="coupon-page">
    <PageHeader title="优惠券管理" eyebrow="COUPON MANAGEMENT" subtitle="满减 / 折扣券与定向发券" back>
      <template #right>
        <view class="header-action" hover-class="eapp-hover" @click="openCreate">
          <AppIcon name="plus" :size="26" color="white" />
          <text>新建</text>
        </view>
      </template>
    </PageHeader>

    <view class="coupon-page__body">
      <view class="overview">
        <view class="overview__item">
          <text class="overview__value eapp-num">{{ stats.total }}</text>
          <text class="overview__label">全部优惠券</text>
        </view>
        <view class="overview__divider" />
        <view class="overview__item">
          <text class="overview__value overview__value--success eapp-num">{{ stats.enabled }}</text>
          <text class="overview__label">启用中</text>
        </view>
        <view class="overview__divider" />
        <view class="overview__item">
          <text class="overview__value eapp-num">{{ stats.disabled }}</text>
          <text class="overview__label">已停用</text>
        </view>
      </view>

      <view class="search-bar">
        <AppIcon name="search" :size="28" color="faint" />
        <input
          v-model="keyword"
          class="search-bar__input"
          placeholder="搜索优惠券名称"
          placeholder-style="color:#98a2b3"
          confirm-type="search"
          @confirm="onSearch"
        />
        <view v-if="keyword" class="search-bar__clear" @click="keyword = ''">
          <AppIcon name="close" :size="24" color="faint" />
        </view>
      </view>

      <view class="filter-row">
        <view
          v-for="item in FILTERS"
          :key="String(item.key)"
          class="filter-chip"
          :class="{ 'filter-chip--active': statusFilter === item.key }"
          hover-class="eapp-hover"
          @click="statusFilter = item.key"
        >
          {{ item.label }}
        </view>
      </view>

      <SectionHeader title="优惠券列表" :subtitle="`${filtered.length} 张`" />

      <view v-if="loading" class="list">
        <SkeletonCard variant="default" />
        <SkeletonCard variant="default" />
      </view>

      <view v-else-if="!filtered.length" class="empty-wrap">
        <EmptyState
          icon="coupon"
          :title="keyword || statusFilter !== 'all' ? '没有符合条件的优惠券' : '还没有优惠券'"
          :desc="keyword || statusFilter !== 'all' ? '换个状态或清空关键词试试' : '创建第一张优惠券，用于满减、折扣与定向发放'"
        >
          <view class="empty-btn" @click="openCreate"><text>新建优惠券</text></view>
        </EmptyState>
      </view>

      <view v-else class="list">
        <view v-for="item in filtered" :key="item.id" class="coupon-card">
          <view class="coupon-card__top">
            <view class="coupon-card__face">
              <text class="coupon-card__face-value eapp-num">{{ faceValue(item) }}</text>
              <text class="coupon-card__face-label">{{ TYPE_LABELS[Number(item.type)] || '优惠' }}</text>
            </view>
            <view class="coupon-card__info">
              <view class="coupon-card__title-line">
                <text class="coupon-card__title">{{ item.name }}</text>
                <StatusTag :text="Number(item.status) === 1 ? '启用' : '停用'" :type="Number(item.status) === 1 ? 'success' : 'neutral'" />
              </view>
              <text class="coupon-card__desc">
                满 {{ Number(item.min_amount || 0) }} 可用 · 已用 {{ item.used_count || 0 }} 张
              </text>
              <text class="coupon-card__desc eapp-num">{{ period(item) }}</text>
            </view>
          </view>

          <view class="coupon-card__tags">
            <text class="mini-tag">{{ STACK_LABELS[item.stack_rule] || '叠加规则未设置' }}</text>
            <text class="mini-tag">{{ TARGET_LABELS[item.target_type] || '全部用户' }}</text>
            <text class="mini-tag eapp-num">发行 {{ item.total_count || '不限' }}</text>
          </view>

          <view class="coupon-card__foot">
            <text class="coupon-card__id eapp-num">#{{ item.id }}</text>
            <view class="coupon-card__actions">
              <view class="mini-action" hover-class="eapp-hover" @click="openEdit(item)">
                <AppIcon name="edit" :size="24" color="primary" /><text>编辑</text>
              </view>
              <view class="mini-action" hover-class="eapp-hover" @click="openSend(item)">
                <AppIcon name="share" :size="24" color="muted" /><text>发券</text>
              </view>
              <view class="mini-action mini-action--danger" hover-class="eapp-hover" @click="onDelete(item)">
                <AppIcon name="close" :size="24" color="danger" /><text>删除</text>
              </view>
            </view>
          </view>
        </view>
      </view>
    </view>

    <CouponForm
      :show="showForm"
      :coupon="editingCoupon"
      :submitting="formSubmitting"
      @close="showForm = false"
      @submit="onSubmitForm"
    />
    <CouponSendPopup
      :show="showSend"
      :coupon="activeCoupon"
      :submitting="sendSubmitting"
      @close="showSend = false"
      @send="onSend"
    />
  </view>
</template>

<style scoped>
.coupon-page {
  min-height: 100vh;
  background: var(--eapp-bg);
}
.coupon-page__body {
  padding: var(--eapp-space-5);
  padding-bottom: calc(env(safe-area-inset-bottom) + 48rpx);
  box-sizing: border-box;
}
.header-action {
  display: flex;
  align-items: center;
  gap: 6rpx;
  height: 60rpx;
  padding: 0 var(--eapp-space-3);
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary);
  color: var(--eapp-primary-on);
  font-size: var(--eapp-fs-sm);
  font-weight: 600;
}
.overview {
  display: flex;
  align-items: center;
  padding: var(--eapp-space-4) var(--eapp-space-2);
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  box-sizing: border-box;
}
.overview__item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4rpx;
}
.overview__value {
  font-size: var(--eapp-fs-lg);
  font-weight: 700;
  color: var(--eapp-text);
}
.overview__value--success {
  color: var(--eapp-success);
}
.overview__label {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}
.overview__divider {
  width: 1rpx;
  height: 52rpx;
  background: var(--eapp-divider);
}
.search-bar {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-2);
  height: 80rpx;
  margin-top: var(--eapp-space-5);
  padding: 0 var(--eapp-space-4);
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius);
  box-sizing: border-box;
}
.search-bar__input {
  flex: 1;
  min-width: 0;
  height: 100%;
  font-size: var(--eapp-fs-base);
  color: var(--eapp-text);
}
.search-bar__clear {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}
.filter-row {
  display: flex;
  gap: var(--eapp-space-2);
  margin: var(--eapp-space-5) 0 calc(var(--eapp-space-6) + var(--eapp-space-2));
}
.filter-chip {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 60rpx;
  border-radius: var(--eapp-radius-pill);
  border: 1rpx solid var(--eapp-border);
  background: var(--eapp-card);
  color: var(--eapp-text-secondary);
  font-size: var(--eapp-fs-sm);
  box-sizing: border-box;
}
.filter-chip--active {
  border-color: var(--eapp-primary);
  background: var(--eapp-primary-soft);
  color: var(--eapp-primary);
  font-weight: 600;
}
.list {
  display: flex;
  flex-direction: column;
  gap: var(--eapp-space-3);
}
.empty-wrap {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
}
.empty-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 72rpx;
  padding: 0 var(--eapp-space-6);
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary);
  color: var(--eapp-primary-on);
  font-size: var(--eapp-fs-sm);
  font-weight: 600;
}
.coupon-card {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-4);
  box-sizing: border-box;
}
.coupon-card__top {
  display: flex;
  gap: var(--eapp-space-3);
}
.coupon-card__face {
  width: 160rpx;
  min-height: 136rpx;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4rpx;
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary-soft);
  border: 1rpx solid var(--eapp-primary-border);
}
.coupon-card__face-value {
  font-size: var(--eapp-fs-lg);
  font-weight: 700;
  color: var(--eapp-primary);
}
.coupon-card__face-label {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-muted);
}
.coupon-card__info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}
.coupon-card__title-line {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-2);
  min-width: 0;
}
.coupon-card__title {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.coupon-card__desc {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-muted);
}
.coupon-card__tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--eapp-space-2);
  margin-top: var(--eapp-space-3);
}
.mini-tag {
  padding: 4rpx 12rpx;
  border-radius: var(--eapp-radius-sm);
  background: var(--eapp-surface-subtle);
  border: 1rpx solid var(--eapp-divider);
  color: var(--eapp-text-muted);
  font-size: var(--eapp-fs-xs);
}
.coupon-card__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--eapp-space-2);
  margin-top: var(--eapp-space-3);
  padding-top: var(--eapp-space-3);
  border-top: 1rpx solid var(--eapp-divider);
}
.coupon-card__id {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}
.coupon-card__actions {
  display: flex;
  gap: var(--eapp-space-2);
}
.mini-action {
  display: flex;
  align-items: center;
  gap: 4rpx;
  height: 54rpx;
  padding: 0 var(--eapp-space-2);
  border-radius: var(--eapp-radius-sm);
  border: 1rpx solid var(--eapp-border);
  background: var(--eapp-surface-subtle);
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-secondary);
  box-sizing: border-box;
}
.mini-action--danger {
  color: var(--eapp-danger);
  border-color: var(--eapp-danger-soft);
  background: var(--eapp-danger-soft);
}
</style>
