<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { getProducts, updateProduct } from '@/api/product'

const { t } = useI18n()
const keyword = ref('')
const list = ref<any[]>([])
const loading = ref(false)

async function loadData() {
  loading.value = true
  try {
    const data: any = await getProducts({ page: 1, size: 20, keyword: keyword.value })
    list.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

async function onToggleSale(item: any) {
  const next = Number(item.status || 0) === 1 ? 0 : 1
  await updateProduct(item.id, { product: { status: next } })
  uni.showToast({ title: '状态已更新', icon: 'success' })
  loadData()
}

function goEdit(id = 0) {
  uni.navigateTo({ url: `/pages/product/edit?id=${id}` })
}

onShow(loadData)
</script>

<template>
  <view class="page">
    <view class="search">
      <up-search v-model="keyword" :placeholder="t('product.search')" @search="loadData" @custom="loadData" />
    </view>
    <view v-if="loading" class="empty">{{ t('common.loading') }}</view>
    <view v-else-if="!list.length" class="empty">{{ t('common.empty') }}</view>
    <view v-for="item in list" :key="item.id" class="card">
      <view class="title">{{ item.title || '-' }}</view>
      <view class="sub">库存：{{ item.stock || 0 }} · 价格：¥{{ Number(item.price || 0).toFixed(2) }}</view>
      <view class="actions">
        <up-button size="mini" type="primary" plain @click="goEdit(item.id)">{{ t('product.edit') }}</up-button>
        <up-button size="mini" type="warning" plain @click="onToggleSale(item)">{{ Number(item.status || 0) === 1 ? t('product.offSale') : t('product.onSale') }}</up-button>
      </view>
    </view>
    <view class="fab">
      <up-button type="primary" shape="circle" @click="goEdit(0)">+ 新建商品</up-button>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; }
.search { margin-bottom: 16rpx; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; margin-bottom: 14rpx; }
.title { font-size: 30rpx; font-weight: 600; }
.sub { margin-top: 8rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.actions { margin-top: 12rpx; display: flex; gap: 10rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
.fab { position: fixed; right: 24rpx; bottom: calc(24rpx + env(safe-area-inset-bottom)); }
</style>
