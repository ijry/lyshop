<template>
  <view style="min-height: 100vh; background: #f5f5f5;">
    <u-navbar :title="groupTitle" :placeholder="true" />

    <view style="padding: 12px 16px;">
      <view v-for="msg in messages" :key="msg.id"
        style="background: #fff; border-radius: 12px; padding: 14px 16px; margin-bottom: 8px; box-shadow: 0 1px 4px rgba(0,0,0,0.03);">
        <view style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px;">
          <view style="display: flex; align-items: center; gap: 6px;">
            <view v-if="!msg.is_read" style="width: 6px; height: 6px; border-radius: 3px; background: #dc2626;" />
            <text style="font-size: 14px; color: #111; font-weight: 600;">{{ msg.title }}</text>
          </view>
          <text style="font-size: 11px; color: #ccc;">{{ msg.created_at?.slice(0,10) }}</text>
        </view>
        <text style="font-size: 13px; color: #666; line-height: 1.5;">{{ msg.content }}</text>
      </view>

      <view v-if="!messages.length" style="text-align: center; padding: 60px 0; color: #999; font-size: 14px;">
        {{ $t('messageList.empty') }}
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { get } from '@/utils/request'

const { t } = useI18n()

const messages = ref<any[]>([])
const group = ref('')

const groupTitles = computed<Record<string, string>>(() => ({
  system: t('messageList.system'), order: t('messageList.order'), marketing: t('messageList.marketing'), im: t('messageList.service')
}))
const groupTitle = computed(() => groupTitles.value[group.value] || t('messageList.title'))

onMounted(async () => {
  const pages = getCurrentPages()
  const query = (pages[pages.length - 1] as any).options
  group.value = query.group || ''
  const data = await get<any>('/api/v1/messages', { group: group.value })
  messages.value = (data?.list || []).filter(
    (m: any) => !group.value || m.group === group.value
  )
})
</script>
