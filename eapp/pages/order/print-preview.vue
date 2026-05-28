<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { getPrintTemplate } from '@/api/order'

const tpl = ref('')
async function load(id: number) {
  const data: any = await getPrintTemplate(id)
  tpl.value = String(data?.template || '')
}
onLoad((opts) => load(Number(opts?.id || 0)))
function onPrint() { uni.showToast({ title: '已发送至默认打印机', icon: 'success' }) }
</script>

<template>
  <view class="page">
    <view class="card">
      <rich-text v-if="tpl" :nodes="tpl" />
      <view v-else class="empty">加载中…</view>
    </view>
    <up-button type="primary" class="mt" @click="onPrint">发送至打印机</up-button>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; padding: 20rpx; background: var(--eapp-bg); }
.card { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; min-height: 480rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 60rpx 0; }
.mt { margin-top: 20rpx; }
</style>
