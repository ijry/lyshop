<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { getMessages, sendMessage } from '@/api/message'

const list = ref<any[]>([])
const loading = ref(false)
const sending = ref(false)
const form = reactive({ title: '', content: '' })

async function loadData() {
  loading.value = true
  try {
    const data: any = await getMessages({ page: 1, size: 20 })
    list.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
  } finally {
    loading.value = false
  }
}

async function onSend() {
  if (!form.title.trim() || !form.content.trim()) {
    uni.showToast({ title: '请填写标题和内容', icon: 'none' })
    return
  }
  sending.value = true
  try {
    await sendMessage({ title: form.title, content: form.content, group: 'system' })
    uni.showToast({ title: '发送成功', icon: 'success' })
    form.title = ''
    form.content = ''
    loadData()
  } finally {
    sending.value = false
  }
}

onShow(loadData)
</script>

<template>
  <view class="page">
    <view class="card">
      <up-input v-model="form.title" placeholder="消息标题" />
      <view class="mt-12rpx" />
      <up-textarea v-model="form.content" placeholder="消息内容" />
      <view class="mt-16rpx" />
      <up-button type="primary" :loading="sending" @click="onSend">发送消息</up-button>
    </view>

    <view class="list">
      <view v-if="loading" class="empty">加载中...</view>
      <view v-else-if="!list.length" class="empty">暂无消息</view>
      <view v-for="item in list" :key="item.id" class="msg-item">
        <view class="title">{{ item.title || '-' }}</view>
        <view class="content">{{ item.content || '-' }}</view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 20rpx; }
.list { margin-top: 14rpx; display: grid; gap: 12rpx; }
.msg-item { background: #fff; border: 1px solid var(--eapp-border); border-radius: 20rpx; padding: 18rpx; }
.title { font-size: 28rpx; font-weight: 600; }
.content { margin-top: 6rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
</style>
