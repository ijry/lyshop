<script setup lang="ts">
import StatusTag from '@/components/common/StatusTag.vue'
defineProps<{ row: any }>()
defineEmits<{ (e: 'click'): void }>()
</script>

<template>
  <view class="card" @click="$emit('click')">
    <view class="head">
      <text class="no">售后 #{{ row.id }}</text>
      <StatusTag :text="row.status_label || row.status || '-'" :type="row.status" />
    </view>
    <view class="line">订单：#{{ row.order_id }} · 类型：{{ row.case_type_label || row.case_type || '-' }}</view>
    <view class="line">原因：{{ row.reason || '-' }}</view>
    <view class="line">退款金额：¥{{ Number(row.refund_amount || 0).toFixed(2) }}</view>
  </view>
</template>

<style scoped>
.card { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; }
.head { display: flex; justify-content: space-between; align-items: center; }
.no { font-size: 26rpx; font-weight: 600; }
.line { margin-top: 6rpx; color: var(--eapp-text-muted); font-size: 23rpx; }
</style>
