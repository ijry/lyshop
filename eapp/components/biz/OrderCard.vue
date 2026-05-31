<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ order: any; selectable?: boolean; selected?: boolean }>()
defineEmits<{ (e: 'click'): void; (e: 'toggle'): void; (e: 'action', key: string): void }>()

const items = computed(() => Array.isArray(props.order?.items) ? props.order.items : [])
const thumbs = computed(() => items.value.slice(0, 3))
const remain = computed(() => Math.max(0, items.value.length - 3))
const itemCount = computed(() => items.value.reduce((s: number, it: any) => s + Number(it.qty || 1), 0))
const firstItem = computed(() => items.value[0])

function money(v: any) { return Number(v || 0).toFixed(2) }

function formatTime(v?: string) {
  if (!v) return ''
  const s = String(v).replace('T', ' ')
  // MM/DD HH:mm
  const m = s.match(/\d{4}-(\d{2})-(\d{2}) (\d{2}:\d{2})/)
  return m ? `${m[1]}/${m[2]} ${m[3]}` : s.slice(5, 16)
}

function maskPhone(v?: string) {
  if (!v) return ''
  return String(v).replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// Left accent bar color by status
const accentColor = computed(() => {
  const s = String(props.order?.status || '')
  if (s === '1') return 'var(--eapp-warning)'   // 待付款 - amber
  if (s === '2') return 'var(--eapp-primary)'   // 待发货 - blue
  if (s === '3') return 'var(--eapp-accent)'    // 已发货 - orange
  if (s === '4') return 'var(--eapp-success)'   // 已完成 - green
  if (s === '6') return 'var(--eapp-text-faint)' // 已取消 - gray
  return 'var(--eapp-border-strong)'
})

// Status badge style
const statusStyle = computed(() => {
  const s = String(props.order?.status || '')
  if (s === '1') return { bg: 'var(--eapp-warning-soft)', color: '#92400e' }
  if (s === '2') return { bg: 'var(--eapp-primary-soft)', color: '#1d4ed8' }
  if (s === '3') return { bg: 'var(--eapp-accent-soft)', color: '#c2410c' }
  if (s === '4') return { bg: 'var(--eapp-success-soft)', color: '#166534' }
  if (s === '6') return { bg: '#f1f5f9', color: '#64748b' }
  return { bg: '#f1f5f9', color: '#64748b' }
})
</script>

<template>
  <view :class="['order-card', selected && 'is-selected']" @click="$emit('click')">
    <!-- Status accent bar -->
    <view class="accent-bar" :style="{ background: accentColor }" />

    <view class="body">
      <!-- Header row -->
      <view class="head">
        <view class="head-left">
          <view v-if="selectable" class="checkbox" @click.stop="$emit('toggle')">
            <view :class="['cb', selected && 'cb-on']">
              <text v-if="selected" class="cb-check">✓</text>
            </view>
          </view>
          <text class="order-no">#{{ order.id }}</text>
          <text v-if="order.created_at" class="order-time">{{ formatTime(order.created_at) }}</text>
        </view>
        <text
          class="status-badge"
          :style="{ background: statusStyle.bg, color: statusStyle.color }"
        >{{ order.status_label || '-' }}</text>
      </view>

      <!-- Products row -->
      <view class="products">
        <view class="thumbs">
          <image v-for="it in thumbs" :key="it.id" :src="it.cover || ''" mode="aspectFill" class="thumb" />
          <view v-if="remain > 0" class="thumb-extra">+{{ remain }}</view>
        </view>
        <view class="product-meta">
          <text v-if="firstItem" class="product-name" :numberOfLines="2">{{ firstItem.title || firstItem.name || '商品' }}</text>
          <text v-if="items.length > 1" class="product-more">等 {{ items.length }} 件商品</text>
          <view class="price-row">
            <text class="qty-label">共 {{ itemCount }} 件</text>
            <text class="amount">¥{{ money(order.pay_amount || order.total_amount) }}</text>
          </view>
        </view>
      </view>

      <!-- Buyer info -->
      <view v-if="order.user_nickname || order.receiver_name || order.receiver_phone" class="buyer-row">
        <view class="buyer-icon">
          <!-- person icon SVG -->
          <text class="icon-dot">·</text>
        </view>
        <text class="buyer-name">{{ order.user_nickname || order.receiver_name || '匿名买家' }}</text>
        <text v-if="order.receiver_phone" class="buyer-phone">{{ maskPhone(order.receiver_phone) }}</text>
        <view class="buyer-spacer" />
        <text v-if="order.remark" class="has-note">有备注</text>
      </view>

      <!-- Divider -->
      <view class="sep" />

      <!-- Actions -->
      <view class="actions" @click.stop>
        <view class="actions-sec">
          <up-button size="mini" plain hairline @click="$emit('action', 'note')">备注</up-button>
          <up-button size="mini" plain hairline @click="$emit('action', 'print')">打单</up-button>
          <up-button v-if="String(order.status) === '1'" size="mini" type="warning" plain hairline @click="$emit('action', 'reprice')">改价</up-button>
          <up-button v-if="String(order.status) === '1'" size="mini" type="warning" plain hairline @click="$emit('action', 'remind_pay')">催付</up-button>
        </view>
        <view class="actions-pri">
          <up-button size="mini" plain hairline @click="$emit('action', 'detail')">详情</up-button>
          <up-button
            v-if="String(order.status) === '2'"
            size="mini"
            type="primary"
            @click="$emit('action', 'ship')"
          >发货</up-button>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
/* Card shell */
.order-card {
  display: flex;
  flex-direction: row;
  background: var(--eapp-card);
  border-radius: 20rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 16rpx rgba(0, 0, 0, 0.06);
}
.order-card.is-selected {
  box-shadow: 0 0 0 3rpx var(--eapp-primary);
  background: var(--eapp-primary-soft);
}

/* Accent bar */
.accent-bar {
  width: 8rpx;
  flex-shrink: 0;
  border-radius: 20rpx 0 0 20rpx;
}

/* Inner body */
.body {
  flex: 1;
  padding: 20rpx 22rpx 16rpx;
  min-width: 0;
}

/* Header */
.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8rpx;
}
.head-left {
  display: flex;
  align-items: center;
  gap: 10rpx;
  overflow: hidden;
  flex: 1;
}
.checkbox {
  flex-shrink: 0;
}
.cb {
  width: 38rpx;
  height: 38rpx;
  border-radius: 10rpx;
  border: 2rpx solid var(--eapp-border-strong);
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}
.cb.cb-on {
  background: var(--eapp-primary);
  border-color: var(--eapp-primary);
}
.cb-check {
  color: #fff;
  font-size: 22rpx;
  line-height: 1;
  font-weight: 700;
}
.order-no {
  font-size: 24rpx;
  font-weight: 700;
  color: var(--eapp-text);
  flex-shrink: 0;
}
.order-time {
  font-size: 20rpx;
  color: var(--eapp-text-faint);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.status-badge {
  font-size: 20rpx;
  border-radius: 999rpx;
  padding: 6rpx 14rpx;
  font-weight: 600;
  white-space: nowrap;
  flex-shrink: 0;
}

/* Products */
.products {
  display: flex;
  align-items: flex-start;
  gap: 14rpx;
  margin-top: 16rpx;
}
.thumbs {
  display: flex;
  gap: 6rpx;
  flex-shrink: 0;
  flex-wrap: nowrap;
}
.thumb {
  width: 84rpx;
  height: 84rpx;
  border-radius: 10rpx;
  background: var(--eapp-bg);
  flex-shrink: 0;
}
.thumb-extra {
  width: 84rpx;
  height: 84rpx;
  border-radius: 10rpx;
  background: var(--eapp-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22rpx;
  color: var(--eapp-text-muted);
  flex-shrink: 0;
}
.product-meta {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}
.product-name {
  font-size: 26rpx;
  color: var(--eapp-text);
  line-height: 1.5;
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}
.product-more {
  font-size: 20rpx;
  color: var(--eapp-text-faint);
}
.price-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-top: 4rpx;
}
.qty-label {
  font-size: 20rpx;
  color: var(--eapp-text-faint);
}
.amount {
  font-size: 32rpx;
  font-weight: 700;
  color: var(--eapp-primary);
}

/* Buyer row */
.buyer-row {
  display: flex;
  align-items: center;
  gap: 8rpx;
  margin-top: 12rpx;
}
.icon-dot {
  color: var(--eapp-text-faint);
  font-size: 24rpx;
}
.buyer-name {
  font-size: 24rpx;
  color: var(--eapp-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 160rpx;
}
.buyer-phone {
  font-size: 22rpx;
  color: var(--eapp-text-faint);
  flex-shrink: 0;
}
.buyer-spacer {
  flex: 1;
}
.has-note {
  font-size: 20rpx;
  color: var(--eapp-warning);
  background: var(--eapp-warning-soft);
  border-radius: 999rpx;
  padding: 4rpx 10rpx;
  flex-shrink: 0;
}

/* Separator */
.sep {
  height: 1rpx;
  background: var(--eapp-border);
  margin: 14rpx 0 12rpx;
}

/* Actions */
.actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8rpx;
}
.actions-sec {
  display: flex;
  gap: 8rpx;
  flex-wrap: wrap;
}
.actions-pri {
  display: flex;
  gap: 8rpx;
  flex-shrink: 0;
}
</style>
