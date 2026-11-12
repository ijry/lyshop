<script setup lang="ts">
import { computed } from 'vue'
import AppIcon from '@/components/common/AppIcon.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import {
  activityMeta,
  activityProgress,
  activityState,
  activityTimeRange,
  type MarketingKind,
} from '@/utils/marketing'

const props = defineProps<{
  activity: any
  kind: MarketingKind
  productCount?: number
}>()

const emit = defineEmits<{
  open: []
  edit: []
  toggle: []
  remove: []
}>()

const meta = computed(() => activityMeta(props.kind))
const state = computed(() => activityState(props.activity))
const progress = computed(() => activityProgress(props.activity))
const range = computed(() => activityTimeRange(props.activity))
const iconColor = computed(() => {
  if (state.value.phase === 'running') return 'success'
  if (state.value.phase === 'upcoming') return 'primary'
  return 'faint'
})
const enabled = computed(() => Number(props.activity?.status || 0) === 1)
</script>

<template>
  <view class="activity-card" hover-class="eapp-hover" @click="emit('open')">
    <view class="activity-card__top">
      <view class="activity-card__icon">
        <AppIcon :name="meta.iconName" :size="34" :color="iconColor" />
      </view>
      <view class="activity-card__heading">
        <view class="activity-card__title-line">
          <text class="activity-card__title">{{ activity.name || '未命名活动' }}</text>
          <StatusTag :text="state.label" :type="state.tone" />
        </view>
        <text class="activity-card__sub eapp-num">{{ meta.title }} · #{{ activity.id }}</text>
      </view>
      <AppIcon name="chevron-right" :size="28" color="faint" />
    </view>

    <view class="activity-card__schedule">
      <AppIcon name="calendar" :size="24" color="faint" />
      <text class="activity-card__range eapp-num">{{ range }}</text>
    </view>

    <view v-if="state.phase === 'running' || state.phase === 'ended'" class="activity-card__progress">
      <view class="activity-card__track">
        <view
          class="activity-card__fill"
          :class="{ 'activity-card__fill--done': state.phase === 'ended' }"
          :style="{ width: `${progress}%` }"
        />
      </view>
      <text class="activity-card__pct eapp-num">{{ progress }}%</text>
    </view>

    <view class="activity-card__divider" />

    <view class="activity-card__foot">
      <view class="activity-card__metric">
        <text class="activity-card__metric-label">活动商品</text>
        <text class="activity-card__metric-value eapp-num">{{ productCount ?? '—' }}</text>
      </view>
      <view class="activity-card__actions">
        <view class="mini-action" hover-class="eapp-hover" @click.stop="emit('edit')">
          <AppIcon name="edit" :size="24" color="primary" />
          <text>编辑</text>
        </view>
        <view class="mini-action" hover-class="eapp-hover" @click.stop="emit('toggle')">
          <AppIcon :name="enabled ? 'lock' : 'check-circle'" :size="24" color="muted" />
          <text>{{ enabled ? '停用' : '启用' }}</text>
        </view>
        <view class="mini-action mini-action--danger" hover-class="eapp-hover" @click.stop="emit('remove')">
          <AppIcon name="close" :size="24" color="danger" />
          <text>删除</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.activity-card {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-5);
  box-sizing: border-box;
}

.activity-card__top {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
}

.activity-card__icon {
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

.activity-card__heading {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.activity-card__title-line {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-2);
  min-width: 0;
}

.activity-card__title {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 320rpx;
}

.activity-card__sub {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}

.activity-card__schedule {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-2);
  margin-top: var(--eapp-space-4);
  padding: var(--eapp-space-2) var(--eapp-space-3);
  background: var(--eapp-surface-subtle);
  border-radius: var(--eapp-radius-sm);
}

.activity-card__range {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-secondary);
}

.activity-card__progress {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
  margin-top: var(--eapp-space-3);
}

.activity-card__track {
  flex: 1;
  height: 10rpx;
  border-radius: var(--eapp-radius-pill);
  background: var(--eapp-divider);
  overflow: hidden;
}

.activity-card__fill {
  height: 100%;
  border-radius: var(--eapp-radius-pill);
  background: var(--eapp-primary);
}

.activity-card__fill--done {
  background: var(--eapp-border-strong);
}

.activity-card__pct {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-muted);
  width: 60rpx;
  text-align: right;
}

.activity-card__divider {
  height: 1rpx;
  background: var(--eapp-divider);
  margin: var(--eapp-space-4) 0 var(--eapp-space-3);
}

.activity-card__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--eapp-space-3);
}

.activity-card__metric {
  display: flex;
  flex-direction: column;
  gap: 2rpx;
  flex-shrink: 0;
}

.activity-card__metric-label {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}

.activity-card__metric-value {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
}

.activity-card__actions {
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
