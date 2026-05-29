<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{ show: boolean; coupon: any }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'send', payload: { count: number }): void }>()

const count = ref(1)

function submit() {
  if (count.value < 1) { uni.showToast({ title: '请输入发送数量', icon: 'none' }); return }
  emit('send', { count: count.value })
}
</script>

<template>
  <up-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="popup-body" v-if="coupon">
      <view class="popup-title">定向发券</view>
      <view class="info">券名：{{ coupon.name }}</view>
      <view class="info">目标：{{ { all: '全部用户', vip_level: '会员等级', new_user: '新用户' }[coupon.target_type as string] || coupon.target_type }}</view>
      <view class="info" v-if="coupon.target_type === 'vip_level'">等级值：{{ coupon.target_value }}</view>
      <view class="mt" />
      <up-input v-model="count" type="number" placeholder="发送数量" />
      <view class="mt-lg" />
      <up-button type="primary" @click="submit">确认发送</up-button>
    </view>
  </up-popup>
</template>

<style scoped>
.popup-body { padding: 24rpx; box-sizing: border-box; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.info { font-size: 26rpx; color: var(--eapp-text-muted); margin-top: 8rpx; }
.mt { margin-top: 12rpx; }
.mt-lg { margin-top: 16rpx; }
</style>
