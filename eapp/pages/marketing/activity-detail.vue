<script setup lang="ts">
import { computed, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import PageHeader from '@/components/biz/PageHeader.vue'
import SectionHeader from '@/components/biz/SectionHeader.vue'
import EmptyState from '@/components/biz/EmptyState.vue'
import AppIcon from '@/components/common/AppIcon.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import SkeletonCard from '@/components/common/SkeletonCard.vue'
import ActivityFormPopup from '@/components/marketing/ActivityFormPopup.vue'
import ActivityProductCard from '@/components/marketing/ActivityProductCard.vue'
import ActivityProductFormPopup from '@/components/marketing/ActivityProductFormPopup.vue'
import { useActivityCenter, useActivityProducts } from '@/composables/useActivityCenter'
import {
  activityMeta,
  activityProgress,
  activityState,
  activityTimeRange,
  formatMoney,
  isActivityKind,
  type MarketingKind,
} from '@/utils/marketing'

const kind = ref<MarketingKind>('seckill')
const activityID = ref(0)
const activityName = ref('活动详情')

const center = useActivityCenter(kind)
const products = useActivityProducts(kind, activityID)

const {
  activities,
  submitting: activitySubmitting,
  load: loadActivities,
  save: saveActivity,
} = center
const {
  rows,
  loading: rowsLoading,
  submitting: productSubmitting,
  keyword,
  productOptions,
  skuOptions,
  filtered: filteredRows,
  stats,
  load: loadRows,
  loadProducts,
  loadSkus,
  save: saveRow,
  remove: removeRow,
} = products

const meta = computed(() => activityMeta(kind.value))
const activity = computed(
  () => activities.value.find((item) => Number(item.id || 0) === activityID.value) || null,
)
const state = computed(() => activityState(activity.value || {}))
const progress = computed(() => activityProgress(activity.value || {}))
const range = computed(() => activityTimeRange(activity.value || {}))

const showActivityForm = ref(false)
const showProductForm = ref(false)
const editingRow = ref<any>(null)

onLoad((options: any) => {
  kind.value = isActivityKind(options?.kind) ? options.kind : 'seckill'
  activityID.value = Number(options?.id || 0)
  activityName.value = options?.name
    ? decodeURIComponent(String(options.name))
    : activityMeta(kind.value).title
  loadProducts()
})

onShow(() => {
  loadActivities()
  loadRows()
})

function openCreateProduct() {
  editingRow.value = null
  skuOptions.value = []
  showProductForm.value = true
}

function openEditProduct(item: any) {
  editingRow.value = item
  loadSkus(Number(item.product_id || 0))
  showProductForm.value = true
}

async function onSelectProduct(productID: number) {
  await loadSkus(productID)
}

async function onSubmitProduct(payload: any) {
  const ok = await saveRow(payload, Number(editingRow.value?.id || 0))
  if (ok) showProductForm.value = false
}

function onRemoveProduct(item: any) {
  uni.showModal({
    title: '移除商品',
    content: `确认将「${item.product_title || `商品 #${item.product_id}`}」移出本场活动？`,
    confirmText: '移除',
    confirmColor: '#b42318',
    success: async (res) => {
      if (res.confirm) await removeRow(Number(item.id || 0))
    },
  })
}

async function onSubmitActivity(payload: any) {
  const ok = await saveActivity(payload, activityID.value)
  if (ok) showActivityForm.value = false
}
</script>

<template>
  <view class="detail-page">
    <PageHeader :title="activityName" :eyebrow="meta.eyebrow" :subtitle="range" back>
      <template #right>
        <view class="eapp-icon-btn" hover-class="eapp-hover" @click="showActivityForm = true">
          <AppIcon name="edit" :size="30" color="primary" />
        </view>
      </template>
    </PageHeader>

    <view class="detail-page__body">
      <view class="summary">
        <view class="summary__head">
          <view class="summary__icon">
            <AppIcon :name="meta.iconName" :size="36" color="primary" />
          </view>
          <view class="summary__heading">
            <text class="summary__title">{{ activity?.name || activityName }}</text>
            <text class="summary__range eapp-num">{{ range }}</text>
          </view>
          <StatusTag :text="state.label" :type="state.tone" />
        </view>

        <view class="summary__progress">
          <view class="summary__track">
            <view class="summary__fill" :style="{ width: `${progress}%` }" />
          </view>
          <text class="summary__pct eapp-num">{{ progress }}%</text>
        </view>

        <view class="summary__metrics">
          <view class="summary__metric">
            <text class="summary__metric-label">活动商品</text>
            <text class="summary__metric-value eapp-num">{{ stats.count }}</text>
          </view>
          <view class="summary__metric">
            <text class="summary__metric-label">累计已售</text>
            <text class="summary__metric-value eapp-num">{{ stats.soldQty }}</text>
          </view>
          <view class="summary__metric">
            <text class="summary__metric-label">预估销售额</text>
            <text class="summary__metric-value eapp-num">{{ formatMoney(stats.gross) }}</text>
          </view>
        </view>
      </view>

      <SectionHeader title="活动商品" :subtitle="`${stats.count} 个 SKU`" />

      <view class="search-bar">
        <AppIcon name="search" :size="28" color="faint" />
        <input
          v-model="keyword"
          class="search-bar__input"
          placeholder="搜索活动商品"
          placeholder-style="color:#98a2b3"
          confirm-type="search"
        />
        <view v-if="keyword" class="search-bar__clear" @click="keyword = ''">
          <AppIcon name="close" :size="24" color="faint" />
        </view>
      </view>

      <view v-if="rowsLoading" class="list">
        <SkeletonCard variant="default" />
      </view>

      <view v-else-if="!filteredRows.length" class="empty-wrap">
        <EmptyState
          icon="product"
          :title="keyword ? '没有匹配的商品' : '还没有活动商品'"
          :desc="keyword ? '换个关键词试试' : '添加商品并设置活动价、限购与库存'"
        >
          <view class="empty-btn" @click="keyword ? (keyword = '') : openCreateProduct()">
            <text>{{ keyword ? '清空搜索' : '添加商品' }}</text>
          </view>
        </EmptyState>
      </view>

      <view v-else class="list">
        <ActivityProductCard
          v-for="item in filteredRows"
          :key="item.id || `${item.product_id}-${item.sku_id}`"
          :row="item"
          :kind="kind"
          @edit="openEditProduct(item)"
          @remove="onRemoveProduct(item)"
        />
      </view>
    </view>

    <view class="footbar">
      <view class="footbar__btn" hover-class="eapp-hover" @click="openCreateProduct">
        <AppIcon name="plus" :size="30" color="white" />
        <text>添加活动商品</text>
      </view>
    </view>

    <ActivityFormPopup
      :show="showActivityForm"
      :kind="kind"
      :activity="activity"
      :submitting="activitySubmitting"
      @close="showActivityForm = false"
      @submit="onSubmitActivity"
    />
    <ActivityProductFormPopup
      :show="showProductForm"
      :kind="kind"
      :row="editingRow"
      :products="productOptions"
      :skus="skuOptions"
      :submitting="productSubmitting"
      @close="showProductForm = false"
      @submit="onSubmitProduct"
      @select-product="onSelectProduct"
    />
  </view>
</template>

<style scoped>
.detail-page {
  min-height: 100vh;
  background: var(--eapp-bg);
}
.detail-page__body {
  padding: var(--eapp-space-5);
  padding-bottom: calc(env(safe-area-inset-bottom) + 160rpx);
  box-sizing: border-box;
}
.summary {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-5);
  margin-bottom: var(--eapp-space-6);
  box-sizing: border-box;
}
.summary__head {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
}
.summary__icon {
  width: 72rpx;
  height: 72rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary-soft);
  border: 1rpx solid var(--eapp-primary-border);
  box-sizing: border-box;
}
.summary__heading {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}
.summary__title {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.summary__range {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-muted);
}
.summary__progress {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
  margin-top: var(--eapp-space-4);
}
.summary__track {
  flex: 1;
  height: 10rpx;
  border-radius: var(--eapp-radius-pill);
  background: var(--eapp-divider);
  overflow: hidden;
}
.summary__fill {
  height: 100%;
  background: var(--eapp-primary);
  border-radius: var(--eapp-radius-pill);
}
.summary__pct {
  width: 64rpx;
  text-align: right;
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-muted);
}
.summary__metrics {
  display: flex;
  margin-top: var(--eapp-space-5);
  padding-top: var(--eapp-space-4);
  border-top: 1rpx solid var(--eapp-divider);
}
.summary__metric {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  min-width: 0;
}
.summary__metric-label {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}
.summary__metric-value {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.search-bar {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-2);
  height: 80rpx;
  padding: 0 var(--eapp-space-4);
  margin-bottom: var(--eapp-space-4);
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
.footbar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  padding: var(--eapp-space-3) var(--eapp-space-5)
    calc(var(--eapp-space-3) + env(safe-area-inset-bottom));
  background: var(--eapp-card);
  border-top: 1rpx solid var(--eapp-border);
  box-sizing: border-box;
}
.footbar__btn {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--eapp-space-2);
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary);
  color: var(--eapp-primary-on);
  font-size: var(--eapp-fs-md);
  font-weight: 600;
}
</style>
