<template>
  <view style="min-height: 100vh; background: #f5f5f5;">
    <u-navbar :title="$t('message.title')" :placeholder="true" />

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
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { get } from '@/utils/request'

const { t } = useI18n()

const unread = ref<Record<string, number>>({})

const groups = computed(() => [
  { key: 'system', title: t('message.system'), desc: t('message.systemDesc'), icon: '🔔', bg: '#fef2f2' },
  { key: 'order', title: t('message.order'), desc: t('message.orderDesc'), icon: '📦', bg: '#eff6ff' },
  { key: 'marketing', title: t('message.marketing'), desc: t('message.marketingDesc'), icon: '🏷️', bg: '#fff7ed' },
  { key: 'im', title: t('message.service'), desc: t('message.serviceDesc'), icon: '💬', bg: '#f0fdf4' },
])

onMounted(async () => {
  const data = await get<any>('/api/v1/messages/unread')
  if (data) unread.value = data
})
</script>
