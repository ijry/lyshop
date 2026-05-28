<script setup lang="ts">
import StatusTag from '@/components/common/StatusTag.vue'
defineProps<{ product: any; selectable?: boolean; selected?: boolean }>()
defineEmits<{ (e: 'click'): void; (e: 'toggle'): void; (e: 'action', key: string): void }>()
</script>

<template>
  <view :class="['product-card', selected ? 'is-selected' : '']" @click="$emit('click')">
    <view v-if="selectable" class="check" @click.stop="$emit('toggle')">{{ selected ? '☑' : '☐' }}</view>
    <image v-if="product.cover" :src="product.cover" mode="aspectFill" class="cover" />
    <view v-else class="cover placeholder">图</view>
    <view class="body">
      <view class="title-row">
        <text class="title">{{ product.title }}</text>
        <StatusTag :text="Number(product.status || 0) === 1 ? '在售' : '仓库'" :type="Number(product.status || 0) === 1 ? 'enabled' : 'disabled'" />
      </view>
      <view class="meta">¥{{ Number(product.price || 0).toFixed(2) }} · 库存 {{ product.stock || 0 }} · 销量 {{ product.sales_count || 0 }}</view>
      <view class="cat">{{ product.category_path_name || '未分类' }}</view>
      <view class="actions" @click.stop>
        <up-button size="mini" plain @click="$emit('action', 'edit')">编辑</up-button>
        <up-button size="mini" type="warning" plain @click="$emit('action', 'toggle_sale')">{{ Number(product.status || 0) === 1 ? '下架' : '上架' }}</up-button>
      </view>
    </view>
  </view>
</template>

<style scoped>
.product-card { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 16rpx; display: flex; gap: 16rpx; align-items: flex-start; position: relative; }
.product-card.is-selected { border-color: var(--eapp-primary); background: var(--eapp-primary-soft); }
.check { position: absolute; top: 14rpx; right: 14rpx; font-size: 28rpx; }
.cover { width: 160rpx; height: 160rpx; border-radius: 14rpx; flex-shrink: 0; background: var(--eapp-bg); display: flex; align-items: center; justify-content: center; color: var(--eapp-text-faint); }
.body { flex: 1; min-width: 0; }
.title-row { display: flex; align-items: center; justify-content: space-between; gap: 10rpx; }
.title { font-size: 28rpx; font-weight: 600; flex: 1; }
.meta { margin-top: 6rpx; color: var(--eapp-text-muted); font-size: 23rpx; }
.cat { margin-top: 4rpx; color: var(--eapp-text-faint); font-size: 22rpx; }
.actions { margin-top: 10rpx; display: flex; gap: 10rpx; }
</style>
