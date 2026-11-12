<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'
import StatusTag from '@/components/common/StatusTag.vue'

defineProps<{ product: any; selectable?: boolean; selected?: boolean }>()
defineEmits<{ (e: 'click'): void; (e: 'toggle'): void; (e: 'action', key: string): void }>()
</script>

<template>
  <view
    :class="['product-card', selected ? 'is-selected' : '']"
    hover-class="eapp-hover"
    :hover-stay-time="80"
    @click="$emit('click')"
  >
    <view v-if="selectable" class="check" @click.stop="$emit('toggle')">
      <view :class="['cb', selected && 'cb-on']">
        <AppIcon v-if="selected" name="check" :size="20" color="white" :stroke-width="2.4" />
      </view>
    </view>
    <image v-if="product.cover" :src="product.cover" mode="aspectFill" class="cover" />
    <view v-else class="cover cover--placeholder">
      <AppIcon name="product" :size="56" color="faint" :stroke-width="1.4" />
    </view>
    <view class="body">
      <view class="title-row">
        <text class="title">{{ product.title }}</text>
        <StatusTag :text="Number(product.status || 0) === 1 ? '在售' : '仓库'" :type="Number(product.status || 0) === 1 ? 'enabled' : 'disabled'" />
      </view>
      <view class="meta eapp-num">
        ¥{{ Number(product.price || 0).toFixed(2) }} · 库存 {{ product.stock || 0 }} · 销量 {{ product.sales_count || 0 }}
      </view>
      <view class="cat">{{ product.category_path_name || '未分类' }}</view>
      <view class="actions" @click.stop>
        <up-button size="mini" plain hairline @click="$emit('action', 'edit')">编辑</up-button>
        <up-button size="mini" type="warning" plain hairline @click="$emit('action', 'toggle_sale')">
          {{ Number(product.status || 0) === 1 ? '下架' : '上架' }}
        </up-button>
      </view>
    </view>
  </view>
</template>

<style scoped>
.product-card {
  position: relative;
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-3);
  display: flex;
  gap: var(--eapp-space-3);
  align-items: flex-start;
  box-sizing: border-box;
}

.product-card.is-selected {
  border-color: var(--eapp-primary);
  background: var(--eapp-primary-soft);
}

.check {
  position: absolute;
  top: 12rpx;
  right: 12rpx;
  z-index: 2;
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

.cover {
  width: 160rpx;
  height: 160rpx;
  border-radius: var(--eapp-radius);
  border: 1rpx solid var(--eapp-border);
  flex-shrink: 0;
  background: var(--eapp-bg);
  box-sizing: border-box;
}

.cover--placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--eapp-surface-subtle);
}

.body {
  flex: 1;
  min-width: 0;
}

.title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10rpx;
  padding-right: 40rpx;
}

.title {
  font-size: var(--eapp-fs-md);
  font-weight: 600;
  color: var(--eapp-text);
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.meta {
  margin-top: 8rpx;
  color: var(--eapp-text-secondary);
  font-size: var(--eapp-fs-sm);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cat {
  margin-top: 4rpx;
  color: var(--eapp-text-faint);
  font-size: var(--eapp-fs-xs);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.actions {
  margin-top: 12rpx;
  display: flex;
  gap: 10rpx;
}
</style>
