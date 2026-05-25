<template>
  <view v-if="show" class="fixed left-0 right-0 bottom-0 z-[999]" style="pointer-events: none;">
    <view class="flex justify-center pb-20rpx" style="pointer-events: auto;">
      <!-- Collapsed: pill showing current preset name -->
      <view v-if="!expanded" class="px-24rpx py-12rpx rounded-full text-24rpx shadow-lg"
        style="background: rgba(0,0,0,0.7); color: #fff; backdrop-filter: blur(8px);"
        @click="expanded = true">
        {{ currentName }} ▾
      </view>
      <!-- Expanded: preset selector panel -->
      <view v-else class="mx-30rpx rounded-2xl shadow-2xl overflow-hidden"
        style="background: rgba(255,255,255,0.96); backdrop-filter: blur(12px); max-width: 680rpx; width: 100%;">
        <view class="px-24rpx pt-20rpx pb-12rpx flex items-center justify-between">
          <text class="text-24rpx text-gray-400">切换演示行业</text>
          <text class="text-24rpx text-gray-400 px-8rpx" @click="expanded = false">✕</text>
        </view>
        <view class="flex flex-wrap px-16rpx pb-20rpx gap-12rpx">
          <view v-for="p in presetList" :key="p.key"
            class="px-24rpx py-14rpx rounded-full text-24rpx"
            :style="p.key === activeKey
              ? 'background: #1d4ed8; color: #fff;'
              : 'background: #f3f4f6; color: #374151;'"
            @click="switchPreset(p.key)">
            {{ p.name }}
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const show = import.meta.env.VITE_MOCK === 'true'
const expanded = ref(false)
const presetList = ref<Array<{ key: string; name: string }>>([])
const activeKey = ref('')
const currentName = ref('')

onMounted(async () => {
  if (!show) return
  const { listPresets, getPresetKey } = await import('@/mock/presets/index')
  presetList.value = listPresets()
  activeKey.value = getPresetKey()
  currentName.value = presetList.value.find(p => p.key === activeKey.value)?.name || activeKey.value
})

function switchPreset(key: string) {
  if (key === activeKey.value) {
    expanded.value = false
    return
  }
  const url = new URL(window.location.href)
  if (key === 'mall') {
    url.searchParams.delete('demo')
  } else {
    url.searchParams.set('demo', key)
  }
  window.location.href = url.toString()
}
</script>
