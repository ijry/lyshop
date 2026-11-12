<script setup lang="ts">
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import PageHeader from '@/components/biz/PageHeader.vue'
import SectionHeader from '@/components/biz/SectionHeader.vue'
import EmptyState from '@/components/biz/EmptyState.vue'
import AppIcon from '@/components/common/AppIcon.vue'
import SkeletonCard from '@/components/common/SkeletonCard.vue'
import ActivityCard from './ActivityCard.vue'
import ActivityFormPopup from './ActivityFormPopup.vue'
import { useActivityCenter } from '@/composables/useActivityCenter'
import {
  ACTIVITY_FILTERS,
  activityDetailUrl,
  activityMeta,
  type MarketingKind,
} from '@/utils/marketing'

const props = defineProps<{ kind: MarketingKind }>()
const meta = computed(() => activityMeta(props.kind))

const {
  productCounts,
  loading,
  submitting,
  keyword,
  phase,
  filtered,
  stats,
  load,
  save,
  toggle,
  remove,
} = useActivityCenter(props.kind)

const showForm = ref(false)
const editing = ref<any>(null)
const hasFilter = computed(() => phase.value !== 'all' || keyword.value.trim().length > 0)

function openCreate() {
  editing.value = null
  showForm.value = true
}

function openEdit(item: any) {
  editing.value = item
  showForm.value = true
}

async function onSubmit(payload: any) {
  const ok = await save(payload, Number(editing.value?.id || 0))
  if (ok) showForm.value = false
}

async function onToggle(item: any) {
  await toggle(item)
}

function onRemove(item: any) {
  uni.showModal({
    title: '删除活动',
    content: `确认删除「${item.name || '未命名活动'}」？活动下的商品配置会一并移除。`,
    confirmText: '删除',
    confirmColor: '#b42318',
    success: async (res) => {
      if (res.confirm) await remove(Number(item.id || 0))
    },
  })
}

function onOpen(item: any) {
  uni.navigateTo({ url: activityDetailUrl(props.kind, item.id, item.name) })
}

function resetFilter() {
  keyword.value = ''
  phase.value = 'all'
}

onShow(load)
</script>

<template>
  <view class="activity-page">
    <PageHeader :title="meta.title" :eyebrow="meta.eyebrow" :subtitle="meta.desc" back>
      <template #right>
        <view class="header-action" hover-class="eapp-hover" @click="openCreate">
          <AppIcon name="plus" :size="26" color="white" />
          <text>新建</text>
        </view>
      </template>
    </PageHeader>

    <view class="activity-page__body">
      <!-- 概览 -->
      <view class="overview">
        <view class="overview__item">
          <text class="overview__value eapp-num">{{ stats.total }}</text>
          <text class="overview__label">全部活动</text>
        </view>
        <view class="overview__divider" />
        <view class="overview__item">
          <text class="overview__value overview__value--success eapp-num">{{ stats.running }}</text>
          <text class="overview__label">进行中</text>
        </view>
        <view class="overview__divider" />
        <view class="overview__item">
          <text class="overview__value overview__value--primary eapp-num">{{ stats.upcoming }}</text>
          <text class="overview__label">未开始</text>
        </view>
        <view class="overview__divider" />
        <view class="overview__item">
          <text class="overview__value eapp-num">{{ stats.disabled }}</text>
          <text class="overview__label">已停用</text>
        </view>
      </view>

      <!-- 搜索 -->
      <view class="search-bar">
        <AppIcon name="search" :size="28" color="faint" />
        <input
          v-model="keyword"
          class="search-bar__input"
          placeholder="搜索活动名称"
          placeholder-style="color:#98a2b3"
          confirm-type="search"
        />
        <view v-if="keyword" class="search-bar__clear" hover-class="eapp-hover" @click="keyword = ''">
          <AppIcon name="close" :size="24" color="faint" />
        </view>
      </view>

      <!-- 状态筛选 -->
      <scroll-view scroll-x class="filter-scroll">
        <view class="filter-row">
          <view
            v-for="item in ACTIVITY_FILTERS"
            :key="item.key"
            class="filter-chip"
            :class="{ 'filter-chip--active': phase === item.key }"
            hover-class="eapp-hover"
            @click="phase = item.key"
          >
            {{ item.label }}
          </view>
        </view>
      </scroll-view>

      <SectionHeader
        title="活动列表"
        :subtitle="`${filtered.length} 场`"
        extra="重置筛选"
        arrow
        @click="hasFilter ? resetFilter() : load()"
      />

      <view v-if="loading" class="list">
        <SkeletonCard variant="default" />
        <SkeletonCard variant="default" />
      </view>

      <view v-else-if="!filtered.length" class="empty-wrap">
        <EmptyState
          :icon="meta.iconName"
          :title="hasFilter ? '没有符合条件的活动' : '还没有活动'"
          :desc="hasFilter ? '换个状态或清空关键词试试' : '创建第一场活动，开始配置活动商品与价格'"
        >
          <view class="empty-btn" hover-class="eapp-hover" @click="hasFilter ? resetFilter() : openCreate()">
            <text>{{ hasFilter ? '重置筛选' : '新建活动' }}</text>
          </view>
        </EmptyState>
      </view>

      <view v-else class="list">
        <ActivityCard
          v-for="item in filtered"
          :key="item.id"
          :activity="item"
          :kind="kind"
          :product-count="productCounts[item.id] || 0"
          @open="onOpen(item)"
          @edit="openEdit(item)"
          @toggle="onToggle(item)"
          @remove="onRemove(item)"
        />
      </view>
    </view>

    <ActivityFormPopup
      :show="showForm"
      :kind="kind"
      :activity="editing"
      :submitting="submitting"
      @close="showForm = false"
      @submit="onSubmit"
    />
  </view>
</template>

<style scoped>
.activity-page {
  min-height: 100vh;
  background: var(--eapp-bg);
  box-sizing: border-box;
}

.activity-page__body {
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
  box-sizing: border-box;
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
  min-width: 0;
}

.overview__value {
  font-size: var(--eapp-fs-lg);
  font-weight: 700;
  color: var(--eapp-text);
}

.overview__value--success {
  color: var(--eapp-success);
}

.overview__value--primary {
  color: var(--eapp-primary);
}

.overview__label {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}

.overview__divider {
  width: 1rpx;
  height: 52rpx;
  background: var(--eapp-divider);
  flex-shrink: 0;
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

.search-bar__placeholder {
  color: var(--eapp-text-faint);
}

.search-bar__clear {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.filter-scroll {
  margin: var(--eapp-space-5) calc(0rpx - var(--eapp-space-5)) calc(var(--eapp-space-6) + var(--eapp-space-2));
  padding: var(--eapp-space-2) var(--eapp-space-5);
  white-space: nowrap;
}

.filter-row {
  display: inline-flex;
  gap: var(--eapp-space-2);
  padding: var(--eapp-space-1) 0;
}

.filter-chip {
  display: inline-flex;
  align-items: center;
  height: 60rpx;
  padding: 0 var(--eapp-space-4);
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
</style>
