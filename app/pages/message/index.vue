<template>
  <view style="min-height: 100vh; background: #f5f5f5;">
    <u-navbar title="消息中心" :placeholder="true" />

    <view style="padding: 12px 16px;">
      <view v-for="g in groups" :key="g.key"
        @click="uni.navigateTo({url:`/pages/message/list?group=${g.key}`})"
        style="display: flex; align-items: center; background: #fff; border-radius: 14px; padding: 16px; margin-bottom: 10px; box-shadow: 0 1px 4px rgba(0,0,0,0.04);">
        <view :style="{ width: '42px', height: '42px', borderRadius: '12px', background: g.bg, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '20px', flexShrink: 0 }">
          {{ g.icon }}
        </view>
        <view style="flex: 1; margin-left: 14px;">
          <text style="font-size: 15px; color: #111; font-weight: 600;">{{ g.title }}</text>
          <text style="font-size: 12px; color: #999; display: block; margin-top: 2px;">{{ g.desc }}</text>
        </view>
        <view style="display: flex; align-items: center; gap: 6px;">
          <view v-if="(unread[g.key] || 0) > 0"
            style="min-width: 18px; height: 18px; background: #dc2626; border-radius: 9px; display: flex; align-items: center; justify-content: center; padding: 0 5px;">
            <text style="color: #fff; font-size: 10px; font-weight: 600;">{{ unread[g.key] }}</text>
          </view>
          <u-icon name="arrow-right" size="14" color="#ccc" />
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'

const unread = ref<Record<string, number>>({})

const groups = [
  { key: 'system', title: '系统通知', desc: '系统公告、升级通知', icon: '🔔', bg: '#fef2f2' },
  { key: 'order', title: '订单消息', desc: '发货、签收、售后通知', icon: '📦', bg: '#eff6ff' },
  { key: 'marketing', title: '营销消息', desc: '优惠券、活动、促销通知', icon: '🏷️', bg: '#fff7ed' },
  { key: 'im', title: '客服消息', desc: '客服回复通知', icon: '💬', bg: '#f0fdf4' },
]

onMounted(async () => {
  const data = await get<any>('/api/v1/messages/unread')
  if (data) unread.value = data
})
</script>
