<script setup lang="ts">
function go(path: string) {
  uni.navigateTo({ url: path })
}

const entries = [
  { key: 'warehouse', title: '仓库管理', icon: '🏠', desc: '仓库信息维护', path: '/pages/wms/warehouse-list' },
  { key: 'stock', title: '库存台账', icon: '📋', desc: '实时库存查询', path: '/pages/wms/stock-ledger' },
  { key: 'doc', title: '出入库单', icon: '📄', desc: '入库/出库单管理', path: '/pages/wms/doc-list' },
  { key: 'movement', title: '库存流水', icon: '🔄', desc: '变动记录查询', path: '/pages/wms/movement-list' },
]
</script>

<template>
  <view class="page">
    <view class="title">仓储管理</view>
    <view class="grid">
      <view v-for="e in entries" :key="e.key" class="card" @click="go(e.path)">
        <text class="card-icon">{{ e.icon }}</text>
        <view class="card-body">
          <text class="card-title">{{ e.title }}</text>
          <text class="card-desc">{{ e.desc }}</text>
        </view>
        <text class="arrow">></text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; }
.title { font-size: 32rpx; font-weight: 700; margin-bottom: 16rpx; }
.grid { display: grid; gap: 12rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 24rpx; display: flex; align-items: center; gap: 16rpx; }
.card-icon { font-size: 40rpx; }
.card-body { flex: 1; min-width: 0; }
.card-title { font-size: 28rpx; font-weight: 600; display: block; }
.card-desc { font-size: 22rpx; color: var(--eapp-text-muted); display: block; margin-top: 4rpx; }
.arrow { color: var(--eapp-text-muted); font-size: 28rpx; }
</style>
