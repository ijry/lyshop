<template>
  <view style="min-height: 100vh; background: #f5f5f5;">
    <u-navbar title="收货地址" :placeholder="true" />

    <view style="padding: 12px 16px;">
      <view v-for="addr in addresses" :key="addr.id"
        style="background: #fff; border-radius: 16px; padding: 16px 20px; margin-bottom: 12px; box-shadow: 0 1px 6px rgba(0,0,0,0.04);">
        <view style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px;">
          <view style="display: flex; align-items: center; gap: 12px;">
            <text style="font-size: 15px; font-weight: 600; color: #111;">{{ addr.name }}</text>
            <text style="font-size: 13px; color: #666;">{{ addr.phone }}</text>
          </view>
          <view v-if="addr.is_default"
            style="background: #fef2f2; color: #dc2626; font-size: 11px; padding: 2px 8px; border-radius: 4px;">
            默认
          </view>
        </view>
        <text style="font-size: 13px; color: #999; line-height: 1.5;">
          {{ addr.province }}{{ addr.city }}{{ addr.district }} {{ addr.detail }}
        </text>
      </view>

      <view v-if="!addresses.length"
        style="text-align: center; padding: 80px 0; color: #999; font-size: 14px;">
        暂无收货地址
      </view>
    </view>

    <!-- Add button -->
    <view style="padding: 16px; position: fixed; bottom: 0; left: 0; right: 0; background: #f5f5f5;">
      <u-button type="primary" text="新增收货地址" shape="circle"
        :custom-style="{background: '#dc2626', borderColor: '#dc2626'}" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/utils/request'

const addresses = ref<any[]>([])

onMounted(async () => {
  const data = await get<any[]>('/api/v1/addresses')
  addresses.value = data || []
})
</script>
