<script setup lang="ts">
import SectionHeader from '@/components/biz/SectionHeader.vue'
import EntryList, { type EntryItem } from '@/components/biz/EntryList.vue'
import AppIcon from '@/components/common/AppIcon.vue'

const entries: EntryItem[] = [
  { key: 'warehouse', title: '仓库管理', desc: '仓库信息维护与启用状态', iconName: 'warehouse', path: '/pages/wms/warehouse-list' },
  { key: 'stock', title: '库存台账', desc: '实时库存查询与安全库存', iconName: 'layers', path: '/pages/wms/stock-ledger' },
  { key: 'doc', title: '出入库单', desc: '入库 / 出库单创建与审核', iconName: 'document', path: '/pages/wms/doc-list' },
  { key: 'movement', title: '库存流水', desc: '库存变动记录与追溯', iconName: 'refresh', path: '/pages/wms/movement-list' },
]

function go(item: EntryItem) {
  if (item.path) uni.navigateTo({ url: item.path })
}
</script>

<template>
  <view class="page">
    <view class="intro">
      <view class="intro__icon">
        <AppIcon name="warehouse" :size="40" color="primary" />
      </view>
      <view class="intro__text">
        <text class="intro__title">仓储管理</text>
        <text class="intro__desc">仓库、库存与单据的一体化管理</text>
      </view>
    </view>

    <view class="group">
      <SectionHeader title="仓储业务" />
      <EntryList :items="entries" @click="go" />
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background: var(--eapp-bg);
  padding: var(--eapp-space-5);
  padding-bottom: calc(env(safe-area-inset-bottom) + 48rpx);
  box-sizing: border-box;
}

.intro {
  display: flex;
  align-items: center;
  gap: var(--eapp-space-3);
  padding: var(--eapp-space-4) var(--eapp-space-5);
  margin-bottom: var(--eapp-space-6);
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  box-sizing: border-box;
}

.intro__icon {
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

.intro__text {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  min-width: 0;
}

.intro__title {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
}

.intro__desc {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
}

.group {
  margin-bottom: var(--eapp-space-6);
}
</style>
