<script setup lang="ts">
import StatusTag from '@/components/common/StatusTag.vue'

defineProps<{ row: any }>()
defineEmits<{ click: [] }>()

function stars(score: number) {
  const s = Math.round(Number(score || 0) * 10) / 10
  return s.toFixed(1)
}

function truncate(text: string, max = 60) {
  const s = String(text || '')
  return s.length > max ? s.slice(0, max) + '...' : s
}
</script>

<template>
  <view class="review-card" @click="$emit('click')">
    <view class="review-card__product">
      <image
        v-if="row.product?.cover || row.order_item?.cover"
        :src="row.product?.cover || row.order_item?.cover"
        mode="aspectFill"
        class="review-card__cover"
      />
      <view class="review-card__product-info">
        <text class="review-card__title">{{ row.product?.title || row.order_item?.title || '-' }}</text>
        <text class="review-card__scores">
          {{ '★' }}{{ stars(row.product_score) }} / {{ '★' }}{{ stars(row.logistics_score) }}
        </text>
      </view>
      <StatusTag
        :text="row.admin_reply ? '已回复' : '待回复'"
        :type="row.admin_reply ? 'success' : 'pending'"
      />
    </view>
    <text class="review-card__content">{{ truncate(row.content) }}</text>
    <view v-if="row.images && row.images.length" class="review-card__images">
      <image
        v-for="(img, idx) in row.images.slice(0, 4)"
        :key="idx"
        :src="img"
        mode="aspectFill"
        class="review-card__img"
      />
    </view>
    <view class="review-card__footer">
      <text class="review-card__user">{{ row.user_nickname || '匿名' }}</text>
      <text v-if="row.appends && row.appends.length" class="review-card__append">
        追评 {{ row.appends.length }}
      </text>
    </view>
  </view>
</template>

<style scoped>
.review-card {
  background: var(--eapp-card);
  border: 1px solid var(--eapp-border);
  border-radius: 22rpx;
  padding: 20rpx;
}
.review-card__product {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 10rpx;
}
.review-card__cover {
  width: 72rpx;
  height: 72rpx;
  border-radius: 12rpx;
  flex-shrink: 0;
}
.review-card__product-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}
.review-card__title {
  font-size: 24rpx;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.review-card__scores {
  font-size: 22rpx;
  color: #f59e0b;
}
.review-card__content {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  font-size: 24rpx;
  color: var(--eapp-text-muted);
  line-height: 1.6;
}
.review-card__images {
  display: flex;
  gap: 8rpx;
  margin-top: 10rpx;
}
.review-card__img {
  width: 100rpx;
  height: 100rpx;
  border-radius: 10rpx;
}
.review-card__footer {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-top: 10rpx;
}
.review-card__user {
  font-size: 22rpx;
  color: var(--eapp-text-muted);
}
.review-card__append {
  font-size: 20rpx;
  color: var(--eapp-primary);
  background: var(--eapp-bg);
  padding: 4rpx 12rpx;
  border-radius: 999rpx;
}
</style>
