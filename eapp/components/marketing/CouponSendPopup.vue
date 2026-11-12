<script setup lang="ts">
import { ref, watch } from 'vue'
import AppIcon from '@/components/common/AppIcon.vue'
import FormSheet from './FormSheet.vue'

const props = withDefaults(
  defineProps<{ show: boolean; coupon?: any; submitting?: boolean }>(),
  { submitting: false },
)
const emit = defineEmits<{ close: []; send: [payload: { count: number }] }>()

const count = ref(10)
const QUICK = [10, 50, 100, 500]

const TARGET_LABELS: Record<string, string> = {
  all: '全部用户',
  vip_level: '会员等级',
  new_user: '新用户',
}

watch(
  () => props.show,
  (visible) => {
    if (visible) count.value = 10
  },
)

function step(delta: number) {
  const next = Number(count.value || 0) + delta
  count.value = Math.max(1, Math.min(99999, next))
}

function onInput(event: any) {
  const value = Number(String(event?.detail?.value || '').replace(/[^0-9]/g, '') || 0)
  count.value = Math.max(0, Math.min(99999, value))
}

function submit() {
  const value = Number(count.value || 0)
  if (value < 1) {
    uni.showToast({ title: '请输入发送数量', icon: 'none' })
    return
  }
  emit('send', { count: value })
}
</script>

<template>
  <FormSheet
    :show="show"
    title="定向发券"
    subtitle="按券设定的目标用户范围发放"
    :submitting="submitting"
    submit-text="确认发送"
    body-height="40vh"
    @close="emit('close')"
    @submit="submit"
  >
    <view v-if="coupon" class="send-coupon">
      <view class="send-coupon__icon">
        <AppIcon name="coupon" :size="36" color="primary" />
      </view>
      <view class="send-coupon__body">
        <text class="send-coupon__name">{{ coupon.name }}</text>
        <text class="send-coupon__meta">
          目标：{{ TARGET_LABELS[String(coupon.target_type)] || '全部用户' }}
        </text>
        <text v-if="coupon.target_type === 'vip_level'" class="send-coupon__meta">
          等级值：{{ coupon.target_value || '-' }}
        </text>
      </view>
    </view>

    <view class="eapp-form-group">
      <view class="eapp-form-group__title">
        <view class="eapp-form-group__bar" />
        <text>发放数量</text>
      </view>
      <view class="stepper">
        <view class="stepper__btn" hover-class="eapp-hover" @click="step(-1)">
          <AppIcon name="minus" :size="30" color="secondary" />
        </view>
        <input
          class="stepper__input eapp-num"
          type="number"
          :value="String(count)"
          @input="onInput"
        />
        <view class="stepper__btn" hover-class="eapp-hover" @click="step(1)">
          <AppIcon name="plus" :size="30" color="secondary" />
        </view>
      </view>

      <view class="quick">
        <view
          v-for="item in QUICK"
          :key="item"
          class="quick__item"
          :class="{ 'quick__item--active': count === item }"
          hover-class="eapp-hover"
          @click="count = item"
        >
          {{ item }} 张
        </view>
      </view>

      <view class="eapp-form-tip">
        <AppIcon name="info" :size="24" color="primary" />
        <text>发放后不可撤销，请确认数量与目标用户范围。</text>
      </view>
    </view>
  </FormSheet>
</template>

<style scoped>
.send-coupon {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
  padding: var(--eapp-space-4);
  margin-bottom: var(--eapp-space-6);
  background: var(--eapp-primary-soft);
  border: 1rpx solid var(--eapp-primary-border);
  border-radius: var(--eapp-radius);
  box-sizing: border-box;
}
.send-coupon__icon {
  width: 72rpx;
  height: 72rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--eapp-radius-sm);
  background: var(--eapp-card);
}
.send-coupon__body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}
.send-coupon__name {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
}
.send-coupon__meta {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-muted);
}
.stepper {
  display: flex;
  align-items: center;
  height: 88rpx;
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius);
  background: var(--eapp-surface-subtle);
  overflow: hidden;
  box-sizing: border-box;
}
.stepper__btn {
  width: 96rpx;
  height: 100%;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--eapp-card);
  border-right: 1rpx solid var(--eapp-border);
}
.stepper__btn:last-child {
  border-right: none;
  border-left: 1rpx solid var(--eapp-border);
}
.stepper__input {
  flex: 1;
  min-width: 0;
  height: 100%;
  text-align: center;
  font-size: var(--eapp-fs-lg);
  font-weight: 700;
  color: var(--eapp-text);
}
.quick {
  display: flex;
  gap: var(--eapp-space-2);
  margin-top: var(--eapp-space-3);
}
.quick__item {
  flex: 1;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-sm);
  background: var(--eapp-card);
  color: var(--eapp-text-secondary);
  font-size: var(--eapp-fs-xs);
  box-sizing: border-box;
}
.quick__item--active {
  border-color: var(--eapp-primary);
  background: var(--eapp-primary-soft);
  color: var(--eapp-primary);
  font-weight: 600;
}
</style>
