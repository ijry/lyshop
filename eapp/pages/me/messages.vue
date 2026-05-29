<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import { reactive, ref, computed } from 'vue'
import { getMessages, sendMessage, markMessageRead } from '@/api/message'

const list = ref<any[]>([])
const loading = ref(false)
const sending = ref(false)
const form = reactive({ title: '', content: '' })
const activeGroup = ref('')
const currentTab = ref(0)

const groupTabs = [
  { key: '', label: '全部' },
  { key: 'system', label: '系统' },
  { key: 'order', label: '订单' },
  { key: 'marketing', label: '营销' },
  { key: 'im', label: '客服' },
]

const filteredList = computed(() => {
  if (!activeGroup.value) return list.value
  return list.value.filter((item: any) => item.group === activeGroup.value)
})

async function loadData() {
  loading.value = true
  try {
    const data: any = await getMessages({ page: 1, size: 50, group: activeGroup.value || undefined })
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

async function onMarkRead(item: any) {
  if (item.is_read) return
  await markMessageRead(item.id)
  item.is_read = 1
  uni.showToast({ title: '已标记已读', icon: 'success' })
}

function priorityClass(p: string) {
  if (p === 'urgent') return 'dot-urgent'
  if (p === 'important') return 'dot-important'
  return ''
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

    <up-tabs
      :list="groupTabs"
      :current="currentTab"
      :scrollable="true"
      keyName="label"
      @click="(item) => { currentTab = item.index; activeGroup = groupTabs[item.index].key; loadData() }"
      :activeStyle="{ color: '#fff', backgroundColor: 'var(--eapp-primary)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :inactiveStyle="{ color: 'var(--eapp-text-muted)', backgroundColor: 'var(--eapp-bg)', borderRadius: '999rpx', height: '56rpx', lineHeight: '56rpx', padding: '0 24rpx' }"
      :itemStyle="{ padding: '0 4rpx', height: '80rpx' }"
      lineColor="transparent"
    />

    <view class="list">
      <view v-if="loading" class="empty">加载中...</view>
      <view v-else-if="!filteredList.length" class="empty">暂无消息</view>
      <view v-for="item in filteredList" :key="item.id" class="msg-item">
        <view class="msg-top">
          <view class="title-row">
            <view v-if="item.priority && item.priority !== 'normal'" :class="['priority-dot', priorityClass(item.priority)]" />
            <text class="title">{{ item.title || '-' }}</text>
          </view>
          <view v-if="!item.is_read" class="read-btn" @click="onMarkRead(item)">标记已读</view>
        </view>
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
.msg-top { display: flex; align-items: center; justify-content: space-between; }
.title-row { display: flex; align-items: center; gap: 8rpx; }
.priority-dot { width: 14rpx; height: 14rpx; border-radius: 50%; flex-shrink: 0; }
.dot-urgent { background: #dc2626; }
.dot-important { background: #f59e0b; }
.title { font-size: 28rpx; font-weight: 600; }
.read-btn { font-size: 22rpx; color: var(--eapp-primary, #2563eb); padding: 4rpx 12rpx; }
.content { margin-top: 6rpx; color: var(--eapp-text-muted); font-size: 24rpx; }
.empty { text-align: center; color: var(--eapp-text-muted); padding: 80rpx 0; }
</style>
