<script setup lang="ts">
import { computed } from 'vue'
import AppIcon from '@/components/common/AppIcon.vue'

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
  const m = s.match(/\d{4}-(\d{2})-(\d{2}) (\d{2}:\d{2})/)
  return m ? `${m[1]}/${m[2]} ${m[3]}` : s.slice(5, 16)
}

function maskPhone(v?: string) {
  if (!v) return ''
  return String(v).replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

/** 状态徽标：低饱和底 + 同色系细边框 + 深色文字 */
const statusStyle = computed(() => {
  const s = String(props.order?.status || '')
  if (s === '1') return { bg: 'var(--eapp-warning-soft)', color: 'var(--eapp-warning)', border: '#ecd6b8' }
  if (s === '2') return { bg: 'var(--eapp-primary-soft)', color: 'var(--eapp-primary)', border: 'var(--eapp-primary-border)' }
  if (s === '3') return { bg: 'var(--eapp-accent-soft)', color: 'var(--eapp-accent)', border: '#c4dee6' }
  if (s === '4') return { bg: 'var(--eapp-success-soft)', color: 'var(--eapp-success)', border: '#c7e3d2' }
  return { bg: 'var(--eapp-surface-subtle)', color: 'var(--eapp-text-muted)', border: 'var(--eapp-border-strong)' }
})
</script>

<template>
  <view
    :class="['order-card', selected && 'is-selected']"
    hover-class="eapp-hover"
    :hover-stay-time="80"
    @click="$emit('click')"
  >
    <view class="body">
      <!-- Header row -->
      <view class="head">
        <view class="head-left">
          <view v-if="selectable" class="checkbox" @click.stop="$emit('toggle')">
            <view :class="['cb', selected && 'cb-on']">
              <AppIcon v-if="selected" name="check" :size="20" color="white" :stroke-width="2.4" />
            </view>
          </view>
          <text class="order-no eapp-num">#{{ order.id }}</text>
          <text v-if="order.created_at" class="order-time eapp-num">{{ formatTime(order.created_at) }}</text>
        </view>
        <view
          class="status-badge"
          :style="{ background: statusStyle.bg, borderColor: statusStyle.border, color: statusStyle.color }"
        >
          <view class="status-badge__dot" :style="{ background: statusStyle.color }" />
          <text class="status-badge__text">{{ order.status_label || '-' }}</text>
        </view>
      </view>

      <!-- Products row -->
      <view class="products">
        <view class="thumbs">
          <image v-for="it in thumbs" :key="it.id" :src="it.cover || ''" mode="aspectFill" class="thumb" />
          <view v-if="remain > 0" class="thumb-extra eapp-num">+{{ remain }}</view>
        </view>
        <view class="product-meta">
          <text v-if="firstItem" class="product-name" :numberOfLines="2">{{ firstItem.title || firstItem.name || '商品' }}</text>
          <text v-if="items.length > 1" class="product-more">等 {{ items.length }} 件商品</text>
          <view class="price-row">
            <text class="qty-label">共 {{ itemCount }} 件</text>
            <text class="amount eapp-num">¥{{ money(order.pay_amount || order.total_amount) }}</text>
          </view>
        </view>
      </view>

      <!-- Buyer info -->
      <view v-if="order.user_nickname || order.receiver_name || order.receiver_phone" class="buyer-row">
        <AppIcon name="user" :size="26" color="faint" />
        <text class="buyer-name">{{ order.user_nickname || order.receiver_name || '匿名买家' }}</text>
        <text v-if="order.receiver_phone" class="buyer-phone eapp-num">{{ maskPhone(order.receiver_phone) }}</text>
        <view class="buyer-spacer" />
        <text v-if="order.remark" class="has-note">有备注</text>
      </view>

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
.order-card {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  overflow: hidden;
  box-sizing: border-box;
}

.order-card.is-selected {
  border-color: var(--eapp-primary);
  background: var(--eapp-primary-soft);
}

.body {
  padding: var(--eapp-space-4) var(--eapp-space-5) var(--eapp-space-3);
  min-width: 0;
  box-sizing: border-box;
}

/* Header */
.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--eapp-space-2);
}

.head-left {
  display: flex;
  align-items: center;
  gap: 10rpx;
  overflow: hidden;
  flex: 1;
  min-width: 0;
}

.checkbox {
  flex-shrink: 0;
  padding: 4rpx;
}

.cb {
  width: 34rpx;
  height: 34rpx;
  border-radius: var(--eapp-radius-sm);
  border: 2rpx solid var(--eapp-border-strong);
  background: var(--eapp-card);
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
}

.cb.cb-on {
  background: var(--eapp-primary);
  border-color: var(--eapp-primary);
}

.order-no {
  font-size: var(--eapp-fs-base);
  font-weight: 700;
  color: var(--eapp-text);
  flex-shrink: 0;
}

.order-time {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  border: 1rpx solid;
  border-radius: var(--eapp-radius-sm);
  padding: 6rpx 12rpx;
  flex-shrink: 0;
  box-sizing: border-box;
}

.status-badge__dot {
  width: 8rpx;
  height: 8rpx;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-badge__text {
  font-size: var(--eapp-fs-xs);
  font-weight: 600;
  white-space: nowrap;
  color: inherit;
}

/* Products */
.products {
  display: flex;
  align-items: flex-start;
  gap: 14rpx;
  margin-top: var(--eapp-space-3);
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
  border-radius: var(--eapp-radius-sm);
  border: 1rpx solid var(--eapp-border);
  background: var(--eapp-bg);
  flex-shrink: 0;
  box-sizing: border-box;
}

.thumb-extra {
  width: 84rpx;
  height: 84rpx;
  border-radius: var(--eapp-radius-sm);
  background: var(--eapp-surface-subtle);
  border: 1rpx solid var(--eapp-border);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
  flex-shrink: 0;
  box-sizing: border-box;
}

.product-meta {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.product-name {
  font-size: var(--eapp-fs-base);
  color: var(--eapp-text);
  line-height: 1.5;
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.product-more {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}

.price-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 8rpx;
  margin-top: 4rpx;
}

.qty-label {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}

.amount {
  font-size: var(--eapp-fs-lg);
  font-weight: 700;
  color: var(--eapp-primary);
}

/* Buyer row */
.buyer-row {
  display: flex;
  align-items: center;
  gap: 8rpx;
  margin-top: var(--eapp-space-2);
}

.buyer-name {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 160rpx;
}

.buyer-phone {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
  flex-shrink: 0;
}

.buyer-spacer {
  flex: 1;
}

.has-note {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-warning);
  background: var(--eapp-warning-soft);
  border: 1rpx solid #ecd6b8;
  border-radius: var(--eapp-radius-sm);
  padding: 2rpx 10rpx;
  flex-shrink: 0;
}

.sep {
  height: 1rpx;
  background: var(--eapp-divider);
  margin: var(--eapp-space-3) 0 var(--eapp-space-2);
}

.actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--eapp-space-1);
}

.actions-sec {
  display: flex;
  gap: var(--eapp-space-1);
  flex-wrap: wrap;
}

.actions-pri {
  display: flex;
  gap: var(--eapp-space-1);
  flex-shrink: 0;
}
</style>
