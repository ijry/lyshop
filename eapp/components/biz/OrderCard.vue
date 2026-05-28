<script setup lang="ts">
import { computed } from 'vue'
import StatusTag from '@/components/common/StatusTag.vue'

const props = defineProps<{ order: any; selectable?: boolean; selected?: boolean }>()
defineEmits<{ (e: 'click'): void; (e: 'toggle'): void; (e: 'action', key: string): void }>()
const thumbs = computed(() => (Array.isArray(props.order?.items) ? props.order.items.slice(0, 3) : []))
const remain = computed(() => Math.max(0, (props.order?.items?.length || 0) - 3))
function money(v: any) { return Number(v || 0).toFixed(2) }
</script>

<template>
  <view :class="['order-card', selected ? 'is-selected' : '']" @click="$emit('click')">
    <view class="head">
      <view class="left">
        <view v-if="selectable" class="check" @click.stop="$emit('toggle')"><text>{{ selected ? '☑' : '☐' }}</text></view>
        <text class="no">#{{ order.id }}</text>
      </view>
      <StatusTag :text="order.status_label || order.status || '-'" :type="order.status" />
    </view>
    <view class="items">
      <view v-for="it in thumbs" :key="it.id" class="item">
        <image v-if="it.cover" :src="it.cover" mode="aspectFill" class="cover" />
        <view v-else class="cover placeholder">商品</view>
      </view>
      <view v-if="remain > 0" class="more">+{{ remain }}</view>
    </view>
    <view class="meta">
      <text class="user">{{ order.user_nickname || order.receiver_name || '匿名买家' }}</text>
      <text class="amount">¥{{ money(order.pay_amount || order.total_amount) }}</text>
    </view>
    <view class="actions" @click.stop>
      <up-button size="mini" plain @click="$emit('action', 'detail')">详情</up-button>
      <up-button v-if="String(order.status) === '1'" size="mini" type="warning" plain @click="$emit('action', 'reprice')">改价</up-button>
      <up-button v-if="String(order.status) === '1'" size="mini" type="primary" plain @click="$emit('action', 'remind_pay')">催付</up-button>
      <up-button v-if="String(order.status) === '2'" size="mini" type="primary" plain @click="$emit('action', 'ship')">发货</up-button>
      <up-button size="mini" plain @click="$emit('action', 'note')">备注</up-button>
      <up-button size="mini" plain @click="$emit('action', 'print')">打单</up-button>
    </view>
  </view>
</template>

<style scoped>
.order-card { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; }
.order-card.is-selected { border-color: var(--eapp-primary); background: var(--eapp-primary-soft); }
.head { display: flex; align-items: center; justify-content: space-between; }
.left { display: flex; align-items: center; gap: 12rpx; }
.check { width: 44rpx; height: 44rpx; display: flex; align-items: center; justify-content: center; font-size: 28rpx; }
.no { font-size: 22rpx; color: var(--eapp-text-muted); }
.items { margin-top: 12rpx; display: flex; align-items: center; gap: 8rpx; }
.cover { width: 96rpx; height: 96rpx; border-radius: 12rpx; background: var(--eapp-bg); display: flex; align-items: center; justify-content: center; color: var(--eapp-text-faint); font-size: 22rpx; }
.more { padding: 0 12rpx; height: 96rpx; display: flex; align-items: center; color: var(--eapp-text-muted); font-size: 24rpx; }
.meta { margin-top: 12rpx; display: flex; align-items: center; justify-content: space-between; }
.user { font-size: 24rpx; color: var(--eapp-text-muted); }
.amount { font-size: 32rpx; color: var(--eapp-primary); font-weight: 700; }
.actions { margin-top: 14rpx; display: flex; flex-wrap: wrap; gap: 10rpx; justify-content: flex-end; }
</style>
